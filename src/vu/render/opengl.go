// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package render

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"log"
	"strings"
	"time"
	"vu/data"
	"vu/render/gl"
)

// opengl implements the Renderer interface wrapping all OpenGL calls.
// See the Renderer interface for comments.  Also see the OpenGL documentation
// for the individual calls.
type opengl struct{}

// Implements Renderer interface.
func (gc *opengl) Init() error {
	gl.Init()
	return gc.validate()
}

// Implements Renderer interface.
func (gc *opengl) Color(r, g, b, a float32)       { gl.ClearColor(r, g, b, a) }
func (gc *opengl) Clear()                         { gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) }
func (gc *opengl) Viewport(width int, height int) { gl.Viewport(0, 0, int32(width), int32(height)) }

// Implements Renderer interface.
func (gc *opengl) Enable(attribute int, enabled bool) {
	switch attribute {
	case BLEND:
		if enabled {
			gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
			gl.Enable(gl.BLEND)
		} else {
			gl.Disable(gl.BLEND)
		}
	case CULL:
		if enabled {
			gl.Enable(gl.CULL_FACE)
		} else {
			gl.Disable(gl.CULL_FACE)
		}
	case DEPTH:
		if enabled {
			gl.Enable(gl.DEPTH_TEST)
		} else {
			gl.Disable(gl.DEPTH_TEST)
		}
	}
}

// Implements Renderer interface.
func (gc *opengl) Render(v *Visible) {
	if v.Mesh != nil && v.Shader != nil {
		gl.BindVertexArray(v.Mesh.Vao)
		gl.UseProgram(v.Shader.Program)
		gc.bindShaderUniforms(v)
		gl.DrawElements(gl.TRIANGLES, int32(len(v.Mesh.F)), gl.UNSIGNED_SHORT, gl.Pointer(nil))

		// cleanup.
		gl.ActiveTexture(gl.TEXTURE0)
		gl.UseProgram(0)
		gl.BindVertexArray(0)
	}
}

// Implements Renderer interface.
func (gc *opengl) BindModel(mesh *data.Mesh) (err error) {
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		log.Printf("opengl:bindModel need to find and fix prior error %X", glerr)
	}

	// Gather the one scene into this one vertex array object.
	gl.GenVertexArrays(1, &(mesh.Vao))
	gl.BindVertexArray(mesh.Vao)

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

	// texture coordatinate, 2 float32's
	if len(mesh.T) > 0 {
		var tbuff uint32
		gl.GenBuffers(1, &tbuff)
		gl.BindBuffer(gl.ARRAY_BUFFER, tbuff)
		gl.BufferData(gl.ARRAY_BUFFER, int64(len(mesh.T)*4), gl.Pointer(&(mesh.T[0])), gl.STATIC_DRAW)
		var tattr uint32 = 2
		gl.VertexAttribPointer(tattr, 2, gl.FLOAT, false, 0, 0)
		gl.EnableVertexAttribArray(tattr)
	}

	// faces data, uint16 in this case, so 2 bytes per element.
	var fbuff uint32
	gl.GenBuffers(1, &fbuff)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, fbuff)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int64(len(mesh.F)*2), gl.Pointer(&(mesh.F[0])), gl.STATIC_DRAW)
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		err = fmt.Errorf("Failed binding model %s\n", mesh.Name)
	}
	return
}

// Implements Renderer interface.
func (gc *opengl) BindGlyphs(mesh *data.Mesh) (err error) {
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		log.Printf("opengl:bindGlyphs need to find and fix prior error %X", glerr)
	}

	// Gather the one scene into this one vertex array object.
	if mesh.Vao == 0 {
		gl.GenVertexArrays(1, &(mesh.Vao))
	}
	gl.BindVertexArray(mesh.Vao)

	// vertex data.
	if mesh.Vbuf == 0 {
		gl.GenBuffers(1, &(mesh.Vbuf))
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.Vbuf)
	gl.BufferData(gl.ARRAY_BUFFER, int64(len(mesh.V)*4), gl.Pointer(&(mesh.V[0])), gl.DYNAMIC_DRAW)
	var vattr uint32 = 0
	gl.VertexAttribPointer(vattr, 4, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(vattr)

	// texture coordatinate, 2 float32's
	if mesh.Tbuf == 0 {
		gl.GenBuffers(1, &(mesh.Tbuf))
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.Tbuf)
	gl.BufferData(gl.ARRAY_BUFFER, int64(len(mesh.T)*4), gl.Pointer(&(mesh.T[0])), gl.DYNAMIC_DRAW)
	var tattr uint32 = 2
	gl.VertexAttribPointer(tattr, 2, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(tattr)

	// faces data, uint16 in this case, so 2 bytes per element.
	if mesh.Fbuf == 0 {
		gl.GenBuffers(1, &mesh.Fbuf)
	}
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mesh.Fbuf)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int64(len(mesh.F)*2), gl.Pointer(&(mesh.F[0])), gl.DYNAMIC_DRAW)
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		err = fmt.Errorf("Failed binding glyphs %s\n", mesh.Name)
	}
	return
}

// Implements Renderer interface.
func (gc *opengl) BindTexture(texture *data.Texture) (err error) {
	if texture == nil {
		return
	}
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		log.Printf("opengl:bindTexture need to find and fix prior error %X", glerr)
	}
	gl.GenTextures(1, &(texture.Tid))
	gl.BindTexture(gl.TEXTURE_2D, texture.Tid)

	// ensure image is in RGBA format
	b := texture.Img.Bounds()
	rgba := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(rgba, rgba.Bounds(), texture.Img, b.Min, draw.Src)
	width, height := int32(b.Dx()), int32(b.Dy())
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Pointer(&(rgba.Pix[0])))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		err = fmt.Errorf("Failed binding texture %s\n", texture.Name)
	}
	return
}

// Implements Renderer interface.
func (gc *opengl) MapTexture(tid int, t *data.Texture) {
	tmap := map[int]uint32{
		0: gl.TEXTURE0,
		1: gl.TEXTURE1,
		2: gl.TEXTURE2,
		3: gl.TEXTURE3,
		4: gl.TEXTURE4,
		5: gl.TEXTURE5,
		6: gl.TEXTURE6,
		7: gl.TEXTURE7,
		8: gl.TEXTURE8,
		9: gl.TEXTURE9,
	}
	gl.ActiveTexture(tmap[tid])
	gl.BindTexture(gl.TEXTURE_2D, t.Tid)
}

// Implements Renderer interface.
func (gc *opengl) BindShader(sh *data.Shader) (pref uint32, err error) {
	pref = gl.CreateProgram()

	// TODO get rid of BindAttribLocation and use layout instead.
	//      this needs GLSL 330 instead of 150
	//      eg: layout(location=0) in vec4 in_position;
	gl.BindAttribLocation(pref, 0, "in_v") // matches vattr in bindModel
	gl.BindAttribLocation(pref, 1, "in_n") // matches nattr in bindModel
	gl.BindAttribLocation(pref, 2, "in_t") // matches tattr in bindModel

	// compile and link the shader program.
	if glerr := gl.BindProgram(pref, sh.Vsh, sh.Fsh); glerr != nil {
		err = fmt.Errorf("Failed to create shader program: %s\n", glerr)
		return
	}

	// initialize the uniform references
	var errmsg string
	for label, _ := range sh.Uniforms {
		if uid := gl.GetUniformLocation(pref, label); uid >= 0 {
			sh.Uniforms[label] = uid
		} else {
			errnum := gl.GetError()
			errmsg += fmt.Sprintf("No %s uniform in shader %X\n", label, errnum)
		}
	}
	if len(errmsg) > 0 {
		err = errors.New(errmsg)
	}
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		log.Printf("opengl:bindShader need to find and fix error %X", glerr)
	}
	return
}

// sTime acts as a reference point for the shaders that need the amount of
// elapsed time in seconds.
var sTime = time.Now()

// bindShaderUniforms sets the uniforms in the shader to the necessary values.
// This is done by an agreed uniform naming convention such that the uniforms
// names used in the shader are unique and apply to a particular set of data.
func (gc *opengl) bindShaderUniforms(v *Visible) {
	for key, ref := range v.Shader.Uniforms {
		switch key {
		case "mvpm":
			gc.bindUniforms(ref, m4, v.Mvp.Pointer())
		case "mvm":
			gc.bindUniforms(ref, m4, v.Mv.Pointer())
		case "nm":
			gc.bindUniforms(ref, m3, v.Mv.M3().Pointer())
		case "l":
			gc.bindUniforms(ref, f4, v.L.X, v.L.Y, v.L.Z, float32(1.0))
		case "ld":
			gc.bindUniforms(ref, f3, v.L.Ld.R, v.L.Ld.G, v.L.Ld.B)
		case "ka":
			gc.bindUniforms(ref, f3, v.Mat.Ka.R, v.Mat.Ka.G, v.Mat.Ka.B)
		case "kd":
			gc.bindUniforms(ref, f3, v.Mat.Kd.R, v.Mat.Kd.G, v.Mat.Kd.B)
		case "ks":
			gc.bindUniforms(ref, f3, v.Mat.Ks.R, v.Mat.Ks.G, v.Mat.Ks.B)
		case "scale":
			gc.bindUniforms(ref, f1, v.Scale)
		case "fd":
			gc.bindUniforms(ref, f1, v.Fade)
		case "alpha":
			if v.Mat != nil {
				gc.bindUniforms(ref, f1, v.Mat.Tr)
			} else {
				gc.bindUniforms(ref, f1, float32(1))
			}
		case "uv":
			gc.bindUniforms(ref, i1, int32(0))
			gc.MapTexture(0, v.Texture)
		case "time":
			gc.bindUniforms(v.Shader.Uniforms["time"], f1, float32(time.Since(sTime).Seconds()))
		case "resolution":
			gc.bindUniforms(v.Shader.Uniforms["resolution"], f2, float32(500), float32(500))
		case "rs":
			gc.bindUniforms(ref, f1, v.RotSpeed)
		}
	}
}

// bindUniforms wraps all the various glUniform calls as a single method.
// It expects the variable parameter list to match the given type.
func (gc *opengl) bindUniforms(uniform int32, utype int, udata ...interface{}) {
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
	case vf1:
		count := udata[0].(int32)
		vptr := udata[1].(*float32)
		gl.Uniform1fv(uniform, count, vptr)
	case vf2:
		count := udata[0].(int32)
		vptr := udata[1].(*float32)
		gl.Uniform2fv(uniform, count, vptr)
	case vi1:
		count := udata[0].(int32)
		vptr := udata[1].(*int32)
		gl.Uniform1iv(uniform, count, vptr)
	case m3:
		mptr := udata[0].(*float32)
		gl.UniformMatrix3fv(uniform, 1, false, mptr)
	case m4:
		mptr := udata[0].(*float32)
		gl.UniformMatrix4fv(uniform, 1, false, mptr)
	}
}

// Current list of supported uniform types.
const (
	i1  = iota // glUniform1i
	f1         // glUniform1f
	f2         // glUniform2f
	f3         // glUniform3f
	f4         // glUniform4f
	m3         // glUniformMatrix3fv
	m4         // glUniformMatrix4fv
	vf1        // glUniform1fv
	vf2        // glUniform2fv
	vi1        // glUniform1iv
)

// validate that OpenGL is available at the right version.  For OpenGL 3.2
// the following lines should be in the report.
//	    GL_VERSION_3_2
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
