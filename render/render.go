// Copyright © 2013-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// Package render provides access to 3D graphics.
// It encapsulates and provides a common interface to graphic card
// programming APIs like OpenGL or DirectX. Render is OS indepdent.
// It relies on OS specific graphics contexts to be created by vu/device.
//
// Package render is provided as part of the vu (virtual universe) 3D engine.
package render

// FUTURE: render PC alternatives include Vulkan, DirectX.
// FUTURE: render alternatives for consoles.
// FUTURE: real time ray tracer renderer implementation... likely need
//         entirely different ray-trace based engine.

import (
	"image"
)

// Context is used to draw 3D model objects within a graphics context.
// The expected usage is along the lines of:
//     • Initialize the graphics layer.
//     • Create some 3D models by binding graphics data to the GPU.
//     • Render the 3D models many times a second.
type Context interface {
	Init() (err error)               // Call first, once at startup.
	Clear()                          // Clear all buffers before rendering.
	Color(r, g, b, a float32)        // Set the default render clear color
	Enable(attr uint32, enable bool) // Enable or disable graphic state.
	Viewport(width int, height int)  // Set the available screen real estate.

	// Rendering works with uniform values and data bound to the GPU.
	BindMesh(vao *uint32, vdata map[uint32]Data, fdata Data) error
	BindShader(vsh, fsh []string, uniforms map[string]int32,
		layouts map[string]uint32) (program uint32, err error)
	BindTexture(tid *uint32, img image.Image) error
	SetTextureMode(tid uint32, clamp bool)
	Render(d *Draw) // Render bound data, textures with bound shaders.

	// BindMap creates a framebuffer object with an associated
	// texture that can be used for shadow maps.
	//   fbo : returned frame buffer object identifier.
	//   tid : returned texture identifier.
	BindMap(fbo, tid *uint32) error

	// BindTarget creates a framebuffer object that can be used
	// as a render target.
	//   fbo : returned frame buffer object identifier.
	//   tid : returned texture identifier.
	//   db  : returned depth buffer render buffer.
	BindTarget(fbo, tid, db *uint32) error

	// Releasing frees up previous bound graphics card data.
	ReleaseMesh(vao uint32)            // Free bound vao reference.
	ReleaseShader(sid uint32)          // Free bound shader reference.
	ReleaseTexture(tid uint32)         // Free bound texture reference.
	ReleaseMap(fbo, tid uint32)        // Free shadow map framebuffer.
	ReleaseTarget(fbo, tid, db uint32) // Free render target framebuffer.
}

// New provides the render implementation as determined by the build.
func New() Context { return newRenderer() }

// Renderer implementation independent constants.
const (
	// Draw modes for vertex data rendering. Used in Draw.SetRefs.
	Triangles = iota // Triangles are the default for 3D models.
	Points           // Points are used for particle effects.
	Lines            // Lines are used for drawing wireframe shapes.

	// BindFrame buffer types.
	DepthBuffer        // For depth only.
	ImageBuffer        // For color and depth.
	LayerSize   = 1024 // Render pass texture size.
)
