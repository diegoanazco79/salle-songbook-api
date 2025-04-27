package main

import (
	v1 "salle-songbook-api/internal/ports/api/http/v1"
	"salle-songbook-api/internal/ports/repository/memory"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	userRepo := memory.NewUserRepository()
	authHandler := v1.NewAuthHandler(userRepo)

	api := r.Group("/api/v1")
	{
		api.POST("/login", authHandler.Login)
	}

	r.Run(":8080")
}
