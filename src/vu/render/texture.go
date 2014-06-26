// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package render

import (
	"image"
)

// Texture deals with 2D pictures that are mapped onto objects. Texture data
// is copied to the graphics card. One or more Textures can be associated with
// a Model and consumed by a Shader.
type Texture interface {
	Name() string        // Unique identifier set on creation.
	Img() image.Image    // Texture image.
	Set(img image.Image) // Set the loaded or generated texture data.
	Bound() bool         // True if the mesh has a GPU reference.
	FreeImg()            // Used to release the image data after binding.
}

// For an overview of opengl textures see:
//    http://www.arcsynthesis.org/gltut/Texturing/Tutorial%2014.html

// ============================================================================

// texture is the default implementation of Texture
type texture struct {
	name   string      // Unique name of the texture.
	img    image.Image // Texture data. Release (set to nil) after GPU binding.
	tid    uint32      // Graphics card texture identifier.
	refs   uint32      // Number of Model references.
	repeat bool        // Repeat the texture when UV greater than 1.
}

// newTexture allocates space for a texture object.
func newTexture(name string) *texture { return &texture{name: name} }

// Implement Texture.
func (t *texture) Name() string        { return t.name }
func (t *texture) Img() image.Image    { return t.img }
func (t *texture) Set(img image.Image) { t.img = img }
func (t *texture) Bound() bool         { return t.tid != 0 }
func (t *texture) FreeImg()            { t.img = nil }
func (t *texture) SetRepeat(on bool)   { t.repeat = on }
