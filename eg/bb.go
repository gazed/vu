// Copyright © 2013-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"log"

	"github.com/gazed/vu"
)

// bb tests the engines handling of billboards, banners and fonts as
// well as combining multiple textures using shaders. Different 3D text
// types are created to show the possibilities for including text in games.
//
// CONTROLS:
//   WS    : move camera            : forward back
//   AD    : spin camera            : left right
//    T    : shut down
func bb() {
	defer catchErrors()
	if err := vu.Run(&bbtag{}); err != nil {
		log.Printf("bb: error starting engine %s", err)
	}
}

// Globally unique "tag" that encapsulates example specific data.
type bbtag struct {
	screenText *vu.Ent // Screen space text.
	scene      *vu.Ent // 3D scene
	ui         *vu.Ent // 2D scene
}

// Create is the startup asset creation.
func (bb *bbtag) Create(eng vu.Eng, s *vu.State) {
	eng.Set(vu.Title("Billboarding & Banners"), vu.Size(400, 100, 800, 600))
	eng.Set(vu.Color(0.3, 0.3, 0.3, 1))

	// 3D scene and camera.
	bb.scene = eng.AddScene()
	bb.scene.Cam().SetClip(0.1, 50).SetFov(60).SetAt(0.5, 2, 2.5)
	sun := bb.scene.AddPart().SetAt(0, 3, -3)
	sun.MakeLight().SetLightColor(0.4, 0.7, 0.9) // need light for gouraud shader.

	// // The floor model gives context to the labels.
	floor := bb.scene.AddPart().SetAt(0, 0, -10)
	floor.MakeModel("diffuse", "msh:floor", "mat:gray")

	// Show a billboarded texture.
	c4 := bb.scene.AddPart().SetAt(0.5, 2, -1).SetScale(0.25, 0.25, 0.25)
	c4.MakeModel("billboarded", "msh:billboard", "tex:core")
	c4.Clamp("core").SetAlpha(0.4)

	// Try banner text with the 3D scene perspective camera.
	font := "lucidiaSu22"
	banner := bb.scene.AddPart().SetScale(0.1, 0.1, 0.1).SetAt(-10, 2, -15)
	banner.MakeLabel("labeled", font).SetStr("Floating Text")

	// Try billboard banner text with the 3D scene perspective camera.
	banner = bb.scene.AddPart().SetScale(0.025, 0.025, 0.025).SetAt(-10, 1, -15)
	banner.MakeLabel("billboarded", font).SetStr("Billboard Text")

	// 2D scene and camera.
	bb.ui = eng.AddScene().SetUI()
	bb.ui.Cam().SetClip(0, 10)

	// 2D static location.
	banner = bb.ui.AddPart().SetAt(100, 100, 0)
	banner.MakeLabel("labeled", font).SetStr("Overlay Text")

	// 3D world to 2D screen location.
	bb.screenText = bb.ui.AddPart()
	bb.screenText.MakeLabel("labeled", font).SetStr("Screen Text")

	// 3D signed distance field font.
	sdf := bb.scene.AddPart().SetScale(0.1, 0.1, 0.1).SetAt(5, -1, -15)
	sdf.MakeLabel("sdf", "lucidiaSdf").SetStr("SDF").SetColor(1, 1, 0)

	// 2D signed distance field font.
	sdf2D := bb.ui.AddPart().SetAt(500, 100, 0).SetScale(0.5, 0.5, 1)
	sdf2D.MakeLabel("sdf", "lucidiaSdf").SetStr("SDF Overlay").SetColor(0, 1, 1)

	// 2D static location with text wrap and txt shader with color.
	banner = bb.ui.AddPart().SetAt(100, 200, 0)
	banner.MakeLabel("labeled", "lucidiaSu16").SetWrap(100).SetColor(1, 0, 1)
	banner.SetStr("A very long pink overlay string that should wrap over at least 3 lines.")
}

// Update is the regular engine callback.
func (bb *bbtag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	run := 10.0  // move so many cubes worth in one second.
	spin := 90.0 // spin so many degrees in one second.
	dt := in.Dt
	cam := bb.scene.Cam()
	for press := range in.Down {
		switch press {
		case vu.KW:
			cam.Move(0, 0, dt*-run, cam.Look)
		case vu.KS:
			cam.Move(0, 0, dt*run, cam.Look)
		case vu.KA:
			cam.SetYaw(cam.Yaw + spin*dt)
		case vu.KD:
			cam.SetYaw(cam.Yaw - spin*dt)
		case vu.KT:
			eng.Shutdown()
		}
	}

	// Use screen coordinates from world coordinates.
	if sx, sy := cam.Screen(5, 2, -15, s.W, s.H); sx == -1 {
		bb.screenText.Cull(true)
	} else {
		bb.screenText.Cull(false)
		bb.screenText.SetAt(float64(sx), float64(sy)-120, 0)
	}
}
