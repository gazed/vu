// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// material is used to colour a mesh. It specifies the surface colour and
// how the surface is lit. Materials are applied to a rendered model by
// a shader.
type material struct {
	name   string  // Unique matrial name.
	tag    uint64  // name and type as a number.
	kd     rgb     // Diffuse colour of the material.
	ka     rgb     // Ambient colour of the material.
	ks     rgb     // Specular colour of the material.
	tr     float32 // Transparency (alpha, dissolve) for the material.
	loaded bool    // True if data has been set.
}

// newMaterial allocates space for material values.
func newMaterial(name string) *material {
	mat := &material{name: name, tag: mat + stringHash(name)<<32}
	mat.kd.R, mat.kd.G, mat.kd.B, mat.tr = 1, 1, 1, 1
	return mat
}

// label, aid, and bid are used to uniquely identify assets.
// Note: aid is the same as bid for CPU local assets.
func (m *material) label() string { return m.name } // asset name
func (m *material) aid() uint64   { return m.tag }  // asset type and name.
func (m *material) bid() uint64   { return m.tag }  // not bound.

// setMaterial creates a new material identified by name.
// Colours can be provided, but if they're not, then the
// default colour is fully transparent black.
func (m *material) setMaterial(kd, ka, ks *rgb, tr float32) {
	m.kd.R, m.kd.G, m.kd.B = kd.R, kd.G, kd.B
	m.ks.R, m.ks.G, m.ks.B = ks.R, ks.G, ks.B
	m.ka.R, m.ka.G, m.ka.B = ka.R, ka.G, ka.B
	m.tr = tr
	m.loaded = true
}

// material
// ===========================================================================
// rgb

// rgb holds colour information where each field is expected to contain
// a value from 0.0 to 1.0. A value of 0 means none of that colour while a value
// of 1.0 means as much as possible of that colour. For example:
//     black := &rgb{0, 0, 0}     white := &rgb{1, 1, 1}
//     red   := &rgb{1, 0, 0}     gray  := &rgb{0.5, 0.5, 0.5}
type rgb struct {
	R float32 // Red.
	G float32 // Green.
	B float32 // Blue.
}

// isUnset returns true if all of the colours are zero.
func (c *rgb) isUnset() bool {
	return c.R == 0 && c.G == 0 && c.B == 0
}
