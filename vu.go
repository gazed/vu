// Copyright © 2015-2024 Galvanized Logic Inc.

// Package vu - the virtual universe engine, provides 3D game support. Vu wraps
// subsystems like rendering, physics, asset loading, audio, etc. to provide
// higher level functionality that includes:
//   - Transform graphs and composite objects.
//   - Timestepped update/render loop.
//   - Access to user input events.
//   - Cameras and transform manipulation.
//   - Loading and controlling groups of graphics and audio assets.
//
// Refer to the vu/eg package for examples of engine functionality.
//
// Vu dependencies are:
//   - Vulkan for graphics card access.      See package vu/render.
//   - OpenAL for sound card access.         See package vu/audio.
//   - WinAPI for Windows display and input. See package vu/device.
package vu

// vu.go is the engine entry point for user apps. It defines how the
// user game communicates with the engine.

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gazed/vu/audio"
	"github.com/gazed/vu/device"
	"github.com/gazed/vu/load"
	"github.com/gazed/vu/render"
)

// NewEngine is called by the game to initialize the engine
// and its subsystems like rendering, physics, and audio. Eg:
//
//	eng, err := vu.NewEngine(vu.Windowed())
//
// The app uses eng to create the initial scenes prior to running
// the engine.
func NewEngine(config ...Attr) (eng *Engine, err error) {
	eng = &Engine{}
	eng.SetFrameLimit(60) // default FPS throttle

	// apply configuration overrides to the defaults.
	cfg := configDefaults
	for _, attr := range config {
		attr(&cfg)
	}

	// create engine systems to handle application data.
	eng.app = newApplication()

	// initialize the device layer needed by the renderer
	eng.dev = device.New(cfg.windowed, cfg.title, cfg.x, cfg.y, cfg.w, cfg.h)
	if err = eng.dev.CreateDisplay(); err != nil {
		eng.dispose() // can't continue without a display.
		return nil, fmt.Errorf("device.CreateDisplay failed %w", err)
	}
	eng.dev.SetResizeHandler(eng.handleResize)

	// initialize the graphic renderer and the display surface.
	eng.rc, err = render.New(render.VULKAN_RENDERER, eng.dev, cfg.title)
	if err != nil {
		eng.dispose() // can't continue without a renderer.
		return nil, fmt.Errorf("render.New failed %w", err)
	}
	eng.rc.SetClearColor(cfg.r, cfg.g, cfg.b, cfg.a)
	GPU = render.GPU // expose GPU type found by the render layer.

	// initialize audio.
	eng.ac = audio.New()
	if err := eng.ac.Init(); err != nil {
		slog.Error("no audio", "error", err)
		eng.ac.DisableAudio()
	}

	// default and fallback assets.
	if err := eng.app.ld.loadDefaultAssets(eng.rc); err != nil {
		eng.dispose() // can't continue without basic assets.
		return nil, fmt.Errorf("render.New failed %w", err)
	}
	return eng, nil
}

// AddScene creates a new application scene graph and camera.
// Scene graphs use zero to indicate that this is a root node.
func (eng *Engine) AddScene(st SceneType) *Entity {
	// expose public AddScene on the engine.
	return eng.app.addScene(st) // app does the real work.
}

// AddSound creates an entity from the named sound asset.
// The name is the sound asset filename without the .wav extension.
//
// Passing the sound identifier to an entity PlaySound() method will
// assigned using SetListener(). Sounds are louder the closer the
// played sound to the sound listener.
func (eng *Engine) AddSound(name string) (sound *Entity) {
	eid := eng.app.sounds.create(eng.app.eids, name) // new sound entity.

	// get the asset for this entity once it has been loaded.
	eng.app.ld.getAsset(assetID(aud, name), eid, eng.app.sounds.assetLoaded)
	return &Entity{app: eng.app, eid: eid} // return the sound entity.
}

// ImportAssets creates assets from the given asset files.
// Expected to be called at least once initialization to
// create the assets referenced by models in a scene.
func (eng *Engine) ImportAssets(assetFilenames ...string) {
	// public wrapper for the underlying loader file importer.
	eng.app.ld.importAssetData(assetFilenames...)
}

// SetFrameLimit throttles the engine to the given frames-per-second
// This reduces GPU usage when the actual FPS is higher than the given limit.
// It will not make the engine faster if the actual FPS is lower than
// the given limit. Throttle limits less than 30FPS and greater than 240FPS
// are ignored.
func (eng *Engine) SetFrameLimit(limit int) {
	if limit >= 30 && limit <= 240 {
		eng.throttle = time.Duration(float64(time.Second) / float64(limit))
	}
}

// GPU is either IntegratedGPU or DiscreteGPU based on what was
// found by the render layer. Exposed to allow the application
// to modify itself for smooth renders on less powerful GPUs.
var GPU render.GPUType // Set on startup after render layer is initialized.

// Expose some render layer values so the render package does
// not always need to be included.
const (
	IntegratedGPU = render.INTEGRATED_GPU
	DiscreteGPU   = render.DISCRETE_GPU
)

// =============================================================================

// Engine controls the engine subsystems and the run loop.
//
//	eng.Run(updater) // Run the engine.
type Engine struct {
	dev *device.Device  // OS specific platform for display and input.
	rc  *render.Context // Render interface.
	ac  *audio.Context  // Audio interface.
	app *application    // User application resources and state

	// Track time each refresh cycle to ensure fixed timestamp updates.
	suspended bool          // true if updating the game state is on hold.
	running   bool          // true if engine is alive.
	throttle  time.Duration // FPS throttle.
}

// Updator is responsible for updating application state each render frame.
// It is implemented by the user app and passed to eng.Run().
type Updator interface {
	// Update allows applications to change state prior to the next render.
	// Update is called each game loop update (many times a second) while the
	// game is running.
	//    eng : the game
	//    i   : user input refreshed prior to each call.
	Update(eng *Engine, i *Input, delta time.Duration)
}

// timestep is how often the state is updated. It is fixed at
// 60 times a second (1s/60 = 0.01666s) so that the game speed
// is constant (independent from computer speed and refresh rate).
// The timestep loop is implemented in a manner such that timesteps must
// be slower than the display refresh rate. See eng.update for details.
var (
	timestep     = time.Duration(16666667) // nanoseconds for 16.7ms
	timestepSecs = timestep.Seconds()
	startTime    = time.Now()
)

// Run the game engine. This method starts the game loop and does not
// return until the game shuts down. The game Update method is called
// each time the game loop updates.
func (eng *Engine) Run(updator Updator) {
	eng.app.updator = updator // application update callback

	// use a fixed timestep to run game updates 60 times a second
	var elapsedTime time.Duration    // accumulate time to trigger timesteps
	previousFrameStart := time.Now() // used to calculate delta time
	eng.running = true

	// loop forever process user input, updating game state, and rendering.
	for eng.running {

		// process user input.
		eng.app.input.Clone(eng.dev.GetInput())
		if !eng.dev.IsRunning() {
			slog.Info("engine shutdown!") // likely user closed window.
			eng.Shutdown()                //
			break                         // exit loop to eng.dispose()
		}

		// run updates while game is not suspended.
		// reset previousFrameStart when resuming (un-pause).
		if !eng.suspended {
			frameStart := time.Now()

			// delta measures the time it takes between frames.
			delta := frameStart.Sub(previousFrameStart)
			elapsedTime += delta

			// handle persistent slowness by dropping updates.
			// fix this by making the updates and render faster.
			if elapsedTime > 3*timestep {
				elapsedTime = timestep // run 1 update and drop the rest
			}

			// run updates at a fixed interval independent of frame rendering.
			// run multiple updates to catch up in cases of periodic slowness.
			for elapsedTime >= timestep {
				elapsedTime -= timestep

				// Simulate physics using a fixed timestep so that
				// each update advances by the same amount.
				eng.app.sim.simulate(eng.app.povs, timestepSecs)

				// FUTURE move particle effects using fixed timestep.
				// eng.app.models.moveParticles(timestepSecs)
			}

			// update the client app before each render frame
			eng.app.updator.Update(eng, eng.app.input, delta)
			if !eng.running {
				slog.Info("app shutdown!") // app called eng.Shutdown()
				break                      // exit loop to eng.dispose()
			}

			// check for any newly created assets.
			eng.app.ld.loadAssets(eng.rc, eng.ac)

			// FUTURE: advance model animations by elapsed time, not at fixed rate like physics.
			// Animation data expects to be played back at a particular frame rate.
			// eng.app.models.animate(delta)

			// render frames outside the fixed timestep.
			// FUTURE: interpolate the render as a fraction between this frame and last.
			eng.app.scenes.setViewMatrixes(eng.rc.Size())
			eng.app.povs.setWorldMatrix(delta)
			eng.app.frame = eng.app.scenes.getFrame(eng.app, eng.app.frame)
			eng.rc.Draw(eng.app.frame, delta)

			// frame complete, remember the start of this frame.
			previousFrameStart = frameStart

			// throttle to rest the CPU/GPU.
			// Requires go1.23+ to get 1ms pecision on windows. See go issue #44343.
			extra := eng.throttle - time.Since(frameStart) // FPS throttle
			extra = extra - extra%10_000                   // round down for wiggle room.
			if extra > 0 {
				time.Sleep(extra)
			}
		}
	}
	eng.dispose()
}

// handleResize processes user window changes.
func (eng *Engine) handleResize() {
	w, h := eng.dev.SurfaceSize()     // display window size - updated
	pw, ph := eng.rc.Size()           // renderer surface size - current
	x, y := eng.dev.SurfaceLocation() // display window upper left corner

	// window has been minimized
	if w == 0 || h == 0 {
		slog.Debug("app minimized: suspending")
		eng.suspended = true
		return
	}
	if eng.suspended {
		slog.Debug("app restored: resuming")
		eng.suspended = false
	}

	// update display surface if size has changed.
	if w != pw || h != ph {
		eng.rc.Resize(w, h)         // request render resize.
		eng.app.scenes.resize(w, h) // update scene cameras.
	}

	// update apps that have registered for resize callbacks.
	if eng.app.resizer != nil {
		eng.app.resizer.Resize(x, y, w, h)
	}
}

// Resizer is responsible for updating an application when the window
// is resized. It is implemented by the user app and set on startup.
type Resizer interface {
	// Resize is called when the window is resized.
	Resize(windowLeft, windowTop int32, windowWidth, windowHeight uint32)
}

// SetResizeListener sets the application callback
// for when the window is resized.
func (eng *Engine) SetResizeListener(resizer Resizer) {
	eng.app.resizer = resizer
}

// ToggleFullscreen switches between a borderless fullscreen window and
// a bordered window.
func (eng *Engine) ToggleFullscreen() {
	eng.dev.ToggleFullscreen()
}

// Mute toggles the sound volume.
// Engine attribute for use in Eng.Set().
func (eng *Engine) Mute(mute bool) {
	gain := 1.0
	if mute {
		gain = 0.0
	}
	eng.ac.SetGain(gain)
}

// MakeMeshes loads application generated mesh data.
func (eng *Engine) MakeMeshes(name string, meshes []load.MeshData) (err error) {
	mids, err := eng.rc.LoadMeshes(meshes) // upload all mesh data.
	if err != nil || len(mids) != len(meshes) {
		return fmt.Errorf("MakeMeshes %s: %w", name, err)
	}
	for i, mid := range mids {
		m := newMesh(fmt.Sprintf("%s%d", name, i))
		m.mid = mid
		eng.app.ld.assets[m.aid()] = m
	}
	labelRange := fmt.Sprintf("msh:%s%d:msh:%s%d", name, mids[0], name, mids[len(mids)-1])
	meshIDRange := fmt.Sprintf("%d:%d", mids[0], mids[len(mids)-1])
	slog.Debug("MakeMeshes", "asset", labelRange, "ids", meshIDRange)
	return nil
}

// Shutdown is an application request to close down the engine.
// Mark the engine as shutdown which will cause the game loop to exit.
func (eng *Engine) Shutdown() {
	eng.running = false
}
func (eng *Engine) dispose() {
	// cleanup up engine subsystem resources.
	if eng.ac != nil {
		eng.ac.Dispose()
		eng.ac = nil
	}
	if eng.rc != nil {
		eng.rc.Dispose()
		eng.rc = nil
	}
	if eng.dev != nil {
		eng.dev.Dispose()
		eng.dev = nil
	}
}
