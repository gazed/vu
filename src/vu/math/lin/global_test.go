// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package lin

// Global helpers used only in testing.

import "fmt"

// Dictate how errors get printed.
const format = "\ngot\n%s\nwanted\n%s"

// Dumps the matrix to a string in row-major order.
//
//   [ x0 y0 z0 ]
//   [ x1 y1 z1 ]
//   [ x2 y2 z2 ]
func (m *M3) Dump() string {
	format := "[%+2.2f, %+2.2f, %+2.2f]\n"
	str := fmt.Sprintf(format, m.X0, m.Y0, m.Z0)
	str += fmt.Sprintf(format, m.X1, m.Y1, m.Z1)
	str += fmt.Sprintf(format, m.X2, m.Y2, m.Z2)
	return str
}

// Dump like M3.Dump().
//
//   [ x0 y0 z0 w0 ]
//   [ x1 y1 z1 w1 ]
//   [ x2 y2 z2 w2 ]
//   [ x3 y3 z3 w3 ]
func (m *M4) Dump() string {
	format := "[%+2.2f, %+2.2f, %+2.2f, %+2.2f]\n"
	str := fmt.Sprintf(format, m.X0, m.Y0, m.Z0, m.W0)
	str += fmt.Sprintf(format, m.X1, m.Y1, m.Z1, m.W1)
	str += fmt.Sprintf(format, m.X2, m.Y2, m.Z2, m.W2)
	str += fmt.Sprintf(format, m.X3, m.Y3, m.Z3, m.W3)
	return str
}

// Convienience methods for getting a vector as a string.
func (v *V3) Dump() string { return fmt.Sprintf("%2.1f", *v) }
func (v *V4) Dump() string { return fmt.Sprintf("%2.1f", *v) }

// Convienience method for getting a quaternion as a string.
func (q *Q) Dump() string { return fmt.Sprintf("%2.1f", *q) }
