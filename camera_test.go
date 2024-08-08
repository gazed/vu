// Copyright © 2014-2024 Galvanized Logic Inc.

package vu

import (
	"testing"

	"github.com/gazed/vu/math/lin"
)

// go test -run Camera
func TestCamera(t *testing.T) {

	// check the camera world position after rotating around X,Y axis.
	// +Y:up +X:right -Z:forward
	t.Run("camera rotations", func(t *testing.T) {
		cam, _, _ := initScene() //
		cam.SetAt(0, 0, 1)       // start at 1 on the Z-axis looking forward to origin.

		// expected world camera position is opposite in order to
		// move the world and keep the camera at 0,0,0
		axismap := []int{
			// pitch:yaw <=> x:y:z camera world position
			+00, +90 /* <=> */, +1, +0, +0, // rotate about Y:+90
			+00, -90 /* <=> */, -1, +0, +0, // rotate about Y:-90
			-90, +00 /* <=> */, +0, +1, +0, // rotate about X:-90
			+90, +00 /* <=> */, +0, -1, +0, // rotate about X:+90
			+00, 180 /* <=> */, +0, +0, +1, // rot about Y:180
			+00, +00 /* <=> */, +0, +0, -1, // no rotation
		}
		var wx, wy, wz float64
		for i := 0; i < len(axismap); i += 5 {
			a0, b0 := axismap[i], axismap[i+1] // spherical angles.
			cam.SetPitch(float64(a0))
			cam.SetYaw(float64(b0))
			cam.updateView()

			// expected rectangular world coordinates for the given angles.
			x0, y0, z0 := axismap[i+2], axismap[i+3], axismap[i+4]
			wx, wy, wz = cam.vm.Wx, cam.vm.Wy, cam.vm.Wz // world camera position

			if x0 != int(wx) || y0 != int(wy) || z0 != int(wz) {
				t.Errorf("expected %d %d %d got %d %d %d", x0, y0, z0, int(wx), int(wy), int(wz))
			}
		}
	})

	// Test perspective and view inverses.
	t.Run("camera inverse projections", func(t *testing.T) {
		cam, _, _ := initScene()
		cam.SetPitch(45)
		cam.Move(0, -15, 15, cam.Lookat())
		cam.updateView()

		// the inverses multiplied with non-inverses should be the identity matrix.
		if !lin.NewM4().Mult(cam.pm, cam.ipm).Aeq(lin.M4I) {
			t.Error("invalid inverse projection matrix")
		}
		if !lin.NewM4().Mult(cam.vm, cam.ivm).Aeq(lin.M4I) {
			t.Error("invalid inverse view matrix")
		}
	})
}

// go test -run Ray
func TestRay(t *testing.T) {

	// Test a ray cast with simple perspective and view inverses.
	// The ray from center screen mouse coordinates should be directly
	// along the -Z axis.
	t.Run("rayz", func(t *testing.T) {
		cam, ww, wh := initScene()
		rx, ry, rz, _ := cam.Ray(ww/2, wh/2, ww, wh) // center of screen.
		ex, ey, ez := 0.0, 0.0, -1.0
		if rx != ex || ry != ey || rz != ez {
			t.Errorf("expected %f %f %f got %f %f %f", ex, ey, ez, rx, ry, rz)
		}
	})

	// Test that the ratio of the rays matches the ratio of the screen.
	t.Run("ray ratios", func(t *testing.T) {
		cam, ww, wh := initScene()
		cam.Move(0, 0, 15, cam.Lookat())

		// shoot and check opposing corner rays.
		blx, bly, _, _ := cam.Ray(0, 0, ww, wh)
		trx, try, _, _ := cam.Ray(ww, wh, ww, wh)
		gotRatio := (try - bly) / (blx - trx)
		expectedRatio := float64(wh) / float64(ww)
		if expectedRatio != gotRatio {
			t.Errorf("Expected %f  got %f", expectedRatio, gotRatio)
		}
	})

	// Test ray sphere intersection.
	t.Run("raycast", func(t *testing.T) {
		cam, ww, wh := initScene()                   // looking down -Z
		rx, ry, rz, _ := cam.Ray(ww/2, wh/2, ww, wh) //
		ray := lin.NewV3().SetS(rx, ry, rz)          // ray: 0,0,-1

		// hit sphere on Z axis.
		s := lin.NewV3().SetS(0, 0, -10)
		if !cam.RayCastSphere(ray, s, 1.0) {
			t.Errorf("expected -z hit")
		}

		// miss sphere off axis.
		s.SetS(0, -2, -10)
		if cam.RayCastSphere(ray, s, 1.0) {
			t.Errorf("expected miss")
		}

		// miss sphere on Z axis opposite direction
		s.SetS(0, 0, 10)
		if cam.RayCastSphere(ray, s, 1.0) {
			t.Errorf("expected +z miss")
		}
	})
}

// =============================================================================
// test utility methods.

// initScene creats a camera with an initialized perspective matrix.
func initScene() (c *Camera, ww, wh int) {
	c = newCamera()
	ww, wh = 1280, 800
	fov, ratio, near, far := 30.0, float64(ww)/float64(wh), 0.1, 500.0
	c.setPerspective(fov, ratio, near, far)
	c.updateView()
	return
}
