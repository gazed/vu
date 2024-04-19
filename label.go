// Copyright © 2017-2024 Galvanized Logic Inc.

package vu

// label.go groups the 2D/3D string rendering code.

import (
	"log/slog"

	"github.com/gazed/vu/load"
)

// AddLabel adds a new label model and transform to the given entity.
// It is intended for single letters, words or small phrases. eg:
//
//	letter2D := scene.AddLabel("text", 0, "shd:icon", "fnt:lucidiaSu18", "tex:color:lucidiaSu18")
//	number3D := scene.AddLabel("text", 0, "shd:sdf", "fnt:lucidiaSu18", "tex:color:lucidiaSdf")
//
// A label requires a texture based shader, font mapping data, and a font
// texture atlas. The mesh is calculated from the string once the font
// assets have loaded.
//
// Signed Distance Fields (SDF) are used to reduce pixelization for 3D labels
// and require an SDF shader with a SDF texture atlas.
func (e *Entity) AddLabel(s string, wrap int, assets ...string) (me *Entity) {
	me = e.addPart() // add a transform node for the label.
	if mod := me.app.models.createLabel(s, wrap, me); mod != nil {
		mod.getAssets(me, assets...)
	}
	return me
}

// LabelSize returns the Label width, height in pixels. Returns 0 if not loaded.
//
// Depends on Ent.AddLabel.
func (e *Entity) LabelSize() (w, h int) {
	if m := e.app.models.get(e.eid); m != nil && m.mtype == labelModel && m.label != nil {
		return m.label.w, m.label.h
	}
	slog.Error("LabelSize needs label", "entity", e.eid)
	return 0, 0
}

// FUTURE: SetText to update a label string and regenerate a new mesh.
// FUTURE: SetWrap to update a label wrap and regenerate a new mesh.

// labelData is an internal call to get label information for the given entity.
func (e *Entity) labelData() (labelStr string, wrap int) {
	if m := e.app.models.get(e.eid); m != nil && m.mtype == labelModel && m.label != nil {
		return m.label.str, m.label.wrap
	}
	slog.Error("labelData needs label", "entity", e.eid)
	return "", 0
}

// setLabelMesh is an internal call to set the underlying mesh
// for the label.
func (e *Entity) setLabelMesh(msh *mesh, sx, sy int) {
	if m := e.app.models.get(e.eid); m != nil && m.mtype == labelModel && m.label != nil {
		m.label.w, m.label.h = sx, sy
		m.mesh = msh
		return
	}
	slog.Error("setLabelMesh needs label", "entity", e.eid)
}

// label entity methods
// =============================================================================
// label data

// label contains the data needed to render graphic strings.
type label struct {
	fnt *font  // Font asset.
	str string // Label text string.

	// Rendered string width and height in pixels.
	// Only valid after all label assets have been loaded since
	// sizes are dependent on font data.
	w, h int // 0 for nil strings or unloaded assets.

	// Sets a wrap amount for the string label in pixels.
	wrap int // Default 0. Negative values ignored.
}

// ============================================================================
// font is font mapping data needed by labels.
// font holds a single bitmapped font. It knows how to pull individual
// character images out of a single image that contains all the characters
// for a font. It is combined with a texture (the font bitmapped image)
// in order to produce displayable strings.
type font struct {
	name  string         // Unique id for a glyph set.
	tag   aid            // Name and type as a number.
	w, h  int            // Width and height of the entire font bitmap image.
	chars map[rune]*char // The "character" image information.
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
//
//	wrap : optional (positive) width in pixels for wrapping text.
//
// The pixel size and mesh data for the resulting string image is returned.
func (f *font) setStr(str string, wrap int) (sx, sy int, md load.MeshData) {
	vx := []float32{} // vec2 vertex data
	uv := []float32{} // vec2 texcoords
	ix := []uint16{}  // triangle indexes

	// gather and arrange the letters for the phrase.
	width, height, fh, cnt := 0, 0, 0, 0
	for _, char := range str {
		c := f.chars[char]
		switch {
		case c != nil:
			fh = c.h // remember font height for wrapping with newlines.
			uv = append(uv, c.uvcs...)
			xo, yo := float32(c.xOffset), float32(c.yOffset)
			if c.w != 0 && c.h != 0 {

				// calculate the x, y positions based on desired locations.
				vx = append(vx,
					float32(width)+xo, float32(height)+yo, // 0,0,
					float32(c.w+width)+xo, float32(height)+yo, // 1,0
					float32(width)+xo, float32(c.h+height)+yo, // 0,1
					float32(c.w+width)+xo, float32(c.h+height)+yo) // 1,1

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

			// create the triangles indexes referring to the points above.
			i0 := uint16(cnt * 4)
			ix = append(ix, i0, i0+2, i0+1, i0+1, i0+2, i0+3)
			cnt++ // count characters rendered.
		case char == '\n':
			// auto wrap at newlines.
			width = 0
			height += fh
		}
	}
	md = make(load.MeshData, load.VertexTypes)
	md[load.Vertexes] = load.F32Buffer(vx, 2)  // vec2
	md[load.Texcoords] = load.F32Buffer(uv, 2) // vec2
	md[load.Indexes] = load.U16Buffer(ix)
	return sx, sy, md
}

// uvs calculates the four UV points for one character.
// The UV's are added to the array so as to match the order the
// vertices are created above in setStr().
//
// Expected to be used for loading fonts from disk.
func (f *font) uvs(x, y, w, h int) []float32 {
	uvcs := []float32{
		float32(x) / float32(f.w),   // 0,0
		float32(y) / float32(f.h),   // ""
		float32(x+w) / float32(f.w), // 1,0
		float32(y) / float32(f.h),   // ""
		float32(x) / float32(f.w),   // 0,1
		float32(y+h) / float32(f.h), // ""
		float32(x+w) / float32(f.w), // 1,1
		float32(y+h) / float32(f.h), // ""
	}
	return uvcs
}

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
