// Copyright © 2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// eid.go defines entity identfiers. See aid.go for asset identifiers.

import (
	"log"
)

// eid is an entity identfier comprised of an id used as a live reference
// to data and an edition used to track when ids are deleted and reused.
// Entity ids are expected to be used as array indicies for component
// data and as such they will not change values over their lifetime.
type eid uint32

// Divide the entity bits into a index id and an edition. The edition
// bits are used to track when an entity has been deleted.
const idBits = 20                     // entity array index : 1048575
const edBits = 12                     // entity edition     :    4096
const maxEntityID = (1 << idBits) - 1 // mask and max active entities.
const maxEdition = (1 << edBits) - 1  // mask and max dispose and reuse.

// id is the value to be used for array lookups.
func (e eid) id() uint32 { return uint32(e & maxEntityID) }

// edtion returns the value that tracks if the id is valid.
func (e eid) edition() uint16 { return uint16((e >> idBits) & maxEdition) }

// eid
// =============================================================================
// eids see:
// http://bitsquid.blogspot.ca/2014/08/building-data-oriented-entity-system.html

// eids handles the creation and deletion of entity identifiers.
// It ensures a limited set of unique identifiers. These identifiers
// are limited so they can be used as indicies into arrays of data.
type eids struct {

	// Starts empty and grows as entities are allocated.
	// Max size is maxEntityID.
	editions []uint16 // track currently used entities.

	// Starts empty and grows as entities are disposed.
	// New entities are allocated from here once it reaches maxFree.
	free []uint32 // tracks entities ready for reuse.
}

// maxFree starts recyling ids once the amount of disposed ids
// reaches the given size.
const maxFree = (1 << (edBits - 1)) // recycling when free reaches 2048.

// create returns a new entity id. Zero is returned for the first entity
// and when all entity identifiers have been allocated.
func (ids *eids) create() eid {
	id := uint32(0)
	if len(ids.free) > maxFree {
		id = ids.free[0]
		ids.free = append(ids.free[:0], ids.free[1:]...)
	} else {
		ids.editions = append(ids.editions, 0)
		if id = uint32(len(ids.editions) - 1); id > maxEntityID {

			// entity ids exhausted if nothing in free list.
			if len(ids.free) == 0 {
				log.Printf("All %d entity identifiers in use", maxEntityID+1)
				return 0 // design error to be caught during development.
			}
			id = ids.free[0]
			ids.free = append(ids.free[:0], ids.free[1:]...)
		}
	}
	return eid(id | uint32(ids.editions[id])<<idBits)
}

// valid entities are those that have been created and not yet disposed.
func (ids *eids) valid(e eid) bool {
	id := e.id()
	if id >= uint32(len(ids.editions)) {
		return false
	}
	return ids.editions[e.id()] == e.edition()
}

// dispose marks an entity as no longer valid. The entity identifer
// is placed for reallocation at a later date. The entity can be reallocated
// maxEdition times before it duplicates a previously generated entity.
func (ids *eids) dispose(e eid) {
	id := e.id()
	ids.editions[id]++              // mark this entity id as invalid.
	ids.free = append(ids.free, id) // queue it up for reallocation.
}

// reset discards all entity information and puts the entity counters
// back to the state when it was first created.
func (ids *eids) reset() {
	ids.editions = []uint16{}
	ids.free = []uint32{}
}
