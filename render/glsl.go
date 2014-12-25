// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package render

// glsl provides pre-made GLSL shaders. Each shader is identified by a
// unique name. These provide some basic shaders to get simple examples
// runing quickly and can be used as starting templates for new shaders.
var glsl = map[string]func() (vsh, fsh []string){
	"colour":  colourShader,
	"flat":    flatShader,
	"gouraud": gouraudShader,
	"phong":   phongShader,
	"uv":      uvShader,
	"bb":      bbShader,
	"bbr":     bbrShader,
	"widget":  widgetShader,
}

// FUTURE: Incorporate some (all?) of the blend algorithms from
//         http://devmaster.net/posts/3040/shader-effects-blend-modes
//         How do other engines handle/blend multiple textures?

// FUTURE: Add edge-detect and emboss shaders, see:
//         http://www.processing.org/tutorials/pshader/

// ===========================================================================

// colourShader shades all verticies the given diffuse colour.
func colourShader() (vsh, fsh []string) {
	vsh = []string{
		"#version 330",
		"layout(location=0) in vec3 in_v;",
		"",
		"uniform mat4 mvpm;", // projection * modelView
		"uniform vec3 kd;",   // material diffuse value
		"out     vec4 v_c;",  // vertex colour
		"void main(void) {",
		"   gl_Position = mvpm * vec4(in_v, 1.0);",
		"	v_c = vec4(kd, 1.0);",
		"}",
	}
	fsh = []string{
		"#version 330",
		"in  vec4 v_c;",  // color from vertex shader
		"out vec4 ffc; ", // final fragment colour.
		"void main(void) {",
		"   ffc = v_c;",
		"}",
	}
	return
}

// ===========================================================================

// flatShader combines a colour with alpha to make transparent objects.
//
// Flat Shading    : The algorithm is used to calculate the colour for the polygon.
//                   The polygon has only one normal and the only part of the above
//                   algorithm used is:
//                            diffuse*(normal . light-direction)
func flatShader() (vsh, fsh []string) {
	vsh = []string{
		"#version 330",
		"layout(location=0) in vec3 in_v;",
		"",
		"uniform mat4  mvpm;",  // projection * modelView
		"uniform vec3  kd;",    // diffuse colour
		"uniform float alpha;", // transparency
		"out     vec4  v_c;",   // vertex colour
		"void main() {",
		"	gl_Position = mvpm * vec4(in_v, 1.0);",
		"	v_c = vec4(kd, alpha);",
		"}",
	}
	fsh = []string{
		"#version 330",
		"in  vec4 v_c;", // color from vertex shader
		"out vec4 ffc;", // final fragment colour
		"void main() {",
		"   ffc = v_c;",
		"}",
	}
	return
}

// gouraudShader is based on
//      http://www.packtpub.com/article/opengl-glsl-4-shaders-basics
//      http://devmaster.net/posts/2974/the-basics-of-3d-lighting
//
// Gouraud Shading : The algorithm is used to calculate a colour at each vertex.
//                   The colours are then interpolated across the polygon.
func gouraudShader() (vsh, fsh []string) {
	vsh = []string{
		"#version 330",
		"layout(location=0) in vec3 in_v;", // vertex coordinates
		"layout(location=1) in vec3 in_n;", // vertex normal
		"",
		"uniform mat4  mvpm;",  // projection * model_view
		"uniform mat4  mvm;",   // model_view
		"uniform mat3  nm;",    // normal matrix
		"uniform vec4  l;",     // untransformed light position
		"uniform vec3  ld;",    // light source intensity
		"uniform vec3  kd;",    // material diffuse value
		"uniform float alpha;", // transparency
		"out     vec4  v_c;",   // vertex colour
		"void main() {",
		"   vec4 vpos = vec4(in_v, 1.0);",
		"   vec3 norm = normalize( nm * in_n);", // Convert normal and position to eye coords
		"   vec4 eyeCoords = mvm * vpos;",
		"   vec3 lightDirection = normalize(vec3(l - eyeCoords));",
		"   ",
		"   vec3 colour = ld * kd * max( dot( lightDirection, norm ), 0.0 );",
		"   v_c = vec4(colour, alpha);", // pass on the amount of diffuse light.
		"   gl_Position = mvpm * vpos;", // pass on the transformed vertex position
		"}",
	}
	fsh = []string{
		"#version 330",
		"in  vec4 v_c;", // color from vertex shader
		"out vec4 ffc;", // final fragment colour
		"void main() {",
		"   ffc = v_c;",
		"}",
	}
	return
}

// ===========================================================================

// phongShader is based on
//       http://www.packtpub.com/article/opengl-glsl-4-shaders-basics
//
// Phong Shading   : The full algorithm is calculated at each point.  This is done
//                   by having a normal at each vertex and interpolating the normals
//                   across the polygon.
func phongShader() (vsh, fsh []string) {
	vsh = []string{
		"#version 330",
		"layout(location=0) in vec3 in_v;", // vertex coordinates
		"layout(location=1) in vec3 in_n;", // vertex normal
		"",
		"uniform mat4  mvpm;",  // projection * model_view
		"uniform mat4  mvm;",   // model_view
		"uniform mat3  nm;",    // normal matrix
		"uniform vec4  l;",     // untransformed light position
		"uniform vec3  ld;",    // light source intensity
		"uniform vec3  ka;",    // material ambient value
		"uniform vec3  kd;",    // material diffuse value
		"uniform vec3  ks;",    // material specular value
		"uniform float alpha;", // transparency
		"out     vec4  v_c;",   // vertex colour
		"void main() {",
		"   vec4 vpos = vec4(in_v, 1.0);",
		"   vec3 norm = normalize( nm * in_n);",
		"   vec4 eyeCoords = mvm * vpos;",
		"   vec3 s = normalize(vec3(l - eyeCoords));",
		"   vec3 v = normalize(-eyeCoords.xyz);",
		"   vec3 r = reflect( -s, norm );",
		"   ",
		"   vec3 la = vec3(1.0);", // FUTURE make la a uniform.
		"   vec3 ambient = la * ka;",
		"   float sDotN = max( dot(s,norm), 0.0 );",
		"   vec3 diffuse = ld * kd * sDotN;",
		"   vec3 spec = vec3(0.0);",
		"   float shininess = 3.0;", // FUTURE make shininess a uniform.",
		"   vec3 ls = vec3(1.0);",   // FUTURE make ls a uniform.",
		"   if( sDotN > 0.0 )",
		"      spec = ls * ks * pow( max( dot(r,v), 0.0 ), shininess );",
		"   ",
		"   vec3 colour = ambient + diffuse + spec;", // combine all the values.
		"   v_c = vec4(colour, alpha);",              // pass on the vertex colour
		"   gl_Position = mvpm * vpos;",              // pass on the transformed vertex position",
		"}",
	}
	fsh = []string{
		"#version 330",
		"in      vec4      v_c;", // input colour
		"out     vec4      ffc;", // final fragment colour
		"void main() {",
		"   ffc = v_c;",
		"}",
	}
	return
}

// ===========================================================================

// uvShader handles a single texture.
func uvShader() (vsh, fsh []string) {
	vsh = []string{
		"#version 330",
		"layout(location=0) in vec3 in_v;", // vertex coordinates
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
		"out     vec4      ffc;",   // final fragment colour
		"void main() {",
		"   ffc = texture(uv, t_uv) * vec4(1.0, 1.0, 1.0, alpha);",
		"}",
	}
	return
}

// ===========================================================================

// bbShader is a billboard shader. Like a uv shader it renders a single texture
// but forces the textured object to always face the camera. See
//     http://www.lighthouse3d.com/opengl/billboarding/billboardingtut.pdf
func bbShader() (vsh, fsh []string) {
	vsh = []string{
		"#version 330",
		"layout(location=0) in vec3 in_v;", // vertex coordinates
		"layout(location=2) in vec2 in_t;", // texture coordinates
		"",
		"uniform mat4  mvpm;",  // projection * model_view
		"uniform vec3  scale;", // scale
		"out     vec2  t_uv;",  // pass uv coordinates through
		"",
		"vec4 billboard(mat4 bb, vec4 vpos, vec3 scale) {",
		"   bb[0][0] = 1.0;",
		"   bb[1][0] = mvpm[0][1];",
		"   bb[2][0] = 0.0;",
		"   bb[0][1] = 0.0;",
		"   bb[1][1] = 1.0;",
		"   bb[2][1] = 0.0;",
		"   bb[0][2] = 0.0;",
		"   bb[1][2] = mvpm[2][1];",
		"   bb[2][2] = 1.0;",
		"   vpos.xyz = vpos.xyz * scale;",
		"   return bb * vpos;",
		"}",
		"",
		"void main() {",
		"   gl_Position = billboard(mvpm, vec4(in_v, 1.0), scale);",
		"   t_uv = in_t;",
		"}",
	}
	fsh = []string{
		"#version 330",
		"in      vec2      t_uv;",
		"uniform sampler2D uv;",
		"uniform float     alpha;", // transparency
		"out     vec4      ffc;",   // final fragment colour
		"void main() {",
		"   ffc = texture(uv, t_uv) * vec4(1.0, 1.0, 1.0, alpha);",
		"}",
	}
	return
}

// ===========================================================================

// bbrShader is a billboard shader that rotates the texture over time.
func bbrShader() (vsh, fsh []string) {
	vsh, _ = bbShader()
	fsh = []string{
		"#version 330",
		"in      vec2      t_uv;",  // texture coordinates from vertex shader
		"uniform sampler2D uv;",    // sampler
		"uniform float     time;",  // current time in seconds
		"uniform float     spin;",  // rotation speed 0 -> 1
		"uniform float     alpha;", // transparency
		"out     vec4      ffc;",   // final fragment colour
		"",
		"void main() {",
		"   float sa = sin(time*spin);",               // calculate rotation
		"   float ca = cos(time*spin);",               // ..
		"   mat2 rot = mat2(ca, -sa, sa, ca);",        // ..
		"   ffc = texture(uv, ((t_uv-0.5)*rot)+0.5);", // rotate around its center
		"   ffc.a = ffc.a*alpha;",
		"}",
	}
	return
}

// ===========================================================================

// widgetShader combines a single texture, colour and alpha.
func widgetShader() (vsh, fsh []string) {
	vsh = []string{
		"#version 330",
		"layout(location=0) in vec3 in_v;", // vertex coordinates
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
		"uniform vec3      kd;",    // mouse over colour
		"uniform float     alpha;", // transparency
		"out     vec4      ffc;",   // final fragment colour
		"void main() {",
		"   ffc = texture(uv, t_uv) * vec4(kd.rgb, alpha);",
		"}",
	}
	return
}
