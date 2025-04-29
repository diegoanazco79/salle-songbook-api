package review

import "salle-songbook-api/internal/core/song"

type ReviewAction string

type Repository interface {
	GetAll() ([]PendingReview, error)
	GetByID(id string) (PendingReview, error)
	Create(PendingReview) (PendingReview, error)
	Delete(id string) error
}

const (
	Create ReviewAction = "create"
	Update ReviewAction = "update"
	Delete ReviewAction = "delete"
)

type SongDataDTO struct {
	Title                  string              `json:"title" bson:"title"`
	Lyrics                 string              `json:"lyrics" bson:"lyrics"`
	LyricsWithGuitarChords string              `json:"lyrics_with_guitar_chords" bson:"lyrics_with_guitar_chords"`
	GuitarChordsVersions   []song.ChordVersion `json:"guitar_chords_versions" bson:"guitar_chords_versions"`
	Author                 string              `json:"author" bson:"author"`
}

type PendingReview struct {
	ID          string       `json:"id" bson:"id"`
	Action      ReviewAction `json:"action" bson:"action"`
	SongID      string       `json:"song_id,omitempty" bson:"song_id,omitempty"`
	NewSongData *SongDataDTO `json:"new_song_data,omitempty" bson:"new_song_data,omitempty"`
	RequestedBy string       `json:"requested_by" bson:"requested_by"`
}
