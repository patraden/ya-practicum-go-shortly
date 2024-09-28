package urlgen

import (
	"math/rand"
)

type UrlGenerator interface {
	GenerateShortUrl(longUrl string) (shortUrl string)
}

// Random generator
type RandUrlGenerator struct {
	length int
	UrlGenerator
}

func NewRandUrlGenerator(len int) *RandUrlGenerator {
	return &RandUrlGenerator{
		length: len,
	}
}

func (g *RandUrlGenerator) GenerateShortUrl(longUrl string) (shortUrl string) {
	bytes := make([]byte, g.length)
	for i := 0; i < g.length; i++ {
		switch rand.Intn(3) {
		case 0:
			bytes[i] = byte(rand.Intn(26) + 65) // A-Z
		case 1:
			bytes[i] = byte(rand.Intn(26) + 97) // a-z
		case 2:
			bytes[i] = byte(rand.Intn(10) + 48) // 0-9
		}
	}
	shortUrl = string(bytes)
	return
}
