package urlgenerator

import (
	"fmt"
	"math/rand"
	"regexp"
)

type RandURLGenerator struct {
	length int
}

func NewRandURLGenerator(l int) *RandURLGenerator {
	return &RandURLGenerator{
		length: l,
	}
}

func (g *RandURLGenerator) GenerateURL(_ string) string {
	charSets := []string{
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ", // A-Z
		"abcdefghijklmnopqrstuvwxyz", // a-z
		"0123456789",                 // 0-9
	}

	bytes := make([]byte, g.length)
	for i := range bytes {
		charSet := charSets[rand.Intn(len(charSets))]
		bytes[i] = charSet[rand.Intn(len(charSet))]
	}

	return string(bytes)
}

func (g *RandURLGenerator) IsValidURL(shortURL string) bool {
	regexPattern := fmt.Sprintf(`^/?[a-zA-Z0-9]{%d}$`, g.length)
	validShortURL := regexp.MustCompile(regexPattern)

	return validShortURL.MatchString(shortURL)
}
