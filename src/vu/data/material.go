// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package data

// Material is used to cover a mesh. It specifies the surface colour and
// how the surface is lit. Materials are expected to be combined with
// a Mesh.
type Material struct {
	Name string  // Unique matrial name.
	Kd   Rgb     // Diffuse colour of the material.
	Ka   Rgb     // Ambient colour of the material.
	Ks   Rgb     // Specular colour of the material.
	Tr   float32 // Transparency (alpha, dissolve) for the material.
}

// newMaterial creates a new material identified by name.  Colours can be
// provided, but if they're not, then the default colour is fully transparent
// black.
func newMaterial(name string, kd, ka, ks *Rgb, tr float32) *Material {
	if ka == nil {
		ka = &Rgb{}
	}
	if ks == nil {
		ks = &Rgb{}
	}
	if kd == nil {
		kd = &Rgb{}
	}
	return &Material{name, *kd, *ka, *ks, tr}
}

// Material
// ===========================================================================
// Rgb

// Rgb holds colour information where each of the fields is expected to contain
// a value from 0.0 to 1.0. A value of 0 means none of that colour while a value
// of 1.0 means as much as possible of that colour. For example:
//     black := &Rgb{0, 0, 0}
//     white := &Rgb{1, 1, 1}
//     red := &Rgb{1, 0, 0}
//     gray := &Rgb{0.5, 0.5, 0.5}
type Rgb struct {
	R float32 // red
	G float32 // green
	B float32 // blue
}
