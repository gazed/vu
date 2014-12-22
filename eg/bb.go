// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"log"

	"github.com/gazed/vu"
)

// bb tests the engines handling of textures, billboarding and banners.
// Similar images are generated using different techniques.
func bb() {
	bb := &bbtag{}
	var err error
	if bb.eng, err = vu.New("Billboarding & Banners", 400, 100, 800, 600); err != nil {
		log.Fatal("bb: error intitializing engine %s", err)
	}
	bb.eng.SetDirector(bb) // get user input through Director.Update()
	bb.create()            // create initial assests.
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
	cam   vu.Camera
	run   float64 // Camera movement speed.
	spin  float64 // Camera spin speed.
}

// create is the startup asset creation.
func (bb *bbtag) create() {
	bb.run = 10   // move so many cubes worth in one second.
	bb.spin = 270 // spin so many degrees in one second.
	bb.scene = bb.eng.AddScene(vu.VP)
	bb.scene.Set2D()
	bb.cam = bb.scene.Cam()
	bb.cam.SetPerspective(60, float64(800)/float64(600), 0.1, 50)
	bb.cam.SetLocation(0.5, 0, 0.25)

	// load the floor model.
	floor := bb.scene.AddPart()
	floor.SetLocation(0, 1, 0)
	floor.SetRole("gouraud").SetMesh("floor").SetMaterial("floor")
	floor.Role().SetUniform("alpha", 1.0)
	floor.Role().SetLightLocation(0, 10, 0)
	floor.Role().SetLightColour(0.4, 0.7, 0.9)

	// Create one image from billboards and textures in overlapping parts.
	b0 := bb.scene.AddPart()
	b0.SetLocation(1, 0, -1)
	b1 := b0.AddPart().SetScale(0.25, 0.25, 0.25)
	b1.SetRole("bbr").SetMesh("billboard").AddTex("core")
	b1.Role().SetUniform("spin", 3.0)
	b1.Role().SetAlpha(0.2)
	b2 := b0.AddPart().SetScale(0.25, 0.25, 0.25)
	b2.SetRole("bbr").SetMesh("billboard").AddTex("core")
	b2.Role().SetUniform("spin", -1.0)
	b2.Role().SetAlpha(0.2)
	b3 := b0.AddPart().SetScale(0.25, 0.25, 0.25)
	b3.SetRole("bbr").SetMesh("billboard").AddTex("halo")
	b3.Role().SetUniform("spin", -2.0)
	b3.Role().SetAlpha(0.2)
	b4 := b0.AddPart().SetScale(0.25, 0.25, 0.25)
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

	// Banner text with an ortho overlay.
	over := bb.eng.AddScene(vu.VO)
	over.Set2D()
	_, _, w, h := bb.eng.Size()
	over.Cam().SetOrthographic(0, float64(w), 0, float64(h), 0, 10)
	banner := over.AddPart()
	banner.SetRole("uv").AddTex("weblySleek22White").SetFont("weblySleek22").SetPhrase("Overlay Text")
	banner.SetLocation(100, 100, 0)

	// Try banner text with the 3D scene perspective camera.
	banner = bb.scene.AddPart().SetScale(0.1, 0.1, 0.1)
	banner.SetRole("uv").AddTex("weblySleek22White").SetFont("weblySleek22").SetPhrase("Floating Text")
	banner.SetLocation(-10, 2, -15)
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
			bb.cam.Move(0, 0, dt*-bb.run)
		case "S":
			bb.cam.Move(0, 0, dt*bb.run)
		case "A":
			bb.cam.Spin(0, dt*bb.spin, 0)
		case "D":
			bb.cam.Spin(0, dt*-bb.spin, 0)
		}
	}
}

// resize handles user screen/window changes.
func (bb *bbtag) resize() {
	x, y, width, height := bb.eng.Size()
	bb.eng.Resize(x, y, width, height)
	bb.cam.SetPerspective(60, float64(width)/float64(height), 0.1, 50)
}
