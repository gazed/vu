// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

// Package shaders used in the engine examples.
// Build and run this file to generate shaders.
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
	content := map[string][]byte{}     // slang content to check for imports
	shds := map[string]os.FileInfo{}   // generated reflection data.
	verts := map[string]os.FileInfo{}  // generated vertex shader.
	frags := map[string]os.FileInfo{}  // generated fragment shader.
	for _, entry := range entries {
		if name, ok := strings.CutSuffix(entry.Name(), ".slang"); ok {
			slangs[name], _ = entry.Info()
			content[name], _ = os.ReadFile(entry.Name())
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

	// check which shaders require recompile, ie:
	// - if any of the module imports are are newer than the shader.
	// - if any generated output is missing
	// - if the slang file is newer than any generated output.
	requiresCompile := map[string]string{}
	for _, name := range shaderNames {
		requiresCompile[name] = "ok"
	}
	for _, name := range shaderNames {
		slangInfo := slangs[name]
		shdInfo, okShd := shds[name]
		_, okVert := verts[name]
		_, okFrag := frags[name]
		switch {
		case strings.HasSuffix(name, "Mod"):
			// modules are not compiled, but the shaders that
			// import them need to be checked.
			for sname, sbytes := range content {
				if strings.Contains(string(sbytes), fmt.Sprintf("import \"%s\";", name)) {
					sinfo := shds[sname]
					if slangInfo.ModTime().After(sinfo.ModTime()) {
						requiresCompile[sname] = "import " + name + " updated"
					}
				}
			}
		case !okShd || !okVert || !okFrag:
			requiresCompile[name] = "missing generated file"
		case slangInfo.ModTime().After(shdInfo.ModTime()):
			// the generated files exist, so they should all have the same time.
			requiresCompile[name] = "slang updated"
		}
	}

	// recompile shaders as needed in alphabetical order.
	for _, name := range shaderNames {
		reason := requiresCompile[name]
		switch {
		case reason == "ok":
			fmt.Printf("%-20s ok\n", name)
		default:
			compileShader(name, slangs[name], reason)
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
