// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

package grid

import (
	"fmt"
	"testing"
)

// go test -run Flow
func TestFlow(t *testing.T) {
	t.Run("flow empty", func(t *testing.T) {
		gridSize := uint(20)
		f := NewFlow(MakeGrid(20, 20), gridSize/2, gridSize/2)
		bl := f.costmap[0][0]
		br := f.costmap[gridSize-1][0]
		tl := f.costmap[0][gridSize-1]
		tr := f.costmap[gridSize-1][gridSize-1]
		if bl != 20 || br != 19 || tl != 19 || tr != 18 {
			t.Errorf("Invalid costmap %d %d %d %d", bl, br, tl, tr)
			dumpCosts(f)
		}

		tl = f.flowmap[0][0]
		bl = f.flowmap[gridSize-1][0]
		tr = f.flowmap[0][gridSize-1]
		br = f.flowmap[gridSize-1][gridSize-1]
		if bl != ne || br != nw || tl != se || tr != sw {
			t.Errorf("Invalid flowmap %d %d %d %d", bl, br, tl, tr)
			f.Dump()
		}

		// uncomment to see visual representations of the grid and flow maps.
		// dumpCosts(f) // dump cost numbers
		// f.Dump()     // dump flow directions.
	})

	t.Run("flow around center block", func(t *testing.T) {
		gridSize := uint(20)
		grid := MakeGrid(20, 20)
		for r := range 10 {
			for c := range 10 {
				grid[r+5][c+5] = SOLID
			}
		}
		f := NewFlow(grid, gridSize-1, gridSize-1)

		yx := f.flowmap[gridSize-2][gridSize-2]
		x0 := f.flowmap[gridSize-1][gridSize-2]
		y0 := f.flowmap[gridSize-2][gridSize-1]
		if yx != se || x0 != east || y0 != south {
			t.Errorf("Invalid flowmap %d %d %d", yx, x0, y0)
			f.Dump()
		}

		// uncomment to see visual representations of the grid and flow maps.
		// dumpCosts(f) // dump cost numbers
		// f.Dump()     // dump flow directions.
	})
}

// ============================================================================
// utility methods

// dumpCosts dumps the gridmap where 0,0 is the bottom left corner and
// size, size is the top right corner.
func dumpCosts(f *Flow) {
	for y := range f.rows {
		for x := range f.cols {
			fmt.Printf("%3d ", f.costmap[y][x])
		}
		fmt.Printf("\n")
	}
}

// ============================================================================
// benchmarking.

// Check flow field creation efficiency.
// Run "go test -bench Flow" to get something like:
// BenchmarkFlowEmpty7-16            536806              2164 ns/op
// BenchmarkFlowEmpty63-16             4651            259793 ns/op
func BenchmarkFlowEmpty7(b *testing.B) {
	g := MakeGrid(7, 7)
	for b.Loop() {
		NewFlow(g, 3, 3)
	}
}
func BenchmarkFlowEmpty63(b *testing.B) {
	g := MakeGrid(63, 63)
	for b.Loop() {
		NewFlow(g, 23, 23)
	}
}
