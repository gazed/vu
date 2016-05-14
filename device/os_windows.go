// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package device

// The microsoft (windows) native layer. This wraps the c functions that
// wrap the microsoft API's (where the real work is done).

// // This is C code and cgo directvies.
//
// #cgo windows CFLAGS: -m64
// #cgo windows,!dx LDFLAGS: -lopengl32 -lgdi32
// #cgo windows,dx LDFLAGS: -ld3d11
// #cgo windows,dx CXXFLAGS: -std=c++11
//
// #include "os_windows.h"
import "C" // must be located here.

import (
	"runtime"
	"unsafe"
)

// OS specific structure to differentiate it from the other native layers.
// Two input structures are continually reused each time rather than allocating
// a new input structure on each readAndDispatch.
type win struct {
	gsu *C.GSEvent
}

// OpenGL related, see: https://code.google.com/p/go-wiki/wiki/LockOSThread
func init() { runtime.LockOSThread() }

// nativeLayer gets a reference to the native operating system.  Each native
// layer implements this factory method. Compiling will leave only the one that
// matches the current platform.
func nativeLayer() native { return &win{gsu: &C.GSEvent{}} }

// Implement native interface.
func (w *win) context(r *nrefs) int64 {
	return int64(C.gs_context((*C.longlong)(&(r.display)), (*C.longlong)(&(r.shell))))
}
func (w *win) display() int64              { return int64(C.gs_display_init()) }
func (w *win) displayDispose(r *nrefs)     { C.gs_display_dispose(C.long(r.display)) }
func (w *win) shell(r *nrefs) int64        { return int64(C.gs_shell(C.long(r.display))) }
func (w *win) shellOpen(r *nrefs)          { C.gs_shell_open(C.long(r.display)) }
func (w *win) shellAlive(r *nrefs) bool    { return uint(C.gs_shell_alive(C.long(r.shell))) == 1 }
func (w *win) isFullscreen(r *nrefs) bool  { return uint(C.gs_fullscreen(C.long(r.display))) == 1 }
func (w *win) toggleFullscreen(r *nrefs)   { C.gs_toggle_fullscreen(C.long(r.display)) }
func (w *win) swapBuffers(r *nrefs)        { C.gs_swap_buffers(C.long(r.shell)) }
func (w *win) setAlphaBufferSize(size int) { C.gs_set_attr_l(C.GS_AlphaSize, C.long(size)) }
func (w *win) setDepthBufferSize(size int) { C.gs_set_attr_l(C.GS_DepthSize, C.long(size)) }
func (w *win) setCursorAt(r *nrefs, x, y int) {
	C.gs_set_cursor_location(C.long(r.display), C.long(x), C.long(y))
}
func (w *win) showCursor(r *nrefs, show bool) {
	tf1 := 0
	if show {
		tf1 = 1
	}
	C.gs_show_cursor(C.long(r.display), C.uchar(tf1))
}

// Implement native interface.
func (w *win) readDispatch(r *nrefs, in *userInput) *userInput {
	w.gsu.event = 0
	w.gsu.mousex = -1
	w.gsu.mousey = -1
	w.gsu.key = 0
	w.gsu.mods = 0
	w.gsu.scroll = 0
	C.gs_read_dispatch(C.long(r.display), w.gsu)

	// transfer/translate the native event into the input buffer.
	in.id = events[int(w.gsu.event)]
	if in.id != 0 {
		in.button = mouseButtons[int(w.gsu.event)]
		in.key = int(w.gsu.key)
		in.scroll = int(w.gsu.scroll)
	} else {
		in.button, in.key, in.scroll = 0, 0, 0
	}
	in.mods = int(w.gsu.mods)
	in.mouseX = int(w.gsu.mousex)
	in.mouseY = int(w.gsu.mousey)
	return in
}

// Implement native interface.
func (w *win) size(r *nrefs) (x int, y int, wx int, hy int) {
	var winx, winy, width, height int32
	C.gs_size(C.long(r.display), (*C.long)(&winx), (*C.long)(&winy), (*C.long)(&width), (*C.long)(&height))
	return int(winx), int(winy), int(width), int(height)
}

// Implement native interface.
func (w *win) setSize(x, y, width, height int) {
	C.gs_set_attr_l(C.GS_ShellX, C.long(x))
	C.gs_set_attr_l(C.GS_ShellY, C.long(y))
	C.gs_set_attr_l(C.GS_ShellWidth, C.long(width))
	C.gs_set_attr_l(C.GS_ShellHeight, C.long(height))
}

// Implement native interface.
func (w *win) setTitle(title string) {
	cstr := C.CString(title)
	defer C.free(unsafe.Pointer(cstr))
	C.gs_set_attr_s(C.GS_AppName, cstr)
}

// Implement native interface.
func (w *win) copyClip(r *nrefs) string {
	if cstr := C.gs_clip_copy(C.long(r.display)); cstr != nil {
		str := C.GoString(cstr)      // make a Go copy.
		C.free(unsafe.Pointer(cstr)) // free the C copy.
		return str
	}
	return ""
}

// Implement native interface.
func (w *win) pasteClip(r *nrefs, s string) {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	C.gs_clip_paste(C.long(r.display), cstr)
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
	C.GS_SysKeyUp:          releasedKey,
	C.GS_ScrollWheel:       scrolled,
	C.GS_WindowResized:     resizedShell,
	C.GS_WindowMoved:       movedShell,
	C.GS_WindowIconified:   iconifiedShell,
	C.GS_WindowUniconified: uniconifiedShell,
	C.GS_WindowActive:      activatedShell,
	C.GS_WindowInactive:    deactivatedShell,
}

// Also map the mice buttons into left and right.
var mouseButtons = map[int]int{
	C.GS_LeftMouseDown:  mouseLeft,
	C.GS_RightMouseDown: mouseRight,
	C.GS_OtherMouseDown: mouseMiddle,
	C.GS_LeftMouseUp:    mouseLeft,
	C.GS_RightMouseUp:   mouseRight,
	C.GS_OtherMouseUp:   mouseMiddle,
}

// Expose the underlying Win key modifier masks.
// Leave the ALT and CMD keys to the OS's.
const (
	shiftKeyMask    = C.GS_ShiftKeyMask
	controlKeyMask  = C.GS_ControlKeyMask
	functionKeyMask = C.GS_FunctionKeyMask
	commandKeyMask  = C.GS_CommandKeyMask
	altKeyMask      = C.GS_AlternateKeyMask
)

// Expose the underlying Win key codes as generic code.
// Each native layer is expected to support the generic codes.
//
// Windows virtual key codes.
// http://msdn.microsoft.com/en-ca/library/windows/desktop/dd375731(v=vs.85).aspx
const (
	key0              = 0x30 // 0 key
	key1              = 0x31 // 1 key
	key2              = 0x32 // 2 key
	key3              = 0x33 // 3 key
	key4              = 0x34 // 4 key
	key5              = 0x35 // 5 key
	key6              = 0x36 // 6 key
	key7              = 0x37 // 7 key
	key8              = 0x38 // 8 key
	key9              = 0x39 // 9 key
	keyA              = 0x41 // A key
	keyB              = 0x42 // B key
	keyC              = 0x43 // C key
	keyD              = 0x44 // D key
	keyE              = 0x45 // E key
	keyF              = 0x46 // F key
	keyG              = 0x47 // G key
	keyH              = 0x48 // H key
	keyI              = 0x49 // I key
	keyJ              = 0x4A // J key
	keyK              = 0x4B // K key
	keyL              = 0x4C // L key
	keyM              = 0x4D // M key
	keyN              = 0x4E // N key
	keyO              = 0x4F // O key
	keyP              = 0x50 // P key
	keyQ              = 0x51 // Q key
	keyR              = 0x52 // R key
	keyS              = 0x53 // S key
	keyT              = 0x54 // T key
	keyU              = 0x55 // U key
	keyV              = 0x56 // V key
	keyW              = 0x57 // W key
	keyX              = 0x58 // X key
	keyY              = 0x59 // Y key
	keyZ              = 0x5A // Z key
	keyF1             = 0x70 // VK_F1            F1 key
	keyF2             = 0x71 // VK_F2            F2 key
	keyF3             = 0x72 // VK_F3            F3 key
	keyF4             = 0x73 // VK_F4            F4 key
	keyF5             = 0x74 // VK_F5            F5 key
	keyF6             = 0x75 // VK_F6            F6 key
	keyF7             = 0x76 // VK_F7            F7 key
	keyF8             = 0x77 // VK_F8            F8 key
	keyF9             = 0x78 // VK_F9            F9 key
	keyF10            = 0x79 // VK_F10           F10 key  ---- on osx-kb
	keyF11            = 0x7A // VK_F11           F11 key
	keyF12            = 0x7B // VK_F12           F12 key
	keyF13            = 0x7C // VK_F13           F13 key
	keyF14            = 0x2C // VK_F14 0x7D      F14 key  0x2C on osx-kb
	keyF15            = 0x91 // VK_F15 0x7E      F15 key  0x91
	keyF16            = 0x13 // VK_F16 0x7F      F16 key  0x13
	keyF17            = 0x80 // VK_F17           F17 key
	keyF18            = 0x81 // VK_F18           F18 key
	keyF19            = 0x82 // VK_F19           F19 key
	keyF20            = 0x83 // VK_F20           F20 key
	keyKeypad0        = 0x60 // VK_NUMPAD0  0x60 Numeric keypad 0 key :: VK_INSERT  0x20 on osx-kb
	keyKeypad1        = 0x61 // VK_NUMPAD1  0x61 Numeric keypad 1 key :: VK_END     0x23 on osx-kb
	keyKeypad2        = 0x62 // VK_NUMPAD2  0x62 Numeric keypad 2 key :: VK_DOWN    0x28 on osx-kb
	keyKeypad3        = 0x63 // VK_NUMPAD3  0x63 Numeric keypad 3 key :: VK_NEXT    0x22 on osx-kb
	keyKeypad4        = 0x64 // VK_NUMPAD4  0x64 Numeric keypad 4 key :: VK_LEFT    0x25 on osx-kb
	keyKeypad5        = 0x65 // VK_NUMPAD5  0x65 Numeric keypad 5 key :: VK_CLEAR   0x0C on osx-kb
	keyKeypad6        = 0x66 // VK_NUMPAD6  0x66 Numeric keypad 6 key :: VK_RIGHT   0x27 on osx-kb
	keyKeypad7        = 0x67 // VK_NUMPAD7  0x67 Numeric keypad 7 key :: VK_HOME    0x26 on osx-kb
	keyKeypad8        = 0x68 // VK_NUMPAD8  0x68 Numeric keypad 8 key :: VK_UP      0x21 on osx-kb
	keyKeypad9        = 0x69 // VK_NUMPAD9  0x69 Numeric keypad 9 key :: VK_PRIOR
	keyKeypadDecimal  = 0x6E // VK_DECIMAL       Decimal key :: VK_DELETE
	keyKeypadMultiply = 0x6A // VK_MULTIPLY      Multiply key
	keyKeypadPlus     = 0x6B // VK_ADD           Add key
	keyKeypadClear    = 0x90 // VK_CLEAR    0x0C CLEAR key :: VK_OEM_CLEAR 0xFE     0x90 on osx-kb
	keyKeypadDivide   = 0x6F // VK_DIVIDE        Divide key
	keyKeypadEnter    = 0x2B // VK_EXECUTE                            :: VK_ENTER on osx-kb
	keyKeypadMinus    = 0x6D // VK_SUBTRACT      Subtract key
	keyKeypadEquals   = 0xE2 //
	keyEqual          = 0xBB //
	keyMinus          = 0xBD // VK_OEM_MINUS     For any country/region, the '-' key // VK_SEPARATOR 0x6C Separator key
	keyLeftBracket    = 0xDB // VK_OEM_4         misc characters; varys: US standard keyboard, the '[{' key
	keyRightBracket   = 0xDD // VK_OEM_6         misc characters; varys: US standard keyboard, the ']}' key
	keyQuote          = 0xC0 // VK_OEM_7         misc characters; varys: US standard keyboard, the 'single/double-quote' key
	keySemicolon      = 0xBA // VK_OEM_1         misc characters; varys: US standard keyboard, the ';:' key
	keyBackslash      = 0xDE // VK_OEM_5         misc characters; varys: US standard keyboard, the '/?' key
	keyGrave          = 0xDF // VK_OEM_3         misc characters; varys: US standard keyboard, the '`~' key
	keySlash          = 0xBF //
	keyComma          = 0xBC // VK_OEM_COMMA     For any country/region, the ',' key
	keyPeriod         = 0xBE // VK_OEM_PERIOD    For any country/region, the '.' key
	keyReturn         = 0x0D // VK_RETURN        ENTER key
	keyTab            = 0x09 // VK_TAB           TAB key
	keySpace          = 0x20 // VK_SPACE         SPACEBAR
	keyDelete         = 0x08 // VK_BACK          BACKSPACE key
	keyForwardDelete  = 0x2E // VK_DELETE        DEL key
	keyEscape         = 0x1B // VK_ESCAPE        ESC key
	keyHome           = 0x24 // VK_HOME          HOME key
	keyPageUp         = 0x21 // VK_PRIOR         PAGE UP key
	keyPageDown       = 0x22 // VK_NEXT          PAGE DOWN key
	keyLeftArrow      = 0x25 // VK_LEFT          LEFT ARROW key
	keyRightArrow     = 0x27 // VK_RIGHT         RIGHT ARROW key
	keyDownArrow      = 0x28 // VK_DOWN          DOWN ARROW key
	keyUpArrow        = 0x26 // VK_UP            UP ARROW key
	keyEnd            = 0x23 // VK_END           END key
	mouseLeft         = 0x01 // VK_LBUTTON Left mouse button (tack on unique values for mouse buttons)
	mouseMiddle       = 0x04 // VK_MBUTTON Middle mouse button (three-button mouse)
	mouseRight        = 0x02 // VK_RBUTTON Right mouse button
)
