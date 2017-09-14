// Copyright © 2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// shadow.go holds code for generating and rendering shadows.

import (
	"log"

	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/render"
)

// SetShadows enables models casting and receiving shadows for
// a scene component. Requests for shadows on 2D scenes are logged.
//
// Depends on Eng.AddScene.
func (e *Ent) SetShadows() *Ent {
	if scene := e.app.scenes.get(e.eid); scene != nil && !scene.is2D() {
		e.app.scenes.createShadows(scene, e.app.ld)
		return e
	}
	log.Printf("SetShadows needs AddScene %d", e.eid)
	return e
}

// =============================================================================

// shadows can be added to a camera controlled scene.
// Each object in the scene will generate a normal draw call and
// a separate second draw call using a castShadow shader.
// The castShadow shader renders the object into a shadowmap
// buffer using the point of view of the scene light.
//
// Each object receiving shadows needs to be rendered using a
// showShadow shadow. The showShadow shader expects the shadow map
// created by castShadow.
type shadows struct {
	bid uint32 // Framebuffer id. Default 0 for default framebuffer.

	// Texture to store shadowmap depth values.
	tex *Texture // Created in newShadows.

	// Bias matrix is used by cast shadow shader reading from the depthBuffer.
	bm *lin.M4 // bias matrix needed for a shadow shader.

	// The view-projection transform for the light. Needed by
	// the depth shader to render objects from the point of view
	// of the light instead of the eye.
	vp *lin.M4 // light view-projection transform.

	caster *shader // shader used to render to depth buffer.
	mv     *lin.M4 // Scratch model-view matrix.
	mvp    *lin.M4 // Scratch model-view-proj matrix.
	sm     *lin.M4 // Scratch matrix.
}

// newShadows creates the resources needed to render shadows.
func newShadows() *shadows {
	s := &shadows{}
	s.vp = &lin.M4{}
	s.bm = &lin.M4{
		Xx: 0.5, Xy: 0.0, Xz: 0.0, Xw: 0.0,
		Yx: 0.0, Yy: 0.5, Yz: 0.0, Yw: 0.0,
		Zx: 0.0, Zy: 0.0, Zz: 0.5, Zw: 0.0,
		Wx: 0.5, Wy: 0.5, Wz: 0.5, Ww: 1.0,
	}
	s.tex = newTexture("shadowMap")
	s.mv = &lin.M4{}
	s.mvp = &lin.M4{}
	return s
}

// Shadows is a generated asset.
func (s *shadows) aid() aid      { return s.tex.aid() }
func (s *shadows) label() string { return s.tex.name }

// drawShadow renders the models shadow from light position.
// Its a sun light so no need to account for orientation.
func (s *shadows) drawShadow(draw *render.Draw, p *pov, sc *scene, m *model) {
	lx, ly, lz := sc.light.wx, sc.light.wy, sc.light.wz

	// calculate draw transforms using lights position instead of camera.
	s.vp.Set(lin.M4I)
	s.vp.TranslateTM(lx, ly, lz) // (light) view
	s.mv.Mult(p.mm, s.vp)        // model-(light) view
	s.vp.Mult(s.vp, sc.cam.pm)   // projection.
	s.mvp.Mult(s.mv, s.vp)       // model-view-projection

	// render using the shader that creates the shadow map.
	p.draw(draw, s.mv, s.mvp, sc, sc.cam, s.bid)
	m.draw(draw, s.caster) // override default shader with shadow castor.
}

// drawShade uses the shadow map to show shadows on an object.
func (s *shadows) drawShade(d *render.Draw, p *pov) {
	d.SetShadowmap(s.tex.tid)
	// Shadow depth bias is the mvp matrix from the light.
	// It is adjusted as needed by shadow maps.
	s.mv.Mult(p.mm, s.vp) // model (light) view.
	s.mv.Mult(s.mv, s.bm) // incorporate shadow bias.
	d.SetDbm(s.mv)
}
