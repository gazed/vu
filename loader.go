// Copyright © 2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/gazed/vu/load"
	"github.com/gazed/vu/render"
)

// loader imports and prepares model and noise data for use.
// It is expected to be run as a goroutine and thus needs to
// be initialized with communication channels.
//
// The loader is used internally by the engine to cache and reuse
// imported data across multiple model and noise instances.
type loader struct {
	loads []*loadReq // Collects asset load requests.

	// loader goroutine communication.
	ld      load.Loader     // asset loader.
	cache   cache           // asset cache.
	control chan msg        // shutdown requests.
	load    chan []*loadReq // asset load requests.
	loaded  chan []*loadReq // loaded asset replies.
	binder  chan msg        // machine loop request channel.
}

// newLoader is expected to be called once on startup.
func newLoader(loaded chan []*loadReq, binder chan msg) *loader {
	l := &loader{loaded: loaded, binder: binder}
	l.ld = load.NewLoader()
	l.cache = newCache()
	l.control = make(chan msg)
	l.load = make(chan []*loadReq)
	return l
}

// queueLoads is called on the engine processing goroutine.
// The models are loaded and returned by calling loadQueued.
func (l *loader) queueLoads(loadRequests []*loadReq) {
	l.loads = append(l.loads, loadRequests...)
}

// loadQueued is called on the engine processing goroutine.
// to send the load requests off for loading. Load requests
// are batched into small chunks so the loader can return
// some assets while others are being loaded.
func (l *loader) loadQueued() {
	batchSize := 100
	for len(l.loads) > 0 {
		if len(l.loads) > batchSize {
			go l.processRequests(l.loads[:batchSize])
			l.loads = l.loads[batchSize:]
		} else {
			go l.processRequests(l.loads)
			l.loads = []*loadReq{}
		}
	}
}

// runLoader loops forever processing all load requests.
// It is started once as a goroutine on engine initialization.
func (l *loader) runLoader() {
	for {
		select {
		case req := <-l.control:
			if req == nil {
				return // exit immediately: channel closed.
			}
			switch m := req.(type) {
			case *shutdown:
				return // exit immediately: the app loop is done.
			case *releaseData:
				l.release(m.data)
			default:
				log.Printf("loader: unknown msg %T", m)
			}
		case requests := <-l.load:
			for _, req := range requests {
				switch a := req.a.(type) {
				case *mesh:
					req.a, req.err = l.loadMesh(a)
				case *texture:
					req.a, req.err = l.loadTexture(a)
				case *shader:
					req.a, req.err = l.loadShader(a)
				case *font:
					req.a, req.err = l.loadFont(a)
				case *animation:
					msh := newMesh(a.name)
					if la, lmsh, ltexs := l.loadAnim(a, msh); la != nil && lmsh != nil {
						req.a = la
						req.msh = lmsh
						req.texs = ltexs
					} else {
						req.a = nil   // return explicit nil for asset interface.
						req.msh = nil // release mesh on fail.
					}
				case *material:
					req.a, req.err = l.loadMaterial(a)
				case *sound:
					req.a, req.err = l.loadSound(a)
				}
			}
			go l.returnAssets(requests)
		}
	}
}

// request is the entry point for all load and unload requests.
// It is expected to be called as a go-routine whereupon it waits
// for the asset loader to process its request.
func (l *loader) processRequests(reqs []*loadReq) { l.load <- reqs }

// returnAssets funnels a group of loaded assets back to the
// engine loop. The engine loop may be busy, so make this
// method, started as a goroutine, wait instead of the loader.
func (l *loader) returnAssets(assets []*loadReq) {
	l.loaded <- assets
}

// loadSound returns a loaded sound immediately if it is cached.
// Otherwise the sound is returned after it is loaded and bound.
func (l *loader) loadSound(s *sound) (*sound, error) {
	data := asset(s)
	if err := l.cache.fetch(&data); err == nil {
		return data.(*sound), nil // only initialized stuff is in the cache.
	}

	// Otherwise the sound needs to be loaded and bound.
	if err := l.importSound(s); err != nil {
		return nil, err
	}
	bindReply := make(chan error)
	l.binder <- &bindData{data: s, reply: bindReply} // request bind.
	if err := <-bindReply; err != nil {              // wait for bind.
		return nil, err
	}
	l.cache.store(s)
	return s, nil
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

// loadShader returns a loaded shader immediately if it is cached.
// Otherwise the shader is returned after it is  loaded and bound.
func (l *loader) loadShader(s *shader) (*shader, error) {
	data := asset(s)
	if err := l.cache.fetch(&data); err == nil {
		s = data.(*shader)
		return s, nil // only initialized stuff is in the cache.
	}

	// Otherwise the shader needs to be loaded and bound.
	if err := l.importShader(s); err != nil {
		return nil, err
	}
	bindReply := make(chan error)
	l.binder <- &bindData{data: s, reply: bindReply} // request bind.
	if err := <-bindReply; err != nil {              // wait for bind.
		return nil, err
	}
	s.bound = true
	l.cache.store(s)
	return s, nil
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
func (l *loader) loadMesh(m *mesh) (*mesh, error) {
	data := asset(m)
	if err := l.cache.fetch(&data); err == nil {
		return data.(*mesh), nil
	}

	// Otherwise the mesh needs to be loaded and bound.
	if err := l.importMesh(m); err != nil {
		return nil, err
	}
	if err := l.bindMesh(m); err != nil {
		return nil, err
	}
	m.bound = true
	l.cache.store(m)
	return m, nil
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
func (l *loader) loadTexture(t *texture) (*texture, error) {
	data := asset(t)
	if err := l.cache.fetch(&data); err == nil {
		return data.(*texture), nil
	}

	// Otherwise the mesh needs to be loaded and bound.
	if err := l.importTexture(t); err != nil {
		return nil, err
	}
	bindReply := make(chan error)
	l.binder <- &bindData{data: t, reply: bindReply} // request bind.
	if err := <-bindReply; err != nil {              // wait for bind.
		return nil, err
	}
	t.bound = true
	l.cache.store(t)
	return t, nil
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
func (l *loader) loadMaterial(m *material) (*material, error) {
	data := asset(m)
	if err := l.cache.fetch(&data); err == nil {
		return data.(*material), nil
	}

	// Otherwise the mesh needs to be loaded (no binding necessary).
	if err := l.importMaterial(m); err != nil {
		return nil, err
	}
	l.cache.store(m)
	return m, nil
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
func (l *loader) loadFont(f *font) (*font, error) {
	data := asset(f)
	if err := l.cache.fetch(&data); err == nil {
		return data.(*font), nil
	}

	// Otherwise the font needs to be loaded and stored.
	if err := l.importFont(f); err != nil {
		return nil, err
	}
	l.cache.store(f)
	return f, nil
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

// loadAnim loads an animated model from disk. This will create
// multiple model assets including a mesh, textures, and animation data.
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
				data := asset(t) // now load any animation textures.
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
		if tex, err = l.loadTexture(tex); tex != nil && err == nil {
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
		log.Printf("loader.dispose unknown %T", d)
	}
}

// loader
// ============================================================================
// loadReq

// loadReq is a request to fetch asset data from persistent store.
// Expected to be created by the engine goroutine and passed on a channel
// to the loader goroutine. The loader loads the asset and passes back
// the request.
type loadReq struct {
	eid uint64 // pov entity identifier.
	a   asset  // asset to be loaded (anm, fnt, mat, msh, shd, snd, tex).
	err error  // true if there was an error with the load.

	// Extra assets generated when loading an animation file.
	msh  *mesh      // only valid after loading an animation.
	texs []*texture // only valid after loading an animation.

	// set by the engine and passed to the loader, but not referenced or changed
	// on the loader goroutine... Accessed by the engine after the load request.
	index int    // texture index used after load complete.
	model *model // model used after load complete. Used for all but snd assets.
	noise *noise // noise used after load complete. Matches snd assets.
}

// loaderAsset
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
