// Copyright © 2014-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package land

import (
	"image"
	"image/color"
	"log"
)

// Tile holds a portion of the overall world map. It is indexed and keyed
// for easy storage and access. Tile is the pre-rendered portion of the map
// at different levels of detail.
type Tile interface {
	Topo() Topo                // Height data.
	Zoom() int                 // Zoom (level of detail) for this tile.
	XY() (tx, ty int)          // Tile X, Y index within the world.
	Key() string               // Unique tile id using zoom and tile XY.
	Set(zoom, tx, ty int) Tile // Repurpose this tile. Data needs resetting.
}

// ============================================================================

// tile is the default implementation of Tile.
type tile struct {
	topo Topo   // Height and region data.
	zoom int    // Zoom level of this tile.
	x, y int    // Tile X and Y within the world at this tile zoom level.
	key  string // Unique key for this tile representing zoom, x, and y.
}

// newTile allocates and initializes a new map tile.
func newTile(topo Topo, zoom, x, y int) *tile {
	return &tile{topo, zoom, x, y, tileKey(uint(zoom), uint(x), uint(y))}
}

// Topo implements Tile.
func (t *tile) Topo() Topo { return t.topo }

// Zoom implements Tile.
func (t *tile) Zoom() int { return t.zoom }

// XY implements Tile.
func (t *tile) XY() (x, y int) { return t.x, t.y }

// Key implements Tile.
func (t *tile) Key() string { return t.key }

// Set implements Tile.
func (t *tile) Set(zoom, x, y int) Tile {
	t.zoom = zoom
	t.x, t.y = x, y
	t.key = tileKey(uint(zoom), uint(x), uint(y))
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

// tile
// ============================================================================
// Utility methods. See:
//   http://msdn.microsoft.com/en-us/library/bb259689.aspx
//   http://www.microimages.com/documentation/TechGuides/76BingStructure.pdf

// tileKey converts the given tile coordinates and zoom level to map key.
// The key is constructed by interleaving the bits from the x and y tile locations.
// The length of the returned key string equals the zoom level.
func tileKey(zoom, tx, ty uint) (key string) {
	mask := uint(1)
	buff := make([]byte, zoom)
	for z := zoom; z > 0; z-- {
		part := byte('0')
		mask = 1 << (z - 1)
		switch {
		case tx&mask != 0 && ty&mask != 0:
			part = '3'
		case tx&mask != 0:
			part = '1'
		case ty&mask != 0:
			part = '2'
		}
		buff[zoom-z] = part
	}
	return string(buff)
}

// keyTile converts the given map key to tile and zoom coordinates.
func keyTile(key string) (zoom, tx, ty uint) {
	tx, ty = 0, 0
	zoom = uint(len(key))
	for cnt := zoom; cnt > 0; cnt-- {
		mask := uint(1) << (cnt - 1)
		switch key[zoom-cnt] {
		case '0':
		case '1':
			tx |= mask
		case '2':
			ty |= mask
		case '3':
			tx |= mask
			ty |= mask
		default:
			log.Printf("Invalid map key %s", key)
		}
	}
	return
}
