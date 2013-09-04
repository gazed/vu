// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package lin

import "math"

// TODO check the math against http://www.cs.princeton.edu/~gewang/projects/darth/stuff/quat_faq.html

// Quaternions (of unit length) represent an angle of rotation and an
// direction (orientation) and are used to track/manipulate 3D object rotations.
// Quaternions behave nicely for mathematical operations other than they are not
// commutative.
type Q struct {
	X float32 // X component of direction vector
	Y float32 // Y component of direction vector
	Z float32 // Z component of direction vector
	W float32 // angle of rotation
}

// QIdentity creates a new quaternion that is the identity quaternion.
func QIdentity() *Q { return &Q{W: 1} }

// Clone creates a copy of a quaternion.  This is needed if when performing
// operations where the original quaternion values need to be preserved.
func (q *Q) Clone() *Q { return &Q{q.X, q.Y, q.Z, q.W} }

// Norm normalizes a quaternion works similar to a vector. This method will not
// do anything if the quaternion is close enough to zero or unit-length.
func (q *Q) Unit() *Q {
	dot := q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W
	if !IsZero(dot) && !IsOne(dot) {
		length := float32(math.Sqrt(float64(dot)))
		q.X /= length
		q.Y /= length
		q.Z /= length
		q.W /= length
	}
	return q
}

// Inverse inverts the calling quaternion.  The inverse of a quaternion is the
// same as the conjugate, as long as the quaternion is unit-length. The calling
// quaternion is returned so that it may be immediately used in another operation.
// Note that the quaternion is normalized by this method before inverting.
func (q *Q) Inverse() *Q {
	q.Unit()
	q.X = -q.X
	q.Y = -q.Y
	q.Z = -q.Z
	return q
}

// Add one quaternion to another.  The calling quaternion q is updated with
// the addition of the input quaternion q2. The calling quaternion is
// returned so that it may be immediately used in another operation.
// The updated quaternion is not normalized.
func (q *Q) Add(q2 *Q) *Q {
	q.X += q2.X
	q.Y += q2.Y
	q.Z += q2.Z
	q.W += q2.W
	return q
}

// Scale multplies the calling quaternion by a scaler value.
// The calling quaternion is returned so that it may be immediately
// used in another operation. The returned quaternion is not normalized.
func (q *Q) Scale(scale float32) *Q {
	q.X *= scale
	q.Y *= scale
	q.Z *= scale
	q.W *= scale
	return q
}

// Mult multiplies the calling quaternion q with the input quaternion q2.
// This applies the rotation of q2 to q.
func (q *Q) Mult(q2 *Q) *Q {
	x := q.W*q2.X + q.X*q2.W + q.Y*q2.Z - q.Z*q2.Y
	y := q.W*q2.Y + q.Y*q2.W + q.Z*q2.X - q.X*q2.Z
	z := q.W*q2.Z + q.Z*q2.W + q.X*q2.Y - q.Y*q2.X
	w := q.W*q2.W - q.X*q2.X - q.Y*q2.Y - q.Z*q2.Z
	q.X, q.Y, q.Z, q.W = x, y, z, w
	return q
}

// QAxisAngle creates a new quaternion that represents the given Axis/Angle.
// The formula for this is:
//
//    qx = ax * sin(angle/2)
//    qy = ay * sin(angle/2)
//    qz = az * sin(angle/2)
//    qw = cos(angle/2)
//
// where: the axis vector is normalised.
// The calling quaternion is returned so that it may be immediately
// used in another operation. The input axis is unchanged.
func QAxisAngle(axis *V3, angleInDegrees float32) *Q {
	angleInRadians := angleInDegrees * PI_OVER_180
	vn := (&V3{axis.X, axis.Y, axis.Z}).Unit()
	halfAngle := float64(angleInRadians * 0.5)
	sinAngle := float32(math.Sin(halfAngle))
	q := &Q{}
	q.X = vn.X * sinAngle
	q.Y = vn.Y * sinAngle
	q.Z = vn.Z * sinAngle
	q.W = float32(math.Cos(halfAngle))
	return q
}

// AxisAngle returns the axis and angle represented by the calling quaternion q.
func (q *Q) AxisAngle() (axis V3, angleInDegrees float32) {
	scale := float32(math.Sqrt(float64(q.X*q.X + q.Y*q.Y + q.Z*q.Z)))
	if IsZero(scale) {
		axis = V3{0, 0, -1}
	} else {
		axis = V3{q.X / scale, q.Y / scale, q.Z / scale}
	}
	angleInDegrees = float32(math.Acos(float64(q.W))*2) / PI_OVER_180
	return
}

// NLerp returns a new quaternion that is the normalized linear interpolation between
// quaternions q and q2 where fraction is expected to be between 0 and 1.
// The input quaternions q and q2 are not changed.
//
// See:
//    http://keithmaggio.wordpress.com/2011/02/15/math-magician-lerp-slerp-and-nlerp/
//    http://number-none.com/product/Understanding Slerp, Then Not Using It/
func (q *Q) Nlerp(q2 *Q, fraction float32) *Q {
	nq := &Q{
		(q2.X-q.X)*fraction + q.X,
		(q2.Y-q.Y)*fraction + q.Y,
		(q2.Z-q.Z)*fraction + q.Z,
		(q2.W-q.W)*fraction + q.W,
	}
	return nq.Unit()
}

// M4 converts a quaternion q to a new rotation matrix.  The quaternion q is
// expected to be normalized before converting it to a matrix.
// A newly allocated 4x4 rotation matrix is returned.
func (q *Q) M4() *M4 {
	x2 := q.X * q.X
	y2 := q.Y * q.Y
	z2 := q.Z * q.Z
	xy := q.X * q.Y
	xz := q.X * q.Z
	yz := q.Y * q.Z
	wx := q.W * q.X
	wy := q.W * q.Y
	wz := q.W * q.Z
	m := &M4{
		1 - 2*(y2+z2), 2 * (xy + wz), 2 * (xz - wy), 0,
		2 * (xy - wz), 1 - 2*(x2+z2), 2 * (yz + wx), 0,
		2 * (xz + wy), 2 * (yz - wx), 1 - 2*(x2+y2), 0,
		0, 0, 0, 1}
	return m
}

// TODO may need to know how to rotate from one vector to another.
// http://stackoverflow.com/questions/1171849/finding-quaternion-representing-the-rotation-from-one-vector-to-another
//     quaternion for double the required rotation to get from u to v is:
//         q.w   == dot(u, v)
//         q.xyz == cross(u, v)
