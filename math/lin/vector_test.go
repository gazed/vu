// SPDX-FileCopyrightText : © 2014-2022 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

package lin

import (
	"testing"
)

// While the functions below are not complicated, they are foundational such that it is
// better to test each one of them then have the bugs discovered later from other code.
// Where applicable, check that the output vector can also be used as one or both
// of the input vectors.

func TestSetV3(t *testing.T) {
	v, a := &V3{}, &V3{1, 2, 3}
	if !v.Set(a).Eq(a) {
		t.Errorf("%s is not the same as %s", v.Dump(), a.Dump())
	}
}
func TestSetV4(t *testing.T) {
	v, a := &V4{}, &V4{1, 2, 3, 4}
	if !v.Set(a).Eq(a) {
		t.Errorf("%s is not the same as %s", v.Dump(), a.Dump())
	}
}

func TestSwapV3(t *testing.T) {
	v, a, vo, ao := &V3{}, &V3{1, 2, 3}, &V3{}, &V3{1, 2, 3}
	v.Swap(a)
	if !v.Eq(ao) || !a.Eq(vo) {
		t.Errorf("%s did not swap with %s", v.Dump(), a.Dump())
	}
}
func TestSwapV4(t *testing.T) {
	v, a, vo, ao := &V4{}, &V4{1, 2, 3, 4}, &V4{}, &V4{1, 2, 3, 4}
	v.Swap(a)
	if !v.Eq(ao) || !a.Eq(vo) {
		t.Errorf("%s did not swap with %s", v.Dump(), a.Dump())
	}
}

func TestMinimumV3(t *testing.T) {
	v, a, want := &V3{1, -2, 3}, &V3{-1, 2, -3}, &V3{-1, -2, -3}
	if !v.Min(v, a).Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}
func TestMinimumV4(t *testing.T) {
	v, a, want := &V4{1, -2, 3, -4}, &V4{-1, 2, -3, 4}, &V4{-1, -2, -3, -4}
	if !v.Min(v, a).Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}

func TestMaxiumumV3(t *testing.T) {
	v, a, want := &V3{1, -2, 3}, &V3{-1, 2, -3}, &V3{1, 2, 3}
	if !v.Max(v, a).Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}
func TestMaxiumumV4(t *testing.T) {
	v, a, want := &V4{1, -2, 3, -4}, &V4{-1, 2, -3, 4}, &V4{1, 2, 3, 4}
	if !v.Max(v, a).Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}

func TestAddV3(t *testing.T) {
	v, want := &V3{1, 2, 3}, &V3{2, 4, 6}
	if !v.Add(v, v).Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}
func TestAddV4(t *testing.T) {
	v, want := &V4{1, 2, 3, 4}, &V4{2, 4, 6, 8}
	if !v.Add(v, v).Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}

func TestSubtractV3(t *testing.T) {
	v, want := &V3{1, 2, 3}, &V3{0, 0, 0}
	if !v.Sub(v, v).Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}
func TestSubtractV4(t *testing.T) {
	v, want := &V4{1, 2, 3, 4}, &V4{0, 0, 0, 0}
	if !v.Sub(v, v).Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}

func TestMultiplyV3(t *testing.T) {
	v, want := &V3{1, 2, 3}, &V3{1, 4, 9}
	if !v.Mult(v, v).Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}
func TestMultiplyV4(t *testing.T) {
	v, want := &V4{1, 2, 3, 4}, &V4{1, 4, 9, 16}
	if !v.Mult(v, v).Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}
func TestMultiplyV3Q(t *testing.T) {
	v, q, want := &V3{1, 2, 3}, &Q{0, 0, 0, 1}, &V3{1, 2, 3}
	if !v.MultQ(v, q).Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
	v, q, want = &V3{1, 0, 0}, NewQ().SetAa(0, 0, 1, Rad(90)).Unit(), &V3{0, 1, 0}
	if !v.MultQ(v, q).Aeq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
	v, q, want = &V3{10, 10, 0}, NewQ().SetAa(1, 0, 0, Rad(45)).Unit(), &V3{10, 7.071067812, 7.071067812}
	if !v.MultQ(v, q).Aeq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}

func TestScaleV3(t *testing.T) {
	v, want := &V3{1, 2, 3}, &V3{2, 4, 6}
	if !v.Scale(v, 2).Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}
func TestScaleV4(t *testing.T) {
	v, want := &V4{1, 2, 3, 4}, &V4{2, 4, 6, 8}
	if !v.Scale(v, 2).Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}

func TestInverseScaleV3(t *testing.T) {
	v, want := &V3{1, 2, 3}, &V3{2, 4, 6}
	if !v.Div(0.5).Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}
func TestInverseScaleV4(t *testing.T) {
	v, want := &V4{1, 2, 3, 4}, &V4{2, 4, 6, 8}
	if !v.Div(0.5).Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}

func TestDotV3(t *testing.T) {
	v, a := &V3{1, 2, 3}, &V3{2, 4, 8}
	if v.Dot(a) != 34 || v.Dot(v) != 14 {
		t.Error("Invalid dot product")
	}
}
func TestDotV4(t *testing.T) {
	v, a := &V4{1, 2, 4, 3}, &V4{2, 4, 8, 9}
	if v.Dot(a) != 69 || v.Dot(v) != 30 {
		t.Error("Invalid dot product")
	}
}

func TestLengthV3(t *testing.T) {
	v := &V3{9, 2, 6}
	if v.Len() != 11 {
		t.Error("Invalid length", v.Len())
	}
}
func TestLengthV4(t *testing.T) {
	v := &V4{6, 6, 6, 6}
	if v.Len() != 12 {
		t.Error("Invalid length", v.Len())
	}
}

func TestDistanceV3(t *testing.T) {
	v, a := &V3{9, 2, 6}, &V3{18, 4, 12}
	if v.Dist(a) != 11 {
		t.Errorf("Invalid distance %f", v.Dist(a))
	}
	if v.Dist(v) != 0 {
		t.Error("Distance with self should be zero.")
	}
}

func TestAngleV3(t *testing.T) {
	v, a, ang := &V3{1, 0, 0}, &V3{0, 1, 0}, 90.0
	if Deg(v.Ang(a)) != ang {
		t.Errorf("Wanted angle %f got  %f", ang, Deg(v.Ang(a)))
	}
}

func TestNormalizeV3(t *testing.T) {
	v, want := &V3{0, 0, 0}, &V3{0, 0, 0}
	if !v.Unit().Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
	v = &V3{5, 6, 7}
	if !Aeq(v.Unit().Len(), 1) {
		t.Errorf("Normalized vectors should have length one")
	}
}
func TestNormalizeV4(t *testing.T) {
	v := &V4{5, 6, 7, 8}
	if !Aeq(v.Unit().Len(), 1) {
		t.Errorf("Normalized vectors should have length one")
	}
}

func TestCrossV3(t *testing.T) {
	v, b, want := &V3{3, -3, 1}, &V3{4, 9, 2}, &V3{-15, -2, 39}
	if !v.Cross(v, b).Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}

func TestLerpV3(t *testing.T) {
	v, b, want := &V3{1, 2, 3}, &V3{5, 6, 7}, &V3{3, 4, 5}
	if !v.Lerp(v, b, 0.5).Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}
func TestLerpV4(t *testing.T) {
	v, b, want := &V4{1, 2, 3, 4}, &V4{5, 6, 7, 8}, &V4{3, 4, 5, 6}
	if !v.Lerp(v, b, 0.5).Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}

func TestPlane(t *testing.T) {
	v, p, q, wantp, wantq := &V3{1, 0, 0}, &V3{}, &V3{}, &V3{0, 1, 0}, &V3{0, 0, 1}
	if v.Plane(p, q); !p.Eq(wantp) || !q.Eq(wantq) {
		t.Errorf("Did not get expected plane vectors for %s", v.Dump())
	}
	v, wantp, wantq = &V3{0, 1, 0}, &V3{-1, 0, 0}, &V3{0, 0, 1}
	if v.Plane(p, q); !p.Eq(wantp) || !q.Eq(wantq) {
		t.Errorf("Did not get expected plane vectors for %s", v.Dump())
	}
	v, wantp, wantq = &V3{0, 0, 1}, &V3{0, -1, 0}, &V3{1, 0, 0}
	if v.Plane(p, q); !p.Eq(wantp) || !q.Eq(wantq) {
		t.Errorf("Did not get expected plane vectors for %s", v.Dump())
	}
}

func TestMultvMV3(t *testing.T) {
	v, m, want := &V3{1, 2, 3},
		&M3{1, 2, 3,
			1, 2, 3,
			1, 2, 3}, &V3{6, 12, 18}
	if !v.MultvM(v, m).Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}
func TestMultvMV4(t *testing.T) {
	v, m, want := &V4{1, 2, 3, 4},
		&M4{1, 2, 3, 4,
			1, 2, 3, 4,
			1, 2, 3, 4,
			1, 2, 3, 4}, &V4{10, 20, 30, 40}
	if !v.MultvM(v, m).Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}

func TestMultMvV3(t *testing.T) {
	v, want, m := &V3{1, 2, 3}, &V3{14, 14, 14},
		&M3{1, 2, 3,
			1, 2, 3,
			1, 2, 3}
	if !v.MultMv(m, v).Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}
func TestMultMvV4(t *testing.T) {
	v, want, m := &V4{1, 2, 3, 4}, &V4{30, 30, 30, 30},
		&M4{1, 2, 3, 4,
			1, 2, 3, 4,
			1, 2, 3, 4,
			1, 2, 3, 4}
	if !v.MultMv(m, v).Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}

func TestCascade(t *testing.T) {
	v, v1, want := &V3{1, 2, 3}, &V3{10, 20, 30}, &V3{-10, -40, -90}
	v.Mult(v, v1).Neg(v)
	if !v.Eq(want) {
		t.Errorf(format, v.Dump(), want.Dump())
	}
}

func TestDirection(t *testing.T) {
	t.Run("right", func(t *testing.T) {
		q := NewQ().SetAa(1, 1, 1, 90) // rotation
		vr0 := NewV3().SetS(1, 0, 0)   // X direction
		vr0.MultQ(vr0, q)

		// same as
		vr1 := NewV3().Right(q)
		if !vr0.Aeq(vr1) {
			t.Errorf("expected same result\n %+v\n %+v", vr0, vr1)
		}
	})

	t.Run("right invert", func(t *testing.T) {
		q := NewQ().SetAa(1, 1, 1, -90) // opposite rotation
		vr0 := NewV3().SetS(1, 0, 0)    // X direction
		vr0.MultQ(vr0, q)

		// same as
		q = NewQ().SetAa(1, 1, 1, 90)
		vr1 := NewV3().RightInverted(q)
		if !vr0.Aeq(vr1) {
			t.Errorf("expected same result\n %+v\n %+v", vr0, vr1)
		}
	})

	t.Run("up", func(t *testing.T) {
		q := NewQ().SetAa(1, 1, 1, 90) // rotation
		vr0 := NewV3().SetS(0, 1, 0)   // Y direction
		vr0.MultQ(vr0, q)

		// same as
		vr1 := NewV3().Up(q)
		if !vr0.Aeq(vr0) {
			t.Errorf("expected same result\n %+v\n %+v", vr0, vr1)
		}
	})

	t.Run("up invert", func(t *testing.T) {
		q := NewQ().SetAa(1, 1, 1, -90) // opposite rotation
		vr0 := NewV3().SetS(0, 1, 0)    // Y direction
		vr0.MultQ(vr0, q)

		// same as
		q = NewQ().SetAa(1, 1, 1, 90)
		vr1 := NewV3().UpInverted(q)
		if !vr0.Aeq(vr1) {
			t.Errorf("expected same result\n %+v\n %+v", vr0, vr1)
		}
	})

	t.Run("forward", func(t *testing.T) {
		q := NewQ().SetAa(1, 1, 1, 90) // rotation
		vr0 := NewV3().SetS(0, 0, 1)   // Z direction
		vr0.MultQ(vr0, q)

		// same as
		vr1 := NewV3().Forward(q)
		if !vr0.Aeq(vr1) {
			t.Errorf("expected same result\n %+v\n %+v", vr0, vr1)
		}
	})

	t.Run("forward invert", func(t *testing.T) {
		q := NewQ().SetAa(1, 1, 1, -90) // opposite rotation
		vr0 := NewV3().SetS(0, 0, 1)    // Z direction
		vr0.MultQ(vr0, q)

		// same as
		q = NewQ().SetAa(1, 1, 1, 90)
		vr1 := NewV3().ForwardInverted(q)
		if !vr0.Aeq(vr1) {
			t.Errorf("expected same result\n %+v\n %+v", vr0, vr1)
		}
	})
}

// unit tests
// ============================================================================
// benchmarking.

// Check golang efficiency for different method signatures and heap/stack
// memory allocation. Run go test -bench=".*Sub*" to get something like:
//     BenchmarkV3Sub	    1000000000	    2.51 ns/op
//     BenchmarkV3SubNew	  50000000	   68.1  ns/op
//     BenchmarkV3SubScalar	 500000000	    3.58 ns/op
//     BenchmarkV3SubNoCall	2000000000	    1.43 ns/op

func BenchmarkV3Sub(b *testing.B) {
	v, a, o := &V3{}, &V3{2, 2, 2}, &V3{1, 1, 1}
	for cnt := 0; cnt < b.N; cnt++ {
		v = v.Sub(a, o)
	}
}
func BenchmarkV3SubNew(b *testing.B) {
	var v *V3
	a, o := &V3{2, 2, 2}, &V3{1, 1, 1}
	for cnt := 0; cnt < b.N; cnt++ {
		v = a.subNew(o)
	}
	v.X = 0 // Otherwise compiler complains about unused variables.
}
func BenchmarkV3SubScalar(b *testing.B) {
	var x, y, z float64
	for cnt := 0; cnt < b.N; cnt++ {
		x, y, z = subScalars(2, 2, 2, 1, 1, 1)
	}
	if x == 1 && y == 1 && z == 1 {
		// Otherwise compiler complains about unused variables.
	}
}
func BenchmarkV3SubNoCall(b *testing.B) {
	var x, y, z float64
	for cnt := 0; cnt < b.N; cnt++ {
		x, y, z = 2-1, 2-1, 2-1
	}
	if x == 1 && y == 1 && z == 1 {
		// Otherwise compiler complains about unused variables.
	}
}

// subNew creates a new V3 that contains the subtraction of vector b from a.
// Used to benchmark how struct allocation affects execution time.
func (a *V3) subNew(b *V3) *V3 { return &V3{a.X - b.X, a.Y - b.Y, a.Z - b.Z} }

// subScalars just does the basic subtraction. Used to benchmark how passing
// lots of float parameter affects execution time.
func subScalars(ax, ay, az, bx, by, bz float64) (x, y, z float64) { return ax - bx, ay - by, az - bx }
