package song

type Repository interface {
	GetAll() ([]Song, error)
	GetByID(id string) (Song, error)
	Create(song Song) (Song, error)
	Update(id string, song Song) (Song, error)
	Delete(id string) error
}

type Song struct {
	ID                     string `json:"id" bson:"id"`
	Title                  string `json:"title" bson:"title"`
	Lyrics                 string `json:"lyrics" bson:"lyrics"`
	LyricsWithGuitarChords string `json:"lyrics_with_guitar_chords" bson:"lyrics_with_guitar_chords"`
	Author                 string `json:"author" bson:"author"`
}
