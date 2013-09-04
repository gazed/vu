// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

// Package data turns file resource data into data objects. Specifically
// data needed by 3D applications. This package:
//    * Defines a data type for each of the various data resources.
//    * Provides caching and fetching of data resources.
//    * Loads data directly from disk for a development build and loads data
//      from a zip file attached to the binary for a production build.
//    * Loads data such that is easily consumed by the rendering
//      and audio components.
//
// The resource data that can be loaded from disk are the types exposed by
// this package:
//        Data                          File           Object
//       ------                        ------         --------
//      bitmapped fonts              : txtfile.fnt --> Glyphs
//      colour and surface data      : txtfile.mtl --> Material
//      vertex data                  : txtfile.obj --> Mesh
//      vertex shader program        : txtfile.vsh -┐
//      fragment shader program      : txtfile.fsh --> Shader
//      audio                        : binfile.wav --> Sound
//      images                       : binfile.png --> Texture
//
// Implementation Notes:
//    * The intent is to eventually have more than one supported file type for
//      a given resource.
//    * Currently intended for smaller 3D applications where data is loaded
//      from disk and kept in memory. There is nothing preventing a more
//      industrial strength (database) back end.
//
// Package data is provided as part of the vu (virtual universe) 3D engine.
package data

// Loader provides methods for loading, caching, and fetching data resources from disk.
// The current resource data types are the types exposed by this  package:
//    *Glyphs
//    *Material
//    *Mesh
//    *Shader
//    *Sound
//    *Texture
//
// Loader methods log attempts to use unsupported data types as development errors.
type Loader interface {

	// Load looks for the named resource on disk and returns a corresponding resource
	// data object. The loader will look for "name" with any of the supported file
	// types for the given resource.
	//
	// Load expects data to be a pointer to one of the resource data types.
	// If found the resource data is copied into the supplied data pointer,
	// otherwise the supplied data pointer is set to nil.
	Load(name string, data interface{})

	// Cache the given resource data so that any single resource only needs
	// to be loaded once. Cache expects data to be one of the valid resource data types
	// and to be uniquely named within its data type.
	Cache(data interface{})

	// Cached returns true if the named data resource has already been cached and ready
	// to use. Cached expects data to be one of the valid resource data types.
	Cached(name string, data interface{}) bool

	// Fetch retrieves a previously cached resource using the given name.
	// Fetch expects the resource data to be a pointer to one of the resource data types.
	// If found the resource data is copied into the supplied data pointer,
	// otherwise the supplied data pointer is set to nil
	Fetch(name string, data interface{})

	// SetDir overrides the default directory location for the given data type.
	// SetDir expects data to be one of the valid resource data types.
	//
	// Note that all directories are expected to be relative to the
	// application location.
	SetDir(dir string, data interface{})

	// Dispose needs to be called to properly terminate and clean up loaded resources.
	Dispose()
}

// Loader interface
// ===========================================================================
// dataLoader Loader implementation.

// NewLoader provides the default loader implmentation.
func NewLoader() Loader {
	return &dataLoader{ldr: newLoader(), dep: newDepot()}
}

// dataLoader is the default implementation of Loader.
type dataLoader struct {
	ldr *loader // Used to load resources.
	dep *depot  // Cached assets (populated by the application).
}

func (dl *dataLoader) Load(name string, data interface{})        { dl.ldr.load(name, data) }
func (dl *dataLoader) Cache(data interface{})                    { dl.dep.cache(data) }
func (dl *dataLoader) Cached(name string, data interface{}) bool { return dl.dep.cached(name, data) }
func (dl *dataLoader) Fetch(name string, data interface{})       { dl.dep.fetch(name, data) }
func (dl *dataLoader) SetDir(dir string, data interface{})       { dl.ldr.setDir(dir, data) }
func (dl *dataLoader) Dispose()                                  { dl.ldr.dispose() }
