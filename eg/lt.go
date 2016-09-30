// Copyright © 2015-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"log"

	"github.com/gazed/vu"
	"github.com/gazed/vu/math/lin"
)

// lt tests the engines handling of some of the engine lighting shaders.
// It also checks the conversion of light position and normal vectors
// needed for proper lighting.
//
// Note the use of the box.obj model that needs 24 verticies to get
// proper lighting on each face. Also note how many more verticies are
// necessary for the sphere.obj model.
//
// CONTROLS:
//   WASD  : move the light position: forward left back right
//   ZX    : move light position    : up down
//   LaRa  : spin the cube          : left right.
func lt() {
	lt := &lttag{}
	if err := vu.New(lt, "Lighting", 400, 100, 800, 600); err != nil {
		log.Printf("lt: error starting engine %s", err)
	}
	defer catchErrors()
}

// Globally unique "tag" that encapsulates example specific data.
type lttag struct {
	cam3D *vu.Camera // 3D main scene camera.
	sun   *vu.Pov    // Light node in Pov hierarchy.
	box   *vu.Pov    // Normal mapped box.
}

// Create is the engine callback for initial asset creation.
func (lt *lttag) Create(eng vu.Eng, s *vu.State) {
	top := eng.Root().NewPov()
	lt.cam3D = top.NewCam()
	lt.cam3D.SetAt(0.5, 2, 0.5)
	lt.sun = top.NewPov().SetAt(0, 2.5, -1.75).SetScale(0.05, 0.05, 0.05)
	lt.sun.NewLight().SetColor(0.5, 0.5, 0.5)

	// Model at the light position.
	lt.sun.NewModel("solid", "msh:sphere", "mat:red")

	// Create solid spheres to test the lighting shaders.
	c4 := top.NewPov().SetAt(-0.5, 2, -2).SetScale(0.25, 0.25, 0.25)
	c5 := top.NewPov().SetAt(0.5, 2, -2).SetScale(0.25, 0.25, 0.25)
	c6 := top.NewPov().SetAt(1.5, 2, -2).SetScale(0.25, 0.25, 0.25)
	c4.NewModel("diffuse", "msh:sphere", "mat:blue")
	c5.NewModel("gouraud", "msh:sphere", "mat:red")
	c6.NewModel("phong", "msh:sphere", "mat:blue")

	// Angle a large flat box with normal map lighting behind the spheres.
	lt.box = top.NewPov().SetAt(0, 2, -10)
	lt.box.SetScale(5, 5, 5).Spin(45, 45, 0)
	lt.box.NewModel("nmap", "msh:box", "mat:tile", "tex:tile", "tex:tile_nrm", "tex:tile_spec")
	lt.resize(s.W, s.H)
}

// Update is the regular engine callback.
func (lt *lttag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	run := 10.0 // move so many units worth in one second.
	if in.Resized {
		lt.resize(s.W, s.H)
	}
	// move the light.
	dt := in.Dt
	speed := run * dt * 0.5
	for press := range in.Down {
		switch press {
		case vu.KW:
			lt.sun.Move(0, 0, -speed, lin.QI) // forward
		case vu.KS:
			lt.sun.Move(0, 0, speed, lin.QI) // back
		case vu.KA:
			lt.sun.Move(-speed, 0, 0, lin.QI) // left
		case vu.KD:
			lt.sun.Move(speed, 0, 0, lin.QI) // right
		case vu.KZ:
			lt.sun.Move(0, speed, 0, lin.QI) // up
		case vu.KX:
			lt.sun.Move(0, -speed, 0, lin.QI) // down
		case vu.KLa:
			lt.box.Spin(0, speed*10, 0)
		case vu.KRa:
			lt.box.Spin(0, -speed*10, 0)
		}
	}
}
func (lt *lttag) resize(ww, wh int) {
	lt.cam3D.SetPerspective(60, float64(ww)/float64(wh), 0.1, 50)
}

// Design notes and references.
//
// General background for lighting and normal mapped lighting:
//   http://http.developer.nvidia.com/CgTutorial/cg_tutorial_chapter08.html
//   http://ogldev.atspace.co.uk/www/tutorial26/tutorial26.html
// Examples with GLSL version 330 shader code.
//   http://www.opengl-tutorial.org/intermediate-tutorials/tutorial-13-normal-mapping/
//   http://www.keithlantz.net/2011/10/tangent-space-normal-mapping-with-glsl/
// More examples and explanations:
//   http://fabiensanglard.net/bumpMapping/index.php
//   http://www.ozone3d.net/tutorials/bump_mapping.php
//   http://www.swiftless.com/tutorials/glsl/8_bump_mapping.html
//   http://sunandblackcat.com/tipFullView.php?l=eng&topicid=30&topic=Phong-Lighting
// Discussion on blending
//   http://blog.selfshadow.com/publications/blending-in-detail/
// Totally shader based so models don't need to load tangents.
//   http://www.thetenthplanet.de/archives/1180 ***nmap shader based on this.
// Good explanation of eyespace.
//   http://pyopengl.sourceforge.net/context/tutorials/shader_6.html
