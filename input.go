// Copyright © 2014-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// input.go wraps the device layer for the engine applications.
// FUTURE: Have the device layer also support text entry.

import (
	"github.com/gazed/vu/device"
)

// Input is used to communicate user feedback to the application.
// User feedback is the current cursor location, current pressed keys,
// mouse buttons, and modifiers. These are sent to the application
// each App.Update() callback.  Overall the keyboard is treated like
// a gamepad controller. Keys and buttons are pressed or not pressed.
//
// The map of keys and mouse buttons that are currently pressed also
// include how long they have been pressed in update ticks. A negative
// value indicates a key release, upon which the total down duration can
// be calculated using the down duration less the KeyReleased timestamp.
type Input struct {
	Mx, My  int         // Current mouse location.
	Down    map[int]int // Keys, buttons with down duration ticks.
	Focus   bool        // True if window is in focus.
	Resized bool        // True if window was resized or moved.
	Scroll  int         // Scroll amount: plus, minus or zero.
	Dt      float64     // Delta time for this update tick.
	Ut      uint64      // Total number of update ticks.
}

// convertInput copies the given device.Pressed input into vu.Input.
// It also adds the delta time and updates the current game time
// in update ticks. It is expected to be called each update.
func (in *Input) convertInput(pressed *device.Pressed, ut uint64, dt float64) {
	in.Mx, in.My = pressed.Mx, pressed.My
	in.Focus = pressed.Focus
	in.Resized = pressed.Resized
	in.Scroll = pressed.Scroll
	in.Dt = dt
	in.Ut = ut

	// Create a key/mouse down map that the application can trash
	// since it is cleared and refilled each update.
	for key := range in.Down {
		delete(in.Down, key)
	}
	for key, val := range pressed.Down {
		in.Down[key] = val
	}
}

// Expose the device package keys as a convenience so the
// device package does not always need including.
// The symbol associated to each key is shown in the comments.
const (
	K0     = device.K0     // 0 48     Standard keyboard numbers.
	K1     = device.K1     // 1 49       "
	K2     = device.K2     // 2 50       "
	K3     = device.K3     // 3 51       "
	K4     = device.K4     // 4 52       "
	K5     = device.K5     // 5 53       "
	K6     = device.K6     // 6 54       "
	K7     = device.K7     // 7 55       "
	K8     = device.K8     // 8 56       "
	K9     = device.K9     // 9 57       "
	KA     = device.KA     // A 65     Standard keyboard letters.
	KB     = device.KB     // B 66       "
	KC     = device.KC     // C 67       "
	KD     = device.KD     // D 68       "
	KE     = device.KE     // E 69       "
	KF     = device.KF     // F 70       "
	KG     = device.KG     // G 71       "
	KH     = device.KH     // H 72       "
	KI     = device.KI     // I 73       "
	KJ     = device.KJ     // J 74       "
	KK     = device.KK     // K 75       "
	KL     = device.KL     // L 76       "
	KM     = device.KM     // M 77       "
	KN     = device.KN     // N 78       "
	KO     = device.KO     // O 79       "
	KP     = device.KP     // P 80       "
	KQ     = device.KQ     // Q 81       "
	KR     = device.KR     // R 82       "
	KS     = device.KS     // S 83       "
	KT     = device.KT     // T 84       "
	KU     = device.KU     // U 85       "
	KV     = device.KV     // V 86       "
	KW     = device.KW     // W 87       "
	KX     = device.KX     // X 88       "
	KY     = device.KY     // Y 89       "
	KZ     = device.KZ     // Z 90       "
	KEqual = device.KEqual // = 61     Standard keyboard punctuation keys.
	KMinus = device.KMinus // - 45       "
	KRBkt  = device.KRBkt  // ] 93       "
	KLBkt  = device.KLBkt  // [ 91       "
	KQt    = device.KQt    // " 34       "
	KSemi  = device.KSemi  // ; 59       "
	KBSl   = device.KBSl   // \ 92       "
	KComma = device.KComma // , 44       "
	KSlash = device.KSlash // / 47       "
	KDot   = device.KDot   // . 46       "
	KGrave = device.KGrave // ~ 126      "
	KRet   = device.KRet   // ⇦ 8678     "
	KTab   = device.KTab   // ⇨ 8680     "
	KSpace = device.KSpace // ▭ 9645     "
	KDel   = device.KDel   // ⇍ 8653     "
	KEsc   = device.KEsc   // ⊶ 8886     "
	KF1    = device.KF1    // α 945    General Function keys.
	KF2    = device.KF2    // β 946      "
	KF3    = device.KF3    // γ 947      "
	KF4    = device.KF4    // δ 948      "
	KF5    = device.KF5    // ε 949      "
	KF6    = device.KF6    // ζ 950      "
	KF7    = device.KF7    // η 951      "
	KF8    = device.KF8    // θ 952      "
	KF9    = device.KF9    // ι 953      "
	KF10   = device.KF10   // κ 954      "
	KF11   = device.KF11   // λ 955      "
	KF12   = device.KF12   // μ 956      "
	KF13   = device.KF13   // ν 957      "
	KF14   = device.KF14   // ξ 958      "
	KF15   = device.KF15   // ο 959      "
	KF16   = device.KF16   // π 960      "
	KF17   = device.KF17   // ρ 961      "
	KF18   = device.KF18   // ς 962      "
	KF19   = device.KF19   // σ 963      "
	KHome  = device.KHome  // ◈ 9672   Specific function keys.
	KPgUp  = device.KPgUp  // ⇑ 8657     "
	KFDel  = device.KFDel  // ⇏ 8655     "
	KEnd   = device.KEnd   // ▣ 9635     "
	KPgDn  = device.KPgDn  // ⇓ 8659     "
	KLa    = device.KLa    // ◀ 9664   Arrow keys
	KRa    = device.KRa    // ▶ 9654     "
	KDa    = device.KDa    // ▼ 9660     "
	KUa    = device.KUa    // ▲ 9650     "
	KKpDot = device.KKpDot // ⊙ 8857   Extended keyboard keypad keys
	KKpMlt = device.KKpMlt // ⊗ 8855     "
	KKpAdd = device.KKpAdd // ⊕ 8853     "
	KKpClr = device.KKpClr // ⊠ 8864     "
	KKpDiv = device.KKpDiv // ⊘ 8856     "
	KKpEnt = device.KKpEnt // ⇐ 8656     "
	KKpSub = device.KKpSub // ⊖ 8854     "
	KKpEql = device.KKpEql // ⊜ 8860     "
	KKp0   = device.KKp0   // ₀ 8320     "
	KKp1   = device.KKp1   // ₁ 8321     "
	KKp2   = device.KKp2   // ₂ 8322     "
	KKp3   = device.KKp3   // ₃ 8323     "
	KKp4   = device.KKp4   // ₄ 8324     "
	KKp5   = device.KKp5   // ₅ 8325     "
	KKp6   = device.KKp6   // ₆ 8326     "
	KKp7   = device.KKp7   // ₇ 8327     "
	KKp8   = device.KKp8   // ₈ 8328     "
	KKp9   = device.KKp9   // ₉ 8329     "
	KLm    = device.KLm    // ◐ 9680   Mouse buttons treated like keys.
	KMm    = device.KMm    // ◓ 9683     "
	KRm    = device.KRm    // ◑ 9681     "
	KCtl   = device.KCtl   // ● 9679   Modifier keys.
	KFn    = device.KFn    // ◍ 9677     "
	KShift = device.KShift // ⇧ 8679     "
	KCmd   = device.KCmd   // ◆ 9670     "
	KAlt   = device.KAlt   // ◇ 9671     "
)

// Keysym returns a single rune representing the given key.
// Zero is returned if there is no rune for the key. This is intended
// to provide a default means of representing each keyboard key with
// a displayable character.
func Keysym(keycode int) rune {
	if symbol, ok := keysym[keycode]; ok {
		return rune(symbol)
	}
	return 0
}

// keysym maps key codes to unicode runes.
// Ensure that font has a character for each of the runes below.
// The symbols are also shown in constant comments so they appear in the godoc.
var keysym = map[int]int{
	K0:     0x0030, // 0 48
	K1:     0x0031, // 1 49
	K2:     0x0032, // 2 50
	K3:     0x0033, // 3 51
	K4:     0x0034, // 4 52
	K5:     0x0035, // 5 53
	K6:     0x0036, // 6 54
	K7:     0x0037, // 7 55
	K8:     0x0038, // 8 56
	K9:     0x0039, // 9 57
	KA:     0x0041, // A 65
	KB:     0x0042, // B 66
	KC:     0x0043, // C 67
	KD:     0x0044, // D 68
	KE:     0x0045, // E 69
	KF:     0x0046, // F 70
	KG:     0x0047, // G 71
	KH:     0x0048, // H 72
	KI:     0x0049, // I 73
	KJ:     0x004A, // J 74
	KK:     0x004B, // K 75
	KL:     0x004C, // L 76
	KM:     0x004D, // M 77
	KN:     0x004E, // N 78
	KO:     0x004F, // O 79
	KP:     0x0050, // P 80
	KQ:     0x0051, // Q 81
	KR:     0x0052, // R 82
	KS:     0x0053, // S 83
	KT:     0x0054, // T 84
	KU:     0x0055, // U 85
	KV:     0x0056, // V 86
	KW:     0x0057, // W 87
	KX:     0x0058, // X 88
	KY:     0x0059, // Y 89
	KZ:     0x005A, // Z 90
	KEqual: 0x003D, // = 61
	KMinus: 0x002D, // - 45
	KRBkt:  0x005D, // ] 93
	KLBkt:  0x005B, // [ 91
	KQt:    0x0022, // " 34
	KSemi:  0x003B, // ; 59
	KBSl:   0x005C, // \ 92
	KComma: 0x002C, // , 44
	KSlash: 0x002F, // / 47
	KDot:   0x002E, // . 46
	KGrave: 0x007E, // ~ 126
	KRet:   0x21E6, // ⇦ 8678
	KTab:   0x21E8, // ⇨ 8680
	KSpace: 0x25AD, // ▭ 9645
	KDel:   0x21CD, // ⇍ 8653
	KEsc:   0x22B6, // ⊶ 8886
	KF1:    0x03B1, // α 945
	KF2:    0x03B2, // β 946
	KF3:    0x03B3, // γ 947
	KF4:    0x03B4, // δ 948
	KF5:    0x03B5, // ε 949
	KF6:    0x03B6, // ζ 950
	KF7:    0x03B7, // η 951
	KF8:    0x03B8, // θ 952
	KF9:    0x03B9, // ι 953
	KF10:   0x03BA, // κ 954
	KF11:   0x03BB, // λ 955
	KF12:   0x03BC, // μ 956
	KF13:   0x03BD, // ν 957
	KF14:   0x03BE, // ξ 958
	KF15:   0x03BF, // ο 959
	KF16:   0x03C0, // π 960
	KF17:   0x03C1, // ρ 961
	KF18:   0x03C2, // ς 962
	KF19:   0x03C3, // σ 963
	KHome:  0x25C8, // ◈ 9672
	KPgUp:  0x21D1, // ⇑ 8657
	KFDel:  0x21CF, // ⇏ 8655
	KEnd:   0x25A3, // ▣ 9635
	KPgDn:  0x21D3, // ⇓ 8659
	KLa:    0x25C0, // ◀ 9664
	KRa:    0x25B6, // ▶ 9654
	KDa:    0x25BC, // ▼ 9660
	KUa:    0x25B2, // ▲ 9650
	KKpDot: 0x2299, // ⊙ 8857
	KKpMlt: 0x2297, // ⊗ 8855
	KKpAdd: 0x2295, // ⊕ 8853
	KKpClr: 0x22A0, // ⊠ 8864
	KKpDiv: 0x2298, // ⊘ 8856
	KKpEnt: 0x21D0, // ⇐ 8656
	KKpSub: 0x2296, // ⊖ 8854
	KKpEql: 0x229C, // ⊜ 8860
	KKp0:   0x2080, // ₀ 8320
	KKp1:   0x2081, // ₁ 8321
	KKp2:   0x2082, // ₂ 8322
	KKp3:   0x2083, // ₃ 8323
	KKp4:   0x2084, // ₄ 8324
	KKp5:   0x2085, // ₅ 8325
	KKp6:   0x2086, // ₆ 8326
	KKp7:   0x2087, // ₇ 8327
	KKp8:   0x2088, // ₈ 8328
	KKp9:   0x2089, // ₉ 8329
	KLm:    0x25D0, // ◐ 9680
	KMm:    0x25D3, // ◓ 9683
	KRm:    0x25D1, // ◑ 9681
	KCtl:   0x25CF, // ● 9679
	KFn:    0x25CD, // ◍ 9677
	KShift: 0x21E7, // ⇧ 8679
	KCmd:   0x25C6, // ◆ 9670
	KAlt:   0x25C7, // ◇ 9671

}

// Runes available in lucidiaSu.
//
// U+0020 ' ' 32     U+0041 'A' 65     U+0061 'a' 97     U+03B2 'β' 946
// U+0021 '!' 33     U+0042 'B' 66     U+0062 'b' 98     U+03B3 'γ' 947
// U+0022 '"' 34     U+0043 'C' 67     U+0063 'c' 99     U+03B4 'δ' 948
// U+0023 '#' 35     U+0044 'D' 68     U+0064 'd' 10     U+03B5 'ε' 9490
// U+0024 '$' 36     U+0045 'E' 69     U+0065 'e' 10     U+03B6 'ζ' 9501
// U+0025 '%' 37     U+0046 'F' 70     U+0066 'f' 10     U+03B7 'η' 9512
// U+0026 '&' 38     U+0047 'G' 71     U+0067 'g' 10     U+03B8 'θ' 9523
// U+0027 ''' 39     U+0048 'H' 72     U+0068 'h' 10     U+03B9 'ι' 9534
// U+0028 '(' 40     U+0049 'I' 73     U+0069 'i' 10     U+03BA 'κ' 9545
// U+0029 ')' 41     U+004A 'J' 74     U+006A 'j' 10     U+03BB 'λ' 9556
// U+002A '*' 42     U+004B 'K' 75     U+006B 'k' 10     U+03BC 'μ' 9567
// U+002B '+' 43     U+004C 'L' 76     U+006C 'l' 10     U+03BD 'ν' 9578
// U+002C ',' 44     U+004D 'M' 77     U+006D 'm' 10     U+03BE 'ξ' 9589
// U+002D '-' 45     U+004E 'N' 78     U+006E 'n' 11     U+03BF 'ο' 9590
// U+002E '.' 46     U+004F 'O' 79     U+006F 'o' 11     U+03C0 'π' 9601
// U+002F '/' 47     U+0050 'P' 80     U+0070 'p' 11     U+03C1 'ρ' 9612
// U+0030 '0' 48     U+0051 'Q' 81     U+0071 'q' 11     U+03C2 'ς' 9623
// U+0031 '1' 49     U+0052 'R' 82     U+0072 'r' 11     U+03C3 'σ' 9634
// U+0032 '2' 50     U+0053 'S' 83     U+0073 's' 11     U+03C4 'τ' 9645
// U+0033 '3' 51     U+0054 'T' 84     U+0074 't' 11     U+03C5 'υ' 9656
// U+0034 '4' 52     U+0055 'U' 85     U+0075 'u' 11     U+03C6 'φ' 9667
// U+0035 '5' 53     U+0056 'V' 86     U+0076 'v' 11     U+03C7 'χ' 9678
// U+0036 '6' 54     U+0057 'W' 87     U+0077 'w' 11     U+03C8 'ψ' 9689
// U+0037 '7' 55     U+0058 'X' 88     U+0078 'x' 12     U+03C9 'ω' 9690
// U+0038 '8' 56     U+0059 'Y' 89     U+0079 'y' 12     U+2080 '₀' 83201
// U+0039 '9' 57     U+005A 'Z' 90     U+007A 'z' 12     U+2081 '₁' 83212
// U+003A ':' 58     U+005B '[' 91     U+007B '{' 12     U+2082 '₂' 83223
// U+003B ';' 59     U+005C '\' 92     U+007C '|' 12     U+2083 '₃' 83234
// U+003C '<' 60     U+005D ']' 93     U+007D '}' 12     U+2084 '₄' 83245
// U+003D '=' 61     U+005E '^' 94     U+007E '~' 12     U+2085 '₅' 83256
// U+003E '>' 62     U+005F '_' 95                       U+2086 '₆' 8326
// U+003F '?' 63     U+0060 '`' 96                       U+2087 '₇' 8327
// U+0040 '@' 64                                         U+2088 '₈' 8328
//                                                       U+2089 '₉' 8329
//
// U+2190 '←' 8592   U+2295 '⊕' 8853   U+25AA '▪' 9642   U+25C8 '◈' 9672
// U+2191 '↑' 8593   U+2296 '⊖' 8854   U+25AB '▫' 9643   U+25C9 '◉' 9673
// U+2192 '→' 8594   U+2297 '⊗' 8855   U+25AC '▬' 9644   U+25CA '◊' 9674
// U+2193 '↓' 8595   U+2298 '⊘' 8856   U+25AD '▭' 9645   U+25CB '○' 9675
// U+21CD '⇍' 8653   U+2299 '⊙' 8857   U+25AE '▮' 9646   U+25CC '◌' 9676
// U+21CF '⇏' 8655   U+229A '⊚' 8858   U+25AF '▯' 9647   U+25CD '◍' 9677
// U+21D0 '⇐' 8656   U+229B '⊛' 8859   U+25B0 '▰' 9648   U+25CE '◎' 9678
// U+21D1 '⇑' 8657   U+229C '⊜' 8860   U+25B1 '▱' 9649   U+25CF '●' 9679
// U+21D2 '⇒' 8658   U+229D '⊝' 8861   U+25B2 '▲' 9650   U+25D0 '◐' 9680
// U+21D3 '⇓' 8659   U+229E '⊞' 8862   U+25B3 '△' 9651   U+25D1 '◑' 9681
// U+21DE '⇞' 8670   U+229F '⊟' 8863   U+25B4 '▴' 9652   U+25D2 '◒' 9682
// U+21DF '⇟' 8671   U+22A0 '⊠' 8864   U+25B5 '▵' 9653   U+25D3 '◓' 9683
// U+21E0 '⇠' 8672   U+22A1 '⊡' 8865   U+25B6 '▶' 9654   U+25D4 '◔' 9684
// U+21E1 '⇡' 8673   U+22B6 '⊶' 8886   U+25B7 '▷' 9655   U+25D5 '◕' 9685
// U+21E2 '⇢' 8674   U+22B7 '⊷' 8887   U+25B8 '▸' 9656   U+25D6 '◖' 9686
// U+21E3 '⇣' 8675   U+2408 '␈' 9224   U+25B9 '▹' 9657   U+25D7 '◗' 9687
// U+21E4 '⇤' 8676   U+240D '␍' 9229   U+25BA '►' 9658   U+25E2 '◢' 9698
// U+21E5 '⇥' 8677   U+241B '␛' 9243   U+25BB '▻' 9659   U+25E3 '◣' 9699
// U+21E6 '⇦' 8678   U+2420 '␠' 9248   U+25BC '▼' 9660   U+25E4 '◤' 9700
// U+21E7 '⇧' 8679   U+2423 '␣' 9251   U+25BD '▽' 9661   U+25E5 '◥' 9701
// U+21E8 '⇨' 8680   U+25A0 '■' 9632   U+25BE '▾' 9662   U+25E6 '◦' 9702
// U+21E9 '⇩' 8681   U+25A1 '□' 9633   U+25BF '▿' 9663   U+25E7 '◧' 9703
// U+21EA '⇪' 8682   U+25A2 '▢' 9634   U+25C0 '◀' 9664   U+25E8 '◨' 9704
// U+2218 '∘' 8728   U+25A3 '▣' 9635   U+25C1 '◁' 9665   U+25E9 '◩' 9705
// U+2219 '∙' 8729   U+25A4 '▤' 9636   U+25C2 '◂' 9666   U+25EA '◪' 9706
// U+2257 '≗' 8791   U+25A5 '▥' 9637   U+25C3 '◃' 9667   U+25EB '◫' 9707
// U+225B '≛' 8795   U+25A6 '▦' 9638   U+25C4 '◄' 9668   U+25EC '◬' 9708
// U+225C '≜' 8796   U+25A7 '▧' 9639   U+25C5 '◅' 9669   U+25ED '◭' 9709
//                   U+25A8 '▨' 9640   U+25C6 '◆' 9670   U+25EE '◮' 9710
//                   U+25A9 '▩' 9641   U+25C7 '◇' 9671
