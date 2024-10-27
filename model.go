// Copyright © 2015-2024 Galvanized Logic Inc.

package vu

// model.go contains code related to renderable models.

import (
	"fmt"
	"image"
	"log/slog"
	"math"
	"strings"

	"github.com/gazed/vu/load"
	"github.com/gazed/vu/render"
)

// AddModel adds a new transform with a model component to the given entity.
func (e *Entity) AddModel(assets ...string) (me *Entity) {
	me = e.addPart() // add a transform node for the model.
	if mod := me.app.models.create(me); mod != nil {
		mod.req = strings.Join(assets, ",")
		mod.getAssets(me, assets...)
	}
	return me
}

// AddInstancedModel adds a model where the immediate children are
// instances of the parent model. The parent model will be rendered
// for each childs transform data.
func (e *Entity) AddInstancedModel(assets ...string) (me *Entity) {
	me = e.addPart() // add a transform node
	if mod := me.app.models.create(me); mod != nil {
		mod.getAssets(me, assets...)
		mod.isInstanced = true
		mod.instanceCount = 0 // until SetInstanceData is called
	}
	return me
}

// SetInstanceData sets the instance data for an instanced model.
func (e *Entity) SetInstanceData(eng *Engine, count uint32, data []load.Buffer) (me *Entity) {
	if mod := e.app.models.get(e.eid); mod != nil && mod.isInstanced {
		var err error
		mod.instanceID, err = eng.rc.LoadInstanceData(data)
		if err != nil {
			slog.Error("SetInstanceData", "error", err)
		}
		mod.instanceCount = count
		return e
	}
	slog.Error("SetInstanceData needs AddInstancedModel", "eid", e.eid)
	return e
}

// SetColor sets the solid color for this model - not the texture
// color information. This affects shaders like pbr0 and label
// that use model uniform "color" The color is passed per object
// instance in the shader push constants.
//
// Depends on Entity.AddModel.
func (e *Entity) SetColor(r, g, b, a float64) *Entity {
	if m := e.app.models.get(e.eid); m != nil {
		if m.mat == nil {
			m.mat = newMaterial(fmt.Sprintf("mat%d", e.eid)) // fake name
		}
		// set or replace the existing material with the values provides.
		m.mat.color = rgba{float32(r), float32(g), float32(b), float32(a)}
		return e
	}
	slog.Error("SetMaterial needs AddModel", "eid", e.eid)
	return e
}

// SetMetallicRoughness sets the PBR material attributes for this model,
// not the texture material information. This affects pbr0 or pbr1 shaders.
// The PBR material is passed per object instance in the shader push constants.
//
// Depends on Entity.AddModel.
func (e *Entity) SetMetallicRoughness(metallic bool, roughness float64) *Entity {
	if m := e.app.models.get(e.eid); m != nil {
		if m.mat == nil {
			m.mat = newMaterial(fmt.Sprintf("mat%d", e.eid)) // fake name
		}
		m.mat.roughness = float32(roughness)
		m.mat.metallic = 0.0
		if metallic {
			m.mat.metallic = 1.0
		}
		return e
	}
	slog.Error("SetMaterial needs AddModel", "eid", e.eid)
	return e
}

// SetModelUniform sets data for the given uniform. The uniform data is
// passed to the shader.
func (e *Entity) SetModelUniform(uniform string, data interface{}) *Entity {
	if m := e.app.models.get(e.eid); m != nil {
		switch uniform {
		case "args4":
			if v, ok := data.([]float32); ok && len(v) == 4 {
				m.uniforms[load.ARGS4] = render.V4S32ToBytes(v[0], v[1], v[2], v[3], m.uniforms[load.ARGS4])
			}
		case "args16":
			if v, ok := data.([]float64); ok && len(v) == 16 {
				m.uniforms[load.ARGS16] = render.V16ToBytes(v, m.uniforms[load.ARGS16])
			}
		default:
			// FUTURE    : add uniforms as needed by shaders.
			// FAR FUTURE: data drive the uniforms based on shader reflection.
			slog.Error("unsupported uniform", "uniform", uniform)
		}
		return e
	}
	slog.Error("SetModelUniform needs AddModel", "eid", e.eid)
	return e
}

// AddUpdatableTexture uploads 2 textures to the GPU. One is for rendering
// and one is for updating.
func (e *Entity) AddUpdatableTexture(eng *Engine, name string, img *image.NRGBA) *Entity {
	if mod := e.app.models.get(e.eid); mod != nil {
		if len(mod.updatable) > 0 {
			slog.Error("AddUpdatableTexture already set", "eid", e.eid)
			return e
		}

		// create 2 texture assets.
		opaque := img.Opaque()
		t1 := newTexture(name + "_a")
		t1.opaque = opaque
		t2 := newTexture(name + "_b")
		t2.opaque = opaque
		mod.updatable = []*texture{t1, t2}

		// upload the initial texture to the GPU
		var err error
		idata := &load.ImageData{
			Width:  uint32(img.Bounds().Size().X),
			Height: uint32(img.Bounds().Size().Y),
			Pixels: []byte(img.Pix),
			Opaque: opaque,
		}
		t1.tid, err = eng.rc.LoadTexture(idata)
		if err != nil {
			slog.Error("AddUpdatableTexture upload1", "err", err)
			return e
		}
		slog.Debug("model", "asset", "tex:"+t1.label(), "tid", t1.tid, "opaque", t1.opaque)
		t2.tid, err = eng.rc.LoadTexture(idata)
		if err != nil {
			slog.Error("AddUpdatableTexture upload2", "err", err)
			return e
		}
		slog.Debug("model", "asset", "tex:"+t2.label(), "tid", t2.tid, "opaque", t2.opaque)

		// TODO fake this
		// m.samplerMap[uniform] = name // remember uniform to texture mapping.
		mod.samplerMap["color"] = t1.label()

		// m.texs only takes one of the 2 textures.
		mod.texs = []*texture{t1}
		return e
	}
	slog.Error("AddUpdatableTexture needs AddModel", "eid", e.eid)
	return e
}

// UpdateTexture uploads the image to the updatable texture and then
// swaps the updatable texture with the render texture.
func (e *Entity) UpdateTexture(eng *Engine, img *image.NRGBA) *Entity {
	if mod := e.app.models.get(e.eid); mod != nil {
		if len(mod.updatable) != 2 {
			slog.Error("UpdateTexture not set", "eid", e.eid)
			return e
		}

		// update the uploadable texture,
		t2 := mod.updatable[1] // updatable is always second
		idata := &load.ImageData{
			Width:  uint32(img.Bounds().Size().X),
			Height: uint32(img.Bounds().Size().Y),
			Pixels: []byte(img.Pix),
			Opaque: img.Opaque(),
		}
		if err := eng.rc.UpdateTexture(t2.tid, idata); err != nil {
			slog.Error("UpdateTexture update", "err", err)
			return e
		}

		// swap textures and render the recently updated texture.
		mod.updatable[0], mod.updatable[1] = mod.updatable[1], mod.updatable[0]
		mod.texs[0] = mod.updatable[0] // always render the first after swap.
		mod.samplerMap["color"] = mod.texs[0].label()
		return e
	}
	slog.Error("UpdateTexture needs AddModel", "eid", e.eid)
	return e
}

// SetLayer helps to order draws - normally 2D UI elements.
// Layer values are 0 (first) to 15 (last). Normally packets
// are drawn in the creation order, but this allows specific ordering.
//
// Depends on Entity.AddModel.
func (e *Entity) SetLayer(layer uint8) *Entity {
	if m := e.app.models.get(e.eid); m != nil {
		m.layer = layer
		return e
	}
	slog.Error("SetLayer needs AddModel", "eid", e.eid)
	return e
}

// FUTURE: AddEffect(...) generate quad for particle effects in geometry stage.

// =============================================================================
// model data

// modelType differentiates models based on their data.
// Each type will have some basic model data, like a mesh or texture,
// plus extra data to get a desired render output.
type modelType int

const (
	basicModel  = iota // standard renderable 2D or 3D model.
	labelModel         // standard model + label/font data.
	actorModel         // standard model + animation joint/bone data.
	effectModel        // standard model + particle effect data.
)

// model transforms groups of application requested assets into
// render draw data. It also provides a consistent API for application
// access and modification of render data.
//
// Generally expected to be accessed through wrapper classes such as
// Model, Particle, Actor, and Label.
type model struct {
	req    string  // original asset request string used for debugging.
	shader *shader // Loaded shader.
	mesh   *mesh   // Mandatory vertex data.

	// textures are mapped to sampler uniforms.
	samplerMap map[string]string // mapping of samplers to textures
	texs       []*texture        // texture assets.
	mat        *material         // material parameters.
	updatable  []*texture        // two updatable textures when set.

	// fntAID is needed for static text labels.
	fntAID aid

	// specific model type data.
	mtype modelType // indicates expected model data.
	label *label    // for a labelModel

	// true if this model will be rendered at each of
	// its child transforms.
	isInstanced   bool   // default false.
	instanceCount uint32 // default false.
	instanceID    uint32 // render instance data ID.

	// FUTURE
	// anim   *actor  // set for an animated model
	// effect *effect // set for a particle effect

	// generic uniforms set the app and passed to the shader.
	uniforms map[load.PacketUniform][]byte

	// packet bucket sort values.
	tocam float64 // distance to camera helps with 3D render order.
	layer uint8   // draw layer 0-15
}

// newModel initializes the data structures and default uniforms.
func newModel(mt modelType) *model {
	return &model{
		mtype:      mt,
		samplerMap: map[string]string{},
		uniforms:   map[load.PacketUniform][]byte{},
	}
}

// getAssets for the current model.
func (m *model) getAssets(me *Entity, assets ...string) {
	for _, attribute := range assets {
		attr := strings.Split(attribute, ":")
		switch len(attr) {
		case 2:
			// most asset hav two fields "asset_type:asset_name"
			name := attr[1]
			switch attr[0] {
			case "msh":
				me.app.ld.getAsset(assetID(msh, name), me.eid, me.app.models.assetLoaded)
			case "mat":
				me.app.ld.getAsset(assetID(mat, name), me.eid, me.app.models.assetLoaded)
			case "shd":
				me.app.ld.getAsset(assetID(shd, name), me.eid, me.app.models.assetLoaded)
			case "fnt":
				m.fntAID = assetID(fnt, name)
				me.app.ld.getAsset(m.fntAID, me.eid, me.app.models.assetLoaded)
			case "anm":
				// TODO get animation bone data from GLB files.
			default:
				slog.Error("undefined model asset", "attr", attr[0], "name", name, "eid", me.eid)
			}
		case 3:
			// textures have three fields  "asset_type:uniform_sampler_name:asset_name"
			uniform := attr[1]
			name := attr[2]
			switch attr[0] {
			case "tex":
				me.app.ld.getAsset(assetID(tex, name), me.eid, me.app.models.assetLoaded)
				m.samplerMap[uniform] = name // remember uniform to texture mapping.
			default:
				slog.Error("undefined model asset", "attr", attr[0], "name", name, "eid", me.eid)
				continue
			}
		default:
			slog.Error("undefined model asset", "attr", attribute, "eid", me.eid)
			continue
		}
	}
}

// addAsset adds the asset to the model.
func (m *model) addAsset(a asset) {
	switch la := a.(type) {
	case *mesh:
		m.mesh = la
	case *material:
		if m.mat == nil {
			// set materials not already set by app.
			m.mat = la
		}
	case *texture:
		// textures are added in the order they are loaded.
		// They will have to need to match the order of the sampler
		// uniforms from the shader config when they are used for rendering.
		m.texs = append(m.texs, la)
	case *font:
		if m.mtype == labelModel && m.label != nil {
			m.label.fnt = la
		}
	case *shader:
		m.shader = la
	default:
		slog.Error("unexepected model asset", "name", a.label())
	}
}

// fillPacket populates a render.Packet for this model returning
// false if required shader information was missing.
func (m *model) fillPacket(packet *render.Packet, pov *pov, cam *Camera) error {
	if m.shader == nil {
		return fmt.Errorf("shader not loaded: %s", m.req)
	}
	packet.ShaderID = m.shader.sid // GPU shader reference

	// check if the model mesh has the necessary data.
	if m.mesh == nil {
		return fmt.Errorf("mesh not loaded: %s", m.req)
	}
	packet.MeshID = m.mesh.mid // GPU mesh reference.

	// check specific needs for different model types.
	switch m.mtype {
	case basicModel:
	case labelModel:
		if m.label == nil || m.label.w <= 0 || len(m.texs) != 1 {
			return fmt.Errorf("label not loaded: %s", m.req)
		}
	case actorModel:
		// FUTURE check animation data.
	case effectModel:
		// FUTURE check particle effect data.
	}

	// handle instanced models where a single mesh is drawn multiple times.
	packet.IsInstanced = false
	if m.isInstanced {
		if m.instanceCount <= 0 {
			return fmt.Errorf("instance data not loaded: %s", m.req)
		}
		packet.IsInstanced = true
		packet.InstanceID = m.instanceID
		packet.InstanceCount = m.instanceCount
	}

	// FUTURE: debug validation that the render layer has the uploaded
	// vertex data for the attributes, ie: the m.mesh.mid references
	// vertex data in the render context and each shader attribute
	// should have a non-zero count for the matching vertex data.
	// Needs code that has access to the render context, ie:
	//   for i := range m.shader.config.Attrs {
	//   	  attr := &m.shader.config.Attrs[i]
	//      eng.rc.HasVertexData(m.mesh.mid, attr.AttrType)
	//   }

	// expect one texture for each sampler. Mismatches happen if:
	// - the texture has not yet loaded.
	// - the app forgot to add a texture to the model
	// - the app added an unnecessary texture was added to a model.
	samplers := m.shader.config.GetSamplerUniforms()
	if len(samplers) != len(m.texs) {
		return fmt.Errorf("texture data not loaded: %s", m.req)
	}

	// add the textures to the packet in the same order as the samplers.
	packet.TextureIDs = packet.TextureIDs[:0] // GPU texture references.
	for _, u := range samplers {
		tstr, ok := m.samplerMap[u.Name] // expect to find a matching uniform.
		if !ok {
			return fmt.Errorf("waiting for textures: %s", m.req)
		}

		// find the texture for this uniform.
		// It must exist since the name was found in the samplerMap.
		found := false
		for _, t := range m.texs {
			if tstr == t.label() {
				found = true
				packet.TextureIDs = append(packet.TextureIDs, t.tid)
				break
			}
		}
		if !found {
			return fmt.Errorf("fix texture map typo: %s", m.req)
		}
	}

	// set the model uniform data expected by the shader.
	// Check that data is available for each uniforms, excluding samplers.
	uniforms := m.shader.config.Uniforms
	for i := range uniforms {
		u := &uniforms[i]
		if u.DataType != load.DataType_SAMPLER && u.Scope == load.ModelScope {
			switch u.PacketUID {
			case load.MODEL:
				packet.Uniforms[load.MODEL] = render.M4ToBytes(pov.mm, packet.Uniforms[load.MODEL])
			case load.SCALE:
				sx, sy, sz := pov.scale()
				packet.Uniforms[load.SCALE] = render.V4SToBytes(sx, sy, sz, 0, packet.Uniforms[load.SCALE])
			case load.COLOR, load.MATERIAL:
				if m.mat == nil {
					return fmt.Errorf("waiting on materials: %s", m.req)
				}
				// Set the model material uniform data.
				r, g, b, a := m.mat.color.r, m.mat.color.g, m.mat.color.b, m.mat.color.a
				packet.Uniforms[load.COLOR] = render.V4S32ToBytes(r, g, b, a, packet.Uniforms[load.COLOR])
				metal, rough := m.mat.metallic, m.mat.roughness
				packet.Uniforms[load.MATERIAL] = render.V4S32ToBytes(metal, rough, 0, 0, packet.Uniforms[load.MATERIAL])
			default:
				// basic uniforms are set using SetModelUniform.
				data, ok := m.uniforms[u.PacketUID]
				if !ok {
					return fmt.Errorf("waiting on uniforms: %s", m.req)
				}
				packet.Uniforms[u.PacketUID] = packet.Uniforms[u.PacketUID][:0]
				packet.Uniforms[u.PacketUID] = append(packet.Uniforms[u.PacketUID], data...)
			}
		}
	}

	// add the eid to help debug packets.
	packet.Tag = uint32(pov.eid) // Use eid for debugging draw calls.

	// set the render packet sorting information.
	packet.Bucket = setBucketType(packet.Bucket, drawOpaque)
	if m.isTransparent() {
		packet.Bucket = setBucketType(packet.Bucket, drawTransparent)
	}
	packet.Bucket = setBucketShader(packet.Bucket, m.shader.sid)
	packet.Bucket = setBucketDistance(packet.Bucket, math.MaxFloat64-m.tocam)
	packet.Bucket = setBucketLayer(packet.Bucket, m.layer)
	return nil // model has all information needed to render.
}

// isTransparent returns true if the model is transparent.
// This is either a property of its base color texture or
// its material alpha value.
func (m *model) isTransparent() bool {
	if m.mat != nil && m.mat.color.a < 1.0 {
		return true
	}
	for _, t := range m.texs {
		if !t.opaque {
			return true
		}
	}
	return false
}

// =============================================================================
// models is the component manager for model data.
type models struct {
	list map[eID]*model // All model objects.
}

// newModels creates the render model component manager.
// Expected to be called once on startup.
func newModels() *models {
	ms := &models{}
	ms.list = map[eID]*model{} // any model in any state.
	return ms
}

// create and track a new model.
func (ms *models) create(e *Entity) *model {
	if _, ok := ms.list[e.eid]; ok {
		return nil // model already exists.
	}
	m := newModel(basicModel)
	ms.list[e.eid] = m
	return m
}

// createLabel creates model data and label data for the given entity.
func (ms *models) createLabel(s string, wrap int, e *Entity) *model {
	if _, ok := ms.list[e.eid]; ok {
		return nil // model already exists.
	}
	m := newModel(labelModel)
	ms.list[e.eid] = m
	m.label = &label{str: s, wrap: wrap}

	// create default white color for the label.
	m.mat = newMaterial(fmt.Sprintf("mat%d", e.eid)) // fake name
	m.mat.color = rgba{1, 1, 1, 1}
	return m
}

// assetsLoaded is called from loader when a model asset
// has finished loading.
func (ms *models) assetLoaded(eid eID, a asset) {
	m := ms.get(eid)
	if m == nil {
		// The model may have been disposed before it finished loading.
		slog.Warn("no model for asset", "eid", eid)
		return
	}
	m.addAsset(a)
}

// get the model for the given entity.
func (ms *models) get(eid eID) *model { return ms.list[eid] }

// getReady returns a model that has all of its assets,
// meaning it can be rendered or have its audio played.
// func (ms *models) getReady(eid eID) *model { return ms.ready[eid] }

// dispose of the model, removing it from all of the maps.
// There is no easy way of knowing when to delete the related assets.
// Leave that to the application.
func (ms *models) dispose(eid eID) {
	delete(ms.list, eid)
}
