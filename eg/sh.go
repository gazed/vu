// Copyright © 2025 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/gazed/vu/device"
)

// sh is used to test and showcase the vu/device package. Just getting a window
// to appear demonstrates that the majority of the functionality is working.
// The remainder of the example dumps keyboard and mouse events showing that
// user input is being processed. See vu/device package for more information.
//
// CONTROLS:
//
//	key   : print out key press info
//	mouse : print out mouse click info
func sh() {
	dev := device.New(false, "eg:sh", 400, 100, 800, 600)
	dev.CreateDisplay()
	device.SetInputHandler(sh_inputHandler)

	// does not return while running
	// Used on platforms that require a render callback loop.
	dev.Run(sh_renderer) // does not return while example is running.
}

type shtag struct{} // Globally unique "tag" for this example.

func sh_renderer() {} // no rendering.

// device level input handler just for this example.
// Normally expected to use input events from eng.Update callback
func sh_inputHandler(event, data int64) {
	switch event {
	case device.EVENT_KEYUP:
	case device.EVENT_KEYDOWN:
		switch data {
		case device.KT:
			fmt.Printf("T key %d\n", data)
		case device.KML:
			fmt.Printf("left mouse click %d\n", data)
		default:
			fmt.Printf("press/click %d\n", data)
		}
	case device.EVENT_SCROLL:
		fmt.Printf("scroll %d\n", data)
	case device.EVENT_MODIFIER:
		fmt.Printf("modifier %d\n", data)
	case device.EVENT_MOVED:
		fmt.Printf("moved %d\n", data)
	case device.EVENT_RESIZED:
		fmt.Printf("resized %d\n", data)
	case device.EVENT_FOCUS_GAINED:
		fmt.Printf("focus gained %d\n", data)
	case device.EVENT_FOCUS_LOST:
		fmt.Printf("focus lost %d\n", data)
	default:
		fmt.Printf("unexpected event type %d\n", data)
	}
}
