// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

package grid

// FUTURE: improve path finding effeciency using any solution that reduces
// the number of nodes: eg: "A∗-plus path planning algorithm for large-scale grid maps" paper.
// Another possibility is converting the grid to a navMesh.
// In both cases use A* on the reduced graph.

import (
	"container/heap" // for priority queue.
	"math"           // for math.Abs
)

// GetPath is a grid based A* path finder. Implementation based on:
//   - https://en.wikipedia.org/wiki/A*_search_algorithm
//   - https://www.redblobgames.com/pathfinding/a-star/implementation.html
//
// The input path slice is reset to zero length and filled with the path
// results. An empty path is returned if no route was found.
//
// Computation time increases proportional to the the grid size.
func GetPath(grid Grid, start, goal Pos, pt PathTraverse, path *[]Pos) {
	startID := grid.Index(start.R, start.C)
	cameFrom := map[uint]Pos{startID: start}
	costSoFar := map[uint]float64{startID: 0}
	frontier := &priorityQueue{queuePosition{pos: start, priority: 0}}
	heap.Init(frontier)
	for frontier.Len() > 0 {
		current := heap.Pop(frontier).(queuePosition).pos
		if current == goal {
			break // success
		}
		for _, next := range findNeighbours(grid, pt.Directions, current) {
			currentID := grid.Index(current.R, current.C)
			nextID := grid.Index(next.R, next.C)
			newCost := costSoFar[currentID] + pt.Cost(current, next)
			if csf, ok := costSoFar[nextID]; !ok || newCost < csf {
				costSoFar[nextID] = newCost
				priority := newCost + pt.Cost(next, goal) // add estimated cost to goal.
				heap.Push(frontier, queuePosition{pos: next, priority: priority})
				cameFrom[nextID] = current
			}
		}
	}

	// unwind cameFrom results to get path from start to goal
	*path = (*path)[:0] // reset to reuse existing memory.
	goalID := grid.Index(goal.R, goal.C)
	if last, ok := cameFrom[goalID]; ok {
		*path = append(*path, last, goal)
		for {
			lastID := grid.Index(last.R, last.C)
			prev := cameFrom[lastID]
			prevID := grid.Index(prev.R, prev.C)
			*path = append(*path, Pos{})   // create space at end of slice...
			copy((*path)[1:], (*path)[0:]) // ...move old values up one spot...
			(*path)[0] = prev              // ...and insert new value at beginning.
			if prevID == startID {
				break // reached start.
			}
			last = prev
		}
	}
}

// Pos is a grid location. Row and Col can be interpreted as needed
// by the grid display, eg: Row:YorZ(vertical) Col:X(horizontal)
type Pos struct {
	R, C uint // Row:Column - internally equivalent to X:Y
}

// =============================================================================
// Utility to find graph neighbours and estimate travel cost.

// Find all traversable grid neighbours from the given grid position.
func findNeighbours(g Grid, directions []int, p Pos) (neighbours []Pos) {
	rows, cols := g.Size()
	for i := 0; i < len(directions); i += 2 {
		dy, dx := directions[i], directions[i+1]
		pr, pc := int(p.R)+dy, int(p.C)+dx
		if pr >= 0 && pr < int(rows) &&
			pc >= 0 && pc < int(cols) &&
			g[uint(pr)][uint(pc)] == EMPTY {
			if dy == 0 || dx == 0 {
				// can always move in the cardinal directions.
				neighbours = append(neighbours, Pos{uint(pr), uint(pc)})
			} else {
				// can only move diagonally if one of the cardinal directions
				// is also traversable.
				if g[p.R][uint(pc)] == EMPTY || g[uint(pr)][p.C] == EMPTY {
					neighbours = append(neighbours, Pos{uint(pr), uint(pc)})
				}
			}
		}
	}
	return neighbours
}

// PathTraverse sets the type of grid traversal combined
// with the associated traversal cost.
type PathTraverse struct {
	Directions []int
	Cost       func(gridA, gridB Pos) float64
}

// FourWay graph traversal.
// Manhattan is used for estimating costs in grids with 4 directions.
var FourWay = PathTraverse{
	Directions: []int{
		-1, +0, // top
		+1, +0, // bot
		+0, -1, // left
		+0, +1, // right
	},
	Cost: func(gridA, gridB Pos) float64 {
		return math.Abs(float64(gridA.R-gridB.R)) + math.Abs(float64(gridA.C-gridB.C))
	},
}

// EightWay graph traversal.
// Chebyshev is used for estimating costs in grids with 8 directions.
var EightWay = PathTraverse{
	Directions: []int{
		-1, +0, // top
		+1, +0, // bot
		+0, -1, // left
		+0, +1, // right
		-1, -1, // top left
		-1, +1, // top right
		+1, -1, // bot left
		+1, +1, // bot right
	},
	Cost: func(gridA, gridB Pos) float64 {
		return max(math.Abs(float64(gridA.R-gridB.R)), math.Abs(float64(gridA.C-gridB.C)))
	},
}

// =============================================================================
// internal heap struct that attachs a priority to a grid location.

// queuePosition wraps a location with a priority value used for sorting.
type queuePosition struct {
	pos      Pos     // grid location.
	priority float64 // sort the grid locations.
}

// priorityQueue is a min-heap of queuePosition.
// Based on example from https://golang.org/pkg/container/heap/
type priorityQueue []queuePosition

// Len implements sort.Interface.
func (h priorityQueue) Len() int { return len(h) }

// Less implements sort.Interface.
func (h priorityQueue) Less(i, j int) bool { return h[i].priority < h[j].priority }

// Swap implements sort.Interface.
func (h priorityQueue) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

// Push implements heap.Interface.
func (h *priorityQueue) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the
	// slice's length, not just its contents.
	*h = append(*h, x.(queuePosition))
}

// Pop implements heap.Interface.
func (h *priorityQueue) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// =============================================================================

// line returns a line by Bresenham's algorithm. Expects uint grid positions.
// Returns a slice of grid positions from point:0 to point:1.
// from: https://rosettacode.org/wiki/Bitmap/Bresenham%27s_line_algorithm#Go
func Line(x0, y0, x1, y1 uint) (line []Pos) {
	// implemented straight from WP pseudocode
	dx := int(x1) - int(x0)
	if dx < 0 {
		dx = -dx
	}
	dy := int(y1) - int(y0)
	if dy < 0 {
		dy = -dy
	}
	var sx, sy int
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}
	err := dx - dy
	for {
		line = append(line, Pos{C: x0, R: y0})
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 = uint(int(x0) + sx)
		}
		if e2 < dx {
			err += dx
			y0 = uint(int(y0) + sy)
		}
	}
	return line
}
