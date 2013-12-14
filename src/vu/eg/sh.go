// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package main

import (
	"fmt"
	"vu/device"
	"vu/render/gl"
)

// sh is used to test and showcase the vu/device package.  Just getting a window
// to appear demonstrates that the majority of the functionality is working.
// The remainder of the example dumps keyboard and mouse events showing that
// user input is being processed.
func sh() {
	sh := &shtag{}
	dev := device.New("Shell", 400, 100, 800, 600)
	gl.Init()
	fmt.Printf("%s %s", gl.GetString(gl.RENDERER), gl.GetString(gl.VERSION))
	fmt.Printf(" GLSL %s\n", gl.GetString(gl.SHADING_LANGUAGE_VERSION))
	dev.Open()
	gl.ClearColor(0.3, 0.6, 0.4, 1.0)
	for dev.IsAlive() {
		sh.update(dev)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		dev.SwapBuffers()
	}
	dev.ShowCursor(true)
	dev.Dispose()
}

// Globally unique "tag" for this example.
type shtag struct{}

// Show all the concurrent user actions.
func (sh *shtag) update(dev device.Device) {
	pressed := dev.Update()
	if len(pressed.Down) > 0 {
		fmt.Print(pressed.Mx, ",", pressed.My, ":")
		fmt.Print("shf:", pressed.Shift, " ctl:", pressed.Control)
		if pressed.Resized {
			fmt.Print("resized:")
		}
		if pressed.Focus {
			fmt.Print("   focus:")
		} else {
			fmt.Print(" nofocus:")
		}
		for key, duration := range pressed.Down {
			fmt.Print(key, ",", duration, ":")
		}
		fmt.Println()
	}
}
