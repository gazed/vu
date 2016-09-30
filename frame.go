// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// frame.go transforms application created data into render draw calls.
// FUTURE: keep frames.scenes up to date rather than regenerate it each update.

import (
	"log"
	"time"

	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/render"
)

// frame links the engine models and cameras to a render frame.
// It generates an ordered list of render system draw data needed for
// one screen. There is always one frame that is being updated while
// the two previous frames are being rendered.
type frame []*render.Draw

// frames creates the list of render draw calls from Pov's, Models,
// and Cameras.
//
// There are three render frames. One for updating, the other two for
// rendering with interpolation.
type frames struct {
	snap  frame      // update draw calls for the next render frame.
	draw  chan frame // frame to be updated returned from machine.
	scene []*Pov     // flattened pov hiearchy updated each frame.

	// Track the number of draw calls, and verticies.
	drawCalls int // Number of models rendered last update.
	verticies int // Number of verticies rendered last update.

	// Scratch variables: reused to reduce garbage collection.
	white *Light  // default light.
	mv    *lin.M4 // Scratch model-view matrix.
	mvp   *lin.M4 // Scratch model-view-proj matrix.
	v0    *lin.V4 // Scratch for location calculations.
}

// newFrames is expected to be called once by engine on startup.
func newFrames() *frames {
	fs := &frames{}
	fs.scene = []*Pov{}   // updated each frame
	fs.white = newLight() // default light.
	fs.mv = &lin.M4{}
	fs.mvp = &lin.M4{}
	fs.v0 = &lin.V4{}
	return fs
}

// render sends the an updated frame off for rendering. If the
// frame is not ready the interpolation data is sent instead.
func (fs *frames) render(machine chan msg, interp float64, ut uint64) {
	if len(fs.snap) > 0 { // is new frame ready?
		machine <- &renderFrame{fr: fs.snap, interp: interp, ut: ut}
		fs.snap = <-fs.draw   // replace the render frame with an old frame.
		fs.snap = fs.snap[:0] // ... and mark it as unpreprepared.
	} else {
		machine <- &renderFrame{fr: nil, interp: interp, ut: ut}
	}
}

// drawFrame fills a frame with draw requests. All transforms
// are expected to have been updated before calling this method.
// Snapshot flattens the Pov hierarchy in the following manner:
//   o depth first traversal of the Pov hiearchy.
//   o light/camera replace the previous light/camera.
//   o layer adds a pre-render pass for the child hierarchy.
//
// The frame memory is recycled in that the Draw records are lazy allocated
// and reused each update.  Note len(frame) is the number of draw calls
// for the most recently prepared frame.
func (fs *frames) drawFrame(eng *engine) {
	fs.snap = fs.snap[:0]   // resize keeping underlying memory.
	fs.scene = fs.scene[:0] // ditto.
	if root := eng.root(); root != nil {
		cam := eng.cams.get(root.id)
		fs.drawCalls, fs.verticies = 0, 0
		fs.scene = fs.updateScene(eng, 0, cam, root, fs.scene)
		fs.snap = fs.updateFrame(eng, fs.scene, fs.snap)
	}
	render.SortDraws(fs.snap)
}

// updateScene recursively turns the Pov hierarchy into a flat list using
// a depth first traversal. Pov's not affecting the rendered scene are culled.
// The traversal is necessary to find which camera affects which Pov's.
//
// Input  : Pov hierarchy figures out which camera is active.
// Problem: find Pov's to be culled - often based on camera distance.
// Output : A flat list of Pov's that need to be rendered.
func (fs *frames) updateScene(eng *engine, rt uint32, cam *Camera, p *Pov, scene []*Pov) []*Pov {
	if !p.Cull {
		culled := false // process children that aren't culled.

		// only calculate distance for visible models.
		if m := eng.models.get(p.id); m != nil && cam != nil {
			px, py, pz := fs.sceneLocation(p, cam.Depth)
			p.toc = cam.Distance(px, py, pz) // may not make sense for 2D screen objects.
			if culled = cam.isCulled(px, py, pz); !culled {
				scene = append(scene, p)
			}
		} else {
			scene = append(scene, p) // Keep non-model nodes.
		}

		// walk scene graph processing children of viable elements.
		if !culled {
			for _, child := range p.kids {
				renderTarget := rt
				if layer := eng.layers.get(child.id); layer != nil {
					renderTarget = layer.bid // update render target layer.
				}
				if camera := eng.cams.get(child.id); camera != nil {
					cam = camera // update camera for culling.
					cam.target = renderTarget
				}
				scene = fs.updateScene(eng, renderTarget, cam, child, scene) // recurse.
			}
		}
	}
	return scene
}

// sceneLocation returns the location in world space for a 3D object,
// and in screen space for a 2D object. Assumes that a 3D objects model
// matrix has been updated.
func (fs *frames) sceneLocation(p *Pov, is3D bool) (px, py, pz float64) {
	if is3D {
		vec := fs.v0.SetS(0, 0, 0, 1)
		vec.MultvM(vec, p.mm)      // Parents location incorporated into mm.
		return vec.X, vec.Y, vec.Z // 3D world space.
	}
	return p.At() // 2D screen pixel space for UI culling.
}

// updateFrame prepares for rendering by converting a sequenced list
// of Pov's into render system draw call requests.
func (fs *frames) updateFrame(eng *engine, viewed []*Pov, f frame) frame {
	var cam *Camera                // default nil camera.
	light := fs.white              // Default light.
	lwx, lwy, lwz := 0.0, 0.0, 0.0 // Light world position.

	// turn pov's, models, and cameras into render draw requests.
	for _, p := range viewed {
		if camera := eng.cams.get(p.id); camera != nil {
			cam = camera // keep the latest camera.
		}

		// keep the latest light.
		if l := eng.lights.get(p.id); l != nil {
			light = l
			if cam != nil {
				lx, ly, lz := p.At()
				vec := fs.v0.SetS(lx, ly, lz, 1)
				vec.MultvM(vec, cam.vm)
				lwx, lwy, lwz = vec.X, vec.Y, vec.Z
			}
		}

		// render all models with loaded assets.
		if model := eng.models.getActive(p.id); model != nil {
			if model.msh != nil && len(model.msh.vdata) > 0 {
				var draw **render.Draw
				if f, draw = fs.getDraw(f); draw != nil {

					// optionally render model shadowmap from light position.
					// Its a sun light so no need to account for orientation.
					if model.castShadow {
						eng.layers.renderShadow(*draw, p, cam, model, lwx, lwy, lwz)
						fs.drawCalls++                           // models rendered stat.
						fs.verticies += model.msh.vdata[0].Len() // verticies rendered stat.
						f, draw = fs.getDraw(f)                  // need new draw call.
					}

					// always render model normally from camera position.
					drawPov(*draw, p, fs.mv, fs.mvp, cam, model, cam.target)
					drawModel(*draw, model, p.mm, eng.layers.shadows)
					drawLight(*draw, light, lwx, lwy, lwz)
					fs.drawCalls++                           // models rendered stat.
					fs.verticies += model.msh.vdata[0].Len() // verticies rendered stat.
				}
			} else {
				log.Printf("Model has no mesh data... %s", model.Shader())
			}
		}
	}
	return f
}

// getDraw returns a render.Draw. The frame is grown as needed and draw
// instances are reused if available. Every frame value up to cap(frame)
// is expected to have already been allocated.
func (fs *frames) getDraw(f frame) (frame, **render.Draw) {
	size := len(f)
	switch {
	case size == cap(f):
		f = append(f, render.NewDraw())
	case size < cap(f): // use previously allocated.
		f = f[:size+1]
		if f[size] == nil {
			f[size] = render.NewDraw()
		}
	}
	return f, &f[size]
}

// frames
// =============================================================================
// utility methods to transform engine data to render draw calls.

// drawPov sets the render data needed for a single draw call.
// The data is copied into a render.Draw instance. One of the key jobs
// of this method is to put each draw request into a particular
// render bucket so that they are drawn in order once sorted.
func drawPov(d *render.Draw, p *Pov, mv, mvp *lin.M4, cam *Camera, m *model, rt uint32) {
	d.SetMv(mv.Mult(p.mm, cam.vm)) // model-view
	d.SetMvp(mvp.Mult(mv, cam.pm)) // model-view-projection
	d.SetPm(cam.pm)                // projection only.
	d.SetScale(p.Scale())
	d.Tag = uint64(p.id)

	// Set the drawing order hints. Overlay trumps transparency since 2D overlay
	// objects can't be sorted by distance anyways.
	bucket := render.Opaque // used to sort the draw data. Lowest first.
	switch {
	case m.castShadow && rt > 0:
		bucket = render.DepthPass // pre-passes first.
	case cam.Overlay > 0:
		bucket = cam.Overlay // OVERLAY draw last.
	case m.alpha < 1:
		bucket = render.Transparent // sort and draw after opaque.
	}
	depth := cam.Depth && m.depth // both must be true for depth rendering.
	tocam := 0.0
	if depth {
		tocam = p.toc
	}
	d.SetHints(bucket, tocam, depth, rt)
}

// drawModel sets the model specific bound data references and
// uniform data needed by the rendering layer.
func drawModel(d *render.Draw, m *model, mm *lin.M4, shadows *layer) {

	// Use any previous render to texture passes.
	if m.layer != nil {
		switch m.layer.attr {
		case render.ImageBuffer:
			// handled as regular texture below.
			// Leave it to the shader to use the right the "uv#" uniform.
		case render.DepthBuffer:
			d.SetShadowmap(m.layer.tex.tid) // texture with depth values.

			// Shadow depth bias is the mvp matrix from the light.
			// It is adjusted as needed by shadow maps.
			m.sm.Mult(mm, m.layer.vp)   // model (light) view.
			m.sm.Mult(m.sm, m.layer.bm) // incorporate shadow bias.
			d.SetDbm(m.sm)
		}
	}

	// use the shadow map texture for models that show shadows.
	if m.hasShadows {
		m.UseLayer(shadows)
	}

	// Set the bound data references.
	d.SetRefs(m.shd.program, m.msh.vao, m.drawMode)
	if total := len(m.texs); total > 0 {
		for cnt, t := range m.texs {
			d.SetTex(total, cnt, t.tid, t.f0, t.fn)
		}
	} else {
		d.SetTex(0, 0, 0, 0, 0) // clear any previous data.
	}

	// Set uniform values. These can be sent as a reference because they
	// are fixed on shader creation.
	d.SetUniforms(m.shd.uniforms) // shader integer uniform references.
	if m.anm != nil && len(m.pose) > 0 {
		d.SetPose(m.pose)
	} else {
		d.SetPose(nil) // clear data.
	}

	// Material transparency.
	d.SetFloats("alpha", float32(m.alpha))

	// Material color uniforms.
	if mat := m.mat; mat != nil {
		drawMaterial(d, mat)
	}

	// For shaders that need elapsed time.
	d.SetFloats("time", float32(time.Since(m.time).Seconds()))

	// Set user specified uniforms.
	for uniform, uvalues := range m.uniforms {
		d.SetFloats(uniform, uvalues...)
	}
}

// drawMaterial sets the data needed by the render system.
func drawMaterial(d *render.Draw, m *material) {
	d.SetFloats("kd", m.kd.R, m.kd.G, m.kd.B)
	d.SetFloats("ks", m.ks.R, m.ks.G, m.ks.B)
	d.SetFloats("ka", m.ka.R, m.ka.G, m.ka.B)
	d.SetFloats("ns", m.ns)
}

// drawLight sets the data needed by the render system.
// In this case the light color.
func drawLight(d *render.Draw, l *Light, px, py, pz float64) {
	d.SetFloats("lp", float32(px), float32(py), float32(pz))    // position
	d.SetFloats("lc", float32(l.R), float32(l.G), float32(l.B)) // color
}
