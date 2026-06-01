// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT
//
// -----------------------------------------------------------------------------
// 2D colored lines
// -----------------------------------------------------------------------------

// shader meta data, attached to vertex main, is passed through reflection.
[__AttributeUsage(_AttributeTargets.Function)]
struct LineAttribute { bool isLine; }; // line shader.

struct VIn {
	float3 position : POSITION0; // float3 implies 3D render pass.
};

import "colorMod";

[shader("vertex")] [Line(true)]
VOut vertMain(VIn i) {
	VOut o;
	o.position = mul(su.proj, mul(su.view, mul(mu.model, float4(i.position, 1.0))));
	return o;
}

[shader("fragment")]
float4 fragMain(VOut i) {
	return mu.color;
}
