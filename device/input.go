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
		i.curr.Down[code] = i.curr.Down[code] + KeyReleased
	}
}

// releaseAll clears the pressed map when the window loses focus
// or other things happen that invalidate the pressed map.
func (i *input) releaseAll() {
	for code, down := range i.curr.Down {
		i.curr.Down[code] = down + KeyReleased
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
	for key := range out.Down {
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
	K0     = key0              // Standard keyboard numbers.
	K1     = key1              //   "
	K2     = key2              //   "
	K3     = key3              //   "
	K4     = key4              //   "
	K5     = key5              //   "
	K6     = key6              //   "
	K7     = key7              //   "
	K8     = key8              //   "
	K9     = key9              //   "
	KA     = keyA              // Standard keyboard letters.
	KB     = keyB              //   "
	KC     = keyC              //   "
	KD     = keyD              //   "
	KE     = keyE              //   "
	KF     = keyF              //   "
	KG     = keyG              //   "
	KH     = keyH              //   "
	KI     = keyI              //   "
	KJ     = keyJ              //   "
	KK     = keyK              //   "
	KL     = keyL              //   "
	KM     = keyM              //   "
	KN     = keyN              //   "
	KO     = keyO              //   "
	KP     = keyP              //   "
	KQ     = keyQ              //   "
	KR     = keyR              //   "
	KS     = keyS              //   "
	KT     = keyT              //   "
	KU     = keyU              //   "
	KV     = keyV              //   "
	KW     = keyW              //   "
	KX     = keyX              //   "
	KY     = keyY              //   "
	KZ     = keyZ              //   "
	KEqual = keyEqual          // Standard keyboard punctuation keys.
	KMinus = keyMinus          //   "
	KRBkt  = keyRightBracket   //   "
	KLBkt  = keyLeftBracket    //   "
	KQt    = keyQuote          //   "
	KSemi  = keySemicolon      //   "
	KBSl   = keyBackslash      //   "
	KComma = keyComma          //   "
	KSlash = keySlash          //   "
	KDot   = keyPeriod         //   "
	KGrave = keyGrave          //   "
	KRet   = keyReturn         //   "
	KTab   = keyTab            //   "
	KSpace = keySpace          //   "
	KDel   = keyDelete         //   "
	KEsc   = keyEscape         //   "
	KF1    = keyF1             // General Function keys.
	KF2    = keyF2             //   "
	KF3    = keyF3             //   "
	KF4    = keyF4             //   "
	KF5    = keyF5             //   "
	KF6    = keyF6             //   "
	KF7    = keyF7             //   "
	KF8    = keyF8             //   "
	KF9    = keyF9             //   "
	KF10   = keyF10            //   "
	KF11   = keyF11            //   "
	KF12   = keyF12            //   "
	KF13   = keyF13            //   "
	KF14   = keyF14            //   "
	KF15   = keyF15            //   "
	KF16   = keyF16            //   "
	KF17   = keyF17            //   "
	KF18   = keyF18            //   "
	KF19   = keyF19            //   "
	KHome  = keyHome           // Specific function keys.
	KPgUp  = keyPageUp         //   "
	KFDel  = keyForwardDelete  //   "
	KEnd   = keyEnd            //   "
	KPgDn  = keyPageDown       //   "
	KLa    = keyLeftArrow      // Arrow keys
	KRa    = keyRightArrow     //   "
	KDa    = keyDownArrow      //   "
	KUa    = keyUpArrow        //   "
	KKpDot = keyKeypadDecimal  // Extended keyboard keypad keys
	KKpMlt = keyKeypadMultiply //   "
	KKpAdd = keyKeypadPlus     //   "
	KKpClr = keyKeypadClear    //   "
	KKpDiv = keyKeypadDivide   //   "
	KKpEnt = keyKeypadEnter    //   "
	KKpSub = keyKeypadMinus    //   "
	KKpEql = keyKeypadEquals   //   "
	KKp0   = keyKeypad0        //   "
	KKp1   = keyKeypad1        //   "
	KKp2   = keyKeypad2        //   "
	KKp3   = keyKeypad3        //   "
	KKp4   = keyKeypad4        //   "
	KKp5   = keyKeypad5        //   "
	KKp6   = keyKeypad6        //   "
	KKp7   = keyKeypad7        //   "
	KKp8   = keyKeypad8        //   "
	KKp9   = keyKeypad9        //   "
	KLm    = mouseLeft         // Mouse buttons treated like keys.
	KMm    = mouseMiddle       //   "
	KRm    = mouseRight        //   "
	KCtl   = controlKey        // Modifier keys.
	KFn    = functionKey       //   "
	KShift = shiftKey          //   "
	KCmd   = commandKey        //   "
	KAlt   = altKey            //   "
)
