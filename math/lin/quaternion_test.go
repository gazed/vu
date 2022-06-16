// Copyright © 2013-2024 Galvanized Logic Inc.

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

func TestDefaultAxisAngle(t *testing.T) {
	q, v, angle, want := &Q{0, 0, 0, 1}, &V3{}, 0.0, &V3{1, 0, 0}
	if v.X, v.Y, v.Z, angle = q.Aa(); !v.Aeq(want) || !Aeq(Deg(angle), 0) {
		t.Errorf("Got axis %s and angle %+2.7f", v.Dump(), Deg(angle))
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

func TestAa(t *testing.T) {
	for deg := 0; deg <= 360; deg++ {
		q := NewQ().SetAa(0, 1, 0, Rad(float64(deg)))
		if !Aeq(q.Len(), 1) {
			t.Errorf("Need unit quat for angle %d", deg)
		}
		x, y, z, rot := q.Aa()
		switch {
		case deg >= 0 && deg < 180:
			if !Aeq(float64(deg), Deg(rot)) { // positive axis.
				t.Errorf("Wanted %d got %+2.5f : %f %f %f %f", deg, Deg(rot), x, y, z, q.W)
			}
		case deg >= 180 && deg <= 360: // negative axis.
			if !Aeq(float64(deg), 360-Deg(rot)) {
				t.Errorf("Wanted %d got %+2.5f : %f %f %f %f", deg, Deg(rot), x, y, z, q.W)
			}
		}
	}
}

func TestSetRotationM(t *testing.T) {
	q2 := &Q{}
	m := NewM3()
	for deg := 0; deg < 360; deg++ {
		q := NewQ().SetAa(0, 1, 0, Rad(float64(deg)))
		m.SetQ(q)
		q2.SetM3(m)
		if !q.Aeq(q2) {
			t.Errorf("SetM deg %d : %s : %s", deg, q.Dump(), q2.Dump())
		}
	}
}

// =============================================================================

// Check when/if the rotation matrix to quaternion is changed.
// BenchmarkSetM3-8     200000000       8.22 ns/op
func BenchmarkSetM3(b *testing.B) {
	q, m := &Q{}, (&M3{}).SetAa(0, 1, 0, Rad(45))
	for cnt := 0; cnt < b.N; cnt++ {
		q.SetM3(m)
	}
}

// BenchmarkSetAa-8   	50000000	        28.4 ns/op
func BenchmarkSetAa(b *testing.B) {
	q := &Q{}
	for cnt := 0; cnt < b.N; cnt++ {
		q.SetAa(0, 1, 0, Rad(45))
	}
}
