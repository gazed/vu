// vcolor vertex shader uses colors assigned to each vertex.

layout(location=0) in vec3 in_v; // Vertex position.
layout(location=3) in vec4 in_c; // Vertex color.

uniform mat4 mvpm;    // Combined model-view-projection matrix.
out     vec4 v_color; // Vertex color passed to fragment shader.

void main(void)
{
   v_color = in_c;
   gl_Position = mvpm * vec4(in_v, 1.0);
}
