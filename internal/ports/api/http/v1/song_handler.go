package v1

import (
	"net/http"

	"salle-songbook-api/internal/core/review"
	"salle-songbook-api/internal/core/song"
	"salle-songbook-api/internal/ports/repository/memory"
	"salle-songbook-api/pkg/response"

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
	response.Success(c, songs, "List of songs retrieved")
}

func (h *SongHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	s, err := h.songRepo.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Song not found", err.Error())
		return
	}
	response.Success(c, s, "Song retrieved")
}

func (h *SongHandler) Create(c *gin.Context) {
	role := c.GetString("role")
	username := c.GetString("username")

	var s song.Song
	if err := c.ShouldBindJSON(&s); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if role == "admin" {
		created, _ := h.songRepo.Create(s)
		response.Created(c, created, "Song created successfully")
	} else {
		pr := review.PendingReview{
			Action: review.Create,
			NewSongData: &review.SongDataDTO{
				Title:                  s.Title,
				Lyrics:                 s.Lyrics,
				LyricsWithGuitarChords: s.LyricsWithGuitarChords,
				Author:                 s.Author,
			},
			RequestedBy: username,
		}
		pending, _ := h.reviewRepo.Create(pr)
		response.Created(c, pending, "Song creation request submitted for review")
	}
}

func (h *SongHandler) Update(c *gin.Context) {
	role := c.GetString("role")
	username := c.GetString("username")
	id := c.Param("id")

	var s song.Song
	if err := c.ShouldBindJSON(&s); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if role == "admin" {
		updated, err := h.songRepo.Update(id, s)
		if err != nil {
			response.Error(c, http.StatusNotFound, "Song not found", err.Error())
			return
		}
		response.Success(c, updated, "Song updated successfully")
	} else {
		pr := review.PendingReview{
			Action: review.Update,
			SongID: id,
			NewSongData: &review.SongDataDTO{
				Title:                  s.Title,
				Lyrics:                 s.Lyrics,
				LyricsWithGuitarChords: s.LyricsWithGuitarChords,
				Author:                 s.Author,
			},
			RequestedBy: username,
		}
		pending, _ := h.reviewRepo.Create(pr)
		response.Created(c, pending, "Song update request submitted for review")
	}
}

func (h *SongHandler) Delete(c *gin.Context) {
	role := c.GetString("role")
	username := c.GetString("username")
	id := c.Param("id")

	if role == "admin" {
		err := h.songRepo.Delete(id)
		if err != nil {
			response.Error(c, http.StatusNotFound, "Song not found", err.Error())
			return
		}
		response.Success(c, nil, "Song deleted successfully")
	} else {
		pr := review.PendingReview{
			Action:      review.Delete,
			SongID:      id,
			RequestedBy: username,
		}
		pending, _ := h.reviewRepo.Create(pr)
		response.Created(c, pending, "Song deletion request submitted for review")
	}
}
