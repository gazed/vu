// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"github.com/gazed/vu/device"
)

// Input is used to communicate current user input to the application.
// This gives the current cursor location, current pressed keys,
// mouse buttons, and modifiers. These are sent to the application
// each App.Update() callback.
//
// The map of keys and mouse buttons that are currently pressed also
// include how long they have been pressed in update ticks. A negative
// value indicates a release. The total down duration can be calculated
// on release using down duration less RELEASED timestamp.
type Input struct {
	Mx, My  int         // Current mouse location.
	Down    map[int]int // Keys, buttons with down duration ticks.
	Focus   bool        // True if window is in focus.
	Resized bool        // True if window was resized or moved.
	Scroll  int         // Scroll amount, if any.
	Dt      float64     // Delta time for this update.
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

	// Create a key/mouse down map that the application can trash.
	// It is expected to be cleared and refilled each update.
	for key, _ := range in.Down {
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
	K_0     = device.K_0     // '0' 48     Standard keyboard numbers.
	K_1     = device.K_1     // '1' 49       "
	K_2     = device.K_2     // '2' 50       "
	K_3     = device.K_3     // '3' 51       "
	K_4     = device.K_4     // '4' 52       "
	K_5     = device.K_5     // '5' 53       "
	K_6     = device.K_6     // '6' 54       "
	K_7     = device.K_7     // '7' 55       "
	K_8     = device.K_8     // '8' 56       "
	K_9     = device.K_9     // '9' 57       "
	K_A     = device.K_A     // 'A' 65     Standard keyboard letters.
	K_B     = device.K_B     // 'B' 66       "
	K_C     = device.K_C     // 'C' 67       "
	K_D     = device.K_D     // 'D' 68       "
	K_E     = device.K_E     // 'E' 69       "
	K_F     = device.K_F     // 'F' 70       "
	K_G     = device.K_G     // 'G' 71       "
	K_H     = device.K_H     // 'H' 72       "
	K_I     = device.K_I     // 'I' 73       "
	K_J     = device.K_J     // 'J' 74       "
	K_K     = device.K_K     // 'K' 75       "
	K_L     = device.K_L     // 'L' 76       "
	K_M     = device.K_M     // 'M' 77       "
	K_N     = device.K_N     // 'N' 78       "
	K_O     = device.K_O     // 'O' 79       "
	K_P     = device.K_P     // 'P' 80       "
	K_Q     = device.K_Q     // 'Q' 81       "
	K_R     = device.K_R     // 'R' 82       "
	K_S     = device.K_S     // 'S' 83       "
	K_T     = device.K_T     // 'T' 84       "
	K_U     = device.K_U     // 'U' 85       "
	K_V     = device.K_V     // 'V' 86       "
	K_W     = device.K_W     // 'W' 87       "
	K_X     = device.K_X     // 'X' 88       "
	K_Y     = device.K_Y     // 'Y' 89       "
	K_Z     = device.K_Z     // 'Z' 90       "
	K_Equal = device.K_Equal // '=' 61     Standard keyboard punctuation keys.
	K_Minus = device.K_Minus // '-' 45       "
	K_RBkt  = device.K_RBkt  // ']' 93       "
	K_LBkt  = device.K_LBkt  // '[' 91       "
	K_Qt    = device.K_Qt    // '"' 34       "
	K_Semi  = device.K_Semi  // ';' 59       "
	K_BSl   = device.K_BSl   // '\' 92       "
	K_Comma = device.K_Comma // ',' 44       "
	K_Slash = device.K_Slash // '/' 47       "
	K_Dot   = device.K_Dot   // '.' 46       "
	K_Grave = device.K_Grave // '~' 126      "
	K_Ret   = device.K_Ret   // '⇦' 8678     "
	K_Tab   = device.K_Tab   // '⇨' 8680     "
	K_Space = device.K_Space // '▭' 9645     "
	K_Del   = device.K_Del   // '⇍' 8653     "
	K_Esc   = device.K_Esc   // '⊶' 8886     "
	K_F1    = device.K_F1    // 'α' 945    General Function keys.
	K_F2    = device.K_F2    // 'β' 946      "
	K_F3    = device.K_F3    // 'γ' 947      "
	K_F4    = device.K_F4    // 'δ' 948      "
	K_F5    = device.K_F5    // 'ε' 949      "
	K_F6    = device.K_F6    // 'ζ' 950      "
	K_F7    = device.K_F7    // 'η' 951      "
	K_F8    = device.K_F8    // 'θ' 952      "
	K_F9    = device.K_F9    // 'ι' 953      "
	K_F10   = device.K_F10   // 'κ' 954      "
	K_F11   = device.K_F11   // 'λ' 955      "
	K_F12   = device.K_F12   // 'μ' 956      "
	K_F13   = device.K_F13   // 'ν' 957      "
	K_F14   = device.K_F14   // 'ξ' 958      "
	K_F15   = device.K_F15   // 'ο' 959      "
	K_F16   = device.K_F16   // 'π' 960      "
	K_F17   = device.K_F17   // 'ρ' 961      "
	K_F18   = device.K_F18   // 'ς' 962      "
	K_F19   = device.K_F19   // 'σ' 963      "
	K_Home  = device.K_Home  // '◈' 9672   Specific function keys.
	K_PgUp  = device.K_PgUp  // '⇑' 8657     "
	K_FDel  = device.K_FDel  // '⇏' 8655     "
	K_End   = device.K_End   // '▣' 9635     "
	K_PgDn  = device.K_PgDn  // '⇓' 8659     "
	K_La    = device.K_La    // '◀' 9664   Arrow keys
	K_Ra    = device.K_Ra    // '▶' 9654     "
	K_Da    = device.K_Da    // '▼' 9660     "
	K_Ua    = device.K_Ua    // '▲' 9650     "
	K_KpDot = device.K_KpDot // '⊙' 8857   Extended keyboard keypad keys
	K_KpMlt = device.K_KpMlt // '⊗' 8855     "
	K_KpAdd = device.K_KpAdd // '⊕' 8853     "
	K_KpClr = device.K_KpClr // '⊠' 8864     "
	K_KpDiv = device.K_KpDiv // '⊘' 8856     "
	K_KpEnt = device.K_KpEnt // '⇐' 8656     "
	K_KpSub = device.K_KpSub // '⊖' 8854     "
	K_KpEql = device.K_KpEql // '⊜' 8860     "
	K_Kp0   = device.K_Kp0   // '₀' 8320     "
	K_Kp1   = device.K_Kp1   // '₁' 8321     "
	K_Kp2   = device.K_Kp2   // '₂' 8322     "
	K_Kp3   = device.K_Kp3   // '₃' 8323     "
	K_Kp4   = device.K_Kp4   // '₄' 8324     "
	K_Kp5   = device.K_Kp5   // '₅' 8325     "
	K_Kp6   = device.K_Kp6   // '₆' 8326     "
	K_Kp7   = device.K_Kp7   // '₇' 8327     "
	K_Kp8   = device.K_Kp8   // '₈' 8328     "
	K_Kp9   = device.K_Kp9   // '₉' 8329     "
	K_Lm    = device.K_Lm    // '◐' 9680   Mouse buttons treated like keys.
	K_Mm    = device.K_Mm    // '◓' 9683     "
	K_Rm    = device.K_Rm    // '◑' 9681     "
	K_Ctl   = device.K_Ctl   // '●' 9679   Modifier keys.
	K_Fn    = device.K_Fn    // '◍' 9677     "
	K_Shift = device.K_Shift // '⇧' 8679     "
	K_Cmd   = device.K_Cmd   // '◆' 9670     "
	K_Alt   = device.K_Alt   // '◇' 9671     "
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
var keysym map[int]int = map[int]int{
	K_0:     0x0030, // '0' 48
	K_1:     0x0031, // '1' 49
	K_2:     0x0032, // '2' 50
	K_3:     0x0033, // '3' 51
	K_4:     0x0034, // '4' 52
	K_5:     0x0035, // '5' 53
	K_6:     0x0036, // '6' 54
	K_7:     0x0037, // '7' 55
	K_8:     0x0038, // '8' 56
	K_9:     0x0039, // '9' 57
	K_A:     0x0041, // 'A' 65
	K_B:     0x0042, // 'B' 66
	K_C:     0x0043, // 'C' 67
	K_D:     0x0044, // 'D' 68
	K_E:     0x0045, // 'E' 69
	K_F:     0x0046, // 'F' 70
	K_G:     0x0047, // 'G' 71
	K_H:     0x0048, // 'H' 72
	K_I:     0x0049, // 'I' 73
	K_J:     0x004A, // 'J' 74
	K_K:     0x004B, // 'K' 75
	K_L:     0x004C, // 'L' 76
	K_M:     0x004D, // 'M' 77
	K_N:     0x004E, // 'N' 78
	K_O:     0x004F, // 'O' 79
	K_P:     0x0050, // 'P' 80
	K_Q:     0x0051, // 'Q' 81
	K_R:     0x0052, // 'R' 82
	K_S:     0x0053, // 'S' 83
	K_T:     0x0054, // 'T' 84
	K_U:     0x0055, // 'U' 85
	K_V:     0x0056, // 'V' 86
	K_W:     0x0057, // 'W' 87
	K_X:     0x0058, // 'X' 88
	K_Y:     0x0059, // 'Y' 89
	K_Z:     0x005A, // 'Z' 90
	K_Equal: 0x003D, // '=' 61
	K_Minus: 0x002D, // '-' 45
	K_RBkt:  0x005D, // ']' 93
	K_LBkt:  0x005B, // '[' 91
	K_Qt:    0x0022, // '"' 34
	K_Semi:  0x003B, // ';' 59
	K_BSl:   0x005C, // '\' 92
	K_Comma: 0x002C, // ',' 44
	K_Slash: 0x002F, // '/' 47
	K_Dot:   0x002E, // '.' 46
	K_Grave: 0x007E, // '~' 126
	K_Ret:   0x21E6, // '⇦' 8678
	K_Tab:   0x21E8, // '⇨' 8680
	K_Space: 0x25AD, // '▭' 9645
	K_Del:   0x21CD, // '⇍' 8653
	K_Esc:   0x22B6, // '⊶' 8886
	K_F1:    0x03B1, // 'α' 945
	K_F2:    0x03B2, // 'β' 946
	K_F3:    0x03B3, // 'γ' 947
	K_F4:    0x03B4, // 'δ' 948
	K_F5:    0x03B5, // 'ε' 949
	K_F6:    0x03B6, // 'ζ' 950
	K_F7:    0x03B7, // 'η' 951
	K_F8:    0x03B8, // 'θ' 952
	K_F9:    0x03B9, // 'ι' 953
	K_F10:   0x03BA, // 'κ' 954
	K_F11:   0x03BB, // 'λ' 955
	K_F12:   0x03BC, // 'μ' 956
	K_F13:   0x03BD, // 'ν' 957
	K_F14:   0x03BE, // 'ξ' 958
	K_F15:   0x03BF, // 'ο' 959
	K_F16:   0x03C0, // 'π' 960
	K_F17:   0x03C1, // 'ρ' 961
	K_F18:   0x03C2, // 'ς' 962
	K_F19:   0x03C3, // 'σ' 963
	K_Home:  0x25C8, // '◈' 9672
	K_PgUp:  0x21D1, // '⇑' 8657
	K_FDel:  0x21CF, // '⇏' 8655
	K_End:   0x25A3, // '▣' 9635
	K_PgDn:  0x21D3, // '⇓' 8659
	K_La:    0x25C0, // '◀' 9664
	K_Ra:    0x25B6, // '▶' 9654
	K_Da:    0x25BC, // '▼' 9660
	K_Ua:    0x25B2, // '▲' 9650
	K_KpDot: 0x2299, // '⊙' 8857
	K_KpMlt: 0x2297, // '⊗' 8855
	K_KpAdd: 0x2295, // '⊕' 8853
	K_KpClr: 0x22A0, // '⊠' 8864
	K_KpDiv: 0x2298, // '⊘' 8856
	K_KpEnt: 0x21D0, // '⇐' 8656
	K_KpSub: 0x2296, // '⊖' 8854
	K_KpEql: 0x229C, // '⊜' 8860
	K_Kp0:   0x2080, // '₀' 8320
	K_Kp1:   0x2081, // '₁' 8321
	K_Kp2:   0x2082, // '₂' 8322
	K_Kp3:   0x2083, // '₃' 8323
	K_Kp4:   0x2084, // '₄' 8324
	K_Kp5:   0x2085, // '₅' 8325
	K_Kp6:   0x2086, // '₆' 8326
	K_Kp7:   0x2087, // '₇' 8327
	K_Kp8:   0x2088, // '₈' 8328
	K_Kp9:   0x2089, // '₉' 8329
	K_Lm:    0x25D0, // '◐' 9680
	K_Mm:    0x25D3, // '◓' 9683
	K_Rm:    0x25D1, // '◑' 9681
	K_Ctl:   0x25CF, // '●' 9679
	K_Fn:    0x25CD, // '◍' 9677
	K_Shift: 0x21E7, // '⇧' 8679
	K_Cmd:   0x25C6, // '◆' 9670
	K_Alt:   0x25C7, // '◇' 9671

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
