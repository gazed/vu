// Copyright © 2013-2024 Galvanized Logic Inc.

package main

import (
	"log/slog"
	"math/rand"
	"time"

	"github.com/gazed/vu"
	"github.com/gazed/vu/math/lin"
)

// cr, collision resolution, demonstrates simulated physics by having balls and
// boxes bounce on a floor. The neat thing is that after the initial locations
// have been set the physics simulation handles all subsequent position updates.
// This example demonstrates:
//   - loading assets.
//   - creating a 3D scene with image and text models.
//   - setting and changing the scene camera location.
//   - adding a light to a scene.
//   - reacting to user input
//   - physics
//
// CONTROLS:
//   - A,D   : move the camera left/right around scene center.
//   - Space : throw a ball from the camera position.
//   - Q     : quit and close window.
func cr() {
	defer catchErrors()
	cr := &crtag{}
	eng, err := vu.NewEngine(
		vu.Windowed(),
		vu.Title("Collision Resolution"),
		vu.Size(200, 200, 800, 600),
		vu.Background(0.15, 0.15, 0.15, 1.0),
	)
	if err != nil {
		slog.Error("cr: engine start", "err", err)
		return
	}

	// seed some randomness for colors and start the engine.
	rand.Seed(time.Now().UTC().UnixNano())

	// Run will call Load once and then call Update each engine tick.
	eng.Run(cr, cr) // does not return while example is running.
}

// Globally unique "tag" that encapsulates example specific data.
type crtag struct {
	scene *vu.Entity
	pos   *lin.V3 // initial location.
	rot   float64 // rotation around origin.
}

// Load is the one time startup engine callback to create initial assets.
func (cr *crtag) Load(eng *vu.Engine) error {

	// import assets from asset files.
	// This creates the assets referenced by the models below.
	eng.ImportAssets("pbr0.shd", "sphere.glb", "box0.glb")

	// New scene with default camera.
	cr.pos = lin.NewV3().SetS(0, 16, 30)
	cr.scene = eng.AddScene(vu.Scene3D)
	cr.scene.Cam().SetAt(cr.pos.X, cr.pos.Y, cr.pos.Z)

	// add one directional light. SetAt sets the direction.
	cr.scene.AddLight(vu.DirectionalLight).SetAt(-1, -2, -2)

	// create a static slab as a base for the other physics objects.
	slab := cr.scene.AddModel("shd:pbr0", "msh:box0", "mat:box0")
	slab.SetScale(50, 10, 50).SetAt(0, -5, 0)
	slab.AddToSimulation(vu.Box(50, 10, 50, vu.StaticSim))

	// create a block of physics cubes.
	cubeSize := 3
	for x := 0; x < cubeSize; x++ {
		for y := 0; y < cubeSize; y++ {
			for z := 0; z < cubeSize; z++ {
				lx := float64(x)*4.05 - 2.0
				ly := float64(y)*4.05 + 12.0
				lz := float64(z)*4.05 - 2.0
				cr.makeBox(lx, ly, lz)
			}
		}
	}
	return nil
}

// Update is the ongoing engine callback.
func (cr *crtag) Update(eng *vu.Engine, in *vu.Input, delta time.Duration) {
	cam := cr.scene.Cam()

	// react to one time press events.
	for press := range in.Pressed {
		switch press {
		case vu.KQ:
			// quit if Q is pressed
			eng.Shutdown()
			return
		case vu.KSpace:
			// throw a ball from the camera's viewpoint.
			atx, aty, atz := cam.At()
			ball := cr.makeBall(atx, aty-2.0, atz)
			lookat := cam.Lookat()
			forward := lin.NewV3().Forward(lookat).Unit()
			throw := lin.NewV3().Scale(forward, -30)
			ball.Push(throw.X, throw.Y, throw.Z)
		}
	}

	// react to continuous press events.
	lookSpeed := 40 * delta.Seconds()
	for press := range in.Down {
		switch press {
		case vu.KA:
			// rotate camera left around center
			cr.rot -= lookSpeed
			transformAroundOrigin := lin.NewT().SetLoc(0, 0, 0).SetAa(0, 1, 0, lin.Rad(cr.rot))
			at := transformAroundOrigin.App(lin.NewV3().Set(cr.pos))
			cam.SetAt(at.X, at.Y, at.Z)
			cam.SetYaw(cr.rot)
		case vu.KD:
			// rotate camera right around center
			cr.rot += lookSpeed
			transformAroundOrigin := lin.NewT().SetLoc(0, 0, 0).SetAa(0, 1, 0, lin.Rad(cr.rot))
			at := transformAroundOrigin.App(lin.NewV3().Set(cr.pos))
			cam.SetAt(at.X, at.Y, at.Z)
			cam.SetYaw(cr.rot)
		}
	}
}

// makeBall creates a visible sphere physics body.
func (cr *crtag) makeBall(lx, ly, lz float64) (ball *vu.Entity) {
	const sphere_radius = 1.2849 // from blender
	ball = cr.scene.AddModel("shd:pbr0", "msh:sphere")
	ball.SetScale(2, 2, 2).SetAt(lx, ly, lz)
	ball.AddToSimulation(vu.Sphere(2*sphere_radius, vu.KinematicSim))
	r, g, b, a, metallic, roughness := randomColor()
	ball.SetColor(r, g, b, a)
	ball.SetMetallicRoughness(metallic, roughness)
	return ball
}

// makeBox creates a visible box physics body.
func (cr *crtag) makeBox(lx, ly, lz float64) (box *vu.Entity) {
	box = cr.scene.AddModel("shd:pbr0", "msh:box0")
	box.SetScale(2, 2, 2).SetAt(lx, ly, lz)
	box.AddToSimulation(vu.Box(2, 2, 2, vu.KinematicSim))
	r, g, b, a, metallic, roughness := randomColor()
	box.SetColor(r, g, b, a)
	box.SetMetallicRoughness(metallic, roughness)
	return box
}

// randomColor generates a random PBR solid color.
func randomColor() (r, g, b, a float64, metallic bool, roughness float64) {
	r = rand.Float64()
	g = rand.Float64()
	b = rand.Float64()
	a = 1.0
	metallic = rand.Float64() >= 0.5
	roughness = rand.Float64()
	return r, g, b, a, metallic, roughness
}
