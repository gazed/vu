// Copyright © 2013-2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// eng.go holds the Eng interface and engine class.

import (
	"fmt"
	"log"
	"time"

	"github.com/gazed/vu/audio"
	"github.com/gazed/vu/device"
	"github.com/gazed/vu/physics"
	"github.com/gazed/vu/render"
)

// Eng provides support for a 3D application conforming to the App interface.
// Eng uses scenes to display application created entities using a camera.
// The application creates one or more scenes, configures the scene cameras,
// and adds render components to each scene.
//
// Eng also controls the overall engine lifetime, sets global engine state,
// and provides application profiling timing.
type Eng interface {
	Shutdown()     // Stop the engine and free allocated resources.
	State() *State // Query engine state. State updated per tick.

	// AddScene creates a new application scene graph and camera.
	AddScene() *Ent // Group rendered objects with a Camera.

	// AddSound loads audio data and returns a unique sound identifier.
	// Passing the sound identifier to an entity PlaySound() method will
	// play the sound at the entities location. Set the single listener
	// location using SetListener(). Sounds are louder the closer the
	// played sound to the sound listener.
	AddSound(name string) uint32

	// Set changes engine wide attributes. It accepts one or more
	// functions that take an EngAttr parameter, ie: vu.Color(0,0,0).
	Set(...EngAttr) // Change one or more engine attributes.

	// Times for the previous update loop. The application
	// can average times over multiple updates.
	Times() *Profile // Per update loop performance metrics.
}

// engine controls the run loop and access to the device layer
// One instance of engine is created by Run.
type engine struct {
	dev device.Device  // OS specific window and rendering context.
	gc  render.Context // Graphics card interface.
	ac  audio.Audio    // Audio card interface.
	app *application   // Application controller implementing Eng.

	// Track time each refresh cycle to ensure fixed timestamp updates.
	startTime  time.Time     // Track start of a game loop display refresh.
	updateTime time.Duration // Grows until big enough to trigger update.
	elapsed    time.Duration // Time since last update. Reset on update.

	// communication with the update goroutine.
	doneUpdate chan *application // update finished.
}

// newEngine initializes the device layers, returning an error
// if any problems are encountered. Called once by Run on startup.
func newEngine(app App) (eng *engine, err error) {
	eng = &engine{}

	// initialize audio.
	eng.ac = audio.New()
	if err = eng.ac.Init(); err != nil {
		log.Printf("No audio. %s.", err)
		eng.ac = &audio.NoAudio{} // Disable audio.
	}
	eng.app = newApplication(app)
	return eng, nil
}

// Init is a one-time callback from the device. It is the
// callback used to start model creation and asset loading.
// Init implements device.App and is called just the once before
// starting the regular callbacks to the Update method.
func (eng *engine) Init(d device.Device) {
	eng.dev = d

	// initialize graphics now that context is available.
	// Graphics can be initialized before device on OSX, but Windows
	// needs a proper context to find the OpenGL functions
	eng.gc = render.New()
	if err := eng.gc.Init(); err != nil {
		log.Printf("Failed starting graphics %s.", err)
		eng.shutdown() // Can't continue without graphics.
	}
	eng.gc.Enable(Blend, true)    // expected application startup state.
	eng.gc.Enable(CullFace, true) // expected application startup state.
	eng.doneUpdate = make(chan *application)

	// Create engine and run initial App.Create on main thread.
	eng.app.state.setScreen(eng.dev.Size())
	eng.app.state.Full = eng.dev.IsFullScreen()
	eng.app.app.Create(eng.app, eng.app.state)
	eng.app.setAttributes(eng)
	eng.startTime = time.Now() // begin tracking elapsed time after init.

	// if App shuts down engine during Create, possibly due to errors,
	// eng.dev will be shutdown as well, stopping any Refresh callback.
}

// Refresh controls the main render loop. It is called by the device layer
// whenever a new frame is needed. It runs fixed timestep updates that are
// separate from the frame refresh rate:.
//    Modern frame display refresh rate is 60fps (0.0167sec) or higher.
//    Application and physics update is fixed at 50fps (0.02sec).
//
// While the game loop is controlled by display refresh callbacks updates
// are still done at fixed timesteps. See:
//    http://gameprogrammingpatterns.com/game-loop.html
//    https://gafferongames.com/post/fix_your_timestep
//    http://lspiroengine.com/?p=378
func (eng *engine) Refresh(dev device.Device) {
	eng.elapsed += time.Since(eng.startTime) // Add time since last Refresh
	eng.startTime = time.Now()               // Reset loop start time tracker.

	// Need the *application instance back from update,
	// blocking if update is not done.
	if eng.app == nil {
		eng.app = <-eng.doneUpdate
	}
	if eng.app.stop {
		eng.shutdown()
		return
	}

	// handle queued up device requests from previous updates since the
	// application goroutine can't access device layer resources.
	eng.app.setAttributes(eng)  // set device attributes from update.
	eng.app.scenes.release(eng) // remove disposed device data.
	eng.app.scenes.rebind(eng)  // refresh changes to GPU assets.
	eng.app.models.rebind(eng)  // refresh changes to GPU assets.
	eng.app.sounds.rebind(eng)  // refresh changes to Audio assets.
	eng.app.sounds.play(eng)    // play sounds requested by application.

	// render a display frame passing in the interpolated time between
	// the most recent update and the previous update.
	eng.render(eng.updateTime.Seconds() / timeStepSecs)

	// run update, if necessary, as a goroutine in the dead time after
	// rendering a frame and before the next frame request.
	eng.update()          // run partly as a goroutine
	eng.dev.SwapBuffers() // can sleep the main thread on MacOS.
}

// render a frame where lerp is the linear interpolation ratio between
// the last two updates.
// fmt.Printf("time since last render %2.4fms\n", eng.elapsed.Seconds()*1000)
func (eng *engine) render(lerp float64) {
	startRender := time.Now() // track this render.
	app := eng.app

	// Use lerp to calculate positions between the two most recent updates.
	// This smooths rendering between the different update and display rates.
	app.povs.renderAt(lerp)
	app.scenes.renderAt(lerp, app.state.W, app.state.H)
	app.frame = app.scenes.draw(app, app.frame)

	// render by sending the frame draw calls to the render context.
	eng.gc.Clear()
	for _, dc := range app.frame {
		eng.gc.Render(dc)
	}
	app.prof.Renders++
	app.prof.Render += time.Since(startRender)
}

// update does some work on the main thread and the reset on a goroutine.
// Access to the *application instance is relinquished until the update
// goroutine completes.
//
// Normal fixed-timestep-code loops while there is more update time available in
// order to handle the case where the update rate is faster than the display.
// However in this case only one update is run and any catchup will occur on
// future display callbacks given that the monitors refresh rate is expected
// to be faster than the game update rate. The slowest monitor refresh rate
// is expected to be 60Hz and the engines fixed timestep update is 50hz.
func (eng *engine) update() {
	app := eng.app
	app.prof.Elapsed += eng.elapsed // Track total time...
	eng.updateTime += eng.elapsed   // Add time for updates.

	// Run one update. A normal fixed timestep runs a loop here to handle
	// update rates run faster than display refreshes. In this engine the
	// update rate is slower than display refreshes, assuming 60Hz refresh
	// minimum, so updates can eventually catch up.
	if eng.updateTime >= timeStep {
		start := time.Now()        // track this update.
		eng.updateTime -= timeStep // Remove time consumed by this update.
		app.ut++                   // Track the total update ticks.

		// Input polling clears the accumulated user input.
		app.input.poll(eng.dev.Down(), app.ut)
		if app.input.Resized {
			app.state.Full = eng.dev.IsFullScreen()
			app.state.setScreen(eng.dev.Size())
			eng.gc.Viewport(app.state.W, app.state.H)
		}

		// The application must finish an update within a reasonable time,
		// Otherwise updates are dropped to avoid the spiral of death where
		// slow updates start delaying future updates. Note that dropping
		// updates effectively slows the game down and the application is
		// informed using the profiling data. To check for slow behaviour
		// try "go build -race" which slows an app quite a bit.
		//
		// This algorithm aggressively drops updates. However its either ok
		// because the application is doing one-off processing like switching
		// levels, or its not ok because the application is doing more than
		// the platform can handle. The later suggests tuning the application.
		//
		// FUTURE: leave some time so that updates can catch up. Would need
		//         to ensure the lerp calculation is still correct.
		for eng.updateTime >= timeStep {
			eng.updateTime -= timeStep // ensure updateTime < timeStep for lerp.
			app.prof.Skipped++         // report the game slowdown.
		}

		// Perform the physics and application update for
		// object locations, particles, animations, entity creation, etc.
		go app.update(app.ut, start, app.prof.Elapsed, eng.doneUpdate)
		eng.app = nil // no access to app while running update.
	}
	eng.elapsed = 0 // reset total elapsed time tracker.
}

// shutdown releases the engine resources allocated on startup.
func (eng *engine) shutdown() {
	if eng.ac != nil {
		eng.ac.Dispose()
		eng.ac = nil
	}
	if eng.dev != nil {
		eng.dev.Dispose()
		eng.dev = nil
	}
	if eng.app != nil {
		eng.app.shutdown()
	}
}

// bind sends data to the graphics or audio card. Data needs to be bound
// before it can be used for rendering or audio. Data needs rebinding
// if it is changed.
//
// Note that rebinding mesh, texture, and sound does not change the device
// reference, it just updates the data on the graphics or audio device.
// The bind is completed before the data is used again.
//
// FUTURE: move all knowledge of asset binding into the render package.
func (eng *engine) bind(a asset) error {
	switch d := a.(type) {
	case *Mesh:
		d.rebind = false // for bind() from loader instead of Mesh.bind().
		return eng.gc.BindMesh(&d.vao, d.vdata, d.faces)
	case *shader:
		var err error
		d.program, err = eng.gc.BindShader(d.vsh, d.fsh, d.uniforms, d.layouts)
		return err
	case *Texture:
		d.rebind = false // for bind() from loader instead of Texture.bind().
		return eng.gc.BindTexture(&d.tid, d.img)
	case *sound:
		return eng.ac.BindSound(&d.sid, &d.did, d.data)
	case *shadows:
		return eng.gc.BindMap(&d.bid, &d.tex.tid)
	case *target:
		err := eng.gc.BindTarget(&d.bid, &d.tex.tid, &d.db)
		return err
	}
	return fmt.Errorf("eng:bind. Unhandled bind request")
}

// release figures out what data to release based on the releaseData type.
//
// FUTURE: move all knowledge of asset binding into the render package.
func (eng *engine) release(asset interface{}) {
	switch a := asset.(type) {
	case *Mesh:
		eng.gc.ReleaseMesh(a.vao)
	case *shader:
		eng.gc.ReleaseShader(a.program)
	case *Texture:
		eng.gc.ReleaseTexture(a.tid)
	case *sound:
		eng.ac.ReleaseSound(a.sid)
	case *shadows:
		eng.gc.ReleaseMap(a.bid, a.tex.tid)
		a.bid, a.tex.tid = 0, 0
	case *target:
		eng.gc.ReleaseTarget(a.bid, a.tex.tid, a.db)
		a.bid, a.tex.tid, a.db = 0, 0, 0
	default:
		log.Printf("machine.release: No bindings for %T", a)
	}
}

// clampTex requests a texture to be non-repeating.
// Expected to be called once when setting up a texture.
func (eng *engine) clampTex(tid uint32) {
	eng.gc.SetTextureMode(tid, true)
}

// engine
// =============================================================================
// engine attributes reeduces the eng API footprint using functional options:
//    http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
//    https://commandcenter.blogspot.ca/2014/01/self-referential-functions-and-design.html
//
// The public functions are called on the update goroutine. Because the
// methods affect the device layer, they are stored and then called when
// back on the main engine thread.

// EngAttr defines an engine attribute that can be used in Eng.Set().
// For example.
//    eng.Set(Color(1,1,1), Mute(true))
type EngAttr func(*engine)

// Color sets the background window clear color.
// Engine attribute for use in Eng.Set().
func Color(r, g, b, a float32) EngAttr {
	return func(eng *engine) { eng.gc.Color(r, g, b, a) }
}

// Title sets the window title. For windowed mode.
// Engine attribute for use in Eng.Set().
func Title(t string) EngAttr {
	return func(eng *engine) { eng.dev.SetTitle(t) }
}

// Size sets the window starting location and size in pixels.
// For windowed mode.
// Engine attribute for use in Eng.Set().
func Size(x, y, w, h int) EngAttr {
	return func(eng *engine) {
		eng.dev.SetSize(x, y, w, h)
		eng.gc.Viewport(w, h)
		eng.app.state.setScreen(eng.dev.Size())
	}
}

// CursorOn hides or shows the cursor.
// Engine attribute for use in Eng.Set().
func CursorOn(show bool) EngAttr {
	return func(eng *engine) { eng.dev.ShowCursor(show) }
}

// CursorAt places the cursor at the window pixel x,y.
// Engine attribute for use in Eng.Set().
func CursorAt(x, y int) EngAttr {
	return func(eng *engine) { eng.dev.SetCursorAt(x, y) }
}

// On enables/disables render attributes like Blend, CullFace, etc...
// Engine attribute for use in Eng.Set().
func On(attr uint32, enabled bool) EngAttr {
	return func(eng *engine) { eng.gc.Enable(attr, enabled) }
}

// ToggleFullScreen flips full screen and windowed mode.
// Engine attribute for use in Eng.Set().
func ToggleFullScreen() EngAttr {
	return func(eng *engine) { eng.dev.ToggleFullScreen() }
}

// Mute toggles the sound volume.
// Engine attribute for use in Eng.Set().
func Mute(mute bool) EngAttr {
	gain := 1.0
	if mute {
		gain = 0.0
	}
	return func(eng *engine) { eng.ac.SetGain(gain) }
}

// Volume sets the sound volume.
// Engine attribute for use in Eng.Set().
func Volume(zeroToOne float64) EngAttr {
	return func(eng *engine) { eng.ac.SetGain(zeroToOne) }
}

// Gravity changes the physics gravity constant.
// Engine attribute for use in Eng.Set().
func Gravity(g float64) EngAttr {
	return func(eng *engine) {
		eng.app.bodies.physics.Set(physics.Gravity(g))
	}
}
