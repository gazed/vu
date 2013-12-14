// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package lin

import (
	"testing"
)

// While the functions being tested are not complicated, they are foundational in that many
// other libraries depend on them. As such they each need a test. Where applicable, tests
// check that the output quaternion can also be used as the input quaternion.

func TestAddQ(t *testing.T) {
	q, want := &Q{1, 2, 3, 4}, &Q{2, 4, 6, 8}
	if !q.Add(q, q).Eq(want) {
		t.Errorf(format, q.Dump(), want.Dump())
	}
}

func TestSubtractQ(t *testing.T) {
	q, want := &Q{1, 2, 3, 4}, &Q{0, 0, 0, 0}
	if !q.Sub(q, q).Eq(want) {
		t.Errorf(format, q.Dump(), want.Dump())
	}
}

func TestInverseQ(t *testing.T) {
	q, qi, want := &Q{0.2, 0.4, 0.5, 0.7}, &Q{}, &Q{-0.2, -0.4, -0.5, 0.7}
	if !qi.Inv(q).Eq(want) {
		t.Errorf(format, q.Dump(), want.Dump())
	}
	if !q.Mult(q, qi).Unit().Aeq(QI) {
		t.Errorf(format, q.Dump(), QI.Dump())
	}
}

func TestNormalizeQ(t *testing.T) {
	q, want := &Q{1, 2, 3, 4}, &Q{0.1825742, 0.3651484, 0.5477226, 0.7302967}
	if !q.Unit().Aeq(want) {
		t.Errorf(format, q.Dump(), want.Dump())
	}
	q, want = &Q{0, 0, 0, 1}, &Q{0, 0, 0, 1}
	if !q.Unit().Eq(want) {
		t.Errorf(format, q.Dump(), want.Dump())
	}
	q, want = &Q{0, 0, 0, 0}, &Q{0, 0, 0, 0}
	if !q.Unit().Eq(want) {
		t.Errorf(format, q.Dump(), want.Dump())
	}
}

func TestScaleQ(t *testing.T) {
	q, want := &Q{1, 2, 3, 4}, &Q{2, 4, 6, 8}
	if !q.Scale(2).Eq(want) {
		t.Errorf(format, q.Dump(), want.Dump())
	}
}

func TestInverseScaleQ(t *testing.T) {
	q, want := &Q{1, 2, 3, 4}, &Q{0.5, 1, 1.5, 2}
	if !q.Div(2).Eq(want) {
		t.Errorf(format, q.Dump(), want.Dump())
	}
}

func TestMultiplyQ(t *testing.T) {
	q, want := &Q{1, 2, 3, 4}, &Q{8, 16, 24, 2}
	if !q.Mult(q, q).Eq(want) {
		t.Errorf(format, q.Dump(), want.Dump())
	}
}

func TestDotQ(t *testing.T) {
	q := &Q{0.1825742, 0.3651484, 0.5477226, 0.7302967}
	if !Aeq(q.Len(), 1) {
		t.Errorf("Dot is not %+2.8f", q.Dot(q))
	}
}

func TestLenQ(t *testing.T) {
	q := &Q{0.1825742, 0.3651484, 0.5477226, 0.7302967}
	if !Aeq(q.Len(), 1) {
		t.Errorf("Len is %+2.7f", q.Len())
	}
}

func TestAngQ(t *testing.T) {
	q, a := NewQ().SetAa(1, 0, 0, Rad(90)), NewQ().SetAa(1, 0, 0, Rad(135))
	angle := Deg(q.Ang(a))
	if !Aeq(angle, 45) {
		t.Errorf("Angle is %+2.7f", angle)
	}
	angle = Deg(q.Ang(q)) // angle between the same.
	if !Aeq(angle, 0) {
		t.Errorf("Angle is %+2.7f", angle)
	}
}

func TestNLerpQ(t *testing.T) {
	q, b, want := (&Q{1, 2, 3, 4}).Unit(), (&Q{8, 2, 6, 10}).Unit(), &Q{0.38151321, 0.25950587, 0.49715611, 0.73480630}
	if !q.Nlerp(q, b, 0.5).Aeq(want) {
		t.Errorf(format, q.Dump(), want.Dump())
	}
	if !Aeq(q.Len(), 1) {
		t.Errorf("Nlerp result should be unit length")
	}
}

func TestMultiplyQV(t *testing.T) {
	q, v, want := &Q{1, 2, 3, 4}, &V3{1, 2, 3}, &Q{4, 8, 12, -14}
	if !q.MultQV(q, v).Eq(want) {
		t.Errorf(format, q.Dump(), want.Dump())
	}
}

func TestGetAxisAngle(t *testing.T) {
	q, v, angle, want := &Q{0.40824829, 0.40824829, 0.40824829, 0.707106781}, &V3{}, 0.0, NewV3().SetS(1, 1, 1).Unit()
	if v.X, v.Y, v.Z, angle = q.Aa(); !v.Aeq(want) || !Aeq(Deg(angle), 90) {
		t.Errorf("Got axis %s and angle %+2.7f", v.Dump(), Deg(angle))
	}
}
func TestDefaultAxisAngle(t *testing.T) {
	q, v, angle, want := &Q{0, 0, 0, 1}, &V3{}, 0.0, &V3{1, 0, 0}
	if v.X, v.Y, v.Z, angle = q.Aa(); !v.Aeq(want) || !Aeq(Deg(angle), 0) {
		t.Errorf("Got axis %s and angle %+2.7f", v.Dump(), Deg(angle))
	}
}

func TestSetAxisAngleQ(t *testing.T) {
	q, want := &Q{}, &Q{0.40824829, 0.40824829, 0.40824829, 0.707106781}
	if !q.SetAa(1, 1, 1, Rad(90)).Aeq(want) {
		t.Errorf(format, q.Dump(), want.Dump())
	}
}

func TestSetRotationM(t *testing.T) {
	q := NewQ().SetAa(1, 1, 1, Rad(90))
	m := NewM3().SetQ(q)
	want := &Q{0.40824831, 0.40824831, 0.40824831, 0.70710677}
	if !q.SetM(m).Aeq(want) {
		t.Errorf(format, q.Dump(), want.Dump())
	}
}

// Rotation 45deg * transform of another rotation of 45deg should give 90deg rotation.
func TestMultTransformQ(t *testing.T) {
	q, want := NewQ().SetAa(1, 1, 1, Rad(45)), &Q{0.40824831, 0.40824831, 0.40824831, 0.70710677}
	transform := NewT().SetLoc(10, 0, 0).SetAa(1, 1, 1, Rad(45))
	if !q.MultT(transform).Aeq(want) {
		t.Errorf(format, q.Dump(), want.Dump())
	}
}
