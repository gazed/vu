// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package render

import (
	"log"
)

// Mesh holds 3D model data in a format that is easily consumed by a rendering
// layer. The data consists of one or more sets of per-vertex data points and
// how the vertex positions are organized into shapes like triangles or lines.
// A Mesh is expected to be referenced by multiple parts and thus does not
// contain any instance information like location or scale. A mesh is most
// often created by the asset pipeline from disk based files that were in turn
// created by tools like Blender.
//
// Note each data buffer must refer to the same number of verticies where the
// number of verticies in one mesh must be less than 65,000. Also note:
//    • Layout location (lloc) 0 is always used for vertex positions.
//    • Per-vertex data currently supports float32 or byte data slices.
//    • There is a maximum of 16 layout locations for one shader.
//    • The data memory is reused when buffer data is updated.
type Mesh interface {
	Name() string // Unique identifier set on creation.
	Size() uint32 // Total bytes used by all buffers.
	Bound() bool  // True if the mesh has a GPU reference.

	// Add and describe sets of per-vertex data indexed by layout location.
	// Layout location (lloc) links to a shader input. The span indicates the
	// number of data points per vertex. InitData must be called once for a
	// given lloc. SetData may be called multiple times to reset data.
	// Some common vertex data buffers and conventions are:
	//    Vertex positions lloc=0 span=3 floats per vertex.
	//    Vertex normals   lloc=1 span=3 floats per vertex.
	//    UV tex coords    lloc=2 span=2 floats per vertex.
	//    Colours          lloc=3 span=4 floats per vertex.
	InitData(lloc, span, usage uint32, normalize bool) Mesh
	SetData(lloc uint32, data interface{}) // Only works on initialized data.

	// Faces are vertex position indicies that form shapes. This data can be used
	// to describe triangles or lines, or ignored for point verticies.
	InitFaces(usage uint32) Mesh // Defaults to STATIC_DRAW
	SetFaces(data []uint16)      // Indicies to vertex positions.
}

// =============================================================================

// mesh is the default implementation of Mesh. It stores the 3D mesh data and
// is expected to be shared among multiple models.
type mesh struct {
	name string // Unique mesh name.
	vao  uint32 // GPU reference for the mesh and all buffers.
	refs uint32 // Number of Model references to this mesh.
	numv int32  // Number of mesh verticies.
	size uint32 // Total bytes used by all data buffers.

	// Per-vertex and vertex index data.
	rebind bool                   // Updated data needs to be sent to GPU.
	faces  *faceData              // Indicies describing the triangle faces.
	vdata  map[uint32]*vertexData // Per-vertex data values.
}

// newMesh allocates space for a mesh structure, including space to
// store buffer data.
func newMesh(name string) *mesh {
	m := &mesh{name: name, rebind: true}
	m.vdata = map[uint32]*vertexData{}
	return m
}

// Implement Mesh.
func (m *mesh) Name() string { return m.name }
func (m *mesh) Size() uint32 { return m.size }
func (m *mesh) Bound() bool  { return m.vao != 0 }
func (m *mesh) InitData(lloc, span, usage uint32, normalize bool) Mesh {
	if _, ok := m.vdata[lloc]; !ok {
		vd := newVertexData(lloc, int32(span))
		vd.normalize = normalize
		vd.setUsage(usage)
		m.vdata[lloc] = vd
	}
	return m
}
func (m *mesh) SetData(lloc uint32, data interface{}) {
	if _, ok := m.vdata[lloc]; ok {
		m.vdata[lloc].set(data)
		m.numv = m.numVerticies() // update when data changes.
		m.size = m.numBytes()     // ditto.
		m.rebind = true
	}
}
func (m *mesh) InitFaces(usage uint32) Mesh {
	if m.faces == nil {
		m.faces = newFaceData()
		m.faces.setUsage(usage)
	}
	return m
}
func (m *mesh) SetFaces(data []uint16) {
	if m.faces != nil {
		m.faces.set(data)
		m.size = m.numBytes() // update when data changes.
		m.rebind = true
	}
}

// numVerticies returns the number of verticies described by the vertex data.
func (m *mesh) numVerticies() int32 {
	if vd, ok := m.vdata[0]; ok && vd.dtype == floatData {
		return int32(len(vd.floats)) / vd.span
	}
	return 0
}

// numBytes returns the total number of data bytes for all data buffers.
func (m *mesh) numBytes() (size uint32) {
	for _, vd := range m.vdata {
		size += vd.size()
	}
	if m.faces != nil {
		size += m.faces.size()
	}
	return
}

// valid returns true if there is a data for layout location 0 and all
// the other vertex buffers describe an equal number of verticies.
func (m *mesh) valid() bool {
	if len(m.vdata) <= 0 { // must have some data.
		return false
	}

	// Vertex positions must exist and must be floats.
	v0 := m.vdata[0]
	if v0 == nil || v0.dtype != floatData || len(v0.floats) <= 0 {
		return false
	}

	// Each vertex buffer must describe the same number of verticies.
	for _, vd := range m.vdata {
		if v0.vcnt != vd.vcnt {
			return false
		}
	}
	return true
}

// hasLocation returns true if there is vertex data for the given
// layout location.
func (m *mesh) hasLocation(lloc uint32) bool {
	_, ok := m.vdata[lloc]
	return ok
}

// mesh
// =============================================================================
// vertexData and faceData abstract away some of the render data details.
// They correspond to render layer buffer data that is eventually consumed
// by shaders. They are package local and exposed/accessed through Mesh.

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
