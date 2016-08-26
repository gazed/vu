// Copyright © 2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package render

import (
	"sort"

	"github.com/gazed/vu/math/lin"
)

// Draw contains the information necessary for a single draw call.
// A draw call needs references for vertex buffers, textures, and a shader
// that have been bound with render.Bind*(). Only the data needed by
// the shader for this draw call has to be set.
//
// Keep in mind that one draw call is often just one object in a scene
// and an entire scene needs to be redrawn 50 or 60 times a second.
// It is thus recommended to reuse allocated Draw instances where feasible
// by resetting the information as necessary before passing to Render.
type Draw interface {
	SetMv(mv *lin.M4)            // Model-View transform.
	SetMvp(mvp *lin.M4)          // Model-View-Projection transform.
	SetPm(pm *lin.M4)            // Projection matrix only.
	SetDbm(dbm *lin.M4)          // Depth bias matrix for shadow maps.
	SetScale(sx, sy, sz float64) // Scaling, per axis.
	SetPose(pose []lin.M4)       // Animation joint/bone transforms.

	// SetHints affects how a draw is rendered.
	//   bucket : Sort order OPAQUE, TRANSPARENT, OVERLAY.
	//   toCam : Distance to Camera.
	//   depth : True to render with depth.
	//   asTex : True to render to texture.
	SetHints(bucket int, tocam float64, depth bool, fbo uint32)

	// SetCounts for bound references
	//   faces  : Number of triangles to be rendered.
	//   verts  : Number of verticies to be rendered.
	SetCounts(faces, verts int)

	// SetRefs of bound references
	//   shader : Program reference
	//   vao    : Vao ref for the mesh vertex buffers.
	//   mode   : POINTS, LINES, TRIANGLES
	SetRefs(shader, vao uint32, mode int)
	Vao() uint32 // The mesh vertex buffers reference.

	// SetTex assigns bound texture information for this draw.
	//   count  : Total number of textures for this draw.
	//   index  : Texture index starting from 0.
	//   tid    : Bound texture reference.
	//   fn, f0 : Used for multiple textures on one mesh.
	SetTex(count, index int, tid, fn, f0 uint32)
	SetShadowmap(tid uint32) // Shadow depth map texture id.

	// Shader uniform data. String keys match the variables expected
	// by the shader source. Each shader variable is expected to have
	// corresponding values in SetFloats.
	SetUniforms(u map[string]int32)          // Variable names:references.
	SetFloats(key string, floats ...float32) // Set variable data.
	Floats(key string) (vals []float32)      // Get variable data.
	SetAlpha(a float64)                      // Transparency.
	SetTime(t float64)                       // Time in seconds.

	// Allow the application to set an object tag. Used for fallback
	// object render sorting where lower values are rendered first.
	Tag() (t uint64) // Retrieve the tag set for this draw.
	SetTag(t uint64) // Set a tag, like a unique model/entity identifier.
	Bucket() int     // Get draw bucket. Lower buckets are rendered first.
}

// NewDraw allocates data needed for a single draw call.
func NewDraw() Draw {
	d := &draw{}
	d.mv = &m4{}
	d.pm = &m4{}
	d.mvp = &m4{}
	d.dbm = &m4{}
	d.scale = &v3{1, 1, 1}
	d.floats = map[string][]float32{} // Float uniform values.
	return d
}

// Draw.
// =============================================================================
// draw implements Draw.

// draw holds the GPU references, shader uinform data, and the
// model-view-projection transforms needed for a single draw call.
type draw struct {
	shader   uint32 // Bind reference for all shader programs.
	vao      uint32 // Bind reference for all vertex buffers.
	mode     int    // POINTS, LINES, TRIANGLES
	numFaces int32  // Number of triangles to be rendered.
	numVerts int32  // Number of verticies to be rendered.
	shtex    uint32 // GPU bound texture shadow depth map.
	texs     []tex  // GPU bound texture references.

	// Rendering hints.
	bucket int     // Render order hint.
	tocam  float64 // Distance to Camera.
	depth  bool    // True to render with depth.
	fbo    uint32  // Framebuffer id. 0 for default.

	// Shader uniform data.
	uniforms map[string]int32     // Expected uniforms and shader references.
	floats   map[string][]float32 // Uniform values.
	alpha    float32              // Shaders alpha value.
	time     float32              // For shaders that need elapsed time.

	// Transform data.
	mv    *m4    // Model View.
	mvp   *m4    // Model View projection.
	pm    *m4    // Projection only.
	dbm   *m4    // Depth bias matrix for shadow maps.
	scale *v3    // Scale X, Y, Z
	pose  []m34  // Per render frame of animation bone data.
	tag   uint64 // Tag for application debugging.
}

// Set transform information.
func (d *draw) SetMv(mv *lin.M4)   { d.mv.tom4(mv) }
func (d *draw) SetMvp(mvp *lin.M4) { d.mvp.tom4(mvp) }
func (d *draw) SetPm(pm *lin.M4)   { d.pm.tom4(pm) }
func (d *draw) SetDbm(dbm *lin.M4) { d.dbm.tom4(dbm) }
func (d *draw) SetScale(sx, sy, sz float64) {
	d.scale.x, d.scale.y, d.scale.z = float32(sx), float32(sy), float32(sz)
}

// Try to reuse allocated model animation data where possible.
func (d *draw) SetPose(pose []lin.M4) {
	if pose == nil {
		d.pose = d.pose[:0]
		return
	}
	need := len(pose)
	if cap(d.pose) < need {
		d.pose = make([]m34, need)
	}
	d.pose = d.pose[0:need]
	for cnt, m4 := range pose {
		(&d.pose[cnt]).tom34(&m4)
	}
}

// SetHints tags this render data with some rendering attributes.
//   bucket: Render order, smallest first.
//   toCam : Distance of object to camera.
//   depth : True to use Z-buffer.
func (d *draw) SetHints(bucket int, toCam float64, depth bool, fbo uint32) {
	d.bucket, d.tocam, d.depth, d.fbo = bucket, toCam, depth, fbo
}

// SetRefs
//   shader: Compiled, linked shader program reference.
//   meshes: Vao buffer reference.
//   mode  : Render mode, triangles or points.
func (d *draw) SetRefs(shader, meshes uint32, mode int) {
	d.shader, d.vao, d.mode = shader, meshes, mode
}

// Vao returns the mesh vao set with SetRefs.
func (d *draw) Vao() uint32 { return d.vao }

// SetCounts specifies how many verticies and how many triangle
// faces for this draw object. This must match the vertex and
// face data.
func (d *draw) SetCounts(faces, verts int) {
	d.numFaces = int32(faces)
	d.numVerts = int32(verts)
}

// SetTex texture references, reusing allocated memory.
func (d *draw) SetTex(count, index int, tid, f0, fn uint32) {
	if count == 0 {
		d.texs = d.texs[:0]
		return
	}
	if cap(d.texs) < count {
		d.texs = make([]tex, count)
	}
	d.texs = d.texs[:count] // ensure length is the same.
	d.texs[index].tid = tid
	d.texs[index].f0 = int32(f0)
	d.texs[index].fn = int32(fn)
}
func (d *draw) SetShadowmap(tid uint32) { d.shtex = tid }

// Set values for the shader uniforms.
func (d *draw) SetUniforms(u map[string]int32) { d.uniforms = u }
func (d *draw) SetAlpha(a float64)             { d.alpha = float32(a) }
func (d *draw) SetTime(t float64)              { d.time = float32(t) }
func (d *draw) SetFloats(key string, floats ...float32) {
	if _, ok := d.floats[key]; ok {
		d.floats[key] = d.floats[key][:0] // reset keeping memory.
	}
	d.floats[key] = append(d.floats[key], floats...)
}
func (d *draw) Floats(key string) (vals []float32) { return d.floats[key] }

// Set an application identifier. Currently used as a last resort
// when sorting draw objects into buckets.
func (d *draw) SetTag(tag uint64) { d.tag = tag }
func (d *draw) Tag() uint64       { return d.tag }
func (d *draw) Bucket() int       { return d.bucket }

// draw
// ===========================================================================
// draws is used to sort a slice of Draw.
type draws []Draw

// Sort parts ordered by bucket first, and distance next.
func (d draws) Len() int      { return len(d) }
func (d draws) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d draws) Less(i, j int) bool {
	di, dj := d[i].(*draw), d[j].(*draw)
	if di.bucket != dj.bucket {
		return di.bucket < dj.bucket // First sort into buckets.
	}
	if di.bucket == Transparent {
		if !lin.Aeq(di.tocam, dj.tocam) {
			return di.tocam > dj.tocam // Sort transparent by distance to camera.
		}
	}
	return di.tag < dj.tag // Sort by eid.
}

// SortDraws sorts draw requests by buckets then by
// distance to camera, and finally by object creation order
// with earlier objects rendered before later objects.
func SortDraws(frame []Draw) { sort.Sort(draws(frame)) }

// =============================================================================

// tex is used to hold a texture reference that is intended
// for only a portion of a model.
type tex struct {
	tid uint32 // GPU bound texture reference.

	// Only set when multiple textures apply to the same model.
	f0, fn int32 // Model face indicies; start and count.
}
