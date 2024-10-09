package service

type URLShortener interface {
	ShortenURL(longURL string) (string, error)
	GetOriginalURL(shortURL string) (string, error)
}
