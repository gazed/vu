// Copyright © 2014-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// light.go holder for lighting related code. Lights currently only
//          have a color.
// FUTURE: handle multiple lights for one scene. Need a shader that
//         incorporates multiple lights into the final color.

import (
	"log"

	"github.com/gazed/vu/render"
)

// MakeLight adds a light component to an entity.
// Light is used by shaders to interact with a models material values.
// Internally light is has a position and orientation in world space.
// Light is defaulted to white 1,1,1. See SetLightColor.
func (e *Ent) MakeLight() *Ent {
	e.app.lights.create(e.eid)
	return e
}

// SetLightColor sets the light component color using
// the r,g,b values from 0-1.
//
// Depends on Ent.MakeLight.
func (e *Ent) SetLightColor(r, g, b float64) *Ent {
	if l := e.app.lights.get(e.eid); l != nil {
		l.r, l.g, l.b = r, g, b
		return e
	}
	log.Printf("SetLightColor needs MakeLight %d", e.eid)
	return e
}

// light entity methods
// =============================================================================
// light data.

// light holds light color for now. Anticipate other future light parameters.
// The lights world position are scatch values updated each render loop.
type light struct {
	r, g, b    float64 // Red, Green, Blue values range from 0 to 1.
	wx, wy, wz float64 // Scratch world position set from pov when rendering.
}

// newLight creates a white light.
func newLight() *light { return &light{r: 1, g: 1, b: 1} }

// SetColor is a convenience method for changing the light color.
func (l *light) SetColor(r, g, b float64) { l.r, l.g, l.b = r, g, b }

// draw turns a light into draw call data needed by the render shaders.
func (l *light) draw(d *render.Draw) {
	d.SetFloats("lp", float32(l.wx), float32(l.wy), float32(l.wz)) // position
	d.SetFloats("lc", float32(l.r), float32(l.g), float32(l.b))    // color
}

// light data.
// =============================================================================
// light component manager.

// lights manages all the active Light instances.
// There's not many lights so not much to optimize.
type lights struct {
	data map[eid]*light // Camera instance data.
}

// newLights creates a light component manager. Expected to be called
// once on startup.
func newLights() *lights { return &lights{data: map[eid]*light{}} }

// get the light associated with the given entity, returning nil
// if there is no such light.
func (ls *lights) get(id eid) *light {
	if light, ok := ls.data[id]; ok {
		return light
	}
	return nil
}

// create and associate a light with the given entity. If there
// already is a light for the entity, dont create one and return
// the existing one instead.
func (ls *lights) create(id eid) *light {
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
