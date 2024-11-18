package utils

import (
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
)

const errLabel = "utils"

func IsServerAddress(addr string) bool {
	return len(strings.Split(addr, ":")) == len(strings.Split("server:port", ":"))
}

// backoff that retries in constant time.
func LinearBackoff(maxElapsedTime time.Duration, interval time.Duration) *backoff.ExponentialBackOff {
	return backoff.NewExponentialBackOff(
		backoff.WithMaxElapsedTime(maxElapsedTime),
		backoff.WithMultiplier(1),
		backoff.WithRandomizationFactor(0),
		backoff.WithInitialInterval(interval),
	)
}

// Generate a random string for paths or query parameters.
func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

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
