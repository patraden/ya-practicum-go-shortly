package utils

import (
	"net/url"
	"strings"
)

func IsURL(longURL string) bool {
	parsedURL, err := url.ParseRequestURI(longURL)
	if err != nil {
		return false
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}

	return true
}

func IsServerAddress(addr string) bool {
	return len(strings.Split(addr, ":")) == len(strings.Split("server:port", ":"))
}
