// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package grid

import (
	"math/rand"
)

// sparse is a minimally populated skirmish grid. The grid is created by
// randomly adding walls around cells, always ensuring that there are no
// dead ends - there must always be at least 2 ways out of any cell.
//
// Basic algorithm : start with empty grid. For each cell in the grid,
// randomly choose one cell and determine if a wall can be placed at one
// of its neighbour cells (Nx, max 4). This is validated by testing each
// of the Nx neighbour cells (max 16) and ensuring no dead end
// (i.e. only 1 floor). Remove that cell from the grid list.
type sparse struct {
	grid // superclass grid
}

// Generate the grid by randomly traversing all the cells and adding random
// walls as long as there are more than two ways out of the cell.
func (s *sparse) Generate(width, depth int) Grid {
	s.create(width, depth, allFloors)

	// randomly check each cell in the grid once.
	candidates := s.cellSlice()
	for len(candidates) > 0 {
		index := rand.Intn(len(candidates))
		u := candidates[index]
		s.addWall(u)
		candidates = append(candidates[:index], candidates[index+1:]...)
	}
	return s
}

// addWall adds a wall to the grid if by doing so the grid remains valid.
// A wall is added randomonly to one of the given cells floors
// if the given cell currently has more than 2 floors.
func (s *sparse) addWall(u *cell) {
	if !u.isWall {
		floors := s.neighbours(u, allFloors)
		if len(floors) > 2 {
			if s.checkLevel(floors) {
				index := rand.Intn(len(floors))
				u = floors[index]
				u.isWall = allWalls
				s.cells[u.x][u.y].isWall = u.isWall
			}
		}
	}
}

// checkLevel ensures that the grid remains valid if any of the given floors
// are made into walls. It does this by putting a temporary wall at each of the
// potential floors and checking that the grid has no dead ends.
//
// This is the part that makes the grid rather sparse since candidates will be
// rejected if any of the floors fails to be a valid wall position.
func (s *sparse) checkLevel(floors []*cell) bool {
	for _, floor := range floors {

		// this temporary wall will be turned back into a floor before exiting.
		floor.isWall = allWalls

		// check around the affected floor for dead ends.
		xmax, ymax := s.Size()
		for xcnt := floor.x - 2; xcnt < floor.x+2; xcnt++ {
			for ycnt := floor.y - 2; ycnt < floor.y+2; ycnt++ {
				if xcnt >= 0 && ycnt >= 0 && xcnt < xmax && ycnt < ymax {
					if len(s.neighbours(s.cells[xcnt][ycnt], allFloors)) < 2 {
						floor.isWall = allFloors
						return false
					}
				}
			}
		}
		floor.isWall = allFloors
	}
	return true
}
