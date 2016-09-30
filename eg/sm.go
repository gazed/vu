// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"log"

	"github.com/gazed/vu"
	"github.com/gazed/vu/math/lin"
)

// sm, shadow map, tests the engines handling of shadow maps
// and multipass rendering. These are hard shadows. More work has to
// be done to soften and blur shadows the further away they are from
// their source.
// FUTURE: upgrade shader to PCSS ie:
//   http://developer.download.nvidia.com/whitepapers/2008/PCSS_Integration.pdf
//
// CONTROLS:
//   WASD  : move light position    : up left down right
func sm() {
	sm := &smtag{}
	if err := vu.New(sm, "Shadow Map", 400, 100, 800, 600); err != nil {
		log.Printf("sm: error starting engine %s", err)
	}
	defer catchErrors()
}

// Globally unique "tag" that encapsulates example specific data.
type smtag struct {
	sun    *vu.Pov
	cube   *vu.Pov
	sphere *vu.Pov
	cam    *vu.Camera
}

// Create is the startup asset creation.
func (sm *smtag) Create(eng vu.Eng, s *vu.State) {
	scene := eng.Root().NewPov()
	sm.cam = scene.NewCam()

	// need a light for shadows.
	sm.sun = scene.NewPov().SetAt(0, 0, 0)
	sm.sun.NewLight().SetColor(0.8, 0.8, 0.8)

	// create a scene that will render a shadow map.
	sm.cam = scene.NewCam()
	sm.cam.SetAt(0, 0, 10)
	sm.cam.SetPerspective(60, float64(s.W)/float64(s.H), 0.1, 50)

	// create a few objects that cast shadows.
	sm.cube = scene.NewPov().SetAt(-1, -1, -4)
	sm.cube.NewModel("gouraud", "msh:box", "mat:gray").Set(vu.CastShadow)
	sm.cube.Spin(45, 45, 0)
	sm.sphere = scene.NewPov().SetAt(1, 1, -4)
	sm.sphere.NewModel("gouraud", "msh:sphere", "mat:red").Set(vu.CastShadow)

	// create a ground block to show shadows.
	ground := scene.NewPov().SetAt(0, 0, -20).SetScale(50, 50, 5)
	ground.NewModel("shadow", "msh:box", "mat:gray", "tex:tile").Set(vu.HasShadows)
}

// Update is the regular engine callback.
func (sm *smtag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	if in.Resized {
		sm.resize(s.W, s.H)
	}
	dt := in.Dt
	rate := 10.0
	for press := range in.Down {
		switch press {
		case vu.KA:
			sm.sun.Move(dt*rate, 0, 0, lin.QI) // left
		case vu.KD:
			sm.sun.Move(-dt*rate, 0, 0, lin.QI) // right
		case vu.KW:
			sm.sun.Move(0, -dt*rate, 0, lin.QI) // shadow up
		case vu.KS:
			sm.sun.Move(0, dt*rate, 0, lin.QI) // shadow down
		}
	}
}
func (sm *smtag) resize(ww, wh int) {
	sm.cam.SetPerspective(60, float64(ww)/float64(wh), 0.1, 50)
}
