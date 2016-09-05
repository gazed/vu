// Copyright © 2015-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"log"
	"math"
	"strings"
	"time"

	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/render"
)

// Model manages rendered 3D objects. Model is the link between loaded
// assets and the data needed by the rendering system. Assets, such as
// shaders, mesh, textures, etc., are specified as unique strings in the
// Load method. Assets are loaded and converted to intermediate data
// which is later converted to render draws for the render system.
// Make is used for Models that are generated instead of loaded.
//
// A model is expected to be attached to a point-of-view (Pov) to give
// it a 3D location and orientation in the render transform hierarchy.
//
// It is the applications responsibility to call enough Model methods
// to ensure the Models shader is provided with its expected asset data.
// Ie: if the shader expects a texture, ensure a texture is loaded.
// Mismatches generate runtime logs.
type Model interface {
	Shader() (name string)       // Rendered models have a shader.
	Load(assets ...string) Model // Create and import assets.
	Make(assets ...string) Model // Create asset placeholders.
	Set(opts ...ModAttr) Model   // Config with functional options.

	// Models shape is described by a set of verticies connected as
	// triangles. Other per-vertex data can be assigned to a Mesh.
	Mesh() Mesh // One mesh per model.

	// Models can have one or more textures applied to a single mesh.
	// Textures are initialized from assets and can be updated with images.
	Tex(index int) Tex          // There can be multiple textures per model.
	LoadTex(i int, name string) // Load replaces existing model Texture.

	// Models are expected to be one thing at a time. The following
	// are particular combinations of the shader, mesh, and textures
	Animator // Animators are animated Models.
	Labeler  // Labels are Models used to display small text phrases.

	// Particle effects are either CPU:application or GPU:shader based.
	// SetEffect sets a CPU controlled particle effect.
	SetEffect(mover ParticleEffect, maxParticles int)

	// Layers are used for shadows or render to texture.
	UseLayer(l Layer) // Use render pass texture.

	// Set/get shader uniform values where id is the shader uniform name.
	Uniform(id string) (value []float32)         // Uniform name/values.
	SetUniform(id string, floats ...interface{}) // Individual values.

	// Alpha is model transparency used in shaders.
	Alpha() (a float64) // 0 for fully transparent, to 1 fully opaque.
	SetAlpha(a float64) // Overrides any material alpha values.
}

// Model
// =============================================================================
// model implements Model.

// model implements Model. It links the applications transform hierarchy
// to the rendering system. The application specifies the model resources.
// These resources are later linked and bound during engine processing.
type model struct {
	shd    *shader         // Mandatory GPU render program.
	texs   []*texture      // Optional: one or more texture images.
	mat    *material       // Optional: material lighting info.
	msh    *mesh           // Mandatory vertex buffer data.
	effect *particleEffect // Optional particle effect.
	loads  []*loadReq      // Assets waiting to be loaded.
	layer  *layer          // Optional previous render pass.

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

	// Rendering attributes .
	castShadow bool // Model to cast a shadow. Default false.
	hasShadows bool // Model to reveal a shadow. Default false.
	depth      bool // Depth buffer on by default.
	drawMode   int  // Render mesh as Triangles, Points, Lines.

	// Shader dependent uniform data.
	alpha    float64              // Transparency between 0 and 1.
	time     time.Time            // Time needed by some shaders.
	uniforms map[string][]float32 // Uniform values.
	sm       *lin.M4              // scratch matrix.
}

// newModel allocates a new model instance setting some common defaults.
// It parses the requested model attributes, creating placeholders for
// the necessary rendering objects. These will be loaded or filled in
// later.
func newModel(shaderName string, attrs ...string) *model {
	m := &model{alpha: 1, depth: true}
	m.shd = newShader(shaderName)
	m.loads = append(m.loads, &loadReq{data: m, a: newShader(shaderName)})
	m.time = time.Now()
	m.uniforms = map[string][]float32{}
	m.sm = &lin.M4{}
	m.Load(attrs...) // the default is to load assets.
	return m
}

// Load is called during App.Update.
func (m *model) Load(attrs ...string) Model {
	for _, attribute := range attrs {
		attr := strings.Split(attribute, ":")
		if len(attr) != 2 {
			continue
		}
		switch attr[0] {
		case "mod":
			m.loadAnim(attr[1]) // animated model
		case "msh":
			m.loadMesh(attr[1]) // static model.
		case "mat":
			m.loadMat(attr[1]) // material for lighting shaders.
		case "tex":
			m.addTex(attr[1]) // texture.
		case "fnt":
			m.loadFont(attr[1]) // font mapping
		}
	}
	return m
}

// Make is called during App.Update.
func (m *model) Make(attrs ...string) Model {
	for _, attribute := range attrs {
		attr := strings.Split(attribute, ":")
		if len(attr) != 2 {
			continue
		}
		switch attr[0] {
		case "msh":
			m.newMesh(attr[1]) // static model that needs data.
		case "tex":
			m.newTex(attr[1]) // texture that needs data.
		}
	}
	return m
}

// Set handles updating Model attributes using functional options.
// These are expected to be used rarely since they are many times
// slower than calling a method. See eng_test benchmark.
func (m *model) Set(options ...ModAttr) Model { // Functional options.
	for _, opt := range options {
		opt(m)
	}
	return m
}

// Alpha is model transparency. This value overrides any material values.
// Alpha is used often enough to separate it from Uniforms and ModAttrs.
func (m *model) Alpha() (a float64) { return m.alpha }
func (m *model) SetAlpha(a float64) { m.alpha = a }

// Each model has one shader.
func (m *model) Shader() string { return m.shd.name }

// Each model may have one mesh.
func (m *model) Mesh() Mesh { return m.msh }

// Material is used to help with coloring for shaders that use lights.
// Overrides existing values if it was the last one set.
func (m *model) loadMat(name string) {
	m.mat = newMaterial(name)
	m.loads = append(m.loads, &loadReq{data: m, a: newMaterial(name)})
}

// Each model has one mesh. The mesh is specified here and
// will be sent for loading and binding later on.
func (m *model) loadMesh(meshName string) {
	if m.msh == nil && m.anm == nil {
		m.msh = newMesh(meshName) // placeholder
		req := &loadReq{data: m, a: newMesh(meshName)}
		m.loads = append(m.loads, req)
	}
}
func (m *model) newMesh(meshName string) {
	if m.msh == nil && m.anm == nil {
		m.msh = newMesh(meshName)
	}
}

// A model may have one more more textures that apply
// to the models mesh.
func (m *model) Tex(index int) Tex {
	if index < 0 || index >= len(m.texs) {
		return nil
	}
	return m.texs[index]
}

func (m *model) LoadTex(index int, name string) {
	if index >= 0 && index < len(m.texs) {
		// Add the set request to a list of textures that need to be loaded.
		// These are handled each update.
		req := &loadReq{data: m, index: index, a: newTexture(name)}
		m.loads = append(m.loads, req)
	}
}
func (m *model) addTex(name string) {
	index := len(m.texs)
	m.texs = append(m.texs, newTexture(name))
	m.loads = append(m.loads, &loadReq{data: m, index: index, a: newTexture(name)})
}
func (m *model) newTex(name string) {
	m.texs = append(m.texs, newTexture(name))
}

// In this case the texture has been generated in the given layer
// and is already on the gpu. Ignored if there is already a
// textured assigned to the model.
//
// Note: The layer texture must be rendered before the model
//       using it is rendered.
func (m *model) UseLayer(l Layer) {
	layer, ok := l.(*layer)
	if m.layer == nil && ok {
		m.layer = layer
		switch m.layer.attr {
		case render.ImageBuffer:
			m.texs = append(m.texs, m.layer.tex)
		case render.DepthBuffer:
			// shadow maps are handled in toDraw.
		}
	}
}

// Wrap the font classes. Fonts are associated with a mesh
// and a font texture.
func (m *model) loadFont(fontName string) *model {
	m.fnt = newFont(fontName)
	m.loads = append(m.loads, &loadReq{data: m, a: newFont(fontName)})
	return m
}
func (m *model) SetStr(phrase string) {
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
}
func (m *model) StrWidth() int { return m.phraseWidth }

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
			log.Print("model.SetUniform: unknown type ", id, ":", value)
		}
	}
	m.uniforms[id] = values
}

// Animation methods wrap animation class.
// FUTURE: handle animation models with multiple textures. Animation models are
//         currently limited to one texture or they have to be processed after
//         other textures to account for the texture index.
func (m *model) loadAnim(animName string) {
	if m.anm == nil && m.msh == nil {
		m.anm = newAnimation(animName)
		m.loads = append(m.loads, &loadReq{data: m, index: len(m.texs), a: newAnimation(animName)})
		m.texs = append(m.texs, newTexture(animName+"0")) // reserve a texture spot.
	}
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
func (m *model) SetEffect(mover ParticleEffect, maxParticles int) {
	if mover != nil {
		m.effect = newParticleEffect(m, mover, maxParticles)
		m.depth = false // only time model depth is set to false.
	}
}

// animate is called to reposition the poses for an animated model.
func (m *model) animate(dt float64) {
	m.frame = m.anm.animate(dt, m.frame, m.move, m.pose)
	nextFrame := int(math.Floor(m.frame + 1))
	if nextFrame >= m.nFrames {
		m.frame -= float64(m.nFrames - 1)
	}
}

func (m *model) queueLoads(requests []*loadReq) ([]*loadReq, bool) {
	if len(m.loads) <= 0 {
		return requests, false
	}

	// Propogate attributes needed for binding, but which where set
	// after the initial load request method call.
	for _, req := range m.loads {
		if t, ok := req.a.(*texture); ok {
			t.repeat = m.texs[req.index].repeat
		}
	}

	requests = append(requests, m.loads...)
	m.loads = m.loads[:0]
	return requests, true
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

// toDraw sets the model specific bound data references and
// uniform data needed by the rendering layer.
func (m *model) toDraw(d *render.Draw, mm *lin.M4) {

	// Use any previous render to texture passes.
	if m.layer != nil {
		switch m.layer.attr {
		case render.ImageBuffer:
			// handled as regular texture below.
			// Leave it to the shader to use the right the "uv#" uniform.
		case render.DepthBuffer:
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

	// Material transparency.
	d.SetFloats("alpha", float32(m.alpha))

	// Material color uniforms.
	if mat := m.mat; mat != nil {
		drawMaterial(d, mat)
	}

	// For shaders that need elapsed time.
	d.SetFloats("time", float32(time.Since(m.time).Seconds()))

	// Set user specified uniforms.
	for uniform, uvalues := range m.uniforms {
		d.SetFloats(uniform, uvalues...)
	}
}

// drawMaterial sets the data needed by the render system.
func drawMaterial(d *render.Draw, m *material) {
	d.SetFloats("kd", m.kd.R, m.kd.G, m.kd.B)
	d.SetFloats("ks", m.ks.R, m.ks.G, m.ks.B)
	d.SetFloats("ka", m.ka.R, m.ka.G, m.ka.B)
	d.SetFloats("ns", m.ns)
}

// drawMesh sets the data needed by the render system.
// In this case the vao is the reference to the mesh data on the GPU.
func drawMesh(d *render.Draw, m *mesh) { d.Vao = m.vao }

// drawLight sets the data needed by the render system.
// In this case the light color.
func drawLight(d *render.Draw, l *Light, px, py, pz float64) {
	d.SetFloats("lp", float32(px), float32(py), float32(pz))    // position
	d.SetFloats("lc", float32(l.R), float32(l.G), float32(l.B)) // color
}

// =============================================================================
// Functional options for Model.

// ModAttr defines a model attribute that can be used in Model.Set().
type ModAttr func(Model)

// DrawMode affects rendered meshes by rendering with Triangles, Lines, Points.
func DrawMode(mode int) ModAttr {
	return func(m Model) { m.(*model).drawMode = mode }
}

// SetDepth toggles Z-Buffer awareness while rendering.
// Generally on, it is often disabled for particle effects.
func SetDepth(enabled bool) ModAttr {
	return func(m Model) { m.(*model).depth = enabled }
}

// CastShadow marks a model that can cast shadows. Casting shadows
// implies a scene with a light and objects that receive shadows.
// It also implies a shadow map capable shader.
func CastShadow() ModAttr {
	return func(m Model) { m.(*model).castShadow = true }
}

// HasShadows marks a models can reveal shadows. Revealing shadows
// implies a scene with a light and objects that cast shadows.
// It also implies a shadow map capable shader.
func HasShadows() ModAttr {
	return func(m Model) { m.(*model).hasShadows = true }
}
