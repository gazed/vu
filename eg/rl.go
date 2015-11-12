// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gazed/vu"
	"github.com/gazed/vu/grid"
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
func rl() {
	rl := &rltag{}
	if err := vu.New(rl, "Random Levels", 400, 100, 800, 600); err != nil {
		log.Printf("rl: error starting engine %s", err)
	}
	defer catchErrors()
}

// Globally unique "tag" that encapsulates example specific data.
type rltag struct {
	ww, wh int            // Window size
	floors map[int]*floor // The random grid
	flr    *floor         // The current floor.
	arrow  vu.Pov         // Camera/player minimap marker.

	// timing values.
	renders int           // number of renders.
	elapsed time.Duration // time since last update.
	update  time.Duration // time of last update.
}

// Create is the engine callback for initial asset creation.
func (rl *rltag) Create(eng vu.Eng, s *vu.State) {
	rl.ww, rl.wh = 800, 600
	rl.floors = make(map[int]*floor)
	rl.setLevel(eng, vu.K_1)
	eng.SetColor(0.15, 0.15, 0.15, 1)
	return
}

// Update is the regular engine callback.
func (rl *rltag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	run := 5.0    // move so many cubes worth in one second.
	spin := 270.0 // spin so many degrees in one second.
	if in.Resized {
		rl.resize(s.W, s.H)
	}

	// pre-process user presses.
	// reduce individual move amounts for multiple move requests.
	dt := in.Dt
	moveDelta := dt * 2
	for press, _ := range in.Down {
		switch press {
		case vu.K_W, vu.K_S, vu.K_Q, vu.K_E:
			moveDelta *= 0.5
		}
	}

	// process user presses.
	for press, down := range in.Down {
		switch press {
		case vu.K_W:
			rl.flr.cam.Move(0, 0, moveDelta*-run, rl.flr.cam.Lookxz())
			rl.arrow.SetLocation(rl.flr.cam.Location())
		case vu.K_S:
			rl.flr.cam.Move(0, 0, moveDelta*run, rl.flr.cam.Lookxz())
			rl.arrow.SetLocation(rl.flr.cam.Location())
		case vu.K_Q:
			rl.flr.cam.Move(moveDelta*-run, 0, 0, rl.flr.cam.Lookxz())
			rl.arrow.SetLocation(rl.flr.cam.Location())
		case vu.K_E:
			rl.flr.cam.Move(moveDelta*run, 0, 0, rl.flr.cam.Lookxz())
			rl.arrow.SetLocation(rl.flr.cam.Location())
		case vu.K_A:
			rl.flr.cam.AdjustYaw(dt * spin)
			rl.arrow.SetRotation(rl.flr.cam.Lookxz())
		case vu.K_D:
			rl.flr.cam.AdjustYaw(dt * -spin)
			rl.arrow.SetRotation(rl.flr.cam.Lookxz())
		case vu.K_1, vu.K_2, vu.K_3, vu.K_4, vu.K_5, vu.K_6, vu.K_7, vu.K_8, vu.K_9, vu.K_0:
			if down == 1 {
				rl.setLevel(eng, press)
			}
		}
	}

	// show some stats to see the effectiveness of culling.
	allModels, allVerts := eng.Modelled()
	renModels, renVerts := eng.Rendered()
	modelStats := fmt.Sprintf("%d  models    culled to %d", allModels, renModels)
	vertexStats := fmt.Sprintf("%d verticies culled to %d", allVerts, renVerts)
	rl.flr.modelStats.SetPhrase(modelStats)
	rl.flr.vertexStats.SetPhrase(vertexStats)

	// http://stackoverflow.com/questions/87304/calculating-frames-per-second-in-a-game
	t := eng.Usage()
	rl.elapsed += t.Elapsed
	rl.update += t.Update
	rl.renders += t.Renders
	if in.Ut%50 == 0 { // average over 50 updates.
		fps := float64(rl.renders) / rl.elapsed.Seconds()
		update := rl.update.Seconds() / 50.0 * 1000
		timings := fmt.Sprintf("FPS %2.2f Update %3.2fms", fps, update)
		rl.flr.times.SetPhrase(timings)
		rl.renders = 0
		rl.elapsed = 0
		rl.update = 0
	}
}

// resize handles user screen/window changes.
func (rl *rltag) resize(ww, wh int) {
	rl.ww, rl.wh = ww, wh
	ratio := float64(ww) / float64(wh)
	for _, flr := range rl.floors {
		flr.cam.SetPerspective(60, ratio, 0.1, 50)
	}
}

// floor tracks all the information for a given level.
type floor struct {
	layout grid.Grid // the floor structure.
	top    vu.Pov    // top of floor transform hierarchy.

	// 3D scene.
	plan  vu.Pov    // how its drawn.
	arrow vu.Pov    // cam minimap location.
	cam   vu.Camera // main 3D camera.

	// 2D user interface including timing stats.
	ui          vu.Camera // overlay 2D camera.
	mmap        vu.Pov    // how its drawn on the minimap.
	mapPart     vu.Pov    // allows the minimap to be moved around.
	modelStats  vu.Model  // Show some render statistics.
	vertexStats vu.Model  // Show some render statistics.
	times       vu.Model  // Show some render statistics.
}

// setLevel switches to the indicated level.
func (rl *rltag) setLevel(eng vu.Eng, keyCode int) {
	if _, ok := rl.floors[keyCode]; !ok {
		var gridSizes = map[int]int{
			vu.K_1: 15,
			vu.K_2: 21,
			vu.K_3: 27,
			vu.K_4: 33,
			vu.K_5: 39,
			vu.K_6: 45,
			vu.K_7: 51,
			vu.K_8: 57,
			vu.K_9: 63,
			vu.K_0: 69,
		}
		var gridType = map[int]grid.Grid{
			vu.K_1: grid.New(grid.DENSE_SKIRMISH),
			vu.K_2: grid.New(grid.DENSE_SKIRMISH),
			vu.K_3: grid.New(grid.SPARSE_SKIRMISH),
			vu.K_4: grid.New(grid.SPARSE_SKIRMISH),
			vu.K_5: grid.New(grid.ROOMS_SKIRMISH),
			vu.K_6: grid.New(grid.ROOMS_SKIRMISH),
			vu.K_7: grid.New(grid.CAVE),
			vu.K_8: grid.New(grid.CAVE),
			vu.K_9: grid.New(grid.DUNGEON),
			vu.K_0: grid.New(grid.DUNGEON),
		}
		flr := &floor{}

		// create the scene
		flr.top = eng.Root().NewPov()
		flr.plan = flr.top.NewPov()
		flr.cam = flr.plan.NewCam()
		flr.cam.SetLocation(1, 0, -1)
		flr.cam.SetPerspective(60, float64(rl.ww)/float64(rl.wh), 0.1, 50)
		flr.cam.SetCull(vu.NewFrontCull(10))

		// create the overlay
		flr.mmap = flr.top.NewPov()
		flr.ui = flr.mmap.NewCam()
		flr.ui.SetUI()
		flr.ui.SetView(vu.XZ_XY)
		flr.ui.SetOrthographic(0, float64(rl.ww), 0, float64(rl.wh), 0, 20)
		flr.mapPart = flr.mmap.NewPov()
		flr.mapPart.SetScale(7.5, 7.5, 7.5)
		flr.mapPart.SetLocation(20, 0, -20)

		// display some rendering statistics.
		flr.modelStats = rl.newText(flr.mmap, 0)
		flr.vertexStats = rl.newText(flr.mmap, 1)
		flr.times = rl.newText(flr.mmap, 2)

		// populate the scenes
		lsize := gridSizes[keyCode]
		flr.layout = gridType[keyCode]
		flr.layout.Generate(lsize, lsize)
		width, height := flr.layout.Size()
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				if flr.layout.IsOpen(x, y) {
					block := flr.mapPart.NewPov().SetLocation(float64(x), 0, float64(-y))
					block.NewModel("alpha").LoadMesh("cube").LoadMat("transparent_gray")
				} else {
					block := flr.plan.NewPov().SetLocation(float64(x), 0, float64(-y))
					block.NewModel("uv").LoadMesh("box").AddTex("tile")
				}
			}
		}
		flr.arrow = flr.mapPart.NewPov().SetLocation(flr.cam.Location())
		flr.arrow.NewModel("solid").LoadMesh("arrow").LoadMat("transparent_blue")
		rl.floors[keyCode] = flr
	}
	if rl.flr != nil {
		rl.flr.plan.SetVisible(false)
		rl.flr.mmap.SetVisible(false)
	}
	rl.flr = rl.floors[keyCode]
	rl.flr.plan.SetVisible(true)
	rl.flr.mmap.SetVisible(true)
	rl.arrow = rl.flr.arrow
}

func (rl *rltag) newText(parent vu.Pov, gap int) vu.Model {
	text := parent.NewPov().SetLocation(10, 0, float64(-rl.wh+40+gap*24))
	text.Spin(-90, 0, 0) // orient to the X-Z plane.
	m := text.NewModel("uv").AddTex("lucidiaSu16White").LoadFont("lucidiaSu16")
	m.SetPhrase(" ")
	return m
}
