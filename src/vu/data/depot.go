// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package data

import (
	"log"
)

// depot is used to cache assets that have been loaded from disk.
// depot is expected to be used along with a loader so that each loaded element
// only needs to be loaded once.
//
// depot only deals with resource data types as defined in this package.
type depot map[int]map[string]interface{}

// newDepot creates a new cache for loaded items.  Expected to be called once
// during application initialization (since a cache works best when there is
// only one instance of it :).
func newDepot() *depot {
	return &depot{
		gly: make(map[string]interface{}),
		mat: make(map[string]interface{}),
		msh: make(map[string]interface{}),
		shd: make(map[string]interface{}),
		snd: make(map[string]interface{}),
		tex: make(map[string]interface{}),
	}
}

// cache an asset based on its type and name. Subsequent calls to cache the
// same asset/name combination are ignored.
func (d *depot) cache(data interface{}) {
	switch data.(type) {
	case *Glyphs:
		glyphset, _ := data.(*Glyphs)
		if glyphset != nil {
			d.cacheData(gly, glyphset.Name, data)
		}
	case *Material:
		material, _ := data.(*Material)
		if material != nil {
			d.cacheData(mat, material.Name, data)
		}
	case *Mesh:
		mesh, _ := data.(*Mesh)
		if mesh != nil {
			d.cacheData(msh, mesh.Name, data)
		}
	case *Shader:
		shader, _ := data.(*Shader)
		if shader != nil {
			d.cacheData(shd, shader.Name, data)
		}
	case *Sound:
		sound, _ := data.(*Sound)
		if sound != nil {
			d.cacheData(snd, sound.Name, data)
		}
	case *Texture:
		texture, _ := data.(*Texture)
		if texture != nil {
			d.cacheData(tex, texture.Name, data)
		}
	default:
		log.Printf("depot.cache: resource type is unknown")
	}

}

// cacheData ensures that data is stored once.
func (d *depot) cacheData(room int, name string, data interface{}) {
	if room >= 0 && room <= tex && name != "" {
		if _, ok := (*d)[room][name]; !ok {
			(*d)[room][name] = data
		}
	}
}

// fetch retrieves a previously cached data resource.  Nil is returned if the
// named data resource is not found.
func (d *depot) fetch(name string, data interface{}) {
	switch data.(type) {
	case **Glyphs:
		dataPtr, _ := data.(**Glyphs)
		*dataPtr = nil
		if stored := (*d)[gly][name]; stored != nil {
			*dataPtr, _ = stored.(*Glyphs)
		}
	case **Material:
		dataPtr, _ := data.(**Material)
		*dataPtr = nil
		if stored := (*d)[mat][name]; stored != nil {
			*dataPtr, _ = stored.(*Material)
		}
	case **Mesh:
		dataPtr, _ := data.(**Mesh)
		*dataPtr = nil
		if stored := (*d)[msh][name]; stored != nil {
			*dataPtr, _ = stored.(*Mesh)
		}
	case **Shader:
		dataPtr, _ := data.(**Shader)
		*dataPtr = nil
		if stored := (*d)[shd][name]; stored != nil {
			*dataPtr, _ = stored.(*Shader)
		}
	case **Sound:
		dataPtr, _ := data.(**Sound)
		*dataPtr = nil
		if stored := (*d)[snd][name]; stored != nil {
			*dataPtr, _ = stored.(*Sound)
		}
	case **Texture:
		dataPtr, _ := data.(**Texture)
		*dataPtr = nil
		if stored := (*d)[tex][name]; stored != nil {
			*dataPtr, _ = stored.(*Texture)
		}
	default:
		log.Printf("depot.fetch: resource type is unknown")
	}
}

// cached returns true if the named asset is in the cache.
func (d *depot) cached(name string, data interface{}) (loaded bool) {
	loaded = false
	switch data.(type) {
	case *Glyphs:
		_, loaded = (*d)[gly][name]
	case *Material:
		_, loaded = (*d)[mat][name]
	case *Mesh:
		_, loaded = (*d)[msh][name]
	case *Shader:
		_, loaded = (*d)[shd][name]
	case *Sound:
		_, loaded = (*d)[snd][name]
	case *Texture:
		_, loaded = (*d)[tex][name]
	}
	return
}
