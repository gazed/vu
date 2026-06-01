// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

module pbrMod;

// information needed to calculate a PBR color.
public struct PBRInfo
{
	public float4 Color;     // fragment rgba
	public float3 World;     // vertex world position.
	public float3 N;         // normalized vertex normal
	public float3 L;         // normalized light direction
	public float3 V;         // normalized view direction
	public float3 I;         // light color * intensity
	public float  Metallic;  // 0.0 or 1.0
	public float  Roughness; // 0.0 to 1.0
};

static const float PI = 3.1415926535897932384626433832795;

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

// PBRColor returns an RGB value for the given light and
// material information.
public float3 PBRColor(PBRInfo pi) {
	float3 H = normalize(pi.V + pi.L);       // half vector
	float nDotH = max(dot(pi.N, H), 0.0);
	float vDotH = max(dot(pi.V, H), 0.0);
	float nDotL = max(dot(pi.N, pi.L), 0.0);
	float nDotV = max(dot(pi.N, pi.V), 0.0);

	// conserve energy so refaction+reflection==1.0
	float3 F = schlickFresnel(vDotH, pi.Metallic, pi.Color.xyz);
	float3 kS = F;			// specular: reflection
	float3 kD = 1.0 - kS;	// diffuse : refraction

	// specular BRDF.
	float3 SpecBRDF_nom  = ggxDistribution(nDotH, pi.Roughness) *
						   F *
						   geomSmith(nDotL, pi.Roughness) *
						   geomSmith(nDotV, pi.Roughness);
	float SpecBRDF_denom = 4.0 * nDotV * nDotL + 0.0001;
	float3 SpecBRDF = SpecBRDF_nom / SpecBRDF_denom;

	// use color for non-metals, metals will already have the color in kS.
	float3 fLambert = lerp(pi.Color.xyz, float3(0.0), pi.Metallic);
	float3 DiffuseBRDF = kD * fLambert / PI;

	// final color value for the given light.
	float alpha = pi.Color.w;
	return (DiffuseBRDF*alpha + SpecBRDF) * pi.I * nDotL;
}

// FinalPBRColor adjusts the color after all the lights have been applied.
public float4 FinalPBRColor(float3 litColor, float4 baseColor) {

	// add fixed (indirect) ambient value as a cheap replacement for global illumination.
	float ambientStrength = 0.0005;
	litColor += (ambientStrength * baseColor).xyz;

	// HDR tone mapping
	litColor = litColor / (litColor + float3(1.0));

	// Gamma correction
	float alpha = baseColor.w;
	return float4(pow(litColor, float3(1.0/2.2)), alpha);
}
