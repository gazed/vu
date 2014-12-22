// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package form

import (
	"github.com/gazed/vu/math/lin"
)

// Section is part of a form. All section dimensions are in pixels.
// The section's top left corner is at x, y. The section label is
// set by the form based on the sections corresponding plan character.
type Section interface {
	Label() string                // Label from plan.
	At() (x, y float64)           // Center location of the section.
	Size() (w, h float64)         // Width and height of the section.
	Bounds() (x, y, w, h float64) // At() and Size() in one call.
	In(x, y int) bool             // True if x and y are in the section.
	Offset() (x, y float64)       // Bottom left corner.
}

// Section interface.
// =============================================================================
// section implementation.

// newSection returns a labelled section with default values.
func newSection(label string) *section { return &section{label: label} }

// section holds values calculated by the layout.
type section struct {
	label    string // user assigned identifier.
	row, col int    // initial section location.

	// The location and size of a section based on constraints.
	x, y float64 // top left corner.
	w, h float64 // width and height.

	// Spans are calculated from the form plan.
	spanx int // span horizontal section. Default 0.
	spany int // span vertical section. Default 0.
}

// Section interface implementation.
func (s *section) At() (x, y float64)           { return s.x, s.y }
func (s *section) Size() (w, h float64)         { return s.w, s.h }
func (s *section) Bounds() (x, y, w, h float64) { return s.x, s.y, s.w, s.h }
func (s *section) Label() string                { return s.label }
func (s *section) In(x, y int) bool {
	hw, hh := lin.Round(s.w*0.5, 0), lin.Round(s.h*0.5, 0)
	fx, fy := float64(x), float64(y)
	return fx >= s.x-hw && fx <= s.x+hw && fy >= s.y-hh && fy <= s.y+hh
}
func (s *section) Offset() (x, y float64) {
	return lin.Round(s.x-s.w*0.5, 0), lin.Round(s.y-s.h*0.5, 0)
}
