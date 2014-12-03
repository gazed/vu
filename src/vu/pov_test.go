// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"testing"
	"vu/math/lin"
)

// Check that the inverse of a perspective view is correct.
func TestInverseVp(t *testing.T) {
	v := newCamera()
	v.Loc.SetS(10, 10, 10)
	v.Rot.SetAa(1, 0, 0, -lin.Rad(90))
	vm := vp(v, &lin.M4{})
	ivm := ivp(v, &lin.M4{})
	if !vm.Mult(vm, ivm).Aeq(lin.M4I) {
		t.Errorf("Matrix times inverse should be identity")
	}
}
