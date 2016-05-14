// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package lin

import (
	"fmt"
	"math"
	"testing"
)

func TestAeqmately(t *testing.T) {
	var f1 = 0.0
	var f2 = 0.000001
	var f3 = -0.0001
	if Aeq(f1, f2) && !Aeq(f1, f3) {
		t.Error("Aeq")
	}
}

func TestApproimatelyZero(t *testing.T) {
	var f1 = 0.0000001
	var f2 = -0.0000001
	var f3 = -0.0001
	if !AeqZ(f1) || !AeqZ(f2) || AeqZ(f3) {
		t.Error("Aeqz")
	}
}

func TestLerp(t *testing.T) {
	if !Aeq(Lerp(10, 5, 0.5), 7.5) {
		t.Error("Lerp")
	}
}

// Check that the results of Atan2 and Atan2F are similar.
func TestAtan2F(t *testing.T) {
	if !Aeq(math.Atan2(1, 0), Atan2F(1, 0)) || !Aeq(math.Atan2(-1, 0), Atan2F(-1, 0)) {
		t.Error("Atan2F")
	}
}

func TestNang(t *testing.T) {
	pos450, neg450 := 7.853981, -7.853981
	pos90, neg90 := 1.570796, -1.570796
	if !Aeq(Nang(pos450), pos90) || !Aeq(Nang(neg450), neg90) {
		t.Error("Nang")
	}
}

func TestClamp(t *testing.T) {
	if Clamp(20, -30, -15) != -15 || Clamp(20, 30, 60) != 30 || Clamp(20, 10, 50) != 20 {
		t.Error("Clamp")
	}
}

func TestRadDeg(t *testing.T) {
	if Deg(Rad(90)) != 90 {
		t.Error("Rad Deg conversion")
	}
}

func TestRound(t *testing.T) {
	f1, f2 := Round(1.48, 0), Round(1.51, 0)
	if f1 != 1.0 || f2 != 2.0 {
		t.Error("Failed to round floats", f1, f2)
	}
	f1, f2 = Round(-0.49, 0), Round(0.49, 0)
	if f1 != 0.0 || f2 != 0.0 {
		t.Error("Failed to round floats", f1, f2)
	}
}

// ============================================================================
// Benchmarking

// Ensure that Atan2F really is faster. Run 'go test -bench=".*"
// For example the last run showed:
//    BenchmarkAtan2     50000000	53.0 ns/op
//    BenchmarkAtan2F 	100000000	12.3 ns/op
func BenchmarkAtan2(b *testing.B) {
	for cnt := 0; cnt < b.N; cnt++ {
		math.Atan2(1, 1)
	}
}
func BenchmarkAtan2F(b *testing.B) {
	for cnt := 0; cnt < b.N; cnt++ {
		Atan2F(1, 1)
	}
}

// ============================================================================
// Test helpers for the other test case files in this package.

// Dictate how errors get printed.
const format = "\ngot\n%s\nwanted\n%s"

// Dumps the matrix to a string.
func (m *M3) Dump() string {
	format := "[%+2.9f, %+2.9f, %+2.9f]\n"
	str := fmt.Sprintf(format, m.Xx, m.Xy, m.Xz)
	str += fmt.Sprintf(format, m.Yx, m.Yy, m.Yz)
	str += fmt.Sprintf(format, m.Zx, m.Zy, m.Zz)
	return str
}

// Dump like M3.Dump().
func (m *M4) Dump() string {
	format := "[%+2.9f, %+2.9f, %+2.9f, %+2.9f]\n"
	str := fmt.Sprintf(format, m.Xx, m.Xy, m.Xz, m.Xw)
	str += fmt.Sprintf(format, m.Yx, m.Yy, m.Yz, m.Yw)
	str += fmt.Sprintf(format, m.Zx, m.Zy, m.Zz, m.Zw)
	str += fmt.Sprintf(format, m.Wx, m.Wy, m.Wz, m.Ww)
	return str
}

// Convienience methods for getting a vector as a string.
func (v *V3) Dump() string { return fmt.Sprintf("%2.9f", *v) }
func (v *V4) Dump() string { return fmt.Sprintf("%2.9f", *v) }

// Convienience method for getting a quaternion as a string.
func (q *Q) Dump() string { return fmt.Sprintf("%2.9f", *q) }
