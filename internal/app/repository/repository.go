package repository

type LinkRepository interface {
	Store(longURL string) (string, error)
	ReStore(shortURL string) (string, error)
}
