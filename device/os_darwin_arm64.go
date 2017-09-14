// Copyright © 2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package device

// The iOS (darwin) native layer. This wraps the c functions that wrap the
// objective-c code that calls the iOS libraries (where the real work is done).
//
// Cross compile test can be run on macOS as follows:
// env GOOS=darwin GOARCH=arm64 CC=/Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/clang CXX=/Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/clang CGO_CFLAGS="-isysroot /Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/Developer/SDKs/iPhoneOS10.3.sdk -arch arm64 -miphoneos-version-min=8.0" CGO_LDFLAGS="-isysroot /Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/Developer/SDKs/iPhoneOS10.3.sdk -arch arm64 -miphoneos-version-min=8.0" CGO_ENABLED=1 go build
// Be-aware that this creates a cross compiled version of the go toolchain
// each time, so it can be a bit slow.

// // The following block is C code and cgo directvies.
//
// #cgo darwin CFLAGS: -x objective-c -fno-common
// #cgo darwin LDFLAGS: -framework Foundation -framework UIKit -framework GLKit -framework OpenGLES
//
// #include <stdlib.h>
// #include <UIKit/UIDevice.h>
// #include <GLKit/GLKit.h>
// #include "os_darwin_arm64.h"
import "C" // must be located here.

import (
	"log"
	"runtime"
	"unsafe"
)

// iOS specific. Otherwise the app freezes within seconds of creation.
func init() { runtime.LockOSThread() }

// Device instance needed to handle callbacks on exported methods.
var dev = &ios{}

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
// The summary is fowarded with each update to the controlling application.
//    d0, d1 are optional event specific data.
//
//export handleInput
func handleInput(event int64, d0, d1 int) {
	if dev.input == nil {
		return // Ignore input before window is up.
	}
	in := dev.input
	switch event {
	case C.devTouchBegin:
		in.recordPress(KLm)
		dev.tx, dev.ty = d0, d1
	case C.devTouchMove:
		dev.tx, dev.ty = d0, d1
	case C.devTouchEnd:
		in.recordRelease(KLm)
		dev.tx, dev.ty = d0, d1
	case C.devResize:
		in.curr.Resized = true
	case C.devFocus:
		in.curr.Focus = d0 == 0
	default:
		log.Printf("handleInput unknown %d", event)
	}
}

// native layer callback functions. native->app calls.
// =============================================================================
// iOS Device implementation. app->native calls.

// ios is the iOS implementation of the Device interface.
// See the Device interface for method descriptions.
type ios struct {
	app   App    // update/render callback.
	input *input // tracks current keys pressed.

	// pixel location of first touch when a touch is active.
	tx, ty int // 0, 0 when no touch is active.
}

// Handled by iOS.
func (os *ios) Dispose()       { C.dev_dispose() }
func (os *ios) Down() *Pressed { return os.input.getPressed(os.tx, os.ty) }
func (os *ios) Size() (x, y, w, h int) {
	var ww, wh, scale int32
	C.dev_size((*C.int)(&ww), (*C.int)(&wh), (*C.int)(&scale))
	return 0, 0, int(ww), int(wh)
}

// Mobile devices rely on touch.
const (
	TouchBegin = C.devTouchBegin // touch.TypeBegin
	TouchMove  = C.devTouchMove  // touch.TypeMove
	TouchEnd   = C.devTouchEnd   // touch.TypeEnd
)

// Ignored for iOS.
func (os *ios) ShowCursor(show bool)   {}
func (os *ios) SetCursorAt(x, y int)   {}
func (os *ios) SetSize(x, y, w, h int) {}
func (os *ios) Copy() string           { return "" }
func (os *ios) Paste(s string)         {}
func (os *ios) SwapBuffers()           {}
func (os *ios) SetTitle(t string)      {}
func (os *ios) IsFullScreen() bool     { return false }
func (os *ios) ToggleFullScreen()      {}

// =============================================================================
// iosLogger

// iosLogger wraps regular logging so the logs appear in the ios console.
type iosLogger struct{}

// Write provides a writer for regular logging.
func (iosLogger) Write(p []byte) (n int, err error) {
	cstr := C.CString(string(p))
	C.dev_log(cstr)
	C.free(unsafe.Pointer(cstr))
	return len(p), nil
}

// Ensure logging is redirected at startup.
func init() {
	log.SetOutput(iosLogger{})
}

// iosLogger
// =============================================================================

// FUTURE: improve design. Currently the engine depends on these constants
//         but there are no plans to use the ios keyboard as a game controller.
//         Also a single touch is being exposed as a left click - ugh.
const (
	K0     = iota // Standard keyboard numbers.
	K1            //
	K2            //
	K3            //
	K4            //
	K5            //
	K6            //
	K7            //
	K8            //
	K9            //
	KA            // Standard keyboard letters.
	KB            //
	KC            //
	KD            //
	KE            //
	KF            //
	KG            //
	KH            //
	KI            //
	KJ            //
	KK            //
	KL            //
	KM            //
	KN            //
	KO            //
	KP            //
	KQ            //
	KR            //
	KS            //
	KT            //
	KU            //
	KV            //
	KW            //
	KX            //
	KY            //
	KZ            //
	KF1           // General Function Keys
	KF2           //
	KF3           //
	KF4           //
	KF5           //
	KF6           //
	KF7           //
	KF8           //
	KF9           //
	KF10          //
	KF11          //
	KF12          //
	KF13          //
	KF14          //
	KF15          //
	KF16          //
	KF17          //
	KF18          //
	KF19          //
	KF20          //
	KKpDot        // Extended keyboard keypad keys
	KKpMlt        //   "
	KKpAdd        //   "
	KKpClr        //   "
	KKpDiv        //   "
	KKpEnt        //   "
	KKpSub        //   "
	KKpEql        //   "
	KKp0          //   "
	KKp1          //   "
	KKp2          //   "
	KKp3          //   "
	KKp4          //   "
	KKp5          //   "
	KKp6          //   "
	KKp7          //   "
	KKp8          //   "
	KKp9          //   "
	KEqual        // Standard keyboard punctuation keys.
	KMinus        //   "
	KLBkt         //   "
	KRBkt         //   "
	KQt           //   "
	KSemi         //   "
	KBSl          //   "
	KComma        //   "
	KSlash        //   "
	KDot          //   "
	KGrave        //   "
	KRet          //   "
	KTab          //   "
	KSpace        //   "
	KDel          //   "
	KEsc          //   "
	KHome         // Specific function keys.
	KPgUp         //   "
	KFDel         //   "
	KEnd          //   "
	KPgDn         //   "
	KLa           // Arrow keys
	KRa           //   "
	KDa           //   "
	KUa           //   "

	// Note that modifier masks are also used as key values
	// since they don't conflict with the standard key codes.
	KCtl   // Modifier masks and key codes.
	KFn    //   "
	KShift //   "
	KCmd   //   "
	KAlt   //   "

	// Mouse buttons are treated like keys. Values don't
	// conflict with other key codes.
	KMm        //   "
	KRm        //   "
	KLm = 0xA0 // Use left mouse for touch for now.
)
