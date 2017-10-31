// Copyright © 2013-2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gazed/vu"
	"github.com/gazed/vu/grid"
	"github.com/gazed/vu/math/lin"
)

// rl, random levels, tests higher level graphics functionality.
// This includes:
//   • vu culling/reducing the total objects rendered based on distance.
//   • vu 2D overlay scene, in this case a minimap.
//   • vu grid generation. Try numbers 0-9.
//   • vu engine statistics.
// rl also tests camera movement that includes holding multiple movement keys
// at the same time. The example does not have collision detection so you can
// literally run through the maze.
//
// This is the one example that displays and tests statistics that can
// be queried from the vu engine.
//
// CONTROLS:
//   WSQE  : move camera            : forward back left right
//   AD    : spin camera            : left right
//   1-9,0 : select level           : larger with higher num. 0 is 10
func rl() {
	defer catchErrors()
	if err := vu.Run(&rltag{}); err != nil {
		log.Printf("rl: error starting engine %s", err)
	}
}

// Globally unique "tag" that encapsulates example specific data.
type rltag struct {
	ww, wh int            // Window size
	floors map[int]*floor // The random grid
	flr    *floor         // The current floor.

	// timing values.
	renders int           // number of renders completed.
	elapsed time.Duration // time since last update.
	update  time.Duration // time of last update.
	render  time.Duration // time of last render.
}

// Create is the engine callback for initial asset creation.
func (rl *rltag) Create(eng vu.Eng, s *vu.State) {
	eng.Set(vu.Title("Random Levels"), vu.Size(400, 100, 800, 600))
	eng.Set(vu.Color(0.15, 0.15, 0.15, 1))
	rl.ww, rl.wh = 800, 600
	rl.floors = make(map[int]*floor)
	rl.setLevel(eng, vu.K1)
	return
}

// Update is the regular engine callback.
func (rl *rltag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	run := 5.0    // move so many cubes worth in one second.
	spin := 270.0 // spin so many degrees in one second.
	rl.ww, rl.wh = s.W, s.H

	// pre-process user presses.
	// reduce individual move amounts for multiple move requests.
	dt := in.Dt
	moveDelta := dt * 2
	for press := range in.Down {
		switch press {
		case vu.KW, vu.KS, vu.KQ, vu.KE:
			moveDelta *= 0.5
		}
	}

	// process user presses.
	for press, down := range in.Down {
		switch press {
		case vu.KW:
			rl.flr.move(0, 0, moveDelta*-run)
		case vu.KS:
			rl.flr.move(0, 0, moveDelta*run)
		case vu.KQ:
			rl.flr.move(moveDelta*-run, 0, 0)
		case vu.KE:
			rl.flr.move(moveDelta*run, 0, 0)
		case vu.KA:
			rl.flr.look(spin * dt)
		case vu.KD:
			rl.flr.look(-spin * dt)
		case vu.K1, vu.K2, vu.K3, vu.K4, vu.K5, vu.K6, vu.K7, vu.K8, vu.K9, vu.K0:
			if down == 1 {
				rl.setLevel(eng, press)
			}
		}
	}

	// show some stats to see the effectiveness of culling.
	stats := eng.Times()
	allModels, allTris, allVerts := stats.Modelled(eng)
	renModels, renTris, renVerts := stats.Rendered(eng)
	modelStats := fmt.Sprintf("%d models    culled to %d", allModels, renModels)
	triStats := fmt.Sprintf("%d triangles culled to %d", allTris, renTris)
	vertexStats := fmt.Sprintf("%d verticies culled to %d", allVerts, renVerts)
	rl.flr.modelStats.Typeset(modelStats)
	rl.flr.triStats.Typeset(triStats)
	rl.flr.vertexStats.Typeset(vertexStats)

	// http://stackoverflow.com/questions/87304/calculating-frames-per-second-in-a-game
	updates := uint64(50)       // Expecting 50 updates/second.
	rl.elapsed += stats.Elapsed // Count total time elapsed.
	rl.renders += stats.Renders // frames rendered.
	rl.render += stats.Render   // total time spent on renders.
	rl.update += stats.Update   // total time spent on updates.
	if stats.Skipped > 0 {
		log.Printf("Slow. Dropped updates: %d", stats.Skipped)
	}
	if in.Ut%updates == 0 {
		fps := float64(rl.renders) / rl.elapsed.Seconds()
		update := (rl.update.Seconds() / float64(updates)) * 1000      // in milliseconds.
		render := ((rl.render).Seconds() / float64(rl.renders)) * 1000 //      "
		timings := fmt.Sprintf("FPS %2.4f Update %2.4fms Render %2.4fms", fps, update, render)
		rl.flr.times.Typeset(timings)
		rl.renders = 0
		rl.elapsed = 0
		rl.update = 0
		rl.render = 0
	}
}

// floor tracks all the information for a given level.
type floor struct {
	layout grid.Grid // the floor structure.

	// 3D scene.
	scene *vu.Ent // top of 3D transform hierarchy.
	plan  *vu.Ent // how its drawn.
	arrow *vu.Ent // cam minimap location.

	// 2D user interface including timing stats.
	ui          *vu.Ent // 2D overlay camera.
	mmap        *vu.Ent // how its drawn on the minimap.
	mapPart     *vu.Ent // allows the minimap to be moved around.
	modelStats  *vu.Ent // Show some render statistics.
	vertexStats *vu.Ent //    "
	triStats    *vu.Ent //    "
	times       *vu.Ent // Show some render statistics.
}

// setLevel switches to the indicated level.
func (rl *rltag) setLevel(eng vu.Eng, keyCode int) {
	if _, ok := rl.floors[keyCode]; !ok {
		var gridSizes = map[int]int{
			vu.K1: 15,
			vu.K2: 21,
			vu.K3: 27,
			vu.K4: 33,
			vu.K5: 39,
			vu.K6: 45,
			vu.K7: 51,
			vu.K8: 57,
			vu.K9: 63,
			vu.K0: 69,
		}
		var gridType = map[int]grid.Grid{
			vu.K1: grid.New(grid.DenseSkirmish),
			vu.K2: grid.New(grid.DenseSkirmish),
			vu.K3: grid.New(grid.SparseSkirmish),
			vu.K4: grid.New(grid.SparseSkirmish),
			vu.K5: grid.New(grid.RoomSkirmish),
			vu.K6: grid.New(grid.RoomSkirmish),
			vu.K7: grid.New(grid.Cave),
			vu.K8: grid.New(grid.Cave),
			vu.K9: grid.New(grid.Dungeon),
			vu.K0: grid.New(grid.Dungeon),
		}
		flr := &floor{}

		// create the scene
		flr.scene = eng.AddScene()
		flr.scene.Cam().SetClip(0.1, 50).SetFov(60).SetAt(1, 0, -1)
		flr.plan = flr.scene.AddPart()
		flr.scene.SetCuller(vu.NewFrontCull(10))

		// create the overlay
		flr.ui = eng.AddScene().SetUI()
		flr.ui.Cam().SetClip(0, 20)
		flr.mmap = flr.ui.AddPart()
		flr.mapPart = flr.mmap.AddPart().SetScale(7, 7, 7).SetAt(20, 20, 0)

		// display some rendering statistics.
		flr.modelStats = rl.newText(flr.mmap, 0)
		flr.triStats = rl.newText(flr.mmap, 1)
		flr.vertexStats = rl.newText(flr.mmap, 2)
		flr.times = rl.newText(flr.mmap, 3)

		// populate the scenes
		lsize := gridSizes[keyCode]
		flr.layout = gridType[keyCode]
		flr.layout.Generate(lsize, lsize)
		width, height := flr.layout.Size()
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				if flr.layout.IsOpen(x, y) {
					// minimap overlay
					block := flr.mapPart.AddPart().SetAt(float64(x), float64(y), 0)
					block.MakeModel("alpha", "msh:cube", "mat:transparent_gray")
				} else {
					// floor level.
					block := flr.plan.AddPart().SetAt(float64(x), 0, float64(-y))
					block.MakeModel("uv", "msh:box", "tex:tile")
				}
			}
		}
		flr.arrow = flr.mapPart.AddPart().SetAt(1, 1, 0)
		flr.arrow.MakeModel("solid", "msh:arrow", "mat:transparent_blue")
		rl.floors[keyCode] = flr
	}
	if rl.flr != nil {
		rl.flr.scene.Cull(true)
		rl.flr.ui.Cull(true)
	}
	rl.flr = rl.floors[keyCode]
	rl.flr.scene.Cull(false)
	rl.flr.ui.Cull(false)
}

// move changes from 3D X,Z coordinates to X,Y needed by the overlay.
func (f *floor) move(x, y, z float64) {
	cam := f.scene.Cam()
	cam.Move(x, y, z, cam.Look)
	cx, _, cz := f.scene.Cam().At()
	f.arrow.SetAt(cx, cz*-1, 0)
}

// look changes rotation about Y into rotation about Z.
func (f *floor) look(spin float64) {
	cam := f.scene.Cam()
	cam.SetYaw(cam.Yaw + spin)
	f.arrow.View().SetAa(0, 0, 1, lin.Rad(cam.Yaw))
}

// newText is a utility method for creating a new text label.
func (rl *rltag) newText(parent *vu.Ent, gap int) *vu.Ent {
	text := parent.AddPart().SetAt(10, float64(rl.wh-(40+gap*24)), 0)
	return text.MakeLabel("txt", "lucidiaSu16")
}
