// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// Stage manager is the top level of what could be considered a "scene graph"
// manager that includes multiple Scene and Part objects.
//    http://www.gamerendering.com/category/scene-management/
//    http://www.realityprime.com/blog/2007/06/scenegraphs-past-present-and-future
// It also acts as a component manager in that it groups data to be processed
// as blocks by different systems.
//
// FUTURE: consider optimizations and enhancements for Stage/Scene/Part/Role.
//         • Different scene graph layouts depending on application.
//           For example quad-trees may reduce the time taken to traverse large
//           terrain based scene graphs or collision detection candidates.
//         • Move to full entity/component/data architecture with cleaner data
//           and processing splits. Part would be reduced to just an entity id.
//           Clearly define more components like particle effects and animations.
//           Profiling should show improvements with large amount of parts
//           and ideally the code is easier to understand and maintain.
//
// Stage Manager
// http://www.emporia.edu/~bartruff/Theatre%20handbook/SMduties.htm
//    "The position has a unique function because it serves the dual function
//     of assistant to the director and production staff during the rehearsal
//     period and then becomes the person in charge of the production during
//     the actual performance."
//

import (
	"log"
	"sort"
	"vu/move"
	"vu/render"
)

// stage (stage manager) prepares for rendering by traversing the scene
// graph in Scene and Part instances. It builds different data groupings
// into separate structures more easily consumed by subsytems like physics
// and render. A single instance of stage is created on startup and all
// scene and part creation passes through it.
type stage struct {
	pid       uint32           // Next unique part identifier.
	bodies    []move.Body      // Bodies being updated by physics.
	backstage []actor          // Parts to be processed for rendering.
	staged    []render.Model   // Models ready for rendering.
	scenes    []*scene         // all scenes in rendering order.
	parts     map[uint32]*part // all parts.

	// Assets provide the data needed by shaders of renderable parts.
	assets *assets // Asset handler injected on creation.
}

// newStage is expected to be created once on engine startup.
func newStage(a *assets) *stage {
	sm := &stage{pid: 0, assets: a}
	sm.backstage = make([]actor, 100)
	sm.backstage = sm.backstage[:0]
	sm.parts = map[uint32]*part{}
	return sm
}

// dispose is called on application exit. All scenes are destroyed
// and all graphics resources are released.
func (sm *stage) dispose() {
	for _, sc := range sm.scenes {
		sc.dispose()
	}
	sm.bodies = []move.Body{}     // part bodies have been disposed.
	sm.backstage = []actor{}      // release references.
	sm.staged = []render.Model{}  // release references.
	sm.parts = map[uint32]*part{} // parts have been disposed.
	sm.scenes = []*scene{}
}

// addScene creates a new scene with its own camera and lighting.
func (sm *stage) addScene(transform int) *scene {
	sc := newScene(transform, sm)
	sm.scenes = append(sm.scenes, sc)
	return sc
}

// remScene disposes given scene and everything within it. While generally
// scenes are created once and last until the application closes, applications
// may need to discard scenes that are no longer needed in order to manage
// resource consumption.
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

// Expected to be called once per update.
func (sm *stage) update(dt float64) {
	sm.restage()              // reset all rendering data lists.
	sm.updateVisibleParts()   // find the parts that should be rendered.
	sm.updateRenderModels(dt) // prepare the part models for rendering.
}

// Walk the transform hierarchy finding parts that should be rendered
// .... produces the sm.backstage list.
func (sm *stage) updateVisibleParts() {
	for _, sc := range sm.scenes {

		// Query each scene for the parts that need rendering.
		if sc.Visible() {
			sm.cullParts(sc, sc.parts) // traverse hierarchy: sets the distance.

			// TODO only sort parts that are ready for rendering... (next loop)
			//      also only sort parts that are transparent.
			if sc.sorted {
				sm.sortParts(sc.parts) // traverse hierarchy: sort based on distance.
			}
		}
	}
}

// cullParts walks the scene part hierarchy, setting the parts distance
// to the camera and marks the part culled or not.
func (sm *stage) cullParts(sc *scene, parts Parts) {
	if parts != nil && len(parts) > 0 {
		for _, p := range parts {
			if p.visible {

				// calculate distance to camera for visible parts.
				// Distance is used for culling and sorting.
				p.toc = sc.cam.Distance(p.Location())
				dress := p.visible
				if sc.cull != nil && p.cullable {
					dress = dress && !sc.cull.Cull(sc, p)
				}
				if dress {
					sm.dress(sc, p)
					sm.cullParts(sc, p.parts) // recurse down the hierarchy.
				}
			}
		}
	}
}

// prepare the part models for rendering. Produces sm.staged list.
func (sm *stage) updateRenderModels(dt float64) {
	for _, a := range sm.backstage {
		parentTransform := a.part.pt

		// parent parts are processed first and thus have valid transforms.
		if parent := sm.part(a.part.parent); parent != nil {
			parentTransform = parent.mm
		}
		a.part.stage(a.scene, parentTransform, dt)
	}
}

// render draws the currently visible parts using the latest
// render data. Nothing in its call graph should change the render data
// Note that GPU shader uniforms will be synced with the current CPU
// render data.
func (sm *stage) render(gc render.Renderer) {
	gc.Clear()
	for _, model := range sm.staged {
		gc.Render(model)
	}
}

// addBody adds the given body to the physics simulation.
// This is called from Part as bodies are associated with a part.
func (sm *stage) addBody(body move.Body) {
	sm.bodies = append(sm.bodies, body)
}

// remBody removes the given body from the physics simulation.
// This is called from Part as bodies are removed from a part.
func (sm *stage) remBody(body move.Body) {
	for index, b := range sm.bodies {
		if b.Eq(body) {
			sm.bodies = append(sm.bodies[:index], sm.bodies[index+1:]...)
			return
		}
	}
}

// restage resets the rendering lists. The underlying memory is kept
// since the lists are used each update cycle.
func (sm *stage) restage() {
	sm.backstage = sm.backstage[:0]
	sm.staged = sm.staged[:0]
}

// dress adds to the list of parts that need to be processed for rendering.
// Parts that need to be rendered are put in a list with their corresponding
// scene informatio.
func (sm *stage) dress(s *scene, p *part) {
	last := len(sm.backstage)
	if len(sm.backstage) < cap(sm.backstage) {
		sm.backstage = sm.backstage[:last+1]
		sm.backstage[last].scene = s
		sm.backstage[last].part = p
	} else {
		sm.backstage = append(sm.backstage, actor{s, p})
	}
}

// stage adds to the list of models that are ready for rendering.
func (sm *stage) stage(m render.Model) {
	sm.staged = append(sm.staged, m)
}

// addPart is called from AddPart (Part or Scene) so that each part can be given
// an identifier and tracked through a map for quick lookup.
func (sm *stage) addPart() *part {
	sm.pid = sm.pid + 1 // first valid part id is 1.
	if sm.pid == 0 {    // go just wraps... no crash.
		log.Printf("partBroker:newPart: dev error. Unique part id wrapped.")
	}
	np := newPart(sm.pid, sm)
	sm.parts[np.pid] = np // remember all parts.
	return np
}

// remPart is called from RemPart (Part or Scene) to remove the part from
// the global part map.
func (sm *stage) remPart(p *part) {
	delete(sm.parts, p.pid)
}

// part returns a part from the global part map, return nil if no such
// part exists.
func (sm *stage) part(pid uint32) *part {
	return sm.parts[pid]
}

// setLast sets the last scene to be rendered. The scene is moved from
// its current position to the end of the list.
func (sm *stage) setLast(s Scene) {
	last := s.(*scene)
	for index, sc := range sm.scenes {
		if s == sc {
			sm.scenes = append(sm.scenes[:index], sm.scenes[index+1:]...)
			sm.scenes = append(sm.scenes, last)
			break
		}
	}
}

// sortParts reorganizes the scene graph based on the distance to the camera.
// This is used for transparent objects so that the further away objects are
// drawn first.
func (sm *stage) sortParts(parts Parts) {
	if parts != nil && len(parts) > 0 {
		for _, childPart := range parts {
			sm.sortParts(childPart.parts)
		}
		sort.Sort(parts)
	}
}

// verify checks that each part has been supplied with the render data
// needed by its shader.
func (sm *stage) verify() error {
	for _, sc := range sm.scenes {
		for _, p := range sc.parts {
			if err := p.verify(); err != nil {
				return err
			}
		}
	}
	return nil
}

// ===========================================================================

// actor is helper data that links a part with a scene. This is
// used to generate the render list data needed to turn visible parts
// into renderable data.
type actor struct {
	scene *scene
	part  *part
}
