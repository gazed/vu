// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package lin

// Vector performs vector size 3 or 4 related math needed for 3D graphics.

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
func (v *V3) AeqZ() bool { return v.Dot(v) < EPSILON }

// AeqZ (~=) almost equals zero returns true if the square length of the vector
// is close enough to zero that it makes no difference.
func (v *V4) AeqZ() bool { return v.Dot(v) < EPSILON }

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

// Dot vector v with input vector v1. Same behaviour as V3.Dot()
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
	return v // v is unchanged.
}

// Unit updates vector v such that its length is 1.
// Same behaviour as V3.Unit()
func (v *V4) Unit() *V4 {
	length := v.Len()
	if length != 0 {
		return v.Div(length)
	}
	return v // v is unchanged.
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
// between v1 and v2. Same behaviour as V3.Lerp()
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

// MultVM updates vector v to be the multiplication of row vector rv
// and matrix m. Vector v may be used as the input vector rv.
// The udpated vector v is returned.
//                   [ x0 y0 z0 ]   [ vx' ]
//    [ vx vy vz ] x [ x1 y1 z1 ] = [ vy' ]
//                   [ x2 y2 z2 ]   [ vz' ]
func (v *V3) MultVM(rv *V3, m *M3) *V3 {
	x := rv.X*m.X0 + rv.Y*m.X1 + rv.Z*m.X2
	y := rv.X*m.Y0 + rv.Y*m.Y1 + rv.Z*m.Y2
	z := rv.X*m.Z0 + rv.Y*m.Z1 + rv.Z*m.Z2
	v.X, v.Y, v.Z = x, y, z
	return v
}

// MultVM updates vector v to be the multiplication of row vector rv
// and matrix m. Same behaviour as V4.MultVM().
//                      [ x0 y0 z0 w0 ]   [ vx' ]
//    [ vx vy vz vw ] x [ x1 y1 z1 w1 ] = [ vy' ]
//                      [ x2 y2 z2 w2 ]   [ vz' ]
//                      [ x3 y3 z3 w3 ]   [ vw' ]
func (v *V4) MultVM(rv *V4, m *M4) *V4 {
	x := rv.X*m.X0 + rv.Y*m.X1 + rv.Z*m.X2 + rv.W*m.X3
	y := rv.X*m.Y0 + rv.Y*m.Y1 + rv.Z*m.Y2 + rv.W*m.Y3
	z := rv.X*m.Z0 + rv.Y*m.Z1 + rv.Z*m.Z2 + rv.W*m.Z3
	w := rv.X*m.W0 + rv.Y*m.W1 + rv.Z*m.W2 + rv.W*m.W3
	v.X, v.Y, v.Z, v.W = x, y, z, w
	return v
}

// MultMV updates vector v to be the multiplication of matrix m and
// column vector cv. Vector v may be used as the input vector cv.
// The udpated vector v is returned.
//    [ x0 y0 z0 ]   [ vx ]
//    [ x1 y1 z1 ] x [ vy ] = [ vx' vy' vz' ]
//    [ x2 y2 z2 ]   [ vz ]
func (v *V3) MultMV(m *M3, cv *V3) *V3 {
	x := m.X0*cv.X + m.Y0*cv.Y + m.Z0*cv.Z
	y := m.X1*cv.X + m.Y1*cv.Y + m.Z1*cv.Z
	z := m.X2*cv.X + m.Y2*cv.Y + m.Z2*cv.Z
	v.X, v.Y, v.Z = x, y, z
	return v
}

// MultMV updates vector v to be the multiplication of matrix m and
// column vector cv.  Same behaviour as V3.MultMV().
//    [ x0 y0 z0 w0 ]   [ vx ]
//    [ x1 y1 z1 w1 ] x [ vy ] = [ vx' vy' vz' vw' ]
//    [ x2 y2 z2 w2 ]   [ vz ]
//    [ x3 y3 z3 w3 ]   [ vw ]
func (v *V4) MultMV(m *M4, cv *V4) *V4 {
	x := m.X0*cv.X + m.Y0*cv.Y + m.Z0*cv.Z + m.W0*cv.W
	y := m.X1*cv.X + m.Y1*cv.Y + m.Z1*cv.Z + m.W1*cv.W
	z := m.X2*cv.X + m.Y2*cv.Y + m.Z2*cv.Z + m.W2*cv.W
	w := m.X3*cv.X + m.Y3*cv.Y + m.Z3*cv.Z + m.W3*cv.W
	v.X, v.Y, v.Z, v.W = x, y, z, w
	return v
}

// Col updates vector v to be filled with the matrix elements from the
// indicated column of matrix m. The updated vector v is returned.
// Vector v is only updated for valid column indices.
func (v *V3) Col(index int, m *M3) *V3 {
	switch index {
	case 0:
		return v.SetS(m.X0, m.X1, m.X2)
	case 1:
		return v.SetS(m.Y0, m.Y1, m.Y2)
	case 2:
		return v.SetS(m.Z0, m.Z1, m.Z2)
	}
	return v
}

// Row updates vector v to be filled with the matrix elements from the
// indicated row of matrix m. The updated vector v is returned.
// Vector v is only updated for valid row indices.
func (v *V3) Row(index int, m *M3) *V3 {
	switch index {
	case 0:
		return v.SetS(m.X0, m.Y0, m.Z0)
	case 1:
		return v.SetS(m.X1, m.Y1, m.Z1)
	case 2:
		return v.SetS(m.X2, m.Y2, m.Z2)
	}
	return v
}

// vector-matrix operations
// ============================================================================
// vector-quaternion operations

// MultVQ  updates vector v to be the rotation of vector a by quaternion q.
func (v *V3) MultVQ(a *V3, q *Q) *V3 {
	v.X, v.Y, v.Z = MultSQ(a.X, a.Y, a.Z, q)
	return v
}

// MultSQ applies rotation q to scalar vector (x,y,z)
// The updated scalar vector (vx,vy,vz) is returned.
func MultSQ(x, y, z float64, q *Q) (vx, vy, vz float64) {
	k0 := q.W*q.W - 0.5

	// k1 = Q.V
	k1 := x*q.X + y*q.Y + z*q.Z

	// (qq-1/2)V+(Q.V)Q
	rx := x*k0 + q.X*k1
	ry := y*k0 + q.Y*k1
	rz := z*k0 + q.Z*k1

	// (Q.V)Q+(qq-1/2)V+q(QxV)
	rx += q.W * (q.Y*z - q.Z*y)
	ry += q.W * (q.Z*x - q.X*z)
	rz += q.W * (q.X*y - q.Y*x)

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
