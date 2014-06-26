// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package grid

import "testing"

func TestDenseLevel(t *testing.T) {
	g := &dense{}
	g.Generate(10, 20)
	w, h := g.Size()
	if w != 11 || h != 21 {
		t.Error("Could not create grid")
	}
	//g.dump() // view grid.
}
