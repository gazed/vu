// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package grid

import "testing"

func TestDungeonGenerate(t *testing.T) {
	d := &dungeon{}
	d.Generate(80, 40)
	w, h := d.Size()
	if w != 81 || h != 41 {
		t.Error("Could not create dungeon")
	}
	// d.dump() // view level.
}
