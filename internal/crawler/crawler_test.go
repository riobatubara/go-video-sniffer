package crawler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCrawlForMediaLinks(t *testing.T) {
	html := `
	<html>
		<body>
			<video src="https://example.com/video.mp4"></video>
			<source src="https://cdn.example.com/stream.m3u8" />
		</body>
	</html>`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	}))
	defer ts.Close()

	links, err := CrawlForMediaLinks(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{
		"https://example.com/video.mp4",
		"https://cdn.example.com/stream.m3u8",
	}

	if len(links) != len(expected) {
		t.Errorf("expected %d links, got %d", len(expected), len(links))
	}

	for i, link := range expected {
		if links[i] != link {
			t.Errorf("expected link %s, got %s", link, links[i])
		}
	}
}
