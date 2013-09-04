// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package lin

import "testing"

func TestIsEqual(t *testing.T) {
	var f1 float32 = 0.0
	var f2 float32 = 0.000001
	var f3 float32 = -0.0001
	if IsEqual(f1, f2) && !IsEqual(f1, f3) {
		t.Fail()
	}
}

func TestIsZero(t *testing.T) {
	var f1 float32 = 0.0000001
	var f2 float32 = -0.0000001
	var f3 float32 = -0.0001
	if !IsZero(f1) || !IsZero(f2) || IsZero(f3) {
		t.Fail()
	}
}

func TestIsOne(t *testing.T) {
	var f1 float32 = 1.0000001
	var f2 float32 = -1.0000001
	var f3 float32 = -1.0001
	if !IsOne(f1) || !IsOne(f2) || IsOne(f3) {
		t.Fail()
	}
}

func TestLerp(t *testing.T) {
	if !IsEqual(Lerp(10, 5, 0.5), 7.5) {
		t.Fail()
	}
}
