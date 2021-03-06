// Copyright © 2014-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gazed/vu"
	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/render"
)

// ps demonstrates a CPU-based particle system and a GPU-based (shader only)
// particle system with support provided by vu/Effect. The CPU particles
// need an update method - see fall and vent below. The GPU-based particles
// are updated by the shader code.
//
// CONTROLS:
//   Tab   : switch effect
//   WS    : move camera            : forward back
//   AD    : spin model             : left right
func ps() {
	defer catchErrors()
	if err := vu.Run(&pstag{}); err != nil {
		log.Printf("ps: error starting engine %s", err)
	}
}

// Globally unique "tag" for this example.
type pstag struct {
	scene   *vu.Ent    // 3D scene with camera.
	random  *rand.Rand // Random number generator.
	effects []*vu.Ent  // Particle effects.
	effect  *vu.Ent    // Active particle effect.
	index   int        // Active particle effect counter.

	// live particles are recalculated each update and
	// shared between the CPU particle effects.
	live []*vu.Particle // scratch particle list.
}

// Create is the engine callback for initial asset creation.
func (ps *pstag) Create(eng vu.Eng, s *vu.State) {
	eng.Set(vu.Title("Particle System"), vu.Size(400, 100, 800, 600))
	eng.Set(vu.Color(0.15, 0.15, 0.15, 1))
	ps.scene = eng.AddScene()
	ps.scene.Cam().SetClip(0.1, 50).SetFov(60).SetAt(0, 0, 2.5)
	ps.live = []*vu.Particle{}
	ps.random = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	// A GPU/shader based particle example using a particle shader.
	// Needs to initialize mesh data.
	gpu := ps.scene.AddPart()
	gpu.MakeEffect("particle", "particle")
	gpu.Cull(true)
	ps.makeParticles(gpu.Mesh())
	ps.effects = append(ps.effects, gpu)

	// A CPU/shader based particle example using an effect shader.
	cpu := ps.scene.AddPart()
	cpu.MakeEffect("effect", "particle").SetMover(ps.fall, 250)
	cpu.Cull(true)
	ps.effects = append(ps.effects, cpu)

	// A colorful exhaust attempt.
	// FUTURE: update textures to look like engine exhaust.
	jet := ps.scene.AddPart()
	jet.MakeEffect("exhaust", "exhaust").SetMover(ps.vent, 1000)
	jet.Cull(true)
	ps.effects = append(ps.effects, jet)

	// Make the first particle effect visible to kick things off.
	ps.effect = ps.effects[ps.index]
	ps.effect.Cull(false)
}

// Update is the engine frequent user-input/state-update callback.
func (ps *pstag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	run := 10.0   // move so many cubes worth in one second.
	spin := 270.0 // spin so many degrees in one second.
	dt := in.Dt
	cam := ps.scene.Cam()
	for press, down := range in.Down {
		switch press {
		case vu.KW:
			cam.Move(0, 0, dt*-run, cam.Look)
		case vu.KS:
			cam.Move(0, 0, dt*run, cam.Look)
		case vu.KA:
			ps.effect.Spin(0, dt*spin, 0)
		case vu.KD:
			ps.effect.Spin(0, dt*-spin, 0)
		case vu.KTab:
			if down == 1 {
				ps.effect.Cull(true) // switch to the next effect.
				ps.index = ps.index + 1
				if ps.index >= len(ps.effects) {
					ps.index = 0
				}
				ps.effect = ps.effects[ps.index]
				ps.effect.Cull(false)
			}
		}
	}
}

// Create GPU based particle vertex buffer data. Example from:
//     http://antongerdelan.net/opengl/particles.html
func (ps *pstag) makeParticles(m *vu.Mesh) {
	pcnt := 300                   // number of particles
	vv := make([]float32, pcnt*3) // vertex location.
	vt := make([]float32, pcnt)   // vertex time.
	var tdiff float32
	var index int
	for cnt := 0; cnt < pcnt; cnt++ {

		// start times
		vt[cnt] = tdiff
		tdiff += 0.01 // spacing for start times is 0.01 seconds

		// start velocities. randomly vary x and z components
		vv[index] = ps.random.Float32() - 0.5   // x
		vv[index+1] = 1                         // y
		vv[index+2] = ps.random.Float32() - 0.5 // z
		index += 3
	}
	m.InitData(0, 3, render.StaticDraw, false).SetData(0, vv)
	m.InitData(1, 1, render.StaticDraw, false).SetData(1, vt)
}

// fall is a CPU particle position updater. It lets particles drift downwards
// at a leisurely pace. Particles are started spread out at a the same height
// and then slowly moved down. Particles that have passed their maximum lifetime
// are removed.
//
// Update the list of particle data and pass back a list of pointers
// to the active particles.
func (ps *pstag) fall(all []*vu.Particle, dt float64) (live []*vu.Particle) {
	emit := 1                        // max particles emitted each update.
	lifespan := float32(1.0 / 200.0) // inverse number of updates to live.
	ps.live = ps.live[:0]            // reset keeping memory.
	for cnt, p := range all {
		switch {
		case p.Alive == 0 && emit > 0: // create particles each update.
			p.Alive, p.Index, emit = 1, float32(cnt), emit-1
			p.X = (ps.random.Float64() - 0.5) // randomly adjust X
			p.Y = 1                           // start at same height.
			p.Z = (ps.random.Float64() - 0.5) // randomly adjust Z
			ps.live = append(ps.live, p)
		case p.Alive > 0: // adjust live particles.
			p.Alive, p.Index = p.Alive-lifespan, float32(cnt)
			p.Y -= 0.01
			ps.live = append(ps.live, p)
		case p.Alive <= 0:
			p.Alive = 0 // reset expired particles.
		}
	}
	return ps.live
}

// vent is a CPU particle position updater. It uses a shader expecting a 2x2
// texture atlas where the textures are assigned according to the particle index.
//
// Update the list of particle data and pass back a list of pointers
// to the active particles.
func (ps *pstag) vent(all []*vu.Particle, dt float64) (live []*vu.Particle) {
	emit := 1                        // max particles emitted each update.
	lifespan := float32(1.0 / 400.0) // inverse number of updates to live.
	ps.live = ps.live[:0]            // reset keeping memory.
	for cnt, p := range all {
		switch {
		case p.Alive == 0 && emit == 1: // create particles slowly.
			p.Alive, p.Index, emit = 1, float32(cnt), emit-1

			// create points in a circular area.
			rad := lin.Rad(ps.random.Float64() * 180)
			r := ps.random.Float64() - 0.5
			p.X = r * math.Cos(rad) //  + x0 which is 0;
			p.Y = r * math.Sin(rad) //  + y0 which is 0;
			p.Z = 0                 // start at same height.
			ps.live = append(ps.live, p)
		case p.Alive > 0: // adjust live particles.
			p.Alive, p.Index = p.Alive-lifespan, float32(cnt)
			p.Z += 0.0005
			p.X = p.X * 1.0005 // move away from center 0.
			p.Y = p.Y * 1.0005 // move away from center 0.
			ps.live = append(ps.live, p)
		case p.Alive <= 0:
			p.Alive = 0 // reset expired particles.
		}
	}
	return ps.live
}
