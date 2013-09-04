// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package lin

import (
	"fmt"
	"testing"
)

// While the functions below are not complicated, they are foundational such that it is
// better to test each one of them then have the bugs discovered later from other code.

func TestClone3v(t *testing.T) {
	v := V3{1, 2, 3}
	v2 := v.Clone()
	got, want := v2.Dump(), "{1.0 2.0 3.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestClone4v(t *testing.T) {
	v := V4{1, 2, 3, 4}
	v2 := v.Clone()
	got, want := v2.Dump(), "{1.0 2.0 3.0 4.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestAdd3v(t *testing.T) {
	v := V3{1, 2, 3}
	v.Add(&V3{5, 6, 7})
	got, want := v.Dump(), "{6.0 8.0 10.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
}
func TestAdd4v(t *testing.T) {
	v := V4{1, 2, 3, 4}
	v.Add(&V4{5, 6, 7, 8})
	got, want := v.Dump(), "{6.0 8.0 10.0 12.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestSubtract3v(t *testing.T) {
	v := V3{1, 2, 3}
	v.Sub(&V3{5, 6, 7})
	got, want := v.Dump(), "{-4.0 -4.0 -4.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
}
func TestSubtract4v(t *testing.T) {
	v := V4{1, 2, 3, 4}
	v.Sub(&V4{5, 6, 7, 8})
	got, want := v.Dump(), "{-4.0 -4.0 -4.0 -4.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestMultiply3v(t *testing.T) {
	v := V3{1, 2, 3}
	v.Mult(&V3{5, 6, 7})
	got, want := v.Dump(), "{5.0 12.0 21.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
}
func TestMultiply4v(t *testing.T) {
	v := V4{1, 2, 3, 4}
	v.Mult(&V4{5, 6, 7, 8})
	got, want := v.Dump(), "{5.0 12.0 21.0 32.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestDot3v(t *testing.T) {
	lhs := V3{2, 4, 8}
	dot := lhs.Dot(&V3{1, 2, 4})
	got, want := fmt.Sprintf("%2.1f", dot), "42.0"
	if got != want {
		t.Errorf(format, got, want)
	}
}
func TestDot4v(t *testing.T) {
	lhs := V4{2, 4, 8, 9}
	dot := lhs.Dot(&V4{1, 2, 4, 3})
	got, want := fmt.Sprintf("%2.1f", dot), "69.0"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestLength3v(t *testing.T) {
	v := V3{2, 4, 8}
	length := v.Len()
	got, want := fmt.Sprintf("%2.1f", length), "9.2"
	if got != want {
		t.Errorf(format, got, want)
	}
}
func TestLength4v(t *testing.T) {
	v := V4{2, 4, 8, 9}
	length := v.Len()
	got, want := fmt.Sprintf("%2.1f", length), "12.8"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestDistance3v(t *testing.T) {
	lhs := V3{1, 1, 1}
	dist := lhs.Dist(&V3{-1, -1, -1})
	got, want := fmt.Sprintf("%2.1f", dist), "3.5"
	if got != want {
		t.Errorf(format, got, want)
	}
}
func TestDistance4v(t *testing.T) {
	lhs := V4{2, 4, 8, 9}
	dist := lhs.Dist(&V4{1, 2, 4, 3})
	got, want := fmt.Sprintf("%2.1f", dist), "7.5"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestNormalize3v(t *testing.T) {
	v := V3{5, 6, 7}
	v.Unit()
	got, want := v.Dump(), "{0.5 0.6 0.7}"
	if got != want {
		t.Errorf(format, got, want)
	}
	v = V3{0, 0, 0}
	v.Unit()
	got, want = v.Dump(), "{0.0 0.0 0.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
	v = V3{0, 0, 0.2}
	shouldBeOne := v.Unit().Len()
	if !IsOne(shouldBeOne) {
		t.Errorf(format, got, want)
	}
}
func TestNormalize4v(t *testing.T) {
	v := V4{5, 6, 7, 8}
	v.Unit()
	got, want := v.Dump(), "{0.4 0.5 0.5 0.6}"
	if got != want {
		t.Errorf(format, got, want)
	}
	v = V4{0, 0, 0, 0}
	v.Unit()
	got, want = v.Dump(), "{0.0 0.0 0.0 0.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestCross3v(t *testing.T) {
	v1 := &V3{3, -3, 1}
	crossProduct := v1.Cross(&V3{4, 9, 2})
	got, want := crossProduct.Dump(), "{-15.0 -2.0 39.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestMultLM3(t *testing.T) {
	v := &V3{1, 2, 3}
	m := &M3{
		1, 2, 3,
		1, 2, 3,
		1, 2, 3}
	v.MultL(m)
	got, want := v.Dump(), "{6.0 12.0 18.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestMultLM4(t *testing.T) {
	v := &V4{1, 2, 3, 4}
	m := &M4{
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4}
	v.MultL(m)
	got, want := v.Dump(), "{10.0 20.0 30.0 40.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestMultiplyV3Q(t *testing.T) {
	v := &V3{1, 2, 3}
	v.MultQ(&Q{0.2, 0.4, 0.5, 0.7})
	got, want := v.Dump(), "{1.3 1.8 2.4}"
	if got != want {
		t.Errorf(format, got, want)
	}

	// test a very small increment.
	v = &V3{0, 0, 0.01}
	v.MultQ(&Q{0.2, 0.4, 0.05, 0.7})
	got, want = v.Dump(), "{0.0 -0.0 -0.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
}
