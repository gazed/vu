#version 450

layout(location=0) out vec4 out_color;

// Samplers
const int COLOR = 0;
layout(set = 1, binding = 0) uniform sampler2D samplers[1];

layout(location=0) in struct in_dto {
    vec2 texcoord;
} dto;

// model uniforms
layout(push_constant) uniform push_constants {
    mat4 model; // 64 bytes

    // fragment shader uniforms
    vec4 color; // 16 bytes: rgba
} mu;

void main() {
    float alpha = texture(samplers[COLOR], dto.texcoord).a;
    out_color = vec4(mu.color.x, mu.color.y, mu.color.z, alpha);
}
