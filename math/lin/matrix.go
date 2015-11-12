// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package lin

// Matrix functions deal with 3x3 and 4x4 matrices expected to be used
// in CPU 3D transform or physics calculations. An example of CPU math is
// providing precalulated per-frame transform matricies to the GPU rather than
// having the GPU calculate identical per-vertex or per-fragment matricies.
// Large scale, time-critical, repetitive math operations are expected to use
// a GPGPU based package, ie. OpenCL.
//
// Note that this matrix implementation does not attempt to be all inclusive.
// Unused matrix methods, like rotation, are excluded since rotations are
// tracked using quaternions.
//
// Row or Column Major order? No matter the convention, the end result of a
// vector point (x, y, z, 1) multiplied with a transform matrix must be:
//   x' = x*Xx + y*Yx + z*Zx + Tx
//   y' = x*Xy + y*Yy + z*Zy + Ty
//	 z' = x*Xz + y*Yz + z*Zz + Tz
// Where x, y, z is the original vector and X, Y, Z are the three axes of the
// coordinate system. Note that expectations can differ per implementation, eg:
//   “For programming purposes, OpenGL matrices are 16-value arrays with base
//    vectors laid out contiguously in memory. The translation components occupy
//    the 13th, 14th, and 15th elements of the 16-element matrix, where indices
//    are numbered from 1 to 16"
// This means the memory layout expected by OpenGL is:
//    Xx, Xy, Xz, Xw, Yx, Yy, Yz, Yw, Zx, Zy, Zz, Zw, Wx, Wy, Wz, Ww
// with the translation values Tx, Ty, Tz at Wx, Wy, Wz. Note that OpenGL
// GLSL shaders interpret each base vector as a column (Column-Major)
// although it is appears as Row-Major when viewed from Golang. Note that
// DirectX HLSL shaders interpret the same memory layout as Row-Major.
// In either case, consistency is key, especially for transforms where it
// is always apply Scale first, then Rotatate, then Translate.
//
// Conforming to the above memory layout, this matrix implementation uses
// explicitly indexed, Row-Major, matrix members as follows:
//          3x3 M3          4x4 M4
//	     [Xx, Xy, Xz]  [Xx, Xy, Xz, Xw]  X-Axis
//	     [Yx, Yy, Yz]  [Yx, Yy, Yz, Yw]  Y-Axis
//	     [Zx, Zy, Zz]  [Zx, Zy, Zz, Zw]  Z-Axis
//	                   [Wx, Wy, Wz, Ww]  Translation vector, Ww == 1.
// This layout allows the entire structure to be passed as a pointer to the
// underlying (C-language) graphics layer.
//
// See appendix G of OpenGL Red Book for matrix algorithms. Also see:
// http://steve.hollasch.net/cgindex/math/matrix/column-vec.html
// http://stackoverflow.com/questions/17784791/4x4-matrix-pre-multiplication-vs-post-multiplication
// http://www.scratchapixel.com/lessons/3d-basic-lessons/lesson-4-geometry/conventions-again-row-major-vs-column-major-vector/

import (
	"log"
	"math"
)

// M3 is a 3x3 matrix where the matrix elements are individually addressable.
type M3 struct {
	Xx, Xy, Xz float64 // indices 0, 1, 2  [00, 01, 02]  X-Axis
	Yx, Yy, Yz float64 // indices 3, 4, 5  [10, 11, 12]  Y-Axis
	Zx, Zy, Zz float64 // indices 6, 7, 8  [20, 21, 22]  Z-Axis
}

// M4 is a 4x4 matrix where the matrix elements are individually addressable.
type M4 struct {
	Xx, Xy, Xz, Xw float64 // indices 0, 1, 2, 3  [00, 01, 02, 03] X-Axis
	Yx, Yy, Yz, Yw float64 // indices 4, 5, 6, 7  [10, 11, 12, 13] Y-Axis
	Zx, Zy, Zz, Zw float64 // indices 8, 9, a, b  [20, 21, 22, 23] Z-Axis
	Wx, Wy, Wz, Ww float64 // indices c, d, e, f  [30, 31, 32, 33]
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

// Eq (==) returns true if all the elements in matrix m have the same value
// as the corresponding elements in matrix a.
func (m *M3) Eq(a *M3) bool {
	return true &&
		m.Xx == a.Xx && m.Xy == a.Xy && m.Xz == a.Xz &&
		m.Yx == a.Yx && m.Yy == a.Yy && m.Yz == a.Yz &&
		m.Zx == a.Zx && m.Zy == a.Zy && m.Zz == a.Zz
}

// Eq (==) returns true if all the elements in matrix m have the same value
// as the corresponding elements in matrix a.
func (m *M4) Eq(a *M4) bool {
	return true &&
		m.Xx == a.Xx && m.Xy == a.Xy && m.Xz == a.Xz && m.Xw == a.Xw &&
		m.Yx == a.Yx && m.Yy == a.Yy && m.Yz == a.Yz && m.Yw == a.Yw &&
		m.Zx == a.Zx && m.Zy == a.Zy && m.Zz == a.Zz && m.Zw == a.Zw &&
		m.Wx == a.Wx && m.Wy == a.Wy && m.Wz == a.Wz && m.Ww == a.Ww
}

// Aeq (~=) almost equals returns true if all the elements in matrix m have
// essentially the same value as the corresponding elements in matrix a.
// Used where equals is unlikely to return true due to float precision.
func (m *M3) Aeq(a *M3) bool {
	return true &&
		Aeq(m.Xx, a.Xx) && Aeq(m.Xy, a.Xy) && Aeq(m.Xz, a.Xz) &&
		Aeq(m.Yx, a.Yx) && Aeq(m.Yy, a.Yy) && Aeq(m.Yz, a.Yz) &&
		Aeq(m.Zx, a.Zx) && Aeq(m.Zy, a.Zy) && Aeq(m.Zz, a.Zz)
}

// Aeq (~=) almost equals returns true if all the elements in matrix m have
// essentially the same value as the corresponding elements in matrix a.
// Same as M3.Aeq().
func (m *M4) Aeq(a *M4) bool {
	return true &&
		Aeq(m.Xx, a.Xx) && Aeq(m.Xy, a.Xy) && Aeq(m.Xz, a.Xz) && Aeq(m.Xw, a.Xw) &&
		Aeq(m.Yx, a.Yx) && Aeq(m.Yy, a.Yy) && Aeq(m.Yz, a.Yz) && Aeq(m.Yw, a.Yw) &&
		Aeq(m.Zx, a.Zx) && Aeq(m.Zy, a.Zy) && Aeq(m.Zz, a.Zz) && Aeq(m.Zw, a.Zw) &&
		Aeq(m.Wx, a.Wx) && Aeq(m.Wy, a.Wy) && Aeq(m.Wz, a.Wz) && Aeq(m.Ww, a.Ww)
}

// SetS (=) explicitly sets the matrix scaler values using the given scalers.
// The source matrix a is unchanged. The updated matrix m is returned.
// 	  Xx, Xy, Xz is the X Axis.
// 	  Yx, Yy, Yz is the Y Axis.
// 	  Zx, Zy, Zz is the Z Axis.
func (m *M3) SetS(Xx, Xy, Xz, Yx, Yy, Yz, Zx, Zy, Zz float64) *M3 {
	m.Xx, m.Xy, m.Xz = Xx, Xy, Xz
	m.Yx, m.Yy, m.Yz = Yx, Yy, Yz
	m.Zx, m.Zy, m.Zz = Zx, Zy, Zz
	return m
}

// Set (=) assigns all the scaler values from matrix a to the
// corresponding scaler values in matrix m.
// The source matrix a is unchanged. The updated matrix m is returned.
func (m *M3) Set(a *M3) *M3 {
	m.Xx, m.Xy, m.Xz = a.Xx, a.Xy, a.Xz
	m.Yx, m.Yy, m.Yz = a.Yx, a.Yy, a.Yz
	m.Zx, m.Zy, m.Zz = a.Zx, a.Zy, a.Zz
	return m
}

// SetM4 (=) updates calling matrix m to be the 3x3 matrix from the top left
// corner of the given 4x4 matrix m4. The source matrix a is unchanged.
// The updated matrix m is returned.
//    [ Xx Xy Xz Xw ]    [ Xx Xy Xz ]
//    [ Yx Yy Yz Yw ] => [ Yx Yy Yz ]
//    [ Zx Zy Zz Zw ]    [ Zx Zy Zz ]
//    [ Wx Wy Wz Ww ]
func (m *M3) SetM4(a *M4) *M3 {
	m.Xx, m.Xy, m.Xz = a.Xx, a.Xy, a.Xz
	m.Yx, m.Yy, m.Yz = a.Yx, a.Yy, a.Yz
	m.Zx, m.Zy, m.Zz = a.Zx, a.Zy, a.Zz
	return m
}

// Set (=) assigns all the elements values from matrix a to the
// corresponding element values in matrix m. The source matrix a is unchanged.
// The updated matrix m is returned.
func (m *M4) Set(a *M4) *M4 {
	m.Xx, m.Xy, m.Xz, m.Xw = a.Xx, a.Xy, a.Xz, a.Xw
	m.Yx, m.Yy, m.Yz, m.Yw = a.Yx, a.Yy, a.Yz, a.Yw
	m.Zx, m.Zy, m.Zz, m.Zw = a.Zx, a.Zy, a.Zz, a.Zw
	m.Wx, m.Wy, m.Wz, m.Ww = a.Wx, a.Wy, a.Wz, a.Ww
	return m
}

// Abs updates m to be the the absolute (non-negative) element values of the
// corresponding element values in matrix a. The source matrix a is unchanged.
// The updated matrix m is returned.
func (m *M3) Abs(a *M3) *M3 {
	m.Xx, m.Xy, m.Xz = math.Abs(a.Xx), math.Abs(a.Xy), math.Abs(a.Xz)
	m.Yx, m.Yy, m.Yz = math.Abs(a.Yx), math.Abs(a.Yy), math.Abs(a.Yz)
	m.Zx, m.Zy, m.Zz = math.Abs(a.Zx), math.Abs(a.Zy), math.Abs(a.Zz)
	return m
}

// Transpose updates m to be the reflection of matrix a over its diagonal.
// This essentially changes row-major order to column-major order
// or vice-versa.
//    [ Xx Xy Xz ]    [ Xx Yx Zx ]
//    [ Yx Yy Yz ] => [ Xy Yy Zy ]
//    [ Zx Zy Zz ]    [ Xz Yz Zz ]
// The input matrix a is not changed. Matrix m may be used as the input parameter.
// The updated matrix m is returned.
func (m *M3) Transpose(a *M3) *M3 {
	t_Xy, t_Xz, t_Yz := a.Xy, a.Xz, a.Yz
	m.Xx, m.Xy, m.Xz = a.Xx, a.Yx, a.Zx
	m.Yx, m.Yy, m.Yz = t_Xy, a.Yy, a.Zy
	m.Zx, m.Zy, m.Zz = t_Xz, t_Yz, a.Zz
	return m
}

// Transpose updates m to be the reflection of matrix a over its diagonal.
//    [ Xx Xy Xz Xw ]    [ Xx Yx Zx Wx ]
//    [ Yx Yy Yz Yw ] => [ Xy Yy Zy Wy ]
//    [ Zx Zy Zz Zw ]    [ Xz Yz Zz Wz ]
//    [ Wx Wy Wz Ww ]    [ Xw Yw Zw Ww ]
// Same behaviour as M3.Transpose()
func (m *M4) Transpose(a *M4) *M4 {
	t_Xy, t_Xz, t_Yz := a.Xy, a.Xz, a.Yz
	t_Xw, t_Yw, t_Zw := a.Xw, a.Yw, a.Zw
	m.Xx, m.Xy, m.Xz, m.Xw = a.Xx, a.Yx, a.Zx, a.Wx
	m.Yx, m.Yy, m.Yz, m.Yw = t_Xy, a.Yy, a.Zy, a.Wy
	m.Zx, m.Zy, m.Zz, m.Zw = t_Xz, t_Yz, a.Zz, a.Wz
	m.Wx, m.Wy, m.Wz, m.Ww = t_Xw, t_Yw, t_Zw, a.Ww
	return m
}

// Add (+) adds matrices a and b storing the results in m.
// Each element of matrix b is added to the corresponding matrix a element.
// It is safe to use the calling matrix m as one or both of the parameters.
// For example the plus.equals operation (+=) is
//     m.Add(m, b)
// The updated matrix m is returned.
func (m *M3) Add(a, b *M3) *M3 {
	m.Xx, m.Xy, m.Xz = a.Xx+b.Xx, a.Xy+b.Xy, a.Xz+b.Xz
	m.Yx, m.Yy, m.Yz = a.Yx+b.Yx, a.Yy+b.Yy, a.Yz+b.Yz
	m.Zx, m.Zy, m.Zz = a.Zx+b.Zx, a.Zy+b.Zy, a.Zz+b.Zz
	return m
}

// Add (+) adds matrices a and b storing the results in m.
// Same behaviour as M3.Add()
func (m *M4) Add(a, b *M4) *M4 {
	m.Xx, m.Xy, m.Xz, m.Xw = a.Xx+b.Xx, a.Xy+b.Xy, a.Xz+b.Xz, a.Xw+b.Xw
	m.Yx, m.Yy, m.Yz, m.Yw = a.Yx+b.Yx, a.Yy+b.Yy, a.Yz+b.Yz, a.Yw+b.Yw
	m.Zx, m.Zy, m.Zz, m.Zw = a.Zx+b.Zx, a.Zy+b.Zy, a.Zz+b.Zz, a.Zw+b.Zw
	m.Wx, m.Wy, m.Wz, m.Ww = a.Wx+b.Wx, a.Wy+b.Wy, a.Wz+b.Wz, a.Ww+b.Ww
	return m
}

// Sub (-) subtracts matrices b from a storing the results in m.
// Each element of matrix b is subtracted from the corresponding matrix a element.
// It is safe to use the calling matrix m as one or both of the parameters.
// For example the minus.equals operation (-=) is
//     m.Sub(m, b)
// The updated matrix m is returned.
func (m *M3) Sub(a, b *M3) *M3 {
	m.Xx, m.Xy, m.Xz = a.Xx-b.Xx, a.Xy-b.Xy, a.Xz-b.Xz
	m.Yx, m.Yy, m.Yz = a.Yx-b.Yx, a.Yy-b.Yy, a.Yz-b.Yz
	m.Zx, m.Zy, m.Zz = a.Zx-b.Zx, a.Zy-b.Zy, a.Zz-b.Zz
	return m
}

// Mult (*) multiplies matrices l and r storing the results in m.
//    [ lXx lXy lXz ] [ rXx rXy rXz ]    [ mXx mXy mXz ]
//    [ lYx lYy lYz ]x[ rYx rYy rYz ] => [ mYx mYy mYz ]
//    [ lZx lZy lZz ] [ rZx rZy rZz ]    [ mZx mZy mZz ]
// It is safe to use the calling matrix m as one or both of the parameters.
// For example (*=) is
//     m.Mult(m, r)
// The updated matrix m is returned.
func (m *M3) Mult(l, r *M3) *M3 {
	xx := l.Xx*r.Xx + l.Xy*r.Yx + l.Xz*r.Zx
	xy := l.Xx*r.Xy + l.Xy*r.Yy + l.Xz*r.Zy
	xz := l.Xx*r.Xz + l.Xy*r.Yz + l.Xz*r.Zz
	yx := l.Yx*r.Xx + l.Yy*r.Yx + l.Yz*r.Zx
	yy := l.Yx*r.Xy + l.Yy*r.Yy + l.Yz*r.Zy
	yz := l.Yx*r.Xz + l.Yy*r.Yz + l.Yz*r.Zz
	zx := l.Zx*r.Xx + l.Zy*r.Yx + l.Zz*r.Zx
	zy := l.Zx*r.Xy + l.Zy*r.Yy + l.Zz*r.Zy
	zz := l.Zx*r.Xz + l.Zy*r.Yz + l.Zz*r.Zz
	m.Xx, m.Xy, m.Xz = xx, xy, xz
	m.Yx, m.Yy, m.Yz = yx, yy, yz
	m.Zx, m.Zy, m.Zz = zx, zy, zz
	return m
}

// Mult updates matrix m to be the multiplication of input matrices l, r.
//    [ lXx lXy lXz lXw ] [ rXx rXy rXz rXw ]    [ mXx mXy mXz mXw ]
//    [ lYx lYy lYz lYw ]x[ rYx rYy rYz rYw ] => [ mYx mYy mYz mYw ]
//    [ lZx lZy lZz lZw ] [ rZx rZy rZz rZw ]    [ mZx mZy mZz mZw ]
//    [ lWx lWy lWz lWw ] [ rWx rWy rWz rWw ]    [ mWx mWy mWz mWw ]
// Same behaviour as M3.Mult()
func (m *M4) Mult(l, r *M4) *M4 {
	xx := l.Xx*r.Xx + l.Xy*r.Yx + l.Xz*r.Zx + l.Xw*r.Wx
	xy := l.Xx*r.Xy + l.Xy*r.Yy + l.Xz*r.Zy + l.Xw*r.Wy
	xz := l.Xx*r.Xz + l.Xy*r.Yz + l.Xz*r.Zz + l.Xw*r.Wz
	xw := l.Xx*r.Xw + l.Xy*r.Yw + l.Xz*r.Zw + l.Xw*r.Ww
	yx := l.Yx*r.Xx + l.Yy*r.Yx + l.Yz*r.Zx + l.Yw*r.Wx
	yy := l.Yx*r.Xy + l.Yy*r.Yy + l.Yz*r.Zy + l.Yw*r.Wy
	yz := l.Yx*r.Xz + l.Yy*r.Yz + l.Yz*r.Zz + l.Yw*r.Wz
	yw := l.Yx*r.Xw + l.Yy*r.Yw + l.Yz*r.Zw + l.Yw*r.Ww
	zx := l.Zx*r.Xx + l.Zy*r.Yx + l.Zz*r.Zx + l.Zw*r.Wx
	zy := l.Zx*r.Xy + l.Zy*r.Yy + l.Zz*r.Zy + l.Zw*r.Wy
	zz := l.Zx*r.Xz + l.Zy*r.Yz + l.Zz*r.Zz + l.Zw*r.Wz
	zw := l.Zx*r.Xw + l.Zy*r.Yw + l.Zz*r.Zw + l.Zw*r.Ww
	wx := l.Wx*r.Xx + l.Wy*r.Yx + l.Wz*r.Zx + l.Ww*r.Wx
	wy := l.Wx*r.Xy + l.Wy*r.Yy + l.Wz*r.Zy + l.Ww*r.Wy
	wz := l.Wx*r.Xz + l.Wy*r.Yz + l.Wz*r.Zz + l.Ww*r.Wz
	ww := l.Wx*r.Xw + l.Wy*r.Yw + l.Wz*r.Zw + l.Ww*r.Ww
	m.Xx, m.Xy, m.Xz, m.Xw = xx, xy, xz, xw
	m.Yx, m.Yy, m.Yz, m.Yw = yx, yy, yz, yw
	m.Zx, m.Zy, m.Zz, m.Zw = zx, zy, zz, zw
	m.Wx, m.Wy, m.Wz, m.Ww = wx, wy, wz, ww
	return m
}

// MultLtR multiplies the transpose of matrix l on left of matrix r
// and stores the result in m. This can be used for saving a method call
// when calculating inverse transforms.
//    [ lXx lYx lZx ] [ rXx rXy rXz ]    [ mXx mXy mXz ]
//    [ lXy lYy lZy ]x[ rYx rYy rYz ] => [ mYx mYy mYz ]
//    [ lXz lYz lZz ] [ rZx rZy rZz ]    [ mZx mZy mZz ]
// It is safe to use the calling matrix m as one or both of the parameters.
// The updated matrix m is returned.
func (m *M3) MultLtR(lt, r *M3) *M3 {
	xx := lt.Xx*r.Xx + lt.Yx*r.Yx + lt.Zx*r.Zx
	xy := lt.Xx*r.Xy + lt.Yx*r.Yy + lt.Zx*r.Zy
	xz := lt.Xx*r.Xz + lt.Yx*r.Yz + lt.Zx*r.Zz
	yx := lt.Xy*r.Xx + lt.Yy*r.Yx + lt.Zy*r.Zx
	yy := lt.Xy*r.Xy + lt.Yy*r.Yy + lt.Zy*r.Zy
	yz := lt.Xy*r.Xz + lt.Yy*r.Yz + lt.Zy*r.Zz
	zx := lt.Xz*r.Xx + lt.Yz*r.Yx + lt.Zz*r.Zx
	zy := lt.Xz*r.Xy + lt.Yz*r.Yy + lt.Zz*r.Zy
	zz := lt.Xz*r.Xz + lt.Yz*r.Yz + lt.Zz*r.Zz
	m.Xx, m.Xy, m.Xz = xx, xy, xz
	m.Yx, m.Yy, m.Yz = yx, yy, yz
	m.Zx, m.Zy, m.Zz = zx, zy, zz
	return m
}

// TranslateTM updates m to be the multiplication of a translation matrix
// T created from x, y, z, and itself. The updated matrix m is returned.
//    [ 1 0 0 0 ]   [ mXx mXy mXz mXw ]     [ mXx  mXy  mXz  mXw  ]
//    [ 0 1 0 0 ] x [ mYx mYy mYz mYw ]  => [ mYx  mYy  mYz  mYw  ]
//    [ 0 0 1 0 ]   [ mZx mZy mZz mZw ]     [ mZx  mZy  mZz  mZw  ]
//    [ x y z 1 ]   [ mWx mWy mWz mWw ]     [ mWx' mWy' mWz' mWw' ]
// Be sure to pick the correct translate (TM or MT) when doing transforms.
func (m *M4) TranslateTM(x, y, z float64) *M4 {
	wx := x*m.Xx + y*m.Yx + z*m.Zx + m.Wx
	wy := x*m.Xy + y*m.Yy + z*m.Zy + m.Wy
	wz := x*m.Xz + y*m.Yz + z*m.Zz + m.Wz
	ww := x*m.Xw + y*m.Yw + z*m.Zw + m.Ww
	m.Wx, m.Wy, m.Wz, m.Ww = wx, wy, wz, ww
	return m
}

// TranslateMT updates m to be the multiplication of itself
// and a translation matrix created from x, y, z.
// The updated matrix m is returned.
//    [ mXx mXy mXz mXw ]   [ 1 0 0 0 ]    [ mXx' mXy' mXz' mXw ]
//    [ mYx mYy mYz mYw ] x [ 0 1 0 0 ] => [ mYx' mYy' mYz' mYw ]
//    [ mZx mZy mZz mZw ]   [ 0 0 1 0 ]    [ mZx' mZy' mZz' mZw ]
//    [ mWx mWy mWz mWw ]   [ x y z 1 ]    [ mWx' mWy' mWz' mWw ]
// Be sure to pick the correct translate (TM or MT) when doing transforms.
func (m *M4) TranslateMT(x, y, z float64) *M4 {
	m.Xx, m.Xy, m.Xz = m.Xx+m.Xw*x, m.Xy+m.Xw*y, m.Xz+m.Xw*z
	m.Yx, m.Yy, m.Yz = m.Yx+m.Yw*x, m.Yy+m.Yw*y, m.Yz+m.Yw*z
	m.Zx, m.Zy, m.Zz = m.Zx+m.Zw*x, m.Zy+m.Zw*y, m.Zz+m.Zw*z
	m.Wx, m.Wy, m.Wz = m.Wx+m.Ww*x, m.Wy+m.Ww*y, m.Wz+m.Ww*z
	return m
}

// Scale (*) each element of matrix m by the given scalar.
// The updated matrix m is returned.
func (m *M3) Scale(s float64) *M3 {
	m.Xx, m.Xy, m.Xz = m.Xx*s, m.Xy*s, m.Xz*s
	m.Yx, m.Yy, m.Yz = m.Yx*s, m.Yy*s, m.Yz*s
	m.Zx, m.Zy, m.Zz = m.Zx*s, m.Zy*s, m.Zz*s
	return m
}

// Scale (*) each element of matrix m by the given scalar.
// The updated matrix m is returned.
func (m *M4) Scale(s float64) *M4 {
	m.Xx, m.Xy, m.Xz, m.Xw = m.Xx*s, m.Xy*s, m.Xz*s, m.Xw*s
	m.Yx, m.Yy, m.Yz, m.Yw = m.Yx*s, m.Yy*s, m.Yz*s, m.Yw*s
	m.Zx, m.Zy, m.Zz, m.Zw = m.Zx*s, m.Zy*s, m.Zz*s, m.Zw*s
	m.Wx, m.Wy, m.Wz, m.Ww = m.Wx*s, m.Wy*s, m.Wz*s, m.Ww*s
	return m
}

// ScaleS (*) scales each column of matrix m using the corresponding scaler
// elements x, y, z. The updated matrix m is returned.
func (m *M3) ScaleS(x, y, z float64) *M3 {
	m.Xx, m.Xy, m.Xz = m.Xx*x, m.Xy*y, m.Xz*z
	m.Yx, m.Yy, m.Yz = m.Yx*x, m.Yy*y, m.Yz*z
	m.Zx, m.Zy, m.Zz = m.Zx*x, m.Zy*y, m.Zz*z
	return m
}

// ScaleV (*) scales each column of matrix m using the given vector v
// for elements for x, y, z. The updated matrix m is returned.
func (m *M3) ScaleV(v *V3) *M3 {
	m.Xx, m.Xy, m.Xz = m.Xx*v.X, m.Xy*v.Y, m.Xz*v.Z
	m.Yx, m.Yy, m.Yz = m.Yx*v.X, m.Yy*v.Y, m.Yz*v.Z
	m.Zx, m.Zy, m.Zz = m.Zx*v.X, m.Zy*v.Y, m.Zz*v.Z
	return m
}

// ScaleSM updates m to be the multiplication of a scale matrix
// created from x, y, z and itself. The updated matrix m is
// returned so that it may be immediately used in another operation.
//    [ x 0 0 ]   [ mXx mXy mXz ]    [ mXx' mXy' mXz' ]
//    [ 0 y 0 ] x [ mYx mYy mYz ] => [ mYx' mYy' mYz' ]
//    [ 0 0 z ]   [ mZx mZy mZz ]    [ mZx' mZy' mZz' ]
// Be sure to pick the correct scale (SM or MS) when doing transforms.
func (m *M3) ScaleSM(x, y, z float64) *M3 {
	m.Xx, m.Xy, m.Xz = m.Xx*x, m.Xy*x, m.Xz*x
	m.Yx, m.Yy, m.Yz = m.Yx*y, m.Yy*y, m.Yz*y
	m.Zx, m.Zy, m.Zz = m.Zx*z, m.Zy*z, m.Zz*z
	return m
}

// ScaleSM updates m to be the multiplication of a scale matrix
// created from x, y, z and itself. Same behaviours as M3.ScaleSM.
//    [ x 0 0 0 ]   [ mXx mXy mXz mXw ]    [ mXx' mXy' mXz' mXw' ]
//    [ 0 y 0 0 ] x [ mYx mYy mYz mYw ] => [ mYx' mYy' mYz' mYw' ]
//    [ 0 0 z 0 ]   [ mZx mZy mZz mZw ]    [ mZx' mZy' mZz' mZw' ]
//    [ 0 0 0 1 ]   [ mWx mWy mWz mWw ]    [ mWx  mWy  mWz  mWw  ]
func (m *M4) ScaleSM(x, y, z float64) *M4 {
	m.Xx, m.Xy, m.Xz, m.Xw = m.Xx*x, m.Xy*x, m.Xz*x, m.Xw*x
	m.Yx, m.Yy, m.Yz, m.Yw = m.Yx*y, m.Yy*y, m.Yz*y, m.Yw*y
	m.Zx, m.Zy, m.Zz, m.Zw = m.Zx*z, m.Zy*z, m.Zz*z, m.Zw*z
	return m
}

// ScaleMS updates m to be the multiplication of m and a scale matrix created
// from x, y, z. The updated matrix m is returned so that it may be immediately
// used in another operation.
//    [ mXx mXy mXz mXw ]   [ x 0 0 0 ]    [ mXx' mXy' mXz' mXw ]
//    [ mYx mYy mYz mYw ] x [ 0 y 0 0 ] => [ mYx' mYy' mYz' mYw ]
//    [ mZx mZy mZz mZw ]   [ 0 0 z 0 ]    [ mZx' mZy' mZz' mZw ]
//    [ mWx mWy mWz mWw ]   [ 0 0 0 1 ]    [ mWx' mWy' mWz' mWw ]
// Be sure to pick the correct scale (SM or MS) when doing transforms.
func (m *M4) ScaleMS(x, y, z float64) *M4 {
	m.Xx, m.Xy, m.Xz = m.Xx*x, m.Xy*y, m.Xz*z
	m.Yx, m.Yy, m.Yz = m.Yx*x, m.Yy*y, m.Yz*z
	m.Zx, m.Zy, m.Zz = m.Zx*x, m.Zy*y, m.Zz*z
	m.Wx, m.Wy, m.Wz = m.Wx*x, m.Wy*y, m.Wz*z
	return m
}

// SetQ converts a quaternion rotation representation to a matrix
// rotation representation. SetQ updates matrix m to be the rotation
// matrix representing the rotation described by unit-quaternion q.
//                       [ mXx mXy mXz ]
//    [ qx qy qz qw ] => [ mYx mYy mYz ]
//                       [ mZx mZy mZz ]
// The parameter q is unchanged. The updated matrix m is returned.
func (m *M3) SetQ(q *Q) *M3 {
	xx, yy, zz := q.X*q.X, q.Y*q.Y, q.Z*q.Z
	xy, xz, yz := q.X*q.Y, q.X*q.Z, q.Y*q.Z
	wx, wy, wz := q.W*q.X, q.W*q.Y, q.W*q.Z
	m.Xx, m.Xy, m.Xz = 1-2*(yy+zz), 2*(xy-wz), 2*(xz+wy)
	m.Yx, m.Yy, m.Yz = 2*(xy+wz), 1-2*(xx+zz), 2*(yz-wx)
	m.Zx, m.Zy, m.Zz = 2*(xz-wy), 2*(yz+wx), 1-2*(xx+yy)
	return m
}

// SetQ converts a quaternion rotation representation to a matrix
// rotation representation. SetQ updates matrix m to be the rotation
// matrix representing the rotation described by unit-quaternion q.
//                       [ mXx mXy mXz 0 ]
//    [ qx qy qz qw ] => [ mYx mYy mYz 0 ]
//                       [ mZx mZy mZz 0 ]
//                       [  0   0   0  1 ]
// The parameter q is unchanged. The updated matrix m is returned.
func (m *M4) SetQ(q *Q) *M4 {
	xx, yy, zz := q.X*q.X, q.Y*q.Y, q.Z*q.Z
	xy, xz, yz := q.X*q.Y, q.X*q.Z, q.Y*q.Z
	wx, wy, wz := q.W*q.X, q.W*q.Y, q.W*q.Z
	m.Xx, m.Xy, m.Xz, m.Xw = 1-2*(yy+zz), 2*(xy-wz), 2*(xz+wy), 0
	m.Yx, m.Yy, m.Yz, m.Yw = 2*(xy+wz), 1-2*(xx+zz), 2*(yz-wx), 0
	m.Zx, m.Zy, m.Zz, m.Zw = 2*(xz-wy), 2*(yz+wx), 1-2*(xx+yy), 0
	m.Wx, m.Wy, m.Wz, m.Ww = 0, 0, 0, 1
	return m
}

// SetSkewSym sets the matrix m to be a skew-symetric matrix based
// on the elements of vector v. Wikipedia states:
//    "A skew-symmetric matrix is a square matrix
//     whose transpose is also its negative."
func (m *M3) SetSkewSym(v *V3) *M3 {
	m.Xx, m.Xy, m.Xz = 0, -v.Z, v.Y
	m.Yx, m.Yy, m.Yz = v.Z, 0, -v.X
	m.Zx, m.Zy, m.Zz = -v.Y, v.X, 0
	return m
}

// Det returns the determinant of matrix m. Determinants are helpful
// when calculating the inverse of transform matrices. Wikipedia states:
//    "The determinant provides important information about [..] a matrix that
//     corresponds to a linear transformation of a vector space [..] the transformation
//     has an inverse operation exactly when the determinant is nonzero."
func (m *M3) Det() float64 {
	return m.Xx*(m.Yy*m.Zz-m.Yz*m.Zy) + m.Xy*(m.Yz*m.Zx-m.Yx*m.Zz) + m.Xz*(m.Yx*m.Zy-m.Yy*m.Zx)
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
		return m.Yy*m.Zz - m.Yz*m.Zy
	case 01:
		return m.Yz*m.Zx - m.Yx*m.Zz // flip to negate.
	case 02:
		return m.Yx*m.Zy - m.Yy*m.Zx
	case 10:
		return m.Xz*m.Zy - m.Xy*m.Zz // flip to negate.
	case 11:
		return m.Xx*m.Zz - m.Xz*m.Zx
	case 12:
		return m.Xy*m.Zx - m.Xx*m.Zy // flip to negate.
	case 20:
		return m.Xy*m.Yz - m.Xz*m.Yy
	case 21:
		return m.Xz*m.Yx - m.Xx*m.Yz // flip to negate.
	case 22:
		return m.Xx*m.Yy - m.Xy*m.Yx
	}
	log.Printf("matrix M3.Cof developer error %d", minor)
	return 0
}

// Adj updates m to be the adjoint matrix of matrix a. The adjoint matrix is
// created by the transpose of the cofactor matrix of the original matrix.
//     [ a.cof(0,0) a.cof(1,0) a.cof(2,0) ]    [ mXx mXy mXz ]
//     [ a.cof(0,1) a.cof(1,1) a.cof(2,1) ] => [ mYx mYy mYz ]
//     [ a.cof(0,2) a.cof(1,2) a.cof(2,2) ]    [ mZx mZy mZz ]
// The updated matrix m is returned.
func (m *M3) Adj(a *M3) *M3 {
	xx, xy, xz := a.Cof(0, 0), a.Cof(1, 0), a.Cof(2, 0)
	yx, yy, yz := a.Cof(0, 1), a.Cof(1, 1), a.Cof(2, 1)
	zx, zy, zz := a.Cof(0, 2), a.Cof(1, 2), a.Cof(2, 2)
	m.Xx, m.Xy, m.Xz = xx, xy, xz
	m.Yx, m.Yy, m.Yz = yx, yy, yz
	m.Zx, m.Zy, m.Zz = zx, zy, zz
	return m
}

// Inv updates m to be the inverse of matrix a. The updated matrix m is returned.
// Matrix m is not updated if the matrix has no inverse.
func (m *M3) Inv(a *M3) *M3 {
	det := a.Det()
	if det != 0 {
		s := 1 / det
		xx, xy, xz := a.Cof(0, 0)*s, a.Cof(1, 0)*s, a.Cof(2, 0)*s
		yx, yy, yz := a.Cof(0, 1)*s, a.Cof(1, 1)*s, a.Cof(2, 1)*s
		zx, zy, zz := a.Cof(0, 2)*s, a.Cof(1, 2)*s, a.Cof(2, 2)*s
		m.Xx, m.Xy, m.Xz = xx, xy, xz
		m.Yx, m.Yy, m.Yz = yx, yy, yz
		m.Zx, m.Zy, m.Zz = zx, zy, zz
	}
	return m
}

// SetAa, set axis-angle, updates m to be a rotation matrix from the
// given axis (ax, ay, az) and angle (in radians). See:
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
	m.Xx = rcos + ax*ax*(1-rcos)
	m.Xy = -az*rsin + ay*ax*(1-rcos)
	m.Xz = ay*rsin + az*ax*(1-rcos)
	m.Yx = az*rsin + ax*ay*(1-rcos)
	m.Yy = rcos + ay*ay*(1-rcos)
	m.Yz = -ax*rsin + az*ay*(1-rcos)
	m.Zx = -ay*rsin + ax*az*(1-rcos)
	m.Zy = ax*rsin + ay*az*(1-rcos)
	m.Zz = rcos + az*az*(1-rcos)
	return m
}

// Ortho sets matrix m with projection values needed to
// transform a 3 dimensional model to a 2 dimensional plane.
// Orthographic projection ignores depth. The input arguments are:
//     left, right:  Vertical clipping planes.
//     bottom, top:  Horizontal clipping planes.
//     near, far  :  Depth clipping planes. The depth values are
//                   negative if the plane is to be behind the viewer
// An orthographic matrix fills the following matrix locations:
//    [ a 0 0 0 ]    [ Xx Xy Xz Xw ]
//    [ 0 b 0 0 ] => [ Yx Yy Yz Yw ]
//    [ 0 0 c 0 ]    [ Zx Zy Zz Zw ]
//    [ d e f 1 ]    [ Wx Wy Wz Ww ]
func (m *M4) Ortho(left, right, bottom, top, near, far float64) *M4 {
	m.Xx = 2 / (right - left)
	m.Xy = 0
	m.Xz = 0
	m.Xw = 0
	m.Yx = 0
	m.Yy = 2 / (top - bottom)
	m.Yz = 0
	m.Yw = 0
	m.Zx = 0
	m.Zy = 0
	m.Zz = -2 / (far - near)
	m.Zw = 0
	m.Wx = -(right + left) / (right - left)
	m.Wy = -(top + bottom) / (top - bottom)
	m.Wz = -(far + near) / (far - near)
	m.Ww = 1
	return m
}

// Persp sets matrix m with projection values needed to
// transform a 3 dimensional model to a 2 dimensional plane.
// Objects that are further away from the viewer will appear smaller.
// The input arguments are:
//    fov        An amount in degrees indicating how much of the
//               scene is visible.
//    aspect     The ratio of height to width of the model.
//    near, far  The depth clipping planes. The depth values are
//               negative if the plane is to be behind the viewer
// A perspective projection matrix fills the following matrix locations:
//    [ a 0 0 0 ]    [ Xx Xy Xz Xw ]
//    [ 0 b 0 0 ] => [ Yx Yy Yz Yw ]
//    [ 0 0 c d ]    [ Zx Zy Zz Zw ]
//    [ 0 0 e 0 ]    [ Wx Wy Wz Ww ]
func (m *M4) Persp(fov, aspect, near, far float64) *M4 {
	f := 1 / float64(math.Tan(Rad(fov)*0.5))
	m.Xx = f / aspect
	m.Yx = 0
	m.Zx = 0
	m.Wx = 0
	m.Xy = 0
	m.Yy = f
	m.Zy = 0
	m.Wy = 0
	m.Xz = 0
	m.Yz = 0
	m.Zz = (far + near) / (near - far)
	m.Wz = 2 * far * near / (near - far)
	m.Xw = 0
	m.Yw = 0
	m.Zw = -1
	m.Ww = 0
	return m
}

// PerspInv sets matrix m to be a new inverse matrix of the given
// perspective matrix values (see NewPersp()).
//   [ a' 0  0  0 ] where a' = 1/a     d' = 1/e    [ Xx Xy Xz Xw ]
//   [ 0  b' 0  0 ]       b' = 1/b     e' = 1/d => [ Yx Yy Yz Yw ]
//   [ 0  0  0  d']       c' = -(c/de)             [ Zx Zy Zz Zw ]
//   [ 0  0  e' c']                                [ Wx Wy Wz Ww ]
// This is used when going from screen x,y coordinates to 3D coordinates.
// as in the case when creating a picking ray from a mouse location.
func (m *M4) PerspInv(fov, aspect, near, far float64) *M4 {
	f := float64(math.Tan(Rad(fov) * 0.5))
	c := 2 * far * near / (near - far)
	m.Xx = f * aspect
	m.Yx = 0
	m.Zx = 0
	m.Wx = 0
	m.Xy = 0
	m.Yy = f
	m.Zy = 0
	m.Wy = 0
	m.Xz = 0
	m.Yz = 0
	m.Zz = 0
	m.Wz = -1
	m.Xw = 0
	m.Yw = 0
	m.Zw = 1 / c
	m.Ww = -((far + near) / (near - far) / (-1 * c))
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
//    [ 1 0 0 ]    [ Xx Xy Xz ]
//    [ 0 1 0 ] => [ Yx Yy Yz ]
//    [ 0 0 1 ]    [ Zx Zy Zz ]
func NewM3I() *M3 { return &M3{Xx: 1, Yy: 1, Zz: 1} }

// NewM4I creates a new 4x4 identity matrix.
//    [ 1 0 0 0 ]    [ Xx Xy Xz Xw ]
//    [ 0 1 0 0 ] => [ Yx Yy Yz Yw ]
//    [ 0 0 1 0 ]    [ Zx Zy Zz Zw ]
//    [ 0 0 0 1 ]    [ Wx Wy Wz Ww ]
func NewM4I() *M4 { return &M4{Xx: 1, Yy: 1, Zz: 1, Ww: 1} }
