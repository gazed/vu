// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package vu

// react is a means of linking user key sequences to code that updates state.
// Applications can bind reactions to pressed key sequences.
// TODO Move from engine to application? Not a lot of value add.

import (
	"time"
)

// Reaction links an application callback to a given user pressed key sequence.
// The "Do" reaction is triggered when the user pressed key sequence matches
// the reaction id.
type Reaction interface {
	// Name uniquely identifies and describes a reaction.  This is generally
	// action oriented word or phrase like "moveUp".
	Name() string

	// Do is the application hander called when the user triggers the pressed
	// key sequence identified by Name.
	Do()

	// Time is the last time this reaction was triggered.
	Time() time.Time

	// SetTime
	SetTime()
}

// Reaction interface
// ===========================================================================
// reaction - Reaction implementation

// NewReaction creates a Reaction.
func NewReaction(name string, do func()) Reaction {
	return &reaction{name, do, time.Now()}
}

// reaction implements Reaction.
type reaction struct {
	name string
	do   func()
	last time.Time // last time a command key was pressed.
}

// Reaction interface implementation.
func (r *reaction) Name() string    { return r.name }
func (r *reaction) Time() time.Time { return r.last }
func (r *reaction) SetTime()        {}
func (r *reaction) Do() {
	r.last = time.Now()
	r.do()
}

// reaction
// ===========================================================================
// reactOnce

// NewReactOnce (metered reaction) wraps a Reaction with a holdoff timer so that
// the reaction is preformed every so often and not every event loop.
//
// This is because key and mouse events are processed many times a
// second and even a quick key press will generate multiple user pressed key
// sequences.  This works great for movement, but some actions need to be gated
// to ensure they are not spammed.
func NewReactOnce(name string, do func()) Reaction {
	ro := &reactOnce{}
	ro.name = name
	ro.do = do
	ro.last = time.Now()
	ro.holdoff, _ = time.ParseDuration("500ms")
	return ro
}

// reactOnce implements Reaction.
type reactOnce struct {
	reaction
	holdoff time.Duration // time in milliseconds before next command key
}

// Reaction interface implementation.
func (ro *reactOnce) SetTime() { ro.last = time.Now() }
func (ro *reactOnce) Do() {
	if time.Now().After(ro.last.Add(ro.holdoff)) {
		ro.last = time.Now()
		ro.do()
	}
}
