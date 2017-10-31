// Copyright © 2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.
// images/robot.png Copyright © 2017 Julien Ruest.

package main

import (
	"fmt"
	"log"
	"math"

	"github.com/gazed/vu"
)

// sz demonstrates screen resize handling both by picking the best fit screen
// ratio and by showing the effect of resizing on a background image.
// The background image uses a core area, which must always be shown,
// and non-core areas which are used to fill different screen ratios.
//
// The point of this example is to get an understanding which common
// ratios an application can support and the calculations needed to
// reposition elements after a resize.
//
// CONTROLS: Resize the window and see how the application responds.
func sz() {
	defer catchErrors()
	if err := vu.Run(&sztag{}); err != nil {
		log.Printf("sz: error starting engine %s", err)
	}
}

// Globally unique "tag" that encapsulates example specific data.
type sztag struct {
	scene  *vu.Ent // 3D scene
	gui    *vu.Ent // 2D scene
	label  *vu.Ent // Ratio description.
	ratio  *vu.Ent // Ratio display object.
	bg     *vu.Ent // Background image.
	button *vu.Ent // Place a button in different ratios.
}

// Create handles startup asset creation.
func (sz *sztag) Create(eng vu.Eng, s *vu.State) {
	eng.Set(vu.Title("Screen Resize"), vu.Size(400, 100, 800, 600))
	eng.Set(vu.Color(0.3, 0.3, 0.3, 1))
	sz.scene = eng.AddScene()
	sz.scene.Cam().SetClip(0.1, 50).SetFov(60)

	// 2D scene.
	sz.gui = eng.AddScene().SetUI()
	sz.gui.Cam().SetClip(0, 10)

	// background image has a core 16x10 ratio.
	sz.bg = sz.gui.AddPart().MakeModel("uv", "msh:icon", "tex:robot7")

	// button will be placed in top left corner.
	sz.button = sz.gui.AddPart().MakeModel("uv", "msh:icon", "tex:core")

	// Show the best fit ratio.
	sz.ratio = sz.gui.AddPart()
	sz.ratio.MakeModel("solid", "msh:square").SetDraw(vu.Lines)
	sz.label = sz.gui.AddPart().MakeLabel("txt", "lucidiaSu22").Typeset("None")
}

// Update is the regular engine callback.
func (sz *sztag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	if in.Resized {
		ww, wh := float64(s.W), float64(s.H)
		cx, cy, ratio := ww*0.5, wh*0.5, ww/wh

		// Show the screen ratio that best fits the window.
		best := -1
		least := math.MaxFloat32
		screenRatios := []float64{4.0 / 3.0, 16.0 / 10.0, 16.0 / 9.0}
		for cnt, r := range screenRatios {
			if math.Abs(ratio-r) < least {
				least = math.Abs(ratio - r)
				best = cnt
			}
		}
		w, h, fit := ww, wh, ""
		sz.label.SetAt(cx, cy, 0)
		switch best {
		case 0:
			fit = "4:3"
			w, h = wh*4/3, wh
			if ratio < 4.0/3.0 {
				w, h = ww, ww*3/4
			}
			sz.ratio.SetColor(1, 0, 0)
		case 1:
			fit = "16:10"
			w, h = wh*1.6, wh
			if ratio < 16.0/10.0 {
				w, h = ww, ww*10/16
			}
			sz.ratio.SetColor(0, 1, 0)
		case 2:
			fit = "16:9"
			w, h = wh*16/9, wh
			if ratio < 16.0/9.0 {
				w, h = ww, ww*9/16
			}
			sz.ratio.SetColor(1, 0, 1)
		}
		x, y := int(cx-w/2), int(cy-h/2)
		sz.ratio.SetAt(cx, cy, 0).SetScale(w-1, h-1, 0)
		txt := fmt.Sprintf("%s\n%d %d %d %d", fit, x, y, int(w), int(h))
		sz.label.Typeset(txt)

		// place a static sized button in the top right corner.
		sz.button.SetAt(float64(x)+40, (h+float64(y))-40, 0).SetScale(64, 64, 1)

		// Display the background image. The image is 1024x640 (16:10)
		// with a core image of 800x500. Scale the image so that the core
		// image matches one of the screen ratio edges.
		sz.bg.SetAt(cx, cy, 0) // always center the image.
		switch {
		case w/h > 1.6: // greater than image ratio...
			scale := 640.0 * (h / 500.0) // ...scale based on height.
			sz.bg.SetScale(scale*1.6, scale, 0)
		case w/h <= 1.6: // less than image ratio...
			scale := 1024.0 * (w / 800.0) // ...scale based on width.
			sz.bg.SetScale(scale, scale*10/16, 0)
		}
	}
}
