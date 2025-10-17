// Copyright © 2025 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

//go:build ios

package ios

// The iOS native layer. This wraps the c functions that wrap the objective-c
// code that calls the iOS libraries (where the real work is done).

// // This block is C code and cgo directives.
//
// #cgo ios CFLAGS: -x objective-c
// #cgo ios LDFLAGS: -framework Foundation -framework UIKit -framework Metal -framework MetalKit -framework QuartzCore
//
// #import <stdlib.h>
// #import <MetalKit/MetalKit.h>
// #import <UIKit/UIKit.h>
// #include "os_ios.h"
import "C" // must be located here.

import "unsafe"

// Run starts the platform render loop and does not return.
func Run(renderCallback func(), inputCallback func(ev int64, d0, d1 int), viewCallback func(view uintptr)) {
	renderCB = renderCallback
	inputCB = inputCallback
	viewCB = viewCallback
	C.dev_run() // does not return!
}

// Called when the display needs renderering.
// Generally related to the display refresh rate.
//
//export renderFrame
func renderFrame() {
	if renderCB != nil {
		renderCB()
	}
}

// Called on user input events.
//
//export handleInput
func handleInput(event int64, d0, d1 int) {
	if inputCB != nil {
		inputCB(event, d0, d1)
	}
}

// Called on startup to save a pointer to the Metal view layer.
//
//export setView
func setView(viewPointer int64) {
	if viewCB != nil {
		viewCB(uintptr(viewPointer))
	}
}

// native layer callback functions. native->app calls.
// =============================================================================
// iOS Device implementation. app->native calls.

// The callbacks for communicating with the vu/device layer.
var (
	renderCB func()                        // render callback function.
	inputCB  func(event int64, d0, d1 int) // input callback function.
	viewCB   func(view uintptr)            // pointer to display view
)

// Handled by iOS.
func SurfaceSize() (w, h uint32) {
	var ww, wh int32
	C.dev_size((*C.int)(&ww), (*C.int)(&wh))
	return uint32(ww), uint32(wh)
}

const (
	EVENT_TOUCH_BEGIN  = C.devTouchBegin
	EVENT_TOUCH_MOVE   = C.devTouchMove
	EVENT_TOUCH_END    = C.devTouchEnd
	EVENT_RESIZED      = C.devResized
	EVENT_FOCUS_GAINED = C.devFocusIn
	EVENT_FOCUS_LOST   = C.devFocusOut
)

// =============================================================================
// IOSConsoleWriter

// IOSWriter wraps regular logging so the logs appear in the console.
type IOSConsoleWriter struct{}

// Write provides a writer for regular logging.
func (IOSConsoleWriter) Write(p []byte) (n int, err error) {
	cstr := C.CString(string(p))
	C.dev_log(cstr)
	C.free(unsafe.Pointer(cstr))
	return len(p), nil
}
