// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package render

// FUTURE: need a directx.go implementation to test the Render API and the
//         graphics layer encapsulation.

import (
	"errors"
	"fmt"
	"image"
	"log"
	"strings"

	"github.com/gazed/vu/render/gl"
)

// opengl is the OpenGL implemntation of Renderer.  See the Renderer interface
// for comments. Also see the OpenGL documentation for the individual calls.
type opengl struct {
	currentShader uint32 // Track the current shader to reduce shader switching.
}

// newRenderer returns an OpenGL implementation of Renderer.
func newRenderer() Renderer { return &opengl{} }

// Render implementation specific constants.
const (

	// Values useed in Renderer.Enable() method.
	BLEND      uint32 = gl.BLEND              // Alpha blending.
	CULL              = gl.CULL_FACE          // Backface culling.
	DEPTH             = gl.DEPTH_TEST         // Z-buffer (depth) awareness.
	POINT_SIZE        = gl.PROGRAM_POINT_SIZE // Enable gl_PointSize in shaders.

	// Vertex data render hints. Used in the Buffer.SetUsage() method.
	STATIC  = gl.STATIC_DRAW  // Data created once and rendered many times.
	DYNAMIC = gl.DYNAMIC_DRAW // Data is continually being updated.
)

// Internal package constants.
const (
	uShort = gl.UNSIGNED_SHORT // Unsigned short data type. Used in Buffer.
)

// Implements Renderer interface.
func (gc *opengl) Init() error {
	gl.Init()
	return gc.validate()
}

// Renderer implementation.
func (gc *opengl) NewModel(s Shader) Model            { return newModel(gc, s) }
func (gc *opengl) NewMesh(name string) Mesh           { return newMesh(name) }
func (gc *opengl) NewTexture(name string) Texture     { return newTexture(name) }
func (gc *opengl) NewShader(name string) Shader       { return newShader(name) }
func (gc *opengl) NewAnimation(name string) Animation { return newAnimation(name) }
func (gc *opengl) Color(r, g, b, a float32)           { gl.ClearColor(r, g, b, a) }
func (gc *opengl) Clear()                             { gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) }
func (gc *opengl) Viewport(width int, height int)     { gl.Viewport(0, 0, int32(width), int32(height)) }

// Renderer implementation.
func (gc *opengl) Enable(attribute uint32, enabled bool) {
	switch attribute {
	case CULL, DEPTH:
		if enabled {
			gl.Enable(attribute)
		} else {
			gl.Disable(attribute)
		}
	case BLEND:
		if enabled {
			gl.Enable(attribute)

			// Using non pre-multiplied alpha colour data so...
			gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
		} else {
			gl.Disable(attribute)
		}
	}
}

// Renderer implementation.
// FUTURE: all kinds of possible optimizations that would need to be
//         profiled before implementing.
//           • group by vao to avoid switching vao's.
//           • group by texture to avoid switching textures.
//           • group by Z/transparent to minimize sorting and switching depth mode.
//           • use interleaved vertex data.
//           • uniform buffers http://www.opengl.org/wiki/Uniform_Buffer_Object.
//           • ... lots more possiblities... leave your fav here.
func (gc *opengl) Render(mod Model) {
	m := mod.(*model)
	if m != nil && m.msh != nil && m.shd != nil {
		gl.Enable(gl.DEPTH_TEST)
		if m.is2D {
			gl.Disable(gl.DEPTH_TEST)
		}
		if m.cull {
			gl.Enable(gl.CULL_FACE)
		} else {
			gl.Disable(gl.CULL_FACE)
		}

		// switch shaders only if necessary.
		if m.shd.program != gc.currentShader {
			gl.UseProgram(m.shd.program)
			gc.currentShader = m.shd.program
		}

		// FUTURE: only need to bind uniforms that have changed.
		m.bindUniforms() // currently bind each time.
		if m.msh.rebind {
			gc.bindMesh(m.msh)
			m.msh.rebind = false
		}

		// bind the data buffers and render.
		gl.BindVertexArray(m.msh.vao)
		switch m.mode {
		case LINES:
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
			fd := m.msh.faces
			gl.DrawElements(gl.LINES, int32(len(fd.data)), gl.UNSIGNED_SHORT, 0)
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		case POINTS:
			gl.Enable(gl.PROGRAM_POINT_SIZE)
			gl.DrawArrays(gl.POINTS, 0, m.msh.numv)
			gl.Disable(gl.PROGRAM_POINT_SIZE)
		case TRIANGLES:
			fd := m.msh.faces
			if len(m.tex) > 1 && m.tex[0].fn > 0 {

				// Some models have multiple texture maps applied to a single set of
				// vertex data.
				for cnt, tex := range m.tex {

					// Use the same texture unit and sampler. Just update which
					// image is being sampled.
					gl.BindTexture(gl.TEXTURE_2D, m.tex[cnt].tid)

					// fn is the number of triangles, 3 indicies per triangle.
					// f0 is the offset in triangles where each triangle has 3 indicies
					//    of 2 bytes (uShort) each.
					gl.DrawElements(gl.TRIANGLES, tex.fn*3, gl.UNSIGNED_SHORT, int64(3*2*tex.f0))
				}
			} else {
				if fd == nil {
					log.Printf("opengl:Render mesh %s has no data", m.msh.name)
				} else {
					gl.DrawElements(gl.TRIANGLES, int32(len(fd.data)), gl.UNSIGNED_SHORT, 0)
				}
			}
		}
		gl.Disable(gl.DEPTH_TEST)
	}
}

// validate that OpenGL is available at the right version. For OpenGL 3.2
// the following lines should be in the report.
//	       [+] glFramebufferTexture
//	       [+] glGetBufferParameteri64v
//	       [+] glGetInteger64i_v
func (gc *opengl) validate() error {
	if report := gl.BindingReport(); len(report) > 0 {
		valid := false
		want := "[+] glFramebufferTexture"
		for _, line := range report {
			if strings.Contains(line, want) {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("Need OpenGL 3.2 or higher.")
		}
	} else {
		return fmt.Errorf("OpenGL unavailable.")
	}
	return nil
}

// bindMesh copies the existing mesh data to the GPU and initializes the vao and
// buffer references.
func (gc *opengl) bindMesh(msh Mesh) error {
	m := msh.(*mesh)
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		return fmt.Errorf("BindMesh needs to find and fix prior error %X", glerr)
	}

	// Reuse existing vao's.
	if m.vao == 0 {
		gl.GenVertexArrays(1, &(m.vao))
	}
	gl.BindVertexArray(m.vao)
	for _, vd := range m.vdata {
		if vd.rebind {
			gc.bindVertexBuffer(vd)
			vd.rebind = false
		}
	}
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		return fmt.Errorf("BindMesh failed to bind vb %s %X", m.name, glerr)
	}
	if m.faces != nil {
		if m.faces.rebind {
			gc.bindFaceBuffer(m.faces)
			m.faces.rebind = false
		}
	}
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		return fmt.Errorf("BindMesh failed to bind fb %s %X", m.name, glerr)
	}
	return nil
}

// bind the buffer data to the GPU.
func (gc *opengl) bindVertexBuffer(vd *vertexData) {
	if vd.ref == 0 {
		gl.GenBuffers(1, &vd.ref)
	}
	bytes := 4 // 4 bytes for float32 (gl.FLOAT)
	switch vd.usage {
	case STATIC:
		switch vd.dtype {
		case floatData:
			gl.BindBuffer(gl.ARRAY_BUFFER, vd.ref)
			gl.BufferData(gl.ARRAY_BUFFER, int64(len(vd.floats)*bytes), gl.Pointer(&(vd.floats[0])), vd.usage)
			gl.VertexAttribPointer(vd.lloc, vd.span, gl.FLOAT, false, 0, 0)
		case byteData:
			gl.BindBuffer(gl.ARRAY_BUFFER, vd.ref)
			gl.BufferData(gl.ARRAY_BUFFER, int64(len(vd.bytes)), gl.Pointer(&(vd.bytes[0])), vd.usage)
			gl.VertexAttribPointer(vd.lloc, vd.span, gl.UNSIGNED_BYTE, vd.normalize, 0, 0)
		}
	case DYNAMIC:
		null := gl.Pointer(uintptr(0))
		switch vd.dtype {
		case floatData:
			gl.BindBuffer(gl.ARRAY_BUFFER, vd.ref)

			// Buffer orphaning, a common way to improve streaming perf. See:
			//         http://www.opengl.org/wiki/Buffer_Object_Streaming
			gl.BufferData(gl.ARRAY_BUFFER, int64(cap(vd.floats)*bytes), null, vd.usage)
			gl.BufferSubData(gl.ARRAY_BUFFER, 0, int64(len(vd.floats)*bytes), gl.Pointer(&(vd.floats[0])))
			gl.VertexAttribPointer(vd.lloc, vd.span, gl.FLOAT, false, 0, 0)
		}
	}
	gl.EnableVertexAttribArray(vd.lloc)
}

// bind the buffer data to the GPU.
func (gc *opengl) bindFaceBuffer(fd *faceData) {
	if len(fd.data) > 0 {
		if fd.ref == 0 {
			gl.GenBuffers(1, &fd.ref)
		}
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, fd.ref)
		bytes := 2 // 2 bytes for uint16 (gl.UNSIGNED_SHORT)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int64(len(fd.data)*bytes), gl.Pointer(&(fd.data[0])), fd.usage)
	}
}

// bindShader compiles the shader and makes it available to the GPU.
func (gc *opengl) bindShader(shdr Shader) (err error) {
	s := shdr.(*shader)
	s.ensureNewLines()
	s.program = gl.CreateProgram()

	// compile and link the shader program.
	if glerr := gl.BindProgram(s.program, s.vsh, s.fsh); glerr != nil {
		err = fmt.Errorf("Failed to create shader program: %s", glerr)
		return
	}

	// initialize the uniform references
	var errmsg string
	for label, _ := range s.uniforms {
		if uid := gl.GetUniformLocation(s.program, label); uid >= 0 {
			s.uniforms[label] = uid
		} else {
			errnum := gl.GetError()
			errmsg += fmt.Sprintf("No %s uniform in shader %X", label, errnum)
		}
	}
	if len(errmsg) > 0 {
		err = errors.New(errmsg)
	}
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		log.Printf("shader:Bind need to find and fix error %X", glerr)
	}
	return
}

// bindTexture makes the texture available on the GPU.
func (gc *opengl) bindTexture(tex Texture) (err error) {
	t := tex.(*texture)
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		log.Printf("opengl:bindTexture need to find and fix prior error %X", glerr)
	}

	if t.tid == 0 {
		gl.GenTextures(1, &(t.tid))
	}
	gl.BindTexture(gl.TEXTURE_2D, t.tid)

	// FUTURE check if RGBA, or NRGBA are alpha pre-multiplied. The docs say yes
	// for RGBA but the data is read from PNG files which are not pre-multiplied
	// and the go png Decode looks like its reading values directly.
	var ptr gl.Pointer
	bounds := t.img.Bounds()
	width, height := int32(bounds.Dx()), int32(bounds.Dy())
	switch imgType := t.img.(type) {
	case *image.RGBA:
		i := t.img.(*image.RGBA)
		ptr = gl.Pointer(&(i.Pix[0]))
	case *image.NRGBA:
		i := t.img.(*image.NRGBA)
		ptr = gl.Pointer(&(i.Pix[0]))
	default:
		return fmt.Errorf("Unsupported image format", imgType)
	}
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, ptr)
	gl.GenerateMipmap(gl.TEXTURE_2D)
	gc.updateTextureMode(tex)
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		err = fmt.Errorf("Failed binding texture %s\n", t.name)
	}
	return
}

// updateTextureMode
func (gc *opengl) updateTextureMode(tex Texture) {
	t := tex.(*texture)
	gl.BindTexture(gl.TEXTURE_2D, t.tid)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAX_LEVEL, 7)
	if t.repeat {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST_MIPMAP_LINEAR)
}

// bindUniform links data to uniforms expected by shaders.
// It expects the variable parameter list types to match the uniform type.
func (gc *opengl) bindUniform(uniform int32, utype, cnt int, udata ...interface{}) {
	switch utype {
	case i1:
		i1 := udata[0].(int32)
		gl.Uniform1i(uniform, i1)
	case f1:
		f1 := udata[0].(float32)
		gl.Uniform1f(uniform, f1)
	case f2:
		f1 := udata[0].(float32)
		f2 := udata[1].(float32)
		gl.Uniform2f(uniform, f1, f2)
	case f3:
		f1 := udata[0].(float32)
		f2 := udata[1].(float32)
		f3 := udata[2].(float32)
		gl.Uniform3f(uniform, f1, f2, f3)
	case f4:
		f1 := udata[0].(float32)
		f2 := udata[1].(float32)
		f3 := udata[2].(float32)
		f4 := udata[3].(float32)
		gl.Uniform4f(uniform, f1, f2, f3, f4)
	case x3:
		mptr := udata[0].(*float32)
		gl.UniformMatrix3fv(uniform, int32(cnt), false, mptr)
	case x34:
		mptr := udata[0].(*float32)
		gl.UniformMatrix3x4fv(uniform, int32(cnt), false, mptr)
	case x4:
		mptr := udata[0].(*float32)
		gl.UniformMatrix4fv(uniform, int32(cnt), false, mptr)
	}
}

// Current list of supported uniform types.
const (
	i1  = iota // glUniform1i
	f1         // glUniform1f
	f2         // glUniform2f
	f3         // glUniform3f
	f4         // glUniform4f
	x3         // glUniformMatrix3fv
	x34        // glUniformMatrix3x4fv
	x4         // glUniformMatrix4fv
)

// useTexture makes the given texture the active texture.
func (gc *opengl) useTexture(sampler, texUnit int32, tex Texture) {
	t := tex.(*texture)
	gc.bindUniform(sampler, i1, 1, texUnit)
	gl.ActiveTexture(gl.TEXTURE0 + uint32(texUnit))
	gl.BindTexture(gl.TEXTURE_2D, t.tid)
}

// Remove graphic resources.
func (gc *opengl) deleteMesh(mid uint32)    { gl.DeleteVertexArrays(1, &mid) }
func (gc *opengl) deleteShader(sid uint32)  { gl.DeleteProgram(sid) }
func (gc *opengl) deleteTexture(tid uint32) { gl.DeleteTextures(1, &tid) }
