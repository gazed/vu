// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package lin

// Vector performs vector math needed for 3D graphics.
//
// The vector math is performed on vectors of size 3 or 4 as this is all
// that is needed for 3D graphics applications.  Likewise all vectors
// values are float32 since this precision is sufficient/expected for
// 3D graphics applications.

import "math"

// V3 is a 3 element vector.  This can also be used as a point.
type V3 struct {
	X float32 // increments as X moves to the right.
	Y float32 // increments as Y moves up.
	Z float32 // increments as Z moves out of the screen (right handed).
}

// V4 is a 4 element vector.  This can also be used as a point.
// Vectors of length 4 often have W:1.
type V4 struct {
	X float32 // increments as X moves to the right.
	Y float32 // increments as Y moves up.
	Z float32 // increments as Z moves out of the screen (right handed).
	W float32 // fourth dimension makes for nice 3D matrix math.
}

// Clone creates a new copy of a vector v.  This is needed when performing
// vector operations where the original vector values need to be preserved.
func (v *V3) Clone() *V3 { return &V3{v.X, v.Y, v.Z} }

// Clone creates a new copy of a vector v. Same behaviour as V3.Clone()
func (v *V4) Clone() *V4 { return &V4{v.X, v.Y, v.Z, v.W} }

// Add the input vector vec to vector v.  Vector v is updated with
// the addition while the input vector vec is unchanged. The updated vector v
// is returned so that it may be immediately used in another operation.
func (v *V3) Add(vec *V3) *V3 {
	v.X += vec.X
	v.Y += vec.Y
	v.Z += vec.Z
	return v
}

// Add the input vector vec to vector v.  Same behaviour as V3.Add().
func (v *V4) Add(vec *V4) *V4 {
	v.X += vec.X
	v.Y += vec.Y
	v.Z += vec.Z
	v.W += vec.W
	return v
}

// Sub subtracts the input vec from vector v.  Vector v is updated
// with the subtraction while the input vector vec is unchanged.
// The updated vector v is returned so that it may be immediately
// used in another operation.
func (v *V3) Sub(vec *V3) *V3 {
	v.X -= vec.X
	v.Y -= vec.Y
	v.Z -= vec.Z
	return v
}

// Subtract the input vec from vector v. Same behaviour as V3.Sub().
func (v *V4) Sub(vec *V4) *V4 {
	v.X -= vec.X
	v.Y -= vec.Y
	v.Z -= vec.Z
	v.W -= vec.W
	return v
}

// Mult multiplies vector v with input vector vec.  Vector v is updated with
// the result of multiplication while the input vector vec is unchanged.
// The updated vector v is returned so that it may be immediately used
// in another operation.
func (v *V3) Mult(vec *V3) *V3 {
	v.X *= vec.X
	v.Y *= vec.Y
	v.Z *= vec.Z
	return v
}

// Mult multiplies vector v with input vector vec. Same behaviour as V3.Mult().
func (v *V4) Mult(vec *V4) *V4 {
	v.X *= vec.X
	v.Y *= vec.Y
	v.Z *= vec.Z
	v.W *= vec.W
	return v
}

// Scale vector v by a constant value. Vector v is updated
// by multiplying each part by the scaler value.  The updated vector v is
// returned so that it may be immediately used in another operation.
func (v *V3) Scale(scaler float32) *V3 {
	v.X *= scaler
	v.Y *= scaler
	v.Z *= scaler
	return v
}

// Scale a vector by a constant value. Same behaviour as V3.Scale()
func (v *V4) Scale(scaler float32) *V4 {
	v.X *= scaler
	v.Y *= scaler
	v.Z *= scaler
	v.W *= scaler
	return v
}

// Dot vector v with input vector vec. Both vectors v and in are unchanged.
// From wikipedia:
//
//    "This operation can be defined either algebraically or geometrically.
//     Algebraically, it is the sum of the products of the corresponding
//     entries of the two sequences of numbers.
//     Geometrically, it is the product of the magnitudes of the two vectors
//     and the cosine of the angle between them."
func (v *V3) Dot(in *V3) float32 {
	return v.X*in.X + v.Y*in.Y + v.Z*in.Z
}

// Dot vector v with input vector vec. Same behaviour as V3.Dot()
func (v *V4) Dot(vec *V4) float32 {
	return v.X*vec.X + v.Y*vec.Y + v.Z*vec.Z + v.W*vec.W
}

// Len returns the length of vector v. Vector length is the
// square root of the dot product.
func (v *V3) Len() float32 {
	dot := v.X*v.X + v.Y*v.Y + v.Z*v.Z
	return float32(math.Sqrt(float64(dot)))
}

// Len returns the length of vector v. Same behaviour as V3.Length()
func (v *V4) Len() float32 {
	dot := v.X*v.X + v.Y*v.Y + v.Z*v.Z + v.W*v.W
	return float32(math.Sqrt(float64(dot)))
}

// Unit normalizes vector v such that its length is between 0-1.  Vector v
// is updated to have the normalized values. The updated Vector v is
// returned so that it may be immediately used in another operation.
// This method will not do anything if the vector is zero or close
// enough to being unit-length.
func (v *V3) Unit() *V3 {
	dot := v.X*v.X + v.Y*v.Y + v.Z*v.Z
	if !IsZero(dot) && !IsOne(dot) {
		length := 1 / float32(math.Sqrt(float64(dot)))
		v.X = v.X * length
		v.Y = v.Y * length
		v.Z = v.Z * length
	}
	return v
}

// Unit normalizes vector v such that its length is between 0-1. Same behaviour as V3.Unit()
func (v *V4) Unit() *V4 {
	dot := v.X*v.X + v.Y*v.Y + v.Z*v.Z + v.W*v.W
	if !IsZero(dot) && !IsOne(dot) {
		length := float32(math.Sqrt(float64(dot)))
		v.X = v.X / length
		v.Y = v.Y / length
		v.Z = v.Z / length
		v.W = v.W / length
	}
	return v
}

// Lerp returns a new vector that is a fraction of the distance (linear interpolation)
// between the input vectors v and v2. The fraction is expected to be
// between 0 and 1.  Vectors v and v2 are unchanged.
func (v *V3) Lerp(v2 *V3, fraction float32) *V3 {
	return &V3{
		(v2.X-v.X)*fraction + v.X,
		(v2.Y-v.Y)*fraction + v.Y,
		(v2.Z-v.Z)*fraction + v.Z,
	}
}

// Lerp returns a new vector that is the fraction of the distance between the
// vectors v and v2.  Same behaviour as V3.Lerp()
func (v *V4) Lerp(v2 *V4, fraction float32) *V4 {
	return &V4{
		(v2.X-v.X)*fraction + v.X,
		(v2.Y-v.Y)*fraction + v.Y,
		(v2.Z-v.Z)*fraction + v.Z,
		(v2.W-v.W)*fraction + v.W,
	}
}

// Nlerp returns a new normalized vector that is the linerar interpolation between
// v and v2.
//
// See:
//    http://keithmaggio.wordpress.com/2011/02/15/math-magician-lerp-slerp-and-nlerp/
//    http://number-none.com/product/Understanding%20Slerp,%20Then%20Not%20Using%20It/
func (v *V3) Nlerp(v2 *V3, fraction float32) *V3 {
	return v.Lerp(v2, fraction).Unit()
}

// Nlerp has the same behaviour as V3.Nlerp()
func (v *V4) Nlerp(v2 *V4, fraction float32) *V4 {
	return v.Lerp(v2, fraction).Unit()
}

// Cross returns a new vector that is the cross product of vectors v and v2.
// A cross product vector is a vector that is perpendicular to both input
// vectors. This is only meaningfull in 3 (or 7) dimensions.
func (v *V3) Cross(v2 *V3) *V3 {
	return &V3{
		v.Y*v2.Z - v.Z*v2.Y,
		v.Z*v2.X - v.X*v2.Z,
		v.X*v2.Y - v.Y*v2.X,
	}
}

// Dist returns the distance between points p and p2. Both points p and p2
// are unchanged.
func (p *V3) Dist(p2 *V3) float32 {
	diffx := p.X - p2.X
	diffy := p.Y - p2.Y
	diffz := p.Z - p2.Z
	diff := diffx*diffx + diffy*diffy + diffz*diffz
	return float32(math.Sqrt(float64(diff)))
}

// Distance returns the distance between points p and p2.
// Same behaviour as V3.Dist()
func (p *V4) Dist(p2 *V4) float32 {
	diffx := p.X - p2.X
	diffy := p.Y - p2.Y
	diffz := p.Z - p2.Z
	diffw := p.W - p2.W
	diff := diffx*diffx + diffy*diffy + diffz*diffz + diffw*diffw
	return float32(math.Sqrt(float64(diff)))
}

// MultL updates vector v to be the multiplication of vector v and matrix m
// where v is a row vector on the left.
// The input matrix m is unchanged.
//                   [ x0 y0 z0 ]   [ vx' ]
//    [ vx vy vz ] x [ x1 y1 z1 ] = [ vy' ]
//                   [ x2 y2 z2 ]   [ vz' ]
func (v *V3) MultL(m *M3) *V3 {
	x := v.X*m.X0 + v.Y*m.X1 + v.Z*m.X2
	y := v.X*m.Y0 + v.Y*m.Y1 + v.Z*m.Y2
	z := v.X*m.Z0 + v.Y*m.Z1 + v.Z*m.Z2
	v.X, v.Y, v.Z = x, y, z
	return v
}

// MultR updates vector v to be the multiplication of vector v and matrix m
// where v is a column vector on the right.
// The input matrix m is unchanged.
//    [ x0 y0 z0 ]   [ vx ]
//    [ x1 y1 z1 ] x [ vy ] = [ vx' vy' vz' ]
//    [ x2 y2 z2 ]   [ vz ]
func (v *V3) MultR(m *M3) *V3 {
	x := v.X*m.X0 + v.Y*m.Y0 + v.Z*m.Z0
	y := v.X*m.X1 + v.Y*m.Y1 + v.Z*m.Z1
	z := v.X*m.X2 + v.Y*m.Y2 + v.Z*m.Z2
	v.X, v.Y, v.Z = x, y, z
	return v
}

// MultL updates vector v to be the multiplication of vector v and matrix m.
// Same behaviour as V3.MultL()
//                      [ x0 y0 z0 w0 ]   [ vx' ]
//    [ vx vy vz vw ] x [ x1 y1 z1 w1 ] = [ vy' ]
//                      [ x2 y2 z2 w2 ]   [ vz' ]
//                      [ x3 y3 z3 w3 ]   [ vw' ]
func (v *V4) MultL(mat *M4) *V4 {
	x := v.X*mat.X0 + v.Y*mat.X1 + v.Z*mat.X2 + v.W*mat.X3
	y := v.X*mat.Y0 + v.Y*mat.Y1 + v.Z*mat.Y2 + v.W*mat.Y3
	z := v.X*mat.Z0 + v.Y*mat.Z1 + v.Z*mat.Z2 + v.W*mat.Z3
	w := v.X*mat.W0 + v.Y*mat.W1 + v.Z*mat.W2 + v.W*mat.W3
	v.X, v.Y, v.Z, v.W = x, y, z, w
	return v
}

// MultR updates vector v to be the multiplication of vector v and matrix m.
// Same behaviour as V3.MultR()
//    [ x0 y0 z0 w0 ]   [ vx ]
//    [ x1 y1 z1 w1 ] x [ vy ] = [ vx' vy' vz' vw' ]
//    [ x2 y2 z2 w2 ]   [ vz ]
//    [ x3 y3 z3 w3 ]   [ vw ]
func (v *V4) MultR(m *M4) *V4 {
	x := v.X*m.X0 + v.Y*m.Y0 + v.Z*m.Z0 + v.W*m.W0
	y := v.X*m.X1 + v.Y*m.Y1 + v.Z*m.Z1 + v.W*m.W1
	z := v.X*m.X2 + v.Y*m.Y2 + v.Z*m.Z2 + v.W*m.W2
	w := v.X*m.X3 + v.Y*m.Y3 + v.Z*m.Z3 + v.W*m.W3
	v.X, v.Y, v.Z, v.W = x, y, z, w
	return v
}

// MultQ updates vector v to be the multiplication of vector v and
// quaternion q. The updated vector contains a vector where
// where the q-rotation has been applied to v.
// The input quaternion q is unchanged.
func (v *V3) MultQ(q *Q) *V3 {
	k0 := q.W*q.W - 0.5

	// k1 = Q.V
	k1 := v.X * q.X
	k1 += v.Y * q.Y
	k1 += v.Z * q.Z

	// (qq-1/2)V+(Q.V)Q
	rx := v.X*k0 + q.X*k1
	ry := v.Y*k0 + q.Y*k1
	rz := v.Z*k0 + q.Z*k1

	// (Q.V)Q+(qq-1/2)V+q(QxV)
	rx += q.W * (q.Y*v.Z - q.Z*v.Y)
	ry += q.W * (q.Z*v.X - q.X*v.Z)
	rz += q.W * (q.X*v.Y - q.Y*v.X)

	//  2((Q.V)Q+(qq-1/2)V+q(QxV))
	v.X = rx + rx
	v.Y = ry + ry
	v.Z = rz + rz
	return v
}
