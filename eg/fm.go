// Copyright © 2014-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"log"

	"github.com/gazed/vu"
	"github.com/gazed/vu/form"
	"github.com/gazed/vu/math/lin"
)

// fm demos and tests the vu/form package by visualizing grid layouts
// created from string based layout plans.
//
// This is the 3rd iteration in an experiment to see if there is any
// application benefit to a 2D UI support package. So far this has been
// NO in that the benefit to the app does not justify the amount of work
// needed in the 2D UI support package.
//
// See vu/form/form.go for more design notes.
func fm() {
	fm := &fmtag{}
	if err := vu.New(fm, "Form Layout", 400, 100, 800, 600); err != nil {
		log.Printf("fm: error starting engine %s", err)
	}
	defer catchErrors()
}

// Encapsulate example specific data with a unique "tag".
type fmtag struct {
	top     vu.Pov    //
	cam     vu.Camera // visible layouts.
	ww, wh  int       // window width and height.
	example int       // current layout example.g
	layouts []*layout // demonstrate multiple layouts.
}

// Create is the engine callback for initial asset creation.
func (fm *fmtag) Create(eng vu.Eng, s *vu.State) {
	fm.cam = eng.Root().NewCam()
	fm.cam.SetUI()
	eng.SetColor(0.95, 0.95, 0.95, 1)

	// create the panel layout examples.
	fm.layouts = append(fm.layouts, fm.simpleLayout(eng, s.W, s.H))
	fm.layouts = append(fm.layouts, fm.spanLayout(eng, s.W, s.H))
	fm.layouts = append(fm.layouts, fm.grabLayout(eng))
	fm.layouts = append(fm.layouts, fm.largeLayout(eng, s.W, s.H))
	fm.layouts = append(fm.layouts, fm.doubleLayout(eng))
	fm.layouts[fm.example].setVisible(true)

	// set non default engine state.
	fm.resize(s.W, s.H)
}

// Update is the regular engine callback.
func (fm *fmtag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	if in.Resized {
		fm.resize(s.W, s.H)
	}
	for press, down := range in.Down {
		switch {

		// switch to the next layout example.
		case press == vu.K_Tab && down == 1:
			fm.layouts[fm.example].setVisible(false)
			fm.example = fm.example + 1
			if fm.example >= len(fm.layouts) {
				fm.example = 0
			}
			fm.layouts[fm.example].setVisible(true)
		}
	}
}

// resize handles user screen/window changes.
func (fm *fmtag) resize(ww, wh int) {
	fm.cam.SetOrthographic(0, float64(ww), 0, float64(wh), 0, 10)
	for _, lo := range fm.layouts {
		lo.resize(ww, wh)
	}
}

// simpleLayout creates a 2x2 form with 10 pixel gaps.
func (fm *fmtag) simpleLayout(eng vu.Eng, ww, wh int) *layout {
	lo := &layout{}
	plan := []string{
		"ab",
		"cd",
	}
	lo.form = form.New(plan, ww, wh, "gap 5 5", "pad 5 5 5 5")
	lo.visualize(eng)
	return lo
}

// spanLayout creates a form where the cell with the same label will
// span rows and columns.
func (fm *fmtag) spanLayout(eng vu.Eng, ww, wh int) *layout {
	lo := &layout{}
	plan := []string{
		"axxb",
		"cxxb",
		"cxxd",
		"eeef",
	}
	lo.form = form.New(plan, ww, wh, "gap 5 5", "pad 5 5 5 5")
	lo.visualize(eng)
	return lo
}

// grabLayout creates a form where the base size dictates the max size
// that non-grabby rows and columns will grow to.
func (fm *fmtag) grabLayout(eng vu.Eng) *layout {
	lo := &layout{}
	plan := []string{
		"abc",
		"def",
		"ghi",
	}
	lo.form = form.New(plan, 200, 200, "grabx 1", "graby 1", "gap 5 5", "pad 5 5 5 5")
	lo.visualize(eng)
	return lo
}

// largeLayout creates a form with multiple spanning sections and
// a single spanning section.
func (fm *fmtag) largeLayout(eng vu.Eng, ww, wh int) *layout {
	lo := &layout{}
	plan := []string{
		"aabbbccd",
		"exxfyyyg",
		"exxhyyyg",
		"ixxhyyyg",
		"jklhmmnn",
	}
	lo.form = form.New(plan, ww, wh, "gap 5 5", "grabx 0", "graby 0", "pad 5 5 5 5")
	lo.visualize(eng)
	return lo
}

// doubleLayout creates a form within a form to create a more complex layout.
func (fm *fmtag) doubleLayout(eng vu.Eng) *layout {
	lo := &layout{}
	plan := []string{
		"abc",
		"def",
		"ghi",
	}
	lo.form = form.New(plan, 200, 200, "grabx 1", "graby 1", "gap 5 5", "pad 5 5 5 5")
	lo.visualize(eng)
	lo.lo = fm.interiorLayout(eng, lo.form.Section("e"))
	return lo
}

// interior layout is part of doubleLayout.
// It creates a second form inside the middle section of the first form.
func (fm *fmtag) interiorLayout(eng vu.Eng, s form.Section) *layout {
	lo := &layout{}
	w, h := s.Size()
	iw, ih := int(lin.Round(w, 0)), int(lin.Round(h, 0))
	plan := []string{
		"pqr",
		"pqr",
		"stu",
	}
	lo.form = form.New(plan, iw, ih, "gap 5 5")
	lo.visualize(eng)
	return lo
}

// =============================================================================

// layout
type layout struct {
	form   form.Form // cell and label position information.
	lo     *layout   // for doubleLayout.
	top    vu.Pov    // single spot for making visible.
	sects  []vu.Pov  // visual representation of a form cell.
	labels []vu.Pov  // cell label.
}

// Called once to create the visual parts of a panel.
func (lo *layout) visualize(eng vu.Eng) {
	lo.top = eng.Root().NewPov()
	lo.setVisible(false)
	lo.sects = make([]vu.Pov, len(lo.form.Sections()))
	lo.labels = make([]vu.Pov, len(lo.form.Sections()))
	for cnt, sect := range lo.form.Sections() {

		// place a box at the section location.
		lo.sects[cnt] = lo.top.NewPov()
		lo.sects[cnt].NewModel("uv").LoadMesh("icon").AddTex("cell")

		// place the cell name in the middle of the cell.
		lo.labels[cnt] = lo.top.NewPov()
		model := lo.labels[cnt].NewModel("uv").AddTex("lucidiaSu16Black")
		if sect.Label() == "" {
			model.LoadFont("lucidiaSu16").SetPhrase("-")
		} else {
			model.LoadFont("lucidiaSu16").SetPhrase(sect.Label())
		}
	}
}

// setVisible hides or shows the form and any child forms.
func (lo *layout) setVisible(vis bool) {
	lo.top.SetVisible(vis)
	if lo.lo != nil {
		lo.lo.setVisible(vis)
	}
}

// resize informs the panel of the size change and updates the visual components
// to match the new cell sizes.
func (lo *layout) resize(ww, wh int) {
	lo.form.Resize(ww, wh)
	for cnt, sect := range lo.form.Sections() {
		x, y, w, h := sect.Bounds()
		lo.sects[cnt].SetScale(w, h, 0)
		lo.sects[cnt].SetLocation(x, y, 0)
		lo.labels[cnt].SetLocation(x, y, 0)
	}
	if lo.lo != nil {
		w, h := lo.form.Section("e").Size()
		iw, ih := int(lin.Round(w, 0)), int(lin.Round(h, 0))
		offx, offy := lo.form.Section("e").Offset()
		lo.lo.resizeChild(offx, offy, iw, ih)
	}
}

// resizeChild resizes the form and offsets it to align with its parent section.
func (lo *layout) resizeChild(offx, offy float64, ww, wh int) {
	lo.form.Resize(ww, wh)
	for cnt, sect := range lo.form.Sections() {
		x, y, w, h := sect.Bounds()
		lo.sects[cnt].SetScale(w, h, 0)
		lo.sects[cnt].SetLocation(x+offx, y+offy, 0)
		lo.labels[cnt].SetLocation(x+offx, y+offy, 0)
	}
}
