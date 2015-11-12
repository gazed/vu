#version 330

layout(location=0) in vec3 in_v;
layout(location=1) in vec3 in_n;

uniform mat4 mvpm;
out     vec4 faceColor;
 
void main(void)
{
   gl_Position = mvpm * vec4(in_v, 1.0);
   vec3 color = vec3(0.7, 0.6, 0.4) * (in_n.x + in_n.y + in_n.z);
   faceColor = vec4(color, 1.0); 
}
