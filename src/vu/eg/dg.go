// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package main

import (
	"fmt"
	"vu/device"
	"vu/render/gl"
)

// Initialize enough of the opengl context that some OpenGL information
// can be dumped to screen along with the bindings.  This is a basic graphics
// package test that checks if the underlying OpenGL functions are available.
// Columns of function names with more [+] (available) than [-] (missing) signs
// will be written the the console if everything is ok.
func dg() {
	gl.Dump() // doesn't need context.

	// Also print the opengl version.
	app := device.New("Dump", 400, 100, 600, 600)
	fmt.Printf("%s %s", gl.GetString(gl.RENDERER), gl.GetString(gl.VERSION))
	fmt.Printf(" GLSL %s\n", gl.GetString(gl.SHADING_LANGUAGE_VERSION))
	app.Dispose()
}
