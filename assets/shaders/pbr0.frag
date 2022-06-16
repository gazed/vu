#version 450

layout(location=0) out vec4 frag_color;

layout(location=0) in struct in_dto {
    vec3 normal;
    vec3 world_pos;
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
    vec4 color;      // 16 bytes: rgba
    vec4 material;   // 16 bytes: x:metallic y:roughness
} mu;

#define PI 3.1415926535897932384626433832795

// uniforms and constants
// =============================================================================
// code

// specular BRDF Fresnel function
vec3 schlickFresnel(float vDotH, float metallic) {
    vec3 F0 = vec3(0.04); // specular color for non-metals

    // use material color for metals
    F0 = mix(F0, vec3(mu.color), metallic);
    vec3 ret = F0 + (1 - F0) * pow(clamp(1.0 - vDotH, 0.0, 1.0), 5);
    return ret;
}

// specular BRDF geometry function
float geomSmith(float dp, float roughness) {
    float k = (roughness + 1.0) * (roughness + 1.0) / 8.0;
    float denom = dp * (1 - k) + k;
    return dp / denom;
}

// specular BRDF normal distribution funtion.
float ggxDistribution(float nDotH, float roughness) {
    float alpha2 = roughness * roughness * roughness * roughness;
    float d = nDotH * nDotH * (alpha2 - 1) + 1;
    float ggxdistrib = alpha2 / (PI * d * d);
    return ggxdistrib;
}

vec3 CalcPBRLighting(light Light, bool IsDirLight, vec3 Normal) {
    vec3 LightIntensity = Light.color.xyz * Light.color.w; // color * intensity
    vec3 l = vec3(0.0);
    float metallic = float(round(mu.material.x)); // 0.0 or 1.0
    float roughness = mu.material.y;              // 0.0 to 1.0

    vec3 PosDir = vec3(Light.pos);
    if (IsDirLight) {
        l = normalize(-PosDir.xyz);
    } else {
        l = PosDir - dto.world_pos;
        float LightToPixelDist = length(l);
        l = normalize(l);
        LightIntensity /= (LightToPixelDist * LightToPixelDist);
    }

    // object normal vector, view vector, half vector.
    vec3 n = Normal;
    vec3 v = normalize(vec3(su.cam) - dto.world_pos);
    vec3 h = normalize(v + l);
    float nDotH = max(dot(n, h), 0.0);
    float vDotH = max(dot(v, h), 0.0);
    float nDotL = max(dot(n, l), 0.0);
    float nDotV = max(dot(n, v), 0.0);

    // conserve energy so refaction+reflection==1.0
    vec3 F = schlickFresnel(vDotH, metallic);
    vec3 kS = F;          // specular: reflection
    vec3 kD = 1.0 - kS;   // diffuse : refraction

    // specular BRDF.
    vec3 SpecBRDF_nom  = ggxDistribution(nDotH, roughness) *
                         F *
                         geomSmith(nDotL, roughness) *
                         geomSmith(nDotV, roughness);
    float SpecBRDF_denom = 4.0 * nDotV * nDotL + 0.0001;
    vec3 SpecBRDF = SpecBRDF_nom / SpecBRDF_denom;

    // use color for non-metals, metals will already have the color in kS.
    vec3 fLambert = mix(vec3(mu.color), vec3(0.0), metallic);
    vec3 DiffuseBRDF = kD * fLambert / PI;

    // final color value for the given light.
    float alpha = mu.color.w;
    vec3 FinalColor = (DiffuseBRDF*alpha + SpecBRDF) * LightIntensity * nDotL;
    return FinalColor;
}

void main() {
    vec3 N = normalize(dto.normal);
    vec3 TotalLight = CalcPBRLighting(su.lights[0], true, N);
    for (int i = 1; i < su.nlights; i++) {
        TotalLight += CalcPBRLighting(su.lights[i], false, N);
    }

    // add fixed (indirect) ambient value as a cheap replacement for global illumination.
    float ambientStrength = 0.0005;
    TotalLight += (ambientStrength * mu.color).xyz;

    // HDR tone mapping
    TotalLight = TotalLight / (TotalLight + vec3(1.0));

    // Gamma correction
    float alpha = mu.color.w;
    frag_color = vec4(pow(TotalLight, vec3(1.0/2.2)), alpha);
}
