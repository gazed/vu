// Copyright © 2015-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// state.go exposes the engine state needed by applications.

// State is used to communicate current engine wide variable settings.
// It is refreshed each update and provided to the application.
// Changing state is done through Eng methods, often Eng.Set().
type State struct {
	X, Y, W, H int     // Window lower left corner and size in pixels.
	R, G, B, A float32 // Background clear color.
	Cursor     bool    // True when cursor is visible.
	CullBacks  bool    // True to set backface culling on.
	Blend      bool    // True for texture blending.
	FullScreen bool    // True when window is full screen.
	Mute       bool    // True when audio is muted.
}

// Screen is a convenience method returning the current window dimensions.
func (s *State) Screen() (x, y, w, h int) { return s.X, s.Y, s.W, s.H }

// Internal convenience methods.
func (s *State) setScreen(x, y, w, h int)    { s.X, s.Y, s.W, s.H = x, y, w, h }
func (s *State) setColor(r, g, b, a float32) { s.R, s.G, s.B, s.A = r, g, b, a }
