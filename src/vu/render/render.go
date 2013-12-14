// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

// Package render provides access to 3D graphics. Render attempts to
// shield the engine from graphic layer implementation details. Its job
// is to transfer 3D data to the graphics card and ensure it is made visible.
//
// Note that vectors and matricies are down converted before being passed
// to the graphics layer. This is because graphics cards are happier with
// 32 bit floats than the 64 bit floats prefered by math and physics.
// No math operations are expected or supported with the 32-bit vector or
// matrix structures. Use vu/math/lin instead.
//
// Package render is provided as part of the vu (virtual universe) 3D engine.
package render

import (
	"vu/data"
	"vu/render/gl"
)

// Renderer is used to draw 3D model objects using a graphics context.
// The expected usage is along the lines of:
//     • Initialize the graphics layer.
//     • Initialize the graphics data by "binding" it to the graphics card.
//       Binding copies data to the graphics card leaving references for
//       later manipulation by the engine.
//     • Enter a loop that calls Render many times a second, completely
//       redrawing all visible objects.
type Renderer interface {
	Init() (err error)                    // Init must be called first and only once.
	Clear()                               // Clear all buffers before rendering.
	Color(r, g, b, a float32)             // Set the default render clear colour
	Enable(attribute uint32, enable bool) // Enable or disable graphic state.
	Viewport(width int, height int)       // Set the available screen real estate.
	MapTexture(tid int, t *data.Texture)  // Link textures with texture units.

	// BindModel makes mesh model data available to the graphics card.
	// This copies the model data into the graphics card memory.
	// Note that the model data is no longer needed in the CPU memory.
	BindModel(mesh *data.Mesh) (err error)

	// BindGlyphs makes bitmap font data available to the graphics card.
	BindGlyphs(mesh *data.Mesh) (err error)

	// BindTexture makes texture data available to the graphics card. This
	// copies the texture data into the graphics card memory. Note that the
	// texture data is no longer needed in the CPU memory.
	BindTexture(texture *data.Texture) (err error)

	// BindShader combines a vertex and fragment shader into a program.
	// Error information is returned when there are compile or link problems.
	BindShader(shader *data.Shader) (programRef uint32, err error)

	// Render draws 3D data. Render is expected to be called by a Mesh to render
	// itself. This will be called frequently as each Scene draws its parts
	// multiple times a second.
	Render(v *Vis)
}

// Enable/Disable graphic state constants. These are the attributes used in the
// Renderer.Enable method.
const (
	BLEND uint32 = gl.BLEND      // Alpha blending.
	CULL         = gl.CULL_FACE  // Backface culling.
	DEPTH        = gl.DEPTH_TEST // Z-buffer (depth) awareness.
)

// Vis contains the information needed to completely render a Part.
// It passes visible scene graph information to the render system.
type Vis struct {
	Mv         *M4            // Model view.
	Mvp        *M4            // Model view projection.
	Scale      *V3            // Model scaling.
	L          *Light         // Light information.
	Alpha      float32        // Model alpha value.
	RotSpeed   float32        // Rotation speed.
	Fade       float32        // Fade value.
	Is2D       bool           // 2D or 3D.
	MeshName   string         // Mesh resource identifier.
	ShaderName string         // Shader identifier.
	MatName    string         // Material resource identifier.
	TexName    string         // Texture resource identifier.
	GlyphName  string         // Glyph resource identifier.
	GlyphText  string         // Banner text.
	GlyphPrev  string         // Previous banner text.
	GlyphWidth int            // Banner size.
	Mesh       *data.Mesh     // Mesh resource.
	Shader     *data.Shader   // Shader resource.
	Mat        *data.Material // Material resource.
	Tex        *data.Texture  // Texture resource.
	Glyph      *data.Glyphs   // Glyph resource.
}

// New provides a default graphics implementation.
func New() Renderer { return &opengl{} }

// NewVis creates and allocates space for a Vis structure.
func NewVis() *Vis {
	vis := &Vis{}
	vis.Mv = &M4{}
	vis.Mvp = &M4{}
	vis.Scale = &V3{}
	vis.L = &Light{}
	return vis
}
