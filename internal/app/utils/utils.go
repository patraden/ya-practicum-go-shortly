package utils

import (
	"crypto/rand"
	"math/big"
	"net/url"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
)

const (
	errLabel = "utils"
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

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
func RandString(n int) string {
	bytes := make([]byte, n)

	if _, err := rand.Read(bytes); err != nil {
		return ``
	}

	for i := range bytes {
		bytes[i] = alphabet[bytes[i]%byte(len(alphabet))]
	}

	return string(bytes)
}

// RandInt generates a random in interval [0, abs(max)).
func RandInt(max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0
	}

	return int(n.Int64())
}

// RandURL generates a random URL with a random domain, path, and query parameters.
func RandURL() string {
	domainLen := 10
	pathLen := 10
	tldLen := 3
	schemaOpt := 2
	paramLen := 5
	schema := "http"

	coin := RandInt(schemaOpt)
	if coin == 0 {
		schema = "https"
	}

	u := &url.URL{
		Scheme: schema,
		Host:   RandString(domainLen) + "." + RandString(tldLen),
		Path:   "resource/" + RandString(pathLen),
	}

	q := u.Query()
	q.Set("id", RandString(paramLen))
	u.RawQuery = q.Encode()

	return u.String()
}
