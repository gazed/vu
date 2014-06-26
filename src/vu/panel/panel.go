// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

// Package panel provides the ability to create a list of 2D graphical
// objects, their properties, their spatial arrangement, and links
// to their associated behaviours. The intent of the package is to make it
// easier to position, and interact with static groups of 2D controls within
// a 3D application.
//
// Rendering, printing, or displaying the panel hierarchy is handled outside
// of this package. This allows an application to pick or create the most
// appropriate device specific means for visualization.
//
// Creating a panel of controls starts with NewPanel and the Panel interface.
// Panel provides the starting point where panel sections and controls may be
// placed.
//
// Package panel is provided as part of the vu (virtual universe) 3D engine.
// It is currently an experimental package. The concern is the complexity of
// the code/interface is not worth the bother of using it in an application.
// See vu/eg/cp.
package panel

// BACKGROUND: be aware of and be inspired from...
//   • Widget/window systems (ie. gtk+) based on vector graphics libraries
//     (ie. cairo. Check out golang draw2d as a cairo option).
//     (also checkout developer.gnome.org/clutter/stable)
//   • Middleware like cegui, libRocket that provide 3D oriented frameworks
//     for rendering 2D geometry to texture surfaces.
//   • Also see a java based version at http://twl.l33tlabs.org
//   • Widget design, ie. http://buoy.sourceforge.net/AboutBuoy.html

// Panel is a top level control panel providing entry points into the widget
// hierarchy. A panel is comprised of controls and control sections.
type Panel interface {
	Section      // Panel is a Section
	Focus() uint // Non-zero is the control id under the mouse.

	// React is expecting to be called with the current user input state.
	// It propgates through all child controls triggering reactions where
	// user input matches the controls triggers.
	React(in *In)

	// Resize indicates that the container window was resized to the given width
	// and height. The change will propogate though all child panels.
	Resize(ww, wh int)
	Size() (ww, wh int) // Fetch the current top level panel size.

	// NewDialog is a Section that can overlay other sections. A new inactive
	// dialog is returned. Dialog sizes are not affected by resizes. Additionally
	// they are positioned in the middle of the current screen.
	NewDialog(w, h int) Section

	// SetDialog activates a dialog. While a dialog is active, and there can only
	// be one active at time, it grabs all the user input until it is dismissed.
	// Passing in a nil Section deactivates the active dialog.
	SetDialog(d Section) // Section that was created with NewDialog.
	Dialog() Section     // Returns the currently active dialog or nil.
	Dialogs() []Section  // Returns all dialogs.
}

// NewPanel creates a top level panel. It expects to be initialized with the
// current size, in pixels, of the application window.
func NewPanel(ww, wh int) Panel { return newPanel(ww, wh) }

// =============================================================================

// newPanel creates a top level control panel.
func newPanel(ww, wh int) *panel {
	p := &panel{section: *newSection(ww, wh)}
	p.dialogs = []Section{}
	p.sw, p.sh = ww, wh
	p.vis = true
	return p
}

// panel implements Panel and is a top level control panel.
type panel struct {
	section           // This is the top most panel.
	focus   uint      // Id of the control the mouse is over, -1 if none.
	dialog  *section  // Currently active dialog. Nil if nothing active.
	dialogs []Section // Dialogs available for activation.
}

// Implements Panel.
func (p *panel) React(in *In) { p.focus = p.react(in) }

// Implements Panel.
func (p *panel) Focus() uint { return p.focus }

// Implements Panel.
func (p *panel) Resize(ww, wh int) {
	p.ww, p.wh = ww, wh
	p.sw, p.sh = ww, wh // Layout the top.
	p.align()           // Propogates window size information.

	// Propogate window size to dialogs.
	for _, dia := range p.dialogs {
		d := dia.(*section)
		d.sx, d.sy = ww/2-d.sw/2, wh/2-d.sh/2
		d.ww, d.wh = ww, wh
		d.align()
	}
}

// Implements Panel.
func (p *panel) Size() (ww, wh int) { return p.ww, p.wh }

// Implements Panel.
func (p *panel) NewDialog(w, h int) Section {
	d := newSection(p.ww, p.wh)         // propogate the window size.
	d.sw, d.sh = w, h                   // fixed dialog size.
	d.sx, d.sy = p.ww/2-h/2, p.wh/2-w/2 // center the dialog.
	p.dialogs = append(p.dialogs, d)
	return d
}

// Implements Panel.
func (p *panel) SetDialog(d Section) {
	if d == nil {
		p.dialog = nil
	} else {
		for _, dia := range p.dialogs {
			if d.Id() == dia.Id() {
				p.dialog = dia.(*section)
				return
			}
		}
	}
}

// Implements Panel.
func (p *panel) Dialog() Section {
	if p.dialog != nil { // because dialog is not an interface.
		return p.dialog
	}
	return nil
}
func (p *panel) Dialogs() []Section { return p.dialogs }

// react overrides section.react so that user input can be directed towards
// an active dialog.
func (p *panel) react(in *In) uint {
	if p.dialog != nil {
		return p.dialog.react(in)
	}
	return p.section.react(in)
}

// =============================================================================

// In is the current user input state. Expected to be updated frequently, it
// holds the current mouse location, recent scrolling amounts, and the down
// duration in update ticks of each key and/or mouse button that is currently
// being pressed by the user.
//
// In is needed by the Panel.React method to find a match between user input
// and triggers on Controls.
type In struct {
	Mx, My int            // Current mouse location.
	Scroll int            // Plus or minus scroll amount. Zero if no scrolling.
	Down   map[string]int // Pressed keys, buttons with down duration.
}

// Set initializes the input state with the given values. Expected to be used
// for transfering values from another input structure. Also expected to be
// read-only, especially in the case of the down map.
func (in *In) Set(mx, my, scroll int, down map[string]int) {
	in.Mx, in.My = mx, my
	in.Scroll = scroll
	in.Down = down
}
