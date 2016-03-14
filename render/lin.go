// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package render

import (
	"github.com/gazed/vu/math/lin"
)

// lin hides the fact that the current underlying graphics implementation
// deals in float32 rather than float64 used by Go and vu/math/lin.
// These are kept package local because it is expected that GPU's will
// transition from 32 to 64 bit and then these 32 bit structures and
// conversions can disappear.
//
// These are data holders only. Please keep all math operations
// restricted to vu/math/lin.

// m3 is a 3x3 float32 matrix that is populated from the more precise
// math/lin float64 representation.
type m3 struct {
	xx, xy, xz float32 // indices 0, 1, 2  [00, 01, 02]  X-Axis
	yx, yy, yz float32 // indices 3, 4, 5  [10, 11, 12]  Y-Axis
	zx, zy, zz float32 // indices 6, 7, 8  [20, 21, 22]  Z-Axis
}

// Pointer accesses the matrix data as an array of floats.
// Used to pass the matrix to native graphic layer.
func (m *m3) Pointer() *float32 { return &(m.xx) }

// m3 updates calling matrix m to be the 3x3 matrix from the top left corner
// of the given 4x4 matrix m4. The source matrix m4 is unchanged.
//    [ x0 y0 z0 w0 ]    [ x0 y0 z0 ]
//    [ x1 y1 z1 w1 ] => [ x1 y1 z1 ]
//    [ x2 y2 z2 w2 ]    [ x2 y2 z2 ]
//    [ x3 y3 z3 w3 ]
func (m *m3) m3(f *m4) *m3 {
	m.xx, m.xy, m.xz = f.xx, f.xy, f.xz
	m.yx, m.yy, m.yz = f.yx, f.yy, f.yz
	m.zx, m.zy, m.zz = f.zx, f.zy, f.zz
	return m
}

// =============================================================================

// m34 is a 3x4 float32 column-major matrix that is populated from the more
// precise math/lin float64 representation. It becomes row-major when Sent to
// the GPU without transposing. It is used as an internal optimization to send
// 4 less floats for each bone transform matrix. The shader is expected be
// aware of this space saving layout.
type m34 struct {
	xx, yx, zx, wx float32 // indices 0, 1, 2, 3  [00, 01, 02, 03]
	xy, yy, zy, wy float32 // indices 4, 5, 6, 7  [10, 11, 12, 13]
	xz, yz, zz, wz float32 // indices 8, 9, a, b  [20, 21, 22, 23]
	// 0, 0, 0, 1 implicit last row.
}

// Pointer is used to access the matrix data as an array of floats.
// Used to pass the matrix to native graphic layer.
func (m *m34) Pointer() *float32 { return &(m.xx) }

// toM4 translates the m34 column-major matrix to a M4 row-major matrix.
func (m *m34) toM4(mm *lin.M4) *lin.M4 {
	mm.Xx, mm.Xy, mm.Xz, mm.Xw = float64(m.xx), float64(m.xy), float64(m.xz), 0
	mm.Yx, mm.Yy, mm.Yz, mm.Yw = float64(m.yx), float64(m.yy), float64(m.yz), 0
	mm.Zx, mm.Zy, mm.Zz, mm.Zw = float64(m.zx), float64(m.zy), float64(m.zz), 0
	mm.Wx, mm.Wy, mm.Wz, mm.Ww = float64(m.wx), float64(m.wy), float64(m.wz), 1
	return mm
}

// tom34 translates the M4 row-major matrix to a m34 column-major matrix.
// This in turn is expected to be reinterpreted as a row-major matrix by the
// GPU shader.
func (m *m34) tom34(mm *lin.M4) *m34 {
	m.xx, m.yx, m.zx, m.wx = float32(mm.Xx), float32(mm.Yx), float32(mm.Zx), float32(mm.Wx)
	m.xy, m.yy, m.zy, m.wy = float32(mm.Xy), float32(mm.Yy), float32(mm.Zy), float32(mm.Wy)
	m.xz, m.yz, m.zz, m.wz = float32(mm.Xz), float32(mm.Yz), float32(mm.Zz), float32(mm.Wz)
	return m
}

// =============================================================================

// m4 is a 4x4 float32 matrix that is populated from the more precise
// math/lin float64 representation.
type m4 struct {
	xx, xy, xz, xw float32 // indices 0, 1, 2, 3  [00, 01, 02, 03] X-Axis
	yx, yy, yz, yw float32 // indices 4, 5, 6, 7  [10, 11, 12, 13] Y-Axis
	zx, zy, zz, zw float32 // indices 8, 9, a, b  [20, 21, 22, 23] Z-Axis
	wx, wy, wz, ww float32 // indices c, d, e, f  [30, 31, 32, 33]
}

// Mvp makes m4 compatible for the Mvp interface.
func (m *m4) Set(mm *lin.M4) Mvp { return m.tom4(mm) }

// Pointer is used to access the matrix data as an array of floats.
// Used to pass the matrix to native graphic layer.
func (m *m4) Pointer() *float32 { return &(m.xx) }

// tom4 turns a math/lin matrix into a matrix that can be used
// by the render system. The input math matrix, mm, is used to fill the values
// in the given render matrix rm.  The updated rm matrix is returned.
func (m *m4) tom4(mm *lin.M4) *m4 {
	m.xx, m.xy, m.xz, m.xw = float32(mm.Xx), float32(mm.Xy), float32(mm.Xz), float32(mm.Xw)
	m.yx, m.yy, m.yz, m.yw = float32(mm.Yx), float32(mm.Yy), float32(mm.Yz), float32(mm.Yw)
	m.zx, m.zy, m.zz, m.zw = float32(mm.Zx), float32(mm.Zy), float32(mm.Zz), float32(mm.Zw)
	m.wx, m.wy, m.wz, m.ww = float32(mm.Wx), float32(mm.Wy), float32(mm.Wz), float32(mm.Ww)
	return m
}

// =============================================================================

// v3 is a float32 based vector that is populated from the more precise
// math/physics float64 representation.
type v3 struct {
	x, y, z float32
}

// =============================================================================

// Mvp exposes the render matrix representation. This is needed by
// applications using the vu/render system, but not the vu engine.
type Mvp interface {
	Set(tm *lin.M4) Mvp // Converts the transform matrix tm to internal data.
	Pointer() *float32  // A pointer to the internal transform data.
}

// NewMvp creates a new internal render transform matrix.
// This is needed by applications using the vu/render system,
// but not the vu engine.
func NewMvp() Mvp { return &m4{} }
