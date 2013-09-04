// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package main

import (
	"fmt"
	"vu/data"
	"vu/device"
	"vu/math/lin"
	"vu/render/gl"
)

// ld tests and demonstrates loading a mesh model that has been exported from
// another tool - a .obj file from Blender in this case. Its really testing the
// "data" package with the key line being:
//        loader.Load("monkey", &mesh)
//
// This example relies on the basic OpenGL calls from package "vu/render/gl" to render
// the imported mesh.
func ld() {
	ld := &ldtag{}
	dev := device.New("Monkey", 400, 100, 800, 600)
	dev.SetResizer(ld)
	ld.initScene()
	dev.Open()
	for dev.IsAlive() {
		dev.ReadAndDispatch()
		ld.render()
		dev.SwapBuffers()
	}
	dev.Dispose()
}

// Globally unique "tag" for this example.
// Also hides any variables shared between methods in this example.
type ldtag struct {
	shaders   uint32
	vao       uint32
	mvpref    int32
	mvp       *lin.M4
	faceCount int32
}

func (ld *ldtag) Resize(w, y, width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
	ld.mvp = lin.M4Perspective(60, float32(width)/float32(height), 0.1, 50)
}

func (ld *ldtag) initScene() {
	gl.Init()
	mesh := &data.Mesh{}
	loader := data.NewLoader()
	loader.Load("monkey", &mesh)
	ld.faceCount = int32(len(mesh.F))

	// Gather the one scene into this one vertex array object.
	gl.GenVertexArrays(1, &ld.vao)
	gl.BindVertexArray(ld.vao)

	// vertex data.
	var vbuff uint32
	gl.GenBuffers(1, &vbuff)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbuff)
	gl.BufferData(gl.ARRAY_BUFFER, int64(len(mesh.V)*4), gl.Pointer(&(mesh.V[0])), gl.STATIC_DRAW)
	var vattr uint32 = 0
	gl.VertexAttribPointer(vattr, 4, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(vattr)

	// normal data.
	var nbuff uint32
	gl.GenBuffers(1, &nbuff)
	gl.BindBuffer(gl.ARRAY_BUFFER, nbuff)
	gl.BufferData(gl.ARRAY_BUFFER, int64(len(mesh.N)*4), gl.Pointer(&(mesh.N[0])), gl.STATIC_DRAW)
	var nattr uint32 = 1
	gl.VertexAttribPointer(nattr, 3, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(nattr)

	// faces data, uint32 in this case, so 4 bytes per element.
	var ebuff uint32
	gl.GenBuffers(1, &ebuff)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebuff)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int64(len(mesh.F)*2), gl.Pointer(&(mesh.F[0])), gl.STATIC_DRAW)

	ld.initShader()
	gl.ClearColor(0.2, 0.2, 0.2, 1.0)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)

	// set the initial perspetive matrix.
	ld.Resize(0, 0, 800, 600)
}

// Compile shaders and link to a program.
func (ld *ldtag) initShader() {
	shader := &data.Shader{}
	loader := data.NewLoader()
	loader.Load("monkey", &shader)
	ld.shaders = gl.CreateProgram()
	gl.BindAttribLocation(ld.shaders, 0, "inPosition")
	gl.BindAttribLocation(ld.shaders, 1, "inNormal")
	if err := gl.BindProgram(ld.shaders, shader.Vsh, shader.Fsh); err != nil {
		fmt.Printf("Failed to create program: %s\n", err)
	}
	ld.mvpref = gl.GetUniformLocation(ld.shaders, "modelViewProjectionMatrix")
	if ld.mvpref < 0 {
		fmt.Printf("No modelViewProjectionMatrix in vertex shader\n")
	}
}

// Draw the scene consisting of one VAO
func (ld *ldtag) render() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(ld.shaders)
	gl.BindVertexArray(ld.vao)

	// Use a modelview and projection matrix
	modelView := lin.M4Scaler(0.5, 0.5, 0.5).Mult(lin.M4Translater(0, 0, -2))
	modelViewProj := modelView.Mult(ld.mvp)
	gl.UniformMatrix4fv(ld.mvpref, 1, false, modelViewProj.Pointer())
	gl.CullFace(gl.BACK)
	gl.DrawElements(gl.TRIANGLES, ld.faceCount, gl.UNSIGNED_SHORT, gl.Pointer(nil))

	// cleanup
	gl.UseProgram(0)
	gl.BindVertexArray(0)
}
