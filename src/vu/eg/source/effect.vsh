#version 330

uniform               mat4 mvpm; // model-view-projection matrix.
layout(location=0) in vec3 in_v; // particle locations.
 
void main(void)
{
   gl_Position = mvpm * vec4(in_v, 1.0); // particle location.
   gl_PointSize = 25.0;                  // particle size in pixels
}
