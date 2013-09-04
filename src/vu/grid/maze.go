// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package grid

import (
	"math/rand"
	"time"
)

// primMaze is a holds a Randomized Prim's Algorithm: (see wikipedia)
type primMaze struct {
	grid // superclass grid
}

// Generate a maze using a Randomized Prim's Algorithm: (see wikipedia)
//   - Start with a grid full of walls.
//   - Pick a cell, mark it as part of the grid. Add the walls of the cell to
//     the wall list.
//   - While there are walls in the list:
//     - Pick a random wall from the list. If the cell on the opposite side
//       isn't in the maze yet:
//       - Make the wall a passage and mark the cell on the opposite side as
//         part of the grid.
//       - Add the neighboring walls of the cell to the wall list.
//     - If the cell on the opposite side already was in the grid, remove the
//       wall from the list.
func (pm *primMaze) Generate(width, depth int) Grid {
	pm.create(width, depth, allWalls)
	randomGenerator := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	// Pick a cell, mark it as part of the grid. Add the walls of the cell to
	// the wall list.
	start := pm.cells[1][1]
	start.isWall = allPassages
	walls := []*cell{pm.north(start), pm.south(start), pm.west(start), pm.east(start)}

	// While there are walls in the list:
	for len(walls) > 0 {

		// Pick a random wall from the list. If the cell on the opposite side
		// isn't in the maze yet...
		randomWall := randomGenerator.Intn(len(walls))
		wall := walls[randomWall]
		if link := pm.link(wall); link != nil {

			// ... then: make the wall a passage and mark the cell on the
			// opposite side as part of the maze.
			pm.cells[wall.x][wall.y].isWall = allPassages
			pm.cells[link.x][link.y].isWall = allPassages

			// Add the neighboring walls of the new Passage  to the wall list.
			newWalls := []*cell{pm.north(link), pm.south(link), pm.west(link), pm.east(link)}
			walls = append(walls, newWalls...)
		} else {
			// ... otherwise: if the cell on the opposite side was already
			// in the maze remove the wall from the list.
			walls = append(walls[:randomWall], walls[randomWall+1:]...)
		}
	}
	return pm
}

// link attempts to return a cell that connects to the existing grid.
// Return nil if no new link can be created.
func (pm *primMaze) link(wall *cell) (u *cell) {
	if wall != nil {
		if u = pm.linkPair(pm.north(wall), pm.south(wall)); u == nil {
			u = pm.linkPair(pm.west(wall), pm.east(wall))
		}
	}
	return
}
func (pm *primMaze) linkPair(sideA *cell, sideB *cell) *cell {
	if sideA != nil && sideB != nil {
		if !sideA.isWall && sideB.isWall {
			return sideB
		}
		if !sideB.isWall && sideA.isWall {
			return sideA
		}
	}
	return nil
}
