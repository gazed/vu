#version 450

// A PBR shader using a texture for the base color and
// fixed material values for metallic and roughness.

layout(location=0) in vec3 position; // vertex world location.
layout(location=1) in vec3 normal;   // vertex normal.
layout(location=2) in vec2 texcoord; // vertex texture coordinates.

layout(location=0) out struct out_dto {
    vec3 normal;
    vec3 world_pos;
    vec2 texcoord;
} dto;

// light is either a directional light
// or a point light with a position.
struct light {
    vec4 pos;   // if w == 1 position, else direction
    vec4 color; // xyz are rgb 0-1 and w is light intensity
};

// scene uniforms
layout(set=0, binding=0) uniform scene_uniforms {
    // vertex shader uniforms
    mat4 proj;       // 64 bytes
    mat4 view;       // 64 bytes

    // fragment shader uniforms
    vec4 cam;        // 16 bytes : local camera location

    // The first light must exist and be directional.
    // The remaining three lights are optional point lights.
    light lights[3]; // 96 bytes
    int nlights;     //  4 bytes : 1 to 4 lights
} su;

// model uniforms max 128 bytes
layout(push_constant) uniform push_constants {
    // vertex shader uniforms
    mat4 model;      // 64 bytes

    // fragment shader uniforms
    vec4 material;   // 16 bytes: x:metallic y:roughness
} mu;

void main() {
    dto.texcoord = texcoord;

    // calcuate unit normal in world space
	// FUTURE create the normal matix once on the CPU and pass in as a uniform.
	mat4 nmat = transpose(inverse(mu.model));
    dto.normal = normalize((nmat * vec4(normal, 0)).xyz);

    // calculate vertex world space position
    dto.world_pos = (mu.model * vec4(position, 1.0)).xyz;
    gl_Position = su.proj * su.view * mu.model * vec4(position, 1.0);
}
