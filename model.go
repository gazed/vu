// Copyright © 2015-2024 Galvanized Logic Inc.

package vu

// model.go contains code related to renderable models.

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/gazed/vu/load"
	"github.com/gazed/vu/render"
)

// AddModel adds a new transform with a model component to the given entity.
func (e *Entity) AddModel(assets ...string) (me *Entity) {
	me = e.addPart() // add a transform node for the model.
	if mod := me.app.models.create(me); mod != nil {
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

// FUTURE AddEffect(...) generate quad for particle effects in geometry stage.

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
	shader *shader // Loaded shader.
	mesh   *mesh   // Mandatory vertex data.

	// textures are mapped to sampler uniforms.
	samplerMap map[string]string // mapping of samplers to textures
	texs       []*texture        // texture assets.
	mat        *material         // material parameters.

	// specific model type data.
	mtype modelType // indicates expected model data.
	label *label    // for a labelModel

	// true if this model will be rendered at each of
	// its child transforms.
	isInstanced   bool   // default false.
	instanceCount uint32 // default false.
	instanceID    uint32 // render instance data ID.

	// TODO anim   *actor  // set for an animated model
	// TODO effect *effect // set for a particle effect

	tocam float64 // distance to camera helps with 3D render order.
}

// newModel initializes the data structures and default uniforms.
func newModel(mt modelType) *model {
	return &model{mtype: mt, samplerMap: map[string]string{}}
}

// canRender returns true if the model has all the assets it needs to render.
func (m *model) canRender() bool {
	if m.shader == nil { // model must have loaded shader.
		return false
	}

	// do we have textures for the required shader samplers
	samplers := m.shader.config.GetSamplerUniforms()
	if len(samplers) != len(m.texs) {
		return false
	}

	// has model vertex data been loaded.
	switch m.mtype {
	case basicModel:
		if m.mesh == nil {
			return false
		}
	case labelModel:
		if m.mesh == nil || m.label == nil {
			return false
		}
		l := m.label
		return len(l.str) > 0 && l.w > 0 && l.h > 0
	case actorModel:
		// FUTURE check animation data.
	case effectModel:
		// FUTURE check particle effect data.
	}

	// instanced models need instance data.
	if m.isInstanced && m.instanceCount <= 0 {
		return false
	}

	// check the shader model level uniforms
	for _, u := range m.shader.config.Uniforms {
		switch u.PacketUID {
		case load.MODEL, load.SCALE:
			// handled already
		case load.COLOR, load.MATERIAL:
			// requires valid mat
			if u.Scope == load.ModelScope && m.mat == nil {
				slog.Warn("model shader requires material", "mesh", m.mesh.name, "shader", m.shader.name)
				return false
			}
		}
	}
	return true
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
				fontAid := assetID(fnt, name)
				me.app.ld.getAsset(fontAid, me.eid, me.app.models.assetLoaded)
				me.app.ld.getLabelMesh(fontAid, me)
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
	ready := m.canRender()
	switch la := a.(type) {
	case *mesh:
		m.mesh = la
	case *material:
		m.mat = la
	case *texture:
		// textures are added in the order they are loaded.
		// They will have to be sorted later to match the
		// order of the sampler uniforms in the shader config.
		// Can't sort here since the shader config might load
		// after the textures.
		m.texs = append(m.texs, la)
	case *font:
		if m.mtype != labelModel || m.label == nil {
			slog.Error("dev: fix non-label loading font data", "font", la.name)
			return
		}
		m.label.fnt = la
	case *shader:
		m.shader = la
	default:
		slog.Error("unexepected model asset", "name", a.label())
	}

	// check if all the assets have been loaded and prep for rendering if ready.
	if m.canRender() != ready {
		m.matchTexturesToSamplers()
		slog.Debug("model ready to render", "shader", m.shader.name, "mesh", m.mesh.name)
	}
}

// matchTexturesToSamplers ensures that the textures are in the
// order expected by the uniform samplers.
func (m *model) matchTexturesToSamplers() {
	orderedTexs := []*texture{}
	for _, u := range m.shader.samplers {
		tstr, ok := m.samplerMap[u.Name] // expect to find a matching uniform.
		if !ok {
			slog.Error("fix typo in application model texture:uniform map")
			return // this needs to be caught and fixed in debug builds.
		}

		// find the texture for this uniform.
		found := false
		for _, t := range m.texs {
			if tstr == t.label() {
				found = true
				orderedTexs = append(orderedTexs, t)
				break
			}
		}
		if !found {
			slog.Error("fix typo in application model texture:name map")
			return // this needs to be caught and fixed in debug builds.
		}
	}
	copy(m.texs, orderedTexs) // copy ordered slice over old.
}

// fillPacket populates a render.Packet for this model.
func (m *model) fillPacket(packet *render.Packet, pov *pov, cam *Camera) {
	packet.ShaderID = m.shader.sid // GPU shader reference
	packet.MeshID = m.mesh.mid     // GPU mesh reference.

	// Rendering hints.
	packet.Tag = uint32(pov.eid) // Use eid for debugging draw calls.
	packet.Bucket = 0            // Used to sort packets. Lower buckets rendered first.

	// copy instanced mesh information into the packet.
	packet.IsInstanced = false
	if m.isInstanced {
		packet.IsInstanced = true
		packet.InstanceID = m.instanceID
		packet.InstanceCount = m.instanceCount
	}

	// copy the ordered textures into the packet.
	packet.TextureIDs = packet.TextureIDs[:0] // GPU texture references.
	for _, tex := range m.texs {
		packet.TextureIDs = append(packet.TextureIDs, tex.tid)
	}

	// Set the model uniform data.
	packet.Data[load.MODEL] = render.M4ToBytes(pov.mm, packet.Data[load.MODEL])

	// Set the model color and material uniform data.
	if m.mat != nil {
		r, g, b, a := m.mat.color.r, m.mat.color.g, m.mat.color.b, m.mat.color.a
		packet.Data[load.COLOR] = render.V4S32ToBytes(r, g, b, a, packet.Data[load.COLOR])

		// Set the model material uniform data.
		metal, rough := m.mat.metallic, m.mat.roughness
		packet.Data[load.MATERIAL] = render.V4S32ToBytes(metal, rough, 0, 0, packet.Data[load.MATERIAL])
	}

	// Set the model material scale uniform data.
	sx, sy, sz := pov.scale()
	packet.Data[load.SCALE] = render.V4SToBytes(sx, sy, sz, 0, packet.Data[load.SCALE])

	// set the render packet sorting information.
	packet.Bucket = setBucketType(packet.Bucket, drawOpaque)
	if m.isTransparent() {
		packet.Bucket = setBucketType(packet.Bucket, drawTransparent)
	}
	packet.Bucket = setBucketShader(packet.Bucket, m.shader.sid)
	packet.Bucket = setBucketDistance(packet.Bucket, m.tocam)
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
	list    map[eID]*model // All model objects.
	loading map[eID]*model // Waiting for assets.
	ready   map[eID]*model // Assets received.
}

// newModels creates the render model component manager.
// Expected to be called once on startup.
func newModels() *models {
	ms := &models{}
	ms.list = map[eID]*model{}    // any model in any state.
	ms.loading = map[eID]*model{} // waiting for initial assets.
	ms.ready = map[eID]*model{}   // assets received.
	return ms
}

// create and track a new model.
func (ms *models) create(e *Entity) *model {
	if _, ok := ms.list[e.eid]; ok {
		return nil // model already exists.
	}
	m := newModel(basicModel)
	ms.list[e.eid] = m
	ms.loading[e.eid] = m
	return m
}

// createLabel creates model data and label data for the given entity.
func (ms *models) createLabel(s string, wrap int, e *Entity) *model {
	if _, ok := ms.list[e.eid]; ok {
		return nil // model already exists.
	}
	m := newModel(labelModel)
	ms.list[e.eid] = m
	ms.loading[e.eid] = m
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
	ms.updateReady(eid, m)
}

// updateReady moves a model to the ready queue
// if all the assets have been loaded.
func (ms *models) updateReady(eid eID, m *model) {
	if m.canRender() {
		delete(ms.loading, eid)
		ms.ready[eid] = m
	}
}

// get loading or loaded models for the given entity.
func (ms *models) get(eid eID) *model { return ms.list[eid] }

// getReady returns a model that has all of its assets,
// meaning it can be rendered or have its audio played.
func (ms *models) getReady(eid eID) *model { return ms.ready[eid] }

// dispose of the model, removing it from all of the maps.
// There is no easy way of knowing when to delete the related assets.
// Leave that to the application.
func (ms *models) dispose(eid eID) {
	delete(ms.list, eid)
	delete(ms.loading, eid)
	delete(ms.ready, eid)
}
