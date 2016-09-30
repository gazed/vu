// Copyright © 2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// body.go helps wrap the vu/physics system.
// DESIGN: Provide a component manager for physics bodies that tracks
//         all physics instances.

import (
	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/physics"
)

// bodies manages all the active physics instances.
// Physics data is kept internally in order to facilitate optimizing
// the per-tick physics update.
type bodies struct {
	physics physics.Physics      // Physics system. Handles forces, collisions.
	shapes  map[eid]physics.Body // Non-colliding physic components.
	solids  map[eid]uint32       // Sparse map of colliding physic components.
	bods    []physics.Body       // Dense array of colliding physics bodies.
	eids    []eid                // Track last entity id to help with deletes.
}

// newBodies creates a manager for a group of physics data. Expectation
// is for a single instance to be created by the engine on startup.
func newBodies() *bodies {
	bs := &bodies{}
	bs.physics = physics.NewPhysics()
	bs.shapes = map[eid]physics.Body{}
	bs.solids = map[eid]uint32{}
	bs.bods = []physics.Body{} // physics bodies...
	bs.eids = []eid{}          // ...and associated entity identifiers.
	return bs
}

// create a new physics body. Guarantees that child Pov's appear later in the
// dense data array since children must be created after their parents.
func (bs *bodies) create(id eid, b physics.Body, t *lin.T) physics.Body {
	if bod, ok := bs.shapes[id]; ok {
		return bod // check non-colliding first.
	}
	if index, ok := bs.solids[id]; ok {
		return bs.bods[index] // check colliding.
	}
	b.SetWorld(t) // physics updates the passed in transform.
	bs.shapes[id] = b
	return b
}

// solidify takes an existing physics body and moves it to the list
// of colliding physics bodies.
func (bs *bodies) solidify(id eid, mass, bounce float64) {
	if b, ok := bs.shapes[id]; ok {
		delete(bs.shapes, id)
		b.SetMaterial(mass, bounce)

		// add the colliding body and update the indicies.
		bs.bods = append(bs.bods, b)
		bs.eids = append(bs.eids, id)
		bs.solids[id] = uint32(len(bs.bods)) - 1
	}
}

// get the physics body for the given id, returning nil if
// it does not exist.
func (bs *bodies) get(id eid) physics.Body {
	if index, ok := bs.solids[id]; ok {
		return bs.bods[index]
	}
	if b, ok := bs.shapes[id]; ok {
		return b
	}
	return nil
}

// dispose deletes the indicated physics body.
func (bs *bodies) dispose(id eid) {
	if _, ok := bs.shapes[id]; ok {
		delete(bs.shapes, id)
		return
	}

	// Otherwise the body is colliding.
	if index, ok := bs.solids[id]; ok {
		delete(bs.solids, id) // delete index from sparse array.

		// Save a mem copy by replacing deleted element with last element.
		// No other indicies need to be updated.
		lastIndex := len(bs.bods) - 1
		lastID := bs.eids[lastIndex]
		bs.eids[index] = bs.eids[lastIndex] // delete by replacing.
		bs.bods[index] = bs.bods[lastIndex] // delete by replacing.
		bs.eids = bs.eids[:lastIndex]       // discard moved last element.
		bs.bods = bs.bods[:lastIndex]       // discard moved last element.
		if id != lastID {
			bs.solids[lastID] = index // update the moved last elements index.
		}
	}
}

// stepVelocities runs physics on all the bodies; adjusting location and orientation.
func (bs *bodies) stepVelocities(eng *engine, dts float64) {
	bs.physics.Step(bs.bods, dts)

	// The associated pov needs to be marked as dirty in order to
	// update the transform.
	for _, eid := range bs.eids {
		if p := eng.povs.get(eid); p != nil {
			p.stable = false
		}
	}
	// FUTURE: benchmark to see if time can be saved by keeping an array of
	//         associated povs for updating the stable flag. How much Quicker
	//         than the current hash lookup?
}
