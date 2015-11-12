// Copyright © 2014-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package physics

import (
	"testing"

	"github.com/gazed/vu/math/lin"
)

func TestCastRayPlane(t *testing.T) {
	r := newBody(NewRay(0, 0.70710678, 0.70710678)) // ray at origin pointing down +Y +Z
	p := newBody(NewPlane(0, 0, 1))                 // normal +Z
	p.World().Loc.SetS(0, 0, 20)                    // move plane 20 +Z
	hit, x, y, z := castRayPlane(r, p)
	cx, cy, cz := 0.0, 20.0, 20.0 // expected contact location.
	if !hit || !lin.Aeq(x, cx) || !lin.Aeq(y, cy) || !lin.Aeq(z, cz) {
		t.Errorf("%t Expected ray-plane hit at %f %f %f, got %f %f %f", hit, cx, cy, cz, x, y, z)
	}
}

func TestCastRotatedRayPlane(t *testing.T) {
	r := newBody(NewRay(0, 0.70710678, -0.70710678)) // ray at origin pointing down +Y -Z
	r.World().Loc.SetS(0, 0, 20)                     // move ray origin +20 on Z axis.
	p := newBody(NewPlane(0, 0, -1))                 // plane at origin with normal -Z
	hit, x, y, z := castRayPlane(r, p)
	cx, cy, cz := 0.0, 20.0, 0.0 // expected contact location.
	if !hit || !lin.Aeq(x, cx) || !lin.Aeq(y, cy) || !lin.Aeq(z, cz) {
		t.Errorf("%t Expected ray-plane hit at %f %f %f, got %f %f %f", hit, cx, cy, cz, x, y, z)
	}
}

func TestCastRaySphere(t *testing.T) {
	r := newBody(NewRay(0.70710678, 0.70710678, 0.70710678)) // 45 degrees from each axis.
	s := newBody(NewSphere(1))                               // sphere of radius 1
	s.World().Loc.SetS(20, 20, 20)
	hit, x, y, z := castRaySphere(r, s)
	cx, cy, cz := 19.4226497, 19.4226497, 19.4226497 // expected contact location.
	if !hit || !lin.Aeq(x, cx) || !lin.Aeq(y, cy) || !lin.Aeq(z, cz) {
		t.Errorf("%t Expected ray-plane hit at %2.7f %2.7f %2.7f, got %2.7f %2.7f %2.7f", hit, cx, cy, cz, x, y, z)
	}
}

func TestCastRotatedRaySphere(t *testing.T) {
	r := newBody(NewRay(0, 0.70710678, -0.70710678)) // ray at origin pointing down +Y -Z
	r.World().Loc.SetS(0, 0, 20)                     // move ray origin +20 on Z axis.
	s := newBody(NewSphere(1))                       // sphere of radius 1.
	s.World().Loc.SetS(0, 20, 0)                     // put sphere up the y-axis.
	hit, x, y, z := castRaySphere(r, s)
	cx, cy, cz := 0.0, 19.2928932, 0.7071068 // expected contact location.
	if !hit || !lin.Aeq(x, cx) || !lin.Aeq(y, cy) || !lin.Aeq(z, cz) {
		t.Errorf("%t Expected ray-plane hit at %2.7f %2.7f %2.7f, got %2.7f %2.7f %2.7f", hit, cx, cy, cz, x, y, z)
	}
}
