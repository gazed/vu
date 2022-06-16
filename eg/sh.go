// Copyright © 2013-2024 Galvanized Logic Inc.

package main

import (
	"fmt"
	"time"

	"github.com/gazed/vu/device"
)

// sh creates a shell to showcase the vu/device package. Just getting a window
// to appear demonstrates that the majority of the functionality is working.
// The remainder of the example dumps keyboard and mouse events showing that
// user input is being processed.
//
// CONTROLS:
//
//	key   : print out key press info
//	mouse : print out mouse click info
func sh() {
	windowed := true
	dev := device.New(windowed, "Shell", 100, 100, 800, 600)
	dev.CreateDisplay()
	for dev.IsRunning() {
		in := dev.GetInput()

		// quit if X is pressed
		if _, ok := in.Pressed[device.KX]; ok {
			dev.Dispose()
		}

		// show pressed keys
		if len(in.Pressed) > 0 {
			fmt.Print("pressed:")
			for k, _ := range in.Pressed {
				fmt.Print(" ", k)
			}
			fmt.Println()
		}

		// show keys held down.
		if len(in.Down) > 0 {
			fmt.Print("down:")
			for k, _ := range in.Down {
				fmt.Print(" ", k)
			}

			// show mouse position when holding down a mouse button.
			_, mouseLeftDown := in.Down[device.KML]
			_, mouseMiddleDown := in.Down[device.KMM]
			_, mouseRightDown := in.Down[device.KMR]
			if mouseLeftDown || mouseMiddleDown || mouseRightDown {
				fmt.Print(" mx,my:", in.Mx, ",", in.My)
			}
			fmt.Println()
		}

		// show released keys
		for k, v := range in.Released {
			fmt.Println("released:", k, " ", v.Milliseconds(), "ms")
		}

		// show scroll amounts
		if in.Scroll != 0 {
			fmt.Println("scroll:", in.Scroll)
		}

		// sleep a bit to simulate a game update/render loop.
		time.Sleep(50 * time.Millisecond)
	}
}

type shtag struct{} // Globally unique "tag" for this example.
