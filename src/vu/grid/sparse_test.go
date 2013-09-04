// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package grid

import "testing"

func TestSparseLevel(t *testing.T) {
	g := &sparse{}
	g.Generate(10, 20)
	w, h := g.Size()
	if w != 11 || h != 21 {
		t.Error("Could not create grid")
	}
	//g.dump() // view level.
}

func TestSparsePassages(t *testing.T) {
	g := &sparse{}
	g.create(7, 7, allPassages)
	passages := g.neighbours(g.cells[0][0], allPassages)
	if len(passages) != 2 {
		t.Error("0,0 should have two passages.")
	}
	passages = g.neighbours(g.cells[6][6], allPassages)
	if len(passages) != 2 {
		t.Error("6,6 should have two passages.")
	}
	passages = g.neighbours(g.cells[5][6], allPassages)
	if len(passages) != 3 {
		t.Error("5,6 should have three passages.")
	}
	passages = g.neighbours(g.cells[5][5], allPassages)
	if len(passages) != 4 {
		t.Error("5,5 should have four passages.")
	}
	g.cells[5][6].isWall = true
	g.cells[5][4].isWall = true
	g.cells[4][5].isWall = true
	g.cells[6][5].isWall = true
	passages = g.neighbours(g.cells[5][5], allPassages)
	if len(passages) != 0 {
		t.Error("5,5 should now have zero passages.")
	}
}
