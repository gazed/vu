// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package lin

// Test combinations of rotations, translations, and scales.
// The standard combination, when there is more than one, is
// first scale, then rotate, then translate.

import "testing"

// Rotating 90 degrees about X (moves points on Z axis to -Y axis)
// and then translate along Z.
func TestRotateTranslate(t *testing.T) {

	// rotate then translate is normally the way to go.  The object ends up
	// rotated at the point it was translated to.
	R, T := QAxisAngle(&V3{X: 1}, 90).M4(), M4Translater(0, 0, 2)
	RT := R.Mult(T)
	vRT := (&V4{0, 0, 0, 1}).MultL(RT)
	got, want := vRT.Dump(), "{0.0 0.0 2.0 1.0}"
	if got != want {
		t.Errorf(format, got, want)
	}

	// Note that this is also RT since TranslateR multiplies on the right.
	RT2 := QAxisAngle(&V3{X: 1}, 90).M4().TranslateR(0, 0, 2)
	rt, rt2 := RT.Dump(), RT2.Dump()
	if rt != rt2 {
		t.Errorf(format, rt, rt2)
	}

	// translate then rotate can be used to rotate an object around an
	// arbitrary axis or point.
	T, R = M4Translater(0, 0, 2), QAxisAngle(&V3{X: 1}, 90).M4()
	TR := T.Mult(R)
	vTR := (&V4{0, 0, 0, 1}).MultL(TR)
	got, want = vTR.Dump(), "{0.0 -2.0 0.0 1.0}"
	if got != want {
		t.Errorf(format, got, want)
	}

	// Note that RT != TR
	tr := TR.Dump()
	if rt == tr {
		t.Errorf(format, rt, tr+"\nShould be different")
	}
}

// Scale the vector v by half and then rotate it 90 degrees about Y
// (moves points on X axis to -Z axis).
func TestRotateScale(t *testing.T) {
	// scale then rotate
	S, R := M4Scaler(0.5, 0, 0), QAxisAngle(&V3{Y: 1}, 90).M4()
	SR := S.Mult(R)
	vSR := (&V4{4, 0, 0, 1}).MultL(SR)
	got, want := vSR.Dump(), "{0.0 0.0 -2.0 1.0}"
	if got != want {
		t.Errorf(format, got, want)
	}

	// Note that this is also SR since ScaleL multiplies on the left.
	SR2 := QAxisAngle(&V3{Y: 1}, 90).M4().ScaleL(0.5, 0, 0)
	sr, sr2 := SR.Dump(), SR2.Dump()
	if sr != sr2 {
		t.Errorf(format, sr, sr2)
	}

	// Note that SR != RS
	R, S = QAxisAngle(&V3{Y: 1}, 90).M4(), M4Scaler(0.5, 0, 0)
	RS := R.Mult(S)
	rs := RS.Dump()
	if sr == rs {
		t.Errorf(format, sr, rs+"\nShould be different")
	}
}

// Scale the vector v, then rotate it, then translate it.
//     v = 4, 0,  0
//     v = 2, 0,  0  : after scaling by 2.
//     v = 0, 0, -2  : after rotating 90 about Y axis (X -> -Z).
//     v = 0, 0,  2  : after translating by 0, 0, 4
func TestScaleRotateTranslate(t *testing.T) {
	S, R, T := M4Scaler(0.5, 0, 0), QAxisAngle(&V3{Y: 1}, 90).M4(), M4Translater(0, 0, 4)
	SRT := S.Mult(R.Mult(T))
	vSRT := (&V4{4, 0, 0, 1}).MultL(SRT)
	got, want := vSRT.Dump(), "{0.0 0.0 2.0 1.0}"
	if got != want {
		t.Errorf(format, got, want)
	}

	// note the faster equivalent
	S, R = M4Scaler(0.5, 0, 0), QAxisAngle(&V3{Y: 1}, 90).M4()
	SRT = S.Mult(R).TranslateR(0, 0, 4)
	vSRT = (&V4{4, 0, 0, 1}).MultL(SRT)
	got, want = vSRT.Dump(), "{0.0 0.0 2.0 1.0}"
	if got != want {
		t.Errorf(format, got, want)
	}

	// note the even faster equivalent
	R = QAxisAngle(&V3{Y: 1}, 90).M4()
	SRT = R.ScaleL(0.5, 0, 0).TranslateR(0, 0, 4)
	vSRT = (&V4{4, 0, 0, 1}).MultL(SRT)
	got, want = vSRT.Dump(), "{0.0 0.0 2.0 1.0}"
	if got != want {
		t.Errorf(format, got, want)
	}
}
