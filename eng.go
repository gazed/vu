// Copyright © 2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"log"
	"math"
	"time"

	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/physics"
	"github.com/gazed/vu/render"
)

// Eng provides support for a 3D application conforming to the App interface.
// Eng holds the root transform hierarchy node, a point-of-view (Pov) where
// models, physics bodies, and noises are attached and processed each update.
// Eng is also used to set global engine state and provides top level timing
// statistics.
type Eng interface {
	Shutdown()     // Stop the engine and free allocated resources.
	Reset()        // Put the engine back to its initial state.
	State() *State // Query engine state. State updated per tick.
	Root() Pov     // Single root transform always exists.

	// Requests to change engine state.
	SetColor(r, g, b, a float32)      // Set background clear colour.
	ShowCursor(show bool)             // Hide or show the cursor.
	SetCursorAt(x, y int)             // Put cursor at the window pixel x,y.
	Enable(attr uint32, enabled bool) // Enable/disable render attributes.
	ToggleFullScreen()                // Flips full screen and windowed mode.
	Mute(mute bool)                   // Toggle sound volume.
	SetVolume(zeroToOne float64)      // Set sound volume.
	SetGravity(g float64)             // Change the gravity constant.

	// Collide checks for collision between two bodies independent
	// of the solver and without updating the the bodies locations.
	Collide(a, b physics.Body) bool

	// Timing is updated each processing loop. The returned update
	// times can flucuate and should be averaged over multiple calls.
	Usage() *Timing                // Per update loop performance metrics.
	Modelled() (models, verts int) // Total render models and verticies.
	Rendered() (models, verts int) // Rendered models and verticies.
}

// App is the application callback interface to the engine. It is implemented
// by the application and registered once on engine creation as follows:
//     err := vu.New(app, "Title", 0, 0, 800, 600) // Reg. app in new Eng.
// Note that it is safe to call Eng methods from goroutines.
type App interface {
	Create(eng Eng, s *State) // Called once after successfull startup.

	// Update allows applications to change state prior to the next render.
	// Update is called many times a second after the initial call to Create.
	//      i : user input refreshed prior to each call.
	//      s : engine state refreshed prior to each call.
	Update(eng Eng, i *Input, s *State) // Process user input.
}

// Eng and App interfaces.
// ===========================================================================
// engine implements Eng.

// engine controls application communication and state updates.
// It is also an entity manager in that it uses a unique entity id to
// group application object instances by component functionality.
// Engine relies on helper classes for the majority of the work:
//     An asset manager for loading application assets.
//     A physics manager for handling forces and collisions.
//     A scene manager for rendering frames.
// Engine expects to be started as a go-routine using the runEngine method.
type engine struct {
	alive   bool               // True until application decides otherwise.
	machine chan msg           // Communicate with device loop.
	stop    chan bool          // Closed or any value means stop the engine.
	data    *appData           // Combination user input and application state.
	sm      *scene             // Scene manager. Creates render frames.
	frame   []render.Draw      // update frame for next frame.
	uf      chan []render.Draw // next update frame returned from machine.
	physics physics.Physics    // Physics manager. Handles forces, collisions.

	// Asset manager. Handles loading assets concurrently.
	loader *loader         // Asset manager.
	loaded chan []*loadReq // Receive loaded models and noises.

	// Sounds are heard by the sound listener at an app set pov.
	soundListener *pov    // Current location of the sound listener.
	sx, sy, sz    float64 // Last location of the sound listener.

	// Group the application entities by component.
	// All entities are Pov (location:orientation) based.
	eid    uint64                  // Next entity id.
	povs   map[uint64]*pov         // Entity transforms.
	cams   map[uint64]*camera      // Camera components.
	models map[uint64]*model       // Visible components.
	lights map[uint64]*light       // Light components.
	noises map[uint64]*noise       // Audible components.
	layers map[uint64]*layer       // (Pre) Render pass components.
	bodies map[uint64]physics.Body // Non-colliding physic components.
	solids map[uint64]physics.Body // Colliding physic components.
	bods   []physics.Body          // Set from solids each update.
	times  *Timing                 // Loop timing statistics.
}

// newEngine is expected to be called once on startup
// from the runEngine() method.
func newEngine(machine chan msg) *engine {
	eng := &engine{alive: true, machine: machine}
	eng.data = newAppData()
	eng.times = &Timing{}
	eng.frame = []render.Draw{}
	eng.Reset()

	// helpers that create and update state.
	eng.physics = physics.NewPhysics()
	eng.loaded = make(chan []*loadReq)
	eng.loader = newLoader(eng.loaded, machine)
	eng.sm = newScene()
	return eng
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
	machine chan msg, uf chan []render.Draw, stop chan bool) {
	defer catchErrors()
	eng := newEngine(machine)
	go eng.loader.runLoader()
	eng.uf = uf
	eng.stop = stop
	eng.data.state.setScreen(wx, wy, ww, wh)
	app.Create(eng, eng.data.state)
	eng.sm.init(eng)
	ut := uint64(0)         // kick off initial update...
	eng.update(app, dt, ut) // queue the initial load asset requests.

	// Initialize timers and kick off the main control timing loop.
	var loopStart time.Time = time.Now()
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
			ut += 1                  // Track the total update ticks.
			updateTimer -= dt        // Remove delta time used.

			// Perform the update, preparing the next render frame.
			eng.update(app, dt, ut) // Update state, physics, etc.
			if eng.alive {          // Application may have quit.
				eng.frame = eng.sm.snapshot(eng, eng.frame)
			}

			// Reset and start counting times for the next update.
			eng.times.Zero()
			eng.times.Update += time.Since(updateStart)
		}

		// A render frame request is sent to the machine. Redraw everything, using
		// interpolation when there is no new frame. Ignore excess render time.
		renderTimer += timeUsed
		if renderTimer >= rt {
			eng.times.Renders += 1

			// Interpolation is the fraction of unused delta time between 0 and 1.
			// ie: State state = currentState*interpolation + previousState * (1.0 - interpolation);
			interpolation := updateTimer.Seconds() / dt.Seconds()
			if len(eng.frame) > 0 {
				eng.machine <- &renderFrame{frame: eng.frame, interp: interpolation, ut: ut}
				eng.frame = <-eng.uf      // immediately get next render frame
				eng.frame = eng.frame[:0] // ... and mark it as unpreprepared.
			} else {
				eng.machine <- &renderFrame{frame: nil, interp: interpolation, ut: ut}
			}
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
		eng.loader.shutdown() // Tell the loader to stop.
		return                // Device/window has closed.
	case loaded := <-eng.loaded:
		for _, req := range loaded {
			if req.err != nil {
				log.Printf("load error: %s", req.err)
				continue
			}
			switch a := req.a.(type) {
			case *mesh:
				if m, ok := req.data.(*model); ok {
					m.msh = a
				}
			case *texture:
				if m, ok := req.data.(*model); ok {
					if req.index < len(m.texs) {
						m.texs[req.index] = a
					}
				}
			case *shader:
				if m, ok := req.data.(*model); ok {
					m.shd = a
				}
			case *font:
				if m, ok := req.data.(*model); ok {
					m.fnt = a
					m.fnt.loaded = true
					if len(m.phrase) > 0 {
						m.phraseWidth = m.fnt.setPhrase(m.msh, m.phrase)
					}
				}
			case *animation:
				if m, ok := req.data.(*model); ok {
					m.anm = a
					m.msh = req.msh
					if req.index < len(m.texs) && len(req.texs) > 0 {
						m.texs[req.index] = req.texs[0]
					}
					m.nFrames = a.maxFrames(0)
					m.pose = make([]lin.M4, len(a.joints))
				}
			case *material:
				if m, ok := req.data.(*model); ok {
					m.mat = a
					if m.alpha == 1.0 {
						m.alpha = a.tr // Copy values so they can be set per model.
					}
					if m.kd.isBlack() {
						m.kd = a.kd // Copy values so they can be set per model.
					}
					m.ks = a.ks // Can't currently be overridden on model.
					m.ka = a.ka // ditto
				}
			case *sound:
				if n, ok := req.data.(*noise); ok {
					n.snds[req.index] = a
					n.loaded = true
				}
			default:
				log.Printf("engine: unknown asset type %T", a)
			}
		}
	default:
		// no channels to process.
	}
}

// update polls user input, runs physics, calls application update, and
// finally refreshes all models resulting in updated transforms.
// The transform hierarchy is now ready to generate a render frame.
func (eng *engine) update(app App, dt time.Duration, ut uint64) {

	// Fetch input from the device thread. Essentially a sequential call.
	eng.machine <- eng.data // blocks until processed by the server.
	<-eng.data.reply        // blocks until processing is finished.
	input := eng.data.input // User input has been refreshed.
	state := eng.data.state // Engine state has been refreshed.
	dts := dt.Seconds()     // delta time as float.

	// Run physics on all the bodies; adjusting location and orientation.
	eng.bods = eng.bods[:0] // reset keeping capacity.
	for _, bod := range eng.solids {
		eng.bods = append(eng.bods, bod)
	}
	eng.physics.Step(eng.bods, dts)

	// Have the application adjust any or all state before rendering.
	input.Dt = dts                // how long to get back to here.
	input.Ut = ut                 // update ticks.
	app.Update(eng, input, state) // application to updates its own state.

	// update assets that the application changed or which need
	// per tick processing. Per-ticks include animated models,
	// particle effects, surfaces, phrases, ...
	if eng.alive {
		eng.updateModels(dts)                // load and bind updated data.
		eng.placeModels(eng.root(), lin.M4I) // update all transforms.
		eng.updateSoundListener()            // reposition sound listener.
	}
}

// updateModels processes any ongoing model updates like animated models
// and CPU particle effects. Any new models are sent off for loading
// and any updated models generate data rebind requests.
func (eng *engine) updateModels(dts float64) {
	for eid, m := range eng.models {
		if len(m.loads) > 0 { // load model assets if necessary.
			eng.loader.queueLoads(m.loads)
			m.loads = m.loads[:0]
		} else if m.loaded() {
			// Handle model data changes from either the Application or
			// from effects, phrase updates, and animations.

			// handle any data updates with rebind requests.
			if pv, ok := eng.povs[eid]; ok && pv.visible {
				if m.effect != nil {
					// udpate and rebind particle effects which can
					// change mesh data.
					m.effect.update(m, dt.Seconds())
				}
				if !m.msh.bound {
					eng.rebind(m.msh)
					m.msh.bound = true
				}
				for _, tex := range m.texs {
					if !tex.bound {
						eng.rebind(tex)
						tex.bound = true
					}
				}
				if m.anm != nil {
					// animations update the bone position matricies.
					// These are bound as uniforms at draw time.
					m.animate(dts)
				}
			}
		}
	}
	for _, n := range eng.noises {
		if len(n.loads) > 0 { // load noise sounds if necessary.
			eng.loader.queueLoads(n.loads)
			n.loads = n.loads[:0]
		}
	}
	eng.loader.loadQueued()
}

// placeModels walks the transform hierarchy updating all the model
// transforms. This is called before rendering passes are done.
func (eng *engine) placeModels(p *pov, parent *lin.M4) {
	p.mm.SetQ(p.rot.Inv(p.at.Rot)) // invert model rotation.
	p.mm.ScaleSM(p.Scale())        // scale is applied first (on left of rotation)
	l := p.at.Loc
	p.mm.TranslateMT(l.X, l.Y, l.Z) // translate is applied last (on right of rotation).
	p.mm.Mult(p.mm, parent)         // model transform + parent transform
	for _, child := range p.children {
		eng.placeModels(child, p.mm) // recursive traversal.
	}
}

// updateSoundListener checks and updates the sound listeners location.
func (eng *engine) updateSoundListener() {
	x, y, z := eng.soundListener.Location()
	if x != eng.sx || y != eng.sy || z != eng.sz {
		eng.sx, eng.sy, eng.sz = x, y, z
		go func(x, y, z float64) {
			eng.machine <- &placeListener{x: x, y: y, z: z}
		}(x, y, z)
	}
}

// release sends a release resource request to the machine.
// Expected to be run as a goroutine so that it can block on the
// send until the machine is ready to process it.
func (eng *engine) release(rd *releaseData) { eng.machine <- rd }

// rebind sends a bind request to the machine and waits for
// the response - making it a synchronous call.
//
// Note that rebinding mesh, texture, and sound does not change
// any model data, it just updates the data on the graphics or audio
// card. This way the model can continue to be rendered while the
// data is being rebound (the binding goroutine is the same as the
// render goroutine, so it is doing one or the other).
func (eng *engine) rebind(data interface{}) {
	bindReply := make(chan error)
	eng.machine <- &bindData{data: data, reply: bindReply} // request bind.
	if err := <-bindReply; err != nil {                    // wait for bind.
		log.Printf("%s", err)
	}
}

// Shutdown is a user request to close down the engine.
// Expected to be called once on Application exit.
func (eng *engine) Shutdown() {
	eng.alive = false
	eng.dispose(eng.root(), POV)
	if eng.machine != nil {
		eng.loader.shutdown()
		eng.machine <- &shutdown{}
	}
}

// Reset removes all entities and sets the engine back to
// its initial state. This allows the application to put the
// engine back in a clean state without restarting.
func (eng *engine) Reset() {
	eng.dispose(eng.root(), POV)
	eng.povs = map[uint64]*pov{}
	eng.cams = map[uint64]*camera{}
	eng.models = map[uint64]*model{}
	eng.lights = map[uint64]*light{}
	eng.layers = map[uint64]*layer{}
	eng.noises = map[uint64]*noise{}
	eng.bodies = map[uint64]physics.Body{}
	eng.solids = map[uint64]physics.Body{}
	eng.eid = 1
	eng.povs[eng.eid] = newPov(eng, eng.eid) // root
	eng.soundListener = eng.povs[eng.eid]
}

// State provides access to current engine state.
func (eng *engine) State() *State { return eng.data.state }

// genid returns the next unique entity id. It craps out and
// starts returning 0 after generating all possible ids.
func (eng *engine) genid() uint64 {
	if eng.eid == math.MaxUint64 {
		return 0
	}
	eng.eid += 1 // first valid id is 1.
	return eng.eid
}

// Implement Eng interface. Returns the top of the transform hierarchy.
func (eng *engine) Root() Pov { return eng.root() }

// pov entities.
func (eng *engine) root() *pov { return eng.povs[1] }
func (eng *engine) newPov(p Pov) Pov {
	if parent, ok := p.(*pov); ok && parent != nil {
		p := newPov(eng, eng.genid())
		eng.povs[p.eid] = p
		p.parent = parent
		parent.children = append(parent.children, p)
		return p
	}
	return nil
}

// camera entities.
func (eng *engine) cam(p Pov) Camera {
	if pv, ok := p.(*pov); ok && pv != nil {
		if cam, ok := eng.cams[pv.eid]; ok {
			return cam
		}
	}
	return nil
}
func (eng *engine) newCam(p Pov) Camera {
	if pv, ok := p.(*pov); ok && pv != nil {
		c := newCamera()
		eng.cams[pv.eid] = c
		return c
	}
	return nil
}

// model entities.
func (eng *engine) model(p Pov) Model {
	if pv, ok := p.(*pov); ok && pv != nil {
		if model, ok := eng.models[pv.eid]; ok {
			return model
		}
	}
	return nil
}
func (eng *engine) newModel(p Pov, shader string) Model {
	if pv, ok := p.(*pov); ok && pv != nil {
		if _, ok := eng.models[pv.eid]; !ok {
			m := newModel(shader)
			eng.models[pv.eid] = m
			return m
		}
	}
	return nil
}

// light entities.
func (eng *engine) light(p Pov) Light {
	if pv, ok := p.(*pov); ok && pv != nil {
		if l, ok := eng.lights[pv.eid]; ok {
			return l
		}
	}
	return nil
}
func (eng *engine) newLight(p Pov) Light {
	if pv, ok := p.(*pov); ok && pv != nil {
		if _, ok := eng.lights[pv.eid]; !ok {
			l := newLight()
			eng.lights[pv.eid] = l
			return l
		}
	}
	return nil
}

// snap shot entities.
func (eng *engine) layer(p Pov) Layer {
	if pv, ok := p.(*pov); ok && pv != nil {
		if l, ok := eng.layers[pv.eid]; ok {
			return l
		}
	}
	return nil
}
func (eng *engine) newLayer(p Pov, attr int) Layer {
	if pv, ok := p.(*pov); ok && pv != nil {
		if _, ok := eng.layers[pv.eid]; !ok {
			l := newLayer(attr)
			eng.loader.bindLayer(l) // synchronously create and bind a fbo.
			eng.layers[pv.eid] = l
			return l
		}
	}
	return nil
}

// body: physics entities.
func (eng *engine) body(p Pov) physics.Body {
	if pv, ok := p.(*pov); ok && pv != nil {
		if body, ok := eng.bodies[pv.eid]; ok {
			return body
		}
		if body, ok := eng.solids[pv.eid]; ok {
			return body
		}
	}
	return nil
}
func (eng *engine) newBody(p Pov, b physics.Body) physics.Body {
	if pv, ok := p.(*pov); ok && pv != nil {
		if _, ok := eng.bodies[pv.eid]; !ok {
			b.SetWorld(pv.at)
			eng.bodies[pv.eid] = b
			return b
		}
	}
	return nil
}

// solid: physics entities.
func (eng *engine) setSolid(p Pov, mass, bounce float64) {
	if pv, ok := p.(*pov); ok && pv != nil {
		if b, okb := eng.bodies[pv.eid]; okb {
			b.SetMaterial(mass, bounce)
			eng.solids[pv.eid] = b
			delete(eng.bodies, pv.eid)
		}
	}
}

// noise: audio entities.
func (eng *engine) noise(p Pov) Noise {
	if pv, ok := p.(*pov); ok && pv != nil {
		if noise, ok := eng.noises[pv.eid]; ok {
			return noise
		}
	}
	return nil
}
func (eng *engine) newNoise(p Pov) Noise {
	if pv, ok := p.(*pov); ok && pv != nil {
		if _, ok := eng.noises[pv.eid]; !ok {
			n := newNoise(eng, pv.eid)
			eng.noises[pv.eid] = n
			return n
		}
	}
	return nil
}

// There is always only one listener. It is associated with the root pov
// by default. This changes the listener location to the given pov.
func (eng *engine) setListener(p Pov) {
	if pv, ok := p.(*pov); ok && pv != nil {
		eng.soundListener = pv
	}
}

// FUTURE: cleaning up resources is not complete. Dispose currently means
// removing entities from the Pov hierarchy and from the eng entity manager,
// yet keeps them in the cache and bound on the GPU/Snd devices. Applications
// often "dispose" parts of the Pov hierarchy only to (re)use the underlying
// data again. Likely need another method/parm for the application to indicate
// data that is to be completely "unloaded".
//
// Note: cached objects know about "loaded" and must be kept in sync
// if device bound data is changed.

// dispose discards the given pov component or the entire pov and all
// its components. Each call recalculates the currently loaded set
// of assets.
func (eng *engine) dispose(p Pov, component int) {
	if pv, ok := p.(*pov); ok && pv != nil {
		switch component {
		case BODY:
			delete(eng.bodies, pv.eid)
			delete(eng.solids, pv.eid)
		case CAMERA:
			delete(eng.cams, pv.eid)
		case MODEL:
			if m, ok := eng.models[pv.eid]; ok {
				eng.disposeModel(m)
				delete(eng.models, pv.eid)
			}
		case NOISE:
			if n, ok := eng.noises[pv.eid]; ok {
				eng.disposeNoise(n)
				delete(eng.noises, pv.eid)
			}
		case LIGHT:
			delete(eng.lights, pv.eid)
		case LAYER:
			if l, ok := eng.layers[pv.eid]; ok {
				eng.disposeLayer(l)
				delete(eng.layers, pv.eid)
			}
		case POV:
			eng.disposePov(pv)
		}
	}
}

// disposePov chops this transform and all of its children out
// of the transform hierarchy. All associated objects are disposed.
func (eng *engine) disposePov(pv *pov) {
	delete(eng.povs, pv.eid)
	eng.dispose(pv, CAMERA)
	eng.dispose(pv, BODY)
	eng.dispose(pv, MODEL)
	eng.dispose(pv, NOISE)
	if pv.parent != nil {
		pv.parent.remChild(pv) // remove the one back reference that matters.
	}
	for _, child := range pv.children {
		child.parent = nil    // avoid unnecessary removing of back references
		eng.disposePov(child) // ... since childs parent is being deleted.
	}
}

// disposeModel releases any references to assets.
func (eng *engine) disposeModel(m *model) {
	m.msh = nil
	m.shd = nil
	m.anm = nil
	m.fnt = nil
	m.mat = nil
	m.texs = []*texture{} // garbage collect all old textures.
}

// disposeNoise releases references to assets.
func (eng *engine) disposeNoise(n *noise) {
	n.snds = []*sound{} // garbage collect the old sounds.
}

// disposeLayer removes the render pass layer, if any, from the given entity.
// No complaints if there is no layer at the given entity. This is safe to
// remove from the GPU since it is not cached.
func (eng *engine) disposeLayer(l *layer) {
	eng.release(&releaseData{data: l}) // dispose of the framebuffer.
}

// Usage returns numbers collected each time through the
// main processing loop. This allows the application to get
// a sense of time usage.
func (eng *engine) Usage() *Timing { return eng.times }

// Modelled returns the total number of models and the total
// number of verticies for all models.
func (eng *engine) Modelled() (models, verts int) {
	models = len(eng.models)
	for _, m := range eng.models {
		if m.msh != nil && len(m.msh.vdata) > 0 {
			verts += m.msh.vdata[0].Len()
		}
	}
	return models, verts
}

// Rendered returns the number of models and the number
// of verticies rendered in the last rendering pass.
func (eng *engine) Rendered() (models, verts int) {
	return eng.sm.renDraws, eng.sm.renVerts
}

// Eng interface implementation to handle requests
// for changing engine state.
func (eng *engine) SetColor(r, g, b, a float32) {
	go func(r, g, b, a float32) {
		eng.machine <- &setColour{r: r, g: g, b: b, a: a}
	}(r, g, b, a)
}
func (eng *engine) ShowCursor(show bool) {
	go func(show bool) { eng.machine <- &showCursor{enable: show} }(show)
}
func (eng *engine) SetCursorAt(x, y int) {
	go func(x, y int) { eng.machine <- &setCursor{cx: x, cy: y} }(x, y)
}
func (eng *engine) Enable(attr uint32, enabled bool) {
	go func(attr uint32, enabled bool) {
		eng.machine <- &enableAttr{attr: attr, enable: enabled}
	}(attr, enabled)
}
func (eng *engine) ToggleFullScreen() {
	go func() { eng.machine <- &toggleScreen{} }()
}
func (eng *engine) Mute(mute bool) {
	gain := 1.0
	if mute {
		gain = 0.0
	}
	go func(gain float64) { eng.machine <- &setVolume{gain: gain} }(gain)
}
func (eng *engine) SetVolume(zeroToOne float64) {
	go func(gain float64) { eng.machine <- &setVolume{gain: zeroToOne} }(zeroToOne)
}

// engine
// ===========================================================================
// Expose/wrap physics shapes.

func (eng *engine) SetGravity(g float64) { eng.physics.SetGravity(g) }

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

// Collide checks if two bodies are intersecting independent of the solver
// and without updating the the bodies locations.
func (eng *engine) Collide(a, b physics.Body) bool {
	return eng.physics.Collide(a, b)
}
