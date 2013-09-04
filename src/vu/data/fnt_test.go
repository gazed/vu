// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package data

import (
	"testing"
)

// Check that font glyphs can be imported.
func TestLoadFnt(t *testing.T) {
	load := &loader{}
	gs, _ := load.fnt("../eg/images", "CordiaNew.fnt")
	if gs == nil {
		t.Fatal("Could not load glyphs")
	}
	if gs.w != 256 || gs.h != 256 || len(gs.glyphs) != 94 {
		t.Errorf("Invalid glyph data: %d %d %d", gs.w, gs.h, len(gs.glyphs))
	}
}
