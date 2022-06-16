// Copyright © 2013-2024 Galvanized Logic Inc.

package main

import (
	"log/slog"
	"time"

	"github.com/gazed/vu"
)

// kc explores treating the keyboard as a controller,
// focusing on button presses and ignoring its ability to input text.
// The image shows a macos keyboard because that's what was available
// at the time. This example demonstrates:
//   - loading assets.
//   - creating a 2D scene with image and text models.
//   - reacting to user input
//
// CONTROLS:
//   - key   : highlight key press
//   - mouse : highlight mouse click
func kc() {
	kc := &kctag{ww: 1200, wh: 340} // match keyboard image size.
	eng, err := vu.NewEngine(
		vu.Windowed(),
		vu.Title("Keyboard Controller"),
		vu.Size(200, 200, int32(kc.ww), int32(kc.wh)),
		vu.Background(0.1, 0.1, 0.5, 1.0),
	)
	if err != nil {
		slog.Error("kc: engine start", "err", err)
		return
	}

	// import assets from asset files.
	// This creates the assets referenced by the models below.
	eng.ImportAssets("icon.shd", "keyboard.png", "keyboard_press.png")  // load some assets
	eng.ImportAssets("label.shd", "lucidiaSu18.fnt", "lucidiaSu18.png") // load more assets

	// create a 2D scene with a camera.
	scene := eng.AddScene(vu.Scene2D)

	// add the keyboard image to the scene.
	kb := scene.AddModel("shd:icon", "msh:icon", "tex:color:keyboard")
	kb.SetScale(float64(kc.ww), float64(kc.wh), 0).SetAt(float64(kc.ww/2), float64(kc.wh/2), 0)

	// add the pressed key focus model to the scene
	kc.focus = scene.AddModel("shd:icon", "msh:icon", "tex:color:keyboard_press")
	kc.focus.SetScale(50, 50, 0) // make it bigger
	kc.focus.Cull(true)          // hide until a key is pressed.

	// Place the key symbols over the keys.
	// map key is key code, map value is key struct
	kc.positions = kc.keyPositions()
	for code, key := range kc.positions {
		if char := vu.Symbol(code); char > 0 {
			cx, cy := key.location(kc.ww, kc.wh)
			letter := scene.AddLabel(string(char), 0, "shd:label", "fnt:lucidiaSu18", "tex:color:lucidiaSu18")
			letter.SetAt(cx-4, cy-9, 0).SetColor(0, 0, 0, 1) // black
		}
	}
	defer catchErrors()
	eng.Run(kc) // does not return while example is running.
}

// Globally unique "tag" that encapsulates example specific data.
type kctag struct {
	ww        int         // window width
	wh        int         // window height
	focus     *vu.Entity  // Hilights first pressed key.
	positions map[int]pos // Screen position for each key.
}

// Update is the application engine callback.
func (kc *kctag) Update(eng *vu.Engine, in *vu.Input, delta time.Duration) {
	// react to one time press events.
	for press := range in.Pressed {
		switch press {
		case vu.KX:
			// quit if X is pressed
			eng.Shutdown()
			return
		}
	}

	// highlight the most recently pressed key
	kc.focus.Cull(true)
	pressed := -1
	lastTime := time.Time{}
	for key, timePressed := range in.Down {
		if timePressed.After(lastTime) {
			pressed = int(key)
			lastTime = timePressed
		}
	}
	if pressed != -1 {
		if position, ok := kc.positions[pressed]; ok {
			cx, cy := position.location(kc.ww, kc.wh)
			kc.focus.SetAt(cx, cy, 0)
			kc.focus.Cull(false)
		}
	}
}

// Position the keys on the keyboard image.
// The positions are in pixels based on the original image.
func (kc *kctag) keyPositions() map[int]pos {
	return map[int]pos{
		// top row
		vu.KEsc: {xoff: 46, yoff: 39},
		vu.KF1:  {xoff: 100, yoff: 39},
		vu.KF2:  {xoff: 154, yoff: 39},
		vu.KF3:  {xoff: 208, yoff: 39},
		vu.KF4:  {xoff: 264, yoff: 39},
		vu.KF5:  {xoff: 320, yoff: 39},
		vu.KF6:  {xoff: 374, yoff: 39},
		vu.KF7:  {xoff: 426, yoff: 39},
		vu.KF8:  {xoff: 480, yoff: 39},
		vu.KF9:  {xoff: 536, yoff: 39},
		vu.KF10: {xoff: 590, yoff: 39},
		vu.KF11: {xoff: 644, yoff: 39},
		vu.KF12: {xoff: 698, yoff: 39},

		// ignore the following keys,
		// ie: don't use print screen as a control and limit function keys to 12.
		// vu.KF13: {xoff: 748, yoff: 39},
		// vu.KF14: {xoff: 824, yoff: 39},
		// vu.KF15: {xoff: 878, yoff: 39},
		// vu.KF16: {xoff: 932, yoff: 39},
		// vu.KF17: {xoff: 1000, yoff: 39},
		// vu.KF18: {xoff: 1054, yoff: 39},
		// vu.KF19: {xoff: 1108, yoff: 39},

		// second row
		vu.KGrave: {xoff: 46, yoff: 82},
		vu.K1:     {xoff: 98, yoff: 82},
		vu.K2:     {xoff: 150, yoff: 82},
		vu.K3:     {xoff: 202, yoff: 82},
		vu.K4:     {xoff: 254, yoff: 82},
		vu.K5:     {xoff: 306, yoff: 82},
		vu.K6:     {xoff: 358, yoff: 82},
		vu.K7:     {xoff: 410, yoff: 82},
		vu.K8:     {xoff: 462, yoff: 82},
		vu.K9:     {xoff: 514, yoff: 82},
		vu.K0:     {xoff: 566, yoff: 82},
		vu.KMinus: {xoff: 620, yoff: 82},
		vu.KEqual: {xoff: 672, yoff: 82},
		vu.KDel:   {xoff: 734, yoff: 82}, // back delete
		// vu.KFn:    {xoff: 824, yoff: 82}, // ignore win: insert
		vu.KHome: {xoff: 874, yoff: 82},
		vu.KPgUp: {xoff: 928, yoff: 82},
		vu.KPClr: {xoff: 1000, yoff: 82}, // num lock
		vu.KPEql: {xoff: 1054, yoff: 82}, // divide
		vu.KPDiv: {xoff: 1106, yoff: 82}, // multiply
		vu.KPMlt: {xoff: 1158, yoff: 82}, // minus

		// third row
		vu.KTab:  {xoff: 58, yoff: 132},
		vu.KQ:    {xoff: 124, yoff: 132},
		vu.KW:    {xoff: 176, yoff: 132},
		vu.KE:    {xoff: 229, yoff: 132},
		vu.KR:    {xoff: 282, yoff: 132},
		vu.KT:    {xoff: 335, yoff: 132},
		vu.KY:    {xoff: 386, yoff: 132},
		vu.KU:    {xoff: 438, yoff: 132},
		vu.KI:    {xoff: 490, yoff: 132},
		vu.KO:    {xoff: 542, yoff: 132},
		vu.KP:    {xoff: 594, yoff: 132},
		vu.KLBkt: {xoff: 646, yoff: 132},
		vu.KRBkt: {xoff: 698, yoff: 132},
		vu.KBSl:  {xoff: 750, yoff: 132},
		vu.KFDel: {xoff: 824, yoff: 132},
		vu.KEnd:  {xoff: 876, yoff: 132},
		vu.KPgDn: {xoff: 928, yoff: 132},
		vu.KP7:   {xoff: 1000, yoff: 132},
		vu.KP8:   {xoff: 1054, yoff: 132},
		vu.KP9:   {xoff: 1106, yoff: 132},
		vu.KPSub: {xoff: 1158, yoff: 132},

		// fourth row
		vu.KA:     {xoff: 137, yoff: 184},
		vu.KS:     {xoff: 190, yoff: 184},
		vu.KD:     {xoff: 242, yoff: 184},
		vu.KF:     {xoff: 294, yoff: 184},
		vu.KG:     {xoff: 346, yoff: 184},
		vu.KH:     {xoff: 398, yoff: 184},
		vu.KJ:     {xoff: 450, yoff: 184},
		vu.KK:     {xoff: 502, yoff: 184},
		vu.KL:     {xoff: 554, yoff: 184},
		vu.KSemi:  {xoff: 608, yoff: 184},
		vu.KQuote: {xoff: 660, yoff: 184},
		vu.KRet:   {xoff: 730, yoff: 184},
		vu.KML:    {xoff: 824, yoff: 184}, // left mouse
		vu.KMM:    {xoff: 876, yoff: 184}, // middle mouse
		vu.KMR:    {xoff: 928, yoff: 184}, // right mouse
		vu.KP4:    {xoff: 1000, yoff: 184},
		vu.KP5:    {xoff: 1054, yoff: 184},
		vu.KP6:    {xoff: 1106, yoff: 184},
		vu.KPAdd:  {xoff: 1158, yoff: 184},

		// fifth row
		vu.KShift: {xoff: 75, yoff: 234},
		vu.KZ:     {xoff: 164, yoff: 234},
		vu.KX:     {xoff: 216, yoff: 234},
		vu.KC:     {xoff: 268, yoff: 234},
		vu.KV:     {xoff: 320, yoff: 234},
		vu.KB:     {xoff: 372, yoff: 234},
		vu.KN:     {xoff: 424, yoff: 234},
		vu.KM:     {xoff: 476, yoff: 234},
		vu.KComma: {xoff: 530, yoff: 234},
		vu.KDot:   {xoff: 582, yoff: 234},
		vu.KSlash: {xoff: 634, yoff: 234},
		vu.KAUp:   {xoff: 876, yoff: 234},
		vu.KP1:    {xoff: 1000, yoff: 234},
		vu.KP2:    {xoff: 1054, yoff: 234},
		vu.KP3:    {xoff: 1106, yoff: 234},

		// sixth row
		vu.KCtl:    {xoff: 58, yoff: 290},
		vu.KAlt:    {xoff: 200, yoff: 290}, // macos: command key
		vu.KCmd:    {xoff: 200, yoff: 290},
		vu.KSpace:  {xoff: 400, yoff: 290},
		vu.KALeft:  {xoff: 824, yoff: 290},
		vu.KADown:  {xoff: 876, yoff: 290},
		vu.KARight: {xoff: 928, yoff: 290},
		vu.KP0:     {xoff: 1026, yoff: 290},
		vu.KPDot:   {xoff: 1105, yoff: 290},
		vu.KPEnt:   {xoff: 1158, yoff: 260},
	}
}

// pos is used to locate each rune on the keyboard image.
type pos struct {
	xoff float64 // x offset in pixels from top left.
	yoff float64 // y offset in pixels from top left.
}

// keyboard image is 1200.0x340.0 and the pixel positions
// of the keys are measured in pixels from the image.
func (p *pos) location(ww, wh int) (sx, sy float64) { return p.xoff, p.yoff }
