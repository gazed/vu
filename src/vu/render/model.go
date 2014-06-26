// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package render

import (
	"fmt"
	"image"
	"log"
	"math"
	"time"
	"vu/math/lin"
)

// Model supplys a shader with data. Model is initialized with a shader and
// provides methods for setting the data expected by the shader. Often data
// consists of a Mesh, one or more Textures, and uniform values. Animated
// models expect joint and key frame data. Model also accepts render
// directives affecting the overall rendering process.
type Model interface {
	Shader() Shader       // One shader must be set on creation.
	SetDrawMode(mode int) // Render directive: TRIANGLES, POINTS, or LINES.
	Set2D()               // Render directive: Turns DEPTH off for this model.

	// Shader uniforms are set using uniform specific methods and through
	// generic SetUniform which takes a uniform name and 1-4 float32 values.
	SetScale(x, y, z float64)            // Model sizing.
	SetMvTransform(mv *lin.M4)           // Model-View transform.
	SetMvpTransform(mvp *lin.M4)         // Model-View-Projection transform.
	SetAlpha(a float64)                  // Set or get the shaders
	Alpha() (a float64)                  // ...alpha uniform value.
	SetUniform(id string, val []float32) // Set or get float32 based
	Uniform(id string) (val []float32)   // ...shader uniform value.

	// Mesh data can be set from a mesh resource using SetMesh, or a set
	// from generated data on a Mesh created from Mesh().
	Name() string            // Model name is the Mesh name, "" if no mesh.
	Mesh() Mesh              // Get existing mesh or lazy creates new Mesh.
	SetMesh(mesh Mesh) Model // Set to given mesh resource.

	// A model may have 0 to 15 textures to match the shader expectations.
	// The index is the same as the texture unit.
	Textures() []Texture                 // Textures can be multiple per
	Texture(index int) Texture           // ...model and are indexed
	AddTexture(t Texture) (index int)    // ...when adding, or
	UseTexture(t Texture, index int)     // ...replacing, or
	RemTexture(index int)                // ...removing, or
	TexMode(index int, mode int)         // ...how they're drawn.
	SetImage(img image.Image, index int) // Directly set texture data.

	// AddModelTexture allows multiple textures to be used on a single mesh.
	// The texture affects faces starting at index f0 and continuing for fn.
	AddModelTexture(t Texture, f0, fn uint32) (index int)

	// Animation data. Each frame contains 1 transform matrix for each joint,
	// and joints contains the indexed parent hierarchy.
	SetAnimation(frames []*lin.M4, joints []int32, numFrames int)
	Animate(dt float64) // Called regularly to interpolate between frames.

	// Verify the availability of the data expected by the shader.
	Verify() error // Return a nil error if expected data matches available.
	Dispose()      // Release all rendering resources.
}

const (
	// Draw mode types for vertex data rendering. One of these is expected
	// when calling Model.SetDrawMode(mode)
	TRIANGLES = iota // Triangles are the default for 3D models.
	POINTS           // Points are used for particle effects.
	LINES            // Lines are used for drawing wireframe shapes.

	// Texture rendering modes. Default is CLAMP
	REPEAT // Textures repeat with UV values greater than 1.
)

// ============================================================================

// model implments Model. It uses render specific knowledge while conforming
// to the generic Model interface. It holds and provides the data needed by
// the shaders.
type model struct {
	gc   graphicsContext // Graphics context injected on creation.
	shd  *shader         // Pipeline renderer for this model.
	msh  *mesh           // Vertex buffer data.
	tex  []*texture      // Texture data. Needed for a uv texture buffer.
	mode int             // How to draw the vertex data.
	is2D bool            // Whether or not to enable depth testing.
	tmap []textureMap    // data to map multiple textures to one vertex mesh.

	// Animation data.
	nFrames int      // number of animation frames.
	frames  []lin.M4 // nFrames*nPoses transform bone positions.
	anim    []m34    // nPoses of a frame of animation data sent to GPU.
	fcnt    float64  // frame counter.
	jparent []int32  // joint parent indicies.

	// Predefined shader uniform values.
	mv    *m4       // Model view.
	mvp   *m4       // Model view projection.
	nm    *m3       // Normal matrix
	scale *v3       // Model scaling.
	alpha float32   // Shaders alpha value.
	time  time.Time // For shaders that need elapsed time.

	// Applicaiton defined shader uniform values.
	uniforms map[string][]float32                 // Render pre-defined.
	common   map[string]func(m *model, ref int32) // Model defined.

	jntM4, tmpM4 *lin.M4 // Per-frame scratch values for animations.
}

type textureMap struct {
	f0, fn int32 // First face index and number of faces.
}

// newModel creates a new model. It needs to be loaded with data.
func newModel(gc Renderer, s Shader) Model {
	m := &model{}
	m.gc = gc.(graphicsContext)
	m.mv = &m4{}
	m.mvp = &m4{}
	m.nm = &m3{}
	m.scale = &v3{}
	m.tex = []*texture{}
	m.anim = []m34{}
	m.uniforms = map[string][]float32{}
	m.time = time.Now()
	m.alpha = 1
	m.setShader(s)
	m.jntM4, m.tmpM4 = &lin.M4{}, &lin.M4{}

	// Provide some common shader uniforms that are needed by most shaders.
	m.common = map[string]func(m *model, ref int32){

		// transform matricies.
		"mvpm": func(m *model, ref int32) { m.gc.bindUniform(ref, x4, 1, m.mvp.Pointer()) },
		"mvm":  func(m *model, ref int32) { m.gc.bindUniform(ref, x4, 1, m.mv.Pointer()) },
		"nm": func(m *model, ref int32) {
			nm := (&m3{}).m3(m.mv)
			m.gc.bindUniform(ref, x3, 1, nm.Pointer())
		},

		// bone position animation data.
		"bpos": func(m *model, ref int32) {
			if len(m.anim) > 0 {
				m.gc.bindUniform(ref, x34, len(m.anim), m.anim[0].Pointer())
			}
		},

		// textures, texture atlases, and multitextures.
		"uv":   func(m *model, ref int32) { m.gc.useTexture(ref, 0, m.tex[0]) },
		"uv0":  func(m *model, ref int32) { m.gc.useTexture(ref, 0, m.tex[0]) },
		"uv1":  func(m *model, ref int32) { m.gc.useTexture(ref, 1, m.tex[1]) },
		"uv2":  func(m *model, ref int32) { m.gc.useTexture(ref, 2, m.tex[2]) },
		"uv3":  func(m *model, ref int32) { m.gc.useTexture(ref, 3, m.tex[3]) },
		"uv4":  func(m *model, ref int32) { m.gc.useTexture(ref, 4, m.tex[4]) },
		"uv5":  func(m *model, ref int32) { m.gc.useTexture(ref, 5, m.tex[5]) },
		"uv6":  func(m *model, ref int32) { m.gc.useTexture(ref, 6, m.tex[6]) },
		"uv7":  func(m *model, ref int32) { m.gc.useTexture(ref, 7, m.tex[7]) },
		"uv8":  func(m *model, ref int32) { m.gc.useTexture(ref, 8, m.tex[8]) },
		"uv9":  func(m *model, ref int32) { m.gc.useTexture(ref, 9, m.tex[9]) },
		"uv10": func(m *model, ref int32) { m.gc.useTexture(ref, 10, m.tex[10]) },
		"uv11": func(m *model, ref int32) { m.gc.useTexture(ref, 11, m.tex[11]) },
		"uv12": func(m *model, ref int32) { m.gc.useTexture(ref, 12, m.tex[12]) },
		"uv13": func(m *model, ref int32) { m.gc.useTexture(ref, 13, m.tex[13]) },
		"uv14": func(m *model, ref int32) { m.gc.useTexture(ref, 14, m.tex[14]) },
		"uv15": func(m *model, ref int32) { m.gc.useTexture(ref, 15, m.tex[15]) },

		// model size, alpha, and elapsed time.
		"scale": func(m *model, ref int32) { m.gc.bindUniform(ref, f3, 1, m.scale.x, m.scale.y, m.scale.z) },
		"alpha": func(m *model, ref int32) { m.gc.bindUniform(ref, f1, 1, m.alpha) },
		"time":  func(m *model, ref int32) { m.gc.bindUniform(ref, f1, 1, float32(time.Since(m.time).Seconds())) },
	}
	return m
}

// Model implementation.
func (m *model) SetAlpha(a float64)                    { m.alpha = float32(a) }
func (m *model) Alpha() (a float64)                    { return float64(m.alpha) }
func (m *model) Set2D()                                { m.is2D = true }
func (m *model) SetUniform(id string, value []float32) { m.uniforms[id] = value }
func (m *model) Uniform(id string) (value []float32)   { return m.uniforms[id] }
func (m *model) SetMvTransform(mv *lin.M4)             { m.mv.tom4(mv) }
func (m *model) SetMvpTransform(mvp *lin.M4)           { m.mvp.tom4(mvp) }

// Model implementation.
func (m *model) SetScale(x, y, z float64) {
	m.scale.x, m.scale.y, m.scale.z = float32(x), float32(y), float32(z)
}

// Model implementation.
func (m *model) AddTexture(tex Texture) (index int) {
	t := tex.(*texture)
	if !t.Bound() {
		if err := m.gc.bindTexture(t); err == nil {
			t.FreeImg()
		} else {
			log.Printf("model.AddTexture: could not bind %s %s", tex.Name(), err)
		}
	}
	t.refs++
	m.tex = append(m.tex, t)
	return len(m.tex)
}

// Model implementation.
func (m *model) AddModelTexture(tex Texture, f0, fn uint32) (index int) {
	m.AddTexture(tex)
	m.tmap = append(m.tmap, textureMap{int32(f0), int32(fn)})
	return len(m.tex)
}

// Model implementation.
func (m *model) UseTexture(t Texture, index int) {
	if tex := m.tex[index]; tex != nil {
		tex.refs-- // do not delete here.
	}
	newt := t.(*texture)
	if !newt.Bound() {
		if err := m.gc.bindTexture(newt); err == nil {
			newt.FreeImg()
		} else {
			log.Printf("model.AddTexture: could not bind %s %s", newt.Name(), err)
		}
	}
	newt.refs++
	m.tex[index] = newt
}

// Model implementation.
func (m *model) RemTexture(index int) {
	if index < len(m.tex) {
		if t := m.tex[index]; t != nil {
			t.refs -= 1
			if 0 >= t.refs {
				m.gc.deleteTexture(t.tid)
				t.tid = 0
			}
			m.tex[index] = nil
			m.tex = append(m.tex[:index], m.tex[index+1:]...)
		}
	}
}
func (m *model) TexMode(index int, mode int) {
	if index < len(m.tex) {
		if t := m.tex[index]; t != nil {
			if mode == REPEAT { // only one mode at the moment.
				t.SetRepeat(true)
			}
			m.gc.updateTextureMode(t)
		}
	}
}

// Model implementation.
func (m *model) Texture(index int) Texture {
	if index < len(m.tex) {
		return m.tex[index]
	}
	return nil // explicitly return nil for nil interface.
}

// Model implementation.
func (m *model) Textures() []Texture {
	textures := []Texture{}
	for _, t := range m.tex {
		textures = append(textures, t)
	}
	return textures
}

// Model implementation.
func (m *model) SetImage(img image.Image, index int) {
	if index < len(m.tex) {
		tex := m.tex[index]
		tex.Set(img)
		if err := m.gc.bindTexture(tex); err == nil {
			tex.FreeImg()
		} else {
			log.Printf("model.SetImage: could not bind %s %s", tex.Name(), err)
		}
	}
}

// setShader is called once on model creation.
func (m *model) setShader(s Shader) {
	if m.shd = s.(*shader); m.shd != nil {
		if !m.shd.Bound() {
			if err := m.gc.bindShader(s); err != nil {
				log.Printf("model.setShader could not bind %s %s", s.Name(), err)
			}
		}
		m.shd.refs++
	}
}

// Model implementation.
func (m *model) Shader() Shader {
	if m.shd != nil {
		return m.shd
	}
	return nil // explicitly return nil for nil interface.
}

// Model implementation.
func (m *model) Name() string {
	if m.msh != nil {
		return m.msh.Name()
	}
	return ""
}

// Model implementation.
func (m *model) SetMesh(modelMesh Mesh) Model {
	m.disposeMesh()
	m.msh = modelMesh.(*mesh)
	if !m.msh.Bound() {
		if err := m.gc.bindMesh(m.msh); err != nil {
			log.Printf("model.SetMesh could not bind %s %s", m.msh.Name(), err)
		}
	}
	m.msh.refs++
	return m
}

// Model implementation.
func (m *model) Mesh() Mesh {
	if m.msh == nil {
		m.msh = newMesh("mesh") // not cached for reuse.
		m.msh.refs++
	}
	return m.msh
}

// // Model implementation.
// // Only overwrite mesh data labelled mesh (not cached mesh data).
// func (m *model) BindMesh() {
// 	if m.msh.rebind && m.msh.name == "mesh" {
// 		if err := m.gc.bindMesh(m.msh); err != nil {
// 			log.Printf("model.NewMesh failed %s", err)
// 		}
// 		m.msh.rebind = false
// 	}
// }

// Model implementation.
func (m *model) SetDrawMode(mode int) {
	switch mode {
	case TRIANGLES, POINTS, LINES:
		m.mode = mode
	}
}

// Model implementation.
// Disposing a graphics asset, means that it needs to be rebound.
// Any cached instances should be freed.
func (m *model) Dispose() {
	m.disposeShader()
	m.disposeMesh()
	for index, _ := range m.tex {
		m.RemTexture(index)
	}
}

// disposeShader releases the shader associated with this model.
func (m *model) disposeShader() {
	if m.shd != nil {
		m.shd.refs -= 1
		if 0 >= m.shd.refs {
			m.gc.deleteShader(m.shd.program)
			m.shd.program = 0
			m.shd = nil
		}
	}
}

// disposeShader releases the mesh data associated with this model.
func (m *model) disposeMesh() {
	if m.msh != nil {
		m.msh.refs -= 1
		if 0 >= m.msh.refs {
			m.gc.deleteMesh(m.msh.vao)
			m.msh.vao = 0
			m.msh = nil
		}
	}
}

// bindUniforms links model data to the uniforms discovered in the model shader.
func (m *model) bindUniforms() {
	for key, ref := range m.shd.uniforms {
		if bindFunc, ok := m.common[key]; ok {
			bindFunc(m, ref)
		} else if floats, ok := m.uniforms[key]; ok {
			switch len(floats) {
			case 1:
				m.gc.bindUniform(ref, f1, 1, floats[0])
			case 2:
				m.gc.bindUniform(ref, f2, 1, floats[0], floats[1])
			case 3:
				m.gc.bindUniform(ref, f3, 1, floats[0], floats[1], floats[2])
			case 4:
				m.gc.bindUniform(ref, f4, 1, floats[0], floats[1], floats[2], floats[3])
			}
		} else {
			log.Printf("No uniform %s for mesh %s shader %s", key, m.msh.Name(), m.shd.Name())
		}
	}
}

// Model implementation.
func (m *model) SetAnimation(frames []*lin.M4, jparents []int32, numFrames int) {
	m.nFrames = numFrames
	m.anim = make([]m34, len(jparents))    // 1 matrix for each joint.
	m.frames = make([]lin.M4, len(frames)) // transform matrices for each frame.
	for cnt, frame := range frames {
		m.frames[cnt].Set(frame)
	}
	m.jparent = m.jparent[:0]
	m.jparent = append(m.jparent, jparents...)
}

// Note that this animates all attributes (position, normal, tangent, bitangent)
// for expository purposes, even though this demo does not use all of them for rendering.
func (m *model) Animate(dt float64) {
	if m.nFrames <= 0 {
		return
	}

	// The frame timer, fcnt, controls the speed of the animation.
	// FUTURE: find a more generic way to time animations.
	m.fcnt += dt * 10 // Increment frame timer.
	m.nFrames = 100

	frame1 := int(math.Floor(m.fcnt))
	frame2 := frame1 + 1
	frameoffset := float64(m.fcnt) - float64(frame1)
	frame1 %= m.nFrames
	frame2 %= m.nFrames

	// Interpolate matrixes between the two closest frames and concatenate with
	// parent matrix if necessary. Concatenate the result with the inverse of the
	// base pose. You would normally do animation blending and inter-frame
	// blending here in a 3D engine.
	nJoints := len(m.frames) / m.nFrames
	for cnt := 0; cnt < nJoints; cnt++ {

		// interpolate between the two closest frames.
		m1, m2 := &m.frames[frame1*nJoints+cnt], &m.frames[frame2*nJoints+cnt]
		m.jntM4.Set(m1).Scale(1-frameoffset).Add(m.jntM4, m.tmpM4.Set(m2).Scale(frameoffset))

		if m.jparent[cnt] >= 0 {

			// parentPose * childPose * childInverseBasePose
			m.jntM4.Mult(m.jntM4, (&m.anim[m.jparent[cnt]]).toM4(m.tmpM4))
		}
		(&m.anim[cnt]).tom34(m.jntM4)
	}
}

// Model implementation.
func (m *model) Verify() error {
	if m.shd == nil {
		return fmt.Errorf("model.Verify: no shader")
	}

	// Check if the expected uniform is supported with model data.
	for label, _ := range m.shd.uniforms {
		if _, ok := m.common[label]; !ok {
			if _, ok := m.uniforms[label]; !ok {
				return fmt.Errorf("model.Verify: no uniform %s in shader %s", label, m.shd.name)
			}
		}
	}

	// Check if the expected attribute is supported with buffer data.
	if m.msh == nil && len(m.shd.attributes) > 0 {
		return fmt.Errorf("model.Verify: expecting %d buffers for shader %s", len(m.shd.attributes), m.shd.name)
	}
	for label, key := range m.shd.attributes {
		if !m.msh.hasLocation(key) {
			return fmt.Errorf("model.Verify: no buffer for attribute %s in shader %s", label, m.shd.name)
		}
	}
	return nil
}
