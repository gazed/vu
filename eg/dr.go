// Copyright © 2024 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"log/slog"
	"time"

	"github.com/gazed/vu"
)

// dr demonstrates display resize handling for a fix aspect ratio.
//
// CONTROLS: Resize the window and see how the application responds.
func dr() {
	defer catchErrors()
	dr := &drtag{ww: 1600.0, wh: 900.0}
	eng, err := vu.NewEngine(
		vu.Windowed(),
		vu.Title("Display Resize"),
		vu.Size(200, 200, int32(dr.ww), int32(dr.wh)),
		vu.Background(0.01, 0.01, 0.01, 1.0),
	)
	if err != nil {
		slog.Error("cr: engine start", "err", err)
		return
	}

	// listen for screen resizes.
	eng.SetResizeListener(dr)

	// Run will call Load once and then call Update each engine tick.
	eng.Run(dr, dr) // does not return while example is running.
}

// Globally unique "tag" that encapsulates example specific data.
type drtag struct {
	ww, wh  float64
	display *vu.Entity // Box representing the desired ratio.
	box     *vu.Entity // Box in top left corner.
}

// Load is the one time startup engine callback to create initial assets.
func (dr *drtag) Load(eng *vu.Engine) error {
	eng.ImportAssets("col2D.shd")

	// create root parts.
	ui := eng.AddScene(vu.Scene2D)

	// display area that matches the desired ratio.
	dr.display = ui.AddModel("shd:col2D", "msh:icon").SetColor(1, 1, 1, 0.1)
	dr.display.SetAt(dr.ww*0.5, dr.wh*0.5, 0).SetScale(dr.ww, dr.wh, 0)

	// box placed in top left corner.
	dr.box = ui.AddModel("shd:col2D", "msh:icon")
	dr.box.SetAt(40, 40, 0).SetScale(64, 64, 1).SetColor(0, 1, 0, 1)
	return nil
}

// Resize is called by the engine when the window size changes.
func (dr *drtag) Resize(windowLeft, windowTop int32, windowWidth, windowHeight uint32) {
	ww, wh := float64(windowWidth), float64(windowHeight)
	cx, cy := ww*0.5, wh*0.5
	fixedRatio := 16.0 / 9.0

	// fit a 16:9 box into the screen.
	ratio := ww / wh
	dw, dh := wh*fixedRatio, wh // max height with black bars on sides
	if ratio < fixedRatio {
		dw, dh = ww, ww*(1.0/fixedRatio) // max width with black bars top and bottom.
	}
	left, top := cx-dw/2, cy-dh/2

	// place a reference static box in the top left corner.
	dr.box.SetAt(left+40, top+40, 0)

	// Show the display area matching the desired ratio.
	dr.display.SetAt(cx, cy, 0) // center the image.
	dr.display.SetScale(dw, dh, 0)
}

// Update is the ongoing engine callback.
func (dr *drtag) Update(eng *vu.Engine, in *vu.Input, delta time.Duration) {

	// react to one time press events.
	for press := range in.Pressed {
		switch press {
		case vu.KQ: // quit if Q is pressed
			eng.Shutdown()
			return
		}
	}
}
