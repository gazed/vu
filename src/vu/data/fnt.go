// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package data

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// fnt reads in a text file describing the UV texture mapping for a
// character set of a particular font.
//
// The glyphs have been created using: www.anglecode.com/products/bmfont.
// The file data format is described at:
//    http://www.angelcode.com/products/bmfont/doc/file_format.html
func (l loader) fnt(directory, filename string) (glyphs *Glyphs, err error) {

	// the header is the first line in the file.
	hfmt := "common lineHeight=%d base=%d scaleW=%d scaleH=%d pages=%d packed=%d alphaChnl=%d redChnl=%d greenChnl=%d blueChnl=%d"
	var lh, b, sw, sh, pgs, pkd, ac, rc, gc, bc int
	var file io.ReadCloser
	if file, err = l.getResource(directory, filename); err != nil {
		return nil, fmt.Errorf("Could not load glyphs from %s\n", filename, err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	// the second header line had the overall attributes.
	reader.ReadString('\n') // ignore the first header line.
	line, _ := reader.ReadString('\n')
	fields := strings.Fields(line)
	line = strings.Join(fields, " ")
	if _, err = fmt.Sscanf(line, hfmt, &lh, &b, &sw, &sh, &pgs, &pkd, &ac, &rc, &gc, &bc); err != nil {
		return nil, fmt.Errorf("Invalid glyph header in %s, %s\n", filename, err)
	}
	glyphSet := &Glyphs{filename, sw, sh, map[rune]*glyph{}}

	// the bulk of the file is one data line per glyph
	dfmt := "char id=%d x=%d y=%d width=%d height=%d xoffset=%d yoffset=%d xadvance=%d page=%d chnl=%d"
	var gid, x, y, w, h, xo, yo, xa, p, c int
	for ; err == nil; line, err = reader.ReadString('\n') {
		fields := strings.Fields(line)
		line = strings.Join(fields, " ")

		// only process lines that match the expected format.
		if _, err := fmt.Sscanf(line, dfmt, &gid, &x, &y, &w, &h, &xo, &yo, &xa, &p, &c); err == nil {
			glyph := &glyph{x, y, w, h, xo, yo, xa, glyphSet.uvs(x, y, w, h)}
			char, _ := utf8.DecodeRune([]byte{byte(gid)})
			glyphSet.glyphs[char] = glyph
		}
	}
	return glyphSet, nil
}
