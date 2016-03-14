// Copyright © 2014-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package land

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"testing"
)

func TestTopoSize(t *testing.T) {
	topo := NewTopo(256, 128)
	if x, y := topo.Size(); x != 256 || y != 128 {
		t.Error("Incorrect section size returned")
	}
}

// Uncomment to see a topo image. Recomment afterwards so as to not generate
// images during automated testing.
func Test2DTopoImageGeneration(t *testing.T) {
	// 	n := newNoise(124)
	// 	topo := NewTopo(256, 256)
	// 	topo.generate(0, 0, 0, n)
	// 	writeImage("target/", "topo0.png", topo.image(-0.25))
	// 	topo.generate(1, 0, 0, n)
	// 	writeImage("target/", "topo00.png", topo.image(-0.25))
}

// tests
// ============================================================================
// test helper functions.

// colorTile paints a tile into an existing image.
func colorTopo(t Topo, img *image.NRGBA, xoff, yoff int) {
	landSplit := 0.25
	var c *color.NRGBA
	for x := range t {
		for y := range t[x] {
			c = t.paint(t[x][y], landSplit)
			img.SetNRGBA(x+xoff, y+yoff, *c)
		}
	}
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
//   BenchmarkTopo1-8             100      14349602 ns/op (0.014 seconds)
//   BenchmarkTopo8-8              50      25198972 ns/op (0.025 seconds)
//   BenchmarkTopo17-8             30      38428378 ns/op (0.038 seconds)

// How long to create a single topo section.
func BenchmarkTopo1(b *testing.B) {
	topo := NewTopo(256, 256)
	n := newNoise(123)
	for cnt := 0; cnt < b.N; cnt++ {
		topo.generate(1, 0, 0, n)
	}
}

// Compare that to a single topo section at zoom level 8 (city size).
func BenchmarkTopo8(b *testing.B) {
	topo := NewTopo(256, 256)
	n := newNoise(123)
	for cnt := 0; cnt < b.N; cnt++ {
		topo.generate(8, 0, 0, n)
	}
}

// Compare that to a single topo section at zoom level 17 (world size).
func BenchmarkTopo17(b *testing.B) {
	topo := NewTopo(256, 256)
	n := newNoise(123)
	for cnt := 0; cnt < b.N; cnt++ {
		topo.generate(17, 0, 0, n)
	}
}
