package extractor

import (
    "encoding/json"
    "fmt"
    "github.com/u2takey/ffmpeg-go"
    "os"
)

type VideoMetadata struct {
    Format struct {
        Duration string `json:"duration"`
        Tags     struct {
            Title string `json:"title"`
        } `json:"tags"`
    } `json:"format"`

    Streams []struct {
        CodecName string `json:"codec_name"`
        Width     int    `json:"width"`
        Height    int    `json:"height"`
    } `json:"streams"`
}

func ExtractMetadata(filePath string) (*VideoMetadata, error) {
    buf, err := ffmpeg.Input(filePath).Output("pipe:", ffmpeg.KwArgs{"f": "ffprobe", "show_format": nil, "show_streams": nil, "print_format": "json"}).Compile().OutputToString()
    if err != nil {
        return nil, fmt.Errorf("ffmpeg metadata extract error: %v", err)
    }

    var meta VideoMetadata
    if err := json.Unmarshal([]byte(buf), &meta); err != nil {
        return nil, fmt.Errorf("json unmarshal error: %v", err)
    }

    return &meta, nil
}
