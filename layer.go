// Copyright © 2015-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// layer.go encapsulates anything related to generating render information
// prior to creating the final screen framebuffer. A layer is often referred
// to as a render pass.
// DESIGN: layers depend on Framebuffers allocated on the GPU and associated
//         with one or more textures.
// Transforms:
//         Creates bound layers.
//         Holds and creates layer data needed by shaders.
// FUTURE: Lots more to do for render targets/passes/layers.
//         Need cleaner design *and* ideally reduced code complexity.
//         Need to better understand the render pass use-cases. See:
//            https://docs.unity3d.com/Manual/SL-RenderPipeline.html
//         for different render scenarios to consider.

import (
	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/render"
)

// Layer is used to render to a 1024x1024 sized frame buffer based texture.
// A layer represents the output of rendering to a texture instead of to
// the default framebuffer for the screen. The layer texture is then used
// as part of the final screen render. There are currently 2 uses for layers:
// a render.ImageBuffer and render.DepthBuffer.
//    ImageBuffer: renders all or parts of a scene to a texture.
//                 Often used to generate a texture which is then used
//                 as a uv texture for another model. Think Portal.
//    DepthBuffer: render objects from the point of view of the
//                 light. Distance (depth) values are written to the
//                 texture and later used by a shadow capabable shader
//                 to draw object shadows.
// Layers are created, associated with one Pov, and used by a Model at
// a potentially different Pov.
type Layer interface{}

// Layer
// =============================================================================
// layer could split into two layer types.

// layer implements Layer. Currently used as an ImageBuffer or DepthBuffer.
type layer struct {
	bid  uint32 // Framebuffer id. Default 0 for default framebuffer.
	attr int    // What type of layer. DepthBuffer or ImageBuffer.

	// Texture to render to. Used for both DepthBuffer and ImageBuffer.
	tex *texture // Created when layer is created.

	// ImageBuffer depth reference. Needed for rendering texture image
	// that uses depth to simulate rendering to a normal framebuffer.
	db uint32 // Valid for ImageBuffer.

	// Bias matrix is used by shadow shader reading from the depthBuffer.
	bm *lin.M4 // bias matrix needed for a shadow shader.

	// The view-projection transform for the light. Needed by
	// the depth shader to render objects from the point of view
	// of the light instead of the eye.
	vp *lin.M4 // light view-projection layer transform.
}

// newLayer creates the framebuffer needed to render to a texture.
func newLayer(attr int) *layer {
	l := &layer{attr: attr}
	l.vp = &lin.M4{}
	l.bm = &lin.M4{
		Xx: 0.5, Xy: 0.0, Xz: 0.0, Xw: 0.0,
		Yx: 0.0, Yy: 0.5, Yz: 0.0, Yw: 0.0,
		Zx: 0.0, Zy: 0.0, Zz: 0.5, Zw: 0.0,
		Wx: 0.5, Wy: 0.5, Wz: 0.5, Ww: 1.0,
	}
	l.tex = newTexture("layer")
	return l
}

// Layer is a generated asset.
func (l *layer) aid() aid      { return l.tex.aid() }
func (l *layer) label() string { return l.tex.name }

// layer
// =============================================================================
// layers

// layers manages all the active layer instances. Currently there is not
// expected to be very many layers.
type layers struct {
	eng  *engine        // Needed to bind layer.
	data map[eid]*layer // Camera instance data.
	bind chan msg       // For binding.

	// Shadow render pass support. A single layer is created that collects
	// data for each model with shadows.
	shadows      *layer  // scratch shadow map.
	shadowShader *shader // shader used to render to depth buffer.
	mv           *lin.M4 // Scratch model-view matrix.
	mvp          *lin.M4 // Scratch model-view-proj matrix.
}

// newLayers creates the layer component manager instance. Expected to
// be called once on startup.
func newLayers(eng *engine) *layers {
	ls := &layers{eng: eng, data: map[eid]*layer{}}
	ls.bind = eng.machine
	ls.mv = &lin.M4{}
	ls.mvp = &lin.M4{}
	return ls
}

// get returns the layer for the given entity, or nil if there is
// no layer for this entity.
func (ls *layers) get(id eid) *layer {
	if l, ok := ls.data[id]; ok {
		return l
	}
	return nil
}

// create ensures there is only one layer per entity.
func (ls *layers) create(id eid, attr int) *layer {
	if l, ok := ls.data[id]; ok {
		return l // Don't allow creating over existing layer.
	}
	l := newLayer(attr)
	ls.bindLayer(l) // synchronously create and bind a fbo.
	ls.data[id] = l
	return l
}

// bindLayer requests a new framebuffer based texture for a view.
func (ls *layers) bindLayer(layer *layer) error {
	bindReply := make(chan error)
	ls.bind <- &bindData{data: layer, reply: bindReply} // request bind.
	return <-bindReply                                  // wait for bind.
}

// dispose removes the render pass layer, if any, from the given entity.
// No complaints if there is no layer at the given entity. This is safe to
// remove from the GPU since it is not cached/shared.
func (ls *layers) dispose(id eid) {
	if l, ok := ls.data[id]; ok {
		delete(ls.data, id)
		ls.eng.release(&releaseData{data: l}) // dispose of the framebuffer.
	}
}

// enableShadows creates a layer specifically for shadow maps.
// It is called on engine startup before other load requests so
// that the return load value is guaranteed to match the request.
// The shadowmap buffer and shadow map render are not tied
// to any specific model or Pov.
func (ls *layers) enableShadows() {
	shaderName := "depth" // engine shader writes shadow map depth values.
	id := assetID(shd, shaderName)
	reqs := map[aid]string{id: shaderName}
	ls.eng.submitLoadReqs(reqs) // eng submits to loading goroutine.
	assets := <-ls.eng.loaded   // wait for load to complete.
	ls.shadowShader = assets[id].(*shader)

	// create the shadow map render buffer.
	l := newLayer(render.DepthBuffer)
	ls.bindLayer(l) // synchronously create and bind a fbo.
	ls.shadows = l
}

// render the models shadow from light position.
// Its a sun light so no need to account for orientation.
func (ls *layers) renderShadow(draw *render.Draw, p *Pov, cam *Camera,
	model *model, lx, ly, lz float64) {
	ls.shadows.vp.Set(lin.M4I)
	ls.shadows.vp.TranslateTM(lx, ly, lz)     // (light) view
	ls.mv.Mult(p.mm, ls.shadows.vp)           // model-(light) view
	ls.shadows.vp.Mult(ls.shadows.vp, cam.pm) // projection.
	ls.mvp.Mult(ls.mv, ls.shadows.vp)         // model-view-projection

	// render the model using the shadow map "depth" shader.
	drawPov(draw, p, ls.mv, ls.mvp, cam, model, ls.shadows.bid)
	shd := model.shd
	model.shd = ls.shadowShader
	drawModel(draw, model, p.mm, ls.shadows)
	model.shd = shd
}
