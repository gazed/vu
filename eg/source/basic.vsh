uniform mat4 mvpm;

layout(location=0) in vec3 in_v;
layout(location=3) in vec4 in_c;

out vec4 ex_Color;

void main(void)
{
   gl_Position = mvpm * vec4(in_v, 1.0);
   ex_Color = in_c;
}
