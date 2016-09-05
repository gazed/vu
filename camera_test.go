// Copyright © 2014-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"math"
	"testing"

	"github.com/gazed/vu/math/lin"
)

// Test a ray cast with simple perspective and view inverses.
// The ray from center screen mouse coordinates should be directly
// along the -Z axis.
func TestRay(t *testing.T) {
	cam, ww, wh := initScene()
	cam.Move(0, 0, 15, cam.Lookat())
	rx, ry, rz := cam.Ray(ww/2, wh/2, ww, wh) // center of screen.
	ex, ey, ez := 0.0, 0.0, -1.0
	if rx != ex || ry != ey || rz != ez {
		t.Errorf("Expected %f %f %f got %f %f %f", ex, ey, ez, rx, ry, rz)
	}
}

// Test a ray cast with perspective inverse and angled view inverse.
func TestAngledRay(t *testing.T) {
	cam, ww, wh := initScene()
	cam.SetPitch(cam.Pitch + 45)
	cam.SetAt(0, -15, 15)
	rx, ry, rz := cam.Ray(ww/2, wh/2, ww, wh) // center of screen.
	ex, ey, ez := 0.0, 0.7071068, -0.7071068
	if !lin.Aeq(rx, ex) || !lin.Aeq(ry, ey) || !lin.Aeq(rz, ez) {
		t.Errorf("Expected %f %f %f got %f %f %f", ex, ey, ez, rx, ry, rz)
	}
}

// Test that the ratio of the rays matches the ratio of the screen.
func TestRayRatios(t *testing.T) {
	cam, ww, wh := initScene()
	cam.Move(0, 0, 15, cam.Lookat())

	// shoot and check opposing corner rays.
	blx, bly, _ := cam.Ray(0, 0, ww, wh)
	trx, try, _ := cam.Ray(ww, wh, ww, wh)
	gotRatio := (try - bly) / (trx - blx)
	expectedRatio := float64(wh) / float64(ww)
	if expectedRatio != gotRatio {
		t.Errorf("Expected %f  got %f", expectedRatio, gotRatio)
	}
}

// Test perspective and view inverses.
func TestInverses(t *testing.T) {
	cam, _, _ := initScene()
	cam.SetPitch(cam.Pitch + 45)
	cam.Move(0, -15, 15, cam.Lookat())

	// the inverses multiplied with non-inverses should be the identity matrix.
	if !lin.NewM4().Mult(cam.pm, cam.ipm).Aeq(lin.M4I) {
		t.Error("Invalid inverse projection matrix")
	}
	if !lin.NewM4().Mult(cam.vm, cam.ivm).Aeq(lin.M4I) {
		t.Error("Invalid inverse view matrix")
	}
}

// Check that the inverse of a perspective view is correct.
func TestInverseVp(t *testing.T) {
	v := newCamera()
	v.at.Loc.SetS(10, 10, 10)
	v.at.Rot.SetAa(1, 0, 0, -lin.Rad(90))
	vm := VP(v.at, lin.NewQ(), &lin.M4{})
	ivm := ivp(v.at, &lin.Q{X: 1, Y: 0, Z: 0, W: 1}, lin.NewQ(), &lin.M4{})
	if !vm.Mult(vm, ivm).Aeq(lin.M4I) {
		t.Errorf("Matrix times inverse should be identity")
	}
}

func TestRoundTrip(t *testing.T) {
	cam, _, _ := initScene()
	cx, cy, cz := 0.0, 0.0, 14.0 // camera location to
	cam.SetAt(cx, cy, cz)        // ...point directly at 0, 0, 0

	// Create the matricies to go between clip and world space.
	toClip := lin.NewM4().Mult(cam.vm, cam.pm)
	toWorld := lin.NewM4().Mult(cam.ipm, cam.ivm)
	if !lin.NewM4().Mult(toClip, toWorld).Aeq(lin.M4I) {
		t.Errorf("Invalid world<->clip matricies")
	}

	// start with world coordinates carefully chosen to give x=1, y=1 clip values
	px, py := 6.002062, 3.751289
	pnt := lin.NewV4().SetS(px, py, 0, 1)
	pnt.MultMv(toClip, pnt)
	if !lin.Aeq(pnt.X/pnt.W, 1) || !lin.Aeq(pnt.Y/pnt.W, 1) {
		t.Errorf("%f %f gave clip %f %f %f, expected (1 1 -0.071429)", px, py, pnt.X, pnt.Y, pnt.Z)
	}

	// now reverse back to world coordinates.
	pnt.MultMv(toWorld, pnt)
	if !lin.Aeq(pnt.X, px) || !lin.Aeq(pnt.Y, py) {
		t.Errorf("got point %f %f %f, expected x=%f y=%f", pnt.X, pnt.Y, pnt.Z, px, py)
	}
}

func TestRayWithSpin(t *testing.T) {
	cam, _, _ := initScene()
	cx, cy, cz := 0.0, -10.0, 14.0             // camera location to
	cam.SetAt(cx, cy, cz)                      // ...point directly at 0, 0, 0
	cam.SetPitch(lin.Deg(math.Atan(-cy / cz))) // 35.53768 degrees
	plane := NewPlane(0, 0, -1)

	ww, wh := 1280, 800
	rx, ry, rz := cam.Ray(0, 0, ww, wh)
	ray := NewRay(rx, ry, rz)
	ray.World().SetLoc(cx, cy, cz)
	hit, hx, hy, hz := Cast(ray, plane)
	ex, ey, ez := -6.191039, -4.755119, 0.0
	if !hit || !lin.Aeq(hx, ex) || !lin.Aeq(hx, ex) || !lin.Aeq(hx, ex) {
		t.Errorf("Hit %t %f %f %f, expected %f %f %f", hit, hx, hy, hz, ex, ey, ez)
	}

	rx, ry, rz = cam.Ray(0, wh, ww, wh)
	ray = NewRay(rx, ry, rz)
	ray.World().SetLoc(cx, cy, cz)
	hit, hx, hy, hz = Cast(ray, plane)
	ex, ey, ez = -9.121797, 7.006131, 0.0
	if !hit || !lin.Aeq(hx, ex) || !lin.Aeq(hx, ex) || !lin.Aeq(hx, ex) {
		t.Errorf("Hit %t %f %f %f, expected %f %f %f", hit, hx, hy, hz, ex, ey, ez)
	}

	rx, ry, rz = cam.Ray(ww, 0, ww, wh)
	ray = NewRay(rx, ry, rz)
	ray.World().SetLoc(cx, cy, cz)
	hit, hx, hy, hz = Cast(ray, plane)
	ex, ey, ez = 6.191039, -4.755119, 0.0
	if !hit || !lin.Aeq(hx, ex) || !lin.Aeq(hx, ex) || !lin.Aeq(hx, ex) {
		t.Errorf("Hit %t %f %f %f, expected %f %f %f", hit, hx, hy, hz, ex, ey, ez)
	}

	rx, ry, rz = cam.Ray(ww, wh, ww, wh)
	ray = NewRay(rx, ry, rz)
	ray.World().SetLoc(cx, cy, cz)
	hit, hx, hy, hz = Cast(ray, plane)
	ex, ey, ez = 9.121797, 7.006131, 0.0
	if !hit || !lin.Aeq(hx, ex) || !lin.Aeq(hx, ex) || !lin.Aeq(hx, ex) {
		t.Errorf("Hit %t %f %f %f, expected %f %f %f", hit, hx, hy, hz, ex, ey, ez)
	}
}

func TestScreen(t *testing.T) {
	cam, _, _ := initScene()
	cx, cy, cz := 0.0, 0.0, 14.0 // camera location to
	cam.SetAt(cx, cy, cz)        // ...point directly at 0, 0, 0

	// center of the world should give the center of the screen.
	px, py, pz := 0.0, 0.0, 0.0
	if x, y := cam.Screen(px, py, pz, 1280, 800); x != 640 || y != 400 {
		t.Errorf("got point %d %d, expected 640, 400", x, y)
	}
}

// =============================================================================
// test utility methods.

// initScene creats a scene with an initialized perspective matrix.
func initScene() (c *Camera, ww, wh int) {
	c = newCamera()
	ww, wh = 1280, 800
	fov, ratio, near, far := 30.0, float64(ww)/float64(wh), 0.1, 500.0
	c.SetPerspective(fov, ratio, near, far)
	return
}
