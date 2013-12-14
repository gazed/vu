// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package render

// shaders provides ready made shaders. These can be overridden with shader
// files from disk.
//
// FUTURE: figure out how to incorporate some (all?) of the blend algorithms
//         from http://devmaster.net/posts/3040/shader-effects-blend-modes
//         How do other engines handle/blend multiple textures?

import (
	"vu/data"
)

// CreateShader creates one of the known shaders. This is expected to be used by
// eng.LoadShader when a shader is not already in the cache. Nil is returned for
// unrecognized shaders. This function contains the list of official shader names.
func CreateShader(shaderName string, sh *data.Shader) {
	sh.Name = shaderName
	switch shaderName {
	case "flat":
		flatShader(sh)
	case "flata":
		flataShader(sh)
	case "gouraud":
		gouraudShader(sh)
	case "phong":
		phongShader(sh)
	case "uv":
		uvShader(sh)
	case "uva":
		uvaShader(sh)
	case "uvra":
		uvraShader(sh)
	case "uvm":
		uvmShader(sh)
	case "wave":
		waveShader(sh)
	case "bb":
		bbShader(sh)
	case "bba":
		bbaShader(sh)
	case "bbra":
		bbraShader(sh)
	default:
		sh.Name = "" // couldn't find the shader.
	}
}

// ===========================================================================

// flatShader is a basic shader used mostly in test programs.
//
// Flat Shading    : The algorithm is used to calculate the colour for the polygon.
//                   The polygon has only one normal and the only part of the above
//                   algorithm used is:
//                            diffuse*(normal . light-direction)
func flatShader(sh *data.Shader) {
	sh.Vsh = []string{
		"#version 150",
		"in      vec4  in_v;",  // vertex coordinates
		"uniform mat4  mvpm;",  // projection * modelView
		"uniform vec3  kd;",    // diffuse colour
		"uniform float alpha;", // transparency
		"out     vec4  v_c;",   // vertex colour
		"void main() {",
		"	gl_Position = mvpm * in_v;",
		"	v_c = vec4(kd, alpha);",
		"}",
	}
	sh.Fsh = []string{
		"#version 150",
		"in  vec4 v_c;", // color from vertex shader
		"out vec4 ffc;", // final fragment colour
		"void main() {",
		"   ffc = v_c;",
		"}",
	}
	sh.Uniforms = map[string]int32{
		"mvpm":  -1,
		"kd":    -1,
		"alpha": -1,
	}
}

// ===========================================================================

// flataShader is a flat shader that fades an object out based on
// distance from the viewer.
func flataShader(sh *data.Shader) {
	flatShader(sh)
	sh.Fsh = []string{
		"#version 150",
		"in      vec4  v_c;", // color from vertex shader
		"uniform float fd;",  // fade distance
		"out     vec4  ffc;", // final fragment colour
		"",
		"float fade(float distance) {",
		"   float z = gl_FragCoord.z / gl_FragCoord.w / distance;",
		"   z = clamp(z, 0.0, 1.0);",
		"   return 1.0 - z;",
		"}",
		"void main() {",
		"   ffc = v_c;",
		"   ffc.a = ffc.a*fade(fd);",
		"}",
	}
	sh.Uniforms["fd"] = -1
}

// ===========================================================================

// gouraudShader is based on
//      http://www.packtpub.com/article/opengl-glsl-4-shaders-basics
//      http://devmaster.net/posts/2974/the-basics-of-3d-lighting
//
// Gouraud Shading : The algorithm is used to calculate a colour at each vertex.
//                   The colours are then interpolated across the polygon.
func gouraudShader(sh *data.Shader) {
	sh.Vsh = []string{
		"#version 150",
		"in      vec4  in_v;",  // vertex coordinates
		"in      vec3  in_n;",  // vertex normal
		"uniform mat4  mvpm;",  // projection * model_view
		"uniform mat4  mvm;",   // model_view
		"uniform mat3  nm;",    // normal matrix
		"uniform vec4  l;",     // untransformed light position
		"uniform vec3  ld;",    // light source intensity
		"uniform vec3  kd;",    // material diffuse value
		"uniform float alpha;", // transparency
		"out     vec4  v_c;",   // vertex colour
		"void main() {",
		"   vec3 norm = normalize( nm * in_n);", // Convert normal and position to eye coords
		"   vec4 eyeCoords = mvm * in_v;",
		"   vec3 lightDirection = normalize(vec3(l - eyeCoords));",
		"   ",
		"   vec3 colour = ld * kd * max( dot( lightDirection, norm ), 0.0 );",
		"   v_c = vec4(colour, alpha);", // pass on the amount of diffuse light.
		"   gl_Position = mvpm * in_v;", // pass on the transformed vertex position
		"}",
	}
	sh.Fsh = []string{
		"#version 150",
		"in  vec4 v_c;", // color from vertex shader
		"out vec4 ffc;", // final fragment colour
		"void main() {",
		"   ffc = v_c;",
		"}",
	}
	sh.Uniforms = map[string]int32{
		"mvpm":  -1,
		"mvm":   -1,
		"nm":    -1,
		"l":     -1,
		"ld":    -1,
		"kd":    -1,
		"alpha": -1,
	}
}

// ===========================================================================

// phongShader is based on
//       http://www.packtpub.com/article/opengl-glsl-4-shaders-basics
//
// Phong Shading   : The full algorithm is calculated at each point.  This is done
//                   by having a normal at each vertex and interpolating the normals
//                   across the polygon.
func phongShader(sh *data.Shader) {
	sh.Vsh = []string{
		"#version 150",
		"in      vec4  in_v;",  // vertex coordinates
		"in      vec3  in_n;",  // vertex normal
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
		"   vec3 norm = normalize( nm * in_n);",
		"   vec4 eyeCoords = mvm * in_v;",
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
		"   gl_Position = mvpm * in_v;",              // pass on the transformed vertex position",
		"}",
	}
	sh.Fsh = []string{
		"#version 150",
		"in      vec4      v_c;", // input colour
		"out     vec4      ffc;", // final fragment colour
		"void main() {",
		"   ffc = v_c;",
		"}",
	}
	sh.Uniforms = map[string]int32{
		"mvpm":  -1,
		"mvm":   -1,
		"nm":    -1,
		"l":     -1,
		"ld":    -1,
		"ka":    -1,
		"kd":    -1,
		"ks":    -1,
		"alpha": -1,
	}
}

// ===========================================================================

// uvShader handles a single texture.
func uvShader(sh *data.Shader) {
	sh.Vsh = []string{
		"#version 150",
		"in      vec4  in_v;", // vertex coordinates
		"in      vec2  in_t;", // texture coordinates
		"uniform mat4  mvpm;", // projection * model_view
		"out     vec2  t_uv;", // pass uv coordinates through
		"void main() {",
		"   gl_Position = mvpm * in_v;",
		"   t_uv = in_t;",
		"}",
	}
	sh.Fsh = []string{
		"#version 150",
		"in      vec2      t_uv;",
		"uniform sampler2D uv;",
		"uniform float     alpha;", // transparency
		"out     vec4      ffc;",   // final fragment colour
		"void main() {",
		"   ffc = texture(uv, t_uv) * vec4(1.0, 1.0, 1.0, alpha);",
		"}",
	}
	sh.Uniforms = map[string]int32{
		"mvpm":  -1,
		"alpha": -1,
		"uv":    -1,
	}
}

// ===========================================================================

// uvaShader is a uvshader that fades an object out based on distance
// from the viewer.
func uvaShader(sh *data.Shader) {
	uvShader(sh)
	sh.Fsh = []string{
		"#version 150",
		"in      vec2      t_uv;",
		"uniform sampler2D uv;",
		"uniform float     fd;",    // fade distance
		"uniform float     alpha;", // transparency
		"out     vec4      ffc;",   // final fragment colour
		"",
		"float fade(float distance) {",
		"   float z = gl_FragCoord.z / gl_FragCoord.w / distance;",
		"   z = clamp(z, 0.0, 1.0);",
		"   return 1.0 - z;",
		"}",
		"void main() {",
		"   ffc = texture(uv, t_uv);",
		"   ffc.a = ffc.a*fade(fd)*alpha;",
		"}",
	}
	sh.Uniforms["fd"] = -1
}

// ===========================================================================

// uvraShader is a uvshader that rotates the texuture over time. It also
// fades an object out based on distance from the viewer.
func uvraShader(sh *data.Shader) {
	uvShader(sh)
	sh.Fsh = []string{
		"#version 150",
		"in      vec2      t_uv;",
		"uniform sampler2D uv;",
		"uniform float     fd;",    // fade distance
		"uniform float     time;",  // current time in seconds
		"uniform float     rs;",    // rotation speed 0 -> 1
		"uniform float     alpha;", // transparency
		"out     vec4      ffc;",   // final fragment colour
		"",
		"float fade(float distance) {",
		"   float z = gl_FragCoord.z / gl_FragCoord.w / distance;",
		"   z = clamp(z, 0.0, 1.0);",
		"   return 1.0 - z;",
		"}",
		"void main() {",
		"   float sa = sin(time*rs);",                 // calculate rotation
		"   float ca = cos(time*rs);",                 // ..
		"   mat2 rot = mat2(ca, -sa, sa, ca);",        // ..
		"   ffc = texture(uv, ((t_uv-0.5)*rot)+0.5);", // rotate around its center
		"   ffc.a = ffc.a*fade(fd)*alpha;",
		"}",
	}
	sh.Uniforms["fd"] = -1
	sh.Uniforms["time"] = -1
	sh.Uniforms["rs"] = -1
}

// ===========================================================================

// uvmShader is a white mask uvshader. This is used to get a white + alpha.
func uvmShader(sh *data.Shader) {
	uvShader(sh)
	sh.Fsh = []string{
		"#version 150",
		"in      vec2      t_uv;",
		"uniform sampler2D uv;",
		"out     vec4      ffc;", // final fragment colour
		"",
		"void main() {",
		"   ffc = texture(uv, t_uv);",
		" 	ffc.r = 1.0;",
		" 	ffc.g = 1.0;",
		" 	ffc.b = 1.0;",
		"}",
	}
	sh.Uniforms = map[string]int32{
		"mvpm": -1,
		"uv":   -1,
	}
}

// ===========================================================================

// bbShader is a billboard shader. Like a uv shader it renders a single texture
// but forces the textured object to always face the camera. See
//     http://www.lighthouse3d.com/opengl/billboarding/billboardingtut.pdf
func bbShader(sh *data.Shader) {
	sh.Vsh = []string{
		"#version 150",
		"in      vec4  in_v;",  // vertex coordinates
		"in      vec2  in_t;",  // texture coordinates
		"uniform mat4  mvpm;",  // projection * model_view
		"uniform vec3  scale;", // scale
		"out     vec2  t_uv;",  // pass uv coordinates through
		"void main() {",
		"   mat4 bb = mvpm;",
		"   bb[0][0] = 1.0;",
		"   bb[1][0] = mvpm[1][0];",
		"   bb[2][0] = 0.0;",
		"   bb[0][1] = 0.0;",
		"   bb[1][1] = 1.0;",
		"   bb[2][1] = 0.0;",
		"   bb[0][2] = 0.0;",
		"   bb[1][2] = mvpm[1][2];",
		"   bb[2][2] = 1.0;",
		"   vec4 vpos = in_v;",
		"   vpos.xyz = vpos.xyz * scale;",
		"   gl_Position = bb * vpos;",
		"   t_uv = in_t;",
		"}",
	}
	sh.Fsh = []string{
		"#version 150",
		"in      vec2      t_uv;",
		"uniform sampler2D uv;",
		"uniform float     alpha;", // transparency
		"out     vec4      ffc;",   // final fragment colour
		"void main() {",
		"   ffc = texture(uv, t_uv) * vec4(1.0, 1.0, 1.0, alpha);",
		"}",
	}
	sh.Uniforms = map[string]int32{
		"mvpm":  -1,
		"alpha": -1,
		"scale": -1,
		"uv":    -1,
	}
}

// ===========================================================================

// bbaShader is a billboard shader that fades an object out based on distance
// from the viewer.
func bbaShader(sh *data.Shader) {
	bbShader(sh)
	sh.Fsh = []string{
		"#version 150",
		"in      vec2      t_uv;",
		"uniform sampler2D uv;",
		"uniform float     fd;",    // fade distance
		"uniform float     alpha;", // transparency
		"out     vec4      ffc;",   // final fragment colour
		"",
		"float fade(float distance) {",
		"   float z = gl_FragCoord.z / gl_FragCoord.w / distance;",
		"   z = clamp(z, 0.0, 1.0);",
		"   return 1.0 - z;",
		"}",
		"void main() {",
		"   ffc = texture(uv, t_uv);",
		"   ffc.a = ffc.a*fade(fd)*alpha;",
		"}",
	}
	sh.Uniforms["fd"] = -1
}

// ===========================================================================

// bbraShader is a billboard shader that rotates the texture over time. It also
// fades an object out based on distance from the viewer.
func bbraShader(sh *data.Shader) {
	bbShader(sh)
	sh.Fsh = []string{
		"#version 150",
		"in      vec2      t_uv;",  // texture coordinates from vertex shader
		"uniform sampler2D uv;",    // sampler
		"uniform float     fd;",    // fade distance
		"uniform float     time;",  // current time in seconds
		"uniform float     rs;",    // rotation speed 0 -> 1
		"uniform float     alpha;", // transparency
		"out     vec4      ffc;",   // final fragment colour
		"",
		"float fade(float distance) {",
		"   float z = gl_FragCoord.z / gl_FragCoord.w / distance;",
		"   z = clamp(z, 0.0, 1.0);",
		"   return 1.0 - z;",
		"}",
		"void main() {",
		"   float sa = sin(time*rs);",                 // calculate rotation
		"   float ca = cos(time*rs);",                 // ..
		"   mat2 rot = mat2(ca, -sa, sa, ca);",        // ..
		"   ffc = texture(uv, ((t_uv-0.5)*rot)+0.5);", // rotate around its center
		"   ffc.a = ffc.a*fade(fd)*alpha;",
		"}",
	}
	sh.Uniforms["fd"] = -1
	sh.Uniforms["time"] = -1
	sh.Uniforms["rs"] = -1
}

// ===========================================================================

// waveShader is a fragment only effect shader.
// This is a modified version of http://glsl.heroku.com/e#8397.0
func waveShader(sh *data.Shader) {
	sh.Vsh = []string{
		"#version 150",
		"in      vec4  in_v;", // vertex coordinates
		"uniform mat4  mvpm;", // projection * model_view
		"",
		"void main() {",
		"	 gl_Position = mvpm * in_v;",
		"}",
	}
	sh.Fsh = []string{
		"#version 150",
		"uniform float time;",       // current time in seconds
		"uniform vec2  resolution;", // viewport size
		"uniform float alpha;",      // transparency
		"out     vec4  ffc;",        // final fragment colour
		"",
		"const float Pi     = 3.14159;",
		"const float fScale = 4.3;",
		"const float fEps   = 0.5;",
		"",
		"void main()  {",
		"	vec2 p = (2.0*gl_FragCoord.xy-resolution)/max(resolution.x,resolution.y);",
		"	for(int i=1; i<100; i++) {",
		"		vec2 newp =p;",
		"		newp.x += 1.5/float(i)*sin(float(i)*p.y+time/40.0+0.3*float(i))+400./20.0;",
		"		newp.y += 0.05/float(i)*sin(float(i)*p.x+time/1.0+0.3*float(i+10))-400./20.0+15.0;",
		"		p = newp;",
		"	}",
		"	vec3 col = vec3(0.5*sin(3.0*p.x)+0.5,0.5*sin(3.0*p.y)+0.5,sin(p.x+p.y));",
		"	vec3 lum = vec3(0.299,0.587,0.114);",
		"	vec3 c = vec3(dot(col*0.2,lum));",
		"	ffc = vec4(c, 1.0);",
		"	ffc.a = ffc.a*alpha;",
		"}",
	}
	sh.Uniforms = map[string]int32{
		"mvpm":       -1,
		"alpha":      -1,
		"time":       -1,
		"resolution": -1,
	}
}
