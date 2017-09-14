// Copyright © 2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"log"

	"github.com/gazed/vu"
)

// sd demonstrates adding a sky-dome to a 3D scene.
// It also demonstrates flying around using the right mouse to change
// direction while moving with keys relative to the chosen direction.
//
// CONTROLS:
//   WS: move camera    : forward back in camera direction.
//   AD: move camera    : left right relative to camera direction.
//   Rm: spin camera    : look around.
func sd() {
	if err := vu.Run(&sdtag{}); err != nil {
		log.Printf("sd: error starting engine %s", err)
	}
	defer catchErrors()
}

// Globally unique "tag" that encapsulates example specific data.
type sdtag struct {
	scene  *vu.Ent // 3D scene
	mx, my int     // Previous mouse positions.
	pitch  float64 // Up-down look direction.
	yaw    float64 // Left-right look direction.
}

// Create handles startup asset creation.
func (sd *sdtag) Create(eng vu.Eng, s *vu.State) {
	eng.Set(vu.Title("Sky Dome"), vu.Size(400, 100, 800, 600))
	eng.Set(vu.Color(0.3, 0.3, 0.3, 1))
	sd.scene = eng.AddScene()
	sd.scene.Cam().SetClip(0.1, 50).SetFov(60)

	// Add a sky dome that is slightly lower than the camera
	// so that the bottom of the dome is rendered as well.
	sd.scene.AddSky().MakeModel("uv", "msh:dome", "tex:sky")

	// Add one block model to give the sky some context.
	block := sd.scene.AddPart().SetAt(0, 0, -10)
	block.MakeModel("nmap", "msh:box", "mat:tile", "tex:tile",
		"tex:tile_nrm", "tex:tile_spec")
}

// Update is the regular engine callback.
func (sd *sdtag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	xdiff, ydiff := in.Mx-sd.mx, in.My-sd.my // mouse move differences...
	sd.mx, sd.my = in.Mx, in.My              // ... from last mouse location.
	var lookSpeed = 0.25
	var runSpeed = 0.15
	cam := sd.scene.Cam()
	for press := range in.Down {
		switch press {
		case vu.KW:
			cam.Move(0, 0, -runSpeed, cam.Lookat())
		case vu.KS:
			cam.Move(0, 0, runSpeed, cam.Lookat())
		case vu.KA:
			cam.Move(-runSpeed, 0, 0, cam.Lookat())
		case vu.KD:
			cam.Move(runSpeed, 0, 0, cam.Lookat())
		case vu.KRm:
			if ydiff != 0 {
				sd.pitch = sd.limitPitch(sd.pitch + float64(ydiff)*lookSpeed)
			}
			if xdiff != 0 {
				sd.yaw += float64(-xdiff) * lookSpeed
			}
			cam.SetPitch(sd.pitch)
			cam.SetYaw(sd.yaw)
		}
	}
}

// limitPitch ensures that look up/down is limited to 90 degrees.
// This helps reduce confusion when looking around.
func (sd *sdtag) limitPitch(pitch float64) float64 {
	switch {
	case pitch > 90:
		return 90
	case pitch < -90:
		return -90
	}
	return pitch
}
