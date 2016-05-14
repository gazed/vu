// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package physics

// // The following block is C code and cgo directvies.
// // It is used to include collision.h definitions.
//
// #cgo CFLAGS: -lm
//
// #include "collision.h"
import "C" // must be located here.

import (
	"log"
	"math"
	"sync"

	"github.com/gazed/vu/math/lin"
)

// Body is a single object contained within a physics simulation.
// Bodies generally relate to a scene node that is displayed to the user.
// Only add bodies that need to participate in physics.
// Bodies that are added to physics are expected to have their movement
// controlled by the physics simulation and not the application.
type Body interface {
	Shape() Shape          // Physics shape for this form.
	World() *lin.T         // Get the location and direction
	SetWorld(world *lin.T) // ...or set the location and direction.

	Eq(b Body) bool           // True if the two bodies are the same.
	Speed() (x, y, z float64) // Current linear velocity.
	Whirl() (x, y, z float64) // Current angular velocity.
	Push(x, y, z float64)     // Add to the body's linear velocity.
	Turn(x, y, z float64)     // Add to the body's angular velocity.
	Stop()                    // Stops linear velocity.
	Rest()                    // Stops angular velocity.

	// SetMaterial associates physical properties with a body. The physical
	// properties are combined with the body's shape to determine its behaviour
	// during collisions. The updated Body is returned.
	//     mass:       use zero mass for unmoving (static/fixed) bodies.
	//     bounciness: total bounciness is determined by multiplying the bounciness
	//                 of the two colliding bodies. If one of the bodies has 0
	//                 bounciness then there is no bounce effect.
	SetMaterial(mass, bounciness float64) Body
}

// Body interface
// ===========================================================================
// body implementation.

// body is the default implementation of the Body interface.
type body struct {
	bid   uint32  // Unique body id for generating pair identfiers.
	shape Shape   // Body shape for collisions.
	world *lin.T  // World transform for the given shape.
	v0    *lin.V3 // Scratch vector.

	guess   *lin.T // Predicted world transform for the given shape.
	movable bool   // Body has mass. It is able to move.

	// Motion data
	imass float64 // Inverse mass is calcuated once on object creation.
	lvel  *lin.V3 // Linear velocity in meters per second.
	lfor  *lin.V3 // Linear forces acting on this body.
	ldamp float64 // Linear damping.
	avel  *lin.V3 // Angular velocity.
	afor  *lin.V3 // Angular forces (torque) acting on this body.
	adamp float64 // Angular damping.
	iit   *lin.V3 // Inverse inertia tensor.
	iitw  *lin.M3 // Inverse inertia tensor world. Tracks oriented inertia amount.

	// Bodys take part in collision resolution. Tracks the extra information
	// needed by the solver. It is initialized and consumed by the solver as needed.
	friction    float64     // Ideally non-zero.
	restitution float64     // Bounciness. Zero to one expected.
	sbod        *solverBody // Body related solver data.

	// Scratch variables are optimizations that avoid creating/destroying
	// temporary objects that are needed each timestep.
	coi    *C.BoxBoxInput   // Scratch box-box collision input.
	cor    *C.BoxBoxResults // Scratch box-box collision output.
	m0, m1 *lin.M3          // Scratch matrices.
	t0     *lin.T           // Scratch transform.
}

// bodyUuid is a cheap simple global id. Allows 4 billion bodies before
// luck takes over. FUTURE: need a body.Dispose() method to allow reuse
// of body ids.
var bodyUUID uint32
var bodyUUIDMutex sync.Mutex // Concurrency safety.

// NewBody returns a new Body structure. The body will
// be positioned, with no rotation, at the origin.
func NewBody(shape Shape) Body { return newBody(shape) }
func newBody(shape Shape) *body {
	b := &body{}
	b.shape = shape
	b.imass = 0                 // no mass, static body by default
	b.friction = 0.5            // good to have some friction
	b.world = lin.NewT().SetI() // world transform
	b.guess = lin.NewT().SetI() // predicted world transform

	// allocate linear and angular motion data
	b.lvel = lin.NewV3()
	b.lfor = lin.NewV3()
	b.avel = lin.NewV3()
	b.afor = lin.NewV3()
	b.iitw = lin.NewM3().Set(lin.M3I)
	b.iit = lin.NewV3()

	// allocate scratch variables
	b.coi = &C.BoxBoxInput{}
	b.cor = &C.BoxBoxResults{}
	b.m0 = &lin.M3{}
	b.m1 = &lin.M3{}
	b.v0 = &lin.V3{}
	b.t0 = lin.NewT()

	// create a unique body identifier
	bodyUUIDMutex.Lock()
	b.bid = bodyUUID
	if bodyUUID++; bodyUUID == 0 {
		log.Printf("Overflow: dev error. Unique body id wrapped.")
	}
	bodyUUIDMutex.Unlock()
	return b
}

// Form interface implementation.
func (b *body) Shape() Shape { return b.shape }

// Allow world to be injected so that it becomes shared data.
// Lazy create the world transform if one was not set.
func (b *body) SetWorld(world *lin.T) { b.world = world }
func (b *body) World() *lin.T {
	if b.world == nil {
		b.world = lin.NewT().SetI()
	}
	return b.world
}

// Body interface implementation.
func (b *body) Eq(a Body) bool           { return b.bid == a.(*body).bid }
func (b *body) Speed() (x, y, z float64) { return b.lvel.X, b.lvel.Y, b.lvel.Z }
func (b *body) Whirl() (x, y, z float64) { return b.avel.X, b.avel.Y, b.avel.Z }
func (b *body) Stop()                    { b.lvel.X, b.lvel.Y, b.lvel.Z = 0, 0, 0 }
func (b *body) Rest()                    { b.avel.X, b.avel.Y, b.avel.Z = 0, 0, 0 }
func (b *body) Push(x, y, z float64) {
	b.lvel.X += x
	b.lvel.Y += y
	b.lvel.Z += z
}
func (b *body) Turn(x, y, z float64) {
	b.avel.X += x
	b.avel.Y += y
	b.avel.Z += z
}
func (b *body) SetMaterial(mass, bounciness float64) Body {
	return b.setMaterial(mass, bounciness)
}
func (b *body) setMaterial(mass, bounciness float64) *body {
	b.imass = 0 // static unless there is mass.
	if !lin.AeqZ(mass) {
		b.imass = 1.0 / mass                 // only need inverse mass
		b.iit = b.shape.Inertia(mass, b.iit) // shape inertia
		if lin.AeqZ(b.iit.X) {               // inverse shape inertia
			b.iit.X = 0
		} else {
			b.iit.X = 1.0 / b.iit.X
		}
		if lin.AeqZ(b.iit.Y) {
			b.iit.Y = 0
		} else {
			b.iit.Y = 1.0 / b.iit.Y
		}
		if lin.AeqZ(b.iit.Z) {
			b.iit.Z = 0
		} else {
			b.iit.Z = 1.0 / b.iit.Z
		}
	}
	b.restitution = bounciness
	b.movable = b.imass != 0
	return b
}

// pairID generates a unique id for bodies a and b.
// The pair id is independent of calling order.
func (b *body) pairID(a *body) uint64 {
	id0, id1 := b.bid, a.bid
	if id0 > id1 {
		id0, id1 = id1, id0 // calling order independence
	}
	return uint64(id0)<<32 + uint64(id1)
}

// applyGravity applies the force of gravity to the total forces
// acting on this body. Static bodies are ignored.
func (b *body) applyGravity(gravity float64) {
	if b.movable {
		b.lfor.Y += gravity
	}
}

// updateInertiaTensor reacalculates the inertia tensor for this body.
func (b *body) updateInertiaTensor() {
	worldBasis, basisTransposed := b.m0, b.m1              // scratch m0, m1
	worldBasis.SetQ(b.world.Rot)                           //
	basisTransposed.Transpose(worldBasis)                  //
	b.iitw.Mult(worldBasis.ScaleV(b.iit), basisTransposed) // scratch m0, m1 free
}

// integrateVelocities updates this bodies linear and angular velocities based
// on the bodies current forces. Static bodies are ignored.
// FUTURE: look up symplectic Euler and see if this is the spot where it
//        should be used (or is already being used).
//             v(t+dt) = v(t) + a(t) * dt
//             x(t+dt) = x(t) + v(t+dt) * dt
func (b *body) integrateVelocities(ts float64) {
	if !b.movable {
		return
	}

	// update linear velocity
	m, v, force := b.imass*ts, b.lvel, b.lfor
	v.X, v.Y, v.Z = v.X+force.X*m, v.Y+force.Y*m, v.Z+force.Z*m

	// update angular velocity
	{ // scratch v0
		torq, a := b.v0, b.avel
		torq.MultMv(b.iitw, b.afor)
		a.X, a.Y, a.Z = a.X+torq.X*ts, a.Y+torq.Y*ts, a.Z+torq.Z*ts
	} // scratch v0 free

	// clamp angular velocity. Collision calculations will fail if its to high.
	avel := b.avel.Len()
	if avel*ts > lin.HalfPi {
		b.avel.Scale(b.avel, lin.HalfPi/ts/avel)
	}
}

// applyDamping adjust linear and angular velocity by their respective
// damping factors.
func (b *body) applyDamping(timestep float64) {
	b.lvel.Scale(b.lvel, math.Pow(1.0-b.ldamp, timestep))
	b.avel.Scale(b.avel, math.Pow(1.0-b.adamp, timestep))
}

// getVelocityInLocalPoint updates vector v to be the linear and angular
// velocity of this body at the given point. The point is expected to be
// in local coordinate space.
func (b *body) getVelocityInLocalPoint(localPoint, v *lin.V3) *lin.V3 {
	return v.Cross(b.avel, localPoint).Add(v, b.lvel)
}

// combinedFriction calculates the combined friction of the two bodies.
// Returned friction value clamped to reasonable range.
func (b *body) combinedFriction(a *body) float64 {
	return lin.Clamp(a.friction*b.friction, -maxFriction, maxFriction)
}

// combinedRestitution calculates the total bounciess of the two
// bodies.
func (b *body) combinedRestitution(a *body) float64 {
	return a.restitution * b.restitution
}

// initSolverBody initializes, and creates if necessary, solver specific
// data structures related to a body. All colliding bodies need solver bodies.
func (b *body) initSolverBody() *solverBody {
	switch {
	case b.sbod == nil && b.movable: // unique to this body.
		b.sbod = newSolverBody(b)
	case b.sbod != nil && b.movable: // reuse existing solver body.
		b.sbod.reset(b)
	case b.sbod == nil && !b.movable: // shared fixed solver body.
		b.sbod = fixedSolverBody()
	}
	return b.sbod
}

// worldAabb updates Abox ab to be the bodies axis-aligned bounding box
// in world coordinates. The updated Abox is returned.
func (b *body) worldAabb(ab *Abox) *Abox { return b.shape.Aabb(b.world, ab, 0) }

// predictedAabb updates Abox ab to be the bodies axis-aligned bounding box
// in the predicted world coordinates.
func (b *body) predictedAabb(ab *Abox, margin float64) *Abox { return b.shape.Aabb(b.guess, ab, margin) }

// updatePredictedTransform provides a guess where the body would appear using
// the current linear and angular velocities within the supplied timestep.
func (b *body) updatePredictedTransform(timestep float64) {
	b.guess.Integrate(b.world, b.lvel, b.avel, timestep)
}

// updateWorldTransform sets the world transform based on the current linear
// and angular velocities. Expected to be called after the solver completes.
func (b *body) updateWorldTransform(timestep float64) {
	b.t0.Integrate(b.world, b.lvel, b.avel, timestep) // scratch t0
	b.world.Set(b.t0)                                 // scratch t0 free
}

// clearForces sets the forces applied to the body back to zero.
func (b *body) clearForces() {
	b.lfor.SetS(0, 0, 0)
	b.afor.SetS(0, 0, 0)
}
