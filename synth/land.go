// Copyright © 2014-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package synth

import (
	"math"
)

// Land provides the ability to procedurally generate terrain height information.
// The higher the level of detail the larger the terrain. Land is created
// using NewLand().
type Land interface {
	TileSize() int     // Land tile width, height. Standard is 256
	Size(zoom int) int // Width and height at the given zoom.

	// Allocate and populate the indicated land tile with height data.
	//    tx,ty: topo/tile index at given zoom.
	NewTile(zoom, tx, ty int) Tile
	Fill(tile Tile)   // (Re)populates a tile with 2D height data.
	Fill3D(tile Tile) // (Re)populates a tile with 3D height data.
}

// NewLand initializes the procedural land generator. The seed determines
// land shape, such that lands created from the same seed will be the same.
// The zoom determines the overall size of the land (limit to 8
// or less pending stress testing). For example if tileSize is 256 then
// increasing the level of detail results in the following sizes:
//    zoom  0 :  256*2^0  = 256m
//    zoom  1 :  256*2^1  = 512
//    zoom  2 :  256*2^2  = 1024     ~1km2
//    zoom  3 :  256*2^3  = 2048     ~4km2
//    zoom  4 :  256*2^4  = 4096
//    zoom  5 :  256*2^5  = 8192     ~64km2 Medium size city.
//    zoom  6 :  256*2^6  = 16384
//    zoom  7 :  256*2^7  = 37768
//    zoom  8 :  256*2^8  = 65536
//    zoom  9 :  256*2^9  = 131072
//    zoom 10 :  256*2^10 = 262144
//    zoom 11 :  256*2^11 = 524288
//    zoom 12 :  256*2^12 = 1048576  ~1,000,000km2 Size of Ontario
//    zoom 13 :  256*2^13 = 2097152
//                          3162000  ~10,000,000km2 Size of Canada
//    zoom 14 :  256*2^14 = 4194304
//    zoom 15 :  256*2^15 = 8388608
//    zoom 16 :  256*2^16 = 16777216
//                          22583000 ~510,000,000km2 Size of Earth
//    zoom 17 :  256*2^17 = 33554472
// Land heights are generated by creating tiles at a particular zoom level.
// It is up to the caller to store/cache or regenerate tiles as needed.
func NewLand(tileSize int, seed int64) Land {
	return newLand(tileSize, seed)
}

// Land interface
// ============================================================================
// land is a default implementation of Land

// land creates the world height map.
// There is a limit to the number of map sections held in memory at any one
// time. Map tiles are expected to be generated as needed for large worlds.
//
// Nice article on how map level of detail can be organized at:
//   http://msdn.microsoft.com/en-us/library/bb259689.aspx
//   http://www.microimages.com/documentation/TechGuides/76BingStructure.pdf
type land struct {
	n    *simplex // expected to be  simplex noise maker.
	seed int64    // for all random calcuations.
	size int      // land tile width and height.
}

// newLand initializes the data needed to create a world land map. The higher
// the zoom, the larger the map and more map tiles needed at each zoom level.
func newLand(size int, seed int64) *land {
	l := &land{}
	l.seed = seed
	l.size = size
	l.n = newSimplex(l.seed)
	return l
}

// Size implements Land.
func (l *land) Size(atZoom int) int { return int(math.Exp2(float64(atZoom))) * l.size }

// TileSize implements Land.
func (l *land) TileSize() int { return l.size }

// NewTile implements Land.
func (l *land) NewTile(atZoom, x, y int) Tile { return l.newTile(atZoom, x, y) }
func (l *land) newTile(atZoom, x, y int) *tile {
	tile := newTile(uint(l.size), uint(l.size), atZoom, x, y)
	l.Fill(tile)
	return tile
}

// Fill implements Land.
// It populates the given tile with height data calculated using the world
// random seed. Note that tile ox,oy are relative to the over all zoom level.
// At zoom 0 they are expected to be 0,0. At zoom 1 they can range from
// 0 to mapSize-1.
func (l *land) Fill(landTile Tile) {
	t, _ := landTile.(*tile)
	if len(t.topo) == l.size && len(t.topo[0]) == l.size {
		t.gen2D(l.n)
	}
}

// Fill3D implements Land.
func (l *land) Fill3D(landTile Tile) {
	t, _ := landTile.(*tile)
	if len(t.topo) == l.size && len(t.topo[0]) == l.size {
		t.gen3D(l.n)
	}
}

// Cube face identifiers used when generating the six sides
// of a 3D cube map. Used in Tile.SetFace.
const (
	XPos = iota // Face on right side of cube.
	XNeg        // Face on left side of cube.
	YPos        // Face on top of cube.
	YNeg        // Face on bottom of cube.
	ZPos        // Face on at back of cube.
	ZNeg        // Face on front of cube.
)
