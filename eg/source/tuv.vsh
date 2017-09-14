layout(location=0) in vec3 in_v;
layout(location=2) in vec2 in_t;

uniform mat4 mvpm;      // Projection * ModelView
out     vec2 uvCoord;

void main() 
{
	gl_Position = mvpm * vec4(in_v, 1.0);
    uvCoord = in_t;
}
