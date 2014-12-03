#version 330
 
// Skeletal shader based on the example data provided in the IQM
// Development kit from http://sauerbraten.org/iqm. Ignores lighting.

layout(location=0) in vec3 in_v;   // vertex coordinates
layout(location=2) in vec2 in_t;   // texture coordinates
layout(location=4) in vec4 joint;  // joint indicies
layout(location=5) in vec4 weight; // joint weights

uniform mat3x4 bpos[100];  // bone positioning transforms. Row-Major!
uniform mat4   mvpm;       // projection * model_view
out     vec2   t_uv;       // pass uv coordinates through

void main() {
   mat3x4 m = bpos[int(joint.x)] * weight.x;
   m += bpos[int(joint.y)] * weight.y;
   m += bpos[int(joint.z)] * weight.z;
   m += bpos[int(joint.w)] * weight.w;

   // Row-Major pre-multiply.
   vec4 mpos = vec4(vec4(in_v, 1.0) * m, 1.0); 
   gl_Position = mvpm * mpos;
   t_uv = in_t;
}
