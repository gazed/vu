// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/gazed/vu/device"
	"github.com/gazed/vu/load"
	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/render"
	"github.com/gazed/vu/render/gl"
)

// tb shows how a basic texture is used in OpenGL.  One texture is loaded and
// rendered on a single mesh. This example is useful in understanding texture
// basics.
//
// This example renders using OpenGL calls from package vu/render/gl.
func tb() {
	tb := new(tbtag)
	dev := device.New("Texture:Basic", 400, 100, 500, 500)
	tb.initScene()
	dev.Open()
	for dev.IsAlive() {
		tb.update(dev)
		tb.drawScene()
		dev.SwapBuffers()
	}
	dev.Dispose()
}

// Globally unique "tag" that encapsulates example specific data.
type tbtag struct {
	shaders uint32
	vao     uint32
	mvpref  int32          // mvp uniform id
	mvp     render.Mvp     // transform matrix for rendering.
	sampler int32          // sampler uniform id
	texture render.Texture // the picture to show.
	tid     uint32         // texture reference.

	// mesh information
	verticies []float32
	faces     []uint8
	tcoords   []float32
}

// update handles user input.
func (tb *tbtag) update(dev device.Device) {
	pressed := dev.Update()
	if pressed.Resized {
		tb.resize(dev.Size())
	}
}

// resize handles user screen/window changes.
func (tb *tbtag) resize(x, y, width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

// initScene is one time initialization that creates a single VAO
func (tb *tbtag) initScene() {
	tb.initData()

	// Bind the OpenGL calls and dump some version info.
	gl.Init()
	fmt.Printf("%s %s", gl.GetString(gl.RENDERER), gl.GetString(gl.VERSION))
	fmt.Printf(" GLSL %s\n", gl.GetString(gl.SHADING_LANGUAGE_VERSION))

	// Gather the one scene into this one vertex array object.
	gl.GenVertexArrays(1, &tb.vao)
	gl.BindVertexArray(tb.vao)

	// vertex data.
	var vbuff uint32
	gl.GenBuffers(1, &vbuff)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbuff)
	gl.BufferData(gl.ARRAY_BUFFER, int64(len(tb.verticies)*4), gl.Pointer(&(tb.verticies[0])), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(0)

	// texture coordatinates
	var tbuff uint32
	gl.GenBuffers(1, &tbuff)
	gl.BindBuffer(gl.ARRAY_BUFFER, tbuff)
	gl.BufferData(gl.ARRAY_BUFFER, int64(len(tb.tcoords)*4), gl.Pointer(&(tb.tcoords[0])), gl.STATIC_DRAW)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(2)

	// faces data.
	var ebuff uint32
	gl.GenBuffers(1, &ebuff)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebuff)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int64(len(tb.faces)), gl.Pointer(&(tb.faces[0])), gl.STATIC_DRAW)

	// create texture and shaders after all the data has been set up.
	tb.initTexture()
	tb.initShader()
	tb.mvp = render.NewMvp().Set(lin.NewM4().Ortho(0, 4, 0, 4, 0, 10))

	// set some state that doesn't need to change during drawing.
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.TEXTURE_2D)
}

// initData creates a flat mesh that the texture is drawn onto.
func (tb *tbtag) initData() {
	tb.verticies = []float32{
		1, 3, 0,
		3, 3, 0,
		1, 1, 0,
		3, 1, 0,
	}
	tb.tcoords = []float32{
		0, 0,
		1, 0,
		0, 1,
		1, 1,
	}
	tb.faces = []uint8{
		0, 2, 1,
		1, 2, 3,
	}
}

// initTexture loads the texture and binds it to the graphics device.
func (tb *tbtag) initTexture() {
	renderer := render.New()
	texture := renderer.NewTexture("image")
	ld := load.NewLoader()
	if img, err := ld.Png(texture.Name()); err == nil {
		texture.Set(img)
		tb.texture = texture
		gl.GenTextures(1, &tb.tid)
		gl.BindTexture(gl.TEXTURE_2D, tb.tid)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)

		// ensure image is in RGBA format
		b := texture.Img().Bounds()
		rgba := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(rgba, rgba.Bounds(), texture.Img(), b.Min, draw.Src)
		width, height := int32(b.Dx()), int32(b.Dy())
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Pointer(&(rgba.Pix[0])))
		if glerr := gl.GetError(); glerr != gl.NO_ERROR {
			fmt.Printf("Failed binding texture image.png\n")
		}
	} else {
		fmt.Println("Could not load image.png file.")
	}

}

// initShader compiles shaders and links them into a shader program.
func (tb *tbtag) initShader() {
	renderer := render.New()
	shader := renderer.NewShader("tuv")
	loader := load.NewLoader()
	vsrc, verr := loader.Vsh(shader.Name())
	fsrc, ferr := loader.Fsh(shader.Name())
	if verr == nil && ferr == nil {
		shader.SetSource(vsrc, fsrc)
		tb.shaders = gl.CreateProgram()
		if err := gl.BindProgram(tb.shaders, shader.Vsh(), shader.Fsh()); err != nil {
			fmt.Printf("Failed to create program: %s\n", err)
		}
		tb.mvpref = gl.GetUniformLocation(tb.shaders, "mvpm")
		tb.sampler = gl.GetUniformLocation(tb.shaders, "uv")
		if tb.mvpref < 0 {
			fmt.Printf("No model-view-projection matrix in vertex shader\n")
		}
	}
}

// drawScene renders the scene consisting of one VAO.
func (tb *tbtag) drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(tb.shaders)
	gl.Uniform1i(tb.sampler, 0)
	gl.ActiveTexture(gl.TEXTURE0 + 0)
	gl.BindTexture(gl.TEXTURE_2D, tb.tid)
	gl.BindVertexArray(tb.vao)
	gl.UniformMatrix4fv(tb.mvpref, 1, false, tb.mvp.Pointer())
	gl.DrawElements(gl.TRIANGLES, int32(len(tb.faces)), gl.UNSIGNED_BYTE, 0)

	// cleanup
	gl.ActiveTexture(0)
	gl.UseProgram(0)
	gl.BindVertexArray(0)
}
