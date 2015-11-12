// Copyright © 2014-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package physics

// caster contains ray casting logic. It is separate from full collision
// tracking and often used to answer the question "what is the user
// clicking on?".

import (
	"math"

	"github.com/gazed/vu/math/lin"
)

// cast is the function prototype for ray casting algorithms. It takes two
// Solids, expecting the first solid to be a ray. It returns the nearest point
// of contact, if any.
//    r   : Ray.
//    f   : Form (shape + location and orientation).
//    hit : true if the ray hit the Geometry.
//    xyz : Point of contact when hit is true.
type cast func(r, f Body) (hit bool, x, y, z float64)

// rayCastAlgorithms holds the algorithms for the supported shapes that
// a ray can be checked against.
//
// FUTURE: Look into adding support for boxes.
var rayCastAlgorithms = map[int]cast{
	PlaneShape:  castRayPlane,
	SphereShape: castRaySphere,
}

// ============================================================================
// ray-plane cast: http://en.wikipedia.org/wiki/Line–plane_intersection

// castRayPlane calculates the point of collision between ray:a and
// plane:b. The contact point is returned if there is an intersection.
func castRayPlane(a, b Body) (hit bool, x, y, z float64) {
	sa, sb := a.Shape().(*ray), b.Shape().(*plane)
	la, lb := a.World().Loc, b.World().Loc
	rdir := b.(*body).v0.SetS(sa.dx, sa.dy, sa.dz).Unit() // ray direction.
	nrm := a.(*body).v0.SetS(sb.nx, sb.ny, sb.nz).Unit()  // plane normal.
	nrm.MultQ(nrm, b.World().Rot)                         // apply world spin to plane normal.
	denom := rdir.Dot(nrm)
	if lin.AeqZ(denom) || denom < 0 {
		return false, 0, 0, 0 // plane is behind ray or ray is parallel to plane.
	}

	// calculate the difference from a point on the plane to the ray origin.
	dx, dy, dz := rdir.X, rdir.Y, rdir.Z
	diff := b.(*body).v0.SetS(lb.X-la.X, lb.Y-la.Y, lb.Z-la.Z)
	dlen := diff.Dot(nrm) / denom
	if dlen < 0 {
		return false, 0, 0, 0
	}

	// Get contact point by scaling the ray direction with the contact distance
	// and adding the ray origin.
	x, y, z = dx*dlen+la.X, dy*dlen+la.Y, dz*dlen+la.Z
	return true, x, y, z
}

// ============================================================================
// ray-sphere cast: http://en.wikipedia.org/wiki/Line–sphere_intersection
// http://www.scratchapixel.com/lessons/3d-basic-lessons/lesson-7-intersecting-simple-shapes/ray-sphere-intersection/

// castRaySphere calculates the point of collision between ray:a and
// sphere:b. The closest contact point is returned if there is an intersection.
func castRaySphere(a, b Body) (hit bool, x, y, z float64) {
	sa, sb := a.Shape().(*ray), b.Shape().(*sphere)
	la, lb := a.World().Loc, b.World().Loc
	sc := a.(*body).v0.SetS(lb.X-la.X, lb.Y-la.Y, lb.Z-la.Z) // sphere center - ray origin
	rdir := b.(*body).v0.SetS(sa.dx, sa.dy, sa.dz).Unit()    // ray direction.
	d0 := rdir.Dot(sc)
	if d0 < 0 {
		return false, 0, 0, 0 // no solutions
	}
	radius2 := sb.R * sb.R
	d1 := sc.Dot(sc) - d0*d0
	if d1 > radius2 {
		return false, 0, 0, 0 // no solutions
	}
	dlen := d0 - math.Sqrt(radius2-d1)

	// Get contact point by scaling the ray direction with the contact distance
	// and adding the ray origin.
	x, y, z = rdir.X*dlen+la.X, rdir.Y*dlen+la.Y, rdir.Z*dlen+la.Z
	return true, x, y, z
}

// ============================================================================
//
// FUTURE:
// https://truesculpt.googlecode.com/hg-history/Release%25200.8/Doc/ray_box_intersect.pdf
// http://www.scratchapixel.com/lessons/3d-basic-lessons/lesson-7-intersecting-simple-shapes/ray-box-intersection/
