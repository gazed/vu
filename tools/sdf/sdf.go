// Copyright © 2017-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// Package sdf generates signed distance field images and font description
// text files from a pair of angel code bmfont font files. It works with a
// set of files generated from bmfont at 8x the expected final resolution.
//
// Some bmfont .bmfc config files are kept as examples and to be able to
// tweak fonts a few years down the road when all the original config
// settings have been forgotten.
//
// Running "sdf name" where:
//     name.png : is the bmfont generated font image atlas.
//     name.fnt : is the bmfont generated font layout text file.
// will produce : name_sdf.png, name_sdf.fnt
package main

// code based on
//    http://www.flatblackfilms.com/distance_field_source.zip
// which is based on code no longer available:
//    http://bitsquid.blogspot.ca/2010/04/distance-field-based-rendering-of.html
// Ported to golang this code deals with PNG/Text files instead of the TGA/XML
// files from the original CSharp code.

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"strings"
)

// sdf generates a signed distance field image from a set of angel code bmfont
// font files. Expects to work with a set of files generated at 8x the final
// size. Ie: a 2048x2048 font image is scaled down to a 256x256 font image.
//
// By convention all font information is contained in to files using the
// same name prefix (so yes, ensure the font image fits in 1 file instead
// of the multiple image page files)
//     name.png - the bmfont generated font image atlas.
//     name.fnt - the bmfont generated font layout text file.
//
// Ensure adding 24=8x3 pixels of padding around the characters to ensure
// the final distance fields do not overlap.
func main() {
	if len(os.Args) >= 2 {
		name := os.Args[1] // first one is program name.

		// Reasonable values for scale and spread
		//    scale : factor of 8 reduces a 2048 image to 256.
		//    spread: output image distance field clamp in pixels.
		//            ie: extend 3 pixels outside the character outline
		//            before clamping to 0.
		scale, spread := 8, 3
		convertImage(name+".png", scale, spread)
		convertFont(name+".fnt", scale)
	}
}

// convertImage reads in the given png image file and creates a new "*_sdf.png"
// file that contains the signed distance field results.
func convertImage(name string, scale, spread int) {
	reader, err := os.Open(name)
	if err != nil {
		log.Printf("Could not find image %s: %s\n", name, err)
		return
	}
	defer reader.Close()
	in, err := png.Decode(reader)
	if err != nil {
		log.Printf("Could not load image %s: %s\n", name, err)
	}
	out := image.NewNRGBA(image.Rect(0, 0, in.Bounds().Max.X/scale, in.Bounds().Max.Y/scale))
	transform(in.(*image.NRGBA), out, scale, spread)

	// write out the transformed image in a new file.
	target := strings.Replace(name, ".png", "_sdf.png", 1)
	writer, err := os.Create(target)
	if err != nil {
		log.Printf("Could not save image %s: %s\n", name, err)
	}
	defer writer.Close()
	png.Encode(writer, out)
}

// transform computes a distance field transform of a high resolution binary source channel
// and returns the result as a low resolution channel.
//    input:   The source channel
//    scale:   The amount the source channel will be scaled down.
//             A value of 8 means the destination image will be 1/8th the size
//             of the source image.
//    spread:  The spread in pixels before the distance field clamps to
//             (zero/one). The value is specified in units of the destination
//             image. The spread in the source image will be spread*scale_down.
func transform(in, out *image.NRGBA, scale, spread int) {
	clamp := float64(spread * scale)
	c := color.NRGBA{} // black. Alpha used to store signed distance value.

	// generate the value for each output pixel based on the larger input image.
	for y := out.Bounds().Min.Y; y < out.Bounds().Max.Y; y++ {
		for x := out.Bounds().Min.X; x < out.Bounds().Max.X; x++ {
			sd := signedDistance(in, x*scale+scale/2, y*scale+scale/2, clamp)
			c.A = uint8(((sd + clamp) / (clamp * 2)) * 255.0)
			out.Set(x, y, c)
		}
	}
}

// signedDistance returns the value to be used instead of the original
// color value.
func signedDistance(i *image.NRGBA, cx, cy int, clamp float64) float64 {
	w, h := i.Bounds().Max.X, i.Bounds().Max.Y
	cd := float64(i.At(cx, cy).(color.NRGBA).A)/255.0 - 0.5 // convert to -0.5 to 0.5

	min_x := cx - int(clamp-1)
	if min_x < 0 {
		min_x = 0
	}
	max_x := cx + int(clamp+1)
	if max_x >= w {
		max_x = w - 1
	}
	distance := clamp
	for dy := 0; dy <= int(clamp+1); dy++ {
		if float64(dy) > distance {
			continue
		}
		if cy-dy >= 0 {
			y1 := cy - dy
			for x := min_x; x <= max_x; x++ {
				if float64(x-cx) > distance {
					continue
				}
				d := float64(i.At(x, y1).(color.NRGBA).A)/255.0 - 0.5
				if cd*float64(d) < 0 {
					d2 := float64((y1-cy)*(y1-cy) + (x-cx)*(x-cx))
					if d2 < distance*distance {
						distance = math.Sqrt(d2)
					}
				}
			}
		}
		if dy != 0 && cy+dy < h {
			y2 := cy + dy
			for x := min_x; x <= max_x; x++ {
				if float64(x-cx) > distance {
					continue
				}
				d := float64(i.At(x, y2).(color.NRGBA).A)/255.0 - 0.5
				if cd*float64(d) < 0 {
					d2 := float64((y2-cy)*(y2-cy) + (x-cx)*(x-cx))
					if d2 < distance*distance {
						distance = math.Sqrt(d2)
					}
				}
			}
		}
	}
	if cd > 0 {
		return distance
	}
	return -distance
}

// Converts the specified font xml by adjusting the size and space layout
// information to match the newly scaled down font file.
func convertFont(name string, scale int) {
	reader, err := os.Open(name)
	if err != nil {
		log.Printf("Could not font file from %s: %s\n", name, err)
		return
	}
	defer reader.Close()
	rdr := bufio.NewReader(reader)

	// TODO process and rewrite the font file...
	target := strings.Replace(name, ".fnt", "_sdf.fnt", 1)
	writer, err := os.Create(target)
	if err != nil {
		log.Printf("Could not save font file %s: %s\n", name, err)
		return
	}
	defer writer.Close()
	wtr := bufio.NewWriter(writer)
	defer wtr.Flush()

	var line string
	for ; err == nil; line, err = rdr.ReadString('\n') {
		fields := strings.Fields(line)    // remove line feeds and replace multiple...
		text := strings.Join(fields, " ") // ...whitespace with a single space.
		if len(fields) <= 0 {
			wtr.WriteString(line)
			continue
		}
		switch fields[0] {
		case "info":
			info := "info face=%q size=%d bold=%d italic=%d charset=%q unicode=%d" +
				" stretchH=%d smooth=%d aa=%d padding=%d,%d,%d,%d spacing=%d,%d outline=%d"
			var f, cs string
			var s, b, i, uc, st, sm, aa, p0, p1, p2, p3, sp0, sp1, ol int
			if _, err = fmt.Sscanf(text, info, &f, &s, &b, &i, &cs, &uc, &st, &sm, &aa, &p0, &p1, &p2, &p3, &sp0, &sp1, &ol); err == nil {
				s = s / scale
				p0, p1, p2, p3 = 0, 0, 0, 0
				text = fmt.Sprintf(info+"\n", f, s, b, i, cs, uc, st, sm, aa, p0, p1, p2, p3, sp0, sp1, ol)
				wtr.WriteString(text)
			} else {
				log.Printf("info error %s", err)
			}
		case "common":
			hfmt := "common lineHeight=%d base=%d scaleW=%d scaleH=%d pages=%d packed=%d alphaChnl=%d redChnl=%d greenChnl=%d blueChnl=%d"
			var lh, b, sw, sh, pgs, pkd, ac, red, gc, bc int
			if _, err = fmt.Sscanf(text, hfmt, &lh, &b, &sw, &sh, &pgs, &pkd, &ac, &red, &gc, &bc); err == nil {
				lh, b = lh/scale, b/scale
				sw, sh = sw/scale, sh/scale
				text = fmt.Sprintf(hfmt, lh, b, sw, sh, pgs, pkd, ac, red, gc, bc)
				wtr.WriteString(text + "\n")
			}
		case "char":
			cfmt := "char id=%d x=%d y=%d width=%d height=%d xoffset=%d yoffset=%d xadvance=%d page=%d chnl=%d"
			var gid, x, y, w, h, xo, yo, xa, p, c int
			if _, err := fmt.Sscanf(text, cfmt, &gid, &x, &y, &w, &h, &xo, &yo, &xa, &p, &c); err == nil {
				p = 0 // expecting only single pages.
				x, y = x/scale, y/scale
				w, h = w/scale, h/scale
				xo, yo = xo/scale, yo/scale
				xa = xa / scale
				text = fmt.Sprintf("char id=%-4d x=%-5d y=%-5d width=%-5d height=%-5d"+
					" xoffset=%-5d yoffset=%-5d xadvance=%-5d page=%-2d chnl=%d\n",
					gid, x, y, w, h, xo, yo, xa, p, c)
				wtr.WriteString(text)
			}
		default:
			wtr.WriteString(text + "\n")
		}
	}
}
