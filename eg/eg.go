// Copyright © 2013-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// Package eg is used to test and demonstrate different aspects of the
// vu (virtual universe) engine. Examples are used both to showcase a
// particular 3D capability and to act as high level test cases for
// the engine. The examples are run using:
//      eg [example name]
// Invoking eg without parameters will list the examples that can be run.
// Please look at each examples source code for possible user actions like
// moving around or tabbing to show different scenes.
//
// The package subdirectories contain asset filess needed by the examples.
package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

// example combines example code with descriptions.
type example struct {
	tag         string // Example identifier.
	description string // Short description of the example.
	function    func() // Function to run the example.
}

// Launch the requested example or list available examples.
// Examples are roughly ordered from simple/basic at the top of the list
// to more complex/interesting at the bottom of the list.
func main() {
	examples := []example{
		{"dg", "Dump OpenGL Bindings", dg},
		{"da", "Dump OpenAL Bindings", da},
		{"tr", "Spinning Triangle", tr},
		{"ld", "Load .obj model", ld},
		{"sh", "Simple Shell", sh},
		{"sz", "Screen Resize", sz},
		{"kc", "Keyboard Controller", kc},
		{"au", "Audio", au},
		{"sf", "Shader Fire", sf},
		{"bb", "Banners & Billboards", bb},
		{"lt", "Lighting", lt},
		{"rl", "Random Levels", rl},
		{"hx", "Hex Grid", hx},
		{"sg", "Scene Graph", sg},
		{"cr", "Collision Resolution", cr},
		{"tm", "Terrain Map", tm},
		{"ps", "Particle System", ps},
		{"rc", "Ray Cast", rc},
		{"sd", "Sky Dome", sd},
		{"ma", "Model Animation", ma},
		{"ff", "Flow Field", ff},
		{"rt", "Ray Trace", rt},
		{"tt", "Render to Texture", tt},
		{"sm", "Shadow Map", sm},
		{"ss", "Super Shapes", ss},
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

// catchErrors is a utility method used by some examples to trap and
// dump runtime errors.
func catchErrors() {
	if r := recover(); r != nil {
		log.Printf("Panic %s: %s Shutting down.", r, debug.Stack())
	}
}
