package main

import (
	"image"
	"image/color"

	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var (
	galleryImage = newGalleryImage()
	imagePreview = kit.ImagePreviewState{Fit: true}
)

// newGalleryImage 生成画廊使用的网格示例图。
func newGalleryImage() kit.ImageResource {
	picture := image.NewRGBA(image.Rect(0, 0, 640, 360))
	for y := 0; y < 360; y++ {
		for x := 0; x < 640; x++ {
			v := uint8(40 + (x/40+y/40)%2*16)
			picture.SetRGBA(x, y, color.RGBA{v, v + 16, v + 32, 255})
		}
	}
	return kit.ImageResource{Source: ui.NewBitmap(picture), Alt: "Grid preview"}
}

func init() {
	register(page{name: "Images", order: 34, view: func(c *ui.Context) {
		kit.ImagePreview(c, galleryImage, &imagePreview).Grow(1)
	}})
}
