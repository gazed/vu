// Copyright © 2015-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"image"
)

// Tex is an texture. An optional, but very common, part of a rendered Model.
// Texture deals with 2D pictures that are mapped onto objects.
// Texture data is copied to the graphics card. One or more
// Textures can be associated with a Model and consumed by a Shader.
type Tex interface {
	Set(img image.Image) // Replace image, nil values ignored.
	Img() image.Image    // Get image, nil if invalid index.
	SetRepeat(on bool)   // True for Repeat, otherwise Clamp.
}

// Tex
// =============================================================================
// texture implements Tex

// texture manages the link between loaded texture assets and textures
// bound to the GPU. When models have more than one texture the f0, fn
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

// Public to satisfy Tex interface.
func (t *texture) Img() image.Image { return t.img }
func (t *texture) Set(img image.Image) {
	t.img, t.bound, t.loaded = img, false, true
}
func (t *texture) SetRepeat(on bool) { t.repeat = on }
