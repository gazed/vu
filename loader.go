// Copyright © 2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/gazed/vu/load"
	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/render"
)

// loader imports and prepares model and noise data for use.
// It is expected to be run as a goroutine and thus needs to
// be initialized with communication channels.
//
// The loader is used internally by the engine to cache and reuse
// imported data across multiple model and noise instances.
type loader struct {
	modLoads []*modelLoad // Collects model load requests.
	texLoads []*texLoad   // Collects texture load requests.

	// loader goroutine communication.
	ld     load.Loader // asset loader.
	cache  cache       // asset cache.
	load   chan msg    // loadModel, loadNoise, bindMesh, bindTexture
	loaded chan msg    // loadModel, loadNoise
	binder chan msg    // machine loop request channel.
}

// newLoader is expected to be called once on startup.
func newLoader(loaded chan msg, binder chan msg) *loader {
	l := &loader{loaded: loaded, binder: binder}
	l.ld = load.NewLoader()
	l.cache = newCache()
	l.load = make(chan msg)
	return l
}

// queueModelLoad is called on the engine processing goroutine.
// The models are loaded and returned by calling loadQueued.
func (l *loader) queueModelLoad(eid uint64, m *model) {
	l.modLoads = append(l.modLoads, &modelLoad{eid: eid, model: m})
}

// queueTextureLoads is called on the engine processing goroutine.
// The textures are loaded and returned by calling loadQueued.
func (l *loader) queueTextureLoads(loads []*texLoad) {
	l.texLoads = append(l.texLoads, loads...)
}

// loadQueued is called on the engine processing goroutine.
// loadQueued sends the load requests off for loading.
// The load requests are batched into small chunks so the loader
// starts returning assets while others are being loaded.
func (l *loader) loadQueued() {
	batchSize := 100
	for len(l.modLoads) > 0 {
		if len(l.modLoads) > batchSize {
			go l.request(l.modLoads[:batchSize])
			l.modLoads = l.modLoads[batchSize:]
		} else {
			go l.request(l.modLoads)
			l.modLoads = []*modelLoad{}
		}
	}
	for len(l.texLoads) > 0 {
		if len(l.texLoads) > batchSize {
			go l.request(l.texLoads[:batchSize])
			l.texLoads = l.texLoads[batchSize:]
		} else {
			go l.request(l.texLoads)
			l.texLoads = []*texLoad{}
		}
	}
}

// runLoader loops forever processing all load requests.
// It is started once on startup as a goroutine.
func (l *loader) runLoader() {
	var req msg
	for {
		req = <-l.load
		switch m := req.(type) {
		case *shutdown:
			return // exit immediately: the app loop is done.
		case *noiseLoad:
			m.noise = l.loadNoise(m.noise)
			go l.returnAsset(m)
		case *releaseData:
			l.release(m.data)
		case []*modelLoad: // multiple model load requests.
			for _, ml := range m {
				if ml.model.shd != nil {
					ml.model = l.loadModel(ml.model)
				} else {
					log.Printf("vu.loader: model for %d was disposed", ml.eid)
					ml.model = nil
				}
			}
			go l.returnAsset(m)

		case []*texLoad: // multiple texture load requests.
			for _, tl := range m {
				tl.tex = l.loadTexture(tl.tex)
			}
			go l.returnAsset(m)
		case nil:
			return // exit immediately: channel closed.
		default:
			log.Printf("loader: unknown msg %t", m)
		}
	}
}

// request is the entry point for all load and unload requests.
// It is expected to be called as a go-routine whereupon it waits
// for the asset loader to process its request.
func (l *loader) request(m msg) { l.load <- m }

// returnAsset funnels all completed loaded assets back to the
// engine loop. The engine loop may be busy, so make this
// method, started as a goroutine, wait instead of the loader.
func (l *loader) returnAsset(asset msg) {
	l.loaded <- asset
}

// loadNoise returns a loaded noise immediately if it is cached.
// Otherwise the noise is returned after it is  loaded and bound.
func (l *loader) loadNoise(n *noise) *noise {
	for index, snd := range n.snds {
		if n.snds[index] = l.loadSound(snd); n.snds[index] == nil {
			return nil
		}
	}
	n.rebind = false
	return n
}

// loadSound returns a loaded sound immediately if it is cached.
// Otherwise the sound is returned after it is loaded and bound.
func (l *loader) loadSound(s *sound) *sound {
	data := asset(s)
	if err := l.cache.fetch(&data); err == nil {
		return data.(*sound) // only initialized stuff is in the cache.
	}

	// Otherwise the sound needs to be loaded and bound.
	if err := l.importSound(s); err != nil {
		log.Printf("Sound load %s: %s", s.name, err)
		return nil // discard load failures
	}
	bindReply := make(chan error)
	l.binder <- &bindData{data: s, reply: bindReply} // request bind.
	if err := <-bindReply; err != nil {              // wait for bind.
		log.Printf("Sound bind %s: %s", s.name, err)
		return nil // discard bind failures
	}
	l.cache.store(s)
	return s
}

// importSound transfers audio data loaded from disk to the sound object.
func (l *loader) importSound(s *sound) error {
	if wh, data, err := l.ld.Wav(s.name); err == nil {
		s.data.Set(wh.Channels, wh.SampleBits, wh.Frequency, wh.DataSize, data)
	} else {
		return fmt.Errorf("loader.loadSound: could not load %s %s", s.label(), err)
	}
	return nil
}

// loadModel returns a loaded model immediately if its parts are cached.
// Otherwise the model is returned after its parts are loaded and bound.
// The entire model must load in order for a valid model to be returned.
// The cache will contain any of the models parts that loaded.
func (l *loader) loadModel(m *model) (loaded *model) {
	if m.shd = l.loadShader(m.shd); m.shd == nil {
		return nil
	}

	// load fonts before mesh to re-bind the mesh
	// after setting the phrase.
	if m.fnt != nil {
		if m.fnt = l.loadFont(m.fnt); m.fnt == nil {
			return nil
		}
		// Lazy create the backing for the phrase using a dynamic mesh.
		// Ie: a mesh that is generated, and not stored in the cache.
		if m.msh == nil {
			m.msh = newMesh("phrase")
			m.msh.generated = true
		} else if !m.msh.generated {
			log.Printf("loader: font on static mesh %s", m.msh.name)
		}
		m.phraseWidth = m.fnt.setPhrase(m.msh, m.phrase)
	}

	// Load meshes and textures for non-animated models.
	if m.anm == nil {
		if m.msh = l.loadMesh(m.msh); m.msh == nil {
			return nil
		}
		for index, tex := range m.texs {
			if m.texs[index] = l.loadTexture(tex); m.texs[index] == nil {
				return nil
			}
		}
	} else {
		// animated model load textures and meshes.
		if m.msh == nil {
			m.msh = newMesh(m.anm.name)
		}
		var texs []*texture
		if m.anm, m.msh, texs = l.loadAnim(m.anm, m.msh); m.anm != nil && m.msh != nil {
			for _, tex := range texs {
				m.texs = append(m.texs, tex)
			}
			m.nFrames = m.anm.maxFrames(0)
			m.pose = make([]lin.M4, len(m.anm.joints))
		} else {
			return nil
		}
	}
	if m.mat != nil {
		if m.mat = l.loadMaterial(m.mat); m.mat == nil {
			return nil
		}
		if m.resetMat {
			m.alpha = m.mat.tr // Copy values so they can be set per model.
			m.kd = m.mat.kd    // ditto
		}
		m.ks = m.mat.ks // Can't currently be overridden on model.
		m.ka = m.mat.ka // ditto
	}
	return m
}

// loadShader returns a loaded shader immediately if it is cached.
// Otherwise the shader is returned after it is  loaded and bound.
func (l *loader) loadShader(s *shader) *shader {
	data := asset(s)
	if err := l.cache.fetch(&data); err == nil {
		s = data.(*shader)
		return s // only initialized stuff is in the cache.
	}

	// Otherwise the shader needs to be loaded and bound.
	if err := l.importShader(s); err != nil {
		log.Printf("Shader load %s: %s", s.name, err)
		return nil // discard load failures
	}
	bindReply := make(chan error)
	l.binder <- &bindData{data: s, reply: bindReply} // request bind.
	if err := <-bindReply; err != nil {              // wait for bind.
		log.Printf("Shader bind %s: %s", s.name, err)
		return nil // discard bind failures
	}
	s.rebind = false
	l.cache.store(s)
	return s
}

// importShader transfers data loaded from disk to the render object.
func (l *loader) importShader(s *shader) error {
	ld := load.NewLoader()

	// first look for .vsh, .fsh disk files.
	vsrc, verr := ld.Vsh(s.name)
	fsrc, ferr := ld.Fsh(s.name)
	if verr == nil && ferr == nil {
		s.setSource(vsrc, fsrc)
		return nil
	}

	// next look for a pre-defined engine shader.
	if sfn, ok := shaderLibrary[s.name]; ok {
		vsrc, fsrc = sfn()
		s.setSource(vsrc, fsrc)
		return nil
	}
	return fmt.Errorf("Could not find shader %s", s.name)
}

// disposeShader clears the given shader from the cache.
// Expected when once the resource is no longer referenced.
func (l *loader) disposeShader(s *shader) { l.cache.remove(s) }

// loadMesh returns a loaded noise immediately if it is cached.
// Otherwise the mesh is returned after it is loaded and bound.
func (l *loader) loadMesh(m *mesh) *mesh {
	if m.generated { // data is generated/changed per update.
		if err := l.bindMesh(m); err != nil {
			log.Printf("Mesh bind %s: %s", m.name, err)
			return nil // discard bind failures
		}
		return m
	}
	data := asset(m)
	if err := l.cache.fetch(&data); err == nil {
		return data.(*mesh)
	}

	// Otherwise the mesh needs to be loaded and bound.
	if err := l.importMesh(m); err != nil {
		log.Printf("Mesh load %s: %s", m.name, err)
		return nil // discard load failures
	}
	if err := l.bindMesh(m); err != nil {
		log.Printf("Mesh bind %s: %s", m.name, err)
		return nil // discard bind failures
	}
	m.rebind = false
	l.cache.store(m)
	return m
}

// bindMesh submits a mesh for binding or rebinding. This transfers
// the mesh data to the GPU.
func (l *loader) bindMesh(m *mesh) error {
	bindReply := make(chan error)
	l.binder <- &bindData{data: m, reply: bindReply} // request bind.
	return <-bindReply                               // wait for bind.
}

// importMesh transfers data loaded from disk to the render object.
func (l *loader) importMesh(m *mesh) error {
	if data, err := l.ld.Obj(m.name); err == nil && len(data) > 0 {
		if len(data[0].V) <= 0 || len(data[0].F) <= 0 {
			return fmt.Errorf("Minimally need vertex and face data for %s", m.name)
		}
		m.initData(0, 3, render.STATIC, false).setData(0, data[0].V)
		if len(data[0].N) > 0 {
			m.initData(1, 3, render.STATIC, false).setData(1, data[0].N)
		}
		if len(data[0].T) > 0 {
			m.initData(2, 2, render.STATIC, false).setData(2, data[0].T)
		}
		m.initFaces(render.STATIC).setFaces(data[0].F)
	} else {
		return fmt.Errorf("assets.loadMesh: could not load %s %s", m.name, err)
	}
	return nil
}

// disposeMesh clears the given mesh from the cache.
// Expected when once the resource is no longer referenced.
func (l *loader) disposeMesh(m *mesh) { l.cache.remove(m) }

// loadTexture returns a loaded texture immediately if it is cached.
// Otherwise the texture is returned after it is loaded and bound.
func (l *loader) loadTexture(t *texture) *texture {
	data := asset(t)
	if err := l.cache.fetch(&data); err == nil {
		return data.(*texture)
	}

	// Otherwise the mesh needs to be loaded and bound.
	if err := l.importTexture(t); err != nil {
		log.Printf("Texture load %s: %s", t.name, err)
		return nil // discard load failures
	}
	bindReply := make(chan error)
	l.binder <- &bindData{data: t, reply: bindReply} // request bind.
	if err := <-bindReply; err != nil {              // wait for bind.
		log.Printf("Texture bind %s: %s", t.name, err)
		return nil // discard bind failures
	}
	t.rebind = false
	l.cache.store(t)
	return t
}

// importTexture transfers data loaded from disk to the render object.
func (l *loader) importTexture(t *texture) error {
	if img, err := l.ld.Png(t.name); err == nil {
		t.set(img)
		return nil
	} else {
		return fmt.Errorf("loader.loadTexture: could not load %s %s", t.name, err)
	}
}

// disposeTexture clears the given texture from the cache.
// Expected when once the resource is no longer referenced.
func (l *loader) disposeTexture(t *texture) { l.cache.remove(t) }

// loadMaterial returns a loaded material immediately if it is cached.
// Otherwise the material is returned after it is loaded and bound.
func (l *loader) loadMaterial(m *material) *material {
	data := asset(m)
	if err := l.cache.fetch(&data); err == nil {
		return data.(*material)
	}

	// Otherwise the mesh needs to be loaded (no binding necessary).
	if err := l.importMaterial(m); err != nil {
		log.Printf("Material load %s: %s", m.name, err)
		return nil // discard load failures
	}
	l.cache.store(m)
	return m
}

// importMaterial transfers data loaded from disk to the render object.
func (l *loader) importMaterial(m *material) error {
	if mtl, err := l.ld.Mtl(m.label()); err == nil {
		kd := &rgb{mtl.KdR, mtl.KdG, mtl.KdB}
		ka := &rgb{mtl.KaR, mtl.KaG, mtl.KaB}
		ks := &rgb{mtl.KsR, mtl.KsG, mtl.KsB}
		m.setMaterial(kd, ka, ks, mtl.Tr)
		return nil
	}
	return fmt.Errorf("loader.loadMaterial: could not load %s", m.name)
}

// disposeMaterial clears the given material from the cache.
// Expected when once the resource is no longer referenced.
func (l *loader) disposeMaterial(m *material) { l.cache.remove(m) }

// loadFont returns a loaded font immediately if it is cached.
// Otherwise the font is returned after it is loaded and bound.
func (l *loader) loadFont(f *font) *font {
	data := asset(f)
	if err := l.cache.fetch(&data); err == nil {
		return data.(*font)
	}

	// Otherwise the font needs to be loaded and stored.
	if err := l.importFont(f); err != nil {
		log.Printf("Font load %s: %s", f.name, err)
		return nil // discard load failures
	}
	l.cache.store(f)
	return f
}

// importFont transfers data loaded from disk to the render object.
func (l *loader) importFont(f *font) error {
	if fnt, err := l.ld.Fnt(f.label()); err == nil {
		f.setSize(fnt.W, fnt.H)
		for _, ch := range fnt.Chars {
			f.addChar(ch.Char, ch.X, ch.Y, ch.W, ch.H, ch.Xo, ch.Yo, ch.Xa)
		}
	} else {
		return fmt.Errorf("assets.loadFont: could not load %s %s", f.label(), err)
	}
	return nil
}

// disposeFont clears the given font from the cache.
// Expected when once the resource is no longer referenced.
func (l *loader) disposeFont(f *font) { l.cache.remove(f) }

// loadAnim loads an animated model from disk. This will
// create multiple components of a model including a mesh,
// textures, and animation data.
func (l *loader) loadAnim(a *animation, m *mesh) (*animation, *mesh, []*texture) {
	data := asset(a)
	if err := l.cache.fetch(&data); err == nil {
		a = data.(*animation) // got the animation.
		data := asset(m)      // now load the mesh.
		if err := l.cache.fetch(&data); err == nil {
			m = data.(*mesh)

			// load all textures based on the animation name.
			texs := []*texture{}
			for cnt := 0; ; cnt++ {
				tname := a.name + strconv.Itoa(cnt)
				t := newTexture(tname)
				data := asset(t) // now load the mesh.
				if err := l.cache.fetch(&data); err == nil {
					texs = append(texs, data.(*texture))
				} else {
					break
				}
			}
			return a, m, texs
		}
	}

	// Otherwise the animation, and mesh, need to be loaded.
	// Textures are loaded, bound, and cached within importAnim.
	var texs []*texture
	var err error
	if texs, err = l.importAnim(a, m); err != nil {
		log.Printf("Animation load %s: %s", a.name, err)
		return nil, nil, nil // discard load failures
	}

	// And the mesh needs to be bound.
	if err := l.bindMesh(m); err != nil {
		log.Printf("Animation bind %s: %s", m.name, err)
		return nil, nil, nil // discard bind failures
	}
	l.cache.store(a)
	l.cache.store(m)
	return a, m, texs
}

// importAnim loads the animation, mesh, and texture for an
// animated model.
func (l *loader) importAnim(a *animation, m *mesh) (texs []*texture, err error) {
	var iqd *load.IqData
	if iqd, err = l.ld.Iqm(a.name); err != nil {
		return nil, err
	}

	// Use the loaded data to initialize a render.Model
	// Vertex position data and face data must be present.
	// All other buffers are optional, but need T, B, W for animation.
	m.initData(0, 3, render.STATIC, false).setData(0, iqd.V)
	m.initFaces(render.STATIC).setFaces(iqd.F)
	if len(iqd.N) > 0 {
		m.initData(1, 3, render.STATIC, false).setData(1, iqd.N)
	}
	if len(iqd.T) > 0 {
		m.initData(2, 2, render.STATIC, false).setData(2, iqd.T)
	}
	if len(iqd.B) > 0 {
		m.initData(4, 4, render.STATIC, false).setData(4, iqd.B)
	}
	if len(iqd.W) > 0 {
		m.initData(5, 4, render.STATIC, true).setData(5, iqd.W)
	}

	// Store the animation data.
	if len(iqd.Frames) > 0 {
		moves := []movement{}
		for _, ia := range iqd.Anims {
			movement := movement{
				name: ia.Name,
				f0:   int(ia.F0),
				fn:   int(ia.Fn),
				rate: float64(ia.Rate)}
			moves = append(moves, movement)
		}
		a.setData(iqd.Frames, iqd.Joints, moves)
	}

	// Get model textures. There may be more than one.
	// Convention: use texture names based on the animation name.
	for cnt, itex := range iqd.Textures {
		tname := a.name + strconv.Itoa(cnt)
		tex := newTexture(tname)
		if tex = l.loadTexture(tex); tex != nil {
			tex.fn, tex.f0 = itex.Fn, itex.F0
			texs = append(texs, tex)
		}
	}
	return texs, nil
}

// release is called when the cached data is no
// longer needed and can be discarded entirely.
func (l *loader) release(data interface{}) {
	switch d := data.(type) {
	case *mesh:
		l.cache.remove(d)
	case *font:
		l.cache.remove(d)
	case *shader:
		l.cache.remove(d)
	case *material:
		l.cache.remove(d)
	case *texture:
		l.cache.remove(d)
	case *animation:
		l.cache.remove(d)
	case *sound:
		l.cache.remove(d)
	default:
		log.Printf("loader.dispose unknown %t", d)
	}
}

// loader
// ============================================================================
// cache

type cache map[uint64]interface{}

// newCache creates a new in-memory cache for loaded items. Expected to be
// called once during application initialization (since a cache works best
// when there is only one instance of it :).
func newCache() cache {
	return make(map[uint64]interface{})
}

// fetch retrieves a previously cached data resource using the given name.
// Fetch expects the resource data to be a pointer to one of the resource
// data types. If found the resource data is copied into the supplied data
// pointer. Otherwise the pointer is unchanged and an error is returned.
func (c cache) fetch(data *asset) (err error) {
	if stored := c[(*data).aid()]; stored != nil {
		*data, _ = stored.(asset)
		return nil
	}
	return fmt.Errorf("cache.fetch: could not fetch asset.")
}

// store an asset based on its type and name. Subsequent calls to cache the
// same asset/name combination are ignored. Cache expects data to be one of
// the valid resource data types and to be uniquely named within its data type.
func (c cache) store(data asset) {
	if data == nil {
		log.Printf("cache.store: invalid cache data %s", data.label())
		return
	}
	if _, ok := c[data.aid()]; !ok {
		c[data.aid()] = data
	} else {
		log.Printf("cache.store: data %s already exists", data.label())
	}
}

// remove an asset based on its type and name.
func (c cache) remove(data asset) {
	if data != nil {
		delete(c, data.aid())
	}
}

// ============================================================================
// asset

// asset describes any data asset that can uniquely identify itself.
type asset interface {
	label() string // Unique identifier set on creation.
	aid() uint64   // Data type and name combined.
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
)

// =============================================================================
// utility methods.

// stringHash turns a string into a number.
// Algorithm based on java String.hashCode().
//     s[0]*31^(n-1) + s[1]*31^(n-2) + ... + s[n-1]
func stringHash(s string) (h uint64) {
	bytes := []byte(s)
	n := len(bytes)
	hash := uint32(0)
	for index, b := range bytes {
		hash += uint32(b) * uint32(math.Pow(31, float64(n-index)))
	}
	return uint64(hash)
}
