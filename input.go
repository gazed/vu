// Copyright © 2014-2024 Galvanized Logic Inc.

package vu

// input.go wraps device package input as a convenience so the
// device package does not always need to be included.
//
// FUTURE: Have the device package support text entry.

import (
	"github.com/gazed/vu/device"
)

// Input is used to communicate user input to the application.
//
// User input is the current cursor location, current pressed keys,
// mouse buttons, and modifiers. These are sent to the game each
// Update() callback. Overall the keyboard is treated like a gamepad
// controller. Keys and buttons are pressed or not pressed.
type Input device.Input

// Clone the data from given Input structure into this one.
func (in *Input) Clone(b *device.Input) {
	in.Mx = b.Mx
	in.My = b.My
	in.Focus = b.Focus
	in.Scroll = b.Scroll

	// clear current keymaps
	for key := range in.Pressed {
		delete(in.Pressed, key)
	}
	for key := range in.Down {
		delete(in.Down, key)
	}
	for key := range in.Released {
		delete(in.Released, key)
	}

	// clone the given device keymaps into these keymaps
	for k, v := range b.Pressed {
		in.Pressed[k] = v
	}
	for k, v := range b.Down {
		in.Down[k] = v
	}
	for k, v := range b.Released {
		in.Released[k] = v
	}
}

// Expose the device package keys as a convenience so the
// device package does not always need to be included.
// The symbol associated to each key is shown in the comments.
//
// Keys are expected to be used for controlling game actions.
// There is no text entry or text layout support.
const (
	K0      = device.K0      // 0 48     Standard keyboard numbers.
	K1      = device.K1      // 1 49       "
	K2      = device.K2      // 2 50       "
	K3      = device.K3      // 3 51       "
	K4      = device.K4      // 4 52       "
	K5      = device.K5      // 5 53       "
	K6      = device.K6      // 6 54       "
	K7      = device.K7      // 7 55       "
	K8      = device.K8      // 8 56       "
	K9      = device.K9      // 9 57       "
	KA      = device.KA      // A 65     Standard keyboard letters.
	KB      = device.KB      // B 66       "
	KC      = device.KC      // C 67       "
	KD      = device.KD      // D 68       "
	KE      = device.KE      // E 69       "
	KF      = device.KF      // F 70       "
	KG      = device.KG      // G 71       "
	KH      = device.KH      // H 72       "
	KI      = device.KI      // I 73       "
	KJ      = device.KJ      // J 74       "
	KK      = device.KK      // K 75       "
	KL      = device.KL      // L 76       "
	KM      = device.KM      // M 77       "
	KN      = device.KN      // N 78       "
	KO      = device.KO      // O 79       "
	KP      = device.KP      // P 80       "
	KQ      = device.KQ      // Q 81       "
	KR      = device.KR      // R 82       "
	KS      = device.KS      // S 83       "
	KT      = device.KT      // T 84       "
	KU      = device.KU      // U 85       "
	KV      = device.KV      // V 86       "
	KW      = device.KW      // W 87       "
	KX      = device.KX      // X 88       "
	KY      = device.KY      // Y 89       "
	KZ      = device.KZ      // Z 90       "
	KEqual  = device.KEqual  // = 61     Standard keyboard punctuation keys.
	KMinus  = device.KMinus  // - 45       "
	KRBkt   = device.KRBkt   // ] 93       "
	KLBkt   = device.KLBkt   // [ 91       "
	KQuote  = device.KQuote  // " 34       "
	KSemi   = device.KSemi   // ; 59       "
	KBSl    = device.KBSl    // \ 92       "
	KComma  = device.KComma  // , 44       "
	KSlash  = device.KSlash  // / 47       "
	KDot    = device.KDot    // . 46       "
	KGrave  = device.KGrave  // ~ 126      "
	KRet    = device.KRet    // ⇦ 8678     "
	KTab    = device.KTab    // ⇨ 8680     "
	KSpace  = device.KSpace  // ▭ 9645     "
	KDel    = device.KDel    // ⇍ 8653     "
	KEsc    = device.KEsc    // ⊶ 8886     "
	KF1     = device.KF1     // α 945    General Function keys.
	KF2     = device.KF2     // β 946      "
	KF3     = device.KF3     // γ 947      "
	KF4     = device.KF4     // δ 948      "
	KF5     = device.KF5     // ε 949      "
	KF6     = device.KF6     // ζ 950      "
	KF7     = device.KF7     // η 951      "
	KF8     = device.KF8     // θ 952      "
	KF9     = device.KF9     // ι 953      "
	KF10    = device.KF10    // κ 954      "
	KF11    = device.KF11    // λ 955      "
	KF12    = device.KF12    // μ 956      "
	KF13    = device.KF13    // ν 957      "
	KF14    = device.KF14    // ξ 958      "
	KF15    = device.KF15    // ο 959      "
	KF16    = device.KF16    // π 960      "
	KF17    = device.KF17    // ρ 961      "
	KF18    = device.KF18    // ς 962      "
	KF19    = device.KF19    // σ 963      "
	KHome   = device.KHome   // ◈ 9672   Specific function keys.
	KPgUp   = device.KPgUp   // ⇑ 8657     "
	KFDel   = device.KFDel   // ⇏ 8655     "
	KEnd    = device.KEnd    // ▣ 9635     "
	KPgDn   = device.KPgDn   // ⇓ 8659     "
	KALeft  = device.KALeft  // ◀ 9664   Arrow keys
	KARight = device.KARight // ▶ 9654     "
	KADown  = device.KADown  // ▼ 9660     "
	KAUp    = device.KAUp    // ▲ 9650     "
	KPDot   = device.KPDot   // ⊙ 8857   Extended keyboard keypad keys
	KPMlt   = device.KPMlt   // ⊗ 8855     "
	KPAdd   = device.KPAdd   // ⊕ 8853     "
	KPClr   = device.KPClr   // ⊠ 8864     "
	KPDiv   = device.KPDiv   // ⊘ 8856     "
	KPEnt   = device.KPEnt   // ⇐ 8656     "
	KPSub   = device.KPSub   // ⊖ 8854     "
	KPEql   = device.KPEql   // ⊜ 8860     "
	KP0     = device.KP0     // ₀ 8320     "
	KP1     = device.KP1     // ₁ 8321     "
	KP2     = device.KP2     // ₂ 8322     "
	KP3     = device.KP3     // ₃ 8323     "
	KP4     = device.KP4     // ₄ 8324     "
	KP5     = device.KP5     // ₅ 8325     "
	KP6     = device.KP6     // ₆ 8326     "
	KP7     = device.KP7     // ₇ 8327     "
	KP8     = device.KP8     // ₈ 8328     "
	KP9     = device.KP9     // ₉ 8329     "
	KML     = device.KML     // ◐ 9680   Mouse buttons treated like keys.
	KMM     = device.KMM     // ◓ 9683     "
	KMR     = device.KMR     // ◑ 9681     "
	KCtl    = device.KCtl    // ● 9679   Modifier keys.
	KFn     = device.KFn     // ◍ 9677     "
	KShift  = device.KShift  // ⇧ 8679     "
	KCmd    = device.KCmd    // ◆ 9670     "
	KAlt    = device.KAlt    // ◇ 9671     "
)
