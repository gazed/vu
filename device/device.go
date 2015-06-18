// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// Package device provides minimal platform/os access to a 3D rendering context
// and user input. Access to user keyboard and mouse input is provided through
// the Update method and Pressed structure. The application is responsible
// for providing any windowing constructs like buttons, controls, dialogs,
// sub-panels, text-boxes, etc.
//
// Package device is provided as part of the vu (virtual universe) 3D engine.
package device

// Big thanks to GLFW (http://www.glfw.org) from which the minimalist API
// philosophy was borrowed along with which OS specific API's mattered.

// Device wraps OS specific functionality. The expected usage is:
//     dev := device.New("title", x, y, width, height)
//     // Application initialization code.
//     dev.Open()
//     for dev.IsAlive() {
//         pressed := dev.Update()
//         // Application update and render code.
//         dev.SwapBuffers()
//     }
//     dev.Dispose()
type Device interface {
	Open()                // Open the window and process events.
	ShowCursor(show bool) // Displays or hides the cursor.
	SetCursorAt(x, y int) // Places the cursor at the given window location.
	Dispose()             // Release OS specific resources.

	// IsAlive returns true as long as the window is able to process user input.
	// Quitting the application window will cause IsAlive to return false.
	IsAlive() bool

	// Size returns the usable graphics context location and size excluding any OS
	// specific window trim. The window x, y (0, 0) coordinates are at the bottom
	// left of the window.
	Size() (x, y, width, height int)
	IsFullScreen() bool // Returns true if window is full screen.
	ToggleFullScreen()  // Flips between full screen and windowed mode.

	// SwapBuffers exchanges the graphic drawing buffers. Expected to be called
	// after completing a render. All rendering contexts are double buffered.
	SwapBuffers()

	// Update returns the current (key/mouse) pressed state. The calling application
	// is expected to:
	//    1. Treat the pressed information as read only.
	//    2. Call this method every update loop. This method is responsible
	//       for regular processing of the native OS window events.
	Update() *Pressed
}

// Pressed is used to communicate current user input. Input mainly consists of
// the list of keys that are currently being pressed and how long they have been
// pressed (measured in update ticks).
// A postitive duration means the key is still being held down.
// A negative duration means that the key has been released since
// the last poll. The total pressed duration prior to release can be
// determined using the difference with KEY_RELEASED.
type Pressed struct {
	Mx, My  int         // Current mouse location.
	Scroll  int         // The amount of scrolling, if any.
	Down    map[int]int // Pressed keys and pressed duration.
	Focus   bool        // True if window has focus.
	Resized bool        // True if window was resized or moved.
}

// Device interfaces
// ===========================================================================
// device provides default Device implementation.

// New provides a newly initialized Device with an underlying window and
// graphics context created, but not yet displayed. The only thing left to
// do is to open the device and start polling it for user input.
func New(title string, x, y, width, height int) Device { return newDevice(title, x, y, width, height) }

// Design note 1: the layers in this package are:
//     device : uses context to process events into something simple.
//     input  : turn user input event stream into pollable structure.
//     native : single point of entry into the native layer.
//     os_darwin : OSX native layer. Wraps the following.
//        os_darwin.m : objective-c code wraps cocoa.
//        os_darwin.h
//        os_darwin_test.m
//     os_windows: Win native layer. Wraps the following.
//        os_windows.c : c code wrapping windows API.
//        os_windows.h
//        os_windows_test.c
//     os_linux  : FUTURE: Linux native layer. Wraps the following.
//        os_linux.cpp : c++ code wrapping linux API.
//        os_linux.h
//        os_linux_test.c
//
// Design note 2: user events need to be processed on the main thread for OSX.
//                See: native::readAndDispatch

// device provides a simplification layer over the more raw context layer.
type device struct {
	os    *nativeOs // Native layer wrapper.
	input *input    // User input handler.
}

// newDevice initializes a OS specific window with a valid render context.
func newDevice(title string, x, y, width, height int) *device {
	d := &device{}
	d.os = newNativeOs()
	d.os.createDisplay(title, x, y, width, height)
	d.os.createShell()
	depthBufferBits, alphaBits := 24, 8 // resonable defaults
	d.os.createContext(depthBufferBits, alphaBits)
	d.input = newInput()
	return d
}

// Access the device specific information in a consistent and general manner.
func (d *device) Open()                           { d.os.openShell() }
func (d *device) Dispose()                        { d.os.dispose() }
func (d *device) IsAlive() bool                   { return d.os.isAlive() }
func (d *device) Size() (x, y, width, height int) { return d.os.size() }
func (d *device) ShowCursor(show bool)            { d.os.showCursor(show) }
func (d *device) SwapBuffers()                    { d.os.swapBuffers() }
func (d *device) IsFullScreen() bool              { return d.os.isFullscreen() }
func (d *device) ToggleFullScreen()               { d.os.toggleFullscreen() }
func (d *device) SetCursorAt(x, y int)            { d.os.setCursorAt(x, y) }
func (d *device) Update() *Pressed {
	return d.input.pollEvents(d.os)
}
