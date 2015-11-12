// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// +build dx
// Use directx only when specifically asked for.

package render

import (
	"image"
)

// FUTURE link this up to directX bindings. Need a directX context from
//        vu/device. Wait for directX 12 to mature a bit.

// directx is the Direct3D implemntation of Renderer. See the Renderer interface
// for comments. See the Direct3D documentation for methods and constants.
type directx struct {
	depthTest bool   // Track current depth setting to reduce state switching.
	shader    uint32 // Track the current shader to reduce shader switching.
}

// newRenderer returns a DirectX implementation of Renderer.
func newRenderer() Renderer {
	gc := &directx{}
	return gc
}

// Render implementation specific constants.
const ()

// Renderer implementation.
func (gc *directx) Init() error {
	return nil
}

// Renderer implementation.
func (gc *directx) Color(r, g, b, a float32)       {}
func (gc *directx) Clear()                         {}
func (gc *directx) Viewport(width int, height int) {}

// Renderer implementation.
func (gc *directx) Enable(attribute uint32, enabled bool) {}

// Render implementation.
func (gc *directx) Render(dr Draw) {}

// Renderer implementation.
func (gc *directx) BindMesh(vao *uint32, vdata map[uint32]Data, fdata Data) error {
	return nil
}

// Renderer implementation.
func (gc *directx) BindVertexBuffer(vdata Data) {}

// Renderer implementation.
func (gc *directx) BindFaceBuffer(fdata Data) {}

// Renderer implementation.
func (gc *directx) BindShader(vsh, fsh []string, uniforms map[string]int32,
	layouts map[string]uint32) (program uint32, err error) {
	return 0, nil
}

// Renderer implementation.
func (gc *directx) BindTexture(tid *uint32, img image.Image, repeat bool) (err error) {
	return nil
}

// Renderer implementation. Remove GPU resources.
func (gc *directx) ReleaseMesh(vao uint32)    {}
func (gc *directx) ReleaseShader(sid uint32)  {}
func (gc *directx) ReleaseTexture(tid uint32) {}
