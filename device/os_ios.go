// SPDX-FileCopyrightText : © 2025 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

//go:build ios

package device

// os_ios.go wraps the ios native layer.

import (
	"fmt"
	"time"
	"unsafe"

	// binding layer
	"github.com/gazed/vu/internal/device/ios"
)

// GetRenderSurfaceInfo exposes the windows API specific information
// needed by the render package to create a rendering surface.
// Called by the platform specific code in the Render package.
//
// from the MoltenVK Runtime UserGuide:
// "When creating a CAMetalLayer to underpin the Vulkan surface to render to,
// it is strongly recommended that you ensure the delegate of the CAMetalLayer
// is the NSView/UIView in which the layer is contained, to ensure correct and
// optimized Vulkan swapchain and refresh timing behavior across multiple display
// screens that might have different properties."
func GetRenderSurfaceInfo(d *Device) (display unsafe.Pointer, err error) {
	if md, ok := d.platform.(*iosDevice); ok {
		return md.display, nil
	}
	return nil, fmt.Errorf("GetRenderSurfaceInfo: invalid device")
}

// =============================================================================
// singleton for the ios platform.
var iosPlatform = &iosDevice{}

// newPlatform returns the platform when running on ios.
func newPlatform() platformAPI { return iosPlatform }

// iosDevice holds display state and implements the Device interface
// for the macos platform.
type iosDevice struct {
	display      unsafe.Pointer // *MTKView fetched from binding layer.
	initCallback func()         // engine initialization callback
	windowed     bool           // window bordered vs full screen
	title        string         // window with border title
	x, y, w, h   uint32         // startup window sizes.
}

func (id *iosDevice) init(windowed bool, title string, x, y, w, h int32) {
	id.windowed = windowed
	id.title = title
	id.x, id.y, id.w, id.h = 0, 0, uint32(w), uint32(h)
}

var startupCallback func()

func startupHandler(viewPointer uintptr) {
	iosPlatform.display = unsafe.Pointer(viewPointer)
	startupCallback() // finish engine initialization.
}

// run does not return.
func (id *iosDevice) run(renderCallback, initCallback func()) {
	id.initCallback = initCallback
	startupCallback = initCallback
	ios.Run(renderCallback, iosInputHandler, startupHandler)
}

// createDisplay does nothing for IOS since the display is created
// immediately in startup.
func (id *iosDevice) createDisplay() error       { return nil }
func (id *iosDevice) surfaceSize() (w, h uint32) { return ios.SurfaceSize() }

// ignored for IOS
func (id *iosDevice) surfaceLocation() (x, y int32) { return 0, 0 }

// ignored for IOS.
func (id *iosDevice) toggleFullscreen() {}

// dispose results in apple main loop killing the process.
// This method does not return
func (id *iosDevice) dispose() {
	id.display = nil
}
func (id *iosDevice) isRunning() bool { return id.display != nil }

// =============================================================================
// startup handling.

// called by the device layer on startup to set the display pointer
// needed by the render layer.
func (id *iosDevice) initHandler(viewPointer uintptr) {
	id.display = unsafe.Pointer(viewPointer)
	id.initCallback() // finish engine initialization.
}

// =============================================================================
// user input handling.
var (
	// shared singleton input data returned to engine.
	input = &Input{
		Pressed:  map[int32]bool{},
		Down:     map[int32]time.Time{},
		Released: map[int32]time.Duration{},
	}

	// collect user events until update loop requests them.
	inputEvents = []int{}

	// resizeHandler processes resize events immediately.
	resizeHandler func() = nil
)

// set by engine on startup.
func (id *iosDevice) setResizeHandler(callback func()) { resizeHandler = callback }

// iosInputHandler consolidates user input until it is requested
// by the engine. Apple user input is delivered by callback.
// The engine requests the user input from the main run loop.
// For Apple devices the main run loop is triggered by a render callback.
func iosInputHandler(event int64, d0, d1 int) {

	// resized events are handled immediately.
	if event == ios.EVENT_RESIZED {
		resizeHandler()
		return
	}

	// default input event engine processing.
	// collect event based input until the engine update loop requests it.
	inputEvents = append(inputEvents, int(event), d0, d1)
}

// hack to provide the mouse location for an active touch.
var lastMouseX int32
var lastMouseY int32

// getInput is called by the update. It process the input collected by
// the event driven input handler.
func (id *iosDevice) getInput() *Input {
	input.reset() // clear the shared input.

	// process collected input events.
	for i := 0; i < len(inputEvents); i += 3 {
		event, d0, d1 := inputEvents[i], inputEvents[i+1], inputEvents[i+2]

		switch event {
		case TOUCH_BEGIN:
			input.keyPressed(TOUCH)
			input.Mx = int32(d0)
			input.My = int32(d1)
			lastMouseX = input.Mx
			lastMouseY = input.My
		case TOUCH_MOVE:
			input.Mx = int32(d0) // update mouse locations.
			input.My = int32(d1)
			lastMouseX = input.Mx
			lastMouseY = input.My
		case TOUCH_END:
			input.keyReleased(TOUCH)
			input.Mx = int32(d0)
			input.My = int32(d1)
			lastMouseX = 0
			lastMouseY = 0
		case ios.EVENT_RESIZED:
			// should have been already handled
			lastMouseX = 0
			lastMouseY = 0
		case ios.EVENT_FOCUS_GAINED:
			input.Focus = true
			lastMouseX = 0
			lastMouseY = 0
		case ios.EVENT_FOCUS_LOST:
			input.loseFocus() // also clears keys.
			lastMouseX = 0
			lastMouseY = 0
		}
	}

	// if there is an existing touch in the down events,
	// then copy the previous mouse position.
	if _, ok := input.Down[TOUCH]; ok {
		input.Mx = lastMouseX
		input.My = lastMouseY
	}

	inputEvents = inputEvents[:0] // clear data, keep memory
	return input                  // return collected input to the engine.
}

// =============================================================================

// Expose the device layer ios user events.
const (
	TOUCH_BEGIN = ios.EVENT_TOUCH_BEGIN
	TOUCH_MOVE  = ios.EVENT_TOUCH_MOVE
	TOUCH_END   = ios.EVENT_TOUCH_END
	TOUCH       = 10
)

// ios fakes the macos keys and will never generate any
// of these events.
const (
	K0 = iota
	K1
	K2
	K3
	K4
	K5
	K6
	K7
	K8
	K9
	KA
	KB
	KC
	KD
	KE
	KF
	KG
	KH
	KI
	KJ
	KK
	KL
	KM
	KN
	KO
	KP
	KQ
	KR
	KS
	KT
	KU
	KV
	KW
	KX
	KY
	KZ
	KF1
	KF2
	KF3
	KF4
	KF5
	KF6
	KF7
	KF8
	KF9
	KF10
	KF11
	KF12
	KF13
	KF14
	KF15
	KF16
	KF17
	KF18
	KF19
	KF20
	KPDot
	KPMlt
	KPAdd
	KPClr
	KPDiv
	KPEnt
	KPSub
	KPEql
	KP0
	KP1
	KP2
	KP3
	KP4
	KP5
	KP6
	KP7
	KP8
	KP9
	KEqual
	KMinus
	KLBkt
	KRBkt
	KQuote
	KSemi
	KBSl
	KComma
	KSlash
	KDot
	KGrave
	KRet
	KTab
	KSpace
	KDel
	KEsc
	KHome
	KPgUp
	KFDel
	KEnd
	KPgDn
	KALeft
	KARight
	KADown
	KAUp

	// Note that modifier masks are also used as key values
	// since they don't conflict with the standard key codes.
	KCtl
	KFn
	KShift
	KCmd
	KAlt

	// Mouse buttons are treated like keys. Values don't
	// conflict with other key codes.
	KML
	KMM
	KMR
)
