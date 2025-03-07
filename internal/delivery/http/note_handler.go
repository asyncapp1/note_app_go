package http

import (
	"net/http"
	"strconv"

	"notes-app/internal/domain"

	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	noteUsecase domain.NoteUsecase
}

func NewNoteHandler(r *gin.RouterGroup, nu domain.NoteUsecase) {
	handler := &NoteHandler{
		noteUsecase: nu,
	}

	r.POST("/notes", handler.Create)
	r.GET("/notes", handler.GetAll)
	r.GET("/notes/:id", handler.GetByID)
	r.PUT("/notes/:id", handler.Update)
	r.DELETE("/notes/:id", handler.Delete)
	r.GET("/notes/query/:query", handler.Query)
}

func (h *NoteHandler) Create(c *gin.Context) {
	var note domain.Note

	// Get user from context and type assert
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	// Type assert the user to domain.User
	userObj, ok := user.(*domain.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type in context"})
		return
	}

	// Bind JSON body to note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.noteUsecase.Create(&note, userObj); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, note)
}

func (h *NoteHandler) GetAll(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*domain.User)

	notes, err := h.noteUsecase.GetAll(userObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notes)
}

func (h *NoteHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user, _ := c.Get("user")
	userObj := user.(*domain.User)

	note, err := h.noteUsecase.GetByID(uint(id), userObj)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var note domain.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("user")
	userObj := user.(*domain.User)

	note.ID = uint(id)
	if err := h.noteUsecase.Update(&note, userObj); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user, _ := c.Get("user")
	userObj := user.(*domain.User)

	if err := h.noteUsecase.Delete(uint(id), userObj); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}

func (h *NoteHandler) Query(c *gin.Context) {
	query := c.Param("query")

	user, _ := c.Get("user")
	userObj := user.(*domain.User)

	notes, err := h.noteUsecase.Query(query, userObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notes)
}
