// SPDX-FileCopyrightText : © 2025 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"log/slog"
	"time"

	"github.com/gazed/vu"
)

// ax checks that engine sound integration works.
// Shows how to use sound as part of the game engine.
//
// CONTROLS: NA
//   - A    : play sound
func ax() {
	ax := &axtag{}
	ww, wh := 1600.0, 900.0

	defer catchErrors()
	eng, err := vu.NewEngine(
		vu.Windowed(),
		vu.Title("Engine Audio"),
		vu.Size(200, 200, int32(ww), int32(wh)),
		vu.Background(0.01, 0.01, 0.01, 1.0),
	)
	if err != nil {
		slog.Error("cr: engine start", "err", err)
		return
	}

	// listen for screen resizes.
	eng.SetResizeListener(ax)

	// Run will call Load once and then call Update each engine tick.
	eng.Run(ax, ax) // does not return while example is running.
}

// Globally unique "tag" that encapsulates example specific data.
type axtag struct {
	model   *vu.Entity // Box representing the desired ratio.
	soundID *vu.Entity // loaded sound asset ready to play
}

// Load is the one time startup engine callback to create initial assets.
func (ax *axtag) Load(eng *vu.Engine) error {

	// load the model and sound assets.
	eng.ImportAssets("col3D.shd", "bloop.wav")

	// create a sound entity from the loaded asset.
	ax.soundID = eng.AddSound("bloop")

	// create root parts.
	scene := eng.AddScene(vu.Scene3D)

	// model as the sound focus.
	ax.model = scene.AddModel("shd:col3D", "msh:cube").SetColor(1, 1, 1, 0.1)
	ax.model.SetAt(0, 0, -4).SetScale(2, 2, 0)
	ax.model.SetListener() // set the sound listener to this models location.
	return nil
}

// Resize is called by the engine when the window size changes.
func (ax *axtag) Resize(windowLeft, windowTop int32, windowWidth, windowHeight uint32) {}

// Update is the ongoing engine callback.
func (ax *axtag) Update(eng *vu.Engine, in *vu.Input, delta time.Duration) {

	// react to one time press events.
	for press := range in.Pressed {
		switch press {
		case vu.KQ: // quit if Q is pressed
			eng.Shutdown()
			return
		case vu.KA: // play sound at the models location.
			ax.model.PlaySound(eng, ax.soundID)
			return
		}
	}
}
