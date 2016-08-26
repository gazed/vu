// Copyright © 2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package synth

import (
	"testing"
)

func TestExpand(t *testing.T) {
	if w := expand(0); w != 0 {
		t.Errorf("Expected 0, got %d", w)
	}
	if w := expand(3); w != 5 {
		t.Errorf("Expected 5, got %d", w)
	}
	if w := expand(5); w != 17 {
		t.Errorf("Expected 17, got %d", w)
	}
	if w := expand(0xffffffff); w != 6148914691236517205 {
		t.Errorf("Expected 6148914691236517205, got 0x%X", w)
	}
}

// Compact is the reverse of expand
func TestCompact(t *testing.T) {
	if c := compact(0); c != 0 {
		t.Errorf("Expected 0, got %d", c)
	}
	if c := compact(5); c != 3 {
		t.Errorf("Expected 3, got %d", c)
	}
	if c := compact(17); c != 5 {
		t.Errorf("Expected 5, got %d", c)
	}
	if c := compact(6148914691236517205); c != 0xffffffff {
		t.Errorf("Expected 0xffffffff, got 0x%X", c)
	}
}

func TestMerge(t *testing.T) {
	if k := ZMerge(3, 5); k != 39 {
		t.Errorf("Expected 39, got %d", k)
	}
	if k := ZMerge(5, 3); k != 27 {
		t.Errorf("Expected 27, got %d", k)
	}
}
func TestSplit(t *testing.T) {
	if a, b := ZSplit(39); a != 3 || b != 5 {
		t.Errorf("Expected 3, 5, got %d %d", a, b)
	}
	if a, b := ZSplit(27); a != 5 || b != 3 {
		t.Errorf("Expected 5, 3, got %d %d", a, b)
	}
}

func TestLabel(t *testing.T) {
	zoom := uint(4) // 4 bits per coordinate for zoom level 4.
	if key := ZLabel(zoom, ZMerge(0, 0)); key != "0000" {
		t.Errorf("Expected 0000, got %s", key)
	}
	if key := ZLabel(zoom, ZMerge(15, 15)); key != "3333" {
		t.Errorf("Expected 3333, got %s", key)
	}
	if key := ZLabel(zoom, ZMerge(7, 8)); key != "2111" {
		t.Errorf("Expected 2111, got %s", key)
	}
}

// ============================================================================
// Benchmarking

// Check ecoding speed. Run 'go test -bench=".*"
// For example the last run showed:
//    BenchmarkMerge-8   200000000            6.45 ns/op
func BenchmarkMerge(b *testing.B) {
	for cnt := 0; cnt < b.N; cnt++ {
		ZMerge(10, 2500)
	}
}
