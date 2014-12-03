// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package grid

import (
	"fmt"
	"testing"
)

func TestFlowId(t *testing.T) {
	f := newFlow(&emptyPlan{})
	id := f.id(5, 7)
	if x, y := f.at(id); x != 5 || y != 7 {
		t.Errorf("Invalid id conversion. Expecting, 5, 7. Got %d %d", x, y)
	}
}

func TestFlowEmpty(t *testing.T) {
	f := newFlow(&emptyPlan{})
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
	if bl != NE || br != NW || tl != SE || tr != SW {
		t.Errorf("Invalid flowmap %d %d %d %d", bl, br, tl, tr)
		printFlowmap(f)
	}

	// uncomment to see visual repsentations of the grid and flow maps.
	// printGridmap(f)
	// printFlowmap(f)
}

func TestFlowBlock(t *testing.T) {
	f := newFlow(&blockedPlan{})
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
	if bl != GOAL || br != WEST || tl != SOUTH || tr != SW {
		t.Errorf("Invalid flowmap %d %d %d %d", bl, br, tl, tr)
		printFlowmap(f)
	}

	// uncomment to see visual repsentations of the grid and flow maps.
	// printGridmap(f)
	// printFlowmap(f)
}

func TestFlowNext(t *testing.T) {
	f := newFlow(&blockedPlan{})
	f.Create(0, 0)
	dx, dy := f.Next(float64(gridSize-1), float64(gridSize-1))
	if dx != -0.7071068 || dy != -0.7071068 {
		t.Errorf("Invalid next direction %f %f", dx, dy)
	}
}

// unit tests.
// ============================================================================
// utility methods

// printGridmap dumps the gridmap where 0,0 is the bottom left corner and
// size, size is the top right corner.
func printGridmap(f *flow) {
	for y := f.ysz - 1; y >= 0; y-- {
		for x := 0; x < f.xsz; x++ {
			fmt.Printf("%3d ", f.goalmap[x][y])
		}
		fmt.Printf("\n")
	}
}

// printFlowmap dumps the flowmap where 0,0 is the bottom left corner and
// size, size is the top right corner.
func printFlowmap(f *flow) {
	for y := f.ysz - 1; y >= 0; y-- {
		for x := 0; x < f.xsz; x++ {
			dir := f.flowmap[x][y]
			switch dir {
			case NORTH:
				print("↑")
			case NE:
				print("↗")
			case EAST:
				print("→")
			case SE:
				print("↘")
			case SOUTH:
				print("↓")
			case SW:
				print("↙")
			case WEST:
				print("←")
			case NW:
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

// utility methods
// ============================================================================
// benchmarking.

// Check flow field creation efficiency.
// Run "go test -bench ." to get something like:
// BenchmarkFlowEmpty	   50000	     64254 ns/op
// BenchmarkFlowBlock	   50000	     45827 ns/op

// No walls or barriers is slightly slower than when there are some walls.
func BenchmarkFlowEmpty(b *testing.B) {
	f := newFlow(&emptyPlan{})
	gx, gy := gridSize/2, gridSize/2
	for cnt := 0; cnt < b.N; cnt++ {
		f.Create(gx, gy)
	}
}
func BenchmarkFlowBlock(b *testing.B) {
	f := newFlow(&blockedPlan{})
	for cnt := 0; cnt < b.N; cnt++ {
		f.Create(0, 0)
	}
}
