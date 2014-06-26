#version 330

uniform sampler2D uv;    // Texture sampler.
uniform float     alpha; // Transparency.
out     vec4      ffc;   // Final fragment colour.

void main(void) 
{
	// generate the texture coordinates.
	vec2 t_uv = vec2(gl_PointCoord.s, 1.0 - gl_PointCoord.t);
	vec4 texel = texture(uv, t_uv);

	// apply overall alpha in addition to the texture alpha.
	texel.a =  texel.a * alpha;
    ffc = texel;
}
