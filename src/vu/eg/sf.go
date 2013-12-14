// Use is governed by a FreeBSD license found in the LICENSE file.
// Copyright © 2013 Galvanized Logic Inc.

package main

import (
	"fmt"
	"time"
	"vu/data"
	"vu/device"
	"vu/math/lin"
	"vu/render"
	"vu/render/gl"
)

// sf demonstrates one example of shader only rendering. This shows the power of shaders using
// an example from shadertoy.com. Specifically:
//       https://www.shadertoy.com/view/Xsl3zN
// For more shader examples also check out:
//       http://glsl.heroku.com
// The real star of this demo though is found in ./shaders/fire.fsh.  Kudos to @301z and
// the other contributors to shadertoy and heroku.
//
// This also demonstrates basic rendering by using OpenGL calls from
// package "vu/render/gl" for rendering.
func sf() {
	sf := new(sftag)
	dev := device.New("Shader Fire", 400, 100, 500, 500)
	sf.initScene()
	dev.Open()
	for dev.IsAlive() {
		sf.update(dev)
		sf.drawScene()
		dev.SwapBuffers()
	}
	dev.Dispose()
}

// Globally unique "tag" that encapsulates example specific data.
type sftag struct {
	vao     uint32
	sTime   time.Time  // start time.
	gTime   int32      // uniform reference to time in seconds since startup.
	sizes   int32      // uniform reference to the viewport sizes vector.
	shaders uint32     // program reference.
	ortho   *lin.M4    // perspective matrix.
	mvp32   *render.M4 // passed to graphics layer.
	mvpref  int32      // mvp uniform id

	// mesh information
	points []float32
	faces  []uint8
}

// update handles user input.
func (sf *sftag) update(dev device.Device) {
	pressed := dev.Update()
	if pressed.Resized {
		sf.resize(dev.Size())
	}
}

// resize handles user screen/window changes.
func (sf *sftag) resize(x, y, width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

// initScene is one time initialization that creates a single VAO
func (sf *sftag) initScene() {
	sf.mvp32 = &render.M4{}
	sf.sTime = time.Now()
	sf.initData()

	// Bind the OpenGL calls and dump some version info.
	gl.Init()
	fmt.Printf("%s %s", gl.GetString(gl.RENDERER), gl.GetString(gl.VERSION))
	fmt.Printf(" GLSL %s\n", gl.GetString(gl.SHADING_LANGUAGE_VERSION))

	gl.GenVertexArrays(1, &sf.vao)
	gl.BindVertexArray(sf.vao)

	// vertex data.
	var vbuff uint32
	gl.GenBuffers(1, &vbuff)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbuff)
	gl.BufferData(gl.ARRAY_BUFFER, int64(len(sf.points)*4), gl.Pointer(&(sf.points[0])), gl.STATIC_DRAW)
	var vattr uint32 = 0
	gl.VertexAttribPointer(vattr, 4, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(vattr)

	// faces data.
	var ebuff uint32
	gl.GenBuffers(1, &ebuff)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebuff)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int64(len(sf.faces)), gl.Pointer(&(sf.faces[0])), gl.STATIC_DRAW)

	// create texture and shaders after all the data has been set up.
	shader := &data.Shader{}
	loader := data.NewLoader()
	loader.Load("fire", &shader)
	sf.shaders = gl.CreateProgram()
	if err := gl.BindProgram(sf.shaders, shader.Vsh, shader.Fsh); err != nil {
		fmt.Printf("Failed to create program: %s\n", err)
	}
	sf.mvpref = gl.GetUniformLocation(sf.shaders, "Mvpm")
	sf.gTime = gl.GetUniformLocation(sf.shaders, "time")
	sf.sizes = gl.GetUniformLocation(sf.shaders, "resolution")
	sf.ortho = lin.NewOrtho(0, 4, 0, 4, 0, 10)

	// set some state that doesn't need to change during drawing.
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
}

// initData creates a flat mesh that the shader renders onto.
func (sf *sftag) initData() {
	sf.points = []float32{
		0, 0, 0, 1,
		4, 0, 0, 1,
		0, 4, 0, 1,
		4, 4, 0, 1,
	}
	sf.faces = []uint8{
		0, 2, 1,
		1, 2, 3,
	}
}

// drawScene renders the shader-only scene.
func (sf *sftag) drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(sf.shaders)
	gl.BindVertexArray(sf.vao)
	timeSinceStart := time.Since(sf.sTime).Seconds()
	gl.Uniform1f(sf.gTime, float32(timeSinceStart))
	gl.Uniform2f(sf.sizes, 500, 500)
	sf.mvp32 = renderMatrix(sf.ortho, sf.mvp32)
	gl.UniformMatrix4fv(sf.mvpref, 1, false, sf.mvp32.Pointer())
	gl.DrawElements(gl.TRIANGLES, int32(len(sf.faces)), gl.UNSIGNED_BYTE, gl.Pointer(nil))

	// cleanup
	gl.UseProgram(0)
	gl.BindVertexArray(0)
}
