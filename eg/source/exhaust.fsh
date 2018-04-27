// Particle shader effect.

in      float     alpha; // Particle transparency.
in      float     index; // Particle index for texture atlas.
uniform sampler2D uv;    // Texture sampler.
out     vec4      ffc;   // Final fragment colour.

void main(void) {
	// generate the texture coordinates.
    float cols = 2;                  // Expecting 2x2 texture atlas.
    float quad = mod(index, 4);      // Random quadrant based on index.
    float u = floor(quad*0.5) * 0.5; // 0 or 0.5
    float v = mod(quad, 2)*0.5;      // 0 or 0.5
	vec2 t_uv = vec2(gl_PointCoord.s*0.5 + u, (1.0 - gl_PointCoord.t)*0.5 + v);
	vec4 texel = texture(uv, t_uv);

	// apply overall alpha in addition to the texture alpha.
	texel.a =  texel.a * alpha;
    ffc = texel;
}
