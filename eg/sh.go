// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"fmt"
	"time"

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

		// slow things down so that the loop is closer
		// to the engine update loop timing.
		time.Sleep(10 * time.Millisecond)
	}
	dev.ShowCursor(true)
	dev.Dispose()
}

// Globally unique "tag" for this example.
type shtag struct{}

// Show all the concurrent user actions.
func (sh *shtag) update(dev device.Device) {
	pressed := dev.Update()
	if pressed.Scroll != 0 {
		fmt.Printf("scroll %d\n", pressed.Scroll)
	}
	if pressed.Resized {
		_, _, ww, wh := dev.Size()
		fmt.Printf("resized %d %d\n", ww, wh)
	}
	if len(pressed.Down) > 0 {
		fmt.Print(pressed.Mx, ",", pressed.My, ":")
		if pressed.Resized {
			fmt.Print(" resized:")
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
}
