// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// Package render provides access to 3D graphics.
// It encapsulates and provides a common interface to graphic card
// programming interfaces like OpenGL and DirectX.
//
// Package render is provided as part of the vu (virtual universe) 3D engine.
package render

import (
	"image"
)

// Renderer is used to draw 3D model objects within a graphics context.
// The expected usage is along the lines of:
//     • Initialize the graphics layer.
//     • Create some 3D models by binding graphics data to the GPU.
//     • Render the 3D models many times a second.
type Renderer interface {
	Init() (err error)               // Call first, once at startup.
	Clear()                          // Clear all buffers before rendering.
	Color(r, g, b, a float32)        // Set the default render clear colour
	Enable(attr uint32, enable bool) // Enable or disable graphic state.
	Viewport(width int, height int)  // Set the available screen real estate.

	// Rendering works with uniform values and data bound to the GPU.
	BindMesh(vao *uint32, vdata map[uint32]Data, fdata Data) (err error)
	BindShader(vsh, fsh []string, uniforms map[string]int32,
		layouts map[string]uint32) (program uint32, err error)
	BindTexture(tid *uint32, img image.Image, repeat bool) (err error)
	Render(d Draw) // Render bound data and textures with bound shaders.

	// Releasing frees up previous bound graphics card data.
	ReleaseMesh(vao uint32)    // Free bound vao reference.
	ReleaseShader(sid uint32)  // Free bound shader reference.
	ReleaseTexture(tid uint32) // Free bound texture reference.
}

// New provides the graphics implementation determined by the
// renderer implementation that was included by the build,
// ie: OpenGL on OSX and Linux.
func New() Renderer { return newRenderer() }

// Renderer independent constants, ie: not specific to OpenGL or DirectX.
const (
	// Draw modes for vertex data rendering. Used in Draw.SetRefs.
	TRIANGLES = iota // Triangles are the default for 3D models.
	POINTS           // Points are used for particle effects.
	LINES            // Lines are used for drawing wireframe shapes.

	// Render buckets. Lower values drawn first. Used in Draw.SetHints.
	OPAQUE      // draw first
	TRANSPARENT // draw after opaque
	OVERLAY     // draw last.
)

// FUTURE: directx.go implementation to test the Render API and the
//         graphics layer encapsulation. Would need a corresponding
//         vu/render/dx package.
// FUTURE: real time ray tracer renderer implementation.
