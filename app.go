// Copyright © 2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// app.go holds all the component managers and runs application state updates.
//        This is the main user facing class.
// DESIGN: keep small by delegating application requests to the components.

import (
	"time"

	"github.com/gazed/vu/physics"
)

// App methods are called by the engine. It is implemented
// by the application and registered once on engine creation as follows:
//	   err := vu.Run(app) // Register app and run engine.
// The application communicates with the engine by calling the Eng
// methods as well as creating and controlling entities.
type App interface {
	Create(eng Eng, s *State) // Called once after successful startup.

	// Update allows applications to change state prior to the next render.
	// Update is called many times a second after the initial call to Create.
	//    i : user input refreshed prior to each call.
	//    s : engine state refreshed prior to each call.
	Update(eng Eng, i *Input, s *State) // Process user input, update state.
}

// app controls the application by implementing Eng and calling the
// application App methods. It is also an entity manager using a unique
// entity id to group similar application data together. Each application
// data type is a component with a component manager. The component managers
// processes the data instances for its type. For example:
//     A physics manager for handling forces and collisions.
//     A scene manager for rendering frames.
//     A pov manager for updating scene graph transforms.
// One app instance is created by engine. It is passed back and forth
// between the application update goroutine and the engine device-aware
// main thread.
type application struct {
	app   App       // Controlling application.
	ld    *loader   // File based asset loading.
	ut    uint64    // Total number of update ticks.
	attrs []EngAttr // Attributes needing setting in main thread.
	stop  bool      // Set by application to terminate engine.

	// Reused and refreshed each update.
	input *Input   // Refreshed each update.
	state *State   // Refreshed each update.
	prof  *Profile // Track render and update times.
	frame frame    // Reusable render frame.

	// Application entities are grouped into components,
	// where each component has a corresponding manager.
	eids   *eids   // Ent id manager.
	scenes *scenes // Scene component.
	povs   *povs   // Transform components.
	models *models // Render components.
	bodies *bodies // Physic components.
	lights *lights // Light components.
	sounds *sounds // Audio components.
}

// newApplication is called once on startup by engine.
func newApplication(callback App) *application {
	app := &application{app: callback}
	app.attrs = []EngAttr{}
	app.input = &Input{Down: map[int]int{}, Dt: timeStepSecs}
	app.state = &State{CullBacks: true, Blend: true}
	app.prof = &Profile{}

	// initialize the component managers.
	app.eids = newEids() // entity id manager.
	app.eids.create()    // mark eid 0 as special.
	app.scenes = newScenes()
	app.povs = newPovs()
	app.ld = newLoader()
	app.lights = newLights()
	app.sounds = newSounds()
	app.bodies = newBodies()
	app.models = newModels()
	app.frame = frame{}
	return app
}

// update application state as well as any object positioning such as cameras,
// physics, particles, animated models, etc. Run as a goroutine with exclusive
// control of the application instance which is returned once updates are done.
//    start time when overall update started.
func (app *application) update(ut uint64,
	start time.Time, elapsed time.Duration,
	done chan *application) {

	// get user input since last update.
	if app.input.Resized {
		app.scenes.resize(app.state.W, app.state.H)
	}
	app.ld.processImports()

	// Update physics and particles using a fixed timestep so that
	// each update advances by the same amount.
	app.bodies.stepVelocities(timeStepSecs)
	app.models.moveParticles(timeStepSecs)

	// Advance model animations by elapsed time, not at fixed rate like physics.
	// Animation data expects to be played back at a particular frame rate.
	app.models.animate(elapsed.Seconds())

	// The application updates its own state as well as creating and
	// deleting game objects, or even shutting down the engine.
	app.app.Update(app, app.input, app.state)
	if !app.stop {

		// Remember where everything was last update so that the display
		// can interpolate between current and previous states.
		app.povs.setPrev()
		app.scenes.setPrev()

		// Reset profile data and start counting times for the next update.
		// This updates time is returned next update.
		app.prof.Zero()
		app.prof.Update += time.Since(start)
	}
	done <- app // return application to the engine thread.
}

// State provides access to current engine state.
func (app *application) State() *State { return app.state }

// Implement Eng interface. Returns the physics instance.
func (app *application) Physics() physics.Physics { return app.bodies.physics }

// Implment Eng interface.
func (app *application) AddScene() *Ent {
	eid := app.eids.create()
	app.povs.create(eid, 0) // pov component for the scene graph root.
	app.scenes.create(eid)
	return &Ent{app: app, eid: eid}
}

// AddSound loads audio data and returns a unique sound identifier.
// Passing the sound indentifier to any Ent.PlaySound() method will
// play the sound at the entities location. Also use an entity location
// to set the unique sound listener location. Played sounds are louder
// the closer the played sound to the sound listener.
//    name : identifies a sound asseet eg: sounds/name.wav
func (app *application) AddSound(name string) uint32 {
	eid := app.eids.create()
	app.sounds.create(app.ld, eid, name)
	return uint32(eid)
}

// Shutdown is an application request to close down the engine.
// Expected to be called once on Application exit.
func (app *application) Shutdown() { app.shutdown() }

// shutdown is called on the update goroutine by the application.
// A shutdown by the device layer just kills everything and stops
// the device layer callbacks to the engine.
func (app *application) shutdown() {
	app.stop = true    // signal engine to shutdown.
	if app.ld != nil { // close the loader.
		app.ld.dispose()
		app.ld = nil
	}
	for _, s := range app.scenes.all {
		app.dispose(s.eid) // Delete all scenes deletes everything.
	}
}

// dispose asks each of the component managers to completely remove any
// knowledge of the given entity. The entity id is recycled.
//
// FUTURE: cleaning up resources is not complete. Dispose currently means
// removing entities from the app entity manager, yet keeps assets in the cache
// and bound on the GPU/Snd devices. Applications often "dispose" entities only
// to (re)use the underlying data again. Likely need another API for the
// application to indicate when assets are to be completely unloaded.
func (app *application) dispose(id eid) {
	dead := []eid{} // need new one each time for recursion.
	if s := app.scenes.get(id); s != nil {
		dead = app.scenes.dispose(id, dead)
	}
	dead = app.povs.dispose(id, dead)
	app.models.dispose(id)
	app.bodies.dispose(id)
	app.lights.dispose(id)
	app.sounds.dispose(id)
	app.eids.dispose(id)
	for _, id := range dead {
		app.dispose(id)
	}
}

// Times returns numbers collected each main update loop.
// This allows the application to get a sense of time and resource usage.
func (app *application) Times() *Profile { return app.prof }

// Set one or more engine attributes.
func (app *application) Set(attrs ...EngAttr) {
	// Save application attribute update requests. Set is called on
	// the update goroutine and the attributes requests need to be .
	// run on the main thread.
	app.attrs = append(app.attrs, attrs...)
}

// setAttributes is called on the main thread to clear any
// outstanding attribute requests. Attributes affect the device layer
// and devices like to be called on the main thread.
func (app *application) setAttributes(eng *engine) {
	for _, attr := range app.attrs {
		attr(eng)
	}

	// Eventually the old attribute functions will be overwritten
	// and garbage collected.
	app.attrs = app.attrs[:0] // reset preserving memory.
}
