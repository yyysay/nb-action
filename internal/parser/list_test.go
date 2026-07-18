package parser

import (
	"os"
	"testing"
)

func TestReadImageList(t *testing.T) {
	content := `
# test images

docker.io/library/alpine

ghcr.io/sagernet/sing-box
`

	file := "test-images.txt"

	err := os.WriteFile(file, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(file)

	images, err := ReadImageList(file)

	if err != nil {
		t.Fatal(err)
	}

	if len(images) != 2 {
		t.Fatalf("got %d images, want 2", len(images))
	}
}
