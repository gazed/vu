// Copyright © 2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"testing"
)

func TestEmptyValid(t *testing.T) {
	ids := &eids{}
	if ids.valid(0) {
		t.Errorf("Expecting invalid for unallocated entity")
	}

}
func TestFirstIsZero(t *testing.T) {
	ids := &eids{}
	if eid := ids.create(); eid != 0 {
		t.Errorf("Expecting first eid to be 0")
	}
}
func TestMaxCreate(t *testing.T) {
	ids := &eids{}
	for cnt := 0; cnt <= maxEntityID; cnt++ {
		if id := ids.create(); int(id) != cnt {
			t.Errorf("Expecting initial ids to be allocated sequentially.")
		}
	}

	// Check that one more than max is caught.
	// Should also generate a design error log.
	if id := ids.create(); id != 0 {
		t.Errorf("Expecting to have exhausted entity ids")
	}
}

func TestMaxCreateWithDispose(t *testing.T) {
	ids := &eids{}
	for cnt := 0; cnt <= maxEntityID; cnt++ {
		ids.create() // create max entities.
	}
	// should have allocated maxEntityID at this point

	// free 2*maxFree entities. Check that the free list can grow
	// larger than the amount that triggers reuse.
	for cnt := 0; cnt < 2*maxFree; cnt++ {
		ids.dispose(eid(cnt)) // should not crash.
	}
	if len(ids.free) != 2*maxFree {
		t.Errorf("Expected freelist %d to be %d", len(ids.free), 2*maxFree)
	}

	// should be able to re-allocate 2*maxFree entities.
	for cnt := 0; cnt < 2*maxFree; cnt++ {
		eid := ids.create()
		if eid == 0 {
			t.Errorf("Expecting to reuse disposed entity ids")
		}
	}

	// Check that one more than max is caught.
	// Should also generate a design error log.
	if eid := ids.create(); eid != 0 {
		t.Errorf("Expecting to have re-exhausted entity ids")
	}
}

// Tests
// =============================================================================
// Benchmarks.

// Hammer eids by creating and deleting as fast as possible.
// More of a stress test than a real usage case.
func BenchmarkCreateDelete(b *testing.B) {
	ids := &eids{}
	var id eid
	for cnt := 0; cnt < b.N; cnt++ {
		id = ids.create()
		ids.dispose(id)
	}
}
