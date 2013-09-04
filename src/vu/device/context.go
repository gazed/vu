// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package device

import (
	"log"
)

// context exposes just enough of the native windowing layer to get a window
// wth a graphics context up and running. Context is expected to be used
// indirectly through Device.
type context struct {
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

// newContext creates and returns a structure that interfaces with the
// native layer.
func newContext() *context {
	c := &context{}
	c.nl = nativeLayer() // one of these in each native layer .go.
	c.nr = &nrefs{}
	return c
}

// createDisplay makes and initializes a new native display instance.
// This represents an application so this call is expected to be
// peformed once at startup.
//
// Window attributes are passed into the method since some native layers
// need them right away.
func (c *context) createDisplay(title string, x, y, width, height int) {
	c.nl.setTitle(title)
	c.nl.setSize(x, y, width, height)
	c.nr.display = c.nl.display()
	if c.nr.display == 0 {
		log.Printf("vu/device.context:createDisplay failed.")
	}
}

// dispose releases any resources used by the application. Expected to be
// called after the user closes the application window.
func (c *context) dispose() { c.nl.displayDispose(c.nr) }

// createShell makes and initializes the underlying application shell.
func (c *context) createShell() {
	c.nr.shell = c.nl.shell(c.nr)
	if c.nr.shell == 0 {
		log.Printf("vu/device.context:createShell failed.")
	}
}

// openShell shows the application window to the user.  Event processing
// (see readAndDispatch) must begin before the window will appear.
func (c *context) openShell() { c.nl.shellOpen(c.nr) }

// isAlive returns true as long as the user hasn't exited the application.
func (c *context) isAlive() bool { return c.nl.shellAlive(c.nr) }

// size returns the current size and location of the applications drawing area.
func (c *context) size() (x, y, w, h int) {
	winx, winy, width, height := c.nl.size(c.nr)
	return int(winx), int(winy), int(width), int(height)
}

// showCursor shows or hides the window cursor.
func (c *context) showCursor(show bool) { c.nl.showCursor(c.nr, show) }

// setCursor places the cursor at the given screen coordinates.
func (c *context) setCursorAt(x, y int) { c.nl.setCursorAt(c.nr, x, y) }

// createContext makes and initializes the OpenGL context.
func (c *context) createContext(depth, alpha int) {
	c.nl.setDepthBufferSize(depth)
	c.nl.setAlphaBufferSize(alpha)
	c.nr.context = c.nl.context(c.nr)
	if c.nr.context == 0 {
		log.Printf("vu/device.context:createContext failed.")
	}
}

// swapBuffers flips the drawing buffers.  Expected to be called at the end
// of each drawwing loop.
func (c *context) swapBuffers() { c.nl.swapBuffers(c.nr) }

// readAndDispatch fetches the next user event, if any.
func (c *context) readAndDispatch() *userEvent { return c.nl.readDispatch(c.nr) }

// context
// ===========================================================================
// userEvent

// userEvent is returned by readAndDispatch.  It is the current input of the
// keyboard and mouse.  It tracks any user input changes.
type userEvent struct {
	id     int // Unique event id: one of the const's below
	mouseX int // Current mouse X position.
	mouseY int // Current mouse Y position.
	button int // Currently pressed mouse button (if any).
	key    int // Current key pressed (if any).
	mods   int // Mask of the current modifier keys (if any).
	scroll int // Scroll amount (if any).
}

// The possible event id's are as follows.
const (
	_ = iota // valid event ids start at one.
	closedShell
	resizedShell
	movedShell
	iconifiedShell
	uniconifiedShell
	activatedShell
	deactivatedShell
	clickedMouse
	releasedMouse
	draggedMouse
	movedMouse
	pressedKey
	releasedKey
	scrolled
)

// userEvent
// ===========================================================================
// native

// native specifies the methods that each of the native layers must implement.
// Each native layer is a CGO wrapper over the underlying native window.
//
// Native code is separated in platform specific files as per
//      http://golang.org/pkg/go/build/
// Supported platforms are osx (darwin), win (windows) , and lin (linux).
//      osx: os_darwin.go  wraps: os_darwin.h, os_darwin.m
//      win: os_windows.go wraps: os_windows.h, os_windows.c
//      lin: TODO
// Each will have a unique implementation of the native interface and will
// only be included when building on their respective platforms.
type native interface {

	// display initializes the OS specific graphics layer. The returned
	// value is a reference of the underlying OS structure. For example:
	//    osx: pointer to NSApplication instance.
	//    win: HWND reference from CreateWindowEx.
	//    lin: TODO
	display() int64

	// displayDispose cleans and releases all resources including the OpenGL
	// context.
	displayDispose(r *nrefs)

	// readDispatch fetches user events.  This must be called inside an fast loop
	// so that the native window system event queue is handled in a timely fashion.
	readDispatch(r *nrefs) *userEvent

	// shell creates the "window" on the given display.  In some cases this is
	// a window and in others it holds device independent attributes.  The supplied
	// Shell structure's id is set to a reference of the underlying OS structure.
	// For example:
	//    osx: pointer to NSWindow
	//    win: HDC (handle) to device context from GetDC(hwnd)
	//    lin: TODO
	shell(r *nrefs) int64

	// shellOpen shows the window (shell) on the given display.  This should be
	// called after the OpenGL context has been created.
	shellOpen(r *nrefs)

	// shellAlive is used to check for the user quitting the application.
	// Return true as long as the user hasn't closed the window.
	shellAlive(r *nrefs) bool

	// size gets the current window size and location.  Applications may need this
	// when a window is resized or moved.  The native layer calls are:
	//    osx: NSRect content = [(id)shell contentRectForFrameRect:frame];
	//    win: GetClientRect(hwnd, &rect);
	//    lin: TODO
	size(r *nrefs) (x, y, w, h int)

	// showCursor hides or shows the cursor.  The cursor is locked to the application
	// window when it is hidden.
	//    osx: [NSCursor unhide]; [NSCursor hide];
	//    win: ReleaseCapture(); SetCapture(hwnd);
	//    lin: TODO
	showCursor(r *nrefs, show bool)

	// setCursorAt places the cursor at the given window coordinates.
	//    osx: CGWarpMouseCursorPosition(point);
	//    win: SetCursorPos(loc.x, loc.y);
	//    lin: TODO
	setCursorAt(r *nrefs, x, y int)

	// context creates an OpenGL context and fills in the context field of the nrefs
	// structure. For example the context field is:
	//    osx: pointer to NSOpenGLContext
	//    win: HGLRC (handle) from wglCreateContext(hdc);
	//    lin: TODO
	//
	// Note that display and shell may be updated when creating a context.  This
	// is due to windows need to re-create a window in order to get a properly
	// initialized context.
	context(r *nrefs) int64

	// swapBuffers flips the opengl front and back buffers. All drawing is done
	// in the back buffer. This is expected to be called each pass through the
	// main loop to display the most recent drawing.
	//    osx: [(id)context flushBuffer]
	//    win: SwapBuffers(hdc)
	//    lin: TODO
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

	// setTitle sets the desired window title. This needs to be called before the
	// window is created.
	setTitle(title string)
}
