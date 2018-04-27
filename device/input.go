// Copyright © 2013-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package device

// Design Notes:
// Events are collected until polled by getPressed.
// The original intent was to collect and process the OS event queues
// concurrently  but, OSX only allows event processng from the main thread.

// input is used to process a user event stream into the Pressed structure
// that is returned to the master application each update.
type input struct {
	curr  *Pressed    // Consolidates current user events into state.
	down  *Pressed    // Clone of curr that is shared with the application.
	queue map[int]int // queue very quick press/release.
}

// newInput creates the memory needed to process user input events.
func newInput() *input {
	i := &input{}
	i.curr = &Pressed{Focus: true, Down: map[int]int{}}
	i.down = &Pressed{Focus: true, Down: map[int]int{}}
	i.queue = map[int]int{}
	return i
}

// getPressed returns the current user input state of buttons pressed
// and keys clicked. It expects to be called each update thread tick.
// Old released keys are removed for the next call.
func (i *input) getPressed(mx, my int) *Pressed {
	i.curr.Mx = mx          // update Current mouse location.
	i.curr.My = my          //   "
	i.updateDurations()     // Keep track of key press duration.
	i.clone(i.curr, i.down) // duplicate in case user modifies map.
	i.clearReleases(i.curr) // No longer relevant once sent to App.
	i.curr.Scroll = 0       // Reset any scrolling.
	i.curr.Resized = false  // Reset any resize indication.
	return i.down           // Return duplicated information.
}

// updateDurations tracks how long keys have been pressed for.
// Expected to be called each update. Ignore released keys.
func (i *input) updateDurations() {
	for key, val := range i.curr.Down {
		if val >= 0 {
			i.curr.Down[key] = val + 1
		}
	}
}

// recordPress tracks new key or mouse down user input events.
// Ignore any key presses unless the window has focus.
func (i *input) recordPress(code int) {
	if code >= 0 && i.curr.Focus {
		if _, ok := i.curr.Down[code]; !ok {
			i.curr.Down[code] = 0
		}
	}
}

// recordRelease tracks key or mouse up user input events.
func (i *input) recordRelease(code int) {
	if dwn, ok := i.curr.Down[code]; ok {
		if dwn == 0 {
			// OSX mouse pad can generated mouse click and release
			// in a single update. Queue up conflicting events.
			// and put it back afterwards - in clearReleases.
			i.queue[code] = 0
		} else {
			i.curr.Down[code] = i.curr.Down[code] + KeyReleased
		}
	}
}

// releaseAll clears the pressed map when the window loses focus
// or other things happen that invalidate the pressed map.
func (i *input) releaseAll() {
	for code, down := range i.curr.Down {
		i.curr.Down[code] = down + KeyReleased
	}
}

// clone the current user input information into the structure that is
// shared with the outside process. Remove any released keys from the map.
// This method is expected to be called by getPressed().
func (i *input) clone(in, out *Pressed) {
	for key := range out.Down {
		delete(out.Down, key)
	}
	for key, val := range in.Down {
		out.Down[key] = val
	}
	out.Mx, out.My = in.Mx, in.My
	out.Focus = in.Focus
	out.Resized = in.Resized
	out.Scroll = in.Scroll
	in.Scroll = 0      // remove previous scroll info.
	in.Resized = false // remove previous resized trigger.
}

// clearReleases removes released key once the information has
// been polled by the application.
func (i *input) clearReleases(in *Pressed) {
	for key, val := range in.Down {
		if val < 0 {
			delete(in.Down, key) // remove released keys.
		}
	}

	// Now safe to put back any key releases that conflicted with
	// an initial key down.
	for code := range i.queue {
		i.curr.Down[code] = i.curr.Down[code] + KeyReleased
		delete(i.queue, code)
	}
}
