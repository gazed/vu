// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// +build !dx
// Use opengl by default.

package render

import (
	"fmt"
	"image"
	"log"
	"strings"

	"github.com/gazed/vu/render/gl"
)

// opengl is the OpenGL implemntation of Renderer.  See the Renderer interface
// for comments. See the OpenGL documentation for OpenGL methods and constants.
type opengl struct {
	depthTest bool   // Track current depth setting to reduce state switching.
	shader    uint32 // Track the current shader to reduce shader switching.
}

// newRenderer returns an OpenGL implementation of Renderer.
func newRenderer() Renderer {
	gc := &opengl{}
	return gc
}

// Renderer implementation specific constants.
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

// Renderer implementation.
func (gc *opengl) Init() error {
	gl.Init()
	return gc.validate()
}

// Renderer implementation.
func (gc *opengl) Color(r, g, b, a float32)       { gl.ClearColor(r, g, b, a) }
func (gc *opengl) Clear()                         { gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) }
func (gc *opengl) Viewport(width int, height int) { gl.Viewport(0, 0, int32(width), int32(height)) }

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

// Render implementation.
// FUTURE: all kinds of possible optimizations that would need to be
//         profiled before implementing.
//           • group by vao to avoid switching vao's.
//           • group by texture to avoid switching textures.
//           • use interleaved vertex data.
//           • uniform buffers http://www.opengl.org/wiki/Uniform_Buffer_Object.
//           • ... lots more possiblities... leave your fav here.
func (gc *opengl) Render(dr Draw) {
	d, ok := dr.(*draw)
	if !ok || d == nil {
		return
	}

	// switch state only if necessary.
	if gc.depthTest != d.depth {
		if d.depth {
			gl.Enable(gl.DEPTH_TEST)
		} else {
			gl.Disable(gl.DEPTH_TEST)
		}
		gc.depthTest = d.depth
	}

	// switch shaders only if necessary.
	if gc.shader != d.shader {
		gl.UseProgram(d.shader)
		gc.shader = d.shader
	}

	// Ask the model to bind its provisioned uniforms.
	// FUTURE: only need to bind uniforms that have changed.
	gc.bindUniforms(d)

	// bind the data buffers and render.
	gl.BindVertexArray(d.vao)
	switch d.mode {
	case LINES:
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		gl.DrawElements(gl.LINES, d.numFaces, gl.UNSIGNED_SHORT, 0)
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	case POINTS:
		gl.Enable(gl.PROGRAM_POINT_SIZE)
		gl.DrawArrays(gl.POINTS, 0, d.numVerts)
		gl.Disable(gl.PROGRAM_POINT_SIZE)
	case TRIANGLES:
		if len(d.texs) > 1 && d.texs[0].fn > 0 {
			// Multiple textures on one model specify which verticies they apply to.
			for _, tex := range d.texs {
				// Use the same texture unit and sampler. Just update which
				// image is being sampled.
				gl.BindTexture(gl.TEXTURE_2D, tex.tid)
				// fn is the number of triangles, 3 indicies per triangle.
				// f0 is the offset in triangles where each triangle has 3 indicies
				//    of 2 bytes (uShort) each.
				gl.DrawElements(gl.TRIANGLES, tex.fn*3, gl.UNSIGNED_SHORT, int64(3*2*tex.f0))
			}
		} else {
			// Single textures are handled with a standard bindUniforms
			gl.DrawElements(gl.TRIANGLES, d.numFaces, gl.UNSIGNED_SHORT, 0)
		}
	}
}

// bindUniforms links model data to the uniforms discovered
// in the model shader.
func (gc *opengl) bindUniforms(d *draw) {
	for key, ref := range d.uniforms {
		switch key {
		case "mvpm":
			gc.bindUniform(ref, x4, 1, d.mvp.Pointer())
		case "mvm":
			gc.bindUniform(ref, x4, 1, d.mv.Pointer())
		case "pm":
			gc.bindUniform(ref, x4, 1, d.pm.Pointer())
		case "nm":
			nm := (&m3{}).m3(d.mv) // normal matrix as subset of model-view.
			gc.bindUniform(ref, x3, 1, nm.Pointer())
		case "uv":
			gc.useTexture(ref, 0, d.texs[0].tid)
		case "uv0":
			gc.useTexture(ref, 0, d.texs[0].tid)
		case "uv1":
			gc.useTexture(ref, 1, d.texs[1].tid)
		case "uv2":
			gc.useTexture(ref, 2, d.texs[2].tid)
		case "uv3":
			gc.useTexture(ref, 3, d.texs[3].tid)
		case "uv4":
			gc.useTexture(ref, 4, d.texs[4].tid)
		case "uv5":
			gc.useTexture(ref, 5, d.texs[5].tid)
		case "uv6":
			gc.useTexture(ref, 6, d.texs[6].tid)
		case "uv7":
			gc.useTexture(ref, 7, d.texs[7].tid)
		case "uv8":
			gc.useTexture(ref, 8, d.texs[8].tid)
		case "uv9":
			gc.useTexture(ref, 9, d.texs[9].tid)
		case "uv10":
			gc.useTexture(ref, 10, d.texs[10].tid)
		case "uv11":
			gc.useTexture(ref, 11, d.texs[11].tid)
		case "uv12":
			gc.useTexture(ref, 12, d.texs[12].tid)
		case "uv13":
			gc.useTexture(ref, 13, d.texs[13].tid)
		case "uv14":
			gc.useTexture(ref, 14, d.texs[14].tid)
		case "uv15":
			gc.useTexture(ref, 15, d.texs[15].tid)
		case "scale":
			gc.bindUniform(ref, f3, 1, d.scale.x, d.scale.y, d.scale.z)
		case "alpha":
			gc.bindUniform(ref, f1, 1, d.alpha)
		case "time":
			gc.bindUniform(ref, f1, 1, d.time)
		case "bpos": // bone position animation data.
			if d.pose != nil && len(d.pose) > 0 {
				gc.bindUniform(ref, x34, len(d.pose), d.pose[0].Pointer())
			}
		default:
			// bind non-standard float based uniforms.
			// known examples are "l", "ld", "kd", "ks", "ka"
			if floats, ok := d.floats[key]; ok {
				switch len(floats) {
				case 1:
					gc.bindUniform(ref, f1, 1, floats[0])
				case 2:
					gc.bindUniform(ref, f2, 1, floats[0], floats[1])
				case 3:
					gc.bindUniform(ref, f3, 1, floats[0], floats[1], floats[2])
				case 4:
					gc.bindUniform(ref, f4, 1, floats[0], floats[1], floats[2], floats[3])
				}
			} else {
				log.Printf("No uniform bound for %s", key)
			}
		}
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

// Renderer implementation.
// BindMesh copies the given mesh data to the GPU
// and initializes the vao and buffer references.
func (gc *opengl) BindMesh(vao *uint32, vdata map[uint32]Data, fdata Data) error {
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		return fmt.Errorf("BindMesh needs to find and fix prior error %X", glerr)
	}

	// Reuse existing vao's.
	if *vao == 0 {
		gl.GenVertexArrays(1, vao)
	}
	gl.BindVertexArray(*vao)
	for _, vbuff := range vdata {
		vd, ok := vbuff.(*vertexData)
		if ok && vd.rebind {
			gc.bindVertexBuffer(vd)
			vd.rebind = false
		}
	}
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		return fmt.Errorf("BindMesh failed to bind vb %X", glerr)
	}
	if fd, ok := fdata.(*faceData); ok {
		if fd.rebind {
			gc.bindFaceBuffer(fd)
			fd.rebind = false
		}
	}
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		return fmt.Errorf("BindMesh failed to bind fb %X", glerr)
	}
	return nil
}

// bindVertexBuffer copies per-vertex data from the CPU to the GPU.
func (gc *opengl) bindVertexBuffer(vdata Data) {
	vd, ok := vdata.(*vertexData)
	if !ok {
		return
	}
	if vd.ref == 0 {
		gl.GenBuffers(1, &vd.ref)
	}
	bytes := 4 // 4 bytes for float32 (gl.FLOAT)
	switch vd.usage {
	case STATIC:
		switch {
		case len(vd.floats) > 0:
			gl.BindBuffer(gl.ARRAY_BUFFER, vd.ref)
			gl.BufferData(gl.ARRAY_BUFFER, int64(len(vd.floats)*bytes), gl.Pointer(&(vd.floats[0])), vd.usage)
			gl.VertexAttribPointer(vd.lloc, vd.span, gl.FLOAT, false, 0, 0)
		case len(vd.bytes) > 0:
			gl.BindBuffer(gl.ARRAY_BUFFER, vd.ref)
			gl.BufferData(gl.ARRAY_BUFFER, int64(len(vd.bytes)), gl.Pointer(&(vd.bytes[0])), vd.usage)
			gl.VertexAttribPointer(vd.lloc, vd.span, gl.UNSIGNED_BYTE, vd.normalize, 0, 0)
		}
	case DYNAMIC:
		var null gl.Pointer // zero.
		switch {
		case len(vd.floats) > 0:
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

// bindFaceBuffer copies triangle face data from the CPU to the GPU.
func (gc *opengl) bindFaceBuffer(fdata Data) {
	fd := fdata.(*faceData)
	if len(fd.data) > 0 {
		if fd.ref == 0 {
			gl.GenBuffers(1, &fd.ref)
		}
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, fd.ref)
		bytes := 2 // 2 bytes for uint16 (gl.UNSIGNED_SHORT)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int64(len(fd.data)*bytes), gl.Pointer(&(fd.data[0])), fd.usage)
	}
}

// Renderer implementation.
// BindShader compiles the shader and makes it available to the GPU.
// It also adds the list of uniforms and vertex layout references to the
// provided maps.
func (gc *opengl) BindShader(vsh, fsh []string, uniforms map[string]int32,
	layouts map[string]uint32) (program uint32, err error) {
	program = gl.CreateProgram()

	// compile and link the shader program.
	if glerr := gl.BindProgram(program, vsh, fsh); glerr != nil {
		err = fmt.Errorf("Failed to create shader program: %s", glerr)
		return
	}

	// initialize the uniform and layout references
	gl.Uniforms(program, uniforms)
	gl.Layouts(program, layouts)
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		log.Printf("shader:Bind need to find and fix error %X", glerr)
	}
	return
}

// Renderer implementation.
// BindTexture makes the texture available on the GPU.
func (gc *opengl) BindTexture(tid *uint32, img image.Image, repeat bool) (err error) {
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		log.Printf("opengl:bindTexture need to find and fix prior error %X", glerr)
	}
	if *tid == 0 {
		gl.GenTextures(1, tid)
	}
	gl.BindTexture(gl.TEXTURE_2D, *tid)

	// FUTURE: check if RGBA, or NRGBA are alpha pre-multiplied. The docs say yes
	// for RGBA but the data is from PNG files which are not pre-multiplied
	// and the go png Decode looks like its reading values directly.
	var ptr gl.Pointer
	bounds := img.Bounds()
	width, height := int32(bounds.Dx()), int32(bounds.Dy())
	switch imgType := img.(type) {
	case *image.RGBA:
		i := img.(*image.RGBA)
		ptr = gl.Pointer(&(i.Pix[0]))
	case *image.NRGBA:
		i := img.(*image.NRGBA)
		ptr = gl.Pointer(&(i.Pix[0]))
	default:
		return fmt.Errorf("Unsupported image format %T", imgType)
	}
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, ptr)
	gl.GenerateMipmap(gl.TEXTURE_2D)
	gc.setTextureMode(*tid, repeat)
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		err = fmt.Errorf("Failed binding texture %d\n", glerr)
	}
	return
}

// setTextureMode is used to switch to a repeating
// texture instead of a 1:1 texture mapping.
func (gc *opengl) setTextureMode(tid uint32, repeat bool) {
	gl.BindTexture(gl.TEXTURE_2D, tid)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAX_LEVEL, 7)
	if repeat {
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
func (gc *opengl) useTexture(sampler, texUnit int32, tid uint32) {
	gc.bindUniform(sampler, i1, 1, texUnit)
	gl.ActiveTexture(gl.TEXTURE0 + uint32(texUnit))
	gl.BindTexture(gl.TEXTURE_2D, tid)
}

// Remove graphic resources.
func (gc *opengl) ReleaseMesh(vao uint32)    { gl.DeleteVertexArrays(1, &vao) }
func (gc *opengl) ReleaseShader(sid uint32)  { gl.DeleteProgram(sid) }
func (gc *opengl) ReleaseTexture(tid uint32) { gl.DeleteTextures(1, &tid) }
