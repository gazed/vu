// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

// Package eg is used to test and demonstrate different aspects of the
// vu (virtual universe) engine.  Thus examples are used to both showcase a
// particular 3D capability and to act as high level test cases for the engine.
// The examples are run using:
//      eg [example name]
// Invoking eg without parameters will list the examples that can be
// run.  Please have a look at the example source code for more information.
//
// Note that the subdirectories contain resource data for the examples.
package main

import (
	"fmt"
	"os"
)

// example combines example code with descriptions.
type example struct {
	tag         string // identifier
	description string // short description of the example
	function    func() // function to run the example
}

// Launch the requested example or list available examples.
// Examples are roughly ordered from simplist at the top of the list to
// more complex/interesting at the bottom of the list.
func main() {
	examples := []example{
		example{"dg", "dg: Dump OpenGL Bindings", dg},
		example{"da", "da: Dump OpenAL Bindings", da},
		example{"sh", "sh: Simple Shell", sh},
		example{"tr", "tr: Spinning Triangle", tr},
		example{"ld", "ld: Load .obj model", ld},
		example{"co", "co: Collision & Motion", co},
		example{"mp", "mp: Mouse Picking", mp},
		example{"rl", "rl: Random Levels ", rl},
		example{"tb", "tb: Texture - Basic", tb},
		example{"bb", "bb: Billborded Texture", bb},
		example{"sg", "sg: Scene Graph", sg},
		example{"au", "au: Audio", au},
		example{"sf", "sf: Shader Fire", sf},
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
