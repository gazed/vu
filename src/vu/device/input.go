// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package device

// input is used to process a user event stream into the Pressed structure
// that can be polled as needed. A goroutine and channel is used to help
// drain the OS event queue each time device.Update is called.
type input struct {
	events chan *userInput // Events processes the readDispatch user events.
	signal chan *Pressed   // Signal triggers the pass back.
	update chan *Pressed   // Passes back the latest pressed information.
	curr   *Pressed        // Consolidates the stream events.
	down   *Pressed        // Used to share with the main process.
}

// newInput creates the memory needed to process events and communicate with
// the main process. The goroutine listening for user input events is started.
func newInput(os *nativeOs) *input {
	i := &input{}
	i.events = make(chan *userInput) // i.processEvents() <- i.events <- os.readAndDispatch
	i.signal = make(chan *Pressed)   // i.processEvents() <- i.signal <- i.pressed()
	i.update = make(chan *Pressed)   // i.pressed()       <- i.update <- i.processEvents()
	i.curr = &Pressed{Focus: true, Down: map[string]int{}}
	i.down = &Pressed{Focus: true, Down: map[string]int{}}
	go i.processEvents()
	return i
}

// processEvents loops forever processing input events into current user input.
func (i *input) processEvents() {
	for {
		select {
		case event := <-i.events:
			i.processEvent(event)
		case <-i.signal:
			i.updateDurations()
			i.clone(i.curr, i.down)
			i.update <- i.down
		}
	}
}

// KEY_RELEASED is used to indicate a key up event has occurred.
// The total duration of a key press can be calculated by the difference
// of Pressed.Down duration with KEY_RELEASED. A user would have to hold
// a key down for 24 hours before the released duration became positive
// (assuming a reasonable update time of 0.02 seconds).
const KEY_RELEASED = -1000000000

// processEvents updates the current input event buffer essentially turning
// the user input stream into a map of what is currently pressed. A duration
// of how long each key has been pressed is recorded in update ticks.
// This method is only expected to be called by i.processEvents().
func (i *input) processEvent(event *userInput) {
	i.curr.Mx, i.curr.My = event.mouseX, event.mouseY
	i.curr.Scroll = event.scroll

	// capture modifier key state.
	if event.mods&shiftKeyMask != 0 {
		i.recordPress(shiftKey)
	} else {
		i.recordRelease(shiftKey)
	}
	if event.mods&controlKeyMask != 0 {
		i.recordPress(controlKey)
	} else {
		i.recordRelease(controlKey)
	}
	if event.mods&functionKeyMask != 0 {
		i.recordPress(functionKey)
	} else {
		i.recordRelease(functionKey)
	}
	if event.mods&commandKeyMask != 0 {
		i.recordPress(commandKey)
	} else {
		i.recordRelease(commandKey)
	}
	if event.mods&altKeyMask != 0 {
		i.recordPress(altKey)
	} else {
		i.recordRelease(altKey)
	}

	// turn key and mouse events into state
	switch event.id {
	case resizedShell, movedShell:
		i.curr.Resized = true
	case activatedShell, uniconifiedShell:
		i.curr.Focus = true
	case deactivatedShell, iconifiedShell:
		i.curr.Focus = false
	case clickedMouse:
		i.recordPress(event.button)
	case releasedMouse:
		i.recordRelease(event.button)
	case pressedKey:
		i.recordPress(event.key)
	case releasedKey:
		i.recordRelease(event.key)
	}
}

// recordPress tracks new key or mouse down user input events.
func (i *input) recordPress(code int) {
	pressed := keyNames[code]
	if _, ok := i.curr.Down[pressed]; !ok {
		i.curr.Down[pressed] = 0
	}
}

// recordRelease tracks key or mouse up user input events.
func (i *input) recordRelease(code int) {
	released := keyNames[code]
	if _, ok := i.curr.Down[released]; ok {
		i.curr.Down[released] = i.curr.Down[released] + KEY_RELEASED
	}
}

// updateDurations tracks how long keys have been pressed for.
// Expected to be called each signal (update loop). Ignore released keys.
func (i *input) updateDurations() {
	for key, val := range i.curr.Down {
		if val >= 0 {
			i.curr.Down[key] = val + 1
		}
	}
}

// clone the current user input information into the structure that is
// shared with the outside process. Remove any released keys from the map.
// This method is expected to be called by i.processEvents() on a signal
// from the outside process.
func (i *input) clone(in, out *Pressed) {
	for key, _ := range out.Down {
		delete(out.Down, key)
	}
	for key, val := range in.Down {
		out.Down[key] = val
		if val < 0 {
			delete(in.Down, key) // remove released keys.
		}
	}
	out.Mx, out.My = in.Mx, in.My
	out.Focus = in.Focus
	out.Resized = in.Resized
	out.Scroll = in.Scroll
	in.Resized = false
}

// latest returns the most recent user input events.
// It is used by the outside process to communicate with the i.processEvents()
// goroutine that is collecting the events.
func (i *input) latest() *Pressed {
	i.signal <- i.down // signal the goroutine to populate with the latest events.
	return <-i.update  // return the updated event information.
}

// pressed
// ===========================================================================
// userInput

// userInput is returned by readAndDispatch.  It is the current input of the
// keyboard and mouse.  It tracks any user input changes.
type userInput struct {
	id     int // Unique event id: one of the const's below
	mouseX int // Current mouse X position.
	mouseY int // Current mouse Y position.
	button int // Currently pressed mouse button (if any).
	key    int // Current key pressed (if any).
	mods   int // Mask of the current modifier keys (if any).
	scroll int // Scroll amount (if any).
}

// userInput
// ===========================================================================
// internal event, key, key-code to string mappings.

// The possible event id's are as follows.
const (
	_ = iota // valid event ids start at one.
	closedShell
	resizedShell
	movedShell
	iconifiedShell
	uniconifiedShell
	activatedShell
	deactivatedShell
	clickedMouse
	releasedMouse
	draggedMouse
	movedMouse
	pressedKey
	releasedKey
	scrolled
)

// Map the modifier key mask values so they don't conflict with regular keys.
const (
	controlKey  = controlKeyMask + 0xFF
	shiftKey    = shiftKeyMask + 0xFF
	functionKey = functionKeyMask + 0xFF
	commandKey  = commandKeyMask + 0xFF
	altKey      = altKeyMask + 0xFF
)

// keyNames holds the names of the individual key presses.
var keyNames map[int]string = map[int]string{
	key_0:              "0",
	key_1:              "1",
	key_2:              "2",
	key_3:              "3",
	key_4:              "4",
	key_5:              "5",
	key_6:              "6",
	key_7:              "7",
	key_8:              "8",
	key_9:              "9",
	key_A:              "A",
	key_B:              "B",
	key_C:              "C",
	key_D:              "D",
	key_E:              "E",
	key_F:              "F",
	key_H:              "H",
	key_G:              "G",
	key_I:              "I",
	key_K:              "K",
	key_J:              "J",
	key_L:              "L",
	key_M:              "M",
	key_N:              "N",
	key_O:              "O",
	key_P:              "P",
	key_Q:              "Q",
	key_R:              "R",
	key_S:              "S",
	key_T:              "T",
	key_U:              "U",
	key_V:              "V",
	key_W:              "W",
	key_X:              "X",
	key_Y:              "Y",
	key_Z:              "Z",
	key_KeypadDecimal:  "KP.",
	key_KeypadMultiply: "KP*",
	key_KeypadPlus:     "KP+",
	key_KeypadClear:    "KPCl",
	key_KeypadDivide:   "KP/",
	key_KeypadEnter:    "KPEnt",
	key_KeypadMinus:    "KP-",
	key_KeypadEquals:   "KP=",
	key_Keypad0:        "KP0",
	key_Keypad1:        "KP1",
	key_Keypad2:        "KP2",
	key_Keypad3:        "KP3",
	key_Keypad4:        "KP4",
	key_Keypad5:        "KP5",
	key_Keypad6:        "KP6",
	key_Keypad7:        "KP7",
	key_Keypad8:        "KP8",
	key_Keypad9:        "KP9",
	key_LeftArrow:      "La",
	key_RightArrow:     "Ra",
	key_DownArrow:      "Da",
	key_UpArrow:        "Ua",
	key_F1:             "F1",
	key_F2:             "F2",
	key_F3:             "F3",
	key_F4:             "F4",
	key_F5:             "F5",
	key_F6:             "F6",
	key_F7:             "F7",
	key_F8:             "F8",
	key_F9:             "F9",
	key_F10:            "F10",
	key_F11:            "F11",
	key_F12:            "F12",
	key_F13:            "F13",
	key_F14:            "F14",
	key_F15:            "F15",
	key_F16:            "F16",
	key_F17:            "F17",
	key_F18:            "F18",
	key_F19:            "F19",
	key_Equal:          "=",
	key_Minus:          "-",
	key_RightBracket:   "]",
	key_LeftBracket:    "[",
	key_Quote:          "Qt",
	key_Semicolon:      ";",
	key_Backslash:      "Bs",
	key_Comma:          ",",
	key_Slash:          "Sl",
	key_Period:         ".",
	key_Grave:          "~",
	key_Return:         "Ret",
	key_Tab:            "Tab",
	key_Space:          "Sp",
	key_Delete:         "Del",
	key_Escape:         "Esc",
	key_Home:           "Home",
	key_PageUp:         "Pup",
	key_ForwardDelete:  "FDel",
	key_End:            "End",
	key_PageDown:       "Pdn",
	mouse_Left:         "Lm", // Treat pressing a mouse button like pressing a key.
	mouse_Right:        "Rm",
	mouse_Middle:       "Mm",
	controlKey:         "Ctl",
	functionKey:        "Fn",
	shiftKey:           "Sh",
	commandKey:         "Cmd",
	altKey:             "Alt",
}
