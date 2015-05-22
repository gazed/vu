// Copyright © 2014-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// Package form is used to divide a 2D area into sections using plans.
// This reduces some of the code complexity involved in creating form based
// 2D UIs. Forms are split into sections using arrays of strings where each
// string character represents a section. For example a form with 4 sections
// can be created as follows:
//    // Create a form with plan "ab"
//    //                         "cd"
//    f := NewForm([]string{"ab", "cd"}, formWidth, formHeight)
//    for cnt, sect := range f.Sections() {
//         x, y, w, h := sect.Bounds()
//         // now use section positions...
//    }
// Each section can be queried for its center location and size and is
// labelled with the corresponding string character from the plan.
//
// Package form is provided as part of the vu (virtual universe) 3D engine.
package form

import (
	"fmt"
	"log"
	"strings"
)

// Form organizes a 2D area into sections. Once a form is created, its
// sections can be queried for their center pixel locations and sizes.
// Resizing a form updates the section centers and sizes. Form is intended
// for static layouts where the form parts are added once and never removed.
//
// Form generates layouts based on the given plan.
type Form interface {
	Resize(w, h int)          // Realign sections to the new dimensions.
	Section(l string) Section // Get Section for plan label l.
	Sections() []Section      // All sections in this form.
}

// Note: This is an experimental package. It is the third attempt at producing
//       a helper API for UI's. Like the first two attempts it looks like the
//       effort and complexity is not worth the benefit gained by the App layer.
// Design:
//   • Focus purely on positioning, leaving other UI parts like buttons
//     windows, and dialogs for some other package.
//   • Key Ideas from: http://www.migcalendar.com/miglayout/mavensite/docs/whitepaper.html
//     • Having a single intelligent layout rather than different types of layout.
//     • Using text based constraint-overrides for clarity.
//
// Background: be aware of and be inspired from...
//   • Widget/window systems (ie. gtk+) based on vector graphics libraries
//     (ie. cairo. Check out golang draw2d as a cairo option).
//     (also checkout developer.gnome.org/clutter/stable)
//   • Middleware like cegui, libRocket that provide 3D oriented frameworks
//     for rendering 2D geometry to texture surfaces.
//   • Also see a java based version at http://twl.l33tlabs.org
//   • Widget design, ie. http://buoy.sourceforge.net/AboutBuoy.html
//   • id tech gui's: http://modwiki.xnet.fi/GUI_scripting
//
// FUTURE work:
//   • Look at constraints as a system of linear equations that are solved
//     to get a layout - mainly Lutteroth's work.
//     http://arxiv.org/pdf/1401.1031.pdf
//     https://www.cs.auckland.ac.nz/~lutteroth/publications/LutterothWeber2008-ModularGUILayout.pdf
//     http://crpit.com/confpapers/CRPITV50Lutteroth.pdf
//   • See if this helps.
//     http://www.w3.org/TR/2012/WD-css3-grid-layout-20121106/

// Form interface.
// =============================================================================
// form implements Form.

// form tracks the individual sections that share a given 2D area.
type form struct {

	// Injected on creation.
	plan []string // grid template.
	rows int      // number of template rows.
	cols int      // number of template columns.
	w, h float64  // absolute width and height.

	// Sections, one per unique plan character.
	sects map[string]*section // generated from plan.
	byrow []Section           // ordered list of sections.

	// Information needed for aligning the grid.
	basew float64      // original width for grabx.
	baseh float64      // original height for graby.
	srows []Section    // ordered list of sections for align.
	scols []Section    // ordered list of sections for align.
	grabx map[int]bool // cols that grab extra space.
	graby map[int]bool // rows that grab extra space.

	// Guides: gaps in pixels between columns and rows.
	gapv int // row spacing between sections.
	gaph int // column spacing between sections.

	// Guides: border padding in pixels at form edges.
	padt int // top border. Default 0.
	padb int // bottom border. Default 0.
	padl int // left border. Default 0.
	padr int // right border. Default 0.
}

// New creates a form based on a character string plan. The plan visually
// describes how the form should be partitioned. Some example plans are:
//    "ab" or "xxaby"
//    "cd"    "xxcdy"
// Identical adjacent characters create a single section spanning multiple
// rows and/or columns.
//
// Form guidelines may be set at creation using separate strings as follows:
//   "gap x y"     Gap in pixels x between columns and y between rows.
//   "pad t b l r" Border padding in pixels: top, bottom, left, right.
//   "grabx c"     Have column number c use extra horizontal space.
//   "graby r"     Have row number r use extra vertical space.
// Note that sizes are in pixels. Column and row numbering starts at 0.
func New(plan []string, w, h int, constraints ...string) Form {
	rows, cols, err := validatePlan(plan)
	if err != nil {
		log.Printf("%s", err)
		return nil
	}
	f := &form{plan: plan, rows: rows, cols: cols}
	f.w = float64(w)
	f.h = float64(h)
	f.sects = map[string]*section{}
	f.grabx = map[int]bool{}
	f.graby = map[int]bool{}

	// Build sections from the plan. Each section has its own letter label.
	// Depends on a nice box shape for spanning sections.
	for row, rowstr := range plan {
		labels := strings.Split(rowstr, "")
		for col, label := range labels {
			if c, ok := f.sects[label]; !ok {
				c := &section{label: label, col: col, row: row}
				f.sects[label] = c
				f.byrow = append(f.byrow, c)
			} else {

				// check for increase of section spanx or spany.
				if col > c.col+c.spanx {
					c.spanx += 1
				}
				if row > c.row+c.spany {
					c.spany += 1
				}
			}
		}
	}

	// create the ordered list by columns for vertical align processing.
	// This inserts spanning sections multiple times to get proper
	// alignment later on.
	for col := 0; col < f.cols; col++ {
		uniques := ""
		for row := 0; row < f.rows; row++ {
			label := f.planLabel(row, col)
			if !strings.Contains(uniques, label) {
				uniques += label
				f.scols = append(f.scols, f.sects[label])
			}
		}
	}
	// Ditto: ordered list by rows for horizontal align processing.
	for row := 0; row < f.rows; row++ {
		uniques := ""
		for col := 0; col < f.cols; col++ {
			label := f.planLabel(row, col)
			if !strings.Contains(uniques, label) {
				uniques += label
				f.srows = append(f.srows, f.sects[label])
			}
		}
	}

	// layout the sections based on the constraints.
	f.parse(constraints)
	f.basew = f.w // remember width for allocating grab space.
	f.baseh = f.h // remember height for allocating grab space.
	f.align()
	return f
}

// planLabel gets the plan label corresponding to the given
// row and column. It is sufficient, not efficient.
func (f *form) planLabel(row, col int) string {
	rowstr := f.plan[row]
	labels := strings.Split(rowstr, "")
	return labels[col]
}

// Resize the form and realign all sections.
func (f *form) Resize(w, h int) {
	f.w, f.h = float64(w), float64(h)
	f.align()
}

// Sections returns all the sections that have been added to the form.
func (f *form) Sections() []Section      { return f.byrow }
func (f *form) Section(l string) Section { return f.sects[l] }

// align readjusts all sections according to specified guide lines and the
// overall form size. Work in floats so rounding errors don't creeep into
// the sizing.
//
// FUTURE: Replace with a linear equation solver. Currently this algorithm
//         handles simple cases and messes up complex plans.
func (f *form) align() {

	// Set all section widths and x center locations
	padl, padr, gaph := float64(f.padl), float64(f.padr), float64(f.gaph)
	sectw := (f.w - padl - padr) / float64(f.cols)
	perGrabx := 0.0
	if f.w > f.basew && len(f.grabx) > 0 {
		sectw = (f.basew - padl - padr) / float64(f.cols)  // reset base width.
		perGrabx = (f.w - f.basew) / float64(len(f.grabx)) // amount per grab col.
	}
	usedx, lastx, nextx := 0.0, 0.0, 0.0
	for _, s := range f.srows {
		sect := s.(*section)
		if sect.col == 0 { // reset total grabbed for each row.
			usedx, lastx, nextx = padl, padl, padl
		}
		lastx = nextx
		spanx := float64(sect.spanx)
		sect.w = sectw*(spanx+1) - gaph // every column gets the base amount.
		if _, ok := f.grabx[sect.col]; ok {
			sect.w += perGrabx // grabbers get base width + part of the extra.
			nextx += perGrabx
		}
		nextx += sectw * (spanx + 1)
		sect.x = usedx + (nextx-lastx)*0.5
		usedx += (nextx - lastx)
	}

	// Set all section heights and y center locations
	padt, padb, gapv := float64(f.padt), float64(f.padb), float64(f.gapv)
	secth := (f.h - padt - padb) / float64(f.rows)
	perGraby := 0.0
	if f.h > f.baseh && len(f.graby) > 0 {
		secth = (f.baseh - padt - padb) / float64(f.rows)  // reset base height.
		perGraby = (f.h - f.baseh) / float64(len(f.graby)) // amount per grap row.
	}
	usedy, lasty, nexty := 0.0, 0.0, 0.0
	for _, s := range f.scols {
		sect := s.(*section)
		if sect.row == 0 {
			usedy, lasty, nexty = padt, padt, padt
		}
		lasty = nexty
		spany := float64(sect.spany)
		sect.h = secth*(spany+1) - gapv // every row get the base amount.
		if _, ok := f.graby[sect.row]; ok {
			sect.h += perGraby // grabbers get base width + part of the extra.
			nexty += perGraby
		}
		nexty += secth * (spany + 1)
		sect.y = f.h - (usedy + (nexty-lasty)*0.5)
		usedy += (nexty - lasty)
	}
}

// parse uses the given string to override the sections default
// behaviour. Note that Sscanf changes values that match the pattern,
// ie. partial successes are possible.
func (f *form) parse(guides []string) {
	var grabx, graby int
	for _, c := range guides {
		grabx, graby = -1, -1
		fmt.Sscanf(c, "grabx %d", &grabx)
		fmt.Sscanf(c, "graby %d", &graby)
		if grabx >= 0 && grabx < f.cols {
			f.grabx[grabx] = true
		}
		if graby >= 0 && graby < f.rows {
			f.graby[graby] = true
		}
		fmt.Sscanf(c, "gap %d %d", &f.gapv, &f.gaph)
		fmt.Sscanf(c, "pad %d %d %d %d", &f.padt, &f.padb, &f.padl, &f.padr)
	}
}

// form
// =============================================================================
// Utility functions.

// Ensure that the string based plan makes sense.
func validatePlan(plan []string) (rows, cols int, err error) {
	rows = len(plan)
	cols = len(plan[0])
	if cols == 0 {
		return 0, 0, fmt.Errorf("Invalid grid plan")
	}
	for _, row := range plan {
		if cols != len(row) {
			return 0, 0, fmt.Errorf("Invalid grid plan")
		}
	}
	return rows, cols, nil
}
