#version 150

in vec4 vertexPosition;
in vec2 uvPoint;

uniform mat4 Mvpm;      // Projection * ModelView

out vec2 uvCoord;

void main() 
{
	gl_Position = Mvpm * vertexPosition;
    uvCoord = uvPoint;
}
