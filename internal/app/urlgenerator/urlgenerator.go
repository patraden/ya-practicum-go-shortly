package urlgenerator

// Short URL generator interface
// Decided to dedicate an interface for this service
// as potentially thoughout the course of development
// there might be different implemnetations
// like random, incremental, hash based etc
type URLGenerator interface {
	GenerateURL(longURL string) string
	IsValidURL(shortURL string) bool
}
