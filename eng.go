// Copyright © 2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"log"
	"math"
	"time"

	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/move"
	"github.com/gazed/vu/render"
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
	Root() Pov     // Single root transform always exists.

	// Timing is updated each processing loop. The returned update
	// times can flucuate and should be averaged over multiple calls.
	Usage() *Timing                // Per update loop performance metrics.
	Modelled() (models, verts int) // Total render models and verticies.
	Rendered() (models, verts int) // Rendered models and verticies.

	// Requests to change engine state.
	SetColor(r, g, b, a float32)      // Set background clear colour.
	ShowCursor(show bool)             // Hide or show the cursor.
	SetCursorAt(x, y int)             // Put cursor at the window pixel x,y.
	Enable(attr uint32, enabled bool) // Enable/disable render attributes.
	ToggleFullScreen()                // Flips full screen and windowed mode.
	Mute(mute bool)                   // Toggle sound volume.
	SetVolume(zeroToOne float64)      // Set sound volume.
	SetGravity(g float64)             // Change the gravity constant.
}

// App is the engine callback expected to be implemented by the application.
// The application is registered on engine creation as follows:
//     err := vu.New(app, "Title", 0, 0, 800, 600)
// Note that it is safe to call Eng methods from goroutines.
type App interface {
	Create(eng Eng, s *State) // One time call after successfull startup.

	// Update allows applications to change state prior to the next render.
	// Update is called many times a second after the initial call to Create.
	//      i : user input refreshed prior to each call.
	//      s : engine state refreshed prior to each call.
	Update(eng Eng, i *Input, s *State) // Process user input.
}

// Eng and App interfaces.
// ===========================================================================
// engine implements Eng.

// engine runs all application communication and all state updates.
// It is a entity manager in that it uses a unique entity id to group
// by component functionality.
//
// Expected to be started as a go-routine using the runEngine method.
type engine struct {
	alive   bool      // True until application decides otherwise.
	machine chan msg  // Communicate with device loop.
	stop    chan bool // Closed or any value means stop the engine.
	data    *appData  // Combination user input and application state.

	// Assets are loaded concurrently.
	loader *loader         // Asset manager.
	loaded chan []*loadReq // Receive models and noises ready to render/play.
	mover  move.Mover      // Physics handles forces, collisions.

	// Use three render frames. One for updating state,
	// the other two are for rendering with interpolation.
	frames [][]render.Draw // 3 render frames.
	vorder []uint64        // Eids for view render order.

	// Sounds are heard by the sound listener at an app set pov.
	soundListener *pov    // Current location of the sound listener.
	sx, sy, sz    float64 // Last location of the sound listener.

	// Lighting currently handles a single light.
	l          *light  // Default light.
	lx, ly, lz float64 // Default light position.

	// Track update times, the number of draw calls, and verticies.
	times    *Timing // Loop timing statistics.
	renDraws int     // Last updates models rendered.
	renVerts int     // Last updates verticies rendered.

	// Group the application entities by component.
	// All entities are Pov (location:orientation) based.
	eid    uint64               // Next entity id.
	povs   map[uint64]*pov      // Entity transforms.
	views  map[uint64]*view     // Cameras components.
	models map[uint64]*model    // Visible components.
	lights map[uint64]*light    // Light components.
	noises map[uint64]*noise    // Audible components.
	bodies map[uint64]move.Body // Non-colliding physic components.
	solids map[uint64]move.Body // Colliding physic components.
	bods   []move.Body          // Set from solids each update.
	mv     *lin.M4              // Scratch model-view matrix.
	mvp    *lin.M4              // Scratch model-view-proj matrix.
}

// newEngine is expected to be called once on startup.
func newEngine(machine chan msg) *engine {
	eng := &engine{alive: true, machine: machine}
	eng.data = newAppData()
	eng.mv = &lin.M4{}
	eng.mvp = &lin.M4{}
	eng.times = &Timing{}
	eng.frames = make([][]render.Draw, 3)
	for cnt, _ := range eng.frames {
		eng.frames[cnt] = []render.Draw{}
	}
	eng.Reset()
	eng.l = newLight()

	// helpers that create and update state.
	eng.mover = move.NewMover()
	eng.loaded = make(chan []*loadReq)
	eng.loader = newLoader(eng.loaded, machine)
	go eng.loader.runLoader()
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

// runEngine calls Create once and Update continuously and regularly.
// The application loop generates device polling and render requests
// for the machine.
func runEngine(app App, wx, wy, ww, wh int, machine chan msg, stop chan bool) {
	defer catchErrors()
	eng := newEngine(machine)
	eng.stop = stop
	eng.data.state.setScreen(wx, wy, ww, wh)
	app.Create(eng, eng.data.state)
	ut := uint64(0)         // kick off initial update...
	eng.update(app, dt, ut) // first update queues the load asset requests.

	// Initialize timers and kick off the main control loop.
	var loopStart time.Time = time.Now()
	var updateStart time.Time
	var timeUsed time.Duration
	var updateTimer time.Duration // Track when to trigger an update.
	var renderTimer time.Duration // Track when to trigger a render.
	var frame []render.Draw       // New render frame, nil if no updated frame.
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

			// Perform the update.
			eng.update(app, dt, ut)     // Update state, physics, etc.
			frame = eng.frames[ut%3]    // Cycle between three render frames.
			frame = eng.genFrame(frame) // ... update the render frame.
			eng.frames[ut%3] = frame    // ... remember the updated frame.

			// Reset and start counting times for the next update.
			eng.times.Zero()
			eng.times.Update += time.Since(updateStart)
		}

		// Interpolation is the fraction of unused delta time between 0 and 1.
		// ie: State state = currentState*interpolation + previousState * (1.0 - interpolation);
		interpolation := updateTimer.Seconds() / dt.Seconds()

		// Redraw everything, using interpolation when there is no new frame.
		// Extra render time is dropped.
		renderTimer += timeUsed
		if renderTimer >= rt {
			eng.times.Renders += 1
			eng.render(frame, interpolation, ut)
			frame = nil                    // mark frame as rendered.
			renderTimer = renderTimer % rt // drop extra render time.
		}
		eng.communicate() // process go-routine messages.
	}
	// Exiting state update.
}

// communicate processes all go-routine channels. Must be non-blocking.
// Incoming messages  are generally responses to asset loading completions
// that were initiated by this engine.
func (eng *engine) communicate() {
	select {
	case <-eng.stop: // closed channels return 0
		eng.loader.control <- &shutdown{}
		return // Exit immediately, The main loop has closed us down.
	case loaded := <-eng.loaded:
		for _, req := range loaded {
			if req.err != nil {
				log.Printf("load error: %s", req.err)
				continue
			}
			switch a := req.a.(type) {
			case *mesh:
				req.model.msh = a
			case *texture:
				if req.index < len(req.model.texs) {
					req.model.texs[req.index] = a
				}
			case *shader:
				req.model.shd = a
			case *font:
				m := req.model
				m.fnt = a
				m.fnt.loaded = true
				if len(m.phrase) > 0 {
					m.phraseWidth = m.fnt.setPhrase(m.msh, m.phrase)
				}
			case *animation:
				m := req.model
				m.anm = a
				m.msh = req.msh
				if req.index < len(m.texs) && len(req.texs) > 0 {
					m.texs[req.index] = req.texs[0]
				}
				m.nFrames = a.maxFrames(0)
				m.pose = make([]lin.M4, len(a.joints))
			case *material:
				m := req.model
				m.mat = a
				if m.alpha == 1.0 {
					m.alpha = a.tr // Copy values so they can be set per model.
				}
				if m.kd.isUnset() {
					m.kd = a.kd // Copy values so they can be set per model.
				}
				m.ks = a.ks // Can't currently be overridden on model.
				m.ka = a.ka // ditto
			case *sound:
				n := req.noise
				n.snds[req.index] = a
				n.loaded = true
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
	eng.mover.Step(eng.bods, dts)

	// Have the application adjust any or all state before rendering.
	input.Dt = dts                // how long to get back to here.
	input.Ut = ut                 // update ticks.
	app.Update(eng, input, state) // application to updates its own state.

	// update assets that the application changed or which need
	// per tick processing. Per-ticks include animated models,
	// particle effects, surfaces, phrases, ...
	if eng.alive {
		eng.updateModels(dts)          // load and bind data.
		eng.place(eng.root(), lin.M4I) // update all transforms.
		eng.updateSoundListener()      // reposition sound listener.
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

// genFrame creates a new render frame. All of the transforms are
// expected to have been placed (updated) before calling this method.
// A frame is rendered using each of the views in creation order.
//
// The frame memory is preserved through the method calls in
// that the Draw records are grown and reused across updates.
func (eng *engine) genFrame(frame []render.Draw) []render.Draw {
	frame = frame[:0] // resize keeping underlying memory.
	eng.renDraws, eng.renVerts = 0, 0
	light := eng.l
	for _, eid := range eng.vorder {
		if view, ok := eng.views[eid]; ok && view.visible {
			if pv, ok := eng.povs[eid]; ok { // expected to be true.
				light, frame = eng.snapshot(pv, view, light, frame)
			}
		}
	}
	render.SortDraws(frame)
	return frame // Note len(frame) is the number of draw calls.
}

// snapshot recursively walks the transform hierarchy for the given view
// and adds render objects for each model it finds.
func (eng *engine) snapshot(p *pov, view *view, light *light, frame []render.Draw) (*light, []render.Draw) {
	processChildren := true
	if p.visible {

		// support one light at a time. It can be repositioned per object.
		if l, ok := eng.lights[p.eid]; ok {
			light = l
			lx, ly, lz := p.Location()
			vec := view.cam.v0.SetS(lx, ly, lz, 1)
			vec.MultvM(vec, view.cam.vm)
			eng.lx, eng.ly, eng.lz = vec.X, vec.Y, vec.Z
		}
		if model, ok := eng.models[p.eid]; ok { // model at this transform.
			if model.loaded() {
				eng.mv.Mult(p.mm, view.cam.vm)    // model-view
				eng.mvp.Mult(eng.mv, view.cam.pm) // model-view-projection

				// Calculate distance to camera for 3D models.
				// Needed for transparency sorting and culling.
				distToCam := 0.0
				var px, py, pz float64
				if view.depth {
					vec := view.cam.v0.SetS(0, 0, 0, 1)
					vec.MultvM(vec, p.mm)
					px, py, pz = vec.X, vec.Y, vec.Z
					distToCam = view.cam.Distance(px, py, pz)
				} else {
					px, py, pz = p.Location() // for UI culling.
				}

				// Create render data for anything that passes culling.
				if view.cull != nil && view.cull.Cull(view.cam, px, py, pz) {
					processChildren = false // don't process children of culled models.
				} else {
					if model.msh != nil && len(model.msh.vdata) > 0 {

						// fill in render and uniform data.
						var draw *render.Draw
						if frame, draw = eng.getDraw(frame); draw != nil {
							(*draw).SetMv(eng.mv)      // model-view
							(*draw).SetMvp(eng.mvp)    // model-view-projection
							(*draw).SetPm(view.cam.pm) // projection only.
							(*draw).SetScale(p.Scale())
							model.toDraw(*draw, p.eid, view.depth, view.overlay, distToCam)
							light.toDraw(*draw, eng.lx, eng.ly, eng.lz)
							eng.renDraws += 1                        // models rendered.
							eng.renVerts += model.msh.vdata[0].Len() // verticies rendered.
						}
					} else {
						log.Printf("Model has no mesh data...")
					}
				}
			} else {
				processChildren = false // only process children of loaded models.
			}
		}
		if processChildren {
			for _, child := range p.children {
				light, frame = eng.snapshot(child, view, light, frame)
			}
		}
	}
	return light, frame
}

// getDraw returns a render.Draw. The frame is grown as needed and
// forms are reused if available. Every frame value up to cap(frame)
// is expected to have already been allocated.
func (eng *engine) getDraw(frame []render.Draw) (f []render.Draw, d *render.Draw) {
	size := len(frame)
	switch {
	case size == cap(frame):
		frame = append(frame, render.NewDraw())
	case size < cap(frame): // use previously allocated.
		frame = frame[:size+1]
		if frame[size] == nil {
			frame[size] = render.NewDraw()
		}
	}
	return frame, &frame[size]
}

// place walks the transform hierarchy updating all the model view transforms.
// This is called before rendering snapshots are taken.
func (eng *engine) place(p *pov, parent *lin.M4) {
	p.mm.SetQ(p.rot.Inv(p.at.Rot)) // invert model rotation.
	p.mm.ScaleSM(p.Scale())        // scale is applied first (on left of rotation)
	l := p.at.Loc
	p.mm.TranslateMT(l.X, l.Y, l.Z) // translate is applied last (on right of rotation).
	p.mm.Mult(p.mm, parent)         // model transform + parent transform
	for _, child := range p.children {
		eng.place(child, p.mm)
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

// render sends a render request to the machine.
func (eng *engine) render(frame []render.Draw, interp float64, ut uint64) {
	eng.machine <- &renderFrame{interp: interp, frame: frame, ut: ut}
}

// release sends a release resource request to the machine.
// Expected to be run as a goroutine so that it can block on the
// send until the machine is ready to process it.
func (eng *engine) release(rd *releaseData) { eng.machine <- rd }

// rebind sends a bind request to the machine.
// Expected to be called as a goroutine so that it can block on the
// send until the machine is ready to process it. It also blocks
// on the reply until the rebind is finished.
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
func (eng *engine) Shutdown() {
	eng.dispose(eng.root(), POV)
	eng.loader.control <- &shutdown{}
	eng.machine <- &shutdown{}
	eng.alive = false
}

// Reset removes all entities and sets the engine back to
// its initial state.
func (eng *engine) Reset() {
	eng.dispose(eng.root(), POV)
	eng.povs = map[uint64]*pov{}
	eng.views = map[uint64]*view{}
	eng.models = map[uint64]*model{}
	eng.lights = map[uint64]*light{}
	eng.noises = map[uint64]*noise{}
	eng.bodies = map[uint64]move.Body{}
	eng.solids = map[uint64]move.Body{}
	eng.eid = 1
	eng.povs[eng.eid] = newPov(eng, eng.eid) // root
	eng.soundListener = eng.povs[eng.eid]
}

// State provides access to current engine state.
func (eng *engine) State() *State { return eng.data.state }

// genid returns the next unique entity id. It craps out and starts returning
// 0 after generating all possible ids.
func (eng *engine) genid() uint64 {
	if eng.eid == math.MaxUint64 {
		return 0
	}
	eng.eid += 1 // first valid id is 1.
	return eng.eid
}

// Implement Eng interface. Group and track all the entity
// types that can be associated with a Pov.
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

// view entities.
func (eng *engine) view(p Pov) View {
	if pv, ok := p.(*pov); ok && pv != nil {
		if view, ok := eng.views[pv.eid]; ok {
			return view
		}
	}
	return nil
}
func (eng *engine) newView(p Pov) View {
	if pv, ok := p.(*pov); ok && pv != nil {
		v := newView()
		eng.views[pv.eid] = v
		eng.vorder = append(eng.vorder, pv.eid)
		return v
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

// body: physics entities.
func (eng *engine) body(p Pov) move.Body {
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
func (eng *engine) newBody(p Pov, b move.Body) move.Body {
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

// dispose discards the given pov aspect or the entire pov and all
// its components. Each call recalculates the currently loaded set
// of assets.
func (eng *engine) dispose(p Pov, aspect int) {
	if pv, ok := p.(*pov); ok && pv != nil {
		switch aspect {
		case BODY:
			delete(eng.bodies, pv.eid)
			delete(eng.solids, pv.eid)
		case VIEW:
			eng.disposeView(pv.eid)
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
		case POV:
			eng.disposePov(pv)
		}
	}
}

// disposePov chops this transform and all of its children out
// of the transform hierarchy. All associated objects are disposed.
func (eng *engine) disposePov(pv *pov) {
	delete(eng.povs, pv.eid)
	eng.dispose(pv, VIEW)
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

// disposeView removes the view, if any, from the given entity.
// No complaints if there is no view at the given entity.
func (eng *engine) disposeView(eid uint64) {
	if _, ok := eng.views[eid]; ok {
		delete(eng.views, eid)
		for index, id := range eng.vorder {
			if eid == id {
				eng.vorder = append(eng.vorder[:index], eng.vorder[index+1:]...)
				return
			}
		}
	}
}

// FUTURE: clean assets... use countAssets to match against cached assets
//      and release assets that are no longer used. This will be a type
//      of garbage collecting triggered by the application.
//      ie: go eng.release(rd)        // for bound data.
//          go eng.loader.request(rd) // to release from cache.

// countAssets walks the current models and noises and returns
// a count of unique resources.
func (eng *engine) countAssets() map[uint64]int {
	assets := map[uint64]int{}
	for _, m := range eng.models {
		if m.msh != nil {
			assets[m.msh.bid()] += 1
		}
		if m.shd != nil {
			assets[m.shd.bid()] += 1
		}
		if m.mat != nil {
			assets[m.mat.bid()] += 1
		}
		if m.fnt != nil {
			assets[m.fnt.bid()] += 1
		}
		if m.anm != nil {
			assets[m.anm.bid()] += 1
		}
		for _, t := range m.texs {
			assets[t.bid()] += 1
		}
	}
	for _, n := range eng.noises {
		for _, s := range n.snds {
			assets[s.bid()] += 1
		}
	}
	return assets
}

// Usage returns numbers collected each time through the
// main processing loop. This allows the application to get
// a sense of how much and where time is being used.
func (eng *engine) Usage() *Timing { return eng.times }

// Modelled returns the number of models and the number
// of verticies for all models.
func (eng *engine) Modelled() (models, verts int) {
	models = len(eng.models)
	for _, m := range eng.models {
		if m.msh != nil && len(m.msh.vdata) > 0 {
			verts += m.msh.vdata[0].Len()
		}
	}
	return models, verts
}

// Modelled returns the number of models and the number
// of verticies rendered in the last rendering pass.
func (eng *engine) Rendered() (models, verts int) {
	return eng.renDraws, eng.renVerts
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
func (eng *engine) SetGravity(g float64) { eng.mover.SetGravity(g) }

// engine
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
