// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

module modTexture;

public struct SceneUniforms
{
	public float4x4 proj;
	public float4x4 view;
};
layout(set=0, binding=0) public ConstantBuffer<SceneUniforms> su;

public static const int COLOR = 0;
layout(set = 1, binding = 0) public uniform Sampler2D samplers[1];

public struct ModelUniforms
{
	public float4x4 model;
}
layout(push_constant) public ModelUniforms mu;

public struct VOut
{
	public float2 texcoord;
	public float4 position : SV_POSITION;
};
