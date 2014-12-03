// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package device

// FUTURE
// Currently the most popular linux distro (ubuntu) is replacing the xserver
// display server with Mir. An example Mir shell can be found at:
//     http://unity.ubuntu.com/mir/demo_shell_8cpp-example.html
// Work on the linux portion will be delayed until there is an
// ubuntu version that fully supports Mir.

// The ubuntu (linux) native layer implementation.
// This wraps the linux windows API's (where the real work is done).

// // The following block is C code and cgo directvies.
//
// #cgo windows LDFLAGS: -lGL
//
// #include "os_linux.h"
import "C" // must be located here.

// OS specific structure to differentiate it from the other native layers.
type lin struct{}

// nativeLayer gets a reference to the native operating system.  Each native
// layer implements this factory method.  Compiling will leave only the one that
// matches the current platform.
func nativeLayer() native { return &lin{} }

// Implements native interface.
func (w lin) context(r *nrefs) int64 {
	return int64(C.gs_context((*C.longlong)(&(r.display)), (*C.longlong)(&(r.shell))))
}
func (w lin) display() int64                               { return 0 }
func (w lin) displayDispose(r *nrefs)                      {}
func (w lin) shell(r *nrefs) int64                         { return 0 }
func (w lin) shellOpen(r *nrefs)                           {}
func (w lin) shellAlive(r *nrefs) bool                     { return true }
func (w lin) swapBuffers(r *nrefs)                         {}
func (w lin) setAlphaBufferSize(size int)                  {}
func (w lin) setDepthBufferSize(size int)                  {}
func (w lin) setCursorAt(r *nrefs, x, y int)               {}
func (w lin) showCursor(r *nrefs, show bool)               {}
func (w lin) readDispatch(r *nrefs) *userInput             { return &userInput{} }
func (w lin) size(r *nrefs) (x int, y int, wx int, hy int) { return }
func (w lin) setSize(x, y, width, height int)              {}
func (w lin) setTitle(title string)                        {}

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
// Ubuntu key codes as shown by "xmodmap -pk".
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

/*
There are 7 KeySyms per KeyCode; KeyCodes range from 8 to 255.

    KeyCode	Keysym (Keysym)	...
    Value  	Value   (Name) 	...

      8
      9    	0xff1b (Escape)	0x0000 (NoSymbol)	0xff1b (Escape)
     10    	0x0031 (1)	0x0021 (exclam)	0x0031 (1)	0x0021 (exclam)
     11    	0x0032 (2)	0x0040 (at)	0x0032 (2)	0x0040 (at)
     12    	0x0033 (3)	0x0023 (numbersign)	0x0033 (3)	0x0023 (numbersign)
     13    	0x0034 (4)	0x0024 (dollar)	0x0034 (4)	0x0024 (dollar)
     14    	0x0035 (5)	0x0025 (percent)	0x0035 (5)	0x0025 (percent)
     15    	0x0036 (6)	0x005e (asciicircum)	0x0036 (6)	0x005e (asciicircum)
     16    	0x0037 (7)	0x0026 (ampersand)	0x0037 (7)	0x0026 (ampersand)
     17    	0x0038 (8)	0x002a (asterisk)	0x0038 (8)	0x002a (asterisk)
     18    	0x0039 (9)	0x0028 (parenleft)	0x0039 (9)	0x0028 (parenleft)
     19    	0x0030 (0)	0x0029 (parenright)	0x0030 (0)	0x0029 (parenright)
     20    	0x002d (minus)	0x005f (underscore)	0x002d (minus)	0x005f (underscore)
     21    	0x003d (equal)	0x002b (plus)	0x003d (equal)	0x002b (plus)
     22    	0xff08 (BackSpace)	0xff08 (BackSpace)	0xff08 (BackSpace)	0xff08 (BackSpace)
     23    	0xff09 (Tab)	0xfe20 (ISO_Left_Tab)	0xff09 (Tab)	0xfe20 (ISO_Left_Tab)
     24    	0x0071 (q)	0x0051 (Q)	0x0071 (q)	0x0051 (Q)
     25    	0x0077 (w)	0x0057 (W)	0x0077 (w)	0x0057 (W)
     26    	0x0065 (e)	0x0045 (E)	0x0065 (e)	0x0045 (E)
     27    	0x0072 (r)	0x0052 (R)	0x0072 (r)	0x0052 (R)
     28    	0x0074 (t)	0x0054 (T)	0x0074 (t)	0x0054 (T)
     29    	0x0079 (y)	0x0059 (Y)	0x0079 (y)	0x0059 (Y)
     30    	0x0075 (u)	0x0055 (U)	0x0075 (u)	0x0055 (U)
     31    	0x0069 (i)	0x0049 (I)	0x0069 (i)	0x0049 (I)
     32    	0x006f (o)	0x004f (O)	0x006f (o)	0x004f (O)
     33    	0x0070 (p)	0x0050 (P)	0x0070 (p)	0x0050 (P)
     34    	0x005b (bracketleft)	0x007b (braceleft)	0x005b (bracketleft)	0x007b (braceleft)
     35    	0x005d (bracketright)	0x007d (braceright)	0x005d (bracketright)	0x007d (braceright)
     36    	0xff0d (Return)	0x0000 (NoSymbol)	0xff0d (Return)
     37    	0xffe3 (Control_L)	0x0000 (NoSymbol)	0xffe3 (Control_L)
     38    	0x0061 (a)	0x0041 (A)	0x0061 (a)	0x0041 (A)
     39    	0x0073 (s)	0x0053 (S)	0x0073 (s)	0x0053 (S)
     40    	0x0064 (d)	0x0044 (D)	0x0064 (d)	0x0044 (D)
     41    	0x0066 (f)	0x0046 (F)	0x0066 (f)	0x0046 (F)
     42    	0x0067 (g)	0x0047 (G)	0x0067 (g)	0x0047 (G)
     43    	0x0068 (h)	0x0048 (H)	0x0068 (h)	0x0048 (H)
     44    	0x006a (j)	0x004a (J)	0x006a (j)	0x004a (J)
     45    	0x006b (k)	0x004b (K)	0x006b (k)	0x004b (K)
     46    	0x006c (l)	0x004c (L)	0x006c (l)	0x004c (L)
     47    	0x003b (semicolon)	0x003a (colon)	0x003b (semicolon)	0x003a (colon)
     48    	0x0027 (apostrophe)	0x0022 (quotedbl)	0x0027 (apostrophe)	0x0022 (quotedbl)
     49    	0x0060 (grave)	0x007e (asciitilde)	0x0060 (grave)	0x007e (asciitilde)
     50    	0xffe1 (Shift_L)	0x0000 (NoSymbol)	0xffe1 (Shift_L)
     51    	0x005c (backslash)	0x007c (bar)	0x005c (backslash)	0x007c (bar)
     52    	0x007a (z)	0x005a (Z)	0x007a (z)	0x005a (Z)
     53    	0x0078 (x)	0x0058 (X)	0x0078 (x)	0x0058 (X)
     54    	0x0063 (c)	0x0043 (C)	0x0063 (c)	0x0043 (C)
     55    	0x0076 (v)	0x0056 (V)	0x0076 (v)	0x0056 (V)
     56    	0x0062 (b)	0x0042 (B)	0x0062 (b)	0x0042 (B)
     57    	0x006e (n)	0x004e (N)	0x006e (n)	0x004e (N)
     58    	0x006d (m)	0x004d (M)	0x006d (m)	0x004d (M)
     59    	0x002c (comma)	0x003c (less)	0x002c (comma)	0x003c (less)
     60    	0x002e (period)	0x003e (greater)	0x002e (period)	0x003e (greater)
     61    	0x002f (slash)	0x003f (question)	0x002f (slash)	0x003f (question)
     62    	0xffe2 (Shift_R)	0x0000 (NoSymbol)	0xffe2 (Shift_R)
     63    	0xffaa (KP_Multiply)	0xffaa (KP_Multiply)	0xffaa (KP_Multiply)	0xffaa (KP_Multiply)	0xffaa (KP_Multiply)	0xffaa (KP_Multiply)	0x1008fe21 (XF86ClearGrab)
     64    	0xffe9 (Alt_L)	0xffe7 (Meta_L)	0xffe9 (Alt_L)	0xffe7 (Meta_L)
     65    	0x0020 (space)	0x0000 (NoSymbol)	0x0020 (space)
     66    	0xffe5 (Caps_Lock)	0x0000 (NoSymbol)	0xffe5 (Caps_Lock)
     67    	0xffbe (F1)	0xffbe (F1)	0xffbe (F1)	0xffbe (F1)	0xffbe (F1)	0xffbe (F1)	0x1008fe01 (XF86Switch_VT_1)
     68    	0xffbf (F2)	0xffbf (F2)	0xffbf (F2)	0xffbf (F2)	0xffbf (F2)	0xffbf (F2)	0x1008fe02 (XF86Switch_VT_2)
     69    	0xffc0 (F3)	0xffc0 (F3)	0xffc0 (F3)	0xffc0 (F3)	0xffc0 (F3)	0xffc0 (F3)	0x1008fe03 (XF86Switch_VT_3)
     70    	0xffc1 (F4)	0xffc1 (F4)	0xffc1 (F4)	0xffc1 (F4)	0xffc1 (F4)	0xffc1 (F4)	0x1008fe04 (XF86Switch_VT_4)
     71    	0xffc2 (F5)	0xffc2 (F5)	0xffc2 (F5)	0xffc2 (F5)	0xffc2 (F5)	0xffc2 (F5)	0x1008fe05 (XF86Switch_VT_5)
     72    	0xffc3 (F6)	0xffc3 (F6)	0xffc3 (F6)	0xffc3 (F6)	0xffc3 (F6)	0xffc3 (F6)	0x1008fe06 (XF86Switch_VT_6)
     73    	0xffc4 (F7)	0xffc4 (F7)	0xffc4 (F7)	0xffc4 (F7)	0xffc4 (F7)	0xffc4 (F7)	0x1008fe07 (XF86Switch_VT_7)
     74    	0xffc5 (F8)	0xffc5 (F8)	0xffc5 (F8)	0xffc5 (F8)	0xffc5 (F8)	0xffc5 (F8)	0x1008fe08 (XF86Switch_VT_8)
     75    	0xffc6 (F9)	0xffc6 (F9)	0xffc6 (F9)	0xffc6 (F9)	0xffc6 (F9)	0xffc6 (F9)	0x1008fe09 (XF86Switch_VT_9)
     76    	0xffc7 (F10)	0xffc7 (F10)	0xffc7 (F10)	0xffc7 (F10)	0xffc7 (F10)	0xffc7 (F10)	0x1008fe0a (XF86Switch_VT_10)
     77    	0xff7f (Num_Lock)	0x0000 (NoSymbol)	0xff7f (Num_Lock)
     78    	0xff14 (Scroll_Lock)	0x0000 (NoSymbol)	0xff14 (Scroll_Lock)
     79    	0xff95 (KP_Home)	0xffb7 (KP_7)	0xff95 (KP_Home)	0xffb7 (KP_7)
     80    	0xff97 (KP_Up)	0xffb8 (KP_8)	0xff97 (KP_Up)	0xffb8 (KP_8)
     81    	0xff9a (KP_Prior)	0xffb9 (KP_9)	0xff9a (KP_Prior)	0xffb9 (KP_9)
     82    	0xffad (KP_Subtract)	0xffad (KP_Subtract)	0xffad (KP_Subtract)	0xffad (KP_Subtract)	0xffad (KP_Subtract)	0xffad (KP_Subtract)	0x1008fe23 (XF86Prev_VMode)
     83    	0xff96 (KP_Left)	0xffb4 (KP_4)	0xff96 (KP_Left)	0xffb4 (KP_4)
     84    	0xff9d (KP_Begin)	0xffb5 (KP_5)	0xff9d (KP_Begin)	0xffb5 (KP_5)
     85    	0xff98 (KP_Right)	0xffb6 (KP_6)	0xff98 (KP_Right)	0xffb6 (KP_6)
     86    	0xffab (KP_Add)	0xffab (KP_Add)	0xffab (KP_Add)	0xffab (KP_Add)	0xffab (KP_Add)	0xffab (KP_Add)	0x1008fe22 (XF86Next_VMode)
     87    	0xff9c (KP_End)	0xffb1 (KP_1)	0xff9c (KP_End)	0xffb1 (KP_1)
     88    	0xff99 (KP_Down)	0xffb2 (KP_2)	0xff99 (KP_Down)	0xffb2 (KP_2)
     89    	0xff9b (KP_Next)	0xffb3 (KP_3)	0xff9b (KP_Next)	0xffb3 (KP_3)
     90    	0xff9e (KP_Insert)	0xffb0 (KP_0)	0xff9e (KP_Insert)	0xffb0 (KP_0)
     91    	0xff9f (KP_Delete)	0xffae (KP_Decimal)	0xff9f (KP_Delete)	0xffae (KP_Decimal)
     92    	0xfe03 (ISO_Level3_Shift)	0x0000 (NoSymbol)	0xfe03 (ISO_Level3_Shift)
     93
     94    	0x003c (less)	0x003e (greater)	0x003c (less)	0x003e (greater)	0x007c (bar)	0x00a6 (brokenbar)	0x007c (bar)
     95    	0xffc8 (F11)	0xffc8 (F11)	0xffc8 (F11)	0xffc8 (F11)	0xffc8 (F11)	0xffc8 (F11)	0x1008fe0b (XF86Switch_VT_11)
     96    	0xffc9 (F12)	0xffc9 (F12)	0xffc9 (F12)	0xffc9 (F12)	0xffc9 (F12)	0xffc9 (F12)	0x1008fe0c (XF86Switch_VT_12)
     97
     98    	0xff26 (Katakana)	0x0000 (NoSymbol)	0xff26 (Katakana)
     99    	0xff25 (Hiragana)	0x0000 (NoSymbol)	0xff25 (Hiragana)
    100    	0xff23 (Henkan_Mode)	0x0000 (NoSymbol)	0xff23 (Henkan_Mode)
    101    	0xff27 (Hiragana_Katakana)	0x0000 (NoSymbol)	0xff27 (Hiragana_Katakana)
    102    	0xff22 (Muhenkan)	0x0000 (NoSymbol)	0xff22 (Muhenkan)
    103
    104    	0xff8d (KP_Enter)	0x0000 (NoSymbol)	0xff8d (KP_Enter)
    105    	0xffe4 (Control_R)	0x0000 (NoSymbol)	0xffe4 (Control_R)
    106    	0xffaf (KP_Divide)	0xffaf (KP_Divide)	0xffaf (KP_Divide)	0xffaf (KP_Divide)	0xffaf (KP_Divide)	0xffaf (KP_Divide)	0x1008fe20 (XF86Ungrab)
    107    	0xff61 (Print)	0xff15 (Sys_Req)	0xff61 (Print)	0xff15 (Sys_Req)
    108    	0xffea (Alt_R)	0xffe8 (Meta_R)	0xffea (Alt_R)	0xffe8 (Meta_R)
    109    	0xff0a (Linefeed)	0x0000 (NoSymbol)	0xff0a (Linefeed)
    110    	0xff50 (Home)	0x0000 (NoSymbol)	0xff50 (Home)
    111    	0xff52 (Up)	0x0000 (NoSymbol)	0xff52 (Up)
    112    	0xff55 (Prior)	0x0000 (NoSymbol)	0xff55 (Prior)
    113    	0xff51 (Left)	0x0000 (NoSymbol)	0xff51 (Left)
    114    	0xff53 (Right)	0x0000 (NoSymbol)	0xff53 (Right)
    115    	0xff57 (End)	0x0000 (NoSymbol)	0xff57 (End)
    116    	0xff54 (Down)	0x0000 (NoSymbol)	0xff54 (Down)
    117    	0xff56 (Next)	0x0000 (NoSymbol)	0xff56 (Next)
    118    	0xff63 (Insert)	0x0000 (NoSymbol)	0xff63 (Insert)
    119    	0xffff (Delete)	0x0000 (NoSymbol)	0xffff (Delete)
    120
    121    	0x1008ff12 (XF86AudioMute)	0x0000 (NoSymbol)	0x1008ff12 (XF86AudioMute)
    122    	0x1008ff11 (XF86AudioLowerVolume)	0x0000 (NoSymbol)	0x1008ff11 (XF86AudioLowerVolume)
    123    	0x1008ff13 (XF86AudioRaiseVolume)	0x0000 (NoSymbol)	0x1008ff13 (XF86AudioRaiseVolume)
    124    	0x1008ff2a (XF86PowerOff)	0x0000 (NoSymbol)	0x1008ff2a (XF86PowerOff)
    125    	0xffbd (KP_Equal)	0x0000 (NoSymbol)	0xffbd (KP_Equal)
    126    	0x00b1 (plusminus)	0x0000 (NoSymbol)	0x00b1 (plusminus)
    127    	0xff13 (Pause)	0xff6b (Break)	0xff13 (Pause)	0xff6b (Break)
    128    	0x1008ff4a (XF86LaunchA)	0x0000 (NoSymbol)	0x1008ff4a (XF86LaunchA)
    129    	0xffae (KP_Decimal)	0xffae (KP_Decimal)	0xffae (KP_Decimal)	0xffae (KP_Decimal)
    130    	0xff31 (Hangul)	0x0000 (NoSymbol)	0xff31 (Hangul)
    131    	0xff34 (Hangul_Hanja)	0x0000 (NoSymbol)	0xff34 (Hangul_Hanja)
    132
    133    	0xffeb (Super_L)	0x0000 (NoSymbol)	0xffeb (Super_L)
    134    	0xffec (Super_R)	0x0000 (NoSymbol)	0xffec (Super_R)
    135    	0xff67 (Menu)	0x0000 (NoSymbol)	0xff67 (Menu)
    136    	0xff69 (Cancel)	0x0000 (NoSymbol)	0xff69 (Cancel)
    137    	0xff66 (Redo)	0x0000 (NoSymbol)	0xff66 (Redo)
    138    	0x1005ff70 (SunProps)	0x0000 (NoSymbol)	0x1005ff70 (SunProps)
    139    	0xff65 (Undo)	0x0000 (NoSymbol)	0xff65 (Undo)
    140    	0x1005ff71 (SunFront)	0x0000 (NoSymbol)	0x1005ff71 (SunFront)
    141    	0x1008ff57 (XF86Copy)	0x0000 (NoSymbol)	0x1008ff57 (XF86Copy)
    142    	0x1005ff73 (SunOpen)	0x0000 (NoSymbol)	0x1005ff73 (SunOpen)
    143    	0x1008ff6d (XF86Paste)	0x0000 (NoSymbol)	0x1008ff6d (XF86Paste)
    144    	0xff68 (Find)	0x0000 (NoSymbol)	0xff68 (Find)
    145    	0x1008ff58 (XF86Cut)	0x0000 (NoSymbol)	0x1008ff58 (XF86Cut)
    146    	0xff6a (Help)	0x0000 (NoSymbol)	0xff6a (Help)
    147    	0x1008ff65 (XF86MenuKB)	0x0000 (NoSymbol)	0x1008ff65 (XF86MenuKB)
    148    	0x1008ff1d (XF86Calculator)	0x0000 (NoSymbol)	0x1008ff1d (XF86Calculator)
    149
    150    	0x1008ff2f (XF86Sleep)	0x0000 (NoSymbol)	0x1008ff2f (XF86Sleep)
    151    	0x1008ff2b (XF86WakeUp)	0x0000 (NoSymbol)	0x1008ff2b (XF86WakeUp)
    152    	0x1008ff5d (XF86Explorer)	0x0000 (NoSymbol)	0x1008ff5d (XF86Explorer)
    153    	0x1008ff7b (XF86Send)	0x0000 (NoSymbol)	0x1008ff7b (XF86Send)
    154
    155    	0x1008ff8a (XF86Xfer)	0x0000 (NoSymbol)	0x1008ff8a (XF86Xfer)
    156    	0x1008ff41 (XF86Launch1)	0x0000 (NoSymbol)	0x1008ff41 (XF86Launch1)
    157    	0x1008ff42 (XF86Launch2)	0x0000 (NoSymbol)	0x1008ff42 (XF86Launch2)
    158    	0x1008ff2e (XF86WWW)	0x0000 (NoSymbol)	0x1008ff2e (XF86WWW)
    159    	0x1008ff5a (XF86DOS)	0x0000 (NoSymbol)	0x1008ff5a (XF86DOS)
    160    	0x1008ff2d (XF86ScreenSaver)	0x0000 (NoSymbol)	0x1008ff2d (XF86ScreenSaver)
    161
    162    	0x1008ff74 (XF86RotateWindows)	0x0000 (NoSymbol)	0x1008ff74 (XF86RotateWindows)
    163    	0x1008ff19 (XF86Mail)	0x0000 (NoSymbol)	0x1008ff19 (XF86Mail)
    164    	0x1008ff30 (XF86Favorites)	0x0000 (NoSymbol)	0x1008ff30 (XF86Favorites)
    165    	0x1008ff33 (XF86MyComputer)	0x0000 (NoSymbol)	0x1008ff33 (XF86MyComputer)
    166    	0x1008ff26 (XF86Back)	0x0000 (NoSymbol)	0x1008ff26 (XF86Back)
    167    	0x1008ff27 (XF86Forward)	0x0000 (NoSymbol)	0x1008ff27 (XF86Forward)
    168
    169    	0x1008ff2c (XF86Eject)	0x0000 (NoSymbol)	0x1008ff2c (XF86Eject)
    170    	0x1008ff2c (XF86Eject)	0x1008ff2c (XF86Eject)	0x1008ff2c (XF86Eject)	0x1008ff2c (XF86Eject)
    171    	0x1008ff17 (XF86AudioNext)	0x0000 (NoSymbol)	0x1008ff17 (XF86AudioNext)
    172    	0x1008ff14 (XF86AudioPlay)	0x1008ff31 (XF86AudioPause)	0x1008ff14 (XF86AudioPlay)	0x1008ff31 (XF86AudioPause)
    173    	0x1008ff16 (XF86AudioPrev)	0x0000 (NoSymbol)	0x1008ff16 (XF86AudioPrev)
    174    	0x1008ff15 (XF86AudioStop)	0x1008ff2c (XF86Eject)	0x1008ff15 (XF86AudioStop)	0x1008ff2c (XF86Eject)
    175    	0x1008ff1c (XF86AudioRecord)	0x0000 (NoSymbol)	0x1008ff1c (XF86AudioRecord)
    176    	0x1008ff3e (XF86AudioRewind)	0x0000 (NoSymbol)	0x1008ff3e (XF86AudioRewind)
    177    	0x1008ff6e (XF86Phone)	0x0000 (NoSymbol)	0x1008ff6e (XF86Phone)
    178
    179    	0x1008ff81 (XF86Tools)	0x0000 (NoSymbol)	0x1008ff81 (XF86Tools)
    180    	0x1008ff18 (XF86HomePage)	0x0000 (NoSymbol)	0x1008ff18 (XF86HomePage)
    181    	0x1008ff73 (XF86Reload)	0x0000 (NoSymbol)	0x1008ff73 (XF86Reload)
    182    	0x1008ff56 (XF86Close)	0x0000 (NoSymbol)	0x1008ff56 (XF86Close)
    183
    184
    185    	0x1008ff78 (XF86ScrollUp)	0x0000 (NoSymbol)	0x1008ff78 (XF86ScrollUp)
    186    	0x1008ff79 (XF86ScrollDown)	0x0000 (NoSymbol)	0x1008ff79 (XF86ScrollDown)
    187    	0x0028 (parenleft)	0x0000 (NoSymbol)	0x0028 (parenleft)
    188    	0x0029 (parenright)	0x0000 (NoSymbol)	0x0029 (parenright)
    189    	0x1008ff68 (XF86New)	0x0000 (NoSymbol)	0x1008ff68 (XF86New)
    190    	0xff66 (Redo)	0x0000 (NoSymbol)	0xff66 (Redo)
    191    	0x1008ff81 (XF86Tools)	0x0000 (NoSymbol)	0x1008ff81 (XF86Tools)
    192    	0x1008ff45 (XF86Launch5)	0x0000 (NoSymbol)	0x1008ff45 (XF86Launch5)
    193    	0x1008ff46 (XF86Launch6)	0x0000 (NoSymbol)	0x1008ff46 (XF86Launch6)
    194    	0x1008ff47 (XF86Launch7)	0x0000 (NoSymbol)	0x1008ff47 (XF86Launch7)
    195    	0x1008ff48 (XF86Launch8)	0x0000 (NoSymbol)	0x1008ff48 (XF86Launch8)
    196    	0x1008ff49 (XF86Launch9)	0x0000 (NoSymbol)	0x1008ff49 (XF86Launch9)
    197
    198
    199    	0x1008ffa9 (XF86TouchpadToggle)	0x0000 (NoSymbol)	0x1008ffa9 (XF86TouchpadToggle)
    200    	0x1008ffb0 (XF86TouchpadOn)	0x0000 (NoSymbol)	0x1008ffb0 (XF86TouchpadOn)
    201    	0x1008ffb1 (XF86TouchpadOff)	0x0000 (NoSymbol)	0x1008ffb1 (XF86TouchpadOff)
    202
    203    	0xff7e (Mode_switch)	0x0000 (NoSymbol)	0xff7e (Mode_switch)
    204    	0x0000 (NoSymbol)	0xffe9 (Alt_L)	0x0000 (NoSymbol)	0xffe9 (Alt_L)
    205    	0x0000 (NoSymbol)	0xffe7 (Meta_L)	0x0000 (NoSymbol)	0xffe7 (Meta_L)
    206    	0x0000 (NoSymbol)	0xffeb (Super_L)	0x0000 (NoSymbol)	0xffeb (Super_L)
    207    	0x0000 (NoSymbol)	0xffed (Hyper_L)	0x0000 (NoSymbol)	0xffed (Hyper_L)
    208    	0x1008ff14 (XF86AudioPlay)	0x0000 (NoSymbol)	0x1008ff14 (XF86AudioPlay)
    209    	0x1008ff31 (XF86AudioPause)	0x0000 (NoSymbol)	0x1008ff31 (XF86AudioPause)
    210    	0x1008ff43 (XF86Launch3)	0x0000 (NoSymbol)	0x1008ff43 (XF86Launch3)
    211    	0x1008ff44 (XF86Launch4)	0x0000 (NoSymbol)	0x1008ff44 (XF86Launch4)
    212    	0x1008ff4b (XF86LaunchB)	0x0000 (NoSymbol)	0x1008ff4b (XF86LaunchB)
    213    	0x1008ffa7 (XF86Suspend)	0x0000 (NoSymbol)	0x1008ffa7 (XF86Suspend)
    214    	0x1008ff56 (XF86Close)	0x0000 (NoSymbol)	0x1008ff56 (XF86Close)
    215    	0x1008ff14 (XF86AudioPlay)	0x0000 (NoSymbol)	0x1008ff14 (XF86AudioPlay)
    216    	0x1008ff97 (XF86AudioForward)	0x0000 (NoSymbol)	0x1008ff97 (XF86AudioForward)
    217
    218    	0xff61 (Print)	0x0000 (NoSymbol)	0xff61 (Print)
    219
    220    	0x1008ff8f (XF86WebCam)	0x0000 (NoSymbol)	0x1008ff8f (XF86WebCam)
    221
    222
    223    	0x1008ff19 (XF86Mail)	0x0000 (NoSymbol)	0x1008ff19 (XF86Mail)
    224    	0x1008ff8e (XF86Messenger)	0x0000 (NoSymbol)	0x1008ff8e (XF86Messenger)
    225    	0x1008ff1b (XF86Search)	0x0000 (NoSymbol)	0x1008ff1b (XF86Search)
    226    	0x1008ff5f (XF86Go)	0x0000 (NoSymbol)	0x1008ff5f (XF86Go)
    227    	0x1008ff3c (XF86Finance)	0x0000 (NoSymbol)	0x1008ff3c (XF86Finance)
    228    	0x1008ff5e (XF86Game)	0x0000 (NoSymbol)	0x1008ff5e (XF86Game)
    229    	0x1008ff36 (XF86Shop)	0x0000 (NoSymbol)	0x1008ff36 (XF86Shop)
    230
    231    	0xff69 (Cancel)	0x0000 (NoSymbol)	0xff69 (Cancel)
    232    	0x1008ff03 (XF86MonBrightnessDown)	0x0000 (NoSymbol)	0x1008ff03 (XF86MonBrightnessDown)
    233    	0x1008ff02 (XF86MonBrightnessUp)	0x0000 (NoSymbol)	0x1008ff02 (XF86MonBrightnessUp)
    234    	0x1008ff32 (XF86AudioMedia)	0x0000 (NoSymbol)	0x1008ff32 (XF86AudioMedia)
    235    	0x1008ff59 (XF86Display)	0x0000 (NoSymbol)	0x1008ff59 (XF86Display)
    236    	0x1008ff04 (XF86KbdLightOnOff)	0x0000 (NoSymbol)	0x1008ff04 (XF86KbdLightOnOff)
    237    	0x1008ff06 (XF86KbdBrightnessDown)	0x0000 (NoSymbol)	0x1008ff06 (XF86KbdBrightnessDown)
    238    	0x1008ff05 (XF86KbdBrightnessUp)	0x0000 (NoSymbol)	0x1008ff05 (XF86KbdBrightnessUp)
    239    	0x1008ff7b (XF86Send)	0x0000 (NoSymbol)	0x1008ff7b (XF86Send)
    240    	0x1008ff72 (XF86Reply)	0x0000 (NoSymbol)	0x1008ff72 (XF86Reply)
    241    	0x1008ff90 (XF86MailForward)	0x0000 (NoSymbol)	0x1008ff90 (XF86MailForward)
    242    	0x1008ff77 (XF86Save)	0x0000 (NoSymbol)	0x1008ff77 (XF86Save)
    243    	0x1008ff5b (XF86Documents)	0x0000 (NoSymbol)	0x1008ff5b (XF86Documents)
    244    	0x1008ff93 (XF86Battery)	0x0000 (NoSymbol)	0x1008ff93 (XF86Battery)
    245    	0x1008ff94 (XF86Bluetooth)	0x0000 (NoSymbol)	0x1008ff94 (XF86Bluetooth)
    246    	0x1008ff95 (XF86WLAN)	0x0000 (NoSymbol)	0x1008ff95 (XF86WLAN)
    247
    248
*/
