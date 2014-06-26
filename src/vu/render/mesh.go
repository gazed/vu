// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package render

// Mesh holds 3D model data in a format that is easily consumed by a rendering
// layer. The data consists of one or more sets of per-vertex data points and
// how the vertex positions are organized into shapes like triangles or lines.
// A Mesh is expected to be referenced by multiple parts and thus does not
// contain any instance information like location or scale.
//
// Note each data buffer must refer to the same number of verticies where the
// number of verticies in one mesh must be less than 65,000. Also note:
//    • Layout location (lloc) 0 is always used for vertex positions.
//    • Per-vertex data currently supports data slices of float32s or bytes.
//    • There is a maximum of 16 layout locations for one shader.
//    • The data memory is reused when buffer data is updated.
type Mesh interface {
	Name() string // Unique identifier set on creation.
	Size() uint32 // Total bytes used by all buffers.
	Bound() bool  // True if the mesh has a GPU reference.

	// Add and describe sets of per-vertex data indexed by layout location.
	// Layout location lloc links to a shader input. The span indicates the
	// number of data points per vertex. InitData must be called once for a
	// given lloc. SetData may be called multiple times to (re)set, update
	// data. Some common vertex data buffers and conventions are:
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
	m := &mesh{name: name}
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
