// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package physics

// // The following block is C code and cgo directvies.
// // It is used to include collision.c code.
//
// #cgo CFLAGS: -lm -std=c99
//
// #include "collision.h"
import "C" // must be located here.

import (
	"log"
	"math"

	"github.com/gazed/vu/math/lin"
)

// collider creates a collision algorithm for each combination of
// basic shapes. These are used during narrowphase detection to
// calculate points of contact between shape primitives.
type collider struct {
	algorithms [][]collide
}

// newCollider initializes the algorithms needed for narrowphase.
//
// FUTURE: Look into adding support for planes and rays.
func newCollider() *collider {
	c := &collider{}
	c.algorithms = make([][]collide, VolumeShapes)
	for cnt := range c.algorithms {
		c.algorithms[cnt] = make([]collide, VolumeShapes)
	}
	c.algorithms[SphereShape][SphereShape] = collideSphereSphere
	c.algorithms[SphereShape][BoxShape] = collideSphereBox
	c.algorithms[BoxShape][SphereShape] = collideBoxSphere
	c.algorithms[BoxShape][BoxShape] = collideBoxBox
	return c
}

// collider
// ============================================================================
// collide

// collide is the function prototype for collision algorithms. It takes two
// shapes and returns the list of contact points between the two shapes.
// An empty list means that there was no contact.
//    a : Body.
//    b : Different body.
//    c : Preallocated point of contact structures to be updated and returned.
type collide func(a, b Body, c []*pointOfContact) (i, j Body, k []*pointOfContact)

// collide
// ============================================================================
// sphere-sphere collision

// collideSphereSphere returns 0 or 1 contact points.
func collideSphereSphere(a, b Body, c []*pointOfContact) (i, j Body, k []*pointOfContact) {
	aa, bb := a.(*body), b.(*body)
	sa, sb := aa.shape.(*sphere), bb.shape.(*sphere)
	la, lb := aa.world.Loc, bb.world.Loc

	// Separation distance between sphere centers in world space.
	dx, dy, dz := la.X-lb.X, la.Y-lb.Y, la.Z-lb.Z
	separation := math.Sqrt(dx*dx + dy*dy + dz*dz)
	if separation > sa.R+sb.R {
		return a, b, c[0:0] // no contact.
	}
	c0 := c[0]
	c0.depth = separation - (sa.R + sb.R) // how much overlap
	c0.normal.SetS(1, 0, 0)               // sphere's have same center
	if separation > lin.Epsilon {         // sphere's have different center
		c0.normal.SetS(dx/separation, dy/separation, dz/separation) // normalize
	}
	c0.point.Scale(c0.normal, sa.R)      // scale unit normal by radius to ...
	c0.point.Add(bb.world.Loc, c0.point) // ... find point of contact on sphere B.
	return a, b, c[0:1]                  // return single contact
}

// sphere-sphere collision
// ============================================================================
// sphere-box collision

// collideSphereBox can handle arbitrarily rotated boxes. It returns
// 0 or 1 points. Collision margins are used so that close enough objects
// are reported as colliding.
//
// Based on bullet physics btSphereBoxCollisionAlgorithm::getSphereDistance
func collideSphereBox(a, b Body, c []*pointOfContact) (i, j Body, k []*pointOfContact) {
	aa, bb := a.(*body), b.(*body)
	sphere, box := aa.shape.(*sphere), bb.shape.(*box)
	scenter := aa.World().Loc
	sradius := sphere.R
	maxContactDistance := 0.1 // contact breaking threshold
	boxMargin := margin

	// Get the box local half extents. Convert sphere's world to the box's local.
	hx, hy, hz := box.Hx, box.Hy, box.Hz
	sx, sy, sz := bb.World().InvS(scenter.X, scenter.Y, scenter.Z)

	// Determine the closest box vertex to the sphere center.
	px, py, pz := sx, sy, sz
	px = math.Min(hx, px)
	px = math.Max(-hx, px)
	py = math.Min(hy, py)
	py = math.Max(-hy, py)
	pz = math.Min(hz, pz)
	pz = math.Max(-hz, pz)

	// use the closest box point to the sphere center as the contact normal
	// (when the box center is outside the sphere)
	intersectionDist := sradius + boxMargin
	contactDist := intersectionDist + maxContactDistance
	nx, ny, nz := sx-px, sy-py, sz-pz

	// No penetration means no collision.
	dsqrd := nx*nx + ny*ny + nz*nz
	if dsqrd > contactDist*contactDist {
		return a, b, c[0:0]
	}

	// Collision occurred, figure out the collision details.
	var distance float64
	if dsqrd <= lin.Epsilon {
		// Handle the sphere center being inside the box. The contact normal is
		// updated to be the normal for the closest box face.
		px, py, pz, nx, ny, nz, distance = sphereBoxPenetration(box, sx, sy, sz)
	} else {
		distance = math.Sqrt(dsqrd)
		nx, ny, nz = nx/distance, ny/distance, nz/distance
	}

	// Apply the box world transform to get back to world space.
	c0 := c[0]
	c0.point.SetS(bb.World().AppS(px+nx*boxMargin, py+ny*boxMargin, pz+nz*boxMargin))
	c0.normal.SetS(bb.World().AppR(nx, ny, nz)) // only need rotation.
	c0.depth = distance - intersectionDist
	return a, b, c[0:1]
}

// sphereBoxPenetration calculates the closest point and normal when the sphere center
// is inside the box. The sphere center is projected onto each of the box faces to find
// the closest.
//
// Based on bullet physics btSphereBoxCollisionAlgorithm::getSpherePenetration
func sphereBoxPenetration(b *box, sx, sy, sz float64) (px, py, pz, nx, ny, nz, depth float64) {
	faceDist := b.Hx - sx
	depth = faceDist
	px, py, pz = b.Hx, sy, sz
	nx, ny, nz = 1, 0, 0
	faceDist = b.Hx + sx
	if faceDist < depth {
		depth = faceDist
		px, py, pz = -b.Hx, sy, sz
		nx, ny, nz = -1, 0, 0
	}
	faceDist = b.Hy - sy
	if faceDist < depth {
		depth = faceDist
		px, py, pz = sx, b.Hy, sz
		nx, ny, nz = 0, 1, 0
	}
	faceDist = b.Hy + sy
	if faceDist < depth {
		depth = faceDist
		px, py, pz = sx, -b.Hy, sz
		nx, ny, nz = 0, -1, 0
	}
	faceDist = b.Hz - sz
	if faceDist < depth {
		depth = faceDist
		px, py, pz = sx, sy, b.Hz
		nx, ny, nz = 0, 0, 1
	}
	faceDist = b.Hz + sz
	if faceDist < depth {
		depth = faceDist
		px, py, pz = sx, sy, -b.Hz
		nx, ny, nz = 0, 0, -1
	}
	depth = -depth // because its inside the box.
	return
}

// collideBoxSphere reverses the collision to be SphereBox.
func collideBoxSphere(a, b Body, c []*pointOfContact) (i, j Body, k []*pointOfContact) {
	return collideSphereBox(b, a, c)
}

// sphere-box collision
// ============================================================================
// box-box collision

// collideBoxBox uses the Separating Axis Test to check for overlap. If there
// is overlap then the axis of least penetration is used as the contact normal.
// For more background see:
//    Real-Time Collision Detection by Christer Ericson. Sections 4.4.1, 5.2.1
//    http://www.jkh.me/files/tutorials/Separating%20Axis%20Theorem%20for%20Oriented%20Bounding%20Boxes.pdf
//    metanetsoftware.com/technique/tutorialA.html
// Up to 4 contact points can be returned.
func collideBoxBox(a, b Body, c []*pointOfContact) (i, j Body, k []*pointOfContact) {
	aa, bb := a.(*body), b.(*body)
	sa, sb := aa.shape.(*box), bb.shape.(*box)

	// Translate box rotation transforms into 4x3 rotation matrix.
	bbi, bbr, m3 := aa.coi, aa.cor, aa.m0
	bbi.orgA[0] = C.btScalar(aa.world.Loc.X)
	bbi.orgA[1] = C.btScalar(aa.world.Loc.Y)
	bbi.orgA[2] = C.btScalar(aa.world.Loc.Z)
	bbi.orgB[0] = C.btScalar(bb.world.Loc.X)
	bbi.orgB[1] = C.btScalar(bb.world.Loc.Y)
	bbi.orgB[2] = C.btScalar(bb.world.Loc.Z)
	bbi.lenA[0] = C.btScalar(sa.Hx + margin)
	bbi.lenA[1] = C.btScalar(sa.Hy + margin)
	bbi.lenA[2] = C.btScalar(sa.Hz + margin)
	bbi.lenB[0] = C.btScalar(sb.Hx + margin)
	bbi.lenB[1] = C.btScalar(sb.Hy + margin)
	bbi.lenB[2] = C.btScalar(sb.Hz + margin)
	m3.SetQ(aa.world.Rot)
	bbi.rotA[0x0], bbi.rotA[0x1], bbi.rotA[0x2] = C.btScalar(m3.Xx), C.btScalar(m3.Xy), C.btScalar(m3.Xz)
	bbi.rotA[0x4], bbi.rotA[0x5], bbi.rotA[0x6] = C.btScalar(m3.Yx), C.btScalar(m3.Yy), C.btScalar(m3.Yz)
	bbi.rotA[0x8], bbi.rotA[0x9], bbi.rotA[0xA] = C.btScalar(m3.Zx), C.btScalar(m3.Zy), C.btScalar(m3.Zz)
	m3.SetQ(bb.world.Rot)
	bbi.rotB[0x0], bbi.rotB[0x1], bbi.rotB[0x2] = C.btScalar(m3.Xx), C.btScalar(m3.Xy), C.btScalar(m3.Xz)
	bbi.rotB[0x4], bbi.rotB[0x5], bbi.rotB[0x6] = C.btScalar(m3.Yx), C.btScalar(m3.Yy), C.btScalar(m3.Yz)
	bbi.rotB[0x8], bbi.rotB[0x9], bbi.rotB[0xA] = C.btScalar(m3.Zx), C.btScalar(m3.Zy), C.btScalar(m3.Zz)
	bbr.ncp, bbr.code = 0, 0
	C.boxBoxClosestPoints(bbi, bbr)

	// Translate the returned c-lang contact information into go-lang Contact information.
	if bbr.code > 0 {
		numContacts := int(bbr.ncp)
		if numContacts < 0 || numContacts > 4 {
			log.Printf("Dev error: should be 0-4 contacts %d.", numContacts)
			numContacts = int(lin.Clamp(0, 4, float64(numContacts)))
		}
		for cnt := 0; cnt < numContacts; cnt++ {
			cc := bbr.bbc[cnt]
			goc := c[cnt]
			goc.depth = float64(cc.d) // depth is 0 for identical centers.
			goc.normal.SetS(float64(cc.n[0]), float64(cc.n[1]), float64(cc.n[2]))
			goc.point.SetS(float64(cc.p[0]), float64(cc.p[1]), float64(cc.p[2]))
		}
		return a, b, c[0:numContacts]
	}
	return a, b, c[0:0] // boxes did not collide.
}

// box-box collision
// ============================================================================
// FUTURE: Implementat collision detection and contact point generation for
//         more complex types. Choose one generic algorithm for convex shapes
//         (anything with faces/edges) and use it for all cases.
//
// FUTURE look at "The Separating Axis Test between Convex Polyhedra" talk
//        as given in the GDC3013 talk by Dirk Gregorius:
//    https://code.google.com/p/box2d/downloads/detail?name=DGregorius_GDC2013.zip&can=2&q=
//        This handles many shapes including tetrahedron, box, convex hull,
//        cylinder, and cone.
//
// FUTURE: improving efficiency by running detection in parallel on the GPU.
