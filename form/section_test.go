// Copyright © 2014-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package form

import (
	"testing"
)

// Check getter consistency.
func TestSection(t *testing.T) {
	s := &section{x: 1, y: 2, w: 3, h: 4}
	x, y, w, h := s.Bounds()
	ax, ay := s.At()
	sw, sh := s.Size()
	if x != 1 || x != ax || y != 2 || y != ay ||
		w != 3 || w != sw || h != 4 || h != sh {
		t.Errorf("Invalid section getters %f %f %f %f", x, y, w, h)
	}
}

// Check section bounds.
func TestIn(t *testing.T) {
	f := New([]string{"ab", "cd"}, 100, 100, "gap 10 10").(*form)
	d := f.sects["d"] // x:75 y:25 w:40 h: 40
	if !(d.In(75, 25) && d.In(95, 45) && d.In(55, 5)) {
		t.Errorf("Invalid in bound check %t %t %t", d.In(75, 25), d.In(95, 45), d.In(55, 5))
	}
	if d.In(96, 45) || d.In(55, 4) {
		t.Errorf("Invalid out bound check %t %t", d.In(95, 46), d.In(55, 4))
	}
}
