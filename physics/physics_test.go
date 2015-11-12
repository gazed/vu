// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package physics

import (
	"fmt"
	"testing"

	"github.com/gazed/vu/math/lin"
)

// Check that broadphase doesn't duplicate comparisons.
// Each compare should cause an overlap.
func TestBroadphaseUniqueCompare(t *testing.T) {
	px, sp := newPhysics(), NewSphere(1)
	bodies := []Body{newBody(sp), newBody(sp), newBody(sp), newBody(sp), newBody(sp)}
	for _, bod := range bodies {
		bod.SetMaterial(1, 0)
	}
	px.broadphase(bodies, px.overlapped)
	if len(px.overlapped) != 10 {
		t.Errorf("Should be 10 unique comparisons for a list of 5. Got %d", len(px.overlapped))
	}
}

// Basic test to check that a sphere will end up above a slab.
// The test uses no restitution (bounciness).
func TestSphereAt(t *testing.T) {
	px := newPhysics()
	slab := newBody(NewBox(100, 25, 100)).SetMaterial(0, 0)
	slab.World().Loc.SetS(0, -25, 0)                // slab below ball at world y==0.
	ball := newBody(NewSphere(1)).SetMaterial(1, 0) //
	ball.World().Loc.SetS(-5, 15, -3)               // ball above slab.
	bodies := []Body{slab, ball}
	for cnt := 0; cnt < 100; cnt++ {
		px.Step(bodies, 0.02)
	}
	ballAt, want := dumpV3(ball.World().Loc), dumpV3(&lin.V3{X: -5, Y: 1, Z: -3})
	if ballAt != want {
		t.Errorf("Ball should be at %s, but its at %s", want, ballAt)
	}
}

// Check that basic collision works independent of general collision resolution.
func TestCollide(t *testing.T) {
	px := newPhysics()
	s0 := newBody(NewSphere(1)).SetMaterial(1, 0)
	s1 := newBody(NewSphere(1)).SetMaterial(1, 0)
	s0.World().Loc.SetS(0, 0, 0)
	s1.World().Loc.SetS(1, 1, 1)
	if !px.Collide(s0, s1) {
		t.Errorf("Expected collision did not happen")
	}
	s0.World().Loc.SetS(-1, -1, -1)
	if px.Collide(s0, s1) {
		t.Errorf("Unexpected collision")
	}
}

// Testing
// ============================================================================
// Utility functions for all package testcases.

func dumpT(t *lin.T) string   { return dumpV3(t.Loc) + dumpQ(t.Rot) }
func dumpQ(q *lin.Q) string   { return fmt.Sprintf("%2.1f", *q) }
func dumpV3(v *lin.V3) string { return fmt.Sprintf("%2.1f", *v) }
func dumpM3(m *lin.M3) string {
	format := "[%+2.1f, %+2.1f, %+2.1f]\n"
	str := fmt.Sprintf(format, m.Xx, m.Xy, m.Xz)
	str += fmt.Sprintf(format, m.Yx, m.Yy, m.Yz)
	str += fmt.Sprintf(format, m.Zx, m.Zy, m.Zz)
	return str
}
