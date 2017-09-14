// Copyright © 2014-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package grid

import (
	"fmt"
	"testing"
)

// Find a path from one corner of an empty room to the opposite
// diagonal corner.
func TestDiagonalPath(t *testing.T) {
	p, dest := newPath(&emptyPlan{}), gridSize-1
	pts := p.Find(0, 0, dest, dest)
	if len(pts) == 0 {
		t.Errorf("No path found")
	}
	dx, dy := pts[len(pts)-2], pts[len(pts)-1]
	if dx != dest || dy != dest {
		t.Errorf("Expecting %d,%d , got %d,%d", dest, dest, dx, dy)
	}
}

// Find a path around a central block in an otherwise empty room.
func TestBlockedPath(t *testing.T) {
	p, dest := newPath(&blockedPlan{}), gridSize-1
	pts := p.Find(0, 0, dest, dest)
	if len(pts) == 0 {
		t.Errorf("No path found")
	}
	dx, dy := pts[len(pts)-2], pts[len(pts)-1]
	if dx != dest || dy != dest {
		t.Errorf("Expecting %d,%d , got %d,%d", dest, dest, dx, dy)
	}
}

// Find a path from one room through a series of corridors to
// another room.
func TestPathWithRooms(t *testing.T) {
	p, ex, ey := newPath(&roomPlan{}), 24, 10 // 25x25 grid.
	pts := p.Find(12, 22, ex, ey)
	if len(pts) == 0 {
		t.Errorf("No path found")
	}
	dx, dy := pts[len(pts)-2], pts[len(pts)-1]
	if dx != ex || dy != ey {
		t.Errorf("Expecting %d,%d , got %d,%d", ex, ey, dx, dy)
	}

	// Uncomment to see a visual representation of the path solution.
	// printPath(pts)
}

// gridSize can be increased to stress test the benchmark.
var gridSize = 20

// dump a visual representation of the calculated path.
func printPath(pts []int) {

	// mark the path.
	for cnt := 0; cnt < len(pts); cnt += 2 {
		x, y := pts[cnt], pts[cnt+1]
		roomMap[x][y] = '-'
	}

	// mark the starting and ending points.
	roomMap[12][22] = 'F' // Y, X flipped for printout.
	roomMap[24][10] = 'T' // Y, X flipped for printout.

	// print the map.
	for cnty := range roomMap {
		for _, sym := range roomMap[cnty] {
			fmt.Printf("%c", sym)
		}
		fmt.Printf("\n")
	}
}

// ============================================================================

// emptyPlan creates a floor plan with no barriers.
type emptyPlan struct{}

func (ep *emptyPlan) Size() (width, depth int) { return gridSize, gridSize }
func (ep *emptyPlan) IsOpen(x, y int) bool {
	return x >= 0 && x < gridSize && y >= 0 && y < gridSize
}

// ============================================================================

// blocked plan creates a floor plan with a large block in the center.
type blockedPlan struct{}

func (bp *blockedPlan) Size() (width, depth int) { return gridSize, gridSize }
func (bp *blockedPlan) IsOpen(x, y int) bool {
	blockSize := 4
	center := gridSize / 2
	bot, top := center-blockSize, center+blockSize
	switch {
	case x >= bot && x <= top && y >= bot && y <= top:
		return false
	default:
		return x >= 0 && x < gridSize && y >= 0 && y < gridSize
	}
}

// ============================================================================

// rooms plan creates a floor plan with 3 rooms connected by corridors.
type roomPlan struct{}

var W = 'W'             // wall
var o = 'o'             // open/corridor
var roomMap = [][]rune{ // 25x25 grid, 0, 0 is top left.
	{W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W},
	{W, o, o, o, o, o, o, o, o, o, o, W, W, W, W, W, W, W, W, W, W, W, W, W, W},
	{W, o, o, o, o, o, o, o, o, o, o, W, W, W, W, W, W, W, W, W, W, W, W, W, W},
	{W, o, o, o, o, o, o, o, o, o, o, W, W, W, o, o, o, o, o, o, o, o, o, o, W},
	{W, o, o, o, o, o, o, o, o, o, o, W, W, W, o, o, o, o, o, o, o, o, o, o, W},
	{W, o, o, o, o, o, o, o, o, o, o, o, o, o, o, o, o, o, o, o, o, o, o, o, W},
	{W, o, o, o, o, o, o, o, o, o, o, W, W, W, o, o, o, o, o, o, o, o, o, o, W},
	{W, o, o, o, o, o, o, o, o, o, o, W, W, W, o, o, o, o, o, o, o, o, o, o, W},
	{W, o, o, o, o, o, o, o, o, o, o, W, W, W, o, o, o, o, o, o, o, o, o, o, W},
	{W, o, o, o, o, o, o, o, o, o, o, W, W, W, o, o, o, o, o, o, o, o, o, o, W},
	{W, o, o, o, o, o, o, o, o, o, o, W, W, W, o, o, o, o, o, o, o, o, o, o, W},
	{W, W, W, W, o, W, W, W, W, W, W, W, W, W, o, o, o, o, o, o, o, o, o, o, W},
	{W, W, W, W, o, W, W, W, W, W, W, W, W, W, o, o, o, o, o, o, o, o, o, o, W},
	{W, W, W, W, o, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, o, W, W, W},
	{W, W, W, W, o, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, o, W, W, W},
	{W, o, o, o, o, o, o, o, o, o, o, W, W, W, W, W, W, W, W, W, W, o, W, W, W},
	{W, o, o, o, o, o, o, o, o, o, o, W, W, W, W, W, W, W, W, W, W, o, W, W, W},
	{W, o, o, o, o, o, o, o, o, o, o, W, W, W, W, W, W, W, W, W, W, o, W, W, W},
	{W, o, o, o, o, o, o, o, o, o, o, W, W, W, W, W, W, W, W, W, W, o, W, W, W},
	{W, o, o, o, o, o, o, o, o, o, o, W, W, W, W, W, W, W, W, W, W, o, W, W, W},
	{W, o, o, o, o, o, o, o, o, o, o, W, W, W, W, W, W, W, W, W, W, o, W, W, W},
	{W, o, o, o, o, o, o, o, o, o, o, W, W, W, W, W, W, W, W, W, W, o, W, W, W},
	{W, o, o, o, o, o, o, o, o, o, o, W, W, W, W, W, W, W, W, W, W, o, W, W, W},
	{W, o, o, o, o, o, o, o, o, o, o, W, W, W, W, W, W, W, W, W, W, W, W, W, W},
	{W, o, o, o, o, o, o, o, o, o, o, W, W, W, W, W, W, W, W, W, W, W, W, W, W},
	{W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W, W},
}

// Mock the Plan interface.
func (rp *roomPlan) Size() (width, depth int) { return len(roomMap), len(roomMap[0]) }
func (rp *roomPlan) IsOpen(x, y int) bool {
	switch {
	case x < 0 || x >= len(roomMap[0]) || y < 0 || y >= len(roomMap):
		return false
	case roomMap[x][y] == W:
		return false
	default:
		return true
	}
}

// unit tests
// ============================================================================
// benchmarking.

// Check A* path finding efficiency.
// Run "go test -bench ." to get something like:
//     BenchmarkAStar	  200000	     11593 ns/op

func BenchmarkAStar(b *testing.B) {
	p := newPath(&emptyPlan{})
	for cnt := 0; cnt < b.N; cnt++ {
		p.Find(0, 0, gridSize-1, gridSize-1)
	}
}
