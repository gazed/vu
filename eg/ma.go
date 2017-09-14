// Copyright © 2014-2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/gazed/vu"
)

// ma, model animation, is an example of loading and animating a model using
// skeletel animation. Load any Inter-Quake-Model (IQM) models found in the model
// directory. This allows the example to function as a model previewer.
//
// CONTROLS:
//   WS    : move camera            : forward back
//   AD    : spin model             : left right
//   0-9   : select animation
//   Tab   : switch model
func ma() {
	if err := vu.Run(&matag{}); err != nil {
		log.Printf("ma: error starting engine %s", err)
	}
	defer catchErrors()
}

// Globally unique "tag" that encapsulates example specific data.
type matag struct {
	scene  *vu.Ent
	ui     *vu.Ent   // 2D user interface.
	label  *vu.Ent   // Animation information display.
	names  []string  // All loaded model names.
	models []*vu.Ent // All loaded models.
	model  *vu.Ent   // Currently selected model.
	index  int       // Index of currently selected model.
	run    float64   // Camera movement speed.
	spin   float64   // Camera spin speed.
}

// Create is the engine callback for initial asset creation.
func (ma *matag) Create(eng vu.Eng, s *vu.State) {
	eng.Set(vu.Title("Model Animation"), vu.Size(400, 100, 800, 600))
	eng.Set(vu.Color(0.15, 0.15, 0.15, 1))
	ma.scene = eng.AddScene()
	ma.scene.Cam().SetClip(0.1, 50).SetFov(60).SetAt(0, 3, 10)

	// load any available IQM models. The loaded model data is fed to
	// the animation capable shader "anim".
	for _, modelFile := range ma.modelFiles() {
		pov := ma.scene.AddPart()
		pov.SetScale(-1, 1, 1)
		if modelFile == "runner" {
			pov.SetScale(-3, 3, 3) // Runner is a bit small.
		}
		pov.Spin(-90, 0, 0) // Have the model face the camera.
		pov.Cull(true)      // Hide initially.

		// Most IQ* files are expected to be animated.
		// Use a "uv" shader to handle IQ* files without animations.
		pov.MakeActor("anim", modelFile)
		ma.models = append(ma.models, pov)
		ma.names = append(ma.names, modelFile)
	}
	ma.model = ma.models[ma.index] // should always have at least one.
	ma.model.Cull(false)

	// Create a banner to show the model name.
	ma.ui = eng.AddScene().SetUI()
	ma.ui.Cam().SetClip(0, 10)
	ma.label = ma.ui.AddPart().SetAt(10, 5, 0)
	ma.label.MakeLabel("txt", "lucidiaSu22")
}

// Update is the recurring callback to update state based on user actions.
func (ma *matag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	run := 10.0 // move so many units worth in one second.
	dt := in.Dt
	cam := ma.scene.Cam()
	for press, down := range in.Down {
		switch press {
		case vu.KW:
			cam.Move(0, 0, dt*run, cam.Look)
		case vu.KS:
			cam.Move(0, 0, dt*-run, cam.Look)
		case vu.KA:
			ma.model.Spin(0, 0, 5)
		case vu.KD:
			ma.model.Spin(0, 0, -5)
		case vu.KTab:
			if down == 1 {

				// switch to the next loaded model.
				ma.model.Cull(true)
				ma.index = ma.index + 1
				if ma.index >= len(ma.models) {
					ma.index = 0
				}
				ma.model = ma.models[ma.index]
				ma.model.Cull(false)
			}
		case vu.K0, vu.K1, vu.K2, vu.K3, vu.K4, vu.K5, vu.K6, vu.K7, vu.K8, vu.K9:
			if down == 1 {
				ma.playAnimation(press)
			}
		}
	}
	ma.showAction()
}

// playAnimation chooses an available animation.
// Animations that are not available are ignored.
func (ma *matag) playAnimation(keyCode int) {
	var actions = map[int]int{
		vu.K0: 0,
		vu.K1: 1,
		vu.K2: 2,
		vu.K3: 3,
		vu.K4: 4,
		vu.K5: 5,
		vu.K6: 6,
		vu.K7: 7,
		vu.K8: 8,
		vu.K9: 9,
	}
	if action, ok := actions[keyCode]; ok {
		ma.model.Animate(action, 0)
	}
}

// showAction updates the information text, the animation
// file, animation sequence name, and the frame numbers.
func (ma *matag) showAction() {
	if names := ma.model.Actions(); len(names) > 0 {
		index, frame, maxFrames := ma.model.Action()
		name := names[index]
		stats := fmt.Sprintf("[%d] %s %d/%d", index, name, frame, maxFrames)
		ma.label.Typeset(ma.names[ma.index] + ":" + stats)
	}
}

// iqmodel groups the 3D assets with the file name of the model file.
type iqmodel struct {
	title string  // IQ file name.
	tr    *vu.Ent // loaded IQ 3D model.
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
