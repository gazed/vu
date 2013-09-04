// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package device

// The microsoft (windows) native layer.  This wraps the c functions that
// wrap the microsoft API's (where the real work is done).

// // The following block is C code and cgo directvies.
//
// #cgo windows CFLAGS: -m64
// #cgo windows LDFLAGS: -lopengl32 -lgdi32
//
// #include "os_windows.h"
import "C" // must be located here.

import (
	"unsafe"
)

// OS specific structure to differentiate it from the other native layers.
type win struct{}

// nativeLayer gets a reference to the native operating system.  Each native
// layer implements this factory method.  Compiling will leave only the one that
// matches the current platform.
func nativeLayer() native { return &win{} }

// Implements native interface.
func (w win) context(r *nrefs) int64 {
	return int64(C.gs_context((*C.longlong)(&(r.display)), (*C.longlong)(&(r.shell))))
}
func (w win) display() int64              { return int64(C.gs_display_init()) }
func (w win) displayDispose(r *nrefs)     { C.gs_display_dispose(C.long(r.display)) }
func (w win) shell(r *nrefs) int64        { return int64(C.gs_shell(C.long(r.display))) }
func (w win) shellOpen(r *nrefs)          { C.gs_shell_open(C.long(r.display)) }
func (w win) shellAlive(r *nrefs) bool    { return uint(C.gs_shell_alive(C.long(r.shell))) == 1 }
func (w win) swapBuffers(r *nrefs)        { C.gs_swap_buffers(C.long(r.shell)) }
func (w win) setAlphaBufferSize(size int) { C.gs_set_attr_l(C.GS_AlphaSize, C.long(size)) }
func (w win) setDepthBufferSize(size int) { C.gs_set_attr_l(C.GS_DepthSize, C.long(size)) }
func (w win) setCursorAt(r *nrefs, x, y int) {
	C.gs_set_cursor_location(C.long(r.display), C.long(x), C.long(y))
}
func (w win) showCursor(r *nrefs, show bool) {
	tf1 := 0
	if show {
		tf1 = 1
	}
	C.gs_show_cursor(C.long(r.display), C.uchar(tf1))
}

// See native interface.
func (w win) readDispatch(r *nrefs) *userEvent {
	gsu := &C.GSEvent{0, -1, -1, 0, 0, 0}
	C.gs_read_dispatch(C.long(r.display), gsu)
	ue := &userEvent{}
	ue.id = events[int(gsu.event)]
	if ue.id != 0 {
		ue.button = mouseButtons[int(gsu.event)]
		ue.key = int(gsu.key)
		ue.mods = int(gsu.mods)
		ue.scroll = int(gsu.scroll)
	}
	ue.mouseX = int(gsu.mousex)
	ue.mouseY = int(gsu.mousey)
	return ue
}

// See native interface.
func (w win) size(r *nrefs) (x int, y int, wx int, hy int) {
	var winx, winy, width, height int32
	C.gs_size(C.long(r.display), (*C.long)(&winx), (*C.long)(&winy), (*C.long)(&width), (*C.long)(&height))
	return int(winx), int(winy), int(width), int(height)
}

// See native interface.
func (w win) setSize(x, y, width, height int) {
	C.gs_set_attr_l(C.GS_ShellX, C.long(x))
	C.gs_set_attr_l(C.GS_ShellY, C.long(y))
	C.gs_set_attr_l(C.GS_ShellWidth, C.long(width))
	C.gs_set_attr_l(C.GS_ShellHeight, C.long(height))
}

// See native interface.
func (w win) setTitle(title string) {
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

// Also map the mice buttons into left and right.
var mouseButtons = map[int]int{
	C.GS_LeftMouseDown:  mouse_Left,
	C.GS_RightMouseDown: mouse_Right,
	C.GS_OtherMouseDown: mouse_Middle,
	C.GS_LeftMouseUp:    mouse_Left,
	C.GS_RightMouseUp:   mouse_Right,
	C.GS_OtherMouseUp:   mouse_Middle,
}

// Expose the underlying OSX key modifier masks.
// Leave the ALT and CMD keys to the OS's.
const (
	shiftKeyMask    = C.GS_ShiftKeyMask
	controlKeyMask  = C.GS_ControlKeyMask
	functionKeyMask = C.GS_FunctionKeyMask
)

// Expose the underlying Win key codes as generic code.
// Each native layer is expected to support the generic codes.
//
// Windows virtual key codes.
// http://msdn.microsoft.com/en-ca/library/windows/desktop/dd375731(v=vs.85).aspx
const (
	key_0              = 0x30 // 0 key
	key_1              = 0x31 // 1 key
	key_2              = 0x32 // 2 key
	key_3              = 0x33 // 3 key
	key_4              = 0x34 // 4 key
	key_5              = 0x35 // 5 key
	key_6              = 0x36 // 6 key
	key_7              = 0x37 // 7 key
	key_8              = 0x38 // 8 key
	key_9              = 0x39 // 9 key
	key_A              = 0x41 // A key
	key_B              = 0x42 // B key
	key_C              = 0x43 // C key
	key_D              = 0x44 // D key
	key_E              = 0x45 // E key
	key_F              = 0x46 // F key
	key_G              = 0x47 // G key
	key_H              = 0x48 // H key
	key_I              = 0x49 // I key
	key_J              = 0x4A // J key
	key_K              = 0x4B // K key
	key_L              = 0x4C // L key
	key_M              = 0x4D // M key
	key_N              = 0x4E // N key
	key_O              = 0x4F // O key
	key_P              = 0x50 // P key
	key_Q              = 0x51 // Q key
	key_R              = 0x52 // R key
	key_S              = 0x53 // S key
	key_T              = 0x54 // T key
	key_U              = 0x55 // U key
	key_V              = 0x56 // V key
	key_W              = 0x57 // W key
	key_X              = 0x58 // X key
	key_Y              = 0x59 // Y key
	key_Z              = 0x5A // Z key
	key_F1             = 0x70 // VK_F1            F1 key
	key_F2             = 0x71 // VK_F2            F2 key
	key_F3             = 0x72 // VK_F3            F3 key
	key_F4             = 0x73 // VK_F4            F4 key
	key_F5             = 0x74 // VK_F5            F5 key
	key_F6             = 0x75 // VK_F6            F6 key
	key_F7             = 0x76 // VK_F7            F7 key
	key_F8             = 0x77 // VK_F8            F8 key
	key_F9             = 0x78 // VK_F9            F9 key
	key_F10            = 0x79 // VK_F10           F10 key  ---- on osx-kb
	key_F11            = 0x7A // VK_F11           F11 key
	key_F12            = 0x7B // VK_F12           F12 key
	key_F13            = 0x7C // VK_F13           F13 key
	key_F14            = 0x2C // VK_F14 0x7D      F14 key  0x2C on osx-kb
	key_F15            = 0x91 // VK_F15 0x7E      F15 key  0x91
	key_F16            = 0x13 // VK_F16 0x7F      F16 key  0x13
	key_F17            = 0x80 // VK_F17           F17 key
	key_F18            = 0x81 // VK_F18           F18 key
	key_F19            = 0x82 // VK_F19           F19 key
	key_F20            = 0x83 // VK_F20           F20 key
	key_Keypad0        = 0x60 // VK_NUMPAD0  0x60 Numeric keypad 0 key :: VK_INSERT  0x20 on osx-kb
	key_Keypad1        = 0x61 // VK_NUMPAD1  0x61 Numeric keypad 1 key :: VK_END     0x23 on osx-kb
	key_Keypad2        = 0x62 // VK_NUMPAD2  0x62 Numeric keypad 2 key :: VK_DOWN    0x28 on osx-kb
	key_Keypad3        = 0x63 // VK_NUMPAD3  0x63 Numeric keypad 3 key :: VK_NEXT    0x22 on osx-kb
	key_Keypad4        = 0x64 // VK_NUMPAD4  0x64 Numeric keypad 4 key :: VK_LEFT    0x25 on osx-kb
	key_Keypad5        = 0x65 // VK_NUMPAD5  0x65 Numeric keypad 5 key :: VK_CLEAR   0x0C on osx-kb
	key_Keypad6        = 0x66 // VK_NUMPAD6  0x66 Numeric keypad 6 key :: VK_RIGHT   0x27 on osx-kb
	key_Keypad7        = 0x67 // VK_NUMPAD7  0x67 Numeric keypad 7 key :: VK_HOME    0x26 on osx-kb
	key_Keypad8        = 0x68 // VK_NUMPAD8  0x68 Numeric keypad 8 key :: VK_UP      0x21 on osx-kb
	key_Keypad9        = 0x69 // VK_NUMPAD9  0x69 Numeric keypad 9 key :: VK_PRIOR
	key_KeypadDecimal  = 0x6E // VK_DECIMAL       Decimal key :: VK_DELETE
	key_KeypadMultiply = 0x6A // VK_MULTIPLY      Multiply key
	key_KeypadPlus     = 0x6B // VK_ADD           Add key
	key_KeypadClear    = 0x90 // VK_CLEAR    0x0C CLEAR key :: VK_OEM_CLEAR 0xFE     0x90 on osx-kb
	key_KeypadDivide   = 0x6F // VK_DIVIDE        Divide key
	key_KeypadEnter    = 0x2B // VK_EXECUTE                            :: VK_ENTER on osx-kb
	key_KeypadMinus    = 0x6D // VK_SUBTRACT      Subtract key
	key_KeypadEquals   = 0x0C //
	key_Equal          = 0xBB //
	key_Minus          = 0xBD // VK_OEM_MINUS     For any country/region, the '-' key // VK_SEPARATOR 0x6C Separator key
	key_LeftBracket    = 0xDB // VK_OEM_4         misc characters; varys: US standard keyboard, the '[{' key
	key_RightBracket   = 0xDD // VK_OEM_6         misc characters; varys: US standard keyboard, the ']}' key
	key_Quote          = 0xDE // VK_OEM_7         misc characters; varys: US standard keyboard, the 'single/double-quote' key
	key_Semicolon      = 0xBA // VK_OEM_1         misc characters; varys: US standard keyboard, the ';:' key
	key_Backslash      = 0xDC // VK_OEM_5         misc characters; varys: US standard keyboard, the '/?' key
	key_Grave          = 0xC0 // VK_OEM_3         misc characters; varys: US standard keyboard, the '`~' key
	key_Slash          = 0xBF //
	key_Comma          = 0xBC // VK_OEM_COMMA     For any country/region, the ',' key
	key_Period         = 0xBE // VK_OEM_PERIOD    For any country/region, the '.' key
	key_Return         = 0x0D // VK_RETURN        ENTER key
	key_Tab            = 0x09 // VK_TAB           TAB key
	key_Space          = 0x20 // VK_SPACE         SPACEBAR
	key_Delete         = 0x08 // VK_BACK          BACKSPACE key
	key_ForwardDelete  = 0x2E // VK_DELETE        DEL key
	key_Escape         = 0x1B // VK_ESCAPE        ESC key
	key_Home           = 0x24 // VK_HOME          HOME key
	key_PageUp         = 0x21 // VK_PRIOR         PAGE UP key
	key_PageDown       = 0x22 // VK_NEXT          PAGE DOWN key
	key_LeftArrow      = 0x25 // VK_LEFT          LEFT ARROW key
	key_RightArrow     = 0x27 // VK_RIGHT         RIGHT ARROW key
	key_DownArrow      = 0x28 // VK_DOWN          DOWN ARROW key
	key_UpArrow        = 0x26 // VK_UP            UP ARROW key
	key_End            = 0x23 // VK_END           END key
	mouse_Left         = 0x01 // VK_LBUTTON Left mouse button (tack on unique values for mouse buttons)
	mouse_Middle       = 0x04 // VK_MBUTTON Middle mouse button (three-button mouse)
	mouse_Right        = 0x02 // VK_RBUTTON Right mouse button
)
