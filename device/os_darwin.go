package device

// os_darwin.go wraps the macos native layer.
// FUTURE be specfic ie: //go:build darwin && arm64

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"

	// binding layer
	"github.com/gazed/vu/internal/device/macos"
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
	if md, ok := d.platform.(*macosDevice); ok {
		return md.display, nil
	}
	return nil, fmt.Errorf("GetRenderSurfaceInfo: invalid device")
}

// =============================================================================
// newPlatform returns the platform when running on a windows system.
func newPlatform() platformAPI { return &macosDevice{} }

// macosDevice holds display state and implements the Device interface
// for the macos platform.
type macosDevice struct {
	display    unsafe.Pointer // *MTKView fetched from binding layer.
	windowed   bool           // window bordered vs full screen
	title      string         // window with border title
	x, y, w, h uint32         // startup window sizes.
}

func (md *macosDevice) init(windowed bool, title string, x, y, w, h int32) {
	// https://stackoverflow.com/questions/25361831/benefits-of-runtime-lockosthread-in-golang
	// "With the Go threading model, calls to C code, assembler code, or blocking
	//  system calls occur in the same thread as the calling Go code, which is
	//  managed by the Go runtime scheduler. The os.LockOSThread() mechanism is
	//  mostly useful when Go has to interface with some foreign library
	//  (a C library for instance). It guarantees that several successive calls
	//  to this library will be done in the same thread."
	runtime.LockOSThread()
	md.windowed = windowed
	md.title = title
	md.x, md.y, md.w, md.h = uint32(x), uint32(y), uint32(w), uint32(h)
}

// run does not return.
func (md *macosDevice) run(renderCallback func()) {
	macos.Run(renderCallback, darwinInputHandler)
}

// createDisplay with the application title.
func (md *macosDevice) createDisplay() error {
	md.display = macos.CreateDisplay(md.title, md.x, md.y, md.w, md.h)
	return nil
}
func (md *macosDevice) surfaceSize() (w, h uint32)    { return macos.SurfaceSize() }
func (md *macosDevice) surfaceLocation() (x, y int32) { return macos.SurfaceLocation() }
func (md *macosDevice) toggleFullscreen()             { macos.ToggleFullscreen() }

// dispose results in apple main loop killing the process.
// This method does not return
func (md *macosDevice) dispose() {
	md.display = nil
	macos.Dispose() // kills the process.
}
func (md *macosDevice) isRunning() bool { return md.display != nil }

// =============================================================================
// user input handling.
var (

	// Allow app, mainly test apps, to override the default input handler.
	inputHandlerHook func(event, data int64) = nil

	// shared singleton input data returned to engine.
	input = &Input{
		Pressed:  map[int32]bool{},
		Down:     map[int32]time.Time{},
		Released: map[int32]time.Duration{},
	}

	// collect user events until update loop requests them.
	inputEvents = []int64{}

	// resizeHandler processes resize events immediately.
	resizeHandler func() = nil
)

// set by engine on startup.
func (md *macosDevice) setResizeHandler(callback func()) { resizeHandler = callback }

// SetInputHandler expected to be called once on startup.
func SetInputHandler(handler func(event, data int64)) {
	inputHandlerHook = handler
}

// darwinInputHandler consolidates user input until it is requested
// by the engine. Apple user input is delivered by callback.
// The engine requests the user input from the main run loop.
// For Apple devices the main run loop is triggered by a render callback.
func darwinInputHandler(event, data int64) {
	// redirect input to the the hook if it exists.
	// Generally this is for testing, or apps that want direct
	// osx input.
	if inputHandlerHook != nil {
		inputHandlerHook(event, data)
		return
	}

	// resized events are handled immediately.
	if event == macos.EVENT_RESIZED || event == macos.EVENT_MOVED {
		resizeHandler()
		return
	}

	// default input event engine processing.
	// collect the event based input until the engine update loop
	// requests it.
	inputEvents = append(inputEvents, event, data)
}

// getInput is called by the update. It process the input collected by
// the event driven darwinInputHandler.
func (md *macosDevice) getInput() *Input {
	input.reset() // clear the shared input.

	// process collected input events.
	for i := 0; i < len(inputEvents); i += 2 {
		event, data := inputEvents[i], inputEvents[i+1]

		switch event {
		case macos.EVENT_KEYUP:
			input.keyReleased(int32(data))
		case macos.EVENT_KEYDOWN:
			input.keyPressed(int32(data))
		case macos.EVENT_SCROLL:
			input.Scroll += int(data)
		case macos.EVENT_MODIFIER:
			if data&macos.KAlt != 0 {
				input.keyPressed(macos.KAlt)
			} else {
				input.keyReleased(macos.KAlt)
			}
		case macos.EVENT_MOVED:
			// should have been already handled
		case macos.EVENT_RESIZED:
			// should have been already handled
		case macos.EVENT_FOCUS_GAINED:
			input.Focus = true
		case macos.EVENT_FOCUS_LOST:
			input.loseFocus() // also clears keys.
		}
	}
	inputEvents = inputEvents[:0] // clear data, keep memory

	// get the current mouse position relative to top left of the window.
	input.Mx, input.My = macos.MousePosition() // relative to bottom left.
	_, ymax := macos.SurfaceSize()             // .. so get size
	input.My = int32(ymax) - input.My          // .. and convert to top left.

	// return collected input to the engine.
	return input
}

// =============================================================================

// Expose the device layer macos user events.
const (
	EVENT_KEYUP        = macos.EVENT_KEYUP
	EVENT_KEYDOWN      = macos.EVENT_KEYDOWN
	EVENT_SCROLL       = macos.EVENT_SCROLL
	EVENT_MODIFIER     = macos.EVENT_MODIFIER
	EVENT_MOVED        = macos.EVENT_MOVED
	EVENT_RESIZED      = macos.EVENT_RESIZED
	EVENT_FOCUS_GAINED = macos.EVENT_FOCUS_GAINED
	EVENT_FOCUS_LOST   = macos.EVENT_FOCUS_LOST
)

// Expose the device layer macos keys.
const (
	K0      = macos.K0
	K1      = macos.K1
	K2      = macos.K2
	K3      = macos.K3
	K4      = macos.K4
	K5      = macos.K5
	K6      = macos.K6
	K7      = macos.K7
	K8      = macos.K8
	K9      = macos.K9
	KA      = macos.KA
	KB      = macos.KB
	KC      = macos.KC
	KD      = macos.KD
	KE      = macos.KE
	KF      = macos.KF
	KG      = macos.KG
	KH      = macos.KH
	KI      = macos.KI
	KJ      = macos.KJ
	KK      = macos.KK
	KL      = macos.KL
	KM      = macos.KM
	KN      = macos.KN
	KO      = macos.KO
	KP      = macos.KP
	KQ      = macos.KQ
	KR      = macos.KR
	KS      = macos.KS
	KT      = macos.KT
	KU      = macos.KU
	KV      = macos.KV
	KW      = macos.KW
	KX      = macos.KX
	KY      = macos.KY
	KZ      = macos.KZ
	KF1     = macos.KF1
	KF2     = macos.KF2
	KF3     = macos.KF3
	KF4     = macos.KF4
	KF5     = macos.KF5
	KF6     = macos.KF6
	KF7     = macos.KF7
	KF8     = macos.KF8
	KF9     = macos.KF9
	KF10    = macos.KF10
	KF11    = macos.KF11
	KF12    = macos.KF12
	KF13    = macos.KF13
	KF14    = macos.KF14
	KF15    = macos.KF15
	KF16    = macos.KF16
	KF17    = macos.KF17
	KF18    = macos.KF18
	KF19    = macos.KF19
	KF20    = macos.KF20
	KPDot   = macos.KPDot
	KPMlt   = macos.KPMlt
	KPAdd   = macos.KPAdd
	KPClr   = macos.KPClr
	KPDiv   = macos.KPDiv
	KPEnt   = macos.KPEnt
	KPSub   = macos.KPSub
	KPEql   = macos.KPEql
	KP0     = macos.KP0
	KP1     = macos.KP1
	KP2     = macos.KP2
	KP3     = macos.KP3
	KP4     = macos.KP4
	KP5     = macos.KP5
	KP6     = macos.KP6
	KP7     = macos.KP7
	KP8     = macos.KP8
	KP9     = macos.KP9
	KEqual  = macos.KEqual
	KMinus  = macos.KMinus
	KLBkt   = macos.KLBkt
	KRBkt   = macos.KRBkt
	KQuote  = macos.KQuote
	KSemi   = macos.KSemi
	KBSl    = macos.KBSl
	KComma  = macos.KComma
	KSlash  = macos.KSlash
	KDot    = macos.KDot
	KGrave  = macos.KGrave
	KRet    = macos.KRet
	KTab    = macos.KTab
	KSpace  = macos.KSpace
	KDel    = macos.KDel
	KEsc    = macos.KEsc
	KHome   = macos.KHome
	KPgUp   = macos.KPgUp
	KFDel   = macos.KFDel
	KEnd    = macos.KEnd
	KPgDn   = macos.KPgDn
	KALeft  = macos.KALeft
	KARight = macos.KARight
	KADown  = macos.KADown
	KAUp    = macos.KAUp

	// Note that modifier masks are also used as key values
	// since they don't conflict with the standard key codes.
	KCtl   = macos.KCtl
	KFn    = macos.KFn
	KShift = macos.KShift
	KCmd   = macos.KCmd
	KAlt   = macos.KAlt

	// Mouse buttons are treated like keys. Values don't
	// conflict with other key codes.
	KML = macos.KML
	KMM = macos.KMM
	KMR = macos.KMR
)
