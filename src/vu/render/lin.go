// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package render

// M3 is a 3x3 float32 matrix that is populated from the more precise
// math/lin float64 representation.
type M3 struct {
	X0, Y0, Z0 float32 // row 1 : indicies 0, 1, 2 [00, 01, 02]
	X1, Y1, Z1 float32 // row 2 : indicies 3, 4, 5 [10, 11, 12]
	X2, Y2, Z2 float32 // row 3 : indicies 6, 7, 8 [20, 21, 22]
}

// M3 updates calling matrix m to be the 3x3 matrix from the top left corner
// of the given 4x4 matrix m4. The source matrix m4 is unchanged.
//    [ x0 y0 z0 w0 ]    [ x0 y0 z0 ]
//    [ x1 y1 z1 w1 ] => [ x1 y1 z1 ]
//    [ x2 y2 z2 w2 ]    [ x2 y2 z2 ]
//    [ x3 y3 z3 w3 ]
func (m *M3) M3(m4 *M4) *M3 {
	m.X0, m.Y0, m.Z0 = m4.X0, m4.Y0, m4.Z0
	m.X1, m.Y1, m.Z1 = m4.X1, m4.Y1, m4.Z1
	m.X2, m.Y2, m.Z2 = m4.X2, m4.Y2, m4.Z2
	return m
}

// M4 is a 4x4 float32 matrix that is populated from the more precise
// math/lin float64 representation.
type M4 struct {
	X0, Y0, Z0, W0 float32 // row 1 : indicies 0, 1, 2, 3 [00, 01, 02, 03]
	X1, Y1, Z1, W1 float32 // row 2 : indicies 4, 5, 6, 7 [10, 11, 12, 13]
	X2, Y2, Z2, W2 float32 // row 3 : indicies 8, 9, a, b [20, 21, 22, 23]
	X3, Y3, Z3, W3 float32 // row 4 : indicies c, d, e, f [30, 31, 32, 33]
}

// Pointer is used to access the matrix data as an array of floats.
// Used to pass the matrix to native graphic layer.
func (m *M3) Pointer() *float32 { return &(m.X0) }

// Pointer is used to access the matrix data as an array of floats.
// Used to pass the matrix to native graphic layer.
func (m *M4) Pointer() *float32 { return &(m.X0) }

// V3 is a float32 based vector that is populated from the more precise
// math/physics float64 representation.
type V3 struct {
	X, Y, Z float32
}
