// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"github.com/gazed/vu/render"
)

// font holds a single bitmapped font. It knows how to pull individual
// character images out of a single image that contains all the characters
// for a font. It has to be combined with a texture (the font bitmapped image)
// in order to produce displayable strings.
type font struct {
	name  string         // Unique id for a glyph set.
	w, h  int            // Width and height of the entire font bitmap image.
	chars map[rune]*char // The "character" image information.
}

// newFont allocates space for a complete font.
func newFont(name string) *font {
	f := &font{name: name}
	f.chars = map[rune]*char{}
	return f
}

// Name implements Font.
func (f *font) Name() string { return f.name }

// SetSize implements Font.
func (f *font) SetSize(w, h int) { f.w, f.h = w, h }

// AddChar implements Font.
func (f *font) AddChar(r rune, x, y, w, h, xo, yo, xa int) {
	uvs := f.uvs(x, y, w, h)
	f.chars[r] = &char{x, y, w, h, xo, yo, xa, uvs}
}

// panel creates a string image for the given string returning the
// verticies, and texture texture (uv) mapping information as a
// buffer slice. The buffer data is expected to be used for populating
// a Mesh.
//
// The width in pixels for the resulting string image is also returned.
func (f *font) Panel(m render.Mesh, phrase string) (width int) {
	vb := []float32{}
	tb := []float32{}
	fb := []uint16{}

	// gather and arrange the letters for the phrase.
	width = 0
	for cnt, char := range phrase {
		if c := f.chars[char]; c != nil {
			tb = append(tb, c.uvcs...)

			// skip spaces.
			xo, yo := float32(c.xOffset), float32(c.yOffset)
			if c.w != 0 && c.h != 0 {
				// calculate the x, y positions based on desired locations.
				xys := []float32{
					float32(width) + xo, yo, 0, // upper left
					float32(c.w+width) + xo, yo, 0, // upper right
					float32(c.w+width) + xo, float32(c.h) + yo, 0, // lower right
					float32(width) + xo, float32(c.h) + yo, 0, // lower left
				}
				vb = append(vb, xys...)
			}
			width += c.xAdvance

			// create the triangles indexes refering to the points created above.
			i0 := uint16(cnt * 4)
			fb = append(fb, i0, i0+1, i0+3, i0+1, i0+2, i0+3)
		}
	}
	m.InitData(0, 3, render.STATIC, false).SetData(0, vb)
	m.InitData(2, 2, render.STATIC, false).SetData(2, tb)
	m.InitFaces(render.STATIC).SetFaces(fb)
	return width
}

// uvs calculates the four UV points for one character.  The x,y coordinates
// are the top left of the character.  Note that the UV's are added to the array
// so as to match the order the vertices are created  in panel().  This makes
// the letters appear the right way up rather than flipped.
//
// Only expected to be used when loading fonts from disk.
func (f *font) uvs(x, y, w, h int) []float32 {
	uvcs := []float32{
		float32(x) / float32(f.w),   // lower left
		float32(y+h) / float32(f.h), // ""
		float32(x+w) / float32(f.w), // lower right
		float32(y+h) / float32(f.h), // ""
		float32(x+w) / float32(f.w), // upper right
		float32(y) / float32(f.h),   // ""
		float32(x) / float32(f.w),   // upper left
		float32(y) / float32(f.h),   // ""
	}
	return uvcs
}

// font
// ===========================================================================
// char

// char represents a single bitmap character.  See
// http://www.angelcode.com/products/bmfont/doc/file_format.html
type char struct {
	x, y     int       // Top left corner.
	w, h     int       // Width and height.
	xOffset  int       // Current position offset for texture to screen.
	yOffset  int       // Current position offset for texture to screen.
	xAdvance int       // Current position advance after drawing character.
	uvcs     []float32 // Character bitmap texture coordinates 0:0, 1:0, 0:1, 1:1.
}
