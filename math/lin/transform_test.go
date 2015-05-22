// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package lin

// Test combinations of rotations, translations, and scales. The standard
// combination, when there is more than one, is first scale, then rotate,
// then translate.

import (
	"fmt"
	"testing"
)

func TestMovementAroundY(t *testing.T) {
	t1 := NewT().SetLoc(5, 0, 0).SetAa(0, 1, 0, Rad(90))
	v, want := &V3{2, 0, 0}, &V3{5, 0, -2}

	// rotates to -Z, and then moves to X:5 giving (5, 0, -2)
	if t1.App(v); !v.Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}

func TestMovementAroundX(t *testing.T) {
	t1 := NewT().SetLoc(5, 0, 0).SetAa(1, 0, 0, Rad(90))
	v, want := &V3{2, 0, 0}, &V3{7, 0, 0}

	// rotate does not affect x values, and then moves to X:5 giving (7, 0, 0)
	if t1.App(v); !v.Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}

func TestMovementAroundZ(t *testing.T) {
	t1 := NewT().SetLoc(5, 0, 0).SetAa(0, 0, 1, Rad(90))
	v, want := &V3{2, 0, 0}, &V3{5, 2, 0}

	// rotates to +Y, and then moves to X:5 giving (5, 2, 0)
	if t1.App(v); !v.Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}

func TestMultTransform(t *testing.T) {
	t1 := NewT().SetLoc(5, 0, 0).SetAa(0, 1, 0, Rad(90)) // move along X, rotate about Y.
	t2 := NewT().SetLoc(5, 0, 0).SetAa(0, 0, 1, Rad(90)) // move along X rotate about Z, (becomes move along Z).
	v, want := &V3{2, 0, 0}, &V3{5, 0, -7}
	if t1.Mult(t1, t2).App(v); !v.Aeq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}

func TestApply(t *testing.T) {
	v, t1, want := &V3{}, NewT().SetLoc(5, 0, 0).SetAa(1, 0, 0, Rad(90)), &V3{6, 0, 0}
	if v.X, v.Y, v.Z = t1.AppS(1, 0, 0); !v.Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
	want = &V3{5, 0, 1} // right hand rule: positive Y to positive Z
	if v.X, v.Y, v.Z = t1.AppS(0, 1, 0); !v.Aeq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
	want = &V3{5, -1, 0} // right hand rule: positive Z turns to -Y
	if v.X, v.Y, v.Z = t1.AppS(0, 0, 1); !v.Aeq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}

func TestApplyInv(t *testing.T) {
	v, t1, want := &V3{}, NewT().SetLoc(5, 0, 0).SetAa(1, 0, 0, Rad(90)), &V3{0, 1, 0}
	if v.X, v.Y, v.Z = t1.InvS(5, 0, 1); !v.Aeq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}

// Rotate a point X=1 90 degrees about the y-axis. This puts it on it on the -z axis
// then translate it along X by 10. It should be at (10, 0 -1)
func TestTransform(t *testing.T) {
	v, transform := NewV3S(1, 0, 0), NewT().SetLoc(10, 0, 0).SetAa(0, 1, 0, Rad(90))
	want := NewV3S(10, 0, -1)
	if transform.App(v); !v.Aeq(want) {
		t.Errorf("Invalid translation: %s", v.Dump())
	}
}

// Ensure the inverse transform puts the point back to where it was.
func TestInverseTransform(t *testing.T) {
	v, transform := NewV3S(1, 0, 0), NewT().SetLoc(10, 0, 0).SetAa(0, 1, 0, Rad(90))
	transform.App(v)
	transform.Inv(v)
	if !Aeq(v.X, 1) || !Aeq(v.Y, 0) || !Aeq(v.Z, 0) {
		t.Errorf("Invalid translation: %s", v.Dump())
	}
}

func TestIntegrateRotateY(t *testing.T) {
	t1, a := NewT(), NewT().SetLoc(0, 0, 0).SetRot(0, 0, 0, 1)
	linv, angv := &V3{0, 0, 0}, &V3{0, 10, 0}
	t1.Integrate(a, linv, angv, 0.02)
	x, y, z, ang := t1.Rot.Aa()
	got := fmt.Sprintf("%f %f %f %f", x, y, z, Deg(ang))
	want := "0.000000 1.000000 0.000000 11.459156"
	if got != want {
		t.Errorf(format, got, want)
	}
}
func TestIntegrateRotateXY(t *testing.T) {
	t1, a := NewT(), NewT().SetLoc(0, 0, 0).SetRot(0, 0, 0, 1)
	linv, angv := &V3{0, 0, 0}, &V3{0.5, 0.5, 0}
	t1.Integrate(a, linv, angv, 0.02)
	x, y, z, ang := t1.Rot.Aa()
	got := fmt.Sprintf("%f %f %f %f", x, y, z, Deg(ang))
	want := "0.707107 0.707107 0.000000 0.810285"
	if got != want {
		t.Errorf(format, got, want)
	}
}

// Test integration using numbers that were pumped through the bullet physics
// simulation integration code.
func TestIntegrateS(t *testing.T) {
	t1, a := NewT(), NewT().SetLoc(-5, 1.388006, -3).SetRot(0.182574, 0.365148, 0.547723, 0.730297)
	linv, angv := &V3{0.516828, -10.105854, 0.000000}, &V3{10.041207, -0.775241, -0.922906}
	t1.Integrate(a, linv, angv, 0.02)
	lx, ly, lz := t1.Loc.GetS()
	rx, ry, rz, rw := t1.Rot.GetS()
	got := fmt.Sprintf("%f %f %f :: %f %f %f %f", lx, ly, lz, rx, ry, rz, rw)
	want := "-4.989663 1.185889 -3.000000 :: 0.253972 0.301044 0.576212 0.716136"
	if got != want {
		t.Errorf(format, got, want)
	}
}

// test applying the transform using AppS and App and numbers
// that were pumped through the bullet physics transform.
func TestApplyBoth(t *testing.T) {
	a := NewT().SetLoc(-5.0, 1.388006, -3.0).SetRot(0.182574, 0.365148, 0.547723, 0.730297)
	want1, want2 := &V3{-4.8, 2.7880069, -1.999998}, &V3{-5.2, -0.0119949, -4.000001}
	v1, v2 := NewV3S(a.AppS(1, 1, 1)), NewV3S(a.AppS(-1, -1, -1))
	if !v1.Aeq(want1) {
		t.Errorf(format, v1.Dump(), want1.Dump())
	}
	if !v2.Aeq(want2) {
		t.Errorf(format, v2.Dump(), want2.Dump())
	}
	v1, v2 = a.App(NewV3S(1, 1, 1)), a.App(NewV3S(-1, -1, -1))
	if !v1.Aeq(want1) {
		t.Errorf(format, v1.Dump(), want1.Dump())
	}
	if !v2.Aeq(want2) {
		t.Errorf(format, v2.Dump(), want2.Dump())
	}
}
