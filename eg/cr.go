// Copyright © 2013-2014 Galvanized Logic Inc.
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
// the physics simulation (vu/move) handles all subsequent position updates.
func cr() {
	cr := &crtag{}
	var err error
	if cr.eng, err = vu.New("Collision Resolution", 400, 100, 800, 600); err != nil {
		log.Printf("cr: error intitializing engine %s", err)
		return
	}
	cr.eng.SetDirector(cr) // get user input through Director.Update()
	cr.create()            // create initial assests.
	defer cr.eng.Shutdown()
	defer catchErrors()
	cr.eng.Action()
}

// Globally unique "tag" that encapsulates example specific data.
type crtag struct {
	eng     vu.Engine
	scene   vu.Scene
	cam     vu.Camera
	striker vu.Part // Move to hit other items.
	run     float64 // Camera movement speed.
	spin    float64 // Camera spin speed.
}

// create is the startup asset creation.
func (cr *crtag) create() {
	cr.run = 10   // move so many cubes worth in one second.
	cr.spin = 270 // spin so many degrees in one second.
	cr.scene = cr.eng.AddScene(vu.VP)
	cr.cam = cr.scene.Cam()
	cr.cam.SetPerspective(60, float64(800)/float64(600), 0.1, 500)
	cr.cam.SetLocation(0, 10, 25)

	// load the static slab.
	slab := cr.scene.AddPart().SetScale(50, 50, 50)
	slab.SetLocation(0, -25, 0)
	slab.SetBody(vu.NewBox(25, 25, 25), 0, 0.4)
	slab.SetRole("gouraud").SetMesh("cube").SetMaterial("floor")
	slab.Role().SetLightLocation(0, 5, -100)
	slab.Role().SetLightColour(0.4, 0.7, 0.9)

	// create a single moving body.
	useBalls := true // Change this to use boxes instead of spheres.
	cr.striker = cr.scene.AddPart()
	cr.striker.SetLocation(15, 15, 0) // -5, 15, -3
	if useBalls {
		cr.getBall(cr.striker)
	} else {
		cr.getBox(cr.striker)
		cr.striker.SetRotation(0.1825742, 0.3651484, 0.5477226, 0.7302967)
	}

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
	cr.eng.Color(0.15, 0.15, 0.15, 1)
	rand.Seed(time.Now().UTC().UnixNano())
	return
}

// Update is the regular engine callback.
func (cr *crtag) Update(in *vu.Input) {
	if in.Resized {
		cr.resize()
	}
	dt := in.Dt
	for press, _ := range in.Down {
		switch press {
		case "W":
			cr.cam.Move(0, 0, dt*-cr.run)
		case "S":
			cr.cam.Move(0, 0, dt*cr.run)
		case "A":
			cr.cam.Spin(0, dt*cr.spin, 0)
		case "D":
			cr.cam.Spin(0, dt*-cr.spin, 0)
		case "B":
			ball := cr.scene.AddPart()
			ball.SetLocation(-2.5+rand.Float64(), 15, -1.5-rand.Float64())
			ball.SetBody(vu.NewSphere(1), 1, 0.9)
			ball.SetRole("gouraud").SetMesh("sphere").SetMaterial("sphere")
			ball.Role().SetLightLocation(0, 5, 0)
			r := rand.Float64()
			g := rand.Float64()
			b := rand.Float64()
			ball.Role().SetLightColour(r, g, b)
		case "Sp":
			cr.striker.Body().Push(-2.5, 0, -0.5)
		}
	}
}

// resize handles user screen/window changes.
func (cr *crtag) resize() {
	x, y, width, height := cr.eng.Size()
	cr.eng.Resize(x, y, width, height)
	cr.cam.SetPerspective(60, float64(width)/float64(height), 0.1, 50)
}

// getBall creates a visible sphere physics body.
func (cr *crtag) getBall(bod vu.Part) {
	bod.SetRole("gouraud").SetMesh("sphere").SetMaterial("sphere")
	bod.SetBody(vu.NewSphere(1), 1, 0.5)
	bod.Role().SetLightLocation(0, 5, 0)
	bod.Role().SetLightColour(0.5, 0.9, 0.7)
}

// getBox creates a visible box physics body.
func (cr *crtag) getBox(bod vu.Part) {
	bod.SetScale(2, 2, 2)
	bod.SetRole("gouraud").SetMesh("cube").SetMaterial("sphere")
	bod.SetBody(vu.NewBox(1, 1, 1), 1, 0)
	bod.Role().SetLightLocation(0, 5, 0)
	bod.Role().SetLightColour(0.9, 0.5, 0.7)
}
