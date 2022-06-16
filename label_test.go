// Copyright © 2024 Galvanized Logic Inc.

package vu

import (
	"testing"

	"github.com/gazed/vu/load"
)

// go test -run Label
func TestLabel(t *testing.T) {
	data, err := load.Font("lucidiaSu18.fnt")
	if err != nil {
		t.Fatalf("unexpected font load error: %s", err)
	}
	fnt := newFont("lucidiaSu18")
	fnt.w, fnt.h = data.W, data.H
	for _, ch := range data.Chars {
		fnt.addChar(ch.Char, ch.X, ch.Y, ch.W, ch.H, ch.Xo, ch.Yo, ch.Xa)
	}

	sx, sy, md := fnt.setStr("X", 0)
	if sx != 10 || sy != 18 {
		t.Errorf("expected size 10:18 got %d:%d", sx, sy)
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

	// dump the label mesh data when debugging.
	// md[load.Vertexes].PrintF32()
	// md[load.Texcoords].PrintF32()
	// md[load.Indexes].PrintU16()
}
