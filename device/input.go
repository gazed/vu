// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package device

// Design note: The original intent was to collect and process the OS
// event queues concurrently. However, OSX only allows event processng
// from the main thread.

// input is used to process a user event stream into the Pressed structure
// that can be polled as needed.
type input struct {
	in   *userInput // Input is processed in a map of pressed keys.
	curr *Pressed   // Consolidates current user events into state.
	down *Pressed   // Clone of curr that is shared with the application.
}

// newInput creates the memory needed to process user input events.
func newInput() *input {
	i := &input{}
	i.in = &userInput{}
	i.curr = &Pressed{Focus: true, Down: map[int]int{}}
	i.down = &Pressed{Focus: true, Down: map[int]int{}}
	return i
}

// pollEvents is called from the main thread as some OS's only allow
// event processing from the main thread. The events are placed in
// the processing queue.
func (i *input) pollEvents(os *nativeOs) *Pressed {
	i.processEvent(os.readDispatch(i.in)) // sample events at twice the update rate
	i.processEvent(os.readDispatch(i.in)) // ...by reading 2 events each update.
	i.updateDurations()
	i.clone(i.curr, i.down)
	return i.down
}

// processEvents updates the current input event buffer essentially turning
// the user input stream into a map of what is currently pressed. A duration
// of how long each key has been pressed is recorded in update ticks.
// This method is only expected to be called by i.pollEvents().
func (i *input) processEvent(event *userInput) {
	i.curr.Mx, i.curr.My = event.mouseX, event.mouseY
	i.curr.Scroll += event.scroll

	// turn key and mouse events into state
	switch event.id {
	case resizedShell, movedShell:
		i.curr.Resized = true
	case activatedShell, uniconifiedShell:
		i.curr.Focus = true
	case deactivatedShell, iconifiedShell:
		i.curr.Focus = false
		i.releaseAll()
	case clickedMouse:
		i.recordPress(event.button)
	case releasedMouse:
		i.recordRelease(event.button)
	case pressedKey:
		i.recordPress(event.key)
	case releasedKey:
		i.recordRelease(event.key)
	default:
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
	if _, ok := i.curr.Down[code]; ok {
		i.curr.Down[code] = i.curr.Down[code] + KEY_RELEASED
	}
}

// releaseAll clears the pressed map when the window loses focus
// or other things happen that invalidate the pressed map.
func (i *input) releaseAll() {
	for code, down := range i.curr.Down {
		i.curr.Down[code] = down + KEY_RELEASED
	}
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

// clone the current user input information into the structure that is
// shared with the outside process. Remove any released keys from the map.
// This method is expected to be called by i.pollEvents().
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
	in.Scroll = 0      // remove previous scroll info.
	in.Resized = false // remove previous resized trigger.
}

// input
// ===========================================================================
// userInput

// userInput is returned by readAndDispatch. It is the current input of the
// keyboard and mouse. It tracks any user input changes.
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
	invalid = iota // valid event ids start at one.
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

// Modifier key values don't conflict with regular keys.
const (
	controlKey  = controlKeyMask
	shiftKey    = shiftKeyMask
	functionKey = functionKeyMask
	commandKey  = commandKeyMask
	altKey      = altKeyMask
)

// Based on the keys on a Mac OSX extended keyboard excluding
// OS specific keys like eject. Most keyboards will support
// some subset of the following keys. Currently pressed keys are
// returned in the Pressed.Down map.
const (
	K_0     = key_0              // Standard keyboard numbers.
	K_1     = key_1              //   "
	K_2     = key_2              //   "
	K_3     = key_3              //   "
	K_4     = key_4              //   "
	K_5     = key_5              //   "
	K_6     = key_6              //   "
	K_7     = key_7              //   "
	K_8     = key_8              //   "
	K_9     = key_9              //   "
	K_A     = key_A              // Standard keyboard letters.
	K_B     = key_B              //   "
	K_C     = key_C              //   "
	K_D     = key_D              //   "
	K_E     = key_E              //   "
	K_F     = key_F              //   "
	K_G     = key_G              //   "
	K_H     = key_H              //   "
	K_I     = key_I              //   "
	K_J     = key_J              //   "
	K_K     = key_K              //   "
	K_L     = key_L              //   "
	K_M     = key_M              //   "
	K_N     = key_N              //   "
	K_O     = key_O              //   "
	K_P     = key_P              //   "
	K_Q     = key_Q              //   "
	K_R     = key_R              //   "
	K_S     = key_S              //   "
	K_T     = key_T              //   "
	K_U     = key_U              //   "
	K_V     = key_V              //   "
	K_W     = key_W              //   "
	K_X     = key_X              //   "
	K_Y     = key_Y              //   "
	K_Z     = key_Z              //   "
	K_Equal = key_Equal          // Standard keyboard punctuation keys.
	K_Minus = key_Minus          //   "
	K_RBkt  = key_RightBracket   //   "
	K_LBkt  = key_LeftBracket    //   "
	K_Qt    = key_Quote          //   "
	K_Semi  = key_Semicolon      //   "
	K_BSl   = key_Backslash      //   "
	K_Comma = key_Comma          //   "
	K_Slash = key_Slash          //   "
	K_Dot   = key_Period         //   "
	K_Grave = key_Grave          //   "
	K_Ret   = key_Return         //   "
	K_Tab   = key_Tab            //   "
	K_Space = key_Space          //   "
	K_Del   = key_Delete         //   "
	K_Esc   = key_Escape         //   "
	K_F1    = key_F1             // General Function keys.
	K_F2    = key_F2             //   "
	K_F3    = key_F3             //   "
	K_F4    = key_F4             //   "
	K_F5    = key_F5             //   "
	K_F6    = key_F6             //   "
	K_F7    = key_F7             //   "
	K_F8    = key_F8             //   "
	K_F9    = key_F9             //   "
	K_F10   = key_F10            //   "
	K_F11   = key_F11            //   "
	K_F12   = key_F12            //   "
	K_F13   = key_F13            //   "
	K_F14   = key_F14            //   "
	K_F15   = key_F15            //   "
	K_F16   = key_F16            //   "
	K_F17   = key_F17            //   "
	K_F18   = key_F18            //   "
	K_F19   = key_F19            //   "
	K_Home  = key_Home           // Specific function keys.
	K_PgUp  = key_PageUp         //   "
	K_FDel  = key_ForwardDelete  //   "
	K_End   = key_End            //   "
	K_PgDn  = key_PageDown       //   "
	K_La    = key_LeftArrow      // Arrow keys
	K_Ra    = key_RightArrow     //   "
	K_Da    = key_DownArrow      //   "
	K_Ua    = key_UpArrow        //   "
	K_KpDot = key_KeypadDecimal  // Extended keyboard keypad keys
	K_KpMlt = key_KeypadMultiply //   "
	K_KpAdd = key_KeypadPlus     //   "
	K_KpClr = key_KeypadClear    //   "
	K_KpDiv = key_KeypadDivide   //   "
	K_KpEnt = key_KeypadEnter    //   "
	K_KpSub = key_KeypadMinus    //   "
	K_KpEql = key_KeypadEquals   //   "
	K_Kp0   = key_Keypad0        //   "
	K_Kp1   = key_Keypad1        //   "
	K_Kp2   = key_Keypad2        //   "
	K_Kp3   = key_Keypad3        //   "
	K_Kp4   = key_Keypad4        //   "
	K_Kp5   = key_Keypad5        //   "
	K_Kp6   = key_Keypad6        //   "
	K_Kp7   = key_Keypad7        //   "
	K_Kp8   = key_Keypad8        //   "
	K_Kp9   = key_Keypad9        //   "
	K_Lm    = mouse_Left         // Mouse buttons treated like keys.
	K_Mm    = mouse_Middle       //   "
	K_Rm    = mouse_Right        //   "
	K_Ctl   = controlKey         // Modifier keys.
	K_Fn    = functionKey        //   "
	K_Shift = shiftKey           //   "
	K_Cmd   = commandKey         //   "
	K_Alt   = altKey             //   "
)
