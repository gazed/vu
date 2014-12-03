// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"testing"
)

// Check particle creation where the number of new particles
// is less than delta-time.
func TestSlowEmit(t *testing.T) {
	e := newEmitter(10)
	total := 0
	for cnt := 0; cnt < 100; cnt++ {
		particles := e.emit(0.02) // 100*0.02 = 2 seconds
		total += len(particles)
	}
	if total != 20 {
		t.Errorf("Expecting 20 particles over 2 seconds. Got %d", total)
	}
}

// Check particle creation where the number of new particles
// is greater than delta-time.
func TestFastEmit(t *testing.T) {
	e := newEmitter(10000)
	particles := e.emit(0.02)
	total := len(particles)
	if total != 200 {
		t.Errorf("Expecting 200 particles over 0.02 seconds. Got %d", total)
	}
}

// Check particle creation where the number of new particles
// is greater than delta-time.
func TestErrorEmit(t *testing.T) {
	e := newEmitter(10000)
	particles := e.emit(2) // calls should be 1 second or less apart.
	total := len(particles)
	if total != 10000 {
		t.Errorf("Expecting 10000 particles over 2 seconds. Got %d", total)
	}
}

// Check that enough space is allocated for particle locations.
func TestEffectSize(t *testing.T) {
	e := newEffect(100, 10, nil)
	if len(e.pb) != 100*3 {
		t.Errorf("Expecting 300 floats, got %d", len(e.pb))
	}
}

// Check that a testMover can be assigned and used.
func TestEffectMover(t *testing.T) {
	e := newEffect(100, 50, testMover) // 50 particles a second means 1 per update.
	e.Update(nil, 0.02)
	e.Update(nil, 0.02)
	if len(e.active) != 2 || e.active[0].Life != 0.04 {
		t.Errorf("Expected two particle, got %d, Expected life at 0.04 got %f",
			len(e.active), e.active[0].Life)
	}
}

func testMover(particles []*EffectParticle, dt float64) []*EffectParticle {
	for _, p := range particles {
		p.Life += dt
	}
	return particles
}
