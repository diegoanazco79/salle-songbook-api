package v1

import (
	"net/http"

	"salle-songbook-api/internal/core/review"
	"salle-songbook-api/internal/core/song"
	"salle-songbook-api/internal/ports/repository/memory"

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
	c.JSON(http.StatusOK, reviews)
}

func (h *ReviewHandler) ApproveReview(c *gin.Context) {
	id := c.Param("id")

	pr, err := h.reviewRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pending review not found"})
		return
	}

	switch pr.Action {
	case review.Create:
		songData := pr.NewSongData.(map[string]interface{})
		newSong := song.Song{
			Title:                  songData["title"].(string),
			Lyrics:                 songData["lyrics"].(string),
			LyricsWithGuitarChords: songData["lyrics_with_guitar_chords"].(string),
			Author:                 songData["author"].(string),
		}
		h.songRepo.Create(newSong)
	case review.Update:
		songData := pr.NewSongData.(map[string]interface{})
		updatedSong := song.Song{
			Title:                  songData["title"].(string),
			Lyrics:                 songData["lyrics"].(string),
			LyricsWithGuitarChords: songData["lyrics_with_guitar_chords"].(string),
			Author:                 songData["author"].(string),
		}
		h.songRepo.Update(pr.SongID, updatedSong)
	case review.Delete:
		h.songRepo.Delete(pr.SongID)
	}

	h.reviewRepo.Delete(id)
	c.JSON(http.StatusOK, gin.H{"message": "Review approved and action executed"})
}

func (h *ReviewHandler) RejectReview(c *gin.Context) {
	id := c.Param("id")

	err := h.reviewRepo.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pending review not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review rejected and deleted"})
}
