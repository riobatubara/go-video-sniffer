package detector

import (
	"reflect"
	"testing"
)

func TestIsVideoURL(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"https://example.com/video.mp4", true},
		{"https://example.com/stream.m3u8", true},
		{"https://example.com/image.jpg", false},
		{"https://example.com/video", false},
	}

	for _, test := range tests {
		if got := IsVideoURL(test.url); got != test.expected {
			t.Errorf("IsVideoURL(%q) = %v; want %v", test.url, got, test.expected)
		}
	}
}

func TestNormalizeAndFilter(t *testing.T) {
	base := "https://example.com/videos"
	rawLinks := []string{
		"video1.mp4",
		"/media/video2.m3u8",
		"https://cdn.example.com/movie.mov",
		"notvideo.txt",
	}

	expected := []string{
		"https://example.com/videos/video1.mp4",
		"https://example.com/media/video2.m3u8",
		"https://cdn.example.com/movie.mov",
	}

	filtered, err := NormalizeAndFilter(base, rawLinks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("got %v; want %v", filtered, expected)
	}
}
