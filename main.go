package main

import (
	"log"
	v1 "salle-songbook-api/internal/ports/api/http/v1"

	"github.com/gin-gonic/gin"

	"salle-songbook-api/configs"
	"salle-songbook-api/internal/ports/api/http/middleware"
	"salle-songbook-api/internal/ports/repository/memory"
	"salle-songbook-api/internal/ports/repository/mongo"
)

func main() {
	// Cargamos configuración
	configs.LoadConfig()

	r := gin.Default()
	r.Use(CORSMiddleware())

	// Inicializamos repositorios
	songRepo := mongo.NewSongMongoRepository()
	reviewRepo := mongo.NewReviewMongoRepository()
	userRepo := memory.NewUserRepository() // Usuarios en memoria todavía

	// Inicializamos handlers
	authHandler := v1.NewAuthHandler(userRepo)
	songHandler := v1.NewSongHandler(songRepo, reviewRepo)
	reviewHandler := v1.NewReviewHandler(reviewRepo, songRepo)

	// Definimos rutas
	api := r.Group("/api/v1")
	{
		api.POST("/login", authHandler.Login)

		// Rutas públicas para canciones (sin middleware)
		api.GET("/songs", songHandler.GetAll)
		api.GET("/songs/:id", songHandler.GetByID)

		// Rutas protegidas para canciones
		songs := api.Group("/songs")
		songs.Use(middleware.AuthMiddleware(), middleware.AdminOrComposerMiddleware())
		{
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

	// Iniciamos el servidor HTTP en el puerto 8080
	log.Println("Server starting on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
