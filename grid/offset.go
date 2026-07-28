// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.

package grid

import (
	"fmt"
)

// OffsetGrid overlaps and offsets a regular grid so that the
// regular grid corners are at the center of the offset-grid tile.
type OffsetGrid [][]uint8

// Offset grid elements are tiles consisting of 4 quadrants.
// Each quadrant is set/unset resulting in 16 distinct tile types.
const (
	QuadTL uint8 = 0b00001000 // top left
	QuadTR uint8 = 0b00000100 // top right
	QuadBL uint8 = 0b00000010 // bottom left
	QuadBR uint8 = 0b00000001 // bottom right
)

// Tile indexes for grid tile types. There are 16 possible tiles:
// 0-15 with the empty tile as default.
const (
	Tile00     uint8 = iota // 0,0,0,0  0 quadrant "empty"
	Tile1Q                  // 0,0,0,1  1 quadrant
	Tile1Q_270              // 0,0,1,0  1 quadrant rotated 270deg
	Tile2Q                  // 0,0,1,1  2 quadrant
	Tile1Q_90               // 0,1,0,0  1 quadrant rotated 90deg
	Tile2Q_90               // 0,1,0,1  2 quadrant rotated 90deg
	TileD2_90               // 0,1,1,0  2 quad diagonal rotated 90deg
	Tile3Q                  // 0,1,1,1  3 quadrant
	Tile1Q_180              // 1,0,0,0  1 quadrant rotated 180deg
	TileD2                  // 1,0,0,1  2 quad diagonal
	Tile2Q_270              // 1,0,1,0  2 quad diagonal rotated 270deg
	Tile3Q_270              // 1,0,1,1  3 quadrant rotated 270deg
	Tile2Q_180              // 1,1,0,0  2 quad diagonal rotated 180deg
	Tile3Q_90               // 1,1,0,1  3 quadrant rotated 90deg
	Tile3Q_180              // 1,1,1,0  3 quadrant rotated 180deg
	Tile4Q                  // 1,1,1,1  4 quadrant "solid"
)

// Symbols for debug grid visualizations using the following unicode
// block elements and spaces to create 2 character strings.
// This displays nicely with the Hack fixed width font.
// - U+2588 █ 	full block
// - U+2584 ▄	lower half block
// - U+2580  ▀  upper half block
//
// Each 2 character string represents one tile.
var OffsetSymbols = map[uint8]string{
	Tile00:     "  ", // space x 2
	Tile1Q:     " ▄", // space + lower half
	Tile1Q_270: "▄ ", // lower half + space
	Tile2Q:     "▄▄", // lower half x 2
	Tile1Q_90:  " ▀", // space + upper half
	Tile2Q_90:  " █", // space+full
	TileD2_90:  "▄▀", // lower half + upper half
	Tile3Q:     "▄█", // lower-half + full
	Tile1Q_180: "▀ ", // upper half + space
	TileD2:     "▀▄", // upper half + lower half
	Tile2Q_270: "█ ", // full+space
	Tile3Q_270: "█▄", // full + lower-half
	Tile2Q_180: "▀▀", // upper half x 2
	Tile3Q_90:  "▀█", // upper-half + full
	Tile3Q_180: "█▀", // full + upper-half
	Tile4Q:     "██", // full x 2
}

// MakeOffsetGrid takes a standard grid and returns the corresponding offset grid.
func MakeOffsetGrid(g Grid, fillEdges bool) (ogrid OffsetGrid) {
	if !g.Valid() {
		return ogrid // uninitialized offset grid.
	}

	// initialize offset-grid to overlap the given grid.
	cols, rows := len(g), len(g[0])
	ogrid = make([][]uint8, rows+1)
	for i := range ogrid {
		ogrid[i] = make([]uint8, cols+1)
	}

	// set the offset-grid tiles based on the given grid.
	for r, row := range g {
		for c := range row {
			set := g[r][c] == SOLID // false:EMPTY true:SOLID

			// set the offset-grid quads to match the grid.
			// Each grid overlaps one quad of 4 overlapping
			// offset-grid tiles
			ogrid[r][c] = setTile(ogrid[r][c], QuadBR, set)
			ogrid[r][c+1] = setTile(ogrid[r][c+1], QuadBL, set)
			ogrid[r+1][c] = setTile(ogrid[r+1][c], QuadTR, set)
			ogrid[r+1][c+1] = setTile(ogrid[r+1][c+1], QuadTL, set)
		}
	}

	// offset-grids overlap regular grids by 1 and thus have extra
	// half tile sized edges not contained in the regular grid.
	// These can be set to on or off with the fillEdges parameter.
	if fillEdges {
		C, R := cols, rows
		// set offset-grid top and bottom row border edges.
		for r := 0; r <= rows; r++ {
			ogrid[r][0] = setTile(ogrid[r][0], QuadTL, fillEdges)
			ogrid[r][0] = setTile(ogrid[r][0], QuadBL, fillEdges)
			ogrid[r][C] = setTile(ogrid[r][C], QuadTR, fillEdges)
			ogrid[r][C] = setTile(ogrid[r][C], QuadBR, fillEdges)
		}
		// set offset-grid left and right column border edges.
		for c := 0; c <= cols; c++ {
			ogrid[0][c] = setTile(ogrid[0][c], QuadTL, fillEdges)
			ogrid[0][c] = setTile(ogrid[0][c], QuadTR, fillEdges)
			ogrid[R][c] = setTile(ogrid[R][c], QuadBL, fillEdges)
			ogrid[R][c] = setTile(ogrid[R][c], QuadBR, fillEdges)
		}
	}
	return ogrid
}

// Set or clear the tile quadrant bits to get one of 16 possible tiles.
func setTile(tile, quad uint8, on bool) uint8 {
	if on {
		return tile | quad // set tile quad bit
	}
	return tile &^ quad // clear tile quad bit
}

// Dump displays the offset-tile grid using the provided ascii tile blocks,
// ie: OffsetSymbols above. Expected to be used for debugging.
func (g OffsetGrid) Dump(symbols map[uint8]string) {
	if symbols == nil {
		symbols = OffsetSymbols // use default
	}
	for _, row := range g {
		for _, val := range row {
			if sym, ok := symbols[val]; ok {
				fmt.Printf("%s", sym)
			} else {
				fmt.Printf("%s", "--")
			}
		}
		fmt.Println()
	}
}
