// Copyright © 2014-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package land

import (
	"image"
	"image/color"
)

// Tile holds a portion of the overall world map and the parameters
// that uniquely indicate the map portion.
type Tile interface {
	Topo() Topo                // Height data.
	Zoom() int                 // Zoom (level of detail) for this tile.
	Origin() (ox, oy int)      // Tile origin.
	Set(zoom, ox, oy int) Tile // Repurpose this tile. Data needs updating.
	SetFace(face int) Tile     // Needed to fill a 3D cube face image.
}

// ============================================================================

// tile is the default implementation of Tile.
type tile struct {
	topo   Topo // Height and region data.
	zoom   int  // Zoom level of this tile.
	ox, oy int  // Origin for this tile.
	face   int  // Cube face.
}

// newTile allocates and initializes a new map tile.
func newTile(topo Topo, zoom, ox, oy int) *tile {
	return &tile{topo: topo, zoom: zoom, ox: ox, oy: oy}
}

// Topo implements Tile.
func (t *tile) Topo() Topo { return t.topo }

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

// image creates a png image of a tile. Expected to be used for debugging.
func (t *tile) image(landSplit float64) *image.NRGBA {
	var c *color.NRGBA
	img := image.NewNRGBA(image.Rect(0, 0, len(t.topo), len(t.topo[0])))
	for x := range t.topo {
		for y := range t.topo[x] {
			c = t.topo.paint(t.topo[x][y], landSplit)
			img.SetNRGBA(x, y, *c)
		}
	}
	return img
}
