// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package main

import (
	"log"
	"math/rand"
	"sort"
	"time"
	"vu"
	"vu/render"
)

// ps demonstrates a CPU-based particle system and a GPU-based (shader only)
// particle system. This demo concentrates on vu/Effect.
func ps() {
	ps := &pstag{}
	var err error
	if ps.eng, err = vu.New("Particle System", 1200, 100, 800, 600); err != nil {
		log.Printf("ps: error intitializing engine %s", err)
		return
	}
	ps.run = 10             // move so many cubes worth in one second.
	ps.spin = 270           // spin so many degrees in one second.
	ps.eng.SetDirector(ps)  // override user input handling.
	defer ps.eng.Shutdown() // shut down the engine.
	defer catchErrors()
	ps.eng.Action()
}

// Globally unique "tag" for this example.
type pstag struct {
	eng    vu.Engine            // 3D engine.
	scene  vu.Scene             // scene graph.
	run    float64              // Camera movement speed.
	spin   float64              // Camera spin speed.
	gshd   vu.Part              // GPU particle effect.
	cshd   vu.Part              // CPU particle effect.
	random *rand.Rand           // Random number generator.
	live   []*vu.EffectParticle // scratch particle list.
}

// Create is the engine one-time initialization callback.
func (ps *pstag) Create(eng vu.Engine) {
	ps.live = []*vu.EffectParticle{}
	ps.random = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	ps.scene = eng.AddScene(vu.VP)
	ps.scene.SetPerspective(60, float64(800)/float64(600), 0.1, 50)
	ps.scene.SetLightLocation(0, 5, 0)
	ps.scene.SetLightColour(0.4, 0.7, 0.9)
	ps.scene.SetLocation(0, 10, 25)

	// background slab for context.
	slab := ps.scene.AddPart().SetLocation(0, -25, 0).SetScale(50, 50, 50)
	slab.SetRole("gouraud").SetMesh("cube").SetMaterial("floor")
	slab.SetBody(vu.Box(25, 25, 25), 0, 0.4)

	// Add the GPU/shader based particle example.
	ps.gshd = ps.scene.AddPart().SetLocation(0, 10, 22)
	ps.gshd.SetVisible(false)
	ps.gshd.SetRole("particle").AddTex("particle")
	ps.gshd.Role().SetDrawMode(vu.POINTS)
	ps.gshd.Role().Set2D()
	ps.makeParticles(ps.gshd.Role().Mesh())

	// Add the CPU/shader based particle example.
	ps.cshd = ps.scene.AddPart().SetLocation(0, 10, 22)
	ps.cshd.SetRole("effect").AddTex("particle").SetDrawMode(vu.POINTS)
	ps.cshd.Role().SetEffect(vu.NewEffect(250, 25, ps.fall))

	// Have a lighter default background.
	eng.Color(0.15, 0.15, 0.15, 1)
}

// Update is the engine frequent user-input/state-update callback.
func (ps *pstag) Update(in *vu.Input) {
	if in.Resized {
		ps.resize()
	}
	dt := in.Dt
	for press, down := range in.Down {
		switch press {
		case "W":
			ps.scene.Move(0, 0, dt*-ps.run)
		case "S":
			ps.scene.Move(0, 0, dt*ps.run)
		case "A":
			ps.scene.Spin(vu.YAxis, dt*ps.spin)
		case "D":
			ps.scene.Spin(vu.YAxis, dt*-ps.spin)
		case "H":
			if down == 1 {
				ps.gshd.SetVisible(!ps.gshd.Visible())
				ps.cshd.SetVisible(!ps.cshd.Visible())
			}
		}
	}
}

// resize sets the view port size.  User resizes are ignored.
func (ps *pstag) resize() {
	x, y, width, height := ps.eng.Size()
	ps.eng.Resize(x, y, width, height)
	ps.scene.SetPerspective(60, float64(width)/float64(height), 0.1, 50)
}

// Create GPU based particle vertex buffer data. Example from:
//     http://antongerdelan.net/opengl/particles.html
func (ps *pstag) makeParticles(m render.Mesh) {
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
	m.InitData(0, 3, render.STATIC, false).SetData(0, vv)
	m.InitData(1, 1, render.STATIC, false).SetData(1, vt)
}

// fall is the CPU particle position updater. It lets particles drift downwards
// at a leisurely pace. Particles are started spread out at a the same height
// and then slowly moved down. Particles that have passed their maximum lifetime
// are removed.
func (ps *pstag) fall(particles []*vu.EffectParticle, dt float64) []*vu.EffectParticle {
	ps.live = ps.live[:0]
	for _, p := range particles {

		// set the initial position for a particle.
		if p.Life == 0 {
			p.X += (ps.random.Float32() - 0.5) // randomly adjust X
			p.Y = 1                            // start at same height.
			p.Z += (ps.random.Float32() - 0.5) // randomly adjust Z
		}
		p.Life += dt
		if p.Life < 4.0 {
			p.Y -= float32(0.01)
			ps.live = append(ps.live, p)
		}
	}
	sort.Sort(Ordered(ps.live))
	return ps.live
}

// Ordered allows particles to be sorted.
type Ordered []*vu.EffectParticle

// Sort particles ordered by Z distance.
func (o Ordered) Len() int           { return len(o) }
func (o Ordered) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o Ordered) Less(i, j int) bool { return o[i].Z < o[j].Z }
