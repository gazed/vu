// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package move

// caster contains ray casting logic. It is separate from full collision
// tracking and often used to answer the question "what is the user
// clicking on".

import (
	"math"
	"vu/math/lin"
)

// cast is the function prototype for ray casting algorithms. It takes two
// Solids, expecting the first solid to be a ray. It returns the nearest point
// of contact, if any.
//    r   : Ray.
//    s   : Solid Shape
//    hit : true if the ray hit the Solid.
//    xyz : Point of contact when hit is true.
type cast func(r, s Solid) (hit bool, x, y, z float64)

// rayCastAlgorithms holds the algorithms for the supported shapes that
// a ray can be checked against.
var rayCastAlgorithms = map[int]cast{
	PlaneShape:  castRayPlane,
	SphereShape: castRaySphere,
}

// ============================================================================
// ray-plane cast: http://en.wikipedia.org/wiki/Line–plane_intersection

// castRayPlane calculates the point of collision between a ray and
// a plane. The contact point is returned if there is an intersection.
func castRayPlane(a, b Solid) (hit bool, x, y, z float64) {
	aa, bb := a.(*solid), b.(*solid)
	sa, sb := aa.shape.(*ray), bb.shape.(*plane)
	la, lb := aa.world.Loc, bb.world.Loc
	dir := bb.v0.SetS(sa.dx, sa.dy, sa.dz).Unit() // ray direction.
	nrm := aa.v0.SetS(sb.nx, sb.ny, sb.nz).Unit() // plane normal.
	nrm.MultQ(nrm, bb.world.Rot)                  // apply spin to plane normal.
	denom := dir.Dot(nrm)
	if lin.AeqZ(denom) || denom < 0 {
		return false, 0, 0, 0 // plane is behind ray or ray is parallel to plane.
	}

	// calculate the difference from a point on the plane to the ray origin.
	dx, dy, dz := dir.X, dir.Y, dir.Z
	diff := bb.v0.SetS(lb.X-la.X, lb.Y-la.Y, lb.Z-la.Z)
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

// castRaySphere calculates the point of collision between a ray and
// a sphere. The closest contact point is returned if there is an intersection.
func castRaySphere(a, b Solid) (hit bool, x, y, z float64) {
	aa, bb := a.(*solid), b.(*solid)
	sa, sb := aa.shape.(*ray), bb.shape.(*sphere)
	la, lb := aa.world.Loc, bb.world.Loc
	sc := aa.v0.SetS(lb.X-la.X, lb.Y-la.Y, lb.Z-la.Z) // sphere center - ray origin
	rdir := bb.v0.SetS(sa.dx, sa.dy, sa.dz).Unit()    // ray direction.
	d0 := sc.Dot(rdir)
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
