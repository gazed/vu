// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

// col3D applies a color to the 3D model

struct VIn {
	float3 position : POSITION0; // float3 implies 3D render pass.
};

import "modColor";

[shader("vertex")]
VOut vertMain(VIn i) {
	VOut o;
	o.position = mul(su.proj, mul(su.view, mul(mu.model, float4(i.position, 1.0))));
	return o;
}

[shader("fragment")]
float4 fragMain(VOut i) {
	return mu.color;
}
