// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package main

import (
	"log"
	"vu"
	"vu/physics"
)

// co demonstrates simple physics by having a ball bounce on a floor.
// Note that the co example sets up some initial locations and momentum
// and vu/physics handles all subsequent movement/location updating.
func co() {
	co := &cotag{}
	var err error
	if co.eng, err = vu.New("Collision", 400, 100, 800, 600); err != nil {
		log.Printf("co: error intitializing engine %s", err)
		return
	}
	co.run = 10            // move so many cubes worth in one second.
	co.spin = 270          // spin so many degrees in one second.
	co.eng.SetDirector(co) // override user input handling.
	co.setStage()
	defer co.eng.Shutdown()
	co.eng.Action()
}

// Globally unique "tag" for this example.
type cotag struct {
	eng   *vu.Eng
	scene vu.Scene
	ball  vu.Part
	run   float32
	spin  float32
}

// Create the shooting range with one block as a target.
func (co *cotag) setStage() {
	co.scene = co.eng.AddScene(vu.VP)
	co.scene.SetPerspective(60, float32(800)/float32(600), 0.1, 50)
	co.scene.SetLightLocation(0, 10, 0)
	co.scene.SetLightColour(0.4, 0.7, 0.9)
	co.scene.SetViewLocation(0, 1, 5)

	// load the floor model.
	floor := co.scene.AddPart()
	floor.SetLocation(0, 1, 0)
	floor.SetFacade("floor", "gouraud", "floor")
	floor.SetBody(0, 0)
	floor.SetShape(physics.Plane(0, 1, 0, 0, 0, 0))

	// create a moving part.
	co.ball = co.scene.AddPart()
	co.ball.SetLocation(0, 10, 0)
	co.ball.SetFacade("sphere", "gouraud", "sphere")
	co.ball.SetBody(10, 100)
	co.ball.SetShape(physics.Sphere(0, 0, 0, 0.15))
	co.ball.ResetMomentum()
	co.ball.SetLinearMomentum(0, 0, -3)

	// set some constant state.
	co.eng.Enable(vu.BLEND, true)
	co.eng.Enable(vu.CULL, true)
	co.eng.Enable(vu.DEPTH, true)
	co.eng.Color(0.1, 0.1, 0.1, 1.0)
	return
}

// Handle engine callbacks.
func (co *cotag) Focus(focus bool) {}
func (co *cotag) Resize(x, y, width, height int) {
	co.eng.ResizeViewport(x, y, width, height)
	ratio := float32(width) / float32(height)
	co.scene.SetPerspective(60, ratio, 0.1, 50)
}
func (co *cotag) Update(pressed []string, gameTime, deltaTime float32) {
	moveDelta := co.eng.Dt
	for _, p := range pressed {
		switch p {
		case "W":
			co.scene.MoveView(0, 0, moveDelta*-co.run)
		case "S":
			co.scene.MoveView(0, 0, moveDelta*co.run)
		case "A":
			co.scene.PanView(vu.YAxis, co.eng.Dt*co.spin)
		case "D":
			co.scene.PanView(vu.YAxis, co.eng.Dt*-co.spin)
		case "Sp":
			co.ball.SetLocation(0, 5, -4)
			co.ball.ResetMomentum()
			co.ball.SetLinearMomentum(0, 1, 0)
		case "/":
			co.ball.SetAngularMomentum(0, 2, -1)
		case ",":
			co.ball.SetAngularMomentum(1, 2, 3)
		}
	}
}
