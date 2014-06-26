// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package land

import (
	"image"
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
// func TestTopoImageGeneration(t *testing.T) {
// 	n := newNoise(124)
// 	topo := NewTopo(256, 256)
// 	topo.generate(0, 0, 0, n)
// 	writeImage("target/", "topo0.png", topo.image(-0.25))
// 	topo.generate(1, 0, 0, n)
// 	writeImage("target/", "topo00.png", topo.image(-0.25))
// }

// ============================================================================
// benchmarks : go test -bench .
//
// Previous runs show that generating single topo sections is expensive and it
// gets worse at the higher zoom levels where even more sections are needed.
//    BenchmarkTopo1	     100	  23572529 ns/op (0.023 seconds)
//    BenchmarkTopo8	      50	  60392420 ns/op (0.060 seconds)
//    BenchmarkTopo17	      20	 106809144 ns/op (0.106 seconds)

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

// benchmarks
// ============================================================================
// test helper functions.

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
