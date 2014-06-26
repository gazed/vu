// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.
//
// Note that models/mrfixit.iqm, models/Head.png, models/Body.png, from
// the IQM development kit, are subject to the following license:
//       Mr. Fixit conceptualized, modeled and animated by geartrooper
//                 a.k.a. John Siar ironantknight@gmail.com
//       Mr. Fixit textures, normal maps and displacement by Acord
//                 a.k.a Anthony Cord tonycord@gmail.com  mask textures by Nieb
//       Mr. Fixit is released under the CC-BY-NC (Creative Commons Attribution Non-Commercial)
//       license. See http://creativecommons.org/licenses/by-nc/3.0/legalcode for more info.
//       Please contact Geartrooper for permission for uses not within the scope of this license.

package main

import (
	"log"
	"vu"
)

// ma, Model Animation, is an example of loading and animating a model using
// skeletel animation. It is based on the example data provided in the IQM
// Development kit from http://sauerbraten.org/iqm.
func ma() {
	ma := &matag{}
	var err error
	if ma.eng, err = vu.New("Model Animation", 400, 100, 800, 600); err != nil {
		log.Fatal("ma: error intitializing engine %s", err)
	}
	ma.eng.SetDirector(ma) // override user input handling.
	if err = ma.eng.Verify(); err != nil {
		log.Fatalf("ma: error initializing model :: %s", err)
	}
	defer ma.eng.Shutdown()
	defer catchErrors()
	ma.eng.Action()
}

// Globally unique "tag" that encapsulates example specific data.
type matag struct {
	eng   vu.Engine
	scene vu.Scene
	model vu.Part
	run   float64 // Camera movement speed.
	spin  float64 // Camera spin speed.
}

// Create is the one time Engine intialization callback.
func (ma *matag) Create(eng vu.Engine) {
	ma.run = 10   // move so many cubes worth in one second.
	ma.spin = 270 // spin so many degrees in one second.
	ma.scene = eng.AddScene(vu.VP)
	ma.scene.SetPerspective(60, float64(800)/float64(600), 0.1, 50)
	ma.scene.SetLocation(0, 4, 10)

	// load the IQM model. The loaded model data is fed to the animation
	// capable shader.
	ma.model = ma.scene.AddPart()
	ma.model.SetRole("anim").SetModelData("mrfixit")
	ma.model.SetScale(-1, 1, 1) // Mirror around Y.
	ma.model.Spin(-90, 90, 0)   // Have the model face the camera.

	// Have a lighter default background.
	eng.Color(0.15, 0.15, 0.15, 1)
}

// Update is the recurring callback to update state based on user actions.
func (ma *matag) Update(in *vu.Input) {
	if in.Resized {
		ma.resize()
	}
	dt := in.Dt
	for press, _ := range in.Down {
		switch press {
		case "W":
			ma.scene.Move(0, 0, dt*ma.run)
		case "S":
			ma.scene.Move(0, 0, dt*-ma.run)
		case "A":
			ma.model.Spin(0, 5, 0)
		case "D":
			ma.model.Spin(0, -5, 0)
		}
	}
}

// resize handles user screen/window changes.
func (ma *matag) resize() {
	x, y, width, height := ma.eng.Size()
	ma.eng.Resize(x, y, width, height)
	ma.scene.SetPerspective(60, float64(width)/float64(height), 0.1, 50)
}
