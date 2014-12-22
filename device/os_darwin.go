// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package device

// The OSX (darwin) native layer. This wraps the c functions that wrap the
// objective-c code that calls the Osx windowing library (where the real
// work is done).

// // The following block is C code and cgo directvies.
//
// #cgo darwin CFLAGS: -x objective-c -fno-common
// #cgo darwin LDFLAGS: -framework Cocoa -framework OpenGL -framework IOKit
//
// #include <stdlib.h>
// #include "os_darwin.h"
import "C" // must be located here.

import (
	"runtime"
	"unsafe"
)

// OS specific structure to differentiate it from the other native layers.
// Two input structures are continually reused each time rather than allocating
// a osx input structure on each readAndDispatch.
type osx struct {
	in  *userInput // Reusable input event buffer.
	in1 *userInput // Alternate reusable input event buffer.
}

// OSX specific. Otherwise the shell will freeze within seconds of creation.
func init() { runtime.LockOSThread() }

// nativeLayer gets a reference to the native operating system. Each native
// layer implements this factory method. Compiling will leave only the one that
// matches the current platform.
func nativeLayer() native { return &osx{&userInput{}, &userInput{}} }

// Implements native interface.
func (o *osx) context(r *nrefs) int64      { return int64(C.gs_context(C.long(r.shell))) }
func (o *osx) display() int64              { return int64(C.gs_display_init()) }
func (o *osx) displayDispose(r *nrefs)     { C.gs_display_dispose(C.long(r.display)) }
func (o *osx) shell(r *nrefs) int64        { return int64(C.gs_shell(C.long(r.display))) }
func (o *osx) shellOpen(r *nrefs)          { C.gs_shell_open(C.long(r.display)) }
func (o *osx) shellAlive(r *nrefs) bool    { return uint(C.gs_shell_alive(C.long(r.shell))) == 1 }
func (o *osx) swapBuffers(r *nrefs)        { C.gs_swap_buffers(C.long(r.context)) }
func (o *osx) setAlphaBufferSize(size int) { C.gs_set_attr_l(C.GS_AlphaSize, C.long(size)) }
func (o *osx) setDepthBufferSize(size int) { C.gs_set_attr_l(C.GS_DepthSize, C.long(size)) }
func (o *osx) setCursorAt(r *nrefs, x, y int) {
	C.gs_set_cursor_location(C.long(r.display), C.long(x), C.long(y))
}
func (o *osx) showCursor(r *nrefs, show bool) {
	tf1 := 0
	if show {
		tf1 = 1
	}
	C.gs_show_cursor(C.uchar(tf1))
}

// See native interface.
func (o *osx) readDispatch(r *nrefs) *userInput {
	gsu := &C.GSEvent{0, -1, -1, 0, 0, 0}
	C.gs_read_dispatch(C.long(r.display), gsu)
	o.in, o.in1 = o.in1, o.in

	// transfer/translate the native event into the input buffer.
	in := o.in
	in.id = events[int(gsu.event)]
	if in.id != 0 {
		in.button = mouseButtons[int(gsu.event)]
		in.key = int(gsu.key)
		in.scroll = int(gsu.scroll)
	} else {
		in.button, in.key, in.scroll = 0, 0, 0
	}
	in.mods = int(gsu.mods) & (controlKeyMask | shiftKeyMask | functionKeyMask | commandKeyMask | altKeyMask)
	in.mouseX = int(gsu.mousex)
	in.mouseY = int(gsu.mousey)
	return in
}

// See native interface.
func (o *osx) size(r *nrefs) (x, y, w, h int) {
	var winx, winy, width, height float32
	C.gs_size(C.long(r.shell), (*C.float)(&winx), (*C.float)(&winy), (*C.float)(&width), (*C.float)(&height))
	return int(winx), int(winy), int(width), int(height)
}

// See native interface.
func (o *osx) setSize(x, y, width, height int) {
	C.gs_set_attr_l(C.GS_ShellX, C.long(x))
	C.gs_set_attr_l(C.GS_ShellY, C.long(y))
	C.gs_set_attr_l(C.GS_ShellWidth, C.long(width))
	C.gs_set_attr_l(C.GS_ShellHeight, C.long(height))
}

// See native interface.
func (o *osx) setTitle(title string) {
	cstr := C.CString(title)
	defer C.free(unsafe.Pointer(cstr))
	C.gs_set_attr_s(C.GS_AppName, cstr)
}

// Transform os specific events to user events.
var events = map[int]int{
	C.GS_LeftMouseDown:     clickedMouse,
	C.GS_RightMouseDown:    clickedMouse,
	C.GS_OtherMouseDown:    clickedMouse,
	C.GS_LeftMouseUp:       releasedMouse,
	C.GS_RightMouseUp:      releasedMouse,
	C.GS_OtherMouseUp:      releasedMouse,
	C.GS_MouseMoved:        movedMouse,
	C.GS_KeyDown:           pressedKey,
	C.GS_KeyUp:             releasedKey,
	C.GS_ScrollWheel:       scrolled,
	C.GS_WindowResized:     resizedShell,
	C.GS_WindowMoved:       movedShell,
	C.GS_WindowIconified:   iconifiedShell,
	C.GS_WindowUniconified: uniconifiedShell,
	C.GS_WindowActive:      activatedShell,
	C.GS_WindowInactive:    deactivatedShell,
}

// Map the mice buttons into left and right.
var mouseButtons = map[int]int{
	C.GS_LeftMouseDown:  mouse_Left,
	C.GS_RightMouseDown: mouse_Right,
	C.GS_OtherMouseDown: mouse_Middle,
	C.GS_LeftMouseUp:    mouse_Left,
	C.GS_RightMouseUp:   mouse_Right,
	C.GS_OtherMouseUp:   mouse_Middle,
}

// Expose the underlying OSX key modifier masks.
const (
	shiftKeyMask    = C.GS_ShiftKeyMask
	controlKeyMask  = C.GS_ControlKeyMask
	functionKeyMask = C.GS_FunctionKeyMask
	commandKeyMask  = C.GS_CommandKeyMask
	altKeyMask      = C.GS_AlternateKeyMask
)

// Expose the underlying OSX key codes as generic code.
// Each native layer is expected to support the generic codes.
const (
	key_0              = 0x1D // kVK_ANSI_0
	key_1              = 0x12 // kVK_ANSI_1
	key_2              = 0x13 // kVK_ANSI_2
	key_3              = 0x14 // kVK_ANSI_3
	key_4              = 0x15 // kVK_ANSI_4
	key_5              = 0x17 // kVK_ANSI_5
	key_6              = 0x16 // kVK_ANSI_6
	key_7              = 0x1A // kVK_ANSI_7
	key_8              = 0x1C // kVK_ANSI_8
	key_9              = 0x19 // kVK_ANSI_9
	key_A              = 0x00 // kVK_ANSI_A
	key_B              = 0x0B // kVK_ANSI_B
	key_C              = 0x08 // kVK_ANSI_C
	key_D              = 0x02 // kVK_ANSI_D
	key_E              = 0x0E // kVK_ANSI_E
	key_F              = 0x03 // kVK_ANSI_F
	key_G              = 0x05 // kVK_ANSI_G
	key_H              = 0x04 // kVK_ANSI_H
	key_I              = 0x22 // kVK_ANSI_I
	key_J              = 0x26 // kVK_ANSI_J
	key_K              = 0x28 // kVK_ANSI_K
	key_L              = 0x25 // kVK_ANSI_L
	key_M              = 0x2E // kVK_ANSI_M
	key_N              = 0x2D // kVK_ANSI_N
	key_O              = 0x1F // kVK_ANSI_O
	key_P              = 0x23 // kVK_ANSI_P
	key_Q              = 0x0C // kVK_ANSI_Q
	key_R              = 0x0F // kVK_ANSI_R
	key_S              = 0x01 // kVK_ANSI_S
	key_T              = 0x11 // kVK_ANSI_T
	key_U              = 0x20 // kVK_ANSI_U
	key_V              = 0x09 // kVK_ANSI_V
	key_W              = 0x0D // kVK_ANSI_W
	key_X              = 0x07 // kVK_ANSI_X
	key_Y              = 0x10 // kVK_ANSI_Y
	key_Z              = 0x06 // kVK_ANSI_Z
	key_F1             = 0x7A // kVK_F1
	key_F2             = 0x78 // kVK_F2
	key_F3             = 0x63 // kVK_F3
	key_F4             = 0x76 // kVK_F4
	key_F5             = 0x60 // kVK_F5
	key_F6             = 0x61 // kVK_F6
	key_F7             = 0x62 // kVK_F7
	key_F8             = 0x64 // kVK_F8
	key_F9             = 0x65 // kVK_F9
	key_F10            = 0x6D // kVK_F10
	key_F11            = 0x67 // kVK_F11
	key_F12            = 0x6F // kVK_F12
	key_F13            = 0x69 // kVK_F13
	key_F14            = 0x6B // kVK_F14
	key_F15            = 0x71 // kVK_F15
	key_F16            = 0x6A // kVK_F16
	key_F17            = 0x40 // kVK_F17
	key_F18            = 0x4F // kVK_F18
	key_F19            = 0x50 // kVK_F19
	key_F20            = 0x5A // kVK_F20
	key_Keypad0        = 0x52 // kVK_ANSI_Keypad0
	key_Keypad1        = 0x53 // kVK_ANSI_Keypad1
	key_Keypad2        = 0x54 // kVK_ANSI_Keypad2
	key_Keypad3        = 0x55 // kVK_ANSI_Keypad3
	key_Keypad4        = 0x56 // kVK_ANSI_Keypad4
	key_Keypad5        = 0x57 // kVK_ANSI_Keypad5
	key_Keypad6        = 0x58 // kVK_ANSI_Keypad6
	key_Keypad7        = 0x59 // kVK_ANSI_Keypad7
	key_Keypad8        = 0x5B // kVK_ANSI_Keypad8
	key_Keypad9        = 0x5C // kVK_ANSI_Keypad9
	key_KeypadDecimal  = 0x41 // kVK_ANSI_KeypadDecimal
	key_KeypadMultiply = 0x43 // kVK_ANSI_KeypadMultiply
	key_KeypadPlus     = 0x45 // kVK_ANSI_KeypadPlus
	key_KeypadClear    = 0x47 // kVK_ANSI_KeypadClear
	key_KeypadDivide   = 0x4B // kVK_ANSI_KeypadDivide
	key_KeypadEnter    = 0x4C // kVK_ANSI_KeypadEnter
	key_KeypadMinus    = 0x4E // kVK_ANSI_KeypadMinus
	key_KeypadEquals   = 0x51 // kVK_ANSI_KeypadEquals
	key_Equal          = 0x18 // kVK_ANSI_Equal
	key_Minus          = 0x1B // kVK_ANSI_Minus
	key_LeftBracket    = 0x21 // kVK_ANSI_LeftBracket
	key_RightBracket   = 0x1E // kVK_ANSI_RightBracket
	key_Quote          = 0x27 // kVK_ANSI_Quote
	key_Semicolon      = 0x29 // kVK_ANSI_Semicolon
	key_Backslash      = 0x2A // kVK_ANSI_Backslash
	key_Grave          = 0x32 // kVK_ANSI_Grave
	key_Slash          = 0x2C // kVK_ANSI_Slash
	key_Comma          = 0x2B // kVK_ANSI_Comma
	key_Period         = 0x2F // kVK_ANSI_Period
	key_Return         = 0x24 // kVK_Return
	key_Tab            = 0x30 // kVK_Tab
	key_Space          = 0x31 // kVK_Space
	key_Delete         = 0x33 // kVK_Delete
	key_ForwardDelete  = 0x75 // kVK_ForwardDelete
	key_Escape         = 0x35 // kVK_Escape
	key_Home           = 0x73 // kVK_Home
	key_PageUp         = 0x74 // kVK_PageUp
	key_PageDown       = 0x79 // kVK_PageDown
	key_LeftArrow      = 0x7B // kVK_LeftArrow
	key_RightArrow     = 0x7C // kVK_RightArrow
	key_DownArrow      = 0x7D // kVK_DownArrow
	key_UpArrow        = 0x7E // kVK_UpArrow
	key_End            = 0x77 // kVK_End
	mouse_Left         = 0xA0 // tack on unique values for mouse buttons
	mouse_Middle       = 0xA1
	mouse_Right        = 0xA2
)
