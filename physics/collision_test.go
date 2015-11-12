// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package physics

import (
	"testing"

	"github.com/gazed/vu/math/lin"
)

func TestCollideSphereSphere(t *testing.T) {
	a, b, cons := NewBody(NewSphere(1)), NewBody(NewSphere(1)), newManifold()
	if _, _, cs := collideSphereSphere(a, b, cons); len(cs) != 1 || cs[0].depth == 2 {
		t.Errorf("Identical spheres at the origin should overlap by %f.", cs[0].depth)
	}

	// check each axis.
	a.World().Loc.SetS(2, 0, 0)
	if _, _, cs := collideSphereSphere(a, b, cons); cs[0].depth != 0 ||
		dumpV3(cs[0].point) != "{1.0 0.0 0.0}" || dumpV3(cs[0].normal) != "{1.0 0.0 0.0}" {
		t.Errorf("Spheres touching at point (1,0,0) do not overlap. %s", dumpV3(cs[0].point))
	}
	a.World().Loc.SetS(0, 2, 0)
	if _, _, cs := collideSphereSphere(a, b, cons); cs[0].depth != 0 ||
		dumpV3(cs[0].point) != "{0.0 1.0 0.0}" || dumpV3(cs[0].normal) != "{0.0 1.0 0.0}" {
		t.Errorf("Spheres touching at point (0,1,0) do not overlap. %s", dumpV3(cs[0].point))
	}
	a.World().Loc.SetS(0, 0, 2)
	if _, _, cs := collideSphereSphere(a, b, cons); cs[0].depth != 0 ||
		dumpV3(cs[0].point) != "{0.0 0.0 1.0}" || dumpV3(cs[0].normal) != "{0.0 0.0 1.0}" {
		t.Errorf("Spheres touching at point (0,0,1) do not overlap. %s", dumpV3(cs[0].point))
	}

	// check just outside and slightly overlapping.
	a.World().Loc.SetS(2.01, 0, 0)
	if _, _, cs := collideSphereSphere(a, b, cons); len(cs) != 0 {
		t.Error("Spheres not touching")
	}
	a.World().Loc.SetS(0, 0, 1.75)
	if _, _, cs := collideSphereSphere(a, b, cons); cs[0].depth != -0.25 ||
		dumpV3(cs[0].point) != "{0.0 0.0 1.0}" || dumpV3(cs[0].normal) != "{0.0 0.0 1.0}" {
		t.Errorf("Spheres touching at point (0,0,1) overlaps by %2.2f %s", cs[0].depth, dumpV3(cs[0].point))
	}
}

func TestCollideSphereBox(t *testing.T) {
	a, b, cons := NewBody(NewSphere(1)), NewBody(NewBox(1, 1, 1)), newManifold()
	if _, _, cs := collideSphereBox(a, b, cons); cs[0].depth != -2.04 ||
		dumpV3(cs[0].point) != "{1.0 0.0 0.0}" || dumpV3(cs[0].normal) != "{1.0 0.0 0.0}" {
		t.Errorf("Sphere touching box at point A %f %s %s", cs[0].depth, dumpV3(cs[0].point), dumpV3(cs[0].normal))
	}
	a.World().Loc.SetS(0, 2, 0)
	if _, _, cs := collideSphereBox(a, b, cons); !lin.Aeq(cs[0].depth, -margin) ||
		dumpV3(cs[0].point) != "{0.0 1.0 0.0}" || dumpV3(cs[0].normal) != "{0.0 1.0 0.0}" {
		t.Errorf("Sphere touching box at point %f %s %s", cs[0].depth, dumpV3(cs[0].point), dumpV3(cs[0].normal))
	}
	a.World().Loc.SetS(0, 0, 2.15)
	if _, _, cs := collideSphereBox(a, b, cons); len(cs) != 0 {
		t.Errorf("Sphere not touching box %f %s %s", cs[0].depth, dumpV3(cs[0].point), dumpV3(cs[0].normal))
	}

	// close enough to be considered in contact.
	a.World().Loc.SetS(0, 0, 2.1)
	if _, _, cs := collideSphereBox(a, b, cons); !lin.Aeq(cs[0].depth, 0.06) ||
		dumpV3(cs[0].point) != "{0.0 0.0 1.0}" || dumpV3(cs[0].normal) != "{0.0 0.0 1.0}" {
		t.Errorf("Sphere close to touching box %f %s %s", cs[0].depth, dumpV3(cs[0].point), dumpV3(cs[0].normal))
	}
}

// Tests that the narrowphase collision lookup finds the algorithm that flips
// the box-sphere to be sphere-box.
func TestCollideBoxSphere(t *testing.T) {
	box, sphere, c, cons := newBody(NewBox(1, 1, 1)), newBody(NewSphere(1)), newCollider(), newManifold()
	sphere.World().Loc.SetS(0, 2, 0)
	algorithm := c.algorithms[box.shape.Type()][sphere.shape.Type()]
	i, j, cs := algorithm(box, sphere, cons)
	ii, jj := i.(*body), j.(*body)
	if ii.shape.Type() != SphereShape || jj.shape.Type() != BoxShape {
		t.Error("Should have flipped the objects into Sphere, Box")
	}
	if !lin.Aeq(cs[0].depth, -margin) || dumpV3(cs[0].point) != "{0.0 1.0 0.0}" || dumpV3(cs[0].normal) != "{0.0 1.0 0.0}" {
		t.Errorf("Contact info should be the same %f %s %s", cs[0].depth, dumpV3(cs[0].point), dumpV3(cs[0].normal))
	}
}

func TestCollideBoxBox(t *testing.T) {
	a, b, cons := NewBody(NewBox(0.5, 0.5, 0.5)), NewBody(NewBox(1, 1, 1)), newManifold()
	if _, _, cs := collideBoxBox(a, b, cons); len(cs) == 0 || cs[0].depth != -1.58 ||
		dumpV3(cs[0].point) != "{-1.0 0.5 0.5}" || dumpV3(cs[0].normal) != "{-1.0 -0.0 -0.0}" {
		depth, point, norm := cs[0].depth, dumpV3(cs[0].point), dumpV3(cs[0].normal)
		t.Errorf("Boxes should collide since one is inside the other %f %s %s", depth, point, norm)
	}

	// just inside of contact range.
	a.World().Loc.SetS(0, 0, 1.49)
	if _, _, cs := collideBoxBox(a, b, cons); len(cs) != 4 || !lin.Aeq(cs[0].depth, -0.09) ||
		dumpV3(cs[0].point) != "{0.5 0.5 1.0}" || dumpV3(cs[0].normal) != "{0.0 0.0 1.0}" {
		depth, point, norm := cs[0].depth, dumpV3(cs[0].point), dumpV3(cs[0].normal)
		t.Errorf("Boxes should collide %f %s %s", depth, point, norm)
	}
	a.World().Loc.SetS(0, 1.49, 0)
	if _, _, cs := collideBoxBox(a, b, cons); len(cs) != 4 || !lin.Aeq(cs[0].depth, -0.09) ||
		dumpV3(cs[0].point) != "{0.5 1.0 0.5}" || dumpV3(cs[0].normal) != "{0.0 1.0 0.0}" {
		depth, point, norm := cs[0].depth, dumpV3(cs[0].point), dumpV3(cs[0].normal)
		t.Errorf("Boxes should collide %f %s %s", depth, point, norm)
	}
	a.World().Loc.SetS(1.49, 0, 0)
	if _, _, cs := collideBoxBox(a, b, cons); len(cs) != 4 || !lin.Aeq(cs[0].depth, -0.09) ||
		dumpV3(cs[0].point) != "{1.0 0.5 0.5}" || dumpV3(cs[0].normal) != "{1.0 0.0 0.0}" {
		depth, point, norm := cs[0].depth, dumpV3(cs[0].point), dumpV3(cs[0].normal)
		t.Errorf("Boxes should collide %f %s %s", depth, point, norm)
	}

	// just outside of contact range.
	a.World().Loc.SetS(0, 0, 1.6)
	if _, _, cs := collideBoxBox(a, b, cons); len(cs) != 0 {
		t.Errorf("Boxes should not collide")
	}
}

func TestCollideBoxBox1(t *testing.T) {
	slab := newBody(NewBox(50, 50, 50)).setMaterial(0, 0)
	slab.World().Loc.SetS(0, -50, 0)
	box := newBody(NewBox(1, 1, 1)).setMaterial(1, 0)
	box.World().Loc.SetS(-5.000000, 1.388000, -3.000000)
	box.World().Rot.SetS(0.182574, 0.365148, 0.547723, 0.730297)
	wantPoint, wantDepth := lin.NewV3S(-5.2, -0.1, -4.0), -0.108
	_, _, cs := collideBoxBox(slab, box, newManifold())
	if !lin.Aeq(cs[0].depth, wantDepth) || dumpV3(cs[0].point) != dumpV3(wantPoint) {
		t.Errorf("Got point %s wanted %s. Got depth %f wanted %f",
			dumpV3(cs[0].point), dumpV3(wantPoint), cs[0].depth, wantDepth)
	}
}

// Testing
// ============================================================================
// Benchmarking

// Ensure that collision is relatively fast. Run go test -bench=".*Collide"
// to get something like:
//     BenchmarkCollideAabb	            500000000	   7.52 ns/op
//     BenchmarkCollideSphereSphere	     20000000	  82.40 ns/op
//     BenchmarkCollideSphereBox	     10000000	 199    ns/op

func BenchmarkCollideAabb(b *testing.B) {
	a, o := &Abox{0, 0, 0, 1, 1, 1}, &Abox{-1, -1, -1, 0.1, 0.1, 0.1}
	for cnt := 0; cnt < b.N; cnt++ {
		a.Overlaps(o)
	}
}
func BenchmarkCollideSphereSphere(b *testing.B) {
	a, o, cs := NewBody(NewSphere(1)), NewBody(NewSphere(1)), newManifold()
	for cnt := 0; cnt < b.N; cnt++ {
		collideSphereSphere(a, o, cs)
	}
}
func BenchmarkCollideSphereBox(b *testing.B) {
	a, o, cs := NewBody(NewSphere(1)), NewBody(NewBox(1, 1, 1)), newManifold()
	for cnt := 0; cnt < b.N; cnt++ {
		collideSphereBox(a, o, cs)
	}
}

// These runs used preallocated scratch structures.
// Allocating structures each time costs 500 ns/op.
//     BenchmarkCollideBoxBox	         10000000	 159 ns/op (cgo call-only, no c-code, go-code)
//     BenchmarkCollideBoxBox	         10000000	 224 ns/op (commented out c-code)
//     BenchmarkCollideBoxBox (cgo impl)  5000000    704 ns/op
func BenchmarkCollideBoxBox(b *testing.B) {
	a, o, cs := NewBody(NewBox(0.5, 0.5, 0.5)), NewBody(NewBox(1, 1, 1)), newManifold()
	for cnt := 0; cnt < b.N; cnt++ {
		collideBoxBox(a, o, cs)
	}
}
