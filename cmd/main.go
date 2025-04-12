package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go-video-sniffer/internal/crawler"
	"go-video-sniffer/internal/detector"
	"go-video-sniffer/internal/downloader"
	"go-video-sniffer/internal/metadata"
)

func saveMetadata(metadata *metadata.VideoMetadata, videoPath string) error {
	metadataPath := videoPath + ".json"

	file, err := os.Create(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(metadata); err != nil {
		return fmt.Errorf("failed to encode metadata: %w", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: video-sniffer <url>")
		os.Exit(1)
	}

	targetURL := os.Args[1]
	outputDir := "output"

	fmt.Println("🔍 Crawling:", targetURL)

	rawLinks, err := crawler.CrawlForMediaLinks(targetURL)
	if err != nil {
		log.Fatalf("❌ Failed to crawl: %v", err)
	}

	if len(rawLinks) == 0 {
		fmt.Println("No media links found.")
		return
	}

	fmt.Printf("🌐 Found %d raw media links\n", len(rawLinks))

	videoLinks, err := detector.NormalizeAndFilter(targetURL, rawLinks)
	if err != nil {
		log.Fatalf("❌ Failed to filter links: %v", err)
	}

	if len(videoLinks) == 0 {
		fmt.Println("No valid video links found.")
		return
	}

	fmt.Printf("🎥 Found %d valid video links\n", len(videoLinks))

	for _, link := range videoLinks {
		fmt.Println("⬇️  Downloading:", link)
		filePath, err := downloader.DownloadVideo(link, outputDir)
		if err != nil {
			log.Printf("⚠️  Failed to download %s: %v", link, err)
			continue
		}

		fmt.Println("📁 Saved to:", filePath)

		fmt.Println("🔍 Extracting metadata...")
		meta, err := metadata.ExtractMetadata(filePath)
		if err != nil {
			log.Printf("⚠️  Failed to extract metadata: %v", err)
			continue
		}

		// Save metadata to a JSON file
		if err := saveMetadata(meta, filePath); err != nil {
			log.Printf("⚠️  Failed to save metadata for %s: %v", filePath, err)
		}

		fmt.Println("🧠 Metadata saved to:", filePath+".json")
		fmt.Println("--------------------------------")
	}
}
