package urlgenerator

type URLGenerator interface {
	GenerateURL(longURL string) (string, error)
	IsValidURL(shortURL string) bool
}
