package uploads

import (
	"fmt"
	"log"
	"mime/multipart"
	"time"

	"gorm.io/gorm"

	"deep-research-agent/agent"
	"deep-research-agent/models"
)

// UploadHandler handles file upload and processing operations
type UploadHandler struct {
	storage       *FileStorage
	db            *gorm.DB
	researchAgent *agent.ResearchAgent
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(db *gorm.DB, agent *agent.ResearchAgent, uploadDir string) (*UploadHandler, error) {
	storage, err := NewFileStorage(db, uploadDir)
	if err != nil {
		return nil, err
	}

	return &UploadHandler{
		storage:       storage,
		db:            db,
		researchAgent: agent,
	}, nil
}

// HandleUpload processes a file upload
func (h *UploadHandler) HandleUpload(fileHeader *multipart.FileHeader) (*models.UploadedFile, error) {
	log.Printf("Processing file upload: %s (size: %d bytes)", fileHeader.Filename, fileHeader.Size)

	// Save file to storage
	uploadedFile, err := h.storage.SaveFile(fileHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	log.Printf("File uploaded successfully: %s (ID: %s)", fileHeader.Filename, uploadedFile.ID)
	return uploadedFile, nil
}

// ProcessFile processes an uploaded file (indexes it for search)
func (h *UploadHandler) ProcessFile(fileID string) error {
	// Get file record
	file, err := h.storage.GetFile(fileID)
	if err != nil {
		return err
	}

	// Update status to processing
	if err := h.storage.UpdateFileStatus(fileID, models.FileStatusProcessing, ""); err != nil {
		return err
	}

	log.Printf("Processing file: %s", file.OriginalName)

	// Index the document in the research agent's storage
	var pageCount int
	switch GetFileExtension(file.OriginalName) {
	case ".pdf":
		pageCount = 10 // Placeholder - would need actual PDF page count
	case ".docx":
		pageCount = 5 // Placeholder - would need actual DOCX page count
	default:
		pageCount = 1
	}

	doc, err := h.researchAgent.IndexDocument(
		file.OriginalName,
		file.Path,
		GetFileExtension(file.OriginalName)[1:], // Remove leading dot
		file.Size,
		pageCount,
	)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to index document: %v", err)
		h.storage.UpdateFileStatus(fileID, models.FileStatusFailed, errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// Update status to indexed
	now := time.Now()
	file.Status = models.FileStatusIndexed
	file.ProcessedAt = &now

	if err := h.db.Save(file).Error; err != nil {
		return fmt.Errorf("failed to update file status: %w", err)
	}

	log.Printf("File processed successfully: %s (Document ID: %s)", file.OriginalName, doc.ID)
	return nil
}

// GetFile retrieves a file by ID
func (h *UploadHandler) GetFile(fileID string) (*models.UploadedFile, error) {
	return h.storage.GetFile(fileID)
}

// ListFiles lists all uploaded files
func (h *UploadHandler) ListFiles(limit, offset int) ([]models.UploadedFile, error) {
	return h.storage.ListFiles(limit, offset)
}

// DeleteFile deletes an uploaded file
func (h *UploadHandler) DeleteFile(fileID string) error {
	log.Printf("Deleting file: %s", fileID)
	return h.storage.DeleteFile(fileID)
}

// GetFilePath returns the path to an uploaded file
func (h *UploadHandler) GetFilePath(fileID string) (string, error) {
	return h.storage.GetFilePath(fileID)
}
