// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package panel

// FUTURE: Create more powerful and flexible layouts,
//         See http://www.miglayout.com for design and API inspirations.

import (
	"log"
)

// Layout is anything that can readjust widget locations and sizes.
// A layout operates on one group (section) of widgets at a time.
type Layout interface {

	// Align the sections child positions given the current window width
	// and height in pixels.
	Align(s Section, ww, wh int)
}

// SizeInfo captures rectangular location and dimension information.
type SizeInfo struct {
	X, Y int // Bottom left corner in pixels relative to parent.
	W, H int // Width and height in pixels.
}

// =============================================================================

// ResizeLayout ensures that each widget is placed at an equivalent location
// and size with respect to its original placement. Explicit positioning
// information must be provided for each widget.
type ResizeLayout struct {
	Ww, Wh int
	Info   map[uint]*SizeInfo
}

// Implements Layout
func (rl *ResizeLayout) Align(s Section, ww, wh int) {
	var nx, ny, nw, nh int
	for _, w := range s.Widgets() {
		if info, ok := rl.Info[w.Id()]; ok {
			nx = int((float64(info.X) / float64(rl.Ww)) * float64(ww))
			ny = int((float64(info.Y) / float64(rl.Wh)) * float64(wh))
			nw = int((float64(info.W) / float64(rl.Ww)) * float64(ww))
			nh = int((float64(info.H) / float64(rl.Wh)) * float64(wh))
			w.(Widget).setAt(nx, ny, nw, nh)
		}
	}
}

// =============================================================================

// RatioLayout ensures that each widget retains its original width to height
// ratio while attempting to grow or shrink according to the new screen size.
// Explicit positioning information must be provided for each widget.
type RatioLayout struct {
	Ww, Wh int
	Info   map[uint]*SizeInfo
}

// Implements Layout
func (rl *RatioLayout) Align(s Section, ww, wh int) {
	rx := float64(ww) / float64(rl.Ww) // Find how the window changed, and ...
	ry := float64(wh) / float64(rl.Wh) // ... take the smallest of the two changes.
	ratio := rx
	if ry < rx {
		ratio = ry
	}
	var nx, ny, nw, nh int
	var centerX, centerY float64
	for _, w := range s.Widgets() {
		if info, ok := rl.Info[w.Id()]; ok {
			centerX = float64(info.X) + float64(info.W)*0.5
			centerY = float64(info.Y) + float64(info.H)*0.5
			nw = int(ratio * float64(info.W))
			nh = int(ratio * float64(info.H))
			nx = int(centerX*ratio - float64(nw)*0.5)
			ny = int(centerY*ratio - float64(nh)*0.5)
			w.(Widget).setAt(nx, ny, nw, nh)
		}
	}
}

// =============================================================================

// FixedLayout ensures that each widget is placed at an relative location
// keeping its original dimensions. Explicit positioning information must
// be provided for each widget.
type FixedLayout struct {
	Info map[uint]*FixedInfo
}

// FixedInfo allows each widget to have a unique fixed location.
type FixedInfo struct {
	X, Y int // Widget corner of the widget. Interpreted by Mode.
	W, H int // Width and height in pixels.
	Mode int // One of BL, BR, TL, TR corner for the X, Y offset.
}

// Implements Layout
func (fl *FixedLayout) Align(s Section, ww, wh int) {
	sec := s.(*section)
	for _, w := range s.Widgets() {
		var x, y int
		if info, ok := fl.Info[w.Id()]; ok {
			switch info.Mode {
			case BL:
				x, y = info.X+sec.sx, info.Y+sec.sy
			case BR:
				x, y = ww-info.W-info.X-sec.sx, info.Y+sec.sy
			case TL:
				x, y = info.X+sec.sx, wh-info.H-info.Y-sec.sy
			case TR:
				x, y = ww-info.W-info.X-sec.sx, wh-info.H-info.Y-sec.sy
			case CL:
				x, y = info.X+sec.sx, info.Y+wh/2-info.H/2+sec.sy
			case CR:
				x, y = ww-info.W-info.X-sec.sx, info.Y+wh/2-info.H/2+sec.sy
			case CT:
				x, y = info.X+ww/2-info.W/2+sec.sx, wh-info.H-info.Y-sec.sy
			case CB:
				x, y = info.X+ww/2-info.W/2+sec.sx, info.Y+sec.sy
			}
			w.(Widget).setAt(x, y, info.W, info.H)
		}
	}
}

// Layout references for layouts that interpret controls X, Y values with
// respect to one of the following locations.
const (
	BL = iota // Bottom left.
	BR        // Bottom right.
	TL        // Top left.
	TR        // Top right.
	CL        // Center left.
	CR        // Center right.
	CT        // Center top.
	CB        // Center bottom.
)

// =============================================================================

// GridLayout resizes each child component so that it fits in the parent size
// that has been divided into a grid. GridLayout ensures that the chosen margin
// is identical on all sides of each child widget.
//
// Use GridLayout to create a grouped box of controls.
// GridLayout currently places widgets starting in the bottom left corner and
// places widgets left to right and up according to the number of columns.
type GridLayout struct {
	Margin  int // Grid margin in pixels.
	Columns int // Will be set to 1 if less than 1.
}

//  Implements Layout
func (gl *GridLayout) Align(s Section, ww, wh int) {
	numWidgets := len(s.Widgets())
	if numWidgets < 1 {
		return // No need to layout when there are no child widgets.
	}
	if gl.Columns < 1 {
		log.Printf("panel.GridLayout invalid number of columns")
		gl.Columns = 1
	}

	// Get the grid size based on the number of widgets and columns.
	cols := gl.Columns
	rows := ((numWidgets - 1) / cols) + 1

	// Calculate the relative parent size to current window size.
	sec := s.(*section)
	rx := float64(sec.sx) / float64(ww)
	ry := float64(sec.sy) / float64(wh)
	rw := float64(sec.sw) / float64(ww)
	rh := float64(sec.sh) / float64(wh)

	// Assign each widget its portion of the parents space.
	for cnt, w := range s.Widgets() {
		row, col := cnt/cols, cnt%cols
		cf0, rf0 := float64(row)/float64(rows), float64(col)/float64(cols)
		cf1, rf1 := float64(cols), float64(rows)
		wx := int((rx+rf0*rw)*float64(ww)) + gl.Margin
		wy := int((ry+cf0*rh)*float64(wh)) + gl.Margin
		ww := int((rw/cf1)*float64(ww)) - gl.Margin
		wh := int((rh/rf1)*float64(wh)) - gl.Margin
		w.(Widget).setAt(wx, wy, ww, wh)
	}
}
