// Copyright © 2013-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// Design Notes:
// Big thanks to GLFW (http://www.glfw.org) from which the minimalist API
// philosophy was borrowed along with which OS specific API's mattered.
// Also thank you to https://github.com/golang/mobile.
// FUTURE: Linux support  : ignore X support and wait for Wayland.
//                          Need to pick one distro for main testing.
// FUTURE: Android support: need access to android hardware.

// Package device provides minimal platform/os access to a 3D rendering context
// and user input. Access to user keyboard and mouse or touch input is provided
// through the Pressed structure. The application is responsible for providing
// windowing constructs like buttons, dialogs, sub-panels, text-boxes, etc.
//
// An application is expected to create a device and then call Run.
// Running a device takes control of the processing, calling the application
// back on the App interface methods.
//
// Package device is provided as part of the vu (virtual universe) 3D engine.
package device

// Device wraps OS specific functionality. The expected usage is for
// the application to create the device as follows:
//     app := &Application{}  // client specific application.
//     dev := device.Run(app) // create device and render context.
// where device.Run does not return but calls the following methods:
//     func (a *Application) Init(dev device.Device) {
//         // one time call.
//     }
//     func (a *Application) Refresh(d device.Device) {
//         // regular call to handle user input,
//         // update state, and render a frame.
//     }
// Commonly applications are closed when the user closes the device window
// or app, stopping calls to Refresh. Calls to Refresh may also stop when
// the application loses focus or is put in the background.
type Device interface {

	// Call Down each App.Refresh to process user input.
	Down() *Pressed // Gets pressed keys since last call.
	Dispose()       // Stop device and release OS specific resources.

	// Returns the size and position of the screen. Mobile devices
	// are always full screen where x=y=0.
	Size() (x, y, w, h int) // in pixels where 0,0 is bottom left.

	// The following methods are for to PC devices and
	// are currently ignored or unnecessary for mobile.

	// SwapBuffers exchanges the graphic drawing buffers.
	// Call after each Refresh.
	SwapBuffers() // A frame has been rendered.

	// Copy/Paste interacts with the system clipboard using strings.
	Copy() string   // Returns nil if no string on clipboard.
	Paste(s string) // Paste the given string onto the clipboard.

	// Windowed computer API's.
	// The x,y (0,0) coordinates are at the bottom left.
	// The w,h are width, height dimensions in pixels.
	SetSize(x, y, w, h int) // Used to restore from preferences.
	SetTitle(t string)      // Expected to be called once on startup.
	IsFullScreen() bool     // Returns true if window is full screen.
	ToggleFullScreen()      // Flips between full screen and windowed.
	ShowCursor(show bool)   // Displays or hides the cursor.
	SetCursorAt(x, y int)   // Places the cursor at the given window location.
}

// Run initializes the device specific layer and starts callbacks
// to the given application. This function does not return.
func Run(app App) {
	// runApp and a Device implementation is defined
	// in each native layer.
	runApp(app) // Does not return!!!
}

// App handled device callbacks. An App instance is provided to
// the device.Run method.
type App interface {
	Init(dev Device) // Called once after device initialization.

	// Refresh the display. Called at the display refresh rate which is
	// usually 60 times/second. The current (key/mouse) pressed state is
	// returned for read only access.
	Refresh(dev Device) // Called repeatedly after Init.
}

// Pressed is used to communicate user input. Input mainly consists
// of the keys that are currently pressed and how long they have been
// pressed (measured in update ticks).
type Pressed struct {
	Mx, My  int  // Current mouse location.
	Scroll  int  // The amount of scrolling, if any.
	Focus   bool // True if window has focus.
	Resized bool // True if window was resized or moved.

	// Pressed keys keyCodes and pressed duration in ticks.
	// A postitive duration means the key is still being held down.
	// A negative duration means that the key has been released since
	// the last poll. The total pressed duration prior to release can
	// be determined using the difference with KEY_RELEASED.
	Down map[int]int // Pressed keys and pressed duration.
}

// KeyReleased is used to indicate a key up event has occurred.
// The total duration of a key press can be calculated by the difference
// of Pressed.Down duration with KEY_RELEASED. A user would have to hold
// a key down for 24 hours before the released duration became positive
// (assuming a reasonable update time of 0.02 seconds).
const KeyReleased = -1000000000
