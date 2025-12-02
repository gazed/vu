// SPDX-FileCopyrightText : © 2022 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

package vu

import (
	"testing"

	"github.com/gazed/vu/math/lin"
)

// go test -run Body
func TestBody(t *testing.T) {

	// go test -run Sim/gravity
	t.Run("gravity on kinematic sphere", func(t *testing.T) {
		app := newApplication()
		scene := app.addScene(Scene3D)
		b1 := scene.AddModel("shd:test", "msh:ball", "mat:ball") // create pov
		b1.AddToSimulation(Sphere(25, KinematicSim))             // create static physics model

		// position before
		b10 := lin.NewV3().SetS(b1.At())
		// run one physics simulation step
		app.sim.simulate(app.povs, timestepSecs)
		// position after
		b11 := lin.NewV3().SetS(b1.At())

		// check movement
		if b10.Eq(b11) {
			t.Errorf("expected ball to drop")
		}
	})

	t.Run("static kinematic overlap", func(t *testing.T) {
		app := newApplication()
		scene := app.addScene(Scene3D)
		b1 := scene.AddModel("shd:test", "msh:ball", "mat:ball") // create pov
		b1.AddToSimulation(Sphere(25, StaticSim))                // create static physics model
		b2 := scene.AddModel("shd:test", "msh:ball", "mat:ball") // create pov
		b2.AddToSimulation(Sphere(25, KinematicSim))             // create physics model
		b2.SetAt(10, 10, 10)                                     // slightly overlapping

		// positions before
		b10, b20 := lin.NewV3().SetS(b1.At()), lin.NewV3().SetS(b2.At())
		// run one physics simulation step
		app.sim.simulate(app.povs, timestepSecs)
		// positions after
		b11, b21 := lin.NewV3().SetS(b1.At()), lin.NewV3().SetS(b2.At())

		// check movement
		if !b10.Eq(b11) {
			t.Errorf("expected static ball to not move")
		}
		if b20.Eq(b21) {
			t.Errorf("expected kinematic ball to move")
		}
	})

	t.Run("kinematic overlap", func(t *testing.T) {
		app := newApplication()
		scene := app.addScene(Scene3D)
		b1 := scene.AddModel("shd:test", "msh:ball", "mat:ball") // create pov
		b1.AddToSimulation(Sphere(25, KinematicSim))             // create physics model
		b2 := scene.AddModel("shd:test", "msh:ball", "mat:ball") // create pov
		b2.AddToSimulation(Sphere(25, KinematicSim))             // create physics model
		b2.SetAt(10, 10, 10)                                     // slightly overlapping

		// positions before
		b10, b20 := lin.NewV3().SetS(b1.At()), lin.NewV3().SetS(b2.At())
		// run one physics simulation step
		app.sim.simulate(app.povs, timestepSecs)
		// positions after
		b11, b21 := lin.NewV3().SetS(b1.At()), lin.NewV3().SetS(b2.At())

		// check movement
		if b10.Eq(b11) || b20.Eq(b21) {
			t.Errorf("expected both balls to move")
		}
	})
}
