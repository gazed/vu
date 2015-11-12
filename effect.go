// Copyright © 2014-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// Some particle-system tutorials and articles, in no particular order.
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
// *This last one is has excellent technical options.
//
// Design considerations hilighted by some of the above:
//    o getting/keeping particle data to/on the GPU. This allows more
//      particles, but moves the particle update programming to the shaders.
//    o keeping all particles part of a single vao.

import (
	"github.com/gazed/vu/render" // Needed to generate per-vertex data.
)

// Effect describes the application provided particle effect that updates
// the list of potential particles and returns the active particle set.
// Delta-time, dt, is the elapsed time in seconds since the last update.
// This is used for CPU particle effects where the application does the
// work of controlling particle lifespans and positions.
//
// An effect, once attached to a Model, positions and updates the model.
type Effect func(all []*EffectParticle, dt float64) (live []*EffectParticle)

// EffectParticle is one of the particles updated by an Effect.
// A set of these are returned by the Effect update method and are rendered
// by the engine.
type EffectParticle struct {
	Index   float32 // Particle number.
	Alive   float32 // Goes from 1:new particle, to 0:dead.
	X, Y, Z float64 // Particle location.
}

// Effect, EffectParticle
// =============================================================================
// effect

// effect turns an application defined particle Effect into points that
// can be rendered. Expected to be used with a relatively small number
// of particles.
type effect struct {
	source Effect    // App control of the particle positions & lifetimes.
	pv     []float32 // Scratch particle verticies: 3 floats per particle.
	pd     []float32 // Scratch particle data: 2 floats per particle.

	// particles is the complete list of all possible particles
	// that are allocated once on creation.
	particles []*EffectParticle
}

// newEffect expects source and model to be non-nil. A point based mesh is
// allocated. Return *effect to allow testing without having to cast.
func newEffect(m *model, source Effect, maxParticles int) *effect {
	if m.msh == nil {
		m.NewMesh("cpu")
		m.msh.loaded = true // mesh data will be set on update.
	}
	floatsPerVertex, floatsPerData := 3, 2
	m.InitMesh(0, uint32(floatsPerVertex), render.DYNAMIC, false)
	m.InitMesh(1, uint32(floatsPerData), render.DYNAMIC, false)

	// Allocate the maximum number of particles.
	particles := []*EffectParticle{}
	for cnt := 0; cnt < maxParticles; cnt++ {
		particles = append(particles, &EffectParticle{})
	}
	return &effect{source: source, particles: particles}
}

// update is called to transform the active particle set into vertex
// point data that can be sent to the GPU for rendering.
func (e *effect) update(m *model, dt float64) {
	// Have the application defined effect return the current
	// set of particles to be rendered.
	activeParticles := e.source(e.particles, dt)

	// Turn the particle positions and data into vertex buffer data.
	e.pv, e.pd = e.pv[:0], e.pd[:0] // keep previous memory.
	for _, p := range activeParticles {
		e.pv = append(e.pv, float32(p.X), float32(p.Y), float32(p.Z))
		e.pd = append(e.pd, p.Index, p.Alive)
	}
	if len(e.pv) > 0 {
		m.SetMeshData(0, e.pv)
		m.SetMeshData(1, e.pd)
	}
}
