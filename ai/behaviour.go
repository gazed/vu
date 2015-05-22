// Copyright © 2014-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package ai

// Behaviour: controls actions for a period of time. These are provided by
// the application and are the building blocks of a behaviour tree. Status
// is maintained by the behaviour as follows:
//     • Uninitialized behaviours should have status INVALID.
//     • Initialized and uncompleted behaviours should have status RUNNING.
//     • Completed behaviours should have status FAILURE or SUCCESS.
type Behaviour interface {
	Status() (status BehaviourState) // Current behaviour status.
	Update() (status BehaviourState) // Update behaviour status.
	Init()                           // Called once on first behaviour update.
	Reset()                          // Set status to INVALID.

	// Observer listens for completed behaviours.
	Observer() (bo BehaviourObserver) // Gets and
	SetObserver(bo BehaviourObserver) // ...sets behaviour observer.
}

// BehaviourBase holds the two fields that every behaviour needs and is
// expected to be embedded in every Behaviour implementation. Eg.
//     type someBehaviour struct {
//        BehaviourBase
//        // some behaviour specific fields.
//     }
type BehaviourBase struct {
	State BehaviourState    // Needed to trigger Init.
	Obs   BehaviourObserver // Observer for this behaviour.
}

// Status returns the behaviour state.
func (bb *BehaviourBase) Status() (status BehaviourState) { return bb.State }

// Observer returns the behaviour observer. Nil is returned if there
// is no current observer.
func (bb *BehaviourBase) Observer() (bo BehaviourObserver) { return bb.Obs }

// SetObserver sets the behaviour observer. Use nil to clear the observer.
func (bb *BehaviourBase) SetObserver(bo BehaviourObserver) { bb.Obs = bo }

// Reset sets the state to INVALID. Allows behaviours to be reset and reused.
func (bb *BehaviourBase) Reset() { bb.State = INVALID }

// =============================================================================

// BehaviourState is a custom type for behaviour state constants.
type BehaviourState int

// BehaviourState values.
const (
	INVALID BehaviourState = iota // Behaviour is not initialized.
	SUCCESS                       // Behaviour has succeeded.
	FAILURE                       // Behaviour has failed.
	RUNNING                       // Behaviour is still processing.
)

// =============================================================================

// BehaviourObserver is the listener interface for behaviour completion.
// It is injected in the BehaviourTree.Start method.
type BehaviourObserver interface {
	Complete(b Behaviour) // Called when behaviour completes.
}
