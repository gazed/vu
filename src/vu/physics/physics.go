// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

// Package physics deals with automatically applying simulated forces, real world based or otherwise,
// to virtual 3D objects known as bodies.  Physics will update bodies location, direction, and motion
// based on forces and collisions with other bodies.
//
// Currently the physics engine is incomplete and untested.  The initial attempt cobbled together
// collision algorithms from any and all sources.
//     http://www.realtimerendering.com/intersections.html.
//     http://www.cs.cmu.edu/~baraff/pbm/rigid2.pdf
//     http://www.ode.org (very nice code structure and layout)
//     http://bulletphysics.org
//     ... and others that are referenced thoughout the package.
// However precious little working physics code made the transition from the source material.
// This was entirely the fault of the package author and in no way the fault of the excellent
// source material. Hence... (TODO) this entire package needs to be replaced by something that
// is fully implemented, tested, and benchmarked (ideally by someone that knows, and likes,
// physics :P). Creating bindings, or porting an open source engine are also possibilities.
//
// Package physics is provided as part of the vu (virtual universe) 3D engine.
package physics

import (
	"vu/math/lin"
)

// Physics is the set of rules that simulates forces acting on bodies in motion.
// The normal usage is to simulate real-life conditions like air resistance and
// gravity, or the lack there-of.
type Physics interface {
	AddBody(b Body)             // Add a body to the physics simulation.
	RemBody(b Body)             // Remove a body from the physics simulation.
	Reset()                     // Clear all bodies from the simulation.
	SetGravity(gravity float32) // SetGravity for the physics simulation.
	SetDamping(damping float32) // SetDamping for all physics motion.
	Collide(b Body)             // Run collision for just the given body.

	// Step the physics simulation one tiny bit forward.  This is expected
	// to be called for each update in the main engine loop.  At the end of
	// a simulation step, all the bodies positions will be updated based on
	// forces acting upon them and/or collision results.  Bodies that are
	// free or have zero mass are not updated.
	//
	// Note that the collision information is both generated and consumed
	// by the physics simulation.  Free bodies that have registered for
	// collision updates are the exception.
	Step(step float32)
}

// Body is a single object that is contained within a physics simulation.
// Bodies generally relate to a scene node that is displayed to the user.
// Only bodies that participate in physics should be added to the physics
// simulation.
type Body interface {
	Id() int                                 // Same id as scene graph node.
	SetShape(shp Shape)                      // Need shape to participate in collision.
	SetLocation(loc *lin.V3)                 // Pointer to shared external location.
	SetRotation(rot *lin.Q)                  // Pointer to shared external rotation.
	LinearMomentum() (x, y, z float32)       // Current linear motion.
	SetLinearMomentum(x, y, z float32)       // Apply initial motion.
	AngularMomentum() (x, y, z float32)      // Current angular motion (think of spin).
	SetAngularMomentum(x, y, z float32)      // Apply initial angular motion.
	ResetMomentum()                          // Remove all linear and angular momentum.
	SetResolver(r func(contacts []*Contact)) // Marks a body as free of physics.
}

// Shape allows collision and collision resolution between different shapes.
// Shapes are only expected to be created by the engine/application but only consumed
// by the physics system.
type Shape interface {
	// Shape vs shape collision returns all the points of contact.
	Collide(s Shape) (contacts []*Contact)

	// Resolve collisions under physics control.
	Bounce(s Shape, c []*Contact, mo1, mo2 *motion)

	// Volume, given a density, can be used to calculate mass as:
	//     mass = density * volume.
	// Shapes with no volume (rays, planes) return 0.
	Volume() float32
}

// Physics, Body, Shape interfaces
// ===========================================================================
// physics provides a default Physics implementation.

// New returns the default physics world instance.
func New() Physics {
	phy := &physics{}
	phy.bodies = map[int]Body{}
	phy.m = &mover{}
	phy.f = &forces{100, 100, 9810} // gravity at 9.81 m/sec*sec
	return phy
}

// physics is the default implementation of the Physics interface.
type physics struct {
	bodies map[int]Body // Bodies in the physics simulations.
	f      *forces      // Forces applied to bodies.
	m      *mover       // Logic for applying forces and resolving collisions.
}

// Physics interface implementation.
func (phy *physics) AddBody(b Body)             { phy.bodies[b.Id()] = b }
func (phy *physics) RemBody(b Body)             { delete(phy.bodies, b.Id()) }
func (phy *physics) Reset()                     { phy.bodies = map[int]Body{} }
func (phy *physics) SetGravity(gravity float32) { phy.f.g = gravity }
func (phy *physics) SetDamping(damping float32) { phy.f.ldamp, phy.f.adamp = damping, damping }
func (phy *physics) Step(timeDelta float32)     { phy.m.step(timeDelta, phy.bodies, phy.f) }
func (phy *physics) Collide(b Body) {
	bod, _ := b.(*body)
	phy.m.checkCollisions([]*body{bod}, phy.bodies)
}

// physics
// ===========================================================================
// body provides a default Body implementation.

// body is the default implementation of the Body interface.
type body struct {
	uid  int     // Unique id. Link to the scene graph node.
	sh   Shape   // Body shape type for collisions.
	now  *motion // Movement information updated each simulation step.
	then *motion // Movement information updated each simulation step.

	// Optional resolver for free bodies.  Setting this will mark the
	// body as free from physics control.
	resolve func(contacts []*Contact)
}

// NewBody returns a stationary physics body.
// Note that zero mass objects can't (won't) be moved.
func NewBody(uid int, size, mass float32) Body {
	b := &body{}
	b.uid = uid
	loc := &lin.V3{0, 0, 0}
	dir := &lin.Q{0, 0, 0, 1}
	b.now = newMotion(size, mass, loc, dir)
	b.then = b.now.clone()
	return b
}

// Body interface implementation.
func (b *body) Id() int                                 { return b.uid }
func (b *body) SetShape(sh Shape)                       { b.sh = sh }
func (b *body) SetLocation(loc *lin.V3)                 { b.now.loc = loc }
func (b *body) SetRotation(dir *lin.Q)                  { b.now.dir = dir }
func (b *body) SetResolver(r func(contacts []*Contact)) { b.resolve = r }
func (b *body) LinearMomentum() (x, y, z float32) {
	lm := b.now.linearMomentum()
	return lm.X, lm.Y, lm.Z
}
func (b *body) SetLinearMomentum(x, y, z float32) {
	b.now.setLinearMomentum(&lin.V3{x, y, z})
}
func (b *body) AngularMomentum() (x, y, z float32) {
	am := b.now.angularMomentum()
	return am.X, am.Y, am.Z
}
func (b *body) SetAngularMomentum(x, y, z float32) {
	b.now.setAngularMomentum(&lin.V3{x, y, z})
}
func (b *body) ResetMomentum() {
	b.now = newMotion(1, 1, b.now.loc, b.now.dir)
	b.then = b.now.clone()
}

// getShape returns the current shape with its origin updated if necessary.
func (b *body) getShape() Shape {
	switch b.sh.(type) {
	case *sphere:
		sh, _ := b.sh.(*sphere)
		l := b.now.loc
		return Sphere(l.X, l.Y, l.Z, sh.radius)
	case *abox:
		sh, _ := b.sh.(*abox)
		l := b.now.loc
		return Abox(sh.x1+l.X, sh.y1+l.Y, sh.z1+l.Z, sh.x2+l.X, sh.y2+l.Y, sh.z2+l.Z)
	default:
		return b.sh
	}
}

// body
// ===========================================================================
// provide default shape implementations.

// Sphere is defined by a center location and a radius.
func Sphere(cx, cy, cz, r float32) Shape { return newSphere(cx, cy, cz, r) }

// Plane is defined by a unit normal vector for the plane and the distance of the
// plane from the origin.
func Plane(nx, ny, nz, dx, dy, dz float32) Shape { return newPlane(nx, ny, nz, dx, dy, dz) }

// Abox (axis aligned bounding box) are specified by a bottom (left) corner and
// a top (right) corner where the bottom corner coordinates are less than the top
// corner coordinates.
func Abox(bx, by, bz, tx, ty, tz float32) Shape { return newAbox(bx, by, bz, tx, ty, tz) }

// Ray is defined by a point of origin and a direction vector.
func Ray(ox, oy, oz, dx, dy, dz float32) Shape { return newRay(ox, oy, oz, dx, dy, dz) }
