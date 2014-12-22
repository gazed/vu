// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package render

import (
	"strings"
)

// Shader is the GPU based program that renders data contained in a Model.
// A shader is a combination of two or more programs that run on
// the graphics card during different stages of the rendering pipeline.
// Different shaders are created for different effects. Shaders commonly need
// buffers of data and uniforms from the CPU side. The data is expected to
// be contained in, and provided by a Model.
type Shader interface {
	Name() string                // Unique identifier set on creation.
	Vsh() []string               // Vertex shader source.
	Fsh() []string               // Fragment shader source.
	Lib() (vsh, fsh []string)    // Find this shaders source in the library.
	SetSource(vsh, fsh []string) // Directly set shader source.
	Bound() bool                 // True if the shader has a GPU reference.
}

// ============================================================================

// shader is the opengl Shader implementation. It encapsulates all the OpenGL
// and GLSL specific knowledge while conforming to the generic Shader interface.
type shader struct {
	name    string   // Unique shader identifier.
	vsh     []string // Vertex shader source, empty if data not loaded.
	fsh     []string // Fragment shader source, empty if data not loaded.
	program uint32   // Compiled program reference. Zero if not compiled.
	refs    uint32   // Number of Model references.

	// Vertex data and uniform expectations are discovered from the shader
	// source. This can be later verified against available data.
	attributes map[string]uint32 // Expected buffer data locations.
	uniforms   map[string]int32  // Expected uniform data.
}

// newShader creates a new shader. It needs to be loaded with source and bound
// to the GPU.
func newShader(name string) Shader {
	sh := &shader{name: name}
	sh.attributes = map[string]uint32{}
	sh.uniforms = map[string]int32{}
	return sh
}

// Implement Shader
func (s *shader) Name() string  { return s.name }
func (s *shader) Vsh() []string { return s.vsh }
func (s *shader) Fsh() []string { return s.fsh }
func (s *shader) Bound() bool   { return s.program != 0 }

// Lib looks in the glsl shader library for a shader with the same
// name as this shader. The shader source is returned if one is found.
func (s *shader) Lib() (vsh, fsh []string) {
	if sfn, ok := glsl[s.name]; ok {
		return sfn()
	}
	return []string{}, []string{}
}

// Shader source is scanned for uniforms and vertex buffer information.
// The uniform references are set on binding and later used by model.go
// to set the uniform values during rendering.
func (s *shader) SetSource(vsh, fsh []string) {
	s.vsh = vsh
	for _, line := range s.vsh {
		if fields := strings.Fields(line); len(fields) > 2 {
			switch fields[0] {
			case "layout(location=0)":
				bid := s.stripId(fields[3])
				s.attributes[bid] = 0
			case "layout(location=1)":
				bid := s.stripId(fields[3])
				s.attributes[bid] = 1
			case "layout(location=2)":
				bid := s.stripId(fields[3])
				s.attributes[bid] = 2
			case "layout(location=3)":
				bid := s.stripId(fields[3])
				s.attributes[bid] = 3
			case "layout(location=4)":
				bid := s.stripId(fields[3])
				s.attributes[bid] = 4
			case "layout(location=5)":
				bid := s.stripId(fields[3])
				s.attributes[bid] = 5
			case "uniform":
				uid := s.stripId(fields[2])
				s.uniforms[uid] = -1
			}
		}
	}
	s.fsh = fsh
	for _, line := range s.fsh {
		if fields := strings.Fields(line); len(fields) > 2 {
			fields := strings.Fields(line)
			switch fields[0] {
			case "uniform":
				uid := s.stripId(fields[2])
				s.uniforms[uid] = -1
			}
		}
	}
}

// stripId is a helper method used by SetSource.
func (s *shader) stripId(id string) string {
	id = strings.Replace(id, ";", "", -1)
	if strings.Contains(id, "[") {
		strs := strings.Split(id, "[")
		return strs[0]
	}
	return id
}

// ensureNewLines properly terminates shader program lines for the shader compiler.
func (s *shader) ensureNewLines() {
	for cnt, line := range s.vsh {
		s.vsh[cnt] = strings.TrimSpace(line) + "\n"
	}
	for cnt, line := range s.fsh {
		s.fsh[cnt] = strings.TrimSpace(line) + "\n"
	}
}
