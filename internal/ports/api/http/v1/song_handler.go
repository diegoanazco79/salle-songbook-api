package v1

import (
	"net/http"

	"salle-songbook-api/internal/core/review"
	"salle-songbook-api/internal/core/song"
	"salle-songbook-api/internal/ports/repository/memory"

	"github.com/gin-gonic/gin"
)

type SongHandler struct {
	songRepo   *memory.SongRepository
	reviewRepo *memory.ReviewRepository
}

func NewSongHandler(songRepo *memory.SongRepository, reviewRepo *memory.ReviewRepository) *SongHandler {
	return &SongHandler{songRepo: songRepo, reviewRepo: reviewRepo}
}

func (h *SongHandler) GetAll(c *gin.Context) {
	songs, _ := h.songRepo.GetAll()
	c.JSON(http.StatusOK, songs)
}

func (h *SongHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	s, err := h.songRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *SongHandler) Create(c *gin.Context) {
	role := c.GetString("role")
	username := c.GetString("username")

	var s song.Song
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if role == "admin" {
		created, _ := h.songRepo.Create(s)
		c.JSON(http.StatusCreated, created)
	} else {
		pr := review.PendingReview{
			Action:      review.Create,
			NewSongData: s,
			RequestedBy: username,
		}
		pending, _ := h.reviewRepo.Create(pr)
		c.JSON(http.StatusAccepted, pending)
	}
}

func (h *SongHandler) Update(c *gin.Context) {
	role := c.GetString("role")
	username := c.GetString("username")
	id := c.Param("id")

	var s song.Song
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if role == "admin" {
		updated, err := h.songRepo.Update(id, s)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}
		c.JSON(http.StatusOK, updated)
	} else {
		pr := review.PendingReview{
			Action:      review.Update,
			SongID:      id,
			NewSongData: s,
			RequestedBy: username,
		}
		pending, _ := h.reviewRepo.Create(pr)
		c.JSON(http.StatusAccepted, pending)
	}
}

func (h *SongHandler) Delete(c *gin.Context) {
	role := c.GetString("role")
	username := c.GetString("username")
	id := c.Param("id")

	if role == "admin" {
		err := h.songRepo.Delete(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
	} else {
		pr := review.PendingReview{
			Action:      review.Delete,
			SongID:      id,
			RequestedBy: username,
		}
		pending, _ := h.reviewRepo.Create(pr)
		c.JSON(http.StatusAccepted, pending)
	}
}
