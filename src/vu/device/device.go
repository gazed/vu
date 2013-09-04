// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

// Package device provides platform/os access to a 3D rendering
// context and user input. It attempts to do the absolute minimum in
// order to get a OS specific window with a 3D rendering context.
// The expectation is that the 3D application itself will provide any
// necessary buttons, controls, and sub-panels.
//
// Package device also allows access to user keyboard and mouse input.
//
// Big thanks to GLFW (http://www.glfw.org) from which the minimalist API
// philosophy was borrowed along with which OS specific API's mattered.
//
// Package device is provided as part of the vu (virtual universe) 3D engine.
package device

//
// Note that a dynamic library is provided for wrapping OSX until golang
// issues 1781, 4069, 5726, 5699 are resolved.  Currently Windows 7 and
// OSX have been tested. Linux (ubuntu) development is planned.

// Device wraps OS specific functionality. It is pretty much the absolute minimum
// that could be done in order to get an OS specific window for running a 3D
// application.  The expected usage is:
//     dev := device.New("title", x, y, width, height)
//     // application initialization code.
//     dev.Open()
//     for dev.IsAlive() {
//         dev.ReadAndDispatch()
//         // application update and render code
//         dev.SwapBuffers()
//     }
//     dev.Dispose()
type Device interface {
	Open()                      // Open displays the window created using device.New().
	ShowCursor(show bool)       // ShowCursor displays or hides the cursor.
	SetCursorAt(x, y int)       // SetCursorAt puts the mouse at the given window location.
	SetResizer(resizer Resizer) // SetResizer is a hook to receive resize callbacks.
	SetFocuser(focuser Focuser) // SetFocuser is a hook to receive window focus callbacks.

	// Dispose is called to release OS specfic resources in an orderly manner.
	// It is called once the window is no longer needed.
	Dispose()

	// IsAlive returns true as long as the window is able to process user input.
	// Generally quitting the application window will cause IsAlive to return false.
	IsAlive() bool

	// Size returns the usable graphics context location and size excluding any OS
	// specific window trim.  The window x, y coordinates are the bottom left of
	// the window.
	Size() (x, y, width, height int)

	// SwapBuffers exchanges the graphic drawing buffers (all rendering contexts
	// are double buffered).
	SwapBuffers()

	// ReadAndDispatch must be called regularly to process OS events.  This returns
	// the currently pressed keys/mouse-buttons (combined with modified keys), and
	// the current mouse location.  For example if the user was holding down the
	// Shift key, the A-key and the left mouse button then the following strings
	// would be passed back from ReadAndDispatch:
	//     ["Lm", "Sh-A"]
	ReadAndDispatch() (down []string, mx, my int)
}

// Resizer defines something that wants to handle window resize and move events.
// Any changes to the location or width and height of a window is passed back
// in Resize.
type Resizer interface {
	Resize(x, y, width, height int)
}

// Focuser defines something that wants to handle window focus events.
type Focuser interface {
	Focus(hasFocus bool)
}

// Device, Resize, Focuser interfaces
// ===========================================================================
// device provides default Device implementation.

// New provides a newly initialized device with an underlying window and
// graphics context created, but not yet opened. The only thing left to do is to
// open the device and start polling it for user input.
func New(title string, x, y, width, height int) Device { return newDevice(title, x, y, width, height) }

// Design notes. The layers in this package are:
//     device : uses context to process events into something simple.
//     context: single point of entry into the native layer.
//     os_darwin : OSX native layer. Wraps the following.
//        os_darwin.m : objective-c code wraps cocoa.
//        os_darwin.h
//        os_darwin_test.m
//     os_windows : Win native layer. Wraps the following.
//        os_windows.c : c code wrapping windows API.
//        os_windows.h
//        os_windows_test.c
//     TODO linux (ubuntu)

// device provides a simplification layer over the more raw context layer.
type device struct {
	ctx     *context
	kmods   []int   // modifier keys
	resizer Resizer // resize event handler.
	focuser Focuser // window focus event handler.
	tracker presser // user input handler.
}

// newDevice initializes a OS specific window with a valid OpenGL context.
func newDevice(title string, x, y, width, height int) *device {
	d := &device{}
	d.ctx = newContext()
	d.ctx.createDisplay(title, x, y, width, height)
	d.ctx.createShell()
	depthBufferBits, alphaBits := 24, 8 // resonable defaults
	d.ctx.createContext(depthBufferBits, alphaBits)
	d.tracker = newPressedTracker(keyNames, mods, modNames)
	return d
}

// Access the device specific information in a consistent and general manner.
func (d *device) Open()                           { d.ctx.openShell() }
func (d *device) Dispose()                        { d.ctx.dispose() }
func (d *device) IsAlive() bool                   { return d.ctx.isAlive() }
func (d *device) Size() (x, y, width, height int) { return d.ctx.size() }
func (d *device) ShowCursor(show bool)            { d.ctx.showCursor(show) }
func (d *device) SwapBuffers()                    { d.ctx.swapBuffers() }
func (d *device) SetResizer(resizer Resizer)      { d.resizer = resizer }
func (d *device) SetFocuser(focuser Focuser)      { d.focuser = focuser }
func (d *device) SetCursorAt(x, y int)            { d.ctx.setCursorAt(x, y) }

// ReadAndDispatch processes OS events that relate to any/all user input.
// This is mouse and keyboard input for the most part.
func (d *device) ReadAndDispatch() (down []string, mx, my int) {
	event := d.ctx.readAndDispatch()
	switch event.id {
	case closedShell:
	case resizedShell:
		if d.resizer != nil {
			d.resizer.Resize(d.ctx.size())
		}
	case movedShell:
		if d.resizer != nil {
			d.resizer.Resize(d.ctx.size())
		}
	case iconifiedShell:
	case uniconifiedShell:
	case activatedShell:
		if d.focuser != nil {
			d.focuser.Focus(true)
		}
	case deactivatedShell:
		if d.focuser != nil {
			d.focuser.Focus(false)
		}
	case clickedMouse:
		d.tracker.pressed(event.button, event.mods, true)
	case releasedMouse:
		d.tracker.pressed(event.button, event.mods, false)
	case draggedMouse:
	case movedMouse:
	case pressedKey:
		d.tracker.pressed(event.key, event.mods, true)
	case releasedKey:
		d.tracker.pressed(event.key, event.mods, false)
	case scrolled:
	}
	return d.tracker.down(), event.mouseX, event.mouseY
}

// device
// ===========================================================================
// presser and key definitions.

// presser defines a class that wants to handle key and mouse presses.
type presser interface {
	pressed(event int, mods int, isPressed bool) // update pressed state with this key event.
	down() []string                              // return the currently pressed key sequences.
}

// mods is the list of valid modifier keys.  Modifier keys are expected to be used
// in conjunction with other normal keys to form key sequences that identify what
// the user is currently pressing.
var mods []int = []int{controlKeyMask, functionKeyMask, shiftKeyMask}

// modNames holds the names of the above modifier keys.
var modNames []string = []string{"Ctl", "Fn", "Sh"}

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
}
