package v1

import (
	"net/http"
	"salle-songbook-api/internal/core/review"
	"salle-songbook-api/internal/core/song"
	"salle-songbook-api/internal/ports/repository/memory"
	"salle-songbook-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewRepo *memory.ReviewRepository
	songRepo   *memory.SongRepository
}

func NewReviewHandler(reviewRepo *memory.ReviewRepository, songRepo *memory.SongRepository) *ReviewHandler {
	return &ReviewHandler{reviewRepo: reviewRepo, songRepo: songRepo}
}

func (h *ReviewHandler) GetAllPendingReviews(c *gin.Context) {
	reviews, _ := h.reviewRepo.GetAll()
	response.Success(c, reviews, "List of pending reviews retrieved")
}

func (h *ReviewHandler) ApproveReview(c *gin.Context) {
	id := c.Param("id")

	pr, err := h.reviewRepo.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Pending review not found", err.Error())
		return
	}

	switch pr.Action {
	case review.Create:
		newSong := song.Song{
			Title:                  pr.NewSongData.Title,
			Lyrics:                 pr.NewSongData.Lyrics,
			LyricsWithGuitarChords: pr.NewSongData.LyricsWithGuitarChords,
			Author:                 pr.NewSongData.Author,
		}
		h.songRepo.Create(newSong)
	case review.Update:
		updatedSong := song.Song{
			Title:                  pr.NewSongData.Title,
			Lyrics:                 pr.NewSongData.Lyrics,
			LyricsWithGuitarChords: pr.NewSongData.LyricsWithGuitarChords,
			Author:                 pr.NewSongData.Author,
		}
		h.songRepo.Update(pr.SongID, updatedSong)
	case review.Delete:
		h.songRepo.Delete(pr.SongID)
	}

	h.reviewRepo.Delete(id)
	response.Success(c, nil, "Review approved and action executed")
}

func (h *ReviewHandler) RejectReview(c *gin.Context) {
	id := c.Param("id")

	err := h.reviewRepo.Delete(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Pending review not found", err.Error())
		return
	}

	response.Success(c, nil, "Review rejected and deleted")
}
