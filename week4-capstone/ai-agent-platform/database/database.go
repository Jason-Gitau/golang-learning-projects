package database

import (
	"ai-agent-platform/models"
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Initialize sets up the database connection and runs migrations
func Initialize(dbPath string, isDevelopment bool) error {
	var err error

	// Configure GORM logger
	logLevel := logger.Silent
	if isDevelopment {
		logLevel = logger.Info
	}

	// Open database connection
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Enable foreign keys for SQLite
	if err := DB.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Run migrations
	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// runMigrations runs all database migrations
func runMigrations() error {
	// AutoMigrate all models
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Agent{},
		&models.Conversation{},
		&models.Message{},
		&models.UsageLog{},
		&models.RateLimit{},
	); err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	// Create indexes for better performance
	if err := createIndexes(); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}

// createIndexes creates additional indexes for performance optimization
func createIndexes() error {
	// User indexes
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_tier ON users(tier)").Error; err != nil {
		return err
	}

	// Agent indexes
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_agents_user_id_status ON agents(user_id, status)").Error; err != nil {
		return err
	}

	// Conversation indexes
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_conversations_user_id_status ON conversations(user_id, status)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_conversations_agent_id_status ON conversations(agent_id, status)").Error; err != nil {
		return err
	}

	// Message indexes
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_messages_conversation_id_created_at ON messages(conversation_id, created_at)").Error; err != nil {
		return err
	}

	// UsageLog indexes
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_usage_logs_user_id_created_at ON usage_logs(user_id, created_at)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_usage_logs_agent_id_created_at ON usage_logs(agent_id, created_at)").Error; err != nil {
		return err
	}

	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// Close closes the database connection
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	return sqlDB.Close()
}
