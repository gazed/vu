// Copyright © 2014-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package grid

// Flow is based mainly on continuum crowds [Treuille 2006]
// See implemetation example at:
//    http://www.youtube.com/watch?v=bovlsENv1g4
// See tutorials and explanations from:
//    http://leifnode.com/2013/12/flow-field-pathfinding/
//    http://computationalbubblegum.blogspot.ca/2014/06/implementing-continuum-crowds-part-1.html
//    http://computationalbubblegum.blogspot.ca/2014/06/implementing-continuum-crowds-part-2.html

import (
	"math"
)

// Flow creates a map to help units move towards their goal. The flow is
// initialized with a map and a goal. Afterwards each unit can query the flow
// for the best diretion to move towards their goal.
type Flow interface {

	// Create a new flow field based on the given goal location.
	// This can be called to update maps as the goal changes.
	Create(goalx, goaly int) // Call once before calling Next.

	// Next, based on the current grid location gx, gy, returns the
	// next grid location nx, ny. 0, 0 is returned if the current
	// location is the goal. 9, 9 is returned if the given location
	// is invalid.
	Next(gx, gy int) (nx, ny int)
}

// NewFlow creates a flow based on a plan.
func NewFlow(p Plan) Flow { return newFlow(p) }

// public interface
// =============================================================================
// private implementaiton

// flow is the default implementation of Flow. It keeps a flow field map
// for the given plan.
type flow struct {
	xsz, ysz   int     // floor plan x,y dimensions.
	costmap    Plan    // base map with impassable and avoidance areas.
	goalmap    [][]int // holds the cost to the goal for each cell.
	flowmap    [][]int // direction to goal for each cell.
	neighbours []int   // scratch for calculating valid neighbours.
	candidates []int   // map crawl candidates. For creating goalmap.
	max        int     // impassable value for flow and goal maps.
}

// Direction constants.
const (
	goal  = iota // goal is 0 meaning no movement necessary or possible.
	north        // x  , y+1
	ne           // x+1, y+1
	east         // x+1, y
	se           // x+1, y-1
	south        // x  , y-1
	sw           // x-1, y-1
	west         // x-1, y
	nw           // x-1, y+1
)

// newFlow creates a flow map towards the given goal using the plan
// as the cost map. Note that using MaxUint8 to mean impassable implies
// flow map sizes less than 512x512. Flow maps really should be limited
// to 100x100 or smaller.
func newFlow(p Plan) *flow {
	f := &flow{max: math.MaxUint8}
	f.xsz, f.ysz = p.Size()
	f.costmap = p
	f.flowmap = make([][]int, f.ysz)
	f.goalmap = make([][]int, f.ysz)
	for y := range f.flowmap {
		f.flowmap[y] = make([]int, f.xsz)
		f.goalmap[y] = make([]int, f.xsz)
	}
	f.neighbours = make([]int, 8) // max neighbours is 8.
	return f
}

// Create implements Flow.
func (f *flow) Create(goalx, goaly int) {
	f.createGoalmap(goalx, goaly) // create goal map from cost map.
	f.createFlowmap(goalx, goaly) // create flow map from goal map.
}

// Next implements Flow.
func (f *flow) Next(gx, gy int) (nx, ny int) {
	switch f.flowmap[gx][gy] {
	case north:
		return 0, 1
	case ne:
		return 1, 1
	case east:
		return 1, 0
	case se:
		return 1, -1
	case south:
		return 0, -1
	case sw:
		return -1, -1
	case west:
		return -1, 0
	case nw:
		return -1, 1
	case f.max:
		return 9, 9 // Invalid.
	}
	return 0, 0
}

// createGoalmap creates the goal map from the cost map.
// This spreads out from the goalnode until each reachable node
// has been processed.
func (f *flow) createGoalmap(goalx, goaly int) {

	// reset all node costs to a large values.
	for y, row := range f.goalmap {
		for x := range row {
			f.goalmap[x][y] = f.max
			f.flowmap[x][y] = f.max
		}
	}

	// set the goal node to value 0 and push it on the open list.
	f.goalmap[goalx][goaly] = 0
	f.candidates = f.candidates[:0] // reset keeping memory.
	f.candidates = append(f.candidates, f.id(goalx, goaly))

	// while there are nodes on the open list.
	for len(f.candidates) > 0 {

		// get the first candidate, removing it from the candidate list.
		candidate := f.candidates[0]
		f.candidates = append(f.candidates[:0], f.candidates[1:]...)
		x, y := f.at(candidate)

		// process the candidates immediate neighbours ignoring diagonals.
		for _, dir := range f.directNeighbours(x, y) {
			endNodeCost := f.goalmap[x][y]
			nx, ny := x, y
			switch dir {
			case north:
				ny = y + 1
			case east:
				nx = x + 1
			case south:
				ny = y - 1
			case west:
				nx = x - 1
			}

			// check and update the cost for each neighbour.
			endNodeCost += f.cost(nx, ny)
			if endNodeCost < f.goalmap[nx][ny] {
				neighbourID := f.id(nx, ny)

				// Set neighbour node cost and add it as a candidate.
				f.goalmap[nx][ny] = endNodeCost
				if !f.alreadyCandidate(neighbourID) {
					f.candidates = append(f.candidates, neighbourID)
				}
			}
		}
	}
}

// createFlowmap creates the flow map from the goal map.
func (f *flow) createFlowmap(goalx, goaly int) {
	for y, row := range f.goalmap {
		for x := range row {
			costToGoal := f.max
			leastCost := f.max

			// ignore goalmaps spots that are impassable
			if f.goalmap[x][y] == f.max {
				f.flowmap[x][y] = f.max
			} else {

				// the direction is the lowest cost of the eight neighbours.
				neighbours := f.findNeighbours(x, y)
				for _, dir := range neighbours {
					switch dir {
					case north:
						costToGoal = f.goalmap[x][y+1]
					case ne:
						costToGoal = f.goalmap[x+1][y+1]
					case east:
						costToGoal = f.goalmap[x+1][y]
					case se:
						costToGoal = f.goalmap[x+1][y-1]
					case south:
						costToGoal = f.goalmap[x][y-1]
					case sw:
						costToGoal = f.goalmap[x-1][y-1]
					case west:
						costToGoal = f.goalmap[x-1][y]
					case nw:
						costToGoal = f.goalmap[x-1][y+1]
					}
					if costToGoal < leastCost {
						leastCost = costToGoal
						f.flowmap[x][y] = dir
					}
				}
			}
		}
	}
	f.flowmap[goalx][goaly] = goal
}

// Find the all neighbours including diagonals. Relies on f.neighbours
// scratch variable. Returns the direction of the valid neighbour.
func (f *flow) findNeighbours(x, y int) []int {
	f.neighbours = f.neighbours[:0] // reset while preserving memory.
	if y+1 < f.ysz {
		f.neighbours = append(f.neighbours, north)
	}
	if x+1 < f.xsz {
		f.neighbours = append(f.neighbours, east)
		if y+1 < f.ysz {
			f.neighbours = append(f.neighbours, ne)
		}
		if y-1 >= 0 {
			f.neighbours = append(f.neighbours, se)
		}
	}
	if y-1 >= 0 {
		f.neighbours = append(f.neighbours, south)
	}
	if x-1 >= 0 {
		f.neighbours = append(f.neighbours, west)
		if y+1 < f.ysz {
			f.neighbours = append(f.neighbours, nw)
		}
		if y-1 >= 0 {
			f.neighbours = append(f.neighbours, sw)
		}
	}
	return f.neighbours
}

// Find the N,S,E,W neighbours. Relies on f.neighbours scratch variable.
// Returns the direction of the valid neighbour.
func (f *flow) directNeighbours(x, y int) []int {
	f.neighbours = f.neighbours[:0] // reset while preserving memory.
	if y+1 < f.ysz {
		f.neighbours = append(f.neighbours, north)
	}
	if x+1 < f.xsz {
		f.neighbours = append(f.neighbours, east)
	}
	if y-1 >= 0 {
		f.neighbours = append(f.neighbours, south)
	}
	if x-1 >= 0 {
		f.neighbours = append(f.neighbours, west)
	}
	return f.neighbours
}

// scan the list of candidates for the given identifier.
func (f *flow) alreadyCandidate(id int) bool {
	for _, candidateID := range f.candidates {
		if id == candidateID {
			return true
		}
	}
	return false
}

// The cost for a plan is very high for walls and 1 for open areas.
func (f *flow) cost(x, y int) int {
	if f.costmap.IsOpen(x, y) {
		return 1 // open area.
	}
	return f.max // wall.
}

// Turn x,y map indicies to unique identifiers.
func (f *flow) id(x, y int) int { return x*f.ysz + y }

// Turn unique identifiers to x,y map indicies.
func (f *flow) at(id int) (x, y int) { return id / f.ysz, id % f.ysz }
