package database

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type MigrationService struct {
	db *gorm.DB
}

type SchemaChange struct {
	ColumnName   string `json:"column_name,omitempty"`
	OldType      string `json:"old_type,omitempty"`
	NewType      string `json:"new_type,omitempty"`
	IsNew        bool   `json:"is_new"`
	IsDeleted    bool   `json:"is_deleted"`
	DefaultValue string `json:"default_value,omitempty"`
}

// MigrationResponse is used for JSON response
type MigrationResponse struct {
	ID            uint           `json:"id"`
	TableName     string         `json:"table_name"`
	Operation     string         `json:"operation"`
	Description   string         `json:"description"`
	SchemaChanges []SchemaChange `json:"schema_changes"` // Changed from string to []SchemaChange
	ExecutedAt    time.Time      `json:"executed_at"`
	Version       string         `json:"version"`
	Status        string         `json:"status"`
	ErrorMessage  string         `json:"error_message"`
}

func NewMigrationService(db *gorm.DB) (*MigrationService, error) {
	if err := db.AutoMigrate(&MigrationHistory{}); err != nil {
		return nil, fmt.Errorf("failed to create migration history table: %v", err)
	}
	return &MigrationService{db: db}, nil
}

func (s *MigrationService) TrackMigration(tableName, operation, description string, changes []SchemaChange) error {
	// Check if similar migration exists
	var existing MigrationHistory
	err := s.db.Where("table_name = ? AND operation = ? AND description = ?",
		tableName, operation, description).First(&existing).Error

	if err == nil {
		// Migration already exists, skip
		return nil
	}

	// Convert changes to JSON string
	changesJSON, err := json.Marshal(changes)
	if err != nil {
		return fmt.Errorf("failed to marshal schema changes: %v", err)
	}

	migration := MigrationHistory{
		TableName:     tableName,
		Operation:     operation,
		Description:   description,
		SchemaChanges: string(changesJSON),
		ExecutedAt:    time.Now(),
		Version:       time.Now().Format("20060102150405"),
		Status:        "SUCCESS",
	}

	tx := s.db.Begin()
	if err := tx.Create(&migration).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create migration history: %v", err)
	}

	return tx.Commit().Error
}

func (s *MigrationService) GetMigrationHistory() ([]MigrationResponse, error) {
	var history []MigrationHistory
	err := s.db.Order("executed_at desc").Find(&history).Error
	if err != nil {
		return nil, err
	}

	// Convert to response format with parsed JSON
	var response []MigrationResponse
	for _, h := range history {
		var changes []SchemaChange
		if err := json.Unmarshal([]byte(h.SchemaChanges), &changes); err != nil {
			return nil, fmt.Errorf("failed to parse schema changes: %v", err)
		}

		response = append(response, MigrationResponse{
			ID:            h.ID,
			TableName:     h.TableName,
			Operation:     h.Operation,
			Description:   h.Description,
			SchemaChanges: changes,
			ExecutedAt:    h.ExecutedAt,
			Version:       h.Version,
			Status:        h.Status,
			ErrorMessage:  h.ErrorMessage,
		})
	}

	return response, nil
}

func (s *MigrationService) GetLatestMigration() (*MigrationResponse, error) {
	var latest MigrationHistory
	err := s.db.Order("executed_at desc").First(&latest).Error
	if err != nil {
		return nil, err
	}

	var changes []SchemaChange
	if err := json.Unmarshal([]byte(latest.SchemaChanges), &changes); err != nil {
		return nil, fmt.Errorf("failed to parse schema changes: %v", err)
	}

	response := &MigrationResponse{
		ID:            latest.ID,
		TableName:     latest.TableName,
		Operation:     latest.Operation,
		Description:   latest.Description,
		SchemaChanges: changes,
		ExecutedAt:    latest.ExecutedAt,
		Version:       latest.Version,
		Status:        latest.Status,
		ErrorMessage:  latest.ErrorMessage,
	}

	return response, nil
}
