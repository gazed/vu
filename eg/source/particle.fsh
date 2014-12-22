#version 400 core

// shader to render simple particle system points 
// http://antongerdelan.net/opengl/particles.html

in      float     opacity;  // particle fade from the vertex shader.
uniform sampler2D uv;       // optional. enable point-sprite coords to use.
out     vec4      ffc;      // final fragment colour.

const vec4 particle_colour = vec4(0.5, 0.4, 0.8, 0.8);

void main () {

	// using point texture coordinates which are pre-defined over the point
	vec2 texcoord = vec2(gl_PointCoord.s, 1.0 - gl_PointCoord.t);
	vec4 texel = texture(uv, texcoord);
	ffc.rgb = particle_colour.rgb * texel.rgb;
	ffc.a = opacity * texel.a;
}
