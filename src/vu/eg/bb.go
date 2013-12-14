// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package main

import (
	"log"
	"vu"
)

// bb tests the engines handling of textures, billboarding and banners.
// This example is used to see what type of effects can be generated
// using simple textures.
func bb() {
	bb := &bbtag{}
	var err error
	if bb.eng, err = vu.New("Billboarding & Banners", 400, 100, 800, 600); err != nil {
		log.Printf("bb: error intitializing engine %s", err)
		return
	}
	bb.run = 10            // move so many cubes worth in one second.
	bb.spin = 270          // spin so many degrees in one second.
	bb.eng.SetDirector(bb) // override user input handling.
	defer bb.eng.Shutdown()
	bb.eng.Action()
}

// Globally unique "tag" that encapsulates example specific data.
type bbtag struct {
	eng   vu.Engine
	scene vu.Scene
	run   float64 // Camera movement speed.
	spin  float64 // Camera spin speed.
}

// Create is the engine intialization callback.
func (bb *bbtag) Create(eng vu.Engine) {
	bb.scene = eng.AddScene(vu.VP)
	bb.scene.SetPerspective(60, float64(800)/float64(600), 0.1, 50)
	bb.scene.SetLightLocation(0, 10, 0)
	bb.scene.SetLightColour(0.4, 0.7, 0.9)
	bb.scene.SetViewLocation(0, 0, 4)
	bb.scene.Set2D()

	// load the floor model.
	floor := bb.scene.AddPart()
	floor.SetLocation(0, 1, 0)
	floor.SetFacade("floor", "gouraud").SetMaterial("floor")

	// load the billboard and textures.
	b1 := bb.scene.AddPart()
	b1.SetFacade("billboard", "bbra").SetTexture("core", 3)
	b1.SetAlpha(0.2)
	b1.SetLocation(0, 0, -2)
	b1.SetScale(0.25, 0.25, 0.25)

	b2 := bb.scene.AddPart()
	b2.SetFacade("billboard", "bbra").SetTexture("core", -1)
	b2.SetAlpha(0.2)
	b2.SetLocation(0, 0, -2)
	b2.SetScale(0.25, 0.25, 0.25)

	b3 := bb.scene.AddPart()
	b3.SetFacade("billboard", "bbra").SetTexture("halo", -2)
	b3.SetAlpha(0.2)
	b3.SetLocation(0, 0, -2)
	b3.SetScale(0.25, 0.25, 0.25)

	b4 := bb.scene.AddPart()
	b4.SetFacade("billboard", "bbra").SetTexture("halo", 1)
	b4.SetAlpha(0.2)
	b4.SetLocation(0, 0, -2)
	b4.SetScale(0.25, 0.25, 0.25)

	// Add a banner overlay
	over := eng.AddScene(vu.VO)
	over.Set2D()
	_, _, w, h := bb.eng.Size()
	over.SetOrthographic(0, float64(w), 0, float64(h), 0, 10)
	over.SetLightLocation(1, 1, 1)
	over.SetLightColour(1, 1, 1)
	banner := over.AddPart()
	banner.SetBanner("Banner Text", "uv", "weblySleek22", "weblySleek22White")
	banner.SetLocation(100, 100, 0)

	// set some constant state.
	bb.eng.Enable(vu.BLEND, true)
	bb.eng.Enable(vu.CULL, true)
	bb.eng.Color(0.1, 0.1, 0.1, 1.0)
}

// Update is the regular engine callback.
func (bb *bbtag) Update(input *vu.Input) {
	if input.Resized {
		bb.resize()
	}
	dt := input.Dt
	for press, _ := range input.Down {
		switch press {
		case "W":
			bb.scene.MoveView(0, 0, dt*-bb.run)
		case "S":
			bb.scene.MoveView(0, 0, dt*bb.run)
		case "A":
			bb.scene.PanView(vu.YAxis, dt*bb.spin)
		case "D":
			bb.scene.PanView(vu.YAxis, dt*-bb.spin)
		}
	}
}

// resize handles user screen/window changes.
func (bb *bbtag) resize() {
	x, y, width, height := bb.eng.Size()
	bb.eng.Resize(x, y, width, height)
	bb.scene.SetPerspective(60, float64(width)/float64(height), 0.1, 50)
}
