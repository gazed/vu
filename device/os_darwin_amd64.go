// Copyright © 2013-2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package device

// The macOS (darwin) native layer. This wraps the objective-c code that
// calls the OSX API's (where the real work is done).
//
// Does the minimum necessary to get a window with an OpenGL graphic context.
// All the window is expected to do is to be able to run in full screen mode,
// quit, show up in the dock, and participate in command-tab app switching.

// // The following block is C code and cgo directvies.
//
// #cgo darwin CFLAGS: -x objective-c -fno-common
// #cgo darwin LDFLAGS: -framework Cocoa -framework Quartz -framework OpenGL
//
// #import <stdlib.h>
// #import <Carbon/Carbon.h>   // for HIToolbox/Events.h kVK_* keycodes
// #import <AppKit/NSEvent.h>  // for NSEventType values
// #import "os_darwin_amd64.h" // native function signatures and constants.
import "C" // must be located just after C code and cgo directives.

import (
	"runtime"
	"unsafe"
)

// macOS specific. Otherwise shell freezes within seconds of creation.
func init() { runtime.LockOSThread() }

// Device instance needed to handle callbacks on exported methods.
var dev = &osx{}

// runApp is the per-device entry method. Compiling will find the one
// that matches the requested or current platform.
func runApp(app App) {
	dev.app = app // The app receiving device callbacks.
	C.dev_run()   // does not return!
}

// prepRender is called from os_darwin_amd64.m after the underlying
// window has been created and before the update render callbacks start.
//
//export prepRender
func prepRender() {
	dev.input = newInput()
	dev.app.Init(dev)
}

// renderFrame is the called from os_darwin_amd64.m each time
// a new render frame is needed.
//
//export renderFrame
func renderFrame(mx, my int64) {
	if dev.input != nil {
		dev.app.Refresh(dev)
	}
}

// handleInput is called from os_darwin_amd64.m to consolidate user input
// events into a consumable summary of which keys are current pressed.
//
//export handleInput
func handleInput(event, data int64) {
	if dev.input == nil {
		return // Ignore input before window is up.
	}
	in := dev.input
	switch event {
	case C.devUp:
		in.recordRelease(int(data))
	case C.devDown:
		in.recordPress(int(data))
	case C.devScroll:
		in.curr.Scroll = int(data)
	case C.devResize:
		in.curr.Resized = true

		// release all down keys on resize
		// to avoid missing key release events.
		in.releaseAll()
	case C.devFocusIn, C.devFocusOut:
		in.curr.Focus = event == C.devFocusIn
	case C.devMod:
		// capture modifier key state.
		if data&KShift != 0 {
			in.recordPress(KShift)
		} else {
			in.recordRelease(KShift)
		}
		if data&KCtl != 0 {
			in.recordPress(KCtl)
		} else {
			in.recordRelease(KCtl)
		}
		if data&KFn != 0 {
			in.recordPress(KFn)
		} else {
			in.recordRelease(KFn)
		}
		if data&KCmd != 0 {
			in.recordPress(KCmd)
		} else {
			in.recordRelease(KCmd)
		}
		if data&KAlt != 0 {
			in.recordPress(KAlt)
		} else {
			in.recordRelease(KAlt)
		}
	}
}

// native layer callback functions. native->app calls.
// =============================================================================
// macOS Device implementation. app->native calls.

// osx is the macOS implementation of the Device interface.
// See the Device interface for method descriptions.
type osx struct {
	app   App    // update/render callback.
	input *input // tracks current keys pressed.
}

// Down returns the current set of user input. Expected to be called
// once per state update (not display refresh).
func (os *osx) Down() *Pressed     { return os.input.getPressed(os.Cursor()) }
func (os *osx) Dispose()           { C.dev_dispose() }
func (os *osx) SwapBuffers()       { C.dev_swap() }
func (os *osx) ToggleFullScreen()  { C.dev_toggle_fullscreen() }
func (os *osx) IsFullScreen() bool { return uint(C.dev_fullscreen()) == 1 }
func (os *osx) SetCursorAt(x, y int) {
	C.dev_set_cursor_location(C.long(x), C.long(y))
}
func (os *osx) Cursor() (x, y int) {
	var mx, my int64
	C.dev_cursor((*C.long)(&mx), (*C.long)(&my))
	return int(mx), int(my)
}
func (os *osx) Size() (x, y, w, h int) {
	var wx, wy, ww, wh int64
	C.dev_size((*C.long)(&wx), (*C.long)(&wy), (*C.long)(&ww), (*C.long)(&wh))
	return int(wx), int(wy), int(ww), int(wh)
}
func (os *osx) SetSize(x, y, w, h int) {
	C.dev_set_size(C.long(x), C.long(y), C.long(w), C.long(h))
}
func (os *osx) SetTitle(title string) {
	cstr := C.CString(title)
	defer C.free(unsafe.Pointer(cstr))
	C.dev_set_title(cstr)
}
func (os *osx) ShowCursor(show bool) {
	trueFalse := 0 // trueFalse needs to be 0 or 1.
	if show {
		trueFalse = 1
	}
	C.dev_show_cursor(C.uchar(trueFalse))
}
func (os *osx) Copy() string {
	if cstr := C.dev_clip_copy(); cstr != nil {
		str := C.GoString(cstr)      // make a Go copy.
		C.free(unsafe.Pointer(cstr)) // free the C copy.
		return str
	}
	return ""
}
func (os *osx) Paste(s string) {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	C.dev_clip_paste(cstr)
}

// macOS Device implementation.
// =============================================================================

// Map the underlying OSX key codes to the Vu device key codes.
//
// Based on the keys on a Mac OSX extended keyboard excluding
// OS specific keys like eject. Most keyboards will support
// some subset of the following keys.
const (
	K0     = C.kVK_ANSI_0              // Standard keyboard numbers.
	K1     = C.kVK_ANSI_1              //
	K2     = C.kVK_ANSI_2              //
	K3     = C.kVK_ANSI_3              //
	K4     = C.kVK_ANSI_4              //
	K5     = C.kVK_ANSI_5              //
	K6     = C.kVK_ANSI_6              //
	K7     = C.kVK_ANSI_7              //
	K8     = C.kVK_ANSI_8              //
	K9     = C.kVK_ANSI_9              //
	KA     = C.kVK_ANSI_A              // Standard keyboard letters.
	KB     = C.kVK_ANSI_B              //
	KC     = C.kVK_ANSI_C              //
	KD     = C.kVK_ANSI_D              //
	KE     = C.kVK_ANSI_E              //
	KF     = C.kVK_ANSI_F              //
	KG     = C.kVK_ANSI_G              //
	KH     = C.kVK_ANSI_H              //
	KI     = C.kVK_ANSI_I              //
	KJ     = C.kVK_ANSI_J              //
	KK     = C.kVK_ANSI_K              //
	KL     = C.kVK_ANSI_L              //
	KM     = C.kVK_ANSI_M              //
	KN     = C.kVK_ANSI_N              //
	KO     = C.kVK_ANSI_O              //
	KP     = C.kVK_ANSI_P              //
	KQ     = C.kVK_ANSI_Q              //
	KR     = C.kVK_ANSI_R              //
	KS     = C.kVK_ANSI_S              //
	KT     = C.kVK_ANSI_T              //
	KU     = C.kVK_ANSI_U              //
	KV     = C.kVK_ANSI_V              //
	KW     = C.kVK_ANSI_W              //
	KX     = C.kVK_ANSI_X              //
	KY     = C.kVK_ANSI_Y              //
	KZ     = C.kVK_ANSI_Z              //
	KF1    = C.kVK_F1                  // General Function Keys
	KF2    = C.kVK_F2                  //
	KF3    = C.kVK_F3                  //
	KF4    = C.kVK_F4                  //
	KF5    = C.kVK_F5                  //
	KF6    = C.kVK_F6                  //
	KF7    = C.kVK_F7                  //
	KF8    = C.kVK_F8                  //
	KF9    = C.kVK_F9                  //
	KF10   = C.kVK_F10                 //
	KF11   = C.kVK_F11                 //
	KF12   = C.kVK_F12                 //
	KF13   = C.kVK_F13                 //
	KF14   = C.kVK_F14                 //
	KF15   = C.kVK_F15                 //
	KF16   = C.kVK_F16                 //
	KF17   = C.kVK_F17                 //
	KF18   = C.kVK_F18                 //
	KF19   = C.kVK_F19                 //
	KF20   = C.kVK_F20                 // TODO check if available on Windows.
	KKpDot = C.kVK_ANSI_KeypadDecimal  // Extended keyboard keypad keys
	KKpMlt = C.kVK_ANSI_KeypadMultiply //   "
	KKpAdd = C.kVK_ANSI_KeypadPlus     //   "
	KKpClr = C.kVK_ANSI_KeypadClear    //   "
	KKpDiv = C.kVK_ANSI_KeypadDivide   //   "
	KKpEnt = C.kVK_ANSI_KeypadEnter    //   "
	KKpSub = C.kVK_ANSI_KeypadMinus    //   "
	KKpEql = C.kVK_ANSI_KeypadEquals   //   "
	KKp0   = C.kVK_ANSI_Keypad0        //   "
	KKp1   = C.kVK_ANSI_Keypad1        //   "
	KKp2   = C.kVK_ANSI_Keypad2        //   "
	KKp3   = C.kVK_ANSI_Keypad3        //   "
	KKp4   = C.kVK_ANSI_Keypad4        //   "
	KKp5   = C.kVK_ANSI_Keypad5        //   "
	KKp6   = C.kVK_ANSI_Keypad6        //   "
	KKp7   = C.kVK_ANSI_Keypad7        //   "
	KKp8   = C.kVK_ANSI_Keypad8        //   "
	KKp9   = C.kVK_ANSI_Keypad9        //   "
	KEqual = C.kVK_ANSI_Equal          // Standard keyboard punctuation keys.
	KMinus = C.kVK_ANSI_Minus          //   "
	KLBkt  = C.kVK_ANSI_LeftBracket    //   "
	KRBkt  = C.kVK_ANSI_RightBracket   //   "
	KQt    = C.kVK_ANSI_Quote          //   "
	KSemi  = C.kVK_ANSI_Semicolon      //   "
	KBSl   = C.kVK_ANSI_Backslash      //   "
	KComma = C.kVK_ANSI_Comma          //   "
	KSlash = C.kVK_ANSI_Slash          //   "
	KDot   = C.kVK_ANSI_Period         //   "
	KGrave = C.kVK_ANSI_Grave          //   "
	KRet   = C.kVK_Return              //   "
	KTab   = C.kVK_Tab                 //   "
	KSpace = C.kVK_Space               //   "
	KDel   = C.kVK_Delete              //   "
	KEsc   = C.kVK_Escape              //   "
	KHome  = C.kVK_Home                // Specific function keys.
	KPgUp  = C.kVK_PageUp              //   "
	KFDel  = C.kVK_ForwardDelete       //   "
	KEnd   = C.kVK_End                 //   "
	KPgDn  = C.kVK_PageDown            //   "
	KLa    = C.kVK_LeftArrow           // Arrow keys
	KRa    = C.kVK_RightArrow          //   "
	KDa    = C.kVK_DownArrow           //   "
	KUa    = C.kVK_UpArrow             //   "

	// Note that modifier masks are also used as key values
	// since they don't conflict with the standard key codes.
	KCtl   = C.NSEventModifierFlagControl  // Modifier masks and key codes.
	KFn    = C.NSEventModifierFlagFunction //   "
	KShift = C.NSEventModifierFlagShift    //   "
	KCmd   = C.NSEventModifierFlagCommand  //   "
	KAlt   = C.NSEventModifierFlagOption   //   "

	// Mouse buttons are treated like keys. Values don't
	// conflict with other key codes.
	KLm = C.devMouseL // Mouse buttons
	KMm = C.devMouseM //   "
	KRm = C.devMouseR //   "
)
