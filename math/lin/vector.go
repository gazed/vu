// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package lin

// Vector performs 3 or 4 element vector related math needed for 3D applications.

import (
	"log"
	"math"
)

// V3 is a 3 element vector. This can also be used as a point.
type V3 struct {
	X float64 // increments as X moves to the right.
	Y float64 // increments as Y moves up from bottom left.
	Z float64 // increments as Z moves out of the screen (right handed view space).
}

// V4 is a 4 element vector. It can be used for points and directions where,
// as a point it would have W:1, and as a direction it would have W:0.
type V4 struct {
	X float64 // increments as X moves to the right.
	Y float64 // increments as Y moves up from bottom left.
	Z float64 // increments as Z moves out of the screen (right handed view space).
	W float64 // fourth dimension makes for nice 3D matrix math.
}

// Eq (==) returns true if each element in the vector v has the same value
// as the corresponding element in vector a.
func (v *V3) Eq(a *V3) bool {
	return v.Z == a.Z && v.Y == a.Y && v.X == a.X
}

// Eq (==) returns true if each element in the vector v has the same value
// as the corresponding element in vector a.
func (v *V4) Eq(a *V4) bool {
	return v.W == a.W && v.Z == a.Z && v.Y == a.Y && v.X == a.X
}

// Aeq (~=) almost-equals returns true if all the elements in vector v have
// essentially the same value as the corresponding elements in vector a.
// Used where a direct comparison is unlikely to return true due to floats.
func (v *V3) Aeq(a *V3) bool {
	return Aeq(v.X, a.X) && Aeq(v.Y, a.Y) && Aeq(v.Z, a.Z)
}

// AeqZ (~=) almost equals zero returns true if the square length of the vector
// is close enough to zero that it makes no difference.
func (v *V3) AeqZ() bool { return v.Dot(v) < Epsilon }

// AeqZ (~=) almost equals zero returns true if the square length of the vector
// is close enough to zero that it makes no difference.
func (v *V4) AeqZ() bool { return v.Dot(v) < Epsilon }

// GetS returns the float64 values of the vector.
func (v *V3) GetS() (x, y, z float64) { return v.X, v.Y, v.Z }

// GetS returns the float64 values of the vector.
func (v *V4) GetS() (x, y, z, w float64) { return v.X, v.Y, v.Z, v.W }

// SetS (=) sets the vector elements to the given values.
// The updated vector v is returned.
func (v *V3) SetS(x, y, z float64) *V3 {
	v.X, v.Y, v.Z = x, y, z
	return v
}

// SetS (=) sets the vector elements to the given values.
// The updated vector v is returned.
func (v *V4) SetS(x, y, z, w float64) *V4 {
	v.X, v.Y, v.Z, v.W = x, y, z, w
	return v
}

// Set (=, copy, clone) sets the elements of vector v to have the same values
// as the elements of vector a. The updated vector v is returned.
func (v *V3) Set(a *V3) *V3 {
	v.X, v.Y, v.Z = a.X, a.Y, a.Z
	return v
}

// Set (=, copy, clone) sets the elements of vector v to have the same values
// as the elements of vector a. The updated vector v is returned.
func (v *V4) Set(a *V4) *V4 {
	v.X, v.Y, v.Z, v.W = a.X, a.Y, a.Z, a.W
	return v
}

// Swap exchanges the element values of vectors v and a.
// The updated vector v is returned. Vector a is also updated.
func (v *V3) Swap(a *V3) *V3 {
	v.X, a.X = a.X, v.X
	v.Y, a.Y = a.Y, v.Y
	v.Z, a.Z = a.Z, v.Z
	return v
}

// Swap exchanges the element values of vectors v and a.
// Same behaviour as V3.Swap().
func (v *V4) Swap(a *V4) *V4 {
	v.X, a.X = a.X, v.X
	v.Y, a.Y = a.Y, v.Y
	v.Z, a.Z = a.Z, v.Z
	v.W, a.W = a.W, v.W
	return v
}

// Min updates the vector v elements to be the minimum of the corresponding
// elements from either vectors a or b. The updated vector v is returned.
func (v *V3) Min(a, b *V3) *V3 {
	v.X, v.Y, v.Z = math.Min(b.X, a.X), math.Min(b.Y, a.Y), math.Min(b.Z, a.Z)
	return v
}

// Min updates the vector v elements to be the minimum of the corresponding
// elements from either vectors a or b. Same behaviour as V3.Min().
func (v *V4) Min(a, b *V4) *V4 {
	v.X, v.Y, v.Z, v.W = math.Min(b.X, a.X), math.Min(b.Y, a.Y), math.Min(b.Z, a.Z), math.Min(b.W, a.W)
	return v
}

// Max updates the vector v elements to be the maxiumum of the corresponding
// elements from either vectors a or b. The updated vector v is returned.
func (v *V3) Max(a, b *V3) *V3 {
	v.X, v.Y, v.Z = math.Max(b.X, a.X), math.Max(b.Y, a.Y), math.Max(b.Z, a.Z)
	return v
}

// Max updates the vector v elements to be the maximum of the corresponding
// elements from either vectors a or b. Same behaviour as V3.Max().
func (v *V4) Max(a, b *V4) *V4 {
	v.X, v.Y, v.Z, v.W = math.Max(b.X, a.X), math.Max(b.Y, a.Y), math.Max(b.Z, a.Z), math.Max(b.W, a.W)
	return v
}

// Abs updates vector v to have the absolute value of the elements of vector a.
// The updated vector v is returned.
func (v *V3) Abs() *V3 {
	v.X, v.Y, v.Z = math.Abs(v.X), math.Abs(v.Y), math.Abs(v.Z)
	return v
}

// Abs updates vector v to have the absolute value of the elements of vector a.
// Same behaviour as V3.Abs().
func (v *V4) Abs() *V4 {
	v.X, v.Y, v.Z, v.W = math.Abs(v.X), math.Abs(v.Y), math.Abs(v.Z), math.Abs(v.W)
	return v
}

// Neg (-) sets vector v to be the negative values of vector a.
// Vector v may be used as the input parameter.
// The updated vector v is returned.
func (v *V3) Neg(a *V3) *V3 {
	v.X, v.Y, v.Z = -a.X, -a.Y, -a.Z
	return v
}

// Neg (-) sets vector v to be the negative values of vector a.
// Same behaviour as V3.Neg().
func (v *V4) Neg(a *V4) *V4 {
	v.X, v.Y, v.Z, v.W = -a.X, -a.Y, -a.Z, -a.W
	return v
}

// Add (+) adds vectors a and b storing the results of the addition in v.
// Vector v may be used as one or both of the parameters.
// For example (+=) is
//     v.Add(v, b)
// The updated vector v is returned.
func (v *V3) Add(a, b *V3) *V3 {
	v.X, v.Y, v.Z = a.X+b.X, a.Y+b.Y, a.Z+b.Z
	return v
}

// Add (+) adds vectors a and b storing the results of the addition in v.
// Same behaviour as V3.Add().
func (v *V4) Add(a, b *V4) *V4 {
	v.X, v.Y, v.Z, v.W = a.X+b.X, a.Y+b.Y, a.Z+b.Z, a.W+b.W
	return v
}

// Sub (-) subtracts vectors b from a storing the results of the subtraction in v.
// Vector v may be used as one or both of the parameters.
// For example (-=) is
//     v.Sub(v, b)
// The updated vector v is returned.
func (v *V3) Sub(a, b *V3) *V3 {
	v.X, v.Y, v.Z = a.X-b.X, a.Y-b.Y, a.Z-b.Z
	return v
}

// Sub (-) subtracts vectors b from a storing the results of the subtraction in v.
// Same behaviour as V3.Sub().
func (v *V4) Sub(a, b *V4) *V4 {
	v.X, v.Y, v.Z, v.W = a.X-b.X, a.Y-b.Y, a.Z-b.Z, a.W-b.W
	return v
}

// Mult (*) multiplies the elements of vectors a and b storing the result in v.
// Vector v may be used as one or both of the parameters. For example (*=) is
// For example (*=) is
//     v.Mult(v, b)
// The updated vector v is returned.
func (v *V3) Mult(a, b *V3) *V3 {
	v.X, v.Y, v.Z = a.X*b.X, a.Y*b.Y, a.Z*b.Z
	return v
}

// Mult (*) multiplies the elements of vectors a and b storing the result in v.
// Same behaviour as V3.Mult().
func (v *V4) Mult(a, b *V4) *V4 {
	v.X, v.Y, v.Z, v.W = a.X*b.X, a.Y*b.Y, a.Z*b.Z, a.W*b.W
	return v
}

// MultQ (*) multiplies a vector by quaternion, effectively applying the
// rotation of quaternion q to vector a and storing the result in v. The input
// vector a, and quaternion q are unchanged.
func (v *V3) MultQ(a *V3, q *Q) *V3 {
	// A implementation based on:
	//   http://molecularmusings.wordpress.com/2013/05/24/a-faster-quaternion-vector-multiplication/
	// It benchmarked about 40% faster than the standard implementation at:
	//   http://www.mathworks.com/help/aeroblks/quaternionrotation.html

	// t = 2 * cross(q.xyz, v)
	c0x, c0y, c0z := 2*(q.Y*a.Z-q.Z*a.Y), 2*(q.Z*a.X-q.X*a.Z), 2*(q.X*a.Y-q.Y*a.X) //cross(q.xyz, v)

	// v' = v + q.w * t + cross(q.xyz, t)
	c1x, c1y, c1z := q.Y*c0z-q.Z*c0y, q.Z*c0x-q.X*c0z, q.X*c0y-q.Y*c0x // cross(q.xyz, t)
	v.X, v.Y, v.Z = a.X+q.W*c0x+c1x, a.Y+q.W*c0y+c1y, a.Z+q.W*c0z+c1z
	return v
}

// Scale (*=) updates the elements in vector v by multiplying the
// corresponding elements in vector a by the given scalar value.
// Vector v may be used as one or both of the vector parameters.
// The updated vector v is returned.
func (v *V3) Scale(a *V3, s float64) *V3 {
	v.X, v.Y, v.Z = a.X*s, a.Y*s, a.Z*s
	return v
}

// Scale (*=) updates the elements in vector v by multiplying the
// corresponding elements in vector a by the given scalar value.
// Same behaviour as V3.Scale().
func (v *V4) Scale(a *V4, s float64) *V4 {
	v.X, v.Y, v.Z, v.W = a.X*s, a.Y*s, a.Z*s, a.W*s
	return v
}

// Div (/= inverse-scale) divides each element in v by the given scalar value.
// The updated vector v is returned. Vector v is not changed if scalar s is zero.
func (v *V3) Div(s float64) *V3 {
	if s != 0 {
		inv := 1 / s
		v.X, v.Y, v.Z = v.X*inv, v.Y*inv, v.Z*inv
	}
	return v
}

// Div (/= inverse-scale) divides each element in v by the given scalar value.
// Same behaviour as V3.Div().
func (v *V4) Div(s float64) *V4 {
	if s != 0 {
		inv := 1 / s
		v.X, v.Y, v.Z, v.W = v.X*inv, v.Y*inv, v.Z*inv, v.W*inv
	}
	return v
}

// Dot vector v with input vector a. Both vectors v and a are unchanged.
// Wikipedia states:
//    "This operation can be defined either algebraically or geometrically.
//     Algebraically, it is the sum of the products of the corresponding
//     entries of the two sequences of numbers. Geometrically, it is the
//     product of the magnitudes of the two vectors and the cosine of
//     the angle between them."
func (v *V3) Dot(a *V3) float64 { return v.X*a.X + v.Y*a.Y + v.Z*a.Z }

// Dot vector v with input vector a. Same behaviour as V3.Dot()
func (v *V4) Dot(a *V4) float64 { return v.X*a.X + v.Y*a.Y + v.Z*a.Z + v.W*a.W }

// Len returns the length of vector v. Vector length is the square root of
// the dot product. The calling vector v is unchanged.
func (v *V3) Len() float64 { return math.Sqrt(v.Dot(v)) }

// LenSqr returns the length of vector v squared.
// The calling vector v is unchanged.
func (v *V3) LenSqr() float64 { return v.Dot(v) }

// Len returns the length of vector v. Same behaviour as V3.Len()
func (v *V4) Len() float64 { return math.Sqrt(v.Dot(v)) }

// LenSqr returns the length of vector v squared.
// Same behaviour as V3.Len()
func (v *V4) LenSqr() float64 { return v.Dot(v) }

// Dist returns the distance between vector end-points v and a
// Both vectors (points) v and a are unchanged.
func (v *V3) Dist(a *V3) float64 { return math.Sqrt(v.DistSqr(a)) }

// DistSqr returns the distance squared between vector end-points v and a.
// Both vectors (points) v and a are unchanged.
func (v *V3) DistSqr(a *V3) float64 {
	dx, dy, dz := a.X-v.X, a.Y-v.Y, a.Z-v.Z
	return dx*dx + dy*dy + dz*dz
}

// Ang returns the angle in radians between vector v and input vector a.
// Ang returns 0 if the magnitude of the two vectors is 0.
func (v *V3) Ang(a *V3) float64 {
	magnitude := math.Sqrt(v.Dot(v) * a.Dot(a))
	if magnitude != 0 {
		return math.Acos(v.Dot(a) / magnitude)
	}
	log.Printf("Dev error. vector.V3:Ang division by zero")
	return 0
}

// Unit updates vector v such that its length is 1.
// Calling vector v is unchanged if its length is zero.
// The updated vector v is returned.
func (v *V3) Unit() *V3 {
	length := v.Len()
	if length != 0 {
		return v.Div(length)
	}
	return v
}

// Unit updates vector v such that its length is 1.
// Same behaviour as V3.Unit()
func (v *V4) Unit() *V4 {
	length := v.Len()
	if length != 0 {
		return v.Div(length)
	}
	return v
}

// Cross updates v to be the cross product of vectors a and b.
// A cross product vector is a vector that is perpendicular to both input
// vectors. This is only meaningful in 3 (or 7) dimensions.
// Input vectors a and b are unchanged. Vector v may be used as either
// input parameter.The updated vector v is returned.
func (v *V3) Cross(a, b *V3) *V3 {
	v.X, v.Y, v.Z = a.Y*b.Z-a.Z*b.Y, a.Z*b.X-a.X*b.Z, a.X*b.Y-a.Y*b.X
	return v
}

// Lerp updates vector v to be a fraction of the distance (linear interpolation)
// between the input vectors a and b. The input ratio is not verified, but is expected
// to be between 0 and 1. Vector v may be used as one of the parameters.
func (v *V3) Lerp(a, b *V3, fraction float64) *V3 {
	v.X = (b.X-a.X)*fraction + a.X
	v.Y = (b.Y-a.Y)*fraction + a.Y
	v.Z = (b.Z-a.Z)*fraction + a.Z
	return v
}

// Lerp updates vector v to be a fraction of the distance (linear interpolation)
// between the input vectors a and b. Same behaviour as V3.Lerp()
func (v *V4) Lerp(a, b *V4, ratio float64) *V4 {
	v.X = (b.X-a.X)*ratio + a.X
	v.Y = (b.Y-a.Y)*ratio + a.Y
	v.Z = (b.Z-a.Z)*ratio + a.Z
	v.W = (b.W-a.W)*ratio + a.W
	return v
}

// Nlerp updates vector v to be a normalized vector that is the linerar interpolation
// between a and b. See:
//    http://keithmaggio.wordpress.com/2011/02/15/math-magician-lerp-slerp-and-nlerp/
//    http://number-none.com/product/Understanding%20Slerp,%20Then%20Not%20Using%20It/
// The calling vector v may be used as either or both of the input parameters.
func (v *V3) Nlerp(a, b *V3, ratio float64) *V3 { return v.Lerp(a, b, ratio).Unit() }

// Nlerp updates vector v to be a normalized vector that is the linerar interpolation
// between a and b. Same behaviour as V3.Lerp()
func (v *V4) Nlerp(a, b *V4, ratio float64) *V4 { return v.Lerp(a, b, ratio).Unit() }

// Plane generates 2 vectors perpendicular to normal vector v.
// The perpendicular vectors and are returned as values of vectors p and q.
//
// Based on bullet physics: btVector3::btPlaneSpace1
func (v *V3) Plane(p, q *V3) {
	squareRootof12 := float64(0.7071067811865475244008443621048490)
	if math.Abs(v.Z) > squareRootof12 {
		a := v.Y*v.Y + v.Z*v.Z
		k := 1 / math.Sqrt(a)

		// p in y-z plane, q = n x p
		p.X, p.Y, p.Z = 0, -v.Z*k, v.Y*k
		q.X, q.Y, q.Z = a*k, -v.X*p.Z, v.X*p.Y
	} else {
		a := v.X*v.X + v.Y*v.Y
		k := 1 / math.Sqrt(a)

		// p in x-y plane, q = n x p
		p.X, p.Y, p.Z = -v.Y*k, v.X*k, 0
		q.X, q.Y, q.Z = -v.Z*p.Y, v.Z*p.X, a*k
	}
	return
}

// vector operations
// ============================================================================
// vector-matrix operations

// MultvM updates vector v to be the multiplication of row vector rv
// and matrix m. Vector v may be used as the input vector rv.
// The udpated vector v is returned.
//                   [ Xx Xy Xz ]
//    [ vx vy vz ] x [ Yx Yy Yz ] = [ vx' vy' vz' ]
//                   [ Zx Zy Zz ]
func (v *V3) MultvM(rv *V3, m *M3) *V3 {
	x := rv.X*m.Xx + rv.Y*m.Yx + rv.Z*m.Zx
	y := rv.X*m.Xy + rv.Y*m.Yy + rv.Z*m.Zy
	z := rv.X*m.Xz + rv.Y*m.Yz + rv.Z*m.Zz
	v.X, v.Y, v.Z = x, y, z
	return v
}

// MultvM updates vector v to be the multiplication of row vector rv
// and matrix m. Same behaviour as V4.MultvM().
//                      [ Xx Xy Xz Xw ]
//    [ vx vy vz vw ] x [ Yx Yy Yz Yw ] = [ vx' vy' vz' vw']
//                      [ Zx Zy Zz Zw ]
//                      [ Wx Wy Wz Ww ]
func (v *V4) MultvM(rv *V4, m *M4) *V4 {
	x := rv.X*m.Xx + rv.Y*m.Yx + rv.Z*m.Zx + rv.W*m.Wx
	y := rv.X*m.Xy + rv.Y*m.Yy + rv.Z*m.Zy + rv.W*m.Wy
	z := rv.X*m.Xz + rv.Y*m.Yz + rv.Z*m.Zz + rv.W*m.Wz
	w := rv.X*m.Xw + rv.Y*m.Yw + rv.Z*m.Zw + rv.W*m.Ww
	v.X, v.Y, v.Z, v.W = x, y, z, w
	return v
}

// MultMv updates vector v to be the multiplication of matrix m and
// column vector cv. Vector v may be used as the input vector cv.
// The udpated vector v is returned.
//    [ Xx Xy Xz ]   [ vx ]   [ vx' ]
//    [ Yx Yy Yz ] x [ vy ] = [ vx' ]
//    [ Zx Zy Zz ]   [ vz ]   [ vz' ]
func (v *V3) MultMv(m *M3, cv *V3) *V3 {
	x := m.Xx*cv.X + m.Xy*cv.Y + m.Xz*cv.Z
	y := m.Yx*cv.X + m.Yy*cv.Y + m.Yz*cv.Z
	z := m.Zx*cv.X + m.Zy*cv.Y + m.Zz*cv.Z
	v.X, v.Y, v.Z = x, y, z
	return v
}

// MultMv updates vector v to be the multiplication of matrix m and
// column vector cv. Same behaviour as V3.MultMv().
//    [ Xx Xy Xz Xw ]   [ vx ]   [ vx' ]
//    [ Yx Yy Yz Yw ] x [ vy ] = [ vy' ]
//    [ Zx Zy Zz Zw ]   [ vz ]   [ vz' ]
//    [ Wx Wy Wz Ww ]   [ vw ]   [ vw' ]
func (v *V4) MultMv(m *M4, cv *V4) *V4 {
	x := m.Xx*cv.X + m.Xy*cv.Y + m.Xz*cv.Z + m.Xw*cv.W
	y := m.Yx*cv.X + m.Yy*cv.Y + m.Yz*cv.Z + m.Yw*cv.W
	z := m.Zx*cv.X + m.Zy*cv.Y + m.Zz*cv.Z + m.Zw*cv.W
	w := m.Wx*cv.X + m.Wy*cv.Y + m.Wz*cv.Z + m.Ww*cv.W
	v.X, v.Y, v.Z, v.W = x, y, z, w
	return v
}

// vector-matrix operations
// ============================================================================
// vector-quaternion operations

// MultvQ  updates vector v to be the rotation of vector a by quaternion q.
func (v *V3) MultvQ(a *V3, q *Q) *V3 {
	v.X, v.Y, v.Z = multSQ(a.X, a.Y, a.Z, q.X, q.Y, q.Z, q.W)
	return v
}

// MultSQ applies rotation q to scalar vector (x,y,z)
// The updated scalar vector (vx,vy,vz) is returned.
func MultSQ(x, y, z float64, q *Q) (vx, vy, vz float64) {
	return multSQ(x, y, z, q.X, q.Y, q.Z, q.W)
}

// MultSQ applies rotation q (qx,qy,qz,qw) to scalar vector (x,y,z)
// The updated scalar vector (vx,vy,vz) is returned.
func multSQ(x, y, z, qx, qy, qz, qw float64) (vx, vy, vz float64) {
	k0 := qw*qw - 0.5

	// k1 = Q.V
	k1 := x*qx + y*qy + z*qz

	// (qq-1/2)V+(Q.V)Q
	rx := x*k0 + qx*k1
	ry := y*k0 + qy*k1
	rz := z*k0 + qz*k1

	// (Q.V)Q+(qq-1/2)V+q(QxV)
	rx += qw * (qy*z - qz*y)
	ry += qw * (qz*x - qx*z)
	rz += qw * (qx*y - qy*x)

	//  2((Q.V)Q+(qq-1/2)V+q(QxV))
	return rx + rx, ry + ry, rz + rz
}

// vector-quaternion operations
// ============================================================================
// vector-transform operations

// AppT updates vector v to be the transform t applied to vector a.  Vector a
// is unchanged.  The updated vector v is returned.
func (v *V3) AppT(t *T, a *V3) *V3 {
	v.X, v.Y, v.Z = t.AppS(a.X, a.Y, a.Z)
	return v
}

// vector-transform operations
// ============================================================================
// convenience functions for allocating vectors.  Nothing else should allocate.

// NewV3 creates a new, all zero, 3D vector.
func NewV3() *V3 { return &V3{} }

// NewV3S creates a new 3D vector using the given scalars.
func NewV3S(x, y, z float64) *V3 { return &V3{x, y, z} }

// NewV4 creates a new, all zero, 4D vector.
func NewV4() *V4 { return &V4{} }

// NewV4S creates a new 3D vector using the given scalars.
func NewV4S(x, y, z, w float64) *V4 { return &V4{x, y, z, w} }
