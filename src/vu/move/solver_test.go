// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package move

import (
	"fmt"
	"testing"
)

func TestNewSolver(t *testing.T) {
	sol := newSolver()
	if sol.info.friction != 0.3 {
		t.Errorf("Solver info must be initialized on startup")
	}
}

// Test that the first contact point produces correct angular velocities.
// The initial values can be setup to correspond to the initial conditions
// from another physics engine. The expected values where generated from
// bullet physics.
func TestBoxHit0(t *testing.T) {

	// create two bodies for the solver.
	slab := newBody(NewBox(50, 50, 50)).setMaterial(0, 0)
	slab.World().Loc.SetS(0, -50, 0)
	slab.updateInertiaTensor()
	box := newBody(NewBox(1, 1, 1)).setMaterial(1, 0)
	box.World().Loc.SetS(-5, 1.388006, -3)
	box.World().Rot.SetS(0.1825742, 0.3651484, 0.5477226, 0.7302967)
	box.lvel.SetS(0, -16.599991, 0)
	box.lfor.SetS(0, -10, 0)
	box.updateInertiaTensor()

	// set up the solver input.
	bodies := map[uint32]*body{0: slab, 1: box}
	points := []*pointOfContact{newPoc()}
	points[0].point.SetS(-5.2, -0.011994, -4)
	points[0].normal.SetS(0, -1, 0)
	points[0].depth = -0.011994
	pair := newContactPair(slab, box)
	pair.mergeContacts(points) // initialize solver info.
	pairs := map[uint64]*contactPair{pair.pid: pair}

	// run the solver once to get updated velocities.
	sol := newSolver()
	sol.solve(bodies, pairs)
	lv, av := box.lvel, box.avel

	// check the linear velocity
	gotlv := fmt.Sprintf("lvel %+.4f %+.4f %+.4f", lv.X, lv.Y, lv.Z)
	wantlv := "lvel +0.5168 -10.1059 +0.0000"
	if gotlv != wantlv {
		t.Errorf("Linv got %s, wanted %s", gotlv, wantlv)
	}

	// check the angular velocity
	gotav := fmt.Sprintf("avel %+.4f %+.4f %+.4f", av.X, av.Y, av.Z)
	wantav := "avel +10.0412 -0.7752 -0.9229"
	if gotav != wantav {
		t.Errorf("Angv got %s, wanted %s", gotav, wantav)
	}

	// check that the transform is updated correctly.
	box.updateWorldTransform(sol.info.timestep)
	bl := box.world.Loc
	gotl := fmt.Sprintf("bloc %f %f %f", bl.X, bl.Y, bl.Z)
	wantl := "bloc -4.989663 1.185889 -3.000000"
	if gotl != wantl {
		t.Errorf("Loc got %s, wanted %s", gotl, wantl)
	}
	qx, qy, qz, qw := box.world.Rot.GetS()
	gotr := fmt.Sprintf("brot %f %f %f %f", qx, qy, qz, qw)
	wantr := "brot 0.253972 0.301044 0.576211 0.716136"
	if gotr != wantr {
		t.Errorf("Rot got %s, wanted %s", gotr, wantr)
	}
}

// Try two contact points.
// The expected values where generated from bullet physics.
func TestBoxHit2(t *testing.T) {

	// create two bodies for the solver.
	slab := newBody(NewBox(50, 50, 50)).setMaterial(0, 0)
	slab.World().Loc.SetS(0, -50, 0)
	slab.updateInertiaTensor()
	box := newBody(NewBox(1, 1, 1)).setMaterial(1, 0)
	box.World().Loc.SetS(-4.966656, 0.913616, -2.962081)
	box.World().Rot.SetS(0.291306, 0.202673, 0.711813, 0.606125)
	box.lvel.SetS(0.575174, -7.106833, 0.947961)
	box.avel.SetS(7.662199, -2.530342, 6.257204)
	box.lfor.SetS(0, -10, 0)
	box.updateInertiaTensor()

	// set up the solver input.
	bodies := map[uint32]*body{0: slab, 1: box}
	points := []*pointOfContact{newPoc(), newPoc()}
	points[0].point.SetS(-4.955563, -0.315041, -1.741308)
	points[0].normal.SetS(0, -1, 0)
	points[0].depth = -0.315041
	points[1].point.SetS(-6.276365, -0.185829, -3.237565)
	points[1].normal.SetS(0, -1, 0)
	points[1].depth = -0.18582
	pair := newContactPair(slab, box)
	pair.mergeContacts(points) // initialize solver info.
	pairs := map[uint64]*contactPair{pair.pid: pair}

	// run the solver once to get updated velocities.
	sol := newSolver()
	sol.solve(bodies, pairs)
	lv, av := box.lvel, box.avel

	// check the linear velocity
	gotlv := fmt.Sprintf("lvel %f %f %f", lv.X, lv.Y, lv.Z)
	wantlv := "lvel 0.538789 0.484830 0.868218"
	if gotlv != wantlv {
		t.Errorf("Linv got %s, wanted %s", gotlv, wantlv)
	}

	// check the angular velocity
	gotav := fmt.Sprintf("avel %f %f %f", av.X, av.Y, av.Z)
	wantav := "avel 0.401297 -0.391900 0.454597"
	if gotav != wantav {
		t.Errorf("Angv got %s, wanted %s", gotav, wantav)
	}
}
