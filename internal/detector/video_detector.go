package detector

import (
	"net/url"
	"path"
	"strings"
)

// Supported video file extensions
var supportedVideoExtensions = []string{
	".mp4", ".webm", ".m3u8", ".mov", ".mpd", ".avi", ".flv", ".mkv",
}

// IsVideoURL returns true if the URL has a video extension
func IsVideoURL(link string) bool {
	for _, ext := range supportedVideoExtensions {
		if strings.Contains(strings.ToLower(link), ext) {
			return true
		}
	}
	return false
}

// NormalizeAndFilter filters out non-video URLs and resolves relative paths
func NormalizeAndFilter(base string, rawLinks []string) ([]string, error) {
	var validLinks []string

	baseURL, err := url.Parse(base)
	if err != nil {
		return nil, err
	}

	for _, raw := range rawLinks {
		parsedURL, err := url.Parse(raw)
		if err != nil {
			continue
		}
		finalURL := baseURL.ResolveReference(parsedURL).String()
		if IsVideoURL(finalURL) {
			validLinks = append(validLinks, finalURL)
		}
	}

	return validLinks, nil
}
