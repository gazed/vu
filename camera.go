// Copyright © 2014-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// camera.go encapsulates view and projection matricies needed for rendering.
// DESIGN:
//   See the following for a first person camera example using quaternions:
//   http://content.gpwiki.org/index.php/OpenGL:Tutorials:Using_Quaternions_to_represent_rotation
//   One way to implement cameras:
//   http://udn.epicgames.com/Three/CameraTechnicalGuide.html
// For mouse-picking or ray-casting See:
//   http://bookofhook.com/mousepick.pdf
//   http://antongerdelan.net/opengl/raycasting.html
//   http://schabby.de/picking-opengl-ray-tracing/
//   (opengl FAQ Picking 20.0.010)
//   http://www.opengl.org/archives/resources/faq/technical/selection.htm
//   http://www.codeproject.com/Articles/625787/Pick-Selection-with-OpenGL-and-OpenCL
//
// FUTURE: Combine and share the Pov and Camera transform knowledge.
//         Ideally this makes the API easier and reduce engine code complexity
//         by removing a set of transform code.
// FUTURE: Make consistent. Look is an accessor Lookat() is a method.

import (
	"github.com/gazed/vu/math/lin"
)

// Camera makes rendered models visible within a frame. A camera is
// associated with a scene where it is used to render the scenes models.
//
// Camera combines a location+orientation using separate up/down angle
// tracking. This allows use as a first-person camera which can limit
// up/down to a given range, often 180deg. Overall orientation is calculated
// by combining Pitch and Yaw. Look is for walking cameras, Lookat is for
// flying cameras.
type Camera struct {
	Pitch float64 // X-axis rotation in degrees. Set using SetPitch.
	Yaw   float64 // Y-axis rotation in degrees. Set using SetYaw.
	Look  *lin.Q  // Y-axis quaternion rotation updated by SetYaw.

	// Position and orientation.
	at   *lin.T // Combined location and Pitch/Yaw orientation.
	prev *lin.T // Previous location and orientation.
	xrot *lin.Q // X-axis rotation: from Pitch.

	// Camera perspective values used to set projection matricies.
	near, far float64 // Frustum clip set by application.
	fov       float64 // Field of view set by application.
	focus     bool    // true if the camera projection needs setting.

	// View transfrom algorithms: 3D perspective uses different
	// transforms than 2D/3D orthographic. Accounts for camera/eye.
	// This affects the view portion of model-view-projection.
	vt viewTransform    // View transform matrix generator.
	it inverseTransform // Inverse view transform matrix generator.

	// Track the view, projection matricies and their inverses.
	// The inverses are needed for Ray casting.
	vm  *lin.M4 // View part of MVP matrix.
	ivm *lin.M4 // Inverse view matrix.
	pm  *lin.M4 // Projection part of MVP matrix.
	ipm *lin.M4 // Inverse projection matrix.

	// Scratch variables needed each update.
	q0  *lin.Q  // Scratch for camera transform calculations.
	v0  *lin.V4 // Scratch for pick ray calculations.
	ray *lin.V3 // Scratch for pick ray calculations.
}

// newCamera creates a default rendering field that is looking down
// the positive Z axis with positive Y up.
func newCamera() *Camera {
	c := &Camera{fov: 60, focus: true} // Default fov.
	c.Look = lin.NewQ().SetAa(0, 1, 0, 0)
	c.at = lin.NewT()
	c.prev = lin.NewT()
	c.xrot = lin.NewQ().SetAa(1, 0, 0, 0)
	c.vm = &lin.M4{}
	c.ivm = (&lin.M4{}).Set(lin.M4I)
	c.pm = &lin.M4{}
	c.ipm = &lin.M4{}
	c.q0 = &lin.Q{}
	c.v0 = &lin.V4{}
	c.ray = &lin.V3{}
	c.vt, c.it = vp, ivp
	c.vt(c.at, c.q0, c.vm)  // initial view transform
	c.it(c.at, c.q0, c.ivm) // inverse view transform.
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

// Lookat returns an orientation which is good for flying around.
// It is a combination of Pitch and Yaw.
func (c *Camera) Lookat() *lin.Q { return c.at.Rot }

// SetPitch sets the rotation around the X axis and updates
// the Look direction. The camera instance is returned.
func (c *Camera) SetPitch(deg float64) *Camera {
	c.Pitch = deg
	c.xrot.SetAa(1, 0, 0, lin.Rad(c.Pitch))
	c.at.Rot.Mult(c.xrot, c.Look)
	return c
}

// SetYaw sets the rotation around the Y axis and updates the
// Look and Lookat directions. The camera instance is returned.
func (c *Camera) SetYaw(deg float64) *Camera {
	c.Yaw = deg
	c.Look.SetAa(0, 1, 0, lin.Rad(c.Yaw))
	c.at.Rot.Mult(c.xrot, c.Look)
	return c
}

// Distance returns the distance squared of the camera to the given Pov.
// Uses the existing Pov world coordinates.
func (c *Camera) Distance(wx, wy, wz float64) float64 {
	dx := wx - c.at.Loc.X
	dy := wy - c.at.Loc.Y
	dz := wz - c.at.Loc.Z
	return float64(dx*dx + dy*dy + dz*dz)
}

// normalizedDistance returns the distance to camera as a 0-1 range.
// This is needed for sorting transparent objects. Expected to be
// called on objects that have already passed culling.
func (c *Camera) normalizedDistance(wx, wy, wz float64) float64 {
	vec := c.v0.SetS(wx, wy, wz, 1)
	vec.MultvM(vec, c.vm)          // apply view matrix.
	vec.MultvM(vec, c.pm)          // apply projection matrix.
	clipz := vec.Z*0.5/vec.W + 0.5 // convert from -1:1 to 0:1
	if clipz < 0 || clipz > 1 {
		return -1 // outside near/far planes.
	}
	return clipz
}

// setPerspective makes the camera use a 3D projection.
// This is the projection part of model-view-projection.
func (c *Camera) setPerspective(fov, ratio, near, far float64) {
	c.pm.Persp(fov, ratio, near, far)
	c.ipm.PerspInv(fov, ratio, near, far)
}

// setOrthographic makes the camera use a 2D projection.
// This is the projection part of model-view-projection.
func (c *Camera) setOrthographic(left, right, bottom, top, near, far float64) {
	c.pm.Ortho(left, right, bottom, top, near, far)
}

// Ray applies inverse transforms to derive world space coordinates
// for a ray projected from the camera through the mouse's mx,my
// screen position given window width and height ww,wh.
func (c *Camera) Ray(mx, my, ww, wh int) (x, y, z float64) {
	c.ray.SetS(0, 0, 0)
	if mx >= 0 && mx <= ww && my >= 0 && my <= wh {
		clipx := float64(2*mx)/float64(ww) - 1 // mx to range -1:1
		clipy := float64(2*my)/float64(wh) - 1 // my to range -1:1
		clip := c.v0.SetS(clipx, clipy, -1, 1)

		// Use inverse perspective to go from clip to eye (view) coordinates.
		eye := clip.MultvM(clip, c.ipm)
		eye.Z = -1 // into the screen
		eye.W = 0  // want a vector, not a point

		// Use inverse view to go from eye (view) to world coordinates.
		world := eye.MultvM(eye, c.ivm)
		c.ray.SetS(world.X, world.Y, world.Z) // ignore the W component.
		c.ray.Unit()                          // return a unit vector.
	}
	return c.ray.X, c.ray.Y, c.ray.Z
}

// Screen applies the camera transform on a 3D point in world space wx,wy,wz
// and returns the 2D screen coordinate sx,sy. The window width and height
// ww,wh are needed. Essentially the reverse of the Ray method and duplicates
// what is done in the rendering pipeline.
// Returns -1,-1 if the point is outside the screen area.
func (c *Camera) Screen(wx, wy, wz float64, ww, wh int) (sx, sy int) {
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

// Camera
// ===========================================================================
// view transforms

// viewTransform creates a transform matrix from location and orientation.
// This is expected to be used for camera transforms. The camera is thought
// of as being at 0,0,0. This means moving the camera forward by x:units
// really means moving the world (everything else) back -x:units. Likewise
// rotating the camera by x:degrees really means rotating the world by -x.
type viewTransform func(*lin.T, *lin.Q, *lin.M4)

// vp perspective projection transform used for Camera.Vt.
func vp(at *lin.T, scr *lin.Q, vm *lin.M4) {
	vm.SetQ(at.Rot)
	vm.TranslateTM(-at.Loc.X, -at.Loc.Y, -at.Loc.Z)
}

// xzxy perspective to ortho view transform used for Camera.Vt.
// Can help transform a 3D map to a 2D overlay.
func xzxy(at *lin.T, scr *lin.Q, vm *lin.M4) {
	rot := scr.SetAa(1, 0, 0, -lin.Rad(90))
	l := at.Loc
	vm.SetQ(rot).ScaleMS(1, 1, 0).TranslateTM(-l.X, -l.Y, -l.Z)
}

// vo orthographic projection transform used for Camera.Vt.
func vo(pov *lin.T, scr *lin.Q, vm *lin.M4) {
	vm.Set(lin.M4I).ScaleMS(1, 1, 0)
}

// inverseTransform creates the inverse transform matrix from
// the location and orientation. Used in ray picking.
type inverseTransform func(*lin.T, *lin.Q, *lin.M4)

// ivp inverse view transform. For ray casting.
func ivp(at *lin.T, scr *lin.Q, vm *lin.M4) {
	rot := scr.Inv(at.Rot)
	vm.SetQ(rot)
	vm.TranslateMT(at.Loc.X, at.Loc.Y, at.Loc.Z)
}

// nv is a null identity view.
func nv(at *lin.T, scr *lin.Q, vm *lin.M4) {
	vm.Set(lin.M4I)
}
