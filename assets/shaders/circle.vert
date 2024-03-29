#version 450

layout(location=0) in vec3 position;
layout(location=1) in vec2 texcoord;

// scene uniforms
layout(set=0, binding=0) uniform scene_uniforms {
    mat4 proj; // 64 bytes
    mat4 view; // 64 bytes
} su;

// model uniforms
layout(push_constant) uniform push_constants {
	mat4 model; // 64 bytes
} mu;

layout(location=0) out struct out_dto {
    vec2 texcoord;
} dto;

void main() {
    dto.texcoord = texcoord;
    gl_Position = su.proj * su.view * mu.model * vec4(position, 1.0);
}
