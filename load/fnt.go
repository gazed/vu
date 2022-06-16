// Copyright © 2013-2024 Galvanized Logic Inc.

package load

// fnt.go imports font texture map data.
// The fnt data was produced with a tool that is no longer available.
// A new font data tool needs to be found.

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Fnt reads in a text file describing the UV texture mapping for
// a character set of a particular font.
//
// The Reader r is expected to be opened and closed by the caller.
// A successful import overwrites the data in FntData.
func Fnt(r io.Reader) (d *FontData, err error) {
	reader := bufio.NewReader(r)
	d = &FontData{}

	// the second header line had the overall attributes.
	reader.ReadString('\n') // ignore the first header line.
	line, _ := reader.ReadString('\n')
	fields := strings.Fields(line)
	line = strings.Join(fields, " ")
	hfmt := "common lineHeight=%d base=%d scaleW=%d scaleH=%d pages=%d packed=%d alphaChnl=%d redChnl=%d greenChnl=%d blueChnl=%d"
	var lh, b, sw, sh, pgs, pkd, ac, red, gc, bc int
	if _, err = fmt.Sscanf(line, hfmt, &lh, &b, &sw, &sh, &pgs, &pkd, &ac, &red, &gc, &bc); err != nil {
		return d, fmt.Errorf("Invalid glyph header %s\n", err)
	}
	d.W, d.H = sw, sh
	d.Chars = d.Chars[:0] // reuse existing memory if available.

	// the bulk of the file is one data line per glyph
	dfmt := "char id=%d x=%d y=%d width=%d height=%d xoffset=%d yoffset=%d xadvance=%d page=%d chnl=%d"
	var gid, x, y, w, h, xo, yo, xa, p, c int
	for ; err == nil; line, err = reader.ReadString('\n') {
		fields := strings.Fields(line)
		line = strings.Join(fields, " ")

		// only process lines that match the expected format.
		if _, err := fmt.Sscanf(line, dfmt, &gid, &x, &y, &w, &h, &xo, &yo, &xa, &p, &c); err == nil {
			d.Chars = append(d.Chars, ChrData{rune(gid), x, y, w, h, xo, yo, xa})
		}
	}
	return d, nil
}
