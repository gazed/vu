// vertex shader nshade based on the direction of the vertex normal.

layout(location=0) in vec3 in_v; // vertex positions.
layout(location=1) in vec3 in_n; // vertex normals.

uniform mat4 pm;      // projection matrix
uniform mat4 vm;      // view matrix
uniform mat4 mm;      // model matrix
out     vec3 v_color; // vertex color based on normal direction.

void main(void)
{
   v_color = vec3(0.7, 0.6, 0.4) * (in_n.x + in_n.y + in_n.z);
   gl_Position = pm * vm * mm * vec4(in_v, 1.0);
}
