// Copyright © 2014-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// mesh.go maps mesh data to the GPU and keeps the GPU reference.

import (
	"github.com/gazed/vu/render"
)

// Mesh is an optional, but very common, part of a rendered Model.
// Mesh holds 3D model data in a format that is easily consumed by rendering.
// The data consists of one or more sets of per-vertex data points and
// how the vertex positions are organized into shapes like triangles or lines.
//
// Meshes are generally loaded from assets, but can also be generated.
// Mesh data is closely tied to a given shader. When generating and refreshing
// vertex data note that InitData must be called once and SetData is called as
// needed to change the data. Data parameters are:
//    lloc     : layout location is the shader input reference.
//    span     : indicates the number of data points per vertex.
//    usage    : StaticDraw or DynamicDraw.
//    normalize: true to convert data to the 0->1 range.
// Some vertex shader data conventions are:
//    Vertex positions lloc=0 span=3_floats_per_vertex.
//    Vertex normals   lloc=1 span=3_floats_per_vertex.
//    UV tex coords    lloc=2 span=2_floats_per_vertex.
//    Color            lloc=3 span=4_floats_per_vertex.
// Note each data buffer must refer to the same number of verticies,
// and the number of verticies in one mesh must be less than 65,000.
//
// A mesh is expected to be referenced by multiple models and thus does not
// contain any instance information like location or scale. A mesh is most
// often created by the asset pipeline from disk based files that were in
// turn created by tools like Blender.
type Mesh struct {
	name   string // Unique mesh name.
	tag    aid    // name and type as a number.
	vao    uint32 // GPU reference for the mesh and all buffers.
	rebind bool   // True if data needs to be sent to the GPU.

	// Per-vertex and vertex index data.
	faces render.Data            // Triangle face indicies.
	vdata map[uint32]render.Data // Per-vertex data values.
}

// newMesh allocates space for a mesh structure,
// including space to store buffer data.
func newMesh(name string) *Mesh {
	m := &Mesh{name: name, tag: assetID(msh, name)}
	m.vdata = map[uint32]render.Data{}
	return m
}

// aid is used to uniquely identify assets.
func (m *Mesh) aid() aid      { return m.tag }  // hashed type and name.
func (m *Mesh) label() string { return m.name } // asset name
func (m *Mesh) counts() (faces, verts int) {
	if m.faces != nil {
		faces = m.faces.Len()
	}
	if m.vdata != nil && len(m.vdata) > 0 {
		verts = m.vdata[0].Len()
	}
	return faces, verts
}

// InitData creates a vertex data buffer.
func (m *Mesh) InitData(lloc, span, usage uint32, normalize bool) *Mesh {
	if _, ok := m.vdata[lloc]; !ok {
		vd := render.NewVertexData(lloc, span, usage, normalize)
		m.vdata[lloc] = vd
	}
	return m
}

// SetData stores data in the specified vertex buffer. May be called one
// or more times after a one-time call to InitData or InitInstances.
// Marks the mesh as needing a rebind.
func (m *Mesh) SetData(lloc uint32, data interface{}) {
	if _, ok := m.vdata[lloc]; ok {
		m.vdata[lloc].Set(data)
		m.rebind = true
	}
}

// InitInstances allocates space for an array of instance transform
// matricies. The instance data needs to be set using SetData where
// the float data is an array of 4x4 transform matricies.
func (m *Mesh) InitInstances(lloc uint32) {
	_, ok0 := m.vdata[lloc]
	_, ok1 := m.vdata[lloc+1]
	_, ok2 := m.vdata[lloc+2]
	_, ok3 := m.vdata[lloc+3]
	if !ok0 && !ok1 && !ok2 && !ok3 {
		vd := render.NewInstancedData(lloc)
		m.vdata[lloc] = vd
		m.vdata[lloc+1] = vd // mark shader location as used.
		m.vdata[lloc+2] = vd //  ""
		m.vdata[lloc+3] = vd //  ""
	}
}

// InitFaces creates a triangle face index buffer.
// Must be called before calling SetFaces.
func (m *Mesh) InitFaces(usage uint32) *Mesh {
	if m.faces == nil {
		m.faces = render.NewFaceData(usage)
	}
	return m
}

// SetFaces stores data for a triangle face index buffer.
// May be called one or more times after a one-time call to InitFaces.
// Marks the mesh as needing a rebind.
func (m *Mesh) SetFaces(data []uint16) {
	if m.faces != nil {
		m.faces.Set(data)
		m.rebind = true
	}
}

// bind updates the texture on the GPU.
func (m *Mesh) bind(eng *engine) error {
	m.rebind = false
	return eng.bind(m)
}

// clone makes a copy of this mesh.
func (m *Mesh) clone() (c *Mesh) {
	c = &Mesh{}
	*c = *m                   // copy by value
	c.faces = m.faces.Clone() // Triangle face indicies.
	c.vdata = map[uint32]render.Data{}
	for k, v := range m.vdata {
		c.vdata[k] = v.Clone() // Per-vertex data values.
	}
	return c
}
