// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

struct VIn {
	float3 position : POSITION0;
	float2 texcoord;
};

struct SceneUniforms
{
	float4x4 proj;
	float4x4 view;
};
layout(set=0, binding=0) ConstantBuffer<SceneUniforms> su;

struct ModelUniforms
{
	float4x4 model;
}
layout(push_constant) ModelUniforms mu;

struct VOut
{
	float2 texcoord;
	float4 position : SV_POSITION;
};

[shader("vertex")]
VOut vertMain(VIn i) {
	VOut o;
	o.texcoord = i.texcoord;
	o.position = mul(su.proj, mul(su.view, mul(mu.model, float4(i.position, 1.0))));
	return o;
}

[shader("fragment")]
float4 fragMain(VOut i) {
	// change texcoord 0:1 to uv -1.0:1.0
	float2 uv = i.texcoord * 2.0 -1;

	// line parameters
	float thickness = 0.01;
	float fade = 0.002;
	float dist = 1.0-length(uv);

	// circle outer edge
	float3 color = float3(smoothstep(0.0, fade, dist));

	// circle inner edge
	color *= float3(smoothstep(thickness, thickness-fade, dist));

	// ignore black pixels completely.
	if ((color.x == 0) && (color.y == 0) && (color.z == 0))
		discard;

	// final color
	return float4(color, 1.0);
}
