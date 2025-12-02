// SPDX-FileCopyrightText : © 2022 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

package vu

import (
	"testing"
)

// check that the entity IDs are properly allocated.
func TestEntityIDs(t *testing.T) {
	t.Run("zero is not a valid entityID", func(t *testing.T) {
		ents := &entities{}
		if ents.valid(0) {
			t.Errorf("expecting invalid for unallocated entity")
		}
		if ents.valid(1) {
			t.Errorf("expecting invalid for unallocated entity")
		}
	})
	t.Run("first valid entityID is one", func(t *testing.T) {
		ents := &entities{}
		if one := ents.create(); one != 1 {
			t.Errorf("expecting first eid to be 1")
		}
	})
	t.Run("disposed eids are not valid", func(t *testing.T) {
		ents := &entities{}
		one := ents.create()
		if !ents.valid(one) {
			t.Errorf("expected valid eid:%d edition:%d", one.id(), one.edition())
		}
		ents.dispose(one)
		if ents.valid(one) {
			t.Errorf("expected invalid eid:%d edition:%d", one.id(), one.edition())
		}
	})
	t.Run("allocate all entityIDs", func(t *testing.T) {
		ents := &entities{}
		for cnt := 1; cnt < maxEntID; cnt++ {
			if id := ents.create(); int(id) != cnt {
				t.Errorf("expecting initial ids to be allocated sequentially.")
			}
		}

		// check that allocating one more than max returns zero.
		if id := ents.create(); id != 0 {
			t.Errorf("expecting to have exhausted entity ids")
		}
	})
	t.Run("allocate more than max using dispose", func(t *testing.T) {
		ents := &entities{}
		for cnt := 1; cnt < maxEntID; cnt++ {
			ents.create() // create max entities.
		}
		// should have allocated maxEntID at this point

		// free 2*maxFree entities. Check that the free list can grow
		// larger than the amount that triggers reuse.
		for cnt := 1; cnt <= 2*maxFree; cnt++ {
			ents.dispose(eID(cnt)) // should not crash.
		}
		if len(ents.free) != 2*maxFree {
			t.Errorf("expected freelist %d to be %d", len(ents.free), 2*maxFree)
		}

		// should be able to reuse the disposed 2*maxFree entity IDs.
		for cnt := 0; cnt < 2*maxFree; cnt++ {
			eid := ents.create()
			if eid == 0 {
				t.Errorf("expecting to reuse disposed entity ids")
			}
		}

		// check that one more than max is caught.
		// Should also generate a design error log.
		if eid := ents.create(); eid != 0 {
			t.Errorf("expecting to have re-exhausted entity ids")
		}
	})
}

// Tests
// =============================================================================
// Benchmarks.

// go test -bench=.
// Hammer eids by creating and deleting as fast as possible.
// More of a stress test than a real usage case.
func BenchmarkCreateDelete(b *testing.B) {
	ents := &entities{}
	var id eID
	for cnt := 0; cnt < b.N; cnt++ {
		id = ents.create()
		ents.dispose(id)
	}
}
