// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"log"

	"github.com/gazed/vu"
	"github.com/gazed/vu/form"
	"github.com/gazed/vu/math/lin"
)

// fm demos and tests the vu/form package by visualizing grid layouts
// created from string based layout plans. Overall this is an experiment
// in trying to provide some application GUI support.
func fm() {
	fm := &fmtag{}
	ww, wh := 800, 600
	var err error
	if fm.eng, err = vu.New("Control Panel", 1200, 100, ww, wh); err != nil {
		log.Printf("fm: error intitializing engine %s", err)
		return
	}
	fm.eng.SetDirector(fm)  // get user input through Director.Update()
	fm.create()             // create initial assests.
	defer fm.eng.Shutdown() // shut down the engine.
	defer catchErrors()
	fm.eng.Action()
}

// Encapsulate example specific data with a unique "tag".
type fmtag struct {
	eng     vu.Engine // 3D engine.
	scene   vu.Scene  // visible layouts.
	ww, wh  int       // window width and height.
	example int       // current layout example.g
	layouts []*layout // demonstrate multiple layouts.
}

// create is the startup asset creation.
func (fm *fmtag) create() {
	fm.scene = fm.eng.AddScene(vu.VO)
	fm.scene.Set2D()
	_, _, fm.ww, fm.wh = fm.eng.Size()

	// create the panel layout examples.
	fm.layouts = append(fm.layouts, fm.simpleLayout())
	fm.layouts = append(fm.layouts, fm.spanLayout())
	fm.layouts = append(fm.layouts, fm.grabLayout())
	fm.layouts = append(fm.layouts, fm.largeLayout())
	fm.layouts = append(fm.layouts, fm.doubleLayout())
	fm.layouts[fm.example].setVisible(true)
	fm.eng.Color(0.95, 0.95, 0.95, 1)
	fm.resize()
}

// Update is the regular engine callback.
func (fm *fmtag) Update(in *vu.Input) {
	if in.Resized {
		fm.resize()
	}
	for press, down := range in.Down {
		switch {
		case press == "Tab" && down == 1:

			// switch to the next layout example.
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
func (fm *fmtag) resize() {
	x, y, ww, wh := fm.eng.Size()
	fm.eng.Resize(x, y, ww, wh)
	fm.ww, fm.wh = ww, wh
	fm.scene.Cam().SetOrthographic(0, float64(ww), 0, float64(wh), 0, 10)
	for _, lo := range fm.layouts {
		lo.resize(ww, wh)
	}
}

// simpleLayout creates a 2x2 form with 10 pixel gaps.
func (fm *fmtag) simpleLayout() *layout {
	lo := &layout{}
	plan := []string{
		"ab",
		"cd",
	}
	lo.form = form.New(plan, fm.ww, fm.wh, "gap 5 5", "pad 5 5 5 5")
	lo.visualize(fm.scene)
	return lo
}

// spanLayout creates a form where the cell with the same label will
// span rows and columns.
func (fm *fmtag) spanLayout() *layout {
	lo := &layout{}
	plan := []string{
		"axxb",
		"cxxb",
		"cxxd",
		"eeef",
	}
	lo.form = form.New(plan, fm.ww, fm.wh, "gap 5 5", "pad 5 5 5 5")
	lo.visualize(fm.scene)
	return lo
}

// grabLayout creates a form where the base size dictates the max size
// that non-grabby rows and columns will grow to.
func (fm *fmtag) grabLayout() *layout {
	lo := &layout{}
	plan := []string{
		"abc",
		"def",
		"ghi",
	}
	lo.form = form.New(plan, 200, 200, "grabx 1", "graby 1", "gap 5 5", "pad 5 5 5 5")
	lo.visualize(fm.scene)
	return lo
}

// largeLayout creates a form with multiple spanning sections and
// a single spanning section.
func (fm *fmtag) largeLayout() *layout {
	lo := &layout{}
	plan := []string{
		"aabbbccd",
		"exxfyyyg",
		"exxhyyyg",
		"ixxhyyyg",
		"jklhmmnn",
	}
	lo.form = form.New(plan, fm.ww, fm.wh, "gap 5 5", "grabx 0", "graby 0", "pad 5 5 5 5")
	lo.visualize(fm.scene)
	return lo
}

// doubleLayout creates a form within a form to create a more complex layout.
func (fm *fmtag) doubleLayout() *layout {
	lo := &layout{}
	plan := []string{
		"abc",
		"def",
		"ghi",
	}
	lo.form = form.New(plan, 200, 200, "grabx 1", "graby 1", "gap 5 5", "pad 5 5 5 5")
	lo.visualize(fm.scene)
	lo.lo = fm.interiorLayout(lo.form.Section("e"))
	return lo
}

// interior layout is part of doubleLayout.
// It creates a second form inside the middle section of the first form.
func (fm *fmtag) interiorLayout(s form.Section) *layout {
	lo := &layout{}
	w, h := s.Size()
	iw, ih := int(lin.Round(w, 0)), int(lin.Round(h, 0))
	plan := []string{
		"pqr",
		"pqr",
		"stu",
	}
	lo.form = form.New(plan, iw, ih, "gap 5 5")
	lo.visualize(fm.scene)
	return lo
}

// =============================================================================

// layout
type layout struct {
	form   form.Form // cell and label position information.
	lo     *layout   // for doubleLayout.
	top    vu.Part   // single spot for making visible.
	sects  []vu.Part // visual representation of a form cell.
	labels []vu.Part // cell label.
}

// Called once to create the visual parts of a panel.
func (lo *layout) visualize(scene vu.Scene) {
	lo.top = scene.AddPart()
	lo.setVisible(false)
	lo.sects = make([]vu.Part, len(lo.form.Sections()))
	lo.labels = make([]vu.Part, len(lo.form.Sections()))
	for cnt, sect := range lo.form.Sections() {

		// place a box at the section location.
		lo.sects[cnt] = lo.top.AddPart()
		lo.sects[cnt].SetRole("uv").SetMesh("icon").AddTex("cell")

		// place the cell name in the middle of the cell.
		lo.labels[cnt] = lo.top.AddPart()
		lo.labels[cnt].SetRole("uv").AddTex("weblySleek16Black")
		if sect.Label() == "" {
			lo.labels[cnt].Role().SetFont("weblySleek16").SetPhrase("-")
		} else {
			lo.labels[cnt].Role().SetFont("weblySleek16").SetPhrase(sect.Label())
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
