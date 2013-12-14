// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package lin

// Matrix functions deal with 3x3 and 4x4 matrices expected to be used
// in 3D transform calculations. See appendix G of OpenGL Red Book for matrix
// algorithms. Unused Matrix methods like rotation are not included
// (use quaternions instead).
//
// Golang stores matrix structures in row-major order. That is structure elements
// are stored sequentially in memory as row1, row2, row3, etc...  This matches
// http://www.opengl.org/archives/resources/faq/technical/transformations.htm#tran0005
// which states:
//   “For programming purposes, OpenGL matrices are 16-value arrays with base
//    vectors laid out contiguously in memory. The translation components occupy
//    the 13th, 14th, and 15th elements of the 16-element matrix, where indices
//    are numbered from 1 to 16"
// This layout allows the entire structure to be passed as a pointer to the
// underlying (C-language) graphics layer.

import (
	"log"
	"math"
)

// M3 is a 3x3 matrix where the matrix elements are individually addressable.
type M3 struct {
	X0, Y0, Z0 float64 // row 1 : indices 0, 1, 2   [00, 01, 02]
	X1, Y1, Z1 float64 // row 2 : indices 3, 4, 5   [10, 11, 12]
	X2, Y2, Z2 float64 // row 3 : indices 6, 7, 8   [20, 21, 22]
}

// M4 is a 4x4 matrix where the matrix elements are individually addressable.
type M4 struct {
	X0, Y0, Z0, W0 float64 // row 1 : indices 0, 1, 2, 3   [00, 01, 02, 03]
	X1, Y1, Z1, W1 float64 // row 2 : indices 4, 5, 6, 7   [10, 11, 12, 13]
	X2, Y2, Z2, W2 float64 // row 3 : indices 8, 9, a, b   [20, 21, 22, 23]
	X3, Y3, Z3, W3 float64 // row 4 : indices c, d, e, f   [30, 31, 32, 33]
}

// M3Z provides a reference zero matrix that can be used
// in calculations. It should never be changed.
var M3Z = &M3{
	0, 0, 0,
	0, 0, 0,
	0, 0, 0}

// M4Z provides a reference zero matrix that can be used
// in calculations. It should never be changed.
var M4Z = &M4{
	0, 0, 0, 0,
	0, 0, 0, 0,
	0, 0, 0, 0,
	0, 0, 0, 0}

// M3I provides a reference identity matrix that can be used
// in calculations. It should never be changed.
var M3I = &M3{
	1, 0, 0,
	0, 1, 0,
	0, 0, 1}

// M4I provides a reference identity matrix that can be used
// in calculations. It should never be changed.
var M4I = &M4{
	1, 0, 0, 0,
	0, 1, 0, 0,
	0, 0, 1, 0,
	0, 0, 0, 1}

// Eq returns true if all the elements in matrix m have the same value
// as the corresponding elements in matrix a.
func (m *M3) Eq(a *M3) bool {
	return true &&
		m.X0 == a.X0 && m.X1 == a.X1 && m.X2 == a.X2 &&
		m.Y0 == a.Y0 && m.Y1 == a.Y1 && m.Y2 == a.Y2 &&
		m.Z0 == a.Z0 && m.Z1 == a.Z1 && m.Z2 == a.Z2
}

// Eq returns true if all the elements in matrix m have the same value
// as the corresponding elements in matrix a.
func (m *M4) Eq(a *M4) bool {
	return true &&
		m.X0 == a.X0 && m.X1 == a.X1 && m.X2 == a.X2 && m.X3 == a.X3 &&
		m.Y0 == a.Y0 && m.Y1 == a.Y1 && m.Y2 == a.Y2 && m.Y3 == a.Y3 &&
		m.Z0 == a.Z0 && m.Z1 == a.Z1 && m.Z2 == a.Z2 && m.Z3 == a.Z3 &&
		m.W0 == a.W0 && m.W1 == a.W1 && m.W2 == a.W2 && m.W3 == a.W3
}

// Aeq (~=) almost equals returns true if all the elements in matrix m have
// essentially the same value as the corresponding elements in matrix a.
// Used where equals is unlikely to return true due to float precision.
func (m *M3) Aeq(a *M3) bool {
	return true &&
		Aeq(m.X0, a.X0) && Aeq(m.X1, a.X1) && Aeq(m.X2, a.X2) &&
		Aeq(m.Y0, a.Y0) && Aeq(m.Y1, a.Y1) && Aeq(m.Y2, a.Y2) &&
		Aeq(m.Z0, a.Z0) && Aeq(m.Z1, a.Z1) && Aeq(m.Z2, a.Z2)
}

// Aeq (~=) almost equals returns true if all the elements in matrix m have
// essentially the same value as the corresponding elements in matrix a.
// Same as M3.Aeq().
func (m *M4) Aeq(a *M4) bool {
	return true &&
		Aeq(m.X0, a.X0) && Aeq(m.X1, a.X1) && Aeq(m.X2, a.X2) && Aeq(m.X3, a.X3) &&
		Aeq(m.Y0, a.Y0) && Aeq(m.Y1, a.Y1) && Aeq(m.Y2, a.Y2) && Aeq(m.Y3, a.Y3) &&
		Aeq(m.Z0, a.Z0) && Aeq(m.Z1, a.Z1) && Aeq(m.Z2, a.Z2) && Aeq(m.Z3, a.Z3) &&
		Aeq(m.W0, a.W0) && Aeq(m.W1, a.W1) && Aeq(m.W2, a.W2) && Aeq(m.W3, a.W3)
}

// SetS (=) explicitly sets the elements values of each of the matrix values.
// The source matrix a is unchanged. The updated matrix m is returned.
func (m *M3) SetS(x0, y0, z0, x1, y1, z1, x2, y2, z2 float64) *M3 {
	m.X0, m.Y0, m.Z0 = x0, y0, z0
	m.X1, m.Y1, m.Z1 = x1, y1, z1
	m.X2, m.Y2, m.Z2 = x2, y2, z2
	return m
}

// Set (=) assigns all the elements values from matrix a to the
// corresponding element values in matrix m.
// The source matrix a is unchanged. The updated matrix m is returned.
func (m *M3) Set(a *M3) *M3 {
	m.X0, m.Y0, m.Z0 = a.X0, a.Y0, a.Z0
	m.X1, m.Y1, m.Z1 = a.X1, a.Y1, a.Z1
	m.X2, m.Y2, m.Z2 = a.X2, a.Y2, a.Z2
	return m
}

// SetM4 (=) updates calling matrix m to be the 3x3 matrix from the top left
// corner of the given 4x4 matrix m4.  The source matrix a is unchanged.
// The updated matrix m is returned.
//    [ x0 y0 z0 w0 ]    [ x0 y0 z0 ]
//    [ x1 y1 z1 w1 ] => [ x1 y1 z1 ]
//    [ x2 y2 z2 w2 ]    [ x2 y2 z2 ]
//    [ x3 y3 z3 w3 ]
func (m *M3) SetM4(a *M4) *M3 {
	m.X0, m.Y0, m.Z0 = a.X0, a.Y0, a.Z0
	m.X1, m.Y1, m.Z1 = a.X1, a.Y1, a.Z1
	m.X2, m.Y2, m.Z2 = a.X2, a.Y2, a.Z2
	return m
}

// Set (=) assigns all the elements values from matrix a to the
// corresponding element values in matrix m.  The source matrix a is unchanged.
// The updated matrix m is returned.
func (m *M4) Set(a *M4) *M4 {
	m.X0, m.Y0, m.Z0, m.W0 = a.X0, a.Y0, a.Z0, a.W0
	m.X1, m.Y1, m.Z1, m.W1 = a.X1, a.Y1, a.Z1, a.W1
	m.X2, m.Y2, m.Z2, m.W2 = a.X2, a.Y2, a.Z2, a.W2
	m.X3, m.Y3, m.Z3, m.W3 = a.X3, a.Y3, a.Z3, a.W3
	return m
}

// Abs updates m to be the the absolute (non-negative) element values of
// the corresponding element values in matrix a.  The source matrix a is unchanged.
// The updated matrix m is returned.
func (m *M3) Abs(a *M3) *M3 {
	m.X0, m.Y0, m.Z0 = math.Abs(a.X0), math.Abs(a.Y0), math.Abs(a.Z0)
	m.X1, m.Y1, m.Z1 = math.Abs(a.X1), math.Abs(a.Y1), math.Abs(a.Z1)
	m.X2, m.Y2, m.Z2 = math.Abs(a.X2), math.Abs(a.Y2), math.Abs(a.Z2)
	return m
}

// Transpose updates m to be the reflection of matrix a over its diagonal.
// This essentially changes column major order to row major order
// or vice-versa.
//    [ x0 y0 z0 ]    [ x0 x1 x2 ]
//    [ x1 y1 z1 ] => [ y0 y1 y2 ]
//    [ x2 y2 z2 ]    [ z0 z1 z2 ]
// The input matrix a is not changed. Matrix m may be used as the input parameter.
// The updated matrix m is returned.
func (m *M3) Transpose(a *M3) *M3 {
	t_y0, t_z0, t_z1 := a.Y0, a.Z0, a.Z1
	m.X0, m.Y0, m.Z0 = a.X0, a.X1, a.X2
	m.X1, m.Y1, m.Z1 = t_y0, a.Y1, a.Y2
	m.X2, m.Y2, m.Z2 = t_z0, t_z1, a.Z2
	return m
}

// Transpose updates m to be the reflection of matrix a over its diagonal.
//    [ x0 y0 z0 w0 ]    [ x0 x1 x2 x3 ]
//    [ x1 y1 z1 w1 ] => [ y0 y1 y2 y3 ]
//    [ x2 y2 z2 w2 ]    [ z0 z1 z2 z3 ]
//    [ x3 y3 z3 w3 ]    [ w0 w1 w2 w3 ]
// Same behaviour as M3.Transpose()
func (m *M4) Transpose(a *M4) *M4 {
	t_y0, t_z0, t_w0 := a.Y0, a.Z0, a.W0
	t_z1, t_w1, t_w2 := a.Z1, a.W1, a.W2
	m.X0, m.Y0, m.Z0, m.W0 = a.X0, a.X1, a.X2, a.X3
	m.X1, m.Y1, m.Z1, m.W1 = t_y0, a.Y1, a.Y2, a.Y3
	m.X2, m.Y2, m.Z2, m.W2 = t_z0, t_z1, a.Z2, a.Z3
	m.X3, m.Y3, m.Z3, m.W3 = t_w0, t_w1, t_w2, a.W3
	return m
}

// Add (+) adds matrices a and b storing the results in m.
// Each element of matrix b is added to the corresponding matrix a element.
// It is safe to use the calling matrix m as one or both of the parameters.
// For example the plus.equals operation (+=) is
//     m.Add(m, b)
// The updated matrix m is returned.
func (m *M3) Add(a, b *M3) *M3 {
	m.X0, m.Y0, m.Z0 = a.X0+b.X0, a.Y0+b.Y0, a.Z0+b.Z0
	m.X1, m.Y1, m.Z1 = a.X1+b.X1, a.Y1+b.Y1, a.Z1+b.Z1
	m.X2, m.Y2, m.Z2 = a.X2+b.X2, a.Y2+b.Y2, a.Z2+b.Z2
	return m
}

// Sub (-) subtracts matrices b from a storing the results in m.
// Each element of matrix b is subtracted from the corresponding matrix a element.
// It is safe to use the calling matrix m as one or both of the parameters.
// For example the minus.equals operation (-=) is
//     m.Sub(m, b)
// The updated matrix m is returned.
func (m *M3) Sub(a, b *M3) *M3 {
	m.X0, m.Y0, m.Z0 = a.X0-b.X0, a.Y0-b.Y0, a.Z0-b.Z0
	m.X1, m.Y1, m.Z1 = a.X1-b.X1, a.Y1-b.Y1, a.Z1-b.Z1
	m.X2, m.Y2, m.Z2 = a.X2-b.X2, a.Y2-b.Y2, a.Z2-b.Z2
	return m
}

// Mult (*) multiplies matrices l and r storing the results in m.
//    [ lx0 ly0 lz0 ] [ rx0 ry0 rz0 ]    [ mx0 my0 mz0 ]
//    [ lx1 ly1 lz1 ]x[ rx1 ry1 rz1 ] => [ mx1 my1 mz1 ]
//    [ lx2 ly2 lz2 ] [ rx2 ry2 rz2 ]    [ mx2 my2 mz2 ]
// It is safe to use the calling matrix m as one or both of the parameters.
// For example (*=) is
//     m.Mult(m, r)
// The updated matrix m is returned.
func (m *M3) Mult(l, r *M3) *M3 {
	x0 := l.X0*r.X0 + l.Y0*r.X1 + l.Z0*r.X2
	y0 := l.X0*r.Y0 + l.Y0*r.Y1 + l.Z0*r.Y2
	z0 := l.X0*r.Z0 + l.Y0*r.Z1 + l.Z0*r.Z2
	x1 := l.X1*r.X0 + l.Y1*r.X1 + l.Z1*r.X2
	y1 := l.X1*r.Y0 + l.Y1*r.Y1 + l.Z1*r.Y2
	z1 := l.X1*r.Z0 + l.Y1*r.Z1 + l.Z1*r.Z2
	x2 := l.X2*r.X0 + l.Y2*r.X1 + l.Z2*r.X2
	y2 := l.X2*r.Y0 + l.Y2*r.Y1 + l.Z2*r.Y2
	z2 := l.X2*r.Z0 + l.Y2*r.Z1 + l.Z2*r.Z2
	m.X0, m.Y0, m.Z0 = x0, y0, z0
	m.X1, m.Y1, m.Z1 = x1, y1, z1
	m.X2, m.Y2, m.Z2 = x2, y2, z2
	return m
}

// Mult updates matrix m to be the multiplication of input matrices l, r.
//    [ lx0 ly0 lz0 lw0 ] [ rx0 ry0 rz0 rw0 ]    [ mx0 my0 mz0 mw0 ]
//    [ lx1 ly1 lz1 lw1 ]x[ rx1 ry1 rz1 rw1 ] => [ mx1 my1 mz1 mw1 ]
//    [ lx2 ly2 lz2 lw2 ] [ rx2 ry2 rz2 rw2 ]    [ mx2 my2 mz2 mw2 ]
//    [ lx3 ly3 lz3 lw3 ] [ rx3 ry3 rz3 rw3 ]    [ mx3 my3 mz3 mw3 ]
// Same behaviour as M3.Mult()
func (m *M4) Mult(l, r *M4) *M4 {
	x0 := l.X0*r.X0 + l.Y0*r.X1 + l.Z0*r.X2 + l.W0*r.X3
	y0 := l.X0*r.Y0 + l.Y0*r.Y1 + l.Z0*r.Y2 + l.W0*r.Y3
	z0 := l.X0*r.Z0 + l.Y0*r.Z1 + l.Z0*r.Z2 + l.W0*r.Z3
	w0 := l.X0*r.W0 + l.Y0*r.W1 + l.Z0*r.W2 + l.W0*r.W3
	x1 := l.X1*r.X0 + l.Y1*r.X1 + l.Z1*r.X2 + l.W1*r.X3
	y1 := l.X1*r.Y0 + l.Y1*r.Y1 + l.Z1*r.Y2 + l.W1*r.Y3
	z1 := l.X1*r.Z0 + l.Y1*r.Z1 + l.Z1*r.Z2 + l.W1*r.Z3
	w1 := l.X1*r.W0 + l.Y1*r.W1 + l.Z1*r.W2 + l.W1*r.W3
	x2 := l.X2*r.X0 + l.Y2*r.X1 + l.Z2*r.X2 + l.W2*r.X3
	y2 := l.X2*r.Y0 + l.Y2*r.Y1 + l.Z2*r.Y2 + l.W2*r.Y3
	z2 := l.X2*r.Z0 + l.Y2*r.Z1 + l.Z2*r.Z2 + l.W2*r.Z3
	w2 := l.X2*r.W0 + l.Y2*r.W1 + l.Z2*r.W2 + l.W2*r.W3
	x3 := l.X3*r.X0 + l.Y3*r.X1 + l.Z3*r.X2 + l.W3*r.X3
	y3 := l.X3*r.Y0 + l.Y3*r.Y1 + l.Z3*r.Y2 + l.W3*r.Y3
	z3 := l.X3*r.Z0 + l.Y3*r.Z1 + l.Z3*r.Z2 + l.W3*r.Z3
	w3 := l.X3*r.W0 + l.Y3*r.W1 + l.Z3*r.W2 + l.W3*r.W3
	m.X0, m.Y0, m.Z0, m.W0 = x0, y0, z0, w0
	m.X1, m.Y1, m.Z1, m.W1 = x1, y1, z1, w1
	m.X2, m.Y2, m.Z2, m.W2 = x2, y2, z2, w2
	m.X3, m.Y3, m.Z3, m.W3 = x3, y3, z3, w3
	return m
}

// MultLtR multiplies the transpose of matrix l on left of matrix r
// and stores the result in m. This can be used for inverse transforms.
//    [ lx0 lx1 lx0 ] [ rx0 ry0 rz0 ]    [ mx0 my0 mz0 ]
//    [ ly0 ly1 ly2 ]x[ rx1 ry1 rz1 ] => [ mx1 my1 mz1 ]
//    [ lz0 lz1 lz2 ] [ rx2 ry2 rz2 ]    [ mx2 my2 mz2 ]
// It is safe to use the calling matrix m as one or both of the parameters.
// The updated matrix m is returned.
func (m *M3) MultLtR(lt, r *M3) *M3 {
	x0 := lt.X0*r.X0 + lt.X1*r.X1 + lt.X2*r.X2
	y0 := lt.X0*r.Y0 + lt.X1*r.Y1 + lt.X2*r.Y2
	z0 := lt.X0*r.Z0 + lt.X1*r.Z1 + lt.X2*r.Z2
	x1 := lt.Y0*r.X0 + lt.Y1*r.X1 + lt.Y2*r.X2
	y1 := lt.Y0*r.Y0 + lt.Y1*r.Y1 + lt.Y2*r.Y2
	z1 := lt.Y0*r.Z0 + lt.Y1*r.Z1 + lt.Y2*r.Z2
	x2 := lt.Z0*r.X0 + lt.Z1*r.X1 + lt.Z2*r.X2
	y2 := lt.Z0*r.Y0 + lt.Z1*r.Y1 + lt.Z2*r.Y2
	z2 := lt.Z0*r.Z0 + lt.Z1*r.Z1 + lt.Z2*r.Z2
	m.X0, m.Y0, m.Z0 = x0, y0, z0
	m.X1, m.Y1, m.Z1 = x1, y1, z1
	m.X2, m.Y2, m.Z2 = x2, y2, z2
	return m
}

// TranslateTM updates m to be the multiplication of a translation matrix
// T created from x, y, z, and itself. The updated matrix m is returned.
//    [ 1 0 0 0 ]   [ x0 y0 z0 w0 ]     [ x0  y0  z0  w0 ]
//    [ 0 1 0 0 ] x [ x1 y1 z1 w1 ]  => [ x1  y1  z1  w1 ]
//    [ 0 0 1 0 ]   [ x2 y2 z2 w2 ]     [ x2  y2  z2  w2 ]
//    [ x y z 1 ]   [ x3 y3 z3 w3 ]     [ x3' y3' z3' w3']
// Be sure to pick the correct translate (TM or MT) when doing transforms.
// Generally its TranslateMT since translate is the last transform
// (given that M4 is in row major order).
func (m *M4) TranslateTM(x, y, z float64) *M4 {
	x3 := x*m.X0 + y*m.X1 + z*m.X2 + m.X3
	y3 := x*m.Y0 + y*m.Y1 + z*m.Y2 + m.Y3
	z3 := x*m.Z0 + y*m.Z1 + z*m.Z2 + m.Z3
	w3 := x*m.W0 + y*m.W1 + z*m.W2 + m.W3
	m.X3, m.Y3, m.Z3, m.W3 = x3, y3, z3, w3
	return m
}

// TranslateMT updates m to be the multiplication of itself
// and a translation matrix created from x, y, z.
// The updated matrix m is returned.
//    [ x0 y0 z0 w0 ]   [ 1 0 0 0 ]    [ x0'  y0' z0' w0 ]
//    [ x1 y1 z1 w1 ] x [ 0 1 0 0 ] => [ x1'  y1' z1' w1 ]
//    [ x2 y2 z2 w2 ]   [ 0 0 1 0 ]    [ x2'  y2' z2' w2 ]
//    [ x3 y3 z3 w3 ]   [ x y z 1 ]    [ x3'  y3' z3' w3 ]
// Be sure to pick the correct translate (TM or MT) when doing transforms.
// Generally its TranslateMT since translate is the last transform
// (given that M4 is in row major order).
func (m *M4) TranslateMT(x, y, z float64) *M4 {
	m.X0, m.Y0, m.Z0 = m.X0+m.W0*x, m.Y0+m.W0*y, m.Z0+m.W0*z
	m.X1, m.Y1, m.Z1 = m.X1+m.W1*x, m.Y1+m.W1*y, m.Z1+m.W1*z
	m.X2, m.Y2, m.Z2 = m.X2+m.W2*x, m.Y2+m.W2*y, m.Z2+m.W2*z
	m.X3, m.Y3, m.Z3 = m.X3+m.W3*x, m.Y3+m.W3*y, m.Z3+m.W3*z
	return m
}

// Scale (*) each element of matrix m by the given scalar.
// The updated matrix m is returned.
func (m *M3) Scale(s float64) *M3 {
	m.X0, m.Y0, m.Z0 = m.X0*s, m.Y0*s, m.Z0*s
	m.X1, m.Y1, m.Z1 = m.X1*s, m.Y1*s, m.Z1*s
	m.X2, m.Y2, m.Z2 = m.X2*s, m.Y2*s, m.Z2*s
	return m
}

// ScaleS (*) scales each column of matrix m using the corresponding vector
// elements for x, y, z.  The updated matrix m is returned.
func (m *M3) ScaleS(x, y, z float64) *M3 {
	m.X0, m.Y0, m.Z0 = m.X0*x, m.Y0*y, m.Z0*z
	m.X1, m.Y1, m.Z1 = m.X1*x, m.Y1*y, m.Z1*z
	m.X2, m.Y2, m.Z2 = m.X2*x, m.Y2*y, m.Z2*z
	return m
}

// ScaleV (*) scales each column of matrix m using the given vector v
// for elements for x, y, z.  The updated matrix m is returned.
func (m *M3) ScaleV(v *V3) *M3 {
	m.X0, m.Y0, m.Z0 = m.X0*v.X, m.Y0*v.Y, m.Z0*v.Z
	m.X1, m.Y1, m.Z1 = m.X1*v.X, m.Y1*v.Y, m.Z1*v.Z
	m.X2, m.Y2, m.Z2 = m.X2*v.X, m.Y2*v.Y, m.Z2*v.Z
	return m
}

// ScaleSM updates m to be the multiplication of a scale matrix
// created from x, y, z and itself. The updated matrix m is
// returned so that it may be immediately used in another operation.
//    [ x 0 0 0 ]   [ x0 y0 z0 w0 ]    [ x0' y0' z0' w0' ]
//    [ 0 y 0 0 ] x [ x1 y1 z1 w1 ] => [ x1' y1' z1' w1' ]
//    [ 0 0 z 0 ]   [ x2 y2 z2 w2 ]    [ x2' y2' z2' w2' ]
//    [ 0 0 0 1 ]   [ x3 y3 z3 w3 ]    [ x3  y3  z3  w3  ]
// Be sure to pick the correct scale (SM or MS) when doing transforms.
// Generally its ScaleSM since scale is the first transform on the left
// (given that M4 is in row major order).
func (m *M4) ScaleSM(x, y, z float64) *M4 {
	m.X0, m.Y0, m.Z0, m.W0 = m.X0*x, m.Y0*x, m.Z0*x, m.W0*x
	m.X1, m.Y1, m.Z1, m.W1 = m.X1*y, m.Y1*y, m.Z1*y, m.W1*y
	m.X2, m.Y2, m.Z2, m.W2 = m.X2*z, m.Y2*z, m.Z2*z, m.W2*z
	return m
}

// ScaleMS updates m to be the multiplication of m and a scale matrix created
// from x, y, z. The updated matrix m is returned so that it may be immediately
// used in another operation.
//    [ x0 y0 z0 w0 ]   [ x 0 0 0 ]    [ x0' y0' z0' w0 ]
//    [ x1 y1 z1 w1 ] x [ 0 y 0 0 ] => [ x1' y1' z1' w1 ]
//    [ x2 y2 z2 w2 ]   [ 0 0 z 0 ]    [ x2' y2' z2' w2 ]
//    [ x3 y3 z3 w3 ]   [ 0 0 0 1 ]    [ x3' y3' z3' w3 ]
// Be sure to pick the correct scale (SM or MS) when doing transforms.
// Generally its ScaleSM since scale is the first transform on the left
// (given that M4 is in row major order).
func (m *M4) ScaleMS(x, y, z float64) *M4 {
	m.X0, m.Y0, m.Z0 = m.X0*x, m.Y0*y, m.Z0*z
	m.X1, m.Y1, m.Z1 = m.X1*x, m.Y1*y, m.Z1*z
	m.X2, m.Y2, m.Z2 = m.X2*x, m.Y2*y, m.Z2*z
	m.X3, m.Y3, m.Z3 = m.X3*x, m.Y3*y, m.Z3*z
	return m
}

// SetQ converts a quaternion rotation representation to a matrix
// rotation representation. SetQ updates matrix m to be the rotation
// matrix representing the rotation described by unit-quaternion q.
//                       [ x0 y0 z0 ]
//    [ qx qy qz qw ] => [ x1 y1 z1 ]
//                       [ x2 y2 z2 ]
// The parameter q is unchanged. The updated matrix m is returned.
func (m *M3) SetQ(q *Q) *M3 {
	xx, yy, zz := q.X*q.X, q.Y*q.Y, q.Z*q.Z
	xy, xz, yz := q.X*q.Y, q.X*q.Z, q.Y*q.Z
	wx, wy, wz := q.W*q.X, q.W*q.Y, q.W*q.Z
	m.X0, m.Y0, m.Z0 = 1-2*(yy+zz), 2*(xy-wz), 2*(xz+wy)
	m.X1, m.Y1, m.Z1 = 2*(xy+wz), 1-2*(xx+zz), 2*(yz-wx)
	m.X2, m.Y2, m.Z2 = 2*(xz-wy), 2*(yz+wx), 1-2*(xx+yy)
	return m
}

// SetQ converts a quaternion rotation representation to a matrix
// rotation representation. SetQ updates matrix m to be the rotation
// matrix representing the rotation described by unit-quaternion q.
//                       [ x0 y0 z0 0 ]
//    [ qx qy qz qw ] => [ x1 y1 z1 0 ]
//                       [ x2 y2 z2 0 ]
//                       [  0  0  0 1 ]
// The parameter q is unchanged. The updated matrix m is returned.
func (m *M4) SetQ(q *Q) *M4 {
	xx, yy, zz := q.X*q.X, q.Y*q.Y, q.Z*q.Z
	xy, xz, yz := q.X*q.Y, q.X*q.Z, q.Y*q.Z
	wx, wy, wz := q.W*q.X, q.W*q.Y, q.W*q.Z
	m.X0, m.Y0, m.Z0, m.W0 = 1-2*(yy+zz), 2*(xy-wz), 2*(xz+wy), 0
	m.X1, m.Y1, m.Z1, m.W1 = 2*(xy+wz), 1-2*(xx+zz), 2*(yz-wx), 0
	m.X2, m.Y2, m.Z2, m.W2 = 2*(xz-wy), 2*(yz+wx), 1-2*(xx+yy), 0
	m.X3, m.Y3, m.Z3, m.W3 = 0, 0, 0, 1
	return m
}

// SetSkewSym sets the matrix m to be a skew-symetric matrix based
// on the elements of vector v. Wikipedia states:
//    "A skew-symmetric matrix is a square matrix whose transpose is
//     also its negative."
func (m *M3) SetSkewSym(v *V3) *M3 {
	m.X0, m.Y0, m.Z0 = 0, -v.Z, v.Y
	m.X1, m.Y1, m.Z1 = v.Z, 0, -v.X
	m.X2, m.Y2, m.Z2 = -v.Y, v.X, 0
	return m
}

// Det returns the determinant of matrix m. Determinants are helpful
// when calculating the inverse of transform matrices. Wikipedia states:
//    "The determinant provides important information about [..] a matrix that
//     corresponds to a linear transformation of a vector space [..] the transformation
//     has an inverse operation exactly when the determinant is nonzero."
func (m *M3) Det() float64 {
	return m.X0*(m.Y1*m.Z2-m.Z1*m.Y2) + m.Y0*(m.Z1*m.X2-m.X1*m.Z2) + m.Z0*(m.X1*m.Y2-m.Y1*m.X2)
}

// Cof returns one of the possible cofactors of a 3x3 matrix given the
// input minor (the row and column removed from the calculation).
// Wikipedia states:
//      "cofactors [...] are useful for computing both the determinant
//       and inverse of square matrices".
func (m *M3) Cof(row, col int) float64 {
	minor := row*10 + col // minor given by the removed row and column.
	switch minor {
	case 00:
		return m.Y1*m.Z2 - m.Z1*m.Y2
	case 01:
		return m.Z1*m.X2 - m.X1*m.Z2 // flip to negate.
	case 02:
		return m.X1*m.Y2 - m.Y1*m.X2
	case 10:
		return m.Z0*m.Y2 - m.Y0*m.Z2 // flip to negate.
	case 11:
		return m.X0*m.Z2 - m.Z0*m.X2
	case 12:
		return m.Y0*m.X2 - m.X0*m.Y2 // flip to negate.
	case 20:
		return m.Y0*m.Z1 - m.Z0*m.Y1
	case 21:
		return m.Z0*m.X1 - m.X0*m.Z1 // flip to negate.
	case 22:
		return m.X0*m.Y1 - m.Y0*m.X1
	}
	log.Printf("matrix M3.Cof developer error %d", minor)
	return 0
}

// Adj updates m to be the adjoint matrix of matrix a.  The adjoint matrix is
// created by the transpose of the cofactor matrix of the original matrix.
//     [ a.cof(0,0) a.cof(1,0) a.cof(2,0) ]    [ x0 y0 z0 ]
//     [ a.cof(0,1) a.cof(1,1) a.cof(2,1) ] => [ x1 y1 z1 ]
//     [ a.cof(0,2) a.cof(1,2) a.cof(2,2) ]    [ x2 y2 z2 ]
// The updated matrix m is returned.
func (m *M3) Adj(a *M3) *M3 {
	x0, y0, z0 := a.Cof(0, 0), a.Cof(1, 0), a.Cof(2, 0)
	x1, y1, z1 := a.Cof(0, 1), a.Cof(1, 1), a.Cof(2, 1)
	x2, y2, z2 := a.Cof(0, 2), a.Cof(1, 2), a.Cof(2, 2)
	m.X0, m.Y0, m.Z0 = x0, y0, z0
	m.X1, m.Y1, m.Z1 = x1, y1, z1
	m.X2, m.Y2, m.Z2 = x2, y2, z2
	return m
}

// Inv updates m to be the inverse of matrix a. The updated matrix m is returned.
// Matrix m is not updated if the matrix has no inverse.
func (m *M3) Inv(a *M3) *M3 {
	det := a.Det()
	if det != 0 {
		s := 1 / det
		x0, y0, z0 := a.Cof(0, 0)*s, a.Cof(1, 0)*s, a.Cof(2, 0)*s
		x1, y1, z1 := a.Cof(0, 1)*s, a.Cof(1, 1)*s, a.Cof(2, 1)*s
		x2, y2, z2 := a.Cof(0, 2)*s, a.Cof(1, 2)*s, a.Cof(2, 2)*s
		m.X0, m.Y0, m.Z0 = x0, y0, z0
		m.X1, m.Y1, m.Z1 = x1, y1, z1
		m.X2, m.Y2, m.Z2 = x2, y2, z2
	}
	return m
}

// SetAa updates m to be a rotation matrix from the given axis (ax, ay, az)
// and angle (in radians). See:
//    http://en.wikipedia.org/wiki/Rotation_matrix#Rotation_matrix_from_axis_and_angle
//    http://web.archive.org/web/20041029003853/...
//    ...http://www.j3d.org/matrix_faq/matrfaq_latest.html#Q38 (*note column order)
// The updated matrix m is returned.
func (m *M3) SetAa(ax, ay, az, ang float64) *M3 {
	alenSqr := ax*ax + ay*ay + az*az
	if alenSqr == 0 {
		log.Printf("quaternion.Q.SetAa Zero length axis.")
		return m
	}

	// ensure normalized unit vector.
	ilen := 1 / math.Sqrt(alenSqr)
	ax, ay, az = ax*ilen, ay*ilen, az*ilen

	// now set the rotation.
	rcos, rsin := math.Cos(ang), math.Sin(ang)
	m.X0 = rcos + ax*ax*(1-rcos)
	m.Y0 = -az*rsin + ay*ax*(1-rcos)
	m.Z0 = ay*rsin + az*ax*(1-rcos)
	m.X1 = az*rsin + ax*ay*(1-rcos)
	m.Y1 = rcos + ay*ay*(1-rcos)
	m.Z1 = -ax*rsin + az*ay*(1-rcos)
	m.X2 = -ay*rsin + ax*az*(1-rcos)
	m.Y2 = ax*rsin + ay*az*(1-rcos)
	m.Z2 = rcos + az*az*(1-rcos)
	return m
}

// methods above do not allocate memory.
// ============================================================================
// convenience functions for allocating matrices. Nothing else should allocate.

// NewM3 creates a new, all zero, 3x3 matrix.
func NewM3() *M3 { return &M3{} }

// NewM4 creates a new, all zero, 4x4 matrix.
func NewM4() *M4 { return &M4{} }

// NewM3I creates a new 3x3 identity matrix.
//    [ x0 y0 z0 ]   [ 1 0 0 ]
//    [ x1 y1 z1 ] = [ 0 1 0 ]
//    [ x2 y2 z2 ]   [ 0 0 1 ]
func NewM3I() *M3 { return &M3{X0: 1, Y1: 1, Z2: 1} }

// NewM4I creates a new 4x4 identity matrix.
//    [ x0 y0 z0 w0 ]   [ 1 0 0 0 ]
//    [ x1 y1 z1 w1 ] = [ 0 1 0 0 ]
//    [ x2 y2 z2 w2 ]   [ 0 0 1 0 ]
//    [ x3 y3 z3 w3 ]   [ 0 0 0 1 ]
func NewM4I() *M4 { return &M4{X0: 1, Y1: 1, Z2: 1, W3: 1} }

// NewOrtho creates a new 4x4 matrix with projection values needed to
// transform a 3 dimensional model to a 2 dimensional plane.
// Orthographic projection ignores depth. The input arguments are:
//     left, right:  Vertical clipping planes.
//     bottom, top:  Horizontal clipping planes.
//     near, far  :  Depth clipping planes. The depth values are
//                   negative if the plane is to be behind the viewer
// An orthographic matrix fills the following matrix locations:
//    [ a 0 0 0 ]
//    [ 0 b 0 0 ]
//    [ 0 0 c 0 ]
//    [ d e f 1 ]
func NewOrtho(left, right, bottom, top, near, far float64) *M4 {
	m := &M4{}
	m.X0 = 2 / (right - left)
	m.Y0 = 0
	m.Z0 = 0
	m.W0 = 0
	m.X1 = 0
	m.Y1 = 2 / (top - bottom)
	m.Z1 = 0
	m.W1 = 0
	m.X2 = 0
	m.Y2 = 0
	m.Z2 = -2 / (far - near)
	m.W2 = 0
	m.X3 = -(right + left) / (right - left)
	m.Y3 = -(top + bottom) / (top - bottom)
	m.Z3 = -(far + near) / (far - near)
	m.W3 = 1
	return m
}

// NewPersp creates a new 4x4 matrix with projection values needed to
// transform a 3 dimentional model to a 2 dimensional plane.
// Objects that are further away from the viewer will appear smaller.
// The input arguments are:
//    fov        An amount in degrees indicating how much of the
//               scene is visible.
//    aspect     The ratio of height to width of the model.
//    near, far  The depth clipping planes. The depth values are
//               negative if the plane is to be behind the viewer
// A perspective matrix fills the following matrix locations:
//    [ a 0 0 0 ]
//    [ 0 b 0 0 ]
//    [ 0 0 c d ]
//    [ 0 0 e 0 ]
func NewPersp(fov, aspect, near, far float64) *M4 {
	m := &M4{}
	f := 1 / float64(math.Tan(Rad(fov)*0.5))
	m.X0 = f / aspect
	m.X1 = 0
	m.X2 = 0
	m.X3 = 0
	m.Y0 = 0
	m.Y1 = f
	m.Y2 = 0
	m.Y3 = 0
	m.Z0 = 0
	m.Z1 = 0
	m.Z2 = (far + near) / (near - far)
	m.Z3 = 2 * far * near / (near - far)
	m.W0 = 0
	m.W1 = 0
	m.W2 = -1
	m.W3 = 0
	return m
}

// NewPerspInv creates a new inverse matrix of the given perspective
// matrix values (see NewPersp()).
//   [ a' 0  0  0 ] where a' = 1/a       d' = 1/e
//   [ 0  b' 0  0 ]       b' = 1/b       e' = 1/d
//   [ 0  0  0  d']       c' = -(c/de)
//   [ 0  0  e' c']
// This is used when going from screen x,y coordinates to 3D coordinates.
// as in the case when creating a picking ray from a mouse location.
func NewPerspInv(fov, aspect, near, far float64) *M4 {
	m := &M4{}
	f := float64(math.Tan(Rad(fov) * 0.5))
	c := 2 * far * near / (near - far)
	m.X0 = f * aspect
	m.X1 = 0
	m.X2 = 0
	m.X3 = 0
	m.Y0 = 0
	m.Y1 = f
	m.Y2 = 0
	m.Y3 = 0
	m.Z0 = 0
	m.Z1 = 0
	m.Z2 = 0
	m.Z3 = -1
	m.W0 = 0
	m.W1 = 0
	m.W2 = 1 / c
	m.W3 = -((far + near) / (near - far) / (-1 * c))
	return m
}
