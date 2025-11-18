package uploads

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"deep-research-agent/models"
)

const (
	// MaxFileSize is the maximum allowed file size (10MB)
	MaxFileSize = 10 * 1024 * 1024

	// UploadDir is the directory where uploaded files are stored
	UploadDir = "uploads/files"
)

// FileStorage manages file upload storage
type FileStorage struct {
	db        *gorm.DB
	uploadDir string
}

// NewFileStorage creates a new file storage manager
func NewFileStorage(db *gorm.DB, uploadDir string) (*FileStorage, error) {
	if uploadDir == "" {
		uploadDir = UploadDir
	}

	// Ensure upload directory exists
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	return &FileStorage{
		db:        db,
		uploadDir: uploadDir,
	}, nil
}

// SaveFile saves an uploaded file to disk and database
func (s *FileStorage) SaveFile(fileHeader *multipart.FileHeader) (*models.UploadedFile, error) {
	// Validate file size
	if fileHeader.Size > MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", MaxFileSize)
	}

	// Validate file type
	ext := filepath.Ext(fileHeader.Filename)
	mimeType := fileHeader.Header.Get("Content-Type")

	if ext != ".pdf" && ext != ".docx" {
		return nil, fmt.Errorf("unsupported file type: %s (only PDF and DOCX are supported)", ext)
	}

	// Generate unique filename
	fileID := uuid.New().String()
	filename := fmt.Sprintf("%s%s", fileID, ext)
	filePath := filepath.Join(s.uploadDir, filename)

	// Open the uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy the file contents
	if _, err := io.Copy(dst, src); err != nil {
		os.Remove(filePath) // Clean up on error
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Create database record
	uploadedFile := &models.UploadedFile{
		ID:           fileID,
		Filename:     filename,
		OriginalName: fileHeader.Filename,
		Size:         fileHeader.Size,
		MimeType:     mimeType,
		Path:         filePath,
		Status:       models.FileStatusUploaded,
		UploadedAt:   time.Now(),
	}

	if err := s.db.Create(uploadedFile).Error; err != nil {
		os.Remove(filePath) // Clean up on error
		return nil, fmt.Errorf("failed to save file record: %w", err)
	}

	return uploadedFile, nil
}

// GetFile retrieves a file record by ID
func (s *FileStorage) GetFile(fileID string) (*models.UploadedFile, error) {
	var file models.UploadedFile
	if err := s.db.Where("id = ?", fileID).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("file not found: %s", fileID)
		}
		return nil, fmt.Errorf("failed to get file: %w", err)
	}
	return &file, nil
}

// ListFiles retrieves all uploaded files
func (s *FileStorage) ListFiles(limit, offset int) ([]models.UploadedFile, error) {
	var files []models.UploadedFile
	query := s.db.Order("uploaded_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&files).Error; err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return files, nil
}

// DeleteFile deletes a file from disk and database
func (s *FileStorage) DeleteFile(fileID string) error {
	// Get file record
	file, err := s.GetFile(fileID)
	if err != nil {
		return err
	}

	// Delete file from disk
	if err := os.Remove(file.Path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file from disk: %w", err)
	}

	// Delete database record
	if err := s.db.Delete(&models.UploadedFile{}, "id = ?", fileID).Error; err != nil {
		return fmt.Errorf("failed to delete file record: %w", err)
	}

	return nil
}

// UpdateFileStatus updates the status of an uploaded file
func (s *FileStorage) UpdateFileStatus(fileID, status, errorMsg string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if errorMsg != "" {
		updates["error"] = errorMsg
	}

	if status == models.FileStatusIndexed || status == models.FileStatusFailed {
		now := time.Now()
		updates["processed_at"] = &now
	}

	if err := s.db.Model(&models.UploadedFile{}).
		Where("id = ?", fileID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update file status: %w", err)
	}

	return nil
}

// GetFilePath returns the full path to an uploaded file
func (s *FileStorage) GetFilePath(fileID string) (string, error) {
	file, err := s.GetFile(fileID)
	if err != nil {
		return "", err
	}
	return file.Path, nil
}

// ValidateFileType checks if a file type is supported
func ValidateFileType(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".pdf" || ext == ".docx"
}

// GetFileExtension returns the file extension
func GetFileExtension(filename string) string {
	return filepath.Ext(filename)
}
