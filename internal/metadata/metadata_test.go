package metadata

import (
	"testing"
)

func TestExtractMetadata(t *testing.T) {
	metadata, err := ExtractMetadata("testdata/test.mp4")
	if err != nil {
		t.Fatalf("failed to extract metadata: %v", err)
	}

	if metadata.Format.Duration == "" {
		t.Errorf("expected duration, got empty")
	}

	foundVideo := false
	for _, stream := range metadata.Streams {
		if stream.CodecType == "video" {
			foundVideo = true
			if stream.Width == 0 || stream.Height == 0 {
				t.Errorf("expected width/height, got 0")
			}
		}
	}
	if !foundVideo {
		t.Errorf("no video stream found")
	}
}
