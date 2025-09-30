// Copyright © 2025 Galvanized Logic Inc.

//go:build darwin || ios

package vu

import (
	"time"
)

// app_engine is set in Run and used in the rendercallback to
// drive the main application loop on Apple devices.
// This is because apple devices control the main loop
// and call back to the engine for rendering.
// (in windows the engine controls the main loop).
var app_engine *Engine

// Run the game engine. This method starts the game loop and does not
// return until the game shuts down. The game Update method is called
// each time the game loop updates.
func (eng *Engine) Run(updator Updator) {
	eng.app.updator = updator // application update callback
	eng.prevFrameStart = time.Now()
	eng.running = true
	app_engine = eng

	// loop forever process user input, updating game state, and rendering.
	eng.dev.Run(renderCallback) // does not return while running

	// apple device's never return execution to this method.
	// The process is killed when the window is closed.
}

// renderCallback is triggered by the application loop at a frequency
// close to the display's refresh rate.
func renderCallback() {
	app_engine.runLoop()
	if !app_engine.running {
		app_engine.dispose()
	}
}
