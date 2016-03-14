// Copyright © 2015-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// FUTURE: enhance design to incorporate/handle HLSL and Vulkan shaders.

import (
	"strings"
)

// shader is an essential part of a rendered Model. It contains the logic
// to control parts of the GPU render pipeline. Shader is OpenGL specific.
// It encapsulates all the OpenGL and GLSL specific knowledge while conforming
// to the generic Shader interface.
type shader struct {
	name    string   // Unique shader identifier.
	tag     uint64   // name and type as a number.
	vsh     []string // Vertex shader source, empty if data not loaded.
	fsh     []string // Fragment shader source, empty if data not loaded.
	program uint32   // Compiled program reference. Zero if not compiled.
	bound   bool     // False if the data needs rebinding.
	loaded  bool     // True if data has been set.

	// Vertex layout data and uniform expectations are discovered from the
	// shader source. This can be verified later against available data.
	layouts  map[string]uint32 // Expected buffer data locations.
	uniforms map[string]int32  // Expected uniform data.
}

// newShader creates a new shader.
// It needs to be loaded with shader source code and bound to the GPU.
func newShader(name string) *shader {
	sh := &shader{name: name, tag: shd + stringHash(name)<<32}
	sh.layouts = map[string]uint32{}
	sh.uniforms = map[string]int32{}
	return sh
}

// label, aid, and bid are used to uniquely identify assets.
func (s *shader) label() string { return s.name }                      // asset name
func (s *shader) aid() uint64   { return s.tag }                       // asset type and name.
func (s *shader) bid() uint64   { return shd + uint64(s.program)<<32 } // asset type and bind ref.

// Shader source is scanned for uniforms and vertex buffer information.
// The uniform references are set on binding and later used by Model
// to set the uniform values during rendering.
func (s *shader) setSource(vsh, fsh []string) {
	s.vsh, s.fsh = vsh, fsh
	s.ensureNewLines()
	s.loaded = len(s.vsh) > 0 && len(s.fsh) > 0
}

// stripId is a helper method used by SetSource to parse GLSL
// shader code.
func (s *shader) stripId(id string) string {
	id = strings.Replace(id, ";", "", -1)
	if strings.Contains(id, "[") {
		strs := strings.Split(id, "[")
		return strs[0]
	}
	return id
}

// ensureNewLines properly terminates shader program lines
// for the shader compiler.
func (s *shader) ensureNewLines() {
	for cnt, line := range s.vsh {
		s.vsh[cnt] = strings.TrimSpace(line) + "\n"
	}
	for cnt, line := range s.fsh {
		s.fsh[cnt] = strings.TrimSpace(line) + "\n"
	}
}

// shader
// =============================================================================
// shaderLibrary - glsl

// DESIGN: Keep the shaders relatively small until there are more/better
//         debugging tools and ways of handling shader code. Would like
//         something better (more concise, more robust) than this current
//         GLSL code within Go code design.

// shaderLibrary provides pre-made GLSL shaders. Each shader is identified
// by a unique name. These provide some basic shaders to get simple examples
// running quickly and can be used as starting templates for new shaders.
var shaderLibrary = map[string]func() (vsh, fsh []string){
	"solid":   solidShader,
	"alpha":   alphaShader,
	"diffuse": diffuseShader,
	"gouraud": gouraudShader,
	"phong":   phongShader,
	"uv":      uvShader,
	"bb":      bbShader,
	"bbr":     bbrShader,
	"anim":    animShader,
	"depth":   depthShader,
	"shadow":  shadowShader,
}

// FUTURE: Add edge-detect and emboss shaders, see:
//         http://www.processing.org/tutorials/pshader/

// ===========================================================================

// solidShader shades all verticies the given diffuse color.
func solidShader() (vsh, fsh []string) {
	vsh = []string{
		"#version 330",
		"layout(location=0) in vec3 in_v;", // verticies
		"",
		"uniform mat4 mvpm;", // model view projection matrix
		"uniform vec3 kd;",   // material diffuse value
		"out     vec4 v_c;",  // vertex color
		"void main(void) {",
		"   gl_Position = mvpm * vec4(in_v, 1.0);",
		"	v_c = vec4(kd, 1.0);",
		"}",
	}
	fsh = []string{
		"#version 330",
		"in  vec4 v_c;",  // color from vertex shader
		"out vec4 ffc; ", // final fragment color.
		"void main(void) {",
		"   ffc = v_c;",
		"}",
	}
	return vsh, fsh
}

// ===========================================================================

// alphaShader combines a color with alpha to make transparent objects.
func alphaShader() (vsh, fsh []string) {
	vsh = []string{
		"#version 330",
		"layout(location=0) in vec3 in_v;", // verticies
		"",
		"uniform mat4  mvpm;",  // model view projection matrix
		"uniform vec3  kd;",    // diffuse color
		"uniform float alpha;", // transparency
		"out     vec4  v_c;",   // vertex color
		"void main() {",
		"	gl_Position = mvpm * vec4(in_v, 1.0);",
		"	v_c = vec4(kd, alpha);",
		"}",
	}
	fsh = []string{
		"#version 330",
		"in  vec4 v_c;", // color from vertex shader
		"out vec4 ffc;", // final fragment color
		"void main() {",
		"   ffc = v_c;",
		"}",
	}
	return vsh, fsh
}

// diffuseShader is based on
//      http://www.packtpub.com/article/opengl-glsl-4-shaders-basics
//      http://devmaster.net/posts/2974/the-basics-of-3d-lighting
//
// Diffuse: The algorithm is used to calculate the color for the polygon.
//          The polygon has only one normal and the only part of the above
//          algorithm used is: diffuse*(normal . light-direction)
func diffuseShader() (vsh, fsh []string) {
	vsh = []string{
		"#version 330",
		"layout(location=0) in vec3 in_v;", // verticies
		"layout(location=1) in vec3 in_n;", // vertex normals
		"",
		"uniform mat4  mvpm;",  // model view projection matrix
		"uniform mat4  mvm;",   // model view matrix
		"uniform mat3  nm;",    // normal matrix
		"uniform vec4  l;",     // light position in camera space.
		"uniform vec3  ld;",    // light source intensity.
		"uniform vec3  kd;",    // material diffuse color.
		"uniform float alpha;", // transparency
		"out     vec4  v_c;",   // vertex color
		"void main() {",
		"   vec4 vpos = vec4(in_v, 1.0);",
		"   vec3 norm = normalize(nm * in_n);", // Convert normal and position to eye coords
		"   vec3 lightDirection = normalize(vec3(l - mvm*vpos));",
		"   vec3 color = ld * kd * max(dot(lightDirection, norm), 0.0);",
		"   v_c = vec4(color, alpha);",  // pass on the amount of diffuse light.
		"   gl_Position = mvpm * vpos;", // pass on the transformed vertex position
		"}",
	}
	fsh = []string{
		"#version 330",
		"in  vec4 v_c;", // interpolated vertex color
		"out vec4 ffc;", // final fragment color
		"void main() {",
		"   ffc = v_c;",
		"}",
	}
	return vsh, fsh
}

// ===========================================================================

// gouraudShader is based on
//      http://www.packtpub.com/article/opengl-glsl-4-shaders-basics
//      http://devmaster.net/posts/2974/the-basics-of-3d-lighting
//
// Gouraud : Calculate a color at each vertex.
//           The color are then interpolated across the polygon.
func gouraudShader() (vsh, fsh []string) {
	vsh = []string{
		"#version 330",
		"layout(location=0) in vec3 in_v;", // verticies
		"layout(location=1) in vec3 in_n;", // vertex normals
		"",
		"uniform mat4  mvpm;",           // model view projection matrix
		"uniform mat4  mvm;",            // model view matrix
		"uniform mat3  nm;",             // normal matrix
		"uniform vec4  l;",              // untransformed light position
		"uniform vec3  ld;",             // light source intensity
		"uniform vec3  ka;",             // material ambient value
		"uniform vec3  kd;",             // material diffuse value
		"uniform vec3  ks;",             // material specular value
		"uniform float alpha;",          // transparency
		"const   vec3  la = vec3(0.3);", // FUTURE make la a uniform.
		"const   vec3  ls = vec3(0.4);", // FUTURE make ls a uniform.
		"const   float shine = 8.0;",    // FUTURE make shine a uniform.
		"out     vec4  v_c;",            // vertex color
		"void main() {",
		"   vec4 vpos = vec4(in_v, 1.0);",
		"   vec3 norm = normalize(nm * in_n);",
		"   vec4 eyeCoords = mvm * vpos;",
		"   vec3 s = normalize(vec3(l - eyeCoords));",
		"   vec3 v = normalize(-eyeCoords.xyz);",
		"   vec3 r = reflect(-s, norm);",
		"   vec3 ambient = la * ka;",
		"   float sDotN = max( dot(s,norm), 0.0 );",
		"   vec3 diffuse = ld * kd * sDotN;",
		"   vec3 spec = vec3(0.0);",
		"   if (sDotN > 0.0)",
		"      spec = ls * ks * pow( max( dot(r,v), 0.0 ), shine );",
		"   vec3 color = ambient + diffuse + spec;", // combine all the values.
		"   v_c = vec4(color, alpha);",              // pass on the vertex color
		"   gl_Position = mvpm * vpos;",             // pass on the transformed vertex
		"}",
	}
	fsh = []string{
		"#version 330",
		"                   in      vec4      v_c;", // interpolated vertex color
		"layout(location=0) out     vec4      ffc;", // final fragment color
		"void main() {",
		"   ffc = v_c;",
		"}",
	}
	return vsh, fsh
}

// ===========================================================================

// phongShader is based on
//       http://www.packtpub.com/article/opengl-glsl-4-shaders-basics
//
// Phong Shading   : Calculate the color intensity for each pixel using
//                   interpolated vertex normals.
func phongShader() (vsh, fsh []string) {
	vsh = []string{
		"#version 330",
		"layout(location=0) in vec3 in_v;", // verticies
		"layout(location=1) in vec3 in_n;", // vertex normals
		"",
		"uniform mat4  mvpm;", // model view projection matrix
		"uniform mat4  mvm;",  // model view matrix
		"uniform mat3  nm;",   // normal matrix
		"uniform vec4  l;",    // untransformed light position
		"out   vec3  v_n;",    // vertex color
		"out   vec3  v_s;",    // vector from vertex to light.
		"out   vec3  v_e;",    // vertex eye position.
		"void main() {",
		"   vec4 vpos = vec4(in_v, 1.0);",
		"   vec4 eyeCoords = mvm * vpos;",
		"   v_n = normalize(nm * in_n);",
		"   v_s = normalize(vec3(l - eyeCoords));",
		"   v_e = normalize(-eyeCoords.xyz);",
		"   gl_Position = mvpm * vpos;", // pass on the transformed vertex position",
		"}",
	}
	fsh = []string{
		"#version 330",
		"in      vec3  v_n;",            // interpolated normal
		"in      vec3  v_s;",            // interpolated vector from vertex to light.
		"in      vec3  v_e;",            // interpolated vector from eye to vertex.
		"uniform vec3  ld;",             // light source intensity
		"uniform vec3  ka;",             // material ambient value
		"uniform vec3  ks;",             // material specular value
		"uniform vec3  kd;",             // material diffuse value
		"uniform float alpha;",          // transparency
		"const   vec3  la = vec3(0.3);", // FUTURE make la a uniform.
		"const   vec3  ls = vec3(0.4);", // FUTURE make ls a uniform.
		"const   float shine = 8.0;",    // FUTURE make shine a uniform.
		"out     vec4  ffc;",            // final fragment color
		"void main() {",
		"   vec3 r = reflect(-v_s, v_n);",
		"   float sDotN = max( dot(v_s,v_n), 0.0 );",
		"   vec3 ambient = la * ka;",
		"   vec3 diffuse = ld * kd * sDotN;",
		"   vec3 spec = vec3(0.0);",
		"   if (sDotN > 0.0)",
		"      spec = ls * ks * pow( max( dot(r,v_e), 0.0 ), shine);",
		"   vec3 color = ambient + diffuse + spec;", // combine all the values.
		"   ffc = vec4(color, alpha);",              // final fragment color
		"}",
	}
	return vsh, fsh
}

// ===========================================================================

// uvShader handles a single texture.
func uvShader() (vsh, fsh []string) {
	vsh = []string{
		"#version 330",
		"layout(location=0) in vec3 in_v;", // verticies
		"layout(location=2) in vec2 in_t;", // texture coordinates
		"",
		"uniform mat4  mvpm;", // projection * model_view
		"out     vec2  t_uv;", // pass uv coordinates through
		"void main() {",
		"   gl_Position = mvpm * vec4(in_v, 1.0);",
		"   t_uv = in_t;",
		"}",
	}
	fsh = []string{
		"#version 330",
		"in      vec2      t_uv;",
		"uniform sampler2D uv;",
		"uniform float     alpha;", // transparency
		"out     vec4      ffc;",   // final fragment color
		"void main() {",
		"   ffc = texture(uv, t_uv);",
		"   ffc.a = ffc.a*alpha;",
		"}",
	}
	return vsh, fsh
}

// ===========================================================================

// bbShader is a billboard shader. Like a uv shader it renders a single texture
// but forces the textured object to always face the camera. See
//     http://www.lighthouse3d.com/opengl/billboarding/billboardingtut.pdf
//     https://code.google.com/p/o3d/source/browse/trunk/samples_webgl/shaders/billboard.shader
//     http://www.sjbaker.org/steve/omniv/alpha_sorting.html
func bbShader() (vsh, fsh []string) {
	vsh = []string{
		"#version 330",
		"layout(location=0) in vec3 in_v;", // vertex coordinates
		"layout(location=2) in vec2 in_t;", // texture coordinates
		"",
		"uniform mat4  mvm;",   // model view matrix
		"uniform mat4  pm;",    // projection matrix
		"uniform vec3  scale;", // scale for each axis
		"out     vec2  t_uv;",  // pass uv coordinates through
		"",
		"void main() {",
		"   gl_Position = pm * (vec4(in_v*scale, 1) + vec4(mvm[3].xyz, 0));",
		"   t_uv = in_t;",
		"}",
	}
	fsh = []string{
		"#version 330",
		"in      vec2      t_uv;",  // interpolated uv coordinates
		"uniform sampler2D uv;",    // texture sampler
		"uniform float     alpha;", // transparency
		"out     vec4      ffc;",   // final fragment color
		"void main() {",
		"   ffc = texture(uv, t_uv) * vec4(1.0, 1.0, 1.0, alpha);",
		"}",
	}
	return vsh, fsh
}

// ===========================================================================

// bbrShader is a billboard shader that rotates the texture over time.
func bbrShader() (vsh, fsh []string) {
	vsh, _ = bbShader()
	fsh = []string{
		"#version 330",
		"in      vec2      t_uv;",  // interpolated uv coordinates
		"uniform sampler2D uv;",    // texture sampler
		"uniform float     time;",  // current time in seconds
		"uniform float     spin;",  // rotation speed 0 -> 1
		"uniform float     alpha;", // transparency
		"out     vec4      ffc;",   // final fragment color
		"",
		"void main() {",
		"   float sa = sin(time*spin);",               // calculate rotation
		"   float ca = cos(time*spin);",               // ..
		"   mat2 rot = mat2(ca, -sa, sa, ca);",        // ..
		"   ffc = texture(uv, ((t_uv-0.5)*rot)+0.5);", // rotate around its center
		"   ffc.a = ffc.a*alpha;",
		"}",
	}
	return vsh, fsh
}

// =============================================================================

// animShader is a bare bones skeletal shader that includes
// uv texture mapping and alpha.
func animShader() (vsh, fsh []string) {
	vsh = []string{
		"#version 330",
		"layout(location=0) in vec3 in_v;",   // verticies
		"layout(location=2) in vec2 in_t;",   // texture coordinates
		"layout(location=4) in vec4 joint;",  // joint indicies
		"layout(location=5) in vec4 weight;", // joint weights
		"uniform mat3x4     bpos[100];",      // bone positioning transforms. Row-Major!
		"uniform mat4       mvpm;",           // model view projection matrix
		"out     vec2       t_uv;",           // pass uv coordinates through
		"",
		"void main() {",
		"   mat3x4 m = bpos[int(joint.x)] * weight.x;", // upto four joints affect vertex.
		"   m += bpos[int(joint.y)] * weight.y;",
		"   m += bpos[int(joint.z)] * weight.z;",
		"   m += bpos[int(joint.w)] * weight.w;",
		"   vec4 mpos = vec4(vec4(in_v, 1.0) * m, 1.0);", // Row-Major pre-multiply.
		"   gl_Position = mvpm * mpos;",
		"   t_uv = in_t;",
		"}",
	}
	fsh = []string{
		"#version 330",
		"in      vec2      t_uv;",  // interpolated uv coordinates
		"uniform sampler2D uv;",    // texture sampler
		"uniform float     alpha;", // transparency
		"out     vec4      ffc; ",  // final fragment color
		"",
		"void main() {",
		"   ffc = texture(uv, t_uv);",
		"   ffc.a = ffc.a*alpha;",
		"}",
	}
	return vsh, fsh
}

// =============================================================================

// depthShader is used to create shadow maps by writing objects depths.
// Expected to be used during the shadow map render pass to render to
// a texture. See:
// http://www.opengl-tutorial.org/intermediate-tutorials/tutorial-16-shadow-mapping
func depthShader() (vsh, fsh []string) {
	vsh = []string{
		"#version 330",
		"layout (location = 0) in vec3 in_v;",
		"uniform mat4          mvpm;",
		"void main() {",
		"    gl_Position = mvpm * vec4(in_v, 1.0);",
		"}",
	}
	fsh = []string{
		"#version 330",
		"layout(location = 0) out float fragdepth;",
		"void main() {",
		"    fragdepth = gl_FragCoord.z;",
		"}",
	}
	return vsh, fsh
}

// =============================================================================

// shadowShader incorporates a shadow depth map into lighting calculations. See:
// http://www.opengl-tutorial.org/intermediate-tutorials/tutorial-16-shadow-mapping
func shadowShader() (vsh, fsh []string) {
	vsh = []string{
		"#version 330 core",
		"layout(location=0) in vec3 in_v;", // verticies
		"layout(location=2) in vec2 in_t;", // texture coordinates
		"uniform mat4       mvpm;",         // model view projection matrix
		"uniform mat4       dbm;",          // depth bias matrix
		"out     vec2       t_uv;",         // pass uv coordinates through
		"out     vec4       s_uv;",         // create shadow uv coordinates
		"void main(){",
		"    gl_Position = mvpm * vec4(in_v, 1.0);",
		"    s_uv = dbm * vec4(in_v, 1.0);",
		"    t_uv = in_t;",
		"}",
	}
	fsh = []string{
		"#version 330 core",
		"in      vec2            t_uv;",      // interpolated uv coordinates
		"in      vec4            s_uv;",      // interpolated shadow uv coordinates
		"uniform sampler2D       uv;",        // object material texture sampler
		"uniform sampler2DShadow sm;",        // shadow map depth texture sampler
		"layout(location = 0) out vec4 ffc;", // final fragment color
		"void main(){",
		"    vec4 lightColor = vec4(1,1,1,1);",        // white light
		"    vec4 diffuseColor = texture(uv, t_uv); ", // object color from texture
		"",
		"", // compare the depth found in the texture at xy
		"", // with the depth at z.
		"    vec2 suv = vec2((s_uv.xy)/s_uv.w);",
		"    float visibility = texture(sm, vec3(suv, (s_uv.z)/s_uv.w));",
		"    visibility = visibility + (1.0-visibility)*0.75;", // map 1-0 to 1-0.75
		"    ffc = visibility * diffuseColor * lightColor;",
		"}",
	}
	return vsh, fsh
}
