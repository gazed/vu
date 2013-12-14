// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package data

import (
	"image"
)

// Texture deals with 2D pictures that are mapped onto objects. Textures are
// copied to the graphics card and expected to be combined with a Mesh.
type Texture struct {
	Name string      // Unique name of the texture.
	Img  image.Image // Texture data.
	Tid  uint32      // Graphics card texture identifier.
}
