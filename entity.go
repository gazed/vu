// Copyright © 2017-2024 Galvanized Logic Inc.

package vu

// entity.go provides unique entity ids that track application created resources.
// An entity component system (ECS) has pros and cons - try to use good design
// over blindly following flavor of the day, eg:
//   https://www.gamedev.net/blogs/entry/2265481-oop-is-dead-long-live-oop/
// This particular attempt does not necessarily count as good design,
// but for some reason made sense to the author at the time.

import (
	"log/slog"
)

// Entity track resoures created by the application. Common application
// entities are scenes and models, eg:
//
//	scene := eng.AddScene(vu.Scene3D)
//	model := scene.AddModel("shd:pbr0", "msh:box0", "mat:box0")
//
// Scenes are used to group models with cameras. Models are used to
// group assets with shaders. Together they are used to render frames.
type Entity struct {
	eid eID          // Unique entity identifier.
	app *application // The application manages all entities.
}

// Dispose all components for this entity.
// If the entity is a scene or a part with child entities,
// then all child entities are also disposed.
func (e *Entity) Dispose(eng *Engine) {
	e.app.dispose(eng, e.eid)
}

// Exists returns true if the entity has been created and
// not yet disposed.
func (e *Entity) Exists() bool {
	return e.app.eids.valid(e.eid)
}

// Entity
// =============================================================================
// eID defines entity identfiers.

// eID is an entity identfier comprised of an id used as a live reference
// to data and an edition used to track when ids are deleted and reused.
// Ent ids are expected to be used as array indicies for component
// data and as such they will not change values over their lifetime.
type eID uint32

// Divide the entity bits into a index id and an edition. The edition
// bits are used to track when an entity has been deleted.
const idBits = 20                    // entity array index : max 1048575
const edBits = 12                    // entity edition     : max    4096
const maxEntID = (1 << idBits) - 1   // mask and max active entities.
const maxEdition = (1 << edBits) - 1 // mask and max dispose and reuse.

// id is the value to be used for array lookups.
func (eid eID) id() uint32 { return uint32(eid & maxEntID) }

// edtion returns the value that tracks if the id is valid.
func (eid eID) edition() uint16 { return uint16((eid >> idBits) & maxEdition) }

// entity
// =============================================================================
// entities see:
// http://bitsquid.blogspot.ca/2014/08/building-data-oriented-entity-system.html

// entities handles the creation and deletion of entity identifiers.
// It ensures a limited set of unique identifiers. These identifiers
// are limited so they can be used as indicies into arrays of data.
type entities struct {

	// Starts empty and grows as entities are allocated.
	// Max size is maxEntID.
	editions []uint16 // track currently used entities.

	// Starts empty and grows as entities are disposed.
	// New entities are allocated from here once maxFree is reached.
	free []uint32 // tracks entities ready for reuse.
}

// maxFree starts recyling ids once the amount of disposed ids
// reaches the given size.
const maxFree = (1 << (edBits - 1)) // recycling when free reaches 2048.

// create returns a new entity id starting at 1.
// Returns zero when all entity identifiers have been allocated.
func (ents *entities) create() eID {
	id := uint32(0)
	if len(ents.free) > maxFree {
		id = ents.free[0]
		ents.free = append(ents.free[:0], ents.free[1:]...)
	} else {
		ents.editions = append(ents.editions, 0)
		if id = uint32(len(ents.editions)); id >= maxEntID {

			// entity ids exhausted if nothing in free list.
			if len(ents.free) == 0 {
				slog.Warn("all entity identifiers in use", "max_entities", maxEntID+1)
				return 0 // design error to be caught during development.
			}
			id = ents.free[0]
			ents.free = append(ents.free[:0], ents.free[1:]...)
		}
	}
	return eID(id | uint32(ents.editions[id-1])<<idBits)
}

// valid entities are those that have been created and not yet disposed.
func (ents *entities) valid(e eID) bool {
	id := e.id()
	if id == 0 {
		return false // id zero is never valid - used to track max allocations.
	}
	if id > uint32(len(ents.editions)) {
		return false
	}
	return ents.editions[id-1] == e.edition()
}

// dispose marks an entity as no longer valid. The entity identifer
// is placed for reallocation at a later date. The entity can be reallocated
// maxEdition times before it duplicates a previously generated entity.
func (ents *entities) dispose(e eID) {
	id := e.id()
	ents.editions[id-1]++             // mark this entity id as invalid.
	ents.free = append(ents.free, id) // queue it up for reallocation.
}
