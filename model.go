// Copyright © 2015-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"image"
	"log"
	"math"
	"time"

	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/render"
)

// Model manages rendered 3D objects. Model is the link between loaded
// assets and the data needed by the rendering system. Assets, such as
// shaders, mesh, textures, etc. are specified as unique strings in the
// Load* methods. Assets are loaded and converted to intermediate data
// which is later converted to render draws for the render system.
//
// A model is expected to be attached to a point-of-view (Pov) to give
// it a 3D location and orientation in the render transform hierarchy.
//
// It is the applications responsibility to ensure that a models shader
// is provided with its expected asset data. Mismatches generate logs.
type Model interface {
	Shader() (name string)     // Rendered models have a shader.
	LoadMat(name string) Model // Optional surface material color.
	Alpha() (a float64)        // Transparency: 1 for full opaque
	SetAlpha(a float64)        // ...0 for fully transparent.
	Color() (r, g, b float64)  // Colors between 0 and 1.
	SetColor(r, g, b float64)  // ...1 for full color.

	// Mesh handles verticies, per-vertex data, and triangle faces.
	// Meshes can be loaded from assets or created/generated.
	// LoadMesh creates a mesh from loaded mesh resource assets.
	LoadMesh(name string) Model // Expects to load static mesh data.
	// NewMesh creates an empty mesh expecting generated data.
	NewMesh(name string) Model // App is responsible for generating data.
	// Generating mesh data is closely tied to a given shader.
	// InitMesh must be called once, SetMesh data may be called as needed.
	//    lloc     : layout location is the shader input reference.
	//    span     : indicates the number of data points per vertex.
	//    usage    : STATIC or DYNAMIC.
	//    normalize: true to convert data to the 0->1 range.
	// Some vertex shader data conventions are:
	//    Vertex positions lloc=0 span=3_floats_per_vertex.
	//    Vertex normals   lloc=1 span=3_floats_per_vertex.
	//    UV tex coords    lloc=2 span=2_floats_per_vertex.
	//    Color            lloc=3 span=4_floats_per_vertex.
	InitMesh(lloc, span, usage uint32, normalize bool) Model
	SetMeshData(lloc uint32, data interface{}) // Only works after InitMesh
	InitFaces(usage uint32) Model              // Defaults to STATIC_DRAW
	SetFaces(data []uint16)                    // Indicies to vertex positions.
	SetDrawMode(mode int) Model                // TRIANGLES, LINES, POINTS.

	// Models can have one or more textures applied to a single mesh.
	// Textures are initialized from assets and can be updated with images.
	AddTex(name string) Model             // Loads and adds a texture.
	NewTex(name string) Model             // Adds new texture. Needs SetImg.
	SetTex(index int, name string)        // Replace/reload texture.
	SetImg(index int, img image.Image)    // Replace image, nil values ignored.
	TexImg(index int) image.Image         // Get image, nil if invalid index.
	SetTexMode(index int, mode int) Model // TEX_CLAMP, TEX_REPEAT.
	UseLayer(l Layer) Model               // Use render pass texture.

	// Animated models can have multiple animated sequences,
	// ie. "moves", that are indexed from 0. Bones can also
	// be used to position other models, ie: attachment points.
	LoadAnim(anim string) Model            // Sets an animated model.
	Animate(action, frame int) bool        // Return true if available.
	Action() (action, frame, maxFrame int) // Current movement info.
	Actions() []string                     // Animation sequence names.
	Pose(bone int) *lin.M4                 // Bone transform: attach point.

	// Fonts are used to display small text phrases using a mesh plane.
	// Fonts imply a texture shader and a texture for this model.
	LoadFont(font string) Model  // Set the character mapping resource.
	SetPhrase(text string) Model // Set the string to display.
	PhraseWidth() int            // Width in pixels, 0 if not loaded.

	// Particle effects are either CPU:application or GPU:shader based.
	// SetEffect sets a CPU controlled particle effect.
	SetEffect(mover ParticleEffect, maxParticles int) Model
	SetDepth(enabled bool) Model // Effects work better ignoring depth.

	// Set/get shader uniform values where id is the shader uniform name.
	Uniform(id string) (value []float32)         // Uniform name/values.
	SetUniform(id string, floats ...interface{}) // Individual values.

	// Models can optionally cast shadows or reveal shadows. Casting
	// shadows implies a scene with a light. Showing shadows implies
	// a shadow map capable shader.
	CastShadow() Model // Toggle casting a shadow. Default false.
	HasShadows() Model // Toggle showing shadows. Default false.
}

// Model
// =============================================================================
// model implements Model.

// model implements Model. It links the applications transform hierarchy
// to the rendering system. The application specifies the model resources.
// These resources are later linked and bound during engine processing.
type model struct {
	shd      *shader         // Mandatory GPU render program.
	texs     []*texture      // Optional: one or more texture images.
	mat      *material       // Optional: material lighting info.
	msh      *mesh           // Mandatory vertex buffer data.
	drawMode int             // TRIANGLES, POINTS, LINES.
	effect   *particleEffect // Optional particle effect.
	loads    []*loadReq      // Assets waiting to be loaded.
	depth    bool            // Depth buffer on by default.
	layer    *layer          // Optional previous render pass.

	// Optional animated model control information.
	anm     *animation // Optional: bone animation info.
	frame   float64    // Frame counter.
	move    int        // Aurrent animation defaults to 0.
	nFrames int        // Number of frames in the current movement.
	pose    []lin.M4   // Pose refreshed each update.

	// Optional font information.
	fnt         *font  // Optional: font layout data.
	phrase      string // Initial pre-load phrase.
	phraseWidth int    // Rendered phrase width in pixels, 0 otherwise.

	// Optional model shadow information.
	castShadow bool // Model to cast a shadow. Default false.
	hasShadows bool // Model to reveal a shadow. Default false.

	// Shader dependent uniform data.
	time     time.Time            // Time needed by some shaders.
	alpha    float32              // Transparency between 0 and 1.
	kd       rgb                  // Diffuse color.
	ka       rgb                  // Ambient color.
	ks       rgb                  // Specular color.
	uniforms map[string][]float32 // Uniform values.
	sm       *lin.M4              // scratch matrix.
}

// newModel allocates a new model instance setting some common defaults.
func newModel(shaderName string) *model {
	m := &model{alpha: 1, depth: true}
	m.shd = newShader(shaderName)
	m.loads = append(m.loads, &loadReq{data: m, a: newShader(shaderName)})
	m.time = time.Now()
	m.uniforms = map[string][]float32{}
	m.sm = &lin.M4{}
	return m
}

// loaded returns true if all the model parts have data.
func (m *model) loaded() bool {
	if m.shd == nil || !m.shd.loaded { // not optional
		return false
	}
	if m.msh == nil || !m.msh.loaded { // not optional
		return false
	}
	for _, tex := range m.texs { // optional
		if !tex.loaded {
			return false
		}
	}
	if m.fnt != nil && !m.fnt.loaded { // optional
		return false
	}
	if m.mat != nil && !m.mat.loaded { // optional
		return false
	}
	if m.anm != nil && !m.anm.loaded { // optional
		return false
	}
	return true
}

// Each model has one shader.
func (m *model) Shader() string { return m.shd.name }

// Alpha is model transparency. This value overrides any material values.
func (m *model) Alpha() (a float64) { return float64(m.alpha) }
func (m *model) SetAlpha(a float64) {
	m.alpha = float32(a)
}

// Color can override values loaded from a material.
func (m *model) Color() (r, g, b float64) {
	return float64(m.kd.R), float64(m.kd.G), float64(m.kd.B)
}
func (m *model) SetColor(r, g, b float64) {
	m.kd.R, m.kd.G, m.kd.B = float32(r), float32(g), float32(b)
}

// Material is used to help with coloring for shaders that use lights.
// Overrides existing values if it was the last one set.
func (m *model) LoadMat(name string) Model {
	m.mat = newMaterial(name)
	m.loads = append(m.loads, &loadReq{data: m, a: newMaterial(name)})
	return m
}

// Each model has one mesh. The mesh is specified here and
// will be sent for loading and binding later on.
func (m *model) LoadMesh(meshName string) Model {
	if m.msh == nil && m.anm == nil {
		m.msh = newMesh(meshName) // placeholder
		req := &loadReq{data: m, a: newMesh(meshName)}
		m.loads = append(m.loads, req)
	}
	return m
}
func (m *model) NewMesh(meshName string) Model {
	if m.msh == nil && m.anm == nil {
		m.msh = newMesh(meshName)
	}
	return m
}
func (m *model) InitMesh(lloc, span, usage uint32, normalize bool) Model {
	if m.msh != nil {
		m.msh.initData(lloc, span, usage, normalize)
	}
	return m
}
func (m *model) SetMeshData(lloc uint32, data interface{}) {
	if m.msh != nil {
		m.msh.setData(lloc, data)
		m.msh.bound = false
	}
}
func (m *model) InitFaces(usage uint32) Model {
	if m.msh != nil {
		m.msh.initFaces(usage)
	}
	return m
}
func (m *model) SetFaces(data []uint16) {
	if m.msh != nil {
		m.msh.setFaces(data)
		m.msh.bound = false
	}
}
func (m *model) SetDrawMode(mode int) Model {
	m.drawMode = mode
	return m
}

// A model may have one more more textures that apply
// to the models mesh.
func (m *model) AddTex(name string) Model {
	index := len(m.texs)
	m.texs = append(m.texs, newTexture(name))
	m.loads = append(m.loads, &loadReq{data: m, index: index, a: newTexture(name)})
	return m
}
func (m *model) SetTex(index int, name string) {
	if index >= 0 && index < len(m.texs) {
		// Add the set request to a list of textures that need to be loaded.
		// These are handled each update.
		req := &loadReq{data: m, index: index, a: newTexture(name)}
		m.loads = append(m.loads, req)
	}
}
func (m *model) NewTex(name string) Model {
	m.texs = append(m.texs, newTexture(name))
	return m
}
func (m *model) SetImg(index int, img image.Image) {
	if img != nil && index >= 0 && index < len(m.texs) {
		m.texs[index].set(img)
	}
}
func (m *model) TexImg(index int) image.Image {
	if index >= 0 && index < len(m.texs) {
		return m.texs[index].img
	}
	return nil
}
func (m *model) SetTexMode(index int, mode int) Model {
	m.texs[index].repeat = false
	if index >= 0 && index < len(m.texs) && mode == TEX_REPEAT {
		m.texs[index].repeat = true
		for _, req := range m.loads {
			if t, ok := req.a.(*texture); ok && req.index == index {
				t.repeat = true
			}
		}
	}
	return m
}

// In this case the texture has been generated in the given layer
// and is already on the gpu. Ignored if there is already a
// textured assigned to the model.
//
// Note: The layer texture must be rendered before the model
//       using it is rendered.
func (m *model) UseLayer(l Layer) Model {
	layer, ok := l.(*layer)
	if m.layer == nil && ok {
		m.layer = layer
		switch m.layer.attr {
		case render.IMAGE_BUFF:
			m.texs = append(m.texs, m.layer.tex)
		case render.DEPTH_BUFF:
			// shadow maps are handled in toDraw.
		}
	}
	return m
}

// Wrap the font classes. Fonts are associated with a mesh
// and a font texture.
func (m *model) LoadFont(fontName string) Model {
	m.fnt = newFont(fontName)
	m.loads = append(m.loads, &loadReq{data: m, a: newFont(fontName)})
	return m
}
func (m *model) SetPhrase(phrase string) Model {
	if m.msh == nil {
		m.msh = newMesh("phrase") // dynamic mesh for phrase backing.
		m.msh.loaded = true       // trigger a rebind in updateModels.
	}
	if len(phrase) > 0 && m.phrase != phrase {
		m.phrase = phrase   // used by loader to set mesh data.
		m.msh.bound = false // mesh will need rebind.
		if m.fnt != nil && m.fnt.loaded {
			m.phraseWidth = m.fnt.setPhrase(m.msh, m.phrase)
		}
	}
	return m
}
func (m *model) PhraseWidth() int { return m.phraseWidth }

// SetUniform combines floats values into a slice of float32's
// that will be passed to the rendering layer and used to set
// shader uniform values.
func (m *model) Uniform(id string) (value []float32) { return m.uniforms[id] }
func (m *model) SetUniform(id string, floats ...interface{}) {
	values := []float32{}
	for _, value := range floats {
		switch v := value.(type) {
		case float32:
			values = append(values, v)
		case float64:
			values = append(values, float32(v))
		case int:
			values = append(values, float32(v))
		default:
			log.Print("vu.model.SetUniform: unknown type ", id, ":", value)
		}
	}
	m.uniforms[id] = values
}

// Animation methods wrap animation class.
// FUTURE: handle animation models with multiple textures. Animation models are
//         currently limited to one texture or they have to be processed after
//         other textures to account for the texture index.
func (m *model) LoadAnim(animName string) Model {
	if m.anm == nil && m.msh == nil {
		m.anm = newAnimation(animName)
		m.loads = append(m.loads, &loadReq{data: m, index: len(m.texs), a: newAnimation(animName)})
		m.texs = append(m.texs, newTexture(animName+"0")) // reserve a texture spot.
	}
	return m
}
func (m *model) Animate(move, frame int) bool {
	if m.anm != nil {
		m.nFrames = m.anm.maxFrames(move)
		m.move = m.anm.isMovement(move)
		if frame < m.nFrames {
			m.frame = float64(frame)
		}
	}
	return move == m.move // was the requested movement available.
}
func (m *model) Action() (move, frame, nFrames int) {
	return m.move, int(math.Floor(m.frame + 1)), m.nFrames
}
func (m *model) Actions() []string {
	if m.anm != nil {
		return m.anm.moveNames()
	}
	return []string{}
}

// Pose returns the bone transform, or the identity matrix
// if there was no transform for the model. The returned matrix
// should not be altered. It is intended for transforming points.
func (m *model) Pose(index int) *lin.M4 {
	if index < len(m.pose) {
		return &m.pose[index]
	}
	return lin.M4I
}

// SetEffect ties the particle effect classes to the model.
func (m *model) SetEffect(mover ParticleEffect, maxParticles int) Model {
	if mover != nil {
		m.effect = newParticleEffect(m, mover, maxParticles)
		m.depth = false // only time model depth is set to false.
	}
	return m
}
func (m *model) SetDepth(enabled bool) Model {
	m.depth = enabled
	return m
}

// animate is called to reposition the poses for an animated model.
func (m *model) animate(dt float64) {
	m.frame = m.anm.animate(dt, m.frame, m.move, m.pose)
	nextFrame := int(math.Floor(m.frame + 1))
	if nextFrame >= m.nFrames {
		m.frame -= float64(m.nFrames - 1)
	}
}

// CastShadow toggles whether or not the model casts a shadow. Default false.
func (m *model) CastShadow() Model {
	m.castShadow = !m.castShadow
	return m
}

// HasShadows toggles whether or not the model reveals shadows. Default false.
func (m *model) HasShadows() Model {
	m.hasShadows = !m.hasShadows
	return m
}

// toDraw sets the model specific bound data references and
// uniform data needed by the rendering layer.
func (m *model) toDraw(d render.Draw, mm *lin.M4) {
	d.SetAlpha(float64(m.alpha)) // 1 : no transparency as the default.

	// Use any previous render to texture passes.
	if m.layer != nil {
		switch m.layer.attr {
		case render.IMAGE_BUFF:
			// handled as regular texture below.
			// Leave it to the shader to use the right the "uv#" uniform.
		case render.DEPTH_BUFF:
			d.SetShadowmap(m.layer.tex.tid) // texture with depth values.

			// Shadow depth bias is the mvp matrix from the light.
			// It is adjusted as needed by shadow maps.
			m.sm.Mult(mm, m.layer.vp)   // model (light) view.
			m.sm.Mult(m.sm, m.layer.bm) // incorporate shadow bias.
			d.SetDbm(m.sm)
		}
	}

	// Set the bound data references.
	d.SetRefs(m.shd.program, m.msh.vao, m.drawMode)
	if total := len(m.texs); total > 0 {
		for cnt, t := range m.texs {
			d.SetTex(total, cnt, t.tid, t.f0, t.fn)
		}
	} else {
		d.SetTex(0, 0, 0, 0, 0) // clear any previous data.
	}

	// Set uniform values. These can be sent as a reference because they
	// are fixed on shader creation.
	d.SetUniforms(m.shd.uniforms) // shader integer uniform references.
	if m.anm != nil && len(m.pose) > 0 {
		d.SetPose(m.pose)
	} else {
		d.SetPose(nil) // clear data.
	}
	d.SetTime(time.Since(m.time).Seconds()) // For shaders that need elapsed time.

	// Set color uniforms.
	d.SetFloats("kd", m.kd.R, m.kd.G, m.kd.B)
	d.SetFloats("ks", m.ks.R, m.ks.G, m.ks.B)
	d.SetFloats("ka", m.ka.R, m.ka.G, m.ka.B)

	// Set user specified uniforms.
	for uniform, uvalues := range m.uniforms {
		d.SetFloats(uniform, uvalues...)
	}
}
