package file

//go:generate easyjson -all record.go
type Record struct {
	UUID     string `json:"uuid"`
	ShortURL string `json:"short_url"`
	LongURL  string `json:"original_url"`
}
