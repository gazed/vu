// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package gl

// The majority of the "bind" functionality interacts with "vu/data" objects
// and is not suitable for this package.  However other utility methods may
// be appropriate.
//
// BindProgram is a utility method and is here because it has no dependencies
// on other packages.

import (
	"errors"
	"fmt"
)

// BindProgram combines a vertex and fragment shader into a program.
// The bulk of this method is error checking and returning error information
// when there are compile or link problems.
func BindProgram(program uint32, vertexSource, fragmentSource []string) error {

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

	// Don't validate since validate checks as if the OpenGL state was ready to use
	// the program and the lack of a current VAO would cause validate to fail.
	// The lack of a check allows bindProgram to work even if BindVertexArray(0)
	// has not been called.
	return nil
}
