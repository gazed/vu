// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package render

import (
	"math"

	"github.com/gazed/vu/math/lin"
)

// Animation is a sequence of positions (frames) for the joints of a model
// where the position of each model vertex can be affected by up to 4 joints.
// Animation data is independent of any given instance, thus making Animation
// safe to cache and reference by multiple models.
// Animation data may contain more than one animation sequence (movement).
type Animation interface {
	Name() string // Unique identifier set on creation.

	// SetData initializes the animation data that has been processed
	// into a series of frames. SetData is expected to be called once
	// during loading/initialization.
	//    frames:   gives the 3D position of all joints.
	//    joints:   number of joints and their parent joints.
	//    movement: range of frames forming a unique motion.
	SetData(frames []*lin.M4, joints []int32, movements []Movement)
	Movements() []string // Return the names of the movements.
}

// Movement is part of an Animation. It allows multiple animated motions to
// be associated with a single model. Each movement refences a sequence of
// frames from the overall animation.
type Movement struct {
	Name   string // Name of the movement.
	F0, Fn int    // First animation frame, number of animation frames.
}

// =============================================================================

// animation implements Animation. It both contains the animation data and
// the knowledge for transforming animation data into a specific pose that
// is sent to the graphics card.
type animation struct {
	name         string     // unique animation name.
	jointCnt     int        // number of joints.
	frames       []lin.M4   // nFrames*nPoses transform bone positions.
	joints       []int32    // joint parent indicies.
	movements    []Movement // frames where animations start and end.
	mnames       []string   // movement names for easy reference.
	jntM4, tmpM4 *lin.M4    // Per-frame scratch values for playing animations.
}

// newAnimation allocates space for animation data and the data structures
// needed to create intermediate poses on the fly.
func newAnimation(name string) *animation {
	a := &animation{name: name}
	a.jntM4, a.tmpM4 = &lin.M4{}, &lin.M4{}
	return a
}

// Name returns the unique animation identifier.
func (a *animation) Name() string { return a.name }

// SetData is expected to be called once during loading/initialization.
func (a *animation) SetData(frames []*lin.M4, joints []int32, movements []Movement) {
	a.jointCnt = len(joints)
	a.movements = movements
	a.mnames = []string{}
	for _, movement := range movements {
		a.mnames = append(a.mnames, movement.Name)
	}
	a.frames = make([]lin.M4, len(frames)) // transform matrices for each frame.
	for cnt, frame := range frames {
		a.frames[cnt].Set(frame)
	}
	a.joints = a.joints[:0]
	a.joints = append(a.joints, joints...)
}

// Movements allows the user to query the name assigned to each distinct
// animation movement.
func (a *animation) Movements() []string { return a.mnames }

// playMovement returns the movement index if it is valid.
// Otherwise 0 is returned.
func (a *animation) playMovement(index int) int {
	if index >= 0 && index < len(a.movements) {
		return index
	}
	return 0
}
func (a *animation) maxFrames(movement int) int {
	if movement >= 0 && movement < len(a.movements) {
		mv := a.movements[movement]
		return mv.Fn
	}
	return 0
}

// animate combines per model instance information with the animation data
// to produce the unique model pose. The pose data is expected to be updated
// on the graphics card each update cycle.
func (a *animation) animate(dt, frame float64, movement int, pose []m34) float64 {
	if len(a.movements) <= 0 {
		return 0
	}
	mv := a.movements[movement]

	// The frame timer, fcnt, controls the speed of the animation.
	// FUTURE: find a more generic way to time animations.
	frame += dt * 20 // Increment frame timer.
	frame1 := int(math.Floor(frame))
	frame2 := frame1 + 1
	frameoffset := float64(frame) - float64(frame1)
	frame1 = (frame1 % (mv.Fn)) + mv.F0
	frame2 = (frame2 % (mv.Fn)) + mv.F0

	// Interpolate matrixes between the two closest frames and concatenate with
	// parent matrix if necessary. Concatenate the result with the inverse of the
	// base pose. Animation blending and inter-frame blending could be done here.
	for cnt := 0; cnt < a.jointCnt; cnt++ {

		// interpolate between the two closest frames.
		m1, m2 := &a.frames[frame1*a.jointCnt+cnt], &a.frames[frame2*a.jointCnt+cnt]
		a.jntM4.Set(m1).Scale(1-frameoffset).Add(a.jntM4, a.tmpM4.Set(m2).Scale(frameoffset))
		if a.joints[cnt] >= 0 {

			// parentPose * childPose * childInverseBasePose
			a.jntM4.Mult(a.jntM4, (&pose[a.joints[cnt]]).toM4(a.tmpM4))
		}
		(&pose[cnt]).tom34(a.jntM4)
	}
	return frame // give back the updated frame position.
}
