// SPDX-FileCopyrightText : © 2024 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

package physics

import (
	"fmt"
	"testing"

	"github.com/gazed/vu/math/lin"
)

// go test -run Hull
func TestConvexHull(t *testing.T) {
	b := NewBox(50, 10, 20, true)
	if len(b.colliders) != 1 {
		t.Fatal("expecting a single collider for a box")
	}
	c := &b.colliders[0]
	if c.ctype != collider_TYPE_CONVEX_HULL {
		t.Fatal("expecting convex hull collider")
	}
}

// check matrix conventions. The physics package uses row-major.
// This means translate*rotate*scale
func TestMatrixOrder(t *testing.T) {
	scale := &lin.M4{
		5.0, 0.0, 0.0, 0.0,
		0.0, 5.0, 0.0, 0.0,
		0.0, 0.0, 5.0, 0.0,
		0.0, 0.0, 0.0, 1.0}
	q := lin.NewQ().SetAa(0, 1, 0, lin.Rad(90)) // rotate 90 around Y
	rotate := lin.NewM4().SetQ(q)
	translate := &lin.M4{
		1.0, 0.0, 0.0, 1.0,
		0.0, 1.0, 0.0, 2.0,
		0.0, 0.0, 1.0, 3.0,
		0.0, 0.0, 0.0, 1.0}
	model := lin.NewM4()
	model.Mult(rotate, scale)
	model.Mult(translate, model)
	expect := &lin.M4{
		+0.0, +0.0, +5.0, +1.0,
		+0.0, +5.0, +0.0, +2.0,
		-5.0, +0.0, +0.0, +3.0,
		+0.0, +0.0, +0.0, +1.0}
	if !model.Aeq(expect) {
		t.Errorf("did not match expected\n%s\n", DumpM4(expect))
	}
}

func DumpM4(m *lin.M4) string {
	format := "  [%+2.9f, %+2.9f, %+2.9f, %+2.9f]\n"
	str := fmt.Sprintf(format, m.Xx, m.Xy, m.Xz, m.Xw)
	str += fmt.Sprintf(format, m.Yx, m.Yy, m.Yz, m.Yw)
	str += fmt.Sprintf(format, m.Zx, m.Zy, m.Zz, m.Zw)
	str += fmt.Sprintf(format, m.Wx, m.Wy, m.Wz, m.Ww)
	return str
}

// =============================================================================
// Benchmarks.

// go test -bench=.
//
// Check how all the lin.New* methods affects performance.
// A couple of results captured are:
//
//	cpu: 13th Gen Intel(R) Core(TM) i7-13700K
//	BenchmarkV3-24          1000000000               0.09693 ns/op
//	BenchmarkV3-24          1000000000               0.08739 ns/op
func BenchmarkV3(b *testing.B) {
	for cnt := 0; cnt < b.N; cnt++ {
		lin.NewV3().Add(
			lin.NewV3().SetS(1, 1, 1),
			lin.NewV3().Scale(
				lin.NewV3().SetS(1, 1, 1).MultMv(
					lin.NewM3I(),
					lin.NewV3().Sub(
						lin.NewV3().SetS(1, 1, 1),
						lin.NewV3().Cross(
							lin.NewV3().SetS(1, 1, 1),
							lin.NewV3().MultMv(lin.NewM3I(), lin.NewV3().SetS(1, 1, 1))))),
				4.0))
	}
}
