#version 150

in vec2 uvCoord;

uniform sampler2D uvSampler;

out vec4 fragColor;
 
void main() 
{
	fragColor = texture(uvSampler, uvCoord); 
}
