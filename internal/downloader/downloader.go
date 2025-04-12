package downloader

import (
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "strings"

    "github.com/u2takey/ffmpeg-go"
    "go-video-sniffer/internal/extractor"
)

// Downloads and saves a video from a URL
func DownloadVideo(videoURL string, baseDomain string) (string, error) {
    parsedURL, err := url.Parse(videoURL)
    if err != nil {
        return "", err
    }

    fileName := extractFilename(parsedURL)
    destDir := filepath.Join("downloads", sanitizeFolder(baseDomain))
    os.MkdirAll(destDir, os.ModePerm)

    destPath := filepath.Join(destDir, fileName)

    if strings.Contains(videoURL, ".m3u8") {
        return downloadStream(videoURL, destPath)
    } else {
        return downloadDirect(videoURL, destPath)
    }
}

// Simple .mp4 download
func downloadDirect(videoURL, destPath string) (string, error) {
    resp, err := http.Get(videoURL)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    out, err := os.Create(destPath)
    if err != nil {
        return "", err
    }
    defer out.Close()

    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return "", err
    }

    return destPath, nil
}

// .m3u8 HLS download via ffmpeg-go
func downloadStream(streamURL, destPath string) (string, error) {
    err := ffmpeg.Input(streamURL).
        Output(destPath, ffmpeg.KwArgs{"c": "copy"}).
        OverWriteOutput().Run()

    if err != nil {
        return "", fmt.Errorf("ffmpeg download error: %v", err)
    }

    return destPath, nil
}

func extractFilename(u *url.URL) string {
    segments := strings.Split(u.Path, "/")
    name := segments[len(segments)-1]
    if name == "" || strings.Contains(name, ".m3u8")
