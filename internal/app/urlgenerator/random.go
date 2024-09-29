package urlgenerator

import (
	"fmt"
	"math/rand"
	"regexp"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/helpers"
)

type RandURLGenerator struct {
	URLGenerator
	length int
}

func NewRandURLGenerator(len int) *RandURLGenerator {
	return &RandURLGenerator{
		length: len,
	}
}

func (g *RandURLGenerator) GenerateURL(longURL string) (string, error) {
	if !helpers.IsURL(longURL) {
		return longURL, fmt.Errorf("invalid long URL")
	}

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
	return string(bytes), nil
}

func (lr *RandURLGenerator) IsValidURL(shortURL string) bool {
	regexPattern := fmt.Sprintf(`^/?[a-zA-Z0-9]{%d}$`, lr.length)
	validShortURL := regexp.MustCompile(regexPattern)
	return validShortURL.MatchString(shortURL)
}