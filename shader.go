// Copyright © 2015-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// shader.go
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
	tag     aid      // name and type as a number.
	vsh     []string // Vertex shader source, empty if data not loaded.
	fsh     []string // Fragment shader source, empty if data not loaded.
	program uint32   // Compiled program reference. Zero if not compiled.

	// Vertex layout data and uniform expectations are discovered from the
	// shader source. This can be verified later against available data.
	layouts  map[string]uint32 // Expected buffer data locations.
	uniforms map[string]int32  // Expected uniform GPU references.
}

// newShader creates a new shader.
// It needs to be loaded with shader source code and bound to the GPU.
func newShader(name string) *shader {
	sh := &shader{name: name, tag: assetID(shd, name)}
	sh.layouts = map[string]uint32{}
	sh.uniforms = map[string]int32{}
	return sh
}

// aid is used to uniquely identify assets.
func (s *shader) aid() aid      { return s.tag }  // hashed type and name.
func (s *shader) label() string { return s.name } // asset name

// Shader source is scanned for uniforms and vertex buffer information.
// The uniform references are set on binding and later used by Model
// to set the uniform values during rendering.
func (s *shader) setSource(vsh, fsh []string) {
	s.vsh, s.fsh = vsh, fsh
	s.ensureNewLines()
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
//
// Naming and index conventions help shader programs interwork with one
// another and the engine.
// Common vertex attributes and indicies are:
//	 layout(location=0) in vec3 in_v; // vertex positions...........always 0.
//	 layout(location=1) in vec3 in_n; // vertex normals.............always 1.
//   layout(location=2) in vec2 in_t; // texture coordinates........always 2
//	 layout(location=6) in mat4 in_m; // instanced transform matrix.always 6-9.
// Common transform matricies set by the engine are:
// 	 uniform mat4 pm; // projection matrix
// 	 uniform mat4 vm; // view matrix
// 	 uniform mat4 mm; // model matrix
var shaderLibrary = map[string]func() (vsh, fsh []string){
	"colored":           func() (vsh, fsh []string) { return coloredShaderV, coloredShaderF },
	"coloredInstanced":  func() (vsh, fsh []string) { return coloredInstancedShaderV, coloredShaderF },
	"textured":          func() (vsh, fsh []string) { return texturedShaderV, texturedShaderF },
	"texturedInstanced": func() (vsh, fsh []string) { return texturedInstancedShaderV, texturedShaderF },
	"labeled":           func() (vsh, fsh []string) { return texturedShaderV, labeledShaderF },
	"sdf":               func() (vsh, fsh []string) { return texturedShaderV, sdfShaderF },
	"billboarded":       func() (vsh, fsh []string) { return billboardedShaderV, texturedShaderF },
	"animated":          func() (vsh, fsh []string) { return animatedShaderV, texturedShaderF },
	"diffuse":           func() (vsh, fsh []string) { return diffuseShaderV, diffuseShaderF },
	"phong":             func() (vsh, fsh []string) { return phongShaderV, phongShaderF },
	"toon":              func() (vsh, fsh []string) { return toonShaderV, toonShaderF },
	"normalMapped":      func() (vsh, fsh []string) { return normalMappedShaderV, normalMappedShaderF },
	"castShadow":        func() (vsh, fsh []string) { return castShadowShaderV, castShadowShaderF },
	"showShadow":        func() (vsh, fsh []string) { return showShadowShaderV, showShadowShaderF },
}

// FUTURE: Add edge-detect and emboss shaders, see:
//         http://www.processing.org/tutorials/pshader/
// ===========================================================================

// coloredShader shades all verticies the given color and alpha value.
var coloredShaderV = []string{
	"layout(location=0) in vec3 in_v;", // vertex positions.
	"",
	"uniform mat4 pm;", // projection matrix
	"uniform mat4 vm;", // view matrix
	"uniform mat4 mm;", // model matrix
	"void main(void) {",
	"   gl_Position = pm * vm * mm * vec4(in_v, 1.0);",
	"}",
}
var coloredShaderF = []string{
	"uniform vec3  kd;",       // material diffuse value
	"uniform float alpha;",    // transparency
	"out     vec4  f_color; ", // fragment color.
	"void main(void) {",
	"   f_color = vec4(kd, alpha);",
	"}",
}

// ===========================================================================

// coloredInstancedShader is an object instance aware coloredShader.
var coloredInstancedShaderV = []string{
	"layout(location=0) in vec3 in_v;", // verticies
	"layout(location=6) in mat4 in_m;", // instanced transform matrix.
	"",
	"uniform mat4  pm;", // projection matrix
	"uniform mat4  vm;", // view matrix
	"void main() {",
	"	gl_Position = pm * vm * in_m * vec4(in_v, 1.0);",
	"}",
}

// ===========================================================================

// texturedShader colors verticies based on texture sampling.
var texturedShaderV = []string{
	"layout(location=0) in vec3 in_v;", // verticies
	"layout(location=2) in vec2 in_t;", // texture coordinates
	"",
	"uniform mat4 pm;",  // projection matrix
	"uniform mat4 vm;",  // view matrix
	"uniform mat4 mm;",  // model matrix
	"out     vec2 v_t;", // pass on vertex texture coordinates.
	"void main() {",
	"   v_t = in_t;",
	"   gl_Position = pm * vm * mm * vec4(in_v, 1.0);",
	"}",
}

// The textures alpha areas are respected. The non-alpha areas can
// be overridden with an additional program alpha value.
var texturedShaderF = []string{
	"in      vec2      v_t;",     // interpoloated vertex texture coordinates.
	"uniform sampler2D uv;",      // texture sampler.
	"uniform float     alpha;",   // transparency overrides texture.
	"out     vec4      f_color;", // final fragment color
	"void main() {",
	"   f_color = texture(uv, v_t);",
	"   f_color.a *= alpha;", // Combine texture alpha with alpha override.
	"}",
}

// ===========================================================================

// texturedInstancedShader is an object instance aware texturedShader.
var texturedInstancedShaderV = []string{
	"layout(location=0) in vec3 in_v;", // verticies
	"layout(location=2) in vec2 in_t;", // texture coordinates
	"layout(location=6) in mat4 in_m;", // instanced world model matricies.
	"",
	"uniform mat4 pm;",  // projection matrix
	"uniform mat4 vm;",  // view matrix
	"out     vec2 v_t;", // pass on vertex texture coordinates.
	"void main() {",
	"   v_t = in_t;",
	"   gl_Position = pm * vm * in_m * vec4(in_v, 1.0);",
	"}",
}

// ===========================================================================

// labeledShader is a texture shader with a color override.
// Letters are the non-alpha portions of the texture.
var labeledShaderF = []string{
	"in      vec2      v_t;",     // interpoloated vertex texture coordinates.
	"uniform sampler2D uv;",      // texture sampler.
	"uniform vec3      kd;",      // text color
	"uniform float     alpha;",   // transparency
	"out     vec4      f_color;", // final fragment color
	"void main() {",
	"   float a = texture(uv, v_t).a;", // text is when a == 1.0
	"   f_color = vec4(kd,  a*alpha);", // alpha for transparent text.
	"}",
}

// ===========================================================================

// sdfShader displays signed distance field fonts.
//           This shader expects a font image file signed distance field
//           values for the font images. SDF images appear blurry.
var sdfShaderF = []string{
	"in      vec2      v_t;",     // interpoloated vertex texture coordinates.
	"uniform sampler2D uv;",      // texture sampler.
	"uniform vec3      kd;",      // text color
	"uniform float     alpha;",   // transparency
	"out     vec4      f_color;", // final fragment color
	"",
	"const float smoothing = 1.0/16.0;",
	"void main() {",
	"   float distance = texture(uv, v_t).a;", // distance field value.
	"   float clamp = smoothstep(0.5 - smoothing, 0.5 + smoothing, distance);",
	"   f_color = vec4(kd, alpha*clamp);",
	"}",
}

// ===========================================================================

// billboardedShader renders a single texture, like a textured shader,
// but forces the textured object to always face the camera. See
//     http://www.lighthouse3d.com/opengl/billboarding/billboardingtut.pdf
//     https://code.google.com/p/o3d/source/browse/trunk/samples_webgl/shaders/billboard.shader
//     http://www.sjbaker.org/steve/omniv/alpha_sorting.html
var billboardedShaderV = []string{
	"layout(location=0) in vec3 in_v;", // vertex coordinates
	"layout(location=2) in vec2 in_t;", // texture coordinates
	"",
	"uniform mat4 pm;",    // projection matrix
	"uniform mat4 vm;",    // view matrix
	"uniform mat4 mm;",    // model matrix
	"uniform vec3 scale;", // scale for each axis
	"out     vec2 v_t;",   // pass uv coordinates through
	"",
	"void main() {",
	"   v_t = in_t;",
	"   mat4 mvm = vm * mm;",
	"   gl_Position = pm * (vec4(in_v*scale, 1) + vec4(mvm[3].xyz, 0));",
	"}",
}

// =============================================================================

// animatedShader is a bare bones skeletal shader that includes
// uv texture mapping and alpha.
var animatedShaderV = []string{
	"layout(location=0) in vec3 in_v;",   // verticies
	"layout(location=2) in vec2 in_t;",   // texture coordinates
	"layout(location=4) in vec4 joint;",  // joint indicies
	"layout(location=5) in vec4 weight;", // joint weights
	"",
	"uniform mat4   pm;",        // projection matrix
	"uniform mat4   vm;",        // view matrix
	"uniform mat4   mm;",        // model matrix
	"uniform mat3x4 bpos[100];", // bone positioning transforms. Row-Major!
	"out     vec2   v_t;",       // pass texture coordinates through
	"",
	"void main() {",
	"   v_t = in_t;",
	"   mat3x4 m = bpos[int(joint.x)] * weight.x;", // up to four joints affect vertex.
	"   m += bpos[int(joint.y)] * weight.y;",
	"   m += bpos[int(joint.z)] * weight.z;",
	"   m += bpos[int(joint.w)] * weight.w;",
	"   vec4 mpos = vec4(vec4(in_v, 1.0) * m, 1.0);", // Row-Major pre-multiply.
	"   gl_Position = pm * vm * mm * mpos;",
	"}",
}

// ===========================================================================

// diffuseShader is based on
//      http://www.packtpub.com/article/opengl-glsl-4-shaders-basics
//      http://devmaster.net/posts/2974/the-basics-of-3d-lighting
//
// Diffuse: The algorithm is used to calculate the color for the polygon.
//          The polygon has only one normal and the only part of the above
//          algorithm used is: diffuse*(normal . light-direction)
var diffuseShaderV = []string{
	"layout(location=0) in vec3 in_v;", // verticies
	"layout(location=1) in vec3 in_n;", // vertex normals
	"",
	"uniform mat4 pm;",            // projection matrix
	"uniform mat4 vm;",            // view matrix
	"uniform mat4 mm;",            // model matrix
	"uniform vec3 lightPosition;", //
	"uniform vec3 lightColor;",    //
	"uniform vec3 kd;",            // material diffuse color.
	"out     vec3 v_color;",       // vertex color
	"void main() {",
	"   mat4 mvm = vm * mm;",          // view-model matrix.
	"   vec4 vmod = vec4(in_v, 1.0);", // vertex in model space.
	"   vec3 nm = normalize((mvm * vec4(in_n, 0)).xyz);", // unit normal in world space.
	"   vec3 lightDir = normalize(lightPosition - vec3(mvm*vmod));",
	"   v_color = lightColor * kd * max(dot(lightDir, nm), 0.0);",
	"   gl_Position = pm * mvm * vmod;", // pass on the transformed vertex position
	"}",
}
var diffuseShaderF = []string{
	"in      vec3  v_color;", // interpolated vertex color
	"uniform float alpha;",   // transparency
	"out     vec4  f_color;", // final fragment color
	"void main() {",
	"   f_color = vec4(v_color, alpha);",
	"}",
}

// ===========================================================================

// phongShader is based on
//       http://www.packtpub.com/article/opengl-glsl-4-shaders-basics
//
// Phong Shading   : Calculate the color intensity for each pixel using
//                   interpolated vertex normals.
var phongShaderV = []string{
	"layout(location=0) in vec3 in_v;", // verticies
	"layout(location=1) in vec3 in_n;", // vertex normals
	"",
	"uniform mat4 pm;",            // projection matrix
	"uniform mat4 vm;",            // view matrix
	"uniform mat4 mm;",            // model matrix
	"uniform vec3 lightPosition;", // light position in world space.
	"",
	"out     vec3 v_n;", // vertex normal
	"out     vec3 v_s;", // vector from vertex to light.
	"out     vec3 v_e;", // vertex eye position.
	"void main() {",
	"   mat4 mvm = vm * mm;",          // view-model matrix.
	"   vec4 vmod = vec4(in_v, 1.0);", // vertex in model space.
	"   vec4 vworld = mvm * vmod;",    // vertex in world space
	"   v_s = normalize(lightPosition - vworld.xyz);", // light vector
	"   v_e = normalize(-vworld.xyz);",                // view vector
	"   v_n = normalize((mvm * vec4(in_n, 0)).xyz);",  // unit normal in world space.
	"   gl_Position = pm * mvm * vmod;",               // vertex in clip space",
	"}",
}
var phongShaderF = []string{
	"in      vec3  v_n;",        // interpolated normal
	"in      vec3  v_s;",        // interpolated vector from vertex to light.
	"in      vec3  v_e;",        // interpolated vector from eye to vertex.
	"uniform vec3  lightColor;", //
	"uniform vec3  ka;",         // material ambient value
	"uniform vec3  ks;",         // material specular value
	"uniform vec3  kd;",         // material diffuse value
	"uniform float alpha;",      // transparency
	"uniform float ns;",         // material specular exponent
	"out     vec4  f_color;",    // final fragment color
	"void main() {",
	"   vec3 r = reflect(-v_s, v_n);",
	"   float sDotN = max( dot(v_s,v_n), 0.0 );",
	"   vec3 ambient = lightColor * ka;",
	"   vec3 diffuse = lightColor * kd * sDotN;",
	"   vec3 spec = vec3(0.0);",
	"   if (sDotN > 0.0)",
	"      spec = lightColor * ks * pow( max( dot(r,v_e), 0.0 ), ns);",
	"   vec3 color = ambient + diffuse + spec;", // combine all the values.
	"   f_color = vec4(color, alpha);",          // final fragment color
	"}",
}

// ===========================================================================

// toonShader is based on https://stackoverflow.com/questions/5795829/
// Creates distinct bands of the objects color based on the light intensity.
var toonShaderV = []string{
	"layout(location=0) in vec3 in_v;", // verticies
	"layout(location=1) in vec3 in_n;", // vertex normals
	"layout(location=2) in vec2 in_t;", // texture coordinates
	"",
	"uniform mat4 pm;",            // projection matrix
	"uniform mat4 vm;",            // view matrix
	"uniform mat4 mm;",            // model matrix
	"uniform vec3 lightPosition;", // light position in world space.
	"",
	"out     vec3 v_s;", // vector from vertex to light.
	"out     vec3 v_n;", // vertex normal
	"out     vec2 v_t;", // pass texture coordinates through
	"void main() {",
	"   mat4 mvm = vm * mm;",          // view-model matrix.
	"   vec4 vmod = vec4(in_v, 1.0);", // vertex in model space.
	"   vec4 vworld = mvm * vmod;",    // vertex in world space
	"   v_s = normalize(lightPosition - vworld.xyz);", // light vector
	"   v_n = normalize((mvm * vec4(in_n, 0)).xyz);",  // unit normal in world space.
	"   v_t = in_t;",                    //
	"   gl_Position = pm * mvm * vmod;", // vertex in clip space,
	"}",
}
var toonShaderF = []string{
	"in      vec3      v_n;",     // interpolated normal
	"in      vec3      v_s;",     // interpolated vector from vertex to light.
	"in      vec2      v_t;",     // texture coordinates.
	"uniform sampler2D uv;",      // texture sampler.
	"out     vec4      f_color;", // final fragment color
	"void main() {",
	"   vec4 c1 = texture(uv, v_t);", // texture color
	"   vec4 c2;",                    // color intensity from light energy.
	"   float lightEnergy = dot(v_s,v_n);",
	"   if (lightEnergy > 0.95)      c2 = vec4(1.0, 1.0, 1.0, 1.0);",
	"   else if (lightEnergy > 0.75) c2 = vec4(0.8, 0.8, 0.8, 1.0);",
	"   else if (lightEnergy > 0.50) c2 = vec4(0.6, 0.6, 0.6, 1.0);",
	"   else if (lightEnergy > 0.25) c2 = vec4(0.4, 0.4, 0.4, 1.0);",
	"   else                       c2 = vec4(0.2, 0.2, 0.2, 1.0);",
	"   f_color = c1 * c2;",
	"}",
}

// ===========================================================================

// normalMapped is a tangent based normal map shader based on code and
// concepts from http://www.thetenthplanet.de/archives/1180
// Generating the cotangent transform for each pixel avoids the need to
// send tangent information for each vertex.
var normalMappedShaderV = []string{
	"layout(location = 0) in vec3 in_v;", // vertex position in modelspace
	"layout(location = 1) in vec3 in_n;", // vertex normal in modespace
	"layout(location = 2) in vec2 in_t;", // vertex uv texture coordinates
	"",
	"uniform mat4 pm;",            // projection matrix
	"uniform mat4 vm;",            // view matrix
	"uniform mat4 mm;",            // model matrix
	"uniform vec3 lightPosition;", // light position.
	"",
	"out vec2 v_t;", // vertex texture coords
	"out vec3 v_n;", // vertex normal
	"out vec3 v_l;", // light to vertex vector
	"out vec3 v_v;", // view vector: ie: camera to vertex
	"",
	"void main() {",
	"   mat4 mvm = vm * mm;", // view-model matrix.
	"	vec4 vmod = vec4(in_v, 1.0);", // vertex in local model space.
	"	vec4 vcam = mvm * vmod;", // vertex in camera view space
	"	v_l = normalize(lightPosition - vcam.xyz);", // normalized vertex to light vector
	"	v_v = -vcam.xyz;", // non-normalized vertex view vector in view space
	"   v_n = (mvm * vec4(in_n, 0)).xyz;", // non-normalized vertex normal in view space.
	"   v_t = in_t;",                      // vertex UV texture coordinates.
	"   gl_Position = pm * mvm * vmod;",   // vertex position in clip space.
	"}",
}

// Fragment shader based on http://www.thetenthplanet.de/archives/1180
var normalMappedShaderF = []string{
	"in vec2 v_t;", // vertex texture coords
	"in vec3 v_n;", // non-normalized vertex normal
	"in vec3 v_v;", // non-normalized view vector
	"in vec3 v_l;", // normalized light to vertex vector
	"",
	"uniform sampler2D uv;",         // base diffuse color
	"uniform sampler2D uv1;",        // normal map texture in tangent space
	"uniform sampler2D uv2;",        // specular texture
	"uniform vec3      lightColor;", //
	"uniform vec3      ka;",         // material ambient color
	"uniform vec3      kd;",         // material diffuse color
	"uniform vec3      ks;",         // material specular color
	"uniform float     ns;",         // material specular shininess.
	"out     vec4      f_color;",
	"",
	// cotangent_frame creates the cotangent frame necessary to
	// map from the cotangent space to world space. Note that the
	// GLSL methods used are only available in the fragment shader.
	// Parameters:",
	//   N  : the interpolated vertex normal in world space
	//   p  : the reverse view vector in world space
	//   tuv: texture coordinates
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
	// perturb_normal transforms a normal from cotangent space
	// Note original code handles different normal map texture encodings.
	//   map = map * 255./127. - 128./127.;      #ifdef WITH_NORMALMAP_UNSIGNED
	//   map.z = sqrt(1. - dot(map.xy, map.xy)); #ifdef WITH_NORMALMAP_2CHANNEL
	//   map.y = -map.y;                         #ifdef WITH_NORMALMAP_GREEN_UP
	// Parameters:
	//   N  : the interpolated vertex normal in world space
	//   V  : the view vector in world space
	//   v_t: texture coordinates
	"vec3 perturb_normal( vec3 N, vec3 V, vec2 tuv ) {",
	"    vec3 map = texture(uv1, tuv).xyz;",       // perturbed normal in cotangent space
	"    map = map * 255./127. - 128./127.;",      // #ifdef WITH_NORMALMAP_UNSIGNED
	"    mat3 TBN = cotangent_frame(N, -V, tuv);", // cotangent to world space transform
	"    return normalize(TBN * map);",            // perturbed normal in world space
	"}",
	"void main() {",
	"    vec3 normal = normalize(v_n);",
	"    normal = perturb_normal(normal, v_v, v_t);", // normal in view space.
	"    vec3 nv_v = normalize(v_v);",
	"    vec3 ambient = lightColor * ka;",
	"    float intensity = max(dot(v_l, normal), 0.0);",
	"    vec3 diffuse = lightColor * kd * intensity;",
	"", // Blinn-Phong half vector.
	"    vec3 halfDir = normalize(v_l + nv_v);",
	"    float specAngle = max(dot(halfDir, normal), 0.0);",
	"    float specFac = pow(clamp(specAngle, 0.0, 1.0), ns);",
	"    vec3 smap = texture(uv2, v_t).rgb;",
	"    vec3 specular = lightColor * ks * smap * specFac;",
	"",                               // Combine into final fragment color.
	"    vec4 t = texture(uv, v_t);", // pure texture color.
	"    vec3 color = ambient*t.rgb + diffuse*t.rgb + specular;", // combine all the values.
	"    f_color = vec4(color, t.a);",                            // final fragment color
	"}",
}

// =============================================================================

// castShadowShader is used to create shadow maps by writing objects depths.
// Expected to be used during the shadow map render pass to render to
// a texture. See:
// http://www.opengl-tutorial.org/intermediate-tutorials/tutorial-16-shadow-mapping
var castShadowShaderV = []string{
	"layout (location = 0) in vec3 in_v;",
	"",
	"uniform mat4 pm;", // projection matrix
	"uniform mat4 vm;", // view matrix
	"uniform mat4 mm;", // model matrix
	"void main() {",
	"    gl_Position = pm * vm * mm * vec4(in_v, 1.0);",
	"}",
}
var castShadowShaderF = []string{
	"layout(location = 0) out float fragdepth;",
	"void main() {",
	"    fragdepth = gl_FragCoord.z;",
	"}",
}

// =============================================================================

// showShadowShader incorporates a shadow depth map into lighting calculations.
// http://www.opengl-tutorial.org/intermediate-tutorials/tutorial-16-shadow-mapping
var showShadowShaderV = []string{
	"layout(location=0) in vec3 in_v;", // verticies
	"layout(location=2) in vec2 in_t;", // texture coordinates
	"",
	"uniform mat4 pm;",  // projection matrix
	"uniform mat4 vm;",  // view matrix
	"uniform mat4 mm;",  // model matrix
	"uniform mat4 dbm;", // depth bias matrix
	"out     vec2 v_t;", // uv texture coordinate
	"out     vec4 v_s;", // shadow uv coordinates
	"void main(){",
	"    gl_Position = pm * vm * mm * vec4(in_v, 1.0);",
	"    v_s = dbm * vec4(in_v, 1.0);",
	"    v_t = in_t;",
	"}",
}
var showShadowShaderF = []string{
	"in      vec2 v_t;", // interpolated uv coordinates
	"in      vec4 v_s;", // interpolated shadow uv coordinates
	"",
	"uniform sampler2D       uv;",      // object material texture sampler
	"uniform sampler2DShadow sm;",      // shadow map depth texture sampler
	"out     vec4            f_color;", // final fragment color
	"void main(){",
	"    vec4 lightColor = vec4(1,1,1,1);",       // white light
	"    vec4 diffuseColor = texture(uv, v_t); ", // object color from texture
	"",
	"", // compare the depth found in the texture at xy
	"", // with the depth at z.
	"    vec2 suv = vec2((v_s.xy)/v_s.w);",
	"    float visibility = texture(sm, vec3(suv, (v_s.z)/v_s.w));",
	"    visibility = visibility + (1.0-visibility)*0.75;", // map 1-0 to 1-0.75
	"    f_color = visibility * diffuseColor * lightColor;",
	"}",
}
