// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.
//
// Note that models/mrfixit.iqm, models/mrfixit0.png, models/mrfixit1.png are from
// the IQM development kit and are subject to the following license:
//       Mr. Fixit conceptualized, modeled and animated by geartrooper
//                 a.k.a. John Siar ironantknight@gmail.com
//       Mr. Fixit textures, normal maps and displacement by Acord
//                 a.k.a Anthony Cord tonycord@gmail.com  mask textures by Nieb
//       Mr. Fixit is released under the CC-BY-NC (Creative Commons Attribution Non-Commercial)
//       license. See http://creativecommons.org/licenses/by-nc/3.0/legalcode for more info.
//       Please contact Geartrooper for permission for uses not within the scope of this license.

package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/gazed/vu"
)

// ma, model animation, is an example of loading and animating a model using
// skeletel animation. It is based on the example data provided in the IQM
// Development kit from http://sauerbraten.org/iqm.
func ma() {
	ma := &matag{}
	var err error
	if ma.eng, err = vu.New("Model Animation", 400, 100, 800, 600); err != nil {
		log.Fatal("ma: error intitializing engine %s", err)
	}
	ma.eng.SetDirector(ma) // get user input through Director.Update()
	ma.create()            // create initial assests.
	if err = ma.eng.Verify(); err != nil {
		log.Fatalf("ma: error initializing model :: %s", err)
	}
	defer ma.eng.Shutdown()
	defer catchErrors()
	ma.eng.Action()
}

// Globally unique "tag" that encapsulates example specific data.
type matag struct {
	eng     vu.Engine
	scene   vu.Scene
	cam     vu.Camera
	models  []*iqmodel
	model   *iqmodel
	title   vu.Part
	current int
	run     float64 // Camera movement speed.
	spin    float64 // Camera spin speed.
}

// create is the startup asset creation.
func (ma *matag) create() {
	ma.run = 10   // move so many cubes worth in one second.
	ma.spin = 270 // spin so many degrees in one second.
	ma.scene = ma.eng.AddScene(vu.VP)
	ma.cam = ma.scene.Cam()
	ma.cam.SetPerspective(60, float64(800)/float64(600), 0.1, 50)
	ma.cam.SetLocation(0, 4, 10)

	// load any available IQM/E models. The loaded model data is fed to
	// the animation capable shader "anim".
	ma.models = []*iqmodel{}
	for _, modelFile := range ma.modelFiles() {
		m := ma.scene.AddPart()

		// Most IQ* files are expected to be animated.
		// Use a "uv" shader to handle IQ* files without animations.
		m.SetRole("anim").SetAnimation(modelFile)
		m.Role().SetCullOff()
		m.SetScale(-1, 1, 1) // Mirror around Y.
		m.Spin(-90, 0, 90)   // Have the model face the camera.
		m.SetVisible(false)
		ma.models = append(ma.models, &iqmodel{modelFile, m})
	}
	ma.model = ma.models[ma.current] // should always have at least one.
	ma.model.mod.SetVisible(true)

	// Create a banner to show the model name.
	over := ma.eng.AddScene(vu.VO)
	over.Set2D()
	_, _, w, h := ma.eng.Size()
	over.Cam().SetOrthographic(0, float64(w), 0, float64(h), 0, 10)
	ma.title = over.AddPart()
	ma.title.SetRole("uv").AddTex("weblySleek22White").SetFont("weblySleek22")
	ma.showMovement(0)

	// Have a lighter default background.
	ma.eng.Color(0.15, 0.15, 0.15, 1)
}

// iqmodel groups the 3D assets with the file name of the model file.
type iqmodel struct {
	title string  // IQ file name.
	mod   vu.Part // loaded IQ 3D model.
}

// modelFiles returns the names of the IQE/IQM files in the models directory.
// Only unique base names are returned.
func (ma *matag) modelFiles() []string {
	uniqueNames := map[string]bool{}
	models := []string{}
	files, _ := ioutil.ReadDir("models")
	for _, f := range files {
		name := f.Name()
		if strings.Contains(name, ".iq") && name[0] != '.' {
			base := name[:len(name)-4]
			if _, ok := uniqueNames[base]; !ok {
				uniqueNames[base] = true
				models = append(models, base)
			}
		}
	}
	return models
}

// Update is the recurring callback to update state based on user actions.
func (ma *matag) Update(in *vu.Input) {
	if in.Resized {
		ma.resize()
	}
	dt := in.Dt
	for press, down := range in.Down {
		switch press {
		case "W":
			ma.cam.Move(0, 0, dt*ma.run)
		case "S":
			ma.cam.Move(0, 0, dt*-ma.run)
		case "A":
			ma.model.mod.Spin(0, 0, 5)
		case "D":
			ma.model.mod.Spin(0, 0, -5)
		case "Tab":
			if down == 1 {

				// switch to the next loaded model.
				ma.model.mod.SetVisible(false)
				ma.current = ma.current + 1
				if ma.current >= len(ma.models) {
					ma.current = 0
				}
				ma.model = ma.models[ma.current]
				ma.model.mod.SetVisible(true)
				ma.showMovement(0)
			}
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			if down == 1 {
				ma.playAnimation(press)
			}
		}
	}
}

// playAnimation chooses an available animation. Animations that are not
// available are ignored.
func (ma *matag) playAnimation(d09 string) {
	digit, _ := strconv.Atoi(d09)
	if ma.model.mod.Role().PlayMovement(digit, nil) {
		ma.showMovement(digit)
	} else {
		ma.showMovement(0)
	}
}

func (ma *matag) showMovement(movement int) {
	if len(ma.model.mod.Role().Movements()) > 0 {
		label := ma.model.mod.Role().Movements()[movement]
		ma.title.Role().SetPhrase(ma.model.title + ":" + label)
	} else {
		ma.title.Role().SetPhrase(ma.model.title)
	}
}

// resize handles user screen/window changes.
func (ma *matag) resize() {
	x, y, width, height := ma.eng.Size()
	ma.eng.Resize(x, y, width, height)
	ma.cam.SetPerspective(60, float64(width)/float64(height), 0.1, 50)
}
