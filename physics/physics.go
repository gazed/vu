// Copyright © 2024 Galvanized Logic Inc.

// Package physics is a real-time simulation of real-world physics.
// Physics applies simulated forces to virtual 3D objects known as bodies.
// Physics updates bodies locations and directions based on forces and
// collisions with other bodies.
//
// Package physics is provided as part of the vu (virtual universe) 3D engine.
package physics

// physics.go exposes the physics API needed by the engine.
// Physics was ported from https://github.com/felipeek/raw-physics.
// The go code matches the file and function names of the oriinal code.
// This is done to help debug porting errors.
//	 vu/physics              : raw-physics/src/physics
//	 body.go                 : ../entity.cpp ../entity.h
//	 broad.go                : broad.cpp broad.h
//	 clipping.go             : clipping.cpp clipping.h
//	 collider.go             : collider.cpp collider.h
//	 epa.go                  : epa.cpp epa.h
//	 gjk.go                  : gjk.cpp gjk.h
//	 pbd.go                  : pbd.cpp pbd.h
//	 pbd_base_constraints.go : pbd_base_constraints.cpp pbd_base_constraints.h
//	 physics_util.go         : physics_util.cpp physics_util.h
//	 support.go              : support.cpp support.h

import (
	"github.com/gazed/vu/math/lin"
)

// Simulate physics for the given timestep. This is expected to be called
// regularly from the main engine loop. At the end of a simulation all of
// the bodies positions and orientations will be updated based on forces
// acting upon them and/or collision results. Fixed, or unmoved bodies,
// or bodies with zero mass are not updated.
func Simulate(bods []Body, timestep float64) {
	bodies = bods
	for i := range bodies {
		b := bodies[i]
		colliders_update(b.colliders, b.world_position, &b.world_rotation)
	}
	const GRAVITY float64 = 10.0
	for i := range bodies {
		bod := &bodies[i]
		if bod.fixed {
			continue // don't bother adding force to fixed bodies.
		}
		position := lin.NewV3()
		force := lin.NewV3().SetS(0.0, -GRAVITY*1.0/bod.inverse_mass, 0.0)
		bod.AddForce(*position, *force, false)
	}
	pbd_simulate(timestep, bodies, 1, 1, true)
	for i := range bodies {
		bod := &bodies[i]
		bod.clear_forces()
	}
}

// bodies are set each call to Simulate.
// Body IDs (bids) are exactly the slice indexes and are
// valid for this one simulation run.
var bodies []Body

type bid uint32 // physics "body" id. Max 4 billion bodies.

// Sphere creates a ball shaped physics body located at the origin.
// The sphere size is defined by the radius.
// The sphere can be static (unmovable) or kinematic (moveable).
func NewSphere(radius float64, static bool) *Body {
	sphereCollider := collider_sphere_create(float32(radius))
	colliders := []collider{sphereCollider}

	world_position := lin.NewV3() // app to call body.SetPosition
	world_rotation := lin.NewQ()  // app to call body.SetRotation
	world_scale := lin.NewV3()    // app to call body.SetScale
	mass := 1.0
	static_friction_coefficient := 0.5
	dynamic_friction_coefficient := 0.5
	restitution_coefficient := 0.0
	return body_create_ex(*world_position, *world_rotation, *world_scale, mass, colliders,
		static_friction_coefficient, dynamic_friction_coefficient, restitution_coefficient, static)
}

// Box creates a box shaped physics body located at the origin.
// The box size is given by the half-extents so that actual size
// is w=2*hx, h=2*hy, d=2*hz.
// The box can be static (unmovable) or kinematic (moveable).
func NewBox(hx, hy, hz float64, static bool) *Body {
	// # Blender 4.0.2 Cube OBJ Y-up Z-forward
	vertexes := []lin.V3{
		{-hx, +hy, +hz}, // vertex 0
		{-hx, -hy, +hz}, // vertex 1
		{-hx, +hy, -hz}, // vertex 2
		{-hx, -hy, -hz}, // vertex 3
		{+hx, +hy, +hz}, // vertex 4
		{+hx, -hy, +hz}, // vertex 5
		{+hx, +hy, -hz}, // vertex 6
		{+hx, -hy, -hz}, // vertex 7
	}
	indexes := []uint32{
		4, 2, 0, // top
		4, 6, 2, // top
		2, 7, 3, // back
		2, 6, 7, // back
		6, 5, 7, // right
		6, 4, 5, // right
		1, 7, 5, // bottom
		1, 3, 7, // bottom
		0, 3, 1, // left
		0, 2, 3, // left
		4, 1, 5, // front
		4, 0, 1, // front
	}

	boxCollider := collider_convex_hull_create(vertexes, indexes)
	colliders := []collider{boxCollider}

	world_position := lin.NewV3()                  // app to call body.SetPosition
	world_rotation := lin.NewQ().SetAa(0, 1, 0, 0) // app to call body.SetRotation
	world_scale := lin.NewV3().SetS(1, 1, 1)       // app to call body.SetScale
	mass := 1.0
	static_friction_coefficient := 0.5
	dynamic_friction_coefficient := 0.5
	restitution_coefficient := 0.0
	return body_create_ex(*world_position, *world_rotation, *world_scale, mass, colliders,
		static_friction_coefficient, dynamic_friction_coefficient, restitution_coefficient, static)
}

// TODO convex hull for a give set of vertexes and indexes.
// func NewConvexHull(...) *Body {}

// v2Int is a 2 element integer vector.
type v2Int struct {
	x uint32
	y uint32
}

// v3Int is a 3 element integer vector.
type v3Int struct {
	x uint32
	y uint32
	z uint32
}

// v4Int is a 4 element integer vector.
type v4Int struct {
	x uint32
	y uint32
	z uint32
	w uint32
}
