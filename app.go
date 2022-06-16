// Copyright © 2017-2024 Galvanized Logic Inc.

package vu

// app.go interfaces with the game or application. It is used to call the
// application from the engine and to hold all application created resources.

import (
	"time"

	"github.com/gazed/vu/render"
)

// application loads and tracks resources created by the application.
// Examples of application created resources are:
//
//	scene component for rendering frames.
//	physics component for handling forces and collisions.
//	transform component for updating scene graph transforms.
//
// One application instance is created by engine on startup.
type application struct {
	updator Updator // Application update callback
	resizer Resizer // Application resize callback
	input   *Input  // User input is refreshed each update.

	// Application resources are grouped by the type of data.
	eids   *entities   // Entity id manager.
	sounds *sounds     // Audio components.
	scenes *scenes     // Scene component, one camera per scene
	povs   *povs       // Transform components.
	models *models     // Render components.
	lights *lights     // Light components.
	sim    *simulation // Physic simulation components.

	// Load assets from files in a separate go-routine.
	ld *assetLoader // looks in local "assets" directory by default.

	// frame holds the scene render packet information.
	frame []render.Pass // reused each render.
}

// Initialize the application data.
func newApplication() (app *application) {
	app = &application{
		input: &Input{
			Pressed:  map[int32]bool{},
			Down:     map[int32]time.Time{},
			Released: map[int32]time.Duration{},
		},

		// initialize the component managers.
		eids:   &entities{},     // entity id manager.
		sounds: newSounds(),     // audio resources.
		scenes: newScenes(),     // scenes to group models.
		povs:   newPovs(),       // model transforms.
		models: newModels(),     // 2D and 3D models.
		lights: newLights(),     // 3D lights.
		sim:    newSimulation(), // physics simulation
	}
	app.ld = newLoader() // start the loader goroutine.
	app.frame = []render.Pass{
		render.NewPass(), // 3D
		render.NewPass(), // 2D
	}
	return app
}

// addScene is used internally for testing.
// It is exposed publicly as eng.AddScene.
func (app *application) addScene(st SceneType) *Entity {
	eid := app.eids.create()
	app.povs.create(eid, 0)    // scene graph root.
	app.scenes.create(eid, st) // 3D or 2D
	return &Entity{app: app, eid: eid}
}

// dispose asks each of the component managers to completely remove any
// knowledge of the given entity. The entity id is recycled.
func (app *application) dispose(eng *Engine, eid eID) {

	// collect the identities that need disposing.
	dead := []eID{}
	if s := app.scenes.get(eid); s != nil {
		dead = app.scenes.dispose(eid, dead)
	}
	dead = app.povs.dispose(eid, dead)
	app.sim.dispose(eid)
	app.lights.dispose(eid)
	app.models.dispose(eid)
	app.sounds.dispose(eng, eid)
	app.eids.dispose(eid)
	for _, eid := range dead {
		app.dispose(eng, eid)
	}
}
