// Copyright © 2015-2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// loader.go gets data from disk. Puts disk loads on worker goroutines.
// FUTURE: handle releaseData requests. See eng.dispose design note.

import (
	"fmt"
	"log"
	"time"

	"github.com/gazed/vu/load"
	"github.com/gazed/vu/render"
)

// loader interfaces between the engine models and the load package.
// loader imports and prepares model and sound data for use.
// loader is used by the engine to cache and reuse imported asset
// data amongst multiple model and noise instances.
type loader struct {
	loc   load.Locator // Locates the asset data on disk.
	cache cache        // asset cache.

	// pending tracks outstanding asset requests. Asset requests are
	// loaded from disk by workers on a goroutine. Syncronization happens
	// by passing ownership of asset data.
	pending   map[aid][]func(asset)
	needAsset chan *diskAsset // assets to be imported from disk.
	haveAsset chan *diskAsset // assets finished importing.
}

// newLoader is called once on startup by the engine.
func newLoader() *loader {
	l := &loader{}
	l.loc = load.NewLocator()
	l.cache = newCache()
	l.pending = map[aid][]func(asset){}

	// allocate enough that ideally avoids blocking and waiting
	// for asset import workers in reasonable scenarios.
	// FUTURE: profile this.
	numWorkers := 5
	l.needAsset = make(chan *diskAsset, 100)
	l.haveAsset = make(chan *diskAsset, 100)
	for wid := 1; wid <= numWorkers; wid++ {
		go importer(wid, l.needAsset, l.haveAsset)
	}
	return l
}

// fetch is called by application created entities that need assets.
// The requested asset is returned using the loaded callback. Cached assets
// are returned immediately while uncached assets are returned after they
// have been imported.
func (l *loader) fetch(a asset, loaded func(asset)) {
	kind := a.aid().kind()
	switch kind {
	case msh, tex, shd, fnt, mat, snd:
		l.fetchOrImport(a, loaded)
	case anm:
		// animations are special in that they need animation data and
		// mesh data where both are imported from the same file.
		m := newMesh(a.label())
		if err := l.cache.fetch(&a); err == nil {
			al := a.(*animation) // got the animation.
			data := asset(m)     // now get the mesh.
			if err := l.cache.fetch(&data); err == nil {
				ml := data.(*Mesh) // got the mesh.
				loaded(al)
				loaded(ml)
				break
			}
		}
		// make animation first asset to make it easier to parse
		// when request is finished.
		l.importAsset(loaded, a, m)
	default:
		log.Printf("loader: unknown request %T", a)
	}
}

// fetchOrImport returns cached assets immediately.
// Otherwise sends the assets off for import.
func (l *loader) fetchOrImport(a asset, loaded func(asset)) {
	if err := l.cache.fetch(&a); err == nil {
		loaded(a) // retrieved from cache.
		return
	}
	l.importAsset(loaded, a) // need to import from disk.
}

// importAsset places an asset import request on the work queue.
// Imports that load more than a single asset from the same disk resource
// will trigger off of the first asset when finished.
func (l *loader) importAsset(callback func(asset), assets ...asset) {
	if len(assets) > 0 {
		a := assets[0]
		if callbacks, ok := l.pending[a.aid()]; !ok {
			l.needAsset <- &diskAsset{assets: assets, loc: l.loc}
			l.pending[a.aid()] = append([]func(asset){}, callback)
		} else {
			l.pending[a.aid()] = append(callbacks, callback)
		}
	}
}

// processImports is run on the main thread each update tick to retreive
// the results of any asset import workers. It fetches finished imports
// within a time window so as to not stall the main loop. This results
// in models becoming visible over time instead of a blank screen.
//
// Asset binding is done on the main thread due to the single
// threaded render context.
func (l *loader) processImports() {
	var timeUsed time.Duration
	start := time.Now()
	timeLimit := 0.01 // 10 milliseconds, about half an update cycle.
	for timeUsed.Seconds() < timeLimit {
		select {
		case done := <-l.haveAsset:
			a := done.assets[0]
			if done.err != nil {
				log.Printf("Failed to load asset %s: %s", a.label(), done.err)
				break // dev error - dev to debug why asset is missing.
			}

			// handle assets that need binding (copy data to GPU).
			switch a.(type) {
			case *Mesh, *shader, *Texture, *sound, *material, *font:
				l.cache.store(a) // cache asset.
			case *animation:
				if len(done.assets) == 2 {
					a2 := done.assets[1] // animation mesh always second.
					l.cache.store(a2)    // cache mesh data.
					l.cache.store(a)     // cache animation data.
				}
			}
			for _, callback := range l.pending[a.aid()] {
				for _, a := range done.assets {
					callback(a)
				}
			}
			delete(l.pending, a.aid())
		default:
			// Called each update so return immediately if there are no assets.
			return
		}
		timeUsed = time.Since(start)
	}
}

// release is called when the cached data is no
// longer needed and can be discarded entirely.
func (l *loader) release(data interface{}) {
	switch d := data.(type) {
	case *Mesh:
		l.cache.remove(d)
	case *font:
		l.cache.remove(d)
	case *shader:
		l.cache.remove(d)
	case *material:
		l.cache.remove(d)
	case *Texture:
		l.cache.remove(d)
	case *animation:
		l.cache.remove(d)
	case *sound:
		l.cache.remove(d)
	default:
		log.Printf("loader.dispose unknown %T", d)
	}
}

// dispose is called when the engine is shutting down.
func (l *loader) dispose() {
	close(l.needAsset) // Close worker queue since no more sends.
	// however don't close the receiving channel in case there
	// are workers trying to write to it.

	// FUTURE: engine is shutting down. Could try to release every asset
	//         in the cache, but this doesn't matter as long as there is
	//         no call to dispose when the device display is closed.
	//         Why release only sometimes?
	//         Also dispose is called on the update goroutine.
	//         Need engine to release the cache on the main thread.
}

// loader
// =============================================================================
// diskAsset import worker.

// diskAsset is an asset that needs to be loaded into memory from persistent
// store by an importer. Currently only animations need more than one asset
// loaded from the same file.
type diskAsset struct {
	assets []asset      // asset instances to be loaded.
	err    error        // used to report asset import errors.
	loc    load.Locator // helper to find assets on disk.
}

// importer imports assets from persistent store. Currently persistent
// storage are just files on disk. This runs as a goroutine taking requests
// for assets from the needAsset channel and returning the loaded asset
// on the fetched channel.
func importer(id int, needAsset <-chan *diskAsset, fetched chan<- *diskAsset) {
	for da := range needAsset {
		importAsset(id, da)
		fetched <- da
	}
}

// importAsset is separate from importer for unit testing.
// The import* methods fill in the data referenced by the diskAsset assets.
func importAsset(id int, da *diskAsset) {
	switch a := da.assets[0].(type) {
	case *Mesh:
		da.err = importMesh(da.loc, a)
	case *shader:
		da.err = importShader(da.loc, a)
	case *Texture:
		da.err = importTexture(da.loc, a)
	case *sound:
		da.err = importSound(da.loc, a)
	case *material:
		da.err = importMaterial(da.loc, a)
	case *font:
		da.err = importFont(da.loc, a)
	case *animation:
		m, ok := da.assets[1].(*Mesh) // mesh is second by convention.
		if !ok {
			da.err = fmt.Errorf("Need mesh to load animation")
			break
		}
		da.err = importAnim(da.loc, a, m)
	default:
		da.err = fmt.Errorf("No import for %T", a)
	}
}

// importSound transfers audio data loaded from disk to the sound object.
func importSound(loc load.Locator, s *sound) error {
	snd := &load.SndData{}
	if err := snd.Load(s.name, loc); err != nil {
		return fmt.Errorf("importSound %s: %s", s.label(), err)
	}
	transferSound(snd, s)
	return nil
}

// importShader transfers data loaded from disk to the render object.
// Disk based files override predefined engine shaders.
func importShader(loc load.Locator, s *shader) error {
	shd := &load.ShdData{}
	if err := shd.Load(s.name, loc); err == nil {
		s.setSource(shd.Vsh, shd.Fsh) // first look for .vsh, .fsh on disk.
		return nil
	}

	// next look for a pre-defined engine shader.
	if sfn, ok := shaderLibrary[s.name]; ok {
		vsrc, fsrc := sfn()
		s.setSource(vsrc, fsrc)
		return nil
	}
	return fmt.Errorf("importShader could not find %s", s.name)
}

// importMesh transfers data loaded from disk to the render object.
func importMesh(loc load.Locator, m *Mesh) error {
	msh := &load.MshData{}
	if err := msh.Load(m.name, loc); err != nil {
		return fmt.Errorf("importMesh %s: %s", m.name, err)
	}
	transferMesh(msh, m)
	return nil
}

// importTexture transfers data loaded from disk to the render object.
func importTexture(loc load.Locator, t *Texture) error {
	img := &load.ImgData{}
	err := img.Load(t.name, loc)
	if err != nil {
		return fmt.Errorf("importTexture %s: %s", t.name, err)
	}
	t.Set(img.Img)
	return nil
}

// importMaterial transfers data loaded from disk to the render object.
func importMaterial(loc load.Locator, m *material) error {
	mtl := &load.MtlData{}
	if err := mtl.Load(m.label(), loc); err == nil {
		transferMaterial(mtl, m)
		return nil
	}
	return fmt.Errorf("importMaterial: %s", m.name)
}

// importFont transfers data loaded from disk to the render object.
func importFont(loc load.Locator, f *font) error {
	fnt := &load.FntData{}
	if err := fnt.Load(f.label(), loc); err != nil {
		return fmt.Errorf("importFont %s: %s", f.label(), err)
	}
	transferFont(fnt, f)
	return nil
}

// importAnim loads the animation and mesh from a single asset resource.
func importAnim(loc load.Locator, a *animation, m *Mesh) (err error) {
	mod := &load.ModData{}
	if err = mod.Load(a.name, loc); err != nil {
		return err
	}
	transferAnim(mod, m, a)
	return nil
}

// diskAsset import worker.
// =============================================================================
// utility methods to move loaded data into formats needed by
// the render system.

// transferMesh moves data from the loading system to the engine instance.
func transferMesh(data *load.MshData, m *Mesh) {
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
func transferAnim(data *load.ModData, m *Mesh, a *animation) {
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
