package process

import (
	"bytes"
	"image/jpeg"
	"os"
	"testing"

	"github.com/unknownmemory/img-processing/internal/shared"
)

func TestTransform(t *testing.T) {
	input, err := os.ReadFile("../../test/fixtures/images/test.jpg")
	if err != nil {
		t.Fatalf("Fixture error: %v", err)
	}

	t.Run("Resize", func(t *testing.T) {
		resize := 0.5

		buffer, _, err := Transform(input, shared.Transformations{
			Resize: &resize,
		})
		if err != nil {
			t.Fatalf("Error during transformation: %v", err)
		}

		original, err := jpeg.DecodeConfig(bytes.NewReader(input))
		if err != nil {
			t.Fatalf("Original image decode: %v", err)
		}

		resized, err := jpeg.DecodeConfig(bytes.NewReader(buffer))
		if err != nil {
			t.Fatalf("Transformed image decode: %v", err)
		}

		if resized.Width >= original.Width || resized.Height >= original.Height {
			t.Fatalf("Expected resized image to be smaller than the original")
		}
	})
}
