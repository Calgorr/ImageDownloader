package handler

import (
	"image"
	"net/http"
	"scaleOps/ImageDownloader/pkg/model"
	"sync"
)

type UrlDownloader interface {
	DownloadImage(urls []string) <-chan model.Image
}

type ImageDownloader struct {
	Client              *http.Client
	ConcurrentDownloads int // Number of concurrent downloads
	BufferSize          int
}

func NewImageDownloader() UrlDownloader {
	return &ImageDownloader{
		Client:              http.DefaultClient,
		ConcurrentDownloads: 25,
		BufferSize:          1000,
	}
}

// DownloadImage downloads images from the provided URLs concurrently.
// The function returns a channel of images and sends the downloaded images to the channel.
// The function will close the channel when all the images have been downloaded.
func (i *ImageDownloader) DownloadImage(urls []string) <-chan model.Image {
	out := make(chan model.Image, i.BufferSize)
	var wg sync.WaitGroup

	for j := 0; j < i.ConcurrentDownloads; j++ {
		wg.Add(1)
		go func(urls []string) {
			defer wg.Done()
			for _, url := range urls {
				resp, err := i.Client.Get(url)
				if err != nil {
					continue
				}
				if resp.StatusCode != http.StatusOK {
					err := resp.Body.Close()
					if err != nil {
						return
					}
					continue
				}

				img, _, err := image.Decode(resp.Body)
				if err != nil {
					continue
				}

				image1 := model.Image{
					Url:   url,
					Image: img,
				}
				out <- image1
			}
		}(urls[j*len(urls)/i.ConcurrentDownloads : (j+1)*len(urls)/i.ConcurrentDownloads])
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
