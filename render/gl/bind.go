// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package gl

// bind contains OpenGL functions that help bind shaders.
// There are no dependencies on anything but OpenGL itself.

import (
	"errors"
	"fmt"
	"strings"
)

// BindProgram combines a vertex and fragment shader into a program.
// The bulk of this method is error checking and returning error information
// when there are compile or link problems.
func BindProgram(program uint32, vertexSource, fragmentSource []string) error {
	glslVersion := GetString(SHADING_LANGUAGE_VERSION)
	vertexSource = addPrelude(vertexSource, glslVersion)
	fragmentSource = addPrelude(fragmentSource, glslVersion)

	// Compile and attach the vertex shader
	var status int32
	vertexShader := CreateShader(VERTEX_SHADER)
	defer DeleteShader(vertexShader)
	ShaderSource(vertexShader, int32(len(vertexSource)), vertexSource, nil)
	CompileShader(vertexShader)
	GetShaderiv(vertexShader, COMPILE_STATUS, &status)
	if status == FALSE {
		errmsg := fmt.Sprintf("Vertex shader compile failed\n")
		var logLength int32
		GetShaderiv(vertexShader, INFO_LOG_LENGTH, &logLength)
		if logLength > 0 {
			log := make([]byte, logLength, logLength)
			GetShaderInfoLog(vertexShader, logLength, &logLength, &(log[0]))
			errmsg += fmt.Sprintf("Vertex shader compile log:\n%s\n", string(log))
		}
		return errors.New(errmsg)
	}
	AttachShader(program, vertexShader)

	// Compile and attach the fragment shader
	fragShader := CreateShader(FRAGMENT_SHADER)
	defer DeleteShader(fragShader)
	ShaderSource(fragShader, int32(len(fragmentSource)), fragmentSource, nil)
	CompileShader(fragShader)
	GetShaderiv(fragShader, COMPILE_STATUS, &status)
	if status == FALSE {
		errmsg := fmt.Sprintf("Fragment shader compile failed\n")
		var logLength int32
		GetShaderiv(fragShader, INFO_LOG_LENGTH, &logLength)
		if logLength > 0 {
			log := make([]byte, logLength, logLength)
			GetShaderInfoLog(fragShader, logLength, &logLength, &(log[0]))
			errmsg += fmt.Sprintf("Fragment shader compile log:\n%s\n", string(log))
		}
		return errors.New(errmsg)
	}
	AttachShader(program, fragShader)

	// Link the program
	LinkProgram(program)
	GetProgramiv(program, LINK_STATUS, &status)
	if status == FALSE {
		errmsg := fmt.Sprintf("Shader link failed\n")
		var logLength int32
		GetProgramiv(program, INFO_LOG_LENGTH, &logLength)
		if logLength > 0 {
			log := make([]byte, logLength, logLength)
			GetProgramInfoLog(program, logLength, &logLength, &(log[0]))
			errmsg += fmt.Sprintf("Shader link log:\n%s\n", string(log))
		}
		return errors.New(errmsg)
	}

	// Don't validate since validate checks as if the OpenGL state is ready for
	// the program and the lack of a current VAO would cause validate to fail.
	// The lack of a check allows bindProgram to work even if BindVertexArray(0)
	// has not been called.
	return nil
}

// addPrelude adds glsl version specific information to shader source.
// This is needed for them to compile both on desktop OpenGL and
// mobile OpenGLES.
func addPrelude(source []string, glslVersion string) []string {
	prelude := []string{}
	switch glslVersion {
	case "OpenGL ES GLSL ES 3.00": // iOS value
		prelude = append(prelude, "#version 300 es\n")
		prelude = append(prelude, "precision highp float;\n")
	default:
		prelude = append(prelude, "#version 330\n")
	}
	src := make([]string, len(source)+len(prelude))
	copy(src[0:], prelude)
	copy(src[len(prelude):], source)
	return src
}

// Uniforms fills in the active uniform names and locations for a compiled
// and linked program. Expected to be called with valid programs, thus
// errors due to invalid programs are ignored.
func Uniforms(program uint32, uniforms map[string]int32) {
	var nUniforms, maxLen int32
	GetProgramiv(program, ACTIVE_UNIFORM_MAX_LENGTH, &maxLen)
	GetProgramiv(program, ACTIVE_UNIFORMS, &nUniforms)
	var size, location, written int32
	var kind uint32
	for i := uint32(0); i < uint32(nUniforms); i++ {
		name := make([]byte, maxLen, maxLen)
		GetActiveUniform(program, i, maxLen, &written, &size, &kind, &(name[0]))
		if written > 0 {
			name = name[:written] // truncate to bytes written.
			location = GetUniformLocation(program, string(name))
			uniform := string(name)
			if parts := strings.Split(uniform, "["); len(parts) > 1 {
				uniform = parts[0] // string array uniforms to just the label.
			}
			uniforms[uniform] = location
		}
	}
}

// Layouts fills in the active attribute names and locations for a compiled
// and linked program. These are the per-vertex data identfiers.
// Expected to be called with valid programs, thus errors due to
// invalid programs are ignored.
func Layouts(program uint32, layouts map[string]uint32) {
	var nAttrs, maxLen int32
	GetProgramiv(program, ACTIVE_ATTRIBUTE_MAX_LENGTH, &maxLen)
	GetProgramiv(program, ACTIVE_ATTRIBUTES, &nAttrs)
	var size, location, written int32
	var kind uint32
	for i := uint32(0); i < uint32(nAttrs); i++ {
		name := make([]byte, maxLen, maxLen)
		GetActiveAttrib(program, i, maxLen, &written, &size, &kind, &(name[0]))
		if written > 0 {
			name = name[:written]
			location = GetAttribLocation(program, string(name))
			layouts[string(name)] = uint32(location)
		}
	}
}
