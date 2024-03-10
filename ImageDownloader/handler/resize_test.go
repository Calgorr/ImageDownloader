package handler

import (
	"bytes"
	"image"
	"testing"

	"github.com/Calgorr/ImageDownloader/ImageDownloader/pkg/model"
	"github.com/bmizerany/assert"
)

func TestResizeImage(t *testing.T) {
	// Initialize ImageResizer with mock values
	imageResizer := &ImageResizer{
		BufferSize: 10,
		Width:      100,
		Height:     100,
	}

	// Create a channel for input images
	inputChannel := make(chan model.Image)
	defer close(inputChannel)

	// Create a sample image for testing
	img, _, _ := image.Decode(bytes.NewReader([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}))
	testImage := model.Image{
		Url:   "test_url",
		Image: img,
	}

	// Start resizing images
	resizedChannel := imageResizer.ResizeImage(inputChannel)

	// Feed the input channel with the test image
	go func() {
		inputChannel <- testImage
	}()

	// Retrieve the resized image
	resizedImage := <-resizedChannel

	// Assert that the resized image is as expected
	assert.Equal(t, testImage.Url, resizedImage.Url)

}
