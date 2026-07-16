// SPDX-FileCopyrightText : © 2025 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

//go:build windows

package vu

// windows allows complete control of the run loop.

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/gazed/vu/device"
)

var accumulator time.Duration

// Run the game engine. This method starts the game loop and does not
// return until the game shuts down. The game Update method is called
// each time the game loop updates.
func (eng *Engine) Run(loader Loader, updater Updater) {
	eng.app.loader = loader   // application load callback
	eng.app.updater = updater // application update callback

	// one time device initialization
	cfg := eng.cfg
	eng.dev = device.New(cfg.windowed, cfg.title, cfg.x, cfg.y, cfg.w, cfg.h)
	if err := eng.initializeDevice(); err != nil {
		slog.Error("initializeDevice", "err", err)
		eng.dispose()
		return
	}

	// one time app load before update loop.
	if err := loader.Load(eng); err != nil {
		slog.Error("loader.Load", "err", err)
		return
	}

	// update the apps scenes and trigger a resize.
	eng.initialResize()

	// start the run loop
	eng.running = true
	eng.prevFrameStart = time.Now()
	var loopStart = eng.prevFrameStart

	// loop forever process user input, updating game state, and rendering.
	// Requires go1.23+ to get 1ms pecision on windows. See go issue #44343.
	for eng.running {
		delta := time.Since(loopStart) // how long since last loop started
		loopStart = time.Now()         // start the new loop.
		accumulator += delta

		// handle persistent slowness by dropping updates.
		// fix this by making the updates and render faster.
		if accumulator > 3*eng.throttle {
			accumulator = eng.throttle // run 1 update and drop the rest
		}

		// run the input poll, game update, and frame render at
		// the requested throttle rate.
		if accumulator >= eng.throttle {
			accumulator -= eng.throttle

			// run the main update/render loop and exit if the engine has quit.
			// if suspended this runs input poll, but ignores game update and frame render.
			if !eng.runLoop() {
				break
			}
		}

		// take breaks as needed.
		if eng.suspended {
			time.Sleep(eng.throttle / 2) // suspended: take a break for half a frame.
		} else {

			// time period was to short to measure. Timer resolution is around 1ms,
			// so take a short break rather than busy wait.
			if delta == 0 {
				time.Sleep(1 * time.Millisecond)
			}
		}
	}
	eng.dispose()
}

// The console just works on windows.
func ConsoleWriter() io.Writer {
	return os.Stderr
}
