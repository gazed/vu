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
// a simplification over regular keyboards, but having alot more
// potential controls than a console controller.
//
// CONTROLS:
//   key   : highlight key press
//   mouse : highlight mouse click
func kc() {
	kc := &kctag{}
	if err := vu.New(kc, "Keyboard Controller", 200, 200, 900, 400); err != nil {
		log.Printf("kc: error starting engine %s", err)
	}
	defer catchErrors()
}

// Globally unique "tag" that encapsulates example specific data.
type kctag struct {
	ui        *vu.Camera  // 2D user interface.
	kb        *vu.Pov     // Keyboard image.
	focus     *vu.Pov     // Hilights first pressed key.
	positions map[int]pos // Screen position for each key.
}

// Create is the startup asset creation.
func (kc *kctag) Create(eng vu.Eng, s *vu.State) {
	top := eng.Root().NewPov()
	kc.ui = top.NewCam().SetUI()
	kc.positions = kc.keyPositions()

	// Create the keyboard image.
	kc.kb = top.NewPov().SetScale(900, 255, 0).SetAt(450, 100+85, 0)
	kc.kb.NewModel("uv", "msh:icon", "tex:keyboard")

	// Pressed key focus
	kc.focus = top.NewPov().SetScale(50, 50, 0)
	kc.focus.NewModel("uv", "msh:icon", "tex:particle")

	// Place the key symbols over the keys.
	for code, key := range kc.positions { // map key is key code, map value is key struct
		if char := vu.Keysym(code); char > 0 {
			cx, cy := key.location()
			letter := top.NewPov().SetAt(cx, cy, 0)
			letter.NewLabel("uv", "lucidiaSu18", "lucidiaSu18Black").SetStr(string(char))
		}
	}

	// Have a lighter default background.
	eng.Set(vu.Color(0.45, 0.45, 0.45, 1))
	kc.resize(s.W, s.H)
}

// Update is the regular engine callback.
func (kc *kctag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	if in.Resized {
		kc.resize(s.W, s.H)
	}

	// hilight the first pressed key.
	kc.focus.Cull = true
	for press := range in.Down {
		kc.focus.Cull = false
		position := kc.positions[press]
		cx, cy := position.location()
		kc.focus.SetAt(cx+6, cy+10, 0)
		break
	}
}
func (kc *kctag) resize(ww, wh int) {
	kc.ui.SetOrthographic(0, float64(ww), 0, float64(wh), 0, 10)
}

// Position the keys on the keyboard image.
func (kc *kctag) keyPositions() map[int]pos {
	return map[int]pos{
		vu.K0:     {col: 9, row: 4, xoff: 0.8, yoff: 0.0},
		vu.K1:     {col: 1, row: 4, xoff: 0.0, yoff: 0.0},
		vu.K2:     {col: 2, row: 4, xoff: 0.0, yoff: 0.0},
		vu.K3:     {col: 3, row: 4, xoff: 0.0, yoff: 0.0},
		vu.K4:     {col: 4, row: 4, xoff: 0.0, yoff: 0.0},
		vu.K5:     {col: 5, row: 4, xoff: 0.0, yoff: 0.0},
		vu.K6:     {col: 5, row: 4, xoff: 0.9, yoff: 0.0},
		vu.K7:     {col: 6, row: 4, xoff: 0.8, yoff: 0.0},
		vu.K8:     {col: 7, row: 4, xoff: 0.8, yoff: 0.0},
		vu.K9:     {col: 8, row: 4, xoff: 0.8, yoff: 0.0},
		vu.KA:     {col: 1, row: 2, xoff: 0.9, yoff: 0.0},
		vu.KB:     {col: 6, row: 1, xoff: 0.1, yoff: 0.0},
		vu.KC:     {col: 4, row: 1, xoff: 0.2, yoff: 0.0},
		vu.KD:     {col: 3, row: 2, xoff: 0.8, yoff: 0.0},
		vu.KE:     {col: 3, row: 3, xoff: 0.5, yoff: 0.0},
		vu.KF:     {col: 4, row: 2, xoff: 0.7, yoff: 0.0},
		vu.KG:     {col: 5, row: 2, xoff: 0.6, yoff: 0.0},
		vu.KH:     {col: 6, row: 2, xoff: 0.5, yoff: 0.0},
		vu.KI:     {col: 8, row: 3, xoff: 0.3, yoff: 0.0},
		vu.KJ:     {col: 7, row: 2, xoff: 0.5, yoff: 0.0},
		vu.KK:     {col: 8, row: 2, xoff: 0.5, yoff: 0.0},
		vu.KL:     {col: 9, row: 2, xoff: 0.5, yoff: 0.0},
		vu.KM:     {col: 8, row: 1, xoff: 0.0, yoff: 0.0},
		vu.KN:     {col: 7, row: 1, xoff: 0.0, yoff: 0.0},
		vu.KO:     {col: 9, row: 3, xoff: 0.1, yoff: 0.0},
		vu.KP:     {col: 10, row: 3, xoff: 0.1, yoff: 0.0},
		vu.KQ:     {col: 1, row: 3, xoff: 0.6, yoff: 0.0},
		vu.KR:     {col: 4, row: 3, xoff: 0.5, yoff: 0.0},
		vu.KS:     {col: 2, row: 2, xoff: 0.8, yoff: 0.0},
		vu.KT:     {col: 5, row: 3, xoff: 0.4, yoff: 0.0},
		vu.KU:     {col: 7, row: 3, xoff: 0.2, yoff: 0.0},
		vu.KV:     {col: 5, row: 1, xoff: 0.1, yoff: 0.0},
		vu.KW:     {col: 2, row: 3, xoff: 0.5, yoff: 0.0},
		vu.KX:     {col: 3, row: 1, xoff: 0.3, yoff: 0.0},
		vu.KY:     {col: 6, row: 3, xoff: 0.3, yoff: 0.0},
		vu.KZ:     {col: 2, row: 1, xoff: 0.2, yoff: 0.0},
		vu.KEqual: {col: 11, row: 4, xoff: 0.6, yoff: 0.0},
		vu.KMinus: {col: 10, row: 4, xoff: 0.7, yoff: 0.0},
		vu.KRBkt:  {col: 12, row: 3, xoff: 0.1, yoff: 0.0},
		vu.KLBkt:  {col: 11, row: 3, xoff: 0.1, yoff: 0.0},
		vu.KQt:    {col: 11, row: 2, xoff: 0.5, yoff: 0.0},
		vu.KSemi:  {col: 10, row: 2, xoff: 0.5, yoff: 0.0},
		vu.KBSl:   {col: 13, row: 3, xoff: 0.0, yoff: 0.0},
		vu.KComma: {col: 9, row: 1, xoff: 0.0, yoff: 0.0},
		vu.KSlash: {col: 11, row: 1, xoff: 0.0, yoff: 0.0},
		vu.KDot:   {col: 10, row: 1, xoff: 0.0, yoff: 0.0},
		vu.KGrave: {col: 0, row: 4, xoff: 0.0, yoff: 0.0},
		vu.KRet:   {col: 12, row: 2, xoff: 0.5, yoff: 0.0},
		vu.KTab:   {col: 0, row: 3, xoff: 0.0, yoff: 0.0},
		vu.KSpace: {col: 7, row: 0, xoff: 0.0, yoff: 0.0},
		vu.KDel:   {col: 12, row: 4, xoff: 0.6, yoff: 0.0},
		vu.KEsc:   {col: 0, row: 5, xoff: 0.0, yoff: 0.0},
		vu.KF1:    {col: 1, row: 5, xoff: 0.0, yoff: 0.0},
		vu.KF2:    {col: 2, row: 5, xoff: 0.0, yoff: 0.0},
		vu.KF3:    {col: 3, row: 5, xoff: 0.0, yoff: 0.0},
		vu.KF4:    {col: 4, row: 5, xoff: 0.0, yoff: 0.0},
		vu.KF5:    {col: 5, row: 5, xoff: 0.0, yoff: 0.0},
		vu.KF6:    {col: 6, row: 5, xoff: 0.0, yoff: 0.0},
		vu.KF7:    {col: 7, row: 5, xoff: 0.0, yoff: 0.0},
		vu.KF8:    {col: 8, row: 5, xoff: 0.0, yoff: 0.0},
		vu.KF9:    {col: 9, row: 5, xoff: 0.0, yoff: 0.0},
		vu.KF10:   {col: 10, row: 5, xoff: 0.0, yoff: 0.0},
		vu.KF11:   {col: 11, row: 5, xoff: 0.0, yoff: 0.0},
		vu.KF12:   {col: 12, row: 5, xoff: 0.0, yoff: 0.0},
		vu.KF13:   {col: 14, row: 5, xoff: 0.0, yoff: 0.0},
		vu.KF14:   {col: 15, row: 5, xoff: 0.0, yoff: 0.0},
		vu.KF15:   {col: 16, row: 5, xoff: 0.0, yoff: 0.0},
		vu.KF16:   {col: 17, row: 5, xoff: 0.0, yoff: 0.0},
		vu.KF17:   {col: 18, row: 5, xoff: 0.0, yoff: 0.0},
		vu.KF18:   {col: 19, row: 5, xoff: 0.0, yoff: 0.0},
		vu.KF19:   {col: 20, row: 5, xoff: 0.0, yoff: 0.0},
		vu.KHome:  {col: 15, row: 4, xoff: 0.0, yoff: 0.0},
		vu.KPgUp:  {col: 16, row: 4, xoff: 0.0, yoff: 0.0},
		vu.KFDel:  {col: 14, row: 3, xoff: 0.0, yoff: 0.0},
		vu.KEnd:   {col: 15, row: 3, xoff: 0.0, yoff: 0.0},
		vu.KPgDn:  {col: 16, row: 3, xoff: 0.0, yoff: 0.0},
		vu.KLa:    {col: 14, row: 0, xoff: 0.0, yoff: 0.0},
		vu.KRa:    {col: 16, row: 0, xoff: 0.0, yoff: 0.0},
		vu.KDa:    {col: 15, row: 0, xoff: 0.0, yoff: 0.0},
		vu.KUa:    {col: 15, row: 1, xoff: 0.0, yoff: 0.0},
		vu.KKpDot: {col: 19, row: 0, xoff: 0.0, yoff: 0.0},
		vu.KKpMlt: {col: 20, row: 4, xoff: 0.0, yoff: 0.0},
		vu.KKpAdd: {col: 20, row: 2, xoff: 0.0, yoff: 0.0},
		vu.KKpClr: {col: 17, row: 4, xoff: 0.0, yoff: 0.0},
		vu.KKpDiv: {col: 19, row: 4, xoff: 0.0, yoff: 0.0},
		vu.KKpEnt: {col: 20, row: 0, xoff: 0.0, yoff: 0.0},
		vu.KKpSub: {col: 20, row: 3, xoff: 0.0, yoff: 0.0},
		vu.KKpEql: {col: 18, row: 4, xoff: 0.0, yoff: 0.0},
		vu.KKp0:   {col: 17, row: 0, xoff: 0.0, yoff: 0.0},
		vu.KKp1:   {col: 17, row: 1, xoff: 0.0, yoff: 0.0},
		vu.KKp2:   {col: 18, row: 1, xoff: 0.0, yoff: 0.0},
		vu.KKp3:   {col: 19, row: 1, xoff: 0.0, yoff: 0.0},
		vu.KKp4:   {col: 17, row: 2, xoff: 0.0, yoff: 0.0},
		vu.KKp5:   {col: 18, row: 2, xoff: 0.0, yoff: 0.0},
		vu.KKp6:   {col: 19, row: 2, xoff: 0.0, yoff: 0.0},
		vu.KKp7:   {col: 17, row: 3, xoff: 0.0, yoff: 0.0},
		vu.KKp8:   {col: 18, row: 3, xoff: 0.0, yoff: 0.0},
		vu.KKp9:   {col: 19, row: 3, xoff: 0.0, yoff: 0.0},
		vu.KLm:    {col: 1, row: 6, xoff: 0.0, yoff: 0.0},
		vu.KMm:    {col: 1, row: 6, xoff: 0.5, yoff: 0.0},
		vu.KRm:    {col: 2, row: 6, xoff: 0.0, yoff: 0.0},
		vu.KCtl:   {col: 0, row: 0, xoff: 0.0, yoff: 0.0},
		vu.KFn:    {col: 14, row: 4, xoff: 0.0, yoff: 0.0},
		vu.KShift: {col: 0, row: 1, xoff: 0.0, yoff: 0.0},
		vu.KCmd:   {col: 3, row: 0, xoff: 0.0, yoff: 0.0},
		vu.KAlt:   {col: 1, row: 0, xoff: 0.6, yoff: 0.0},
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
