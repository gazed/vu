// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

// noise uses fractal brownian noise from modNoise to color a model.

struct VIn {
	float3 position : POSITION0; // float3 implies 3D render pass.
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
	float4	 f4;	 // noise parameters. x:scale y:frequency
};
layout(push_constant) ModelUniforms mu;

struct VOut
{
	float4 position : SV_POSITION;
	float3 vworld;	// vertex world position.
};

[shader("vertex")]
VOut vertMain(VIn i) {
	VOut o;
	o.vworld = mul(mu.model, float4(i.position, 1.0)).xyz;
	o.position = mul(su.proj, mul(su.view, mul(mu.model, float4(i.position, 1.0))));
	return o;
}

import "noise3DMod";

[shader("fragment")]
float4 fragMain(VOut i) {
	float scale = mu.f4.x;
	float frequency = mu.f4.y;
	float n = fbm3D(scale*i.vworld, frequency);
	if (n < 0.5) {
		return float4(0.1*n, 0.1*n, 1.0*n, n); // blue oceans
	}
	if (n < 0.53) {
		return float4(1.0*n, 1.0*n, 0.1*n, n); // yellow beaches
	}
	if (n < 0.7) {
		return float4(0.1*n, 1.0*n, 0.1*n, n); // green forest
	}
	return float4(n, n, n, n); // white mountains
}
