// Copyright © 2016-2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package render

import (
	"sort"

	"github.com/gazed/vu/math/lin"
)

// Draw holds the GPU references, shader uniforms, and the
// model-view-projection transforms needed for a single draw call.
// The GPU references are expected to have been obtained using
// render.Bind*() methods.
//
// Only the data expected by the shader for this draw call has to be set.
// Draw provides setters to transform model data to the form needed
// for rendering. Often this means truncating float64 to float32.
//
// Keep in mind that one draw call is often just one object in a scene
// and an entire scene needs to be redrawn 50 or 60 times a second.
// It is thus recommended to reuse allocated Draw instances where feasible
// by resetting the information as necessary before passing to Render.
type Draw struct {
	Shader uint32 // Bind reference for all shader programs.
	Vao    uint32 // Bind reference for all vertex buffers.
	Mode   int    // Points, Lines, Triangles
	Texs   []tex  // GPU bound texture references.
	Shtex  uint32 // GPU bound texture shadow depth map.
	Tag    uint32 // Application tag for debugging. Commonly an Entity id.

	// Rendering hints.
	Bucket    uint64 // Used to sort draws. Lower buckets rendered first.
	Depth     bool   // True to render with depth.
	Scissor   bool   // True to render with within scissor dimensions.
	Sx, Sy    int32  // Start of scissor area. First 2 parameters.
	Sw, Sh    int32  // Scissor width and height. Last 2 parameters.
	Fbo       uint32 // Framebuffer id. 0 for default.
	FaceCnt   int32  // Number of triangles to be rendered.
	VertCnt   int32  // Number of verticies to be rendered.
	Instances int32  // Postive instance count for instanced mesh.

	// Shader uniform data. The Uniforms are queried from the shader.
	// The UniformData is set by the Engine and App.
	Uniforms    map[string]int32    // Shader uniform variables, references.
	UniformData map[int32][]float32 // Uniform values indexed by ref.

	// Animation pose matrices are converted to a single slice.
	NumPoses int       // Number of animation bone data matricies.
	Poses    []float32 // All the bone data matricies.
}

// NewDraw allocates data needed for a single draw call.
// Scale is initialized to all 1's. Everything else is default.
func NewDraw() *Draw {
	d := &Draw{}
	d.UniformData = map[int32][]float32{} // Float uniform values.
	d.Poses = []float32{}
	return d
}

// Reset clears old draw data so the draw call can be reused.
func (d *Draw) Reset() {
	d.SetPoses(nil)      // Clear animation.
	d.SetTex(0, 0, 0, 0) // Clear texture info.
	d.Scissor = false    // Clear scissor
	d.Instances = 0      // Clear instance data.
}

// SetUniformData for the shader.
func (d *Draw) SetUniformData(ref int32, floats ...float32) {
	if _, ok := d.UniformData[ref]; ok {
		d.UniformData[ref] = d.UniformData[ref][:0] // reset keeping memory.
	}
	d.UniformData[ref] = append(d.UniformData[ref], floats...)
}

// SetM4Data sets the uniform matrix data given the shader reference.
func (d *Draw) SetM4Data(ref int32, m *lin.M4) {
	d.UniformData[ref] = M4ToData(m, d.UniformData[ref])
}

// SetPoses sets the Animation joint/bone transforms.
// It tries to reuse allocated model animation data where possible.
func (d *Draw) SetPoses(poses []lin.M4) {
	d.Poses = d.Poses[:0] // Reuse allocated memory later.
	if poses == nil {
		return
	}
	d.NumPoses = len(poses)
	for _, m4 := range poses {
		d.Poses = appendPoses(&m4, d.Poses) // Append pose data.
	}
}

// SetCounts specifies how many verticies and how many triangle
// faces for this draw object. This must match the vertex and
// face data.
//   faces  : Number of triangles to be rendered.
//   verts  : Number of verticies to be rendered.
func (d *Draw) SetCounts(faces, verts int) {
	d.FaceCnt = int32(faces)
	d.VertCnt = int32(verts)
}

// SetRefs of bound references
//   shader : Compiled, linked shader Program reference.
//   vao    : Vao ref for the mesh vertex buffers.
//   mode   : Render mode: Points, Lines, Triangles
func (d *Draw) SetRefs(shader, meshes uint32, mode int) {
	d.Shader, d.Vao, d.Mode = shader, meshes, mode
}

// SetTex assigns bound texture information for this draw
// reusing previously allocated memory where possible.
//   count  : Total number of textures for this draw.
//   index  : Texture index starting from 0.
//   tid    : Bound texture reference.
func (d *Draw) SetTex(count, index, order int, tid uint32) {
	if count == 0 {
		d.Texs = d.Texs[:0]
		return
	}
	if cap(d.Texs) < count {
		d.Texs = make([]tex, count)
	}
	d.Texs = d.Texs[:count] // ensure length is the same.
	d.Texs[index].tid = tid
	d.Texs[index].order = order
}

// SetShadowmap sets the texture id of the shadow map.
func (d *Draw) SetShadowmap(tid uint32) { d.Shtex = tid }

// Draw
// ===========================================================================
// draws is used to sort a slice of Draw.

type draws []*Draw

// Sort draws where the Bucket must be set to ensure proper draw order.
func (d draws) Len() int      { return len(d) }
func (d draws) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d draws) Less(i, j int) bool {
	if d[i].Bucket == d[j].Bucket {
		// Break ties by prefering earlier entities.
		// Prevents z-battle screen flickering.
		return d[i].Tag < d[j].Tag
	}
	return d[i].Bucket > d[j].Bucket
}

// SortDraws sorts draw requests by buckets then by
// distance to camera, and finally by object creation order
// with earlier objects rendered before later objects.
func SortDraws(frame []*Draw) { sort.Sort(draws(frame)) }

// =============================================================================

// tex is used to hold a texture reference that is intended
// for only a portion of a model.
type tex struct {
	tid   uint32 // GPU bound texture reference.
	order int    // Shader uniform identifier for multiple textures.
}

// =============================================================================

// M4ToData transforms a lin.M4 matrix to an array of floats needed by
// the rendering system. The passed in slice is set to 0 len before
// appending the matrix data in the order needed by the shader.
func M4ToData(m *lin.M4, d []float32) []float32 {
	d = d[:0]
	d = append(d, float32(m.Xx), float32(m.Xy), float32(m.Xz), float32(m.Xw))
	d = append(d, float32(m.Yx), float32(m.Yy), float32(m.Yz), float32(m.Yw))
	d = append(d, float32(m.Zx), float32(m.Zy), float32(m.Zz), float32(m.Zw))
	d = append(d, float32(m.Wx), float32(m.Wy), float32(m.Wz), float32(m.Ww))
	return d
}

// appendPoses appends a 4x4 matrix as a 3x4 float32 column-major matrix.
// The matrix data is interpreted as row-major when sent to the GPU.
// The last column of the 4x4 matrix is not sent as an optimization.
// The shader is expected be aware of this space saving layout.
//    xx, yx, zx, wx float32 // indices 0, 1, 2, 3  [00, 01, 02, 03]
//    xy, yy, zy, wy float32 // indices 4, 5, 6, 7  [10, 11, 12, 13]
//    xz, yz, zz, wz float32 // indices 8, 9, a, b  [20, 21, 22, 23]
//                           //         0, 0, 0, 1 implicit last row.
func appendPoses(m *lin.M4, d []float32) []float32 {
	d = append(d, float32(m.Xx), float32(m.Yx), float32(m.Zx), float32(m.Wx))
	d = append(d, float32(m.Xy), float32(m.Yy), float32(m.Zy), float32(m.Wy))
	d = append(d, float32(m.Xz), float32(m.Yz), float32(m.Zz), float32(m.Wz))
	return d
}
