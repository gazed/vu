// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// Package render provides access to 3D graphics. It makes data visible by
// sending 3D data to the graphics card. The main steps involved are:
//     • Create a Renderer.
//     • Create one or more Models, each Model associated with a Shader.
//     • Populate the Models with Meshes, Textures, and other Shader data.
//     • Rapidly and forever, call Renderer.Render(m) for each Model m.
// Package render is provided as part of the vu (virtual universe) 3D engine.
package render

// Renderer is used to draw 3D model objects within a graphics context.
// The expected usage is along the lines of:
//     • Initialize the graphics layer.
//     • Create 3D models using combinations of graphics data.
//     • Loop, rendering the 3D models many times a second.
type Renderer interface {
	Init() (err error)               // Call first, once at startup.
	Clear()                          // Clear all buffers before rendering.
	Color(r, g, b, a float32)        // Set the default render clear colour
	Enable(attr uint32, enable bool) // Enable or disable graphic state.
	Viewport(width int, height int)  // Set the available screen real estate.

	// Graphics data is encapsulated (combined and managed) in a Model.
	NewModel(s Shader) Model            // Model encapsulates the following:
	NewShader(name string) Shader       //    Shader program.
	NewMesh(name string) Mesh           //    Per vertex data.
	NewTexture(name string) Texture     //    Image data.
	NewAnimation(name string) Animation //    Animation data.
	Render(m Model)                     // Render draws a Model.
}

// New provides a default graphics implementation.
func New() Renderer { return newRenderer() }

// =============================================================================

// graphicsContext hides the existence of renderer methods that are local to
// this package. Internally classes that implement Renderer also implement
// graphicsContext.
type graphicsContext interface {
	Renderer // a graphicsContext is a Renderer

	// Binding data ensures the data is available on the graphics card.
	bindMesh(m Mesh) error
	bindShader(s Shader) error
	bindTexture(t Texture) error
	bindUniform(uniform int32, utype, num int, udata ...interface{})
	updateTextureMode(tex Texture)

	// Deleting frees up previous bound graphics card data. These are accessed
	// through the Model.Dispose methods.
	deleteMesh(mid uint32)
	deleteShader(sid uint32)
	deleteTexture(tid uint32)

	// useTexture makes the given bound texture t the active texture and
	// assigns it to the given texture unit (0-15). Sampler is the texture
	// sampler shader reference.
	useTexture(sampler, texUnit int32, t Texture)
}
