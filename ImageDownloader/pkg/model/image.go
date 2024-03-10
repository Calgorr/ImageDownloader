package model

import "image"

type Image struct {
	Url         string
	ContentType string
	Image       image.Image
}
