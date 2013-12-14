// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package main

import (
	"log"
	"math/rand"
	"time"
	"vu"
)

// cr, collision resolution, demonstrates simulated physics by having balls bounce
// on a floor. The neat thing is that after the initial locations have been set
// the physics simulation (vu/move) handles all subsequent position updates.
func cr() {
	cr := &crtag{}
	var err error
	if cr.eng, err = vu.New("Collision Resolution", 400, 100, 800, 600); err != nil {
		log.Printf("cr: error intitializing engine %s", err)
		return
	}
	cr.run = 10            // move so many cubes worth in one second.
	cr.spin = 270          // spin so many degrees in one second.
	cr.eng.SetDirector(cr) // override user input handling.
	defer cr.eng.Shutdown()
	cr.eng.Action()
}

// Globally unique "tag" that encapsulates example specific data.
type crtag struct {
	eng   vu.Engine
	scene vu.Scene
	bod   vu.Part
	run   float64 // Camera movement speed.
	spin  float64 // Camera spin speed.
}

// Create is the engine intialization callback. Create the physics objects.
func (cr *crtag) Create(eng vu.Engine) {
	cr.scene = eng.AddScene(vu.VP)
	cr.scene.SetPerspective(60, float64(800)/float64(600), 0.1, 50)
	cr.scene.SetLightLocation(0, 5, 0)
	cr.scene.SetLightColour(0.4, 0.7, 0.9)
	cr.scene.SetViewLocation(0, 10, 25)

	// load the static slab.
	slab := cr.scene.AddPart()
	slab.SetFacade("cube", "gouraud").SetMaterial("floor")
	slab.SetScale(50, 50, 50)
	slab.SetBody(vu.Box(25, 25, 25), 0, 0.4)
	slab.SetLocation(0, -25, 0)

	// create a single moving body.
	useBalls := true
	cr.bod = cr.scene.AddPart()
	cr.bod.SetLocation(15, 15, 0) // -5, 15, -3
	if useBalls {
		cr.getBall(cr.bod)
	} else {
		cr.getBox(cr.bod)
		cr.bod.SetRotation(0.1825742, 0.3651484, 0.5477226, 0.7302967)
	}

	// Box can be used as a ball replacement.

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
				bod.SetLocation(lx, ly, lz)
				if useBalls {
					cr.getBall(bod)
				} else {
					cr.getBox(bod)
				}
			}
		}
	}

	// set some constant state.
	rand.Seed(time.Now().UTC().UnixNano())
	cr.eng.Enable(vu.BLEND, true)
	cr.eng.Enable(vu.CULL, true)
	cr.eng.Enable(vu.DEPTH, true)
	cr.eng.Color(0.1, 0.1, 0.1, 1.0)
	return
}

// Update is the regular engine callback.
func (cr *crtag) Update(input *vu.Input) {
	if input.Resized {
		cr.resize()
	}
	dt := input.Dt
	for press, _ := range input.Down {
		switch press {
		case "W":
			cr.scene.MoveView(0, 0, dt*-cr.run)
		case "S":
			cr.scene.MoveView(0, 0, dt*cr.run)
		case "A":
			cr.scene.PanView(vu.YAxis, dt*cr.spin)
		case "D":
			cr.scene.PanView(vu.YAxis, dt*-cr.spin)
		case "B":
			ball := cr.scene.AddPart()
			ball.SetFacade("sphere", "gouraud").SetMaterial("sphere")
			ball.SetBody(vu.Sphere(1), 1, 0.9)
			ball.SetLocation(-2.5+rand.Float64(), 15, -1.5-rand.Float64())
		case "Sp":
			cr.bod.Push(-2.5, 0, -0.5)
		}
	}
}

// resize handles user screen/window changes.
func (cr *crtag) resize() {
	x, y, width, height := cr.eng.Size()
	cr.eng.Resize(x, y, width, height)
	cr.scene.SetPerspective(60, float64(width)/float64(height), 0.1, 50)
}

// getBall creates a visible sphere physics body.
func (cr *crtag) getBall(bod vu.Part) {
	bod.SetFacade("sphere", "gouraud").SetMaterial("sphere")
	bod.SetBody(vu.Sphere(1), 1, 0.5)
}

// getBox creates a visible box physics body.
func (cr *crtag) getBox(bod vu.Part) {
	bod.SetScale(2, 2, 2)
	bod.SetFacade("cube", "gouraud").SetMaterial("sphere")
	bod.SetBody(vu.Box(1, 1, 1), 1, 0)
}
