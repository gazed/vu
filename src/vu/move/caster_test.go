// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package move

import (
	"testing"
	"vu/math/lin"
)

func TestCastRayPlane(t *testing.T) {
	r := newSolid(NewRay(0, 0.70710678, 0.70710678))
	p := newSolid(NewPlane(0, 0, 1))
	p.World().Loc.SetS(0, 0, 20)
	hit, x, y, z := castRayPlane(r, p)
	cx, cy, cz := 0.0, 20.0, 20.0 // expected contact location.
	if !hit || !lin.Aeq(x, cx) || !lin.Aeq(y, cy) || !lin.Aeq(z, cz) {
		t.Errorf("%t Expected ray-plane hit at %f %f %f, got %f %f %f", hit, cx, cy, cz, x, y, z)
	}
}

func TestCastRaySphere(t *testing.T) {
	r := newSolid(NewRay(0.70710678, 0.70710678, 0.70710678))
	s := newSolid(NewSphere(1))
	s.World().Loc.SetS(20, 20, 20)
	hit, x, y, z := castRaySphere(r, s)
	cx, cy, cz := 19.4226497, 19.4226497, 19.4226497 // expected contact location.
	if !hit || !lin.Aeq(x, cx) || !lin.Aeq(y, cy) || !lin.Aeq(z, cz) {
		t.Errorf("%t Expected ray-plane hit at %2.7f %2.7f %2.7f, got %2.7f %2.7f %2.7f", hit, cx, cy, cz, x, y, z)
	}
}
