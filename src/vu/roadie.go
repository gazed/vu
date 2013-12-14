// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package vu

import (
	"log"
	"vu/audio"
	"vu/data"
	"vu/render"
)

// roadie is a resource manager.  It mediates between the loading and caching
// of the data resources and the consumption of the resources by subsystems
// like render and audio.  Roadie is responsible for ensuring that the loaded
// resources are prepared before consumption. Preparation consists of binding
// the resource which means copying the resource data onto the graphics or audio
// hardware.
type roadie struct {
	ld data.Loader     // Data load and cache subsystem.
	gc render.Renderer // Graphics subsystem injected on creation.
	ac audio.Audio     // Audio subsystem injected on creation.
}

// newRoadie is expected to be called once during engine initializes.
func newRoadie(ac audio.Audio, gc render.Renderer) *roadie {
	r := &roadie{}
	r.ld = data.NewLoader()
	r.gc = gc
	r.ac = ac
	return r
}

// dispose cleans up.  Roadie is only responsible for the data package.
func (r *roadie) dispose() { r.ld.Dispose() }

// useSound gets the named sound resource, lazy loading it if necessary.
func (r *roadie) useSound(sound string, s **data.Sound) {
	if !r.ld.Cached(sound, *s) {
		r.loadSound(sound, *s) // lazy load and bind.
	}
	if r.ld.Fetch(sound, s); *s == nil {
		log.Printf("roadie.useSound: could not fetch %s\n", sound)
	}
	return
}

// useMesh lazy loads and bind if necessary, and from then on just fetches the
// initialized resource from the cache.
func (r *roadie) useMesh(mesh string, m **data.Mesh) {
	if !r.ld.Cached(mesh, *m) {
		r.loadMesh(mesh, *m) // lazy load and bind.
	}
	if r.ld.Fetch(mesh, m); *m == nil {
		log.Printf("roadie.useMesh: could not fetch %s\n", mesh)
	}
	return
}

// useShader lazy loads and bind if necessary, and from then on just fetches the
// initialized resource from the cache.
func (r *roadie) useShader(shader string, s **data.Shader) {
	if !r.ld.Cached(shader, *s) {
		r.loadShader(shader, *s) // lazy load and bind.
	}
	if r.ld.Fetch(shader, s); *s == nil {
		log.Printf("roadie.useShader: could not fetch %s\n", shader)
	}
	return
}

// useMaterial lazy loads and bind if necessary, and from then on just fetches the
// initialized resource from the cache.
func (r *roadie) useMaterial(material string, m **data.Material) {
	if !r.ld.Cached(material, *m) {
		r.loadMaterial(material, *m) // lazy load and bind.
	}
	if r.ld.Fetch(material, m); *m == nil {
		log.Printf("roadie.useMaterial: could not fetch %s\n", material)
	}
	return
}

// useTexture lazy loads and bind if necessary, and from then on just fetches the
// initialized resource from the cache.
func (r *roadie) useTexture(texture string, t **data.Texture) {
	if !r.ld.Cached(texture, *t) {
		r.loadTexture(texture, *t) // lazy load and bind.
	}
	if r.ld.Fetch(texture, t); *t == nil {
		log.Printf("roadie.useTexture: could not fetch %s\n", texture)
	}
	return
}

// useGlyphs lazy loads and bind if necessary, and from then on just fetches the
// initialized resource from the cache.
func (r *roadie) useGlyphs(glyphset string, g **data.Glyphs) {
	if !r.ld.Cached(glyphset, *g) {
		r.loadGlyphs(glyphset, *g) // lazy load and bind.
	}
	if r.ld.Fetch(glyphset, g); *g == nil {
		log.Printf("roadie.useGlyphs: could not fetch %s\n", glyphset)
	}
	return
}

// loaded checks if the named type exists in the ld. Data is expected to
// be a pointer to one of the valid depot types.
func (r *roadie) loaded(name string, data interface{}) bool { return r.ld.Cached(name, data) }

// load places the given resource into the resource cache.
// This is expected to be used for small hand-crafted resources.
func (r *roadie) load(data interface{}) { r.ld.Cache(data) }

// loadTexture imports the named resource and makes it available in the
// resource cache.
func (r *roadie) loadTexture(name string, t *data.Texture) {
	if r.ld.Load(name, &t); t != nil {
		if err := r.gc.BindTexture(t); err != nil {
			log.Printf("roadie.loadTexture: could not bind texture %s %s\n", t.Name, err)
			return
		}
		r.ld.Cache(t)
	} else {
		log.Printf("roadie.loadTexture: could not load texture %s\n", name)
	}
}

// loadMesh imports the named resource and makes it available in the
// resource cache.
func (r *roadie) loadMesh(name string, m *data.Mesh) {
	if r.ld.Load(name, &m); m != nil {
		if err := r.gc.BindModel(m); err == nil {
			r.ld.Cache(m)
		} else {
			log.Printf("roadie.loadMesh: could not initialize mesh %s %s %s\n", name, m.Name, err)
		}
	}
}

// loadMaterial imports the named resource and makes it available in the
// resource cache.
func (r *roadie) loadMaterial(name string, m *data.Material) {
	if r.ld.Load(name, &m); m != nil {
		r.ld.Cache(m)
	}
}

// loadGlyphs imports the named resource and makes it available in the
// resource cache.
func (r *roadie) loadGlyphs(name string, g *data.Glyphs) {
	if r.ld.Load(name, &g); g != nil {
		r.ld.Cache(g)
	}
}

// loadSound imports the named resource and makes it available in the
// resource cache.
func (r *roadie) loadSound(name string, sound *data.Sound) {
	if r.ld.Load(name, &sound); sound != nil {
		if err := r.ac.BindSound(sound); err != nil {
			log.Printf("roadie.loadSound: could not load sound %s %s %s\n", name, sound.Name, err)
			return
		}
		r.ld.Cache(sound)
	} else {
		log.Printf("roadie.loadSound: could not load sound %s %s\n", name, sound.Name)
	}
}

// loadShader imports the named resource and makes it available in the
// resource cache.
func (r *roadie) loadShader(name string, shader *data.Shader) {
	var err error
	if render.CreateShader(name, shader); shader != nil && shader.Name != "" {

		// override the default shader source with anything found on disk.
		ds := &data.Shader{}
		if r.ld.Load(name, &ds); ds != nil {
			if ds.Vsh != nil {
				shader.Vsh = ds.Vsh
			}
			if ds.Fsh != nil {
				shader.Fsh = ds.Fsh
			}
		}
		shader.EnsureNewLines()
		if shader.Program, err = r.gc.BindShader(shader); err != nil {
			log.Printf("roadie.loadShader: could not bind %s\n%s", name, err)
			return
		}
		r.ld.Cache(shader)
	} else {
		log.Printf("roadie.loadShader: unknown shader %s", name)
	}
}
