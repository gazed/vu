// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package grid

import (
	"math/rand"
)

// dense is a corridor only skirmish grid. It is a Prim's maze where the
// dead-ends have been eliminated.  Additionally each side of the grid is
// guaranteed to have one exit to the level exterior.
type dense struct {
	grid // superclass grid.
}

// Generate a maze using a Prim's maze as the basis.  Make a skirmish
// friendly level by knocking out a wall at any dead end and then chopping
// some outside exits if necessary.
func (d *dense) Generate(width, depth int) Grid {
	maze := &primMaze{}
	maze.Generate(width, depth)
	d.cells = maze.cells

	// randomly traverse the grid removing dead ends.
	candidates := d.cellSlice()
	for len(candidates) > 0 {
		index := rand.Intn(len(candidates))
		cell := candidates[index]
		d.fixDeadEnd(cell)
		candidates = append(candidates[:index], candidates[index+1:]...)
	}
	d.ensureExits()
	return d
}

// fixDeadEnd checks if the given cell is a dead end and creates another
// floor if it is.
func (d *dense) fixDeadEnd(u *cell) {
	if !u.isWall {
		walls := d.neighbours(u, allWalls)
		if len(walls) > 2 {
			index := rand.Intn(len(walls))
			u = walls[index]
			u.isWall = allFloors
		}
	}
}

// ensureExits makes sure there is an outside exit on each side.
// The corners are left alone.
func (d *dense) ensureExits() {
	var north, south, east, west []*cell
	xmax, ymax := d.Size()
	for x := 1; x < xmax-1; x++ {
		if d.cells[x][ymax-1].isWall {
			north = append(north, d.cells[x][ymax-1])
		}
		if d.cells[x][0].isWall {
			south = append(south, d.cells[x][0])
		}
	}
	for y := 1; y < ymax-1; y++ {
		if d.cells[xmax-1][y].isWall {
			east = append(east, d.cells[xmax-1][y])
		}
		if d.cells[0][y].isWall {
			west = append(west, d.cells[0][y])
		}
	}
	d.ensureExit(south, xmax-2)
	d.ensureExit(north, xmax-2)
	d.ensureExit(west, ymax-2)
	d.ensureExit(east, ymax-2)
}

// ensureExit chops a hole in the given side.  Sometimes the hole chopped
// can be a dead-end.  Chopping an additional hole in the holes neighbouring
// walls guarantees an exit.
func (d *dense) ensureExit(side []*cell, max int) {
	if len(side) == max {
		index := rand.Intn(len(side))
		u := side[index]
		u.isWall = allFloors

		// ensure the chop gets into the grid by chopping again if necessary.
		walls := d.neighbours(u, allWalls)
		if len(walls) == 3 {
			wallIndex := rand.Intn(len(walls))
			u := walls[wallIndex]
			u.isWall = allFloors
		}
	}
}
