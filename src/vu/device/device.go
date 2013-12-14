// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

// Package device provides platform/os access to a 3D rendering
// context and user input. It provides only what is necessary
// to get a OS specific window with a 3D rendering context.
// The application layer is expected to provide any necessary
// buttons, controls and sub-panels.
//
// Package device provides access to user keyboard and mouse input through
// the Update method and Pressed structure.
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

// Pressed is used to communicate current user input.
// This is the list of keys that are currently being pressed and
// how long they have been pressed (measured in update ticks).
//
// A postitive duration means the key is still being held down.
// A negative duration means that the key has been released since
// the last poll. The duration prior to release can be determined by
// its difference with KEY_RELEASE.
type Pressed struct {
	Mx, My  int            // Current mouse location.
	Down    map[string]int // Pressed keys and pressed duration.
	Shift   bool           // True if the shift modifier is pressed.
	Control bool           // True if the control modifier is pressed.
	Focus   bool           // True if window has focus.
	Resized bool           // True if window was resized or moved.
}

// Device interfaces
// ===========================================================================
// device provides default Device implementation.

// New provides a newly initialized Device with an underlying window and
// graphics context created, but not yet displayed. The only thing left to do
// is to open the device and start polling it for user input.
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
//     os_linux  : Linux native layer. FUTURE: Wraps the following.
//        os_linux.cpp : c++ code wrapping linux API.
//        os_linux.h
//        os_linux_test.c
//
// Design note 2: user events need to be processed on the main thread for OSX.
//                See: native::readAndDispatch

// device provides a simplification layer over the more raw context layer.
type device struct {
	os      *nativeOs // Native layer wrapper.
	pressed *input    // User input handler.
}

// newDevice initializes a OS specific window with a valid OpenGL context.
func newDevice(title string, x, y, width, height int) *device {
	d := &device{}
	d.os = newNativeOs()
	d.os.createDisplay(title, x, y, width, height)
	d.os.createShell()
	depthBufferBits, alphaBits := 24, 8 // resonable defaults
	d.os.createContext(depthBufferBits, alphaBits)
	d.pressed = newInput(d.os)
	return d
}

// Access the device specific information in a consistent and general manner.
func (d *device) Open()                           { d.os.openShell() }
func (d *device) Dispose()                        { d.os.dispose() }
func (d *device) IsAlive() bool                   { return d.os.isAlive() }
func (d *device) Size() (x, y, width, height int) { return d.os.size() }
func (d *device) ShowCursor(show bool)            { d.os.showCursor(show) }
func (d *device) SwapBuffers()                    { d.os.swapBuffers() }
func (d *device) SetCursorAt(x, y int)            { d.os.setCursorAt(x, y) }
func (d *device) Update() *Pressed {
	d.os.readAndDispatch(d.pressed.events)
	return d.pressed.latest()
}
