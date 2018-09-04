// Copyright © 2013-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/gazed/vu"
)

// cr, collision resolution, demonstrates simulated physics by having balls bounce
// on a floor. The neat thing is that after the initial locations have been set
// the physics simulation handles all subsequent position updates.
//
// Alternatively set useBalls to false and "go build" to have the demo use cubes.
//
// CONTROLS:
//   WASD  : move the camera.
//   B     : generate new falling spheres.
//   Space : accellerate the striker sphere.
func cr() {
	defer catchErrors()
	if err := vu.Run(&crtag{}); err != nil {
		log.Printf("cr: error initializing engine %s", err)
	}
}

// Globally unique "tag" that encapsulates example specific data.
type crtag struct {
	scene   *vu.Ent
	striker *vu.Ent // Move to hit other items.
}

// Create is the engine callback for initial asset creation.
func (cr *crtag) Create(eng vu.Eng, s *vu.State) {
	eng.Set(vu.Title("Collision Resolution"), vu.Size(400, 100, 800, 600))

	// New scene with default camera.
	cr.scene = eng.AddScene()
	cr.scene.Cam().SetClip(0.1, 100).SetFov(60).SetAt(0, 10, 25)
	sun := cr.scene.MakeLight(vu.DirectionalLight).SetLightColor(0.8, 0.8, 0.8)
	sun.SetAt(0, 10, 10)

	// load the static slab.
	slab := cr.scene.AddPart().SetScale(50, 50, 50).SetAt(0, -25, 0)
	slab.MakeBody(vu.Box(25, 25, 25)).SetSolid(0, 0.4)
	slab.MakeModel("diffuse", "msh:box", "mat:gray")

	// create a single moving body.
	cr.striker = cr.scene.AddPart().SetAt(15, 15, 0)
	useBalls := true // Flip to use boxes instead of spheres.
	if useBalls {
		cr.getBall(cr.striker)
	} else {
		cr.getBox(cr.striker)
	}
	cr.striker.SetUniform("kd", rand.Float64(), rand.Float64(), rand.Float64())

	// create a block of physics bodies.
	cubeSize := 3
	startX := -5 - cubeSize/2
	startY := -5
	startZ := -3 - cubeSize/2
	for k := 0; k < cubeSize; k++ {
		for i := 0; i < cubeSize; i++ {
			for j := 0; j < cubeSize; j++ {
				bod := cr.scene.AddPart()
				lx := float64(2*i + startX)
				ly := float64(20 + 2*k + startY)
				lz := float64(2*j + startZ)
				bod.SetAt(lx, ly, lz)
				if useBalls {
					cr.getBall(bod)
				} else {
					cr.getBox(bod)
				}
			}
		}
	}

	// set non default engine state.
	eng.Set(vu.Color(0.15, 0.15, 0.15, 1))
	rand.Seed(time.Now().UTC().UnixNano())
}

// Update is the regular engine callback.
func (cr *crtag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	run := 10.0   // move so many cubes worth in one second.
	spin := 270.0 // spin so many degrees in one second.
	dt := in.Dt
	cam := cr.scene.Cam()
	for press := range in.Down {
		switch press {
		case vu.KW:
			cam.Move(0, 0, dt*-run, cam.Look)
		case vu.KS:
			cam.Move(0, 0, dt*run, cam.Look)
		case vu.KA:
			cam.SetYaw(cam.Yaw + spin*dt)
		case vu.KD:
			cam.SetYaw(cam.Yaw - spin*dt)
		case vu.KB:
			ball := cr.scene.AddPart()
			ball.SetAt(-2.5+rand.Float64(), 15, -1.5-rand.Float64())
			ball.MakeBody(vu.Sphere(1)).SetSolid(1, 0.9)
			m := ball.MakeModel("phong", "msh:sphere", "mat:red")
			m.SetUniform("kd", rand.Float64(), rand.Float64(), rand.Float64())
		case vu.KSpace:
			cr.striker.Push(-2.5, 0, -0.5)
		}
	}
}

// getBall creates a visible sphere physics body.
func (cr *crtag) getBall(p *vu.Ent) {
	p.MakeBody(vu.Sphere(1)).SetSolid(1, 0.5)
	p.MakeModel("phong", "msh:sphere", "mat:red")
}

// getBox creates a visible box physics body.
func (cr *crtag) getBox(p *vu.Ent) {
	p.SetScale(2, 2, 2)
	p.MakeBody(vu.Box(1, 1, 1)).SetSolid(1, 0)
	p.MakeModel("phong", "msh:box", "mat:red")
}
