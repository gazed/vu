// Copyright © 2014-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package ai

import (
	"fmt"
	"testing"
)

func TestFlowId(t *testing.T) {
	f := newGridFlow(&emptyGrid{})
	id := f.id(5, 7)
	if x, y := f.at(id); x != 5 || y != 7 {
		t.Errorf("Invalid id conversion. Expecting, 5, 7. Got %d %d", x, y)
	}
}

func TestFlowEmpty(t *testing.T) {
	f := newGridFlow(&emptyGrid{})
	gx, gy := gridSize/2, gridSize/2
	f.Create(gx, gy)

	bl := f.goalmap[0][0]
	br := f.goalmap[gridSize-1][0]
	tl := f.goalmap[0][gridSize-1]
	tr := f.goalmap[gridSize-1][gridSize-1]
	if bl != 20 || br != 19 || tl != 19 || tr != 18 {
		t.Errorf("Invalid gridmap %d %d %d %d", bl, br, tl, tr)
		printGridmap(f)
	}

	bl = f.flowmap[0][0]
	br = f.flowmap[gridSize-1][0]
	tl = f.flowmap[0][gridSize-1]
	tr = f.flowmap[gridSize-1][gridSize-1]
	if bl != ne || br != nw || tl != se || tr != sw {
		t.Errorf("Invalid flowmap %d %d %d %d", bl, br, tl, tr)
		printFlowmap(f)
	}

	// uncomment to see visual representations of the grid and flow maps.
	// printGridmap(f)
	// printFlowmap(f)
}

func TestFlowBlock(t *testing.T) {
	f := newGridFlow(&blockedGrid{})
	f.Create(0, 0)

	bl := f.goalmap[0][0]
	br := f.goalmap[gridSize-1][0]
	tl := f.goalmap[0][gridSize-1]
	tr := f.goalmap[gridSize-1][gridSize-1]
	if bl != 0 || br != 19 || tl != 19 || tr != 38 {
		t.Errorf("Invalid gridmap %d %d %d %d", bl, br, tl, tr)
		printGridmap(f)
	}

	bl = f.flowmap[0][0]
	br = f.flowmap[gridSize-1][0]
	tl = f.flowmap[0][gridSize-1]
	tr = f.flowmap[gridSize-1][gridSize-1]
	if bl != goal || br != west || tl != south || tr != sw {
		t.Errorf("Invalid flowmap %d %d %d %d", bl, br, tl, tr)
		printFlowmap(f)
	}

	// uncomment to see text representations of the grid and flow maps.
	// printGridmap(f)
	// printFlowmap(f)
}

func TestFlowNext(t *testing.T) {
	f := newGridFlow(&blockedGrid{})
	f.Create(0, 0)
	dx, dy := f.Next(gridSize-1, gridSize-1)
	if dx != -1 || dy != -1 {
		t.Errorf("Invalid next direction %d %d", dx, dy)
	}
}

// unit tests.
// ============================================================================
// utility methods

// printGridmap dumps the gridmap where 0,0 is the bottom left corner and
// size, size is the top right corner.
func printGridmap(f *gridFlow) {
	for y := f.ysz - 1; y >= 0; y-- {
		for x := 0; x < f.xsz; x++ {
			fmt.Printf("%3d ", f.goalmap[x][y])
		}
		fmt.Printf("\n")
	}
}

// printFlowmap dumps the flowmap where 0,0 is the bottom left corner and
// size, size is the top right corner.
func printFlowmap(f *gridFlow) {
	for y := f.ysz - 1; y >= 0; y-- {
		for x := 0; x < f.xsz; x++ {
			dir := f.flowmap[x][y]
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
			case f.max:
				print("X") // marks impassable areas:
			default:
				print("o") // marks the goal.
			}
		}
		println()
	}
}

// emptyGrid creates a floor plan with no barriers.
type emptyGrid struct{}

func (ep *emptyGrid) Size() (width, depth int) { return gridSize, gridSize }
func (ep *emptyGrid) IsOpen(x, y int) bool {
	return x >= 0 && x < gridSize && y >= 0 && y < gridSize
}

// ============================================================================

// blocked plan creates a floor plan with a large block in the center.
type blockedGrid struct{}

func (bp *blockedGrid) Size() (width, depth int) { return gridSize, gridSize }
func (bp *blockedGrid) IsOpen(x, y int) bool {
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

// utility methods
// ============================================================================
// benchmarking.

// Check flow field creation efficiency.
// Run 'go test -bench "Flow*"' to get something like:
// BenchmarkFlowEmpty-8   	   30000	     45274 ns/op
// BenchmarkFlowBlock-8   	   50000	     31162 ns/op

func BenchmarkFlowEmpty(b *testing.B) {
	f := newGridFlow(&emptyGrid{})
	gx, gy := gridSize/2, gridSize/2
	for cnt := 0; cnt < b.N; cnt++ {
		f.Create(gx, gy)
	}
}
func BenchmarkFlowBlock(b *testing.B) {
	f := newGridFlow(&blockedGrid{})
	for cnt := 0; cnt < b.N; cnt++ {
		f.Create(0, 0)
	}
}
