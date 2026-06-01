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

static const int COLOR = 0;
layout(set=1, binding=0) uniform Sampler2D samplers[1];

struct ModelUniforms
{
	float4x4 model;
	float3	 scale;
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

	// for row_major flip the multiplication order, ie:
	// row_major float4x4 mvm = mul(mu.model, su.view);
	// o.position = mul((float4(i.position*mu.scale, 1.0) + float4(mvm[3].xyz, 0.0)), su.proj);

	// works with default col_major
	float4x4 mvm = mul(su.view, mu.model);
	o.position = mul(su.proj, (float4(i.position*mu.scale, 1.0) + float4(mvm[0].w, mvm[1].w, mvm[2].w, 0.0)));
	return o;
}

[shader("fragment")]
float4 fragMain(VOut i) {
	return samplers[COLOR].Sample(i.texcoord);
}
