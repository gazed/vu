// SPDX-FileCopyrightText : © 2024-2025 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

package load

// ttf.go loads truetype fonts into an internal data format needed
// to create the font atlases and corresponding font layout information.
// Cobbled togehter based on some minimal atlas examples:
// - https://github.com/udhos/ratlas (uses golang freetype instead of x/image)
// - https://gist.github.com/baines/b0f9e4be04ba4e6f56cab82eef5008ff  (C + freetype)

import (
	"fmt"
	"image"
	"image/draw"
	"log/slog"

	// DEBUG to dump atlas image as png.
	// "image/png"
	// "os"

	// coded using golang.org/x/image v0.20.0 (expect API changes)
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// TTFRunes can be overridden by the application.
// Default: attampt to load basic runes plus some symbols.
var TTFRunes = []rune(" ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890`~!@#$%^&*()[]{}/=?+\\|-_.>,<'\";:")

// Ttf reads the truetype font and generates the atlas image and atlas
// character mapping data.
func Ttf(ttfBytes []byte, size int) (atlas *FontAtlas, err error) {
	f, err := opentype.Parse(ttfBytes)
	if err != nil {
		return nil, fmt.Errorf("openttype parse:%w", err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return nil, fmt.Errorf("openttype face:%w", err)
		return nil, err
	}

	// A reasonable amount of runes with a reasonable font size should easily
	// fit into a 512x512 image. FUTURE: make this configurable.
	imgSize := 512 // image width and height in pixels

	atlas = &FontAtlas{}
	img := image.NewNRGBA(image.Rect(0, 0, imgSize, imgSize))
	penx, peny := 0, 0
	lineHeight := face.Metrics().Height.Round()
	XAscent := face.Metrics().Ascent.Round()
	for _, r := range TTFRunes {
		bounds, _, ok := face.GlyphBounds(r)
		if !ok {
			slog.Error("failed to find rune", "rune", r, "char", string(r))
			continue
		}
		minX := bounds.Min.X.Floor()
		minY := bounds.Min.Y.Floor()
		maxX := bounds.Max.X.Ceil()
		maxY := bounds.Max.Y.Ceil()
		glyphWidth := maxX - minX + 2 // width padding.
		glyphHeight := maxY - minY
		descent := int(float32(maxY) + (float32(bounds.Min.Y)/64.0 - float32(minY)))
		bearingX := int(float32(bounds.Min.X) / 64.0)

		// advance to the next line if necessary.
		if penx+glyphWidth >= imgSize {
			penx = 0
			peny += lineHeight
			if peny >= imgSize {
				return nil, fmt.Errorf("atlas image to small %d", imgSize)
			}
		}

		// create glyph image
		dst := image.NewNRGBA(image.Rect(0, 0, glyphWidth, glyphHeight))
		d := &font.Drawer{
			Dot:  fixed.P(-minX+1, -minY), // glyph origin with width pad adjustment.
			Dst:  dst,
			Src:  image.White,
			Face: face,
		}
		dr, mask, maskp, xadvance, _ := d.Face.Glyph(d.Dot, r)
		draw.DrawMask(d.Dst, dr, d.Src, image.Point{}, mask, maskp, draw.Over)

		// copy glyph image to atlas image aliging the glyphs
		// on the baseline within identical height boxes.
		base := maxY - descent + (XAscent + minY)
		draw.Draw(img, image.Rect(penx, peny+base, penx+glyphWidth, peny+base+glyphHeight), dst, image.Point{}, draw.Src)

		// capture the glyph atlas position for rendering text.
		xoff := bearingX
		yoff := 0 // descent is built into the font placement.
		g := Glyph{r, penx, peny, glyphWidth, lineHeight, xoff, yoff, xadvance.Round()}
		atlas.Glyphs = append(atlas.Glyphs, g)
		penx += glyphWidth
	}
	atlas.Img.Pixels = []byte(img.Pix)
	atlas.Img.Width = uint32(img.Bounds().Size().X)
	atlas.Img.Height = uint32(img.Bounds().Size().Y)
	atlas.Img.Opaque = true
	atlas.NRGBA = img

	// DEBUG dump atlas image as png.
	// atlasPng, _ := os.Create("atlas.png")
	// defer atlasPng.Close()
	// png.Encode(atlasPng, img)

	return atlas, nil
}
