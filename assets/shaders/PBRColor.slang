// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT
//
// A Physically Based Rendering shader using a solid color instead of textures.
// Requires one directional light.

struct VIn {
	float3 position : POSITION0; // float3 implies 3D render pass.
	float3 normal;				 // vertex normal.
};

import "modPBR";

// use DescriptorSet when exceeding 128 byte PushConstant size limit.
struct ModelUniforms
{
	float4x4 model;  // model to world transform
	float4	 color;  // rgba
	float4	 pbrMR;  // x:metallic y:roughness
};
layout(push_constant) ModelUniforms mu;

struct VOut
{
	float4 position : SV_POSITION;
	float3 normal;	// normalized vertex normal.
	float3 vworld;	// vertex world position.
};

[shader("vertex")]
VOut vertMain(VIn i) {
	VOut o;

	// calcuate unit normal in world space
	float3x3 mod3x3 = (float3x3)mu.model;
	o.normal = normalize(mul(mod3x3, i.normal));

	// calculate vertex world space position
	o.vworld = mul(mu.model, float4(i.position, 1.0)).xyz;
	o.position = mul(su.proj, mul(su.view, mul(mu.model, float4(i.position, 1.0))));
	return o;
}

[shader("fragment")]
float4 fragMain(VOut o) {
	FragIn f;
	f.normal = o.normal;
	f.vworld = o.vworld;

	PBRMat material;
	material.color = mu.color;					  // rgba
	material.metallic = float(round(mu.pbrMR.x)); // 0.0 or 1.0;
	material.roughness = mu.pbrMR.y;			  // 0.0 to 1.0
	return PBRFrag(f, material);
}
