// Copyright © 2017-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// ent.go holds types and methods common to all entities.

import (
	"log"
)

// Ent is an application created entity. Common application entities
// are scenes and models. For example:
//    scene := eng.AddScene()
//    part := scene.AddPart().MakeModel("uv", "msh:icon", "tex:goal")
// Scenes are used to group models with cameras. Models are used to
// group assets with shaders. Together they are used to render frames.
//
// Entities have methods that expect one or more components (data) to have
// been established by the application before being called. For example a
// label can only be Typeset once the entities label data has been created
// using MakeLabel, eg:
//    banner := scene.AddPart().MakeLabel("txt", font)
//    banner.Typeset("Floating Text") // Typeset implies MakeLabel.
// Methods that work with all entities are: Dispose, Exists.
// Other Ent methods will generate a log for missing component dependencies
// when applied to the wrong entity.
//
// Components are different types of data that can be associated with an entity.
// Entity methods, listed below, work with specific components as follows:
//
// Scene : Eng.AddScene creates a scene entity with a Camera.
//         Cam, Is2DSetCuller, SetOrtho, SetOver, SetUI.
// Part  : AddPart creates a new entity with a point-of-view
// and a scene graph node.
//         AddPart, At, SetAt, World, Move, Cull, Culled,
//         View, SetView, SetAa, Spin, SetSpin, Scale, SetScale,
// Model : MakeModel attaches model data to a part entity.
//         MakeModel, Load, Mesh, GenMesh, Tex, GenTex, SetTex, SetFirst,
//         SetUniform, Alpha, SetAlpha, SetColor, SetDraw, Clamp.
// Actor : MakeActor attaches an animated model with a part entity.
//          MakeActor, Animate, Action, Actions, Pose.
// Label : MakeLabel attaches a string model with a part entity.
//         MakeLabel, Typeset, SetWrap, Size.
// Body  : MakeBody attaches a physics body with a part entity.
//         MakeBody, Body, DisposeBody, SetSolid, Cast, Push.
// Light : MakeLight attaches light data to a part entity.
//         MakeLight, SetLightColor.
// Sound : Eng.AddSound creates a new sound entity.
//         PlaySound, SetListener.
// Shadow: SetShadows adds shadows to an existing scene entity.
//         SetShadows.
// Target: AsTex controls rendering a scene to a texture for an
// existing scene entity.
//         AsTex.
type Ent struct {
	eid eid          // Unique entity identifier.
	app *application // Manager of all component managers.
}

// Dispose all components for this entity.
// If the entity is a scene or a part with child entities,
// then all child entities are also disposed.
func (e *Ent) Dispose() {
	e.app.dispose(e.eid)
}

// Exists returns true if the entity has been created and
// not yet disposed.
func (e *Ent) Exists() bool {
	return e.app.eids.valid(e.eid)
}

// Ent
// =============================================================================
// eid defines entity identfiers. See aid.go for asset identifiers.

// eid is an entity identfier comprised of an id used as a live reference
// to data and an edition used to track when ids are deleted and reused.
// Ent ids are expected to be used as array indicies for component
// data and as such they will not change values over their lifetime.
type eid uint32

// Divide the entity bits into a index id and an edition. The edition
// bits are used to track when an entity has been deleted.
const idBits = 20                    // entity array index : 1048575
const edBits = 12                    // entity edition     :    4096
const maxEntID = (1 << idBits) - 1   // mask and max active entities.
const maxEdition = (1 << edBits) - 1 // mask and max dispose and reuse.

// id is the value to be used for array lookups.
func (e eid) id() uint32 { return uint32(e & maxEntID) }

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
	// Max size is maxEntID.
	editions []uint16 // track currently used entities.

	// Starts empty and grows as entities are disposed.
	// New entities are allocated from here once it reaches maxFree.
	free []uint32 // tracks entities ready for reuse.
}

// newEids creates and returns a new entity id manager.
func newEids() *eids {
	return &eids{editions: []uint16{}, free: []uint32{}}
}

// maxFree starts recyling ids once the amount of disposed ids
// reaches the given size.
const maxFree = (1 << (edBits - 1)) // recycling when free reaches 2048.

// create returns a new entity id starting at 1.
// Zero is returned when all entity identifiers have been allocated.
func (ids *eids) create() eid {
	id := uint32(0)
	if len(ids.free) > maxFree {
		id = ids.free[0]
		ids.free = append(ids.free[:0], ids.free[1:]...)
	} else {
		ids.editions = append(ids.editions, 0)
		if id = uint32(len(ids.editions)); id >= maxEntID {

			// entity ids exhausted if nothing in free list.
			if len(ids.free) == 0 {
				log.Printf("All %d entity identifiers in use", maxEntID+1)
				return 0 // design error to be caught during development.
			}
			id = ids.free[0]
			ids.free = append(ids.free[:0], ids.free[1:]...)
		}
	}
	return eid(id | uint32(ids.editions[id-1])<<idBits)
}

// valid entities are those that have been created and not yet disposed.
func (ids *eids) valid(e eid) bool {
	id := e.id()
	if id >= uint32(len(ids.editions)) {
		return false
	}
	return ids.editions[e.id()-1] == e.edition()
}

// dispose marks an entity as no longer valid. The entity identifer
// is placed for reallocation at a later date. The entity can be reallocated
// maxEdition times before it duplicates a previously generated entity.
func (ids *eids) dispose(e eid) {
	id := e.id()
	ids.editions[id-1]++            // mark this entity id as invalid.
	ids.free = append(ids.free, id) // queue it up for reallocation.
}
