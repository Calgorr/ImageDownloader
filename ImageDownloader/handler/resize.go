package handler

import (
	imageprocessing "github.com/Calgorr/ImageDownloader/ImageDownloader/pkg/image_processing"
	"github.com/Calgorr/ImageDownloader/ImageDownloader/pkg/model"
)

type Resizer interface {
	ResizeImage(image <-chan model.Image) <-chan model.Image
}

type ImageResizer struct {
	Width      uint
	Height     uint
	BufferSize int
}

func NewImageResizer(width uint, height uint) Resizer {
	return &ImageResizer{
		Width:      width,
		Height:     height,
		BufferSize: 1000,
	}
}

func (i *ImageResizer) ResizeImage(image <-chan model.Image) <-chan model.Image {
	out := make(chan model.Image, i.BufferSize)
	go func() {
		for img := range image {
			img.Image = imageprocessing.Resize(img.Image, i.Width, i.Height)
			out <- img
		}
		close(out)
	}()
	return out
}
