// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

struct VIn {
	// vertex attributes
	float3 position : POSITION0;
	float2 texcoord;

	// instance attributes
	float3 i_position; // model world position
	float3 i_color;    // color to add to the base texture
	float  i_scale;    // model scale factor
};

struct SceneUniforms
{
	row_major float4x4 proj;
	row_major float4x4 view;
};
layout(set=0, binding=0) ConstantBuffer<SceneUniforms> su;

static const int COLOR = 0;
layout(set = 1, binding = 0) uniform Sampler2D samplers[1];

struct VOut
{
	float4 position : SV_POSITION;
	float2 texcoord;
	float3 color;
};

[shader("vertex")]
VOut vertMain(VIn i) {
	VOut o;
	o.color = i.i_color;
	o.texcoord = i.texcoord;

	// billboard calculation to have image face the camera.
	// converted from glsl so works with row-major input matricies ... reversed mult order.
	row_major float4x4 mm = float4x4(float4(1.0, 0.0, 0.0, 0.0),
									 float4(0.0, 1.0, 0.0, 0.0),
									 float4(0.0, 0.0, 1.0, 0.0),
									 float4(i.i_position, 1.0));
	row_major float4x4 mvm = mul(mm, su.view);
	o.position = mul((float4(i.position*i.i_scale, 1.0) + float4(mvm[3].xyz, 0.0)), su.proj);
	return o;
}

[shader("fragment")]
float4 fragMain(VOut i) {
	float4 color = samplers[COLOR].Sample(i.texcoord);

	// modify the color, keeping the base_color alpha.
	float3 colorized = clamp(color.xyz * i.color, 0.0, 1.0);
	return float4(colorized, color.w);
}
