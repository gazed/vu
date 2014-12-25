// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gazed/vu/audio"
	"github.com/gazed/vu/load"
	"github.com/gazed/vu/render"
)

// Asset ensures that resources get loaded and are ready for use.
// As a resource manager, asset mediates between data loading done by
// vu/load and asset consumption by subsystems like vu/render and vu/audio.
//
// Expected usage is to create one instance of assets at startup since assets
// caches loaded data.
type assets struct {
	ld load.Loader     // Data load and cache subsystem.
	gc render.Renderer // Graphics subsystem injected on creation.
	ac audio.Audio     // Audio subsystem injected on creation.
	d  *depot          // asset cache.
}

// newAssets is expected to be called once during engine initialization.
func newAssets(ac audio.Audio, gc render.Renderer) *assets {
	a := &assets{}
	a.ld = load.NewLoader()
	a.gc = gc
	a.ac = ac
	a.d = newDepot()
	return a
}

// getFont fetches from the asset cache, lazy loading and caching if necessary.
// An error is logged, and nil is returned if the asset can not be located.
func (a *assets) getFont(name string) *font {
	f := newFont(name)
	data := asset(f)
	if err := a.d.fetch(fnt, &data); err == nil {
		return data.(*font)
	}
	if err := a.loadFont(f); err == nil {
		a.d.cache(fnt, f)
		return f
	}
	log.Printf("assets.getFont: could not fetch %s", name)
	return nil
}

// loadFont transfers data loaded from disk to the render object.
func (a *assets) loadFont(f *font) error {
	if fnt, err := a.ld.Fnt(f.Name()); err == nil {
		f.SetSize(fnt.W, fnt.H)
		for _, ch := range fnt.Chars {
			f.AddChar(ch.Char, ch.X, ch.Y, ch.W, ch.H, ch.Xo, ch.Yo, ch.Xa)
		}
	} else {
		return fmt.Errorf("assets.loadFont: could not load %s %s", f.Name(), err)
	}
	return nil
}

// getShader fetches from the asset cache, lazy loading and caching if necessary.
// An error is logged, and nil is returned if the asset can not be located.
func (a *assets) getShader(name string) render.Shader {
	var err error
	s := a.gc.NewShader(name)
	data := s.(asset)
	if err = a.d.fetch(shd, &data); err == nil {
		return data.(render.Shader)
	}
	if err = a.loadShader(s); err == nil {
		a.d.cache(shd, s)
		return s
	}
	log.Printf("assets.getShader: could not fetch %s %s", name, err)
	return nil
}
func (a *assets) remShader(s render.Shader) {
	a.d.remove(shd, s)
}

// loadShader transfers data loaded from disk to the render object.
func (a *assets) loadShader(s render.Shader) error {

	// first look for .vsh, .fsh disk files.
	vsrc, verr := a.ld.Vsh(s.Name())
	fsrc, ferr := a.ld.Fsh(s.Name())
	if verr == nil && ferr == nil {
		s.SetSource(vsrc, fsrc)
		return nil
	}

	// next look for a pre-defined engine shader.
	if vsrc, fsrc = s.Lib(); len(vsrc) > 0 && len(fsrc) > 0 {
		s.SetSource(vsrc, fsrc)
		return nil
	}
	return fmt.Errorf("Could not find shader %s", s.Name())
}

// getTexture fetches from the asset cache, lazy loading and caching if necessary.
// An error is logged, and nil is returned if the asset can not be located.
func (a *assets) getTexture(name string) render.Texture {
	var err error
	t := a.gc.NewTexture(name)
	data := t.(asset)
	if err = a.d.fetch(tex, &data); err == nil {
		return data.(render.Texture)
	}
	if err = a.loadTexture(t); err == nil {
		a.d.cache(tex, t)
		return t
	}
	log.Printf("assets.getTexture: could not fetch %s %s", name, err)
	return nil
}
func (a *assets) remTexture(t render.Texture) {
	a.d.remove(tex, t)
}

// loadTexture transfers data loaded from disk to the render object.
func (a *assets) loadTexture(t render.Texture) error {
	if img, err := a.ld.Png(t.Name()); err == nil {
		t.Set(img)
	} else {
		return fmt.Errorf("assets.loadTexture: could not load %s %s", t.Name(), err)
	}
	return nil
}

// getMaterial fetches from the asset cache, lazy loading and caching if necessary.
// An error is logged, and nil is returned if the asset can not be located.
func (a *assets) getMaterial(name string) *material {
	m := newMaterial(name)
	data := asset(m)
	if err := a.d.fetch(mat, &data); err == nil {
		return data.(*material)
	}
	if err := a.loadMaterial(m); err == nil {
		a.d.cache(mat, m)
		return m
	}
	log.Printf("assets.getMaterial: could not fetch %s", name)
	return nil
}

// loadMaterial transfers data loaded from disk to the render object.
func (a *assets) loadMaterial(m *material) error {
	if mtl, err := a.ld.Mtl(m.Name()); err == nil {
		kd := &rgb{mtl.KdR, mtl.KdG, mtl.KdB}
		ka := &rgb{mtl.KaR, mtl.KaG, mtl.KaB}
		ks := &rgb{mtl.KsR, mtl.KsG, mtl.KsB}
		m.SetMaterial(kd, ka, ks, mtl.Tr)
	} else {
		return fmt.Errorf("assets.loadMaterial: could not load %s %s", m.Name(), err)
	}
	return nil
}

// getMesh fetches from the asset cache, lazy loading and caching if necessary.
// An error is logged, and nil is returned if the asset can not be located.
func (a *assets) getMesh(name string) (model render.Mesh) {
	var err error
	m := a.gc.NewMesh(name)
	data := m.(asset)
	if err = a.d.fetch(msh, &data); err == nil {
		return data.(render.Mesh)
	}
	if err = a.loadMesh(m); err == nil {
		a.d.cache(msh, m)
		return m
	}
	log.Printf("assets.getMesh: could not fetch %s %s", name, err)
	return nil
}

// remMesh clears the given mesh from the cache. Expected to be called
// once the mesh resource is no longer referenced.
func (a *assets) remMesh(m render.Mesh) {
	a.d.remove(msh, m)
}

// newMesh creates an empty mesh whose vertex and buffer data is
// expected to be generated later.
func (a *assets) newMesh(name string) render.Mesh {
	return a.gc.NewMesh(name)
}

// loadMesh transfers data loaded from disk to the render object.
func (a *assets) loadMesh(m render.Mesh) error {
	if data, err := a.ld.Obj(m.Name()); err == nil && len(data) > 0 {
		if len(data[0].V) <= 0 || len(data[0].F) <= 0 {
			return fmt.Errorf("Minimally need vertex and face data for %s", m.Name())
		}
		m.InitData(0, 3, render.STATIC, false).SetData(0, data[0].V)
		if len(data[0].N) > 0 {
			m.InitData(1, 3, render.STATIC, false).SetData(1, data[0].N)
		}
		if len(data[0].T) > 0 {
			m.InitData(2, 2, render.STATIC, false).SetData(2, data[0].T)
		}
		m.InitFaces(render.STATIC).SetFaces(data[0].F)
	} else {
		return fmt.Errorf("assets.loadMesh: could not load %s %s", m.Name(), err)
	}
	return nil
}

// newModel creates a new render Model.
func (a *assets) newModel(s render.Shader) render.Model { return a.gc.NewModel(s) }

// getModel loads a complete model from disk.
// Model's are composite objects created from separately cached items.
// The first model of each type is cached in order to reference the components
// needed to build other model instances. Cached models should never be used
// directly as they contain per instance data.
func (a *assets) getModel(name string, m render.Model) render.Model {
	var err error
	if err, data := a.d.fetchModel(name); err == nil {
		model := data.(render.Model) // get the reference model.
		m.SetMesh(model.Mesh())      // reuse the mesh
		if model.Animation() != nil {
			m.SetAnimation(model.Animation()) // reuse the animation.
		}
		for _, t := range model.Textures() { // reuse the textures.
			m.AddTexture(t)
		}
		return m
	}
	if err = a.loadModel(name, m); err == nil { // create and cache textures.
		a.d.cache(msh, m.Mesh()) // cache the mesh
		if m.Animation() != nil {
			a.d.cache(anm, m.Animation()) // cache the animation.
		}
		a.d.cache(mod, m) // cache the reference model
		return m
	}
	log.Printf("assets.getModel: could not fetch %s %s", name, err)
	return nil
}

// loadModel loads the binary IQM file before the text based IQE file.
// The loaded textures are cached, but the mesh and animation are not.
func (a *assets) loadModel(name string, m render.Model) (err error) {
	var iqd *load.IqData
	if iqd, err = a.ld.Iqm(name); err != nil {
		if iqd, err = a.ld.Iqe(name); err != nil {
			return fmt.Errorf("assets.loadModel: could not load %s %s", m.Name(), err)
		}
	}

	// Use the loaded data to initialize a render.Model
	// Vertex position data and face data must be present.
	msh := a.gc.NewMesh(name)
	msh.InitData(0, 3, render.STATIC, false).SetData(0, iqd.V)
	msh.InitFaces(render.STATIC).SetFaces(iqd.F)

	// load the optional vertex buffers.
	if len(iqd.N) > 0 {
		msh.InitData(1, 3, render.STATIC, false).SetData(1, iqd.N)
	}
	if len(iqd.T) > 0 {
		msh.InitData(2, 2, render.STATIC, false).SetData(2, iqd.T)
	}
	if len(iqd.B) > 0 {
		msh.InitData(4, 4, render.STATIC, false).SetData(4, iqd.B)
	}
	if len(iqd.W) > 0 {
		msh.InitData(5, 4, render.STATIC, true).SetData(5, iqd.W)
	}

	// store the animation data.
	if len(iqd.Frames) > 0 {
		anim := a.gc.NewAnimation(name)
		movements := []render.Movement{}
		for _, ia := range iqd.Anims {
			movement := render.Movement{ia.Name, int(ia.F0), int(ia.Fn)}
			movements = append(movements, movement)
		}
		anim.SetData(iqd.Frames, iqd.Joints, movements)
		m.SetAnimation(anim)
	}
	m.SetMesh(msh)

	// Get the textures for model. There may be more than one. Replace the
	// desired texture name to one based on the model name.
	for cnt, mtex := range iqd.Textures {
		tname := name + strconv.Itoa(cnt)
		if t := a.getTexture(tname); t != nil {
			m.AddModelTexture(t, mtex.F0, mtex.Fn)
		}
	}
	return nil
}

// remModel clears the given model from the cache. Expected to be called
// once the model resource is no longer referenced.
func (a *assets) remModel(m render.Model) {
	a.d.remove(mod, m)
}

// getSound fetches from the asset cache, lazy loading and caching if necessary.
// An error is logged, and nil is returned if the asset can not be located.
func (a *assets) getSound(name string) audio.Sound {
	var err error
	s := a.ac.NewSound(name)
	data := s.(asset)
	if err = a.d.fetch(snd, &data); err == nil {
		return data.(audio.Sound)
	}
	if err = a.loadSound(s); err == nil {
		if err = s.Bind(); err == nil {
			a.d.cache(snd, s)
			return s
		}
	}
	log.Printf("assets.getSound: could not fetch %s %s", name, err)
	return nil
}

// loadSound transfers data loaded from disk to the render object.
func (a *assets) loadSound(s audio.Sound) error {
	if wh, data, err := a.ld.Wav(s.Name()); err == nil {
		s.SetData(wh.Channels, wh.SampleBits, wh.Frequency, wh.DataSize, data)
	} else {
		return fmt.Errorf("assets.loadSound: could not load %s %s", s.Name(), err)
	}
	return nil
}

// assets
// ============================================================================
// depot

type depot map[int]map[string]interface{}

// newDepot creates a new in-memory cache for loaded items. Expected to be
// called once during application initialization (since a cache works best
// when there is only one instance of it :).
func newDepot() *depot {
	return &depot{
		fnt: make(map[string]interface{}),
		mat: make(map[string]interface{}),
		msh: make(map[string]interface{}),
		shd: make(map[string]interface{}),
		snd: make(map[string]interface{}),
		tex: make(map[string]interface{}),
		anm: make(map[string]interface{}),
		mod: make(map[string]interface{}),
	}
}

// fetch retrieves a previously cached data resource using the given name.
// Fetch expects the resource data to be a pointer to one of the resource
// data types. If found the resource data is copied into the supplied data
// pointer. Otherwise the pointer is unchanged and an error is returned.
func (d *depot) fetch(dataType int, data *asset) (err error) {
	if stored := (*d)[dataType][(*data).Name()]; stored != nil {
		*data, _ = stored.(asset)
		return nil
	}
	return fmt.Errorf("depot.fetch: could not fetch asset.")
}

func (d *depot) fetchModel(name string) (err error, data asset) {
	if stored := (*d)[mod][name]; stored != nil {
		data, _ = stored.(asset)
		return nil, data
	}
	return fmt.Errorf("depot.fetch: could not fetch model asset."), nil
}

// cache an asset based on its type and name. Subsequent calls to cache the
// same asset/name combination are ignored. Cache expects data to be one of
// the valid resource data types and to be uniquely named within its data type.
func (d *depot) cache(dataType int, data asset) {
	if data != nil && data.Name() != "" {
		name := data.Name()
		if _, ok := (*d)[dataType][name]; !ok {
			(*d)[dataType][name] = data
		} else {
			log.Printf("depot.cache: data %s already exists", data.Name())
		}
	} else {
		log.Printf("depot.cache: invalid cache data")
	}
}

// remove an asset based on its type and name.
func (d *depot) remove(dataType int, data asset) {
	if data != nil {
		delete((*d)[dataType], data.Name())
	}
}

// depot
// ============================================================================
// asset

// asset describes any data asset that can uniquely identify itself.
type asset interface {
	Name() string // Unique identifier set on creation.
}

// Data types. Not expected to be used outside of this file.
// Cached data is reusable, except for model's.
const (
	fnt = iota // font
	shd        // shader
	mat        // material
	msh        // mesh
	tex        // texture
	snd        // sound
	anm        // animation
	mod        // model - for reference, don't use directly.
)
