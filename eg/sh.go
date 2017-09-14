// Copyright © 2013-2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/gazed/vu/device"
	"github.com/gazed/vu/render/gl"
)

// sh is used to test and showcase the vu/device package. Just getting a window
// to appear demonstrates that the majority of the functionality is working.
// The remainder of the example dumps keyboard and mouse events showing that
// user input is being processed. See vu/device package for more information.
//
// CONTROLS:
//   key   : print out key press info
//   mouse : print out mouse click info
func sh() {
	device.Run(&shtag{}) // does not return from here...
}

type shtag struct{} // Globally unique "tag" for this example.

// Init is a one-time callback before rendering updates.
func (sh *shtag) Init(dev device.Device) {
	gl.Init()
	dev.SetTitle("Shell")
	dev.SetSize(400, 100, 800, 600)
	gl.Viewport(0, 0, 800, 600)
	fmt.Printf("%s %s", gl.GetString(gl.RENDERER), gl.GetString(gl.VERSION))
	fmt.Printf(" GLSL %s\n", gl.GetString(gl.SHADING_LANGUAGE_VERSION))
	gl.ClearColor(0.3, 0.6, 0.4, 1.0)
}

// Refresh is called each update tick. In this case it
// prints Show all the concurrent user actions to the console.
func (sh *shtag) Refresh(dev device.Device) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	pressed := dev.Down()
	if pressed.Scroll != 0 {
		fmt.Printf("scroll %d\n", pressed.Scroll)
	}
	if pressed.Resized {
		x, y, ww, wh := dev.Size()
		fmt.Printf("resized %d %d %d %d\n", x, y, ww, wh)
		gl.Viewport(0, 0, int32(ww), int32(wh))
	}
	if len(pressed.Down) > 0 {
		if pressed.Focus {
			fmt.Print("   focus:")
		} else {
			fmt.Print(" nofocus:")
		}
		fmt.Print(pressed.Mx, ",", pressed.My, ":")
		for key, downTicks := range pressed.Down {
			fmt.Print(key, ",", downTicks, ":")
		}
		fmt.Println()

		// demo clipboard copy/paste.
		if down, ok := pressed.Down[device.KC]; ok && down == 1 {
			fmt.Printf("\"%s\" ", dev.Copy())
		}
		if down, ok := pressed.Down[device.KP]; ok && down == 1 {
			dev.Paste("Sample clipboard text")
		}

		// toggle windowed mode if W is pressed.
		if down, ok := pressed.Down[device.KW]; ok && down == 1 {
			dev.ToggleFullScreen()
		}

		// quit if X is pressed
		if down, ok := pressed.Down[device.KX]; ok && down == 1 {
			dev.Dispose()
		}
	}
	dev.SwapBuffers()
}
