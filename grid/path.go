// Copyright © 2014-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package grid

// path is an A* implementation that works with 2D level layouts that conform
// to Plan. See:
//     http://en.wikipedia.org/wiki/A*_search_algorithm
//     http://www.policyalmanac.org/games/aStarTutorial.htm
//     http://www.gamasutra.com/view/feature/131505/toward_more_realistic_pathfinding.php?print=1
// Other pathfinding links:
//     http://www.ai-blog.net/archives/000152.html (discusses navigation meshes)
//     http://grail.cs.washington.edu/projects/crowd-flows/ (flowfield algorithm)

// Design Notes:
//   • This A* implementation has been lightly optimized for short routes
//     on small grids.
//   • Scratch variables avoid reallocating memory for sequential path lookups,
//     but may not be the way to go if concurrent access is needed.
//   • Check all changes against the benchmark provided.

// Path finds a route between two points within a given Plan.
type Path interface {

	// Find calculates the path from one point to another. A path of
	// x, y points is returned upon success. The returned path will be
	// empty if there was no way to get to the destination point.
	// The to: tx, ty, and from: fx, fy points are expected to be
	// traversable spots.
	Find(fx, fy, tx, ty int) (path []int)
}

// NewPath creates a new path finder for the given Plan p.
func NewPath(p Plan) Path { return newPath(p) }

// =============================================================================

// path is the default implementation of Path.
type path struct {
	fp           Plan          // floor plan.
	xsz, ysz     int           // floor plan x,y dimensions
	candidates   []*node       // nodes to be considered.
	route        []int         // scratch for returning path points.
	neighbours   []*node       // scratch for calculating current neighbours.
	orthMoveCost int           // cost for moving up,down, left, right.
	diagMoveCost int           // cost for moving diagonally.
	nodes        map[int]*node // reuse nodes across multiple calls.
}

// newPath is used by test cases to get an initialized path instance.
func newPath(fp Plan) *path {
	p := &path{}
	p.fp = fp
	p.orthMoveCost = 1 // heuristic horizontal, vertical move cost.
	p.diagMoveCost = 2 // heuristic diagonal move cost.
	p.xsz, p.ysz = fp.Size()
	p.candidates = make([]*node, 50) // guess for small paths.
	p.neighbours = make([]*node, 8)  // max neighbours is 8.
	p.route = []int{}
	p.nodes = make(map[int]*node, p.xsz*p.ysz)
	return p
}

// nodeState values.
const (
	isChecked   = 1 << iota // Node has been tried in path.
	isCandidate = 1 << iota // Node available for trial in path.
)

// id calculates a unique id for a given x, y value. Each node is labelled
// with an id that can be used as map lookups.
func (p *path) id(x, y int) int { return x*p.xsz + y }

// Find calculates the path from one point to another. A path of intermediate
// points is returned upon success. The returned path will be empty if there
// was no way to get to the destination point. The from, to points are expected
// to be valid spots within the path's initialized floor plan.
func (p *path) Find(fx, fy, tx, ty int) (path []int) {
	if !p.fp.IsOpen(fx, fy) || !p.fp.IsOpen(tx, ty) {
		return p.route[:0] // no path found, return empty list.
	}

	// reset any previous path data.
	p.candidates = p.candidates[:0]
	for _, n := range p.nodes {
		n.state = 0
	}

	// create the initial candidate set from the start node.
	start := newNode(p.id(fx, fy), fx, fy, 0)
	start.projCost = start.heuristic(tx, ty)
	p.candidates = []*node{start}
	p.nodes[start.id] = start

	// continue while there are still nodes to be tried.
	destinationNode := p.id(tx, ty)
	for len(p.candidates) > 0 {
		current := p.closest()
		if current.id == destinationNode {
			return p.traceBack(p.route[0:0], current) // backtrack to get path from end point.
		}

		// query neighbours for possible path nodes.
		for _, neighbour := range p.neighbourNodes(current) {
			if neighbour.state != isCandidate {
				neighbour.from = current.id
				p.candidates = append(p.candidates, neighbour)
				neighbour.state = isCandidate
			}

			// Update the projected cost for all neighbours since candidates may be
			// revisited.
			neighbour.projCost = neighbour.pathCost + neighbour.heuristic(tx, ty)
		}
	}
	return p.route[0:0] // no path found, return empty list.
}

// closest returns the node with the lowest f score from the list of open nodes.
// The returned node is removed from the open node list and added to the closed.
func (p *path) closest() *node {
	closest := p.candidates[0]
	index := 0
	for cnt, n := range p.candidates {
		if n.projCost < closest.projCost {
			closest = n
			index = cnt
		}
	}
	p.candidates = append(p.candidates[:index], p.candidates[index+1:]...) // remove closest.
	closest.state = isChecked
	return closest
}

// neighbourNodes creates the valid neighbour nodes for the given node. The nodes
// path distance variables are not set. Neighbours are valid if they are
//    • inside the floor plan.
//    • passable.
//    • a diagonal with two passable adjacent neighbours.
func (p *path) neighbourNodes(n *node) []*node {
	p.neighbours = p.neighbours[0:0] // reset while preserving memory.
	x, y := n.x, n.y
	var xplus, xminus, yplus, yminus bool // horizontal/vertical grid spots.
	if xplus = p.fp.IsOpen(x+1, y); xplus {
		p.addNeighbour(n.x+1, n.y, n.pathCost+p.orthMoveCost)
	}
	if xminus = p.fp.IsOpen(x-1, y); xminus {
		p.addNeighbour(n.x-1, n.y, n.pathCost+p.orthMoveCost)
	}
	if yplus = p.fp.IsOpen(x, y+1); yplus {
		p.addNeighbour(n.x, n.y+1, n.pathCost+p.orthMoveCost)
	}
	if yminus = p.fp.IsOpen(x, y-1); yminus {
		p.addNeighbour(n.x, n.y-1, n.pathCost+p.orthMoveCost)
	}
	if xminus && yminus && p.fp.IsOpen(x-1, y-1) { // diagonal: xminus, yminus must be passable.
		p.addNeighbour(n.x-1, n.y-1, n.pathCost+p.diagMoveCost)
	}
	if xminus && yplus && p.fp.IsOpen(x-1, y+1) { // diagonal: xminus, yplus must be passable.
		p.addNeighbour(n.x-1, n.y+1, n.pathCost+p.diagMoveCost)
	}
	if xplus && yminus && p.fp.IsOpen(x+1, y-1) { // diagonal: xplus, yminus must be passable.
		p.addNeighbour(n.x+1, n.y-1, n.pathCost+p.diagMoveCost)
	}
	if xplus && yplus && p.fp.IsOpen(x+1, y+1) { // diagonal: xplus, yplus must be passable.
		p.addNeighbour(n.x+1, n.y+1, n.pathCost+p.diagMoveCost)
	}
	return p.neighbours
}

// addNeighbour creates the new neighbour node and adds it to the current neighbour list
// as along as the node is not on the checked list.
func (p *path) addNeighbour(x, y, cost int) {
	id := p.id(x, y)
	if n, ok := p.nodes[id]; ok {
		if n.state == 0 {
			n.pathCost = cost
			p.neighbours = append(p.neighbours, n)
		} else if n.state == isCandidate {
			p.neighbours = append(p.neighbours, n)
		}
	} else {
		n := newNode(id, x, y, cost)
		p.nodes[id] = n
		p.neighbours = append(p.neighbours, n)
	}
}

// traceBack uses recursion to get the from-to path that was discovered.
func (p *path) traceBack(route []int, prev *node) []int {
	if prev.from != -1 {
		n := p.nodes[prev.from]
		route = p.traceBack(route, n)
	}
	route = append(route, prev.x, prev.y)
	return route
}

// =============================================================================

// node keeps distance calculations for a given grid location.
type node struct {
	id       int // unique node identifer based on x, y.
	x, y     int // node location.
	pathCost int // cost to get here from starting point (g)
	projCost int // cost fromStart+aBestGuess to final point (f=g+h).
	from     int // track the previous node in the path.
	state    int // track the candidate and checked state.
}

// newNode creates a new node at the given grid location. The distance values
// are initialized to 0.
func newNode(id, x, y, cost int) *node { return &node{id: id, x: x, y: y, pathCost: cost, from: -1} }

// heuristic returns an appoximation of the distance between the current
// node and the given point.
func (n *node) heuristic(x, y int) int {
	dx, dy := x-n.x, y-n.y
	return dx*dx + dy*dy
}
