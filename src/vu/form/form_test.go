// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package form

import (
	"testing"
)

// // Check that a section can be added with no parameters.
func TestForm(t *testing.T) {
	plan := []string{
		"ab",
		"cd",
	}
	f := New(plan, 100, 100).(*form)
	d := f.sects["d"]
	if !(d.x == 75 && d.y == 25 && d.w == 50 && d.h == 50) {
		t.Errorf("Invalid layout %f %f %f %f", d.x, d.y, d.w, d.h)
	}
}

func TestGap(t *testing.T) {
	plan := []string{
		"ab",
		"cd",
	}
	f := New(plan, 100, 100, "gap 10 10").(*form)
	d := f.sects["d"]
	if !(d.x == 75 && d.y == 25 && d.w == 40 && d.h == 40) {
		t.Errorf("Invalid layout %f %f %f %f", d.x, d.y, d.w, d.h)
	}
}

func TestPad(t *testing.T) {
	plan := []string{
		"ab",
		"cd",
	}
	f := New(plan, 100, 100, "pad 10 10 10 10").(*form)
	d := f.sects["d"]
	if !(d.x == 70 && d.y == 30 && d.w == 40 && d.h == 40) {
		t.Errorf("Invalid layout %f %f %f %f", d.x, d.y, d.w, d.h)
	}
}

func TestSpan(t *testing.T) {
	plan := []string{
		"xxa",
		"xxb",
	}
	f := New(plan, 100, 100).(*form)
	x := f.sects["x"]
	if !(int(x.x) == 33 && x.y == 50 && int(x.w) == 66 && x.h == 100) {
		t.Errorf("Invalid layout %f %f %f %f", x.x, x.y, x.w, x.h)
	}
}

func TestGrowCorner(t *testing.T) {
	plan := []string{
		"ab",
		"cd",
	}
	f := New(plan, 100, 100, "grabx 0", "graby 0").(*form)
	f.Resize(200, 200)
	a := f.sects["a"]
	if !(a.x == 75 && a.y == 125 && a.w == 150 && a.h == 150) {
		t.Errorf("Invalid resize %f %f %f %f", a.x, a.y, a.w, a.h)
	}
	d := f.sects["d"]
	if !(d.x == 175 && d.y == 25 && d.w == 50 && d.h == 50) {
		t.Errorf("Invalid resize %f %f %f %f", d.x, d.y, d.w, d.h)
	}
}

func TestGrowCenter(t *testing.T) {
	plan := []string{
		"abc",
		"def",
		"ghi",
	}
	f := New(plan, 300, 300, "grabx 1", "graby 1").(*form)
	f.Resize(600, 600)
	a := f.sects["a"]
	if !(a.x == 50 && a.y == 550 && a.w == 100 && a.h == 100) {
		t.Errorf("Invalid resize %f %f %f %f", a.x, a.y, a.w, a.h)
	}
	e := f.sects["e"]
	if !(e.x == 300 && e.y == 300 && e.w == 400 && e.h == 400) {
		t.Errorf("Invalid resize %f %f %f %f", e.x, e.y, e.w, e.h)
	}
	i := f.sects["i"]
	if !(i.x == 550 && i.y == 50 && i.w == 100 && i.h == 100) {
		t.Errorf("Invalid resize %f %f %f %f", i.x, i.y, i.w, i.h)
	}
}
