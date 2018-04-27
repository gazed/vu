// land is an experimental surface shader.
// It expects land textures to be stacked in an atlas and it will blend
// a blend texture with the next texture down in the atlas.

layout(location=0) in vec3  in_v;   // vertex coordinates
layout(location=1) in vec3  in_n;   // vertex normal
layout(location=2) in vec4  in_t;   // vertex texture uv coordinates + base/ratio.

uniform float ratio;  // texture to texture atlas ratio.
uniform mat4  pm;     // projection transform matrix.
uniform mat4  vm;     // view transform matrix
uniform mat4  mm;     // model transform matrix
out     vec3  f_nm;   // output vertex normal.
out     vec2  tuv0;   // uv coordinates
out     vec2  tuv1;   // uv coordinates
out     float weight; // texture blend weighting

void main() {
   mat4 mvm = vm * mm;
   gl_Position = pm * mvm * vec4(in_v, 1.0);
   f_nm = normalize((mvm * vec4(in_n, 0)).xyz);
   float blend = in_t.z;
   tuv0 = vec2(in_t.x, in_t.y+(blend*ratio));
   tuv1 = vec2(in_t.x, in_t.y+((blend+1)*ratio));
   weight = in_t.w;
}
