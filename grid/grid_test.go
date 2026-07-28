// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

package grid

import (
	"reflect"
	"testing"
)

// go test -run Grid
func TestGrid(t *testing.T) {
	t.Run("empty grid", func(t *testing.T) {
		g := MakeGrid(7, 7)
		g.Fill(EMPTY)
		expectedGrid := Grid{
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
		}
		if !reflect.DeepEqual(g, expectedGrid) {
			t.Errorf("expected empty grid")
		}
	})
	t.Run("solid grid", func(t *testing.T) {
		g := MakeGrid(3, 5)
		g.Fill(SOLID)
		expectedGrid := Grid{
			{1, 1, 1, 1, 1},
			{1, 1, 1, 1, 1},
			{1, 1, 1, 1, 1},
		}
		if !reflect.DeepEqual(g, expectedGrid) {
			t.Errorf("expected solid grid")
		}
	})
	t.Run("grid ID", func(t *testing.T) {
		g := MakeGrid(3, 3)
		if got := g.ID(0, 0); got != 0 {
			t.Errorf("expected 0 got %d", got)
		}
		if got := g.ID(1, 2); got != 5 {
			t.Errorf("expected 5 got %d", got)
		}
		if got := g.ID(2, 2); got != 8 {
			t.Errorf("expected 8 got %d", got)
		}
	})
	t.Run("dump grid", func(t *testing.T) {
		g := Grid{
			{1, 1, 0, 1, 1},
			{1, 0, 0, 0, 1},
			{0, 0, 1, 0, 0},
			{1, 0, 0, 0, 1},
			{1, 1, 0, 1, 1},
		}
		gridSymbols := map[uint16]string{0: "  ", 1: "██"}
		g.Dump(gridSymbols)
	})
}
