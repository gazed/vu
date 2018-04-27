in      vec2      v_t;     // interpolated texture coordinates.
uniform sampler2D uv;      // texture sampler.
out     vec4      f_color; // final fragment color.

void main() {
	f_color = texture(uv, v_t);
}
