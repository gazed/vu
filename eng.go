// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// Package vu, virtual universe, provides 3D application support. Vu wraps
// subsystems like rendering, physics, data loading, audio, etc. to provide
// higher level functionality that includes:
//    • Scene graphs and composite objects.
//    • Timestepped update/render loop.
//    • Callback access to user input events.
//    • Cameras and transform manipulation.
//    • Linking loaded assets to render and audio systems.
// Refer to the vu/eg examples package for working code samples that test
// and demo engine functionality.
//
// Vu dependencies are:
//    • OpenGL for graphics card access.        See package vu/render.
//    • OpenAL for sound card access.           See package vu/audio.
//    • Cocoa  for OSX windowing and input.     See pacakge vu/device.
//    • WinAPI for Windows windowing and input. See package vu/device.
package vu

import (
	"time"

	"github.com/gazed/vu/audio"
	"github.com/gazed/vu/device"
	"github.com/gazed/vu/move"
	"github.com/gazed/vu/render"
)

// Engine initializes and provides runtime suport for a 3D application.
// Interaction with the application is through the Director interface.
type Engine interface {
	SetDirector(d Director) // Enable application callbacks.
	Verify() error          // Check model data against shader expectations.
	Action()                // Kick off the main update loop.

	// Shuddown stops the engine and the underlying graphics context while
	// Reset removes all allocated resources while keeping the graphics context.
	Shutdown() // Stop the engine and free allocated resources.
	Reset()    // Put the engine back to its initial state.

	// The application window/viewport is queried and controlled as follows:
	Size() (x, y, width, height int)  // Get the current viewport size.
	Resize(x, y, width, height int)   // Resize the current viewport.
	Color(r, g, b, a float32)         // Set background clear colour.
	ShowCursor(show bool)             // Hide or show the cursor.
	SetCursorAt(x, y int)             // Place cursor at the x,y window location.
	Enable(attr uint32, enabled bool) // Enable/disable global graphic attributes.

	// Scenes group visible objects with a camera. Scenes are drawn in the
	// order they are created, unless modified using SetLastScene.
	AddScene(transform int) Scene // Add a scene.
	RemScene(s Scene)             // Remove and dispose a scene.
	SetLastScene(s Scene)         // Put the given scene last in the scene list.

	// PlaceSoundListener sets the 3D location of the entity that can hear sounds.
	// Sounds that are played at other locations will be heard more faintly as
	// the distance between the played sound and listener increases.
	PlaceSoundListener(x, y, z float64) // Create a sound listener.
	Mute(mute bool)                     // Toggle game sound.
}

// Director is the engine callback to the application.
// Director is expected to be implemented by the application
// and registered with the engine as follows:
//     eng, _ = vu.New("Title", 0,0,800,600) // App creates Engine.
//     eng.SetDirector(app)                  // App registers as a Director.
type Director interface {

	// Update allows applications to change state prior to the next render.
	// Update is called many times a second once the application calls eng.Action.
	// Applications commonly create some resources prior to starting Updates.
	Update(i *Input) // Application expected to return quickly.
}

// Engine constants used as input to various methods.
const (
	// Global graphic state constants. See Engine.Enable(const, bool).
	BLEND = render.BLEND // Alpha blending. Enabled by default.
	CULL  = render.CULL  // Backface culling. Enabled by default.
	DEPTH = render.DEPTH // Z-buffer awareness. Enabled by default.

	// 3D Direction constants. Primarily used for panning or
	// rotating a camera view. See Camera.Spin.
	XAxis = iota // Affect only the X axis.
	YAxis        // Affect only the Y axis.
	ZAxis        // Affect only the Z axis.

	// Camera transform choices. Used in Eng.AddScene(transform).
	VP    = iota // Perspective view transform.
	VO           // Orthographic view transform.
	VF           // First person view transform with up/down angle.
	XZ_XY        // Perspective to Ortho view transform.

	// Per-part rendering constants. See Role.SetDrawMode(mode int).
	TRIANGLES = render.TRIANGLES // Triangles are the norm.
	POINTS    = render.POINTS    // Used for particle effects.
	LINES     = render.LINES     // Used for drawing squares and boxes.

	// Texture rendering directives. See Role.SetTexMode()
	TEX_REPEAT = render.REPEAT // Repeat texture when UV greater than 1.

	// User input key released indicator. Total time down, in update
	// ticks, is key down ticks minus RELEASED. See Director.Update().
	RELEASED = device.KEY_RELEASED
)

// Engine, Director, and public API
// ===========================================================================
// engine implements Engine.

// The engine is where everything starts. Engine is top of the hierarchy
// of the application window and provides access to the capabilities of the
// sub-components.
//
// Engine initializes the underlying subsystems and, for the most part, wraps
// access to subsystem functionality.
type engine struct {
	gc     render.Renderer     // Graphics card interface layer.
	ac     audio.Audio         // Audio card interface layer.
	dev    device.Device       // Os specific window and rendering context.
	in     *Input              // Propogates device input to the application.
	aud    audio.SoundListener // Audio listener.
	assets *assets             // Data resource manager.
	mover  move.Mover          // Physics handles forces, collisions.
	stage  *stage              // Rendering culls and draws the scene graph.
	app    Director            // Application callbacks.
}

// New creates a 3D engine and application window. The expected usage is:
//      if eng, err = vu.New("Title", 100, 100, 800, 600); err != nil {
//          log.Printf("Failed to initialize engine %s", err)
//          return
//      }
//      defer eng.Shutdown() // Close down nicely.
//      eng.SetDirector(app) // Enable application update callbacks.
//         ....              // application initialization.
//      eng.Action()         // Start update callbacks (does not return).
// A miniumum window width of 100 and height of 100 is enforced.
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
	eng.in = &Input{Down: map[string]int{}}

	// initialize the os specific shell, graphics context, and
	// user input monitor.
	eng.dev = device.New(name, x, y, width, height)

	// initialize the audio layer.
	eng.ac = audio.New()
	eng.aud = eng.ac.NewSoundListener()
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
	eng.Enable(BLEND, true)
	eng.Enable(CULL, true)
	eng.Color(0, 0, 0, 1)
	eng.assets = newAssets(eng.ac, eng.gc)
	eng.stage = newStage(eng.assets)
	eng.mover = move.NewMover()
	eng.gc.Viewport(width, height)
	eng.dev.Open()
	return eng, err
}

// Shutdown stops the engine and frees up any allocated resources.
func (eng *engine) Shutdown() {
	if eng.stage != nil {
		eng.stage.dispose()
		eng.stage = nil
	}
	if eng.ac != nil {
		eng.ac.Shutdown()
		eng.ac = nil
	}
	if eng.dev != nil {
		eng.dev.Dispose()
		eng.dev = nil
	}
	eng.app = nil
}

// Reset cleans up graphics resources and puts the engine back to its initial
// state. All application created parts are destroyed which results in all
// resources being removed up as nothing is left that references them.
func (eng *engine) Reset() {
	if eng.stage != nil {
		eng.stage.dispose()
	}
	eng.stage = newStage(eng.assets)
}

// SetDirector establishes the application update callback receiver.
func (eng *engine) SetDirector(director Director) {
	eng.app = director
}

// Verify can be optionally called after SetDirector to check the initial
// resource loading and model creation.
func (eng *engine) Verify() error { return eng.stage.verify() }

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
	for eng.dev != nil && eng.dev.IsAlive() {

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

// ===========================================================================
// Start the real work of delegating the update and render calls down
// to the helper classes and systems.

// update delegates application state updates to the stage manager.
// update is expected to be called from the engine Action loop.
func (eng *engine) update(ut uint64, dt float64) {
	eng.in.convertInput(eng.dev.Update(), dt) // get user input.
	if eng.app != nil {
		eng.app.Update(eng.in) // applications turn for state updates.
		if eng.stage != nil {  // App may have shutdown engine.
			eng.mover.Step(eng.stage.bodies, dt) // update physics state.
			eng.stage.update(dt)                 // prepare render state.
		}
	}
}

// render delegates application rendering to the stage manager.
// Render is expected to be called only from the engine Action loop.
func (eng *engine) render() {
	if eng.stage != nil { // App may have shutdown engine.
		eng.stage.render(eng.gc)
		eng.dev.SwapBuffers()
	}
}

// Implements Engine.
func (eng *engine) AddScene(transform int) Scene { return eng.stage.addScene(transform) }
func (eng *engine) RemScene(s Scene)             { eng.stage.remScene(s) }

// SetLastScene is expected to be used for overlay scenes.
// It moves the indicated scene to be the last scene rendered.
func (eng *engine) SetLastScene(s Scene) { eng.stage.setLast(s) }

// ===========================================================================
// Expose/wrap device level information.

// GetSize returns the application viewport area in pixels.  This excludes any
// OS specific window trim.  The window x, y coordinates are the bottom left of
// the window.
func (eng *engine) Size() (x, y, width, height int) { return eng.dev.Size() }

// Resize needs to be called on window resize to adjust the graphics viewport.
// The engine starts the resize by informing the application during update,
// but leaves viewport resizing, using this method, under application control.
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
// distance between the played sound and listener increases. The location is
// often the same as the main camera.
func (eng *engine) PlaceSoundListener(x, y, z float64) { eng.aud.SetLocation(x, y, z) }

// Mute turns the game sound on (mute == false) or off (mute == true).
func (eng *engine) Mute(mute bool) { eng.ac.Mute(mute) }

// ===========================================================================
// Expose/wrap physics shapes.

// NewBox creates a box shaped physics body located at the origin.
// The box size is w=2*hx, h=2*hy, d=2*hz. Used in Part.SetBody()
func NewBox(hx, hy, hz float64) move.Body {
	return move.NewBody(move.NewBox(hx, hy, hz))
}

// NewSphere creates a ball shaped physics body located at the origin.
// The sphere size is defined by the radius. Used in Part.SetBody()
func NewSphere(radius float64) move.Body {
	return move.NewBody(move.NewSphere(radius))
}

// NewRay creates a ray located at the origin and pointing in the
// direction dx, dy, dz. Used in Part.SetForm()
func NewRay(dx, dy, dz float64) move.Body {
	return move.NewBody(move.NewRay(dx, dy, dz))
}

// SetRay updates the ray direction.
func SetRay(ray move.Body, x, y, z float64) {
	move.SetRay(ray, x, y, z)
}

// NewPlane creates a plane located on the origin and oriented by the
// plane normal nx, ny, nz. Used in Part.SetForm()
func NewPlane(nx, ny, nz float64) move.Body {
	return move.NewBody(move.NewPlane(nx, ny, nz))
}

// Cast checks if a ray r intersects the given Body b, returning the
// nearest point of intersection if there is one. The point of contact
// x, y, z is valid when hit is true.
func Cast(ray, b move.Body) (hit bool, x, y, z float64) { return move.Cast(ray, b) }
