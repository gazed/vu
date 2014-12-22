// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package lin

import (
	"testing"
)

func TestSetEqualsM3(t *testing.T) {
	m, a := &M3{},
		&M3{11, 12, 13,
			21, 22, 23,
			31, 32, 33}
	if !m.Set(a).Eq(a) {
		t.Errorf(format, m.Dump(), a.Dump())
	}
}
func TestSetEqualsM4(t *testing.T) {
	m, a := &M4{},
		&M4{11, 12, 13, 14,
			21, 22, 23, 24,
			31, 32, 33, 34,
			41, 42, 43, 44}
	if !m.Set(a).Eq(a) {
		t.Errorf(format, m.Dump(), a.Dump())
	}
}

func TestSetM3(t *testing.T) {
	m, m4, want := &M3{},
		&M4{11, 12, 13, 14,
			21, 22, 23, 24,
			31, 32, 33, 34,
			41, 42, 43, 44},
		&M3{11, 12, 13,
			21, 22, 23,
			31, 32, 33}
	if !m.SetM4(m4).Eq(want) {
		t.Errorf(format, m.Dump(), want.Dump())
	}
}

func TestAbsM3(t *testing.T) {
	m, want :=
		&M3{-11, -12, +13,
			+21, -22, +23,
			+31, -32, -33},
		&M3{11, 12, 13,
			21, 22, 23,
			31, 32, 33}
	if !m.Abs(m).Eq(want) {
		t.Errorf(format, m.Dump(), want.Dump())
	}
}

func TestTransposeM3(t *testing.T) {
	m, want :=
		&M3{1, 2, 3,
			4, 5, 6,
			7, 8, 9},
		&M3{1, 4, 7,
			2, 5, 8,
			3, 6, 9}
	if !m.Transpose(m).Eq(want) {
		t.Errorf(format, m.Dump(), want.Dump())
	}
}
func TestTransposeM4(t *testing.T) {
	m, want :=
		&M4{11, 12, 13, 14,
			21, 22, 23, 24,
			31, 32, 33, 34,
			41, 42, 43, 44},
		&M4{11, 21, 31, 41,
			12, 22, 32, 42,
			13, 23, 33, 43,
			14, 24, 34, 44}
	if !m.Transpose(m).Eq(want) {
		t.Errorf(format, m.Dump(), want.Dump())
	}
}

func TestAddM3(t *testing.T) {
	m, want :=
		&M3{11, 12, 13,
			21, 22, 23,
			31, 32, 33},
		&M3{22, 24, 26,
			42, 44, 46,
			62, 64, 66}
	if !m.Add(m, m).Eq(want) {
		t.Errorf(format, m.Dump(), want.Dump())
	}
}

func TestAddM4(t *testing.T) {
	m, want :=
		&M4{11, 12, 13, 14,
			21, 22, 23, 24,
			31, 32, 33, 34,
			41, 42, 43, 44},
		&M4{22, 24, 26, 28,
			42, 44, 46, 48,
			62, 64, 66, 68,
			82, 84, 86, 88}
	if !m.Add(m, m).Eq(want) {
		t.Errorf(format, m.Dump(), want.Dump())
	}
}

func TestSubM3(t *testing.T) {
	m :=
		&M3{-11, -12, +13,
			+21, -22, +23,
			+31, -32, -33}
	if !m.Sub(m, m).Eq(M3Z) {
		t.Errorf(format, m.Dump(), M3Z.Dump())
	}
}

func TestMultiplyM3(t *testing.T) {
	m, want :=
		&M3{1, 2, 3,
			4, 5, 6,
			7, 8, 9},
		&M3{30, 36, 42,
			66, 81, 96,
			102, 126, 150}
	if !m.Mult(m, m).Eq(want) {
		t.Errorf(format, m.Dump(), want.Dump())
	}
}

func TestMultiplyM4(t *testing.T) {
	m, want :=
		&M4{1, 2, 3, 4,
			5, 6, 7, 8,
			9, 10, 11, 12,
			13, 14, 15, 16},
		&M4{90, 100, 110, 120,
			202, 228, 254, 280,
			314, 356, 398, 440,
			426, 484, 542, 600}
	if !m.Mult(m, m).Eq(want) {
		t.Errorf(format, m.Dump(), want.Dump())
	}
}

func TestMultLtR(t *testing.T) {
	m, want :=
		&M3{1, 2, 3,
			4, 5, 6,
			7, 8, 9},
		&M3{66, 78, 90,
			78, 93, 108,
			90, 108, 126}
	if !m.MultLtR(m, m).Eq(want) {
		t.Errorf(format, m.Dump(), want.Dump())
	}
}

func TestTranslateTM(t *testing.T) {
	m, want :=
		&M4{1, 2, 3, 4,
			1, 2, 3, 4,
			1, 2, 3, 4,
			1, 2, 3, 4},
		&M4{1, 2, 3, 4,
			1, 2, 3, 4,
			1, 2, 3, 4,
			7, 14, 21, 28}
	if !m.TranslateTM(1, 2, 3).Eq(want) {
		t.Errorf(format, m.Dump(), want.Dump())
	}
}

func TestTranslateMT(t *testing.T) {
	m, want :=
		&M4{1, 2, 3, 4,
			1, 2, 3, 4,
			1, 2, 3, 4,
			1, 2, 3, 4},
		&M4{5, 10, 15, 4,
			5, 10, 15, 4,
			5, 10, 15, 4,
			5, 10, 15, 4}
	if !m.TranslateMT(1, 2, 3).Eq(want) {
		t.Errorf(format, m.Dump(), want.Dump())
	}
}

func TestScaleM3SM(t *testing.T) {
	m, want :=
		&M3{1, 2, 3,
			1, 2, 3,
			1, 2, 3},
		&M3{1, 2, 3,
			2, 4, 6,
			3, 6, 9}
	if !m.ScaleSM(1, 2, 3).Eq(want) {
		t.Errorf(format, m.Dump(), want.Dump())
	}
}

func TestScaleM4SM(t *testing.T) {
	m, want :=
		&M4{1, 2, 3, 4,
			1, 2, 3, 4,
			1, 2, 3, 4,
			1, 2, 3, 4},
		&M4{1, 2, 3, 4,
			2, 4, 6, 8,
			3, 6, 9, 12,
			1, 2, 3, 4}
	if !m.ScaleSM(1, 2, 3).Eq(want) {
		t.Errorf(format, m.Dump(), want.Dump())
	}
}

func TestScaleMS(t *testing.T) {
	m, want :=
		&M4{1, 2, 3, 4,
			1, 2, 3, 4,
			1, 2, 3, 4,
			1, 2, 3, 4},
		&M4{1, 4, 9, 4,
			1, 4, 9, 4,
			1, 4, 9, 4,
			1, 4, 9, 4}
	if !m.ScaleMS(1, 2, 3).Eq(want) {
		t.Errorf(format, m.Dump(), want.Dump())
	}
}

func TestSetQ(t *testing.T) {
	m, q, want := &M3{}, &Q{0.2, 0.4, 0.5, 0.7},
		&M3{+0.18, -0.54, +0.76,
			+0.86, +0.42, +0.12,
			-0.36, +0.68, +0.60}
	if !m.SetQ(q).Aeq(want) {
		t.Errorf(format, m.Dump(), want.Dump())
	}

	// check identity quaternion
	q, want = &Q{0, 0, 0, 1},
		&M3{1, 0, 0,
			0, 1, 0,
			0, 0, 1}
	if !m.SetQ(q).Eq(want) {
		t.Errorf(format, m.Dump(), want.Dump())
	}
}

func TestSetSkewSymetric(t *testing.T) {
	m, mi, v := &M3{}, &M3{}, &V3{1, 2, 3}
	m.SetSkewSym(v)            // the skew symetric matrix
	mi.Transpose(m)            // its transpose (which is its negative)
	if !m.Add(m, mi).Eq(M3Z) { // add the negative should be zero.
		t.Errorf(format, m.Dump(), M3Z.Dump())
	}
}

// See http://www.wikihow.com/Inverse-a-3X3-Matrix
func TestDeterminantM3(t *testing.T) {
	m :=
		&M3{1, 2, 3,
			4, 5, 6,
			7, 8, 9}
	if m.Det() != 0 {
		t.Error("No inverse possible for m, determinant should be 0")
	}
	m =
		&M3{1, 2, 3,
			0, 1, 4,
			5, 6, 0}
	if m.Det() != 1 {
		t.Error("Inverse possible for m, determinant should be non-zero")
	}
}

// Also tests all possible permutations of M3.Cofac (cofactor).
// See http://www.wikihow.com/Inverse-a-3X3-Matrix
func TestAdjointM3(t *testing.T) {
	m, want :=
		&M3{1, 2, 3,
			0, 1, 4,
			5, 6, 0},
		&M3{-24, 18, 5,
			20, -15, -4,
			-5, 4, 1}
	if !m.Adj(m).Eq(want) {
		t.Errorf(format, m.Dump(), want.Dump())
	}
}

// See http://www.wikihow.com/Inverse-a-3X3-Matrix
func TestInvM3(t *testing.T) {
	m, a := &M3{},
		&M3{1, 2, 3,
			0, 1, 4,
			5, 6, 0}
	m.Inv(a)
	if !NewM3().Mult(m, a).Eq(M3I) {
		t.Errorf(format, m.Dump(), a.Dump())
	}
}

func TestSetAxisAngle(t *testing.T) {
	m, want := &M3{},
		&M3{1, 0, 0, // rotation 90 degrees around X.
			0, 0, -1,
			0, 1, 0}
	if !m.SetAa(1, 0, 0, Rad(90)).Aeq(want) {
		t.Errorf(format, m.Dump(), want.Dump())
	}

	// same check with quaternion.
	q := NewQ().SetAa(1, 0, 0, Rad(90))
	if !m.SetQ(q).Aeq(want) {
		t.Errorf(format, m.Dump(), want.Dump())
	}
}

func TestOrthographicM4(t *testing.T) {
	m, want := NewM4().Ortho(2, 3, 4, 5, 6, 7),
		&M4{+2, +0, +0, +0,
			+0, +2, +0, +0,
			+0, +0, -2, +0,
			-5, -9, -13, 1}
	if !m.Aeq(want) {
		t.Errorf(format, m.Dump(), want.Dump())
	}
}

func TestPerspectiveInv(t *testing.T) {
	m := &M4{}
	p := NewM4().Persp(45, 800.0/600.0, 0.1, 50)
	ip := NewM4().PerspInv(45, 800.0/600.0, 0.1, 50)
	if !m.Mult(p, ip).Aeq(M4I) {
		t.Errorf(format, m.Dump(), M4I.Dump())
	}
}

// unit tests
// ============================================================================
// benchmarking.

// Check the time is saved by using the reference identity matrix instead of
// creating a new one. Run 'go test -bench=".*"' to get something like:
//     BenchmarkRefMI 2000000000	 0.72 ns/op
//     BenchmarkNewMI	20000000    95.9  ns/op
func BenchmarkRefMI(b *testing.B) {
	var m *M4
	for cnt := 0; cnt < b.N; cnt++ {
		m = M4I
	}
	m.Xx = 0 // make the compiler happy.
}
func BenchmarkNewMI(b *testing.B) {
	var m *M4
	for cnt := 0; cnt < b.N; cnt++ {
		m = &M4{Xx: 1, Yy: 1, Zz: 1, Ww: 1}
	}
	m.Xx = 0 // make the compiler happy.
}
