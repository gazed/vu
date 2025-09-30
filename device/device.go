package device

// device.go wraps the platform specific functionality.

import (
	"time"
)

// New creates the platform for the current host.
func New(windowed bool, title string, x int32, y int32, w int32, h int32) *Device {
	// newPlatform is implemented by each platform.
	// FUTURE: provide platforms for linux, macos, etc.
	d := &Device{platform: newPlatform()}

	// just init the underlying data structs. Use CreateDisplay() as the next step.
	d.platform.init(windowed, title, x, y, w, h)
	return d
}

// Device provides access to a platform specific user input
// and display surface.
type Device struct {
	platform platformAPI
}

// does not return while running
// Used on platforms that require a render callback loop.
func (d *Device) Run(renderCallback func()) {
	d.platform.run(renderCallback)
}

// CreateDisplay initializes a window. Called once at initialization.
func (d *Device) CreateDisplay() error {
	return d.platform.createDisplay()
}

// Returns the size of the surface in pixels.
func (d *Device) SurfaceSize() (w, h uint32) {
	return d.platform.surfaceSize()
}

// Returns the upper corner location of the surface in pixels.
// Consistent with the values provided in device.New().
func (d *Device) SurfaceLocation() (x, y int32) {
	return d.platform.surfaceLocation()
}

// GetInput updates the user input. Expected to be called frequently,
// ie: once before each game update.
func (d *Device) GetInput() *Input {
	return d.platform.getInput()
}

// Dispose releases all platform specific resources.
func (d *Device) Dispose() {
	d.platform.dispose()
}

// IsRunning returns true if the app is still running.
func (d *Device) IsRunning() bool {
	return d.platform.isRunning()
}

// SetResizeHandler registers a callback to do stuff any time
// the window is resized.
func (d *Device) SetResizeHandler(callback func()) {
	d.platform.setResizeHandler(callback)
}

// ToggleFullscreen toggles between windowed with border
// and windowed fullscreen with no border.
func (d *Device) ToggleFullscreen() {
	d.platform.toggleFullscreen()
}

// platformAPI is the interface that each platform must implement.
// One platform will be active on startup.
type platformAPI interface {

	// initialize the platform.
	init(windowed bool, title string, x, y, w, h int32)
	run(renderCallback func())

	// exposed as Device public methods.
	createDisplay() error             // see CreateDisplay
	surfaceSize() (w, h uint32)       // see SurfaceSize
	surfaceLocation() (x, y int32)    // see SurfaceLocation
	getInput() *Input                 // see GetInput
	dispose()                         // see Dispose
	isRunning() bool                  // see IsRunning
	setResizeHandler(callback func()) // see SetResizeHandler
	toggleFullscreen()                // see ToggleFullscreen
}

// =============================================================================

// Input gathers user input and organizes it for mapping to game actions.
// Input data is shared with the app and the app should treat the data
// as read only.
type Input struct {
	Mx, My int32 // Mouse location relative to top left.
	Scroll int   // Scroll amount: positive, negative, or Zero if no scrolling.
	Focus  bool  // True if window has focus.

	// Pressed are keys that were pressed since last request.
	Pressed map[int32]bool //

	// Down are keys that are currently being pressed.
	// This includes keys that were just pressed and which
	// are being held down. The value can be used to calculate
	// how long the key has been down.
	// Down keys are cleared when the key is
	// released or when the windows looses focus.
	Down map[int32]time.Time // time when key was pressed.

	// Released are keys that were just released since last request.
	// Keys are also released when the window loses focus.
	Released map[int32]time.Duration // total time down.

	// internal signal for when the user has closed the window.
	shutdown bool // true when user closes window.
}

// reset prepares input data for a refresh.
func (in *Input) reset() {
	in.Mx = 0     // to be refreshed with current mouse location.
	in.My = 0     // ""
	in.Scroll = 0 // no scrolling happening.
	// Focus is kept or set in getInput.

	// clear the Pressed and Released as they are a one time notification.
	// The Down keys are kept until they are released.
	for k, _ := range in.Pressed {
		delete(in.Pressed, k)
	}
	for k, _ := range in.Released {
		delete(in.Released, k)
	}
}

// keyReleased records keys that are no longer down in the last poll.
func (in *Input) keyReleased(k int32) {
	if v, ok := in.Down[k]; ok {
		in.Released[k] = time.Since(v) // one time event.
		delete(in.Down, k)             // key is no longer down.
	}
}

// keyPressed records keys that have just been pressed in the last poll.
// Ignore repeat pressed events for keys that are already down.
func (in *Input) keyPressed(k int32) {
	if _, ok := in.Down[k]; !ok {
		in.Pressed[k] = true    // one time event.
		in.Down[k] = time.Now() // key is down.
	}
}

// loseFocus clears the down keys since once focus is lost the key
// release events will be missed.
func (in *Input) loseFocus() {
	input.reset()
	for k, v := range in.Down {
		in.Released[k] = time.Since(v) // inform app
		delete(in.Down, k)             // key is no longer down.
	}
	in.Focus = false
}

// allMouseButtonsReleased returns true if none of the mouse buttons
// are currently held down.
func (in *Input) allMouseButtonsReleased() bool {
	_, leftMouseDown := in.Down[KML]
	_, middleMouseDown := in.Down[KMM]
	_, rightMouseDown := in.Down[KMR]
	return !leftMouseDown && !middleMouseDown && !rightMouseDown
}
