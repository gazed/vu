// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package lin

import (
	"fmt"
	"testing"
)

// While the functions being tested are not complicated, they are foundational in that many
// other libraries depend on them.  As such they each need a test.

func TestCloneQ(t *testing.T) {
	q := &Q{1, 2, 3, 4}
	q2 := q.Clone()
	got, want := q.Dump(), "{1.0 2.0 3.0 4.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
	if &q == &q2 {
		t.Error("q and q2 should be distinct")
	}
}

func TestQIdentity(t *testing.T) {
	q := QIdentity()
	got, want := q.Dump(), "{0.0 0.0 0.0 1.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestNormalizeQ(t *testing.T) {
	q := Q{1, 2, 3, 4}
	q.Unit()
	got, want := q.Dump(), "{0.2 0.4 0.5 0.7}"
	if got != want {
		t.Errorf(format, got, want)
	}
	q = Q{0, 0, 0, 1}
	q.Unit()
	got, want = q.Dump(), "{0.0 0.0 0.0 1.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
	q = Q{0, 0, 0, 0}
	q.Unit()
	got, want = q.Dump(), "{0.0 0.0 0.0 0.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestInverseQ(t *testing.T) {
	q := Q{0.2, 0.4, 0.5, 0.7}
	q = *q.Inverse()
	got, want := q.Dump(), "{-0.2 -0.4 -0.5 0.7}"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestAddQ(t *testing.T) {
	q := Q{0.2, 0.4, 0.5, 0.7}
	q.Add(&Q{1, 2, 3, 4})
	got, want := q.Dump(), "{1.2 2.4 3.5 4.7}"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestScaleQ(t *testing.T) {
	q := Q{0.2, 0.4, 0.5, 0.7}
	q.Scale(2)
	got, want := q.Dump(), "{0.4 0.8 1.0 1.4}"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestMultiplyQ(t *testing.T) {
	q1, q2 := (&Q{0, 1, 0, 2}).Unit(), (&Q{1, 0, 0, 2}).Unit()
	q1.Mult(q2)
	got, want := q1.Dump(), "{0.4 0.4 -0.2 0.8}"
	if got != want {
		t.Errorf(format, got, want)
	}

	// show that q1*q2 does not equal q2*q1
	q1, q2 = (&Q{0, 1, 0, 2}).Unit(), (&Q{1, 0, 0, 2}).Unit()
	q2.Mult(q1)
	got, want = q2.Dump(), "{0.4 0.4 0.2 0.8}"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestQAxisAngleQ(t *testing.T) {
	q := QAxisAngle(&V3{1, 1, 1}, 90)
	got, want := q.Dump(), "{0.4 0.4 0.4 0.7}"
	if got != want {
		t.Errorf(format, got, want)
	}

	q = QAxisAngle(&V3{0, 0, -1}, 0)
	got, want = q.Dump(), "{0.0 0.0 -0.0 1.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestAxisAngleQ(t *testing.T) {
	q := Q{0.4, 0.4, 0.4, 0.7}
	axis, angle := q.AxisAngle()
	got, want := axis.Dump(), "{0.6 0.6 0.6}"
	if got != want {
		t.Errorf(format, got, want)
	}
	got, want = fmt.Sprintf("%2.1f", angle), "91.1"
	if got != want {
		t.Errorf(format, got, want)
	}

	// test the identity quaternion
	q = Q{0, 0, 0, 1}
	axis, angle = q.AxisAngle()
	got, want = axis.Dump(), "{0.0 0.0 -1.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
	got, want = fmt.Sprintf("%2.1f", angle), "0.0"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestNLerpQ(t *testing.T) {
	q := (&Q{0.2, 0.4, 0.5, 0.7}).Nlerp(&Q{8, 2, 6, 10}, 0.5)
	got, want := q.Dump(), "{0.5 0.2 0.4 0.7}"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestM4Q(t *testing.T) {
	q := Q{0.2, 0.4, 0.5, 0.7}
	m := q.M4()
	expect := M4{
		+0.18, +0.86, -0.36, +0.00,
		-0.54, +0.42, +0.68, +0.00,
		+0.76, +0.12, +0.60, +0.00,
		+0.00, +0.00, +0.00, +1.00}
	got, want := m.Dump(), expect.Dump()
	if got != want {
		t.Errorf(format, got, want)
	}

	// check identity quaternion
	q = Q{0, 0, 0, 1}
	m = q.M4()
	expect = M4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
	got, want = m.Dump(), expect.Dump()
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestMovementAroundY(t *testing.T) {
	loc := &V3{}          // at origin
	dir := &Q{0, 0, 0, 1} // no rotation

	// move 1 unit along X
	loc.Add((&V3{1, 0, 0}).MultQ(dir))
	got, want := loc.Dump(), "{1.0 0.0 0.0}"
	if got != want {
		t.Errorf(format, got, want)
	}

	// rotate 90 degrees around Y and move 1 unit along X (should be moving along -Z)
	//dir = QAxisAngle(&V3{0, 1, 0}, 90).Mult(dir)
	dir = dir.Mult(QAxisAngle(&V3{0, 1, 0}, 90))
	loc.Add((&V3{1, 0, 0}).MultQ(dir))
	got, want = loc.Dump(), "{1.0 0.0 -1.0}"
	if got != want {
		t.Errorf(format, got, want)
	}

	// rotate -90 degrees around Y and move 1 unit along X (should be moving along X again)
	dir = QAxisAngle(&V3{0, 1, 0}, -90).Mult(dir)
	loc.Add((&V3{1, 0, 0}).MultQ(dir))
	got, want = loc.Dump(), "{2.0 0.0 -1.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestMovementAroundX(t *testing.T) {
	loc := &V3{}          // at origin
	dir := &Q{0, 0, 0, 1} // no rotation

	// move 1 unit along Z
	loc.Add((&V3{0, 0, 1}).MultQ(dir))
	got, want := loc.Dump(), "{0.0 0.0 1.0}"
	if got != want {
		t.Errorf(format, got, want)
	}

	// rotate 90 degrees around X and move 1 unit along Z (should be moving along -Y)
	dir = QAxisAngle(&V3{1, 0, 0}, 90).Mult(dir)
	loc.Add((&V3{0, 0, 1}).MultQ(dir))
	got, want = loc.Dump(), "{0.0 -1.0 1.0}"
	if got != want {
		t.Errorf(format, got, want)
	}

	// rotate -90 degrees around X and move 1 unit along Z (should be moving along Z again)
	dir = QAxisAngle(&V3{1, 0, 0}, -90).Mult(dir)
	loc.Add((&V3{0, 0, 1}).MultQ(dir))
	got, want = loc.Dump(), "{0.0 -1.0 2.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
}
