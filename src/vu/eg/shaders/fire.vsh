#version 150
#
in      vec4 vertexPosition;
uniform mat4 Mvpm;            // Projection * ModelView

void main() {
	gl_Position = Mvpm * vertexPosition;
}
