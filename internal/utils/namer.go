package utils

import (
    "fmt"
    "os"
    "path/filepath"
    "regexp"
    "strings"

    "go-video-sniffer/internal/extractor"
)

func GenerateFinalName(meta *extractor.VideoMetadata, originalPath string, domain string) (string, error) {
    title := meta.Format.Tags.Title
    if title == "" {
        title = "Untitled_Video"
    }

    // Clean title: remove bad chars
    title = sanitize(title)

    resolution := "unknown"
    if len(meta.Streams) > 0 && meta.Streams[0].Width > 0 {
        res := fmt.Sprintf("%dx%d", meta.Streams[0].Width, meta.Streams[0].Height)
        resolution = res
    }

    finalFileName := fmt.Sprintf("%s_%s.mp4", title, resolution)
    finalDir := filepath
