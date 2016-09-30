// Copyright © 2015-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// loader.go
// FUTURE: benchmark running the file imports on worker goroutines.
//         This will require tracking outstanding load requests in
//         loader and tracking load requests waiting for the outstanding
//         loads. Once imports are done they can be cached and then
//         matched against outstanding load requests.

import (
	"fmt"
	"log"

	"github.com/gazed/vu/load"
	"github.com/gazed/vu/render"
)

// loader interfaces between the engine models and the load package.
// loader imports and prepares model and noise data for use.
// It is expected to be run as a goroutine and thus needs to
// be initialized with communication channels.
//
// The loader is used internally by the engine to cache and reuse
// imported data for multiple model and noise instances.
type loader struct {
	assets map[aid]string // Track asset existence.

	// loader goroutine communication.
	loc    load.Locator        // Locates the asset data on disk.
	cache  cache               // asset cache.
	stop   chan bool           // shutdown requests.
	load   chan map[aid]string // asset load requests.
	loaded chan map[aid]asset  // loaded asset replies.
	bind   chan msg            // machine loop request channel.
}

// newLoader is expected to be called once on startup by the engine.
func newLoader(reqs chan map[aid]string, done chan map[aid]asset,
	bind chan msg, stop chan bool) *loader {
	l := &loader{bind: bind, assets: map[aid]string{}}
	l.loc = load.NewLocator()
	l.cache = newCache()
	l.stop = stop
	l.load = reqs
	l.loaded = done
	return l
}

// runLoader loops forever processing all load requests.
// It is started once as a goroutine on engine initialization
// and is stopped when the engine shuts down.
func runLoader(machine chan msg, load chan map[aid]string,
	loaded chan map[aid]asset, stop chan bool) {
	l := newLoader(load, loaded, machine, stop)
	defer catchErrors()
	for {
		select {
		case <-l.stop: // Stop on any value. Closed channels return 0
			return // exit immediately.
		case requests := <-l.load:

			// FUTURE: spawn the load requests off to worker goroutines.
			assets := map[aid]asset{}
			for id, name := range requests {
				a := id.dataType()
				switch a {
				case anm:
					msh := newMesh(name)
					anm := newAnimation(name)
					if la, lm := l.loadAnim(anm, msh); la == nil || lm == nil {
						log.Printf("Animation %s failed to load", name) // dev error.
					} else {
						assets[la.aid()] = la
						assets[lm.aid()] = lm
					}
				case msh:
					if lm, err := l.loadMesh(newMesh(name)); err != nil {
						log.Printf("Mesh %s failed to load %s", name, err) // dev error.
					} else {
						assets[lm.aid()] = lm
					}
				case tex:
					if lt, err := l.loadTexture(newTexture(name)); err != nil {
						log.Printf("Texture %s failed to load %s", name, err) // dev error.
					} else {
						assets[lt.aid()] = lt
					}
				case shd:
					if ls, err := l.loadShader(newShader(name)); err != nil {
						log.Printf("Shader %s failed to load %s", name, err) // dev error.
					} else {
						assets[ls.aid()] = ls
					}
				case fnt:
					if lf, err := l.loadFont(newFont(name)); err != nil {
						log.Printf("Font %s failed to load %s", name, err) // dev error.
					} else {
						assets[lf.aid()] = lf
					}
				case mat:
					if lm, err := l.loadMaterial(newMaterial(name)); err != nil {
						log.Printf("Material %s failed to load %s", name, err) // dev error.
					} else {
						assets[lm.aid()] = lm
					}
				case snd:
					if ls, err := l.loadSound(newSound(name)); err != nil {
						log.Printf("Sound %s failed to load %s", name, err) // dev error.
					} else {
						assets[ls.aid()] = ls
					}
				default:
					log.Printf("loader: unknown request %T", a)
					// FUTURE: handle releaseData requests. See eng.dispose design note.
				}
			}
			l.returnAssets(assets)
		}
	}
}

// returnAssets funnels a group of loaded assets back to the
// engine loop. The engine loop may be busy so this method
// starts a goroutine so the loader is not blocked.
func (l *loader) returnAssets(loaded map[aid]asset) {
	go func(loaded map[aid]asset) {
		l.loaded <- loaded
	}(loaded)
}

// shutdown is called on the engine processing goroutine
// to shutdown the loader.
func (l *loader) shutdown() { l.stop <- true }

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
	l.bind <- &bindData{data: s, reply: bindReply} // request bind.
	if err := <-bindReply; err != nil {            // wait for bind.
		return nil, err
	}
	l.cache.store(s)
	return s, nil
}

// importSound transfers audio data loaded from disk to the sound object.
func (l *loader) importSound(s *sound) error {
	snd := &load.SndData{}
	if err := snd.Load(s.name, l.loc); err != nil {
		return fmt.Errorf("loader.loadSound: could not load %s %s", s.label(), err)
	}
	transferSound(snd, s)
	return nil
}

// loadShader returns a loaded shader immediately if it is cached.
// Otherwise the shader is returned after it is loaded and bound.
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
	l.bind <- &bindData{data: s, reply: bindReply} // request bind.
	if err := <-bindReply; err != nil {            // wait for bind.
		return nil, err
	}
	l.cache.store(s)
	return s, nil
}

// importShader transfers data loaded from disk to the render object.
// Disk based files override predefined engine shaders.
func (l *loader) importShader(s *shader) error {
	shd := &load.ShdData{}
	if err := shd.Load(s.name, l.loc); err == nil {
		s.setSource(shd.Vsh, shd.Fsh) // first look for .vsh, .fsh on disk.
		return nil
	}

	// next look for a pre-defined engine shader.
	if sfn, ok := shaderLibrary[s.name]; ok {
		vsrc, fsrc := sfn()
		s.setSource(vsrc, fsrc)
		return nil
	}
	return fmt.Errorf("Could not find shader %s", s.name)
}

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
	l.cache.store(m)
	return m, nil
}

// bindMesh submits a mesh for binding or rebinding. This transfers
// the mesh data to the GPU.
func (l *loader) bindMesh(m *mesh) error {
	bindReply := make(chan error)
	l.bind <- &bindData{data: m, reply: bindReply} // request bind.
	return <-bindReply                             // wait for bind.
}

// importMesh transfers data loaded from disk to the render object.
func (l *loader) importMesh(m *mesh) error {
	msh := &load.MshData{}
	if err := msh.Load(m.name, l.loc); err != nil {
		return fmt.Errorf("loader.loadMesh: could not load %s %s", m.name, err)
	}
	transferMesh(msh, m)
	return nil
}

// loadTexture returns a loaded texture immediately if it is cached.
// Otherwise the texture is returned after it is loaded and bound.
func (l *loader) loadTexture(t *texture) (*texture, error) {
	data := asset(t)
	if err := l.cache.fetch(&data); err == nil {
		return data.(*texture), nil
	}

	// Otherwise the texture needs to be loaded and bound.
	if err := l.importTexture(t); err != nil {
		return nil, err
	}
	bindReply := make(chan error)
	l.bind <- &bindData{data: t, reply: bindReply} // request bind.
	if err := <-bindReply; err != nil {            // wait for bind.
		return nil, err
	}
	l.cache.store(t)
	return t, nil
}

// importTexture transfers data loaded from disk to the render object.
func (l *loader) importTexture(t *texture) error {
	img := &load.ImgData{}
	err := img.Load(t.name, l.loc)
	if err != nil {
		return fmt.Errorf("loader.loadTexture: could not load %s %s", t.name, err)
	}
	t.Set(img.Img)
	return nil
}

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
	mtl := &load.MtlData{}
	if err := mtl.Load(m.label(), l.loc); err == nil {
		transferMaterial(mtl, m)
		return nil
	}
	return fmt.Errorf("loader.loadMaterial: could not load %s", m.name)
}

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
	fnt := &load.FntData{}
	if err := fnt.Load(f.label(), l.loc); err != nil {
		return fmt.Errorf("loader.loadFont: could not load %s %s", f.label(), err)
	}
	transferFont(fnt, f)
	return nil
}

// loadAnim loads an animated model from disk. This will create
// multiple model assets including a mesh, textures, and animation data.
func (l *loader) loadAnim(a *animation, m *mesh) (*animation, *mesh) {
	data := asset(a)
	if err := l.cache.fetch(&data); err == nil {
		a = data.(*animation) // got the animation.
		data := asset(m)      // now load the mesh.
		if err := l.cache.fetch(&data); err == nil {
			m = data.(*mesh)
			return a, m
		}
	}

	// Otherwise the animation and mesh need to be loaded.
	// Textures are loaded, bound, and cached within importAnim.
	var err error
	if err = l.importAnim(a, m); err != nil {
		log.Printf("Animation load %s: %s", a.name, err)
		return nil, nil // discard load failures
	}

	// And the mesh needs to be bound.
	if err := l.bindMesh(m); err != nil {
		log.Printf("Animation bind %s: %s", m.name, err)
		return nil, nil // discard bind failures
	}
	l.cache.store(a)
	l.cache.store(m)
	return a, m
}

// importAnim loads the animation, mesh, and texture for an
// animated model.
func (l *loader) importAnim(a *animation, m *mesh) (err error) {
	mod := &load.ModData{}
	if err = mod.Load(a.name, l.loc); err != nil {
		return err
	}
	transferAnim(mod, m, a)
	return nil
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
// the completed request.
type loadReq struct {
	eid uint64 // pov entity identifier.
	a   asset  // asset to be loaded (anm, fnt, mat, msh, shd, snd, tex, fbo).
	err error  // true if there was an error with the load.

	// Extra assets generated when loading an animation file.
	msh  *mesh      // only valid after loading an animation.
	texs []*texture // only valid after loading an animation.

	// engine specific data used to quickly access the model or noise instance
	// to store the loaded asset. Not to be used by the loader goroutine.
	// Accessed by the engine after the load request.
	index int         // texture index used after load complete.
	data  interface{} // model or noise. Needed by engine after load complete.
}

// loadReq
// =============================================================================
// utility transfer methods to move data from the load system to the engine.

// transferMesh moves data from the loading system to the engine instance.
func transferMesh(data *load.MshData, m *mesh) {
	m.InitData(0, 3, render.StaticDraw, false).SetData(0, data.V)
	m.InitFaces(render.StaticDraw).SetFaces(data.F)
	if len(data.N) > 0 {
		m.InitData(1, 3, render.StaticDraw, false).SetData(1, data.N)
	}
	if len(data.T) > 0 {
		m.InitData(2, 2, render.StaticDraw, false).SetData(2, data.T)
	}
}

// transferMaterial moves data from the loading system to the engine instance.
func transferMaterial(data *load.MtlData, m *material) {
	m.kd.R, m.kd.G, m.kd.B = data.KdR, data.KdG, data.KdB
	m.ks.R, m.ks.G, m.ks.B = data.KsR, data.KsG, data.KsB
	m.ka.R, m.ka.G, m.ka.B = data.KaR, data.KaG, data.KaB
	m.tr = data.Alpha
	m.ns = data.Ns
}

// transferSound moves data from the loading system to the engine instance.
func transferSound(data *load.SndData, s *sound) {
	a := data.Attrs
	s.data.Set(a.Channels, a.SampleBits, a.Frequency, a.DataSize, data.Data)
}

// transferFont moves data from the loading system to the engine instance.
func transferFont(data *load.FntData, f *font) {
	f.setSize(data.W, data.H)
	for _, ch := range data.Chars {
		f.addChar(ch.Char, ch.X, ch.Y, ch.W, ch.H, ch.Xo, ch.Yo, ch.Xa)
	}
}

// transferAnim moves data from the loading system to the engine instance.
// Animation data is a combination of mesh data and animation data.
func transferAnim(data *load.ModData, m *mesh, a *animation) {
	transferMesh(&data.MshData, m)
	if len(data.Blends) > 0 {
		m.InitData(4, 4, render.StaticDraw, false).SetData(4, data.Blends)
	}
	if len(data.Weights) > 0 {
		m.InitData(5, 4, render.StaticDraw, true).SetData(5, data.Weights)
	}

	// Store the animation data.
	if len(data.Frames) > 0 {
		moves := []movement{}
		for _, ia := range data.Movements {
			movement := movement{
				name: ia.Name,
				f0:   int(ia.F0),
				fn:   int(ia.Fn),
				rate: float64(ia.Rate)}
			moves = append(moves, movement)
		}
		a.setData(data.Frames, data.Joints, moves)
	}
}

// utility transfer methods.
// =============================================================================
// cache

// cache reuses loaded assets.
type cache map[aid]interface{}

// newCache creates a new in-memory cache for loaded items. Expected to be
// called once during application initialization (since a cache works best
// when there is only one instance of it :)
func newCache() cache {
	return make(map[aid]interface{})
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
	return fmt.Errorf("cache.fetch: could not fetch asset")
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
