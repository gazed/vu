// Copyright © 2015-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// Package vu - virtual universe, provides 3D application support. Vu wraps
// subsystems like rendering, physics, asset loading, audio, etc. to provide
// higher level functionality that includes:
//    • Transform graphs and composite objects.
//    • Timestepped update/render loop.
//    • Access to user input events.
//    • Cameras and transform manipulation.
//    • Loading and controlling groups of graphics and audio assets.
// Refer to the vu/eg package for examples of engine functionality.
//
// Vu dependencies are:
//    • OpenGL (ES) for graphics card access.   See package vu/render.
//    • OpenAL for sound card access.           See package vu/audio.
//    • Cocoa  for OSX windowing and input.     See package vu/device.
//    • Xcode  for iOS cross-compiling & build. See package vu/device.
//    • WinAPI for Windows windowing and input. See package vu/device.
package vu

// vu.go holds package constants and functions.

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/gazed/vu/device"
	"github.com/gazed/vu/physics"
	"github.com/gazed/vu/render"
)

// Run creates the engine and initializes the underlying device layers.
// It does not return, transferring control to the main render loop.
// The render loop calls the application through the App interface.
// Run is expected to be called once on application startup.
//    app  : used by engine to communicate with App.
func Run(app App) (err error) {
	if app == nil {
		return fmt.Errorf("No application. Shutting down.")
	}
	eng, err := newEngine(app) // engine operator starts engine.
	if err != nil {
		return err
	}
	defer eng.shutdown() // ensure shutdown happens no matter what.
	defer catchErrors()  // dump any panics.
	device.Run(eng)      // This method does not return!!!

	// if the device is killed, for example the user closes the window,
	// then this next line won't even be reached and the defers won't
	// be run since as the OS kills the process.
	return nil // Should not even reach here.
}

// Vu engine startup.
// ===========================================================================
// Wrap physics for cleaner application API's.

// Body wraps physics.Body and groups the physics wrappers together
// in the documentation. Allows engine users to access common physics
// methods without including physics package.
type Body physics.Body

// Box creates a box shaped physics body located at the origin.
// The box size is given by the half-extents so that actual size
// is w=2*hx, h=2*hy, d=2*hz.
func Box(hx, hy, hz float64) Body {
	return physics.NewBody(physics.NewBox(hx, hy, hz))
}

// Sphere creates a ball shaped physics body located at the origin.
// The sphere size is defined by the radius.
func Sphere(radius float64) Body {
	return physics.NewBody(physics.NewSphere(radius))
}

// Ray creates a ray located at the origin and pointing in the
// direction dx, dy, dz.
func Ray(dx, dy, dz float64) Body {
	return physics.NewBody(physics.NewRay(dx, dy, dz))
}

// Plane creates a plane located on the origin and oriented by the
// plane normal nx, ny, nz.
func Plane(nx, ny, nz float64) Body {
	return physics.NewBody(physics.NewPlane(nx, ny, nz))
}

// physics wrappers.
// =============================================================================
// public constants

// Engine constants needed as input to methods as noted.
const (

	// Global graphic state constants. See Eng.State
	Blend       = render.Blend       // Alpha blending. Enabled by default.
	CullFace    = render.CullFace    // Backface culling. Enabled by default.
	DepthTest   = render.DepthTest   // Z-buffer awareness. Enabled by default.
	StaticDraw  = render.StaticDraw  // Created once, rendered many times.
	DynamicDraw = render.DynamicDraw // Data continually being updated.

	// Per-model rendering constants for Model DrawMode option.
	Triangles = render.Triangles // Triangles are the norm.
	Points    = render.Points    // Used for particle effects.
	Lines     = render.Lines     // Used for drawing lines and boxes.

	// KeyReleased indicator. Total time down, in update ticks,
	// is key down ticks minus KeyReleased. See App.Update.
	KeyReleased = device.KeyReleased
)

// public constants
// =============================================================================
// internal constants and utility methods.

// constants to ensure reasonable input values.
const (
	maxWindowTitle    = 40  // Max number of characters for a window title.
	minWindowSize     = 100 // Miniumum pixels for a window width or height.
	minWindowPosition = 0   // Bottom left corner of the screen.

	// timeStep is how often the state is updated. It is fixed at
	// 50 times a second (1s/50 = 0.02s) so that the game speed is
	// constant (independent from computer speed and refresh rate).
	// The timestep loop is implemented in a manner such that timeSteps must
	// be slower than the display refresh rate. See eng.update for details.
	timeStep = time.Duration(20 * time.Millisecond) // 0.02s, 50fps int64
)

// timeStepSecs is the float representation of the fixed timeStep.
var timeStepSecs = timeStep.Seconds()

// catchErrors should be defered at the top of each goroutine so that
// errors can be logged in production loads as required by the application.
func catchErrors() {
	if r := recover(); r != nil {
		log.Printf("Panic %s: %s Shutting down.", r, debug.Stack())
		os.Exit(-1)
	}
}
