// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"github.com/gazed/vu/math/lin"
)

// Camera tracks the location and orientation of a camera as well as an
// associated projection transform.
type Camera interface {
	Location() (x, y, z float64)    // Get, or
	SetLocation(x, y, z float64)    // ...Set the camera location.
	Rotation() (x, y, z, w float64) // Get, or
	SetRotation(x, y, z, w float64) // ...Set the view orientation.
	Move(x, y, z float64)           // Adjust current camera location.
	Spin(x, y, z float64)           // Rotate degrees about the given axis.
	Tilt() (up float64)             // Get, or
	SetTilt(up float64)             // ...Set the camera tilt angle.
	SetTransform(transform int)     // Set a view transform.

	// Use one of the following to create a projection transform.
	SetPerspective(fov, ratio, near, far float64)                // 3D.
	SetOrthographic(left, right, bottom, top, near, far float64) // 2D.

	// Ray applies inverse transforms to derive world space coordinates for
	// a ray projected from the camera through the mouse's mx, my screen
	// position given window width and height ww, wh.
	Ray(mx, my, ww, wh int) (x, y, z float64)

	// Screen calculates screen coordinates sx, sy for world coordinates
	// wx, wy, wz and window width and height ww, wh.
	Screen(wx, wy, wz float64, ww, wh int) (sx, sy int)

	// Distance returns the distance squared of the camera to the given point.
	Distance(px, py, pz float64) float64
}

// ===========================================================================

// camera combines a location+direction (pov) with a separate up/down angle
// tracking.  This allows use as a FPS camera which can limit up/down to
// a given range (often 180).
type camera struct {
	pov               // Location/direction.
	up  float64       // The up/down angle in degrees.  Limit this to +90/-90
	vt  viewTransform // The assigned view transform matrix.

	// Track of the view, projection matricies and thier inverses.
	vm  *lin.M4 // View part of MVP matrix.
	ivm *lin.M4 // Inverse view matrix.
	pm  *lin.M4 // Projection part of MVP matrix.
	ipm *lin.M4 // Inverse projection matrix.
	q0  *lin.Q  // Scratch for camera transform calculations.
	v0  *lin.V4 // Scratch for pick ray calculations.
	ray *lin.V3 // Scratch for pick ray calculations.
}

// newCamera creates a default point of view that is looking down the positive
// Z axis.
func newCamera() *camera {
	c := &camera{}
	c.pov = newPov()
	c.vm = &lin.M4{}
	c.ivm = (&lin.M4{}).Set(lin.M4I)
	c.pm = &lin.M4{}
	c.ipm = &lin.M4{}
	c.q0 = &lin.Q{}
	c.v0 = &lin.V4{}
	c.ray = &lin.V3{}
	return c
}

// setViewTransform turns the publicly visible view transform choices
// into a view tranform function.
func (c *camera) SetTransform(transform int) {
	switch transform {
	case VP:
		c.vt = vp
	case VO:
		c.vt = vo
	case VF:
		c.vt = vf
	case XZ_XY:
		c.vt = xz_xy
	default:
		c.vt = vp
	}
}

// transform applies the view transform to the scene camera and returns
// the result in the supplied matrix.
func (c *camera) transform(vm *lin.M4) *lin.M4 { return c.vt(c, vm) }

func (c *camera) Rotation() (x, y, z, w float64) {
	return c.pov.Rot.X, c.pov.Rot.Y, c.pov.Rot.Z, c.pov.Rot.W
}
func (c *camera) SetRotation(x, y, z, w float64) {
	c.pov.Rot.X, c.pov.Rot.Y, c.pov.Rot.Z, c.pov.Rot.W = x, y, z, w
	c.updateViewTransform()
}
func (c *camera) Location() (x, y, z float64) {
	return c.pov.Loc.X, c.pov.Loc.Y, c.pov.Loc.Z
}
func (c *camera) SetLocation(x, y, z float64) {
	c.pov.Loc.X, c.pov.Loc.Y, c.pov.Loc.Z = x, y, z
	c.updateViewTransform()
}
func (c *camera) Move(x, y, z float64) {
	c.pov.Move(x, y, z)
	c.updateViewTransform()
}
func (c *camera) Spin(x, y, z float64) {
	c.pov.Spin(x, y, z)
	c.updateViewTransform()
}
func (c *camera) Tilt() (up float64) { return c.up }
func (c *camera) SetTilt(up float64) {
	c.up = up
	c.transform(c.vm)
	// FUTURE: calculate inverse view calculation for views with tilt.
}

// Distance returns the distance squared of the camera to the given point.
func (c *camera) Distance(px, py, pz float64) float64 {
	dx := px - c.Loc.X
	dy := py - c.Loc.Y
	dz := pz - c.Loc.Z
	return float64(dx*dx + dy*dy + dz*dz)
}

// SetPerspective sets the scene to use a 3D perspective
func (c *camera) SetPerspective(fov, ratio, near, far float64) {
	c.pm.Persp(fov, ratio, near, far)
	c.ipm.PerspInv(fov, ratio, near, far)
	c.transform(c.vm)
	ivp(c, c.ivm)
}

// SetOrthographic sets the scene to use a 2D orthographic perspective.
func (c *camera) SetOrthographic(left, right, bottom, top, near, far float64) {
	c.pm.Ortho(left, right, bottom, top, near, far)
	c.transform(c.vm)

	// Inverse matrix currently ignored for Orthographic. Ortho views are
	// expected to match the screen pixel sizes.
	c.ipm.Set(lin.M4I)
}

// updateViewTransform ensures that the view and inverse-view transform are
// kept upto date each time the camera moves. Calculating once per move should
// be quicker than calculating later for each object in the scene.
func (c *camera) updateViewTransform() {
	c.transform(c.vm) // view transform.
	ivp(c, c.ivm)     // inverse view transform.
}

// Ray applies inverse transforms to derive world space coordinates for
// a ray projected from the camera through the mouse's screen position. See:
//     http://bookofhook.com/mousepick.pdf
//     http://antongerdelan.net/opengl/raycasting.html
//     http://schabby.de/picking-opengl-ray-tracing/
//     (opengl FAQ Picking 20.0.010)
//     http://www.opengl.org/archives/resources/faq/technical/selection.htm
//     http://www.codeproject.com/Articles/625787/Pick-Selection-with-OpenGL-and-OpenCL
func (c *camera) Ray(mx, my, ww, wh int) (x, y, z float64) {
	c.ray.SetS(0, 0, 0)
	if mx >= 0 && mx <= ww && my >= 0 && my <= wh {
		clipx := float64(2*mx)/float64(ww) - 1 // mx to range -1:1
		clipy := float64(2*my)/float64(wh) - 1 // my to range -1:1
		clip := c.v0.SetS(clipx, clipy, -1, 1)

		// use the inverse perspective to go from clip to eye (view) coordinates
		eye := clip.MultvM(clip, c.ipm)
		eye.Z = -1 // into the screen
		eye.W = 0  // want a vector, not a point

		// use the inverse view to go from eye (view) coordinates to world coordinates.
		world := eye.MultvM(eye, c.ivm)
		c.ray.SetS(world.X, world.Y, world.Z) // ignore the W component.
		c.ray.Unit()                          // ensure that a unit vector is returned.
	}
	return c.ray.X, c.ray.Y, c.ray.Z
}

// Screen applies the camera transform on a 3D point in world space wx, wy, wz
// and returns the 2D screen coordinate sx, sy. The window width and height
// ww, wh are also needed. Essentially the reverse of the Ray method and
// duplicates what is done in the rendering pipeline.
func (c *camera) Screen(wx, wy, wz float64, ww, wh int) (sx, sy int) {
	vec := c.v0.SetS(wx, wy, wz, 1)
	vec.MultvM(vec, c.vm)
	vec.MultvM(vec, c.pm)
	clipx := vec.X/vec.W + 1 // range -1:1 to 0:2
	clipy := vec.Y/vec.W + 1 // range -1:1 to 0:2
	sx = int(lin.Round(clipx*0.5*float64(ww), 0))
	sy = int(lin.Round(clipy*0.5*float64(wh), 0))
	return
}

// view
// ===========================================================================
// view transforms

// View transforms create a view matrix based on the point of view of the
// camera. Since the camera is always at 0, 0, 0 then moving forward by
// x:units really means moving the world (everything else) back -x:units.
// Same with rotating one way by x:degrees really means rotating the world
// -x:degrees.
type viewTransform func(*camera, *lin.M4) *lin.M4

// vp view transform implementation.
func vp(c *camera, vm *lin.M4) *lin.M4 {
	vm.SetQ(c.Rot)
	return vm.TranslateTM(-c.Loc.X, -c.Loc.Y, -c.Loc.Z)
}

// vo view transform implementation.
func vo(c *camera, vm *lin.M4) *lin.M4 {
	return vm.Set(lin.M4I).ScaleMS(1, 1, 0)
}

// vf view transform implementation.
func vf(c *camera, vm *lin.M4) *lin.M4 {
	rot := c.q0.SetAa(1, 0, 0, lin.Rad(-c.up))
	rot.Mult(rot, c.Rot)
	return vm.SetQ(rot).TranslateTM(-c.Loc.X, -c.Loc.Y, -c.Loc.Z)
}

// xz_xy view transform implementation.
func xz_xy(c *camera, vm *lin.M4) *lin.M4 {
	rot := c.q0.SetAa(1, 0, 0, -lin.Rad(90))
	return vm.SetQ(rot).ScaleMS(1, 1, 0).TranslateTM(-c.Loc.X, -c.Loc.Y, -c.Loc.Z)
}

// inverse vp view transform. Experimental for ray casting... only one view
// inverse for now, need better design to incorporate more if needed.
func ivp(c *camera, vm *lin.M4) *lin.M4 {
	rot := c.q0.Inv(c.Rot)
	vm.SetQ(rot)
	return vm.TranslateMT(c.Loc.X, c.Loc.Y, c.Loc.Z)
}
