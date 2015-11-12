// Copyright © 2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"image"
)

// texture deals with 2D pictures that are mapped onto objects.
// Texture data is copied to the graphics card. One or more
// Textures can be associated with a Model and consumed by a Shader.
// Models can have more than one texture. In this case the f0, fn
// fields are used to indicate which model faces apply to this texture.
type texture struct {
	name   string      // Unique name of the texture.
	tag    uint64      // Name and type as a number.
	img    image.Image // Texture data.
	tid    uint32      // Graphics card texture identifier.
	repeat bool        // Repeat the texture when UV greater than 1.
	bound  bool        // False if the data needs rebinding.
	loaded bool        // True if data has been set.

	// First face index and number of faces.
	// Used for multiple uv textures for the same model.
	f0, fn uint32 // Non-zero if texture only applies to particular faces.
}

// newTexture allocates space for a texture object.
func newTexture(name string) *texture {
	return &texture{name: name, tag: tex + stringHash(name)<<32}
}

// label, aid, and bid are used to uniquely identify assets.
func (t *texture) label() string { return t.name }                  // asset name
func (t *texture) aid() uint64   { return t.tag }                   // asset type and name.
func (t *texture) bid() uint64   { return tex + uint64(t.tid)<<32 } // asset type and bind ref.

// set texture image data and render attributes.
func (t *texture) set(img image.Image) {
	t.img = img
	t.bound = false
	t.loaded = true
}
func (t *texture) setRepeat(on bool) { t.repeat = on }
