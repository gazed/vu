// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

// col2D applies a color to the 2D model

struct VIn {
	float2 position : POSITION0; // float2 implies 2D render pass.
};

import "modColor";

[shader("vertex")]
VOut vertMain(VIn i) {
	VOut o;
	o.position = mul(su.proj, mul(su.view, mul(mu.model, float4(i.position.xy, 0.0, 1.0))));
	return o;
}

[shader("fragment")]
float4 fragMain(VOut i) {
	return mu.color;
}
