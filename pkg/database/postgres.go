package database

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"notes-app/internal/domain"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Initialize migration service
	migrationService, err := NewMigrationService(db)
	if err != nil {
		log.Fatal("Failed to initialize migration service:", err)
	}

	// Auto migrate both User and Note models
	err = db.AutoMigrate(&domain.User{}, &domain.Note{})
	if err != nil {
		log.Fatal(err)
	}

	// Get current schema from both models using reflection
	modelTypes := []reflect.Type{
		reflect.TypeOf(domain.User{}),
		reflect.TypeOf(domain.Note{}),
	}

	for _, modelType := range modelTypes {
		tableName := db.NamingStrategy.TableName(modelType.Name())

		// Check if table exists
		var tableExists bool
		err = db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = ?)", tableName).Scan(&tableExists).Error
		if err != nil {
			log.Fatal("Failed to check table existence:", err)
		}

		if !tableExists {
			continue
		}

		// Track schema changes for each model
		currentSchema := make(map[string]string)
		for i := 0; i < modelType.NumField(); i++ {
			field := modelType.Field(i)
			currentSchema[field.Tag.Get("json")] = getFieldType(field.Type)
		}

		// Get existing columns from database
		var columns []struct {
			ColumnName string
			DataType   string
		}
		err = db.Raw(`
			SELECT column_name, data_type 
			FROM information_schema.columns 
			WHERE table_name = ?
		`, tableName).Scan(&columns).Error
		if err != nil {
			log.Fatal("Failed to get column information:", err)
		}

		// Create map of existing columns
		existingColumns := make(map[string]string)
		for _, col := range columns {
			existingColumns[col.ColumnName] = col.DataType
		}

		// Track only new or changed columns
		var changes []SchemaChange
		for columnName, dataType := range currentSchema {
			if existingType, exists := existingColumns[columnName]; !exists {
				// New column
				changes = append(changes, SchemaChange{
					ColumnName: columnName,
					NewType:    dataType,
					IsNew:      true,
				})
			} else if existingType != dataType {
				// Changed column type
				changes = append(changes, SchemaChange{
					ColumnName: columnName,
					OldType:    existingType,
					NewType:    dataType,
					IsNew:      false,
				})
			}
		}

		// Only track migration if there are actual changes
		if len(changes) > 0 {
			err = migrationService.TrackMigration(
				tableName,
				"ALTER",
				fmt.Sprintf("Schema changes detected: %d modifications", len(changes)),
				changes,
			)
			if err != nil {
				log.Printf("Warning: Failed to track migration: %v", err)
			}
		}
	}

	return db
}

// Helper function to convert Go types to database types
func getFieldType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "numeric"
	case reflect.String:
		return "text"
	default:
		if t.String() == "time.Time" {
			return "timestamp"
		}
		return "text"
	}
}
