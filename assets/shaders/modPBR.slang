// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT
//
module modPBR;

public struct SceneUniforms
{
	// vertex shader uniforms
	public float4x4 proj;
	public float4x4 view;

	// fragment shader uniforms
	public float4 cam;		  // local camera location

	// fragment shader supports up to a max of 5 lights. Ordered so that
	// sun lights are first, followed by point lights, followed by spot lights.
	// eg: 1 sun light, 2 point lights, 1 spot light for a total of 4 lights.
	//
	// the number of each type of light must match the order and type stored in lights.
	// lightCnt[0] = number of sun lights.
	// lightCnt[1] = number of point lights.
	// lightCnt[2] = number of spot lights.
	public int4 lightCnt;
	public Light lights[5];
};
layout(set=0, binding=0) public ConstantBuffer<SceneUniforms> su;

// Light is either a directional light, point light, or spot light.
// - a sun light uses: color, intensity,		 , direction.
// - point lights use: color, intensity, position,			, attenuation.
// - spot lights use : color, intensity, position, direction, attenuation, cutoff.
public struct Light {
	public float4 color; // XYZ are rgb 0-1, W is light intensity
	public float4 pos;	 // XYZ is the world space light position, W is attenuation.
	public float4 dir;	 // XYZ is world space light direction, W is the cos(cutoff) angle in radians.
};

#define PI 3.1415926535897932384626433832795

// specular BRDF Fresnel function
float3 schlickFresnel(float vDotH, float metallic, float3 color) {
	float3 F0 = float3(0.04); // specular color for non-metals

	// use material color for metals
	F0 = lerp(F0, color, metallic);
	float3 ret = F0 + (1 - F0) * pow(clamp(1.0 - vDotH, 0.0, 1.0), 5);
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

public float3 PBRLighting(Light light, float3 l, float3 intensity, float3 Normal, float3 vworld, PBRMat material) {
	float4 color	 = material.color;
	float  metallic  = material.metallic;
	float  roughness = material.roughness;

	// object normal vector, view vector, half vector.
	float3 n = Normal;
	float3 v = normalize(su.cam.xyz - vworld);
	float3 h = normalize(v + l);
	float nDotH = max(dot(n, h), 0.0);
	float vDotH = max(dot(v, h), 0.0);
	float nDotL = max(dot(n, l), 0.0);
	float nDotV = max(dot(n, v), 0.0);

	// conserve energy so refaction+reflection==1.0
	float3 F = schlickFresnel(vDotH, metallic, color.xyz);
	float3 kS = F;			// specular: reflection
	float3 kD = 1.0 - kS;	// diffuse : refraction

	// specular BRDF.
	float3 SpecBRDF_nom  = ggxDistribution(nDotH, roughness) *
						   F *
						   geomSmith(nDotL, roughness) *
						   geomSmith(nDotV, roughness);
	float SpecBRDF_denom = 4.0 * nDotV * nDotL + 0.0001;
	float3 SpecBRDF = SpecBRDF_nom / SpecBRDF_denom;

	// use color for non-metals, metals will already have the color in kS.
	float3 fLambert = lerp(color.xyz, float3(0.0), metallic);
	float3 DiffuseBRDF = kD * fLambert / PI;

	// final color value for the given light.
	float alpha = color.w;
	float3 FinalColor = (DiffuseBRDF*alpha + SpecBRDF) * intensity * nDotL;
	return FinalColor;
}

public struct FragIn
{
	public float3 normal;  // normalized vertex normal.
	public float3 vworld;  // vertex world position.
};

public struct PBRMat
{
	public float4 color;	 // rgba
	public float  metallic;  // 0.0 or 1.0;
	public float  roughness; // 0.0 to 1.0
};

// common fragment shader code.
public float4 PBRFrag(FragIn f, PBRMat material) {
	float3 N = normalize(f.normal); // re-normalized interpolated vertex normal.

	// setup light calculation variables.
	int sunLights = su.lightCnt.x;
	int pointLights = su.lightCnt.y;
	int spotLights = su.lightCnt.z;
	float3 totalLight = float3(0);
	float3 lightDir;
	float  lightDist;
	float3 intensity;
	Light  lightData;

	// by convention, directional lights are first. generally, expect one.
	for (int i = 0; i < sunLights; i++) {
		lightData = su.lights[i];
		lightDir = normalize(lightData.pos.xyz);
		intensity = lightData.color.xyz * lightData.color.w; // color * intensity
		totalLight += PBRLighting(su.lights[i], lightDir, intensity, N, f.vworld, material);
	}

	// point lights, if any, are next.
	int pointIndex = sunLights;
	for (int i = pointIndex; i < pointIndex+pointLights; i++) {
		lightData = su.lights[i];
		lightDir = lightData.pos.xyz - f.vworld;
		float lightDist = length(lightDir); // get distance before normalization.
		lightDir = normalize(lightDir);		// normalize light direction.
		intensity = lightData.color.xyz * lightData.color.w; // color * intensity
		intensity /= (lightDist * lightDist);				 // drops off by distance squared.
		totalLight += PBRLighting(su.lights[i], lightDir, intensity, N, f.vworld, material);
	}

	// spot lights, if any, are last.
	int spotIndex = sunLights+pointLights;
	for (int i = spotIndex; i < spotIndex+spotLights; i++) {
		lightData = su.lights[i];
		lightDir = normalize(lightData.pos.xyz - f.vworld);
		float beamEdge = dot(lightDir, normalize(-lightData.dir.xyz));
		float cutoff = lightData.dir.w;
		intensity = float3(0); // outside spot light beam
		if (beamEdge > cutoff) {
			// inside spot light beam
			// TODO add intensity fall off with distance..
			intensity = lightData.color.xyz * lightData.color.w; // color * intensity

			// add blur to beam edge.
			float falloff =  1.0 - ((1.0 - beamEdge)/(1.0 - cutoff));
			intensity = intensity*falloff;
		}
		totalLight += PBRLighting(lightData, lightDir, intensity, N, f.vworld, material);
	}

	// add fixed (indirect) ambient value as a cheap replacement for global illumination.
	float ambientStrength = 0.0005;
	totalLight += (ambientStrength * material.color).xyz;

	// HDR tone mapping
	totalLight = totalLight / (totalLight + float3(1.0));

	// Gamma correction
	float alpha = material.color.w;
	return float4(pow(totalLight, float3(1.0/2.2)), alpha);
}
