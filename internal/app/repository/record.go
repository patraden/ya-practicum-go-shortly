package repository

type FileRecord struct {
	ID       int    `json:"uuid"`
	ShortURL string `json:"short_url"`
	LongURL  string `json:"original_url"`
}
