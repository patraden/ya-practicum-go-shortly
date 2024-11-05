package dto

//go:generate easyjson -all url.go
type ShortenURLRequest struct {
	LongURL string `json:"url"`
}

type ShortenedURLResponse struct {
	ShortURL string `json:"result"`
}
