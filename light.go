// Copyright © 2014-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// light.go defines light color.
// FUTURE: handle multiple lights for one scene. Need a shader that
//         incorporates multiple lights into the final color.

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

// =============================================================================

// lights manages all the active Light instances.
// There's not many lights so not much to optimize.
type lights struct {
	data map[eid]*Light // Camera instance data.
}

// newLights creates a light component manager. Expected to be called
// once on startup.
func newLights() *lights { return &lights{data: map[eid]*Light{}} }

// get the light associated with the given entity, returning nil
// if there is no such light.
func (ls *lights) get(id eid) *Light {
	if light, ok := ls.data[id]; ok {
		return light
	}
	return nil
}

// create and associate a light with the given entity. If there
// already is a light for the entity, dont create one and return
// the existing one instead.
func (ls *lights) create(id eid) *Light {
	if l, ok := ls.data[id]; ok {
		return l // Don't allow creating over existing camera.
	}
	l := newLight()
	ls.data[id] = l
	return l
}

// dispose the light associated for the given entity. Do nothing
// if no such light exists.
func (ls *lights) dispose(id eid) { delete(ls.data, id) }
