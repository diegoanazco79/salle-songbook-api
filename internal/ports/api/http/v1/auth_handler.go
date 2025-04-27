package v1

import (
	"net/http"
	"salle-songbook-api/internal/ports/repository/memory"
	"salle-songbook-api/pkg/token"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userRepo *memory.UserRepository
}

func NewAuthHandler(repo *memory.UserRepository) *AuthHandler {
	return &AuthHandler{userRepo: repo}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetByUsername(req.Username)
	if err != nil || user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	tokenString, err := token.GenerateToken(user.Username, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
