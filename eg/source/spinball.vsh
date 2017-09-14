layout(location=0) in vec3 in_v;
layout(location=2) in vec2 in_t;

uniform mat4  mvm;   // model view matrix
uniform mat4  pm;    // projection matrix
uniform vec3  scale; // scale
out     vec2  t_uv;  // pass uv coordinates through

void main() {
   gl_Position = pm * (vec4(in_v*scale, 1) + vec4(mvm[3].xyz, 0));
   t_uv = in_t;
}
