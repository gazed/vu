// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package load

import (
	"testing"
)

func TestLoadFnt(t *testing.T) {
	load := newLoader().setDir(src, "../eg/source")
	f, _ := load.fnt("CordiaNew")
	if f == nil {
		t.Fatal("Could not load glyphs")
	}
	if f.W != 256 || f.H != 256 || len(f.Chars) != 94 {
		t.Errorf("Invalid font data: %d %d %d", f.W, f.H, len(f.Chars))
	}
}
