// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package grid

import "testing"

// Used to view level while tweaking algorithm.
func TestPrimGenerate(t *testing.T) {
	g := &primMaze{}
	g.Generate(10, 20)
	w, h := g.Size()
	if w != 11 || h != 21 {
		t.Error("Could not create grid")
	}
	// g.dump() // view level.
}
