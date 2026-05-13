// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT
//
// A Physically Based Rendering shader using one texture.
// Requires one directional light.

struct VIn {
	float3 position : POSITION0; // float3 implies 3D render pass.
	float2 texcoord;			 // vertex texture coordinates.
	float3 normal;				 // vertex normal.
};

import "modPBR";

static const int COLOR = 0;
layout(set = 1, binding = 0) uniform Sampler2D samplers[1];

// use DescriptorSet when exceeding 128 byte PushConstant size limit.
struct ModelUniforms
{
	float4x4 model;  // model to world transform
	float4	 pbrMR;  // x:metallic y:roughness
};
layout(push_constant) ModelUniforms mu;

struct VOut
{
	float4 position : SV_POSITION;
	float3 normal;	 // normalized vertex normal.
	float3 vworld;	 // vertex world position.
	float2 texcoord;
};

[shader("vertex")]
VOut vertMain(VIn i) {
	VOut o;
	o.texcoord = i.texcoord;

	// calcuate unit normal in world space
	float3x3 mod3x3 = (float3x3)mu.model;
	o.normal = normalize(mul(mod3x3, i.normal));

	// calculate vertex world space position
	o.vworld = mul(mu.model, float4(i.position, 1.0)).xyz;
	o.position = mul(su.proj, mul(su.view, mul(mu.model, float4(i.position, 1.0))));
	return o;
	return o;
}

[shader("fragment")]
float4 fragMain(VOut o) {
	FragIn f;
	f.normal = o.normal;
	f.vworld = o.vworld;

	PBRMat material;
	material.color = samplers[COLOR].Sample(o.texcoord);
	material.metallic = float(round(mu.pbrMR.x)); // 0.0 or 1.0;
	material.roughness = mu.pbrMR.y;			  // 0.0 to 1.0
	return PBRFrag(f, material);
}
