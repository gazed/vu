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
	dr := &drtag{}
	ww, wh := 1600.0, 900.0

	defer catchErrors()
	eng, err := vu.NewEngine(
		vu.Windowed(),
		vu.Title("Display Resize"),
		vu.Size(200, 200, int32(ww), int32(wh)),
		vu.Background(0.01, 0.01, 0.01, 1.0),
	)
	if err != nil {
		slog.Error("cr: engine start", "err", err)
		return
	}

	eng.ImportAssets("col2D.shd")

	// create root parts.
	ui := eng.AddScene(vu.Scene2D)

	// display area that matches the desired ratio.
	dr.display = ui.AddModel("shd:col2D", "msh:icon").SetColor(1, 1, 1, 0.1)
	dr.display.SetAt(ww*0.5, wh*0.5, 0).SetScale(ww, wh, 0)

	// box placed in top left corner.
	dr.box = ui.AddModel("shd:col2D", "msh:icon")
	dr.box.SetAt(40, 40, 0).SetScale(64, 64, 1).SetColor(0, 1, 0, 1)

	// listen for screen resizes.
	eng.SetResizeListener(dr)
	eng.Run(dr) // does not return while example is running.
}

// Globally unique "tag" that encapsulates example specific data.
type drtag struct {
	display *vu.Entity // Box representing the desired ratio.
	box     *vu.Entity // Box in top left corner.
}

// Resize is called by the engine when the window size changes.
func (dr *drtag) Resize(windowLeft, windowTop int32, windowWidth, windowHeight uint32) {
	ww, wh := float64(windowWidth), float64(windowHeight)
	cx, cy := ww*0.5, wh*0.5

	// fit a 16:9 box into the screen.
	ratio := ww / wh
	dw, dh := wh*16.0/9.0, wh // max height with black bars on sides
	if ratio < 16.0/9.0 {
		dw, dh = ww, ww*9.0/16.0 // max width with black bars top and bottom.
	}
	left, top := cx-dw/2, cy-dh/2

	// place a reference static box in the top left corner.
	dr.box.SetAt(left+40, top+40, 0)

	// Show the display area matching the desired ratio.
	dr.display.SetAt(cx, cy, 0) // center the image.
	dr.display.SetScale(dw, dh, 0)
}

// Update is the regular engine callback.
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
