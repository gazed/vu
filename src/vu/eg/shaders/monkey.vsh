#version 150

uniform mat4 modelViewProjectionMatrix;
uniform mat3 normalMatrix;

in vec4 inPosition;
in vec3 inNormal;
out vec4 faceColor;
 
void main(void)
{
   gl_Position = modelViewProjectionMatrix * inPosition;
   vec3 color = vec3(0.7, 0.6, 0.4) * (inNormal.x + inNormal.y + inNormal.z);
   faceColor = vec4(color, 1.0); 
}
