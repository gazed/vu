// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package panel

// Widget is a rectangular area of the display. It can have a background
// image and can be hidden. Widget locations are set by the parent and the
// current absolute window location can be queried.
//
// Widget is intended as the base class for anything that can cover a 2D
// area of the screen.
type Widget interface {
	Id() uint           // Unique runtime identifier.
	Over(x, y int) bool // Is the screen x, y within the widget.

	// At returns the center location and half-extents in absolute
	// screen/window pixels. It is a convenience method for rendering.
	At() (cx, cy, hx, hy float64)

	// Allow widgets to appear or disappear from the UI.
	Visible() bool      // Visible widgets are rendered so,
	SetVisible(is bool) // ...invisible widgets can be placeholders.

	// Optional background image for the widget. The image is sized to
	// exactly fill the widget. The img is a resource identifier matching
	// the renderers resource identifier conventions.
	Img() (img string)        // Img identifies the widgets...
	SetImg(img string) Widget // ...background image.

	// Internal interfaces.
	react(in *In) uint        // process user input and return the focus.
	setAt(sx, sy, sw, sh int) // relocate/resize the frame.
	at() (sx, sy, sw, sh int) // get the frame size.
}

// FUTURE: AddAnimation, RemAnimation for Widget.

// =============================================================================

// widget implements Widget and the common methods of the framework interface.
// As a base class it provides the simple common support methods needed by the
// super classes.
type widget struct {
	id  uint   // Unique identifier.
	vis bool   // Render/interact with only if visible.
	img string // Optional background image.

	// Absolute screen locations calculated on resizes.
	sx, sy int // Absolute screen bottom left corner in pixels.
	sw, sh int // Absolute screen width and height in pixels.
}

// widgetId is a global widget counter so that each widget can have a unique
// identifier in a given application instance. Not intended to be used across
// application restarts. This value only grows so the max screen widgets for
// an application is 2^32-1. The expection for widgets per application is in
// the 100's to 1000's range.
var widgetId uint

// newWidget is a package internal method for creating widgets. This is the
// only method that should be incrementing widgetId.
func newWidget() widget {
	widgetId++                  // first widget starts at 1...
	return widget{id: widgetId} // ... so that 0 means "no widget".
}

// Widget implementation.
func (w *widget) Id() uint { return w.id }

// Widget implementation.
func (w *widget) Over(x, y int) bool {
	return x > w.sx && x < w.sx+w.sw && y > w.sy && y < w.sy+w.sh
}

// Widget implementation.
func (w *widget) Visible() bool { return w.vis }

// Widget implementation.
func (w *widget) SetVisible(is bool) { w.vis = is }

// Widget implementation.
func (w *widget) Img() (img string) { return w.img }

// Widget implementation.
func (w *widget) SetImg(img string) Widget {
	w.img = img
	return w
}

// Widget implementation.
func (w *widget) At() (cx, cy, hx, hy float64) {
	cx = float64(w.sx) + float64(w.sw)*0.5        // center x.
	cy = float64(w.sy) + float64(w.sh)*0.5        // center y.
	hx, hy = float64(w.sw)*0.5, float64(w.sh)*0.5 // half-extents.
	return
}

// Widget implementation. Expecting window coordinates in pixels.
func (w *widget) setAt(sx, sy, sw, sh int) { w.sx, w.sy, w.sw, w.sh = sx, sy, sw, sh }
func (w *widget) at() (sx, sy, sw, sh int) { return w.sx, w.sy, w.sw, w.sh }

// Widget implementation. Override: default does nothing.
func (w *widget) react(in *In) uint { return 0 }
