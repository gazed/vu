package device

// os_windows.go wraps the the microsoft windows native layer.

import (
	"fmt"
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"github.com/gazed/vu/internal/device/win"
	"golang.org/x/sys/windows"
)

// newPlatform returns the platform when running on a windows system.
func newPlatform() platformAPI { return &windowsDevice{} }

// windowsDevice holds display state and implements the Device interface
// for the windows platform.
type windowsDevice struct {
	hinstance win.HINSTANCE // module instance
	hwnd      win.HWND      // window handle
	windowed  bool          // window bordered vs full screen
	title     string        // window with border title
	windx     int32         // window with border bottom left corner
	windy     int32         // window with border bottom left corner
	windw     int32         // window with border width
	windh     int32         // window with border height
	fullw     int32         // fullscreen width
	fullh     int32         // fullscreen height
}

// newPlatform gets the platform specific window and input handler.
func (wd *windowsDevice) init(windowed bool, title string, x int32, y int32, w int32, h int32) {
	// https://stackoverflow.com/questions/25361831/benefits-of-runtime-lockosthread-in-golang
	// "With the Go threading model, calls to C code, assembler code, or blocking
	//  system calls occur in the same thread as the calling Go code, which is
	//  managed by the Go runtime scheduler. The os.LockOSThread() mechanism is
	//  mostly useful when Go has to interface with some foreign library
	//  (a C library for instance). It guarantees that several successive calls
	//  to this library will be done in the same thread."
	runtime.LockOSThread()
	wd.windowed = windowed
	wd.title = title
	wd.windx = x
	wd.windy = y
	wd.windw = w
	wd.windh = h
	wd.fullw = win.GetSystemMetrics(win.SM_CXSCREEN)
	wd.fullh = win.GetSystemMetrics(win.SM_CYSCREEN)
}

// GetRenderSurfaceInfo exposes the windows API specific information
// needed by the render package to create a rendering surface.
// Called by the platform specific code in the Render package.
func GetRenderSurfaceInfo(d *Device) (hinst windows.Handle, hwnd windows.HWND, err error) {
	if wd, ok := d.platform.(*windowsDevice); ok {
		return windows.Handle(wd.hinstance), windows.HWND(wd.hwnd), nil
	}
	return 0, 0, fmt.Errorf("GetRenderSurfaceInfo: invalid device")
}

// User input data is refreshed each call to PollInput
// This is shared with the App.
var input = &Input{
	Pressed:  map[int32]bool{},
	Down:     map[int32]time.Time{},
	Released: map[int32]time.Duration{},
}
var userShutdown bool = false // set to true if the user closes the window.

// resizeHandler processes resize events immediately since the windows loop
// shuts down on MINIMIZED events
var resizeHandler func() = nil

func (wd *windowsDevice) setResizeHandler(callback func()) { resizeHandler = callback }

// Device interface: createDisplay
func (wd *windowsDevice) createDisplay() error {

	// get the application instance.
	wd.hinstance = win.GetModuleHandle(nil)
	if wd.hinstance == 0 {
		return fmt.Errorf("GetModuleHandle failed %d", win.GetLastError())
	}
	classname := syscall.StringToUTF16Ptr("vuwin")

	// treat some constants as pointers because thats what windows wants.
	arrow := (unsafe.Pointer)((uintptr)(win.IDC_ARROW))
	appIcon := (unsafe.Pointer)((uintptr)(win.IDI_APPLICATION))

	// create the window.
	wc := win.WNDCLASSEX{}
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.Style = win.CS_HREDRAW | win.CS_VREDRAW | win.CS_OWNDC
	wc.LpfnWndProc = syscall.NewCallback(winProcessMsg)
	wc.CbClsExtra = 0
	wc.CbWndExtra = 0
	wc.HInstance = wd.hinstance
	wc.HIcon = win.LoadIcon(0, (*uint16)(appIcon))
	wc.HCursor = win.LoadCursor(0, (*uint16)(arrow))
	wc.HbrBackground = (win.HBRUSH)(win.COLOR_WINDOW + 1)
	wc.LpszMenuName = nil
	wc.LpszClassName = classname
	wc.HIconSm = 0
	if win.RegisterClassEx(&wc) == 0 {
		return fmt.Errorf("RegisterClassEx failed %d", win.GetLastError())
	}

	var style uint32
	var wx, wy, ww, wh int32
	styleEx := uint32(win.WS_EX_APPWINDOW)
	if wd.windowed {
		style = uint32(win.WS_OVERLAPPED | win.WS_SYSMENU | win.WS_CAPTION)

		// adjust the window dimensions to accommodate the window frame.
		style = style | win.WS_THICKFRAME | win.WS_MINIMIZEBOX | win.WS_MAXIMIZEBOX
		border := win.RECT{0, 0, 0, 0}
		win.AdjustWindowRectEx(&border, style, false, styleEx)
		wd.windw += border.Right - border.Left
		wd.windh += border.Bottom - border.Top
		wd.windx += border.Left
		wd.windy += border.Top
		wx, wy, ww, wh = wd.windx, wd.windy, wd.windw, wd.windh
	} else {
		style = uint32(win.WS_POPUP | win.WS_VISIBLE)
		wx, wy, ww, wh = 0, 0, wd.fullw, wd.fullh
	}

	// create the application window.
	wintitle := syscall.StringToUTF16Ptr(wd.title)
	var lpParam unsafe.Pointer
	wd.hwnd = win.CreateWindowEx(
		styleEx,
		classname, // must match the classname used in RegisterClassEx
		wintitle,  // visible window title
		style,
		wx,
		wy,
		ww,
		wh,
		win.HWND(0),
		win.HMENU(0),
		wd.hinstance,
		lpParam,
	)
	if wd.hwnd == 0 {
		return fmt.Errorf("CreateWindowEx failed %d", win.GetLastError())
	}

	// show the window.
	show := win.SW_SHOWMAXIMIZED
	if wd.windowed {
		show = win.SW_SHOW
	}
	win.ShowWindow(wd.hwnd, int32(show))
	win.SetForegroundWindow(wd.hwnd)
	return nil
}

// Windows callback procedure. This method is mostly microsoft magic
// as each event has its own behaviour and different return codes.
func winProcessMsg(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) uintptr {
	switch msg {
	case win.WM_ERASEBKGND:
		// repaint window, generally on window resize.
		return 1 // app will redraw the window.
	case win.WM_CLOSE:
		// user has requested shutdown. Closing the window immediately here
		// means the app should be ready for shutdown without notification.
		input.shutdown = true
		return 0
	case win.WM_DESTROY:
		win.PostQuitMessage(0) // generates WM_QUIT message
		return 0
	case win.WM_SIZE:
		if resizeHandler != nil {
			resizeHandler()
		}
	case win.WM_EXITSIZEMOVE:
		if resizeHandler != nil {
			resizeHandler()
		}
		return 0
	case win.WM_ACTIVATE:
		// window is gaining or losing focus.
		if win.LOWORD(uint32(wParam)) == win.WA_INACTIVE {
			input.loseFocus()
			return 0
		}
		return 0 // nothing to do if focus is gained.

	// keys presses
	case win.WM_KEYDOWN:
		input.keyPressed(int32(wParam))
		return 0
	case win.WM_KEYUP:
		input.keyReleased(int32(wParam))
		return 0
	case win.WM_SYSKEYDOWN:
		// allow some syskeys for use in games
		if wParam == KAlt || wParam == KF10 {
			input.keyPressed(int32(wParam))
			return 0
		}
		return win.DefWindowProc(hwnd, msg, wParam, lParam)
	case win.WM_SYSKEYUP:
		// allow some syskeys for use in games
		if wParam == KAlt || wParam == KF10 {
			input.keyReleased(int32(wParam))
			return 0
		}
		return win.DefWindowProc(hwnd, msg, wParam, lParam)
	case win.WM_SYSCOMMAND:
		if (wParam & 0xfff0) == win.SC_KEYMENU {
			// ignore windows system commands like F10 - menu
			return 0
		}

	// mouse buttons
	case win.WM_LBUTTONDOWN:
		win.SetCapture(hwnd)
		input.keyPressed(KML)
		return 0
	case win.WM_LBUTTONUP:
		input.keyReleased(KML)
		if input.allMouseButtonsReleased() {
			win.ReleaseCapture()
		}
		return 0
	case win.WM_MBUTTONDOWN:
		win.SetCapture(hwnd)
		input.keyPressed(KMM)
		return 0
	case win.WM_MBUTTONUP:
		input.keyReleased(KMM)
		if input.allMouseButtonsReleased() {
			win.ReleaseCapture()
		}
		return 0
	case win.WM_RBUTTONDOWN:
		win.SetCapture(hwnd)
		input.keyPressed(KMR)
		return 0
	case win.WM_RBUTTONUP:
		input.keyReleased(KMR)
		if input.allMouseButtonsReleased() {
			win.ReleaseCapture()
		}
		return 0
	case win.WM_MOUSEWHEEL:
		// normalize the mouse delta from the high word
		if delta := int16(wParam >> 16); delta != 0 {
			if delta < 0 {
				input.Scroll = -1 // scroll backward
			} else {
				input.Scroll = 1 // scroll forward
			}
		}
		return 0
	}

	// Pass all unhandled messages to DefWindowProc
	return win.DefWindowProc(hwnd, msg, wParam, lParam)
}

// dispose implement Device.
// Destroy the application window. Attempt to remove the rendering context
// and the device context as well.
func (wd *windowsDevice) dispose() {
	if wd.hwnd != 0 {
		win.DestroyWindow(wd.hwnd)
		wd.hwnd = 0
	}
}

// surfaceSize implements Device.
func (wd *windowsDevice) surfaceSize() (w, h uint32) {
	var rect win.RECT
	win.GetClientRect(wd.hwnd, &rect)
	w = uint32(rect.Right - rect.Left)
	h = uint32(rect.Bottom - rect.Top)
	return w, h
}

// isRunning implement Device.
func (wd *windowsDevice) isRunning() bool { return wd.hwnd != 0 }

// getInput implements Device.
// Handles all outstanding messages and returns.
func (wd *windowsDevice) getInput() *Input {
	input.reset()        // clear the shared input.
	if !wd.isRunning() { // only continue if the window is up.
		return input
	}
	var msg win.MSG
	for win.PeekMessage(&msg, win.WM_NULL, 0, 0, win.PM_REMOVE) {

		// FUTURE call TranslateMessage(&msg) if the characters
		// are needed instead of just the virtual key codes.

		win.DispatchMessage(&msg) // goes to winProcessMsg
	}
	if input.shutdown {
		wd.dispose()
		return input
	}

	// return collected input to the app. The app needs to
	// check the IsRunning().
	if !input.shutdown {
		// get focus and mouse coordinates
		input.Focus = wd.hwnd == win.GetActiveWindow()
		input.Mx, input.My = wd.cursorLocation()
	}
	return input // singleton for collecting the latest user input.
}

// Get the current mouse position relative to the bottom left corner
// of the application window.
func (wd *windowsDevice) cursorLocation() (mx, my int32) {
	var point win.POINT
	win.GetCursorPos(&point)
	win.ScreenToClient(wd.hwnd, &point)
	var rect win.RECT
	win.GetClientRect(wd.hwnd, &rect)
	return point.X, point.Y
}

// switch between a window with a border and a fullscreen
// window with no border. Expected to be called using F11.
func (wd *windowsDevice) toggleFullscreen() {
	wd.windowed = !wd.windowed
	if wd.windowed {
		// enter bordered window.
		style := uint32(win.WS_OVERLAPPEDWINDOW | win.WS_VISIBLE | win.WS_CLIPCHILDREN | win.WS_MAXIMIZE)
		win.SetWindowLongPtr(wd.hwnd, win.GWL_STYLE, uintptr(style))
		wx, wy, ww, wh := wd.windx, wd.windy, wd.windw, wd.windh
		win.SetWindowPos(wd.hwnd, 0, wx, wy, ww, wh, win.SWP_FRAMECHANGED|win.SWP_SHOWWINDOW)
	} else {
		// enter fullscreen window.
		style := uint32(win.WS_POPUP | win.WS_VISIBLE)
		win.SetWindowLongPtr(wd.hwnd, win.GWL_STYLE, uintptr(style))
		ww, wh := wd.fullw, wd.fullh
		win.SetWindowPos(wd.hwnd, 0, 0, 0, ww, wh, win.SWP_FRAMECHANGED|win.SWP_SHOWWINDOW)
	}
}

// =============================================================================

// Windows virtual key codes. Map Windows key codes to Vu key codes.
//
//	http://msdn.microsoft.com/en-ca/library/windows/desktop/dd375731(v=vs.85).aspx
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
	KF10 = 0x79 // VK_F10       F10 key  ---- on macos
	KF11 = 0x7A // VK_F11       F11 key
	KF12 = 0x7B // VK_F12       F12 key
	KF13 = 0x7C // VK_F13       F13 key
	KF14 = 0x2C // VK_F14 0x7D  F14 key  0x2C on macos
	KF15 = 0x91 // VK_F15 0x7E  F15 key  0x91
	KF16 = 0x13 // VK_F16 0x7F  F16 key  0x13
	KF17 = 0x80 // VK_F17       F17 key
	KF18 = 0x81 // VK_F18       F18 key
	KF19 = 0x82 // VK_F19       F19 key
	KF20 = 0x83 // VK_F20       F20 key

	// Keypad keys
	KPDot = 0x6E // VK_DECIMAL   Decimal key    :: VK_DELETE
	KPMlt = 0x6A // VK_MULTIPLY  Multiply key
	KPAdd = 0x6B // VK_ADD       Add key
	KPClr = 0x90 // VK_CLEAR     0x0C CLEAR key :: VK_OEM_CLEAR 0xFE 0x90 on macos
	KPDiv = 0x6F // VK_DIVIDE    Divide key
	KPEnt = 0x2B // VK_EXECUTE                  :: VK_ENTER on macos
	KPSub = 0x6D // VK_SUBTRACT      Subtract key
	KPEql = 0xE2 //
	KP0   = 0x60 // VK_NUMPAD0  0x60 keypad 0 key :: VK_INSERT  0x20 on macos
	KP1   = 0x61 // VK_NUMPAD1  0x61 keypad 1 key :: VK_END     0x23 on macos
	KP2   = 0x62 // VK_NUMPAD2  0x62 keypad 2 key :: VK_DOWN    0x28 on macos
	KP3   = 0x63 // VK_NUMPAD3  0x63 keypad 3 key :: VK_NEXT    0x22 on macos
	KP4   = 0x64 // VK_NUMPAD4  0x64 keypad 4 key :: VK_LEFT    0x25 on macos
	KP5   = 0x65 // VK_NUMPAD5  0x65 keypad 5 key :: VK_CLEAR   0x0C on macos
	KP6   = 0x66 // VK_NUMPAD6  0x66 keypad 6 key :: VK_RIGHT   0x27 on macos
	KP7   = 0x67 // VK_NUMPAD7  0x67 keypad 7 key :: VK_HOME    0x26 on macos
	KP8   = 0x68 // VK_NUMPAD8  0x68 keypad 8 key :: VK_UP      0x21 on macos
	KP9   = 0x69 // VK_NUMPAD9  0x69 keypad 9 key :: VK_PRIOR

	// Misc and Punctuation keys.
	KEqual = 0xBB         //
	KMinus = 0xBD         // VK_OEM_MINUS  For any country/region, the '-' key // VK_SEPARATOR 0x6C Separator key
	KLBkt  = 0xDB         // VK_OEM_4      misc characters; varys: US keyboard, the '[{' key
	KRBkt  = 0xDD         // VK_OEM_6      misc characters; varys: US keyboard, the ']}' key
	KQuote = win.VK_OEM_7 // VK_OEM_7      misc characters; varys: US keyboard, the 'single/double-quote' key
	KSemi  = 0xBA         // VK_OEM_1      misc characters; varys: US keyboard, the ';:' key
	KBSl   = win.VK_OEM_5 // VK_OEM_5      misc characters; varys: US keyboard, the '/?' key
	KComma = 0xBC         // VK_OEM_COMMA  For any country/region, the ',' key
	KSlash = 0xBF         //
	KDot   = 0xBE         // VK_OEM_PERIOD For any country/region, the '.' key
	KGrave = win.VK_OEM_3 // misc characters; varys: US keyboard, the '`~' key
	KRet   = 0x0D         // VK_RETURN     ENTER key
	KTab   = 0x09         // VK_TAB        TAB key
	KSpace = 0x20         // VK_SPACE      SPACEBAR
	KDel   = 0x08         // VK_BACK       BACKSPACE key
	KEsc   = 0x1B         // VK_ESCAPE     ESC key

	// Control keys.
	KHome   = 0x24           // VK_HOME    HOME key
	KPgUp   = 0x21           // VK_PRIOR   PAGE UP key
	KFDel   = 0x2E           // VK_DELETE  DEL key
	KEnd    = 0x23           // VK_END     END key
	KPgDn   = 0x22           // VK_NEXT    PAGE DOWN key
	KALeft  = 0x25           // VK_LEFT    LEFT ARROW key
	KARight = 0x27           // VK_RIGHT   RIGHT ARROW key
	KADown  = 0x28           // VK_DOWN    DOWN ARROW key
	KAUp    = 0x26           // VK_UP      UP ARROW key
	KCtl    = win.VK_CONTROL // modifier masks and key codes.
	KFn     = 0              // Did not find on windows.
	KShift  = win.VK_SHIFT
	KCmd    = win.VK_LWIN | win.VK_RWIN
	KAlt    = win.VK_MENU

	// Mouse buttons are treated like keys.
	// Values don't conflict with other key codes.
	KML = win.VK_LBUTTON // 0x01 Left mouse button
	KMM = win.VK_MBUTTON // 0x04 Middle mouse button
	KMR = win.VK_RBUTTON // 0x02 Right mouse button
)
