#version 450

layout(location=0) out vec4 out_color;

// model uniforms
layout(push_constant) uniform push_constants {
    mat4 model; // 64 bytes

    // fragment shader uniforms
    vec4 color; // 16 bytes: rgba
} mu;


void main() {
    out_color = vec4(mu.color.xyz, 1.0);
}
