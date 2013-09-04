// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package main

import (
	"log"
	"vu"
)

// bb tests the engines handling of textures and billboarding.
// This is more of a case to see what type of effects can be generated using
// a very simple texture scheme.
func bb() {
	bb := &bbtag{}
	var err error
	if bb.eng, err = vu.New("Texture:Load", 400, 100, 800, 600); err != nil {
		log.Printf("bb: error intitializing engine %s", err)
		return
	}
	bb.run = 10            // move so many cubes worth in one second.
	bb.spin = 270          // spin so many degrees in one second.
	bb.eng.SetDirector(bb) // override user input handling.
	bb.stagePlay()
	defer bb.eng.Shutdown()
	bb.eng.Action()
}

// Globally unique "tag" for this example.
type bbtag struct {
	eng   *vu.Eng
	scene vu.Scene
	run   float32
	spin  float32
}

func (bb *bbtag) stagePlay() {
	bb.scene = bb.eng.AddScene(vu.VP)
	bb.scene.SetPerspective(60, float32(800)/float32(600), 0.1, 50)
	bb.scene.SetLightLocation(0, 10, 0)
	bb.scene.SetLightColour(0.4, 0.7, 0.9)
	bb.scene.SetViewLocation(0, 0, 4)

	// load the floor model.
	floor := bb.scene.AddPart()
	floor.SetLocation(0, 1, 0)
	floor.SetFacade("floor", "gouraud", "floor")

	// load the billboard and textures.
	b1 := bb.scene.AddPart()
	b1.SetFacade("billboard", "bbra", "alpha")
	b1.SetTexture("core", 3)
	b1.SetLocation(0, 0, -2)
	b1.SetScale(0.25, 0.25, 0.25)

	b2 := bb.scene.AddPart()
	b2.SetFacade("billboard", "bbra", "alpha")
	b2.SetTexture("core", -1)
	b2.SetLocation(0, 0, -2)
	b2.SetScale(0.25, 0.25, 0.25)

	b3 := bb.scene.AddPart()
	b3.SetFacade("billboard", "bbra", "alpha")
	b3.SetTexture("halo", -2)
	b3.SetLocation(0, 0, -2)
	b3.SetScale(0.25, 0.25, 0.25)

	b4 := bb.scene.AddPart()
	b4.SetFacade("billboard", "bbra", "alpha")
	b4.SetTexture("halo", 1)
	b4.SetLocation(0, 0, -2)
	b4.SetScale(0.25, 0.25, 0.25)

	// set some constant state.
	bb.eng.Enable(vu.BLEND, true)
	bb.eng.Enable(vu.CULL, true)
	bb.eng.Enable(vu.DEPTH, true)
	bb.eng.Color(0.1, 0.1, 0.1, 1.0)
	return
}

// Handle engine callbacks.
func (bb *bbtag) Focus(focus bool) {}
func (bb *bbtag) Resize(x, y, width, height int) {
	bb.eng.ResizeViewport(x, y, width, height)
	ratio := float32(width) / float32(height)
	bb.scene.SetPerspective(60, ratio, 0.1, 50)
}
func (bb *bbtag) Update(pressed []string, gt, dt float32) {
	moveDelta := bb.eng.Dt
	for _, p := range pressed {
		switch p {
		case "W":
			bb.scene.MoveView(0, 0, moveDelta*-bb.run)
		case "S":
			bb.scene.MoveView(0, 0, moveDelta*bb.run)
		case "A":
			bb.scene.PanView(vu.YAxis, bb.eng.Dt*bb.spin)
		case "D":
			bb.scene.PanView(vu.YAxis, bb.eng.Dt*-bb.spin)
		}
	}
}
