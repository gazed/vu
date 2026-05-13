// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

// -----------------------------------------------------------------------------
// image2D draws the given texture on a 2D model.

struct VIn {
	float2 position : POSITION0;
	float2 texcoord;
};

import "modTexture";

[shader("vertex")]
VOut vertMain(VIn i) {
	VOut o;
	o.texcoord = i.texcoord;
	o.position = mul(su.proj, mul(su.view, mul(mu.model, float4(i.position.xy, 0.0, 1.0))));
	return o;
}

[shader("fragment")]
float4 fragMain(VOut i) {
	return samplers[COLOR].Sample(i.texcoord);
}
