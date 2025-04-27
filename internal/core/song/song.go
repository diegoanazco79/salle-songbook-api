package song

type Song struct {
	ID                     string `json:"id"`
	Title                  string `json:"title"`
	Lyrics                 string `json:"lyrics"`
	LyricsWithGuitarChords string `json:"lyrics_with_guitar_chords"`
	Author                 string `json:"author"`
}
