// Copyright © 2015-2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// texture.go wraps images that are assocated with a model and sent
//            to the GPU for use by the shader.

import (
	"image"
)

// Texture manages the link between loaded texture assets and textures
// bound to the GPU.
//
// Texture is an an optional, but very common, part of a rendered model.
// Texture deals with 2D pictures that are mapped onto objects.
// Texture data is copied to the graphics card. One or more textures
// can be associated with a model entity and consumed by a shader.
type Texture struct {
	name   string      // Unique name of the texture.
	tag    aid         // Name and type as a number.
	img    image.Image // Texture data.
	tid    uint32      // Graphics card texture identifier.
	rebind bool        // True if data needs to be sent to the GPU.
	clamp  bool        // Set to True to trigger a one time clamp.

	// When models have more than one texture the f0, fn
	// fields are used to indicate which model faces apply to this texture.
	// First face index and number of faces.
	f0, fn uint32 // Non-zero if texture only applies to particular faces.
}

// newTexture allocates space for a texture object.
func newTexture(name string) *Texture {
	return &Texture{name: name, tag: assetID(tex, name)}
}

// aid is used to uniquely identify assets.
func (t *Texture) aid() aid      { return t.tag }  // hashed type and name.
func (t *Texture) label() string { return t.name } // asset name

// Set replaces the texture image - no questions asked.
// Marks the texture as needing to be updated on the GPU.
func (t *Texture) Set(img image.Image) {
	t.img = img
	t.rebind = true
}

// bind updates the texture on the GPU.
func (t *Texture) bind(eng *engine) error {
	t.rebind = false
	return eng.bind(t)
}
