// Copyright © 2015-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// profile.go - consolidate engine profiling data.
// FUTURE : lots can be done here, but it needs to be balanced against engine
//          code clutter. Variations in execution times are expected to be mostly
//          influenced by device capability, ie: mobile devices are less capable
//          and slower than consoles or desktops.

import (
	"fmt"
	"time"
)

// Profile is used to collect timing values while the application is
// running. The numbers are reset each update. Applications are expected
// to track and smooth these per-update values over a number of updates.
//
// In general expect things to go slower on slower machines and slower
// as the number of models, draw calls, and verticies increases.
//
// FPS = Renders/Elapsed. When all is well the  number of renders matches
// the monitor refresh rate, which is 60fps, for most flat screens.
// Skipped updates indicate when the program is overwhelming the platform
// it is running on and also means the game is slowing down.
type Profile struct {
	Elapsed time.Duration // Total time since last update.
	Update  time.Duration // Time used by last update.

	// Key indication that the application is asking to much from the
	// underlying platform. A non-zero value means the game is running
	// slower than desired.
	Skipped int // Updates skipped since last update.

	// The number of renders completed and the real time used.
	Renders int           // Renders completed since last update.
	Render  time.Duration // Render time used since last update.
}

// Zero all time and counter values. Called by the engine after each
// application update callback.
func (p *Profile) Zero() {
	p.Elapsed, p.Update, p.Skipped = 0, 0, 0
	p.Render, p.Renders = 0, 0
}

// Dump current amount of update loop time, tracked in milliseconds,
// to the console. Times are expected to be reset each update.
// Expected to be used for development debugging.
func (p *Profile) Dump() {
	milliseconds := 1000.0
	e := p.Elapsed.Seconds() * milliseconds
	u := p.Update.Seconds() * milliseconds
	fmt.Printf("E:%2.4f U:%2.4f #:%d\n", e, u, p.Renders)
}

// Rendered returns the total number of models and the number
// of models rendered in the last rendering pass.
func (p *Profile) Rendered(eng Eng) (models, rendered int) {
	models = eng.(*application).models.stats()
	rendered = len(eng.(*application).frame)
	return models, rendered
}
