#version 150

uniform mat4 Mvpm;

in vec4 in_Position;
in vec4 in_Color;
out vec4 ex_Color;
 
void main(void)
{
   gl_Position = Mvpm * in_Position;
   ex_Color = in_Color;
}
