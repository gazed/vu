// Copyright © 2025 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

//go:build !ios

package macos

// The macOS (darwin) c-go binding layer. This wraps the objective-c code that
// calls the OSX API's (where the real work is done).
//
// This is the minimum necessary to get a window with a graphic context.
// The window is expected to do is to be able to run in full screen mode,
// quit, show up in the dock, and participate in command-tab app switching.
//
// All c-go calls are encapsulated in this package.

// // The following block is C code and cgo directvies.
//
// #cgo darwin CFLAGS: -x objective-c
// #cgo darwin LDFLAGS: -framework Cocoa -framework Metal -framework MetalKit -framework QuartzCore
//
// #import <stdlib.h>
// #import <Carbon/Carbon.h>   // for HIToolbox/Events.h kVK_* keycodes
// #import <AppKit/NSEvent.h>  // for NSEventType values
// #import "os_darwin.h"       // native function signatures and constants.
import "C" // must be located just after C code and cgo directives.

import "unsafe"

// =============================================================================
// macOS Device implementation. app->native calls.

// The callbacks for communicating with the vu/device layer.
var (
	renderer func()                  // render callback function.
	inputCB  func(event, data int64) // input callback function.
)

func CreateDisplay(title string, x, y, w, h uint32) (view unsafe.Pointer) {
	cstr := C.CString(title)
	defer C.free(unsafe.Pointer(cstr))

	// C.long returned by dev_init
	viewLong := C.dev_init(cstr, C.long(x), C.long(y), C.long(w), C.long(h))
	viewPtr := uintptr(viewLong)   // change to pointer
	return unsafe.Pointer(viewPtr) // change to unsafe.Pointer
}

// Run starts the platform render loop and does not return.
func Run(renderCallback func(), inputCallback func(ev, dat int64)) {
	renderer = renderCallback
	inputCB = inputCallback
	C.dev_run()
}

func SurfaceSize() (w, h uint32) {
	var wx, wy, ww, wh int64
	C.dev_size((*C.long)(&wx), (*C.long)(&wy), (*C.long)(&ww), (*C.long)(&wh))
	return uint32(ww), uint32(wh)
}
func SurfaceLocation() (x, y int32) {
	var wx, wy, ww, wh int64
	C.dev_size((*C.long)(&wx), (*C.long)(&wy), (*C.long)(&ww), (*C.long)(&wh))
	return int32(wx), int32(wy)
}
func SetTitle(title string) {
	cstr := C.CString(title)
	defer C.free(unsafe.Pointer(cstr))
	C.dev_set_title(cstr)
}

func Dispose()          { C.dev_dispose() }
func ToggleFullscreen() { C.dev_toggle_fullscreen() }

// get the current mouse position.
func MousePosition() (x, y int32) {
	var mx, my int64
	C.dev_cursor((*C.long)(&mx), (*C.long)(&my))
	return int32(mx), int32(my)
}

// Called when the display needs renderering.
// Generally related to the display refresh rate.
//
//export renderFrame
func renderFrame() {
	if renderer != nil {
		renderer()
	}
}

// Called on user input events.
//
//export handleInput
func handleInput(event, data int64) {
	if inputCB != nil {
		inputCB(event, data)
	}
}

// macOS Device implementation.
// =============================================================================

// Map the underlying OSX key codes to the Vu device key codes.
//
// Based on the keys on a Mac OSX extended keyboard excluding
// OS specific keys like eject. Most keyboards will support
// some subset of the following keys.
const (
	EVENT_KEYUP        = C.devUp
	EVENT_KEYDOWN      = C.devDown
	EVENT_SCROLL       = C.devScroll
	EVENT_MODIFIER     = C.devMod
	EVENT_MOVED        = C.devMoved
	EVENT_RESIZED      = C.devResized
	EVENT_FOCUS_GAINED = C.devFocusIn
	EVENT_FOCUS_LOST   = C.devFocusOut

	K0      = C.kVK_ANSI_0              // Standard keyboard numbers.
	K1      = C.kVK_ANSI_1              //
	K2      = C.kVK_ANSI_2              //
	K3      = C.kVK_ANSI_3              //
	K4      = C.kVK_ANSI_4              //
	K5      = C.kVK_ANSI_5              //
	K6      = C.kVK_ANSI_6              //
	K7      = C.kVK_ANSI_7              //
	K8      = C.kVK_ANSI_8              //
	K9      = C.kVK_ANSI_9              //
	KA      = C.kVK_ANSI_A              // Standard keyboard letters.
	KB      = C.kVK_ANSI_B              //
	KC      = C.kVK_ANSI_C              //
	KD      = C.kVK_ANSI_D              //
	KE      = C.kVK_ANSI_E              //
	KF      = C.kVK_ANSI_F              //
	KG      = C.kVK_ANSI_G              //
	KH      = C.kVK_ANSI_H              //
	KI      = C.kVK_ANSI_I              //
	KJ      = C.kVK_ANSI_J              //
	KK      = C.kVK_ANSI_K              //
	KL      = C.kVK_ANSI_L              //
	KM      = C.kVK_ANSI_M              //
	KN      = C.kVK_ANSI_N              //
	KO      = C.kVK_ANSI_O              //
	KP      = C.kVK_ANSI_P              //
	KQ      = C.kVK_ANSI_Q              //
	KR      = C.kVK_ANSI_R              //
	KS      = C.kVK_ANSI_S              //
	KT      = C.kVK_ANSI_T              //
	KU      = C.kVK_ANSI_U              //
	KV      = C.kVK_ANSI_V              //
	KW      = C.kVK_ANSI_W              //
	KX      = C.kVK_ANSI_X              //
	KY      = C.kVK_ANSI_Y              //
	KZ      = C.kVK_ANSI_Z              //
	KF1     = C.kVK_F1                  // General Function Keys
	KF2     = C.kVK_F2                  //
	KF3     = C.kVK_F3                  //
	KF4     = C.kVK_F4                  //
	KF5     = C.kVK_F5                  //
	KF6     = C.kVK_F6                  //
	KF7     = C.kVK_F7                  //
	KF8     = C.kVK_F8                  //
	KF9     = C.kVK_F9                  //
	KF10    = C.kVK_F10                 //
	KF11    = C.kVK_F11                 //
	KF12    = C.kVK_F12                 //
	KF13    = C.kVK_F13                 //
	KF14    = C.kVK_F14                 //
	KF15    = C.kVK_F15                 //
	KF16    = C.kVK_F16                 //
	KF17    = C.kVK_F17                 //
	KF18    = C.kVK_F18                 //
	KF19    = C.kVK_F19                 //
	KF20    = C.kVK_F20                 //
	KPDot   = C.kVK_ANSI_KeypadDecimal  // Extended keyboard keypad keys
	KPMlt   = C.kVK_ANSI_KeypadMultiply //   "
	KPAdd   = C.kVK_ANSI_KeypadPlus     //   "
	KPClr   = C.kVK_ANSI_KeypadClear    //   "
	KPDiv   = C.kVK_ANSI_KeypadDivide   //   "
	KPEnt   = C.kVK_ANSI_KeypadEnter    //   "
	KPSub   = C.kVK_ANSI_KeypadMinus    //   "
	KPEql   = C.kVK_ANSI_KeypadEquals   //   "
	KP0     = C.kVK_ANSI_Keypad0        //   "
	KP1     = C.kVK_ANSI_Keypad1        //   "
	KP2     = C.kVK_ANSI_Keypad2        //   "
	KP3     = C.kVK_ANSI_Keypad3        //   "
	KP4     = C.kVK_ANSI_Keypad4        //   "
	KP5     = C.kVK_ANSI_Keypad5        //   "
	KP6     = C.kVK_ANSI_Keypad6        //   "
	KP7     = C.kVK_ANSI_Keypad7        //   "
	KP8     = C.kVK_ANSI_Keypad8        //   "
	KP9     = C.kVK_ANSI_Keypad9        //   "
	KEqual  = C.kVK_ANSI_Equal          // Standard keyboard punctuation keys.
	KMinus  = C.kVK_ANSI_Minus          //   "
	KLBkt   = C.kVK_ANSI_LeftBracket    //   "
	KRBkt   = C.kVK_ANSI_RightBracket   //   "
	KQuote  = C.kVK_ANSI_Quote          //   "
	KSemi   = C.kVK_ANSI_Semicolon      //   "
	KBSl    = C.kVK_ANSI_Backslash      //   "
	KComma  = C.kVK_ANSI_Comma          //   "
	KSlash  = C.kVK_ANSI_Slash          //   "
	KDot    = C.kVK_ANSI_Period         //   "
	KGrave  = C.kVK_ANSI_Grave          //   "
	KRet    = C.kVK_Return              //   "
	KTab    = C.kVK_Tab                 //   "
	KSpace  = C.kVK_Space               //   "
	KDel    = C.kVK_Delete              //   "
	KEsc    = C.kVK_Escape              //   "
	KHome   = C.kVK_Home                // Specific function keys.
	KPgUp   = C.kVK_PageUp              //   "
	KFDel   = C.kVK_ForwardDelete       //   "
	KEnd    = C.kVK_End                 //   "
	KPgDn   = C.kVK_PageDown            //   "
	KALeft  = C.kVK_LeftArrow           // Arrow keys
	KARight = C.kVK_RightArrow          //   "
	KADown  = C.kVK_DownArrow           //   "
	KAUp    = C.kVK_UpArrow             //   "

	// Note that modifier masks are also used as key values
	// since they don't conflict with the standard key codes.
	KCtl   = C.NSEventModifierFlagControl  // Modifier masks and key codes.
	KFn    = C.NSEventModifierFlagFunction //   "
	KShift = C.NSEventModifierFlagShift    //   "
	KCmd   = C.NSEventModifierFlagCommand  //   "
	KAlt   = C.NSEventModifierFlagOption   //   "

	// Mouse buttons are treated like keys. Values don't
	// conflict with other key codes.
	KML = C.devMouseL // Mouse buttons
	KMM = C.devMouseM //   "
	KMR = C.devMouseR //   "
)

// =============================================================================
// MacOSConsoleWriter

// macosWriter wraps regular logging so the logs appear in the console.
type MacOSConsoleWriter struct{}

// Write provides a writer for regular logging.
func (MacOSConsoleWriter) Write(p []byte) (n int, err error) {
	cstr := C.CString(string(p))
	C.dev_log(cstr)
	C.free(unsafe.Pointer(cstr))
	return len(p), nil
}
