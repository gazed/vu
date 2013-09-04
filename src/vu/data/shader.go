// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package data

import (
	"strings"
)

// Shader controls rendering. A shader is a set of programs that run on the
// graphics card. Different shaders are created for different effects. Shaders commonly
// need input from the CPU side. These input variables are called uniforms.
type Shader struct {
	Name     string           // Unique name (file name) of the shader.
	Vsh      []string         // Vertex shader source.
	Fsh      []string         // Fragment shader source.
	Program  uint32           // Reference to the compiled vertex and fragment program.
	Uniforms map[string]int32 // Uniform data is required.
}

// EnsureNewLines makes sure that shader program lines of code are properly terminated
// for the shader compiler.
func (s *Shader) EnsureNewLines() {
	for cnt, line := range s.Vsh {
		s.Vsh[cnt] = strings.TrimSpace(line) + "\n"
	}
	for cnt, line := range s.Fsh {
		s.Fsh[cnt] = strings.TrimSpace(line) + "\n"
	}
}
