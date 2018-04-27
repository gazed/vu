// Copyright © 2013-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package render

import (
	"log"
)

// Data carries the buffer data that is bound/copied to the GPU.
// Data is expected to be instances from NewVertexData or NewFaceData.
// A model is often comprised of multiple sets of data.
type Data interface {
	Set(data interface{}) // Copy data in. Invalid types are logged.
	Len() int             // Number of elements.
	Size() uint32         // Number of bytes
	Clone() Data          // Duplicate.
}

// NewVertexData creates and specifies usage for a set of vertex data.
// Vertex data can be the vertex positions or per-vertex data points
// like normals or UV texture information.
// Data can now be loaded and updated using Data.Set().
//     lloc      : shader layout location index.
//     span      : values per vertex.
//     usage     : StaticDraw or DynamicDraw
//     normalize : true to normalize data to the 0-1 range.
func NewVertexData(lloc, span, usage uint32, normalize bool) Data {
	vd := &vertexData{}
	vd.floats = []float32{}
	vd.bytes = []byte{}
	vd.span = int32(span)
	vd.lloc = lloc
	vd.usage = usage
	vd.normalize = normalize
	return vd
}

// NewInstancedData creates instanced vertex data for storing 4x4 transform
// matricies. The shader lloc location value consumes 4 values in total,
// ie: if lloc is 3, then 3,4,5,6 are used to represent the 4 vectors of
// the matrix.
func NewInstancedData(lloc uint32) Data {
	vd := &vertexData{}
	vd.floats = []float32{}
	vd.bytes = []byte{} // not used.
	vd.span = 4         // 4 floats per vector in the matrix.
	vd.lloc = lloc
	vd.usage = StaticDraw
	vd.instanced = true
	return vd
}

// NewFaceData creates and specifies usagefor a set of triangle faces.
// Triangle faces contain vertex indicies ordered to draw triangles.
// Data can now be loaded and updated using Data.Set().
//     usage     : STATIC or DYNAMIC
func NewFaceData(usage uint32) Data {
	fd := &faceData{}
	fd.data = []uint16{}
	fd.usage = usage
	return fd
}

// Design Note:
// vertexData and faceData abstract away some of the render data details.
// They correspond to render layer buffer data that is eventually consumed
// by shaders.

// Data public
// =============================================================================
// vertexData

// vertexData contains per vertex information. A vertex buffer can hold
// one of float32 or byte data, but not both.
type vertexData struct {
	normalize bool      // Normalize to 0-1 range. Default false.
	span      int32     // Elements per vertex
	ref       uint32    // Vertex GPU buffer reference.
	lloc      uint32    // Shader layout location.
	usage     uint32    // STATIC_DRAW, DYNAMIC_DRAW.
	vcnt      int       // Number of verticies covered by this data.
	rebind    bool      // Data was updated and needs GPU rebind.
	instanced bool      // Data is array of matricies for instanced meshes.
	floats    []float32 // Vertex buffer arranged as [][span]float32
	bytes     []byte    // Vertex buffer arranged as [][span]byte
}

// Set makes a copy of the given data, replacing any existing data, and marks
// the data as needing to be resent to the GPU.
func (vd *vertexData) Set(data interface{}) {
	vd.vcnt = 0
	switch d := data.(type) {
	case []float32:
		vd.floats = vd.floats[:0]           // keep allocated memory.
		vd.floats = append(vd.floats, d...) // copy in new data.
		vd.vcnt = len(vd.floats) / int(vd.span)
		vd.rebind = true
	case []byte:
		vd.bytes = vd.bytes[:0]           // keep allocated memory.
		vd.bytes = append(vd.bytes, d...) // copy in new data.
		vd.vcnt = len(vd.bytes) / int(vd.span)
		vd.rebind = true
	default:
		log.Printf("vertexData.Set: invalid data type %t", d)
	}
}

// Size returns the current buffer data size in bytes.
func (vd *vertexData) Size() uint32 {
	if len(vd.floats) > 0 {
		return uint32(len(vd.floats)) * 4
	}
	if len(vd.bytes) > 0 {
		return uint32(len(vd.bytes))
	}
	return 0
}

// Len returns the number of verticies where one vertex is a set of points.
func (vd *vertexData) Len() int { return vd.vcnt }

// Clone returns a copy of the Data, including any GPU refs.
func (vd *vertexData) Clone() Data {
	c := &vertexData{}
	*c = *vd // copy by value.
	c.floats = make([]float32, len(vd.floats))
	copy(c.floats, vd.floats)
	c.bytes = make([]byte, len(vd.bytes))
	copy(c.bytes, vd.bytes)
	return c
}

// vertexData
// =============================================================================
// faceData

// faceData contains the vertex draw order. The values specify the
// order the GPU should render/processes the vertex data.
type faceData struct {
	data   []uint16 // Vertex buffer arranged as [][span]uint16.
	ref    uint32   // Vertex GPU buffer reference.
	usage  uint32   // STATIC_DRAW, DYNAMIC_DRAW.
	rebind bool     // True when data has changed and needs rebinding.
}

// Set makes a copy of the given data, replacing any existing data, and marks
// the data as needing to be resent to the GPU. Data is expected as []uint16.
func (fd *faceData) Set(data interface{}) {
	switch d := data.(type) {
	case []uint16:
		fd.data = fd.data[:0]           // keep allocated memory.
		fd.data = append(fd.data, d...) // copy in new data.
		fd.rebind = true                // Set to false when rebound.
	default:
		log.Printf("faceData.Set: invalid data type %t", d)
	}
}

// Size returns the size of the face data in bytes.
func (fd *faceData) Size() uint32 { return uint32(len(fd.data)) * 2 }

// Len returns the number of face indicies.
func (fd *faceData) Len() int { return len(fd.data) }

// Clone returns a copy of the Data, including any GPU refs.
func (fd *faceData) Clone() Data {
	c := &faceData{}
	*c = *fd // copy by value
	c.data = make([]uint16, len(fd.data))
	copy(c.data, fd.data)
	return c
}
