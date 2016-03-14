// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"log"

	"github.com/gazed/vu"
)

// kc explores treating the keyboard as a controller.
// It assigns each keyboard key a unique symbol.
//
// Treating the keyboard like a complicated console controller
// basically means ignoring its ability to input text. Overall
// a simplification over regluar keyboards, but having lot more
// potential controls than a console controller.
func kc() {
	kc := &kctag{}
	if err := vu.New(kc, "Keyboard Controller", 200, 200, 900, 400); err != nil {
		log.Printf("kc: error starting engine %s", err)
	}
	defer catchErrors()
}

// Globally unique "tag" that encapsulates example specific data.
type kctag struct {
	ui        vu.Camera   // 2D user interface.
	kb        vu.Pov      // Keyboard image.
	focus     vu.Pov      // Hilights first pressed key.
	positions map[int]pos // Screen position for each key.
}

// Create is the startup asset creation.
func (kc *kctag) Create(eng vu.Eng, s *vu.State) {
	top := eng.Root().NewPov()
	kc.ui = top.NewCam()
	kc.ui.SetUI()
	kc.positions = kc.keyPositions()

	// Create the keyboard image.
	kc.kb = top.NewPov().SetScale(900, 255, 0).SetLocation(450, 100+85, 0)
	kc.kb.NewModel("uv").LoadMesh("icon").AddTex("keyboard")

	// Pressed key focus
	kc.focus = top.NewPov().SetScale(50, 50, 0)
	kc.focus.NewModel("uv").LoadMesh("icon").AddTex("particle")

	// Place the key symbols over the keys.
	font := "lucidiaSu18"
	fontColor := "lucidiaSu18Black"
	for code, key := range kc.positions { // map key is key code, map value is key struct
		if char := vu.Keysym(code); char > 0 {
			cx, cy := key.location()
			letter := top.NewPov().SetLocation(cx, cy, 0)
			model := letter.NewModel("uv")
			model.AddTex(fontColor).LoadFont(font).SetPhrase(string(char))
		}
	}

	// Have a lighter default background.
	eng.SetColor(0.45, 0.45, 0.45, 1)
	kc.resize(s.W, s.H)
}

// Update is the regular engine callback.
func (kc *kctag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	if in.Resized {
		kc.resize(s.W, s.H)
	}

	// hilight the first pressed key.
	kc.focus.SetVisible(false)
	for press, _ := range in.Down {
		kc.focus.SetVisible(true)
		position := kc.positions[press]
		cx, cy := position.location()
		kc.focus.SetLocation(cx+6, cy+10, 0)
		break
	}
}
func (kc *kctag) resize(ww, wh int) {
	kc.ui.SetOrthographic(0, float64(ww), 0, float64(wh), 0, 10)
}

// Position the keys on the keyboard image.
func (kc *kctag) keyPositions() map[int]pos {
	return map[int]pos{
		vu.K_0:     pos{col: 9, row: 4, xoff: 0.8, yoff: 0.0},
		vu.K_1:     pos{col: 1, row: 4, xoff: 0.0, yoff: 0.0},
		vu.K_2:     pos{col: 2, row: 4, xoff: 0.0, yoff: 0.0},
		vu.K_3:     pos{col: 3, row: 4, xoff: 0.0, yoff: 0.0},
		vu.K_4:     pos{col: 4, row: 4, xoff: 0.0, yoff: 0.0},
		vu.K_5:     pos{col: 5, row: 4, xoff: 0.0, yoff: 0.0},
		vu.K_6:     pos{col: 5, row: 4, xoff: 0.9, yoff: 0.0},
		vu.K_7:     pos{col: 6, row: 4, xoff: 0.8, yoff: 0.0},
		vu.K_8:     pos{col: 7, row: 4, xoff: 0.8, yoff: 0.0},
		vu.K_9:     pos{col: 8, row: 4, xoff: 0.8, yoff: 0.0},
		vu.K_A:     pos{col: 1, row: 2, xoff: 0.9, yoff: 0.0},
		vu.K_B:     pos{col: 6, row: 1, xoff: 0.1, yoff: 0.0},
		vu.K_C:     pos{col: 4, row: 1, xoff: 0.2, yoff: 0.0},
		vu.K_D:     pos{col: 3, row: 2, xoff: 0.8, yoff: 0.0},
		vu.K_E:     pos{col: 3, row: 3, xoff: 0.5, yoff: 0.0},
		vu.K_F:     pos{col: 4, row: 2, xoff: 0.7, yoff: 0.0},
		vu.K_G:     pos{col: 5, row: 2, xoff: 0.6, yoff: 0.0},
		vu.K_H:     pos{col: 6, row: 2, xoff: 0.5, yoff: 0.0},
		vu.K_I:     pos{col: 8, row: 3, xoff: 0.3, yoff: 0.0},
		vu.K_J:     pos{col: 7, row: 2, xoff: 0.5, yoff: 0.0},
		vu.K_K:     pos{col: 8, row: 2, xoff: 0.5, yoff: 0.0},
		vu.K_L:     pos{col: 9, row: 2, xoff: 0.5, yoff: 0.0},
		vu.K_M:     pos{col: 8, row: 1, xoff: 0.0, yoff: 0.0},
		vu.K_N:     pos{col: 7, row: 1, xoff: 0.0, yoff: 0.0},
		vu.K_O:     pos{col: 9, row: 3, xoff: 0.1, yoff: 0.0},
		vu.K_P:     pos{col: 10, row: 3, xoff: 0.1, yoff: 0.0},
		vu.K_Q:     pos{col: 1, row: 3, xoff: 0.6, yoff: 0.0},
		vu.K_R:     pos{col: 4, row: 3, xoff: 0.5, yoff: 0.0},
		vu.K_S:     pos{col: 2, row: 2, xoff: 0.8, yoff: 0.0},
		vu.K_T:     pos{col: 5, row: 3, xoff: 0.4, yoff: 0.0},
		vu.K_U:     pos{col: 7, row: 3, xoff: 0.2, yoff: 0.0},
		vu.K_V:     pos{col: 5, row: 1, xoff: 0.1, yoff: 0.0},
		vu.K_W:     pos{col: 2, row: 3, xoff: 0.5, yoff: 0.0},
		vu.K_X:     pos{col: 3, row: 1, xoff: 0.3, yoff: 0.0},
		vu.K_Y:     pos{col: 6, row: 3, xoff: 0.3, yoff: 0.0},
		vu.K_Z:     pos{col: 2, row: 1, xoff: 0.2, yoff: 0.0},
		vu.K_Equal: pos{col: 11, row: 4, xoff: 0.6, yoff: 0.0},
		vu.K_Minus: pos{col: 10, row: 4, xoff: 0.7, yoff: 0.0},
		vu.K_RBkt:  pos{col: 12, row: 3, xoff: 0.1, yoff: 0.0},
		vu.K_LBkt:  pos{col: 11, row: 3, xoff: 0.1, yoff: 0.0},
		vu.K_Qt:    pos{col: 11, row: 2, xoff: 0.5, yoff: 0.0},
		vu.K_Semi:  pos{col: 10, row: 2, xoff: 0.5, yoff: 0.0},
		vu.K_BSl:   pos{col: 13, row: 3, xoff: 0.0, yoff: 0.0},
		vu.K_Comma: pos{col: 9, row: 1, xoff: 0.0, yoff: 0.0},
		vu.K_Slash: pos{col: 11, row: 1, xoff: 0.0, yoff: 0.0},
		vu.K_Dot:   pos{col: 10, row: 1, xoff: 0.0, yoff: 0.0},
		vu.K_Grave: pos{col: 0, row: 4, xoff: 0.0, yoff: 0.0},
		vu.K_Ret:   pos{col: 12, row: 2, xoff: 0.5, yoff: 0.0},
		vu.K_Tab:   pos{col: 0, row: 3, xoff: 0.0, yoff: 0.0},
		vu.K_Space: pos{col: 7, row: 0, xoff: 0.0, yoff: 0.0},
		vu.K_Del:   pos{col: 12, row: 4, xoff: 0.6, yoff: 0.0},
		vu.K_Esc:   pos{col: 0, row: 5, xoff: 0.0, yoff: 0.0},
		vu.K_F1:    pos{col: 1, row: 5, xoff: 0.0, yoff: 0.0},
		vu.K_F2:    pos{col: 2, row: 5, xoff: 0.0, yoff: 0.0},
		vu.K_F3:    pos{col: 3, row: 5, xoff: 0.0, yoff: 0.0},
		vu.K_F4:    pos{col: 4, row: 5, xoff: 0.0, yoff: 0.0},
		vu.K_F5:    pos{col: 5, row: 5, xoff: 0.0, yoff: 0.0},
		vu.K_F6:    pos{col: 6, row: 5, xoff: 0.0, yoff: 0.0},
		vu.K_F7:    pos{col: 7, row: 5, xoff: 0.0, yoff: 0.0},
		vu.K_F8:    pos{col: 8, row: 5, xoff: 0.0, yoff: 0.0},
		vu.K_F9:    pos{col: 9, row: 5, xoff: 0.0, yoff: 0.0},
		vu.K_F10:   pos{col: 10, row: 5, xoff: 0.0, yoff: 0.0},
		vu.K_F11:   pos{col: 11, row: 5, xoff: 0.0, yoff: 0.0},
		vu.K_F12:   pos{col: 12, row: 5, xoff: 0.0, yoff: 0.0},
		vu.K_F13:   pos{col: 14, row: 5, xoff: 0.0, yoff: 0.0},
		vu.K_F14:   pos{col: 15, row: 5, xoff: 0.0, yoff: 0.0},
		vu.K_F15:   pos{col: 16, row: 5, xoff: 0.0, yoff: 0.0},
		vu.K_F16:   pos{col: 17, row: 5, xoff: 0.0, yoff: 0.0},
		vu.K_F17:   pos{col: 18, row: 5, xoff: 0.0, yoff: 0.0},
		vu.K_F18:   pos{col: 19, row: 5, xoff: 0.0, yoff: 0.0},
		vu.K_F19:   pos{col: 20, row: 5, xoff: 0.0, yoff: 0.0},
		vu.K_Home:  pos{col: 15, row: 4, xoff: 0.0, yoff: 0.0},
		vu.K_PgUp:  pos{col: 16, row: 4, xoff: 0.0, yoff: 0.0},
		vu.K_FDel:  pos{col: 14, row: 3, xoff: 0.0, yoff: 0.0},
		vu.K_End:   pos{col: 15, row: 3, xoff: 0.0, yoff: 0.0},
		vu.K_PgDn:  pos{col: 16, row: 3, xoff: 0.0, yoff: 0.0},
		vu.K_La:    pos{col: 14, row: 0, xoff: 0.0, yoff: 0.0},
		vu.K_Ra:    pos{col: 16, row: 0, xoff: 0.0, yoff: 0.0},
		vu.K_Da:    pos{col: 15, row: 0, xoff: 0.0, yoff: 0.0},
		vu.K_Ua:    pos{col: 15, row: 1, xoff: 0.0, yoff: 0.0},
		vu.K_KpDot: pos{col: 19, row: 0, xoff: 0.0, yoff: 0.0},
		vu.K_KpMlt: pos{col: 20, row: 4, xoff: 0.0, yoff: 0.0},
		vu.K_KpAdd: pos{col: 20, row: 2, xoff: 0.0, yoff: 0.0},
		vu.K_KpClr: pos{col: 17, row: 4, xoff: 0.0, yoff: 0.0},
		vu.K_KpDiv: pos{col: 19, row: 4, xoff: 0.0, yoff: 0.0},
		vu.K_KpEnt: pos{col: 20, row: 0, xoff: 0.0, yoff: 0.0},
		vu.K_KpSub: pos{col: 20, row: 3, xoff: 0.0, yoff: 0.0},
		vu.K_KpEql: pos{col: 18, row: 4, xoff: 0.0, yoff: 0.0},
		vu.K_Kp0:   pos{col: 17, row: 0, xoff: 0.0, yoff: 0.0},
		vu.K_Kp1:   pos{col: 17, row: 1, xoff: 0.0, yoff: 0.0},
		vu.K_Kp2:   pos{col: 18, row: 1, xoff: 0.0, yoff: 0.0},
		vu.K_Kp3:   pos{col: 19, row: 1, xoff: 0.0, yoff: 0.0},
		vu.K_Kp4:   pos{col: 17, row: 2, xoff: 0.0, yoff: 0.0},
		vu.K_Kp5:   pos{col: 18, row: 2, xoff: 0.0, yoff: 0.0},
		vu.K_Kp6:   pos{col: 19, row: 2, xoff: 0.0, yoff: 0.0},
		vu.K_Kp7:   pos{col: 17, row: 3, xoff: 0.0, yoff: 0.0},
		vu.K_Kp8:   pos{col: 18, row: 3, xoff: 0.0, yoff: 0.0},
		vu.K_Kp9:   pos{col: 19, row: 3, xoff: 0.0, yoff: 0.0},
		vu.K_Lm:    pos{col: 1, row: 6, xoff: 0.0, yoff: 0.0},
		vu.K_Mm:    pos{col: 1, row: 6, xoff: 0.5, yoff: 0.0},
		vu.K_Rm:    pos{col: 2, row: 6, xoff: 0.0, yoff: 0.0},
		vu.K_Ctl:   pos{col: 0, row: 0, xoff: 0.0, yoff: 0.0},
		vu.K_Fn:    pos{col: 14, row: 4, xoff: 0.0, yoff: 0.0},
		vu.K_Shift: pos{col: 0, row: 1, xoff: 0.0, yoff: 0.0},
		vu.K_Cmd:   pos{col: 3, row: 0, xoff: 0.0, yoff: 0.0},
		vu.K_Alt:   pos{col: 1, row: 0, xoff: 0.6, yoff: 0.0},
	}
}

// pos is used to locate each rune on the keyboard image.
type pos struct {
	col  int     // keyboard position.. 21 columns.
	row  int     // keyboard position.. 6 rows + 1 mouse row.
	xoff float64 // column offustment.
	yoff float64 // row offustment.
}

// location gives a positions x, y location in screen pixels.
func (p *pos) location() (x, y float64) {
	xspan := 41.0
	yspan := 38.0
	x = 25.0 + (float64(p.col)+p.xoff)*xspan
	y = 85.0 + (float64(p.row)+p.yoff)*yspan
	if p.col > 13 {
		x += 12 // first gap
	}
	if p.col > 16 {
		x += 12 // second gap
	}
	return x, y
}
