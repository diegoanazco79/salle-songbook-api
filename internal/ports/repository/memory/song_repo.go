package memory

import (
	"errors"
	"salle-songbook-api/internal/core/song"

	"github.com/google/uuid"
)

type SongRepository struct {
	songs map[string]song.Song
}

func NewSongRepository() *SongRepository {
	return &SongRepository{
		songs: make(map[string]song.Song),
	}
}

func (r *SongRepository) GetAll() ([]song.Song, error) {
	result := make([]song.Song, 0, len(r.songs))
	for _, s := range r.songs {
		result = append(result, s)
	}
	return result, nil
}

func (r *SongRepository) GetByID(id string) (song.Song, error) {
	s, ok := r.songs[id]
	if !ok {
		return song.Song{}, errors.New("song not found")
	}
	return s, nil
}

func (r *SongRepository) Create(s song.Song) (song.Song, error) {
	s.ID = uuid.NewString()
	r.songs[s.ID] = s
	return s, nil
}

func (r *SongRepository) Update(id string, s song.Song) (song.Song, error) {
	_, ok := r.songs[id]
	if !ok {
		return song.Song{}, errors.New("song not found")
	}
	s.ID = id
	r.songs[id] = s
	return s, nil
}

func (r *SongRepository) Delete(id string) error {
	_, ok := r.songs[id]
	if !ok {
		return errors.New("song not found")
	}
	delete(r.songs, id)
	return nil
}
