// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"github.com/gazed/vu/device"
)

// Input is used to communicate current user input to the application.
// This gives the current cursor location, current pressed keys,
// mouse buttons, and modifiers. These are sent to the application
// using the Director.Update() callback.
//
// The map of keys and mouse buttons that are currently pressed also
// include how long they have been pressed in update ticks. A negative
// value indicates a release. The total down duration can then be
// calculated by down duration less RELEASED timestamp.
type Input struct {
	Mx, My  int            // Current mouse location.
	Down    map[string]int // Keys, buttons with down duration ticks.
	Focus   bool           // True if window is in focus.
	Resized bool           // True if window was resized or moved.
	Scroll  int            // Scroll amount, if any.
	Dt      float64        // Delta time used for updates.
	Gt      float64        // Game time: total number of update ticks.
}

// convertInput copies the given device.Pressed input into vu.Input.
// Its also adds the delta time and updates the current game time
// in update ticks (it is expected to be called each update).
func (in *Input) convertInput(pressed *device.Pressed, dt float64) {
	in.Mx, in.My = pressed.Mx, pressed.My
	in.Focus = pressed.Focus
	in.Resized = pressed.Resized
	in.Scroll = pressed.Scroll
	in.Dt = dt
	in.Gt += 1

	// create a key/mouse down map that the application can trash.
	// It is expected to be cleared and refilled each update.
	for key, _ := range in.Down {
		delete(in.Down, key)
	}
	for key, val := range pressed.Down {
		in.Down[key] = val
	}
}
