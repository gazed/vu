// Copyright © 2015-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// model.go contains code related to the model component.

import (
	"log"
	"math"
	"strings"
	"time"

	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/render"
)

// MakeModel adds a model component to an entity.
// The following assets can be loaded into a model. The first parameter
// is the shader asset name. The remaining assets can be in any order
// and are prefixed with an asset type. Names must be unique within
// an asset type.
//    shd: first parameter, "name" of shader - one per model.
//    msh: "msh:name" - one per model.
//    tex: "tex:name" - 0 to 14 per model.
//    mat: "mat:name" - 0 or 1 per model.
//
// A model manages rendered 3D objects. It is the link between loaded
// assets and the data needed by the rendering system for the shader.
// Assets, such as shaders, mesh, textures, etc., are specified as unique
// strings. The assets are located, loaded, and converted to intermediate
// data which is later converted to render draws for the render system.
//
// It is the applications responsibility to call enough model methods
// to ensure the models shader is provided with its expected asset data.
// Ie: if the shader expects a texture, ensure a texture is loaded.
//
// A model is attached to an entity with a point-of-view to give
// it a 3D location and orientation.
func (e *Ent) MakeModel(shader string, attrs ...string) *Ent {
	mod := e.app.models.create(e)
	e.app.models.loadAssets(e, mod, append(attrs, "shd:"+shader)...)
	return e
}

// MakeInstancedModel is similar to MakeModel except this model is marked
// as having multiple instances that will be rendered in a single draw call.
// An instanced model needs child Parts to be rendered. Each child Part
// provides the instance transform data (postion,rotation,scale) for one
// instance of this model.
//
// This model cannot be a composite model in that its child Parts do not
// have models - only transform data.
// Works only with the default SetDraw mode of Triangles.
func (e *Ent) MakeInstancedModel(shader string, attrs ...string) *Ent {
	mod := e.app.models.create(e)
	mod.isInstanced = true
	e.app.models.loadAssets(e, mod, append(attrs, "shd:"+shader)...)
	return e
	// Design Note: https://learnopengl.com/Advanced-OpenGL/Instancing
}

// Load more model assets after the model has been created.
// The assets can be in any order and are prefixed with an asset type.
// Names must be unique within an asset type.
//    msh: "msh:name" - one per model.
//    tex: "txt:name" - 0 to 14 per model.
//    mat: "mat:name" - 0 to 1 per model.
//
// Depends on Ent.MakeModel.
func (e *Ent) Load(assets ...string) *Ent {
	if m := e.app.models.get(e.eid); m != nil {
		e.app.models.loadAssets(e, m, assets...)
		return e
	}
	log.Printf("Load needs MakeModel %d", e.eid)
	return e
}

// Mesh returns the vertex data for an entity with a model component.
// The Mesh will be marked for rebinding as the expected reason
// to get a mesh is to modify it.
//
// Depends on Ent.MakeModel. Returns nil if missing model component.
func (e *Ent) Mesh() *Mesh {
	if m := e.app.models.get(e.eid); m != nil {
		e.app.models.rebinds[e.eid] = m
		return m.msh
	}
	log.Printf("Mesh needs MakeModel %d", e.eid)
	return nil
}

// GenMesh is used to create a mesh where the data is filled by the
// application instead of the loader. Does nothing if the original
// MakeModel call already loaded a mesh.
//
// Depends on Ent.MakeModel. Returns nil if missing model component.
func (e *Ent) GenMesh(name string) *Mesh {
	if m := e.app.models.get(e.eid); m != nil {
		if m.msh == nil {
			m.msh = newMesh(name)
			m.track[msh] = 1
			e.app.models.rebinds[e.eid] = m
		}
		return m.msh
	}
	log.Printf("GenMesh needs MakeModel %d", e.eid)
	return nil
}

// Tex returns the texture for an entity with a model component and
// a texture at the given index. The Texture will be marked for rebinding
// as the expected reason  to get a texture is to modify it.
//
// Depends on Ent.MakeModel. Returns nil if missing model component.
func (e *Ent) Tex(index int) *Texture {
	if m := e.app.models.get(e.eid); m != nil {
		e.app.models.rebinds[e.eid] = m
		if index >= 0 && index < len(m.texs) {
			return m.texs[index]
		}
	}
	log.Printf("Tex needs MakeModel %d", e.eid)
	return nil
}

// GenTex is used to create a mesh where the data is filled by the
// application instead of the loader. Can be called on an existing
// model entity.
//
// Depends on Ent.MakeModel. Returns nil if missing model component.
func (e *Ent) GenTex(name string) *Texture {
	if m := e.app.models.get(e.eid); m != nil {
		t := newTexture(name)
		m.texs = append(m.texs, t)
		m.tpos[name] = len(m.texs) - 1
		m.track[tex] = m.track[tex] + 1
		e.app.models.rebinds[e.eid] = m
		return t
	}
	log.Printf("GenTex needs MakeModel %d", e.eid)
	return nil
}

// SetTex assigns the model a texture that has been generated
// from a scene and which already exists on the GPU.
// Ignored if there is already a texture assigned to the model.
//
// Depends on Ent.MakeModel.
//     scene : depends on Eng.AddScene.
func (e *Ent) SetTex(scene *Ent) *Ent {
	m := e.app.models.get(e.eid)
	t := e.app.scenes.getTarget(scene.eid)
	if len(m.texs) == 0 && t != nil {
		m.texs = append(m.texs, t.tex)
		m.track[tex] = 1
		return e
	}
	log.Printf("SetTex needs MakeModel %d and AddScene %d", e.eid, scene.eid)
	return e
}

// SetFirst ensures that the named texture is at texture position 0.
// Position 0 is the texture position most shaders use to color a model.
// This allows models to have multiple textures and change from displaying
// one texture to another. The order of the remaining textures is not
// guaranteed.
//
// Nothing happens if the texture is invalid. Expected to be called
// after textures have been loaded.
//
// Depends on Ent.MakeModel.
func (e *Ent) SetFirst(name string) *Ent {
	if m := e.app.models.get(e.eid); m != nil {
		m.setFirstTexture(name)
		return e
	}
	log.Printf("SetFirst needs MakeModel %d", e.eid)
	return e
}

// SetUniform combines floats values into a slice of float32's
// that will be passed to rendering and used to set shader uniform values.
//
// Depends on Ent.MakeModel.
func (e *Ent) SetUniform(id string, floats ...interface{}) *Ent {
	if m := e.app.models.get(e.eid); m != nil {
		m.setUniforms(id, floats...)
		return e
	}
	log.Printf("SetUniform needs MakeModel %d", e.eid)
	return e
}

// Alpha returns the alpha value for this model entity.
//
// Depends on Ent.MakeModel. Returns 0 if there is no model component.
func (e *Ent) Alpha() float64 {
	if m := e.app.models.get(e.eid); m != nil {
		return m.alpha()
	}
	log.Printf("Alpha needs MakeModel %d", e.eid)
	return 0
}

// SetAlpha sets the alpha value for this model entity.
//
// Depends on Ent.MakeModel. Sets nothing if there is no model component.
func (e *Ent) SetAlpha(a float64) *Ent {
	if m := e.app.models.get(e.eid); m != nil {
		m.setUniforms("alpha", a)
		return e
	}
	log.Printf("SetAlpha needs MakeModel %d", e.eid)
	return e
}

// SetColor sets a model component material color where
// the r,g,b values are from 0-1.
//
// Depends on Ent.MakeModel.
func (e *Ent) SetColor(r, g, b float64) *Ent {
	if m := e.app.models.get(e.eid); m != nil {
		m.setUniforms("kd", r, g, b)
		return e
	}
	log.Printf("SetColor needs MakeModel %d", e.eid)
	return e
}

// SetDraw affects rendered meshes by rendering with Triangles, Lines, Points.
// The default mode is Triangles.
//
// Depends on Ent.MakeModel.
func (e *Ent) SetDraw(mode int) *Ent {
	if m := e.app.models.get(e.eid); m != nil {
		m.mode = mode
		return e
	}
	log.Printf("SetDraw needs MakeModel %d", e.eid)
	return e
}

// Clamp a texture instead of using the default repeating texture.
// Expected to be called after a model with texture assets has been defined,
// but before the model assets have been loaded. Ie: something that is set
// before first use, not for flipping back and forth.
func (e *Ent) Clamp(name string) *Ent {
	// FUTURE: would be nice to know whether or not to clamp as the asset
	//         is being created.
	e.app.models.clamps[e.eid] = append(e.app.models.clamps[e.eid], name)
	return e
}

// model entity methods.
// =============================================================================
// model data

// model transforms groups of application requested assets into
// render draw data. It also provides a consistent API for application
// access and modification of render data.
//
// Generally expected to be accessed through wrapper classes such as
// Model, Particle, Actor, and Label.
type model struct {
	// Base asset instances will be nil until loaded.
	msh   *Mesh      // Mandatory vertex buffer data.
	shd   *shader    // Mandatory GPU render program.
	texs  []*Texture // Optional: one or more texture images.
	mat   *material  // Optional: material lighting info.
	tocam float64    // distance to camera helps with 3D render order.
	mode  int        // Render as Triangles, Points, Lines, or Instanced.

	// Mark as instanced, meanings all child Pov's are treated
	// as the same model with different transforms.
	isInstanced bool // True renders all instances in one draw call.

	// Mark as a effect since GPU effects don't have effect data.
	isEffect bool // Mark GPU and CPU particle effects.

	// Allow textures to change order. Must be able to do this
	// before assets are back from loading.
	tpos map[string]int // Texture ordering to match shader.

	// match requested assets with loaded loaded assets to know when
	// to render. Asset types are initialized to 0 in constructor.
	track map[int]int // expected number of each asset type.

	time     time.Time            // Shader uniform.
	uniforms map[string][]float32 // Shader uniform names and data.
}

// newModel initializes the data structures and default uniforms.
func newModel() *model {
	m := &model{}
	m.uniforms = map[string][]float32{}
	m.track = map[int]int{shd: 0, msh: 0, tex: 0, mat: 0, fnt: 0, anm: 0}
	m.tpos = map[string]int{}
	m.setUniforms("alpha", 1)
	m.time = time.Now()
	return m
}

// alpha returns the uniform alpha value, returning the default 1
// if there is none.
func (m *model) alpha() float64 {
	if values, ok := m.uniforms["alpha"]; ok && len(values) > 0 {
		return float64(values[0])
	}
	return 1
}

// setUniforms is a helper method for externally visible SetUniform.
func (m *model) setUniforms(id string, floats ...interface{}) {
	values, ok := m.uniforms[id]
	if !ok {
		values = []float32{}
	}
	values = values[:0] // reset preserving memory.
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

// setFirstTexture ensures the named texture is the first
// texture, order == 0, selected by a shader.
func (m *model) setFirstTexture(name string) {
	swap := m.tpos[name]
	for key, index := range m.tpos {
		if index == 0 {
			m.tpos[key] = swap
			m.tpos[name] = 0
			break
		}
	}
}

// draw turns a Model instance into draw call data.
// Expects to be called after pov.draw()
func (m *model) draw(d *render.Draw, shd *shader) {
	if shd == nil {
		shd = m.shd
	}
	d.SetRefs(shd.program, m.msh.vao, m.mode)
	d.SetCounts(m.msh.counts())
	for cnt, t := range m.texs {
		d.SetTex(len(m.texs), cnt, m.tpos[t.name], t.tid, t.f0, t.fn)
	}

	// Set uniform names and references needed by the shader.
	d.SetUniforms(shd.uniforms)   // shader integer uniform references.
	if mat := m.mat; mat != nil { // Material color uniforms.
		mat.draw(d)
	}

	// Set the shader data reference by the uniform names.
	for uniform, uvalues := range m.uniforms {
		d.SetFloats(uniform, uvalues...)
	}
	d.SetFloats("time", float32(time.Since(m.time).Seconds()))

	// If instanced then get the current number of instances from
	// the mesh data.
	if m.isInstanced {
		d.Instances = int32(m.msh.instances)
	}

	// Update bucket draw order information for transparent objects.
	d.Bucket = setDist(d.Bucket, m.tocam)
	alpha := m.uniforms["alpha"][0]
	switch {
	case alpha == 1:
		d.Bucket = setOpaque(d.Bucket)
	case alpha >= 0 && alpha < 1:
		d.Bucket = setTransparent(d.Bucket)
	default:
		log.Printf("Alpha outside range 0->1 %f", alpha) // Dev error.
	}
}

// drawEffect converts a model particle effect into a draw call.
func (m *model) drawEffect(d *render.Draw) {
	d.SetRefs(m.shd.program, m.msh.vao, render.Points)
	d.Depth = false
	d.SetCounts(m.msh.counts())
	d.SetTex(1, 0, 0, m.texs[0].tid, 0, 0)
	d.SetUniforms(m.shd.uniforms)
	d.SetFloats("time", float32(time.Since(m.time).Seconds()))
	d.Bucket = setOpaque(d.Bucket)
	for uniform, uvalues := range m.uniforms {
		d.SetFloats(uniform, uvalues...)
	}
}

// =============================================================================

// models is the component manager for model data. Its job is to turn model
// data into render.draw calls. The lifecycle of a model starts with an asset
// load request. This is a space and colon separated string specified by
// the application
//    "shd:name msh:name tex:name ..." Used to locate disk assets.
type models struct {
	all     map[eid]*model // All model objects.
	rebinds map[eid]*model // Models need asset rebinds.
	loading map[eid]*model // Waiting for assets.
	ready   map[eid]*model // Assets received.

	// Optional data information for specific model types.
	// Processing for each data type depends on a specific set of assets.
	actors  map[eid]*actor  // Depends on anm asset.
	acting  map[eid]*actor  // Actors ready for animations.
	effects map[eid]*effect // CPU particle effect.
	labels  map[eid]*label  // Depends on fnt asset.
	idata   []float32       // Scratch for instance model transform data.

	// Texture clamps are textures that need clamping. Need to remember
	// which ones because the request often happens when the texture
	// asset is away for loading. They are processed once the texture
	// is loaded.
	clamps map[eid][]string // Texture clamps.
}

// newModels creates the render model component manager.
// Expected to be called once on startup.
func newModels() *models {
	ms := &models{}
	ms.all = map[eid]*model{}      // any model in any state.
	ms.loading = map[eid]*model{}  // waiting for initial assets.
	ms.ready = map[eid]*model{}    // assets received.
	ms.rebinds = map[eid]*model{}  // model needing rebinds.
	ms.actors = map[eid]*actor{}   // optional actor data.
	ms.acting = map[eid]*actor{}   // live actors data.
	ms.effects = map[eid]*effect{} // optional particle data.
	ms.labels = map[eid]*label{}   // optional label data.
	ms.clamps = map[eid][]string{} // optional texture clamps
	return ms
}

// create is called by the App on the update goroutine. A new model object
// placeholder is created and the asset loading kicked off.
func (ms *models) create(e *Ent) *model {
	if _, ok := ms.all[e.eid]; !ok {
		m := newModel()
		ms.all[e.eid] = m
		ms.loading[e.eid] = m
		// leave control of when to load assets to the caller.
		// Assets from cache may load and immediately call loaded.
		return m
	}
	return nil
}

// createLabel adds extra label data in addition to the model data.
func (ms *models) createLabel(e *Ent, assets ...string) *model {
	m := ms.create(e)
	if _, ok := ms.labels[e.eid]; !ok {
		ms.labels[e.eid] = &label{str: " "} // nil strings invalidates mesh.
	}
	m.msh = newMesh("phrase")    // dynamic mesh for phrase backing.
	m.track[msh] = 1             // track generated mesh.
	m.setUniforms("kd", 1, 1, 1) // default color is white.
	e.app.models.loadAssets(e, m, assets...)
	return m
}

// createActor adds extra actor data in addition to the model data.
func (ms *models) createActor(e *Ent, assets ...string) *model {
	m := ms.create(e)
	if _, ok := ms.actors[e.eid]; !ok {
		ms.actors[e.eid] = &actor{}
	}
	e.app.models.loadAssets(e, m, assets...)
	return m
}

// createEffect adds particle effect tracking data in addition to the model data.
func (ms *models) createEffect(e *Ent, assets ...string) *model {
	m := e.app.models.create(e)
	m.msh = newMesh("effect")
	m.mode = Points
	m.isEffect = true
	e.app.models.loadAssets(e, m, assets...)
	return m
}

// loadAssets adds assets to the current model.
// The defers are used to ensure fetching happends after the assets are
// tracked. Otherwise cached assets can immediately return in a loaded
// callback before all tracking data is updated.
func (ms *models) loadAssets(e *Ent, m *model, assets ...string) {
	ld, id := e.app.ld, e.eid
	callback := func(a asset) { ms.loaded(id, a, e.app) }
	for _, attribute := range assets {
		attr := strings.Split(attribute, ":")
		if len(attr) != 2 {
			continue
		}
		name := attr[1]
		switch attr[0] {
		case "msh": // static model.
			m.track[msh] = 1
			defer ld.fetch(newMesh(name), callback)
		case "mat": // material for lighting shaders.
			m.track[mat] = 1
			defer ld.fetch(newMaterial(name), callback)
		case "tex": // texture.
			m.tpos[name] = len(m.tpos)
			m.track[tex] = m.track[tex] + 1
			defer ld.fetch(newTexture(name), callback)
		case "shd": // shader.
			m.track[shd] = 1
			defer ld.fetch(newShader(name), callback)
		case "fnt":
			m.track[fnt] = 1
			defer ld.fetch(newFont(name), callback)
		case "anm":
			m.track[anm] = 1
			m.track[msh] = 1
			defer ld.fetch(newAnimation(name), callback)
		default:
			log.Printf("Unknown model asset %s %s for %d", attr[0], name, id)
		}
	}
}

// loaded is called from the update goroutine when a model asset
// has finished loading.
func (ms *models) loaded(eid eid, a asset, app *application) {
	m := ms.get(eid)
	if m == nil {
		// The model may have been disposed before it finished loading.
		log.Printf("No model for %d asset %s", eid, a.label())
		return
	}
	switch la := a.(type) {
	case *shader:
		m.shd = la
		ms.rebinds[eid] = m
	case *Texture:
		m.texs = append(m.texs, la)
		ms.processTexture(eid, m, la)
		ms.rebinds[eid] = m
	case *Mesh:
		m.msh = la
		if m.isInstanced {
			// clone the underlying mesh to get separate vao and data instances.
			// Each instanced model stores the transform data for its children
			// so it can't use the cached mesh.
			// The cached mesh remains available for non-instanced use.
			m.msh = m.msh.clone()
			m.msh.vao = 0          // get new vao on upcoming rebind.
			m.msh.InitInstances(6) // allocate instance transform space.
			ms.setInstanced(eid, m, app)
		}
		ms.rebinds[eid] = m
	case *material:
		m.mat = la
		if m.alpha() == 1.0 {
			// Set with material alpha if app has not
			// overridden the default alpha value.
			m.setUniforms("alpha", m.mat.tr)
		}
	case *animation:
		a := ms.getActor(eid)
		if a == nil {
			log.Printf("Animation data for actor %d", eid)
			return
		}
		a.anm = la
		a.nFrames = a.anm.maxFrames(0)
		a.pose = make([]lin.M4, len(a.anm.joints))
	case *font:
		l := ms.getLabel(eid)
		if l == nil {
			log.Printf("Font data for label %d", eid)
			return
		}
		l.fnt = la
		ms.updateLabel(eid, l, m)
	default:
		log.Printf("Unexepected model asset: %s", la.label())
	}

	// move to the ready queue if all the assets have been loaded.
	if ms.canRender(m, eid) {
		delete(ms.loading, eid)
		ms.ready[eid] = m
		if a, ok := ms.actors[eid]; ok {
			ms.acting[eid] = a
		}
	}
}

// canRender returns true when the model has all its assets.
// Models have 1 asset per type, except for textures which can have
// zero or more per model.
func (ms *models) canRender(m *model, eid eid) bool {
	if m.shd == nil {
		return false // always need a shader for a model.
	}
	if m.track[msh] != 0 {
		if m.msh == nil {
			return false
		}
	}
	if m.track[tex] != len(m.texs) {
		return false
	}
	if m.track[mat] != 0 && m.mat == nil {
		return false
	}
	if m.track[anm] != 0 {
		a := ms.getActor(eid)
		if a == nil || a.anm == nil {
			return false
		}
	}
	if m.track[fnt] != 0 {
		l := ms.getLabel(eid)
		if l == nil || l.fnt == nil {
			return false
		}
	}
	return true
}

// updateInstanced is called each refresh to check if any of the instanced
// models child parts has moved. Ie, if a part has a loaded parent model that
// is instanced, then the parent model mesh needs to be updated with the
// childs new transform data.
func (ms *models) updateInstanced(app *application, changed []eid) {
	needsUpdating := map[eid]*model{}
	for _, eid := range changed {
		if node := app.povs.getNode(eid); node.parent != 0 {
			if m := app.models.getReady(node.parent); m != nil && m.isInstanced {
				needsUpdating[node.parent] = m
			}
		}
	}
	for eid, m := range needsUpdating {
		app.models.setInstanced(eid, m, app)
	}
}

// setInstanced updates the transform data for an instanced model.
func (ms *models) setInstanced(eid eid, m *model, app *application) {
	if !m.isInstanced {
		log.Printf("Called on non-instanced model")
		return
	}
	if m.msh == nil {
		log.Printf("Called on non-mesh model")
		return
	}
	n := app.povs.getNode(eid)
	if n == nil {
		log.Printf("No node %d", eid)
		return
	}
	ms.idata = ms.idata[:0] // reset keeping capacity.
	if len(n.kids) > 0 {    // need kids to render.
		for _, kidEid := range n.kids {
			if kid := app.povs.get(kidEid); kid != nil {
				m4 := kid.wm // world transform matrix.

				// Copy same as memory order, as expected by the shader.
				ms.idata = append(ms.idata, float32(m4.Xx), float32(m4.Xy), float32(m4.Xz), float32(m4.Xw)) // X-Axis
				ms.idata = append(ms.idata, float32(m4.Yx), float32(m4.Yy), float32(m4.Yz), float32(m4.Yw)) // Y-Axis
				ms.idata = append(ms.idata, float32(m4.Zx), float32(m4.Zy), float32(m4.Zz), float32(m4.Zw)) // Z-Axis
				ms.idata = append(ms.idata, float32(m4.Wx), float32(m4.Wy), float32(m4.Wz), float32(m4.Ww)) // Transform.
			}
		}
		m.msh.SetData(6, ms.idata)    // convention: 6 matches models.loaded method
		m.msh.instances = len(n.kids) // Non-zero when mesh is to be drawn instanced.
		ms.rebinds[eid] = m
	}
}

// get loading or loaded models for the given entity.
func (ms *models) get(eid eid) *model      { return ms.all[eid] }
func (ms *models) getLabel(eid eid) *label { return ms.labels[eid] }
func (ms *models) getActor(eid eid) *actor { return ms.actors[eid] }
func (ms *models) getReady(eid eid) *model { return ms.ready[eid] }

// animate updates the animations. Needs to be called each update tick.
// Animations are always updated even if they are not rendered.
func (ms *models) animate(dt float64) {
	for _, a := range ms.acting {
		a.frame = a.anm.animate(dt, a.frame, a.move, a.pose)
		nextFrame := int(math.Floor(a.frame + 1))
		if nextFrame >= a.nFrames {
			a.frame -= float64(a.nFrames - 1)
		}
	}
}

// setEffect adds CPU particle data to a model entity.
func (ms *models) setEffect(eid eid, mover Mover, maxParticles int) {
	if m := ms.get(eid); m != nil {
		eff := newEffect(m.msh, mover, maxParticles)
		eff.move(m.msh, eff.parts, 0.1)
		ms.effects[eid] = eff
		ms.rebinds[eid] = m
	}
}

// moveParticles updates the particle effects.
// Particles are always updated even if they are not rendered.
func (ms *models) moveParticles(dt float64) {
	for eid, eff := range ms.effects {
		if m, ok := ms.ready[eid]; ok {
			eff.move(m.msh, eff.parts, dt)
			ms.rebinds[eid] = m
		}
	}
}

// Create the backing mesh based on the given label and queue the
// label mesh for a rebind.
func (ms *models) updateLabel(eid eid, l *label, m *model) {
	if l.fnt != nil && m != nil && m.msh != nil && l.str != "" {
		l.w, l.h = l.fnt.setStr(m.msh, l.str, l.wrap)
		ms.rebinds[eid] = m
	}
}

// FUTURE: design a better way to process application clamp requests
//         for textures that are loading.
func (ms *models) processTexture(eid eid, m *model, t *Texture) {
	if clamps, ok := ms.clamps[eid]; ok {
		for _, name := range clamps {
			if t.name == name {
				t.clamp = true
			}
		}
	}
}

// rebind is called from main thread each loop to move asset data to the GPU.
// Assets that have trickled in from the loader are rebound. Each update can
// add assets for binding and this method, run on the main thread, processes them.
func (ms *models) rebind(eng *engine) {
	for eid, m := range ms.rebinds {
		if m.shd != nil && m.shd.program == 0 {
			if err := eng.bind(m.shd); err != nil {
				log.Printf("Bind shader %s failed: %s", m.shd.name, err)
				return // dev error - asset should be bindable.
			}
		}
		if m.msh != nil && m.msh.rebind {
			if err := m.msh.bind(eng); err != nil {
				log.Printf("Bind mesh %s failed: %s", m.msh.name, err)
			}
		}
		for _, t := range m.texs {
			if t.rebind {
				if err := t.bind(eng); err != nil {
					log.Printf("Bind texture %s failed : %s", t.name, err)
				}
			}
			if t.clamp {
				eng.clampTex(t.tid)
				t.clamp = false // HACK: only do once.
			}
		}
		// model removed once all outstanding rebinds are completed.
		delete(ms.rebinds, eid)
	}
}

// draw populates the render.Draw data depending on the model data
// for this entity.
func (ms *models) draw(eid eid, m *model, d *render.Draw) {
	if _, ok := ms.effects[eid]; ok || m.isEffect {
		m.drawEffect(d)
	} else {
		m.draw(d, nil)

		// Set animation pose data for actors.
		if a, ok := ms.acting[eid]; ok && len(a.pose) > 0 {
			d.SetPose(a.pose)
		}
	}
}

// dispose of the model, removing it from all of the maps.
// There is no easy way of knowing when to delete the related assets.
// Leave that to the application.
func (ms *models) dispose(eid eid) {
	delete(ms.all, eid)
	delete(ms.loading, eid)
	delete(ms.ready, eid)
	delete(ms.rebinds, eid)
	delete(ms.actors, eid)
	delete(ms.acting, eid)
	delete(ms.effects, eid)
	delete(ms.labels, eid)
	delete(ms.clamps, eid)
}

// stats returns the number of all models. Used by profile.go.
func (ms *models) stats() (models int) { return len(ms.all) }
