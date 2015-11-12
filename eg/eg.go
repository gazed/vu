// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.


// Package eg is used to test and demonstrate different aspects of the
// vu (virtual universe) engine. Examples are used both to showcase a
// particular 3D capability and to act as high level test cases for
// the engine. The examples are run using:
//      eg [example name]
// Invoking eg without parameters will list the examples that can be run.
// Please look at each examples source code for possible user actions.
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
		example{"dg", "dg: Dump OpenGL Bindings", dg},
		example{"da", "da: Dump OpenAL Bindings", da},
		example{"tr", "tr: Spinning Triangle", tr},
		example{"ld", "ld: Load .obj model", ld},
		example{"sh", "sh: Simple Shell", sh},
		example{"kc", "kc: Keyboard Controller", kc},
		example{"au", "au: Audio", au},
		example{"sf", "sf: Shader Fire", sf},
		example{"bb", "bb: Banners & Billboards", bb},
		example{"lt", "lt: Lighting", lt},
		example{"rl", "rl: Random Levels", rl},
		example{"sg", "sg: Scene Graph", sg},
		example{"cr", "cr: Collision Resolution", cr},
		example{"tm", "tm: Terrain Map", tm},
		example{"fm", "fm: Form Layout", fm},
		example{"ps", "ps: Particle System", ps},
		example{"rc", "rc: Ray Cast", rc},
		example{"ma", "ma: Model Animation", ma},
		example{"ff", "ff: Flow Field", ff},
		example{"rt", "rt: Ray Trace", rt},
		example{"tt", "tt: Render to Texture", tt},
		example{"sm", "sm: Shadow Map", sm},
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
