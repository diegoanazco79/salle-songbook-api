package memory

import (
	"errors"
	"salle-songbook-api/internal/core/review"

	"github.com/google/uuid"
)

type ReviewRepository struct {
	reviews map[string]review.PendingReview
}

func NewReviewRepository() *ReviewRepository {
	return &ReviewRepository{
		reviews: make(map[string]review.PendingReview),
	}
}

func (r *ReviewRepository) Create(pr review.PendingReview) (review.PendingReview, error) {
	pr.ID = uuid.NewString()
	r.reviews[pr.ID] = pr
	return pr, nil
}

func (r *ReviewRepository) GetAll() ([]review.PendingReview, error) {
	result := make([]review.PendingReview, 0, len(r.reviews))
	for _, pr := range r.reviews {
		result = append(result, pr)
	}
	return result, nil
}

func (r *ReviewRepository) GetByID(id string) (review.PendingReview, error) {
	pr, ok := r.reviews[id]
	if !ok {
		return review.PendingReview{}, errors.New("pending review not found")
	}
	return pr, nil
}

func (r *ReviewRepository) Delete(id string) error {
	_, ok := r.reviews[id]
	if !ok {
		return errors.New("pending review not found")
	}
	delete(r.reviews, id)
	return nil
}
