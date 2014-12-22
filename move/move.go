// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.
//
// Huge thanks to bullet physics for showing what a physics engine is all about
// in cool-hard-code reality rather than theory. Methods and files that were
// derived from bullet physics are commented to indicate their origin.
// Bullet physics, for the most part, has the following license:
//
//   Bullet Continuous Collision Detection and Physics Library
//   Copyright (c) 2003-2006 Erwin Coumans  http://continuousphysics.com/Bullet/
//
//   This software is provided 'as-is', without any express or implied warranty.
//   In no event will the authors be held liable for any damages arising from the use of this software.
//   Permission is granted to anyone to use this software for any purpose,
//   including commercial applications, and to alter it and redistribute it freely,
//   subject to the following restrictions:
//
//   1. The origin of this software must not be misrepresented; you must not claim that you wrote the original software.
//      If you use this software in a product, an acknowledgment in the product documentation would be appreciated but is not required.
//   2. Altered source versions must be plainly marked as such, and must not be misrepresented as being the original software.
//   3. This notice may not be removed or altered from any source distribution.

// Move is a real-time simulation of real-world physics.  Move automatically
// applies simulated forces to virtual 3D objects known as bodies. Move
// updates bodies locations and directions based on forces and collisions
// with other bodies.
//
// Bodies are created using NewBody(shape). For example:
//    box    := NewBody(NewBox(hx, hy, hz))
//    sphere := NewBody(NewSphere(radius))
// Creating and storing bodies is the responsibility of the calling application.
// Bodies are moved with frequent and regular calls to Mover.Step(). Regulating
// the calls to Step() is the responsibility of the calling application. Once
// Step() has completed, the bodies updated location and direction are available
// in Body.World().
//
// Package move is provided as part of the vu (virtual universe) 3D engine.
package move

// See the open source physics engines:
//     www.bulletphysics.com
//     www.ode.org
// There is a 2D engine physics engine architecture overview at
//     http://gamedev.tutsplus.com/series/custom-game-physics-engine
// For regulating physics timesteps see:
//     http://gafferongames.com/game-physics/fix-your-timestep
// Other physics references:
//     http://www.geometrictools.com/Source/Physics.html

// Mover simulates forces acting on moving bodies. Expected usage
// is to simulate real-life conditions like air resistance and gravity,
// or the lack thereof.
type Mover interface {
	SetGravity(gravity float64) // Default is 10m/s.
	SetMargin(margin float64)   // Default is 0.04.

	// Step the physics simulation one tiny bit forward. This is expected
	// to be called regularly from the main engine loop. At the end of
	// a simulation step, all the bodies positions will be updated based on
	// forces acting upon them and/or collision results. Unmoving/unmoved
	// bodies, or bodies with zero mass are not updated.
	Step(bodies []Body, timestep float64)

	// Collide checks for collision between bodies a, b independent of
	// the current physics simulation. Bodies positions and velocities
	// are not updated. Provided for occasional or one-off checks.
	Collide(a, b Body) bool
}

// Mover interface
// ===========================================================================
// mover: default Mover implementation.

// mover is the default implementation of the Mover interface.
// It coordinates the physics pipeline by calling broadphase, narrowphase
// and solver.
type mover struct {
	gravity    float64                 // Force in m/s. Default is 10m/s.
	col        *collider               // Checks for collisions, updates collision contacts.
	sol        *solver                 // Resolves collisions, updates bodies locations.
	overlapped map[uint64]*contactPair // Overlapping pairs. Updated during broadphase.

	// scratch variables keep memory so that temp variables
	// don't have to be continually allocated and garbage collected
	abA, abB *Abox             // Scratch broadphase axis aligned bounding boxes.
	mf0      []*pointOfContact // Scratch narrowphase manifold.
}

// NewMover creates and returns a mover instance. Generally expected
// to be called once per application that needs a physics simulation.
func NewMover() Mover { return newMover() }
func newMover() *mover {
	mov := &mover{}
	mov.gravity = -10
	mov.col = newCollider()
	mov.sol = newSolver()
	mov.overlapped = map[uint64]*contactPair{}
	mov.mf0 = newManifold()
	mov.abA = &Abox{}
	mov.abB = &Abox{}
	return mov
}

// margin is a gap to smooth out collision detection.
var margin float64 = 0.04

// maxFriction is used to limit the amount of friction that
// can be applied to the combined friction of colliding bodies.
var maxFriction float64 = 100.0

// Mover interface implementation.
// Step the physics simulation forward by delta time (timestep).
// Note that the body.iitw is initialized once the first pass completes.
func (mov *mover) Step(bodies []Body, timestep float64) {

	// apply forces (e.g. gravity) to bodies and predict body locations
	mov.predictBodyLocations(bodies, timestep)

	// update overlapped pairs
	mov.broadphase(bodies, mov.overlapped)
	if len(mov.overlapped) > 0 {

		// collide overlapped pairs
		if colliding := mov.narrowphase(mov.overlapped); len(colliding) > 0 {
			mov.sol.info.timestep = timestep

			// resolve all colliding pairs
			mov.sol.solve(colliding, mov.overlapped)
		}
	}

	// adjust body locations based on velocities
	mov.updateBodyLocations(bodies, timestep)
	mov.clearForces(bodies)
}

// Mover interface implementation.
func (mov *mover) SetGravity(gravity float64)        { mov.gravity = gravity }
func (mov *mover) SetMargin(collisionMargin float64) { margin = collisionMargin }

// predictBodyLocations applies motion to moving/awake bodies as if there
// was nothing else around.
//
// Based on bullet btSimpleDynamicsWorld::predictUnconstraintMotion
func (mov *mover) predictBodyLocations(bodies []Body, dt float64) {
	var b *body
	for _, bb := range bodies {
		b = bb.(*body)
		b.guess.Set(b.world)
		if b.movable {

			// Fg = m*a. Apply gravity as if mass was 1.
			// FUTURE: use bodies mass when applying gravity.
			b.applyGravity(mov.gravity)    // updates forces.
			b.integrateVelocities(dt)      // applies forces to velocities.
			b.applyDamping(dt)             // damps velocities.
			b.updatePredictedTransform(dt) // applies velocities to prediction transform.
		}
	}
}

// broadphase checks for overlaps using the axis aligned bounding box
// for each body.
//
// FUTURE: create a broadphase bounding volume hierarchy to help with dealing
//         with a much larger number of bodies. Especially non-colliding bodies.
func (mov *mover) broadphase(bodies []Body, pairs map[uint64]*contactPair) {
	var bodyA, bodyB *body
	var uniques []Body
	var pairId uint64
	for cnt1, B1 := range bodies {
		bodyA = B1.(*body)
		uniques = bodies[cnt1+1:]
		for _, B2 := range uniques {
			bodyB = B2.(*body)

			// FUTURE: Add masking feature that allows bodies to only collide
			//         with other bodies that have matching mask types.

			// check as long as one of the bodies can move.
			if bodyA.movable || bodyB.movable {
				pairId = bodyA.pairId(bodyB)
				_, existing := pairs[pairId]
				if existing {
					abA := bodyA.predictedAabb(mov.abA, margin)
					abB := bodyB.predictedAabb(mov.abB, margin)
					overlaps := abA.Overlaps(abB)
					if !overlaps {
						// Remove existing
						delete(pairs, pairId)
					}
					// Otherwise Hold existing
				} else {
					abA := bodyA.worldAabb(mov.abA)
					abB := bodyB.worldAabb(mov.abB)
					overlaps := abA.Overlaps(abB)
					if overlaps {
						// Add new
						pairs[pairId] = newContactPair(bodyA, bodyB)
					}
					// Otherwise ignore non-overlapping pair
				}
			}
		}
	}
}

// narrowphase checks for actual collision. If bodies are colliding,
// then the persistent collision information for the bodies is updated.
// This includes the contact, normal, and depth information.
// Return all colliding bodies.
func (mov *mover) narrowphase(pairs map[uint64]*contactPair) (colliding map[uint32]*body) {
	colliding = map[uint32]*body{}
	scrManifold := mov.mf0 // scatch mf0
	for _, cpair := range pairs {
		bodyA, bodyB := cpair.bodyA, cpair.bodyB
		algorithm := mov.col.algorithms[bodyA.shape.Type()][bodyB.shape.Type()]
		bA, bB, manifold := algorithm(bodyA, bodyB, scrManifold)
		cpair.bodyA, cpair.bodyB = bA.(*body), bB.(*body) // handle potential body swaps.

		// bodies are colliding if there are contact points in the manifold.
		// Update any contact points and prepare for the solver.
		if len(manifold) > 0 {
			colliding[bodyA.bid] = bodyA
			colliding[bodyB.bid] = bodyB
			cpair.refreshContacts(bodyA.world, bodyB.world)
			cpair.mergeContacts(manifold)
		}
	} // scratch mf0 free
	return colliding
}

// updateBodyLocations applies the updated linear and angular velocities to the
// the bodies current position.
func (mov *mover) updateBodyLocations(bodies []Body, timestep float64) {
	var b *body
	for _, bb := range bodies {
		b = bb.(*body)
		if b.movable {
			b.updateWorldTransform(timestep)
			b.updateInertiaTensor()
		}
	}
}

// clearFoces removes any forces acting on bodies. This allows for the forces
// to be changed each simulation step.
func (mov *mover) clearForces(bodies []Body) {
	var b *body
	for _, bb := range bodies {
		b = bb.(*body)
		b.clearForces()
	}
}

// Collide returns true if the two shapes, a, b are touching or overlapping.
func (mov *mover) Collide(a, b Body) (hit bool) {
	aa, bb := a.(*body), b.(*body)
	algorithm := mov.col.algorithms[aa.shape.Type()][bb.shape.Type()]
	_, _, manifold := algorithm(aa, bb, mov.mf0)
	return len(manifold) > 0
}

// Cast checks if a ray r intersects the given Form f, giving back the
// nearest point of intersection if there is one. The point of contact
// x, y, z is valid when hit is true.
func Cast(ray, b Body) (hit bool, x, y, z float64) {
	if ray != nil && b != nil && b.Shape() != nil {
		if alg, ok := rayCastAlgorithms[b.Shape().Type()]; ok {
			return alg(ray, b)
		}
	}
	return false, 0, 0, 0
}
