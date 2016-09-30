// Copyright © 2015-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"testing"

	"github.com/gazed/vu/math/lin"
)

// TestChildWorldTransform checks that a child object gets its
// world space location.
func TestChildWorldTransform(t *testing.T) {
	eng := newEngine(nil)
	parent := eng.Root().NewPov().SetAt(0, -8, 0).SetScale(4, 4, 4)
	parent.Spin(-90, 0, 0)
	child := parent.NewPov().SetAt(0, 0.78, 0.01).SetScale(0.1, 0.1, 0.1)
	eng.povs.updateWorldTransforms()
	if x, y, z := child.World(); !lin.Aeq(x, 0) || !lin.Aeq(y, -7.96) || !lin.Aeq(z, -3.12) {
		t.Errorf("Expecting 0.0, -7.96, -3.12, got %f, %f %f", x, y, z)
	}
	if eng != nil {
		eng.Shutdown()
	}
}

// Looked at reducing the API footprints using functional options:
//   http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
//   https://commandcenter.blogspot.ca/2014/01/self-referential-functions-and-design.html
// However it had a noticable affect on performance as it was roughly 100 times more
// expensive to make state changes. Currently functional options are used for
// the Engine but not for Pov, Camera, or high usage Model attributes.
//     BenchmarkSetFunction-8      20000000            90.1 ns/op
//     BenchmarkSetDirect-8        2000000000           1.64 ns/op
//
// NOTE: historical. Supporting code has been deleted for this benchmark.
//       No longer implemented on Pov or other state crucial classes.
//       Kept as a design note and reminder.
// func BenchmarkSetFunction(b *testing.B) {
//    eng := newEngine(nil)
//    p := eng.Root().NewPov()
//    for cnt := 0; cnt < b.N; cnt++ {
//        p.Set(At(0, -8, 0), Scale(4, 4, 4)) // functional options.
//    }
// }
//
func BenchmarkSetDirect(b *testing.B) {
	eng := newEngine(nil)
	p := eng.Root().NewPov()
	for cnt := 0; cnt < b.N; cnt++ {
		p.SetAt(0, -8, 0)   // options using method calls.
		p.SetScale(4, 4, 4) // options using method calls.
	}
}
