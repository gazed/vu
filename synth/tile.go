// Copyright © 2014-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package synth

import (
	"math"
)

// Tile holds a portion of the overall world map and the parameters
// that uniquely indicate the map portion.
type Tile interface {
	Size() (x, y int)          // Tile width and height dimensions.
	Zoom() int                 // Zoom (level of detail) for this tile.
	Topo() [][]float64         // Size x Size height data points.
	Origin() (ox, oy int)      // Tile origin.
	Set(zoom, ox, oy int) Tile // Repurpose this tile. Data needs updating.
	SetFace(face int) Tile     // Needed to fill a 3D cube face image.
}

// ============================================================================

// tile is the default implementation of Tile.
type tile struct {
	topo   [][]float64 // Height values for venly distributed grid points.
	zoom   int         // Zoom level of this tile.
	ox, oy int         // Origin for this tile.
	face   int         // Cube face.
}

// newTile allocates and initializes a new map tile.
func newTile(width, height uint, zoom, ox, oy int) *tile {
	topo := make([][]float64, width)
	for x := range topo {
		topo[x] = make([]float64, height)
	}
	return &tile{topo: topo, zoom: zoom, ox: ox, oy: oy}
}

// Size is the width and height of the topology section.
func (t *tile) Size() (x, y int)  { return len(t.topo), len(t.topo[0]) }
func (t *tile) Topo() [][]float64 { return t.topo }

// Zoom implements Tile.
func (t *tile) Zoom() int { return t.zoom }

// Origin implements Tile.
func (t *tile) Origin() (x, y int) { return t.ox, t.oy }

// Set implements Tile.
func (t *tile) Set(zoom, ox, oy int) Tile {
	t.zoom = zoom
	t.ox, t.oy = ox, oy
	return t
}

// SetFace implements Tile.
func (t *tile) SetFace(face int) Tile {
	if face >= XPos && face <= ZNeg {
		t.face = face
	}
	return t
}

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
func (t *tile) gen2D(n *simplex) {
	freq := 2.0           // overall size.
	gain := 0.55          // range of heights.
	octaves := 6 + t.zoom // feature sharpness
	lacunarity := 2.0     // feature scatter
	zexp := 1.0 / math.Exp2(float64(t.zoom))
	size := float64(len(t.topo))
	flip := len(t.topo) - 1
	for x := range t.topo {
		for y := range t.topo[x] {
			total := 0.0
			nfreq := freq / size
			amplitude := gain
			for o := 0; o < octaves; o++ {
				xval := float64(x+t.ox) * nfreq
				yval := float64(y+t.ox) * nfreq
				total += n.Gen2D(xval*zexp, yval*zexp) * amplitude
				nfreq *= lacunarity
				amplitude *= gain
			}
			t.topo[x][flip-y] = total // Put 0,0 at bottom left.
		}
	}
}

// gen3D is used to create 1 side of a 6 side cube map. The images
// are generated such that they align with the other cube side images.
// Face is one of XPos, XNeg, YPos, YNeg, ZPos, ZNeg. The given value
// of xo,yo are applied based on the plane.
// The images are generated with 0,0 in the bottom left corner.
func (t *tile) gen3D(n *simplex) {
	freq := 2.0           // overall size.
	gain := 0.55          // range of heights.
	octaves := 6 + t.zoom // feature sharpness
	lacunarity := 2.0     // feature scatter
	exp := int(math.Exp2(float64(t.zoom)))
	zexp := 1.0 / float64(exp<<1) // exp2(zoom)
	size, flip := len(t.topo), len(t.topo)-1
	inv := 1.0 / float64(size)

	// iterate over a plane by fixing the planes loop.
	xo, xn, yo, yn, zo, zn := 0, size-1, 0, size-1, 0, size-1
	switch t.face {
	case XNeg:
		xo, xn, yo, zo = 0, 0, t.ox, t.oy
	case XPos:
		xo, xn, yo, zo = int(exp)*size-1, 0, t.ox, t.oy
	case YNeg:
		yo, yn, xo, zo = 0, 0, t.ox, t.oy
	case YPos:
		yo, yn, xo, zo = int(exp)*size-1, 0, t.ox, t.oy
	case ZNeg:
		zo, zn, xo, yo = 0, 0, t.ox, t.oy
	case ZPos:
		zo, zn, xo, yo = int(exp)*size-1, 0, t.ox, t.oy
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
					total += n.Gen3D(xval*zexp, yval*zexp, zval*zexp) * amplitude
					nfreq *= lacunarity
					amplitude *= gain
				}

				// Put 0,0 at bottom left.
				switch t.face {
				case XNeg, XPos:
					t.topo[y][flip-z] = total
				case YNeg, YPos:
					t.topo[x][flip-z] = total
				case ZNeg, ZPos:
					t.topo[x][flip-y] = total
				}
			}
		}
	}
}
