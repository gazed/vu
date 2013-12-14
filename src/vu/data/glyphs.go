// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package data

// Glyphs holds a single bitmapped font. It knows how to pull individual
// character images out of a single image that contains all the characters
// for a font. It has to be combined with a texture (the font bitmapped image)
// in order to produce displayable strings.
//
// Glyphs can be used with any image of the same font and font size with identical
// bitmapped layouts, i.e. different colours of the same font and font size
// can use the same Glyph instance.
type Glyphs struct {
	Name   string          // Unique id for a glyph set.
	w, h   int             // Width and height of the entire font bitmap image.
	glyphs map[rune]*glyph // The "characters".
}

// Panel creates a string image for the given string using a flat mesh.
// The mesh has the necessary texture (uv) mapping information and is ready
// to be combined with a texture corresponding to the bitmapped font.
// Note that the mesh has not yet been bound to the graphics card.
//
// The width in pixels for the resulting string image is returned.
func (gs *Glyphs) Panel(b *Mesh, phrase string) (panelWidth int) {
	if b == nil {
		b = &Mesh{Name: "banner"}
	} else {
		b.V = []float32{}
		b.T = []float32{}
		b.F = []uint16{}
	}

	// gather and arrange the letters for the phrase.
	panelWidth = 0
	for cnt, char := range phrase {
		g := gs.glyphs[char]
		b.T = append(b.T, g.uvcs...)

		// skip spaces.
		if g.w != 0 && g.h != 0 {
			// calculate the x, y positions based on desired locations.
			xys := []float32{
				float32(panelWidth + g.xOffset), // upper left
				float32(g.yOffset), 0, 1,        // ""
				float32(g.w + panelWidth + g.xOffset), // upper right
				float32(g.yOffset), 0, 1,              // ""
				float32(g.w + panelWidth + g.xOffset), // lower right
				float32(g.h + g.yOffset), 0, 1,        // ""
				float32(panelWidth + g.xOffset), // lower left
				float32(g.h + g.yOffset), 0, 1,  // ""
			}
			b.V = append(b.V, xys...)
		}
		panelWidth += g.xAdvance

		// create the triangles indexes refering to the points created above.
		i0 := uint16(cnt * 4)
		b.F = append(b.F, i0, i0+1, i0+3, i0+1, i0+2, i0+3)
	}
	return panelWidth
}

// uvs calculates the four UV points for one glyph.  The x,y coordinates
// are the top left of the glyph.  Note that the UV's are added to the array
// so as to match the order the vertices are created later on.  This makes
// the letters appear the right way up rather than flipped.
//
// Only expected to be used when loading glyphsets from disk.
func (gs *Glyphs) uvs(x, y, w, h int) []float32 {
	uvcs := []float32{
		float32(x) / float32(gs.w),   // lower left
		float32(y+h) / float32(gs.h), // ""
		float32(x+w) / float32(gs.w), // lower right
		float32(y+h) / float32(gs.h), // ""
		float32(x+w) / float32(gs.w), // upper right
		float32(y) / float32(gs.h),   // ""
		float32(x) / float32(gs.w),   // upper left
		float32(y) / float32(gs.h),   // ""
	}
	return uvcs
}

// Glyphs
// ===========================================================================
// glyph

// glyph represents a single bitmap character.  See
// http://www.angelcode.com/products/bmfont/doc/file_format.html
type glyph struct {
	x, y     int       // Top left corner.
	w, h     int       // Width and height.
	xOffset  int       // Current position offset for texture to screen.
	yOffset  int       // Current position offset for texture to screen.
	xAdvance int       // Current position advance after drawing character.
	uvcs     []float32 // Character bitmap texture coordinates 0:0, 1:0, 0:1, 1:1.
}
