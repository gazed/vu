// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package grid

import (
	"math/rand"
)

// cave holds a cavern like grid structure.
type cave struct {
	grid // superclass grid
}

// Generate a cave using cellular automata generation algorithm.
// The algorithm based on:
//     http://www.roguebasin.com/index.php?title=Cellular_Automata_Method_
//            for_Generating_Random_Cave-Like_Levels
//
// FUTURE: Could use flood filling to detect isolated caves if they
//         are not a useful gameplay feature.
func (c *cave) Generate(width, depth int) Grid {
	c.create(width, depth, allFloors)
	scratch := make([][]*cell, len(c.cells))
	for x, _ := range scratch {
		scratch[x] = make([]*cell, len(c.cells[0]))
		for y, _ := range c.cells[x] {
			scratch[x][y] = &cell{x, y, allFloors}
		}
	}

	// randomly fill the map with walls.
	for x, _ := range c.cells {
		for y, _ := range c.cells[x] {
			c.cells[x][y].isWall = rand.Intn(100) < 40 // 40% chance of wall.
		}
	}

	// iterate, treating each wall with a cellular automaton live/die depending on
	// the number of neighbours.
	iterations := 4
	makeWall := func(w3x3, w5x5 int) bool { return w3x3 >= 5 || w5x5 <= 3 }
	for cnt := 0; cnt < iterations; cnt++ {
		c.runGeneration(scratch, makeWall)
	}

	// run three more generations to clean up single walls.
	iterations = 3
	makeWall = func(w3x3, w5x5 int) bool { return w3x3 >= 5 }
	for cnt := 0; cnt < iterations; cnt++ {
		c.runGeneration(scratch, makeWall)
	}
	return c
}

// runGeneration applies the cell automation rule to the current grid.
// Results are stored in a temporary grid and then copied back once the
// generation has finished.
func (c *cave) runGeneration(scratch [][]*cell, makeWall func(w3x3, w5x5 int) bool) {
	for x, _ := range c.cells {
		for y, cell := range c.cells[x] {

			// A tile is a wall if the 3x3 region around it has at least 5 walls.
			// Otherwise it is a floor.
			w3, w5 := c.wallCount(cell)
			scratch[x][y].isWall = makeWall(w3, w5)
		}
	}

	// copy the generation back into main grid.
	for x, _ := range c.cells {
		for y, _ := range c.cells[x] {
			c.cells[x][y].isWall = scratch[x][y].isWall
		}
	}
}

// wallCount returns the number of walls contained in the 3x3 and 5x5 regions
// around the given cell. The wall count includes the given cell.
func (c *cave) wallCount(u *cell) (w3x3, w5x5 int) {
	maxx, maxy := len(c.cells), len(c.cells[0])
	for cx := u.x - 2; cx <= u.x+2; cx++ {
		for cy := u.y - 2; cy <= u.y+2; cy++ {
			if cx < maxx && cy < maxy && cx >= 0 && cy >= 0 {
				if c.cells[cx][cy].isWall {
					w5x5++
					if cx >= u.x-1 && cx <= u.x+1 && cy >= u.y-1 && cy <= u.y+1 {
						w3x3++
					}
				}
			} else { // outside counts as walls.
				w5x5++
				if cx >= u.x-1 && cx <= u.x+1 {
					w3x3++
				}
			}
		}
	}
	return
}
