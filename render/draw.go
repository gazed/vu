// Copyright © 2016 Galvanized Logic Inc.
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
	Tag    uint64 // Application tag for debugging. Commonly an Entity id.

	// Shader uniform data.
	Uniforms map[string]int32     // Expected uniforms and shader references.
	Floats   map[string][]float32 // Uniform values.

	// Rendering hints.
	Bucket  int     // Used to sort draws. Lower buckets rendered first.
	Tocam   float64 // Distance to Camera for sorting by distance.
	Depth   bool    // True to render with depth.
	Fbo     uint32  // Framebuffer id. 0 for default.
	FaceCnt int32   // Number of triangles to be rendered.
	VertCnt int32   // Number of verticies to be rendered.

	// Transform data.
	Mv   *m4   // Model View.
	Mvp  *m4   // Model View projection.
	Pm   *m4   // Projection only.
	Dbm  *m4   // Depth bias matrix for shadow maps.
	Pose []m34 // Per render frame of animation bone data.
}

// NewDraw allocates data needed for a single draw call.
// Scale is initialized to all 1's. Everything else is default.
func NewDraw() *Draw {
	d := &Draw{}
	d.Mv = &m4{}
	d.Pm = &m4{}
	d.Mvp = &m4{}
	d.Dbm = &m4{}
	d.Floats = map[string][]float32{} // Float uniform values.
	return d
}

// SetMv sets the Model-View transform.
func (d *Draw) SetMv(mv *lin.M4) { d.Mv.tom4(mv) }

// SetMvp sets the Model-View-Projection transform.
func (d *Draw) SetMvp(mvp *lin.M4) { d.Mvp.tom4(mvp) }

// SetPm sets the Projection matrix.
func (d *Draw) SetPm(pm *lin.M4) { d.Pm.tom4(pm) }

// SetDbm sets the depth bias matrix for shadow maps.
func (d *Draw) SetDbm(dbm *lin.M4) { d.Dbm.tom4(dbm) }

// SetScale sets the scaling factors per axis.
func (d *Draw) SetScale(sx, sy, sz float64) {
	d.SetFloats("scale", float32(sx), float32(sy), float32(sz))
}

// SetPose sets the Animation joint/bone transforms.
// It tries to reuse allocated model animation data where possible.
func (d *Draw) SetPose(pose []lin.M4) {
	if pose == nil {
		d.Pose = d.Pose[:0] // Keep the data, but hide it.
		return
	}
	need := len(pose)
	if cap(d.Pose) < need {
		d.Pose = make([]m34, need) // Resize if necessary.
	}
	d.Pose = d.Pose[0:need]
	for cnt, m4 := range pose {
		(&d.Pose[cnt]).tom34(&m4) // Copy data in.
	}
}

// SetHints affects how a draw is rendered.
//   bucket : Sort order Opaque, Transparent, Overlay.
//   toCam : Distance to Camera.
//   depth : True to render with depth, ie: use Z-buffer.
//   fbo   : Frame buffer object for render to texture.
func (d *Draw) SetHints(bucket int, toCam float64, depth bool, fbo uint32) {
	d.Bucket, d.Tocam, d.Depth, d.Fbo = bucket, toCam, depth, fbo
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
//   fn, f0 : Used for multiple textures on one mesh.
func (d *Draw) SetTex(count, index int, tid, f0, fn uint32) {
	if count == 0 {
		d.Texs = d.Texs[:0]
		return
	}
	if cap(d.Texs) < count {
		d.Texs = make([]tex, count)
	}
	d.Texs = d.Texs[:count] // ensure length is the same.
	d.Texs[index].tid = tid
	d.Texs[index].f0 = int32(f0)
	d.Texs[index].fn = int32(fn)
}

// SetShadowmap sets the texture id of the shadow map.
func (d *Draw) SetShadowmap(tid uint32) { d.Shtex = tid }

// SetUniforms for the shader. String keys match the variables expected
// by the shader source. Each shader variable is expected to have
// corresponding values in SetFloats.
func (d *Draw) SetUniforms(u map[string]int32) { d.Uniforms = u }

// SetFloats sets the named shader uniform variable data.
func (d *Draw) SetFloats(key string, floats ...float32) {
	if _, ok := d.Floats[key]; ok {
		d.Floats[key] = d.Floats[key][:0] // reset keeping memory.
	}
	d.Floats[key] = append(d.Floats[key], floats...)
}

// Draw
// ===========================================================================
// draws is used to sort a slice of Draw.

type draws []*Draw

// Sort parts ordered by bucket first, and distance next.
func (d draws) Len() int      { return len(d) }
func (d draws) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d draws) Less(i, j int) bool {
	di, dj := d[i], d[j]
	if di.Bucket != dj.Bucket {
		return di.Bucket < dj.Bucket // First sort into buckets.
	}
	if di.Bucket == Transparent {
		if !lin.Aeq(di.Tocam, dj.Tocam) {
			return di.Tocam > dj.Tocam // Sort transparent by distance to camera.
		}
	}
	return di.Tag < dj.Tag // Sort by eid.
}

// SortDraws sorts draw requests by buckets then by
// distance to camera, and finally by object creation order
// with earlier objects rendered before later objects.
func SortDraws(frame []*Draw) { sort.Sort(draws(frame)) }

// =============================================================================

// tex is used to hold a texture reference that is intended
// for only a portion of a model.
type tex struct {
	tid uint32 // GPU bound texture reference.

	// Only set when multiple textures apply to the same model.
	f0, fn int32 // Model face indicies; start and count.
}
