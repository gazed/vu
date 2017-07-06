// Copyright © 2014-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// font.go encapsulates knowledge for displaying strings on screen.
// FUTURE: create 3D fonts from system fonts on the fly
//         or make a separate tool like:
//         http://www.angelcode.com/products/bmfont/

import (
	"github.com/gazed/vu/render"
)

// Labeler is a Model that displays a small text phrase. The Model
// combines a quad mesh, font mapping data, and a bitmapped font texture
// to display a string. Controlling a Labeler amounts to setting the
// string value and centering it using its width in pixels.
//
// Intended for single words or small phrases. Use another methods
// for large text layouts. The default color is white: 1,1,1.
type Labeler interface {
	SetStr(text string) Labeler // Set the string to display.
	SetWrap(w int) Labeler      // Set the string wrap length in pixels.
	StrSize() (w, h int)        // Width, Height in pixels, 0 if not loaded.
	StrColor(r, g, b float64)   // Label color where each value is from 0-1
}

// Labeler
// =============================================================================
// font is font mapping data needed by Labeler.

// font is an optional part of a rendered Model.
// font holds a single bitmapped font. It knows how to pull individual
// character images out of a single image that contains all the characters
// for a font. It is combined with a texture (the font bitmapped image)
// in order to produce displayable strings.
type font struct {
	name  string         // Unique id for a glyph set.
	tag   aid            // Name and type as a number.
	w, h  int            // Width and height of the entire font bitmap image.
	chars map[rune]*char // The "character" image information.

	// scratch variables reused to create rendered text phrases.
	vb []float32 // verticies.
	tb []float32 // texture mapping "uv" values.
	fb []uint16  // triangle face indicies.
}

// newFont allocates space for font mapping data.
func newFont(name string) *font {
	f := &font{name: name, tag: assetID(fnt, name)}
	f.chars = map[rune]*char{}
	return f
}

// aid is used to uniquely identify assets.
func (f *font) aid() aid      { return f.tag }  // hashed type and name.
func (f *font) label() string { return f.name } // asset name

// set font mapping data. Expected to be called by loader
// as fonts are loaded from disk.
func (f *font) setSize(w, h int) { f.w, f.h = w, h }
func (f *font) addChar(r rune, x, y, w, h, xo, yo, xa int) {
	uvs := f.uvs(x, y, w, h)
	f.chars[r] = &char{x, y, w, h, xo, yo, xa, uvs}
}

// setStr creates an image for the given string returning
// the verticies, and texture texture (uv) mapping information as
// a buffer slice.
//    wrap : optional (positive) width limit before the text wraps.
//
// The width in pixels for the resulting string image is returned.
func (f *font) setStr(m *mesh, str string, wrap int) (sx, sy int) {
	vb := f.vb[:0] // reset keeping allocated memory.
	tb := f.tb[:0] //  ""
	fb := f.fb[:0] //  ""

	// gather and arrange the letters for the phrase.
	width, height, fh, cnt := 0, 0, 0, 0
	for _, char := range str {
		c := f.chars[char]
		switch {
		case c != nil:
			fh = c.h // remember font height for wrapping with newlines.
			tb = append(tb, c.uvcs...)
			xo, yo := float32(c.xOffset), float32(c.yOffset)
			if c.w != 0 && c.h != 0 {

				// calculate the x, y positions based on desired locations.
				vb = append(vb,
					float32(width)+xo, float32(-height)+yo, 0, // upper left
					float32(c.w+width)+xo, float32(-height)+yo, 0, // upper right
					float32(c.w+width)+xo, float32(c.h-height)+yo, 0, // lower right
					float32(width)+xo, float32(c.h-height)+yo, 0) // lower left

				// keep track of the max size in pixels.
				if sx < c.w+width {
					sx = c.w + width
				}
				if sy < c.h+height {
					sy = c.h + height
				}
			}
			width += c.xAdvance
			if wrap > 0 && (width > wrap && char == ' ') {
				width = 0
				height += c.h
			}

			// create the triangles indexes referring to the points created above.
			i0 := uint16(cnt * 4)
			fb = append(fb, i0, i0+1, i0+3, i0+1, i0+2, i0+3)
			cnt += 1 // count characters rendered.
		case char == '\n':
			// auto wrap at newlines.
			width = 0
			height += fh
		}
	}
	m.InitData(0, 3, render.StaticDraw, false).SetData(0, vb)
	m.InitData(2, 2, render.StaticDraw, false).SetData(2, tb)
	m.InitFaces(render.StaticDraw).SetFaces(fb)
	f.vb = vb // reuse the allocated memory.
	f.tb = tb //   ""
	f.fb = fb //   ""
	return sx, sy
}

// uvs calculates the four UV points for one character. The x,y coordinates
// are the top left of the character. Note that the UV's are added to the array
// so as to match the order the vertices are created in panel(). This makes
// the letters appear the right way up rather than flipped.
//
// Expected to be used for loading fonts from disk.
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

// char represents a single bitmap character. Works well with
// http://www.angelcode.com/products/bmfont/doc/file_format.html
type char struct {
	x, y     int       // Top left corner.
	w, h     int       // Width and height.
	xOffset  int       // Current position offset for texture to screen.
	yOffset  int       // Current position offset for texture to screen.
	xAdvance int       // Current position advance after drawing character.
	uvcs     []float32 // Character bitmap texture coordinates 0:0, 1:0, 0:1, 1:1.
}
