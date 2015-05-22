// Copyright © 2014-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package ai

import (
	"testing"
)

// Test a success over 2 update cycles.
func TestBehaviourSuccess(t *testing.T) {
	mo.status = INVALID
	bt := NewBehaviourTree()
	b := &mockBehaviour{stopat: 2, finalStatus: SUCCESS}
	bt.Start(b, mo)
	bt.Tick()
	bt.Tick()
	if mo.status != SUCCESS {
		t.Errorf("Expected %d got %d", SUCCESS, mo.status)
	}
}

// Test a failure that takes 3 update cycles.
func TestBehaviourFail(t *testing.T) {
	mo.status = INVALID
	bt := NewBehaviourTree()
	b := &mockBehaviour{stopat: 3, finalStatus: FAILURE}
	bt.Start(b, mo)
	bt.Tick()
	bt.Tick()
	bt.Tick()
	if mo.status != FAILURE {
		t.Errorf("Expected %d got %d", FAILURE, mo.status)
	}
}

// Test resetting a behaviour.
func TestBehaviourReset(t *testing.T) {
	mo.status = INVALID
	bt := NewBehaviourTree()
	b := &mockBehaviour{stopat: 1, finalStatus: SUCCESS}
	bt.Start(b, mo)
	bt.Tick()
	if mo.status != SUCCESS || b.Status() != SUCCESS {
		t.Errorf("Expected %d got %d", SUCCESS, mo.status)
	}
	b.Reset()
	if b.Status() != INVALID {
		t.Errorf("Expected %d got %d", INVALID, b.Status())
	}
}

func TestSequenceSuccess(t *testing.T) {
	mo.status = INVALID
	bt := NewBehaviourTree()
	behaviours := []Behaviour{
		&mockBehaviour{stopat: 1, finalStatus: SUCCESS},
		&mockBehaviour{stopat: 1, finalStatus: SUCCESS},
		&mockBehaviour{stopat: 1, finalStatus: SUCCESS},
	}
	seq := NewSequence(bt, behaviours)
	bt.Start(seq, mo)
	bt.Tick()
	if mo.status != SUCCESS {
		t.Errorf("Expected %d got %d", SUCCESS, mo.status)
	}
}

func TestSequenceFailure(t *testing.T) {
	mo.status = INVALID
	bt := NewBehaviourTree()
	behaviours := []Behaviour{
		&mockBehaviour{stopat: 1, finalStatus: SUCCESS},
		&mockBehaviour{stopat: 1, finalStatus: SUCCESS},
		&mockBehaviour{stopat: 1, finalStatus: FAILURE},
	}
	seq := NewSequence(bt, behaviours)
	bt.Start(seq, mo)
	bt.Tick()
	if mo.status != FAILURE {
		t.Errorf("Expected %d got %d", FAILURE, mo.status)
	}
}

func TestSelectorSuccess(t *testing.T) {
	mo.status = INVALID
	bt := NewBehaviourTree()
	behaviours := []Behaviour{
		&mockBehaviour{stopat: 1, finalStatus: FAILURE},
		&mockBehaviour{stopat: 1, finalStatus: FAILURE},
		&mockBehaviour{stopat: 1, finalStatus: SUCCESS},
	}
	sel := NewSelector(bt, behaviours)
	bt.Start(sel, mo)
	bt.Tick()
	if mo.status != SUCCESS {
		t.Errorf("Expected %d got %d", SUCCESS, mo.status)
	}
}

func TestSelectorFailure(t *testing.T) {
	mo.status = INVALID
	bt := NewBehaviourTree()
	behaviours := []Behaviour{
		&mockBehaviour{stopat: 1, finalStatus: FAILURE},
		&mockBehaviour{stopat: 1, finalStatus: FAILURE},
		&mockBehaviour{stopat: 1, finalStatus: FAILURE},
	}
	sel := NewSelector(bt, behaviours)
	bt.Start(sel, mo)
	bt.Tick()
	if mo.status != FAILURE {
		t.Errorf("Expected %d got %d", FAILURE, mo.status)
	}
}

func TestSequenceSelector(t *testing.T) {
	mo.status = INVALID
	bt := NewBehaviourTree()

	// a sequence that succeeds.
	behaviours := []Behaviour{
		&mockBehaviour{stopat: 1, finalStatus: SUCCESS},
		&mockBehaviour{stopat: 1, finalStatus: SUCCESS},
		&mockBehaviour{stopat: 1, finalStatus: SUCCESS},
	}
	seq := NewSequence(bt, behaviours)

	// a selector that will eventually call the sequence and succeed.
	behaviours = []Behaviour{
		&mockBehaviour{stopat: 1, finalStatus: FAILURE},
		&mockBehaviour{stopat: 1, finalStatus: FAILURE},
		seq, // Will succeed.
	}
	sel := NewSelector(bt, behaviours)

	// process the sequence.
	bt.Start(sel, mo)
	bt.Tick()
	if mo.status != SUCCESS {
		t.Errorf("Expected %d got %d", SUCCESS, mo.status)
	}

	// test that reset works.
	sel.Reset()
	if sel.Status() != INVALID || seq.Status() != INVALID {
		t.Errorf("Expected %d got %d", INVALID, sel.Status())
	}
}

// =============================================================================
// Utility methods.

var mo = &mockObserver{}

type mockObserver struct{ status BehaviourState }

func (mo *mockObserver) Complete(b Behaviour) { mo.status = b.Status() }

// =============================================================================

// mockBehaviour is a simple behaviour that waits
// a few update cycles before finishing.
type mockBehaviour struct {
	BehaviourBase
	finalStatus BehaviourState // return this status on completion.
	stopat      int            // stop when counter reaches here.
	counter     int            // counts the update cycles.
}

func (mb *mockBehaviour) Init() { mb.State = RUNNING }
func (mb *mockBehaviour) Update() (status BehaviourState) {
	mb.counter++
	if mb.counter >= mb.stopat {
		mb.State = mb.finalStatus
	}
	return mb.State
}
