#version 330

in      vec2      t_uv;
uniform sampler2D uv;
out     vec4      ffc;   // final fragment colour

void main() {
   ffc = texture(uv, t_uv);
}
