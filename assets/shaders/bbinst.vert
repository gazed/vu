#version 450

// vertex attributes
layout(location=0) in vec3 position; // vertex position.
layout(location=1) in vec2 texcoord; // vertex texture coordinates.

// instance attributes
layout(location=2) in vec3 i_locus;  // model world position 
layout(location=3) in vec3 i_color;  // color to add to the base texture
layout(location=4) in float i_scale; // model scale factor 

// scene uniforms
layout(set=0, binding=0) uniform scene_uniforms {
    mat4 proj; // 64 bytes
    mat4 view; // 64 bytes
} su;

layout(location=0) out struct out_dto {
    vec3 color;
    vec2 texcoord;
} dto;

void main() {
    dto.color = i_color;
    dto.texcoord = texcoord;

    // construct the model transform.
    mat4 mm = mat4(vec4(1.0, 0.0, 0.0, 0.0),
                   vec4(0.0, 1.0, 0.0, 0.0),
                   vec4(0.0, 0.0, 1.0, 0.0),
                   vec4(i_locus, 1.0));
    mat4 mvm = su.view * mm;
    gl_Position = su.proj * (vec4(position*i_scale, 1) + vec4(mvm[3].xyz, 0));
}
