package handler

import (
	"bytes"
	"image"
	"testing"

	"github.com/Calgorr/ImageDownloader/ImageDownloader/pkg/model"
	"github.com/bmizerany/assert"
)

func TestResizeImage(t *testing.T) {
	imageResizer := &ImageResizer{
		BufferSize: 10,
		Width:      100,
		Height:     100,
	}

	inputChannel := make(chan model.Image)
	defer close(inputChannel)

	img, _, _ := image.Decode(bytes.NewReader([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}))
	testImage := model.Image{
		Url:   "test_url",
		Image: img,
	}

	resizedChannel := imageResizer.ResizeImage(inputChannel)

	go func() {
		inputChannel <- testImage
	}()

	resizedImage := <-resizedChannel

	assert.Equal(t, testImage.Url, resizedImage.Url)

}
