package main

import (
	"context"
	v1 "salle-songbook-api/internal/ports/api/http/v1"

	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"salle-songbook-api/configs"
	"salle-songbook-api/internal/ports/api/http/middleware"
	"salle-songbook-api/internal/ports/repository/memory"
	"salle-songbook-api/internal/ports/repository/mongo"
)

var ginLambda *ginadapter.GinLambda

func init() {
	// Cargamos configuración
	configs.LoadConfig()

	r := gin.Default()

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

	// Adaptamos Gin a AWS Lambda handler
	ginLambda = ginadapter.New(r)
}

func main() {
	// Ejecutamos como lambda handler
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}
