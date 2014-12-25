// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package land

// The topo height information is created using a combination of perlin noise
// (see noise) with fractional brownian motion as explained at:
//   https://code.google.com/p/fractalterraingeneration/wiki/Fractional_Brownian_Motion

import (
	"image"
	"image/color"
	"math"
)

// Topo (topology) is height map information representing one part of a world.
// A small world may consist of only one zoom level (4 topos). A topo provides
// height information at a overall resolution based on a zoom value. Larger zoom
// values mean larger maps and hence more height maps. Each topo section at a given
// zoom level has an implied x, y index based on the topo size and zoom level.
type Topo [][]float64

// NewTopo allocates space for a single topology section. This is independent of
// zoom level.
func NewTopo(xwidth, yheight uint) Topo {
	t := make(Topo, xwidth)
	for x := range t {
		t[x] = make([]float64, yheight)
	}
	return t
}

// Size is the width and height of the topology section.
func (t Topo) Size() (x, y int) { return len(t), len(t[0]) }

// generate a topology section by creating height values between 1 and -1 for each
// element. Zoom, xoff and yoff can be used to create adjacent topolgy sections
// that seemlessly stitch together. The final topology section is the result of
// multiple generated simplex noise (2nd gen Perlin) combined as fractional
// brownian motion.
//
// Zoom, xoff, yoff are combined to indicate the location of the topology section
// within an overall world. Zoom increases the number of land sections needed
// for the world.
//   Zoom level 0 :     1 topology sections
//   Zoom level 1 :     4 topology sections indexed 0,0 to 1,1
//   Zoom level 2 :    16 topology sections indexed 0,0 to 3,3
//   Zoom level 3 :    64 topology sections indexed 0,0 to 7,7
//   Zoom level 4 :   256 topology sections indexed 0,0 to 15,15
//   Zoom level 5 :  1024 topology sections indexed 0,0 to 31,31
//   Zoom level 6 :  4096 topology sections indexed 0,0 to 63,63
//   Zoom level 7 : 16384 topology sections indexed 0,0 to 127,127
//   Zoom level 8 : 65536 topology sections indexed 0,0 to 255,255
func (t Topo) generate(zoom, xoff, yoff uint, n *noise) {
	freq := 2.0           // overall size.
	gain := 0.55          // range of heights.
	octaves := 6 + zoom*2 // feature sharpness
	lacunarity := 2.0     // feature scatter
	zexp := math.Exp2(float64(zoom))
	size := float64(len(t))
	for y := range t {
		for x := range t[y] {
			total := 0.0
			xfreq := freq / size
			yfreq := freq / size
			amplitude := gain
			for o := uint(0); o < octaves; o++ {
				xval := float64(x)*xfreq + float64(xoff)*size*xfreq
				yval := float64(y)*yfreq + float64(yoff)*size*yfreq
				total += n.generate(xval/zexp, yval/zexp) * amplitude
				xfreq *= lacunarity
				yfreq *= lacunarity
				amplitude *= gain
			}
			t[x][y] = total
		}
	}
}

// paint associates a colour for the indicated section value.
// Note: this is for debugging only. Colour should be based on
// land type, not height.
func (t Topo) paint(height, landSplit float64) (c *color.NRGBA) {
	c = &color.NRGBA{255, 255, 255, 255}
	switch {
	case height > landSplit: // ground is uniform green for now.
		c = &color.NRGBA{0, 255, 0, 255}

	// shallower water is lighter.
	case height < landSplit:
		h := 255 + (height * 255)
		f := math.Exp(height)
		c = &color.NRGBA{uint8(h * f * 0.5), uint8(h * f), 255, 255}
	}
	return c
}

// image creates a colour 2D image from the height map where
// heights below landSplit represent water and heights above landSplit
// represent land.
func (t Topo) image(landSplit float64) *image.NRGBA {
	var c *color.NRGBA
	img := image.NewNRGBA(image.Rect(0, 0, len(t), len(t[0])))
	for x := range t {
		for y := range t[x] {
			c = t.paint(t[x][y], landSplit)
			img.SetNRGBA(x, y, *c)
		}
	}
	return img
}
