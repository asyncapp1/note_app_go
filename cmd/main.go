package main

import (
	"log"

	"notes-app/internal/delivery/http"
	"notes-app/internal/delivery/http/middleware"
	"notes-app/internal/repository"
	"notes-app/internal/usecase"
	"notes-app/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	db := database.NewPostgresDB()

	// Initialize services
	migrationService, err := database.NewMigrationService(db)
	if err != nil {
		log.Fatal("Failed to initialize migration service:", err)
	}

	// Repositories
	noteRepo := repository.NewNoteRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Usecases
	noteUsecase := usecase.NewNoteUsecase(noteRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)

	r := gin.Default()

	// Public routes group
	public := r.Group("")
	http.NewAuthHandler(public, userUsecase)

	// Protected routes
	protected := r.Group("")
	protected.Use(middleware.AuthMiddleware(userUsecase))
	{
		http.NewNoteHandler(protected, noteUsecase)
		http.NewMigrationHandler(protected, migrationService)
	}

	if err := r.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
