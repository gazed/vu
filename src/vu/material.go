// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package vu

// material is used to colour a mesh. It specifies the surface colour and
// how the surface is lit. Materials are expected to be combined with
// a Mesh.
type material struct {
	name string  // Unique matrial name.
	kd   rgb     // Diffuse colour of the material.
	ka   rgb     // Ambient colour of the material.
	ks   rgb     // Specular colour of the material.
	tr   float32 // Transparency (alpha, dissolve) for the material.
}

func newMaterial(name string) *material {
	mat := &material{name: name}
	mat.kd.R, mat.kd.G, mat.kd.B, mat.tr = 1, 1, 1, 1
	return mat
}

// Name implements Material
func (m *material) Name() string { return m.name }

// setMaterial creates a new material identified by name.  Colours can be
// provided, but if they're not, then the default colour is fully transparent
// black.
func (m *material) SetMaterial(kd, ka, ks *rgb, tr float32) {
	m.kd.R, m.kd.G, m.kd.B = kd.R, kd.G, kd.B
	m.ks.R, m.ks.G, m.ks.B = ks.R, ks.G, ks.B
	m.ka.R, m.ka.G, m.ka.B = ka.R, ka.G, ka.B
	m.tr = tr
}

func (m *material) SetKd(r, g, b float32) { m.kd.R, m.kd.G, m.kd.B = r, g, b }
func (m *material) SetKs(r, g, b float32) { m.ks.R, m.ks.G, m.ks.B = r, g, b }
func (m *material) SetKa(r, g, b float32) { m.ka.R, m.ka.G, m.ka.B = r, g, b }
func (m *material) SetTr(tr float32)      { m.tr = tr }
