package dto

//go:generate easyjson -all url.go
type URLRequest struct {
	LongURL string `json:"url"`
}

type URLResponse struct {
	ShortURL string `json:"result"`
}
