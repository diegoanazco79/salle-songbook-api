package main

import (
	"salle-songbook-api/configs"
	"salle-songbook-api/internal/ports/api/http/middleware"
	"salle-songbook-api/internal/ports/repository/mongo"

	v1 "salle-songbook-api/internal/ports/api/http/v1"
	"salle-songbook-api/internal/ports/repository/memory"

	"github.com/gin-gonic/gin"
)

func main() {
	configs.LoadConfig() // 👈 importantísimo cargar configuración

	r := gin.Default()

	songRepo := mongo.NewSongMongoRepository()
	reviewRepo := mongo.NewReviewMongoRepository()
	userRepo := memory.NewUserRepository() // los users siguen en memoria por ahora

	authHandler := v1.NewAuthHandler(userRepo)
	songHandler := v1.NewSongHandler(songRepo, reviewRepo)
	reviewHandler := v1.NewReviewHandler(reviewRepo, songRepo)

	api := r.Group("/api/v1")
	{
		api.POST("/login", authHandler.Login)

		songs := api.Group("/songs")
		songs.Use(middleware.AuthMiddleware(), middleware.AdminOrComposerMiddleware())
		{
			songs.GET("", songHandler.GetAll)
			songs.GET("/:id", songHandler.GetByID)
			songs.POST("", songHandler.Create)
			songs.PUT("/:id", songHandler.Update)
			songs.DELETE("/:id", songHandler.Delete)
		}

		pendingReviews := api.Group("/pending-reviews")
		pendingReviews.Use(middleware.AuthMiddleware(), middleware.OnlyAdminMiddleware())
		{
			pendingReviews.GET("", reviewHandler.GetAllPendingReviews)
			pendingReviews.POST("/:id/approve", reviewHandler.ApproveReview)
			pendingReviews.POST("/:id/reject", reviewHandler.RejectReview)
		}
	}

	r.Run(":8080")
}
