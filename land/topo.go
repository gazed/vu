// Copyright © 2014-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package land

// The topo height information is created using a combination of noise
// (see noise.go) with fractional brownian motion as explained at:
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
// zoom level has an implied x,y index based on the topo size and zoom level.
type Topo [][]float64

// NewTopo allocates space for a single topology section.
// This is independent of zoom level.
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
// that seemlessly stitch together.
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
func (t Topo) generate(zoom, xoff, yoff int, n *noise) {
	freq := 2.0         // overall size.
	gain := 0.55        // range of heights.
	octaves := 6 + zoom // feature sharpness
	lacunarity := 2.0   // feature scatter
	zexp := 1.0 / math.Exp2(float64(zoom))
	size := float64(len(t))
	flip := len(t) - 1
	for x := range t {
		for y := range t[x] {
			total := 0.0
			nfreq := freq / size
			amplitude := gain
			for o := 0; o < octaves; o++ {
				xval := float64(x+xoff) * nfreq
				yval := float64(y+yoff) * nfreq
				total += n.generate2D(xval*zexp, yval*zexp) * amplitude
				nfreq *= lacunarity
				amplitude *= gain
			}
			t[x][flip-y] = total // Put 0,0 at bottom left.
		}
	}
}

// generate3D is used to create 1 side of a 6 side cube map. The images
// are generated such that they align with the other cube side images.
// Face is one of XPos, XNeg, YPos, YNeg, ZPos, ZNeg. The given value
// of xo,yo are applied based on the plane.
// The images are generated with 0,0 in the bottom left corner.
func (t Topo) generate3D(zoom, face, xoff, yoff int, n *noise) {
	freq := 2.0         // overall size.
	gain := 0.55        // range of heights.
	octaves := 6 + zoom // feature sharpness
	lacunarity := 2.0   // feature scatter
	exp := int(math.Exp2(float64(zoom)))
	zexp := 1.0 / float64(exp<<1) // exp2(zoom)
	size, flip := len(t), len(t)-1
	inv := 1.0 / float64(size)

	// iterate over a plane by fixing the planes loop.
	xo, xn, yo, yn, zo, zn := 0, size-1, 0, size-1, 0, size-1
	switch face {
	case XNeg:
		xo, xn, yo, zo = 0, 0, xoff, yoff
	case XPos:
		xo, xn, yo, zo = int(exp)*size-1, 0, xoff, yoff
	case YNeg:
		yo, yn, xo, zo = 0, 0, xoff, yoff
	case YPos:
		yo, yn, xo, zo = int(exp)*size-1, 0, xoff, yoff
	case ZNeg:
		zo, zn, xo, yo = 0, 0, xoff, yoff
	case ZPos:
		zo, zn, xo, yo = int(exp)*size-1, 0, xoff, yoff
	}
	for x := 0; x <= xn; x++ {
		for y := 0; y <= yn; y++ {
			for z := 0; z <= zn; z++ {
				total := 0.0
				nfreq := freq * inv
				amplitude := gain
				for o := 0; o < octaves; o++ {
					xval := float64(x+xo) * nfreq
					yval := float64(y+yo) * nfreq
					zval := float64(z+zo) * nfreq
					total += n.generate3D(xval*zexp, yval*zexp, zval*zexp) * amplitude
					nfreq *= lacunarity
					amplitude *= gain
				}

				// Put 0,0 at bottom left.
				switch face {
				case XNeg, XPos:
					t[y][flip-z] = total
				case YNeg, YPos:
					t[x][flip-z] = total
				case ZNeg, ZPos:
					t[x][flip-y] = total
				}
			}
		}
	}
}

// image creates a color 2D image from the height map where
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

// paint associates a color for the indicated section value.
// Note: this is for debugging only. Color probably should be a
// combination of land type and height.
func (t Topo) paint(height, landSplit float64) (c *color.NRGBA) {
	c = &color.NRGBA{255, 255, 255, 255}
	switch {
	case height > landSplit: // ground is uniform green for now.
		c = &color.NRGBA{0, 255, 0, 255}

	// shallower water is lighter.
	case height < landSplit:
		c = &color.NRGBA{10, 100, 200, 255}
	}
	return c
}
