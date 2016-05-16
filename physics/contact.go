// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package physics

import (
	"math"

	"github.com/gazed/vu/math/lin"
)

// contactPair contains information about two bodies that are close or
// contacting. The bodies may be overlapping (pre-solver) or in resting contact
// (post-solver). Contacts are created, if necessary, during broad phase,
// checked by narrow phase, and updated by the solver.
type contactPair struct {
	bodyA *body             // (A, 0) partner body.
	bodyB *body             // (B, 1) reference body for normal and point.
	pid   uint64            // Unique pair identifier.
	pocs  []*pointOfContact // The current points of contact.
	valid bool              // Broadphase check for deleted bodies.

	// The following fields are used only by the solver.
	processingLimit float64 // Bodies outside this range are ignored.
	breakingLimit   float64 // Bodies outside this range are not contacting.

	// scratch variables are optimizations that avoid creating/destroying
	// temporary objects that are needed each timestep.
	v0, v1, v2 *lin.V3 // Scratch vectors.
}

// newContactPair creates a contact between two bodies. Expected to be
// used for contacting bodies not already being tracked by physics.
func newContactPair(bodyA, bodyB *body) *contactPair {
	con := &contactPair{}
	con.bodyA, con.bodyB = bodyA, bodyB
	if bodyA != nil && bodyB != nil {
		con.pid = bodyA.pairID(bodyB)
	}
	con.pocs = newManifold() // allocate space for 4 contacts.
	con.pocs = con.pocs[0:0] // start with zero contacts.
	con.breakingLimit = 0.02
	con.processingLimit = lin.Large
	con.v0 = lin.NewV3()
	con.v1 = lin.NewV3()
	con.v2 = lin.NewV3()
	return con
}

// refreshContacts updates the solver information for existing points.
// Any changes to the world transforms are applied to the existing points
// and invalid points are discarded.
//
// based on bullet btPersistentManifold::refreshContactPoints
func (con *contactPair) refreshContacts(wtA, wtB *lin.T) {

	// update the solver contact information using the latest world transforms.
	for _, poc := range con.pocs {
		poc.sp.worldA.AppT(wtA, poc.sp.localA)
		poc.sp.worldB.AppT(wtB, poc.sp.localB)
		{ // scratch v0
			poc.sp.distance = con.v0.Sub(poc.sp.worldA, poc.sp.worldB).Dot(poc.sp.normalWorldB)
		} // scratch v0 free
	}

	// remove invalid points.
	valid, distSqr := 0, 0.0
	for index, poc := range con.pocs {
		if poc.sp.distance > con.breakingLimit {
			// removing invalid contact due to distance (projected on contact-normal direction)
		} else {
			{ // scratch v0, v1, v2
				projection := con.v0.Sub(poc.sp.worldA, con.v1.Scale(poc.sp.normalWorldB, poc.sp.distance))
				distSqr = con.v2.Sub(poc.sp.worldB, projection).LenSqr()
			} // scratch v0, v1, v2 free
			if distSqr > con.breakingLimit*con.breakingLimit {
				// removing invalid contact since relative movement orthogonal to normal exceeds margin
			} else {
				con.pocs[valid].set(con.pocs[index]) // keep valid points.
				valid++
			}
		}
	}
	con.pocs = con.pocs[:valid] // reduce size to only the valid points.
}

// mergeContacts merges the newly discovered contact points with the existing
// contact points. This matters most with shapes that can produce multiple
// contact points, ie. box/box collision.
func (con *contactPair) mergeContacts(points []*pointOfContact) {
	if len(points) > 0 {
		for _, poc := range points {
			poc.prepForSolver(con)
			index := con.closestPoint(poc)
			switch {

			// first try to replace a similar current point with the new point.
			case index >= 0:
				con.pocs[index].set(poc)

			// otherwise add the new point if there is space.
			case len(con.pocs) < 4:
				index = len(con.pocs)
				con.pocs = con.pocs[0 : len(con.pocs)+1]
				con.pocs[index].set(poc)
				con.pocs[index].sp.warmImpulse = 0

			// last resort: replace a point giving the best contact coverage based on area.
			default:
				index := con.largestArea(con.pocs, poc)
				con.pocs[index].set(poc)
			}
		}
	}

	// ... else no new points.
	// Note that the previous points are kept when there are no new points.
	// This is essentially resting contact and the contacting pair needs to be
	// put back through the solver in order to adjust for the force of gravity.
}

// closestPoint will return the manifold index of the closest existing
// point to the given point. Negative 1 is returned if no close points are
// found. The solver contact point information must be initialized before
// calling this method.
//
// Based on bullet btPersistentManifold::getCacheEntry
func (con *contactPair) closestPoint(point *pointOfContact) int {
	shortestDistance := con.breakingLimit * con.breakingLimit
	nearestPointIndex := -1
	{ // scratch v0
		diff := con.v0
		for index, poc := range con.pocs {
			diff.Sub(poc.sp.localA, point.sp.localA)
			diffLenSquared := diff.Dot(diff)
			if diffLenSquared < shortestDistance {
				shortestDistance = diffLenSquared
				nearestPointIndex = index
			}
		}
	} // scratch v0 free
	return nearestPointIndex
}

// largestArea calculates areas using the existing points and the new
// point. The return value is the best index to insert the new point.
// This is expected to be called when there are 4 existing manifold points
// and none of them were close enough to the new contact point.
// Always returns an index from 0-3.
//
// Based on bullet btPersistentManifold::sortCachedPoints
func (con *contactPair) largestArea(existingPoints []*pointOfContact, point *pointOfContact) int {
	var a0, a1, a2, a3 float64
	eps := existingPoints

	// calculate the area's based on four points.
	{ // scratch v0, v1, v2
		a0 = con.area(point.sp.localA, eps[1].sp.localA, eps[2].sp.localA, eps[3].sp.localA)
		a1 = con.area(point.sp.localA, eps[0].sp.localA, eps[2].sp.localA, eps[3].sp.localA)
		a2 = con.area(point.sp.localA, eps[0].sp.localA, eps[1].sp.localA, eps[3].sp.localA)
		a3 = con.area(point.sp.localA, eps[0].sp.localA, eps[1].sp.localA, eps[2].sp.localA)
	} // scratch v0, v1, v2 free
	largestArea := lin.AbsMax(a0, a1, a2, a3)
	return largestArea
}

// area returns the value of the largest area from 3 possible areas created
// from the 4 given points.
//
// Based on bullet btPersistentManifold::calcArea4Points
func (con *contactPair) area(p0, p1, p2, p3 *lin.V3) float64 {
	v0, v1, vx := con.v0, con.v1, con.v2
	l0 := vx.Cross(v0.Sub(p0, p1), v1.Sub(p2, p3)).LenSqr()
	l1 := vx.Cross(v0.Sub(p0, p2), v1.Sub(p1, p3)).LenSqr()
	l2 := vx.Cross(v0.Sub(p0, p3), v1.Sub(p1, p2)).LenSqr()
	return math.Max(math.Max(l0, l1), l2)
}

// contactPair
// ============================================================================
// pointOfContact

// pointOfContact describes one point of contact between two bodies.
// It holds the point of contact, a contact normal vector, and contact depth.
// This is sufficient information to allow the contacting objects to be
// separated if necessary. Additional information is added and used by
// the solver.
//
// The point of contact on the other shape can be calculated by
// (point + normal*depth)
type pointOfContact struct {
	point  *lin.V3      // Point of contact on B in world coordinates.
	normal *lin.V3      // Unit normal point on B in world coordinates.
	depth  float64      // Penetration depth
	sp     *solverPoint // Initialized on creation, used by the solver.

	// Scratch variables are optimizations that avoid creating/destroying
	// temporary objects that are needed each timestep.
	v0 *lin.V3 // scratch vector.
}

// newPoc allocates space for, and returns, a pointOfContact structure.
func newPoc() *pointOfContact {
	poc := &pointOfContact{}
	poc.point = lin.NewV3()
	poc.normal = lin.NewV3()
	poc.sp = newSolverPoint()
	poc.v0 = lin.NewV3()
	return poc
}

// newManifold creates space for up to 4 points of contact.
func newManifold() []*pointOfContact {
	return []*pointOfContact{newPoc(), newPoc(), newPoc(), newPoc()}
}

// prepForSolver creates/updates the contact point data needed by the solver.
// Information from the contactPair and pointOfContact are used to initialize
// the corresponding solver point. Returns true if the point was prepared and
// false if the point was discarded.
func (poc *pointOfContact) prepForSolver(con *contactPair) {
	sp := poc.sp
	sp.distance = poc.depth
	{ // scratch v0
		sp.worldA.Set(poc.point).Add(sp.worldA, poc.v0.Scale(poc.normal, poc.depth))
		sp.localA = con.bodyA.world.Inv(sp.localA.Set(sp.worldA))
	} // scratch v0, v1 free
	sp.worldB.Set(poc.point)
	sp.localB = con.bodyB.world.Inv(sp.localB.Set(poc.point))
	sp.normalWorldB.Set(poc.normal)
	sp.combinedFriction = con.bodyA.combinedFriction(con.bodyB)
	sp.combinedRestitution = con.bodyA.combinedRestitution(con.bodyB)

	// note that sp.lateralFrictionDir is recalculated each time in the
	// solver setup.
}

// set updates poc to have a copy of the given pointOfContact information.
func (poc *pointOfContact) set(cp *pointOfContact) {
	poc.point.Set(cp.point)
	poc.normal.Set(cp.normal)
	poc.depth = cp.depth
	poc.sp.set(cp.sp)
}
