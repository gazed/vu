// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package load

import (
	"testing"
)

// Uses vu/eg resource directories.
func TestLoadFnt(t *testing.T) {
	f := &FntData{}
	l := NewLocator().Dir("FNT", "../eg/source")
	if err := f.Load("lucidiaSu16", l); err != nil {
		t.Fatal("Could not load glyphs")
	}
	if f.W != 256 || f.H != 256 || len(f.Chars) != 247 {
		t.Errorf("Invalid font data: %d %d %d", f.W, f.H, len(f.Chars))
	}
}
