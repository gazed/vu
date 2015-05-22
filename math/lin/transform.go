// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package lin

import "math"

// T is a 3D transform for rotation and translation. It excludes scaling and
// shear information. T is used as a simplification and optimization instead
// of keeping all transform information in a 4x4 matrix.
//
// T supports linear algebra operations that are similar to those supported
// by V3, V4, M3, M4, and Q.  The main ones are:
//      Multiply two transforms together to produce a composite transform.
//      Apply a transform or inverse transform to a vector.
type T struct {
	Loc *V3 // Location (translation, origin).
	Rot *Q  // Rotation (direction, orientation).
}

// Equals (==) returns true of all elements of transform t have the same value as
// the corresponding element of transform a.
func (t *T) Eq(a *T) bool { return t.Rot.Eq(a.Rot) && t.Loc.Eq(a.Loc) }

// Aeq (~=) almost-equals returns true if all the elements in transform t have
// essentially the same value as the corresponding elements in transform a.
// Used where a direct comparison is unlikely to return true due to floats.
func (t *T) Aeq(a *T) bool { return t.Rot.Aeq(a.Rot) && t.Loc.Aeq(a.Loc) }

// Set (=, copy, clone) assigns all the elements values from transform a to the
// corresponding element values in transform t. The updated transform t is returned.
func (t *T) Set(a *T) *T {
	t.Loc.Set(a.Loc)
	t.Rot.Set(a.Rot)
	return t
}

// SetI updates transform t to be the identity transform.
// The updated transform t is returned.
func (t *T) SetI() *T {
	t.Loc.SetS(0, 0, 0)
	t.Rot.Set(QI)
	return t
}

// SetVQ (=) sets the transform t based on the given quaternion rotation and
// translation location. The updated transform t is returned.
func (t *T) SetVQ(loc *V3, rot *Q) *T {
	t.Loc.Set(loc)
	t.Rot.Set(rot)
	return t
}

// SetAa updates transform t to have the rotation specified by the given
// axis and angle in radians. The updated transform t is returned.
func (t *T) SetAa(ax, ay, az, ang float64) *T {
	t.Rot.SetAa(ax, ay, az, ang)
	return t
}

// SetLoc updates transform t to have the location speccified  by lx, ly, lz.
// The updated transform t is returned.
func (t *T) SetLoc(lx, ly, lz float64) *T {
	t.Loc.X, t.Loc.Y, t.Loc.Z = lx, ly, lz
	return t
}

// SetRot updates transform t to have the rotation speccified  by x, y, z, w.
// The updated transform t is returned.
func (t *T) SetRot(x, y, z, w float64) *T {
	t.Rot.X, t.Rot.Y, t.Rot.Z, t.Rot.W = x, y, z, w
	return t
}

// Mult (*) updates the transform t to be the product of the
// transforms a and b. Transform t may be used as one or both of
// the input transforms. The updated transform t is returned.
func (t *T) Mult(a, b *T) *T {
	tx, ty, tz := t.Loc.GetS() // preserve original translation.
	t.Loc.MultvQ(b.Loc, a.Rot) // apply rotation to incoming translation.
	t.Loc.X, t.Loc.Y, t.Loc.Z = t.Loc.X+tx, t.Loc.Y+ty, t.Loc.Z+tz
	t.Rot.Mult(a.Rot, b.Rot)
	return t
}

// App applies its tranform to vector v. The updated vector v is returned.
func (t *T) App(v *V3) *V3 {
	v.MultvQ(v, t.Rot) // apply rotation.
	v.Add(v, t.Loc)    // apply translation.
	return v
}

// AppS applies transform t, rotation then translation, to input scalar
// vector (x,y,z) returning the transformed scalar vector (vx,vy,vz).
func (t *T) AppS(x, y, z float64) (vx, vy, vz float64) {
	vx, vy, vz = MultSQ(x, y, z, t.Rot)             // apply rotation.
	return vx + t.Loc.X, vy + t.Loc.Y, vz + t.Loc.Z // apply translation.
}

// AppR applies just the transform rotation to input vector (x,y,z)
// returning the rotated vector (vx,vy,vz)
func (t *T) AppR(x, y, z float64) (vx, vy, vz float64) {
	return MultSQ(x, y, z, t.Rot) // apply rotation.
}

// Inv updates vector v to be the inverse transform t applied
// to vector a.  The updated vector v is returned.
func (t *T) Inv(v *V3) *V3 {
	v.Sub(v, t.Loc)                            // apply inverse translation.
	ix, iy, iz := -t.Rot.X, -t.Rot.Y, -t.Rot.Z // apply inverse rotation.
	v.X, v.Y, v.Z = multSQ(v.X, v.Y, v.Z, ix, iy, iz, t.Rot.W)
	return v
}

// InvS applies the inverse transform t, inverse translation, then inverse
// rotation, to input vector (x,y,z) returning the transformed vector (vx,vy,vz).
func (t *T) InvS(x, y, z float64) (vx, vy, vz float64) {
	vx, vy, vz = x-t.Loc.X, y-t.Loc.Y, z-t.Loc.Z // apply inverse translation.
	ix, iy, iz := -t.Rot.X, -t.Rot.Y, -t.Rot.Z   // apply inverse rotation.
	return multSQ(vx, vy, vz, ix, iy, iz, t.Rot.W)
}

// Integrate updates transform t to be the linear integration of
// transform a with the given linear velocity linv, and angular velocity angv
// over the given amount of time dt. Transforms t and a must be distinct.
// The input vectors linv, angv are not changed.
// The updated transform t is returned.
//
// Based on bullet physics: btTransformUtil::integrateTransform.
func (t *T) Integrate(a *T, linv, angv *V3, dt float64) *T {

	// add interpolated linear velocity to current velocity.
	t.Loc.X = a.Loc.X + linv.X*dt
	t.Loc.Y = a.Loc.Y + linv.Y*dt
	t.Loc.Z = a.Loc.Z + linv.Z*dt

	// add interpolated angular velocity to current rotation. Google:
	//    "Practical Parameterization of Rotations Using the Exponential Map",
	//    F. Sebastian Grassia
	angularMotionLimit := 0.5 * HALF_PI
	angLen := angv.Len()
	if angLen*dt > angularMotionLimit {
		angLen = angularMotionLimit / dt // limit the angular motion
	}
	fac := 0.0
	if angLen < 0.001 {
		// Taylor's expansions of sync function
		fac = 0.5*dt - dt*dt*dt*0.020833333333*angLen*angLen
	} else {
		fac = math.Sin(0.5*angLen*dt) / angLen
	}

	// apply s rotation to existing rotation r
	rx, ry, rz, rw := a.Rot.X, a.Rot.Y, a.Rot.Z, a.Rot.W
	sx, sy, sz, sw := angv.X*fac, angv.Y*fac, angv.Z*fac, math.Cos(angLen*dt*0.5)
	t.Rot.X = rw*sx + rx*sw - ry*sz + rz*sy
	t.Rot.Y = rw*sy + rx*sz + ry*sw - rz*sx
	t.Rot.Z = rw*sz - rx*sy + ry*sx + rz*sw
	t.Rot.W = rw*sw - rx*sx - ry*sy - rz*sz
	t.Rot.Unit()
	return t
}

// ============================================================================
// convenience functions for allocating transforms. Nothing else should allocate.

// NewT creates and returns a transform at the origin with no rotation.
func NewT() *T {
	return &T{&V3{}, &Q{0, 0, 0, 1}}
}
