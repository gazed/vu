// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package lin

// Matrix functions deal with 3x3 and 4x4 matricies expected to be used
// in 3D model calculations. See appendix G of OpenGL Red Book for matrix
// algorithms. Unused Matrix methods like rotation (use quaternions instead)
// are not included.
//
// Where possible the calling matrix is updated instead of creating a new
// one. The pointer to the calling/updated matrix is is passed back so the
// matrix may be immediately used for a second transform.
//
// Golang stores structures in row-major order. That is structure elements
// are stored sequentially in memory as row1, row2, row3, etc...  Additionally
// http://www.opengl.org/archives/resources/faq/technical/transformations.htm#tran0005
// states:
//   “For programming purposes, OpenGL matrices are 16-value arrays with base
//    vectors laid out contiguously in memory. The translation components occupy
//    the 13th, 14th, and 15th elements of the 16-element matrix, where indices
//    are numbered from 1 to 16 as described in section 2.11.2 of the OpenGL
//    2.1 Specification."
// Following this layout the entire structure can be passed as a pointer to the
// underlying (C-language) graphics layer.

import "math"

// M3 is a 3x3 matrix where the matrix elements are individually addressable.
type M3 struct {
	X0, Y0, Z0 float32 // row 1 : indicies 0, 1, 2
	X1, Y1, Z1 float32 // row 2 : indicies 3, 4, 5
	X2, Y2, Z2 float32 // row 3 : indicies 6, 7, 8
}

// M4 is a 4x4 matrix where the matrix elements are individually addressable.
type M4 struct {
	X0, Y0, Z0, W0 float32 // row 1 : indicies 0, 1, 2, 3
	X1, Y1, Z1, W1 float32 // row 2 : indicies 4, 5, 6, 7
	X2, Y2, Z2, W2 float32 // row 3 : indicies 8, 9, a, b
	X3, Y3, Z3, W3 float32 // row 4 : indicies c, d, e, f
}

// Pointer is used to access the matrix data as an array of float32.
// The pointer is expected to be used in native graphic layer calls since
// the native layer expects memory to contain the data as a sequential
// number of bytes.
func (m *M3) Pointer() *float32 { return &(m.X0) }

// Pointer is used to access the matrix data as an array of float32.
// Same behaviour as M3.Pointer()
func (m *M4) Pointer() *float32 { return &(m.X0) }

// M3Identity creates a new 3x3 identity matrix.
// The new matrix is returned so that it may be immediately
// used in another operation.
//
//   [ x0 y0 z0 ]    [ 1 0 0 ]
//   [ x1 y1 z1 ] => [ 0 1 0 ]
//   [ x2 y2 z2 ]    [ 0 0 1 ]
func M3Identity() *M3 {
	return &M3{X0: 1, Y1: 1, Z2: 1}
}

// M4Identity creates a new 4x4 identity matrix.
// Same behaviour as M3Identity()
//
//   [ x0 y0 z0 w0 ]    [ 1 0 0 0 ]
//   [ x1 y1 z1 w1 ] => [ 0 1 0 0 ]
//   [ x2 y2 z2 w2 ]    [ 0 0 1 0 ]
//   [ x3 y3 z3 w3 ]    [ 0 0 0 1 ]
func M4Identity() *M4 {
	return &M4{X0: 1, Y1: 1, Z2: 1, W3: 1}
}

// Clone creates a new copy of a matrix m3.  This is used when performing
// matrix operations where the original matrix values need to be preserved.
func (m *M3) Clone() *M3 {
	return &M3{
		m.X0, m.Y0, m.Z0,
		m.X1, m.Y1, m.Z1,
		m.X2, m.Y2, m.Z2,
	}
}

// Clone has the same behaviour as M3.Clone()
func (m *M4) Clone() *M4 {
	return &M4{
		m.X0, m.Y0, m.Z0, m.W0,
		m.X1, m.Y1, m.Z1, m.W1,
		m.X2, m.Y2, m.Z2, m.W2,
		m.X3, m.Y3, m.Z3, m.W3,
	}
}

// M3 returns a new 3x3 matrix created from the top left corner of the
// source 4x4 matrix m. The source matrix m4 is unchanged.
//
//   [ x0 y0 z0 w0 ]    [ x0 y0 z0 ]
//   [ x1 y1 z1 w1 ] => [ x1 y1 z1 ]
//   [ x2 y2 z2 w2 ]    [ x2 y2 z2 ]
//   [ x3 y3 z3 w3 ]
func (m *M4) M3() *M3 {
	return &M3{
		m.X0, m.Y0, m.Z0,
		m.X1, m.Y1, m.Z1,
		m.X2, m.Y2, m.Z2,
	}
}

// Transpose reflects the matrix m over its diagonal.
// Essentially changes column major order to row major order
// or vice-versa. The updated matrix m is returned so that it
// may be immediately used in another operation.
//
//   [ x0 y0 z0 ]    [ x0 x1 x2 ]
//   [ x1 y1 z1 ] => [ y0 y1 y2 ]
//   [ x2 y2 z2 ]    [ z0 z1 z2 ]
func (m *M3) Transpose() *M3 {
	t_y0 := m.Y0
	t_z0 := m.Z0
	t_z1 := m.Z1
	m.X0, m.Y0, m.Z0 = m.X0, m.X1, m.X2
	m.X1, m.Y1, m.Z1 = t_y0, m.Y1, m.Y2
	m.X2, m.Y2, m.Z2 = t_z0, t_z1, m.Z2
	return m
}

// Transpose reflects the matrix m over its diagonal.
// Same behaviour as M3.Transpose()
//
//   [ x0 y0 z0 w0 ]    [ x0 x1 x2 x3 ]
//   [ x1 y1 z1 w1 ] => [ y0 y1 y2 y3 ]
//   [ x2 y2 z2 w2 ]    [ z0 z1 z2 z3 ]
//   [ x3 y3 z3 w3 ]    [ w0 w1 w2 w3 ]
func (m *M4) Transpose() *M4 {
	t_y0 := m.Y0
	t_z0 := m.Z0
	t_w0 := m.W0
	t_z1 := m.Z1
	t_w1 := m.W1
	t_w2 := m.W2
	m.X0, m.Y0, m.Z0, m.W0 = m.X0, m.X1, m.X2, m.X3
	m.X1, m.Y1, m.Z1, m.W1 = t_y0, m.Y1, m.Y2, m.Y3
	m.X2, m.Y2, m.Z2, m.W2 = t_z0, t_z1, m.Z2, m.Z3
	m.X3, m.Y3, m.Z3, m.W3 = t_w0, t_w1, t_w2, m.W3
	return m
}

// Mult multiplies a matrix m (on the left) by another matrix r (on the right)
// and return the result in the calling matrix m.
//
// Matrix multiplication always takes rows from the left matrix m
// and multiplies (dot products) them against columns from the right
// matrix to get the new matrix values.
func (m *M3) Mult(r *M3) *M3 {
	x0 := m.X0*r.X0 + m.Y0*r.X1 + m.Z0*r.X2
	y0 := m.X0*r.Y0 + m.Y0*r.Y1 + m.Z0*r.Y2
	z0 := m.X0*r.Z0 + m.Y0*r.Z1 + m.Z0*r.Z2
	x1 := m.X1*r.X0 + m.Y1*r.X1 + m.Z1*r.X2
	y1 := m.X1*r.Y0 + m.Y1*r.Y1 + m.Z1*r.Y2
	z1 := m.X1*r.Z0 + m.Y1*r.Z1 + m.Z1*r.Z2
	x2 := m.X2*r.X0 + m.Y2*r.X1 + m.Z2*r.X2
	y2 := m.X2*r.Y0 + m.Y2*r.Y1 + m.Z2*r.Y2
	z2 := m.X2*r.Z0 + m.Y2*r.Z1 + m.Z2*r.Z2
	m.X0, m.Y0, m.Z0 = x0, y0, z0
	m.X1, m.Y1, m.Z1 = x1, y1, z1
	m.X2, m.Y2, m.Z2 = x2, y2, z2
	return m
}

// Mult multiplies a matrix m (on the left) by another matrix r (on the right)
// and return the result in the calling matrix m.
// Same behaviour as M3.Mult()
func (m *M4) Mult(r *M4) *M4 {
	x0 := m.X0*r.X0 + m.Y0*r.X1 + m.Z0*r.X2 + m.W0*r.X3
	y0 := m.X0*r.Y0 + m.Y0*r.Y1 + m.Z0*r.Y2 + m.W0*r.Y3
	z0 := m.X0*r.Z0 + m.Y0*r.Z1 + m.Z0*r.Z2 + m.W0*r.Z3
	w0 := m.X0*r.W0 + m.Y0*r.W1 + m.Z0*r.W2 + m.W0*r.W3
	x1 := m.X1*r.X0 + m.Y1*r.X1 + m.Z1*r.X2 + m.W1*r.X3
	y1 := m.X1*r.Y0 + m.Y1*r.Y1 + m.Z1*r.Y2 + m.W1*r.Y3
	z1 := m.X1*r.Z0 + m.Y1*r.Z1 + m.Z1*r.Z2 + m.W1*r.Z3
	w1 := m.X1*r.W0 + m.Y1*r.W1 + m.Z1*r.W2 + m.W1*r.W3
	x2 := m.X2*r.X0 + m.Y2*r.X1 + m.Z2*r.X2 + m.W2*r.X3
	y2 := m.X2*r.Y0 + m.Y2*r.Y1 + m.Z2*r.Y2 + m.W2*r.Y3
	z2 := m.X2*r.Z0 + m.Y2*r.Z1 + m.Z2*r.Z2 + m.W2*r.Z3
	w2 := m.X2*r.W0 + m.Y2*r.W1 + m.Z2*r.W2 + m.W2*r.W3
	x3 := m.X3*r.X0 + m.Y3*r.X1 + m.Z3*r.X2 + m.W3*r.X3
	y3 := m.X3*r.Y0 + m.Y3*r.Y1 + m.Z3*r.Y2 + m.W3*r.Y3
	z3 := m.X3*r.Z0 + m.Y3*r.Z1 + m.Z3*r.Z2 + m.W3*r.Z3
	w3 := m.X3*r.W0 + m.Y3*r.W1 + m.Z3*r.W2 + m.W3*r.W3
	m.X0, m.Y0, m.Z0, m.W0 = x0, y0, z0, w0
	m.X1, m.Y1, m.Z1, m.W1 = x1, y1, z1, w1
	m.X2, m.Y2, m.Z2, m.W2 = x2, y2, z2, w2
	m.X3, m.Y3, m.Z3, m.W3 = x3, y3, z3, w3
	return m
}

// M4Translater creates a new translation matrix based on the given
// x, y, z values.  The new  matrix is returned so that it may be immediately
// used in another operation.
//
//   [ x0 y0 z0 w0 ]    [ 1 0 0 0 ]
//   [ x1 y1 z1 w1 ] => [ 0 1 0 0 ]
//   [ x2 y2 z2 w2 ]    [ 0 0 1 0 ]
//   [ x3 y3 z3 w3 ]    [ x y z 1 ]
func M4Translater(x, y, z float32) *M4 {
	return &M4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		x, y, z, 1,
	}
}

// TranslateL translates a matrix by treating the input x, y, z as if
// it was a translation matrix used on the left hand side of a matrix
// multiplication with the calling matrix m.  Only the translation
// coordinates are affected.  The updated matrix m is returned so that
// it may be immediately used in another operation.
//
//   [ 1 0 0 0 ]   [ x0 y0 z0 w0 ]     [ x0  y0  z0  w0 ]
//   [ 0 1 0 0 ] x [ x1 y1 z1 w1 ]  => [ x1  y1  z1  w1 ]
//   [ 0 0 1 0 ]   [ x2 y2 z2 w2 ]     [ x2  y2  z2  w2 ]
//   [ x y z 1 ]   [ x3 y3 z3 w3 ]     [ x3' y3' z3' w3 ]
//
// x, y, z are the coordinates of a translation vector.
//
// TranstateL is used to avoid the creation time of a translation matrix.
// Be sure to pick the correct translate (L or R) when doing transforms.
// Generally its TranslateR since translate is the last multiplication
// (given that M4 is in row major order).
func (m *M4) TranslateL(x, y, z float32) *M4 {
	m.X3 += x*m.X0 + y*m.X1 + z*m.X2
	m.Y3 += x*m.Y0 + y*m.Y1 + z*m.Y2
	m.Z3 += x*m.Z0 + y*m.Z1 + z*m.Z2
	m.W3 += x*m.W0 + y*m.W1 + z*m.W2
	return m
}

// TranslateR translates a matrix by treating the input x, y, z as if
// it was a translation matrix used on the right hand side of a matrix
// multiplication with the calling matrix m.  The updated matrix m is
// returned so that it may be immediately used in another operation.
//
//   [ x0 y0 z0 w0 ]   [ 1 0 0 0 ]    [ x0'  y0' z0' w0 ]
//   [ x1 y1 z1 w1 ] x [ 0 1 0 0 ] => [ x1'  y1' z1' w1 ]
//   [ x2 y2 z2 w2 ]   [ 0 0 1 0 ]    [ x2'  y2' z2' w2 ]
//   [ x3 y3 z3 w3 ]   [ x y z 1 ]    [ x3'  y3' z3' w3 ]
//
// x, y, z are the coordinates of a translation vector.
//
// Be sure to pick the correct translate (L or R) when doing transforms.
// Generally its TranslateR since translate is the last multiplication
// (given that M4 is in row major order).
func (m *M4) TranslateR(x, y, z float32) *M4 {
	m.X0 += m.W0 * x
	m.Y0 += m.W0 * y
	m.Z0 += m.W0 * z
	m.X1 += m.W1 * x
	m.Y1 += m.W1 * y
	m.Z1 += m.W1 * z
	m.X2 += m.W2 * x
	m.Y2 += m.W2 * y
	m.Z2 += m.W2 * z
	m.X3 += m.W3 * x
	m.Y3 += m.W3 * y
	m.Z3 += m.W3 * z
	return m
}

// M4Scaler creates a new scale matrix from the given vector. The new matrix
// is returned so that it may be immediately used in another operation.
//
//   [ x0 y0 z0 w0 ]    [ x 0 0 0 ]
//   [ x1 y1 z1 w1 ] => [ 0 y 0 0 ]
//   [ x2 y2 z2 w2 ]    [ 0 0 z 0 ]
//   [ x3 y3 z3 w3 ]    [ 0 0 0 1 ]
func M4Scaler(x, y, z float32) *M4 {
	return &M4{X0: x, Y1: y, Z2: z, W3: 1}
}

// ScaleL scales a matrix by treating the input x, y, z as if it was a scaler
// matrix used on the left hand side in a matrix multiplication with the
// calling matrix. Notice that the translation vector is not touched.
// The updated matrix m is returned so that it may be immediately
// used in another operation.
//
//   [ x 0 0 0 ]   [ x0 y0 z0 w0 ]     [ x0' y0' z0' w0' ]
//   [ 0 y 0 0 ] x [ x1 y1 z1 w1 ]  => [ x1' y1' z1' w1' ]
//   [ 0 0 z 0 ]   [ x2 y2 z2 w2 ]     [ x2' y2' z2' w2' ]
//   [ 0 0 0 1 ]   [ x3 y3 z3 w3 ]     [ x3  y3  z3  w3  ]
//
// Be sure to pick the correct scale (L or R) when doing transforms.
// Generally its ScaleL since scale is the first multiplication on the left
// (given that M4 is in row major order).
func (m *M4) ScaleL(xScale, yScale, zScale float32) *M4 {
	m.X0 *= xScale
	m.Y0 *= xScale
	m.Z0 *= xScale
	m.W0 *= xScale
	m.X1 *= yScale
	m.Y1 *= yScale
	m.Z1 *= yScale
	m.W1 *= yScale
	m.X2 *= zScale
	m.Y2 *= zScale
	m.Z2 *= zScale
	m.W2 *= zScale
	return m
}

// ScaleR scales a matrix by treating the input x, y, z as if it was a scaler
// matrix used on the right hand side in a matrix multiplication with the
// calling matrix.  The updated matrix m is returned so that it may be
// immediately used in another operation.
//
//   [ x0 y0 z0 w0 ]    [ x 0 0 0 ]    [ x0' y0' z0' w0 ]
//   [ x1 y1 z1 w1 ] x  [ 0 y 0 0 ] => [ x1' y1' z1' w1 ]
//   [ x2 y2 z2 w2 ]    [ 0 0 z 0 ]    [ x2' y2' z2' w2 ]
//   [ x3 y3 z3 w3 ]    [ 0 0 0 1 ]    [ x3' y3' z3' w3 ]
//
// Be sure to pick the correct scale (L or R) when doing transforms.
// Generally its ScaleL since scale is the first multiplication on the left
// (given that M4 is in row major order).
func (m *M4) ScaleR(xScale, yScale, zScale float32) *M4 {
	m.X0 *= xScale
	m.Y0 *= yScale
	m.Z0 *= zScale
	m.X1 *= xScale
	m.Y1 *= yScale
	m.Z1 *= zScale
	m.X2 *= xScale
	m.Y2 *= yScale
	m.Z2 *= zScale
	m.X3 *= xScale
	m.Y3 *= yScale
	m.Z3 *= zScale
	return m
}

// M4Orthographic creates a new 4x4 matrix with projection values needed to
// transform a 3 dimensional model to a 2 dimensional plane.
// Orthographic projection ignores depth. The input arguments are:
//     left, right:  Vertical clipping planes.
//     bottom, top:  Horizontal clipping planes.
//     near, far  :  Depth clipping planes. The depth values are
//                   negative if the plane is to be behind the viewer
//
// An orthographic matrix fills the following matrix locations:
//   [ a 0 0 0 ]
//   [ 0 b 0 0 ]
//   [ 0 0 c 0 ]
//   [ d e f 1 ]
func M4Orthographic(left, right, bottom, top, near, far float32) *M4 {
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

// M4Perspective creates a new 4x4 matrix with projection values needed to
// transform a 3 dimentional model to a 2 dimensional plane.
// Objects that are further away from the viewer will appear smaller.
// The input arguments are:
//    fov        An amount in degrees indicating how much of the
//               scene is visible.
//    apsect     The ratio of hieght to width of the model.
//    near, far  The depth clipping planes. The depth values are
//               negative if the plane is to be behind the viewer
//
// A perspective matrix fills the following matrix locations:
//   [ a 0 0 0 ]
//   [ 0 b 0 0 ]
//   [ 0 0 c d ]
//   [ 0 0 e 0 ]
func M4Perspective(fov, aspect, near, far float32) *M4 {
	m := &M4{}
	f := 1 / float32(math.Tan(float64(fov)*PI_OVER_180*0.5))
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

// M4PerspectiveI creates a new inverse matrix of the given perspective
// matrix values (see M4Perspective).
//   [ a' 0  0  0 ] where a' = 1/a       d' = 1/e
//   [ 0  b' 0  0 ]       b' = 1/b       e' = 1/d
//   [ 0  0  0  d']       c' = -(c/de)
//   [ 0  0  e' c']
//
// This is used when going from screen x,y coordinates to 3D coordinates.
// as in the case when creating a picking ray from a mouse location.
func M4PerspectiveI(fov, aspect, near, far float32) *M4 {
	m := &M4{}
	f := float32(math.Tan(float64(fov) * PI_OVER_180 * 0.5))
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

// IModelView creates an inverse model-view matrix.
// The idea is to seprately reverse the rotation and translation effects
// that make up the model-view matrix.  This is done by transposing the
// rotation portion and reversing the translation.
//   [ x0 y0 z0 w0 ]   [ rx0 ry0 rz0 0 ]    [ rx0 rx1 rx2 0 ]
//   [ x1 y1 z1 w1 ] = [ rx1 ry1 rz1 0 ] => [ ry0 ry1 ry2 0 ]
//   [ x2 y2 z2 w2 ]   [ rx2 ry2 rz2 0 ]    [ rz0 rz1 rz2 0 ]
//   [ x3 y3 z3 w3 ]   [  tx  ty  tz 1 ]    [ -ix -iy -iz 1 ]
//
// This is used when going from screen x,y coordinates to 3D coordinates.
// as in the case when creating a picking ray from a mouse location.
//
// This inverse does not take scaling into account and will not work with
// model objects that have had scaling applied.  A full inverse will be
// needed in that case (likely find another way to do it since a full
// matrix inverse is expensive).
func (m *M4) IModelView() *M4 {
	r := m.M3().Transpose()
	t := &V3{m.X3, m.Y3, m.Z3}
	t.MultL(r)
	m.X0, m.Y0, m.Z0, m.W0 = r.X0, r.Y0, r.Z0, 0
	m.X1, m.Y1, m.Z1, m.W1 = r.X1, r.Y1, r.Z1, 0
	m.X2, m.Y2, m.Z2, m.W2 = r.X2, r.Y2, r.Z2, 0
	m.X3, m.Y3, m.Z3, m.W3 = -t.X, -t.Y, -t.Z, 1
	return m
}
