#version 330

layout(location=2) in vec2 in_t;
layout(location=0) in vec4 in_v;

uniform mat4 mvpm;
uniform vec3 scale;
out     vec2 t_uv;

// billboard ensures the object is camera facing.
vec4 billboard(mat4 mvpm, vec4 vpos, vec3 scale) {
   mat4 bb = mvpm;
   bb[0][0] = 1.0;
   bb[1][0] = mvpm[0][1];
   bb[2][0] = 0.0;
   bb[0][1] = 0.0;
   bb[1][1] = 1.0;
   bb[2][1] = 0.0;
   bb[0][2] = 0.0;
   bb[1][2] = mvpm[2][1];
   bb[2][2] = 1.0;
   vpos.xyz = vpos.xyz * scale;
   return bb * vpos;
}

void main() {
   t_uv = in_t;
   gl_Position = billboard(mvpm, in_v, scale);
}


