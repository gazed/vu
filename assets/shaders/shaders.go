// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

// Package shaders used in the engine examples.
// Run "go build; ./shaders" to create .spv byte code files and .shd reflection data.
//
// PBR shaders are based on the youtube tutorial43 from:
//   - https://github.com/emeiri/ogldev/
package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// slangc is kindof slow so only recompile shaders that have changed.
// Expected to be run from the asset/shader directory.
func main() {
	entries, err := os.ReadDir(".")
	if err != nil {
		slog.Error("failed to read directory", "error", err)
		return
	}

	// find the "slang" shaders and the corresponding generated files.
	shaderNames := []string{}
	slangs := map[string]os.FileInfo{} // shaders
	shds := map[string]os.FileInfo{}   // generated reflection data.
	verts := map[string]os.FileInfo{}  // generated vertex shader.
	frags := map[string]os.FileInfo{}  // generated fragment shader.
	for _, entry := range entries {
		if name, ok := strings.CutSuffix(entry.Name(), ".slang"); ok {
			slangs[name], _ = entry.Info()
			shaderNames = append(shaderNames, name)
		}
		if name, ok := strings.CutSuffix(entry.Name(), ".shd"); ok {
			shds[name], _ = entry.Info()
		}
		if name, ok := strings.CutSuffix(entry.Name(), "_vert.spv"); ok {
			verts[name], _ = entry.Info()
		}
		if name, ok := strings.CutSuffix(entry.Name(), "_frag.spv"); ok {
			frags[name], _ = entry.Info()
		}
	}

	// compile shaders alpabetically (case insensitive)
	sort.Slice(shaderNames, func(i, j int) bool {
		return strings.ToLower(shaderNames[i]) < strings.ToLower(shaderNames[j])
	})

	// recompile shaders if any generated output is missing or if
	// the slang file is newer than any generated output.
	for _, name := range shaderNames {
		slangInfo := slangs[name]
		shdInfo, okShd := shds[name]
		_, okVert := verts[name]
		_, okFrag := frags[name]
		switch {
		case strings.HasSuffix(name, "Mod"): // modules are not for standalone compiles.
			// FUTURE: check for changes to shader modules and
			//         recompile the shaders that import them.
			continue
		case !okShd || !okVert || !okFrag:
			compileShader(name, slangInfo, "missing generated file")
		case slangInfo.ModTime().After(shdInfo.ModTime()):
			// the generated files exist, so they should all have the same time.
			compileShader(name, slangInfo, "slang file updated")
		default:
			fmt.Printf("%-20s ok\n", name)
		}
	}
}

// compile the shader. "slangc" is installed as part of the VulkanSDK.
func compileShader(name string, s os.FileInfo, reason string) {
	fmt.Printf("%-20s %s : compiling\n", name, reason)

	// currently 3 commands needed to generate shader output.
	cmds := []*exec.Cmd{
		// slangc billboard.slang -target spirv -o billboard_vert.spv -entry vertMain
		exec.Command("slangc", s.Name(), "-target", "spirv", "-o", fmt.Sprintf("%s_vert.spv", name), "-entry", "vertMain"),
		// slangc billboard.slang -target spirv -o billboard_frag.spv -entry fragMain
		exec.Command("slangc", s.Name(), "-target", "spirv", "-o", fmt.Sprintf("%s_frag.spv", name), "-entry", "fragMain"),
		//  sh -c "slangc billboard.slang -reflection-json billboard.shd -target spirv > /dev/null"
		exec.Command("slangc", s.Name(), "-reflection-json", name+".shd", "-target", "spirv"),
	}
	for _, cmd := range cmds {
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("error compiling shader\n%s\n", string(output))
			return // abort on any problem
		}
	}
}
