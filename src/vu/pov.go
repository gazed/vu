// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package vu

// pov combines location and direction to give a "point of view".

import (
	"vu/math/lin"
)

// pov combines a location and an orientation.  It is intended to be
// embedded in other structures as a base class.  While not visible
// outside the package, it confers some useful public methods to structures
// that embed it.
//
// pov is used to build up Views and ViewTransforms to get particular types
// of cameras and camera behaviours.
type pov struct {
	loc *lin.V3 // location/postion - where we are.
	dir *lin.Q  // direction/orientation - which way we're facing.
}

// Location provides safe access for the location coordinates.
func (p *pov) Location() (x, y, z float32) {
	if p.loc == nil {
		p.loc = &lin.V3{0, 0, 0}
	}
	return p.loc.X, p.loc.Y, p.loc.Z
}

// SetLocation directly places an object at the supplied coordinates.
func (p *pov) SetLocation(x, y, z float32) {
	if p.loc == nil {
		p.loc = &lin.V3{0, 0, 0}
	}
	p.loc.X, p.loc.Y, p.loc.Z = x, y, z
}

// Rotation provides safe access to the current quaternion values.
// Intended to be used with SetRotation to set one object facing the same
// direction as some other object.
func (p *pov) Rotation() (x, y, z, w float32) {
	if p.dir == nil {
		p.dir = &lin.Q{0, 0, 0, 1}
	}
	return p.dir.X, p.dir.Y, p.dir.Z, p.dir.W
}

// SetRotation puts quaternion values direction in the spots direction.
// Intended to be used with Rotation to set one object facing the same
// direction as some other object.
func (p *pov) SetRotation(x, y, z, w float32) {
	if p.dir == nil {
		p.dir = &lin.Q{0, 0, 0, 1}
	}
	p.dir.X, p.dir.Y, p.dir.Z, p.dir.W = x, y, z, w
}

// AxisAngle gets the facing direction returning an axis defined by
// x, y, z and an angle around the axis specified in degrees.
func (p *pov) AxisAngle() (x, y, z, degrees float32) {
	axis, angle := p.dir.AxisAngle()
	return axis.X, axis.Y, axis.Z, angle
}

// SetAxisAngle using an axis defined by x, y, z and an angle specified
// in degrees.
func (p *pov) SetAxisAngle(x, y, z, degrees float32) {
	p.dir = lin.QAxisAngle(&lin.V3{x, y, z}, degrees)
}

// Move increments the current position with respect to the current
// orientation.
func (p *pov) Move(x, y, z float32) {
	dist := (&lin.V3{x, y, z}).MultQ(p.dir)
	p.loc.X += dist.X
	p.loc.Y += dist.Y
	p.loc.Z += dist.Z
}

// RotateX changes the current orientation by some amount of x-axis rotation.
func (p *pov) RotateX(degrees float32) {
	rotation := lin.QAxisAngle(&lin.V3{1, 0, 0}, degrees)
	p.dir = rotation.Mult(p.dir)
}

// RotateY changes the current orientation by some amount of y-axis rotation.
func (p *pov) RotateY(degrees float32) {
	rotation := lin.QAxisAngle(&lin.V3{0, 1, 0}, degrees)
	p.dir = p.dir.Mult(rotation)
}

// RotateZ changes the current orientation by some amount of z-axis rotation.
func (p *pov) RotateZ(degrees float32) {
	rotation := lin.QAxisAngle(&lin.V3{0, 0, 1}, degrees)
	p.dir = rotation.Mult(p.dir)
}

// pov
// ===========================================================================
// view

// view combines a location direction (pov) with a separate up/down angle
// tracking.  This allows use as a FPS camera which can limit up/down to
// a given range (often 180).
type view struct {
	pov               // Location/direction.
	pl  *lin.V3       // Previous location. Used in collision resolution.
	up  float32       // The up/down angle in degrees.  Limit this to +90/-90
	vt  viewTransform // The assigned view transform matrix.
}

// newView creates a default point of view that is looking down the positive
// Z axis.
func newView() *view {
	v := &view{}
	v.loc = &lin.V3{0, 0, 0}
	v.pl = &lin.V3{0, 0, 0}
	v.dir = &lin.Q{0, 0, 0, 1}
	return v
}

func (v *view) PreviousLocation() (x, y, z float32) {
	if v.pl == nil {
		v.pl = &lin.V3{0, 0, 0}
	}
	return v.pl.X, v.pl.Y, v.pl.Z
}

// move increments the current position with respect to the current
// orientation and remembers the previous position.
func (v *view) Move(x, y, z float32) {
	v.pl.X, v.pl.Y, v.pl.Z = v.loc.X, v.loc.Y, v.loc.Z
	v.pov.Move(x, y, z)
}

// view
// ===========================================================================
// view transforms

// View transforms create a view matrix based on a point of view.
type viewTransform func(*view) *lin.M4

// View transform choices. These are used when adding a new scene. The view
// transform dictates the overall behaviour of the camera.
//    VP    is a perspective view transform. It is the reverse order of the
//          model transforms so that the scene revolves around the camera.
//    VO    is a orthographic view transform.
//    VF    is a first person view transform that uses the up/down angle.
//    XZ_XY is an overlay view that transforms an X,-Z mapped object to X,Y.
//          Useful for turning 3D floor plans into 2D mini-maps.
const (
	VP    = iota // Perspective view transform.
	VO           // Orthographic view transform.
	VF           // First person view transform.
	XZ_XY        // Perspective to Ortho view transform.
)

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
func vp(v *view) *lin.M4 {
	view := v.dir.Clone().Inverse().M4()
	return view.TranslateL(-v.loc.X, -v.loc.Y, -v.loc.Z)
}

// vo view transform implementation.
func vo(v *view) *lin.M4 {
	return lin.M4Scaler(1, 1, 0).TranslateR(5, 5, 0)
}

// vf view transform implementation.
func vf(v *view) *lin.M4 {
	updownRot := lin.QAxisAngle(&lin.V3{1, 0, 0}, v.up)
	view := updownRot.Mult(v.dir.Clone().Inverse()).M4()
	return view.TranslateL(-v.loc.X, -v.loc.Y, -v.loc.Z)
}

// xz_xy view transform implementation.
func xz_xy(v *view) *lin.M4 {
	rotx := lin.QAxisAngle(&lin.V3{1, 0, 0}, 90).M4()
	view := rotx.Mult(lin.M4Scaler(1, 1, 0))
	return view.TranslateL(-v.loc.X, -v.loc.Y, -v.loc.Z)
}
