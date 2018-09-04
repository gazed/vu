// Copyright © 2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package ai

import (
	"fmt"
	"math"
	"testing"
)

// Find a path from one corner of an empty room to the opposite
// diagonal corner.
func TestDiagonalRoute(t *testing.T) {
	graph := &emptyGraph{}
	start, goal := newGridPoint(0, 0), newGridPoint(gridSize-1, gridSize-1)
	path := []Point{} // results
	Find(graph, start, goal, &path)
	if len(path) == 0 {
		t.Errorf("No path found")
	}
	pathEnd := path[len(path)-1]
	if goal.ID() != pathEnd.ID() {
		t.Errorf("Expecting %d got %d", goal.ID(), pathEnd.ID())
	}
	// printPath(graph, path) // Uncomment to see text visualization.
}

// Find a path from one room to another through a series of corridors.
func TestRouteWithRooms(t *testing.T) {
	graph := &roomGraph{}
	start, goal := newGridPoint(18, 7), newGridPoint(18, 16)
	path := []Point{} // results
	Find(graph, start, goal, &path)
	if len(path) == 0 {
		t.Errorf("No path found")
	}
	// printPath(graph, path) // Uncomment to see text visualization.
}

// ============================================================================
// test utility methods and data.

// gridPoint is a valid location on a small grid.
type gridPoint int64 // Depends on global gridSize.
var gridSize = 20    // Bounds of the underlying grid.

// newGridPoint creates a gridPoint for the given grid x,y values.
func newGridPoint(x, y int) gridPoint {
	return gridPoint(x*gridSize + y)
}

// Implement Point interface.
func (p gridPoint) ID() int64 { return int64(p) }
func (p gridPoint) XY() (x, y int) {
	return int(p) / gridSize, int(p) % gridSize
}

// ============================================================================

// emptyGraph is a Graph for a floor plan with no barriers.
type emptyGraph struct{}

func (g *emptyGraph) Neighbours(at Point) (pts []Point) {
	x, y := at.(gridPoint).XY()
	if x >= 0 && x < gridSize && y >= 0 && y < gridSize {
		if x+1 < gridSize {
			pts = append(pts, newGridPoint(x+1, y))
		}
		if x-1 >= 0 {
			pts = append(pts, newGridPoint(x-1, y))
		}
		if y+1 < gridSize {
			pts = append(pts, newGridPoint(x, y+1))
		}
		if y-1 >= 0 {
			pts = append(pts, newGridPoint(x, y-1))
		}
	}
	return pts
}
func (g *emptyGraph) Cost(a, b Point) float64 { return 1.0 }
func (g *emptyGraph) Estimate(a, b Point) float64 {
	ax, ay := a.(gridPoint).XY()
	bx, by := b.(gridPoint).XY()
	return math.Abs(float64(ax-bx)) + math.Abs(float64(ay-by))
}

// ============================================================================

// roomGraph creates a floor plan with 3 rooms connected by corridors.
type roomGraph struct{}

var W = 'W'             // wall
var o = '.'             // open/corridor
var roomMap = [][]rune{ // gridSize*gridSize, 0, 0 is top left.
	{W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W},
	{W, o, o, o, o, o, o, o, W, W, W, W, W, W, W, W, W, W, W, W},
	{W, o, o, o, o, o, o, o, W, W, W, W, W, W, W, W, W, W, W, W},
	{W, o, o, o, o, o, o, o, W, W, W, o, o, o, o, o, o, o, o, W},
	{W, o, o, o, o, o, o, o, W, W, W, o, o, o, o, o, o, o, o, W},
	{W, o, o, o, o, o, o, o, o, o, o, o, o, o, o, o, o, o, o, W},
	{W, o, o, o, o, o, o, o, W, W, W, o, o, o, o, o, o, o, o, W},
	{W, o, o, o, o, o, o, o, W, W, W, o, o, o, o, o, o, o, o, W},
	{W, o, o, o, o, o, o, o, W, W, W, o, o, o, o, o, o, o, o, W},
	{W, W, W, W, o, W, W, W, W, W, W, o, o, o, o, o, o, o, o, W},
	{W, W, W, W, o, W, W, W, W, W, W, W, W, W, W, W, o, W, W, W},
	{W, W, W, W, o, W, W, W, W, W, W, W, W, W, W, W, o, W, W, W},
	{W, o, o, o, o, o, o, o, W, W, W, W, W, W, W, W, o, W, W, W},
	{W, o, o, o, o, o, o, o, W, W, W, W, W, W, W, W, o, W, W, W},
	{W, o, o, o, o, o, o, o, W, W, W, W, W, W, W, W, o, W, W, W},
	{W, o, o, o, o, o, o, o, W, W, W, W, W, W, W, W, o, W, W, W},
	{W, o, o, o, o, o, o, o, W, W, W, W, W, W, W, W, o, W, W, W},
	{W, o, o, o, o, o, o, o, W, W, W, W, W, W, W, W, o, W, W, W},
	{W, o, o, o, o, o, o, o, W, W, W, W, W, W, W, W, o, W, W, W},
	{W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W},
}

// Mock the Graph interface.
// func (r *roomGraph) Size() (width, depth int) { return len(roomMap), len(roomMap[0]) }
func (g *roomGraph) Neighbours(at Point) (pts []Point) {
	roomSize := len(roomMap)
	x, y := at.(gridPoint).XY()
	if x >= 0 && x < roomSize && y >= 0 && y < roomSize {
		if x+1 < roomSize && roomMap[x+1][y] == o {
			pts = append(pts, newGridPoint(x+1, y))
		}
		if x-1 >= 0 && roomMap[x-1][y] == o {
			pts = append(pts, newGridPoint(x-1, y))
		}
		if y+1 < roomSize && roomMap[x][y+1] == o {
			pts = append(pts, newGridPoint(x, y+1))
		}
		if y-1 >= 0 && roomMap[x][y-1] == o {
			pts = append(pts, newGridPoint(x, y-1))
		}
	}
	return pts
}

func (g *roomGraph) Cost(a, b Point) float64 {
	x, y := b.(gridPoint).XY() // cost depends on destination.
	switch {
	case x < 0 || x >= len(roomMap[0]) || y < 0 || y >= len(roomMap):
		return 1000.0 // outside the map.
	case roomMap[x][y] == W:
		return 100.0 // costly to move into a wall.
	default:
		return 1.0
	}
}
func (g *roomGraph) Estimate(a, b Point) float64 {
	ax, ay := a.(gridPoint).XY()
	bx, by := b.(gridPoint).XY()
	return math.Abs(float64(ax-bx)) + math.Abs(float64(ay-by))
}

// ============================================================================

// dump a visual representation of the calculated path.
func printPath(graph Graph, pts []Point) {
	roomMap := make([][]rune, gridSize) // 0, 0 is top left.
	for cnt := range roomMap {
		roomMap[cnt] = make([]rune, gridSize)
	}
	for x := 0; x < gridSize; x++ {
		for y := 0; y < gridSize; y++ {
			cost := graph.Cost(newGridPoint(0, 0), newGridPoint(x, y))
			roomMap[x][y] = '.'
			if cost > 1.0 {
				roomMap[x][y] = 'X'
			}
		}
	}

	// mark the route on the map.
	for _, pt := range pts {
		x, y := pt.(gridPoint).XY()
		roomMap[x][y] = '*'
	}

	// print the map.
	for cnty := range roomMap {
		for _, sym := range roomMap[cnty] {
			fmt.Printf("%c ", sym)
		}
		fmt.Printf("\n")
	}
}

// test utility methods and data.
// ============================================================================
// benchmarking.

// Check A* path finding efficiency.
// Run "go test -bench AStar" to get something like:
// BenchmarkAStar-8   	   10000	    115709 ns/op
func BenchmarkAStar(b *testing.B) {
	graph := &roomGraph{}
	start, goal := newGridPoint(18, 7), newGridPoint(18, 16)
	path := []Point{} // results
	for cnt := 0; cnt < b.N; cnt++ {
		Find(graph, start, goal, &path)
	}
}
