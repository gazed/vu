// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package physics

import (
	"fmt"
	"testing"

	"github.com/gazed/vu/math/lin"
)

// Check unique pair unique ids. Assign fixed body ids for an easy visual check.
func TestPairID(t *testing.T) {
	b0, b1 := newBody(NewSphere(1)), newBody(NewSphere(1))
	b0.bid, b1.bid = 1, 2
	con := newContactPair(b0, b1)
	pid0, pid1 := b0.pairID(b1), b1.pairID(b0)
	if pid0 != 0x100000002 || pid1 != 0x100000002 || con.pid != 0x100000002 {
		t.Error("Pair id's should be the same regardless of body order")
	}
}

func TestClosestPoint(t *testing.T) {
	b0, b1 := newBody(NewBox(0.5, 0.5, 0.5)), newBody(NewBox(1, 1, 1))
	b0.World().Loc.SetS(0, 0, 1.49)
	con := newContactPair(b0, b1)
	_, _, pocs := collideBoxBox(b0, b1, con.pocs[0:4])
	con.pocs = pocs
	for _, cp := range pocs {
		cp.prepForSolver(con)
	}
	cp0, cp2 := pocs[0], pocs[2]
	if con.closestPoint(cp0) != 0 || con.closestPoint(cp2) != 2 {
		t.Errorf("Should have found the matching points %d %d", con.closestPoint(cp0), con.closestPoint(cp2))
	}
}

func TestPrepForSolver(t *testing.T) {
	ball := newBody(NewSphere(1)).SetMaterial(1, 0).(*body)
	ball.World().Loc.SetS(-5, 0.99, -3)
	box := newBody(NewBox(50, 50, 50)).SetMaterial(0, 0).(*body)
	box.World().Loc.SetS(0, -50, 0)
	con := newContactPair(ball, box)
	_, _, con.pocs = collideSphereBox(ball, box, con.pocs[0:4])
	if len(con.pocs) != 1 {
		t.Errorf("Should have a single contact point.")
	}

	// prepare one of the collision points and check the results.
	cp0 := con.pocs[0]
	cp0.prepForSolver(con)
	got := ""
	got += fmt.Sprintf("LocalA %s WorldA %s\n", dumpV3(cp0.sp.localA), dumpV3(cp0.sp.worldA))
	got += fmt.Sprintf("LocalB %s WorldB %s\n", dumpV3(cp0.sp.localB), dumpV3(cp0.sp.worldB))
	got += fmt.Sprintf("NormalB  %s\n", dumpV3(cp0.sp.normalWorldB))
	got += fmt.Sprintf("LatFric  %s\n", dumpV3(cp0.sp.lateralFrictionDir))
	got += fmt.Sprintf("Friction %f\n", cp0.sp.combinedFriction)
	got += fmt.Sprintf("Bounce   %f\n", cp0.sp.combinedRestitution)
	got += fmt.Sprintf("Distance %f\n", cp0.sp.distance)
	want := "" +
		"LocalA {0.0 -1.0 0.0} WorldA {-5.0 -0.0 -3.0}\n" +
		"LocalB {-5.0 50.0 -3.0} WorldB {-5.0 0.0 -3.0}\n" +
		"NormalB  {0.0 1.0 0.0}\n" +
		"LatFric  {0.0 0.0 0.0}\n" +
		"Friction 0.250000\n" +
		"Bounce   0.000000\n" +
		"Distance -0.050000\n"
	if got != want {
		t.Errorf("Got \n%s", got)
	}
}

func TestLargestArea(t *testing.T) {
	con := &contactPair{}
	con.v0, con.v1, con.v2 = lin.NewV3(), lin.NewV3(), lin.NewV3()

	// Existing points: essentially 14,0,+-1, 16,0,+-1
	manifold := newManifold()
	manifold[0].sp.localA.SetS(13.993946, 25.000000, -0.999210) // 14,0,-1
	manifold[1].sp.localA.SetS(14.006243, 25.000000, 0.979937)  // 14,0,1
	manifold[2].sp.localA.SetS(15.989870, 25.000000, 0.996212)  // 16,0,1
	manifold[3].sp.localA.SetS(15.993749, 25.000000, -0.999743) // 16,0,-1

	// new point A should replace existing point 0.
	ptA := newPoc()
	ptA.sp.localA.SetS(14.024626, 25.000000, -1.020002) // 14,0,-1
	if index := con.largestArea(manifold, ptA); index != 0 {
		t.Errorf("Wrong replacement ptA for best contact area %d", index)
	}

	// new point A should replace existing point 1.
	ptB := newPoc()
	ptB.sp.localA.SetS(14.008444, 25.000000, 0.979925) // 14,0,1
	if index := con.largestArea(manifold, ptB); index != 1 {
		t.Errorf("Wrong replacement ptB for best contact area %d", index)
	}
}
