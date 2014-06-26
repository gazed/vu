// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package render

// vertexData and faceData abstract away some of the render data details.
// They correspond to render layer buffer data that is eventually consumed
// by shaders. They are package local and exposed/accessed through Mesh.

import (
	"log"
)

// vertexData contains per vertex information. A vertex buffer can hold
// one of float32 or byte data, but not both.
type vertexData struct {
	floats    []float32 // Vertex buffer arranged as [][span]float32
	bytes     []byte    // Vertex buffer arranged as [][span]byte
	span      int32     // Elements per vertex
	dtype     uint32    // Data type is one of floatData or byteData.
	ref       uint32    // Vertex GPU buffer reference.
	lloc      uint32    // Shader layout location.
	usage     uint32    // STATIC_DRAW, DYNAMIC_DRAW.
	normalize bool      // Normalize to 0-1 range. Default false.
	vcnt      int       // Number of verticies covered by this data.
	rebind    bool      // Updated data needs to be sent to GPU.
}

// Tag the type data type. Used as values for vertexData.dtype.
const (
	floatData = iota // most vertex data are floats.
	byteData         // byte vertex data are animated bone indexes and weights.
)

// newVertexData is used within the package to get the object rather than
// the interface.
func newVertexData(lloc uint32, span int32) *vertexData {
	vd := &vertexData{}
	vd.floats = []float32{}
	vd.bytes = []byte{}
	vd.span = span
	vd.lloc = lloc
	vd.usage = STATIC
	return vd
}

// set makes a copy of the given data, replacing any existing data, and marks
// the data as needed to be resent to the GPU.
func (vd *vertexData) set(data interface{}) {
	vd.vcnt = 0
	switch d := data.(type) {
	case []float32:
		vd.floats = vd.floats[:0]           // keep allocated memory.
		vd.floats = append(vd.floats, d...) // copy in new data.
		vd.dtype = floatData
		vd.vcnt = len(vd.floats) / int(vd.span)
		vd.rebind = true
	case []byte:
		vd.bytes = vd.bytes[:0]           // keep allocated memory.
		vd.bytes = append(vd.bytes, d...) // copy in new data.
		vd.dtype = byteData
		vd.vcnt = len(vd.bytes) / int(vd.span)
		vd.rebind = true // Set to false when rebound.
	default:
		log.Printf("vData.SetData: dev error : invalid data type")
	}
}

// setUsage tells the graphics card how the data is expected to be read, updated.
func (vd *vertexData) setUsage(usage uint32) {
	switch usage {
	case STATIC, DYNAMIC:
		vd.usage = usage
	}
}

// size returns the current buffer data size in bytes.
func (vd *vertexData) size() uint32 {
	if len(vd.floats) > 0 {
		return uint32(len(vd.floats)) * 4
	}
	if len(vd.bytes) > 0 {
		return uint32(len(vd.bytes))
	}
	return 0
}

// vertexData
// ============================================================================
// faceData

// faceData contains the vertex draw order.
type faceData struct {
	data   []uint16 // Vertex buffer arranged as [][span]uint16.
	ref    uint32   // Vertex GPU buffer reference.
	usage  uint32   // STATIC_DRAW, DYNAMIC_DRAW.
	rebind bool     // Updated data needs to be sent to GPU.
}

// newFaceData is used within the package to get the object rather than
// the interface.
func newFaceData() *faceData {
	fd := &faceData{}
	fd.data = []uint16{}
	fd.usage = STATIC
	return fd
}

// set makes a copy of the given data, replacing any existing data, and marks
// the data as needed to be resent to the GPU.
func (fd *faceData) set(data []uint16) {
	fd.data = fd.data[:0]              // keep allocated memory.
	fd.data = append(fd.data, data...) // copy in new data.
	fd.rebind = true                   // Set to false when rebound.
}

// setUsage tells the graphics card how the data is expected to be read, updated.
func (fd *faceData) setUsage(usage uint32) {
	switch usage {
	case STATIC, DYNAMIC:
		fd.usage = usage
	}
}

// size returns the current buffer data size in bytes.
func (fd *faceData) size() uint32 { return uint32(len(fd.data)) * 2 }
