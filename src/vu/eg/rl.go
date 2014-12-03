// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"log"
	"vu"
	"vu/grid"
)

// rl tests higher level graphics functionality. This includes:
//   • vu culling/reducing the total objects rendered based on distance. See:
//      • plan.SetCuller()
//   • vu 2D overlay scene, in this case a minimap.  See:
//    	• flr.mmap = rl.eng.AddScene(vu.XZtoXY)
//    	• flr.mmap.Set2D()
//    	• flr.mmap.SetOrthographic()
//   • vu grid generation. See: vu/grid.
// rl also tests camera movement that includes holding multiple movement keys
// at the same time. The example does not have collision detection so you can
// literally run through the maze.
//
// rl also tests vu/grid by generating different types and size of grids using
// the number keys 0-9.
func rl() {
	rl := &rltag{}
	var err error
	if rl.eng, err = vu.New("Random Levels", 400, 100, 800, 600); err != nil {
		log.Printf("rl: error intitializing engine %s", err)
		return
	}
	rl.eng.SetDirector(rl) // get user input through Director.Update()
	rl.create()            // create initial assests.
	if err = rl.eng.Verify(); err != nil {
		log.Fatalf("rl: error initializing model :: %s", err)
	}
	defer rl.eng.Shutdown()
	defer catchErrors()
	rl.eng.Action()
}

// Globally unique "tag" that encapsulates example specific data.
type rltag struct {
	eng           vu.Engine
	floors        map[string]*floor // The random grid
	flr           *floor            // The current floor.
	width, height int               // Window size
	arrow         vu.Part           // Camera/player minimap marker.
	run           float64           // Camera movement speed.
	spin          float64           // Camera spin speed.
}

// create is the startup asset creation.
func (rl *rltag) create() {
	rl.run = 5    // move so many cubes worth in one second.
	rl.spin = 270 // spin so many degrees in one second.
	rl.width, rl.height = 800, 600
	rl.floors = make(map[string]*floor)
	rl.setLevel("1")
	rl.eng.Color(0.15, 0.15, 0.15, 1)
	return
}

// Update is the regular engine callback.
func (rl *rltag) Update(in *vu.Input) {
	if in.Resized {
		rl.resize()
	}

	// pre-process user presses.
	// reduce individual move amounts for multiple move requests.
	dt := in.Dt
	moveDelta := dt * 2
	for press, _ := range in.Down {
		switch press {
		case "W", "S", "Q", "E":
			moveDelta *= 0.5
		}
	}

	// process user presses.
	for press, down := range in.Down {
		switch press {
		case "W":
			rl.flr.cam.Move(0, 0, moveDelta*-rl.run)
			rl.arrow.SetLocation(rl.flr.cam.Location())
		case "S":
			rl.flr.cam.Move(0, 0, moveDelta*rl.run)
			rl.arrow.SetLocation(rl.flr.cam.Location())
		case "Q":
			rl.flr.cam.Move(moveDelta*-rl.run, 0, 0)
			rl.arrow.SetLocation(rl.flr.cam.Location())
		case "E":
			rl.flr.cam.Move(moveDelta*rl.run, 0, 0)
			rl.arrow.SetLocation(rl.flr.cam.Location())
		case "A":
			rl.flr.cam.Spin(0, dt*rl.spin, 0)
			rl.arrow.SetRotation(rl.flr.cam.Rotation())
		case "D":
			rl.flr.cam.Spin(0, dt*-rl.spin, 0)
			rl.arrow.SetRotation(rl.flr.cam.Rotation())
		case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0":
			if down == 1 {
				rl.setLevel(press)
			}
		}
	}
}

// resize handles user screen/window changes.
func (rl *rltag) resize() {
	x, y, width, height := rl.eng.Size()
	rl.eng.Resize(x, y, width, height)
	rl.width = width
	rl.height = height
	ratio := float64(width) / float64(height)
	for _, flr := range rl.floors {
		flr.cam.SetPerspective(60, ratio, 0.1, 50)
	}
}

// floor tracks all the information for a given level.
type floor struct {
	layout  grid.Grid // the floor structure.
	arrow   vu.Part   // cam minimap location.
	plan    vu.Scene  // how its drawn.
	mmap    vu.Scene  // how its drawn on the minimap.
	cam     vu.Camera // main 3D camera.
	mapPart vu.Part   // allows the minimap to be moved around.
}

// setLevel switches to the indicated level.
func (rl *rltag) setLevel(id string) {
	if _, ok := rl.floors[id]; !ok {
		var gridSizes = map[string]int{
			"1": 15,
			"2": 21,
			"3": 27,
			"4": 33,
			"5": 39,
			"6": 45,
			"7": 51,
			"8": 57,
			"9": 63,
			"0": 69,
		}
		var gridType = map[string]grid.Grid{
			"1": grid.New(grid.DENSE_SKIRMISH),
			"2": grid.New(grid.DENSE_SKIRMISH),
			"3": grid.New(grid.SPARSE_SKIRMISH),
			"4": grid.New(grid.SPARSE_SKIRMISH),
			"5": grid.New(grid.ROOMS_SKIRMISH),
			"6": grid.New(grid.ROOMS_SKIRMISH),
			"7": grid.New(grid.CAVE),
			"8": grid.New(grid.CAVE),
			"9": grid.New(grid.DUNGEON),
			"0": grid.New(grid.DUNGEON),
		}
		flr := &floor{}

		// create the scene
		flr.plan = rl.eng.AddScene(vu.VP)
		flr.plan.SetSorted(true)
		flr.plan.SetCuller(vu.NewFacingCuller(10))
		flr.cam = flr.plan.Cam()
		flr.cam.SetLocation(1, 0, -1)
		flr.cam.SetPerspective(60, float64(rl.width)/float64(rl.height), 0.1, 50)

		// create the overlay
		flr.mmap = rl.eng.AddScene(vu.XZ_XY)
		flr.mmap.Set2D()
		flr.mmap.Cam().SetOrthographic(-0.2, 100, -0.2, 75, 0, 10)
		flr.mapPart = flr.mmap.AddPart()
		flr.mapPart.SetLocation(3, 0, -3)

		// populate the scenes
		lsize := gridSizes[id]
		flr.layout = gridType[id]
		flr.layout.Generate(lsize, lsize)
		width, height := flr.layout.Size()
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				if flr.layout.IsOpen(x, y) {
					block := flr.mapPart.AddPart()
					block.SetLocation(float64(x), 0, float64(-y))
					block.SetRole("flat").SetMesh("cube").SetMaterial("gray")
				} else {
					block := flr.plan.AddPart()
					block.SetLocation(float64(x), 0, float64(-y))
					block.SetRole("gouraud").SetMesh("cube").SetMaterial("cube")
					block.Role().SetLightLocation(0, 10, 0)
					block.Role().SetLightColour(0.4, 0.7, 0.9)
				}
			}
		}
		flr.arrow = flr.mapPart.AddPart()
		flr.arrow.SetRole("flat").SetMesh("arrow").SetMaterial("blue")
		flr.arrow.SetLocation(flr.cam.Location())
		rl.floors[id] = flr
	}
	if rl.flr != nil {
		rl.flr.plan.SetVisible(false)
		rl.flr.mmap.SetVisible(false)
	}
	rl.flr = rl.floors[id]
	rl.flr.plan.SetVisible(true)
	rl.flr.mmap.SetVisible(true)
	rl.arrow = rl.flr.arrow
}
