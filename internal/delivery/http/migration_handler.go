package http

import (
	"net/http"
	"notes-app/pkg/database"

	"github.com/gin-gonic/gin"
)

type MigrationHandler struct {
	migrationService *database.MigrationService
}

func NewMigrationHandler(r *gin.RouterGroup, ms *database.MigrationService) {
	handler := &MigrationHandler{
		migrationService: ms,
	}

	// Register routes
	r.GET("/migrations", handler.GetMigrationHistory)
	r.GET("/migrations/latest", handler.GetLatestMigration)
}

func (h *MigrationHandler) GetMigrationHistory(c *gin.Context) {
	history, err := h.migrationService.GetMigrationHistory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch migration history: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_migrations": len(history),
		"migrations":       history,
	})
}

func (h *MigrationHandler) GetLatestMigration(c *gin.Context) {
	latest, err := h.migrationService.GetLatestMigration()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch latest migration: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, latest)
}
