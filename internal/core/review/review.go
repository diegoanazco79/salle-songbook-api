package review

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
	Title                  string `json:"title"`
	Lyrics                 string `json:"lyrics"`
	LyricsWithGuitarChords string `json:"lyrics_with_guitar_chords"`
	Author                 string `json:"author"`
}

type PendingReview struct {
	ID          string       `json:"id"`
	Action      ReviewAction `json:"action"`
	SongID      string       `json:"song_id,omitempty"`
	NewSongData *SongDataDTO `json:"new_song_data,omitempty"`
	RequestedBy string       `json:"requested_by"`
}
