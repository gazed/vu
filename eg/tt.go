// Copyright © 2015-2016 Galvanized Logic. All rights reserved.
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
//   ws    : spin frame             : left right
//   AD    : spin model             : left right
//   T     : shut down
func tt() {
	tt := &toTex{}
	if err := vu.New(tt, "Render to Texture", 400, 100, 800, 600); err != nil {
		log.Printf("tt: error starting engine %s", err)
	}
	defer catchErrors()
}

// Globally unique "tag" that encapsulates example specific data.
type toTex struct {
	cam0   *vu.Camera // Camera for rendering monkey to texture scene.
	monkey *vu.Pov    // Allow user to spin monkey.
	cam1   *vu.Camera // Camera for rendering texture frame.
	frame  *vu.Pov    // Allow user to spin frame.
}

// Create is the startup asset creation.
func (tt *toTex) Create(eng vu.Eng, s *vu.State) {

	// create a scene that will render the blender monkey model to a texture.
	scene0 := eng.Root().NewPov()
	tt.cam0 = scene0.NewCam()
	scene0.NewLayer() // render scene to texture.
	background := scene0.NewPov().SetAt(0, 0, -10).SetScale(100, 100, 1)
	background.NewModel("uv", "msh:icon", "tex:wall")
	tt.monkey = scene0.NewPov().SetAt(0, 0, -5)
	tt.monkey.NewModel("monkey", "msh:monkey", "mat:gray")

	// create a scene consisting of a single quad. The quad will display
	// the rendered texture from scene0. The texture image is flipped
	// so reorient it using flipped texture coordinates in flipboard.
	scene1 := eng.Root().NewPov()
	tt.cam1 = scene1.NewCam()
	tt.frame = scene1.NewPov().SetAt(0, 0, -0.5).SetScale(0.25, 0.25, 0.25)
	model := tt.frame.NewModel("uv", "msh:flipboard")
	model.UseLayer(scene0.Layer()) // use rendered texture from scene0.

	// set camera perspectives and default background color.
	tt.resize(s.W, s.H)
}

// Update is the regular engine callback.
func (tt *toTex) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	spin := 270.0 // spin so many degrees in one second.
	if in.Resized {
		tt.resize(s.W, s.H)
	}
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
func (tt *toTex) resize(ww, wh int) {
	tt.cam0.SetPerspective(60, float64(1024)/float64(1024), 0.1, 50) // Image size.
	tt.cam1.SetPerspective(60, float64(ww)/float64(wh), 0.1, 50)     // Screen size.
}
