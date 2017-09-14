// Copyright © 2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// target.go hold code for rendering scenes to texture targets.

import (
	"log"
)

// AsTex controls the scene render output. Use on=true to render the scene
// to a texture buffer. The scene can now be used as a texture for a model
// in another scene. The default is to render the scene to the display buffer.
//
// Depends on Eng.AddScene.
func (e *Ent) AsTex(on bool) *Ent {
	if scene := e.app.scenes.get(e.eid); scene != nil {
		e.app.scenes.setTarget(scene, on)
		scene.setProjection(e.app.state.W, e.app.state.H)
		return e
	}
	log.Printf("AsTex needs AddScene %d", e.eid)
	return e
}

// =============================================================================

// target is for rendering to a framebuffer texture instead of
// the normal display buffer. It is a means of generating a texture
// from rendering instead of importing a texture. One use is for portals
// to other scenes: one scene is drawn to a texture and displayed
// in another scene.
type target struct {
	bid uint32 // Framebuffer id. Default 0 for default framebuffer.

	// Texture to render to. Used for both DepthBuffer and ImageBuffer.
	tex *Texture // Created when target is created.

	// ImageBuffer depth reference. Needed for rendering texture image
	// that uses depth to simulate rendering to a normal framebuffer.
	db uint32 // Valid for ImageBuffer.
}

// newTarget creates the framebuffer needed to render to a texture.
func newTarget() *target {
	t := &target{}
	t.tex = newTexture("renderTarget")
	return t
}

// Target is a generated asset.
func (t *target) aid() aid      { return t.tex.aid() }
func (t *target) label() string { return t.tex.name }
