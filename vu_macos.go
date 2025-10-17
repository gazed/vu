// Copyright © 2025 Galvanized Logic Inc.

//go:build darwin && !ios

package vu

// macos allows device initialization before taking over
// the run loop and calling rendering callbacks.

import (
	"io"
	"log/slog"
	"time"

	"github.com/gazed/vu/device"
	"github.com/gazed/vu/internal/device/macos"
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
func (eng *Engine) Run(loader Loader, updator Updator) {
	app_engine = eng          // singleton ref to the engine.
	eng.app.loader = loader   // application load callback
	eng.app.updator = updator // application update callback

	// one time device initialization
	cfg := eng.cfg
	eng.dev = device.New(cfg.windowed, cfg.title, cfg.x, cfg.y, cfg.w, cfg.h)
	if err := eng.initializeDevice(); err != nil {
		slog.Error("initializeDevice", "err", err)
		eng.dispose()
		return
	}

	// one time app load before update loop.
	if err := eng.app.loader.Load(eng); err != nil {
		slog.Error("loader.Load", "err", err)
		return
	}

	// update the apps scenes and trigger a resize.
	eng.initialResize()

	// start the run loop
	eng.prevFrameStart = time.Now()
	eng.running = true

	// loop forever process user input, updating game state, and rendering.
	eng.dev.Run(renderCallback, nil) // does not return while running

	// apple device's never return execution to this method.
	// The process is killed when the window is closed.
}

// renderCallback is triggered by the application loop at a frequency
// close to the display's refresh rate.
func renderCallback() {
	if app_engine == nil {
		slog.Error("renderCallback: app_engine not set")
		return
	}
	app_engine.runLoop()
	if !app_engine.running {
		app_engine.dispose()
	}
}

// Get a macos console writer that conforms to the io.Writer interface.
func ConsoleWriter() io.Writer {
	return &macos.MacOSConsoleWriter{}
}
