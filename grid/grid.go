// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

// Package grid provides a utilities for constructing 2D areas.
package grid

import (
	"fmt"
	"log/slog"
	"math"
)

// Grid is a rectangular area where each grid element can be set
// to a value (65535). The grid is invalid if any dimension is 0.
type Grid [][]uint16 // Defaults to all EMPTY

const EMPTY = 0 // convention for an traversable empty location.
const SOLID = 1 // convention for an blocking wall location.

// MakeGrid returns a grid for the given dimensions.
func MakeGrid(rows, columns uint) (grid Grid) {
	if rows <= 0 || columns <= 0 {
		return grid // invalid, return un-inialized grid.
	}
	grid = make(Grid, rows)
	for r := range grid {
		grid[r] = make([]uint16, columns)
	}
	return grid
}

// Size returns the grid size as the number of
// rows (depth) and columns (width).
func (g Grid) Size() (rows, cols uint) {
	if len(g) == 0 {
		return 0, 0 // invalid grid.
	}
	return uint(len(g)), uint(len(g[0]))
}

// In returns true if the cell is within the grid boundaries.
func (g Grid) In(r, c uint) bool {
	rows, cols := g.Size()
	return r < rows && c < cols
}

// Valid returns true if the grid size is non-zero and rectangular.
func (g Grid) Valid() bool {
	rows, cols := g.Size()
	if rows == 0 || cols == 0 {
		return false
	}

	// require all columns to be the same size.
	for _, row := range g {
		if cols != uint(len(row)) {
			return false // not rectangular
		}
	}
	return true
}

// Index returns a single value for the given grid indices.
// Returns 0 and logs error for invalid grid locations.
func (g Grid) Index(row, col uint) uint {
	if g.In(row, col) {
		_, cols := g.Size()
		return row*cols + col
	}
	slog.Error("Grid:Index invalid indices")
	return math.MaxUint // invalid - developer error to get here.
}

// Indices returns the grid indices from the given single index.
func (g Grid) Indices(index uint) (row, col uint) {
	_, cols := g.Size()
	row = index / cols // integer division
	col = index % cols // modulo
	if g.In(row, col) {
		return row, col
	}
	slog.Error("Grid:Indices invalid index")
	return math.MaxUint, math.MaxUint // invalid - developer error to get here.
}

// Fill all grid elements to the given value.
func (g Grid) Fill(val uint16) {
	for r, row := range g {
		for c := range row {
			g[r][c] = val
		}
	}
}

// Get the value of the given grid element.
// Returns 0 and logs an error for invalid grid locations.
func (g Grid) Get(r, c uint) uint16 {
	if g.In(r, c) {
		return g[r][c]
	}
	slog.Error("Grid:Get invalid location")
	return 0
}

// Set the given grid element to the given value.
func (g Grid) Set(r, c uint, val uint16) {
	if g.In(r, c) {
		g[r][c] = val
	}
}

// SetPositions sets all the grid elements identified by positions
// to the given value.  Invalid positions are ignored.
func (g Grid) SetPositions(positions []Pos, val uint16) {
	for _, p := range positions {
		if g.In(p.R, p.C) {
			g[p.R][p.C] = val
		}
	}
}

// GridSymbols is the default for grid debug dumps.
var GridSymbols = map[uint16]string{0: "  ", 1: "██"}

// Dump prints a grid for debugging purposes.
// This looks better with some fixed fonts than others.
// The grid is dumped such that the 0,0 is at the top left,
// ie: 0,0 will be the first character of the first line.
func (g Grid) Dump(symbols map[uint16]string) {
	if symbols == nil {
		symbols = GridSymbols // use default
	}
	for _, row := range g {
		for _, set := range row {
			if sym, ok := symbols[set]; ok {
				fmt.Printf("%s", sym)
			} else {
				fmt.Printf("%s", "--")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
