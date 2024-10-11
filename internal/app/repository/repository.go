package repository

type URLRepository interface {
	AddURL(shortURL string, longURL string) error
	GetURL(shortURL string) (string, error)
	DelURL(shortURL string) error
}
