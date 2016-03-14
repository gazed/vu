// Copyright © 2014-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// Design notes:
// See the following for a first person camera example using quaternions:
// http://content.gpwiki.org/index.php/OpenGL:Tutorials:Using_Quaternions_to_represent_rotation
// One way to implement cameras:
// http://udn.epicgames.com/Three/CameraTechnicalGuide.html

import (
	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/render"
)

// Camera is necessary to render models. A camera is attached to a point
// of view (Pov) where it renders all models in that Pov's hierarchy.
// Camera tracks the location and orientation of a camera as well as an
// associated projection transform. Keeping its location and orientation
// separate from the transform hierarchy allows a camera to be positioned
// independently from the models.
type Camera interface {
	Location() (x, y, z float64)    // Get, or
	SetLocation(x, y, z float64)    // ...Set the camera location.
	Move(x, y, z float64, q *lin.Q) // Adjust location along orientation.

	// Orientation is calculated from pitch and yaw. Lookat can be used
	// for flying cameras. Lookxz is good for walking cameras.
	Lookat() *lin.Q          // Get the XYZ view orientation.
	Lookxz() *lin.Q          // Get quaternion rotation about Y.
	Pitch() (deg float64)    // Looking up/down. Get or...
	SetPitch(deg float64)    // ...Set the X rotation in degrees,
	AdjustPitch(deg float64) // ...adjust rotation around X axis.
	Yaw() (deg float64)      // Spinning around. Get or...
	SetYaw(deg float64)      // ...Set the Y rotation in degrees,
	AdjustYaw(deg float64)   // ...adjust rotation around Y axis.

	// SetCull sets a method that reduces the number of Models rendered
	// each update. It can be application supplied or engine supplied
	// ie: NewFacingCuller.
	SetCull(c Cull)        // Set to nil to turn off culling.
	SetDepth(enabled bool) // True for 3D camera. 2D cams ignore depth.
	SetLast(index int)     // For sequencing UI cameras. Higher is later.
	SetUI()                // UI camera: 2D, no depth, drawn last.

	// Set one of the possible view transfrom algorithms. This affects
	// the view portion of model-view-projection.
	SetView(vt ViewTransform) // Update the view and inverse view.

	// Use one of the following to create a projection transform.
	// This is the projection part of model-view-projection.
	SetPerspective(fov, ratio, near, far float64)                // 3D.
	SetOrthographic(left, right, bottom, top, near, far float64) // 2D.

	// Ray applies inverse transforms to derive world space coordinates for
	// a ray projected from the camera through the mouse's mx,my screen
	// position given window width and height ww,wh.
	Ray(mx, my, ww, wh int) (x, y, z float64)

	// Screen calculates screen coordinates sx,sy for world coordinates
	// wx,wy,wz and window width and height ww,wh.
	Screen(wx, wy, wz float64, ww, wh int) (sx, sy int)

	// Distance returns the distance squared of the camera to the given point.
	Distance(px, py, pz float64) float64
}

// Camera
// ===========================================================================
// camera implements Camera

// camera combines a location+direction (Pov) with a separate up/down angle
// tracking. This allows use as a FPS camera which can limit up/down to
// a given range (often 180).
type camera struct {
	at      *lin.T        // Location/direction and Y spin.
	xdeg    float64       // X-axis rotation in degrees.
	ydeg    float64       // Y-axis rotation in degrees.
	xrot    *lin.Q        // X-axis rotation: pitch.
	yrot    *lin.Q        // Y-axis rotation: yaw.
	vt      ViewTransform // Assigned view transform matrix generator.
	depth   bool          // True for 3D depth processing.
	cull    Cull          // Set by application.
	overlay int           // Set render bucket with OVERLAY or greater.
	target  uint32        // render layer target. Default 0.

	// Track the view, projection matricies and their inverses.
	vm  *lin.M4 // View part of MVP matrix.
	ivm *lin.M4 // Inverse view matrix.
	pm  *lin.M4 // Projection part of MVP matrix.
	ipm *lin.M4 // Inverse projection matrix.

	// Scratch variables needed each update.
	q0  *lin.Q  // Scratch for camera transform calculations.
	qx  *lin.Q  // Scratch for camera transform calculations.
	v0  *lin.V4 // Scratch for pick ray calculations.
	ray *lin.V3 // Scratch for pick ray calculations.
}

// newCamera creates a default rendering field that is looking down
// the positive Z axis with positive Y up.
func newCamera() *camera {
	c := &camera{depth: true}
	c.vt = VP
	c.at = lin.NewT()
	c.vm = &lin.M4{}
	c.ivm = (&lin.M4{}).Set(lin.M4I)
	c.pm = &lin.M4{}
	c.ipm = &lin.M4{}
	c.q0 = &lin.Q{}
	c.xrot = lin.NewQ().SetAa(1, 0, 0, 0)
	c.yrot = lin.NewQ().SetAa(0, 1, 0, 0)
	c.v0 = &lin.V4{}
	c.ray = &lin.V3{}
	return c
}

// SetView accepts the given camera ViewTransform.
func (c *camera) SetView(vt ViewTransform) { c.vt = vt }

// transform applies the view transform to the scene camera
// and returns the result. The input matrix is not changed.
func (c *camera) transform(vm *lin.M4) *lin.M4 { return c.vt(c.at, c.q0, vm) }

// isCulled applies the camera cull algorithm to the given location.
func (c *camera) isCulled(px, py, pz float64) bool {
	if c.cull != nil {
		return c.cull.Culled(c, px, py, pz)
	}
	return false
}

// updateTransform ensures that the view and inverse-view transform are
// kept in sync each time the camera moves. Calculating once per move should
// be quicker than calculating later for each object in the scene.
func (c *camera) updateTransform() {
	c.transform(c.vm)            // view transform.
	ivp(c.at, c.qx, c.q0, c.ivm) // inverse view transform.
}
func (c *camera) Location() (x, y, z float64) {
	return c.at.Loc.X, c.at.Loc.Y, c.at.Loc.Z
}
func (c *camera) SetLocation(x, y, z float64) {
	c.at.Loc.X, c.at.Loc.Y, c.at.Loc.Z = x, y, z
	c.updateTransform()
}

// Lookat returns a direction good for flying around.
func (c *camera) Lookat() *lin.Q { return c.at.Rot }

// Lookxz returns a direction that works for walking around.
func (c *camera) Lookxz() *lin.Q { return c.yrot }

// Move relative to the given orientation.
// Use Lookat() to fly. Use Lookxz() to run along XZ.
func (c *camera) Move(x, y, z float64, q *lin.Q) {
	dx, dy, dz := lin.MultSQ(x, y, z, q)
	c.at.Loc.X += dx
	c.at.Loc.Y += dy
	c.at.Loc.Z += dz
	c.updateTransform()
}

// Pitch gets the degrees of rotation around the X axis.
func (c *camera) Pitch() (deg float64) { return c.xdeg }
func (c *camera) SetPitch(deg float64) {
	c.xdeg = deg
	c.xrot.SetAa(1, 0, 0, lin.Rad(c.xdeg))
	c.at.Rot.Mult(c.xrot, c.yrot)
	c.updateTransform()
}

// AdjustPitch updates the cameras rotation about the X axis
// as well as the overall camera orientation.
func (c *camera) AdjustPitch(deg float64) {
	c.SetPitch(c.xdeg + deg)
}

// Yaw gets the rotation around the Y axis.
func (c *camera) Yaw() (deg float64) { return c.ydeg }
func (c *camera) SetYaw(deg float64) {
	c.ydeg = deg
	c.yrot.SetAa(0, 1, 0, lin.Rad(c.ydeg))
	c.at.Rot.Mult(c.xrot, c.yrot)
	c.updateTransform()
}

// AdjustYaw updates the cameras rotation about the Y axis
// as well as the overall camera orientation.
func (c *camera) AdjustYaw(deg float64) {
	c.SetYaw(c.ydeg + deg)
}

// Distance returns the distance squared of the camera to the given point.
func (c *camera) Distance(px, py, pz float64) float64 {
	dx := px - c.at.Loc.X
	dy := py - c.at.Loc.Y
	dz := pz - c.at.Loc.Z
	return float64(dx*dx + dy*dy + dz*dz)
}

// Implement Camera interface.
func (c *camera) SetDepth(enabled bool) { c.depth = enabled }
func (c *camera) SetCull(cull Cull)     { c.cull = cull }
func (c *camera) SetLast(index int)     { c.overlay = render.OVERLAY + index }
func (c *camera) SetUI() {
	c.overlay = render.OVERLAY // Draw last.
	c.depth = false            // 2D rendering.
	c.SetView(VO)              // orthographic view transform.
}

// SetPerspective makes the camera use a 3D projection.
func (c *camera) SetPerspective(fov, ratio, near, far float64) {
	c.pm.Persp(fov, ratio, near, far)
	c.ipm.PerspInv(fov, ratio, near, far)
	c.updateTransform()
}

// SetOrthographic makes the camera use a 2D projection.
func (c *camera) SetOrthographic(left, right, bottom, top, near, far float64) {
	c.pm.Ortho(left, right, bottom, top, near, far)
	c.transform(c.vm)

	// Inverse matrix currently ignored for Orthographic.
	// Ortho views are expected to match the screen pixel sizes.
	c.ipm.Set(lin.M4I)
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

		// Use the inverse perspective to go from clip to eye (view) coordinates.
		eye := clip.MultvM(clip, c.ipm)
		eye.Z = -1 // into the screen
		eye.W = 0  // want a vector, not a point

		// Use the inverse view to go from eye (view) coordinates to world coordinates.
		world := eye.MultvM(eye, c.ivm)
		c.ray.SetS(world.X, world.Y, world.Z) // ignore the W component.
		c.ray.Unit()                          // ensure that a unit vector is returned.
	}
	return c.ray.X, c.ray.Y, c.ray.Z
}

// Screen applies the camera transform on a 3D point in world space wx,wy,wz
// and returns the 2D screen coordinate sx,sy. The window width and height
// ww,wh are also needed. Essentially the reverse of the Ray method and
// duplicating what is done in the rendering pipeline.
func (c *camera) Screen(wx, wy, wz float64, ww, wh int) (sx, sy int) {
	vec := c.v0.SetS(wx, wy, wz, 1)
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

// camera
// ===========================================================================
// view transforms

// ViewTransform creates a transform matrix from location and orientation.
// This is expected to be used for camera transforms. The camera is thought
// of as being at 0,0,0. This means moving the camera forward by x:units
// really means moving the world (everything else) back -x:units. Likewise
// rotating the camera by x:degrees really means rotating the world by -x.
type ViewTransform func(*lin.T, *lin.Q, *lin.M4) *lin.M4

// VP perspective projection transform.
func VP(at *lin.T, scr *lin.Q, vm *lin.M4) *lin.M4 {
	vm.SetQ(at.Rot)
	return vm.TranslateTM(-at.Loc.X, -at.Loc.Y, -at.Loc.Z)
}

// VO orthographic projection transform.
func VO(pov *lin.T, scr *lin.Q, vm *lin.M4) *lin.M4 {
	return vm.Set(lin.M4I).ScaleMS(1, 1, 0)
}

// XZ_XY perspective to ortho view transform.
// Can help transform a 3D map to a 2D overlay.
func XZ_XY(at *lin.T, scr *lin.Q, vm *lin.M4) *lin.M4 {
	rot := scr.SetAa(1, 0, 0, -lin.Rad(90))
	l := at.Loc
	return vm.SetQ(rot).ScaleMS(1, 1, 0).TranslateTM(-l.X, -l.Y, -l.Z)
}

// inverse vp view transform. For ray casting... only one view
// inverse for now. Need better design if more are needed.
func ivp(at *lin.T, xrot, scr *lin.Q, vm *lin.M4) *lin.M4 {
	rot := scr.Inv(at.Rot)
	vm.SetQ(rot)
	return vm.TranslateMT(at.Loc.X, at.Loc.Y, at.Loc.Z)
}
