package v1

import (
	"net/http"
	"salle-songbook-api/internal/core/song"
	"salle-songbook-api/internal/ports/repository/memory"

	"github.com/gin-gonic/gin"
)

type SongHandler struct {
	repo *memory.SongRepository
}

func NewSongHandler(repo *memory.SongRepository) *SongHandler {
	return &SongHandler{repo: repo}
}

func (h *SongHandler) GetAll(c *gin.Context) {
	songs, _ := h.repo.GetAll()
	c.JSON(http.StatusOK, songs)
}

func (h *SongHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	s, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *SongHandler) Create(c *gin.Context) {
	var s song.Song
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created, _ := h.repo.Create(s)
	c.JSON(http.StatusCreated, created)
}

func (h *SongHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var s song.Song
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := h.repo.Update(id, s)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *SongHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.repo.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}
