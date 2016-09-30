// Copyright © 2014-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// camera.go encapsulates view and projection matricies needed for rendering.
// DESIGN:
//   See the following for a first person camera example using quaternions:
//   http://content.gpwiki.org/index.php/OpenGL:Tutorials:Using_Quaternions_to_represent_rotation
//   One way to implement cameras:
//   http://udn.epicgames.com/Three/CameraTechnicalGuide.html
// For mousepicking See:
//   http://bookofhook.com/mousepick.pdf
//   http://antongerdelan.net/opengl/raycasting.html
//   http://schabby.de/picking-opengl-ray-tracing/
//   (opengl FAQ Picking 20.0.010)
//   http://www.opengl.org/archives/resources/faq/technical/selection.htm
//   http://www.codeproject.com/Articles/625787/Pick-Selection-with-OpenGL-and-OpenCL
// FUTURE: attach a Camera to a top level Pov and control its location, orientation
//   through a separate camera controller? Will this make the API easier and
//   reduce engine code complexity by removing a separate camera transform?
//   Assign models to a camera until the next camera is found in
//   the pov hierarchy traversal.

import (
	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/render"
)

// Camera makes rendered models visible within a frame. A camera is
// attached to a point of view (Pov) where it renders all models in that
// Pov's hierarchy. Cameras location and orientation is independent from
// the Pov allowing a camera to be positioned independently from the models.
//
// Camera combines a location+orientation using separate up/down angle
// tracking. This allows use as a first-person camera which can limit
// up/down to a given range, often 180deg. Overall orientation is calculated
// by combining Pitch and Yaw. Look is for walking cameras, Lookat is for
// flying cameras.
type Camera struct {
	Pitch   float64 // X-axis rotation in degrees.
	Yaw     float64 // Y-axis rotation in degrees.
	Look    *lin.Q  // Y-axis quaternion rotation from Yaw.
	Depth   bool    // True by default for 3D depth processing.
	Overlay int     // Set render order bucket with OVERLAY or greater.

	// SetCull sets a method that reduces the number of Models rendered
	// each update. It can be application supplied or engine supplied.
	Cull Culler // Set by application, ie: c.Cull = vu.NewFacingCuller.

	// Set one of the possible view transfrom algorithms. This affects
	// the view portion of model-view-projection.
	Vt ViewTransform // Assigned view transform matrix generator.

	// Internal values.
	at     *lin.T // Combined Pitch/Yaw orientation.
	xrot   *lin.Q // X-axis rotation: from Pitch.
	target uint32 // render layer target. Default 0.

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
func newCamera() *Camera {
	c := &Camera{Depth: true}
	c.Vt = VP
	c.Look = lin.NewQ().SetAa(0, 1, 0, 0)
	c.at = lin.NewT()
	c.xrot = lin.NewQ().SetAa(1, 0, 0, 0)
	c.vm = &lin.M4{}
	c.ivm = (&lin.M4{}).Set(lin.M4I)
	c.pm = &lin.M4{}
	c.ipm = &lin.M4{}
	c.q0 = &lin.Q{}
	c.v0 = &lin.V4{}
	c.ray = &lin.V3{}
	return c
}

// transform applies the view transform to the scene camera
// and returns the result. The input matrix is not changed.
func (c *Camera) transform(vm *lin.M4) *lin.M4 { return c.Vt(c.at, c.q0, vm) }

// isCulled applies the camera cull algorithm to the given location.
func (c *Camera) isCulled(px, py, pz float64) bool {
	if c.Cull != nil {
		return c.Cull.Culled(c, px, py, pz)
	}
	return false
}

// updateTransform ensures that the view and inverse-view transform are
// kept in sync each time the camera moves. Calculating once per move should
// be quicker than calculating later for each object in the scene.
func (c *Camera) updateTransform() {
	c.transform(c.vm)            // view transform.
	ivp(c.at, c.qx, c.q0, c.ivm) // inverse view transform.
}

// At returns the cameras current location in world space.
func (c *Camera) At() (x, y, z float64) {
	return c.at.Loc.GetS()
}

// SetAt positions the camera in world space.
func (c *Camera) SetAt(x, y, z float64) {
	c.at.Loc.SetS(x, y, z)
	c.updateTransform()
}

// Move adjusts the camera location relative to the given orientation.
// For orientation, use Lookat() to fly, use Look to run along XZ.
func (c *Camera) Move(x, y, z float64, q *lin.Q) {
	dx, dy, dz := lin.MultSQ(x, y, z, q)
	c.at.Loc.X += dx
	c.at.Loc.Y += dy
	c.at.Loc.Z += dz
	c.updateTransform()
}

// Lookat returns an orientation which is good for flying around.
// It is a combination of Pitch and Yaw.
func (c *Camera) Lookat() *lin.Q { return c.at.Rot }

// SetPitch sets the degrees of rotation around the X axis.
func (c *Camera) SetPitch(deg float64) {
	c.Pitch = deg
	c.xrot.SetAa(1, 0, 0, lin.Rad(c.Pitch))
	c.at.Rot.Mult(c.xrot, c.Look)
	c.updateTransform()
}

// SetYaw sets the rotation around the Y axis.
func (c *Camera) SetYaw(deg float64) {
	c.Yaw = deg
	c.Look.SetAa(0, 1, 0, lin.Rad(c.Yaw))
	c.at.Rot.Mult(c.xrot, c.Look)
	c.updateTransform()
}

// Distance returns the distance squared of the camera to the given point.
func (c *Camera) Distance(px, py, pz float64) float64 {
	dx := px - c.at.Loc.X
	dy := py - c.at.Loc.Y
	dz := pz - c.at.Loc.Z
	return float64(dx*dx + dy*dy + dz*dz)
}

// SetUI configures the camera to be 2D: no depth, drawn last.
func (c *Camera) SetUI() *Camera {
	c.Overlay = render.Overlay // Draw last.
	c.Depth = false            // 2D rendering.
	c.Vt = VO                  // orthographic view transform.
	return c
}

// SetPerspective makes the camera use a 3D projection.
// This is the projection part of model-view-projection.
func (c *Camera) SetPerspective(fov, ratio, near, far float64) {
	c.pm.Persp(fov, ratio, near, far)
	c.ipm.PerspInv(fov, ratio, near, far)
	c.updateTransform()
}

// SetOrthographic makes the camera use a 2D projection.
// This is the projection part of model-view-projection.
func (c *Camera) SetOrthographic(left, right, bottom, top, near, far float64) {
	c.pm.Ortho(left, right, bottom, top, near, far)
	c.transform(c.vm)

	// Inverse matrix currently ignored for Orthographic.
	// Ortho views are expected to match the screen pixel sizes.
	c.ipm.Set(lin.M4I)
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

// camera
// ===========================================================================
// view transforms

// ViewTransform creates a transform matrix from location and orientation.
// This is expected to be used for camera transforms. The camera is thought
// of as being at 0,0,0. This means moving the camera forward by x:units
// really means moving the world (everything else) back -x:units. Likewise
// rotating the camera by x:degrees really means rotating the world by -x.
type ViewTransform func(*lin.T, *lin.Q, *lin.M4) *lin.M4

// VP perspective projection transform used for Camera.Vt.
func VP(at *lin.T, scr *lin.Q, vm *lin.M4) *lin.M4 {
	vm.SetQ(at.Rot)
	return vm.TranslateTM(-at.Loc.X, -at.Loc.Y, -at.Loc.Z)
}

// VO orthographic projection transform used for Camera.Vt.
func VO(pov *lin.T, scr *lin.Q, vm *lin.M4) *lin.M4 {
	return vm.Set(lin.M4I).ScaleMS(1, 1, 0)
}

// XzXy perspective to ortho view transform used for Camera.Vt.
// Can help transform a 3D map to a 2D overlay.
func XzXy(at *lin.T, scr *lin.Q, vm *lin.M4) *lin.M4 {
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

// =============================================================================

// cams manages all the active Camera instances.
// There's not many cameras so not much to optimize.
type cams struct {
	data map[eid]*Camera // Camera instance data.
}

// newCams creates the camera component manager and is expected to
// be called once on startup.
func newCams() *cams { return &cams{data: map[eid]*Camera{}} }

// get returns the Camera associated with the given entity or nil
// if there is no camera.
func (cs *cams) get(id eid) *Camera {
	if cam, ok := cs.data[id]; ok {
		return cam
	}
	return nil
}

// create makes a new camera and associates it with the given entity.
// If there already is a camera for the given entity, nothing is created
// and the existing camera is returned.
func (cs *cams) create(id eid) *Camera {
	if c, ok := cs.data[id]; ok {
		return c // Don't allow creating over existing camera.
	}
	c := newCamera()
	cs.data[id] = c
	return c
}

// dispose removes the camera associated with the given entity.
// Nothing happens if there is no camera.
func (cs *cams) dispose(id eid) { delete(cs.data, id) }
