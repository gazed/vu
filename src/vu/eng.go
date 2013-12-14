// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

// Package vu (virtual universe) provides 3D application support. Vu wraps
// the individual subsystems like rendering, physics, data loading,
// audio, etc. to provide higher level functionality that includes:
//    • Scene graphs and composite objects.
//    • Timestepped update/render loop.
//    • Regular user input updates.
//    • Cameras and transform manipulation.
//
// The vu/eg (examples) package provides relatively small working code samples
// of engine functionality both for testing and demonstration purposes.
//
// Vu dependencies are:
//    • OpenGL for graphics card access.        See package vu/render.
//    • OpenAL for sound card access.           See package vu/audio.
//    • Cocoa  for OSX windowing and input.     See pacakge vu/device.
//    • WinAPI for Windows windowing and input. See package vu/device.
package vu

import (
	"time"
	"vu/audio"
	"vu/data"
	"vu/device"
	"vu/math/lin"
	"vu/move"
	"vu/render"
)

// Engine initializes and runs 3D application support. Interaction with the application
// is through the Director interface callbacks.
type Engine interface {
	SetDirector(director Director)    // Enable application callbacks.
	Action()                          // Kick off the main update loop.
	Shutdown()                        // Stop the engine and free allocated resources.
	AddScene(transform int) Scene     // Add a scene.
	RemScene(s Scene)                 // Remove a scene.
	SetOverlay(s Scene)               // Mark a scene as the overlay scene.
	Size() (x, y, width, height int)  // Get the current viewport size.
	Resize(x, y, width, height int)   // Resize the current viewport.
	Color(r, g, b, a float32)         // Set background clear colour.
	Enable(attr uint32, enabled bool) // Enable/disable global graphic attributes.
	ShowCursor(show bool)             // Hide or show the cursor.
	SetCursorAt(x, y int)             // Put the cursor at the given window location.

	// PlaceSoundListener sets the 3D location of the entity that can hear sounds.
	// Sounds that are played at other locations will be heard more faintly as the
	// distance between the played sound and listener increases.
	PlaceSoundListener(x, y, z float64)     // Create a sound listener.
	UseSound(sound string) audio.SoundMaker // Create a sound maker.
	Mute(mute bool)                         // Toggle game sound.
}

// Director gives the application a chance to react at key moments.
// It is expected to be used by the application as follows:
//     eng, _ = vu.New("Title", 800, 600) // App creates Engine.
//     eng.SetDirector(app)               // App registers as a Director.
type Director interface {
	// Create is called to populate the initial Scenes and Parts.
	Create(eng Engine)

	// Update is called many times a second to update application state.
	// The engine will use the updated state in the next render.
	// It is expected that this method returns quickly.
	//
	// User input is provided as a map of currently pressed keys,
	// mouse buttons and their pressed durations.
	Update(i *Input)
}

// Input is used to communicate current user input to the application.
// This gives the current cursor location, current pressed keys,
// mouse buttons, and modifiers.
//
// The map of keys and mouse buttons that are currently pressed also
// include how long they have been pressed in update ticks. A negative
// value indicates a release where the duration can be calculated by
// (RELEASED - duration).
type Input struct {
	Mx, My  int            // Current mouse location.
	Down    map[string]int // Pressed keys, buttons with duration.
	Shift   bool           // True if shift modifier is currently pressed.
	Control bool           // True if control modifier is currently pressed.
	Focus   bool           // True if window is in focus.
	Resized bool           // True if window was resized or moved.
	Dt      float64        // Delta time used for updates.
	Gt      float64        // Game time is the total number of updates.
}

// 3D Direction constants. Primarily used for panning or rotating a camera view.
// See Scene.PanView.
const (
	XAxis = iota // Affect only the X axis.
	YAxis        // Affect only the Y axis.
	ZAxis        // Affect only the Z axis.
)

// Global graphic state constants. These are attributes used in the
// Engine.Enable method.
const (
	BLEND = render.BLEND // Alpha blending.
	CULL  = render.CULL  // Backface culling.
	DEPTH = render.DEPTH // Z-buffer awareness.
)

// Engine, Director, and public API
// ===========================================================================
// engine implements Engine.

// Eng, the engine, is where everything starts. Eng is top of the hierarchy
// of the application window and provides access to the capabilities of the
// sub-components.
//
// Eng initializes the underlying subsystems and, for the most part, wraps
// access to subsystem functionality.
type engine struct {
	gc  render.Renderer     // Graphics card interface layer.
	ac  audio.Audio         // Audio card interface layer.
	dev device.Device       // Os specific window and rendering context.
	aud audio.SoundListener // Audio listener.
	res *roadie             // Data resource organizer.
	man *stage              // Stage manager.
}

// New creates a 3D engine. The expected usage is:
//      if eng, err = vu.New("Title", 100, 100, 800, 600); err != nil {
//          log.Printf("Failed to initialize engine %s", err)
//          return
//      }
//      defer eng.Shutdown() // Close down nicely.
//      eng.SetDirector(app) // Enable application callbacks.
//                           // Application 3D setup and initialization.
//      eng.Action()         // Run application loop (does not return).
func New(name string, x, y, width, height int) (e Engine, err error) {
	if name == "" {
		name = "Title"
	}
	if width < 100 {
		width = 100
	}
	if height < 100 {
		height = 100
	}
	eng := &engine{}

	// initialize the os specific shell, graphics context, and
	// user input monitor.
	eng.dev = device.New(name, x, y, width, height)

	// initialize the audio layer.
	eng.aud = audio.NewSoundListener()
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
	eng.dev.Open()
	e = eng
	return
}

// Shutdown stops the engine and frees up any allocated resources.
func (eng *engine) Shutdown() {
	if eng.ac != nil {
		eng.ac.Shutdown()
	}
	if eng.res != nil {
		eng.res.dispose()
	}
	if eng.dev != nil {
		eng.dev.Dispose()
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
func (eng *engine) Action() {
	ut := uint64(0) // update ticks counts the number of updates.

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
	for eng.dev.IsAlive() {

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
		for updateTime >= float64(dt) {
			eng.update(ut, dt)        // update state, physics and animations.
			updateTime -= float64(dt) // track the used delta time.
			ut += 1                   // track the total updates
		}

		// FUTURE interpolate the state based on the remaining delta time.  Right now
		//        the rendering is done on un-interpolated state which may be slightly
		//        behind where it should be.
		// interpolatedTime := updateTime / dt;  // fraction of unused delta time between 0 and 1.
		// State state = currentState*interpolatedTime + previousState * ( 1.0 - interpolatedTime );

		// redraw everything based on the current state.
		eng.render()
	}
}

// update delegates application state updates to Stage.
// update is expected to be called from the engine Action loop.
func (eng *engine) update(ut uint64, dt float64) {
	pressed := eng.dev.Update()
	eng.man.runUpdate(pressed, dt)
}

// render draws the currently visible parts. The visible list is updated
// as soon as it is availble. Render is called from the engine Action loop.
func (eng *engine) render() {
	eng.gc.Clear()
	is3D := true
	for _, vis := range eng.man.vis {
		eng.populate(vis)
		if is3D != !vis.Is2D {
			is3D = !vis.Is2D
			eng.gc.Enable(render.DEPTH, is3D)
		}
		eng.gc.Render(vis)
	}
	eng.dev.SwapBuffers()
}

// SetDirector establishes application callbacks.
func (eng *engine) SetDirector(director Director) {
	eng.man = newStage(director)
	director.Create(eng)
}

// populate ensures that the resources identified in the scene graph are properly loaded
// and made available for rendering. Populate will log developer errors if requested
// resources can not be located.
func (eng *engine) populate(vis *render.Vis) {
	if len(vis.ShaderName) > 0 && (vis.Shader == nil || vis.ShaderName != vis.Mesh.Name) {
		if vis.Shader == nil {
			vis.Shader = &data.Shader{}
		}
		eng.res.useShader(vis.ShaderName, &vis.Shader)
	}
	if len(vis.MatName) > 0 && (vis.Mat == nil || vis.MatName != vis.Mat.Name) {
		if vis.Mat == nil {
			vis.Mat = &data.Material{}
		}
		eng.res.useMaterial(vis.MatName, &vis.Mat)
	}
	if len(vis.TexName) > 0 && (vis.Tex == nil || vis.TexName != vis.Tex.Name) {
		if vis.Tex == nil {
			vis.Tex = &data.Texture{}
		}
		eng.res.useTexture(vis.TexName, &vis.Tex)
	}
	if len(vis.GlyphName) > 0 && (vis.Glyph == nil || vis.GlyphName != vis.Glyph.Name) {
		if vis.Glyph == nil {
			vis.Glyph = &data.Glyphs{}
		}
		eng.res.useGlyphs(vis.GlyphName, &vis.Glyph)
	}

	// Banner meshes are created instead of loaded, and will reuse available bindings.
	if vis.MeshName == "banner" {
		if vis.Mesh == nil {
			vis.Mesh = &data.Mesh{}
		}
		if vis.GlyphText != vis.GlyphPrev {
			vis.GlyphWidth = vis.Glyph.Panel(vis.Mesh, vis.GlyphText)
			eng.res.gc.BindGlyphs(vis.Mesh)
			vis.GlyphPrev = vis.GlyphText
		}
	} else if len(vis.MeshName) > 0 && (vis.Mesh == nil || vis.MeshName != vis.Mesh.Name) {

		// Regular meshes are loaded from disk and create a new binding.
		if vis.Mesh == nil {
			vis.Mesh = &data.Mesh{}
		}
		eng.res.useMesh(vis.MeshName, &vis.Mesh)
	}

	// If there is a material alpha then it overrides the default alpha value.
	// Non-default alpha values override the material alpha.
	if vis.Alpha == 1 && vis.Mat != nil {
		vis.Alpha = vis.Mat.Tr
	}
}

// ===========================================================================
// Expose/wrap device level information.

// GetSize returns the application viewport area in pixels.  This excludes any
// OS specific window trim.  The window x, y coordinates are the bottom left of
// the window.
func (eng *engine) Size() (x, y, width, height int) { return eng.dev.Size() }

// Resize needs to be called on window resize to adjust the graphics viewport.
func (eng *engine) Resize(x, y, width, height int) { eng.gc.Viewport(width, height) }

// ShowCursor hides and locks the cursor for the current window.
func (eng *engine) ShowCursor(show bool) { eng.dev.ShowCursor(show) }

// SetCursorAt puts the cursor at the given window location. Often this is used
// by the application when the cursor is hidden and the mouse movements are being
// tracked. Setting the cursor to the middle of the screen ensures movement doesn't
// get stuck at the screen edges.
func (eng *engine) SetCursorAt(x, y int) {
	eng.dev.SetCursorAt(x, y)
}

// RELEASED is used to indicate a released key or button.
const RELEASED = device.KEY_RELEASED

// ===========================================================================
// Expose/wrap graphic and audio controls.

// Color sets the default background clear color. This color will appear if nothing
// else is drawn over it.
func (eng *engine) Color(r, g, b, a float32) { eng.gc.Color(r, g, b, a) }

// Enable or disable global graphics attributes.
// Current valid values are: CULL, BLEND, DEPTH
func (eng *engine) Enable(attribute uint32, enabled bool) { eng.gc.Enable(attribute, enabled) }

// PlaceSoundListener sets the 3D location of the entity that can hear sounds.
// Sounds that are played at other locations will be heard more faintly as the
// distance between the played sound and listener increases.
func (eng *engine) PlaceSoundListener(x, y, z float64) { eng.aud.SetLocation(x, y, z) }

// UseSound creates a SoundMaker that is linked to the given sound resource.
func (eng *engine) UseSound(sound string) audio.SoundMaker {
	s := &data.Sound{}
	eng.res.useSound(sound, &s)
	return audio.NewSoundMaker(s)
}

// Mute turns the game sound on (mute == false) or off (mute == true).
func (eng *engine) Mute(mute bool) { eng.ac.Mute(mute) }

// RenderMatrix turns a math/lin matrix into a matrix that can be used
// by the render system. The input math matrix, mm, is used to fill the values
// in the given render matrix rm.  The updated rm matrix is returned.
func RenderMatrix(mm *lin.M4, rm *render.M4) *render.M4 {
	rm.X0, rm.Y0, rm.Z0, rm.W0 = float32(mm.X0), float32(mm.Y0), float32(mm.Z0), float32(mm.W0)
	rm.X1, rm.Y1, rm.Z1, rm.W1 = float32(mm.X1), float32(mm.Y1), float32(mm.Z1), float32(mm.W1)
	rm.X2, rm.Y2, rm.Z2, rm.W2 = float32(mm.X2), float32(mm.Y2), float32(mm.Z2), float32(mm.W2)
	rm.X3, rm.Y3, rm.Z3, rm.W3 = float32(mm.X3), float32(mm.Y3), float32(mm.Z3), float32(mm.W3)
	return rm
}

// ===========================================================================
// Expose/wrap physics.

// Box creates a box shaped physics body located at the origin.
// The box size is w=2*hx, h=2*hy, d=2*hz.
func Box(hx, hy, hz float64) move.Body {
	return move.NewBody(move.NewBox(hx, hy, hz))
}

// Sphere creates a ball shaped physics body located at the origin.
// The sphere size is set by the radius.
func Sphere(radius float64) move.Body {
	return move.NewBody(move.NewSphere(radius))
}

// ===========================================================================
// Expose/wrap scene manager.

func (eng *engine) AddScene(transform int) Scene { return eng.man.addScene(transform) }
func (eng *engine) RemScene(s Scene)             { eng.man.remScene(s) }
func (eng *engine) SetOverlay(s Scene)           { eng.man.setOverlay(s) }
