package song

type Repository interface {
	GetAll() ([]Song, error)
	GetByID(id string) (Song, error)
	Create(song Song) (Song, error)
	Update(id string, song Song) (Song, error)
	Delete(id string) error
}

type Song struct {
	ID                     string `json:"id"`
	Title                  string `json:"title"`
	Lyrics                 string `json:"lyrics"`
	LyricsWithGuitarChords string `json:"lyrics_with_guitar_chords"`
	Author                 string `json:"author"`
}
