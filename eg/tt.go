// Copyright © 2015-2017 Galvanized Logic. All rights reserved.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"log"

	"github.com/gazed/vu"
)

// tt demonstrates rendering to a scene to a texture, and then
// displaying the scene on a quad. Background info at:
//   http://www.opengl-tutorial.org/intermediate-tutorials/tutorial-14-render-to-texture/
//   http://processors.wiki.ti.com/index.php/Render_to_Texture_with_OpenGL_ES
//   http://in2gpu.com/2014/09/24/render-to-texture-in-opengl/
//   http://www.lighthouse3d.com/tutorials/opengl_framebuffer_objects/
// This is another example of multi-pass rendering and can be used for
// generating live in-game portals.
//
// CONTROLS: none
//   WS    : spin frame             : left right
//   AD    : spin model             : left right
//   T     : shut down
func tt() {
	if err := vu.Run(&toTex{}); err != nil {
		log.Printf("tt: error starting engine %s", err)
	}
	defer catchErrors()
}

// Globally unique "tag" that encapsulates example specific data.
type toTex struct {
	monkey *vu.Ent // Allow user to spin monkey.
	frame  *vu.Ent // Allow user to spin frame.
}

// Create is the startup asset creation.
func (tt *toTex) Create(eng vu.Eng, s *vu.State) {
	eng.Set(vu.Title("Render to Texture"), vu.Size(400, 100, 800, 600))

	// scene that will render the blender monkey model to a square texture.
	scene0 := eng.AddScene().AsTex(true)
	scene0.Cam().SetClip(0.1, 50).SetFov(60)
	background := scene0.AddPart().SetAt(0, 0, -10).SetScale(100, 100, 1)
	background.MakeModel("uv", "msh:icon", "tex:wall")
	tt.monkey = scene0.AddPart().SetAt(0, 0, -5)
	tt.monkey.MakeModel("monkey", "msh:monkey", "mat:gray")

	// scene consisting of a single quad. The quad will display
	// the rendered texture from scene0. The texture image is flipped
	// so reorient it using flipped texture coordinates in flipboard.
	scene1 := eng.AddScene()
	scene1.Cam().SetClip(0.1, 50).SetFov(60)
	tt.frame = scene1.AddPart().SetAt(0, 0, -0.5).SetScale(0.25, 0.25, 0.25)
	model := tt.frame.MakeModel("uv", "msh:flipboard")
	model.SetTex(scene0) // use rendered texture from scene0.
}

// Update is the regular engine callback.
func (tt *toTex) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	spin := 270.0 // spin so many degrees in one second.
	dt := in.Dt
	for press := range in.Down {
		switch press {
		case vu.KW:
			tt.frame.Spin(0, dt*-spin, 0)
		case vu.KS:
			tt.frame.Spin(0, dt*+spin, 0)
		case vu.KA:
			tt.monkey.Spin(0, dt*-spin, 0)
		case vu.KD:
			tt.monkey.Spin(0, dt*+spin, 0)
		case vu.KT:
			eng.Shutdown()
		}
	}
}
