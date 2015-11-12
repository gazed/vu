// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"github.com/gazed/vu/render"
)

// mesh holds 3D model data in a format that is easily consumed by a rendering
// layer. The data consists of one or more sets of per-vertex data points and
// how the vertex positions are organized into shapes like triangles or lines.
// A mesh is expected to be referenced by multiple models and thus does not
// contain any instance information like location or scale. A mesh is most
// often created by the asset pipeline from disk based files that were in turn
// created by tools like Blender.
//
// Note each data buffer must refer to the same number of verticies,
// and the number of verticies in one mesh must be less than 65,000.
type mesh struct {
	name   string // Unique mesh name.
	tag    uint64 // name and type as a number.
	vao    uint32 // GPU reference for the mesh and all buffers.
	bound  bool   // False if the data needs rebinding.
	loaded bool   // True if data has been set.

	// Per-vertex and vertex index data.
	faces render.Data            // Triangle face indicies.
	vdata map[uint32]render.Data // Per-vertex data values.
}

// newMesh allocates space for a mesh structure,
// including space to store buffer data.
func newMesh(name string) *mesh {
	m := &mesh{name: name, tag: msh + stringHash(name)<<32}
	m.vdata = map[uint32]render.Data{}
	return m
}

// label, aid, and bid are used to uniquely identify assets.
func (m *mesh) label() string { return m.name }                  // asset name
func (m *mesh) aid() uint64   { return m.tag }                   // asset type and name.
func (m *mesh) bid() uint64   { return msh + uint64(m.vao)<<32 } // asset type and bind ref.

// initData creates a vertex data buffer.
func (m *mesh) initData(lloc, span, usage uint32, normalize bool) *mesh {
	if _, ok := m.vdata[lloc]; !ok {
		vd := render.NewVertexData(lloc, span, usage, normalize)
		m.vdata[lloc] = vd
	}
	return m
}

// setData stores data in the specified vertex buffer.
func (m *mesh) setData(lloc uint32, data interface{}) {
	if _, ok := m.vdata[lloc]; ok {
		m.vdata[lloc].Set(data)
		m.loaded = true
	}
}

// initFaces creates a triangle face index buffer.
func (m *mesh) initFaces(usage uint32) *mesh {
	if m.faces == nil {
		m.faces = render.NewFaceData(usage)
	}
	return m
}

// setFaces stores data for a triangle face index buffer.
func (m *mesh) setFaces(data []uint16) {
	if m.faces != nil {
		m.faces.Set(data)
	}
}
