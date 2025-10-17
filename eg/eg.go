// Copyright © 2013-2024 Galvanized Logic Inc.

// Package eg is used to test and demo different aspects of the vu engine.
// The examples are run using:
//
//	eg [example name]
//
// Invoking eg without parameters will list the examples that can be run.
// Please look at each examples source code for possible user actions like
// moving around or tabbing to show different scenes.
// Assets are located in vu/assets.
//
// One time setup to compile the shader code:
//   - cd vu/assets/shaders
//   - go generate
package main

import (
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/gazed/vu/load"
)

// configLogging keeps the default logging settings.
// This runs before init() and is overridden in eg_debug.go by debug builds.
var configLogging func() = func() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn})))
}

// Launch the requested example or list available examples.
// Examples are roughly ordered from simple/basic at the top of the list
// to more complex/interesting at the bottom of the list.
func main() {
	configLogging()

	// change the default asset location to work for the examples.
	load.SetAssetDir(".spv", "../assets/shaders")
	load.SetAssetDir(".shd", "../assets/shaders")
	load.SetAssetDir(".png", "../assets/images")
	load.SetAssetDir(".wav", "../assets/audio")
	load.SetAssetDir(".glb", "../assets/models")
	load.SetAssetDir(".ttf", "../assets/fonts")
	load.SetAssetDir(".yaml", "../assets/data")

	examples := []example{
		{"au", "Audio", au},                // test audio bindings.
		{"ax", "Engine Audio", ax},         // test engine audio.
		{"vv", "Vulkan Version", vv},       // test vulkan availability and version
		{"kc", "Keyboard Controller", kc},  // test 2D render
		{"dr", "Display Ratio", dr},        // test window resizing
		{"mh", "Monkey Heads", mh},         // test 3D render with GLB files and PBR shaders.
		{"cr", "Collision Resolution", cr}, // test physics package.
		{"ps", "Primitive Shapes", ps},     // test drawing primitive shapes.
		{"is", "Instanced Stars", is},      // test drawing instanced models

		// FUTURE: {"ma", "Model Animation", ma}, // test animated models
	}

	// run the first matching example.
	for _, arg := range os.Args {
		for _, eg := range examples {
			if arg == eg.tag {
				eg.function()
				os.Exit(0)
			}
		}
	}

	// print usage if nothing was run.
	fmt.Printf("Usage: eg [example]\n")
	fmt.Printf("Examples are:\n")
	for _, example := range examples {
		fmt.Printf("   %s: %s \n", example.tag, example.description)
	}
}

// example combines example code with descriptions.
type example struct {
	tag         string // Example identifier.
	description string // Short description of the example.
	function    func() // Function to run the example.
}

// catchErrors is a utility method used by some examples to trap and
// dump runtime errors.
func catchErrors() {
	if r := recover(); r != nil {
		slog.Error("error", "msg", r)
		debug.PrintStack()
	}
}
