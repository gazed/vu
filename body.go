// Copyright © 2016-2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// body.go integrates the physics package into the engine.
// DESIGN: Expose entity related physics methods and provide a physics
//         body component manager.
// FUTURE: bodies currently share the pov's transform data. This only
//         works because pov is storing pointers to the transform data.

import (
	"log"

	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/physics"
)

// MakeBody creates non-colliding body associated with this Entity.
// Bodies are generally set on top level pov transforms which always
// have valid world coordindates.
//   b: a physics body eg: vu.Box, vu.Sphere.
func (e *Ent) MakeBody(b Body) *Ent {
	p := e.app.povs.get(e.eid)
	e.app.bodies.create(e.eid, b, p.tn)
	return e
}

// Body returns the physics body for this entity, returning nil
// if no physics body exists.
func (e *Ent) Body() Body { return e.app.bodies.get(e.eid) }

// DisposeBody removes the physics body from the given entity.
// Does nothing if there was no physics body.
func (e *Ent) DisposeBody() { e.app.bodies.dispose(e.eid) }

// SetSolid makes the existing physics Body collide.
//
// Depends on Ent.MakeBody.
func (e *Ent) SetSolid(mass, bounce float64) {
	if body := e.app.bodies.get(e.eid); body != nil {
		e.app.bodies.solidify(e.eid, mass, bounce)
		return
	}
	log.Printf("SetSolid needs MakeBody %d", e.eid)
}

// Cast checks if the ray intersects the given entity, returning
// the point of intersection if there is one. The point of contact
// x, y, z is valid when hit is true.
//
// Depends on Ent.MakeBody.
func (e *Ent) Cast(ray Body) (hit bool, x, y, z float64) {
	if body := e.app.bodies.get(e.eid); body != nil {
		return physics.Cast(ray, body)
	}
	log.Printf("Cast needs MakeBody %d", e.eid)
	return false, 0, 0, 0
}

// Push adds to the body's linear velocity.
// It is a wrapper for physics.Body.Push
//
// Depends on Ent.MakeBody.
func (e *Ent) Push(x, y, z float64) {
	if body := e.app.bodies.get(e.eid); body != nil {
		body.Push(x, y, z)
		return
	}
	log.Printf("Push needs MakeBody %d", e.eid)
}

// body entity methods
// =============================================================================
// bodies is the body component manager

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
	bs.physics = physics.NewPhysics()  // underlying physics engine.
	bs.shapes = map[eid]physics.Body{} // non-colliding bodies.
	bs.solids = map[eid]uint32{}       // Sparse map of colliding bodies.
	bs.bods = []physics.Body{}         // Dense array of colliding bodies...
	bs.eids = []eid{}                  // ...and associated entity identifiers.
	return bs
}

// create a new physics body. Guarantees that child pov's appear later in the
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
		b.SetProps(mass, bounce)

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

// stepVelocities runs physics on all the bodies; adjusting location and
// orientation. Physics has references to update the pov transform vectors.
func (bs *bodies) stepVelocities(dts float64) {
	bs.physics.Step(bs.bods, dts)
}
