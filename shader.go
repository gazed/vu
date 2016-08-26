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

// stripID is a helper method used by SetSource to parse GLSL
// shader code.
func (s *shader) stripID(id string) string {
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
	"uvc":     uvcShader,
	"bump":    bumpShader,
	"nmap":    nmapShader,
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
		"uniform vec3  lp;",    // light position in world space.
		"uniform vec3  lc;",    // light color.
		"uniform vec3  kd;",    // material diffuse color.
		"uniform float alpha;", // transparency
		"out     vec4  v_c;",   // vertex color
		"void main() {",
		"   vec4 vpos = vec4(in_v, 1.0);",                    // vertex in model space.
		"   vec3 nm = normalize((mvm * vec4(in_n, 0)).xyz);", // unit normal in world space.
		"   vec3 lightDir = normalize(lp - vec3(mvm*vpos));",
		"   vec3 color = lc * kd * max(dot(lightDir, nm), 0.0);",
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
		"uniform mat4  mvpm;",  // model view projection matrix
		"uniform mat4  mvm;",   // model view matrix
		"uniform vec3  lp;",    // light position in world space
		"uniform vec3  lc;",    // light color
		"uniform vec3  ka;",    // material ambient value
		"uniform vec3  kd;",    // material diffuse value
		"uniform vec3  ks;",    // material specular value
		"uniform float ns;",    // material specular exponent
		"uniform float alpha;", // transparency
		"out     vec4  v_c;",   // vertex color
		"void main() {",
		"   vec4 vmod = vec4(in_v, 1.0);",                    // vertex in model space.
		"   vec3 nm = normalize((mvm * vec4(in_n, 0)).xyz);", // unit normal in world space.
		"   vec4 vworld = mvm * vmod;",                       // vertex in world space
		"   vec3 s = normalize(lp - vworld.xyz);",            // light vector
		"   vec3 v = normalize(-vworld.xyz);",                // view vector
		"   vec3 r = reflect(-s, nm);",                       // light vec reflected around normal.
		"   vec3 ambient = lc * ka;",
		"   float sDotN = max( dot(s, nm), 0.0 );",
		"   vec3 diffuse = lc * kd * sDotN;",
		"   vec3 spec = vec3(0.0);",
		"   if (sDotN > 0.0)",
		"      spec = lc * ks * pow( max( dot(r,v), 0.0 ), ns );",
		"   vec3 color = ambient + diffuse + spec;", // combine all the values.
		"   v_c = vec4(color, alpha);",              // pass on the vertex color
		"   gl_Position = mvpm * vmod;",             // pass on the transformed vertex
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
		"uniform vec3  lp;",   // light position in world space.
		"out   vec3  v_n;",    // vertex normal
		"out   vec3  v_s;",    // vector from vertex to light.
		"out   vec3  v_e;",    // vertex eye position.
		"void main() {",
		"   vec4 vmod = vec4(in_v, 1.0);",                // vertex in model space.
		"   vec4 vworld = mvm * vmod;",                   // vertex in world space
		"   v_s = normalize(lp - vworld.xyz);",           // light vector
		"   v_e = normalize(-vworld.xyz);",               // view vector
		"   v_n = normalize((mvm * vec4(in_n, 0)).xyz);", // unit normal in world space.
		"   gl_Position = mvpm * vmod;",                  // vertex in clip space",
		"}",
	}
	fsh = []string{
		"#version 330",
		"in      vec3  v_n;",   // interpolated normal
		"in      vec3  v_s;",   // interpolated vector from vertex to light.
		"in      vec3  v_e;",   // interpolated vector from eye to vertex.
		"uniform vec3  lc;",    // light color
		"uniform vec3  ka;",    // material ambient value
		"uniform vec3  ks;",    // material specular value
		"uniform vec3  kd;",    // material diffuse value
		"uniform float alpha;", // transparency
		"uniform float ns;",    // material specular exponent
		"out     vec4  ffc;",   // final fragment color
		"void main() {",
		"   vec3 r = reflect(-v_s, v_n);",
		"   float sDotN = max( dot(v_s,v_n), 0.0 );",
		"   vec3 ambient = lc * ka;",
		"   vec3 diffuse = lc * kd * sDotN;",
		"   vec3 spec = vec3(0.0);",
		"   if (sDotN > 0.0)",
		"      spec = lc * ks * pow( max( dot(r,v_e), 0.0 ), ns);",
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
		"   ffc.a *= alpha;",
		"}",
	}
	return vsh, fsh
}

// ===========================================================================

// uvcShader handles a single texture and incorporates a single color.
func uvcShader() (vsh, fsh []string) {
	vsh, _ = uvShader()
	fsh = []string{
		"#version 330",
		"in      vec2      t_uv;",
		"uniform sampler2D uv;",
		"uniform vec3      kd;",    // material diffuse value
		"uniform float     alpha;", // transparency
		"out     vec4      ffc;",   // final fragment color
		"void main() {",
		"   ffc = texture(uv, t_uv);",
		"   ffc.xyz += kd;",
		"   ffc.a *= alpha;",
		"}",
	}
	return vsh, fsh
}

// ===========================================================================

// non-tangent based bump map code based on Vu uv shader and
// concepts from http://www.swiftless.com/tutorials/glsl/8_bump_mapping.html
func bumpShader() (vsh, fsh []string) {
	vsh = []string{
		"#version 330",
		"layout(location=0) in vec3 in_v;", // verticies
		"layout(location=2) in vec2 in_t;", // texture coordinates
		"uniform               mat4 mvpm;", // projection * model_view
		"out                   vec2 t_uv;", // pass uv coordinates through
		"void main() {",
		"   gl_Position = mvpm * vec4(in_v, 1.0);", // transformed vertex
		"   t_uv = in_t;",                          // uv coordinates
		"}",
	}
	fsh = []string{
		"#version 330",
		"in      vec2      t_uv;", // uv coordinates
		"uniform sampler2D uv;",   // model texture
		"uniform sampler2D uv1;",  // texture with normals
		"out     vec4      ffc;",  // final fragment color
		"void main() {",
		"   vec3 normal = normalize(texture(uv1, t_uv).rgb * 2.0 - 1.0);", // Normal from texture.
		"   vec3 light_pos = normalize(vec3(1.0, 1.0, 1.5));",             // Fixed light position.
		"   float diffuse = max(dot(normal, light_pos), 0.0);",            // Calculate the lighting diffuse value
		"   vec3 color = diffuse * texture(uv, t_uv).rgb;",                // Color from model texture.
		"   ffc = vec4(color, 1.0);",                                      // Color of current pixel
		"}",
	}
	return vsh, fsh
}

// ===========================================================================

// nmap is a tangent based normal map shader based on code and concepts from
//   http://www.thetenthplanet.de/archives/1180
// Generating the cotangent transform for each pixel Avoids the need to
// send tangent information for each vertex.
func nmapShader() (vsh, fsh []string) {
	vsh = []string{
		"#version 330 core",
		"layout(location = 0) in vec3 in_v;", // vertex position in modelspace
		"layout(location = 1) in vec3 in_n;", // vertex normal in modespace
		"layout(location = 2) in vec2 in_t;", // vertex uv texture coordinates
		"",
		"uniform mat4 mvpm;", // modelViewProject matrix for clipspace transform.
		"uniform mat4 mvm;",  // modelView matrix (and 3x3 normal matrix). Local -> view space.
		"uniform vec3 lp;",   // directional light.
		"out vec2 t_uv;",     // vertex texture coords
		"out vec3 v_n;",      // vertex normal
		"out vec3 v_l;",      // light to vertex vector
		"out vec3 v_v;",      // view vector: ie: camera to vertex
		"",
		"void main() {",
		"	vec4 vmod = vec4(in_v, 1.0);", // vertex in local model space.
		"	vec4 vcam = mvm * vmod;", // vertex in camera view space
		"	v_l = normalize(lp - vcam.xyz);", // normalized vertex to light vector
		"	v_v = -vcam.xyz;", // non-normalized vertex view vector in view space
		"   v_n = (mvm * vec4(in_n, 0)).xyz;", // non-normalized vertex normal in view space.
		"   t_uv = in_t;",                     // vertex UV texture coordinates.
		"   gl_Position = mvpm * vmod;",       // vertex position in clip space.
		"}",
	}

	// Fragment shader based on http://www.thetenthplanet.de/archives/1180
	// cotangent_frame creates the cotangent frame necessary to
	// map from the cotangent space to world space. Note that the
	// GLSL methods used are only available in the fragment shader.
	// Parameters:",
	//   N  : the interpolated vertex normal in world space
	//   p  : the reverse view vector in world space
	//   tuv: texture coordinates
	//
	// perturb_normal transforms a normal from cotangent space
	// Note original code handles different normal map texture encodings.
	//   map = map * 255./127. - 128./127.;      #ifdef WITH_NORMALMAP_UNSIGNED
	//   map.z = sqrt(1. - dot(map.xy, map.xy)); #ifdef WITH_NORMALMAP_2CHANNEL
	//   map.y = -map.y;                         #ifdef WITH_NORMALMAP_GREEN_UP
	// Parameters:
	//   N  : the interpolated vertex normal in world space
	//   V  : the view vector in world space
	//   tuv: texture coordinates
	fsh = []string{
		"#version 330 core",
		"in vec2 t_uv;",          // vertex texture coords
		"in vec3 v_n;",           // non-normalized vertex normal
		"in vec3 v_v;",           // non-normalized view vector
		"in vec3 v_l;",           // normalized light to vertex vector
		"uniform sampler2D uv;",  // base diffuse color
		"uniform sampler2D uv1;", // normal map texture in tangent space
		"uniform sampler2D uv2;", // specular texture
		"uniform vec3      lc;",  // light color
		"uniform vec3      ka;",  // material ambient color
		"uniform vec3      kd;",  // material diffuse color
		"uniform vec3      ks;",  // material specular color
		"uniform float     ns;",  // material specular shininess.
		"out vec4 ffc;",
		"",
		"mat3 cotangent_frame( vec3 N, vec3 p, vec2 tuv ) {",
		"", // get edge vectors of the pixel triangle
		"    vec3 dp1 = dFdx( p );",
		"    vec3 dp2 = dFdy( p );",
		"    vec2 duv1 = dFdx( tuv );",
		"    vec2 duv2 = dFdy( tuv );",
		"", // solve the linear system
		"    vec3 dp2perp = cross( dp2, N );",
		"    vec3 dp1perp = cross( N, dp1 );",
		"    vec3 T = dp2perp * duv1.x + dp1perp * duv2.x;",
		"    vec3 B = dp2perp * duv1.y + dp1perp * duv2.y;",
		"", // construct a scale-invariant frame
		"    float invmax = inversesqrt( max( dot(T,T), dot(B,B) ) );",
		"    return mat3( T * invmax, B * invmax, N );",
		"}",
		"vec3 perturb_normal( vec3 N, vec3 V, vec2 tuv ) {",
		"    vec3 map = texture(uv1, tuv).xyz;",       // perturbed normal in cotangent space
		"    map = map * 255./127. - 128./127.;",      // #ifdef WITH_NORMALMAP_UNSIGNED
		"    mat3 TBN = cotangent_frame(N, -V, tuv);", // cotangent to world space transform
		"    return normalize(TBN * map);",            // perturbed normal in world space
		"}",
		"void main() {",
		"    vec3 normal = normalize(v_n);",
		"    normal = perturb_normal(normal, v_v, t_uv);", // normal in view space.
		"    vec3 nv_v = normalize(v_v);",
		"    vec3 ambient = lc * ka;",
		"    float intensity = max(dot(v_l, normal), 0.0);",
		"    vec3 diffuse = lc * kd * intensity;",
		"", // Blinn-Phong half vector.
		"    vec3 halfDir = normalize(v_l + nv_v);",
		"    float specAngle = max(dot(halfDir, normal), 0.0);",
		"    float specFac = pow(clamp(specAngle, 0.0, 1.0), ns);",
		"    vec3 smap = texture(uv2, t_uv).rgb;",
		"    vec3 specular = lc * ks * smap * specFac;",
		"", // Combine into final fragment color.
		"    vec3 base = texture(uv, t_uv).rgb;",                   // pure texture color.
		"    vec3 color = ambient*base + diffuse*base + specular;", // combine all the values.
		"    ffc = vec4(color, 1.0);",                              // final fragment color
		" }",
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
		"   mat3x4 m = bpos[int(joint.x)] * weight.x;", // up to four joints affect vertex.
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
