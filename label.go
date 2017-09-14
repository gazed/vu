// Copyright © 2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// label.go groups the 2D/3D string rendering code.
//          Also see font.go.

import (
	"log"
)

// MakeLabel adds a label component to an entity. It is intended for single
// words or small phrases.
//   shader: a font aware shader like "txt" or "sdf".
//   font  : identifies the font mapping file and font texture file, eg:
//           "lucidiaSu22" -> source/lucidiaSu22.fnt images/lucidiaSu22.png
// Internally it is a model with a quad mesh that displays a series of
// characters mapped from a font texture file. Manipulating a label amounts
// to setting the string value and centering it using its width in pixels.
// The default color is white: 1,1,1.
//
// Consider using a signed-distance-field (sdf) shader plus matching
// texture data to reduce pixelization for 3D labels.
func (e *Ent) MakeLabel(shader, font string) *Ent {
	e.app.models.createLabel(e, "shd:"+shader, "tex:"+font, "fnt:"+font)
	return e
}

// Typeset the given string, setting the labels string and
// regenerating the existing mesh. Will cause a mesh rebind.
//
// Depends on Ent.MakeLabel.
func (e *Ent) Typeset(s string) *Ent {
	l := e.app.models.getLabel(e.eid)
	v := e.app.models.get(e.eid)
	if l != nil && v != nil {
		if l.str != s && len(s) > 0 {
			l.str = s
			e.app.models.updateLabel(e.eid, l, v)
		}
		return e
	}
	log.Printf("Typeset needs MakeLabel %d", e.eid)
	return e
}

// SetWrap sets the string wrap length in pixels.
//
// Depends on Ent.MakeLabel.
func (e *Ent) SetWrap(w int) *Ent {
	l := e.app.models.getLabel(e.eid)
	v := e.app.models.get(e.eid)
	if l != nil && v != nil {
		if l.wrap != w {
			l.wrap = w
			e.app.models.updateLabel(e.eid, l, v)
		}
		return e
	}
	log.Printf("SetWrap needs MakeLabel %d", e.eid)
	return e
}

// Size returns the Label width, height in pixels. Returns 0 if not loaded
// or the Label is the empty string.
//
// Depends on Ent.MakeLabel.
func (e *Ent) Size() (w, h int) {
	if label := e.app.models.getLabel(e.eid); label != nil {
		return label.w, label.h
	}
	log.Printf("Size needs MakeLabel %d", e.eid)
	return 0, 0
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
	wrap int // Default 0. Non positive values ignored.
}
