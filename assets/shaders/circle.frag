#version 450

layout(location=0) out vec4 out_color;

layout(location=0) in struct in_dto {
    vec2 texcoord;
} dto;

void main() {
    // change texcoord 0:1 to uv -1.0:1.0
    vec2 uv = dto.texcoord * 2.0 -1;

    // line parameters
    float thickness = 0.01;
    float fade = 0.002;
    float dist = 1.0-length(uv);

    // circle outer edge
    vec3 color = vec3(smoothstep(0.0, fade, dist));

    // circle inner edge
    color *= vec3(smoothstep(thickness, thickness-fade, dist));

    // ignore black pixels completely.
    if ((color.x == 0) && (color.y == 0) && (color.z == 0))
        discard;

    // final color
    out_color = vec4(color, 1.0);
}
