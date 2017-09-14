in      vec2      uvCoord;
uniform sampler2D uv;
out     vec4      fragColor;

void main()
{
	fragColor = texture(uv, uvCoord);
}
