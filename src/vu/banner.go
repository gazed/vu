// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package vu

import (
	"log"
	"vu/data"
)

// banner is a type of facade that displays bitmapped text strings on a flat
// mesh panel.
type banner struct {
	facade              // See facade.go.
	glys   *data.Glyphs // Optional, bit rendered font. A texture (font image file) must be specified.
	text   string       // Text to write on the banner.
	pixels int          // The number of horizontal pixels the banner uses.
}

// newBanner creates a banner with the supplied text and appearance attributes.
//    uid : unique identifier for this banner.
//    res : resource loader
//    text: the banner text.
//    shader, glyphs, texture : resource identifiers.
func newBanner(uid int, res *roadie, text, shader, glyphs, texture string) *banner {
	b := &banner{}
	b.uid = uid
	b.res = res
	b.msh = &data.Mesh{Name: "banner"}
	b.setShader(shader)
	b.setGlyphs(glyphs)
	b.setTexture(texture, 0)
	b.update(text)
	return b
}

// width gives back the number of pixels covered by the current banner text.
func (b *banner) width() int { return b.pixels }

// update revises the existing banner to show the new text string.
func (b *banner) update(text string) {
	if len(text) <= 0 {
		log.Printf("part:Banner.update ignoring empty text")
		return
	}
	b.text = text
	b.pixels = b.glys.Panel(b.msh, text)
	b.res.gc.BindGlyphs(b.msh)
}

// glyphs provides safe access to the glyphs unique label. Return an
// empty string for glyphs that have not yet been initialized.
func (b *banner) glyphs() string {
	if b.glys == nil {
		return ""
	}
	return b.glys.Name
}

// setGlyphs initializes the surface glyph information from one of the
// preloaded glyph sets.
func (b *banner) setGlyphs(glyphs string) {
	b.glys = b.res.useGlyphs(glyphs)
}
