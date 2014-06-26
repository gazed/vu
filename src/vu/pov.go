// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package vu

// pov combines location and direction (orientation) to give a "point of view".

import (
	"vu/math/lin"
)

// pov combines a location and an orientation. It is intended to be
// embedded in other structures as a base class. While not visible
// outside the package, it confers some useful public methods to structures
// that embed it.
//
// pov is used to build up views and viewTransforms to get particular types
// of cameras and camera behaviours.
type pov struct {
	loc *lin.V3 // location/postion - where we are.
	dir *lin.Q  // rotation/direction/orientation - which way we're facing.
}

// Location provides safe access for the location coordinates.
func (p *pov) Location() (x, y, z float64) {
	if p.loc == nil {
		p.loc = &lin.V3{0, 0, 0}
	}
	return p.loc.X, p.loc.Y, p.loc.Z
}

// SetLocation directly places an object at the supplied coordinates.
func (p *pov) SetLocation(x, y, z float64) {
	if p.loc == nil {
		p.loc = &lin.V3{0, 0, 0}
	}
	p.loc.X, p.loc.Y, p.loc.Z = x, y, z
}

// Rotation provides safe access to the current quaternion values.
// Intended to be used with SetRotation to set one object facing the same
// direction as some other object.
func (p *pov) Rotation() (x, y, z, w float64) {
	if p.dir == nil {
		p.dir = &lin.Q{0, 0, 0, 1}
	}
	return p.dir.X, p.dir.Y, p.dir.Z, p.dir.W
}

// SetRotation puts quaternion values direction in the spots direction.
// Intended to be used with Rotation to set one object facing the same
// direction as some other object.
func (p *pov) SetRotation(x, y, z, w float64) {
	if p.dir == nil {
		p.dir = &lin.Q{0, 0, 0, 1}
	}
	p.dir.X, p.dir.Y, p.dir.Z, p.dir.W = x, y, z, w
}

// AxisAngle gets the facing direction returning an axis defined by
// x, y, z and an angle around the axis specified in degrees.
func (p *pov) AxisAngle() (x, y, z, degrees float64) {
	return p.dir.Aa()
}

// SetAxisAngle using an axis defined by x, y, z and an angle specified
// in degrees.
func (p *pov) SetAxisAngle(x, y, z, degrees float64) {
	p.dir.SetAa(x, y, z, degrees)
}

// Move increments the current position with respect to the current
// orientation, i.e. adds the distance travelled in the current direction
// to the current location.
func (p *pov) Move(x, y, z float64) {
	dx, dy, dz := lin.MultSQ(x, y, z, p.dir)
	p.loc.X += dx
	p.loc.Y += dy
	p.loc.Z += dz
}

// Spin rotates the current direction by the given number degrees around each
// axis.  Generally this is called with one direction change at a time.
func (p *pov) Spin(x, y, z float64) {
	if x != 0 {
		rotation := lin.NewQ().SetAa(1, 0, 0, lin.Rad(x))
		p.dir.Mult(rotation, p.dir)
	}
	if y != 0 {
		rotation := lin.NewQ().SetAa(0, 1, 0, lin.Rad(y))
		p.dir.Mult(p.dir, rotation)
	}
	if z != 0 {
		rotation := lin.NewQ().SetAa(0, 0, 1, lin.Rad(z))
		p.dir.Mult(rotation, p.dir)
	}
}

// pov
// ===========================================================================
// view

// view combines a location direction (pov) with a separate up/down angle
// tracking.  This allows use as a FPS camera which can limit up/down to
// a given range (often 180).
type view struct {
	pov               // Location/direction.
	up  float64       // The up/down angle in degrees.  Limit this to +90/-90
	vt  viewTransform // The assigned view transform matrix.
	q0  *lin.Q        // scratch quaternion.
}

// newView creates a default point of view that is looking down the positive
// Z axis.
func newView() *view {
	v := &view{}
	v.loc = &lin.V3{0, 0, 0}
	v.dir = &lin.Q{0, 0, 0, 1}
	v.q0 = &lin.Q{}
	return v
}

// view
// ===========================================================================
// view transforms

// View transforms create a view matrix based on a point of view.
// This is the camera, and since the camera is always at 0, 0, 0 then
// moving forward by x:units really means moving the world (everything else)
// back -x:units.  Same with rotating one way by x:degrees really means
// rotating the world -x:degrees.
type viewTransform func(*view, *lin.M4) *lin.M4

// getViewTransform turns the publicly visible view transform choices
// into a view tranform function.
func getViewTransform(transform int) viewTransform {
	switch transform {
	case VP:
		return vp
	case VO:
		return vo
	case VF:
		return vf
	case XZ_XY:
		return xz_xy
	default:
		return vp
	}
}

// vp view transform implementation.
func vp(v *view, vm *lin.M4) *lin.M4 {
	vm.SetQ(v.dir)
	return vm.TranslateTM(-v.loc.X, -v.loc.Y, -v.loc.Z)
}

// vo view transform implementation.
func vo(v *view, vm *lin.M4) *lin.M4 {
	return vm.Set(lin.M4I).ScaleMS(1, 1, 0)
}

// vf view transform implementation.
func vf(v *view, vm *lin.M4) *lin.M4 {
	rot := v.q0.SetAa(1, 0, 0, lin.Rad(-v.up))
	rot.Mult(rot, v.dir)
	return vm.SetQ(rot).TranslateTM(-v.loc.X, -v.loc.Y, -v.loc.Z)
}

// xz_xy view transform implementation.
func xz_xy(v *view, vm *lin.M4) *lin.M4 {
	rot := v.q0.SetAa(1, 0, 0, -lin.Rad(90))
	return vm.SetQ(rot).ScaleMS(1, 1, 0).TranslateTM(-v.loc.X, -v.loc.Y, -v.loc.Z)
}
