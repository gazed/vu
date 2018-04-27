layout(location=0) in vec3 in_v;
layout(location=2) in vec2 in_t;

uniform mat4 mvpm; // Model-View-Projection transform matrix.
out     vec2 v_t;  // texture map coordinates.

void main() {
    v_t = in_t;
	gl_Position = mvpm * vec4(in_v, 1.0);
}
