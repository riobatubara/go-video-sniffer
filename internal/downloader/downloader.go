// package downloader

// import (
//     "fmt"
//     "io"
//     "net/http"
//     "net/url"
//     "os"
//     "path/filepath"
//     "strings"

//     "github.com/u2takey/ffmpeg-go"
//     "go-video-sniffer/internal/extractor"
// )

// // Downloads and saves a video from a URL
// func DownloadVideo(videoURL string, baseDomain string) (string, error) {
//     parsedURL, err := url.Parse(videoURL)
//     if err != nil {
//         return "", err
//     }

//     fileName := extractFilename(parsedURL)
//     destDir := filepath.Join("downloads", sanitizeFolder(baseDomain))
//     os.MkdirAll(destDir, os.ModePerm)

//     destPath := filepath.Join(destDir, fileName)

//     if strings.Contains(videoURL, ".m3u8") {
//         return downloadStream(videoURL, destPath)
//     } else {
//         return downloadDirect(videoURL, destPath)
//     }
// }

// // Simple .mp4 download
// func downloadDirect(videoURL, destPath string) (string, error) {
//     resp, err := http.Get(videoURL)
//     if err != nil {
//         return "", err
//     }
//     defer resp.Body.Close()

//     out, err := os.Create(destPath)
//     if err != nil {
//         return "", err
//     }
//     defer out.Close()

//     _, err = io.Copy(out, resp.Body)
//     if err != nil {
//         return "", err
//     }

//     return destPath, nil
// }

// // .m3u8 HLS download via ffmpeg-go
// func downloadStream(streamURL, destPath string) (string, error) {
//     err := ffmpeg.Input(streamURL).
//         Output(destPath, ffmpeg.KwArgs{"c": "copy"}).
//         OverWriteOutput().Run()

//     if err != nil {
//         return "", fmt.Errorf("ffmpeg download error: %v", err)
//     }

//     return destPath, nil
// }

// func extractFilename(u *url.URL) string {
//     segments := strings.Split(u.Path, "/")
//     name := segments[len(segments)-1]
//     if name == "" || strings.Contains(name, ".m3u8")

package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// sanitizeFilename removes unsafe characters from a filename
func sanitizeFilename(url string) string {
	name := filepath.Base(url)
	name = strings.Split(name, "?")[0]
	return strings.ReplaceAll(name, " ", "_")
}

// DownloadVideo downloads the video from the given URL and saves it to outputDir
func DownloadVideo(videoURL string, outputDir string) (string, error) {
	resp, err := http.Get(videoURL)
	if err != nil {
		return "", fmt.Errorf("failed to download %s: %w", videoURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status code %d for %s", resp.StatusCode, videoURL)
	}

	filename := sanitizeFilename(videoURL)
	outputPath := filepath.Join(outputDir, filename)

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create output dir: %w", err)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return outputPath, nil
}
