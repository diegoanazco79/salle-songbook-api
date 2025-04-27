package review

type ReviewAction string

const (
	Create ReviewAction = "create"
	Update ReviewAction = "update"
	Delete ReviewAction = "delete"
)

type PendingReview struct {
	ID          string       `json:"id"`
	Action      ReviewAction `json:"action"`
	SongID      string       `json:"song_id,omitempty"`
	NewSongData interface{}  `json:"new_song_data,omitempty"`
	RequestedBy string       `json:"requested_by"`
}
