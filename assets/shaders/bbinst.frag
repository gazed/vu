#version 450

layout(location=0) out vec4 out_color;

// samplers
const int COLOR = 0;
layout(set = 1, binding = 0) uniform sampler2D samplers[1];

layout(location=0) in struct in_dto {
    vec3 color;
    vec2 texcoord;
} dto;

void main() {
    vec4 base_color = texture(samplers[COLOR], dto.texcoord);
    vec3 colorized = clamp(base_color.xyz * dto.color, 0.0, 1.0);

    // modify the color, keeping the base_color alpha.
    out_color = vec4(colorized, base_color.w);
}
