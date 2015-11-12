// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package device

import (
	"log"
)

// native specifies the methods that each of the native layers must implement.
// Each native layer is a CGO wrapper over the underlying native window.
//
// Native code is separated in platform specific files as per
//      http://golang.org/pkg/go/build/
// Supported platforms are osx (darwin), win (windows) , and lin (linux).
//      osx: os_darwin.go  wraps: os_darwin.h, os_darwin.m
//      win: os_windows.go wraps: os_windows.h, os_windows.c
//      FUTURE lin: os_linux.go   wraps: os_linux.h, os_linux.c
// Each will have a unique implementation of the native interface and will
// only be included when building on their respective platforms.
//
// Note that OSX expects all window access to be done on the main thread.
type native interface {

	// display initializes the OS specific graphics layer. The returned
	// value is a reference of the underlying OS structure. For example:
	//    osx: pointer to NSApplication instance.
	//    win: HWND reference from CreateWindowEx.
	//    lin: FUTURE
	display() int64

	// displayDispose cleans and releases all resources including the OpenGL
	// context.
	displayDispose(r *nrefs)

	// readDispatch fetches user events. This must be called inside a fast loop
	// so that the native window system event queue is handled in a timely fashion.
	// The input userInput is filled and returned. This will return with the
	// latest modifier keys and mouse location even if there are no user events
	// to process.
	readDispatch(r *nrefs, in *userInput) *userInput

	// shell creates the "window" on the given display. In some cases this is
	// a window and in others it holds device independent attributes. The supplied
	// Shell structure's id is set to a reference of the underlying OS structure.
	// For example:
	//    osx: pointer to NSWindow
	//    win: HDC (handle) to device context from GetDC(hwnd)
	//    lin: FUTURE
	shell(r *nrefs) int64

	// shellOpen shows the window (shell) on the given display.  This should be
	// called after the OpenGL context has been created.
	shellOpen(r *nrefs)

	// shellAlive is used to check for the user quitting the application.
	// Return true as long as the user hasn't closed the window.
	shellAlive(r *nrefs) bool

	// size gets the current window size and location of the bottom left corner.
	// Applications may need this when a window is resized or moved.
	// The native layer calls are:
	//    osx: NSRect content = [(id)shell contentRectForFrameRect:frame];
	//    win: GetClientRect(hwnd, &rect);
	//    lin: FUTURE
	size(r *nrefs) (x, y, w, h int)

	// showCursor hides or shows the cursor. The cursor is locked to the
	// application window when it is hidden.
	//    osx: [NSCursor unhide]; [NSCursor hide];
	//    win: ReleaseCapture(); SetCapture(hwnd);
	//    lin: FUTURE
	showCursor(r *nrefs, show bool)

	// setCursorAt places the cursor at the given window coordinates.
	//    osx: CGWarpMouseCursorPosition(point);
	//    win: SetCursorPos(loc.x, loc.y);
	//    lin: FUTURE
	setCursorAt(r *nrefs, x, y int)

	// context creates an OpenGL context and fills in the context field of the
	// nrefs structure. For example the context field is:
	//    osx: pointer to NSOpenGLContext
	//    win: HGLRC (handle) from wglCreateContext(hdc);
	//    lin: FUTURE
	//
	// Note that display and shell may be updated when creating a context. This
	// is due to windows need to re-create a window in order to get a properly
	// initialized context.
	context(r *nrefs) int64

	// swapBuffers flips the opengl front and back buffers. All drawing is done
	// in the back buffer. This is expected to be called each pass through the
	// main loop to display the most recent drawing.
	//    osx: [(id)context flushBuffer]
	//    win: SwapBuffers(hdc)
	//    lin: FUTURE
	swapBuffers(r *nrefs)

	// setAlphaBufferSize sets the desired size of the OpenGL alpha buffer.
	// This needs to be called before the OpenGL context is created.
	setAlphaBufferSize(size int)

	// setAlphaBufferSize sets the desired size of the OpenGL depth buffer.
	// This needs to be called before the OpenGL context is created.
	setDepthBufferSize(size int)

	// setSize sets the desired window size. This needs to be called before the
	// window is created.
	setSize(x, y, width, height int)

	// setTitle sets the desired window title. This needs to be called before
	// the window is created.
	setTitle(title string)

	// isFullscreen returns true when application is is full screen mode.
	isFullscreen(r *nrefs) bool

	// toggleFullscreen flips between application windowed and full screen mode.
	// Must be called after starting processing with readDispatch().
	toggleFullscreen(r *nrefs)
}

// native
// ===========================================================================
// nativeOs wraps a native implementation.

// nativeOs exposes just enough of the native windowing layer to get a window
// wth a graphics context up and running. Native is expected to be used
// indirectly through Device.
type nativeOs struct {
	nl native // native layer support
	nr *nrefs // references to native layer objects.
}

// nrefs keeps and passes pointers/handles to the native layer window,
// shell, and drawing context objects. The different native layers need one
// or more of the references depending on the call.
type nrefs struct {
	display int64 // native os display
	shell   int64 // native os shell
	context int64 // native os opengl context
}

// newNative creates and returns a structure that interfaces with the
// native layer.
func newNativeOs() *nativeOs {
	os := &nativeOs{}
	os.nl = nativeLayer() // one of these in each native layer os_*.go.
	os.nr = &nrefs{}
	return os
}

// createDisplay makes and initializes a new native display instance.
// This represents an application so this call is expected to be
// peformed once at startup.
//
// Window attributes are passed into the method since some native layers
// need them right away.
func (os *nativeOs) createDisplay(title string, x, y, width, height int) {
	os.nl.setTitle(title)
	os.nl.setSize(x, y, width, height)
	os.nr.display = os.nl.display()
	if os.nr.display == 0 {
		log.Printf("vu/device.native:createDisplay failed.")
	}
}

// dispose releases any resources used by the application. Expected to be
// called after the user closes the application window.
func (os *nativeOs) dispose() { os.nl.displayDispose(os.nr) }

// createShell makes and initializes the underlying application shell.
func (os *nativeOs) createShell() {
	os.nr.shell = os.nl.shell(os.nr)
	if os.nr.shell == 0 {
		log.Printf("vu/device.native:createShell failed.")
	}
}

// openShell shows the application window to the user. Event processing
// (see readAndDispatch) must begin before the window will appear.
func (os *nativeOs) openShell() { os.nl.shellOpen(os.nr) }

// isAlive returns true as long as the user hasn't exited the application.
func (os *nativeOs) isAlive() bool { return os.nl.shellAlive(os.nr) }

// size returns the current size and location of the applications drawing area.
func (os *nativeOs) size() (x, y, w, h int) {
	winx, winy, width, height := os.nl.size(os.nr)
	return int(winx), int(winy), int(width), int(height)
}

// showCursor shows or hides the window cursor.
func (os *nativeOs) showCursor(show bool) { os.nl.showCursor(os.nr, show) }

// setCursor places the cursor at the given screen coordinates.
func (os *nativeOs) setCursorAt(x, y int) { os.nl.setCursorAt(os.nr, x, y) }

// createContext makes and initializes the OpenGL context.
func (os *nativeOs) createContext(depth, alpha int) {
	os.nl.setDepthBufferSize(depth)
	os.nl.setAlphaBufferSize(alpha)
	os.nr.context = os.nl.context(os.nr)
	if os.nr.context == 0 {
		log.Printf("vu/device.native:createContext failed.")
	}
}

// isFullscreen returns true when application is is full screen mode.
func (os *nativeOs) isFullscreen() bool { return os.nl.isFullscreen(os.nr) }

// toggleFullscreen flips between application windowed and full screen mode.
func (os *nativeOs) toggleFullscreen() { os.nl.toggleFullscreen(os.nr) }

// swapBuffers flips the drawing buffers.
// Expected to be called at the end of each drawwing loop.
func (os *nativeOs) swapBuffers() { os.nl.swapBuffers(os.nr) }

// readDispatch polls the next user event from the native OS.
// Reading the event must be done on the main thread, but the
// events can be processed in a separate thread.
func (os *nativeOs) readDispatch(in *userInput) *userInput {
	return os.nl.readDispatch(os.nr, in)
}
