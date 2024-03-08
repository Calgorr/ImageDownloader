package imageprocessing

import (
	"image"

	"github.com/nfnt/resize"
)

func Resize(img image.Image, newWidth, newHeight uint) image.Image {
	resizedImg := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
	return resizedImg
}
