// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// eng.go holds all the component managers and runs application state updates.
//        This is the main user facing class.
// DESIGN: keep small by delegating application requests to the components.

import (
	"log"
	"time"

	"github.com/gazed/vu/physics"
)

// Eng provides support for a 3D application conforming to the App interface.
// Eng provides the root transform hierarchy node, a point-of-view (Pov) where
// models, physics bodies, and noises are attached and processed each update.
// Eng is also used to set global engine state and provides top level timing
// statistics.
type Eng interface {
	Shutdown()     // Stop the engine and free allocated resources.
	Reset()        // Put the engine back to its initial state.
	State() *State // Query engine state. State updated per tick.
	Root() *Pov    // Single root of the transform hierarchy.

	// Physics returns the vu/physics system manager.
	Physics() physics.Physics // Allows setting of physics attributes.

	// Set changes engine wide attributes. It accepts one or more
	// functions that take an EngAttr parameter, ie: vu.Color(0,0,0).
	Set(...EngAttr) // Update one or more engine attributes.

	// Timing gives application feedback. It is updated each processing loop.
	// The returned update times should be averaged over multiple calls.
	Usage() *Timing // Per update loop performance metrics.
}

// App is the application callback interface to the engine. It is implemented
// by the application and registered once on engine creation as follows:
//     err := vu.New(app, "Title", 0, 0, 800, 600) // Reg. app in new Eng.
// The App communicates with the engine using the Eng methods.
type App interface {
	Create(eng Eng, s *State) // Called once after successful startup.

	// Update allows applications to change state prior to the next render.
	// Update is called many times a second after the initial call to Create.
	//    i : user input refreshed prior to each call.
	//    s : engine state refreshed prior to each call.
	Update(eng Eng, i *Input, s *State) // Process user input, update state.
}

// Eng and App interfaces.
// ===========================================================================
// engine implements Eng.

// engine controls application communication and state updates.
// It is also an entity manager in that it uses a unique entity id
// to group application object instances by component functionality.
// Engine relies on helper classes for the majority of the work:
//     An asset manager for loading application assets.
//     A physics manager for handling forces and collisions.
//     A scene manager for rendering frames.
// Engine expects to be started as a go-routine using the runEngine method.
type engine struct {
	alive bool     // True until application decides otherwise.
	data  *appData // Combination user input and application state.

	// Communicate with the vu device goroutine.
	machine   chan msg   // Communicate with device loop.
	bindReply chan error // Receive machine bind responses.
	stop      chan bool  // Closed or any value to stop machine.

	// Communicate with the asset load manager goroutine.
	load     chan map[aid]string // Request asset load.
	loaded   chan map[aid]asset  // Receive newly loaded assets.
	stopLoad chan bool           // Send or close to stop loader.

	// Application entities are grouped into components.
	// All entities are Pov (location:orientation) based.
	ids    *eids   // Entity id manager.
	povs   *povs   // Entity transform component.
	cams   *cams   // Camera component.
	models *models // Visible component.
	bodies *bodies // Physic component.
	frames *frames // Render component.
	sounds *sounds // Audio component.
	lights *lights // Light component.
	layers *layers // Pre-render-pass component
	times  *Timing // Update loop timing statistics.
}

// newEngine is expected to be called once on startup
// from the runEngine() method.
func newEngine(machine chan msg) *engine {
	eng := &engine{alive: true, machine: machine}
	eng.ids = &eids{}
	eng.data = newAppData()
	eng.povs = newPovs()
	eng.times = &Timing{}
	eng.Reset() // allocate data components.

	// init comunications with the load goroutine.
	eng.load = make(chan map[aid]string)
	eng.loaded = make(chan map[aid]asset)
	eng.stopLoad = make(chan bool)

	// used to synchronize bind requests with the machine.
	eng.bindReply = make(chan error)
	return eng
}

// Reset removes all entities and sets the engine back to
// its initial state. This allows the application to put the
// engine back in a clean state without restarting.
func (eng *engine) Reset() {
	if eng.root() != nil {
		eng.disposePov(eng.root().id)
	}
	eng.ids.reset()
	eng.povs.reset()
	eng.povs.create(eng, eng.ids.create(), nil) // root
	eng.cams = newCams()
	eng.models = newModels(eng)
	eng.lights = newLights()
	eng.layers = newLayers(eng)
	eng.sounds = newSounds(eng)
	eng.sounds.setListener(eng.povs.get(0))
	eng.bodies = newBodies()
	eng.frames = newFrames()
}

// Shutdown is a user request to close down the engine.
// Expected to be called once on Application exit.
func (eng *engine) Shutdown() {
	eng.alive = false
	eng.disposePov(eng.root().id)
	if eng.machine != nil {
		eng.stopLoad <- true
		eng.machine <- &shutdown{}
	}
}

// main application loop timing constants.
const (
	// delta time is how often the state is updated. It is fixed at
	// 50 times a second (50/1s = 20ms) so that the game speed is
	// constant (independent from computer speed and refresh rate).
	dt = time.Duration(20 * time.Millisecond) // 0.02s, 50fps

	// render time limits how often render frames are generated.
	// Most modern flat monitors are 60fps.
	rt = time.Duration(10 * time.Millisecond) // 0.01s, 100fps

	// capTime guards against slow updates and the spiral of death.
	// Ignore any updating and rendering time that was more than 200ms.
	capTime = time.Duration(200 * time.Millisecond) // 0.2s
)

// runEngine is the main application timing loop. It calls Create once
// on startup and Update on a regular basis. The application callbacks
// allows the application to initiate object creation for rendering and
// to consume user input from device polling.
func runEngine(app App, wx, wy, ww, wh int,
	machine chan msg, draw chan frame, stop chan bool) {
	defer catchErrors()
	eng := newEngine(machine)
	go runLoader(eng.machine, eng.load, eng.loaded, eng.stopLoad)
	eng.frames.draw = draw
	eng.stop = stop
	eng.data.state.setScreen(wx, wy, ww, wh)
	app.Create(eng, eng.data.state)
	eng.layers.enableShadows() // Call before other load requests.
	ut := uint64(0)            // kick off initial update...
	eng.update(app, dt, ut)    // queue the initial load asset requests.

	// Initialize timers and kick off the main control timing loop.
	loopStart := time.Now()
	var updateStart time.Time
	var timeUsed time.Duration
	var updateTimer time.Duration // Track when to trigger an update.
	var renderTimer time.Duration // Track when to trigger a render.
	for eng.alive {
		timeUsed = time.Since(loopStart) // Count previous loop.
		eng.times.Elapsed += timeUsed    // Track total time.
		if timeUsed > capTime {          // Avoid slow update death.
			timeUsed = capTime
		}
		loopStart = time.Now()

		// Trigger update based on current elapsed time.
		// This advances state at a constant rate (dt).
		updateTimer += timeUsed
		for updateTimer >= dt {
			updateStart = time.Now() // Time the update.
			ut++                     // Track the total update ticks.
			updateTimer -= dt        // Remove delta time used.

			// Perform the update, preparing the next render frame.
			eng.update(app, dt, ut) // Update state, physics, etc.
			if eng.alive {          // Application may have quit.
				eng.frames.drawFrame(eng)
			}

			// Reset and start counting times for the next update.
			eng.times.Zero()
			eng.times.Update += time.Since(updateStart)
		}

		// A render frame request is sent to the machine. Redraw everything, using
		// interpolation when there is no new frame. Ignore excess render time.
		renderTimer += timeUsed
		if renderTimer >= rt {
			eng.times.Renders++

			// Interpolation is the fraction of unused delta time between 0 and 1.
			// ie: State state = currentState*interpolation + previousState * (1.0 - interpolation);
			interpolation := updateTimer.Seconds() / dt.Seconds()
			eng.frames.render(eng.machine, interpolation, ut)
			renderTimer = renderTimer % rt // drop extra render time.
		}
		eng.communicate() // process go-routine messages.
	}
}

// communicate processes all go-routine channels. Must be non-blocking.
// Incoming messages are generally responses to asset loading requests
// initiated by engine.
func (eng *engine) communicate() {
	select {
	case <-eng.stop: // closed channels return 0
		eng.stopLoad <- true // Tell the loader to stop.
		return               // Device/window has closed.
	case assets := <-eng.loaded:
		eng.models.finishLoads(assets)
		eng.sounds.finishLoads(assets)
	default:
		// don't block when there are no channels to process.
	}
}

// update polls user input, runs physics, calls application update,
// and finally refreshes all models resulting in updated transforms.
// The transform hierarchy is now ready to generate a render frame.
func (eng *engine) update(app App, dt time.Duration, ut uint64) {

	// Fetch input from the device thread. Essentially a sequential call.
	eng.machine <- eng.data // blocks until processed by the server.
	<-eng.data.reply        // blocks until processing is finished.
	input := eng.data.input // User input has been refreshed.
	state := eng.data.state // Engine state has been refreshed.
	dts := dt.Seconds()     // delta time as float.

	// update the location and orientation of any physics bodies.
	eng.bodies.stepVelocities(eng, dts) // Marks povs as dirty.

	// Have the application adjust any or all state before rendering.
	input.Dt = dts                // how long to get back to here.
	input.Ut = ut                 // update ticks.
	app.Update(eng, input, state) // application to updates its own state.

	// update assets that the application changed or which need
	// per tick processing. Per-ticks include animated models,
	// particle effects, surfaces, phrases, ...
	if eng.alive {
		eng.models.refresh(dts) // check for new load requests.
		eng.sounds.refresh()    // check for new load requests.
		eng.povs.updateWorldTransforms()
		eng.sounds.repositionSoundListener()
	}
}

// release sends a release resource request to the machine.
// Expected to be run as a goroutine so that its this method that blocks
// until the machine is ready to process it.
func (eng *engine) release(rd *releaseData) { eng.machine <- rd }

// rebind sends a bind request to the machine and waits for
// the response - making it a synchronous call.
//
// Note that rebinding mesh, texture, and sound does not change
// any model data, it just updates the data on the graphics or audio
// card. This way the model can continue to be rendered while the
// data is being rebound: the binding goroutine is the same as the
// render goroutine, so it is doing one or the other.
func (eng *engine) rebind(assets []asset) {
	eng.machine <- &bindData{data: assets, reply: eng.bindReply} // request binds.
	if err := <-eng.bindReply; err != nil {                      // wait for binds.
		log.Printf("%s", err)
	}
}

// submitLoadReqs is called on the engine processing goroutine to send
// the load requests off for loading. A new go routine is started so
// that engine is not blocked while the load request waits for processing.
func (eng *engine) submitLoadReqs(reqs map[aid]string) {
	go func(load chan map[aid]string, reqs map[aid]string) {
		load <- reqs
	}(eng.load, reqs)
}

// clampTex requests a texture to be non-repeating.
// Expected to be called once when setting up a texture.
func (eng *engine) clampTex(tid uint32) {
	eng.machine <- &clampTex{tid: tid}
}

// State provides access to current engine state.
func (eng *engine) State() *State { return eng.data.state }

// Implement Eng interface. Returns the top of the transform hierarchy.
func (eng *engine) Root() *Pov { return eng.root() }

// Implement Eng interface. Returns the physics instance.
func (eng *engine) Physics() physics.Physics { return eng.bodies.physics }

// pov entities. newPov can only be called from an existing Pov
// so parent is never nil.
func (eng *engine) newPov(parent *Pov) *Pov {
	return eng.povs.create(eng, eng.ids.create(), parent)
}
func (eng *engine) root() *Pov { return eng.povs.get(0) }

// camera entities.
func (eng *engine) cam(id eid) *Camera    { return eng.cams.get(id) }
func (eng *engine) newCam(id eid) *Camera { return eng.cams.create(id) }

// model entities.
func (eng *engine) model(id eid) Model { return eng.models.get(id) }
func (eng *engine) newModel(id eid, shader string, assets ...string) Model {
	return eng.models.create(id, shader, assets...)
}

// light entities.
func (eng *engine) light(id eid) *Light    { return eng.lights.get(id) }
func (eng *engine) newLight(id eid) *Light { return eng.lights.create(id) }

// layer render pass entities.
func (eng *engine) layer(id eid) Layer              { return eng.layers.get(id) }
func (eng *engine) newLayer(id eid, attr int) Layer { return eng.layers.create(id, attr) }

// body: physics entities.
func (eng *engine) body(id eid) physics.Body { return eng.bodies.get(id) }
func (eng *engine) newBody(id eid, b physics.Body, p *Pov) physics.Body {
	return eng.bodies.create(id, b, p.T)
}
func (eng *engine) setSolid(id eid, mass, bounce float64) {
	eng.bodies.solidify(id, mass, bounce)
}

// sound: audio entities.
func (eng *engine) addSound(id eid, name string) { eng.sounds.create(id, name) }
func (eng *engine) playSound(id eid, index int)  { eng.sounds.play(id, index) }

// FUTURE: cleaning up resources is not complete. Dispose currently means
// removing entities from the Pov hierarchy and from the eng entity manager,
// yet keeps them in the cache and bound on the GPU/Snd devices. Applications
// often "dispose" parts of the Pov hierarchy only to (re)use the underlying
// data again. Likely need another API for the application to indicate data
// that is to be completely "unloaded".
//
// Note: cached objects know about "loaded" and must be kept in sync
// if device bound data is changed.

// dispose discards the given pov component or the entire pov and all
// its components. Each call recalculates the currently loaded set
// of assets.
func (eng *engine) dispose(id eid, component int) {
	switch component {
	case PovBody:
		eng.bodies.dispose(id)
	case PovCam:
		eng.cams.dispose(id)
	case PovModel:
		eng.models.dispose(id)
	case PovSound:
		eng.sounds.dispose(id)
	case PovLight:
		eng.lights.dispose(id)
	case PovLayer:
		eng.layers.dispose(id)
	case PovNode:
		eng.disposePov(id)
	}
}

// disposePov chops this transform and all of its children out
// of the transform hierarchy. All associated objects are disposed.
func (eng *engine) disposePov(id eid) {
	eng.povs.dispose(id)
	eng.cams.dispose(id)
	eng.bodies.dispose(id)
	eng.models.dispose(id)
	eng.sounds.dispose(id)
	eng.lights.dispose(id)
	eng.layers.dispose(id)
}

// Usage returns numbers collected each time through the
// main processing loop. This allows the application to get
// a sense of time usage.
func (eng *engine) Usage() *Timing { return eng.times }

// Set one or more engine attributes.
func (eng *engine) Set(attrs ...EngAttr) {
	for _, attr := range attrs {
		attr(eng)
	}
}

// engine
// ===========================================================================
// expose/wrap physics shapes.

// NewBox creates a box shaped physics body located at the origin.
// The box size is given by the half-extents so that actual size
// is w=2*hx, h=2*hy, d=2*hz.
func NewBox(hx, hy, hz float64) physics.Body {
	return physics.NewBody(physics.NewBox(hx, hy, hz))
}

// NewSphere creates a ball shaped physics body located at the origin.
// The sphere size is defined by the radius.
func NewSphere(radius float64) physics.Body {
	return physics.NewBody(physics.NewSphere(radius))
}

// NewRay creates a ray located at the origin and pointing in the
// direction dx, dy, dz.
func NewRay(dx, dy, dz float64) physics.Body {
	return physics.NewBody(physics.NewRay(dx, dy, dz))
}

// SetRay updates the ray direction.
func SetRay(ray physics.Body, x, y, z float64) {
	physics.SetRay(ray, x, y, z)
}

// SetPlane updates the plane normal.
func SetPlane(plane physics.Body, x, y, z float64) {
	physics.SetPlane(plane, x, y, z)
}

// NewPlane creates a plane located on the origin and oriented by the
// plane normal nx, ny, nz.
func NewPlane(nx, ny, nz float64) physics.Body {
	return physics.NewBody(physics.NewPlane(nx, ny, nz))
}

// Cast checks if a ray r intersects the given Body b, returning the
// nearest point of intersection if there is one. The point of contact
// x, y, z is valid when hit is true.
func Cast(ray, b physics.Body) (hit bool, x, y, z float64) {
	return physics.Cast(ray, b)
}

// engine physics
// ===========================================================================
// engine attributes reeduces the eng API footprint using functional options:
//    http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
//    https://commandcenter.blogspot.ca/2014/01/self-referential-functions-and-design.html
// See eng_test notes about slowness. Usable on engine since the messages are
// being sent across a channel. Currently the messages are blocking.

// EngAttr defines an engine attribute that can be used in Eng.Set().
type EngAttr func(Eng)

// Color sets the background window clear color.
// Engine attribute expected to be used in Eng.Set().
func Color(r, g, b, a float32) EngAttr {
	return func(e Eng) {
		e.(*engine).machine <- &setColor{r: r, g: g, b: b, a: a}
	}
}

// CursorOn hides or shows the cursor.
// Engine attribute expected to be used in Eng.Set().
func CursorOn(show bool) EngAttr {
	return func(e Eng) {
		e.(*engine).machine <- &showCursor{enable: show}
	}
}

// CursorAt places the cursor at the window pixel x,y.
// Engine attribute expected to be used in Eng.Set().
func CursorAt(x, y int) EngAttr {
	return func(e Eng) {
		e.(*engine).machine <- &setCursor{cx: x, cy: y}
	}
}

// On enables/disables render attributes like Blend, CullFace, etc...
// Engine attribute expected to be used in Eng.Set().
func On(attr uint32, enabled bool) EngAttr {
	return func(e Eng) {
		e.(*engine).machine <- &enableAttr{attr: attr, enable: enabled}
	}
}

// ToggleFullScreen flips full screen and windowed mode.
// Engine attribute expected to be used in Eng.Set().
func ToggleFullScreen() EngAttr {
	return func(e Eng) { e.(*engine).machine <- &toggleScreen{} }
}

// Mute toggles the sound volume.
// Engine attribute expected to be used in Eng.Set().
func Mute(mute bool) EngAttr {
	gain := 1.0
	if mute {
		gain = 0.0
	}
	return func(e Eng) {
		e.(*engine).machine <- &setVolume{gain: gain}
	}
}

// Volume sets the sound volume.
// Engine attribute expected to be used in Eng.Set().
func Volume(zeroToOne float64) EngAttr {
	return func(e Eng) {
		e.(*engine).machine <- &setVolume{gain: zeroToOne}
	}
}

// Gravity changes the physics gravity constant.
// Engine attribute expected to be used in Eng.Set().
func Gravity(g float64) EngAttr {
	return func(e Eng) {
		e.(*engine).bodies.physics.Set(physics.Gravity(g))
	}
}
