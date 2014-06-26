// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package vu

import (
	"vu/math/lin"
)

// Culler can be attached to a scene in order to reduce the number
// of items sent for rendering.
type Culler interface {
	Cull(sc Scene, p Part) bool
}

// ============================================================================

// NewRadiusCuller returns a scene culler that removes parts in a radius that
// is around the camera. This is used for an overhead view like a minimap
// where everything around the camera must be drawn.
func NewRadiusCuller(r float64) Culler {
	if r < 0 {
		r = 0
	}
	return &radiusCuller{r}
}

// A scene with radius culling needs the scene's parts distance to camera
// to be have been calculated.
type radiusCuller struct {
	radius float64
}

// Culler implmentation.
func (rc *radiusCuller) Cull(s Scene, p Part) bool {
	prt := p.(*part)
	return prt.toc > rc.radius*rc.radius
}

// ============================================================================

// NewFacingCuller returns a scene culler that removes parts in a radius that
// is in front of the camera. This is used for first person view where only
// whats in front of the camera is important.
func NewFacingCuller(r float64) Culler {
	if r < 0 {
		r = 0
	}
	return &facingCuller{r}
}

// A scene with facing culling needs the scene's parts distance to camera
// to be have been calculated.
//
// Cull places the cull area in front of the camera by moving the
// center up radius units in facing direction. Don't move it all the way
// up so that stuff above or below still exists when looking up/down.
type facingCuller struct {
	radius float64
}

// Culler implmentation.
func (fc *facingCuller) Cull(s Scene, p Part) bool {
	prt := p.(*part)
	scn := s.(*scene)
	toc := prt.toc // get the current distance to camera.
	cam := scn.cam

	// project the camera location along the lookat vector.
	fudgeFactor := float64(0.8) // don't move all the way up.
	lookAt := lin.NewQ().SetAa(1, 0, 0, lin.Rad(cam.up))
	lookAt.Mult(cam.dir, lookAt)
	cx, cy, cz := lin.MultSQ(0, 0, -fc.radius*fudgeFactor, lookAt) // distance
	cx, cy, cz = cx+cam.loc.X, cy+cam.loc.Y, cz+cam.loc.Z          // added to location.

	// cull the part if its to far away.
	toc = prt.distanceTo(cx, cy, cz)
	return toc > fc.radius*fc.radius
}
