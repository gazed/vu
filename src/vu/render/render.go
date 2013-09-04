// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

// Package render provides an inteface to the graphics card via an underlying
// abstraction layer (currently OpenGL). Render (attempts to) shield the engine
// from graphic layer implementation details.
//
// Package render is provided as part of the vu (virtual universe) 3D engine.
package render

import (
	"vu/data"
	"vu/math/lin"
)

// Renderer is used to draw 3D model objects using a graphics context.
// The expected usage is along the lines of:
//     * Initialize the graphics layer.
//     * Initialize the graphics data by "binding" it to the graphics card.
//       Binding copies data to the graphics card leaving references for
//       later manipulation by the engine.
//     * Enter a loop that calls Render many times a second, completely
//       redrawing all visible objects.
type Renderer interface {
	Init() (err error)                   // Init must be called first, and only once.
	Clear()                              // Clear all buffers before rendering.
	Color(r, g, b, a float32)            // Color sets the default render clear colour
	Enable(attribute int, enable bool)   // Enable or disable graphic state.
	Viewport(width int, height int)      // Viewport sets the available screen real estate.
	MapTexture(tid int, t *data.Texture) // MapTexture links textures with texture units.

	// BindModel makes mesh model data available to the graphics engine.
	// This copies the model data onto the graphics card memory (and so the
	// model data is no longer needed in the CPU memory).
	BindModel(mesh *data.Mesh) (err error)

	// BindGlyphs makes bitmap font data available to the graphics engine.
	BindGlyphs(mesh *data.Mesh) (err error)

	// BindTexture makes texture data available to the graphics engine.  This
	// copies the texture data onto the graphics card memory (and so the texture
	// data is no longer needed in the CPU memory).
	BindTexture(texture *data.Texture) (err error)

	// BindShader combines a vertex and fragment shader into a program.
	// The bulk of this method is error checking and returning error information
	// when there are compile or link problems.
	BindShader(shader *data.Shader) (programRef uint32, err error)

	// Render is expected to be called by a Mesh to render itself.  This will be
	// called frequently as a Scene draw all its parts.
	Render(v *Visible)
}

// Enable/Disable graphic state constants. These are the attributes used in the
// Renderer.Enable method.
const (
	BLEND int = iota // alpha blending.
	CULL             // backface culling.
	DEPTH            // z-buffer (depth) awareness.
)

// Visible contains the information needed to completely render a mesh. It is
// used in place of a very long string of parameters to the Renderer.Render
// method.
type Visible struct {
	Mv       *lin.M4        // Model view.
	Mvp      *lin.M4        // Model view projection.
	Mesh     *data.Mesh     // Mesh data.
	Shader   *data.Shader   // Shader program.
	Scale    float32        // Billboards need scale separated out.
	L        *Light         // Light information.
	Mat      *data.Material // Material information.
	Texture  *data.Texture  // Optional texture.
	RotSpeed float32        // Rotation speed for texture.
	Fade     float32        // Fade distance for fade shaders.
}

// New provides a default graphics implementation
func New() Renderer { return &opengl{} }
