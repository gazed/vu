// Copyright © 2017-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// scene.go gathers application created data and transforms the data
//          into render draw calls.

import (
	"log"

	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/render"
)

// SetUI configures a scene to be 2D and drawn over 3D scenes
// using an orthographic projection.
//
// Depends on initialization with Eng.AddScene.
func (e *Ent) SetUI() *Ent {
	if s := e.app.scenes.get(e.eid); s != nil {
		s.overlay = 1    // Draw over 3D scenes.
		s.isOrtho = true // orthographic projection
		s.cam.vt = vo    // orthographic view transform.
		s.cam.it = nv    // no inverse view transform needed.
		return e
	}
	log.Printf("SetUI needs AddScene %d", e.eid)
	return e
}

// Is2D returns true if the entity is a UI scene.
// 2D scenes are used to draw information over a 3D scene.
//
// Depends on Eng.AddScene.
func (e *Ent) Is2D() bool {
	if s := e.app.scenes.get(e.eid); s != nil {
		return s.is2D()
	}
	log.Printf("Is2D needs AddScene %d", e.eid)
	return false
}

// Cam returns the camera instace for a scene, returning nil
// if the entity is not a scene.
//
// Depends on Eng.AddScene.
func (e *Ent) Cam() *Camera {
	if s := e.app.scenes.get(e.eid); s != nil {
		return s.cam
	}
	log.Printf("Cam needs AddScene %d", e.eid)
	return nil
}

// SetScissor restricts the scene to the given screen area.
// Scissor can be disabled by setting all values to 0.
// Any negative parameter values causes the call to be ignored.
//
// Depends on Eng.AddScene.
func (e *Ent) SetScissor(x, y, w, h int) *Ent {
	if s := e.app.scenes.get(e.eid); s != nil {
		switch {
		case x < 0 || y < 0 || w < 0 || h < 0:
			// ignore.
		case x == 0 && y == 0 && w == 0 && h == 0:
			s.scissor = false
		default:
			s.scissor = true
			s.sx, s.sy = int32(x), int32(y)
			s.sw, s.sh = int32(w), int32(h)
		}
		return e
	}
	log.Printf("SetCuller needs AddScene %d", e.eid)
	return e
}

// SetCuller sets a method that reduces the number of Models rendered
// each update. It can be application supplied or engine supplied.
//
// Depends on Eng.AddScene.
func (e *Ent) SetCuller(culler Culler) {
	if s := e.app.scenes.get(e.eid); s != nil {
		s.culler = culler
		return
	}
	log.Printf("SetCuller needs AddScene %d", e.eid)
}

// SetOrtho configures the scene to use a 3D orthographic projection.
//
// Depends on Eng.AddScene.
func (e *Ent) SetOrtho() *Ent {
	if s := e.app.scenes.get(e.eid); s != nil {
		s.isOrtho = true // orthographic projection
		return e
	}
	log.Printf("SetOrtho needs AddScene %d", e.eid)
	return e
}

// SetOver draws this scene over scenes with lower values.
// Default value for a 3D scene is 0, a 2D scene is 1.
//
// Depends on Eng.AddScene.
func (e *Ent) SetOver(over uint8) *Ent {
	if s := e.app.scenes.get(e.eid); s != nil {
		s.overlay = over
		return e
	}
	log.Printf("SetOver needs AddScene %d", e.eid)
	return e
}

// scene entity methods.
// =============================================================================
// scene data

// scene contains application created resources used to render screen images.
// A scene groups one camera with a group of application created entities.
// Scene is created by the Application calling Eng.AddScene().
type scene struct {
	eid eid    // Scene and top level scene graph node.
	fbo uint32 // Render target. Default 0: display buffer.

	// Cam is this scenes camera data. Guaranteed to be non-nil.
	cam     *Camera // Created automatically with a new scene.
	isOrtho bool    // 2D or 3D orthographic projection

	// scissor the scene to be drawn within the following area.
	scissor bool // Set true to scissor the scene.
	sx, sy  int32
	sw, sh  int32

	// Overlay determines if this scene is a 3D or 2D. Default 0 means 3D.
	// Any other value means 2D. In all cases higher Overlay values are
	// drawn over scenes with lower Overlay values.
	overlay uint8  // Default 0. Add 1 to draw after other scenes.
	culler  Culler // Set by application, ie: c.Cull = vu.NewFacingCuller.
}

// newScene creates a new transform hiearchy branch with its own camera.
func newScene(eid eid, cam *Camera) *scene {
	s := &scene{eid: eid, cam: cam}
	return s
}

// is2D returns true if this scene ignores depth while rendering.
// Generally true for Overlay scenes or 2D applications.
func (s *scene) is2D() bool { return s.overlay != 0 }

// setProjection updates scenes cameras projection matrix to match
// the latest application window size.
func (s *scene) setProjection(ww, wh int) {
	w, h := float64(ww), float64(wh)
	c := s.cam
	switch {
	case s.is2D(), s.isOrtho:
		c.setOrthographic(0, w, 0, h, c.near, c.far)
	case s.fbo > 0:
		c.setPerspective(c.fov, 1.0, c.near, c.far)
	default:
		c.setPerspective(c.fov, w/h, c.near, c.far)
	}
	c.focus = true
}

// draw sets the scene attributes that affect every object within
// the scene.
func (s *scene) draw(d *render.Draw) {
	d.Fbo = s.fbo // 0 for standard display back buffer.
	d.Depth = !s.is2D()
	d.Bucket = setBucket(uint8(s.fbo), s.overlay)

	// Set the optional scissor for all draw calls in the scene.
	if s.scissor {
		d.Scissor = true
		d.Sx, d.Sy = s.sx, s.sy
		d.Sw, d.Sh = s.sw, s.sh
	}
}

// scene data.
// =============================================================================
// scene component manager.

// scenes manages all the Scene instances.
// There's not many scenes so not much to optimize.
type scenes struct {
	all      map[eid]*scene   // Scene instance data.
	shadows  map[eid]*shadows // Optional scene shadows.
	targets  map[eid]*target  // Optional scene render targets.
	skys     map[eid]*sky     // Optional sky dome.
	rebinds  map[eid][]asset  // Scene assets needing rebinds.
	released []asset          // Scene assets being disposed.

	// Scratch variables: reused each update.
	parts []uint32 // flattened pov hiearchy.
	v0    *lin.V4  // Scratch for location calculations.
	t0    *lin.T   // Scratch transform for 'tweening interpolation.
}

// newScenes creates the scene component manager and is expected to
// be called once on startup.
func newScenes() *scenes {
	ss := &scenes{}
	ss.all = map[eid]*scene{}
	ss.shadows = map[eid]*shadows{}
	ss.targets = map[eid]*target{}
	ss.skys = map[eid]*sky{}
	ss.rebinds = map[eid][]asset{}
	ss.parts = []uint32{} // updated each frame
	ss.v0 = &lin.V4{}     // scratch
	ss.t0 = lin.NewT()    // scratch
	return ss
}

// get returns the Scene associated with the given entity.
func (ss *scenes) get(id eid) *scene { return ss.all[id] }

// getTarget returns the render to texture target data.
func (ss *scenes) getTarget(id eid) *target { return ss.targets[id] }

// create makes a new scene and associates it with the given entity.
// Nothing is created if there already is a scene for the given entity.
func (ss *scenes) create(eid eid) *scene {
	scene, ok := ss.all[eid]
	if !ok {
		ss.all[eid] = newScene(eid, newCamera())
	}
	return scene // Don't allow creating over existing scene.
}

// createSky adds the data needed for a scene sky dome.
func (ss *scenes) createSky(e *Ent) *sky {
	if scene := ss.get(e.eid); scene != nil && !scene.is2D() {
		if _, ok := ss.skys[e.eid]; !ok {
			ss.skys[e.eid] = newSky(e.app)
		}
		return ss.skys[e.eid]
	}
	return nil
}

// createShadows adds shadow related data to a scene, queuing
// create requests for GPU assets.
func (ss *scenes) createShadows(s *scene, ld *loader) {
	if _, ok := ss.shadows[s.eid]; !ok {
		shadows := newShadows()
		ss.shadows[s.eid] = shadows

		// import castShadow shader, ignoring asset cache.
		shadows.caster = newShader("castShadow")
		if err := importShader(ld.loc, shadows.caster); err != nil {
			log.Printf("Missing castShadow shader: %s", err) // dev error.
		}
		ss.rebinds[s.eid] = append(ss.rebinds[s.eid], shadows.caster)
		ss.rebinds[s.eid] = append(ss.rebinds[s.eid], shadows)
	}
}

// setTarget enables or disables rendering to a texture target.
// Requests to create or delete the necessary GPU resources are queued.
func (ss *scenes) setTarget(s *scene, on bool) {
	t, ok := ss.targets[s.eid]
	if on && !ok {
		t := newTarget()
		ss.targets[s.eid] = t
		ss.rebinds[s.eid] = append(ss.rebinds[s.eid], t)
	}
	if ok && !on {
		ss.released = append(ss.released, t)
		delete(ss.targets, s.eid)
		s.fbo = 0 // draw to default display.
	}
}

// draw converts the Pov transform hierarchy to draw requests.
//
// The frame memory is recycled in that the Draw records are lazy allocated
// and reused each update. Note len(frame) is the number of draw calls
// for the most recently prepared frame.
func (ss *scenes) draw(app *application, fr frame) frame {
	fr = fr[:0] // reset keeping underlying memory.
	for _, sc := range ss.all {
		if n := app.povs.getNode(sc.eid); n != nil && !n.cull {
			index := app.povs.index[sc.eid]
			ss.parts = ss.filter(app, sc, index, ss.parts[:0])
			fr = ss.drawScene(app, sc, ss.parts, fr)
		}
	}
	render.SortDraws(fr)
	return fr
}

// filter recursively turns the Pov hierarchy into a flat list using a depth
// first traversal. Pov's not affecting the rendered scene are excluded.
//
// Problem: find Pov's to be culled - often based on camera distance.
// Output : A flat list of Pov's that need to be rendered.
func (ss *scenes) filter(app *application, sc *scene, index uint32, parts []uint32) []uint32 {
	p := app.povs.povs[index]
	n := app.povs.nodes[index]
	if culled := n.cull; !culled {
		w := p.tw.Loc
		if ready := app.models.getReady(p.eid); ready != nil {
			if sc.culler == nil || !sc.culler.Culled(sc.cam, w.X, w.Y, w.Z) {
				parts = append(parts, index)
				if !sc.is2D() {
					// save distance to camera for transparency sorting.
					ready.tocam = sc.cam.normalizedDistance(w.X, w.Y, w.Z)
				}
			}
		}

		// recurse scene graph processing children of viable elements.
		if !culled {
			for _, kid := range n.kids {
				if ki, ok := app.povs.index[kid]; ok {
					parts = ss.filter(app, sc, ki, parts)
				}
			}
		}
	}
	return parts
}

// drawScene prepares for rendering by converting a sequenced list
// of pov's into render draw call requests.
func (ss *scenes) drawScene(app *application, sc *scene, parts []uint32, f frame) frame {
	if sky, ok := ss.skys[sc.eid]; ok {
		// Optional skydome for 3D scenes.
		f = sky.draw(app, sc, f)
	}
	shadows, drawShadows := ss.shadows[sc.eid]

	// turn all the pov's, models, and cameras into render draw requests.
	var draw **render.Draw
	for _, index := range parts {
		p := &(app.povs.povs[index])

		// generate draw calls for all models with loaded assets.
		if m := app.models.getReady(p.eid); m != nil && m.msh != nil {
			if m.msh.vao <= 0 {
				log.Printf("bad mesh vao %d %s", p.eid, m.msh.name)
				continue // Engine error: should have generated vao.
			}
			if f, draw = f.getDraw(); draw != nil {
				sc.draw(*draw) // Apply scene attributes to the draw call.

				// Scenes with shadows generates extra draw calls that are sorted
				// so they are executed before the models that need the shadow map.
				// FUTURE: cast shadows for all lights.
				if drawShadows {
					lx, ly, lz := app.lights.position(sc.eid, app)

					// render from the lights position into the shadow map.
					shadows.drawShadow(*draw, p, sc, m, lx, ly, lz)
					f, draw = f.getDraw() // Need new draw call.
					sc.draw(*draw)        // Apply scene attributes to new draw call.
				}

				// render model normally from scene camera.
				// This sets the expected shader uniforms into the draw call.
				app.models.draw(p.eid, m, *draw, p, sc)
				app.lights.draw(sc.eid, *draw, app)
				if drawShadows {
					// render using the shadow map.
					shadows.drawShade(*draw, p)
				}
			}
		}
	}
	return f
}

// setPrev saves the previous locations and orientations.
// It needs to be called each state update.
func (ss *scenes) setPrev() {
	for _, scene := range ss.all {
		c := scene.cam
		c.prev.Set(c.at)
	}
}

// setRenderTransforms calculates the current render frame camera locations
// and orientations. Called each render frame.
func (ss *scenes) setRenderTransforms(lerp float64, w, h int) {
	for _, scene := range ss.all {
		c := scene.cam

		// update projection matrix if necessary.
		if c.focus {
			scene.setProjection(w, h)
			if sky, ok := ss.skys[scene.eid]; ok {
				sky.cam.pm.Set(c.pm)
			}
		}

		// update location and orientations if necessary
		if !c.focus && c.prev.Eq(c.at) {
			// Cameras that haven't moved during update already
			// have the correct transform matricies.
			continue
		}
		c.focus = false
		ss.t0.Loc.Set(c.at.Loc)
		ss.t0.Rot.Set(c.at.Rot)
		// FUTURE interpolate small camera position changes
		// ss.t0.Loc.Lerp(c.prev.Loc, c.at.Loc, lerp)
		// ss.t0.Rot.Nlerp(c.prev.Rot, c.at.Rot, lerp)

		// Set the view transform. Updates c.vm.
		c.vt(ss.t0, c.q0, c.vm) // view transform

		// Inverse only matters for perspective view transforms.
		c.it(ss.t0, c.q0, c.ivm) // inverse view transform.

		// update associated sky dome camera to use the same rotation
		// but ignore the main scene location.
		if sky, ok := ss.skys[scene.eid]; ok {
			sky.cam.prev.Set(sky.cam.at)
			sky.cam.SetPitch(c.Pitch)
			sky.cam.SetYaw(c.Yaw)
			c = sky.cam

			// interpolate the camera position and orientation.
			// This helps smooth motion when there are more renders than updates.
			ss.t0.Loc.Lerp(c.prev.Loc, c.at.Loc, lerp)
			ss.t0.Rot.Nlerp(c.prev.Rot, c.at.Rot, lerp)

			// Set the view transform. Updates c.vm. Ignore inverse.
			c.vt(ss.t0, c.q0, c.vm) // view transform
		}
	}
}

// rebind is called on the main engine thread to transfer or create
// GPU related data needed by the scenes.
func (ss *scenes) rebind(eng *engine) {
	for eid, assets := range ss.rebinds {
		for _, a := range assets {
			if err := eng.bind(a); err != nil {
				log.Printf("Bind scene asset %s failed: %s", a.label(), err)
			}
			if t, ok := a.(*target); ok {
				ss.all[eid].fbo = t.bid
			}
		}
		ss.rebinds[eid] = ss.rebinds[eid][:0] // reset preserving memory.
	}
}

// release is called on the main engine thread to remove
// GPU related data for the scenes.
func (ss *scenes) release(eng *engine) {
	for _, asset := range ss.released {
		eng.release(asset)
	}
	ss.released = ss.released[:0] // reset keeping allocated memory.
}

// resize adjusts all the cameras to the latest application window size.
// Called by the engine each time the application window changes size.
func (ss *scenes) resize(ww, wh int) {
	for _, scene := range ss.all {
		scene.setProjection(ww, wh)
	}
}

// dispose removes the scene data associated with the given entity.
// Nothing happens if there is no scene data. Returns a list of
// eids that need other components disposed.
func (ss *scenes) dispose(eid eid, dead []eid) []eid {
	delete(ss.all, eid)
	delete(ss.rebinds, eid)
	if t, ok := ss.shadows[eid]; ok {
		ss.released = append(ss.released, t)
		delete(ss.shadows, eid)
	}
	if t, ok := ss.targets[eid]; ok {
		ss.released = append(ss.released, t)
		delete(ss.targets, eid)
	}
	if sky, ok := ss.skys[eid]; ok {
		delete(ss.skys, eid)
		dead = append(dead, sky.eid)
	}
	return dead
}
