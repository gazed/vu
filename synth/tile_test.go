// Copyright © 2014-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package synth

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"testing"
)

func TestTileSize(t *testing.T) {
	tile := newTile(256, 128, 0, 0, 0)
	if x, y := tile.Size(); x != 256 || y != 128 {
		t.Error("Incorrect section size returned")
	}
}

// Uncomment to see a topo image. Recomment afterwards so as to not generate
// images during automated testing.
func Test2DTileImageGeneration(t *testing.T) {
	// n := newSimplex(124)
	// tile := newTile(256, 256, 0, 0, 0)
	// img := image.NewNRGBA(image.Rect(0, 0, len(tile.topo), len(tile.topo[0])))
	// tile.gen2D(n)
	// writeImage("target/", "topo0.png", colorTile(tile.topo, img))
	// tile.Set(1, 0, 0)
	// tile.gen2D(n)
	// colorTile(tile.topo, img)
	// writeImage("target/", "topo00.png", img)
}

// tests
// ============================================================================
// test helper functions.

// colorTile paints a tile into an existing image and returns
// the input image.
func colorTile(t [][]float64, img *image.NRGBA) *image.NRGBA {
	landSplit := 0.25
	var c *color.NRGBA
	for x := range t {
		for y := range t[x] {
			c = paint(t[x][y], landSplit)
			img.SetNRGBA(x, y, *c)
		}
	}
	return img
}

// paint associates a color for the indicated section value.
// Note: this is for debugging only. Color probably should be a
// combination of land type and height.
func paint(height, landSplit float64) (c *color.NRGBA) {
	c = &color.NRGBA{255, 255, 255, 255}
	switch {
	case height > landSplit: // ground is uniform green for now.
		c = &color.NRGBA{0, 255, 0, 255}

	// shallower water is lighter.
	case height < landSplit:
		c = &color.NRGBA{10, 100, 200, 255}
	}
	return c
}

// write generates a png of the given image.
func writeImage(dir, name string, img *image.NRGBA) {
	os.Mkdir(dir, 0777)
	f, err := os.Create(dir + name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}

// test helper functions.
// ============================================================================
// benchmarks : go test -bench .
//
// Previous runs show that generating single topo sections is expensive and it
// gets worse at the higher zoom levels where even more sections are needed.
//   BenchmarkTile1-8             100      14349602 ns/op (0.014 seconds)
//   BenchmarkTile8-8              50      25198972 ns/op (0.025 seconds)
//   BenchmarkTile17-8             30      38428378 ns/op (0.038 seconds)

// How long to create a single topo section.
func BenchmarkTile1(b *testing.B) {
	tile := newTile(256, 256, 1, 0, 0)
	n := newSimplex(123)
	for cnt := 0; cnt < b.N; cnt++ {
		tile.gen2D(n)
	}
}

// Compare that to a single topo section at zoom level 8 (city size).
func BenchmarkTile8(b *testing.B) {
	tile := newTile(256, 256, 8, 0, 0)
	n := newSimplex(123)
	for cnt := 0; cnt < b.N; cnt++ {
		tile.gen2D(n)
	}
}

// Compare that to a single topo section at zoom level 17 (world size).
func BenchmarkTile17(b *testing.B) {
	tile := newTile(256, 256, 17, 0, 0)
	n := newSimplex(123)
	for cnt := 0; cnt < b.N; cnt++ {
		tile.gen2D(n)
	}
}
