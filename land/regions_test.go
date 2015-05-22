// Copyright © 2014-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package land

import (
	"testing"
)

func TestRegions(t *testing.T) {
	if r := Regions(8, 4, 1234); r[2][4] != 1 {
		t.Errorf("Expected 1, but got %d", r[2][4])
	}
}

// ============================================================================
// benchmarks : go test -bench .
//
// Previous runs show that generating regions gets slower the more regions
// there are to generate.
//    BenchmarkRegions8_4	  100000	     20253 ns/op
//    BenchmarkRegions256_2	     500	   4579463 ns/op
//    BenchmarkRegions256_8	     100	  12141833 ns/op

// How long to create 4 regions on a small 8x8 map.
func BenchmarkRegions8_4(b *testing.B) {
	for cnt := 0; cnt < b.N; cnt++ {
		Regions(8, 4, 1234)
	}
}

// How long to create 2 regions on a standard 256x256 map.
func BenchmarkRegions256_2(b *testing.B) {
	for cnt := 0; cnt < b.N; cnt++ {
		Regions(256, 2, 1234)
	}
}

// How long to create 8 regions on a standard 256x256 map.
func BenchmarkRegions256_8(b *testing.B) {
	for cnt := 0; cnt < b.N; cnt++ {
		Regions(256, 8, 1234)
	}
}
