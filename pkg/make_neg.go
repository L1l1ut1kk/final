package pkg

import (
	"image"
	"image/color"
)

func CreateNegativeImage(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	negative := image.NewRGBA(bounds)

	// set the negative colors
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			negative.Set(x, y, color.RGBA{uint8(255 - r>>8), uint8(255 - g>>8), uint8(255 - b>>8), uint8(a >> 8)})
		}
	}

	return negative
}
