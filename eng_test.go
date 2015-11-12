// Copyright © 2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"testing"

	"github.com/gazed/vu/math/lin"
)

// TestChildWorldTransform checks that a child object can calculate its
// world space location.
func TestChildWorldTransform(t *testing.T) {
	eng := newEngine(nil)
	parent := eng.Root().NewPov().SetLocation(0, -8, 0).SetScale(4, 4, 4)
	parent.Spin(-90, 0, 0)
	child := parent.NewPov().SetLocation(0, 0.78, 0.01).SetScale(0.1, 0.1, 0.1)

	// call placeModels to initialize the model transform matrix needed by World.
	eng.placeModels(eng.root(), lin.M4I) // update all transforms.
	if x, y, z := child.World(); !lin.Aeq(x, 0) || !lin.Aeq(y, -7.96) || !lin.Aeq(z, -3.12) {
		t.Errorf("Expecting %f %f %f, got %f, %f %f", 0.0, -7.96, -3.12, x, y, z)
	}
	if eng != nil {
		eng.Shutdown()
	}
}
