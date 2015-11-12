// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package physics

import (
	"testing"

	"github.com/gazed/vu/math/lin"
)

// Check that each pair has a unique id.
func TestUuid(t *testing.T) {
	b0, b1 := newBody(NewSphere(1)), newBody(NewSphere(1))
	if b1.bid-b0.bid != 1 {
		t.Error("Body id's should be incrementing")
	}
}

func TestSphereProperties(t *testing.T) {
	b := newBody(NewSphere(1)).SetMaterial(0.5, 0.8).(*body)
	if b.movable != true || !lin.Aeq(b.imass, 2) {
		t.Errorf("Expecting movable body with mass %f", b.imass)
	}
	if dumpV3(b.iit) != "{5.0 5.0 5.0}" {
		t.Errorf("Expecting initial inverse inertia %s", dumpV3(b.iit))
	}
}
func TestBoxProperties(t *testing.T) {
	b := newBody(NewBox(100, 1, 100)).SetMaterial(0, 0.1).(*body)
	if b.movable == true || b.imass != 0.0 {
		t.Errorf("Expecting stationary body with no mass.")
	}
}
func TestApplyGravity(t *testing.T) {
	b := newBody(NewSphere(1)).SetMaterial(0.5, 0.8).(*body)
	want := "{0.0 10.0 0.0}"
	if b.applyGravity(10); dumpV3(b.lfor) != want {
		t.Errorf("Expected forces %s, got %s", want, dumpV3(b.lfor))
	}
}
func TestUpdateInertiaTensor(t *testing.T) {
	b := newBody(NewSphere(1)).SetMaterial(0.5, 0.8).(*body)
	b.updateInertiaTensor()
	want := &lin.M3{
		Xx: 5.0, Xy: 0.0, Xz: 0.0,
		Yx: 0.0, Yy: 5.0, Yz: 0.0,
		Zx: 0.0, Zy: 0.0, Zz: 5.0}
	if dumpM3(b.iitw) != dumpM3(want) {
		t.Errorf("Expecting updated inverse inertia tensor world \n%s", dumpM3(b.iitw))
	}
}
func TestIntegrateVelocities(t *testing.T) {
	b := newBody(NewSphere(1)).SetMaterial(0.5, 0.8).(*body)
	b.lfor.SetS(1, 1, 1)
	b.lvel.SetS(2, 2, 2)
	if b.integrateVelocities(0.2); dumpV3(b.lvel) != "{2.4 2.4 2.4}" {
		t.Errorf("Expecting change in linear velocity %s.", dumpV3(b.lvel))
	}
	b.afor.SetS(1, 1, 1)
	b.avel.SetS(3, 3, 3)
	b.updateInertiaTensor()
	if b.integrateVelocities(0.2); dumpV3(b.avel) != "{4.0 4.0 4.0}" {
		t.Errorf("Expecting change in angular velocity %s.", dumpV3(b.avel))
	}
}
func TestIntegrateLinearVelocity(t *testing.T) {
	b := newBody(NewSphere(1)).SetMaterial(1, 0.8).(*body)
	b.lfor.SetS(0, -10, 0)
	b.lvel.SetS(0, 0, 0)
	if b.integrateVelocities(0.02); dumpV3(b.lvel) != "{0.0 -0.2 0.0}" {
		t.Errorf("Expecting new linear velocity %s.", dumpV3(b.lvel))
	}
}
func TestApplyDamping(t *testing.T) {
	b := newBody(NewSphere(1)).SetMaterial(0.5, 0.8).(*body)
	b.lvel.SetS(2, 2, 2)
	b.avel.SetS(3, 3, 3)
	b.ldamp, b.adamp = 0.5, 0.5
	b.applyDamping(0.2)
	if dumpV3(b.lvel) != "{1.7 1.7 1.7}" || dumpV3(b.avel) != "{2.6 2.6 2.6}" {
		t.Errorf("Expecting velocity damping %s %s.", dumpV3(b.lvel), dumpV3(b.avel))
	}
}
func TestGetVelocityInLocalPoint(t *testing.T) {
	b := newBody(NewSphere(1)).SetMaterial(0.5, 0.8).(*body)
	b.lvel.SetS(2, 2, 2)
	b.avel.SetS(3, 3, 3)
	v, p, want := lin.NewV3(), lin.NewV3S(1, 1, 1), "{2.0 2.0 2.0}"
	if b.getVelocityInLocalPoint(p, v); dumpV3(v) != want {
		t.Errorf("Expecting local velocity %s, got %s", dumpV3(v), want)
	}
}
func TestUpdatePredictedTransform(t *testing.T) {
	b := newBody(NewSphere(1)).SetMaterial(0.5, 0.8).(*body)
	b.lvel.SetS(2, 2, 2)
	b.avel.SetS(3, 3, 3)
	want := "{0.4 0.4 0.4}{0.3 0.3 0.3 0.9}"
	if b.updatePredictedTransform(0.2); dumpT(b.guess) != want {
		t.Errorf("Expecting transform prediction %s, got %s", want, dumpT(b.guess))
	}
	if dumpT(b.world) != "{0.0 0.0 0.0}{0.0 0.0 0.0 1.0}" {
		t.Errorf("World transform should not have changed %s", dumpT(b.world))
	}
}
func TestUpdateWorldTransform(t *testing.T) {
	b := newBody(NewSphere(1)).SetMaterial(0.5, 0.8).(*body)
	b.lvel.SetS(2, 2, 2)
	b.avel.SetS(3, 3, 3)
	want := "{0.4 0.4 0.4}{0.3 0.3 0.3 0.9}"
	if b.updateWorldTransform(0.2); dumpT(b.world) != want {
		t.Errorf("Expecting transform prediction %s, got %s", want, dumpT(b.world))
	}
}
