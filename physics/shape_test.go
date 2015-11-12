// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package physics

import (
	"testing"

	"github.com/gazed/vu/math/lin"
)

func TestBox(t *testing.T) {
	bx := Shape(NewBox(1, 1, 1)) // compiler checks Shape interface.
	if bx.Type() != BoxShape {
		t.Error("Invalid box shape")
	}
}

func TestBoxAabb(t *testing.T) {
	bx := Shape(NewBox(1, 1, 1))
	ab := bx.Aabb(lin.NewT().SetI(), &Abox{}, 0.01)
	if ab.Sx != -1.01 || ab.Sy != -1.01 || ab.Sz != -1.01 || ab.Lx != 1.01 || ab.Ly != 1.01 || ab.Lz != 1.01 {
		t.Error("Invalid bounding box for Box")
	}
}

func TestBoxVolume(t *testing.T) {
	bx := Shape(NewBox(1, 1, 1))
	if bx.Volume() != 8 {
		t.Errorf("Expected box volume 8, got %f", bx.Volume())
	}
}

func TestBoxInertia(t *testing.T) {
	bx, inertia, want := Shape(NewBox(1, 1, 1)), lin.NewV3(), "{0.7 0.7 0.7}"
	if bx.Inertia(1, inertia); dumpV3(inertia) != want {
		t.Errorf("Expected box inertia %s, got %s", want, dumpV3(inertia))
	}
}

func TestSphere(t *testing.T) {
	sp := Shape(NewSphere(1)) // compiler checks Shape interface.
	if sp.Type() != SphereShape {
		t.Error("Invalid sphere shape")
	}
}

func TestSphereAabb(t *testing.T) {
	sp := Shape(NewSphere(1))
	ab := sp.Aabb(lin.NewT().SetI(), &Abox{}, 0.01)
	if ab.Sx != -1.01 || ab.Sy != -1.01 || ab.Sz != -1.01 || ab.Lx != 1.01 || ab.Ly != 1.01 || ab.Lz != 1.01 {
		t.Error("Invalid bounding box for Sphere")
	}
}

func TestSphereVolume(t *testing.T) {
	sp := Shape(NewSphere(1.25))
	if !lin.Aeq(sp.Volume(), 6.13592315) {
		t.Errorf("Expected sphere mass 6.13592315, got %2.8f", sp.Volume())
	}
}

func TestSphereInertia(t *testing.T) {
	sp, inertia, want := Shape(NewSphere(1.25)), lin.NewV3(), "{0.6 0.6 0.6}"
	if sp.Inertia(1, inertia); dumpV3(inertia) != want {
		t.Errorf("Expected sphere inertia %s, got %s", want, dumpV3(inertia))
	}
}

func TestAboxOverlap(t *testing.T) {
	var a, b, c, d *Abox
	a, b = &Abox{0, 0, 0, 1, 1, 1}, &Abox{-1, -1, -1, 0, 0, 0}
	if a.Overlaps(b) {
		t.Error("Touching at a point, but not overlapping")
	}
	b = &Abox{-1, -1, -1, 0.1, 0.0, 0.0}
	c = &Abox{-1, -1, -1, 0.0, 0.1, 0.0}
	d = &Abox{-1, -1, -1, 0.0, 0.0, 0.1}
	if a.Overlaps(b) || a.Overlaps(c) || a.Overlaps(d) {
		t.Error("Touching along edges, but not overlapping")
	}
	b = &Abox{-1, -1, -1, 0.1, 0.1, 0.0}
	c = &Abox{-1, -1, -1, 0.0, 0.1, 0.1}
	d = &Abox{-1, -1, -1, 0.1, 0.0, 0.1}
	if a.Overlaps(b) || a.Overlaps(c) || a.Overlaps(d) {
		t.Error("Touching along faces, but not overlapping")
	}
	b = &Abox{-1, -1, -1, 0.1, 0.1, 0.1}
	if !a.Overlaps(b) || !b.Overlaps(a) {
		t.Error("Overlapping")
	}
}
