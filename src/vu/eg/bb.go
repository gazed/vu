// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package main

import (
	"log"
	"vu"
)

// bb tests the engines handling of textures, billboarding and banners.
// Similar images are generated using different techniques.
func bb() {
	bb := &bbtag{}
	var err error
	if bb.eng, err = vu.New("Billboarding & Banners", 400, 100, 800, 600); err != nil {
		log.Fatal("bb: error intitializing engine %s", err)
	}
	bb.eng.SetDirector(bb) // override user input handling.
	if err = bb.eng.Verify(); err != nil {
		log.Fatalf("bb: error initializing model :: %s", err)
	}
	defer bb.eng.Shutdown()
	defer catchErrors()
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
	bb.run = 10   // move so many cubes worth in one second.
	bb.spin = 270 // spin so many degrees in one second.
	bb.scene = eng.AddScene(vu.VP)
	bb.scene.SetPerspective(60, float64(800)/float64(600), 0.1, 50)
	bb.scene.SetLightLocation(0, 10, 0)
	bb.scene.SetLightColour(0.4, 0.7, 0.9)
	bb.scene.SetLocation(0.5, 0, 0.25)
	bb.scene.Set2D()

	// load the floor model.
	floor := bb.scene.AddPart()
	floor.SetRole("gouraud").SetMesh("floor").SetMaterial("floor")
	floor.Role().SetUniform("alpha", 1.0)
	floor.SetLocation(0, 1, 0)

	// Create one image from billboards and textures in overlapping parts.
	b1 := bb.scene.AddPart().SetLocation(1, 0, -1).SetScale(0.25, 0.25, 0.25)
	b1.SetRole("bbr").SetMesh("billboard").AddTex("core")
	b1.Role().SetUniform("spin", 3.0)
	b1.Role().SetAlpha(0.2)
	b2 := bb.scene.AddPart().SetLocation(1, 0, -1).SetScale(0.25, 0.25, 0.25)
	b2.SetRole("bbr").SetMesh("billboard").AddTex("core")
	b2.Role().SetUniform("spin", -1.0)
	b2.Role().SetAlpha(0.2)
	b3 := bb.scene.AddPart().SetLocation(1, 0, -1).SetScale(0.25, 0.25, 0.25)
	b3.SetRole("bbr").SetMesh("billboard").AddTex("halo")
	b3.Role().SetUniform("spin", -2.0)
	b3.Role().SetAlpha(0.2)
	b4 := bb.scene.AddPart().SetLocation(1, 0, -1).SetScale(0.25, 0.25, 0.25)
	b4.SetRole("bbr").SetMesh("billboard").AddTex("halo")
	b4.Role().SetUniform("spin", 1.0)
	b4.Role().SetAlpha(0.2)

	// Create a similar image using a single multi-texture shader.
	c4 := bb.scene.AddPart()
	c4.SetRole("spinball").SetMesh("billboard").AddTex("core")
	c4.Role().AddTex("core").AddTex("halo").AddTex("halo")
	c4.Role().SetAlpha(0.4)
	c4.SetLocation(0.5, 0, -1)
	c4.SetScale(0.25, 0.25, 0.25)

	// Add a banner overlay
	over := eng.AddScene(vu.VO)
	over.Set2D()
	_, _, w, h := bb.eng.Size()
	over.SetOrthographic(0, float64(w), 0, float64(h), 0, 10)
	banner := over.AddPart()
	banner.SetRole("uv").AddTex("weblySleek22White").SetFont("weblySleek22").SetPhrase("Banner Text")
	banner.SetLocation(100, 100, 0)
}

// Update is the regular engine callback.
func (bb *bbtag) Update(in *vu.Input) {
	if in.Resized {
		bb.resize()
	}
	dt := in.Dt
	for press, _ := range in.Down {
		switch press {
		case "W":
			bb.scene.Move(0, 0, dt*-bb.run)
		case "S":
			bb.scene.Move(0, 0, dt*bb.run)
		case "A":
			bb.scene.Spin(vu.YAxis, dt*bb.spin)
		case "D":
			bb.scene.Spin(vu.YAxis, dt*-bb.spin)
		}
	}
}

// resize handles user screen/window changes.
func (bb *bbtag) resize() {
	x, y, width, height := bb.eng.Size()
	bb.eng.Resize(x, y, width, height)
	bb.scene.SetPerspective(60, float64(width)/float64(height), 0.1, 50)
}
