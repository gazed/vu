// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package ai

// Resources for understanding behaviour trees:
//    http://en.wikipedia.org/wiki/Behavior_Trees
//    http://www.indiedb.com/groups/indievault/tutorials/game-ai-behavior-tree
//    http://aigamedev.com/broadcasts/session-second-generation-bt/
//    http://www.gamasutra.com/blogs/ChrisSimpson/20140717/221339/Behavior_trees_for_AI_How_they_work.php

import (
	"container/list"
	"log"
)

// BehaviourTree processes behaviours. Multiple behaviours may be started
// where each is a tree of behaviours composed of Sequences and/or Selectors.
type BehaviourTree interface {

	// Start processing behaviour b and associate its completion status
	// with behaviour observer bo.
	Start(b Behaviour, bo BehaviourObserver) // Process a behaviour.

	// Stop informs the given behaviours observer of completion without
	// waiting for the next update tick. The behaviour itself will be
	// removed next update tick.
	Stop(b Behaviour) // Stop processing a behaviour

	// Tick updates each active behaviour. Completed behaviours
	// send notifications through their observers.
	Tick() // Expected to be called each regular update cycle.
}

// NewBehaviourTree creates an empty behaviour tree. It must be initialized
// with behaviours using the Start method and then updated regularly using
// the Tick method.
func NewBehaviourTree() BehaviourTree {
	return &behaviourTree{behaviours: list.New()}
}

// =============================================================================

// behaviourTree implements BehaviourTree.
type behaviourTree struct {
	behaviours *list.List
}

// Start pushes a behaviour onto the processing list and associates
// it with the given observer.
func (bt *behaviourTree) Start(b Behaviour, bo BehaviourObserver) {
	if bo != nil {
		b.SetObserver(bo)
	}
	bt.behaviours.PushFront(b)
}

// Stop immediately propogrates completion status to the parent behaviour
// observer.  Completed behaviours will be removed from the tree next tick.
func (bt *behaviourTree) Stop(b Behaviour) {
	status := b.Status()
	if status != FAILURE && status != SUCCESS {
		log.Printf("behaviourTree.Stop: status must be FAILURE or SUCCESS %d.", status)
	}

	// Inform the behaviour observer of the completion.
	if b.Observer() != nil {
		b.Observer().Complete(b)
	}
}

// Tick updates all active behaviours.
func (bt *behaviourTree) Tick() {
	bt.behaviours.PushBack(nil) // Nil marker for this update tick.
	for bt.step() {
		// Step wll process behaviours until finding the
		// nil marker inserted above.
	}
}

// step processes each of the currently queued behaviours and stops when
// it hits the nil marker.
func (bt *behaviourTree) step() bool {
	elem := bt.behaviours.Front()
	bt.behaviours.Remove(elem)
	if elem.Value == nil {
		return false // Found the nil marker. This update tick is done.
	}
	behaviour := elem.Value.(Behaviour)
	tick(behaviour)
	if behaviour.Status() != RUNNING && behaviour.Observer() != nil {
		behaviour.Observer().Complete(behaviour)
	} else {

		// Still running to put it back in the queue for processing next tick.
		// It will be behind the nil marker.
		bt.behaviours.PushBack(behaviour)
	}
	return true
}

// =============================================================================

// tick ensures that a behaviour is properly initialized, updated, and
// closed down when complete. It is used by the behaviourTree to update
// a behaviour.
func tick(b Behaviour) {
	if b.Status() == INVALID {
		b.Init()
	}
	b.Update()
}

// =============================================================================
// sequence is a Behaviour.

// NewSequence creates a Behaviour and adds it to the BehaviourTree. A sequence
// runs its list of behaviours until one fails, in which case the sequence fails.
// The sequence succeeds if all of its behaviours succeed. The list of sequence
// behaviours is processed from lowest index to highest index.
func NewSequence(bt BehaviourTree, behaviours []Behaviour) Behaviour {
	return &sequence{bt: bt, behaviours: behaviours}
}

// sequence implements a sequence Behaviour.
type sequence struct {
	BehaviourBase
	bt         BehaviourTree // Injected on creation.
	behaviours []Behaviour   // Ordered Behaviours, 0 first.
	current    int           // Currently processing behaviour.
}

// A sequence is running while it is processing its child behaviours.
func (seq *sequence) Init() {
	seq.State = RUNNING
	if len(seq.behaviours) > 0 {
		seq.current = 0
		behaviour := seq.behaviours[seq.current]
		seq.bt.Start(behaviour, seq)
	}
}
func (seq *sequence) Update() (status BehaviourState) { return seq.State }
func (seq *sequence) Reset() {
	seq.State = INVALID
	for _, b := range seq.behaviours {
		b.Reset()
	}
}

// Complete handles child completion through the BehaviourObserver interface.
// Either the sequence completes or the next child begins processing.
func (seq *sequence) Complete(b Behaviour) {
	if b.Status() == FAILURE {
		seq.State = FAILURE
		seq.bt.Stop(seq)
		return
	}
	if b.Status() != SUCCESS { // Completion means FAILURE or SUCCESS.
		log.Printf("sequence.Complete: invalid completion status %d", b.Status())
	}
	if len(seq.behaviours) <= seq.current+1 {
		seq.State = SUCCESS
		seq.bt.Stop(seq)
		return
	}
	seq.current++ // Process next behaviour.
	seq.bt.Start(seq.behaviours[seq.current], seq)
}

// =============================================================================
// selector is a Behaviour.

// NewSelector creates a Behaviour and adds it to the BehaviourTree. A selector
// runs its list of behaviours until one succeeds, in which case the selector
// succeeds. The selector fails if all of its behaviours fail. The list of
// selector behaviours is processed from lowest index to highest index.
func NewSelector(bt BehaviourTree, behaviours []Behaviour) Behaviour {
	return &selector{bt: bt, behaviours: behaviours}
}

// selector implements a selector Behaviour.
type selector struct {
	BehaviourBase
	bt         BehaviourTree // Injected on creation.
	behaviours []Behaviour   // Behaviours in priority order, 0 first.
	current    int           // Currently processing behaviour.
}

// A selector is running while it is processing its child behaviours.
func (sel *selector) Init() {
	sel.State = RUNNING
	if len(sel.behaviours) > 0 {
		sel.current = 0
		behaviour := sel.behaviours[sel.current]
		sel.bt.Start(behaviour, sel)
	}
}
func (sel *selector) Update() (status BehaviourState) { return sel.State }
func (sel *selector) Reset() {
	sel.State = INVALID
	for _, b := range sel.behaviours {
		b.Reset()
	}
}

// Complete handles child completion through the BehaviourObserver interface.
// Either the selector completes or the next child begins processing.
func (sel *selector) Complete(b Behaviour) {
	if b.Status() == SUCCESS {
		sel.State = SUCCESS
		sel.bt.Stop(sel)
		return
	}
	if b.Status() != FAILURE { // completion means FAILURE or SUCCESS.
		log.Printf("selector.Complete: invalid completion status %d", b.Status())
	}
	if len(sel.behaviours) <= sel.current+1 {
		sel.State = FAILURE
		sel.bt.Stop(sel)
		return
	}
	sel.current++ // Process next behaviour.
	sel.bt.Start(sel.behaviours[sel.current], sel)
}
