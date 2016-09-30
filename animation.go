// Copyright © 2015-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// animation.go encapsulates the manipulation of animated model data.

import (
	"math"

	"github.com/gazed/vu/math/lin"
)

// Animator is a Model that can have multiple animated actions.
// Actions are indexed from 0.
//
// An animated model is a rigged model animation system. An animation
// is a sequence of positions (frames) for the joints (bones) of a model
// where the position of each vertex can be affected by up to 4 joints.
// Animation data is independent of any given instance, thus making it
// safe to cache and be referenced by multiple models. Animation data
// may contain more than one animation sequence (action).
type Animator interface {
	Animate(action, frame int) bool        // Return true if available.
	Action() (action, frame, maxFrame int) // Current movement info.
	Actions() []string                     // Animation sequence names.
}

// Animator
// =============================================================================
// animation underlyes Animator and is controlled through Model.

// animation contains the animation data and the knowledge for transforming
// animation data into a specific pose that is sent to the graphics card.
type animation struct {
	name     string     // unique animation name.
	tag      aid        // name and type as a number.
	jointCnt int        // number of joints.
	frames   []lin.M4   // nFrames*nPoses transform bone positions.
	joints   []int32    // joint parent indicies.
	moves    []movement // frames where animations start and end.
	mnames   []string   // movement names for easy reference.

	// Per-frame scratch value for playing animations.
	jnt0 *lin.M4 // Reused each update to calculate joint (bone) positions.
	jnt1 *lin.M4 // Ditto.
}

// newAnimation allocates space for animation data and the data structures
// needed to create intermediate poses on the fly.
func newAnimation(name string) *animation {
	a := &animation{name: name, tag: assetID(anm, name)}
	a.jnt0, a.jnt1 = &lin.M4{}, &lin.M4{}
	return a
}

// aid is used to uniquely identify assets.
func (a *animation) aid() aid      { return a.tag }  // hashed type and name.
func (a *animation) label() string { return a.name } // asset name

// setData initializes the animation data that has been processed
// into a series of frames. SetData is expected to be called once
// during loading/initialization.
//    frames  : gives the 3D position of all joints.
//    joints  : number of joints and their parent joints.
//    movement: range of frames forming a unique motion.
func (a *animation) setData(frames []*lin.M4, joints []int32, movements []movement) {
	a.jointCnt = len(joints)
	a.moves = movements
	a.mnames = []string{}
	for _, movement := range a.moves {
		a.mnames = append(a.mnames, movement.name)
	}
	a.frames = make([]lin.M4, len(frames)) // transform matrices for each frame.
	for cnt, frame := range frames {
		a.frames[cnt].Set(frame)
	}
	a.joints = a.joints[:0]
	a.joints = append(a.joints, joints...)
}

// setRate changes the number of frames per second for the given
// animation movement.
//    movement: the affected animation movement, indexed from 0 up.
//    rate    : frames per second. Often 24.
func (a *animation) setRate(movement int, rate float64) {
	if movement >= 0 && movement < len(a.moves) {
		a.moves[movement].rate = rate
	}
}

// moveNames allows the user to query the name assigned
// to each distinct animation movement. The slice index can be
// used as the movement parameter in other methods.
func (a *animation) moveNames() []string { return a.mnames }

// isMovement returns the movement index if it is valid.
// Otherwise 0 is returned.
func (a *animation) isMovement(movement int) int {
	if movement >= 0 && movement < len(a.moves) {
		return movement
	}
	return 0
}

// maxFrames returns the number of frames in the current movement.
// Return 0 for unrecognized movements.
//    movement: the affected animation movement, indexed from 0 up.
func (a *animation) maxFrames(movement int) int {
	if movement >= 0 && movement < len(a.moves) {
		mv := a.moves[movement]
		return mv.fn
	}
	return 0
}

// animate combines per model instance information with the animation data
// to produce the unique model pose. The pose data is expected to be updated
// on the graphics card each update tick.
//    dt      : time since last update. Generally 0.02sec.
//    frame   : the current frame position.
//    movement: the affected animation movement, indexed from 0 up.
//    pose    : interpolated data at the fractional frame position.
// Returns the new fractional frame position.
func (a *animation) animate(dt, frame float64, movement int, pose []lin.M4) float64 {
	if len(a.moves) <= 0 {
		return 0
	}
	mv := a.moves[movement]

	// The frame timer, fcnt, controls the speed of the animation.
	frame1 := int(math.Floor(frame))
	frame2 := frame1 + 1
	frameoffset := frame - float64(frame1)
	frame1 = (frame1 % (mv.fn)) + mv.f0
	frame2 = (frame2 % (mv.fn)) + mv.f0

	// Interpolate matrixes between the two closest frames and concatenate with
	// parent matrix if necessary. Concatenate the result with the inverse of the
	// base pose. FUTURE: blending and inter-frame blending could be done here.
	for cnt := 0; cnt < a.jointCnt; cnt++ {

		// interpolate between the two closest frames.
		m1, m2 := &a.frames[frame1*a.jointCnt+cnt], &a.frames[frame2*a.jointCnt+cnt]
		a.jnt0.Set(m1).Scale(1-frameoffset).Add(a.jnt0, a.jnt1.Set(m2).Scale(frameoffset))
		if a.joints[cnt] >= 0 {

			// parentPose * childPose * childInverseBasePose
			a.jnt0.Mult(a.jnt0, &pose[a.joints[cnt]])
		}
		(&pose[cnt]).Set(a.jnt0)
	}
	return frame + dt*mv.rate // return incremented frame position.
}

// anim
// =============================================================================
// movement

// movement is part of an Animation. It allows multiple animated motions to
// be associated with a single model. Each movement refences a sequence of
// frames from the overall animation.
type movement struct {
	name   string  // Name of the movement.
	f0, fn int     // First animation frame, number of animation frames.
	rate   float64 // Animation frames per second. Often 24.
}
