// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"log"

	"github.com/gazed/vu"
)

// bb tests the engines handling of billboards and banners as well
// as combining multiple textures using shaders.
func bb() {
	bb := &bbtag{}
	if err := vu.New(bb, "Billboarding & Banners", 400, 100, 800, 600); err != nil {
		log.Printf("bb: error starting engine %s", err)
	}
	defer catchErrors()
}

// Globally unique "tag" that encapsulates example specific data.
type bbtag struct {
	cam        vu.Camera // Allow the user to move the camera.
	ui         vu.Camera // 2D user interface.
	run        float64   // Camera movement speed.
	spin       float64   // Camera spin speed.
	screenText vu.Pov    // Screen space text.
}

// Create is the startup asset creation.
func (bb *bbtag) Create(eng vu.Eng, s *vu.State) {
	bb.run = 10   // move so many cubes worth in one second.
	bb.spin = 270 // spin so many degrees in one second.
	top := eng.Root().NewPov()
	view := top.NewView()
	bb.cam = view.Cam()
	bb.cam.SetLocation(0.5, 2, 2.5)
	sun := top.NewPov().SetLocation(0, 3, -3)
	sun.NewLight().SetColour(0.4, 0.7, 0.9)

	// Load the floor model.
	floor := top.NewPov()
	floor.NewModel("gouraud").LoadMesh("floor").LoadMat("floor")

	// Create a single image from multiple textures using a shader.
	c4 := top.NewPov().SetLocation(0.5, 2, -1).SetScale(0.25, 0.25, 0.25)
	model := c4.NewModel("spinball").LoadMesh("billboard")
	model.AddTex("core").AddTex("core").AddTex("halo").AddTex("halo")
	model.SetAlpha(0.4)

	// Try banner text with the 3D scene perspective camera.
	font := "weblySleek22"
	banner := top.NewPov().SetScale(0.1, 0.1, 0.1).SetLocation(-10, 3, -15)
	banner.NewModel("uv").AddTex(font + "White").LoadFont(font).SetPhrase("Floating Text")

	// Try billboard banner text with the 3D scene perspective camera.
	banner = top.NewPov().SetScale(0.025, 0.025, 0.025).SetLocation(-10, 2, -15)
	banner.NewModel("bb").AddTex(font + "White").LoadFont(font).SetPhrase("Billboard Text")

	// Banner text with an ortho overlay.
	v2D := eng.Root().NewPov()
	view2D := v2D.NewView()
	view2D.SetUI()
	bb.ui = view2D.Cam()

	// 2D static location.
	banner = v2D.NewPov().SetLocation(100, 100, 0)
	banner.NewModel("uv").AddTex(font + "White").LoadFont(font).SetPhrase("Overlay Text")

	// 3D world to 2D screen location.
	bb.screenText = v2D.NewPov()
	bb.screenText.NewModel("uv").AddTex(font + "White").LoadFont(font).SetPhrase("Screen Text")
	bb.resize(s.W, s.H)
}

// Update is the regular engine callback.
func (bb *bbtag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	if in.Resized {
		bb.resize(s.W, s.H)
	}
	dt := in.Dt
	for press, _ := range in.Down {
		switch press {
		case "W":
			bb.cam.Move(0, 0, dt*-bb.run, bb.cam.Lookxz())
		case "S":
			bb.cam.Move(0, 0, dt*bb.run, bb.cam.Lookxz())
		case "Q":
			bb.cam.Move(dt*-bb.run, 0, 0, bb.cam.Lookxz())
		case "E":
			bb.cam.Move(dt*bb.run, 0, 0, bb.cam.Lookxz())
		case "A":
			bb.cam.AdjustYaw(dt * bb.spin)
		case "D":
			bb.cam.AdjustYaw(dt * -bb.spin)
		case "T":
			eng.Shutdown()
		}
	}

	// Use screen coordinates from world coordinates.
	if sx, sy := bb.cam.Screen(5, 2, -15, s.W, s.H); sx == -1 {
		bb.screenText.SetVisible(false)
	} else {
		bb.screenText.SetVisible(true)
		bb.screenText.SetLocation(float64(sx), float64(sy), 0)
	}
}
func (bb *bbtag) resize(ww, wh int) {
	bb.cam.SetPerspective(60, float64(ww)/float64(wh), 0.1, 50)
	bb.ui.SetOrthographic(0, float64(ww), 0, float64(wh), 0, 10)
}
