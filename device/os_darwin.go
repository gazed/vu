// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package device

// The OSX (darwin) native layer. This wraps the c functions that wrap the
// objective-c code that calls the OSX window library (where the real
// work is done).

// // The following block is C code and cgo directvies.
//
// #cgo darwin CFLAGS: -x objective-c -fno-common
// #cgo darwin LDFLAGS: -framework Cocoa -framework OpenGL
//
// #include <stdlib.h>
// #include "os_darwin.h"
import "C" // must be located here.

import (
	"runtime"
	"unsafe"
)

// OS specific structure to differentiate it from the other native layers.
type osx struct {
	gsu *C.GSEvent
}

// OSX specific. Otherwise the shell will freeze within seconds of creation.
func init() { runtime.LockOSThread() }

// nativeLayer gets a reference to the native operating system. Each native
// layer implements this factory method. Compiling will leave only the one that
// matches the current platform.
func nativeLayer() native { return &osx{gsu: &C.GSEvent{}} }

// Implement native interface.
func (o *osx) context(r *nrefs) int64      { return int64(C.gs_context(C.long(r.shell))) }
func (o *osx) display() int64              { return int64(C.gs_display_init()) }
func (o *osx) displayDispose(r *nrefs)     { C.gs_display_dispose(C.long(r.display)) }
func (o *osx) shell(r *nrefs) int64        { return int64(C.gs_shell(C.long(r.display))) }
func (o *osx) shellOpen(r *nrefs)          { C.gs_shell_open(C.long(r.display)) }
func (o *osx) shellAlive(r *nrefs) bool    { return uint(C.gs_shell_alive(C.long(r.shell))) == 1 }
func (o *osx) isFullscreen(r *nrefs) bool  { return uint(C.gs_fullscreen(C.long(r.display))) == 1 }
func (o *osx) toggleFullscreen(r *nrefs)   { C.gs_toggle_fullscreen(C.long(r.display)) }
func (o *osx) swapBuffers(r *nrefs)        { C.gs_swap_buffers(C.long(r.context)) }
func (o *osx) setAlphaBufferSize(size int) { C.gs_set_attr_l(C.GS_AlphaSize, C.long(size)) }
func (o *osx) setDepthBufferSize(size int) { C.gs_set_attr_l(C.GS_DepthSize, C.long(size)) }
func (o *osx) setCursorAt(r *nrefs, x, y int) {
	C.gs_set_cursor_location(C.long(r.display), C.long(x), C.long(y))
}
func (o *osx) showCursor(r *nrefs, show bool) {
	trueFalse := 0 // trueFalse needs to be 0 or 1.
	if show {
		trueFalse = 1
	}
	C.gs_show_cursor(C.uchar(trueFalse))
}

// Implement native interface.
func (o *osx) readDispatch(r *nrefs, in *userInput) *userInput {
	o.gsu.event = 0
	o.gsu.mousex = -1
	o.gsu.mousey = -1
	o.gsu.key = 0
	o.gsu.scroll = 0

	// o.gsu.mods retain the modifier key state between calls.
	C.gs_read_dispatch(C.long(r.display), o.gsu)

	// transfer/translate the native event into the input buffer.
	in.id = events[int(o.gsu.event)]
	if in.id != 0 {
		in.button = mouseButtons[int(o.gsu.event)]
		in.key = int(o.gsu.key)
		in.scroll = int(o.gsu.scroll)
	} else {
		in.button, in.key, in.scroll = 0, 0, 0
	}
	in.mods = int(o.gsu.mods) & (controlKeyMask | shiftKeyMask | functionKeyMask | commandKeyMask | altKeyMask)
	in.mouseX = int(o.gsu.mousex)
	in.mouseY = int(o.gsu.mousey)
	return in
}

// Implement native interface.
func (o *osx) size(r *nrefs) (x, y, w, h int) {
	var winx, winy, width, height float32
	C.gs_size(C.long(r.shell), (*C.float)(&winx), (*C.float)(&winy), (*C.float)(&width), (*C.float)(&height))
	return int(winx), int(winy), int(width), int(height)
}

// Implement native interface.
func (o *osx) setSize(x, y, width, height int) {
	C.gs_set_attr_l(C.GS_ShellX, C.long(x))
	C.gs_set_attr_l(C.GS_ShellY, C.long(y))
	C.gs_set_attr_l(C.GS_ShellWidth, C.long(width))
	C.gs_set_attr_l(C.GS_ShellHeight, C.long(height))
}

// Implement native interface.
func (o *osx) setTitle(title string) {
	cstr := C.CString(title)
	defer C.free(unsafe.Pointer(cstr))
	C.gs_set_attr_s(C.GS_AppName, cstr)
}

// Implement native interface: nrefs unused, needed by other platforms.
func (o *osx) copyClip(r *nrefs) string {
	if cstr := C.gs_clip_copy(); cstr != nil {
		str := C.GoString(cstr)      // make a Go copy.
		C.free(unsafe.Pointer(cstr)) // free the C copy.
		return str
	}
	return ""
}

// Implement native interface: nrefs unused, needed by other platforms.
func (o *osx) pasteClip(r *nrefs, s string) {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	C.gs_clip_paste(cstr)
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
	C.GS_LeftMouseDown:  mouseLeft,
	C.GS_RightMouseDown: mouseRight,
	C.GS_OtherMouseDown: mouseMiddle,
	C.GS_LeftMouseUp:    mouseLeft,
	C.GS_RightMouseUp:   mouseRight,
	C.GS_OtherMouseUp:   mouseMiddle,
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
	key0              = 0x1D // kVK_ANSI_0
	key1              = 0x12 // kVK_ANSI_1
	key2              = 0x13 // kVK_ANSI_2
	key3              = 0x14 // kVK_ANSI_3
	key4              = 0x15 // kVK_ANSI_4
	key5              = 0x17 // kVK_ANSI_5
	key6              = 0x16 // kVK_ANSI_6
	key7              = 0x1A // kVK_ANSI_7
	key8              = 0x1C // kVK_ANSI_8
	key9              = 0x19 // kVK_ANSI_9
	keyA              = 0x00 // kVK_ANSI_A
	keyB              = 0x0B // kVK_ANSI_B
	keyC              = 0x08 // kVK_ANSI_C
	keyD              = 0x02 // kVK_ANSI_D
	keyE              = 0x0E // kVK_ANSI_E
	keyF              = 0x03 // kVK_ANSI_F
	keyG              = 0x05 // kVK_ANSI_G
	keyH              = 0x04 // kVK_ANSI_H
	keyI              = 0x22 // kVK_ANSI_I
	keyJ              = 0x26 // kVK_ANSI_J
	keyK              = 0x28 // kVK_ANSI_K
	keyL              = 0x25 // kVK_ANSI_L
	keyM              = 0x2E // kVK_ANSI_M
	keyN              = 0x2D // kVK_ANSI_N
	keyO              = 0x1F // kVK_ANSI_O
	keyP              = 0x23 // kVK_ANSI_P
	keyQ              = 0x0C // kVK_ANSI_Q
	keyR              = 0x0F // kVK_ANSI_R
	keyS              = 0x01 // kVK_ANSI_S
	keyT              = 0x11 // kVK_ANSI_T
	keyU              = 0x20 // kVK_ANSI_U
	keyV              = 0x09 // kVK_ANSI_V
	keyW              = 0x0D // kVK_ANSI_W
	keyX              = 0x07 // kVK_ANSI_X
	keyY              = 0x10 // kVK_ANSI_Y
	keyZ              = 0x06 // kVK_ANSI_Z
	keyF1             = 0x7A // kVK_F1
	keyF2             = 0x78 // kVK_F2
	keyF3             = 0x63 // kVK_F3
	keyF4             = 0x76 // kVK_F4
	keyF5             = 0x60 // kVK_F5
	keyF6             = 0x61 // kVK_F6
	keyF7             = 0x62 // kVK_F7
	keyF8             = 0x64 // kVK_F8
	keyF9             = 0x65 // kVK_F9
	keyF10            = 0x6D // kVK_F10
	keyF11            = 0x67 // kVK_F11
	keyF12            = 0x6F // kVK_F12
	keyF13            = 0x69 // kVK_F13
	keyF14            = 0x6B // kVK_F14
	keyF15            = 0x71 // kVK_F15
	keyF16            = 0x6A // kVK_F16
	keyF17            = 0x40 // kVK_F17
	keyF18            = 0x4F // kVK_F18
	keyF19            = 0x50 // kVK_F19
	keyF20            = 0x5A // kVK_F20
	keyKeypad0        = 0x52 // kVK_ANSI_Keypad0
	keyKeypad1        = 0x53 // kVK_ANSI_Keypad1
	keyKeypad2        = 0x54 // kVK_ANSI_Keypad2
	keyKeypad3        = 0x55 // kVK_ANSI_Keypad3
	keyKeypad4        = 0x56 // kVK_ANSI_Keypad4
	keyKeypad5        = 0x57 // kVK_ANSI_Keypad5
	keyKeypad6        = 0x58 // kVK_ANSI_Keypad6
	keyKeypad7        = 0x59 // kVK_ANSI_Keypad7
	keyKeypad8        = 0x5B // kVK_ANSI_Keypad8
	keyKeypad9        = 0x5C // kVK_ANSI_Keypad9
	keyKeypadDecimal  = 0x41 // kVK_ANSI_KeypadDecimal
	keyKeypadMultiply = 0x43 // kVK_ANSI_KeypadMultiply
	keyKeypadPlus     = 0x45 // kVK_ANSI_KeypadPlus
	keyKeypadClear    = 0x47 // kVK_ANSI_KeypadClear
	keyKeypadDivide   = 0x4B // kVK_ANSI_KeypadDivide
	keyKeypadEnter    = 0x4C // kVK_ANSI_KeypadEnter
	keyKeypadMinus    = 0x4E // kVK_ANSI_KeypadMinus
	keyKeypadEquals   = 0x51 // kVK_ANSI_KeypadEquals
	keyEqual          = 0x18 // kVK_ANSI_Equal
	keyMinus          = 0x1B // kVK_ANSI_Minus
	keyLeftBracket    = 0x21 // kVK_ANSI_LeftBracket
	keyRightBracket   = 0x1E // kVK_ANSI_RightBracket
	keyQuote          = 0x27 // kVK_ANSI_Quote
	keySemicolon      = 0x29 // kVK_ANSI_Semicolon
	keyBackslash      = 0x2A // kVK_ANSI_Backslash
	keyGrave          = 0x32 // kVK_ANSI_Grave
	keySlash          = 0x2C // kVK_ANSI_Slash
	keyComma          = 0x2B // kVK_ANSI_Comma
	keyPeriod         = 0x2F // kVK_ANSI_Period
	keyReturn         = 0x24 // kVK_Return
	keyTab            = 0x30 // kVK_Tab
	keySpace          = 0x31 // kVK_Space
	keyDelete         = 0x33 // kVK_Delete
	keyForwardDelete  = 0x75 // kVK_ForwardDelete
	keyEscape         = 0x35 // kVK_Escape
	keyHome           = 0x73 // kVK_Home
	keyPageUp         = 0x74 // kVK_PageUp
	keyPageDown       = 0x79 // kVK_PageDown
	keyLeftArrow      = 0x7B // kVK_LeftArrow
	keyRightArrow     = 0x7C // kVK_RightArrow
	keyDownArrow      = 0x7D // kVK_DownArrow
	keyUpArrow        = 0x7E // kVK_UpArrow
	keyEnd            = 0x77 // kVK_End
	mouseLeft         = 0xA0 // tack on unique values for mouse buttons
	mouseMiddle       = 0xA1
	mouseRight        = 0xA2
)
