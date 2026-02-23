// SPDX-FileCopyrightText : © 2015-2024 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

package vu

// light.go holder for lighting related code.
// FUTURE: different types of lights and more lights.

import (
	"log/slog"
	"math"

	"github.com/gazed/vu/render"
)

// Types of lights
const (
	SunLight   = iota // sun rays   : position implies direction.
	PointLight        // light bulb : position and direction.
	SpotLight         // flash light: position and direction and cone cutoff angle.
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

// SetCutoff assigns the cone cutoff angle, in radians, that
// is required for spot lights. The default is 0.5 radians, about 28.6deg.
//
// Depends on Entity.AddLight.
func (e *Entity) SetCutoff(angleInRadians float32) *Entity {
	if l := e.app.lights.get(e.eid); l != nil {
		l.cutoff = float32(math.Cos(float64(angleInRadians)))
		return e
	}
	slog.Error("SetCone needs AddLight", "eid", e.eid)
	return e
}

// =============================================================================
// light data.

// lights have color and intensity. Spot lights have a cutoff angle.
// The position and orientation of lights are handle by pov components.
type light struct {
	kind        int     // light type
	r, g, b     float32 // Light color: values are 0 to 1.
	attenuation float32 // light attenuation.
	intensity   float32 // light intensity.
	cutoff      float32 // spot light cone.
}

// newLight creates a white light.
func newLight(lightType int) (l *light) {
	cutoffAngle := float32(math.Cos(0.5)) // default spotlight cutoff angle ~28.6deg.
	return &light{
		kind:      lightType,
		r:         1.0,         // default white color.
		g:         1.0,         //   ""
		b:         1.0,         //   ""
		intensity: 5.0,         // default intensity.
		cutoff:    cutoffAngle, // default cutoff angle.
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
	light.R = l.r
	light.G = l.g
	light.B = l.b
	light.Intensity = l.intensity

	// position for all lights.
	px, py, pz := pov.at()
	light.Px = float32(px)
	light.Py = float32(py)
	light.Pz = float32(pz)
	light.Attenuation = l.attenuation

	// direction for spotlight.
	dx, dy, dz, _ := pov.tn.Rot.Aa()
	light.Dx = float32(dx)
	light.Dy = float32(dy)
	light.Dz = float32(dz)
	light.Cutoff = l.cutoff
}

// dispose the light associated for the given entity. Do nothing
// if no such light exists.
func (ls *lights) dispose(id eID) {
	delete(ls.data, id)
}
