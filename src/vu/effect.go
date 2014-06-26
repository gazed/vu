// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

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
//
// Design considerations hilighted by some of the above:
//    o getting/keeping particle data to/on the GPU. This allows more particles,
//      but moves the particle update programming to the shaders.
//    o keeping all particles part of a single vao.

import (
	"vu/render" // Needed to generate per-vertex data.
)

// Effect tracks the active particles over their lifespans.
// An effect is combined with a EffectMover that adjusts active particles
// positions. The EffectMover will be called each update.
type Effect interface {

	// Update will emit new particles, remove particles that are past their
	// life time, and update the attributes of currently active particles.
	// Delta-time, dt, is the elapsed time in seconds, eg: 0.02.
	Update(m render.Mesh, dt float64) // Moves particles, refreshes render data.
}

// NewEffect creates a new particle effect that will not exceed maxParticles
// and which will generate new particles at the given rate per second.
// The effect relies on an EffectMover supplied by the application which
// updates particles location and life information.
func NewEffect(maxParticles, rate int, mover EffectMover) Effect {
	return newEffect(maxParticles, rate, mover)
}

// EffectMover is the Application provided particle updater expected by
// NewEffect(). The EffectMover is expected to update the currently active
// particles each time it is called. Delta-time, dt, is the elapsed time
// in seconds.
type EffectMover func(activeParticles []*EffectParticle, dt float64) []*EffectParticle

// EffectParticle corresponds to a particle updated by the EffectMover method.
type EffectParticle struct {
	Life    float64 // Time alive in seconds.
	X, Y, Z float32 // Particle location.
}

// Effect, EffectMover, EffectParticle
// =============================================================================
// effect

// effect is called each update to emit new particles and update existing
// particles and emitter properties.  Effect controls each particles lifetime,
// position, rotation, colour. It can also control the location and direction
// of the emitter.
type effect struct {
	source *emitter          // location, orientation, and generator points.
	plife  []float32         // Particle lifetime buffer.
	pb     []float32         // scratch buffer: 4 floats per particle.
	active []*EffectParticle // Currently active particles.
	move   EffectMover       // Set by App to control particle position & lifetime.
	max    int               // Maximum number of particles.
}

//newEffect allows testing to use *effect without having to cast.
func newEffect(maxParticles, rate int, em EffectMover) *effect {
	e := &effect{}
	e.move = em
	e.max = maxParticles
	e.source = &emitter{rate: rate}

	// Create enough capacity to hold the maximum number of particles.
	floatsPerParticle := 3
	e.pb = make([]float32, e.max*floatsPerParticle)
	return e
}

// Effect interface implementation.
func (e *effect) Update(m render.Mesh, dt float64) {

	// emit more particles
	particles := e.source.emit(dt)
	if len(e.active)+len(particles) < e.max {
		e.active = append(e.active, particles...)
	} else {
		// Ignore developer error where particles are not being removed
		// as fast as they are generated.
	}

	// have the application mover update the location, and remove particles
	// that have outlived their time.
	if e.move != nil {
		e.active = e.move(e.active, dt) // update particles
		e.pb = e.pb[:0]                 // keep previous memory.
		for _, p := range e.active {
			e.pb = append(e.pb, p.X, p.Y, p.Z)
		}
	}

	// update the data.
	if len(e.pb) > 0 && m != nil {
		floatsPerParticle := uint32(3)
		m.InitData(0, floatsPerParticle, render.DYNAMIC, false).SetData(0, e.pb)
	}

	// FUTURE: make particle updates concurrent.
}

// =============================================================================

// emitter creates new particles at a pre-determined rate.
type emitter struct {
	rate      int               // particles per second.
	time      float64           // tracks elapsed delta time for slow particle rates.
	particles []*EffectParticle // reusable storage for new particles.
}

// FUTURE emit particles from random verticies of a mesh.
//	  mesh string  // verticies are emit points, normals are directions.

// newEmitter intializes an new emitter with the given particle rate per second
// and the max particle life life parameters.
func newEmitter(rate int) *emitter {
	e := &emitter{}
	e.rate = rate
	e.particles = make([]*EffectParticle, rate)
	return e
}

// emit creates new particles.
// The particle pointers that emit returns must be referenced by something
// else before calling emit again. Note that emit must be called at least once
// per second to get the desired rate. It is expected to be called each update.
func (e *emitter) emit(dt float64) []*EffectParticle {
	e.particles = e.particles[:0] // reset temporary particle storage.
	e.time += dt                  // handle slow particle rates.
	if need := int(e.time * float64(e.rate)); need > 0 {
		if need > e.rate {
			need = e.rate // don't overflow particle storage.
			// Developer error: emit not being called
			// fast enough or dt is to large.
		}
		for cnt := 0; cnt < need; cnt++ {
			p := &EffectParticle{}
			e.particles = append(e.particles, p)
		}
		e.time = 0
	}
	return e.particles
}
