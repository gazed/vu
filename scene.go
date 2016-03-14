// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// FUTURE: Continue to enhance support for multiple render passes.

import (
	"log"

	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/render"
)

// scene manager helps the engine create frames (render draw call lists)
// from Pov's, Models, and Cameras. It is an integral engine helper class
// that encapsulates and interacts with the render system to create a
// render frame from engine component information.
//
// There are three render frames. One for updating, the other two for
// rendering with interpolation.
type scene struct {
	scene []*pov // flattened pov hiearchy updated each frame.
	white *light // default light.

	// Multiple render pass support.
	pass         *layer  // default disabled render pass layer.
	shadowMap    *layer  // scratch shadow map.
	shadowShader *shader // shadow map specific shader.

	// Track update times, the number of draw calls, and verticies.
	renDraws int // Number of models rendered last update.
	renVerts int // Number of verticies rendered last update.

	// Scratch variables: reused to reduce garbage collection.
	mv  *lin.M4 // Scratch model-view matrix.
	mvp *lin.M4 // Scratch model-view-proj matrix.
	v0  *lin.V4 // Scratch for location calculations.
}

// newScene is expected to be called once by engine on startup.
func newScene() *scene {
	s := &scene{}
	s.scene = []*pov{}   // updated each frame
	s.pass = &layer{}    // default render target.
	s.white = newLight() // default light.
	s.mv = &lin.M4{}
	s.mvp = &lin.M4{}
	s.v0 = &lin.V4{}
	return s
}

// init is called once on engine startup to initialize the scene manager.
// Create an internal shadow map render buffer that is not tied to the
// scene graph. Create an internal shadow map specific shader not tied
// to any specific model.
func (sm *scene) init(eng *engine) {
	l := newLayer(render.DEPTH_BUFF) // internal shadow map render buffer.
	eng.loader.bindLayer(l)          // synchronously create and bind an fbo.
	sm.shadowMap = l

	// create a shadow map specific shader.
	var err error
	sms := newShader("depth")             // shadow map specific shader.
	sms, err = eng.loader.loadShader(sms) // synchronously create and bind.
	if sms == nil {
		log.Printf("scene.init: problem loading depth shader %s", err)
	}
	sm.shadowShader = sms
}

// snapshot fills a frame with draw requests. All transforms are expected
// to have been updated before calling this method.
// Snapshot flattens the Pov hierarchy in the following manner:
//   o depth first traversal of the Pov hiearchy.
//   o light/camera replace the previous light/camera.
//   o layer adds a pre-render pass for the child hierarchy.
//
// The frame memory is recycled in that the Draw records are lazy allocated
// and reused each update.  Note len(frame) is the number of draw calls
// for the most recently prepared frame.
func (sm *scene) snapshot(eng *engine, frame []render.Draw) []render.Draw {
	frame = frame[:0]       // resize keeping underlying memory.
	sm.scene = sm.scene[:0] // ditto.
	if root := eng.root(); root != nil {
		cam, _ := eng.cams[root.eid]
		sm.renDraws, sm.renVerts = 0, 0
		sm.scene = sm.updateScene(eng, 0, cam, root, sm.scene)
		frame = sm.updateFrame(eng, sm.scene, frame)
	}
	render.SortDraws(frame)
	return frame
}

// updateScene recursively turns the Pov hierarchy into a flat list using
// a depth first traversal. Pov's not affecting the rendered scene are culled.
// The important thing is to generate a consistent list of Pov's that are
// easily processed into draw calls for the render frame.
// Note that the render layer target is updated on the camera here so that
// it affects the relevant children, and not others.
func (sm *scene) updateScene(eng *engine, rt uint32, cam *camera, p *pov, scene []*pov) []*pov {
	if p.visible {
		culled := false // process children that aren't culled.

		// only calculate distance for visible models.
		if _, ok := eng.models[p.eid]; ok && cam != nil {
			px, py, pz := sm.sceneLocation(p, cam.depth)
			p.toc = cam.Distance(px, py, pz) // may not make sense for 2D screen objects.
			if culled = cam.isCulled(px, py, pz); !culled {
				scene = append(scene, p)
			}
		} else {
			scene = append(scene, p) // Keep non-model nodes.
		}

		// walk scene graph processing children of viable elements.
		if !culled {
			for _, child := range p.children {
				renderTarget := rt
				if layer, ok := eng.layers[child.eid]; ok {
					renderTarget = layer.bid // update render target layer.
				}
				if camera, ok := eng.cams[child.eid]; ok {
					cam = camera // update camera for culling.
					cam.target = renderTarget
				}
				scene = sm.updateScene(eng, renderTarget, cam, child, scene) // recurse.
			}
		}
	}
	return scene
}

// sceneLocation returns the location in world space for a 3D object,
// and in screen space for a 2D object. Assumes that a 3D objects model
// matrix has been updated.
func (sm *scene) sceneLocation(p *pov, is3D bool) (px, py, pz float64) {
	if is3D {
		vec := sm.v0.SetS(0, 0, 0, 1)
		vec.MultvM(vec, p.mm)      // Parents location incorporated into mm.
		return vec.X, vec.Y, vec.Z // 3D world space.
	} else {
		return p.Location() // 2D screen pixel space for UI culling.
	}
}

// updateFrame prepares for rendering by converting a sequenced list
// of Pov's into render system draw call requests.
func (sm *scene) updateFrame(eng *engine, viewed []*pov, frame []render.Draw) []render.Draw {
	var cam *camera                // default nil camera.
	light := sm.white              // Default light.
	lwx, lwy, lwz := 0.0, 0.0, 0.0 // Light world position.

	// turn pov's, models, and cameras into render draw requests.
	for _, p := range viewed {
		if camera, ok := eng.cams[p.eid]; ok {
			cam = camera // keep the latest camera.
		}

		// keep the latest light.
		if l, ok := eng.lights[p.eid]; ok {
			light = l
			if cam != nil {
				lx, ly, lz := p.Location()
				vec := sm.v0.SetS(lx, ly, lz, 1)
				vec.MultvM(vec, cam.vm)
				lwx, lwy, lwz = vec.X, vec.Y, vec.Z
			}
		}

		// render all models with loaded assets.
		if model, ok := eng.models[p.eid]; ok && model.loaded() {
			if model.msh != nil && len(model.msh.vdata) > 0 {
				var draw *render.Draw

				// optionally render model shadowmap from light position.
				// Its a sun light so no need to account for orientation.
				if model.castShadow {
					sm.shadowMap.vp.Set(lin.M4I)
					sm.shadowMap.vp.TranslateTM(lwx, lwy, lwz)    // (light) view
					sm.mv.Mult(p.mm, sm.shadowMap.vp)             // model-(light) view
					sm.shadowMap.vp.Mult(sm.shadowMap.vp, cam.pm) // projection.
					sm.mvp.Mult(sm.mv, sm.shadowMap.vp)           // model-view-projection

					// render the model using the shadow map "depth" shader.
					if frame, draw = sm.getDraw(frame); draw != nil {
						sm.toDraw(*draw, p, cam, model, sm.shadowMap.bid)
						shd := model.shd
						model.shd = sm.shadowShader
						model.toDraw(*draw, p.mm)
						model.shd = shd

						// capture statistics.
						sm.renDraws += 1                        // models rendered.
						sm.renVerts += model.msh.vdata[0].Len() // verticies rendered.
					}
				}

				// render model normally from camera position.
				if frame, draw = sm.getDraw(frame); draw != nil {
					sm.toDraw(*draw, p, cam, model, cam.target)
					model.toDraw(*draw, p.mm)
					light.toDraw(*draw, lwx, lwy, lwz)

					// capture statistics.
					sm.renDraws += 1                        // models rendered.
					sm.renVerts += model.msh.vdata[0].Len() // verticies rendered.
				}
			} else {
				log.Printf("Model has no mesh data... %s", model.Shader())
			}
		}
	}
	return frame
}

// toDraw sets the render data needed for a single draw call.
// The data is copied into a render.Draw instance. One of the key jobs
// of this method is to put each draw request into a particular
// render bucket so that they are drawn in order once sorted.
func (sm *scene) toDraw(d render.Draw, p *pov, cam *camera, m *model, rt uint32) {
	d.SetMv(sm.mv.Mult(p.mm, cam.vm))    // model-view
	d.SetMvp(sm.mvp.Mult(sm.mv, cam.pm)) // model-view-projection
	d.SetPm(cam.pm)                      // projection only.
	d.SetScale(p.Scale())
	d.SetTag(p.eid)

	// Set the drawing order hints. Overlay trumps transparency since 2D overlay
	// objects can't be sorted by distance anyways.
	bucket := render.OPAQUE // used to sort the draw data. Lowest first.
	switch {
	case m.castShadow && rt > 0:
		bucket = render.DEPTH_PASS // pre-passes first.
	case cam.overlay > 0:
		bucket = cam.overlay // OVERLAY draw last.
	case m.alpha < 1:
		bucket = render.TRANSPARENT // sort and draw after opaque.
	}
	depth := cam.depth && m.depth // both must be true for depth rendering.
	tocam := 0.0
	if depth {
		tocam = p.toc
	}
	d.SetHints(bucket, tocam, depth, rt)

	// use the shadow map texture for models that show shadows.
	if m.hasShadows {
		m.UseLayer(sm.shadowMap)
	}
}

// getDraw returns a render.Draw. The frame is grown as needed and draw
// instances are reused if available. Every frame value up to cap(frame)
// is expected to have already been allocated.
func (sm *scene) getDraw(frame []render.Draw) (f []render.Draw, d *render.Draw) {
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
