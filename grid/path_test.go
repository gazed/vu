// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

package grid

import (
	"testing"
)

// go test -run Path
func TestPathing(t *testing.T) {

	// should look similar to "dump grid" with an extra border.
	// Note that the offset grid is drawn with the OffsetSymbols.
	t.Run("four-way path", func(t *testing.T) {
		g := Grid{
			{1, 1, 0, 1, 1},
			{1, 0, 0, 0, 1},
			{0, 0, 1, 0, 0},
			{1, 0, 0, 0, 1},
			{1, 1, 0, 1, 1},
		}
		start, goal := Pos{2, 0}, Pos{2, 4}
		path := []Pos{}
		GetPath(g, start, goal, FourWay, &path)
		if len(path) != 7 {
			t.Errorf("expected path length 7 got %d", len(path))
		}
	})

	t.Run("eight-way path", func(t *testing.T) {
		g := Grid{
			{1, 1, 0, 1, 1},
			{1, 0, 0, 0, 1},
			{0, 0, 1, 0, 0},
			{1, 0, 0, 0, 1},
			{1, 1, 0, 1, 1},
		}
		start, goal := Pos{2, 0}, Pos{2, 4}
		path := []Pos{}
		GetPath(g, start, goal, EightWay, &path)
		if len(path) != 5 {
			t.Errorf("expected path length 5 got %d", len(path))
		}
	})
}

// ============================================================================
// Run "go test -bench AStar" to get something like:
//   - BenchmarkAStarSmallEightWay-16            527217              2130 ns/op
//   - BenchmarkAStarMediumFourWay-16             17601             68335 ns/op
//   - BenchmarkAStarMediumEightWay-16            12343             97505 ns/op
//
// Benchmarking takeaways:
//   - larger grids are slower - so reduce grid size or keep searches localized,
//     ie: not one corner of the grid to another.
//   - FourWay is faster than EightWay.
func BenchmarkAStarSmallEightWay(b *testing.B) {
	start, goal := Pos{2, 0}, Pos{2, 4}
	path := []Pos{}
	for b.Loop() {
		GetPath(smallGrid, start, goal, EightWay, &path)
	}
}
func BenchmarkAStarMediumFourWay(b *testing.B) {
	rows, cols := uint(len(mediumGrid)), uint(len(mediumGrid[0]))
	start, goal := Pos{0, 0}, Pos{rows - 1, cols - 1}
	path := []Pos{}
	for b.Loop() {
		GetPath(mediumGrid, start, goal, FourWay, &path)
	}
}
func BenchmarkAStarMediumEightWay(b *testing.B) {
	rows, cols := uint(len(mediumGrid)), uint(len(mediumGrid[0]))
	start, goal := Pos{0, 0}, Pos{rows - 1, cols - 1}
	path := []Pos{}
	for b.Loop() {
		GetPath(mediumGrid, start, goal, EightWay, &path)
	}
}

var smallGrid = Grid{
	{1, 1, 0, 1, 1},
	{1, 0, 0, 0, 1},
	{0, 0, 1, 0, 0},
	{1, 0, 0, 0, 1},
	{1, 1, 0, 1, 1},
}
var mediumGrid = Grid{
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
}
