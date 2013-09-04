// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package data

import (
	"image"
	"image/png"
	"io"
)

// png loads a PNG image into a format that is easily consumed by rendering layer.
// For example the image height and width are:
//    width := img.Bounds().Max.X - img.Bounds().Min.X
//    height := img.Bounds().Max.Y - img.Bounds().Min.Y
//
// And the image bytes in RGBA format is:
//    rgba, _ := img.(*image.RGBA)
//
// Any errors loading the image results in a nil return value, and the
// underlying error is returned.
func (l loader) png(directory, imageFile string) (img image.Image, err error) {
	var file io.ReadCloser
	if file, err = l.getResource(directory, imageFile); err == nil {
		defer file.Close()
		if img, err = png.Decode(file); err == nil {
			return
		}
	}
	return
}
