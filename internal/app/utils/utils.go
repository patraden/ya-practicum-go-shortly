package utils

import (
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
)

const errLabel = "utils"

// IsServerAddress checks if the provided address is a valid server address with the format "server:port".
func IsServerAddress(addr string) bool {
	return len(strings.Split(addr, ":")) == len(strings.Split("server:port", ":"))
}

// LinearBackoff creates a backoff policy that retries in constant time.
func LinearBackoff(maxElapsedTime time.Duration, interval time.Duration) *backoff.ExponentialBackOff {
	return backoff.NewExponentialBackOff(
		backoff.WithMaxElapsedTime(maxElapsedTime),
		backoff.WithMultiplier(1),
		backoff.WithRandomizationFactor(0),
		backoff.WithInitialInterval(interval),
	)
}

// RandomString generates a random string of length n consisting of uppercase and lowercase letters.
func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

// RandURL generates a random URL with a random domain, path, and query parameters.
func RandURL() string {
	domainLen := 10
	pathLen := 10
	tldLen := 3
	schemaOpt := 2
	paramLen := 5

	schema := "http"
	if rand.Intn(schemaOpt) == 0 {
		schema = "https"
	}

	u := &url.URL{
		Scheme: schema,
		Host:   RandomString(domainLen) + "." + RandomString(tldLen),
		Path:   "resource/" + RandomString(pathLen),
	}

	q := u.Query()
	q.Set("id", RandomString(paramLen))
	u.RawQuery = q.Encode()

	return u.String()
}
