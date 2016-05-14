// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"log"

	"github.com/gazed/vu"
)

// bb tests the engines handling of billboards, banners and fonts as
// well as combining multiple textures using shaders. Different 3D text
// types are created to show the possibilities for including text in games.
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
	screenText vu.Pov    // Screen space text.
}

// Create is the startup asset creation.
func (bb *bbtag) Create(eng vu.Eng, s *vu.State) {
	top := eng.Root().NewPov()
	bb.cam = top.NewCam()
	bb.cam.SetLocation(0.5, 2, 2.5)
	sun := top.NewPov().SetLocation(0, 3, -3)
	sun.NewLight().SetColor(0.4, 0.7, 0.9)

	// Load the floor model.
	floor := top.NewPov()
	floor.NewModel("gouraud").LoadMesh("floor").LoadMat("gray")

	// Create a single image from multiple textures using a shader.
	c4 := top.NewPov().SetLocation(0.5, 2, -1).SetScale(0.25, 0.25, 0.25)
	model := c4.NewModel("spinball").LoadMesh("billboard")
	model.AddTex("core").AddTex("core").AddTex("halo").AddTex("halo")
	model.SetAlpha(0.4)

	// Try banner text with the 3D scene perspective camera.
	font := "lucidiaSu22"
	banner := top.NewPov().SetScale(0.1, 0.1, 0.1).SetLocation(-10, 3, -15)
	banner.NewModel("uv").AddTex(font + "White").LoadFont(font).SetPhrase("Floating Text")

	// Try billboard banner text with the 3D scene perspective camera.
	banner = top.NewPov().SetScale(0.025, 0.025, 0.025).SetLocation(-10, 2, -15)
	banner.NewModel("bb").AddTex(font + "White").LoadFont(font).SetPhrase("Billboard Text")

	// Banner text with an ortho overlay.
	v2D := eng.Root().NewPov()
	bb.ui = v2D.NewCam()
	bb.ui.SetUI()

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
	run := 10.0   // move so many cubes worth in one second.
	spin := 270.0 // spin so many degrees in one second.
	if in.Resized {
		bb.resize(s.W, s.H)
	}
	dt := in.Dt
	for press := range in.Down {
		switch press {
		case vu.KW:
			bb.cam.Move(0, 0, dt*-run, bb.cam.Lookxz())
		case vu.KS:
			bb.cam.Move(0, 0, dt*run, bb.cam.Lookxz())
		case vu.KQ:
			bb.cam.Move(dt*-run, 0, 0, bb.cam.Lookxz())
		case vu.KE:
			bb.cam.Move(dt*run, 0, 0, bb.cam.Lookxz())
		case vu.KA:
			bb.cam.AdjustYaw(dt * spin)
		case vu.KD:
			bb.cam.AdjustYaw(dt * -spin)
		case vu.KT:
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
