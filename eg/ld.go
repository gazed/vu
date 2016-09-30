// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/gazed/vu/device"
	"github.com/gazed/vu/load"
	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/render"
	"github.com/gazed/vu/render/gl"
)

// ld loads a mesh model that has been exported from another tool.
// A OBJ file from Blender is used to test the vu/load package.
// This example renders using OpenGL from package vu/render/gl.
// See other examples use of the vu:Pov interface for a much easier way
// to load and render models when done as part of the vu engine.
//
// CONTROLS: NA
func ld() {
	ld := &ldtag{}
	dev := device.New("Load Model", 400, 100, 800, 600)
	ld.initScene()
	dev.Open()
	for dev.IsAlive() {
		ld.update(dev)
		ld.render()
		dev.SwapBuffers()
	}
	dev.Dispose()
}

// Globally unique "tag" that encapsulates example specific data.
type ldtag struct {
	shaders   uint32
	vao       uint32
	mvpref    int32
	persp     *lin.M4    // perspective matrix.
	mvp64     *lin.M4    // scratch for transform calculations.
	mvp       render.Mvp // transform matrix for rendering.
	faceCount int32
	loc       load.Locator
}

// update handles user input
func (ld *ldtag) update(dev device.Device) {
	pressed := dev.Update()
	if pressed.Resized {
		ld.resize(dev.Size())
	}
}

// resize handles user screen/window changes.
func (ld *ldtag) resize(x, y, width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
	ld.persp.Persp(60, float64(width)/float64(height), 0.1, 50)
}

// initScene is called once on startup to load the 3D data.
func (ld *ldtag) initScene() {
	ld.persp = lin.NewM4()
	ld.mvp64 = lin.NewM4()
	ld.mvp = render.NewMvp()
	ld.loc = load.NewLocator()
	gl.Init()
	mesh := &load.MshData{}
	mesh.Load("monkey", ld.loc)
	ld.faceCount = int32(len(mesh.F))

	// Gather the one scene into this one vertex array object.
	gl.GenVertexArrays(1, &ld.vao)
	gl.BindVertexArray(ld.vao)

	// vertex data.
	var vbuff uint32
	gl.GenBuffers(1, &vbuff)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbuff)
	gl.BufferData(gl.ARRAY_BUFFER, int64(len(mesh.V)*4), gl.Pointer(&(mesh.V[0])), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(0)

	// normal data.
	var nbuff uint32
	gl.GenBuffers(1, &nbuff)
	gl.BindBuffer(gl.ARRAY_BUFFER, nbuff)
	gl.BufferData(gl.ARRAY_BUFFER, int64(len(mesh.N)*4), gl.Pointer(&(mesh.N[0])), gl.STATIC_DRAW)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(1)

	// faces data, uint32 in this case, so 4 bytes per element.
	var fbuff uint32
	gl.GenBuffers(1, &fbuff)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, fbuff)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int64(len(mesh.F)*2), gl.Pointer(&(mesh.F[0])), gl.STATIC_DRAW)

	// final state setup before launch.
	ld.initShader()
	gl.ClearColor(0.2, 0.2, 0.2, 1.0)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
	ld.resize(0, 0, 800, 600)
}

// initShader compiles shaders and links them into a shader program.
func (ld *ldtag) initShader() {
	shader := &load.ShdData{}
	if err := shader.Load("monkey", ld.loc); err == nil {
		ld.shaders = gl.CreateProgram()
		if err := gl.BindProgram(ld.shaders, shader.Vsh, shader.Fsh); err != nil {
			fmt.Printf("Failed to create program: %s\n", err)
		}
		ld.mvpref = gl.GetUniformLocation(ld.shaders, "mvpm")
		if ld.mvpref < 0 {
			fmt.Printf("No modelViewProjectionMatrix in vertex shader\n")
		}
	}
}

// render draws the scene consisting of one VAO
func (ld *ldtag) render() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(ld.shaders)
	gl.BindVertexArray(ld.vao)

	// use a model-view-projection matrix
	ld.mvp64.Set(lin.M4I).ScaleSM(0.5, 0.5, 0.5).TranslateMT(0, 0, -2)
	ld.mvp.Set(ld.mvp64.Mult(ld.mvp64, ld.persp))
	gl.UniformMatrix4fv(ld.mvpref, 1, false, ld.mvp.Pointer())
	gl.DrawElements(gl.TRIANGLES, ld.faceCount, gl.UNSIGNED_SHORT, 0)

	// cleanup
	gl.UseProgram(0)
	gl.BindVertexArray(0)
}
