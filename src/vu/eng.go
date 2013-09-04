// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

// Package vu (virtual universe) provides 3D application support. Vu wraps
// the individual subsystems like rendering, physics, data loading,
// audio, etc. to provide higher level functionality that includes:
//     scene graph,
//     timestepped update/render loop,
//     user input handling and notification,
//     composite objects like Banners that combine font and mesh data,
//     cameras and transform manipulation,
//     shaders.
//
// The vu/eg (examples) package provides relatively small working code samples
// of engine functionality.  Some of the examples are based solely on the engine
// subsystem support packages, and others are based on the engine package.
//
// Effort has been made to avoid third party packages other than where there was
// no choice. Dependencies include:
//      OpenGL for access to the graphics card (vu/render/gl).
//      OpenAL for access to the sound card (vu/audio).
//      Cocoa  for access to OSX windowing and input (vu/device).
//      WinAPI for access to Win windowing and input (vu/device).
package vu

import (
	"time"
	"vu/audio"
	"vu/data"
	"vu/device"
	"vu/physics"
	"vu/render"
)

// Eng, the engine, is where everything starts. Eng is top of the hierarchy
// of the application window and provides access to the capabilities of the
// sub-components.
//
// Eng initializes the underlying subsystems and, for the most part, wraps
// access to subsystem functionality.
type Eng struct {
	device  device.Device       // Os specific window and rendering context.
	gc      render.Renderer     // Graphics card interface layer.
	ac      audio.Audio         // Audio card interface layer.
	auditor audio.SoundListener // Audio listener.
	phy     physics.Physics     // Physics controls/applies forces and handles collisions.
	gid     int                 // Globally incremented unique id provider.
	res     *roadie             // Data resource organizer.
	scenes  []*scene            // Scene graph.
	overlay int                 // Overlay is the last scene drawn.
	app     Director            // Application callbacks.
	Xm, Ym  int                 // Mouse position updated each loop.
	Xp, Yp  int                 // Previous mouse position updated each loop.
	Gt, Dt  float32             // Game time, delta time, updated each loop.
}

// New creates a 3D engine. The expected usage is:
//      if eng, err = vu.New("Title", 100, 100, 800, 600); err != nil {
//          log.Printf("Failed to initialize engine %s", err)
//          return
//      }
//      defer eng.Shutdown() // close down nicely.
//      // application 3D setup and initialization.
//      eng.Action()         // runs application loop (does not return).
func New(name string, x, y, width, height int) (eng *Eng, err error) {
	if name == "" {
		name = "Title"
	}
	if width < 100 {
		width = 100
	}
	if height < 100 {
		height = 100
	}
	eng = &Eng{}
	eng.phy = physics.New()

	// initialize the os specific shell, graphics context, and
	// user input monitor.
	eng.device = device.New(name, x, y, width, height)

	// initialize the audio layer.
	eng.auditor = audio.NewSoundListener()
	eng.ac = audio.New()
	if err = eng.ac.Init(); err != nil {
		eng.Shutdown()
		return // failed to initialize audio layer
	}

	// initialize the graphics layer.
	eng.gc = render.New()
	if err = eng.gc.Init(); err != nil {
		eng.Shutdown()
		return // failed to initialize graphics layer.
	}
	eng.res = newRoadie(eng.ac, eng.gc)
	eng.gc.Viewport(width, height)
	eng.device.Open()
	return
}

// Shutdown stops the engine and frees up any allocated resources.
func (eng *Eng) Shutdown() {
	if eng.ac != nil {
		eng.ac.Shutdown()
	}
	if eng.res != nil {
		eng.res.dispose()
	}
	if eng.device != nil {
		eng.device.Dispose()
	}
}

// Action is the main update/render loop. This regulates game update/render frequency
// and is based on:
//     http://gafferongames.com/game-physics/fix-your-timestep
//     http://www.koonsolo.com/news/dewitters-gameloop
//     http://sacredsoftware.net/tutorials/Animation/TimeBasedAnimation.xhtml
// The loop runs until the application closes.
//
// The application state is updated a variable number of times each loop in order
// that each state update is the same fixed timestep interval.
func (eng *Eng) Action() {
	// game time is the total elapsed time.  This is the total time spent
	// updating state.  Game-time/delta-time gives the number of game ticks.
	gt := float64(0)

	// delta time is how often the state is updated.  It is fixed at
	// 50 times a second (50/1000ms = 0.02) so that the game speed is constant
	// (independent from computer speed and refresh rate).
	dt := float64(0.02)

	// update time tracks the time available for updating state.  It carries
	// any unused update time into the next loop.  At the start of each loop
	// available time (based on rendering) is added.  Slow rendering causes
	// more time added on for updates and fast rendering results less time
	// for updates per loop, causing potentially no updates in a given loop.
	updateTime := float64(0)

	// elapsedTime tracks how long one frame/loop took.  This will be
	// capped if updating and rendering took a very long time in order to
	// avoid a spiral of death where even more updating is attempted when
	// things are running slow.
	elapsedTime := float64(0)

	// capTime guards against unreasonably slow updates and the spiral of death.
	// Essentially ignore any updating and rendering time that was more than 200ms.
	const capTime = float64(0.2)
	lastTime := time.Now() // the computer time updated every frame/game-loop

	// 3D loops are forever (but really only last until the user wimps out)
	for eng.device.IsAlive() {

		// how long since the last time through the loop.  The more time the loop
		// took, the more updates will need to be performed.
		elapsedTime = time.Since(lastTime).Seconds()
		lastTime = time.Now()
		if elapsedTime > capTime {
			elapsedTime = capTime
		}

		// ease up on the CPU if the render speed is over 100fps.
		if elapsedTime < 0.01 {
			time.Sleep(time.Duration((0.01-elapsedTime)*1000) * time.Millisecond)
		}

		// run updates based on how long the previous loop took.  This advances
		// state at a constant rate (dt).
		updateTime += elapsedTime
		for updateTime >= dt {
			eng.update(gt, dt) // update state, physics and animations.
			updateTime -= dt   // track the used delta time.
			gt += dt           // track the elapsed total time
		}

		// TODO interpolate the state based on the remaining delta time.  Right now
		//      the rendering is done on un-interpolated state which may be slightly
		//      behind where it should be.
		// interpolatedTime := updateTime / dt;  // fraction of unused delta time between 0 and 1.
		// State state = currentState*interpolatedTime + previousState * ( 1.0 - interpolatedTime );

		// redraw everything based on the current state.
		eng.render()
	}
}

// update is expected to be called from the engine Action loop. It will
// process any user feedback and then move game state forward.
func (eng *Eng) update(gameTime, deltaTime float64) {
	eng.Gt = float32(gameTime)
	eng.Dt = float32(deltaTime)

	// Advance the physics simulation for moving bodies.
	eng.phy.Step(eng.Dt)

	// Get user button/key presses and turn them into action requests.  The action
	// requests are processed by the application to change the game state.
	pressedSequences, mx, my := eng.device.ReadAndDispatch()
	eng.Xp, eng.Yp = eng.Xm, eng.Ym
	eng.Xm, eng.Ym = mx, my
	if eng.app != nil {
		eng.app.Update(pressedSequences, eng.Gt, eng.Dt)
	}

	// keep the scene parts sorted by distance inorder for transparency to work. The
	// further away objects have to be drawn first.
	for _, sc := range eng.scenes {
		if sc.radius > 0 {
			sc.cullParts(sc.parts)
		}
		if sc.sorted {
			sc.sortParts(sc.parts)
		}
	}
}

// render is called from the engine Action loop.  Its job is to render the active
// scenes in the order given (if they are visible), ensuring that an special overlay
// scene is drawn last.
func (eng *Eng) render() {
	eng.gc.Clear()
	var overlay *scene
	for _, sc := range eng.scenes {
		if sc.Visible() {
			if sc.uid == eng.overlay {
				overlay = sc
				continue
			}
			sc.render(eng.gc)
		}
	}
	if overlay != nil {
		overlay.render(eng.gc)
	}
	eng.device.SwapBuffers()
}

// SetDirector establishes application callbacks.
func (eng *Eng) SetDirector(director Director) {
	eng.app = director
	eng.device.SetFocuser(eng.app)
	eng.device.SetResizer(eng.app)
}

// guid generates unique id's that are valid for a currently running application.
// These are used to identify/tag resources. They are not meant to be persisted
// across application restarts.
func (eng *Eng) guid() int {
	eng.gid += 1
	return eng.gid
}

// ===========================================================================
// Expose/wrap device level information.

// GetSize returns the application viewport area in pixels.  This excludes any
// OS specific window trim.  The window x, y coordinates are the bottom left of
// the window.
func (eng *Eng) Size() (x, y, width, height int) { return eng.device.Size() }

// ResizeViewport needs to be called on window resize to adjust the graphics
// viewport.  This is expected to be called by the Application from the
// Director.Resize callback.
func (eng *Eng) ResizeViewport(x, y, width, height int) { eng.gc.Viewport(width, height) }

// ShowCursor hides and locks the cursor for the current window.
func (eng *Eng) ShowCursor(show bool) { eng.device.ShowCursor(show) }

// SetCursorAt puts the cursor at the given window location.  This is used
// when the cursor is hidden and the mouse movements are being tracked.
// Setting the cursor to the middle of the screen ensures movement doesn't
// get stuck.
func (eng *Eng) SetCursorAt(x, y int) {
	diffx, diffy := eng.Xm-x, eng.Ym-y
	eng.Xp, eng.Yp = eng.Xp+diffx, eng.Yp+diffy
	eng.Xm, eng.Ym = x, y
	eng.device.SetCursorAt(x, y)
}

// ===========================================================================
// Expose/wrap graphic and audio controls.

// Color sets the default background clear color. This color will appear if nothing
// else is drawn over it.
func (eng *Eng) Color(r, g, b, a float32) { eng.gc.Color(r, g, b, a) }

// Enable or disable global graphics attributes.
// Current valid values are: CULL, BLEND, DEPTH
func (eng *Eng) Enable(attribute int, enabled bool) { eng.gc.Enable(attribute, enabled) }

// Global graphic state constants. These are the attributes used in the
// eng.Enable method.
const (
	BLEND = render.BLEND // alpha blending.
	CULL  = render.CULL  // backface culling.
	DEPTH = render.DEPTH // z-buffer awareness.
)

// BindModel allows externally created meshes to be bound to the underlying
// graphics card.
func (eng *Eng) BindModel(mesh *data.Mesh) error { return eng.gc.BindModel(mesh) }

// AuditorLocation sets the 3D location of the entity that can hear sounds.  Sounds that are
// played at other locations will be heard more faintly as the distance between the played
// sound and listener increases.
func (eng *Eng) AuditorLocation(x, y, z float32) { eng.auditor.SetLocation(x, y, z) }

// UseSound creates a SoundMaker that is linked to the given sound resource.
func (eng *Eng) UseSound(sound string) audio.SoundMaker {
	s := eng.res.useSound(sound)
	return audio.NewSoundMaker(s)
}

// Mute turns the game sound on (mute == false) or off (mute == true).
func (eng *Eng) Mute(mute bool) { eng.ac.Mute(mute) }

// ===========================================================================
// Expose/wrap data loading capabilities.

// Loaders imports the named resource and makes it available in the resource cache.
// These are only expected to be used for externally created resources since most of
// the available data resources will be lazy loaded as Parts are created.
//
// These methods can also be used to pre-load resources, which will avoid lazy loading.
func (eng *Eng) LoadTexture(name string)  { eng.res.loadTexture(name) }
func (eng *Eng) LoadMesh(name string)     { eng.res.loadMesh(name) }
func (eng *Eng) LoadMaterial(name string) { eng.res.loadMaterial(name) }
func (eng *Eng) LoadGlyphs(name string)   { eng.res.loadGlyphs(name) }
func (eng *Eng) LoadSound(name string)    { eng.res.loadSound(name) }
func (eng *Eng) LoadShader(name string)   { eng.res.loadShader(name) }
func (eng *Eng) Load(data interface{})    { eng.res.load(data) }

// Loaded checks if the named data type is already initialized and in the cache.
func (eng *Eng) Loaded(name string, data interface{}) bool { return eng.res.loaded(name, data) }

// ===========================================================================
// Expose scene graph capabilities.

// AddScene creates a new scene with its own camera and lighting.
func (eng *Eng) AddScene(transform int) Scene {
	if eng.scenes == nil {
		eng.scenes = []*scene{}
	}
	sc := newScene(eng, eng.guid(), transform)
	eng.scenes = append(eng.scenes, sc)
	return sc
}

// SetOverlay marks a scene to be the one drawn after all the other
// screens.  This is expected to be a heads-up-display 2D scene.
func (eng *Eng) SetOverlay(s Scene) {
	if sc, _ := s.(*scene); sc != nil {
		eng.overlay = sc.uid
	}
}

// RemScene disposes given scene and everything within it. While generally scenes are
// created once and last until the application closes, applications may need to discard
// scenes that are no longer needed in order to manage resource consumption.
func (eng *Eng) RemScene(s Scene) {
	if sc, _ := s.(*scene); sc != nil {
		for index, existing := range eng.scenes {
			if sc.uid == existing.uid {
				existing.dispose()
				eng.scenes = append(eng.scenes[:index], eng.scenes[index+1:]...)
				return
			}
		}
	}
}

// 3D Direction constants. Primarily used for panning or rotating a camera view.
// See scene.PanView
const (
	XAxis = iota // affect only the X axis.
	YAxis        // affect only the Y axis.
	ZAxis        // affect only the Z axis.
)

// Eng
// ===========================================================================
// Director

// Director gives the application a chance to react at key moments.
// It is expected to be used by the application as follows:
//     eng, _ = vu.New("Title", 800, 600)
//     eng.SetDirector(app)
// where app implements the Director interface.
type Director interface {
	// Update is called many times each second to update application state.
	// The engine will use the updated state in the next render. It is expected
	// that this method returns relatively quickly.
	//
	// User input is provided as a series of strings where each string represents
	// a key or key-sequence that is currently being pressed by the user.
	Update(userInput []string, gameTime, deltaTime float32)

	// Resize notifies when the user has changed the rendering context area.
	// This will be called once each time the user alters the size of the
	// application window.
	Resize(x, y, width, height int)

	// Focus notifies when the window gains or loses focus.  Called when the
	// user switches to another application.
	Focus(focus bool)
}
