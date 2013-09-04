// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package physics

import (
	"fmt"
	"testing"
	"vu/math/lin"
)

// Dictate how errors get printed.
const format = "\ngot\n%s\nwanted\n%s"

func TestNewSphere(t *testing.T) {
	s := newSphere(1, 2, 3, 4)
	if s.center.X != 1 || s.center.Y != 2 || s.center.Z != 3 || s.radius != 4 {
		t.Error("Improperly created new sphere")
	}
}

func TestSphereCollideSphere(t *testing.T) {
	s1 := newSphere(1, 1, 1, 1)
	s2 := newSphere(1, 1, 1, 1)
	c := s1.collideSphere(s2)
	if c == nil {
		t.Error("Identical spheres should touch")
	}
	got := fmt.Sprintf("%2.2f %2.2f", c[0].Normal, c[0].Depth)
	want := "{1.00 0.00 0.00} 2.00"
	if got != want {
		t.Errorf(format, got, want)
	}

	// check spheres that just touch at a point.
	s1 = newSphere(1, 1, 0, 1)
	s2 = newSphere(1, -1, 0, 1)
	c = s1.collideSphere(s2)
	if c == nil {
		t.Error("Spheres should touch at y=0 on x-axis")
	}
	got = fmt.Sprintf("%2.2f %2.2f", c[0].Normal, c[0].Depth)
	want = "{0.00 1.00 0.00} 0.00"
	if got != want {
		t.Errorf(format, got, want)
	}

	// check partially overlapping spheres
	s1 = newSphere(1, 1, 1, 1)
	s2 = newSphere(2, 2, 2, 1)
	c = s1.collideSphere(s2)
	if c == nil {
		t.Error("Spheres should overlap. Distance is 1 < radii of 2")
	}
	got = fmt.Sprintf("%2.2f %2.2f", c[0].Normal, c[0].Depth)
	want = "{-0.58 -0.58 -0.58} 0.27"
	if got != want {
		t.Errorf(format, got, want)
	}

	// check non overlapping spheres
	s1 = newSphere(1, 1, 1, 1)
	s2 = newSphere(-1, -1, -1, 1)
	c = s1.collideSphere(s2)
	if c != nil {
		t.Error("Spheres should not touch. Distance is 3.5 > radii of 2")
	}
}

func TestSphereBounce(t *testing.T) {
	s1 := newSphere(1, 1, 1, 1)
	mo1 := newMotion(1, s1.Volume()*10, &lin.V3{}, lin.QIdentity())
	mo1.setLinearMomentum(&lin.V3{-1, -1, -1})
	s2 := newSphere(2, 2, 2, 1)
	mo2 := newMotion(1, s2.Volume()*20, &lin.V3{}, lin.QIdentity())
	mo2.setLinearMomentum(&lin.V3{1, 1, 1})
	c := s1.collideSphere(s2)
	s1.bounceSphere(s2, c[0], mo1, mo2)
	got := fmt.Sprintf("%2.2f %2.2f", mo1.linearMomentum(), mo2.linearMomentum())

	// The two balls are heading opposite their original direction and the total
	// momentum has been conserved.
	want := "&{1.33 1.33 1.33} &{-0.67 -0.67 -0.67}"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestSphereCollidePlane(t *testing.T) {

	// plane intesects middle of sphere
	s := newSphere(1, 1, 1, 1)
	p := newPlane(1, 0, 0, 1, 1, 1)
	c := s.collidePlane(p)
	if c == nil {
		t.Error("Sphere and plane should overlap.")
	}

	// plane intesects edge of sphere
	s = newSphere(1, 1, 1, 1)
	p = newPlane(0, 0, 1, 1, 1, 0)
	c = s.collidePlane(p)
	if c == nil {
		t.Error("Sphere and plane should overlap.")
	}

	// plane does not intesect sphere
	s = newSphere(1, 1, 1, 1)
	p = newPlane(0, 1, 0, 0, -2, 0)
	c = s.collidePlane(p)
	if c != nil {
		t.Error("Sphere and plane should not overlap.")
	}
}

// Also tests reversing the object test types.
func TestSphereCollideRay(t *testing.T) {
	s := newSphere(0, 0, -2, 1)
	r0 := newRay(0, 0, 2, 0, 0, 4)  // 0 0 2 -  0  0 -2 center
	r1 := newRay(0, 0, 2, 1, 0, 4)  // 0 0 2 - -1  0 -2 left
	r2 := newRay(0, 0, 2, -1, 0, 4) // 0 0 2 -  1  0 -2 right
	r3 := newRay(0, 0, 2, 0, -1, 4) // 0 0 2 -  0  1 -2 top
	r4 := newRay(0, 0, 2, 0, 1, 4)  // 0 0 2 -  0 -1 -2 bottom
	for _, ray := range []*ray{r0, r1, r2, r3, r4} {
		c := s.Collide(ray)
		if c == nil {
			t.Error("Sphere and ray should overlap at one point.")
		}
	}

	// test some misses
	r6 := newRay(0, 0, 2, 1.1, 0, 4)  // miss left
	r7 := newRay(0, 0, 2, -1.1, 0, 4) // miss right
	r8 := newRay(0, 0, 2, 0, -1.1, 4) // miss top
	r9 := newRay(0, 0, 2, 0, 1.1, 4)  // miss bottom
	for _, ray := range []*ray{r6, r7, r8, r9} {
		c := s.Collide(ray)
		if c != nil {
			t.Error("Sphere and ray should miss.")
		}
	}
}

func TestSphereCollideAbox(t *testing.T) {

	// abox intesects
	s := newSphere(-1.5, 0, 0, 1)
	a := newAbox(-1, -1, -1, 1, 1, 1)
	c := s.collideAbox(a)
	if c == nil {
		t.Error("Sphere and abox should overlap.")
	}
	got := fmt.Sprintf("%2.2f %2.2f", c[0].Normal, c[0].Depth)
	want := "{-1.00 0.00 0.00} 0.50"
	if got != want {
		t.Errorf(format, got, want)
	}

	// intersects opposite face so normal should be flipped.
	s = newSphere(1.5, 0, 0, 1)
	a = newAbox(-1, -1, -1, 1, 1, 1)
	c = s.collideAbox(a)
	if c == nil {
		t.Error("Sphere and abox should overlap.")
	}
	got = fmt.Sprintf("%2.2f %2.2f", c[0].Normal, c[0].Depth)
	want = "{1.00 0.00 0.00} 0.50"
	if got != want {
		t.Errorf(format, got, want)
	}

	// abox intesects edge of sphere
	s = newSphere(0, 2, 0, 1)
	a = newAbox(-1, -1, -1, 1, 1, 1)
	c = s.collideAbox(a)
	if c == nil {
		t.Error("Sphere and abox should overlap.")
	}
	got = fmt.Sprintf("%2.2f %2.2f", c[0].Normal, c[0].Depth)
	want = "{0.00 1.00 0.00} 0.00"
	if got != want {
		t.Errorf(format, got, want)
	}

	// abox does not intesect sphere
	s = newSphere(1, 1, 1, 1)
	a = newAbox(-1, -1, -1, 0, 0, 0)
	c = s.collideAbox(a)
	if c != nil {
		t.Error("Sphere and abox should not overlap.")
	}
}

func TestPlaneBounce(t *testing.T) {
	s := newSphere(1, 1, 1, 1)
	mo := newMotion(1, s.Volume()*10, &lin.V3{}, lin.QIdentity())
	mo.setLinearMomentum(&lin.V3{-1, -1, -1})
	p := newPlane(1, 0, 0, 1, 1, 1)
	c := s.collidePlane(p)
	s.bouncePlane(p, c[0], mo, nil)
	got := fmt.Sprintf("%2.2f", mo.linearMomentum())

	// The ball is heading as a reflection around the plane normal.
	want := "&{0.90 -0.90 -0.90}"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestRayBounce(t *testing.T) {
	s := newSphere(1, 1, 1, 1)
	mo1 := newMotion(1, s.Volume()*10, &lin.V3{}, lin.QIdentity())
	mo1.setLinearMomentum(&lin.V3{-1, -1, -1})
	ray := newRay(0, 0, 0, 1, 1, 1)
	mo2 := newMotion(1, 1, &lin.V3{}, lin.QIdentity())
	mo2.setLinearMomentum(&lin.V3{-1, -1, -1})
	contacts := s.Collide(ray)
	s.Bounce(ray, contacts, mo1, mo2)
	got := fmt.Sprintf("%2.2f %2.2f", mo1.linearMomentum(), mo2.linearMomentum())

	// The ball is now the sum of both momentums.
	want := "&{-2.00 -2.00 -2.00} &{0.00 0.00 0.00}"
	if got != want {
		t.Errorf(format, got, want)
	}
}

// TODO Fix. this test shows the sphere::abox collision is not generating the
//      proper collision point data.
// check a sphere that has collided with the corners of two aboxes.
// Should get back 4 contact locations as the sphere intersects two sides
// of each abox.
//func TestSphereCollideMultipleAbox(t *testing.T) {
//	a1 := newAbox(1, 1, 1, 3, 3, 3) // center at 2, 2, 2
//	a2 := newAbox(3, 1, 1, 5, 3, 3) // center at 4, 2, 2
//	s := newSphere(3, 0.5, 0.5, 1)
//	c1 := s.collideAbox(a1)
//	c2 := s.collideAbox(a2)
//
//	if len(c1) != 2 || len(c2) != 2 {
//		println("c1", len(c1), c1[0].Normal.X, c1[0].Normal.Y, c1[0].Normal.Z)
//		println("c2", len(c2), c2[0].Normal.X, c2[0].Normal.Y, c2[0].Normal.Z)
//		t.Error("Not detecting Sphere overlapping Abox corner.")
//	}
//}

func TestAboxCollideAbox(t *testing.T) {
	// touching
	a1 := newAbox(1, 1, 1, 3, 3, 3) // center at 2, 2, 2
	a2 := newAbox(3, 1, 1, 5, 3, 3) // center at 4, 2, 2
	c1 := a1.collideAbox(a2)
	if len(c1) >= 1 {
		t.Error("Aboxes were touching, not overlapping.")
	}

	// overlapping on positive X axis
	a1 = newAbox(1, 1, 1, 3, 3, 3) // center at 2, 2, 2
	a2 = newAbox(2, 1, 1, 4, 3, 3) // center at 3, 2, 2
	c1 = a1.collideAbox(a2)
	if len(c1) != 1 || c1[0].Depth != 1 || c1[0].Normal.X != 1 {
		t.Error("Aboxes should overlap on the positive x axis.")
	}
}
