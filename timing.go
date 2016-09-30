// Copyright © 2015-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// timing.go
// FUTURE : lots can be done here, but it needs to be balanced against engine
//          code clutter. Variations in timing are expected to be mostly
//          influenced by device capability, ie: mobile devices are less capable
//          than consoles or desktops.
//      o Rendering times are not currently captured from the machine goroutine.
//      o Loading times are not currently captured from the loader goroutine.

import (
	"fmt"
	"time"
)

// Timing is used to collect main processing loop numbers while the
// the application loop is active. The numbers are reset each update.
// Applications are expected to track and smooth these per-update
// values over a number of updates.
//
// Timing gives a rough idea. Expect things to go slower the more
// models, draw calls, and the greater number of verticies.
//
// FPS = Renders/Elapsed. This is how many render requests were generated.
// Actual number of renders likely matches the monitor refresh rate
// which is 60/sec for most flat screen monitors.
type Timing struct {
	Elapsed time.Duration // Total loop time since last update.
	Update  time.Duration // Time used for previous state update.
	Renders int           // Render requests since last update.
}

// Zero all time and counter values.
func (t *Timing) Zero() {
	t.Update, t.Elapsed, t.Renders = 0, 0, 0
}

// Dump current amount of update loop time tracked in milliseconds.
// Times are expected to be reset each update.
func (t *Timing) Dump() {
	milliseconds := 1000.0
	e := t.Elapsed.Seconds() * milliseconds
	u := t.Update.Seconds() * milliseconds
	fmt.Printf("E:%2.4f U:%2.4f #:%d\n", e, u, t.Renders)
}

// Modelled returns the total number of models and the total
// number of verticies for all models.
func (t *Timing) Modelled(eng Eng) (models, verts int) {
	return eng.(*engine).models.counts()
}

// Rendered returns the number of models and the number
// of verticies rendered in the last rendering pass.
func (t *Timing) Rendered(eng Eng) (models, verts int) {
	return eng.(*engine).frames.drawCalls, eng.(*engine).frames.verticies
}
