// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package vu

import (
	"fmt"
	"testing"
	"vu/math/lin"
)

// Dictate how errors get printed.
const format = "\ngot\n%s\nwanted\n%s"

func TestPickDir(t *testing.T) {
	view := &lin.M4{
		1.00, +0.00, +0.00, +0.00,
		0.00, +1.00, +0.00, +0.00,
		0.00, +0.00, +1.00, +0.00,
		0.00, -2.00, -2.00, +1.00}
	mv := &lin.M4{
		1.00, +0.00, +0.00, +0.00,
		0.00, +1.00, +0.00, +0.00,
		0.00, +0.00, +1.00, +0.00,
		0.00, +1.50, -4.00, +1.00}
	mv.Mult(view)

	// try the two corners and the center of the screen.
	sx, sy := 0, 0
	w := PickDir(sx, sy, 45, 800, 600, 0.1, 50, mv)
	got := fmt.Sprintf("%+4d %+4d w %+2.5f %+2.5f %+2.5f", sx, sy, w.X, w.Y, w.Z)
	want := "  +0   +0 w -0.45450 -0.34087 -0.82294"
	if got != want {
		t.Errorf(format, got, want)
	}
	sx, sy = 800, 600
	w = PickDir(sx, sy, 45, 800, 600, 0.1, 50, mv)
	got = fmt.Sprintf("%+4d %+4d w %+2.5f %+2.5f %+2.5f", sx, sy, w.X, w.Y, w.Z)
	want = "+800 +600 w +0.45450 +0.34087 -0.82294"
	if got != want {
		t.Errorf(format, got, want)
	}
	sx, sy = 400, 300
	w = PickDir(sx, sy, 45, 800, 600, 0.1, 50, mv)
	got = fmt.Sprintf("%+4d %+4d w %+2.5f %+2.5f %+2.5f", sx, sy, w.X, w.Y, w.Z)
	want = "+400 +300 w +0.00000 +0.00000 -1.00000"
	if got != want {
		t.Errorf(format, got, want)
	}

	o := &lin.V3{0, -2, -2}
	d := &lin.V3{0, 0, -1}
	o.Add(d.Scale(3))
}

// Dump like lin.M3.Dump().
//
//   [ x0 y0 z0 w0 ]
//   [ x1 y1 z1 w1 ]
//   [ x2 y2 z2 w2 ]
//   [ x3 y3 z3 w3 ]
func dump(m *lin.M4) string {
	format := "[%+2.2f, %+2.2f, %+2.2f, %+2.2f]\n"
	str := fmt.Sprintf(format, m.X0, m.Y0, m.Z0, m.W0)
	str += fmt.Sprintf(format, m.X1, m.Y1, m.Z1, m.W1)
	str += fmt.Sprintf(format, m.X2, m.Y2, m.Z2, m.W2)
	str += fmt.Sprintf(format, m.X3, m.Y3, m.Z3, m.W3)
	return str
}
