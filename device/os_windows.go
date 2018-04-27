// Copyright © 2013-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package device

// The microsoft (windows) native layer. This wraps the c functions that
// wrap the microsoft API's (where the real work is done).

// // C code and cgo directvies.
//
// #cgo windows CFLAGS: -m64
// #cgo windows,!dx LDFLAGS: -lopengl32 -lgdi32
// #cgo windows,dx LDFLAGS: -ld3d11
// #cgo windows,dx CXXFLAGS: -std=c++11
//
// #include "os_windows.h" // native function signatures and constants.
import "C" // must be located here.

import (
	"runtime"
	"time"
	"unsafe"
)

// OpenGL related, see: https://code.google.com/p/go-wiki/wiki/LockOSThread
func init() { runtime.LockOSThread() }

// Device instance needed to handle callbacks on exported methods.
var dev = &win{}

// runApp is the per-device entry method. Compiling will find the one
// that matches the requested or current platform.
func runApp(app App) {
	dev.app = app // The app receiving device callbacks.
	C.dev_run()   // does not return!
}

// prepRender is called from os_windows.c after the underlying window
// has been created and before the update render callbacks start.
//
//export prepRender
func prepRender() {
	dev.input = newInput()
	dev.app.Init(dev)
}

// renderFrame is the called from os_Windows.c as soon as the previous
// frame has been rendered. Actual frame rate is still limited by the
// monitor.
//
//export renderFrame
func renderFrame(mx, my int64) {
	if dev.input != nil {
		dev.app.Refresh(dev)
	}
}

// handleInput is called from os_windows.c to consolidate user input
// events into a consumable summary of which keys are current pressed.
// The summary is fowarded each update to the controlling application.
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
	}
}

// native layer callback functions. native->app calls.
// =============================================================================
// winOS Device implementation. app->native calls.

// win is the macOS implementation of the Device interface.
// See the Device interface for method descriptions.
type win struct {
	app      App       // update/render callback.
	input    *input    // tracks current keys pressed.
	lastSwap time.Time // helps throttle very fast apps.
}

// Implement the Device interface. See docs in device.go
// Mostly call the underlying native layer.
func (os *win) Down() *Pressed     { return os.input.getPressed(os.Cursor()) }
func (os *win) Dispose()           { C.dev_dispose() }
func (os *win) ToggleFullScreen()  { C.dev_toggle_fullscreen() }
func (os *win) IsFullScreen() bool { return uint(C.dev_fullscreen()) == 1 }
func (os *win) SwapBuffers() {
	elapsed := time.Since(os.lastSwap)
	os.lastSwap = time.Now()
	C.dev_swap()

	// Throttle unreasonable refresh rates since most monitors only
	// refresh at 60 or 120 times a second. Generally an update is
	// every 20ms, so start throttling if the app is twice that.
	// Use a smaller sleep time since sleep is not exact.
	if elapsed/time.Millisecond < 10 {
		time.Sleep(5 * time.Millisecond)
	}
}
func (os *win) SetCursorAt(x, y int) {
	C.dev_set_cursor_location(C.long(x), C.long(y))
}
func (os *win) Cursor() (x, y int) {
	var mx, my int32
	C.dev_cursor((*C.long)(&mx), (*C.long)(&my))
	return int(mx), int(my)
}
func (os *win) Size() (x, y, w, h int) {
	var winx, winy, width, height int32
	C.dev_size((*C.long)(&winx), (*C.long)(&winy), (*C.long)(&width), (*C.long)(&height))
	return int(winx), int(winy), int(width), int(height)
}
func (os *win) SetSize(x, y, w, h int) {
	C.dev_set_size(C.long(x), C.long(y), C.long(w), C.long(h))

	// resising doesn't trigger an OS resize event so inform
	// the application directly.
	os.input.curr.Resized = true
}
func (os *win) SetTitle(title string) {
	cstr := C.CString(title)
	defer C.free(unsafe.Pointer(cstr))
	C.dev_set_title(cstr)
}
func (os *win) ShowCursor(show bool) {
	trueFalse := 0 // trueFalse needs to be 0 or 1.
	if show {
		trueFalse = 1
	}
	C.dev_show_cursor(C.uchar(trueFalse))
}
func (os *win) Copy() string {
	if cstr := C.dev_clip_copy(); cstr != nil {
		str := C.GoString(cstr)      // make a Go copy.
		C.free(unsafe.Pointer(cstr)) // free the C copy.
		return str
	}
	return ""
}
func (os *win) Paste(s string) {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	C.dev_clip_paste(cstr)
}

// winOS Device implementation.
// =============================================================================

// Expose the underlying Win key codes to the Vu device key codes
// supported by each of the native layers.
//
// Windows virtual key codes.
// http://msdn.microsoft.com/en-ca/library/windows/desktop/dd375731(v=vs.85).aspx
const (
	// keyboard numbers.
	K0 = 0x30 // 0 key
	K1 = 0x31 // 1 key
	K2 = 0x32 // 2 key
	K3 = 0x33 // 3 key
	K4 = 0x34 // 4 key
	K5 = 0x35 // 5 key
	K6 = 0x36 // 6 key
	K7 = 0x37 // 7 key
	K8 = 0x38 // 8 key
	K9 = 0x39 // 9 key

	// keyboard letters.
	KA = 0x41 // A key
	KB = 0x42 // B key
	KC = 0x43 // C key
	KD = 0x44 // D key
	KE = 0x45 // E key
	KF = 0x46 // F key
	KG = 0x47 // G key
	KH = 0x48 // H key
	KI = 0x49 // I key
	KJ = 0x4A // J key
	KK = 0x4B // K key
	KL = 0x4C // L key
	KM = 0x4D // M key
	KN = 0x4E // N key
	KO = 0x4F // O key
	KP = 0x50 // P key
	KQ = 0x51 // Q key
	KR = 0x52 // R key
	KS = 0x53 // S key
	KT = 0x54 // T key
	KU = 0x55 // U key
	KV = 0x56 // V key
	KW = 0x57 // W key
	KX = 0x58 // X key
	KY = 0x59 // Y key
	KZ = 0x5A // Z key

	// Function Keys
	KF1  = 0x70 // VK_F1        F1 key
	KF2  = 0x71 // VK_F2        F2 key
	KF3  = 0x72 // VK_F3        F3 key
	KF4  = 0x73 // VK_F4        F4 key
	KF5  = 0x74 // VK_F5        F5 key
	KF6  = 0x75 // VK_F6        F6 key
	KF7  = 0x76 // VK_F7        F7 key
	KF8  = 0x77 // VK_F8        F8 key
	KF9  = 0x78 // VK_F9        F9 key
	KF10 = 0x79 // VK_F10       F10 key  ---- on osx-kb
	KF11 = 0x7A // VK_F11       F11 key
	KF12 = 0x7B // VK_F12       F12 key
	KF13 = 0x7C // VK_F13       F13 key
	KF14 = 0x2C // VK_F14 0x7D  F14 key  0x2C on osx-kb
	KF15 = 0x91 // VK_F15 0x7E  F15 key  0x91
	KF16 = 0x13 // VK_F16 0x7F  F16 key  0x13
	KF17 = 0x80 // VK_F17       F17 key
	KF18 = 0x81 // VK_F18       F18 key
	KF19 = 0x82 // VK_F19       F19 key
	KF20 = 0x83 // VK_F20       F20 key

	// Keypad keys
	KKpDot = 0x6E // VK_DECIMAL   Decimal key    :: VK_DELETE
	KKpMlt = 0x6A // VK_MULTIPLY  Multiply key
	KKpAdd = 0x6B // VK_ADD       Add key
	KKpClr = 0x90 // VK_CLEAR     0x0C CLEAR key :: VK_OEM_CLEAR 0xFE 0x90 on osx-kb
	KKpDiv = 0x6F // VK_DIVIDE    Divide key
	KKpEnt = 0x2B // VK_EXECUTE                  :: VK_ENTER on osx-kb
	KKpSub = 0x6D // VK_SUBTRACT      Subtract key
	KKpEql = 0xE2 //
	KKp0   = 0x60 // VK_NUMPAD0  0x60 keypad 0 key :: VK_INSERT  0x20 on osx-kb
	KKp1   = 0x61 // VK_NUMPAD1  0x61 keypad 1 key :: VK_END     0x23 on osx-kb
	KKp2   = 0x62 // VK_NUMPAD2  0x62 keypad 2 key :: VK_DOWN    0x28 on osx-kb
	KKp3   = 0x63 // VK_NUMPAD3  0x63 keypad 3 key :: VK_NEXT    0x22 on osx-kb
	KKp4   = 0x64 // VK_NUMPAD4  0x64 keypad 4 key :: VK_LEFT    0x25 on osx-kb
	KKp5   = 0x65 // VK_NUMPAD5  0x65 keypad 5 key :: VK_CLEAR   0x0C on osx-kb
	KKp6   = 0x66 // VK_NUMPAD6  0x66 keypad 6 key :: VK_RIGHT   0x27 on osx-kb
	KKp7   = 0x67 // VK_NUMPAD7  0x67 keypad 7 key :: VK_HOME    0x26 on osx-kb
	KKp8   = 0x68 // VK_NUMPAD8  0x68 keypad 8 key :: VK_UP      0x21 on osx-kb
	KKp9   = 0x69 // VK_NUMPAD9  0x69 keypad 9 key :: VK_PRIOR

	// Misc and Punctuation keys.
	KEqual = 0xBB //
	KMinus = 0xBD // VK_OEM_MINUS  For any country/region, the '-' key // VK_SEPARATOR 0x6C Separator key
	KLBkt  = 0xDB // VK_OEM_4      misc characters; varys: US keyboard, the '[{' key
	KRBkt  = 0xDD // VK_OEM_6      misc characters; varys: US keyboard, the ']}' key
	KQt    = 0xC0 // VK_OEM_7      misc characters; varys: US keyboard, the 'single/double-quote' key
	KSemi  = 0xBA // VK_OEM_1      misc characters; varys: US keyboard, the ';:' key
	KBSl   = 0xDE // VK_OEM_5      misc characters; varys: US keyboard, the '/?' key
	KComma = 0xBC // VK_OEM_COMMA  For any country/region, the ',' key
	KSlash = 0xBF //
	KDot   = 0xBE // VK_OEM_PERIOD For any country/region, the '.' key
	KGrave = 0xDF // VK_OEM_3      misc characters; varys: US keyboard, the '`~' key
	KRet   = 0x0D // VK_RETURN     ENTER key
	KTab   = 0x09 // VK_TAB        TAB key
	KSpace = 0x20 // VK_SPACE      SPACEBAR
	KDel   = 0x08 // VK_BACK       BACKSPACE key
	KEsc   = 0x1B // VK_ESCAPE     ESC key

	// Control keys.
	KHome  = 0x24         // VK_HOME    HOME key
	KPgUp  = 0x21         // VK_PRIOR   PAGE UP key
	KFDel  = 0x2E         // VK_DELETE  DEL key
	KEnd   = 0x23         // VK_END     END key
	KPgDn  = 0x22         // VK_NEXT    PAGE DOWN key
	KLa    = 0x25         // VK_LEFT    LEFT ARROW key
	KRa    = 0x27         // VK_RIGHT   RIGHT ARROW key
	KDa    = 0x28         // VK_DOWN    DOWN ARROW key
	KUa    = 0x26         // VK_UP      UP ARROW key
	KCtl   = C.VK_CONTROL // modifier masks and key codes.
	KFn    = 0            // Did not find on windows.
	KShift = C.VK_SHIFT
	KCmd   = C.VK_LWIN | C.VK_RWIN
	KAlt   = C.VK_MENU

	// Mouse buttons are treated like keys.
	// Values don't conflict with other key codes.
	KLm = C.devMouseL // Mouse buttons
	KMm = C.devMouseM //   "
	KRm = C.devMouseR //   "
)
