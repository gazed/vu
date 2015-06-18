// Copyright © 2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"log"

	"github.com/gazed/vu"
	"github.com/gazed/vu/math/lin"
)

// lt tests the engines handling of the engine lighting shaders
// and the conversion of light position and normal vectors needed
// for proper lighting.
func lt() {
	lt := &lttag{}
	if err := vu.New(lt, "Lighting", 400, 100, 800, 600); err != nil {
		log.Printf("lt: error starting engine %s", err)
	}
	defer catchErrors()
}

// Globally unique "tag" that encapsulates example specific data.
type lttag struct {
	cam3D vu.Camera // 3D main scene camera.
	sun   vu.Pov    // Light node in Pov hierarchy.
	run   float64   // Camera movement speed.
	spin  float64   // Camera spin speed.
}

// Create is the engine callback for initial asset creation.
func (lt *lttag) Create(eng vu.Eng, s *vu.State) {
	lt.run = 10   // move so many cubes worth in one second.
	lt.spin = 270 // spin so many degrees in one second.
	top := eng.Root().NewPov()
	view := top.NewView()
	lt.cam3D = view.Cam()
	lt.cam3D.SetLocation(0.5, 2, 0.5)
	lt.sun = top.NewPov().SetLocation(0, 2.5, -1.75).SetScale(0.05, 0.05, 0.05)
	lt.sun.NewLight().SetColour(0.4, 0.7, 0.9)

	// Model at the light position.
	lt.sun.NewModel("solid").LoadMesh("sphere").LoadMat("sphere")

	// Create solid spheres to test the lighting shaders.
	c4 := top.NewPov().SetLocation(-0.5, 2, -2).SetScale(0.25, 0.25, 0.25)
	c4.NewModel("diffuse").LoadMesh("sphere").LoadMat("floor")
	c5 := top.NewPov().SetLocation(0.5, 2, -2).SetScale(0.25, 0.25, 0.25)
	c5.NewModel("gouraud").LoadMesh("sphere").LoadMat("floor")
	c6 := top.NewPov().SetLocation(1.5, 2, -2).SetScale(0.25, 0.25, 0.25)
	c6.NewModel("phong").LoadMesh("sphere").LoadMat("floor")
	lt.resize(s.W, s.H)
}

// Update is the regular engine callback.
func (lt *lttag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	if in.Resized {
		lt.resize(s.W, s.H)
	}
	// move the light.
	dt := in.Dt
	speed := lt.run * dt * 0.5
	for press, _ := range in.Down {
		switch press {
		case vu.K_W:
			lt.sun.Move(0, 0, -speed, lin.QI) // forward
		case vu.K_S:
			lt.sun.Move(0, 0, speed, lin.QI) // back
		case vu.K_A:
			lt.sun.Move(-speed, 0, 0, lin.QI) // left
		case vu.K_D:
			lt.sun.Move(speed, 0, 0, lin.QI) // right
		case vu.K_Q:
			lt.sun.Move(0, speed, 0, lin.QI) // up
		case vu.K_E:
			lt.sun.Move(0, -speed, 0, lin.QI) // down
		}
	}
}
func (lt *lttag) resize(ww, wh int) {
	lt.cam3D.SetPerspective(60, float64(ww)/float64(wh), 0.1, 50)
}
