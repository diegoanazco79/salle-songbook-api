package mongo

import (
	"context"
	"errors"
	"time"

	"salle-songbook-api/configs"
	"salle-songbook-api/internal/core/song"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SongMongoRepository struct {
	collection *mongo.Collection
}

func NewSongMongoRepository() *SongMongoRepository {
	db := configs.AppConfig.Client.Database(configs.AppConfig.DatabaseName)
	collection := db.Collection("songs")

	// Crear índice único en title
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"title": 1},
		Options: options.Index().SetUnique(true),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		panic("Failed to create index on songs collection: " + err.Error())
	}

	return &SongMongoRepository{collection: collection}
}

func (r *SongMongoRepository) GetAll() ([]song.Song, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var songs []song.Song
	for cursor.Next(ctx) {
		var s song.Song
		if err := cursor.Decode(&s); err != nil {
			return nil, err
		}
		songs = append(songs, s)
	}

	return songs, nil
}

func (r *SongMongoRepository) GetByID(id string) (song.Song, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return song.Song{}, err
	}

	var s song.Song
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&s)
	if err != nil {
		return song.Song{}, errors.New("song not found")
	}

	return s, nil
}

func (r *SongMongoRepository) Create(s song.Song) (song.Song, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.ID = uuid.NewString()

	doc := bson.M{
		"id":                        s.ID,
		"title":                     s.Title,
		"lyrics":                    s.Lyrics,
		"lyrics_with_guitar_chords": s.LyricsWithGuitarChords,
		"author":                    s.Author,
	}

	res, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return song.Song{}, errors.New("song title already exists")
		}
		return song.Song{}, err
	}

	s.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return s, nil
}

func (r *SongMongoRepository) Update(id string, s song.Song) (song.Song, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return song.Song{}, err
	}

	update := bson.M{
		"$set": bson.M{
			"title":                     s.Title,
			"lyrics":                    s.Lyrics,
			"lyrics_with_guitar_chords": s.LyricsWithGuitarChords,
			"author":                    s.Author,
		},
	}

	result, err := r.collection.UpdateByID(ctx, objID, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return song.Song{}, errors.New("song title already exists")
		}
		return song.Song{}, err
	}
	if result.MatchedCount == 0 {
		return song.Song{}, errors.New("song not found")
	}

	s.ID = id
	return s, nil
}

func (r *SongMongoRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("song not found")
	}

	return nil
}
