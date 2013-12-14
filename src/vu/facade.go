// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package vu

// facade combines the an objects graphic resources into a viewable shape.
// It handles the appearance of an object. A facade is expected to be attached
// to a part within a scene in order to be rendered.  Geneally a facade is
// specified by the application and consumed by the engine rendering subsystem.
//
// The facades resources are lazy-loaded from the resource depot.
type facade struct {
	mesh   string  // Facade frame.
	shader string  // Shape painter.
	mat    string  // Optional material (includes alpha value)
	tex    string  // Optional texture.
	alpha  float64 // Optional override alpha value between 0 & 1.
	rots   float64 // Texture rotation speed.

	glyphs string // Optional for a banner. A texture, and text must be specified.
	text   string // Banner text.
}

// newFacade creates a facade with the given mesh, shader, and material.
// Textures can optionally be included later using addTexture().  Each facade
// is given a reference to the roadie which is used to load and cache the
// needed resources.
//
// Usage: facades are created by the application using Part.SetFacade.
func newFacade(mesh, shader string) *facade {
	f := &facade{}
	f.mesh = mesh
	f.shader = shader
	f.alpha = 1
	return f
}

// newBanner creates a banner facade.
//    text: the banner text.
//    shader, glyphs, texture : resource identifiers.
func newBanner(text, shader, glyphs, texture string) *facade {
	f := newFacade("banner", shader)
	f.glyphs = glyphs
	f.tex = texture
	f.text = text
	return f
}
