// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package vu

import (
	"vu/device"
	"vu/move"
	"vu/render"
)

// stage is the stage manager. A stage manager is responsible for interacting
// with the Director to ensure eveything requested by the director is properly
// staged and rendered. Stage is responsible for all aspects of the scene graph
// including:
//     Creation and updating by the application.
//     Updating locations using physics.
//     Culling and calculating the visible parts for rendering.
//
// Stage runs as a goroutine that communicates with, and is regulated by, the
// main engine. Stage orchestrates the application callbacks to the Director
// which result in an updated scene graph.
type stage struct {
	app     Director      // Application callbacks.
	mov     move.Mover    // Physics controls/applies forces and handles collisions.
	scenes  []*scene      // Scene graph.
	overlay *scene        // Overlay is the last scene drawn.
	Mx, My  int           // Mouse position updated each loop.
	Px, Py  int           // Previous mouse position updated each loop.
	in      *Input        // Used to propogate device user input.
	solids  []*part       // Parts with Bodies going to be updated by physics.
	bodies  []move.Body   // Bodies being updated by physics.
	vis     []*render.Vis // Visible list being updated.
}

func newStage(director Director) *stage {
	sm := &stage{}
	sm.app = director
	sm.mov = move.NewMover()
	sm.in = &Input{}
	sm.solids = []*part{}
	sm.bodies = []move.Body{}
	return sm
}

// addScene creates a new scene with its own camera and lighting.
func (sm *stage) addScene(transform int) Scene {
	if sm.scenes == nil {
		sm.scenes = []*scene{}
	}
	sc := newScene(transform, sm)
	sm.scenes = append(sm.scenes, sc)
	return sc
}

// remScene disposes given scene and everything within it. While generally scenes are
// created once and last until the application closes, applications may need to discard
// scenes that are no longer needed in order to manage resource consumption.
func (sm *stage) remScene(s Scene) {
	if sc, _ := s.(*scene); sc != nil {
		for index, existing := range sm.scenes {
			if sc == existing {
				existing.dispose()
				sm.scenes = append(sm.scenes[:index], sm.scenes[index+1:]...)
				return
			}
		}
	}
}

// setOverlay marks a scene to be the one drawn after all the other
// screens.  This is expected to be a heads-up-display 2D scene.
func (sm *stage) setOverlay(s Scene) {
	if sc, _ := s.(*scene); sc != nil {
		sm.overlay = sc
	}
}

// runUpdate runs one update producing a new sm.visible at the end.
func (sm *stage) runUpdate(pressed *device.Pressed, dt float64) {
	sm.Px, sm.Py = sm.Mx, sm.My
	sm.Mx, sm.My = pressed.Mx, pressed.My

	// customize the device input for the application.
	sm.in.Mx, sm.in.My = sm.Mx, sm.My
	sm.in.Shift = pressed.Shift
	sm.in.Control = pressed.Control
	sm.in.Down = pressed.Down // yes, this is a reference.
	sm.in.Focus = pressed.Focus
	sm.in.Resized = pressed.Resized
	sm.in.Dt = dt
	sm.in.Gt += 1
	sm.app.Update(sm.in) // application state changes based on input.
	sm.moveScene(dt)     // physics updates.
	sm.stageScene()      // prepare scene graph for render.
}

// moveScene updates all the moving parts that are controlled by the
// physics system. Physics is run and then the new positions and directions
// are used by the part.
func (sm *stage) moveScene(dt float64) {
	sm.mov.Step(sm.bodies, dt)
	for _, b := range sm.bodies {
		if b.IsMovable() {
			p := b.Data().(*part)
			p.loc.Set(b.World().Loc)
			p.dir.Set(b.World().Rot)
		}
	}
}

// stageScene culls and sorts the scenes to calculate the list of visible
// parts that need to be rendered.
func (sm *stage) stageScene() {

	// Keep the scene parts sorted by distance for transparency to work.
	// The further away objects have to be drawn first.
	sm.vis = sm.vis[:0] // reset, keep the underlying memory.
	for _, sc := range sm.scenes {
		if sc.radius > 0 || sc.sorted {
			sc.setDistances(sc.parts)
			if sc.radius > 0 {
				sc.cullParts(sc.parts)
			}
			if sc.sorted {
				sc.sortParts(sc.parts)
			}
		}
	}

	// Calculate the parts that need rendering, drawing the overlay scene last.
	for _, sc := range sm.scenes {
		if sc.Visible() {
			if sc == sm.overlay {
				continue
			}
			sc.stage(&sm.vis)
		}
	}
	if sm.overlay != nil && sm.overlay.Visible() {
		sm.overlay.stage(&sm.vis)
	}
}

// Add implements move.World.Add
func (sm *stage) Add(body move.Body) {
	sm.bodies = append(sm.bodies, body)
}

// Rem implements move.World.Rem
func (sm *stage) Rem(body move.Body) {
	for index, b := range sm.bodies {
		if b.Id() == body.Id() {
			sm.bodies = append(sm.bodies[:index], sm.bodies[index+1:]...)
			return
		}
	}
}
