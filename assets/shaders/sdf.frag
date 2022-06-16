#version 450

layout(location=0) out vec4 out_color;

// Samplers
const int COLOR = 0;
layout(set = 1, binding = 0) uniform sampler2D samplers[1];

// FUTURE: convert to model push constant Uniforms
const vec3  text_color = vec3(0.0, 0.0, 0.0);
const float text_alpha = 1.0;

layout(location=0) in struct dto {
    vec2 texcoord;
} in_dto;

const float smoothing = 1.0/16.0;

void main() {
    float distance = texture(samplers[COLOR], in_dto.texcoord).a;
    float clamp = smoothstep(0.5 - smoothing, 0.5 + smoothing, distance);
    out_color = vec4(text_color, text_alpha*clamp);
}
