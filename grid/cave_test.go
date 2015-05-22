// Copyright © 2014-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package grid

import "testing"

// Used to view level while tweaking algorithm.
func TestCaveGenerate(t *testing.T) {
	c := &cave{}
	c.Generate(80, 40)
	w, h := c.Size()
	if w != 81 || h != 41 {
		t.Error("Could not create grid")
	}
	// c.dump() // view level.
}
