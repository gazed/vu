// Copyright © 2015-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// texture.go
// FUTURE: some textures are used both as repeating and not-repeating.
//         Currently the texture needs to be loaded twice so that one
//         can be marked as repeating.
//         Look at setting repeat texture attribute as part of the draw
//         call so that the texture data can be reused.

import (
	"image"
)

// Tex is an texture. An optional, but very common, part of a rendered Model.
// Texture deals with 2D pictures that are mapped onto objects.
// Texture data is copied to the graphics card. One or more
// Textures can be associated with a Model and consumed by a Shader.
type Tex interface {
	Set(img image.Image) // Replace image, nil values ignored.
}

// Tex
// =============================================================================
// texture implements Tex

// texture manages the link between loaded texture assets and textures
// bound to the GPU. When models have more than one texture the f0, fn
// fields are used to indicate which model faces apply to this texture.
type texture struct {
	name string      // Unique name of the texture.
	tag  aid         // Name and type as a number.
	img  image.Image // Texture data.
	tid  uint32      // Graphics card texture identifier.

	// First face index and number of faces.
	// Used for multiple uv textures for the same model.
	f0, fn uint32 // Non-zero if texture only applies to particular faces.
}

// newTexture allocates space for a texture object.
func newTexture(name string) *texture {
	return &texture{name: name, tag: assetID(tex, name)}
}

// aid is used to uniquely identify assets.
func (t *texture) aid() aid      { return t.tag }  // hashed type and name.
func (t *texture) label() string { return t.name } // asset name

// Public to satisfy Tex interface.
func (t *texture) Set(img image.Image) { t.img = img }

// =============================================================================
// texture

// texid is needed to track textures prior to them being loaded.
// Currently only needed/used by model.
type texid struct {
	name string // Unique name of the asset.
	id   aid    // Name and type as a number.
}

func newTexid(name string, id aid) *texid {
	return &texid{name: name, id: id}
}

// Implement asset interface.
func (t *texid) aid() aid      { return t.id }   // hashed type and name.
func (t *texid) label() string { return t.name } // asset name
