package validation 

import(
	"errors"
	"net/url"
	"strings"
)

var ErrInvalidYoutubeURL = errors.New("invalid YouTube URL")

func IsValidYoutubeURL(rawURL string) bool {
	parsedURL,err := url.ParseRequestURI(rawURL)

	if err != nil {
		return false
	}

	if parsedURL.Scheme != "https" {
		return false
	}

	if parsedURL.User != nil {
		return false
	}

	host := strings.ToLower(parsedURL.Hostname())

	AllowedHosts := map[string]bool{
        "youtube.com":             true,
        "www.youtube.com":         true,
        "m.youtube.com":           true,
        "youtu.be":                true,
        "www.youtube-nocookie.com": true,		
	}

	return AllowedHosts[host]
}