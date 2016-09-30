// Copyright © 2014-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// material.go handles model surface color data.

// material is used to color a mesh. It specifies the surface color and
// how the surface is lit. Materials are applied to a rendered model by
// a shader.
type material struct {
	name string  // Unique matrial name.
	tag  aid     // name and type as a number.
	kd   rgb     // Diffuse color of the material.
	ka   rgb     // Ambient color of the material.
	ks   rgb     // Specular color of the material.
	ns   float32 // Specular sharpness.
	tr   float32 // Transparency (alpha, dissolve) for the material.
}

// newMaterial allocates space for material values.
func newMaterial(name string) *material {
	mat := &material{name: name, tag: assetID(mat, name)}
	mat.kd.R, mat.kd.G, mat.kd.B, mat.tr = 1, 1, 1, 1
	return mat
}

// aid is used to uniquely identify assets.
func (m *material) aid() aid      { return m.tag }  // hashed type and name.
func (m *material) label() string { return m.name } // asset name

// material
// ===========================================================================
// rgb

// rgb holds color information where each field is expected to contain
// a value from 0.0 to 1.0. A value of 0 means none of that color while a value
// of 1 means as much as possible of that color. For example:
//     black := &rgb{0, 0, 0}     white := &rgb{1, 1, 1}
//     red   := &rgb{1, 0, 0}     gray  := &rgb{0.5, 0.5, 0.5}
type rgb struct {
	R float32 // Red.
	G float32 // Green.
	B float32 // Blue.
}

// isBlack returns true if all of the color are zero.
// This may indicate the material has not be set.
func (c *rgb) isBlack() bool { return c.R == 0 && c.G == 0 && c.B == 0 }
