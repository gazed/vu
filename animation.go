// Copyright © 2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"math"

	"github.com/gazed/vu/math/lin"
)

// animation is a sequence of positions (frames) for the joints of a model
// where the position of each model vertex can be affected by up to 4 joints.
// Animation data is independent of any given instance, thus making Animation
// safe to cache and reference by multiple models.
// Animation data may contain more than one animation sequence (movement).
//
// animation contains the animation data and the knowledge for transforming
// animation data into a specific pose that is sent to the graphics card.
type animation struct {
	name     string     // unique animation name.
	tag      uint64     // name and type as a number.
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
	a := &animation{name: name, tag: anm + stringHash(name)<<32}
	a.jnt0, a.jnt1 = &lin.M4{}, &lin.M4{}
	return a
}

// label, aid, and bid are used to uniquely identify assets.
// Note: aid is the same as bid for CPU local assets.
func (a *animation) label() string { return a.name } // asset name
func (a *animation) aid() uint64   { return a.tag }  // asset type and name.
func (a *animation) bid() uint64   { return a.tag }  // not bound.

// setData initializes the animation data that has been processed
// into a series of frames. SetData is expected to be called once
// during loading/initialization.
//    frames  : gives the 3D position of all joints.
//    joints  : number of joints and their parent joints.
//    movement: range of frames forming a unique motion.
// setData is expected to be called once during loading/initialization.
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
func (a *animation) setRate(movement int, rate float64) {
	if movement >= 0 && movement < len(a.moves) {
		a.moves[movement].rate = rate
	}
}

// moveNames allows the user to query the name
// assigned to each distinct animation movement.
func (a *animation) moveNames() []string { return a.mnames }

// playMovement returns the movement index if it is valid.
// Otherwise 0 is returned.
func (a *animation) playMovement(index int) int {
	if index >= 0 && index < len(a.moves) {
		return index
	}
	return 0
}

// maxFrames returns the number of frames in the current movement.
// Return 0 for unrecognized movements.
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
