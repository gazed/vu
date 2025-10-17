// Copyright © 2025 Galvanized Logic Inc.

//go:build windows

package vu

// windows allows complete control of the run loop.

import (
	"io"
	"os"
	"time"
)

// Run the game engine. This method starts the game loop and does not
// return until the game shuts down. The game Update method is called
// each time the game loop updates.
func (eng *Engine) Run(loader Loader, updator Updator) {
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
	if err := loader.Load(eng); err != nil {
		slog.Error("loader.Load", "err", err)
		return
	}

	// update the apps scenes and trigger a resize.
	eng.initialResize()

	// start the run loop
	eng.prevFrameStart = time.Now()
	eng.running = true
	var frameStart time.Time

	// loop forever process user input, updating game state, and rendering.
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

// The console just works on windows.
func ConsoleWriter() io.Writer {
	return os.Stderr
}
