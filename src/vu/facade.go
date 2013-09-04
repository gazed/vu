// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package vu

import (
	"vu/data"
)

// facade combines the an objects graphic resources into a viewable shape.
// It handles the appearance of an object. A facade is expected to be attached
// to a part within a scene in order to be rendered.  Geneally a facade is
// specified by the application and consumed by the engine rendering subsystem.
//
// The facades resources are lazy-loaded from the resource depot.
type facade struct {
	uid   int            // Unique id, expected to be same as linked part id.
	res   *roadie        // Resource manager.
	msh   *data.Mesh     // Facade frame.
	shadr *data.Shader   // Shape painter.
	mat   *data.Material // Optional, can be nil, alone, or combined with texture.
	tex   *data.Texture  // Optional texture.
	rots  float32        // Oexture rotation speed.
}

// newFacade creates a facade with the given mesh, shader, and material.
// Textures can optionally be included later using addTexture().  Each facade
// is given a reference to the roadie which is used to load and cache the
// needed resources.
//
// Usage: facades are created by the application using Part.SetFacade. The part
// gives the facade the its (not-so) unique id which is the same as the part uid.
func newFacade(uid int, r *roadie, mesh, shader, material string) *facade {
	f := &facade{}
	f.uid = uid
	f.res = r
	f.setMesh(mesh)
	f.setShader(shader)
	f.setMaterial(material)
	return f
}

// alpha gets the the current transparency attribute. Zero is returned if there is no
// transparency for this facade (or if the transparency is actually 0).
func (f *facade) alpha() float32 {
	if f.mat == nil {
		return 0
	}
	return f.mat.Tr
}

// setAlpha is used to set the transparency attribute.  This is expected to be a
// value between 0 and 1.0.
func (f *facade) setAlpha(a float32) {
	if f.mat != nil && a >= 0 && a <= 1 {
		f.mat.Tr = a
	}
}

// mesh provides safe access to the current mesh label associated with this facade.
// Return an empty string for uninitialized meshes.
func (f *facade) mesh() string {
	if f.msh == nil {
		return ""
	}
	return f.msh.Name
}

// setMesh initializes the surface mesh from one of the preloaded mesh.
func (f *facade) setMesh(mesh string) { f.msh = f.res.useMesh(mesh) }

// shader provides safe access to the current shader label associated with this
// facade. Return an empty string for uninitialized shaders.
func (f *facade) shader() string {
	if f.shadr == nil {
		return ""
	}
	return f.shadr.Name
}

// setShader initializes the surface shader from one of the preloaded shaders.
func (f *facade) setShader(shader string) { f.shadr = f.res.useShader(shader) }

// material provides safe access to the unique material label. Return an
// empty string for uninitialized materials.
func (f *facade) material() string {
	if f.mat == nil {
		return ""
	}
	return f.mat.Name
}

// setMaterial initializes the surface material from one of the
// preloaded materials.
func (f *facade) setMaterial(material string) { f.mat = f.res.useMaterial(material) }

// texture returns information on the desired texture.  Return an empty string label
// for uninitialized textures.
func (f *facade) texture(index int) (texture string, rotSpeed float32) {
	if f.tex == nil {
		return "", 0
	}
	return f.tex.Name, f.rots
}

// setTexture uses the given texture.
func (f *facade) setTexture(texture string, rotSpeed float32) {
	if t := f.res.useTexture(texture); t != nil {
		f.tex = t
		f.rots = rotSpeed
	}
}
