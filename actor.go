// Copyright © 2017-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// actor.go deals with animated models.
// FUTURE: handle animation models with multiple textures.
//         Animation models are currently limited to one texture.

import (
	"log"
	"math"

	"github.com/gazed/vu/math/lin"
)

// MakeActor creates an animated model associated with this Entity.
// Actor is a specialized static Model that combines a mesh with animation
// data and an animation aware shader.
// An animation is a sequence of positions (frames) for the joints (bones) of
// a model where the position of each vertex can be affected by up to 4 joints.
//   shader: animation aware shader eg: "anim".
//   actor : name of the actor animation data file and actor texture.
//           eg: "actor" finds models/actor.anm and images/actor0.png
func (e *Ent) MakeActor(shader, actor string) *Ent {
	e.app.models.createActor(e, "shd:"+shader, "anm:"+actor, "tex:"+actor+"0")
	return e
}

// Animate requests a particular animation for the model.
// Returns true if the requested animation was available.
//
// Depends on initialization with Ent.MakeActor and loaded animation data.
func (e *Ent) Animate(move, frame int) bool {
	actor := e.app.models.getActor(e.eid)
	if actor == nil || actor.anm == nil {
		log.Printf("Animate needs MakeActor %d and loaded data", e.eid)
		return false
	}
	actor.nFrames = actor.anm.maxFrames(move)
	actor.move = actor.anm.isMovement(move)
	if frame < actor.nFrames {
		actor.frame = float64(frame)
	}
	return move == actor.move // true if the requested movement is available.
}

// Action returns the current animation information. Animations consist
// of a number of different movements, each with a number of frames.
//    move    the currently selected animation.
//    frame   the current animation frame.
//    nFrames the number of frames in the selected animation movement.
//
// Depends on initialization with Ent.MakeActor.
func (e *Ent) Action() (move, frame, nFrames int) {
	if actor := e.app.models.getActor(e.eid); actor != nil {
		return actor.move, int(math.Floor(actor.frame + 1)), actor.nFrames
	}
	log.Printf("Action needs MakeActor %d", e.eid)
	return 0, 0, 0
}

// Actions returns the names of the different animations available to
// this Actor. An empty list is returned if the animation data has not
// yet been loaded.
//
// Depends on initialization with Ent.MakeActor and loaded animation data.
func (e *Ent) Actions() []string {
	if actor := e.app.models.getActor(e.eid); actor != nil && actor.anm != nil {
		return actor.anm.moveNames()
	}
	log.Printf("Actions needs MakeActor %d and loaded data", e.eid)
	return []string{}
}

// Pose returns the bone transform, or the identity matrix
// if there was no transform for the model. The returned matrix
// should not be altered. It is intended for transforming points.
//
// Depends on initialization with Ent.MakeActor.
func (e *Ent) Pose(index int) *lin.M4 {
	if a := e.app.models.getActor(e.eid); a != nil && index < len(a.pose) {
		return &a.pose[index]
	}
	log.Printf("Pose needs MakeActor%d", e.eid)
	return &lin.M4{}
}

// actor entity methods.
// =============================================================================
// actor animation control data.

// actor holds extra animation data. This is combined with the
// animation asset each display frame.
type actor struct {
	anm     *animation // Animation asset.
	frame   float64    // Frame counter.
	move    int        // Current animation defaults to 0.
	nFrames int        // Number of frames in the current movement.
	pose    []lin.M4   // Pose refreshed each update.
}
