// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"fmt"
	"time"

	"github.com/gazed/vu/device"
	"github.com/gazed/vu/load"
	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/render"
	"github.com/gazed/vu/render/gl"
)

// sf demonstrates one example of shader only rendering. This shows the power
// of shaders using an example from shadertoy.com. Specifically:
//       https://www.shadertoy.com/view/Xsl3zN
// For more shader examples also check out:
//       http://glsl.heroku.com
// The real star of this demo is found in ./source/fire.fsh. Kudos to @301z
// and the other contributors to shadertoy and heroku.
// This example renders using OpenGL calls from package vu/render/gl.
// See other examples and vu/shader.go for engine supported shaders in action.
//
// CONTROLS: NA
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
	mvp     render.Mvp // transform matrix for rendering.
	mvpref  int32      // mvp uniform id

	// mesh information
	verticies []float32
	faces     []uint8
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
	gl.BufferData(gl.ARRAY_BUFFER, int64(len(sf.verticies)*4), gl.Pointer(&(sf.verticies[0])), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(0)

	// faces data.
	var ebuff uint32
	gl.GenBuffers(1, &ebuff)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebuff)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int64(len(sf.faces)), gl.Pointer(&(sf.faces[0])), gl.STATIC_DRAW)

	// create texture and shaders after all the data has been set up.
	shader := &load.ShdData{}
	if err := shader.Load("fire", load.NewLocator()); err == nil {
		sf.shaders = gl.CreateProgram()
		if err := gl.BindProgram(sf.shaders, shader.Vsh, shader.Fsh); err != nil {
			fmt.Printf("Failed to create program: %s\n", err)
		}
		sf.mvpref = gl.GetUniformLocation(sf.shaders, "mvpm")
		sf.gTime = gl.GetUniformLocation(sf.shaders, "time")
		sf.sizes = gl.GetUniformLocation(sf.shaders, "screen")
		sf.mvp = render.NewMvp().Set(lin.NewM4().Ortho(0, 4, 0, 4, 0, 10))

		// set some state that doesn't need to change during drawing.
		gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	}
}

// initData creates a flat mesh that the shader renders onto.
func (sf *sftag) initData() {
	sf.verticies = []float32{
		0, 0, 0,
		4, 0, 0,
		0, 4, 0,
		4, 4, 0,
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
	gl.UniformMatrix4fv(sf.mvpref, 1, false, sf.mvp.Pointer())
	gl.DrawElements(gl.TRIANGLES, int32(len(sf.faces)), gl.UNSIGNED_BYTE, 0)

	// cleanup
	gl.UseProgram(0)
	gl.BindVertexArray(0)
}
