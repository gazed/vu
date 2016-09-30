// Copyright © 2015-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// model.go groups a shader with its related assets. Key app facing API.
// FUTURE: Would ditching the interface and exposing the struct make the
//         application experience better or worse? Possibly need to break
//         Model into specialized types.

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
	Make(assets ...string) Model // Create msh or tex placeholders.
	Set(opts ...ModAttr) Model   // Config with functional options.

	// Mesh allows setting vertex data on meshes. Triggers rebind.
	// Needed for meshes created with Make.
	Mesh() Mesh // One mesh per model.

	// Tex allows setting a texture image. Triggers rebind.
	// Textures are in the order they were added. OrderTex puts the
	// named texture at the given index and shifts the remaining.
	Tex(index int) Tex           // Multiple textures per model.
	ClampTex(name string) Model  // Make texture non-repeating.
	OrderTex(name string, i int) // Order textures for shader.

	// A Model is expected to be one thing at a time. The following
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
	// Assets are added when loading and removed when loaded.
	assets  map[aid]string  // Model is comprised of named assets.
	clamps  map[string]bool // Collects requests for repeating textures.
	rebinds []asset         // Collects rebind requests.

	// Loaded asset instances.
	shd    *shader         // Mandatory GPU render program.
	texs   []*texture      // Optional: one or more texture images.
	tids   []*texid        // Texture placeholder to track ordering.
	mat    *material       // Optional: material lighting info.
	msh    *mesh           // Mandatory vertex buffer data.
	effect *particleEffect // Optional particle effect.
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

	// Rendering attributes.
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
	m := &model{alpha: 1, depth: true, assets: map[aid]string{}}
	m.assets[assetID(shd, shaderName)] = shaderName
	m.clamps = map[string]bool{}

	m.time = time.Now()
	m.uniforms = map[string][]float32{}
	m.sm = &lin.M4{}
	m.Load(attrs...) // default is to load assets.
	return m
}

// Load is called during App.Update.
func (m *model) Load(attrs ...string) Model {
	for _, attribute := range attrs {
		attr := strings.Split(attribute, ":")
		if len(attr) != 2 {
			continue
		}
		name := attr[1]
		switch attr[0] {
		case "mod": // animated model
			m.assets[assetID(anm, name)] = name
			textureName := name + "0"
			aid := assetID(tex, textureName)
			m.assets[aid] = textureName
			m.tids = append(m.tids, newTexid(textureName, aid))
		case "msh": // static model.
			m.assets[assetID(msh, name)] = name
		case "mat": // material for lighting shaders.
			m.assets[assetID(mat, name)] = name
		case "tex": // texture.
			aid := assetID(tex, name)
			m.assets[aid] = name
			m.tids = append(m.tids, newTexid(name, aid))
		case "fnt": // font mapping
			m.assets[assetID(fnt, name)] = name
			m.msh = newMesh("phrase") // dynamic mesh for phrase backing.
		}
	}
	return m
}

// Make is used to create objects where the data is filled by the
// application instead of the loader. Called during App.Update.
func (m *model) Make(attrs ...string) Model {
	for _, attribute := range attrs {
		attr := strings.Split(attribute, ":")
		if len(attr) != 2 {
			continue
		}
		name := attr[1]
		switch attr[0] {
		case "msh": // static model that needs generated data.
			m.newMesh(name)
		case "tex": // texture that needs generated data.
			m.texs = append(m.texs, newTexture(name))
		}
	}
	return m
}

// Set handles updating Model attributes using functional options.
// These are expected to be used rarely since they are many times
// slower than calling a method. See eng_test benchmark.
func (m *model) Set(options ...ModAttr) Model { // Functional options.
	for _, opt := range options {
		opt(m) // apply option by calling function with model instance.
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
func (m *model) Mesh() Mesh {
	if m.msh != nil {
		m.rebinds = append(m.rebinds, m.msh)
	}
	return m.msh
}
func (m *model) newMesh(meshName string) {
	if m.msh == nil && m.anm == nil {
		m.msh = newMesh(meshName)
	}
}

// A model may have one more more textures that apply
// to the models mesh.
func (m *model) Tex(i int) Tex {
	if i < 0 || i >= len(m.texs) {
		return nil
	}
	m.rebinds = append(m.rebinds, m.texs[i])
	return m.texs[i]
}

// ClampTex for when a repeating texture is not going to work,
// as is often the case with rotating textures.
func (m *model) ClampTex(textureName string) Model {
	m.clamps[textureName] = true
	return m
}

// OrderTex tracks the order of textures for the model.
// The texture order must match that expected by the shader.
// Changing a visible texture can be accomplished by switching texture order.
// Note that texture order needs to be maintained if if there are outstanding
// load requests. This allows apps to load and order textures during the
// original create.
func (m *model) OrderTex(name string, i int) {
	if i < 0 || i >= len(m.tids) {
		return
	}
	at := -1
	for cnt, tb := range m.tids {
		if tb.name == name {
			at = cnt
			break
		}
	}
	if at == -1 || at == i {
		return // couldn't find it or its already in its spot.
	}
	a := m.tids[at]
	m.tids = append(m.tids[:at], m.tids[at+1:]...)                     // cut
	m.tids = append(m.tids[:i], append([]*texid{a}, m.tids[i:]...)...) // insert
	if len(m.tids) == len(m.texs) {
		t := m.texs[at]
		m.texs = append(m.texs[:at], m.texs[at+1:]...)                       // cut
		m.texs = append(m.texs[:i], append([]*texture{t}, m.texs[i:]...)...) // insert
	}
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

// SetStr changes the text phrase and causes a mesh rebind.
func (m *model) SetStr(phrase string) Labeler {
	if len(phrase) > 0 && m.phrase != phrase {
		m.phrase = phrase // used by loader to set mesh data.
		m.rebinds = append(m.rebinds, m.msh)
		if m.fnt != nil {
			m.phraseWidth = m.fnt.setPhrase(m.msh, m.phrase)
		}
	}
	return m
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
// FUTURE: handle animation models with multiple textures.
//         Animation models are currently limited to one texture.
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

// =============================================================================
// Functional options for Model.

// ModAttr defines optional model attributes that can be used in Model.Set().
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
var CastShadow = func(m Model) { m.(*model).castShadow = true }

// HasShadows marks a model that reveals shadows. Revealing shadows
// implies a scene with a light and objects that cast shadows.
// It also implies a shadow map capable shader.
var HasShadows = func(m Model) { m.(*model).hasShadows = true }

// functional options.
// =============================================================================
// models

// models is a component that models manages all the active Models instances.
// Models exist in either load or render state.
//
// FUTURE: Look at putting each model attribute into their own indexed array.
//         This is how the Data oriented gurus roll. Benchmark first!
type models struct {
	eng     *engine        // Needed for binding and other machine stuff.
	data    map[eid]*model // All models.
	load    map[eid]*model // New models need to be run through loader.
	loading map[eid]*model // Models waiting for asset loads.
	active  map[eid]*model // Models that can be rendered.
	rebinds []asset        // Scratch slice for assets whose data has changed.
}

// newModels creates the model component manager.
// Expected to be called once on startup.
func newModels(eng *engine) *models {
	ms := &models{eng: eng}
	ms.data = map[eid]*model{}
	ms.load = map[eid]*model{}
	ms.loading = map[eid]*model{}
	ms.active = map[eid]*model{}
	return ms
}

// get loading or loaded models for the given entity.
func (ms *models) get(id eid) *model {
	if mod, ok := ms.data[id]; ok {
		return mod
	}
	return nil
}

// getActive returns fully loaded models ready for rendering.
func (ms *models) getActive(id eid) *model {
	if mod, ok := ms.active[id]; ok {
		return mod
	}
	return nil
}

// create a new model and run it through the loader to
// ensure it has all of its assets.
func (ms *models) create(id eid, shader string, assets ...string) *model {
	if _, ok := ms.data[id]; !ok {
		m := newModel(shader, assets...)
		ms.data[id] = m
		ms.load[id] = m
		return m
	}
	return nil
}

// dispose of the model, removing it from all of the maps.
func (ms *models) dispose(id eid) {
	if m, ok := ms.data[id]; ok {
		m.msh = nil
		m.shd = nil
		m.anm = nil
		m.fnt = nil
		m.mat = nil
		m.texs = []*texture{} // garbage collect all old textures.
	}
	delete(ms.data, id)
	delete(ms.load, id)
	delete(ms.loading, id)
	delete(ms.active, id)
}

func (ms *models) counts() (models, verts int) {
	models = len(ms.data)
	for _, m := range ms.data {
		if m.msh != nil && len(m.msh.vdata) > 0 {
			verts += m.msh.vdata[0].Len()
		}
	}
	return models, verts
}

// process any ongoing model updates like animated models
// and CPU particle effects. Any new models are sent off for loading
// and any updated models generate data rebind requests.
func (ms *models) refresh(dts float64) {
	ms.queueLoads()

	// Process ongoing activity on active models.
	// FUTURE Don't traverse all active models. Have separate
	//        lists for animated models and particle models.
	for _, m := range ms.active {

		if m.effect != nil {
			// udpate and rebind particle effects first since
			// they change mesh data and then need rebinding.
			m.effect.update(m, dt.Seconds())
		}
		if m.anm != nil {
			// animations update the bone position matricies.
			// These are bound as uniforms at draw time.
			m.animate(dts)
		}

		// handle any data updates with rebind requests.
		if len(m.rebinds) > 0 {
			ms.rebinds = append(ms.rebinds, m.rebinds...)
			m.rebinds = m.rebinds[:0] // reset keeping memory.
		}
	}

	// handle all rebind requests at once.
	if len(ms.rebinds) > 0 {
		ms.eng.rebind(ms.rebinds)
		ms.rebinds = ms.rebinds[:0] // reset keeping memory.
	}
}

// queueLoads ensures new models are passed through the loading system.
// Overall there are few assets used by lots of models.
func (ms *models) queueLoads() {
	if len(ms.load) > 0 {
		reqs := map[aid]string{}
		for id, m := range ms.load {
			for aid, name := range m.assets {

				// filter out duplicate load requests.
				if _, ok := reqs[aid]; !ok {
					reqs[aid] = name
				}
			}
			delete(ms.load, id)
			ms.loading[id] = m
		}
		if len(reqs) > 0 {
			ms.eng.submitLoadReqs(reqs)
		}
	}
}

// finishLoads processes models waiting for assets to be loaded.
// It is called by the engine when a new batch of loaded assets
// has been received from the loader.
func (ms *models) finishLoads(assets map[aid]asset) {
	for eid, m := range ms.loading {
		for aid := range m.assets {
			if a, ok := assets[aid]; ok {
				switch at := a.(type) {
				case *mesh:
					m.msh = at
					delete(m.assets, aid)
				case *texture:
					if len(m.texs) == 0 {
						m.texs = make([]*texture, len(m.tids))
					}
					for cnt, tid := range m.tids {
						if tid.id == a.aid() {
							m.texs[cnt] = at
							if _, ok := m.clamps[at.name]; ok {
								ms.eng.clampTex(at.tid)
							}
						}
					}
					delete(m.assets, aid)
				case *shader:
					m.shd = at
					delete(m.assets, aid)
				case *font:
					m.fnt = at
					if len(m.phrase) > 0 {
						m.phraseWidth = m.fnt.setPhrase(m.msh, m.phrase)
					}
					delete(m.assets, aid)
				case *animation:
					m.anm = at
					m.nFrames = at.maxFrames(0)
					m.pose = make([]lin.M4, len(at.joints))

					// Mesh created when loading animation data.
					if at, ok := assets[assetID(msh, m.anm.name)]; ok {
						m.msh = at.(*mesh)
					}
					delete(m.assets, aid)
				case *material:
					m.mat = at

					// Override with material if not directly set by app.
					if m.alpha == 1.0 {
						m.alpha = float64(at.tr)
					}
					delete(m.assets, aid)
				}
			}
		}
		if len(m.assets) == 0 {
			delete(ms.loading, eid)
			ms.active[eid] = m
		}
	}
}
