// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/gazed/vu/device"
	"github.com/gazed/vu/load"
	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/render"
	"github.com/gazed/vu/render/gl"
)

// tr demonstrates basic OpenGL by drawing a triangle. If anything this example
// shows that basic OpenGL is not all that basic. Check out the following methods:
//   • initData()   creates a triangle mesh that includes verticies, faces, and colors.
//   • initScene()  makes the data available to the graphics card.
//   • initShader() uses render.BindProgram to load and prepare the shader programs.
//   • drawScene()  is called to render and spin the triangle.
// Note that the OpenGL go bindings in vu.render/gl have been generated.
//
// CONTROLS: NA
func tr() {
	tag := &trtag{}
	dev := device.New("Triangle", 400, 100, 600, 600)
	tag.initScene()
	dev.Open()
	tag.resize(600, 600)
	for dev.IsAlive() {
		dev.Update()
		tag.drawScene()
		dev.SwapBuffers()
	}
	dev.Dispose()
}

// Globally unique "tag" that encapsulates example specific data.
type trtag struct {
	shaders     uint32
	vao         uint32
	rotateAngle float64
	mvpRef      int32
	persp       *lin.M4    // Perspective matrix.
	mvp64       *lin.M4    // Scratch for transform calculations.
	mvp         render.Mvp // transform matrix for rendering.
	lastTime    time.Time  // Controls triangle rotation speed.

	// mesh information
	verticies []float32
	faces     []uint8
	color     []float32
}

// A tetrahedron (triangular pyramid) with sides of length 2 and
// centered at the origin (in order to look nice), would have points:
//    (0,  2/√3,  −1/√6)
//    (±1, −1/√3, −1/√6)
//    (0,  0,     3/√6)
func (tag *trtag) initData() {
	tag.verticies = []float32{
		0, float32(2.0 / math.Sqrt(3)), float32(-1.0 / math.Sqrt(6)), // Red
		-1, float32(-1.0 / math.Sqrt(3)), float32(-1.0 / math.Sqrt(6)), // Green
		1, float32(-1.0 / math.Sqrt(3)), float32(-1.0 / math.Sqrt(6)), // Blue
		0, 0, float32(3.0 / math.Sqrt(6)), // White
	}
	tag.faces = []uint8{
		0, 2, 1, // front face  RGB
		2, 0, 3, // left face   BRW
		1, 3, 0, // right face  GRW
		1, 2, 3, // bottom face GBW
	}
	tag.color = []float32{
		1, 0, 0, 1, // vertex 0 is red
		0, 1, 0, 1, // vertex 1 is green
		0, 0, 1, 1, // vertex 2 is blue
		1, 1, 1, 1, // vertex 3 is white
	}
}

// resize sets the view port size. User resizes are ignored.
func (tag *trtag) resize(width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
	tag.persp.Persp(60, float64(width)/float64(height), 0.1, 50)
}

// initScene is one time initialization that creates a single VAO
func (tag *trtag) initScene() {
	tag.mvp64 = lin.NewM4()
	tag.persp = lin.NewM4().Persp(60, float64(600)/float64(600), 0.1, 50)
	tag.mvp = render.NewMvp()
	tag.initData()

	// Bind the OpenGL calls and dump some version info.
	gl.Init()
	fmt.Printf("%s %s", gl.GetString(gl.RENDERER), gl.GetString(gl.VERSION))
	fmt.Printf(" GLSL %s\n", gl.GetString(gl.SHADING_LANGUAGE_VERSION))

	// Gather the one scene into this one vertex array object.
	gl.GenVertexArrays(1, &tag.vao)
	gl.BindVertexArray(tag.vao)

	// create shaders
	tag.initShader()
	gl.UseProgram(tag.shaders)

	// vertex data.
	var vbuff uint32
	gl.GenBuffers(1, &vbuff)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbuff)
	gl.BufferData(gl.ARRAY_BUFFER, int64(len(tag.verticies)*4), gl.Pointer(&(tag.verticies[0])), gl.STATIC_DRAW)
	vattr := uint32(gl.GetAttribLocation(tag.shaders, "in_v"))
	gl.EnableVertexAttribArray(vattr)
	gl.VertexAttribPointer(vattr, 3, gl.FLOAT, false, 0, 0)

	// color data.
	var cbuff uint32
	gl.GenBuffers(1, &cbuff)
	gl.BindBuffer(gl.ARRAY_BUFFER, cbuff)
	gl.BufferData(gl.ARRAY_BUFFER, int64(len(tag.color)*4), gl.Pointer(&(tag.color[0])), gl.STATIC_DRAW)
	cattr := uint32(gl.GetAttribLocation(tag.shaders, "in_c"))
	gl.EnableVertexAttribArray(cattr)
	gl.VertexAttribPointer(cattr, 4, gl.FLOAT, false, 0, 0)

	// faces data.
	var ebuff uint32
	gl.GenBuffers(1, &ebuff)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebuff)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int64(len(tag.faces)), gl.Pointer(&(tag.faces[0])), gl.STATIC_DRAW)

	// set some state that doesn't need to change during drawing.
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
}

// initShader compiles shaders and links them into a shader program.
func (tag *trtag) initShader() {
	shader := &load.ShdData{}
	if err := shader.Load("basic", load.NewLocator()); err == nil {
		tag.shaders = gl.CreateProgram()
		if err := gl.BindProgram(tag.shaders, shader.Vsh, shader.Fsh); err != nil {
			fmt.Printf("Failed to create program: %s\n", err)
		}
		tag.mvpRef = gl.GetUniformLocation(tag.shaders, "mvpm")
		if tag.mvpRef < 0 {
			fmt.Printf("No model-view-projection matrix in vertex shader\n")
		}
	}
}

// drawScene renders the 3D models consisting of one VAO
func (tag *trtag) drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	tag.checkError("gl.Clear")
	gl.UseProgram(tag.shaders)
	tag.checkError("gl.UseProgram")
	gl.BindVertexArray(tag.vao)
	tag.checkError("gl.BindVertexArray")

	// Use a modelview matrix and quaternion to rotate the 3D object.
	tag.mvp64.SetQ(lin.NewQ().SetAa(0, 1, 0, lin.Rad(-tag.rotateAngle)))
	tag.mvp64.TranslateMT(0, 0, -4)
	tag.mvp.Set(tag.mvp64.Mult(tag.mvp64, tag.persp))
	gl.UniformMatrix4fv(tag.mvpRef, 1, false, tag.mvp.Pointer())
	if err := gl.GetError(); err != 0 {
		fmt.Printf("gl.UniformMatrix error %d\n", err)
	}
	gl.DrawElements(gl.TRIANGLES, int32(len(tag.faces)), gl.UNSIGNED_BYTE, 0)
	if err := gl.GetError(); err != 0 {
		fmt.Printf("gl.DrawElements error %d\n", err)
	}

	// cleanup
	gl.UseProgram(0)
	tag.checkError("gl.UseProgram-0")
	gl.BindVertexArray(0)
	tag.checkError("gl.BindVertexArray-0")

	// rotate based on time... not on how fast the computer runs.
	if time.Now().Sub(tag.lastTime).Seconds() > 0.01 {
		tag.rotateAngle++
		tag.lastTime = time.Now()
	}
}

// checkError helps to debug OpenGL errors by printing out when the occur.
func (tag *trtag) checkError(txt string) {
	cnt := 0
	err := gl.GetError()
	for err != 0 {
		fmt.Printf("%s error %d::%d\n", txt, cnt, err)
		err = gl.GetError()
		cnt++
		if cnt > 10 {
			os.Exit(-1)
		}
	}
}
