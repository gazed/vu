// Copyright © 2025 Galvanized Logic Inc.

//go:build windows

package vu

import (
	"time"
)

// Run the game engine. This method starts the game loop and does not
// return until the game shuts down. The game Update method is called
// each time the game loop updates.
func (eng *Engine) Run(updator Updator) {
	eng.app.updator = updator // application update callback
	eng.prevFrameStart = time.Now()
	eng.running = true

	// loop forever process user input, updating game state, and rendering.
	var frameStart time.Time
	for eng.running {
		if !eng.suspended {
			frameStart = time.Now()
		}

		// run the main update/render loop and exit if the engine has quit.
		if !eng.runLoop() {
			break
		}

		// throttle to rest the CPU/GPU.
		// Requires go1.23+ to get 1ms pecision on windows. See go issue #44343.
		if !eng.suspended {
			extra := eng.throttle - time.Since(frameStart) // FPS throttle
			extra = extra - extra%10_000                   // round down for wiggle room.
			if extra > 0 {
				time.Sleep(extra)
			}
		}
	}
	eng.dispose()
}
