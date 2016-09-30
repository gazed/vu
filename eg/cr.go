// Copyright © 2013-2016 Galvanized Logic Inc.
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
// Set useBalls to false and "go build" to have the demo use cubes.
//
// CONTROLS:
//   WASD  : move the light position: forward left back right
//   ZX    : move light position    : up down
//   LaRa  : spin the cube          : left right.
func cr() {
	cr := &crtag{}
	if err := vu.New(cr, "Collision Resolution", 400, 100, 800, 600); err != nil {
		log.Printf("cr: error initializing engine %s", err)
	}
	defer catchErrors()
}

// Globally unique "tag" that encapsulates example specific data.
type crtag struct {
	top     *vu.Pov
	cam     *vu.Camera
	striker *vu.Pov // Move to hit other items.
}

// Create is the engine callback for initial asset creation.
func (cr *crtag) Create(eng vu.Eng, s *vu.State) {
	cr.top = eng.Root().NewPov()
	sun := cr.top.NewPov().SetAt(0, 10, 10)
	sun.NewLight().SetColor(0.8, 0.8, 0.8)
	cr.cam = cr.top.NewCam()
	cr.cam.SetPerspective(60, float64(800)/float64(600), 0.1, 500)
	cr.cam.SetAt(0, 10, 25)

	// load the static slab.
	slab := cr.top.NewPov().SetScale(50, 50, 50).SetAt(0, -25, 0)
	slab.NewBody(vu.NewBox(25, 25, 25))
	slab.SetSolid(0, 0.4)
	slab.NewModel("diffuse", "msh:box", "mat:gray")

	// create a single moving body.
	cr.striker = cr.top.NewPov()
	cr.striker.SetAt(15, 15, 0)
	useBalls := true // Flip to use boxes instead of spheres.
	if useBalls {
		cr.getBall(cr.striker)
	} else {
		cr.getBox(cr.striker)
	}
	cr.striker.Model().SetUniform("kd", rand.Float64(), rand.Float64(), rand.Float64())

	// create a block of physics bodies.
	cubeSize := 3
	startX := -5 - cubeSize/2
	startY := -5
	startZ := -3 - cubeSize/2
	for k := 0; k < cubeSize; k++ {
		for i := 0; i < cubeSize; i++ {
			for j := 0; j < cubeSize; j++ {
				bod := cr.top.NewPov()
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
	if in.Resized {
		cr.cam.SetPerspective(60, float64(s.W)/float64(s.H), 0.1, 50)
	}
	dt := in.Dt
	for press := range in.Down {
		switch press {
		case vu.KW:
			cr.cam.Move(0, 0, dt*-run, cr.cam.Look)
		case vu.KS:
			cr.cam.Move(0, 0, dt*run, cr.cam.Look)
		case vu.KA:
			cr.cam.SetYaw(cr.cam.Yaw + spin*dt)
		case vu.KD:
			cr.cam.SetYaw(cr.cam.Yaw - spin*dt)
		case vu.KB:
			ball := cr.top.NewPov()
			ball.SetAt(-2.5+rand.Float64(), 15, -1.5-rand.Float64())
			ball.NewBody(vu.NewSphere(1))
			ball.SetSolid(1, 0.9)
			m := ball.NewModel("gouraud", "msh:sphere", "mat:red")
			m.SetUniform("kd", rand.Float64(), rand.Float64(), rand.Float64())
		case vu.KSpace:
			body := cr.striker.Body()
			body.Push(-2.5, 0, -0.5)
		}
	}
}

// getBall creates a visible sphere physics body.
func (cr *crtag) getBall(p *vu.Pov) {
	p.NewBody(vu.NewSphere(1))
	p.SetSolid(1, 0.5)
	p.NewModel("gouraud", "msh:sphere", "mat:red")
}

// getBox creates a visible box physics body.
func (cr *crtag) getBox(p *vu.Pov) {
	p.SetScale(2, 2, 2)
	p.NewBody(vu.NewBox(1, 1, 1))
	p.SetSolid(1, 0)
	p.NewModel("gouraud", "msh:box", "mat:red")
}
