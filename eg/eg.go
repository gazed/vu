// Copyright © 2013-2016 Galvanized Logic Inc.
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
// The package subdirectories contain resource data needed by the examples.
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
		{"dg", "dg: Dump OpenGL Bindings", dg},
		{"da", "da: Dump OpenAL Bindings", da},
		{"tr", "tr: Spinning Triangle", tr},
		{"ld", "ld: Load .obj model", ld},
		{"sh", "sh: Simple Shell", sh},
		{"kc", "kc: Keyboard Controller", kc},
		{"au", "au: Audio", au},
		{"sf", "sf: Shader Fire", sf},
		{"bb", "bb: Banners & Billboards", bb},
		{"lt", "lt: Lighting", lt},
		{"rl", "rl: Random Levels", rl},
		{"hx", "hx: Hex Grid", hx},
		{"sg", "sg: Scene Graph", sg},
		{"cr", "cr: Collision Resolution", cr},
		{"tm", "tm: Terrain Map", tm},
		{"fm", "fm: Form Layout", fm},
		{"ps", "ps: Particle System", ps},
		{"rc", "rc: Ray Cast", rc},
		{"ma", "ma: Model Animation", ma},
		{"ff", "ff: Flow Field", ff},
		{"rt", "rt: Ray Trace", rt},
		{"tt", "tt: Render to Texture", tt},
		{"sm", "sm: Shadow Map", sm},
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
		fmt.Printf("   %s \n", example.description)
	}
}

// catchErrors is a utility method used by some examples to trap and
// dump runtime errors.
func catchErrors() {
	if r := recover(); r != nil {
		log.Printf("Panic %s: %s Shutting down.", r, debug.Stack())
	}
}
