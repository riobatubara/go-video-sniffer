package downloader

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadVideo(t *testing.T) {
	videoContent := "fake video bytes"
	videoPath := "/video/test.mp4"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == videoPath {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(videoContent))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	outputDir := "test_output"
	defer os.RemoveAll(outputDir) // Clean up after test

	downloadURL := ts.URL + videoPath
	outputFilePath, err := DownloadVideo(downloadURL, outputDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outputFilePath)
	if err != nil {
		t.Fatalf("failed to read downloaded file: %v", err)
	}

	if string(data) != videoContent {
		t.Errorf("file content mismatch: got %s, want %s", string(data), videoContent)
	}
}
