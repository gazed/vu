// Copyright © 2013-2024 Galvanized Logic Inc.

package lin

// quaternion.go deals with quaternion math specifically for linear algebra rotations.
// For a nice explanation of quaternions see http://3dgep.com/understanding-quaternions

import (
	"log/slog"
	"math"
)

// Q is a unit length quaternion representing an angle of rotation.
// Q is used to track/manipulate 3D object rotations. Quaternions behave
// nicely for mathematical operations other than they are not commutative.
type Q struct {
	X float64 // X component of direction vector.
	Y float64 // Y component of direction vector.
	Z float64 // Z component of direction vector.
	W float64 // Angle of rotation.
}

// QI provides a reference identity matrix that can be used
// in calculations. It should never be changed.
var QI = &Q{0, 0, 0, 1}

// Eq (==) returns true if each element in the quaternion q has the same value
// as the corresponding element in quaterion r.
func (q *Q) Eq(r *Q) bool {
	return q.W == r.W && q.Z == r.Z && q.Y == r.Y && q.X == r.X
}

// Aeq (~=) almost-equals returns true if all the elements in quaternion q have
// essentially the same value as the corresponding elements in quaternion r.
// Used where a direct comparison is unlikely to return true due to floats.
func (q *Q) Aeq(r *Q) bool {
	return Aeq(q.X, r.X) && Aeq(q.Y, r.Y) && Aeq(q.Z, r.Z) && Aeq(q.W, r.W)
}

// GetS returns the component parts of a quaternion.
func (q *Q) GetS() (x, y, z, w float64) { return q.X, q.Y, q.Z, q.W }

// SetS (=) explicitly sets each of the quaternion values to the given values.
// The updated quaternion q is returned.
func (q *Q) SetS(x, y, z, w float64) *Q {
	q.X, q.Y, q.Z, q.W = x, y, z, w
	return q
}

// Set (=) assigns all the elements values from quaternion r to the corresponding
// element values in quaternion q. The updated quaternion q is returned.
func (q *Q) Set(r *Q) *Q {
	q.X, q.Y, q.Z, q.W = r.X, r.Y, r.Z, r.W
	return q
}

// Inv updates q to be inverse of quaternion r. The updated q is returned.
// The inverse of a quaternion is the same as the conjugate,
// as long as the quaternion is unit-length.
func (q *Q) Inv(r *Q) *Q {
	q.X, q.Y, q.Z, q.W = -r.X, -r.Y, -r.Z, r.W
	return q
}

// Add (+) quaternions r and s returning the result in quaternion q.
func (q *Q) Add(r, s *Q) *Q {
	q.X, q.Y, q.Z, q.W = r.X+s.X, r.Y+s.Y, r.Z+s.Z, r.W+s.W
	return q
}

// Neg (-) returns the negative of quaternion q where each element is negated.
// The updated q is returned.
func (q *Q) Neg() *Q {
	q.X, q.Y, q.Z, q.W = -q.X, -q.Y, -q.Z, -q.W
	return q
}

// Sub (-) subtracts quaternion s from r returning the difference in quaternion q.
func (q *Q) Sub(r, s *Q) *Q {
	q.X, q.Y, q.Z, q.W = r.X-s.X, r.Y-s.Y, r.Z-s.Z, r.W-s.W
	return q
}

// Scale (*=) quaternion q by s returning the result in quaternion q.
func (q *Q) Scale(s float64) *Q {
	q.X, q.Y, q.Z, q.W = q.X*s, q.Y*s, q.Z*s, q.W*s
	return q
}

// Div (/= inverse-scale) divides each element in q by the given scalar value.
// Scale values of zero are logged as an error and q is not scaled.
// The updated q is returned.
func (q *Q) Div(s float64) *Q {
	if s == 0 {
		slog.Error("quaternion:q.Div: division by zero")
	} else {
		s := 1 / s
		q.X, q.Y, q.Z, q.W = q.X*s, q.Y*s, q.Z*s, q.W*s
	}
	return q
}

// Mult (*) multiplies quaternions r and s returning the result in q.
// This applies the rotation of s to r giving q, leaving r and s unchanged.
// It is safe to use the calling quaternion q as one or both of the parameters.
// For example (*=) is
//
//	q.Mult(q, s)
//
// The updated calling quaternion q is returned.
func (q *Q) Mult(r, s *Q) *Q {
	x := r.W*s.X + r.X*s.W - r.Y*s.Z + r.Z*s.Y
	y := r.W*s.Y + r.X*s.Z + r.Y*s.W - r.Z*s.X
	z := r.W*s.Z - r.X*s.Y + r.Y*s.X + r.Z*s.W
	w := r.W*s.W - r.X*s.X - r.Y*s.Y - r.Z*s.Z
	q.X, q.Y, q.Z, q.W = x, y, z, w
	return q
}

// Unit normalizes quaternion q to have length 1.
// The normalized (unit length) q is returned.  Quaternion q is not
// updated if the length of quaternion q is zero.
func (q *Q) Unit() *Q {
	qlen := q.Len()
	if qlen != 0 {
		q.Scale(1 / qlen)
	}
	return q
}

// Dot returns the dot product of the quaternions q and r.
// Quaternion q may be used as the input parameter.
// For example (Dot=), the length squared, is
//
//	q.Dot(q)
func (q *Q) Dot(r *Q) float64 { return q.X*r.X + q.Y*r.Y + q.Z*r.Z + q.W*r.W }

// Len returns the length of the quaternion q.
func (q *Q) Len() float64 { return math.Sqrt(q.Dot(q)) }

// Ang returns the angle in radians between quaternions q and r.
// For the formula to calculate angles between quaternions see:
//   - http://math.stackexchange.com/questions/90081/quaternion-distance
func (q *Q) Ang(r *Q) float64 {
	qdotr := q.Dot(r)
	return math.Acos(2*(qdotr*qdotr) - 1)
}

// Nlerp updates q to be the normalized linear interpolation between
// quaternions r and s where ratio is expected to be between 0 and 1.
// The input quaternions r and s are not changed. See:
//   - http://keithmaggio.wordpress.com/2011/02/15/math-magician-lerp-slerp-and-nlerp/
//   - http://number-none.com/product/Understanding Slerp, Then Not Using It/
//
// The updated calling quaternion q is returned.
func (q *Q) Nlerp(r, s *Q, ratio float64) *Q {
	q.X = (s.X-r.X)*ratio + r.X
	q.Y = (s.Y-r.Y)*ratio + r.Y
	q.Z = (s.Z-r.Z)*ratio + r.Z
	q.W = (s.W-r.W)*ratio + r.W
	return q.Unit() // normalize the linear interpolation for a rotation.
}

// Aa gets the rotation of quaternion q as an axis and angle.
// The axis (x, y, z) and the angle in radians is returned.
// The return elements will be zero if the length of the quaternion is 0.
// See:
//
//	http://web.archive.org/web/20041029003853/...
//	...http://www.j3d.org/matrix_faq/matrfaq_latest.html#Q57
func (q *Q) Aa() (ax, ay, az, angle float64) {
	sinSqr := 1.0 - q.W*q.W
	if AeqZ(sinSqr) {
		return 1, 0, 0, 2 * math.Acos(q.W)
	}
	sin := 1.0 / math.Sqrt(sinSqr)
	return q.X * sin, q.Y * sin, q.Z * sin, 2 * math.Acos(q.W)
}

// SetAa set axis-angle, updates q to have the rotation of the given
// axis (ax, ay, az) and angle (in radians). See:
//
//	http://web.archive.org/web/20041029003853/...
//	...http://www.j3d.org/matrix_faq/matrfaq_latest.html#Q56
//
// The updated quaternion q is returned.
// The quaternion q is set to 0,0,0,1 if the axis length is 0.
//
// Ensure SetAa returns quaternions consistent with SetM3.
// Convention 5 from "Consistent Representations of and Conversions
// Between 3D Rotations":
//
//	"The rotation angle ω is limited to the interval [0, π].
//	 For angles in the range ]π, 2π[, the sign of the unit axis vector nˆ
//	 must be reversed, and ω replaced by 2π − ω. For angles outside the range
//	 [0, 2π[, the angle must first be reduced to the interval [0, 2π[ by
//	 adding or subtracting the appropriate integer multiple of 2π."
func (q *Q) SetAa(ax, ay, az, angle float64) *Q {
	alenSqr := ax*ax + ay*ay + az*az
	if alenSqr == 0 {
		q.X, q.Y, q.Z, q.W = 0, 0, 0, 1
		return q
	}

	// Convention 5 matches output of SetM3.
	tau := 2 * math.Pi
	for angle > tau {
		angle -= tau // reduce multiples of 2π
	}
	if angle > math.Pi {
		angle = tau - angle        // restrict to 180...
		ax, ay, az = -ax, -ay, -az // ...and flip rotation axis.
	}
	s := math.Sin(angle*0.5) / math.Sqrt(alenSqr)
	q.X, q.Y, q.Z, q.W = ax*s, ay*s, az*s, math.Cos(angle*0.5)
	return q
}

// RotateTo returns a rotation to rotate from direction vector v1
// to direction vector v2.
//
// from https://stackoverflow.com/questions/1171849/finding-quaternion-representing-the-rotation-from-one-vector-to-another
func (q *Q) RotateTo(v1, v2 *V3) *Q {
	k_cos_theta := v1.Dot(v2)
	k := math.Sqrt(v1.LenSqr() * v2.LenSqr())
	if k_cos_theta/k == -1 {
		// 180 degree rotation around orthogonal vector.
		o := NewV3().Cross(v1, v2) // orthogonal vector
		return q.SetS(0, o.X, o.Y, o.Z).Unit()
	}
	axis := NewV3().Cross(v1, v2)
	angle := k_cos_theta + k
	return q.SetS(axis.X, axis.Y, axis.Z, angle).Unit()
}

// Roll, Pitch, Yaw from one of the answers to
// https://stackoverflow.com/questions/5782658/extracting-yaw-from-a-quaternion
// switch the answer Z:Yaw, Y:Pitch, X:Roll to be X:Pitch, Y:Yaw, Z:Roll

// Z Roll (Bank) in radians
func (q *Q) Roll() float64 {
	return math.Atan2(2.0*(q.Z*q.W+q.X*q.Y), -1.0+2.0*(q.W*q.W+q.X*q.X))
}

// X Pitch (Attitude) in radians
func (q *Q) Pitch() float64 {
	return math.Atan2(2.0*(q.Z*q.Y+q.W*q.X), 1.0-2.0*(q.X*q.X+q.Y*q.Y))
}

// Y Yaw (Heading) in radians
func (q *Q) Yaw() float64 {
	return math.Asin(2.0 * (q.Y*q.W - q.Z*q.X))
}

// ============================================================================
// quaternion-vector operations

// MultQV multiplies quaternion r and vector v and returns the result in
// quaternion q. The upated quaternion q is returned.
func (q *Q) MultQV(r *Q, v *V3) *Q {
	x := +r.W*v.X + r.Y*v.Z - r.Z*v.Y
	y := +r.W*v.Y + r.Z*v.X - r.X*v.Z
	z := +r.W*v.Z + r.X*v.Y - r.Y*v.X
	w := -r.X*v.X - r.Y*v.Y - r.Z*v.Z
	q.X, q.Y, q.Z, q.W = x, y, z, w
	return q
}

// ============================================================================
// quaternion-matrix operations

// SetM3 updates quaternion q to be the rotation of matrix m. See
//
//	http://www.flipcode.com/documents/matrfaq.html#Q55
//	http://www.euclideanspace.com/maths/geometry/rotations/conversions/matrixToQuaternion/
//	https://d3cw3dd2w32x2b.cloudfront.net/wp-content/uploads/2015/01/matrix-to-quat.pdf
//
// The updated q is returned.
//
// SetM3 outputs quaternions that are consistent with SetAa.
func (q *Q) SetM3(m *M3) *Q {
	t := 0.0
	if m.Zz < 0 {
		if m.Xx > m.Yy {
			t = 1.0 + m.Xx - m.Yy - m.Zz
			q.SetS(t, m.Xy+m.Yx, m.Zx+m.Xz, m.Zy-m.Yz)
		} else {
			t = 1.0 - m.Xx + m.Yy - m.Zz
			q.SetS(m.Xy+m.Yx, t, m.Yz+m.Zy, m.Xz-m.Zx)
		}
	} else {
		if m.Xx < -m.Yy {
			t = 1.0 - m.Xx - m.Yy + m.Zz
			q.SetS(m.Zx+m.Xz, m.Yz+m.Zy, t, m.Yx-m.Xy)
		} else {
			t = 1.0 + m.Xx + m.Yy + m.Zz
			q.SetS(m.Zy-m.Yz, m.Xz-m.Zx, m.Yx-m.Xy, t)
		}
	}
	q.Scale(0.5 / math.Sqrt(t))
	if q.W < 0 {
		return q.Scale(-1)
	}
	return q
}

// ============================================================================
// quaternion-transform operations

// SetT updates quaternion q to have the rotation in transform t.
// The updated quaternion q is returned.
func (q *Q) SetT(t *T) *Q {
	q.Set(t.Rot)
	return q
}

// ============================================================================
// convenience functions for allocating quaternions. Nothing else should allocate.
// methods above do not allocate memory.

// NewQ creates a new, all zero, quaternion.
func NewQ() *Q { return &Q{} }

// NewQI creates a new identity quaternion.
func NewQI() *Q { return &Q{W: 1} }
