// Copyright © 2013-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

// FUTURE: upgrade shader to PCSS ie:
//   http://developer.download.nvidia.com/whitepapers/2008/PCSS_Integration.pdf

import (
	"log"

	"github.com/gazed/vu"
	"github.com/gazed/vu/math/lin"
)

// sm tests the engines handling of shadow maps and multipass rendering.
// These are hard shadows. More work has to be done to soften and blur
// shadows the further away they are from their source.
//
// CONTROLS:
//   WASD  : move light position    : forward left back right
func sm() {
	defer catchErrors()
	if err := vu.Run(&smtag{}); err != nil {
		log.Printf("sm: error starting engine %s", err)
	}
}

// Globally unique "tag" that encapsulates example specific data.
type smtag struct {
	scene  *vu.Ent
	sun    *vu.Ent
	cube   *vu.Ent
	sphere *vu.Ent
}

// Create is the startup asset creation.
func (sm *smtag) Create(eng vu.Eng, s *vu.State) {
	eng.Set(vu.Title("Shadow Map"), vu.Size(400, 100, 800, 600))

	// create a scene that will render a shadow map.
	sm.scene = eng.AddScene().SetShadows()
	sm.scene.Cam().SetClip(0.1, 50).SetFov(60).SetAt(0, 0, 10)

	// need a light for shadows.
	sm.sun = sm.scene.MakeLight(vu.DirectionalLight).SetLightColor(0.8, 0.8, 0.8)

	// create a few objects that cast shadows.
	sm.cube = sm.scene.AddPart().SetAt(-1, -1, -4)
	sm.cube.MakeModel("phong", "msh:box", "mat:gray")
	sm.cube.Spin(45, 45, 0)
	sm.sphere = sm.scene.AddPart().SetAt(1, 1, -4)
	sm.sphere.MakeModel("phong", "msh:sphere", "mat:red")

	// create a ground block to show shadows.
	ground := sm.scene.AddPart().SetAt(0, 0, -20).SetScale(50, 50, 5)
	ground.MakeModel("showShadow", "msh:box", "mat:gray", "tex:tile")
}

// Update is the regular engine callback.
func (sm *smtag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
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
