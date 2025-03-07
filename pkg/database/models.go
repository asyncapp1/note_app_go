package database

import (
	"time"
)

// MigrationHistory tracks all database migrations and changes
type MigrationHistory struct {
	ID            uint      `gorm:"primaryKey"`
	TableName     string    `gorm:"not null"`
	Operation     string    `gorm:"not null"` // CREATE, ALTER, DROP, etc.
	Description   string    `gorm:"not null"`
	SchemaChanges string    `gorm:"type:text"` // JSON string of changes
	ExecutedAt    time.Time `gorm:"not null"`
	Version       string    `gorm:"not null"`
	Status        string    `gorm:"not null"` // SUCCESS, FAILED
	ErrorMessage  string    `gorm:"type:text"`
}
