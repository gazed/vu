// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package grid

import "testing"

func TestCaveGenerate(t *testing.T) {
	c := &cave{}
	c.Generate(80, 40)
	w, h := c.Size()
	if w != 81 || h != 41 {
		t.Error("Could not create grid")
	}
	// c.dump() // view level.
}
