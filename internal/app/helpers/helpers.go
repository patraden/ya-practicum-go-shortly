package helpers

import "net/url"

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
