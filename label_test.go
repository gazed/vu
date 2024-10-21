// Copyright © 2024 Galvanized Logic Inc.

package vu

import (
	"image"
	"testing"

	// DEBUG to dump text block image as png.
	// "image/png"
	// "os"

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

// go test -run Text
func TestTextBlock(t *testing.T) {
	atlas, err := load.TTFont("18:lucon.ttf")
	if err != nil {
		t.Fatalf("unexpected font load error: %s", err)
	}
	f := newFont("lucon18")
	f.img = atlas.NRGBA
	f.w, f.h = int(atlas.Img.Width), int(atlas.Img.Height)
	for _, g := range atlas.Glyphs {
		f.addChar(g.Char, g.X, g.Y, g.W, g.H, g.Xo, g.Yo, g.Xa)
	}

	imgSize := 256 // image width and height in pixels
	img := image.NewNRGBA(image.Rect(0, 0, imgSize, imgSize))
	t.Run("success", func(t *testing.T) {
		if err := f.writeText("string1: 111", 0, 1, img); err != nil {
			t.Errorf("string1 should have worked %s", err)
		}
		if err := f.writeText("string2: 222", 200, 2, img); err != nil {
			t.Errorf("string2 should have worked %s", err)
		}

		// DEBUG to dump text block image as png.
		// textPng, _ := os.Create("textblock.png")
		// defer textPng.Close()
		// png.Encode(textPng, img)
	})

	t.Run("failure", func(t *testing.T) {
		if err := f.writeText("x1", 0, 1, nil); err == nil {
			t.Errorf("x1 should have failed %s", err)
		}
		if err := f.writeText("x2", -3, 1, nil); err == nil {
			t.Errorf("x2 should have failed %s", err)
		}
		if err := f.writeText("x3", 0, imgSize+10, nil); err == nil {
			t.Errorf("x3 should have failed %s", err)
		}
	})
}
