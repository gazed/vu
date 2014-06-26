// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package grid

import "testing"

func TestLargeRoomsLevel(t *testing.T) {
	g := &rooms{}
	g.Generate(20, 20)
	w, h := g.Size()
	if w != 21 || h != 21 {
		t.Error("Could not create grid")
	}
	// g.dump() // view level.
}

// check splits with the default 4 as the split minimum size.
func TestRoomsSplitSpots(t *testing.T) {
	g := &rooms{grid{}, 4, 7}
	rm := room{0, 0, 9, 7}
	spots := g.splitSpots(&rm, topBottom)
	if spots != 0 {
		t.Error("Only one spot for a 9x7 top/bottom split")
	}
	rm = room{7, 7, 9, 7}
	spots = g.splitSpots(&rm, leftRight)
	if spots != 2 {
		t.Error("Three possibilities for a 9x7 left/right split")
	}
	rm = room{0, 0, 21, 10}
	spots = g.splitSpots(&rm, leftRight)
	if spots != 14 {
		t.Error("Fourteen spots for a 21x21")
	}
}

func TestRoomsSplitRoom(t *testing.T) {
	g := &rooms{grid{}, 4, 7}
	rm := room{0, 0, 7, 7}
	rm1, rm2 := g.splitRoom(&rm, 3, topBottom)
	if rm1.x != 0 || rm1.y != 0 || rm2.x != 0 || rm2.y != 3 {
		t.Error("Improper top/bottom split position")
	}
	if rm1.w != 7 || rm1.h != 4 || rm2.w != 7 || rm2.h != 4 {
		t.Error("Improper top/bottom split size")
	}
	rm = room{7, 7, 21, 21}
	rm1, rm2 = g.splitRoom(&rm, 12, leftRight)
	if rm1.x != 7 || rm1.y != 7 || rm2.x != 19 || rm2.y != 7 {
		t.Error("Improper left/right split position")
	}
	if rm1.w != 13 || rm1.h != 21 || rm2.w != 9 || rm2.h != 21 {
		t.Error("Improper left/right split size")
	}
}
