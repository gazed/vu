// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package grid

import "testing"

// Used to view level while tweaking algorithm.
func TestSparseLevel(t *testing.T) {
	g := &sparse{}
	g.Generate(10, 20)
	w, h := g.Size()
	if w != 11 || h != 21 {
		t.Error("Could not create grid")
	}
	//g.dump() // view level.
}

func TestSparseFloors(t *testing.T) {
	g := &sparse{}
	g.create(7, 7, allFloors)
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
	g.cells[5][6].isWall = true
	g.cells[5][4].isWall = true
	g.cells[4][5].isWall = true
	g.cells[6][5].isWall = true
	floors = g.neighbours(g.cells[5][5], allFloors)
	if len(floors) != 0 {
		t.Error("5,5 should now have zero floors.")
	}
}
