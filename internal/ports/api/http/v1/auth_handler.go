package v1

import (
	"net/http"
	"salle-songbook-api/internal/ports/repository/memory"
	"salle-songbook-api/pkg/response"
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
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.userRepo.GetByUsername(req.Username)
	if err != nil || user.Password != req.Password {
		response.Error(c, http.StatusUnauthorized, "Invalid credentials", "unauthorized")
		return
	}

	tokenString, expiresAt, err := token.GenerateToken(user.Username, string(user.Role))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Could not generate token", err.Error())
		return
	}

	response.Success(c, gin.H{
		"token":      tokenString,
		"expires_at": expiresAt.Unix(),
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	}, "Login successful")
}
