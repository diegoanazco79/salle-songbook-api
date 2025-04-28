package mongo

import (
	"context"
	"errors"
	"time"

	"salle-songbook-api/configs"
	"salle-songbook-api/internal/core/review"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewMongoRepository struct {
	collection *mongo.Collection
}

func NewReviewMongoRepository() *ReviewMongoRepository {
	db := configs.AppConfig.Client.Database(configs.AppConfig.DatabaseName)
	collection := db.Collection("pending_reviews")
	return &ReviewMongoRepository{collection: collection}
}

func (r *ReviewMongoRepository) GetAll() ([]review.PendingReview, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reviews []review.PendingReview
	for cursor.Next(ctx) {
		var pr review.PendingReview
		if err := cursor.Decode(&pr); err != nil {
			return nil, err
		}
		reviews = append(reviews, pr)
	}

	return reviews, nil
}

func (r *ReviewMongoRepository) GetByID(id string) (review.PendingReview, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var pr review.PendingReview
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&pr)
	if err != nil {
		return review.PendingReview{}, errors.New("pending review not found")
	}

	return pr, nil
}

func (r *ReviewMongoRepository) Create(pr review.PendingReview) (review.PendingReview, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pr.ID = uuid.NewString()

	doc := bson.M{
		"id":            pr.ID,
		"action":        pr.Action,
		"song_id":       pr.SongID,
		"new_song_data": pr.NewSongData,
		"requested_by":  pr.RequestedBy,
	}

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return review.PendingReview{}, err
	}

	return pr, nil
}

func (r *ReviewMongoRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("pending review not found")
	}

	return nil
}
