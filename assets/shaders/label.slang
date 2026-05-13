// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

// label colors the visible pixels of the given image.
// Expected to be used for colored text.

struct VIn {
	float2 position : POSITION0;
	float2 texcoord;
};

struct SceneUniforms
{
	float4x4 proj;
	float4x4 view;
};
layout(set=0, binding=0) ConstantBuffer<SceneUniforms> su;

static const int COLOR = 0;
layout(set = 1, binding = 0) uniform Sampler2D samplers[1];

struct ModelUniforms
{
	float4x4 model;
	float4	 color;
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
	o.position = mul(su.proj, mul(su.view, mul(mu.model, float4(i.position.xy, 0.0, 1.0))));
	return o;
}

[shader("fragment")]
float4 fragMain(VOut i) {
	float alpha = samplers[COLOR].Sample(i.texcoord).a;
	return float4(mu.color.x, mu.color.y, mu.color.z, alpha);
}
