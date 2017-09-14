// Copyright © 2014-2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"github.com/gazed/vu/render"
)

// effect.go groups particle effect code.
//
// Note: some particle-system tutorials and articles, in no particular order.
// Google "opengl particle system". The links demonstrate a wide range
// of implementation choices.
//    http://www.opengl-tutorial.org/intermediate-tutorials/billboards-particles/particles-instancing/
//    http://en.wikibooks.org/wiki/OpenGL_Programming/Particle_systems
//    http://directtovideo.wordpress.com/2009/10/06/a-thoroughly-modern-particle-system/
//    http://antongerdelan.net/opengl/particles.html
//    http://prideout.net/blog/?p=63
//    http://ogldev.atspace.co.uk/www/tutorial28/tutorial28.html
//    http://natureofcode.com/book/chapter-4-particle-systems
//   *http://stackoverflow.com/questions/17397724/point-sprites-for-particle-system
//   *This last one discusses some excellent technical options.
// DESIGN considerations hilighted by some of the above:
//    o getting/keeping particle data to/on the GPU. This allows more
//      particles, but moves the particle update programming to the shaders.
//    o keeping all particles part of a single vao.

// MakeEffect adds an particle effect component to an entity.
// Effect is a model that uses point based vertex data to display,
// move, create, and destroy, a bunch of small images.
//   shader : particle effect capable shader.
//   texture: particle texture applied to each particle.
func (e *Ent) MakeEffect(shader, texture string) *Ent {
	e.app.models.createEffect(e, "shd:"+shader, "tex:"+texture)
	return e
}

// SetMover ties the optional particle effect updater to this entity.
// Ignored if there is already a particle effect mover for this entity.
func (e *Ent) SetMover(mover Mover, maxParticles int) *Ent {
	e.app.models.setEffect(e.eid, mover, maxParticles)
	return e
}

// effect entity methods.
// =============================================================================
// Mover and Particle

// Mover is used to update CPU particle effects. It is the application
// that does the work of controlling particle lifespans and positions.
// All particles are passed in and the active particles are returned.
//    dt: delta-time is the elapsed time in seconds since the last update.
type Mover func(all []*Particle, dt float64) (live []*Particle)

// Particle is one of the particles updated by a Mover.
// A slice of these are returned by the Mover update method
// to be rendered by the engine.
type Particle struct {
	Index   float32 // Particle number.
	Alive   float32 // Goes from 1:new particle, to 0:dead.
	X, Y, Z float64 // Particle location.
}

// Mover and Particle
// =============================================================================
// effect

// effect turns an application defined particle effect into points that
// can be rendered. Expected to be used with a relatively small number
// of particles.
type effect struct {
	mover Mover     // App control of the particle positions & lifetimes.
	pv    []float32 // Scratch particle verticies: 3 floats per particle.
	pd    []float32 // Scratch particle data: 2 floats per particle.

	// parts is the complete list of all possible particles
	// that are allocated once on creation.
	parts []*Particle
}

// newEffect expects effect and model to be non-nil. A point based mesh
// is allocated. Return *effect to allow testing without having to cast.
func newEffect(m *Mesh, mov Mover, maxParticles int) *effect {
	floatsPerVertex, floatsPerData := 3, 2
	m.InitData(0, uint32(floatsPerVertex), render.DynamicDraw, false)
	m.InitData(1, uint32(floatsPerData), render.DynamicDraw, false)

	// Allocate the maximum number of particles.
	parts := []*Particle{}
	for cnt := 0; cnt < maxParticles; cnt++ {
		parts = append(parts, &Particle{})
	}
	return &effect{mover: mov, parts: parts}
}

// move is called to refresh the active particle set and the vertex
// point data that can be sent to the GPU for rendering.
func (e *effect) move(m *Mesh, parts []*Particle, dt float64) {
	// Have the application defined effect return the current
	// set of particles to be rendered.
	activeParticles := e.mover(e.parts, dt)

	// Turn the particle positions and data into vertex buffer data.
	e.pv, e.pd = e.pv[:0], e.pd[:0] // keep previous memory.
	for _, p := range activeParticles {
		e.pv = append(e.pv, float32(p.X), float32(p.Y), float32(p.Z))
		e.pd = append(e.pd, p.Index, p.Alive)
	}
	if len(e.pv) > 0 {
		m.SetData(0, e.pv)
		m.SetData(1, e.pd)
	}
}
