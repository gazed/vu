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
	t.Run("grid Index", func(t *testing.T) {
		g := MakeGrid(3, 3)
		if got := g.Index(0, 0); got != 0 {
			t.Errorf("expected 0 got %d", got)
		}
		if r, c := g.Indices(0); r != 0 && c != 0 {
			t.Errorf("expected 0,0 got %d:%d", r, c)
		}
		if got := g.Index(1, 2); got != 5 {
			t.Errorf("expected 5 got %d", got)
		}
		if r, c := g.Indices(5); r != 1 && c != 2 {
			t.Errorf("expected 1,2 got %d:%d", r, c)
		}
		if got := g.Index(2, 2); got != 8 {
			t.Errorf("expected 8 got %d", got)
		}
		if r, c := g.Indices(8); r != 2 && c != 2 {
			t.Errorf("expected 2,2 got %d:%d", r, c)
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
		g.Dump(GridSymbols)
	})
}
