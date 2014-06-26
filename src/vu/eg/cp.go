// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package main

import (
	"log"
	"vu"
	"vu/panel"
)

// cp (control panel) demos and tests the experimental vu/panel overlay package
// and panel rendering provided by vu/overlay.go. The panel API is trialed by
// this demo to determine its overall usefullness (unclear at the moment, still
// feels clunky).
func cp() {
	cp := &cptag{}
	ww, wh := 600, 800
	var err error
	if cp.eng, err = vu.New("Control Panel", 1200, 100, ww, wh); err != nil {
		log.Printf("wo: error intitializing engine %s", err)
		return
	}
	cp.eng.SetDirector(cp)  // override user input handling.
	defer cp.eng.Shutdown() // shut down the engine.
	defer catchErrors()
	cp.eng.Action()
}

// Encapsulate example specific data with a unique "tag".
type cptag struct {
	eng vu.Engine // 3D engine.

	// Widget related fields.
	over     panel.Panel     // main overlay control panel.
	diag     panel.Section   // dialog control panel.
	input    *panel.In       // reusable overlay input.
	pindex   []int           // picture index
	controls []panel.Control // picture controls.
	chars    []string        // trial pictures.
}

// Create is the engine intialization callback.
// Layout an initial panel with a few controls, and secondary (dialog).
func (cp *cptag) Create(eng vu.Engine) {
	_, _, ww, wh := eng.Size()

	// Create the controls on the simulated main screen overlay.
	// Use a ratio layout.
	cp.input = &panel.In{}
	cp.over = panel.NewPanel(ww, wh)
	cp.over.SetImg("panel")
	c1 := cp.over.AddControl().AddReaction("Lm", cp.newDialog)
	s1 := cp.over.AddSection()
	ts := cp.over.AddSection()
	cp.over.SetLayout(&panel.RatioLayout{ww, wh, map[uint]*panel.SizeInfo{
		c1.Id(): &panel.SizeInfo{270, 625, 80, 80},
		s1.Id(): &panel.SizeInfo{20, 20, 150, 150},
		ts.Id(): &panel.SizeInfo{200, 200, 200, 200},
	}})
	ts.SetImg("image") // test showing an image in a Section.
	ts.SetVisible(true)

	// Try a grid layout.
	s1.SetLayout(&panel.GridLayout{Columns: 2, Margin: 0})
	s1.AddControl().SetImg("image")
	s1.AddControl().SetImg("image")
	hv := s1.AddControl()
	hv.SetImg("image")
	hv.SetImgOnHover(true)

	// Mimic a popup dialog with a second overlay.
	cp.diag = cp.over.NewDialog(400, 400)
	cp.diag.SetVisible(true)
	cp.diag.SetImg("darkbg")
	s2 := cp.diag.AddSection()
	s3 := cp.diag.AddSection()
	cp.diag.SetLayout(&panel.FixedLayout{map[uint]*panel.FixedInfo{
		s2.Id(): &panel.FixedInfo{0, 50, 396, 346, panel.BL},
		s3.Id(): &panel.FixedInfo{0, 0, 400, 50, panel.BL},
	}})

	// Example character portrait panel.
	s2.SetLayout(&panel.GridLayout{Columns: 4, Margin: 4})
	c2 := s2.AddControl().AddReaction("Lm", func() { cp.show(0, 1) }).AddReaction("Rm", func() { cp.show(0, -1) })
	c3 := s2.AddControl().AddReaction("Lm", func() { cp.show(1, 1) }).AddReaction("Rm", func() { cp.show(1, -1) })
	c4 := s2.AddControl().AddReaction("Lm", func() { cp.show(2, 1) }).AddReaction("Rm", func() { cp.show(2, -1) })
	c5 := s2.AddControl().AddReaction("Lm", func() { cp.show(3, 1) }).AddReaction("Rm", func() { cp.show(3, -1) })

	// Ok, cancel buttons at the bottom of the character portrait panel.
	c6 := s3.AddControl().AddReaction("Lm", cp.closeDialog).SetImg("cancel")
	c7 := s3.AddControl().AddReaction("Rm", cp.closeDialog).SetImg("ok")
	s3.SetLayout(&panel.FixedLayout{map[uint]*panel.FixedInfo{
		c6.Id(): &panel.FixedInfo{5, 5, 80, 40, panel.BL},
		c7.Id(): &panel.FixedInfo{5, 5, 80, 40, panel.BR},
	}})

	// Give the buttons some images to cycle through.
	cp.controls = []panel.Control{c2, c3, c4, c5}
	cp.chars = []string{"char1", "char2", "char3", "char4"}
	cp.pindex = []int{0, 1, 2, 3}
	c2.SetImg(cp.chars[cp.pindex[0]])
	c3.SetImg(cp.chars[cp.pindex[1]])
	c4.SetImg(cp.chars[cp.pindex[2]])
	c5.SetImg(cp.chars[cp.pindex[3]])

	// Render the initial panel.
	cp.eng.SetPanel(cp.over)
	return
}

// Update is the regular engine callback.
func (cp *cptag) Update(in *vu.Input) {
	if in.Resized {
		cp.resize()
	}
	cp.input.Set(in.Mx, in.My, in.Scroll, in.Down)
	cp.over.React(cp.input)
}

// ============================================================================
// Button callbacks.

// newDialog flips from the first panel set to the second.
func (cp *cptag) newDialog() { cp.over.SetDialog(cp.diag) }

// closeDialog flips back to the first panel set.
func (cp *cptag) closeDialog() { cp.over.SetDialog(nil) }

// show the next character image.
func (cp *cptag) show(i, change int) {
	cp.pindex[i] = cp.next(cp.pindex[i], change, 4)
	cp.controls[i].SetImg(cp.chars[cp.pindex[i]])
}

// next cycles forwards or backwards through a limited set of indicies.
func (cp *cptag) next(i, change, max int) int {
	switch change {
	case 1:
		i = (i + 1) % max
	case -1:
		if i = i - 1; i < 0 {
			i = max - 1
		}
	}
	return i
}

// resize handles user screen/window changes.
func (cp *cptag) resize() {
	x, y, ww, wh := cp.eng.Size()
	cp.eng.Resize(x, y, ww, wh)
	cp.over.Resize(ww, wh)
}
