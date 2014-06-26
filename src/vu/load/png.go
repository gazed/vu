// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package load

import (
	"image"
	"image/png"
	"io"
)

// png loads a PNG image into an image that is easily consumed by rendering layer.
// The image height, width, and bytes in (N)RGBA format are:
//    width  := img.Bounds().Max.X - img.Bounds().Min.X
//    height := img.Bounds().Max.Y - img.Bounds().Min.Y
//    rgba, _ := img.(*image.(N)RGBA)
// Note that golang NRGBA are images with an alpha channel, but without alpha
// pre-multiplication. RGBA are images originally without an alpha channel,
// but assigned an alpha of 1 when read in.
// Any errors loading the image results in a nil return value, and the
// underlying error is returned.
func (l *loader) png(name string) (i image.Image, err error) {
	var file io.ReadCloser
	if file, err = l.getResource(l.dir[img], name+".png"); err == nil {
		defer file.Close()
		if i, err = png.Decode(file); err == nil {
			return
		}
	}
	return
}
