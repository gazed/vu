// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package vu

// Stage manager is the top level of what could be considered a "scene graph"
// manager that includes multiple Scene and Part objects.
//    http://www.gamerendering.com/category/scene-management/
//    http://www.realityprime.com/blog/2007/06/scenegraphs-past-present-and-future
//
// http://www.emporia.edu/~bartruff/Theatre%20handbook/SMduties.htm (Stage Manager)
//    "The position has a unique function because it serves the dual function
//     of assistant to the director and production staff during the rehearsal
//     period and then becomes the person in charge of the production during
//     the actual performance."

import (
	"vu/device"
	"vu/move"
	"vu/panel"
	"vu/render"
)

// stage (stage manager) handles updating and rendering for all scenes.
// It organizes and correlates application provided scene graph information
// into separate structures more easily consumed by subsytems like physics
// and rendering.
//
// Each update allows the application (Director interface) and physics to
// alter render state. The stage manager uses the updated state to prepare
// the lists of rendered items, culling and sorting as needed. Rendering is
// accomplished by passing the list of visible items to the rendering system.
//
// The stage manager ensures the proper interworking between the following
// subsystems:
//    device  : stage forwards user input to the application (Director).
//    physics : stage tracks the parts/bodies participating in physics.
//    render  : stage tracks the currently rendered parts and effects.
//    panel   : stage renders the overlay control panel last.
//
// The primary responsibility of the stage manager is to work with Scenes
// in culling objects that do not need to be rendered this frame.
type stage struct {
	app    Director       // Application callbacks.
	scenes []*scene       // Transform hierarchy of Scene/Parts from Application.
	last   *scene         // Optional last scene drawn.
	panel  *overlay       // Control panel overlay set by Application.
	Mx, My int            // Mouse position updated each loop.
	Px, Py int            // Previous mouse position updated each loop.
	in     *Input         // Used to propogate device user input.
	mov    move.Mover     // Physics controls/applies forces, handles collisions.
	bodies []move.Body    // Bodies being updated by physics.
	staged []render.Model // Tracks currently rendered models.
	assets *assets        // Scene asset handler injected on creation.
}

// newStage is expected to be created once on engine startup.
func newStage(director Director, props *assets) *stage {
	sm := &stage{}
	sm.app = director
	sm.mov = move.NewMover()
	sm.in = &Input{}
	sm.bodies = []move.Body{}
	sm.assets = props
	return sm
}

// dispose is called on application exit. All scenes are destroyed
// and all graphics resources are released.
func (sm *stage) dispose() {
	if sm.panel != nil {
		sm.panel.dispose()
		sm.panel = nil
	}
	for _, sc := range sm.scenes {
		sc.dispose()
	}
}

// addScene creates a new scene with its own camera and lighting.
func (sm *stage) addScene(transform int) *scene {
	if sm.scenes == nil {
		sm.scenes = []*scene{}
	}
	sc := newScene(transform, sm, sm.assets)
	if sm.last == nil {
		sm.scenes = append(sm.scenes, sc)
	} else {
		n := len(sm.scenes)
		last := sm.scenes[n-1]
		sm.scenes = append(sm.scenes[:n-1], sc, last)
	}
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

// setPanel replaces the overlay control panel with the given panel.
// The input panel may be nil to remove an existing control panel.
func (sm *stage) setPanel(p panel.Panel) {
	if sm.panel != nil {
		sm.panel.dispose()
		sm.panel = nil
	}
	if p != nil {
		sm.panel = newOverlay(newScene(VP, sm, sm.assets), p)
	}
}

// setLast sets the last scene to be rendered. The scene is moved from
// its current position to the end of the list. Set the last scene to nil
// turns off having a special last scene.
func (sm *stage) setLast(s Scene) {
	if s == nil {
		sm.last = nil // turn off having a last scene.
	} else {
		sm.last = s.(*scene)
		for index, sc := range sm.scenes {
			if s == sc {
				sm.scenes = append(sm.scenes[:index], sm.scenes[index+1:]...)
				sm.scenes = append(sm.scenes, sm.last)
			}
		}
	}
}

// updateState produces an updated set of render data. Nothing in its
// call graph should be doing rendering. Note that GPU data may be created
// or updated, just not rendered.
func (sm *stage) updateState(pressed *device.Pressed, dt float64) {
	sm.Px, sm.Py = sm.Mx, sm.My           // remember previous mouse...
	sm.Mx, sm.My = pressed.Mx, pressed.My // ...and current mouse.

	// customize the device input for the application.
	sm.in.Mx, sm.in.My = sm.Mx, sm.My
	sm.in.Down = pressed.Down // yes, this is a reference passed to the App.
	sm.in.Focus = pressed.Focus
	sm.in.Resized = pressed.Resized
	sm.in.Scroll = pressed.Scroll
	sm.in.Dt = dt
	sm.in.Gt += 1
	sm.app.Update(sm.in) // application state changes based on input.
	sm.moveScenes(dt)    // physics updates.
	sm.stageScenes(dt)   // prepare scene graph for render.
}

// moveScene updates all the moving parts that are controlled by the
// physics system. Physics is run and then the new positions and directions
// are used by the part.
func (sm *stage) moveScenes(dt float64) {
	sm.mov.Step(sm.bodies, dt)
	for _, b := range sm.bodies {
		if b.IsMovable() {
			p := b.Data().(*part)
			p.loc.Set(b.World().Loc)
			p.dir.Set(b.World().Rot)
		}
	}
}

// stageScenes culls and sorts the scenes to calculate the list of visible
// parts that need to be rendered. Traverse all parts and recreate the list
// of parts that need to be rendered. Expected to be called once per update.
//
// FUTURE: possible culling and scene graph optimizations are possible here.
//         For example quad-trees may reduce the time taken to traverse large
//         terrain based scene graphs or collision detection candidates.
func (sm *stage) stageScenes(dt float64) {
	sm.staged = sm.staged[:0] // reset, keep the underlying memory.
	for _, sc := range sm.scenes {

		// Keep the scene parts sorted by distance for transparency to work.
		// The further away objects have to be drawn first.
		if sc.cull != nil || sc.sorted {
			sc.setDistances(sc.parts)
			if sc.cull != nil {
				sc.cullParts(sc.parts)
			}
			if sc.sorted {
				sc.sortParts(sc.parts)
			}
		}

		// Query each scene for the parts that need rendering.
		if sc.Visible() {
			sc.stage(dt) // fills sm.staged
		}
	}

	// Prepare the control panel overlay. These parts are always rendered last
	// so their visible parts are appended after other visible parts.
	if sm.panel != nil {
		overlayParts := sm.panel.update()
		for _, p := range overlayParts {
			sm.staged = append(sm.staged, p.role.model)
		}
	}
}

// renderVisible draws the currently visible parts using the latest
// render data. Nothing in its call graph should change the render data
// Note that GPU shader uniforms will be synced with the current CPU
// render data.
func (sm *stage) renderVisible(gc render.Renderer) {
	gc.Clear()
	for _, model := range sm.staged {
		gc.Render(model)
	}
}

// verify runs verify on each scene.
func (sm *stage) verify() error {
	for _, sc := range sm.scenes {
		if err := sc.verify(); err != nil {
			return err
		}
	}
	return nil
}

// stage implements feedback. Track per-update-cycle which bodies are part of
// physics and which models are to be rendered. Physics callbacks are made
// during the Application Update. Staging render models is done each update
// during scene graph traversal.
func (sm *stage) track(body move.Body) { sm.bodies = append(sm.bodies, body) }
func (sm *stage) stage(m render.Model) { sm.staged = append(sm.staged, m) }
func (sm *stage) release(body move.Body) {
	for index, b := range sm.bodies {
		if b.Id() == body.Id() {
			sm.bodies = append(sm.bodies[:index], sm.bodies[index+1:]...)
			return
		}
	}
}

// stage manager
// ============================================================================
// feedback

// feedback is an internal transform hierarchy monitor. Feedback is injected
// into the transform hierarchy to get feedback on render and physics updates
// as the hierarchy is traversed.
type feedback interface {
	track(b move.Body)    // Add to the physics simulation.
	release(b move.Body)  // Remove from the physics simulation.
	stage(m render.Model) // Add to current rendering cycle.
}
