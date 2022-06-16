#version 450

layout(location=0) out vec4 out_color;

// samplers
const int COLOR = 0;
layout(set = 1, binding = 0) uniform sampler2D samplers[1];

layout(location=0) in struct in_dto {
    vec2 texcoord;
} dto;

void main() {
    out_color = texture(samplers[COLOR], dto.texcoord);
}
