// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package grid

import "testing"

func TestBadNew(t *testing.T) {
	g := New(42)
	if g != nil {
		t.Error("Should be nil with invalid grid type")
	}
}

func TestGridGenerate0(t *testing.T) {
	g := &grid{}
	g.create(0, 0, allWalls)
	w, h := g.Size()
	if w != 7 || h != 7 {
		t.Error("Could not create maze")
	}
}

func TestGridGenerate(t *testing.T) {
	g := &grid{}
	g.create(10, 20, allWalls)
	w, h := g.Size()
	if w != 11 || h != 21 {
		t.Error("Could not create maze")
	}
}

func TestGridSize(t *testing.T) {
	g := &grid{}
	g.create(10, 20, allWalls)
	w, h := g.Size()
	if w != 11 || h != 21 {
		t.Error("Not the right size")
	}
}

func TestGridBand(t *testing.T) {
	g := &grid{}
	g.create(21, 21, allWalls)
	if g.Band(0, 0) != 0 || g.Band(21, 21) != 0 {
		t.Error("Incorrect band 0")
	}
	if g.Band(1, 1) != 1 || g.Band(20, 18) != 1 {
		println(g.Band(1, 1), g.Band(19, 18))
		t.Error("Incorrect band 1")
	}
	if g.Band(10, 10) != 10 {
		t.Error("Incorrect center band")
	}
}

func TestGridCells(t *testing.T) {
	g := &grid{}
	g.create(10, 20, allWalls)
	cells := g.cellSlice()
	if len(cells) != 11*21 {
		t.Error("Not the right amount of cells")
	}
}

func TestGridNorth(t *testing.T) {
	g := &grid{}
	g.create(10, 20, allWalls)
	if g.north(g.cells[1][20]) != nil {
		t.Error("Invalid north should return nil")
	}
	if g.north(g.cells[0][19]).y != 20 {
		t.Error("North should return valid number")
	}
}
func TestGridWest(t *testing.T) {
	g := &grid{}
	g.create(7, 7, allWalls)
	if g.west(g.cells[2][4]).x != 1 {
		t.Error("Invalid west-x should valid 1")
	}
	if g.west(g.cells[2][4]).y != 4 {
		t.Error("Invalid west-y should return 4")
	}
}

func TestNeighbours(t *testing.T) {
	g := &grid{}
	g.create(0, 0, allFloors)
	w, h := g.Size()
	if w != 7 || h != 7 {
		t.Error("Could not create maze")
	}
	floors := g.neighbours(g.cells[0][0], allFloors)
	if len(floors) != 2 {
		t.Error("0,0 should have two floors.")
	}
	floors = g.neighbours(g.cells[6][6], allFloors)
	if len(floors) != 2 {
		t.Error("6,6 should have two floors.")
	}
	floors = g.neighbours(g.cells[5][6], allFloors)
	if len(floors) != 3 {
		t.Error("5,6 should have three floors.")
	}
	floors = g.neighbours(g.cells[5][5], allFloors)
	if len(floors) != 4 {
		t.Error("5,5 should have four floors.")
	}

	// check once for no neighbouring walls
	listWalls := g.neighbours(g.cells[0][0], allWalls)
	if len(listWalls) != 0 {
		t.Error("0,0 should have 0 walls.")
	}
}
