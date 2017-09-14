layout(location=0) in vec3 in_v;

uniform mat4 mvpm; // Projection * ModelView

void main() {
	gl_Position = mvpm * vec4(in_v, 1.0);
}
