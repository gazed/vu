// Copyright © 2013-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"github.com/gazed/vu/device"
	"github.com/gazed/vu/render/gl"
)

// Basic graphics package test that checks if the underlying OpenGL functions
// are available. Columns of function names marked [+]:available or [ ]:missing
// are written the the console.
//
// Windows does not find the opengl functions without a graphics context.
//
// CONTROLS: NA
func dg() {
	device.Run(&dgtag{}) // Does not return.

	// on macOS could just run gl.Dump() instead of having to
	// create the device.
}

type dgtag struct{}

// Init is a one-time callback before rendering updates.
func (tag *dgtag) Init(dev device.Device) {
	gl.Dump()
	dev.Dispose()
}

// Refresh application state and render a new frame.
func (tag *dgtag) Refresh(dev device.Device) {}
