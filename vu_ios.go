// Copyright © 2025 Galvanized Logic Inc.

//go:build ios

package vu

// ios takes over the run loop immediately on startup and uses
// callbacks to complete the one time engine initialization
// and the ongoing rendering.

import (
	"io"
	"log/slog"
	"time"

	"github.com/gazed/vu/device"
	"github.com/gazed/vu/internal/device/ios"
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

	// initialize the device data.
	cfg := eng.cfg
	eng.dev = device.New(cfg.windowed, cfg.title, cfg.x, cfg.y, cfg.w, cfg.h)

	// loop forever process user input, updating game state, and rendering.
	eng.dev.Run(renderCallback, loadCallback) // does not return while running

	// apple devices never return execution to this method.
	// The process is killed when the window is closed.
}

// loadCallback is called once on startup after the device
// display has been created.
func loadCallback() {
	if app_engine == nil {
		slog.Error("loadCallback: app_engine not set")
		return
	}
	eng := app_engine

	// one time device initialization
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

	// prep for the run loop renderCallbacks
	eng.prevFrameStart = time.Now()
	eng.running = true
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

// Get an ios console writer that conforms to the io.Writer interface.
func ConsoleWriter() io.Writer {
	return &ios.IOSConsoleWriter{}
}
