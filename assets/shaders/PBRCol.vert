#version 450

// A PBR shader using material values instead of textures.
// Requires one directional light.

layout(location=0) in vec3 position; // vertex world location.
layout(location=1) in vec3 normal;   // vertex normal.

layout(location=0) out struct out_dto {
    vec3 normal;   // normalized vertex normal.
    vec3 v_world;  // vertex world position.
} dto;

// light is either a directional light, point light, or spot light.
// - a sun light uses: color, intensity,         , direction.
// - point lights use: color, intensity, position,          , attenuation.
// - spot lights use : color, intensity, position, direction, attenuation, cutoff.
struct light {
    vec4 color; // XYZ are rgb 0-1, W is light intensity
    vec4 pos;   // XYZ is the world space light position, W is attenuation.
    vec4 dir;   // XYZ is world space light direction, W is the cos(cutoff) angle in radians.
};

// scene uniforms currently fit in 512 bytes,
// and align on a 256 bytes boundary.
layout(set=0, binding=0) uniform scene_uniforms {
    // vertex shader uniforms
    mat4 proj;       //  64 bytes
    mat4 view;       //  64 bytes

    // fragment shader uniforms
    vec4 cam;        //  16 bytes : local camera location

	// fragment shader supports up to a max of 5 lights. Ordered so that
	// sun lights are first, followed by point lights, followed by spot lights.
	// eg: 1 sun light, 2 point lights, 1 spot light for a total of 4 lights.
	//
	// the number of each type of light must match the order and type stored in lights.
	// lightCnt[0] = number of sun lights.
	// lightCnt[1] = number of point lights.
	// lightCnt[2] = number of spot lights.
    ivec4 lightCnt;  //  16        bytes
    light lights[5]; // 240 (5*48) bytes
} su;

// model uniforms max 128 bytes
layout(push_constant) uniform push_constants {
    // vertex shader uniforms
    mat4 model;      // 64 bytes

    // fragment shader uniforms
    vec4 color;      // 16 bytes: rgba
    vec4 material;   // 16 bytes: x:metallic y:roughness
} mu;

void main() {

    // calcuate unit normal in world space
	// https://www.lighthouse3d.com/tutorials/glsl-12-tutorial/the-normal-matrix/
	// FUTURE create the normal matix once on the CPU and pass in as a uniform.
	mat4 nmat = transpose(inverse(mu.model));
    dto.normal = normalize((nmat * vec4(normal, 0)).xyz);


    // calculate vertex world space position
    dto.v_world = (mu.model * vec4(position, 1.0)).xyz;
    gl_Position = su.proj * su.view * mu.model * vec4(position, 1.0);
}
