// package extractor

// import (
//     "encoding/json"
//     "fmt"
//     "github.com/u2takey/ffmpeg-go"
//     "os"
// )

// type VideoMetadata struct {
//     Format struct {
//         Duration string `json:"duration"`
//         Tags     struct {
//             Title string `json:"title"`
//         } `json:"tags"`
//     } `json:"format"`

//     Streams []struct {
//         CodecName string `json:"codec_name"`
//         Width     int    `json:"width"`
//         Height    int    `json:"height"`
//     } `json:"streams"`
// }

// func ExtractMetadata(filePath string) (*VideoMetadata, error) {
//     buf, err := ffmpeg.Input(filePath).Output("pipe:", ffmpeg.KwArgs{"f": "ffprobe", "show_format": nil, "show_streams": nil, "print_format": "json"}).Compile().OutputToString()
//     if err != nil {
//         return nil, fmt.Errorf("ffmpeg metadata extract error: %v", err)
//     }

//     var meta VideoMetadata
//     if err := json.Unmarshal([]byte(buf), &meta); err != nil {
//         return nil, fmt.Errorf("json unmarshal error: %v", err)
//     }

//     return &meta, nil
// }

package metadata

import (
	"bytes"
	"encoding/json"
	"fmt"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// VideoMetadata holds common metadata fields
type VideoMetadata struct {
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`

	Streams []struct {
		CodecType string `json:"codec_type"`
		CodecName string `json:"codec_name"`
		Width     int    `json:"width,omitempty"`
		Height    int    `json:"height,omitempty"`
	} `json:"streams"`
}

// ExtractMetadata runs ffprobe and parses metadata
func ExtractMetadata(filePath string) (*VideoMetadata, error) {
	buf := bytes.NewBuffer(nil)

	err := ffmpeg.Input(filePath).
		Format("json").
		WithCustomCommand("ffprobe").
		WithOutput(buf, ffmpeg.KwArgs{"show_format": "1", "show_streams": "1"}).
		Run()

	if err != nil {
		return nil, fmt.Errorf("ffprobe failed: %w", err)
	}

	var metadata VideoMetadata
	if err := json.Unmarshal(buf.Bytes(), &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return &metadata, nil
}
