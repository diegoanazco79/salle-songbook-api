package song

type Repository interface {
	GetAll() ([]Song, error)
	GetByID(id string) (Song, error)
	Create(song Song) (Song, error)
	Update(id string, song Song) (Song, error)
	Delete(id string) error
}

type ChordVersion struct {
	ID          string `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Content     string `json:"content" bson:"content"`
}

type Song struct {
	ID                     string         `json:"id" bson:"id"`
	Title                  string         `json:"title" bson:"title"`
	Lyrics                 string         `json:"lyrics" bson:"lyrics"`
	LyricsWithGuitarChords string         `json:"lyrics_with_guitar_chords" bson:"lyrics_with_guitar_chords"`
	GuitarChordsVersions   []ChordVersion `json:"guitar_chords_versions" bson:"guitar_chords_versions"`
	Author                 string         `json:"author" bson:"author"`
}
