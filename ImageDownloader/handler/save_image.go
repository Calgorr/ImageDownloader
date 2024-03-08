package handler

import (
	"context"
	"time"

	"github.com/Calgorr/ImageDownloader/ImageDownloader/pkg/model"
	"github.com/Calgorr/ImageDownloader/ImageDownloader/pkg/repository"
)

type SaveImage interface {
	SaveImage(image <-chan model.Image) <-chan bool
}

type ImageSaver struct {
	ImageRepo  repository.ImageRepository
	BufferSize int
}

func NewImageSaver(imageRepo repository.ImageRepository) SaveImage {
	return &ImageSaver{
		ImageRepo:  imageRepo,
		BufferSize: 1000,
	}
}

func (i *ImageSaver) SaveImage(image <-chan model.Image) <-chan bool {
	out := make(chan bool, i.BufferSize)
	timer := time.NewTimer(1 * time.Second)
	go func() {
		defer close(out)
		var images []*model.Image
		for {
			select {
			case img, ok := <-image:
				if !ok {
					i.flushImages(images, out)
					return
				}
				images = append(images, &img)
				if len(images) >= repository.BatchSize {
					i.flushImages(images, out)
					images = []*model.Image{}
				}
			case <-timer.C:
				if len(images) > 0 {
					i.flushImages(images, out)
					images = []*model.Image{}
				}
			}
		}
	}()
	return out
}

func (i *ImageSaver) flushImages(images []*model.Image, out chan<- bool) {
	err := i.ImageRepo.InsertImage(context.Background(), images)
	if err != nil {
		out <- false
		return
	}
	out <- true
}
