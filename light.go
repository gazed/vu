// Copyright © 2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"github.com/gazed/vu/render"
)

// Light is attached to a Pov to give it a position in world space.
// It is used by shaders to interact with a models material values.
// Light is defaulted to white 1, 1, 1. Valid r, g, b colour values
// are between 0 and 1.
type Light interface {
	Colour() (r, g, b float64)       // Get light colour.
	SetColour(r, g, b float64) Light // Set light colour.
}

// Light
// =============================================================================
// light implements Light.

// light is used to set shader uniform values.
// Primarly shaders that care about lighting.
type light struct {
	r, g, b float64 // light colour.
}

// newLight creates a white light.
func newLight() *light {
	l := &light{r: 1, g: 1, b: 1}
	return l
}

// Implement Light interface.
func (l *light) Colour() (r, g, b float64) { return l.r, l.g, l.b }
func (l *light) SetColour(r, g, b float64) Light {
	l.r, l.g, l.b = r, g, b
	return l
}

// toDraw sets all the data references and uniform data needed
// by the rendering layer.
func (l *light) toDraw(d render.Draw, px, py, pz float64) {
	d.SetFloats("l", float32(px), float32(py), float32(pz), 1)
	d.SetFloats("ld", float32(l.r), float32(l.g), float32(l.b))
}
