// Copyright © 2014-2024 Galvanized Logic Inc.

package vu

import (
	"testing"

	"github.com/gazed/vu/math/lin"
)

// go test -run Sim
func TestSimulation(t *testing.T) {

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

// go test -run Bug
// Used to debug physics port.
// Adds test bodies to simulation and then loops the simulation.
// Match the results to the same loop output in raw-physics.
func TestBug(t *testing.T) {
	// app := newApplication()
	// scene := app.addScene(Scene3D)

	// slab := scene.AddModel("shd:pbr0", "msh:box0", "mat:box0")
	// slab.SetScale(50, 10, 50).SetAt(0, -5, 0)
	// slab.AddToSimulation(Box(50, 10, 50, StaticSim))

	// // const sphere_radius = 1.2849 // from blender
	// // b := scene.AddModel("shd:pbr0", "msh:sphere", "mat:sphere")
	// // b.SetScale(2, 2, 2).SetAt(0, 6.5, 0) // <- position of infinite loop
	// // b.SetAa(1, 0, 0, lin.Rad(-45))
	// // b.AddToSimulation(Sphere(sphere_radius, KinematicSim))

	// b := scene.AddModel("shd:pbr0", "msh:sphere", "mat:sphere")
	// b.SetScale(2, 2, 2).SetAt(0, 7, 0)
	// b.SetAa(1, 0, 0, lin.Rad(-45))
	// b.AddToSimulation(Box(2, 2, 2, KinematicSim))

	// // run one physics simulation step
	// for i := 0; i < 120; i++ {
	// 	app.sim.simulate(app.povs, timestepSecs)
	// 	pbod := (*physics.Body)(b.Body())
	// 	pos := pbod.Position()
	// 	vel := pbod.Velocity()
	// 	rot := pbod.Rotation()
	// 	fmt.Printf("%03d pos:%+9.6f %+9.6f %+9.6f vel:%+9.6f %+9.6f %+9.6f rot:%+9.6f %+9.6f %+9.6f %+9.6f\n",
	// 		i, pos.X, pos.Y, pos.Z, vel.X, vel.Y, vel.Z, rot.X, rot.Y, rot.Z, rot.W)
	// }
}
