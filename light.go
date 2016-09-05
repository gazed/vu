// Copyright © 2014-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// Light is used by shaders to interact with a models material values.
// Light is attached to a Pov to give it a position in world space.
// Light is defaulted to white 1,1,1. Valid R,G,B color values
// are from 0 to 1.
type Light struct {
	R, G, B float64 // Red, Green, Blue values range from 0 to 1.
}

// newLight creates a white light.
func newLight() *Light { return &Light{R: 1, G: 1, B: 1} }

// SetColor is a convenience method for changing the light color.
func (l *Light) SetColor(r, g, b float64) { l.R, l.G, l.B = r, g, b }
