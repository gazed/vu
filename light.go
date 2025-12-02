// SPDX-FileCopyrightText : © 2015-2024 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

package vu

// light.go holder for lighting related code.
// FUTURE: different types of lights and more lights.

import (
	"log/slog"

	"github.com/gazed/vu/render"
)

// Types of lights
const (
	DirectionalLight = iota // no falloff ie: the sun
	PointLight              // falloff ie: light bulb
)

// AddLight adds a light to a scene.
// Lighting aware shaders combine the light with the models material
// values. Objects in the scene with lighting capable shaders will
// be influenced by the light.
//
// Each light is has a position and orientation in world space.
// The default position is 0,0,0 and the default direction is 0,0,-1.
// A lights default color is white 1,1,1.
//
// Depends on Eng.AddScene.
func (e *Entity) AddLight(lightType int) *Entity {

	// lights must be set directly on scenes
	if s := e.app.scenes.get(e.eid); s != nil {
		me := e.AddPart() // add a transform node for the light.
		if l := me.app.lights.create(me, e, lightType); l != nil {
			return me
		}
	}
	slog.Error("AddLight requires scene", "eid", e.eid)
	return e
}

// SetLight assigns a color and intensity to the light. The lights color is combined
// with objects material color to produce a final color. The R,G,B light values
// are between 0 (no color), and 1 (full color).
//
// Depends on Entity.AddLight.
func (e *Entity) SetLight(r, g, b, intensity float32) *Entity {
	if l := e.app.lights.get(e.eid); l != nil {
		l.r, l.g, l.b = r, g, b
		l.intensity = intensity
		return e
	}
	slog.Error("SetLight needs AddLight", "eid", e.eid)
	return e
}

// =============================================================================
// light data.

// light holds light color for now. Anticipate other future light parameters.
// The lights world position are scatch values updated each render loop.
type light struct {
	kind      int     // light type
	r, g, b   float32 // Light color: values are 0 to 1.
	intensity float32 //
}

// newLight creates a white light.
func newLight(lightType int) (l *light) {
	return &light{
		kind:      lightType,
		r:         1.0, // default white color.
		g:         1.0, //   ""
		b:         1.0, //   ""
		intensity: 5.0, // default intensity.
	}
}

// =============================================================================
// light component manager.

// lights manages all the active Light instances.
// There's not many lights so not much to optimize.
type lights struct {
	data map[eID]*light // All light instance data.
}

// newLights creates a light component manager. Expected to be called
// once on startup.
func newLights() *lights {
	return &lights{data: map[eID]*light{}}
}

// get the light associated with the given entity, returning nil
// if there is no such light.
func (ls *lights) get(id eID) *light {
	if light, ok := ls.data[id]; ok {
		return light
	}
	return nil
}

// create and associate a light with the given entity. If there
// already is a light for the entity, don't create one and return
// the existing one instead.
func (ls *lights) create(light, scene *Entity, lightType int) (l *light) {
	l = newLight(lightType)
	ls.data[light.eid] = l // All lights.
	return l
}

// fillPass populates the render.Pass data with the light information
// at the given pov.
func (ls *lights) fillLight(light *render.Light, lid eID, pov *pov) {
	l := ls.data[lid]
	px, py, pz := pov.at()
	light.X = float32(px)
	light.Y = float32(py)
	light.Z = float32(pz)
	light.R = l.r
	light.G = l.g
	light.B = l.b
}

// dispose the light associated for the given entity. Do nothing
// if no such light exists.
func (ls *lights) dispose(id eID) {
	delete(ls.data, id)
}
