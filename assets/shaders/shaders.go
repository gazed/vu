// SPDX-FileCopyrightText : © 2022-2024 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

// Package shaders contains the builtin engine shaders.
// Run "go generate" to create the shader .spv byte code files.
//
// Shaders are linked directly to the engine code in a few ways:
//   - vu/load/shd.go expects specific names for attributes and
//     uniform types. The shd.go code may need to be updated for new shaders.
//   - The order of attributes and uniforms in the shader description
//     files (*.shd) matter and relate directly to the layout values
//     in the shader code.
//   - The shader push constants block only guarantees upto 128 bytes.
//
// PBR shaders are based on the youtube tutorial43 from:
//   - https://github.com/emeiri/ogldev/
package shaders

// =============================================================================
// run "go generate" to create or update the shader byte code.

// 3D shaders
//go:generate glslc bbinst.vert -o bbinst.vert.spv
//go:generate glslc bbinst.frag -o bbinst.frag.spv
//go:generate glslc bboard.vert -o bboard.vert.spv
//go:generate glslc bboard.frag -o bboard.frag.spv
//go:generate glslc circle.vert -o circle.vert.spv
//go:generate glslc circle.frag -o circle.frag.spv
//go:generate glslc col3D.vert -o col3D.vert.spv
//go:generate glslc col3D.frag -o col3D.frag.spv
//go:generate glslc lines.vert -o lines.vert.spv
//go:generate glslc lines.frag -o lines.frag.spv
//go:generate glslc PBRCol.vert -o PBRCol.vert.spv
//go:generate glslc PBRCol.frag -o PBRCol.frag.spv
//go:generate glslc PBRTex.vert -o PBRTex.vert.spv
//go:generate glslc PBRTex.frag -o PBRTex.frag.spv
//go:generate glslc tex3D.vert -o tex3D.vert.spv
//go:generate glslc tex3D.frag -o tex3D.frag.spv
//go:generate glslc sdf.vert -o sdf.vert.spv
//go:generate glslc sdf.frag -o sdf.frag.spv

// 2D shaders
//go:generate glslc col2D.vert -o col2D.vert.spv
//go:generate glslc col2D.frag -o col2D.frag.spv
//go:generate glslc icon.vert -o icon.vert.spv
//go:generate glslc icon.frag -o icon.frag.spv
//go:generate glslc label.vert -o label.vert.spv
//go:generate glslc label.frag -o label.frag.spv
//go:generate glslc lines2D.vert -o lines2D.vert.spv
//go:generate glslc lines2D.frag -o lines2D.frag.spv
