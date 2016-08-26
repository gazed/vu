// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package load

import (
	"image/png"
	"io"
)

// Png populates image data using the given reader.
// The Reader r is expected to be opened and closed by the caller.
// A successful import replaces the image in ImgData with a new image.
func Png(r io.Reader, d *ImgData) (err error) {
	d.Img, err = png.Decode(r)
	return err
}
