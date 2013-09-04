// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package physics

import (
	"fmt"
	"testing"
	"vu/math/lin"
)

func TestNewMotion(t *testing.T) {
	m := newMotion(10, 20, &lin.V3{}, lin.QIdentity())
	got, want := m.dump(), ""+
		"mass            : 20.00\n"+
		"location        : &{0.00 0.00 0.00}\n"+
		"linear momentum : &{0.00 0.00 0.00}\n"+
		"linear velocity : &{0.00 0.00 0.00}\n"+
		"direction       : &{0.00 0.00 0.00 1.00}\n"+
		"angular momentum: &{0.00 0.00 0.00}\n"+
		"angular velocity: &{0.00 0.00 0.00}\n"
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestClone(t *testing.T) {
	mo := newMotion(1, 2, &lin.V3{}, lin.QIdentity())
	mo.setLocation(&lin.V3{0, 1, 2})
	mo.setRotation(&lin.Q{3, 4, 5, 1})
	mo.setLinearMomentum(&lin.V3{1, 2, 3})
	mo.setAngularMomentum(&lin.V3{4, 5, 6})
	m2 := mo.clone()
	got, want := m2.dump(), mo.dump()
	if got != want {
		t.Errorf(format, got, want)
	}
	if &mo == &m2 {
		t.Error("mo and m2 should be distinct")
	}
}

// Dump is a testing utility method.
func (mo *motion) dump() string {
	str := fmt.Sprintf("mass            : %2.2f\n", mo.mass)
	str += fmt.Sprintf("location        : %2.2f\n", mo.loc)
	str += fmt.Sprintf("linear momentum : %2.2f\n", mo.linm)
	str += fmt.Sprintf("linear velocity : %2.2f\n", mo.linv)
	str += fmt.Sprintf("direction       : %2.2f\n", mo.dir)
	str += fmt.Sprintf("angular momentum: %2.2f\n", mo.angm)
	str += fmt.Sprintf("angular velocity: %2.2f\n", mo.angv)
	return str
}
