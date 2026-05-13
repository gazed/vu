// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT
//
module modColor;

public struct SceneUniforms
{
	public float4x4 proj;
	public float4x4 view;
};
layout(set=0, binding=0) public ConstantBuffer<SceneUniforms> su;

public struct ModelUniforms
{
	public float4x4 model;
	public float4	color;
};
layout(push_constant) public ModelUniforms mu;

public struct VOut
{
	public float4 position : SV_POSITION;
};
