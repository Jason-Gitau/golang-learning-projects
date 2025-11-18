package storage

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/user/web-scraper/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database handles all database operations
type Database struct {
	db *gorm.DB
	mu sync.RWMutex
}

// NewDatabase creates a new database connection
func NewDatabase(dbPath string) (*Database, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(&models.ScrapedPage{}, &models.Link{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Printf("Database initialized at: %s", dbPath)

	return &Database{db: db}, nil
}

// SavePage saves a scraped page to the database
func (d *Database) SavePage(result *models.ScrapeResult) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Create the page record
	page := models.ScrapedPage{
		URL:         result.URL,
		Title:       result.Title,
		Description: result.Description,
		LinkCount:   len(result.Links),
		StatusCode:  result.StatusCode,
		RetryCount:  result.RetryCount,
		ScrapedAt:   result.ScrapedAt,
		Duration:    result.Duration.Milliseconds(),
	}

	if result.Error != nil {
		page.Error = result.Error.Error()
	}

	// Use FirstOrCreate to avoid duplicates
	var existingPage models.ScrapedPage
	if err := d.db.Where("url = ?", page.URL).First(&existingPage).Error; err == nil {
		// Update existing page
		page.ID = existingPage.ID
		if err := d.db.Save(&page).Error; err != nil {
			return fmt.Errorf("failed to update page: %w", err)
		}

		// Delete old links
		d.db.Where("page_id = ?", page.ID).Delete(&models.Link{})
	} else {
		// Create new page
		if err := d.db.Create(&page).Error; err != nil {
			return fmt.Errorf("failed to create page: %w", err)
		}
	}

	// Save links
	if len(result.Links) > 0 {
		links := make([]models.Link, len(result.Links))
		for i, linkData := range result.Links {
			links[i] = models.Link{
				PageID: page.ID,
				URL:    linkData.URL,
				Text:   linkData.Text,
			}
		}

		if err := d.db.Create(&links).Error; err != nil {
			return fmt.Errorf("failed to create links: %w", err)
		}
	}

	return nil
}

// GetPage retrieves a page by URL
func (d *Database) GetPage(url string) (*models.ScrapedPage, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var page models.ScrapedPage
	if err := d.db.Where("url = ?", url).First(&page).Error; err != nil {
		return nil, err
	}

	return &page, nil
}

// GetAllPages retrieves all scraped pages
func (d *Database) GetAllPages() ([]models.ScrapedPage, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var pages []models.ScrapedPage
	if err := d.db.Order("scraped_at DESC").Find(&pages).Error; err != nil {
		return nil, err
	}

	return pages, nil
}

// GetPageLinks retrieves all links for a specific page
func (d *Database) GetPageLinks(pageID uint) ([]models.Link, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var links []models.Link
	if err := d.db.Where("page_id = ?", pageID).Find(&links).Error; err != nil {
		return nil, err
	}

	return links, nil
}

// GetStatistics retrieves scraping statistics
func (d *Database) GetStatistics() (*models.Statistics, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	stats := &models.Statistics{}

	// Count total pages
	var totalPages int64
	d.db.Model(&models.ScrapedPage{}).Count(&totalPages)
	stats.TotalPages = int(totalPages)

	// Count successful pages (status code 200)
	var successfulPages int64
	d.db.Model(&models.ScrapedPage{}).Where("status_code = ? AND error = ?", 200, "").Count(&successfulPages)
	stats.SuccessfulPages = int(successfulPages)

	// Count failed pages
	var failedPages int64
	d.db.Model(&models.ScrapedPage{}).Where("status_code != ? OR error != ?", 200, "").Count(&failedPages)
	stats.FailedPages = int(failedPages)

	// Count total links
	var totalLinks int64
	d.db.Model(&models.Link{}).Count(&totalLinks)
	stats.TotalLinks = int(totalLinks)

	// Calculate total and average duration
	var totalDuration int64
	d.db.Model(&models.ScrapedPage{}).Select("SUM(duration)").Scan(&totalDuration)
	stats.TotalDuration = time.Duration(totalDuration) * time.Millisecond

	if stats.TotalPages > 0 {
		avgDuration := totalDuration / int64(stats.TotalPages)
		stats.AverageDuration = time.Duration(avgDuration) * time.Millisecond
	}

	// Get time range
	var firstPage, lastPage models.ScrapedPage
	d.db.Order("scraped_at ASC").First(&firstPage)
	d.db.Order("scraped_at DESC").First(&lastPage)

	if firstPage.ID != 0 {
		stats.StartTime = firstPage.ScrapedAt
	}
	if lastPage.ID != 0 {
		stats.EndTime = lastPage.ScrapedAt
	}

	return stats, nil
}

// DeleteAllPages deletes all scraped pages and their links
func (d *Database) DeleteAllPages() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Delete all links first
	if err := d.db.Exec("DELETE FROM links").Error; err != nil {
		return fmt.Errorf("failed to delete links: %w", err)
	}

	// Delete all pages
	if err := d.db.Exec("DELETE FROM scraped_pages").Error; err != nil {
		return fmt.Errorf("failed to delete pages: %w", err)
	}

	log.Println("All data deleted from database")
	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
