// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package physics

// TODO one day, with a real physics engine, do something like the pipeline in
// Bullet Physics, i.e.:
//     forward dynamics
//        apply gravity
//        predict transforms
//     broadphase
//        detect pairs using AABB's
//     narrowphase
//        compute contacts
//     forward dynamics
//        solve constraints
//        integrate position.

import (
	"vu/math/lin"
)

// mover is a placeholder structure that groups the physics simulation code.
type mover struct{}

// step runs the physics simulation and pipeline.
func (m *mover) step(deltaTime float32, bodies map[int]Body, f *forces) {
	moved := m.move(deltaTime, bodies, f)
	m.checkCollisions(moved, bodies)
}

// move moves bodies that are under the control of physics and returns a list
// of all bodies that have moved. The moved list includes free bodies that are
// controlled outside of physics.
func (m *mover) move(deltaTime float32, bodies map[int]Body, f *forces) (moved []*body) {
	for _, b := range bodies {
		bod, _ := b.(*body)
		if bod.resolve != nil { // if free
			nl := bod.now.loc
			pl := bod.then.loc
			if nl.X != pl.X || nl.Y != pl.Y || nl.Z != pl.Z {
				moved = append(moved, bod)
			}
			bod.then = bod.now.clone()
		} else if bod.then.isMoving() || bod.now.isMoving() {
			moved = append(moved, bod)
			bod.then = bod.now.clone()
			bod.now.integrate(0, deltaTime)
		}
	}
	return moved
}

// checkCollisions allows collision checking independent of moving the physics bodies.
func (m *mover) checkCollisions(moved []*body, bodies map[int]Body) {
	collided := m.collide(moved, bodies)
	m.resolve(collided)
}

// collide finds all collision pairs by checking each body that has
// moved against every other body.
//
// TODO check if looking only at nearby bodies saves any time, ie. does the extra
//      bother of checking what's near outweigh the collision calculation?
// TODO should only check each pair once there are lots of moving bodies.
// TODO could check AABB collision first to cull obvious misses.
func (m *mover) collide(moved []*body, bodies map[int]Body) (collided []*collisionPair) {
	for _, bod1 := range moved {
		for _, b2 := range bodies {
			bod2, _ := b2.(*body)
			if bod1 == bod2 {
				continue
			}
			sh1 := bod1.getShape()
			sh2 := bod2.getShape()
			if sh1 != nil && sh2 != nil {
				if cons := sh1.Collide(sh2); cons != nil {
					cp := &collisionPair{cons, bod1, bod2, sh1, sh2}
					collided = append(collided, cp)
				}
			}
		}
	}
	return
}

// resolve collisions.
func (m *mover) resolve(collided []*collisionPair) {
	freeResolvers := map[int]*freeResolver{}
	for _, pair := range collided {

		// handle bodies that are not under the control of physics forces.
		// Consolidate all the contact points for a given body.
		if pair.b1.resolve != nil {
			uid := pair.b1.uid
			bid := pair.b2.uid
			for _, c := range pair.contacts {
				c.Bid = bid
			}
			if fr, ok := freeResolvers[uid]; ok {
				fr.contacts = append(fr.contacts, pair.contacts...)
			} else {
				freeResolvers[uid] = &freeResolver{pair.b1, pair.contacts}
			}
		} else {

			// handle bodies that are only controlled by physics forces.
			// TODO get rid of the hacks, especially how objects are brought to rest.
			pair.s1.Bounce(pair.s2, pair.contacts, pair.b1.now, nil)
			lm := pair.b1.now.linearMomentum()

			// correct the position since bounce only adjusts the velocities.
			c := pair.contacts[0]
			pair.b1.now.loc.Add(c.Normal.Unit().Scale(c.Depth))

			// if the motion is less than gravity, then bring the object to rest.
			if lm.Y <= 11 && pair.b1.now.loc.Y < 5 {
				pair.b1.now.f.gravity(0)
				pair.b1.now.setLinearMomentum(&lin.V3{0, 0, 0})
				pair.b1.now.setLinearMomentum(&lin.V3{0, 0, 0})
				pair.b1.now.loc.Y = 0.2
			}
		}
	}

	// Resolve the consolidate contacts for each body.
	for _, fr := range freeResolvers {
		fr.b.resolve(fr.contacts)
	}
}

// collisions are always between two bodies.
type collisionPair struct {
	contacts []*Contact
	b1, b2   *body
	s1, s2   Shape
}

// freeResolver is used to track the resolvers for objects that are not under
// the control of physics.  It is used and consumed in the resolve method.
type freeResolver struct {
	b        *body
	contacts []*Contact
}
