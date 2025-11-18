package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"deep-research-agent/models"
)

// ResearchDB handles database operations for research sessions
type ResearchDB struct {
	db *gorm.DB
}

// NewResearchDB creates a new database connection
func NewResearchDB(dbPath string) (*ResearchDB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate schemas
	if err := db.AutoMigrate(
		&models.ResearchSession{},
		&models.Document{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &ResearchDB{db: db}, nil
}

// SaveSession saves a research session to the database
func (r *ResearchDB) SaveSession(query models.ResearchQuery, result *models.ResearchResult, status string) (*models.ResearchSession, error) {
	// Serialize query and result to JSON
	queryData, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	var resultData []byte
	if result != nil {
		resultData, err = json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal result: %w", err)
		}
	}

	// Calculate progress
	progress := 0.0
	if result != nil && result.TotalSteps > 0 {
		progress = float64(result.SuccessfulSteps) / float64(result.TotalSteps) * 100
	}

	session := &models.ResearchSession{
		ID:        uuid.New().String(),
		Query:     query.Query,
		QueryData: string(queryData),
		Result:    string(resultData),
		Status:    status,
		Progress:  progress,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if status == "completed" || status == "failed" {
		now := time.Now()
		session.CompletedAt = &now
	}

	if err := r.db.Create(session).Error; err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	return session, nil
}

// UpdateSession updates an existing research session
func (r *ResearchDB) UpdateSession(sessionID string, result *models.ResearchResult, status string) error {
	resultData, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	progress := 0.0
	if result != nil && result.TotalSteps > 0 {
		progress = float64(result.SuccessfulSteps) / float64(result.TotalSteps) * 100
	}

	updates := map[string]interface{}{
		"result":     string(resultData),
		"status":     status,
		"progress":   progress,
		"updated_at": time.Now(),
	}

	if status == "completed" || status == "failed" {
		now := time.Now()
		updates["completed_at"] = &now
	}

	if err := r.db.Model(&models.ResearchSession{}).Where("id = ?", sessionID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	return nil
}

// GetSession retrieves a research session by ID
func (r *ResearchDB) GetSession(sessionID string) (*models.ResearchSession, error) {
	var session models.ResearchSession
	if err := r.db.Where("id = ?", sessionID).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("session not found: %s", sessionID)
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return &session, nil
}

// ListSessions retrieves all research sessions
func (r *ResearchDB) ListSessions(limit, offset int) ([]models.ResearchSession, error) {
	var sessions []models.ResearchSession
	query := r.db.Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	return sessions, nil
}

// ListSessionsByStatus retrieves sessions filtered by status
func (r *ResearchDB) ListSessionsByStatus(status string, limit, offset int) ([]models.ResearchSession, error) {
	var sessions []models.ResearchSession
	query := r.db.Where("status = ?", status).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("failed to list sessions by status: %w", err)
	}

	return sessions, nil
}

// DeleteSession deletes a research session by ID
func (r *ResearchDB) DeleteSession(sessionID string) error {
	if err := r.db.Where("id = ?", sessionID).Delete(&models.ResearchSession{}).Error; err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

// DeleteOldSessions deletes sessions older than the specified age
func (r *ResearchDB) DeleteOldSessions(maxAge time.Duration) (int64, error) {
	cutoff := time.Now().Add(-maxAge)
	result := r.db.Where("created_at < ?", cutoff).Delete(&models.ResearchSession{})
	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete old sessions: %w", result.Error)
	}
	return result.RowsAffected, nil
}

// IndexDocument saves a document to the index
func (r *ResearchDB) IndexDocument(filename, filePath, fileType string, fileSize int64, pageCount int) (*models.Document, error) {
	doc := &models.Document{
		ID:         uuid.New().String(),
		Filename:   filename,
		FilePath:   filePath,
		FileType:   fileType,
		FileSize:   fileSize,
		PageCount:  pageCount,
		Indexed:    true,
		UploadedAt: time.Now(),
	}

	if err := r.db.Create(doc).Error; err != nil {
		return nil, fmt.Errorf("failed to index document: %w", err)
	}

	return doc, nil
}

// GetDocument retrieves a document by ID
func (r *ResearchDB) GetDocument(docID string) (*models.Document, error) {
	var doc models.Document
	if err := r.db.Where("id = ?", docID).First(&doc).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("document not found: %s", docID)
		}
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	return &doc, nil
}

// GetDocumentByFilename retrieves a document by filename
func (r *ResearchDB) GetDocumentByFilename(filename string) (*models.Document, error) {
	var doc models.Document
	if err := r.db.Where("filename = ?", filename).First(&doc).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("document not found: %s", filename)
		}
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	return &doc, nil
}

// ListDocuments retrieves all indexed documents
func (r *ResearchDB) ListDocuments(limit, offset int) ([]models.Document, error) {
	var docs []models.Document
	query := r.db.Order("uploaded_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&docs).Error; err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}

	return docs, nil
}

// DeleteDocument deletes a document by ID
func (r *ResearchDB) DeleteDocument(docID string) error {
	if err := r.db.Where("id = ?", docID).Delete(&models.Document{}).Error; err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	return nil
}

// GetSessionCount returns the total number of sessions
func (r *ResearchDB) GetSessionCount() (int64, error) {
	var count int64
	if err := r.db.Model(&models.ResearchSession{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count sessions: %w", err)
	}
	return count, nil
}

// GetDocumentCount returns the total number of documents
func (r *ResearchDB) GetDocumentCount() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Document{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}
	return count, nil
}

// Close closes the database connection
func (r *ResearchDB) Close() error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

// ParseSessionResult deserializes the result from a session
func ParseSessionResult(session *models.ResearchSession) (*models.ResearchResult, error) {
	if session.Result == "" {
		return nil, nil
	}

	var result models.ResearchResult
	if err := json.Unmarshal([]byte(session.Result), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return &result, nil
}

// ParseSessionQuery deserializes the query from a session
func ParseSessionQuery(session *models.ResearchSession) (*models.ResearchQuery, error) {
	if session.QueryData == "" {
		return nil, nil
	}

	var query models.ResearchQuery
	if err := json.Unmarshal([]byte(session.QueryData), &query); err != nil {
		return nil, fmt.Errorf("failed to unmarshal query: %w", err)
	}

	return &query, nil
}
