#version 330

in      vec2      t_uv;
uniform sampler2D uv0;
uniform sampler2D uv1;
uniform sampler2D uv2;
uniform sampler2D uv3;
uniform float     time;
out     vec4      ffc;

// sping calculates rotated uv coordinates.
vec2 spin(vec2 coords, float now, float rotSpeed) {
   float sa = sin(now*rotSpeed);
   float ca = cos(now*rotSpeed);
   mat2 rot = mat2(ca, -sa, sa, ca);
   return ((coords-0.5)*rot)+0.5;
}

void main() {
   vec4 t0 = texture(uv0, spin(t_uv, time, 1));
   vec4 t1 = texture(uv1, spin(t_uv, time, -0.75));
   vec4 t2 = texture(uv2, spin(t_uv, time, 1.5));
   vec4 t3 = texture(uv3, spin(t_uv, time, -2));
   ffc = mix(mix(t0, t1, 0.5), mix(t2, t3, 0.5), 0.5);
}

