// Copyright © 2014-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"github.com/gazed/vu/math/lin"
)

// Cull is attached to a Camera in order to reduce the number
// of items sent for rendering.
type Cull interface {

	// Culled returns true if a model represented by point, px, py, pz
	// should be culled using the given camera.
	Culled(cam Camera, px, py, pz float64) bool
}

// =============================================================================

// NewFrontCull returns a culler that keeps objects in a radius directly
// in front of the camera. Objects behind the camera and far away from
// the camera are culled.
func NewFrontCull(r float64) Cull {
	if r < 0 {
		r = 0
	}
	return &frontCull{radius: r, rr: r * r}
}

// frontCull removes everything that is not in a circular radius
// projected in front of the camera.
type frontCull struct {
	radius float64 // regular radius.
	rr     float64 // radius squared.
}

// Culler implmentation.
// Project the part location back along the lookat vector.
func (fc *frontCull) Culled(cam Camera, px, py, pz float64) bool {
	fudgeFactor := float64(0.8) // don't move all the way up.
	cx, cy, cz := lin.MultSQ(0, 0, -fc.radius*fudgeFactor, cam.Lookat())
	px, py, pz = px-cx, py-cy, pz-cz // move part location back.
	toc := cam.Distance(px, py, pz)  // cull if outside radius.
	return toc > fc.rr
}

// =============================================================================

// NewRadiusCull returns a culler that removes objects outside a given
// radius from the camera. Can be used for show objects around the camera
// for top down minimaps.
func NewRadiusCull(r float64) Cull {
	if r < 0 {
		r = 0
	}
	return &radiusCull{rr: r * r}
}

// radiusCull removes everything far away from the camera.
type radiusCull struct {
	rr float64 // radius squared.
}

// Culler implmentation. True true if the given location is
// within the culler radius of the camera.
func (rc *radiusCull) Culled(cam Camera, px, py, pz float64) bool {
	toc := cam.Distance(px, py, pz)
	return toc > rc.rr
}
