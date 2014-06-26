// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package load

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// FontData holds UV texture mapping information for a font.
// It is an intermediate format that is intended for populating
// render assets.
type FntData struct {
	W, H  int       // width and height
	Chars []ChrData // character data.
}

// ChrData holds UV texture mapping information for one character.
// It is an intermediate format that is intended for for populating
// render assets.
type ChrData struct {
	Char       rune // character.
	X, Y, W, H int  // character bit size.
	Xo, Yo, Xa int  // character offset.
}

// fnt reads in a text file describing the UV texture mapping for a
// character set of a particular font.
//
// The glyphs have been created using: www.anglecode.com/products/bmfont.
// The file data format is described at:
//    http://www.angelcode.com/products/bmfont/doc/file_format.html
func (l *loader) fnt(name string) (data *FntData, err error) {
	filename := name + ".fnt"

	// the header is the first line in the file.
	var file io.ReadCloser
	if file, err = l.getResource(l.dir[src], filename); err != nil {
		return nil, fmt.Errorf("Could not load glyphs from %s\n", filename, err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	// the second header line had the overall attributes.
	reader.ReadString('\n') // ignore the first header line.
	line, _ := reader.ReadString('\n')
	fields := strings.Fields(line)
	line = strings.Join(fields, " ")
	hfmt := "common lineHeight=%d base=%d scaleW=%d scaleH=%d pages=%d packed=%d alphaChnl=%d redChnl=%d greenChnl=%d blueChnl=%d"
	var lh, b, sw, sh, pgs, pkd, ac, rc, gc, bc int
	if _, err = fmt.Sscanf(line, hfmt, &lh, &b, &sw, &sh, &pgs, &pkd, &ac, &rc, &gc, &bc); err != nil {
		return nil, fmt.Errorf("Invalid glyph header in %s, %s\n", filename, err)
	}
	data = &FntData{sw, sh, []ChrData{}}

	// the bulk of the file is one data line per glyph
	dfmt := "char id=%d x=%d y=%d width=%d height=%d xoffset=%d yoffset=%d xadvance=%d page=%d chnl=%d"
	var gid, x, y, w, h, xo, yo, xa, p, c int
	for ; err == nil; line, err = reader.ReadString('\n') {
		fields := strings.Fields(line)
		line = strings.Join(fields, " ")

		// only process lines that match the expected format.
		if _, err := fmt.Sscanf(line, dfmt, &gid, &x, &y, &w, &h, &xo, &yo, &xa, &p, &c); err == nil {
			char, _ := utf8.DecodeRune([]byte{byte(gid)})
			data.Chars = append(data.Chars, ChrData{char, x, y, w, h, xo, yo, xa})
		}
	}
	return data, nil
}
