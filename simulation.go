// Copyright © 2024 Galvanized Logic Inc.

package vu

// simulation.go integrates physics into the engine.

import (
	"log/slog"

	"github.com/gazed/vu/physics"
)

// Body wraps physics.Body allowing the engine app to access common physics
// methods without including the physics package.
type Body *physics.Body

// Simulated bodies can either be kinematic-movable or static-immovable.
const (
	KinematicSim bool = false // movable
	StaticSim    bool = true  // immovable
)

// Sphere creates a ball shaped physics body located at the origin.
// The sphere size is defined by the radius.
func Sphere(radius float64, static bool) Body { return physics.NewSphere(radius, static) }

// Box creates a box shaped physics body located at the origin.
// The box size is given by the half-extents so that actual size
// is w=2*hx, h=2*hy, d=2*hz.
func Box(hx, hy, hz float64, static bool) Body { return physics.NewBox(hx, hy, hz, static) }

// AddToSimulation
// Bodies are generally set on top level pov transforms which always
// have valid world coordindates.
//
//	b: a physics body eg: vu.Box, vu.Sphere.
func (e *Entity) AddToSimulation(b Body) *Entity {
	if p := e.app.povs.get(e.eid); p == nil {
		slog.Error("AddToSimulation requires existing pov")
		return e
	}
	e.app.sim.create(e.eid, b)
	return e
}

// Push adds to the body's linear velocity.
// It is a wrapper for physics.Body.Push
//
// Depends on AddToSimultation.
func (e *Entity) Push(x, y, z float64) {
	if body := e.app.sim.get(e.eid); body != nil {
		(*physics.Body)(body).Push(x, y, z)
		return
	}
	slog.Error("Push needs AddToSimulation", "entity_id", e.eid)
}

// Body returns the physics body for this entity, returning nil
// if no physics body exists.
func (e *Entity) Body() Body { return e.app.sim.get(e.eid) }

// DisposeBody removes the physics body from the given entity.
// Does nothing if there was no physics body.
func (e *Entity) DisposeBody() { e.app.sim.dispose(e.eid) }

// =============================================================================
// simulation is the physics component manager

// simulation manages all the active physics instances.
type simulation struct {
	bids   map[eID]uint32 // Sparse mapping of eid to bid.
	bodies []physics.Body // Dense array of physics bodies, indexed by bid.
	eids   []eID          // Dense array of eids indexed by bid.
}

// newSimulation creates a manager for a group of physics data. Expectation
// is for a single instance to be created by the engine on startup.
func newSimulation() *simulation {
	sim := &simulation{}
	sim.bodies = []physics.Body{} // Dense array of physics bodies...
	sim.eids = []eID{}            // ...and associated entity identifiers.
	sim.bids = map[eID]uint32{}   // map entity ids to body ids.
	return sim
}

// create a new physics body. Guarantees that child pov's appear later in
// the dense data array since children must be created after their parents.
func (sim *simulation) create(id eID, b Body) Body {
	if bid, ok := sim.bids[id]; ok {
		return &sim.bodies[bid] // already exists.
	}
	bid := len(sim.bodies)              // body id is the array index.
	sim.bodies = append(sim.bodies, *b) // save body - indexed by bid
	sim.eids = append(sim.eids, id)     //  ""       - indexed by bid
	sim.bids[id] = uint32(bid)          // map eid to bid.
	return b
}

// get the physics body for the given id, returning nil if
// it does not exist.
func (sim *simulation) get(id eID) Body {
	if bid, ok := sim.bids[id]; ok {
		return &sim.bodies[bid]
	}
	return nil
}

// dispose deletes the indicated physics body.
func (sim *simulation) dispose(eid eID) {
	if index, ok := sim.bids[eid]; ok {
		delete(sim.bids, eid) // delete index from sparse array.

		// Save a mem copy by replacing deleted element with last element.
		// No other indicies need to be updated.
		lastIndex := len(sim.bodies) - 1
		lastID := sim.eids[lastIndex]             // eid of last index.
		sim.eids[index] = sim.eids[lastIndex]     // delete by replacing.
		sim.bodies[index] = sim.bodies[lastIndex] // delete by replacing.
		sim.eids = sim.eids[:lastIndex]           // discard moved last element.
		sim.bodies = sim.bodies[:lastIndex]       // discard moved last element.
		if eid != lastID {
			// if the deleted element wasn't the last...
			sim.bids[lastID] = index // ...update the moved element index.
		}
	}
}

// simulate runs physics on all the bodies; adjusting location and orientation.
// Expected to be called on regular timesteps from the main game loop.
func (sim *simulation) simulate(ps *povs, timestep float64) {

	// update simulation body transforms with povs that may have
	// been changed by the app.
	for i := range sim.bodies {
		bod := &sim.bodies[i]
		eid := sim.eids[i]
		p := ps.get(eid)
		if p == nil {
			slog.Error("physics body with no pov", "eid", eid)
			continue
		}
		bod.SetPosition(*p.tn.Loc)
		bod.SetRotation(*p.tn.Rot)
		bod.SetScale(*p.sw)
	}

	// run the physics simulation.
	physics.Simulate(sim.bodies, timestep)

	// apply any physics transform changes to the povs
	for i := range sim.bodies {
		bod := &sim.bodies[i]
		eid := sim.eids[i]
		p := ps.get(eid)
		if p == nil {
			slog.Error("physics body with no pov", "eid", eid)
			continue
		}
		p.tn.Loc.Set(bod.Position())
		p.tn.Rot.Set(bod.Rotation())
		// physics does not change scale.
		ps.updateWorld(p, eid)
	}
}
