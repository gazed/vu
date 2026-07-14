// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT
//
// A Physically Based Rendering shader using one base color texture.
// Requires at least one directional light.

struct VIn {
	float3 position : POSITION0; // float3 implies 3D render pass.
	float2 texcoord;			 // vertex texture coordinates.
	float3 normal;				 // vertex normal.
};

import "lightMod"; // for Light struct
import "pbrMod";   // for PBRInfo struct and PBRColor func.

struct SceneUniforms
{
	float4x4 proj;		// projection matrix
	float4x4 view;		// view matrix
	float4	 cam;		// local camera location
	Light	 lights[4]; // max 4 lights.
};
layout(set=0, binding=0) ConstantBuffer<SceneUniforms> su;

static const int COLOR = 0;
layout(set=1, binding=0) uniform Sampler2D samplers[1];

// model data as push constants limited to 128 bytes.
struct ModelUniforms
{
	float4x4 model;  // model to world transform
	float4	 pbrMR;  // x:metallic y:roughness
};
layout(push_constant) ModelUniforms mu;

struct VOut
{
	float4 position : SV_POSITION;
	float2 texcoord;
	float3 world;  // vertex world position.
	float3 norm;   // normalized vertex normal
	float3 view;   // normalized view vector.
};

[shader("vertex")]
VOut vertMain(VIn i) {
	VOut o;
	o.position = mul(su.proj, mul(su.view, mul(mu.model, float4(i.position, 1.0))));
	o.texcoord = i.texcoord;

	// calculate lighting information once per vertex.
	o.world = mul(mu.model, float4(i.position, 1.0)).xyz;
	o.norm	= normalize(mul((float3x3)mu.model, i.normal));
	o.view	= normalize(su.cam.xyz - o.world);
	return o;
}

[shader("fragment")]
float4 fragMain(VOut i) {
	PBRInfo info;
	info.Color = samplers[COLOR].Sample(i.texcoord);
	info.World = i.world;
	info.N = normalize(i.norm);
	info.V = normalize(i.view);
	info.Metallic = float(round(mu.pbrMR.x)); // 0.0 or 1.0
	info.Roughness = mu.pbrMR.y;			  // 0.0 to 1.0

	// calculate the affect of lights on the final color
	float3 totalLight = float3(0.0);
	LightInfo li;
	li.Color = info.Color;
	li.World = info.World;
	for (int i = 0; i < 4; i++) {
		if (updateLightInfo(su.lights[i], li)) {
			info.L = li.L;
			info.I = li.I;
			totalLight += PBRColor(info);
		}
	}
	return FinalPBRColor(totalLight, info.Color);
}
