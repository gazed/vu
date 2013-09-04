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
// The remainder of the example dumps and keyboard and mouse events showing that
// user input is being processed.
func sh() {
	sh := &shtag{}
	dev := device.New("Shell", 400, 100, 800, 600)
	dev.SetFocuser(sh)
	gl.Init()
	fmt.Printf("%s %s", gl.GetString(gl.RENDERER), gl.GetString(gl.VERSION))
	fmt.Printf(" GLSL %s\n", gl.GetString(gl.SHADING_LANGUAGE_VERSION))
	dev.Open()
	gl.ClearColor(0.3, 0.6, 0.4, 1.0)
	for dev.IsAlive() {
		pressed, _, _ := dev.ReadAndDispatch()
		sh.React(pressed)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		dev.SwapBuffers()
	}
	dev.ShowCursor(true)
	dev.Dispose()
}

// Globally unique "tag" for this example.
type shtag struct{}

// Show all the concurrent user actions.
func (sh *shtag) React(pressed []string) {
	if len(pressed) > 0 {
		for _, p := range pressed {
			print(p, ":")
		}
		println()
	}
}

func (sh *shtag) Focus(focus bool) { println("window has focus", focus) }
