#version 450

layout(location=0) out vec4 frag_color;

layout(location=0) in struct in_dto {
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
    vec4 cam;        //  16 bytes : world camera location

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

vec3 PBRLighting(light Light, vec3 l, vec3 intensity, vec3 Normal) {
    float metallic = float(round(mu.material.x)); // 0.0 or 1.0
    float roughness = mu.material.y;              // 0.0 to 1.0

    // object normal vector, view vector, half vector.
    vec3 n = Normal;
    vec3 v = normalize(vec3(su.cam) - dto.v_world);
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
    vec3 FinalColor = (DiffuseBRDF*alpha + SpecBRDF) * intensity * nDotL;
    return FinalColor;
}

void main() {
    vec3 N = dto.normal; // normalized vertex normal.

	// setup light calculation variables.
	int sunLights = su.lightCnt.x;
	int pointLights = su.lightCnt.y;
	int spotLights = su.lightCnt.z;
    vec3 totalLight = vec3(0);
	vec3 lightDir;
	float lightDist;
	vec3 intensity;
	light lightData;

	// by convention, directional lights are first. generally, expect one.
    for (int i = 0; i < sunLights; i++) {
		lightData = su.lights[i];
        lightDir = normalize(lightData.pos.xyz);
        intensity = lightData.color.xyz * lightData.color.w; // color * intensity
		totalLight += PBRLighting(su.lights[i], lightDir, intensity, N);
    }

	// point lights, if any, are next.
	int pointIndex = sunLights;
    for (int i = pointIndex; i < pointIndex+pointLights; i++) {
		lightData = su.lights[i];
        lightDir = lightData.pos.xyz - dto.v_world;
        float lightDist = length(lightDir); // get distance before normalization.
        lightDir = normalize(lightDir);     // normalize light direction.
        intensity = lightData.color.xyz * lightData.color.w; // color * intensity
        intensity /= (lightDist * lightDist);                // drops off by distance squared.
        totalLight += PBRLighting(su.lights[i], lightDir, intensity, N);
    }

	// spot lights, if any, are last.
	int spotIndex = sunLights+pointLights;
    for (int i = spotIndex; i < spotIndex+spotLights; i++) {
	lightData = su.lights[i];
        lightDir = normalize(lightData.pos.xyz - dto.v_world);
		float beamEdge = dot(lightDir, normalize(-lightData.dir.xyz));
		float cutoff = lightData.dir.w;
		intensity = vec3(0); // outside spot light beam
		if (beamEdge > cutoff) {
			// inside spot light beam
			// TODO add intensity fall off with distance..
			intensity = lightData.color.xyz * lightData.color.w; // color * intensity

			// add blur to beam edge.
			float falloff =  1.0 - ((1.0 - beamEdge)/(1.0 - cutoff));
			intensity = intensity*falloff;
		}
        totalLight += PBRLighting(lightData, lightDir, intensity, N);
    }

    // add fixed (indirect) ambient value as a cheap replacement for global illumination.
    float ambientStrength = 0.0005;
    totalLight += (ambientStrength * mu.color).xyz;

    // HDR tone mapping
    totalLight = totalLight / (totalLight + vec3(1.0));

    // Gamma correction
    float alpha = mu.color.w;
    frag_color = vec4(pow(totalLight, vec3(1.0/2.2)), alpha);
}
