// Copyright © 2014-2024 Galvanized Logic Inc.

package vu

import (
	"testing"

	"github.com/gazed/vu/math/lin"
)

// go test -run Camera
func TestCamera(t *testing.T) {

	// Test perspective and view inverses.
	t.Run("camera inverse projections", func(t *testing.T) {
		cam, _, _ := initScene()
		cam.SetPitch(cam.Pitch + 45)
		cam.Move(0, -15, 15, cam.Lookat())
		cam.vt(cam.at, cam.vm)  // view transform - vp by default
		cam.it(cam.at, cam.ivm) // inverse view transform - ivp by default.

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
		rx, ry, rz := cam.Ray(ww/2, wh/2, ww, wh) // center of screen.
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
		blx, bly, _ := cam.Ray(0, 0, ww, wh)
		trx, try, _ := cam.Ray(ww, wh, ww, wh)
		gotRatio := (try - bly) / (blx - trx)
		expectedRatio := float64(wh) / float64(ww)
		if expectedRatio != gotRatio {
			t.Errorf("Expected %f  got %f", expectedRatio, gotRatio)
		}
	})

	// Test ray sphere intersection.
	t.Run("raycast", func(t *testing.T) {
		cam, ww, wh := initScene()                           // looking down -Z
		ray := lin.NewV3().SetS(cam.Ray(ww/2, wh/2, ww, wh)) // ray: 0,0,-1

		// hit sphere on Z axis.
		s := lin.NewV3().SetS(0, 0, -10)
		hit, _, _, _ := cam.RayCastSphere(ray, s, 1.0)
		if !hit {
			t.Errorf("expected -z hit")
		}

		// miss sphere off axis.
		s.SetS(0, -2, -10)
		hit, _, _, _ = cam.RayCastSphere(ray, s, 1.0)
		if hit {
			t.Errorf("expected miss")
		}

		// miss sphere on Z axis opposite direction
		s.SetS(0, 0, 10)
		hit, _, _, _ = cam.RayCastSphere(ray, s, 1.0)
		if hit {
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
	return
}
