// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/gazed/vu/device"
	"github.com/gazed/vu/render/gl"
)

// Initialize enough of the opengl context that some OpenGL information
// can be dumped to screen along with the bindings. This is a basic graphics
// package test that checks if the underlying OpenGL functions are available.
// Columns of function names marked [+]:available or [ ]:missing will
// be written the the console.
//
// CONTROLS: NA
func dg() {
	app := device.New("Dump", 400, 100, 600, 600)
	gl.Dump() // gets graphic context to properly bind.
	fmt.Printf("%s %s", gl.GetString(gl.RENDERER), gl.GetString(gl.VERSION))
	fmt.Printf(" GLSL %s\n", gl.GetString(gl.SHADING_LANGUAGE_VERSION))
	app.Dispose()
}
