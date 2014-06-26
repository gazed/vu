// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package load

import (
	"image"
	"image/color"
	"testing"
)

func TestLoadPng(t *testing.T) {
	load := newLoader().setDir(img, "../eg/images")
	if img, _ := load.png("xxx"); img != nil {
		t.Error("Image should be nil for bad files")
	}
	if img, _ := load.png("image"); img == nil {
		t.Error("Could not load image file")
	}
}

// Prove that RGBA images have alpha of 1 (255) and that
// NRGBA images are not alpha pre-multiplied.  Images have been
// created from gimp, one with an alpha channel and one without.
func TestLoadPngAlpha(t *testing.T) {
	load := newLoader().setDir(img, "../eg/images")
	img, _ := load.png("redNoAlpha")
	c := img.(*image.RGBA).At(200, 100)
	img, _ = load.png("redAlpha")
	a := img.(*image.NRGBA).At(200, 100)

	// check the red and alpha values in both cases.
	cr, ar := c.(color.RGBA).R, a.(color.NRGBA).R
	ca, aa := c.(color.RGBA).A, a.(color.NRGBA).A
	if cr != 255 || ar != 255 {
		t.Errorf("Got unexpected reds %d %d", cr, ar)
	}
	if ca != 255 || aa != 46 {
		t.Errorf("Got unexpected alphas %d %d", ca, aa)
	}
}
