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
	"runtime"
	"time"

	"github.com/gazed/vu/audio"
	"github.com/gazed/vu/device"
	"github.com/gazed/vu/load"
	"github.com/gazed/vu/render"
)

// init is called once on package load. Needed because
// underlying platforms insist that the windows are created
// on the main startup thread.
func init() { runtime.LockOSThread() }

// NewEngine is called by the game to initialize as much of
// the engine as the underlying platform allows before entering
// the main run loop. Eg: creating an ios display enters a run
// loop that does not return.
//
//	eng, err := vu.NewEngine(vu.Windowed())
func NewEngine(config ...Attr) (eng *Engine, err error) {
	eng = &Engine{}
	eng.SetFrameLimit(60) // default FPS throttle

	// apply configuration overrides to the defaults.
	eng.cfg = configDefaults
	for _, attr := range config {
		attr(&eng.cfg)
	}

	// create app to hold application created objects and resources.
	eng.app = newApplication()

	// initialize audio.
	eng.ac = audio.New()
	if err := eng.ac.Init(); err != nil {
		slog.Error("no audio", "error", err)
		eng.ac.DisableAudio()
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
//	eng.Run(loader, updater) // Run the engine.
type Engine struct {
	cfg Config          // engine configuration settings.
	dev *device.Device  // OS specific platform for display and input.
	rc  *render.Context // Render interface.
	ac  *audio.Context  // Audio interface.
	app *application    // User application resources and state

	// Track time each refresh cycle to ensure fixed timestamp updates.
	suspended      bool          // true if updating the game state is on hold.
	running        bool          // true if engine is alive.
	throttle       time.Duration // FPS throttle.
	prevFrameStart time.Time     // used to calculate delta time
	elapsedTime    time.Duration // accumulate time to trigger timesteps
}

// Loader is responsible for creating the initial application objects.
// The loader is implemented by the user app and passed to eng.Run().
// It is called once on startup.
type Loader interface {
	// Load allows applications to change state prior to the next render.
	// Update is called each game loop update (many times a second) while the
	// game is running.
	//    eng : the game engine.
	// Returning an error will stop the engine.
	Load(eng *Engine) error
}

// Updator is responsible for updating application state each render frame.
// The updated is implemented by the user app and passed to eng.Run().
// It is called once per engine update.
type Updator interface {
	// Update allows applications to change state prior to the next render.
	// Update is called each game loop update (many times a second) while the
	// game is running.
	//    eng  : the game engine.
	//    i    : user input refreshed prior to each call.
	//    delta: elapsed time since last call.
	Update(eng *Engine, i *Input, delta time.Duration)
}

// initializeDevice is called once on startup.
// The darwin systems can terminate the process on dispose.
func (eng *Engine) initializeDevice() (err error) {

	// initialize the device layer needed by the renderer
	if err := eng.dev.CreateDisplay(); err != nil {
		return fmt.Errorf("device.CreateDisplay failed %w", err)
	}
	eng.dev.SetResizeHandler(eng.handleResize)

	// initialize the graphic renderer and the display surface.
	cfg := eng.cfg
	eng.rc, err = render.New(render.VULKAN_RENDERER, eng.dev, cfg.title)
	if err != nil {
		return fmt.Errorf("render.New failed %w", err)
	}
	eng.rc.SetClearColor(cfg.r, cfg.g, cfg.b, cfg.a)
	GPU = render.GPU // expose GPU type found by the render layer.

	// default and fallback assets.
	if err := eng.app.ld.loadDefaultAssets(eng.rc); err != nil {
		return fmt.Errorf("render.New failed %w", err)
	}
	return nil
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

// runLoop is called two different ways.
// 1. directly  on windows
// 2. callbacks on apple devices
func (eng *Engine) runLoop() (running bool) {

	// process user input.
	eng.app.input.Clone(eng.dev.GetInput())
	if !eng.dev.IsRunning() {
		slog.Info("engine shutdown!") // likely user closed window.
		eng.Shutdown()                //
		return false                  // stop running
	}

	// suspend if focus is lost.
	eng.suspended = !eng.app.input.Focus

	// ignore updates and rendering while suspended. IOS in particular
	// causes Vulkan errors when rendering to an app without focus.
	if eng.suspended {
		return true // continue running to process input.
	}

	// render a frame.
	frameStart := time.Now()

	// delta measures the time it takes between frames.
	delta := frameStart.Sub(eng.prevFrameStart)
	eng.elapsedTime += delta
	eng.prevFrameStart = frameStart // remember for next frame.

	// handle persistent slowness by dropping updates.
	// fix this by making the updates and render faster.
	if eng.elapsedTime > 3*timestep {
		eng.elapsedTime = timestep // run 1 update and drop the rest
	}

	// run updates at a fixed interval independent of frame rendering.
	// run multiple updates to catch up in cases of periodic slowness.
	for eng.elapsedTime >= timestep {
		eng.elapsedTime -= timestep

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
		return false               //
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
	return true // continue running.
}

// WindowSize can be called once the display has initialized.
// This is any time after or during the Loader.Load() callback.
// Mainly needed for ios devices where the size is not known
// until the display has been created.
func (eng *Engine) WindowSize() (x, y, w, h uint32) {
	if eng.dev == nil {
		return 0, 0, 0, 0
	}
	xi, yi := eng.dev.SurfaceLocation()
	w, h = eng.dev.SurfaceSize()
	return uint32(xi), uint32(yi), w, h
}

// initialResize is called one time after the display surface has
// been initialized and after the app has created the initial scenes.
func (eng *Engine) initialResize() {
	w, h := eng.rc.Size()       // renderer surface size.
	eng.app.scenes.resize(w, h) // update scene cameras.
	eng.handleResize()          // call app resize.
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

// LoadErrors returns true if there were assets that failed to load.
// Ideally completely debug any loading problems prior to shipping.
func (eng *Engine) LoadErrors() bool { return eng.app.ld.failed > 0 }

// Mute toggles the sound volume.
// Engine attribute for use in Eng.Set().
func (eng *Engine) Mute(mute bool) {
	gain := 1.0
	if mute {
		gain = 0.0
	}
	eng.ac.SetGain(gain)
}

// MakeMeshes loads application generated mesh assets.
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

// MakeTextures loads application generated texture assets.
// FUTURE: have the render layer upload all the textures at once, similar to rc.LoadMeshes.
func (eng *Engine) MakeTextures(name string, textureData []*load.ImageData) (err error) {
	for i := range textureData {
		tid, err := eng.rc.LoadTexture(textureData[i])
		if err != nil {
			slog.Debug("MakeTextures error", "name", name, "index", i, "err", err)
			continue
		}
		t := newTexture(fmt.Sprintf("%s%d", name, i))
		t.tid = tid
		eng.app.ld.assets[t.aid()] = t
		slog.Debug("MakeTextures", "asset", "tex:"+t.label(), "tid", t.tid, "opaque", t.opaque)
	}
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
