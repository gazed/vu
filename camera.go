// Copyright © 2014-2024 Galvanized Logic Inc.

package vu

// camera.go holds the view and projection matricies needed for rendering.

import (
	"math"

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
}

// newCamera creates a default rendering field that is looking
// down the negative Z axis with positive Y up.
func newCamera() *Camera {
	c := &Camera{fov: 90, focus: true} // Default fov.
	c.Look = lin.NewQ().SetAa(0, 1, 0, 0)
	c.at = lin.NewT()
	c.prev = lin.NewT()
	c.xrot = lin.NewQ().SetAa(1, 0, 0, 0)
	c.vm = &lin.M4{}
	c.ivm = (&lin.M4{}).Set(lin.M4I)
	c.pm = &lin.M4{}
	c.ipm = &lin.M4{}
	c.vt, c.it = vp, ivp
	c.vt(c.at, c.vm)  // initial view transform
	c.it(c.at, c.ivm) // inverse view transform.
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

// Ray applies inverse transforms to derive world space coordinates
// for a ray projected from the camera through the mouse's mx,my
// screen position given window width and height ww,wh.
func (c *Camera) Ray(mx, my, ww, wh int) (x, y, z float64) {
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
	}
	return ray.X, ray.Y, ray.Z
}

// RayCastSphere calculates the point of collision between a ray
// originating from the camera and a sphere in world space.
// The closest contact point is returned if there is an intersection.
// Expecting the ray to be a unit vector, see: camera.Ray().
//
//	http://en.wikipedia.org/wiki/Line–sphere_intersection
func (c *Camera) RayCastSphere(ray, sphere *lin.V3, radius float64) (hit bool, x, y, z float64) {
	rx, ry, rz := c.At() // ray origin is the camera world location.

	// vector from ray origin to sphere center
	rs := lin.NewV3().SetS(sphere.X-rx, sphere.Y-ry, sphere.Z-rz)

	// distance between the center of the sphere and the ray.
	// If the distance is larger than the radius there is no intersection.
	d0 := ray.Dot(rs)
	if d0 < 0 {
		return false, 0, 0, 0 // no hit
	}
	r2 := radius * radius
	d1 := rs.Dot(rs) - d0*d0
	if d1 > r2 {
		return false, 0, 0, 0 // no hit
	}

	// Get contact point by scaling the ray direction with
	// the contact distance and adding the ray origin.
	dlen := d0 - math.Sqrt(r2-d1)
	x, y, z = rx+dlen*ray.X, ry+dlen*ray.Y, rz+dlen*ray.Z
	return true, x, y, z
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

// Camera
// ===========================================================================
// view transforms

// viewTransform creates a transform matrix from location and orientation.
// This is expected to be used for camera transforms. The camera is thought
// of as being at 0,0,0. Moving the camera forward by x:units really means
// moving the world (everything else) back -x:units. Likewise rotating the
// camera by x:degrees really means rotating the world by -x.
type viewTransform func(*lin.T, *lin.M4)

// vp perspective projection transform used for Camera.Vt.
func vp(at *lin.T, vm *lin.M4) {
	vm.SetQ(at.Rot)
	vm.TranslateTM(-at.Loc.X, -at.Loc.Y, -at.Loc.Z)
}

// xzxy perspective to ortho view transform used for Camera.Vt.
// Can help transform a 3D map to a 2D overlay.
func xzxy(at *lin.T, vm *lin.M4) {
	l := at.Loc
	rot := lin.NewQ().SetAa(1, 0, 0, -lin.Rad(90))
	vm.SetQ(rot).ScaleMS(1, 1, 0).TranslateTM(-l.X, -l.Y, -l.Z)
}

// vo orthographic projection transform used for Camera.Vt.
func vo(pov *lin.T, vm *lin.M4) {
	vm.Set(lin.M4I).ScaleMS(1, 1, 0)
}

// inverseTransform creates the inverse transform matrix from
// the location and orientation. Used in ray picking.
type inverseTransform func(*lin.T, *lin.M4)

// ivp inverse view transform. For ray casting.
func ivp(at *lin.T, vm *lin.M4) {
	vm.SetQ(lin.NewQ().Inv(at.Rot))
	vm.TranslateMT(at.Loc.X, at.Loc.Y, at.Loc.Z)
}

// nv is a null identity view.
func nv(at *lin.T, vm *lin.M4) {
	vm.Set(lin.M4I)
}
