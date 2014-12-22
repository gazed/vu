// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package land

import (
	"testing"
)

func TestSameSeed(t *testing.T) {
	n0 := newNoise(123)
	n1 := newNoise(123)
	if n0.generate(0, 0) != n1.generate(0, 0) {
		t.Error("Two generators with the same seed should produce the same result")
	}
}

// ============================================================================
// benchmarks : go test -bench .
//
// A previous run show that generating a single noise value is relatively quick.
// Of course a single map tile needs 256x256 noise values multiple times over.
//   BenchmarkNoise	100000000	        28.8 ns/op

// How long does creating a single random noise value take.
func BenchmarkNoise(b *testing.B) {
	n := newNoise(123)
	for cnt := 0; cnt < b.N; cnt++ {
		n.generate(10, 10)
	}
}
