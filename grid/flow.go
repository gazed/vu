// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

package grid

// flow is based mainly on continuum crowds [Treuille 2006]

import (
	"math"
	"slices"
)

// Flow generates a flow field map for a grid.
// Flow fields helps move units towards their goal. The flow is initialized
// with a grid and a goal. Afterwards the flow can be queried for the best
// direction to move towards the goal.
type Flow struct {
	cols, rows int     // grid x,y dimensions.
	grid       Grid    // base map with impassable and avoidance areas.
	costmap    [][]int // holds the cost to the goal for each cell.
	flowmap    [][]int // direction to goal for each cell.
	neighbours []int   // scratch for calculating valid neighbours.
	candidates []uint  // map crawl candidates. For creating costmap.
}

const (
	// maxFlowSize also means impassable, implying that flow map sizes are
	// less than 512x512. Flow maps should be limited to 100x100 or smaller
	// for performance reasons.
	maxFlowSize = math.MaxUint8
)

// NewFlow creates a flow map towards the given goal using the plan
// as the cost map. Note that using MaxUint8 to mean impassable implies
// flow map sizes less than 512x512. Flow maps really should be limited
// to 100x100 or smaller.
func NewFlow(g Grid, goalRow, goalCol uint) *Flow {
	f := &Flow{}
	c, r := g.Size()
	f.cols, f.rows = int(c), int(r)
	f.grid = g
	f.flowmap = make([][]int, f.rows)
	f.costmap = make([][]int, f.rows)
	for y := range f.flowmap {
		f.flowmap[y] = make([]int, f.cols)
		f.costmap[y] = make([]int, f.cols)
	}
	f.neighbours = make([]int, 8)     // scratch: max neighbours is 8.
	f.createCostmap(goalRow, goalCol) // create goal map from cost map.
	f.createFlowmap(goalRow, goalCol) // create flow map from goal map.
	return f
}

// Direction constants.
const (
	goal  = iota // default 0 means no movement necessary or possible.
	north        // x  , y-1
	ne           // x+1, y-1
	east         // x+1, y
	se           // x+1, y+1
	south        // x  , y+1
	sw           // x-1, y+1
	west         // x-1, y
	nw           // x-1, y-1
)

// Next, based on the current grid location gx, gy, returns the
// next grid location nx, ny. 0, 0 is returned if the current
// location is the goal. 9, 9 is returned if the given location
// is invalid.
func (f *Flow) Next(row, col int) (nx, ny int) {
	switch f.flowmap[row][col] {
	case north:
		return +0, -1 // up
	case ne:
		return +1, -1 // up right
	case east:
		return +1, +0 // right
	case se:
		return +1, +1 // down right
	case south:
		return +0, +1 // down
	case sw:
		return -1, +1 // down left
	case west:
		return -1, +0 // left
	case nw:
		return -1, -1 // up left
	case maxFlowSize:
		return 9, 9 // Invalid.
	}
	return 0, 0 // reached the goal.
}

// createCostmap creates the goal map from the grid (cost map).
// This spreads out from the goalnode until each reachable node
// has been processed.
func (f *Flow) createCostmap(goalRow, goalCol uint) {

	// reset all node costs to a large values.
	for y, row := range f.costmap {
		for x := range row {
			f.costmap[y][x] = maxFlowSize
			f.flowmap[y][x] = maxFlowSize
		}
	}

	// set the goal node to value 0 and push it on the open list.
	f.costmap[goalRow][goalCol] = 0
	f.candidates = f.candidates[:0] // reset keeping memory.
	f.candidates = append(f.candidates, f.grid.Index(goalRow, goalCol))

	// while there are nodes on the open list.
	for len(f.candidates) > 0 {

		// get the first candidate, removing it from the candidate list.
		candidate := f.candidates[0]
		f.candidates = append(f.candidates[:0], f.candidates[1:]...)
		y, x := f.grid.Indices(candidate)

		// process the candidates immediate neighbours ignoring diagonals.
		for _, dir := range f.directNeighbours(int(y), int(x)) {
			endNodeCost := f.costmap[y][x]
			nx, ny := x, y
			switch dir {
			case north:
				ny = y - 1
			case east:
				nx = x + 1
			case south:
				ny = y + 1
			case west:
				nx = x - 1
			}

			// check and update the cost for each neighbour.
			endNodeCost += f.cost(ny, nx)
			if endNodeCost < f.costmap[ny][nx] {
				neighbourID := f.grid.Index(ny, nx)

				// Set neighbour node cost and add it as a candidate.
				f.costmap[ny][nx] = endNodeCost
				if !slices.Contains(f.candidates, neighbourID) {
					f.candidates = append(f.candidates, neighbourID)
				}
			}
		}
	}
}

// createFlowmap creates the flow map from the goal map.
func (f *Flow) createFlowmap(goalRow, goalCol uint) {
	neighbours := []int{ // x:y neighbour directions.
		+0, -1, // north
		+1, -1, // ne
		+1, +0, // east
		+1, +1, // se
		+0, +1, // south
		-1, +1, // sw
		-1, +0, // west
		-1, -1, // nw
	}
	for y, row := range f.costmap {
		for x := range row {
			costToGoal := maxFlowSize
			leastCost := maxFlowSize

			// ignore costmap spots that are impassable
			if f.costmap[y][x] == maxFlowSize {
				f.flowmap[y][x] = maxFlowSize
				continue
			}

			// check the cost at each valid neighbour
			// and take the direction with the lowest cost.
			for i := 0; i < len(neighbours); i += 2 {
				nx, ny := neighbours[i], neighbours[i+1]
				gx, gy := uint(x+nx), uint(y+ny)
				if !f.grid.In(gy, gx) {
					continue // invalid grid location.
				}
				costToGoal = f.costmap[gy][gx]
				if costToGoal < leastCost {
					leastCost = costToGoal

					// direction is the neighbour index.
					f.flowmap[y][x] = i/2 + 1 // direction.
				}
			}
		}
	}
	f.flowmap[goalRow][goalCol] = goal
}

// Find the N,S,E,W neighbours. Relies on f.neighbours scratch variable.
// Returns the direction of the valid neighbour.
func (f *Flow) directNeighbours(y, x int) []int {
	f.neighbours = f.neighbours[:0] // reset while preserving memory.
	if y+1 < f.rows {
		f.neighbours = append(f.neighbours, south)
	}
	if x+1 < f.cols {
		f.neighbours = append(f.neighbours, east)
	}
	if y-1 >= 0 {
		f.neighbours = append(f.neighbours, north)
	}
	if x-1 >= 0 {
		f.neighbours = append(f.neighbours, west)
	}
	return f.neighbours
}

// The cost for a plan is very high for walls and 1 for open areas.
func (f *Flow) cost(y, x uint) int {
	if f.grid[y][x] == EMPTY {
		return 1 // traversable area.
	}
	return maxFlowSize // wall.
}

// Dump the flowmap where 0,0 is the top left corner.
func (f *Flow) Dump() {
	for y := range f.rows {
		for x := range f.cols {
			dir := f.flowmap[y][x]
			switch dir {
			case north:
				print("↑")
			case ne:
				print("↗")
			case east:
				print("→")
			case se:
				print("↘")
			case south:
				print("↓")
			case sw:
				print("↙")
			case west:
				print("←")
			case nw:
				print("↖")
			case maxFlowSize:
				print("X") // marks impassable areas:
			default:
				print("o") // marks the goal.
			}
		}
		println()
	}
	println()
}
