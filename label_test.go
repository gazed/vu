// Copyright © 2024 Galvanized Logic Inc.

package vu

import (
	"testing"

	"github.com/gazed/vu/load"
)

// go test -run Label
func TestLabel(t *testing.T) {
	atlas, err := load.TTFont("18:lucon.ttf")
	if err != nil {
		t.Fatalf("unexpected font load error: %s", err)
	}
	f := newFont("lucon18")
	f.w, f.h = int(atlas.Img.Width), int(atlas.Img.Height)
	for _, g := range atlas.Glyphs {
		f.addChar(g.Char, g.X, g.Y, g.W, g.H, g.Xo, g.Yo, g.Xa)
	}
	sx, sy, md := f.setStr("X", 0)
	if sx != 13 || sy != 18 {
		t.Errorf("expected size 13:18 got %d:%d", sx, sy)
	}
	vb := md[load.Vertexes]
	if len(vb.Data) != 8*4 && vb.Count != 4 {
		t.Errorf("vertexes expected 32 bytes and 4 vec2 got %d bytes %d vec2", len(vb.Data), vb.Count)
	}
	tb := md[load.Texcoords]
	if len(tb.Data) != 8*4 && tb.Count != 4 {
		t.Errorf("texcoords expected 32 bytes and 4 vec2 got %d bytes %d vec2", len(tb.Data), tb.Count)
	}
	ib := md[load.Indexes]
	if len(ib.Data) != 6*2 && ib.Count != 6 {
		t.Errorf("indexes expected 12 bytes and 6 indexes got %d bytes %d vec2", len(ib.Data), ib.Count)
	}

	// DEBUG: dump the label mesh data.
	// md[load.Vertexes].PrintF32()
	// md[load.Texcoords].PrintF32()
	// md[load.Indexes].PrintU16()
}
