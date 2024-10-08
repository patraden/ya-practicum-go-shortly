package urlgen

import (
	"math/rand"
)

type URLGenerator interface {
	GenerateShortURL(longURL string) (shortURL string)
}

// Random generator
type RandURLGenerator struct {
	length int
	URLGenerator
}

func NewRandURLGenerator(len int) *RandURLGenerator {
	return &RandURLGenerator{
		length: len,
	}
}

func (g *RandURLGenerator) GenerateShortURL(longURL string) (shortURL string) {
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
	shortURL = string(bytes)
	return
}
