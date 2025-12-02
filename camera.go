// SPDX-FileCopyrightText : © 2014-2024 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

package vu

// camera.go holds the view and projection matricies needed for rendering.

import (
	"fmt"

	"github.com/gazed/vu/math/lin"
)

// Camera makes rendered models visible within a frame. A camera is
// associated with a scene where it is used to render the scenes models.
//
// Camera combines a location+orientation using separate pitch angle
// tracking. This allows use as a first-person camera which can limit
// up/down to a given range (often 180 degrees). Overall orientation is
// calculated by combining Pitch and Yaw. Look is for walking cameras,
// Lookat is for flying cameras.
type Camera struct {

	// Position and orientation.
	at   *lin.T // Combined location and Pitch/Yaw orientation.
	yrot *lin.Q // Y-axis quaternion rotation updated by SetYaw.
	xrot *lin.Q // X-axis quaternion rotation updated by SetPitch.

	// View matrix. The V part of the MVP transform matrix
	vm  *lin.M4 // camera view matrix
	ivm *lin.M4 // Inverse camera view matrix.

	// Projection matrix. The P part of the MVP transform matrix.
	near, far float64 // Frustum clip set by application.
	fov       float64 // Field of view set by application.
	focus     bool    // true if the camera projection needs setting.
	pm        *lin.M4 // Projection matrix.
	ipm       *lin.M4 // Inverse projection matrix.
}

// newCamera creates a default rendering field that is looking
// down the negative Z axis with positive Y up.
func newCamera() *Camera {
	c := &Camera{fov: 90, focus: true} // Default fov.
	c.at = lin.NewT()
	c.yrot = lin.NewQ().SetAa(0, 1, 0, 0)
	c.xrot = lin.NewQ().SetAa(0, 0, 0, 0)
	c.vm = &lin.M4{}
	c.ivm = &lin.M4{}
	c.pm = &lin.M4{}
	c.ipm = &lin.M4{}
	return c
}

// SetClip sets the near and far clipping planes for perspective
// and orthographic cameras.
func (c *Camera) SetClip(near, far float64) *Camera {
	c.near, c.far, c.focus = near, far, true
	return c
}

// SetFov sets the field of view for perspective projection cameras.
// Ignored for orthographic projection cameras.
func (c *Camera) SetFov(deg float64) *Camera {
	c.fov, c.focus = deg, true
	return c
}

// At returns the cameras current location in world space.
func (c *Camera) At() (x, y, z float64) {
	return c.at.Loc.GetS()
}

// SetAt positions the camera in world space
// The camera instance is returned.
func (c *Camera) SetAt(x, y, z float64) *Camera {
	c.at.Loc.SetS(x, y, z)
	return c
}

// Move adjusts the camera location relative to the given orientation.
// For orientation, use Lookat() to fly, use Look to run along XZ.
func (c *Camera) Move(x, y, z float64, q *lin.Q) {
	dx, dy, dz := lin.MultSQ(x, y, z, q)
	c.at.Loc.X += dx
	c.at.Loc.Y += dy
	c.at.Loc.Z += dz
}

// Lookat returns the current camera rotation.
// The rotation is created from Pitch and Yaw.
func (c *Camera) Lookat() *lin.Q { return c.at.Rot }

// SetLook directly sets the camera orientation to the given quaternion.
// The orientation will be overwritten with new calls to SetPitch and SetYaw.
func (c *Camera) SetLook(q *lin.Q) { c.at.Rot.Set(q) }

// SetPitch sets the rotation around the X axis and updates
// the Look direction. The camera instance is returned.
func (c *Camera) SetPitch(deg float64) *Camera {
	c.xrot.SetAa(1, 0, 0, lin.Rad(deg))
	c.at.Rot.Mult(c.xrot, c.yrot).Unit()
	return c
}

// SetYaw sets the rotation around the Y axis and updates the
// Look and Lookat directions. The camera instance is returned.
func (c *Camera) SetYaw(deg float64) *Camera {
	c.yrot.SetAa(0, 1, 0, lin.Rad(deg))
	c.at.Rot.Mult(c.xrot, c.yrot).Unit()
	return c
}

// Ray applies inverse transforms to derive world space coordinates
// for a ray projected from the camera through the mouse's mx,my
// screen position given window width and height ww,wh.
func (c *Camera) Ray(mx, my, ww, wh int) (x, y, z float64, err error) {
	ray := lin.NewV3().SetS(0, 0, 0)
	if mx >= 0 && mx <= ww && my >= 0 && my <= wh {
		clipx := float64(2*mx)/float64(ww) - 1 // mx to range -1:1
		clipy := float64(2*my)/float64(wh) - 1 // my to range -1:1
		clip := lin.NewV4().SetS(clipx, clipy, -1, 1)

		// Use inverse perspective to go from clip to eye (view) coordinates.
		eye := clip.MultvM(clip, c.ipm)
		eye.Z = -1 // into the screen
		eye.W = 0  // want a vector, not a point

		// Use inverse view to go from eye (view) to world coordinates.
		world := eye.MultvM(eye, c.ivm)
		ray.SetS(world.X, world.Y, world.Z) // ignore the W component.
		ray.Unit()                          // return a unit vector.
		return ray.X, ray.Y, ray.Z, nil
	}
	return 0, 0, 0, fmt.Errorf("mouse not in window")
}

// RayCastSphere checks for collision between a ray originating from the camera
// and a sphere in world space. The ray must be a unit vector, see: camera.Ray().
//   - see: http://en.wikipedia.org/wiki/Line–sphere_intersection
func (c *Camera) RayCastSphere(ray, sphere *lin.V3, radius float64) (hit bool) {

	// vector from ray origin to sphere center
	rox, roy, roz := c.At() // ray origin is the camera world location.
	rs := lin.NewV3().SetS(sphere.X-rox, sphere.Y-roy, sphere.Z-roz)

	// distance between the center of the sphere and the ray.
	// If the distance is larger than the radius there is no intersection.
	d0 := rs.Dot(ray)
	if d0 < 0 {
		return false // no hit
	}
	d1 := rs.Dot(rs) - d0*d0
	if d1 > radius*radius {
		return false // no hit
	}
	return true // hit

	// FUTURE could get contact points.
	// dlen := d0 - math.Sqrt(radius*radius-d1)
	// d0 - dlen // ray length point 1
	// d0 + dlen // ray length point 2
}

// RayCastDisk checks for collision between a ray originating from the camera
// and a circle in world space. The ray must be a unit vector, see: camera.Ray().
// See: https://en.wikipedia.org/wiki/Line%E2%80%93plane_intersection
func (c *Camera) RayCastDisk(ray, center, normal *lin.V3, radius float64) (hit bool) {

	// check if the ray intersects the plane of the disk.
	// If this is zero then they are parallel (no intersection)
	// or the ray is entirely in the plane.
	if denom := ray.Dot(normal); !lin.AeqZ(denom) {

		// only intersect camera facing faces.
		if denom < 0 {

			// vector from ray origin to disk center
			rox, roy, roz := c.At() // ray origin is the camera world location.
			rayDisk := lin.NewV3().SetS(center.X-rox, center.Y-roy, center.Z-roz)

			// intersection point of ray with disk plane
			t := rayDisk.Dot(normal) / denom
			px, py, pz := rox+ray.X*t, roy+ray.Y*t, roz+ray.Z*t

			// is the intersection point within the disk radius.
			v := lin.NewV3().SetS(px-center.X, py-center.Y, pz-center.Z)
			d2 := v.Dot(v)
			if d2 < radius*radius {
				return true
			}
		}
	}
	return false
}

// Screen applies the camera transform on a 3D point in world space wx,wy,wz
// and returns the 2D screen coordinate sx,sy. The window width and height
// ww,wh are needed. Essentially the reverse of the Ray method and duplicates
// what is done in the rendering pipeline.
// Returns -1,-1 if the point is outside the screen area.
func (c *Camera) Screen(wx, wy, wz float64, ww, wh int) (sx, sy int) {
	vec := lin.NewV4().SetS(wx, wy, wz, 1)
	vec.MultvM(vec, c.vm)          // apply view matrix.
	vec.MultvM(vec, c.pm)          // apply projection matrix.
	clipx := vec.X*0.5/vec.W + 0.5 // convert to range 0:1
	clipy := vec.Y*0.5/vec.W + 0.5 // convert to range 0:1
	clipz := vec.Z*0.5/vec.W + 0.5 // convert to range 0:1
	if clipx < 0 || clipx > 1 || clipy < 0 || clipy > 1 || clipz < 0 || clipz > 1 {
		return -1, -1 // outside the screen area.
	}
	sx = int(lin.Round(clipx*float64(ww), 0))
	sy = int(lin.Round(clipy*float64(wh), 0))
	return sx, sy
}

// distance returns the distance squared of the camera to the given Pov.
// Uses the existing Pov world coordinates.
func (c *Camera) distance(wx, wy, wz float64) float64 {
	dx := wx - c.at.Loc.X
	dy := wy - c.at.Loc.Y
	dz := wz - c.at.Loc.Z
	return float64(dx*dx + dy*dy + dz*dz)
}

// setPerspective makes the camera use a 3D projection.
// This is the projection part of model-view-projection.
func (c *Camera) setPerspective(fov, ratio, near, far float64) {
	c.pm.PerspectiveProjection(fov, ratio, near, far)
	c.ipm.PerspectiveInverse(fov, ratio, near, far)
}

// setOrthographic makes the camera use a 2D projection.
// This is the projection part of model-view-projection.
func (c *Camera) setOrthographic(left, right, bottom, top, near, far float64) {
	c.pm.OrthographicProjection(left, right, bottom, top, near, far)
}

// updateView recalulates the view matricies.
func (c *Camera) updateView() {
	// Set the view transform matrix
	c.vm.SetQ(c.at.Rot)
	c.vm.TranslateTM(-c.at.Loc.X, -c.at.Loc.Y, -c.at.Loc.Z)

	// Set the view inverse transform matrix
	c.ivm.SetQ(lin.NewQ().Inv(c.at.Rot))
	c.ivm.TranslateMT(c.at.Loc.X, c.at.Loc.Y, c.at.Loc.Z)
}
