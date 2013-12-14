// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package main

import (
	"log"
	"time"
	"vu"
	"vu/grid"
)

// rl tests higher level graphics functionality. This includes:
//    vu culling/reducing the total objects rendered based on distance. See:
//   		 flr.plan.SetVisibleRadius(10)
//   		 flr.plan.SetVisibleDirection(true)
//    vu 2D overlay scene, in this case a minimap.  See:
//   		flr.mmap = rl.eng.AddScene(vu.XZtoXY)
//   		flr.mmap.Set2D()
//   		flr.mmap.SetOrthographic(-0.2, 100, -0.2, 75, 0, 10)
//    vu grid generation. See:
//   		 &grid.Dense{},
//   		 &grid.Sparse{},
//   		 &grid.Rooms{},
//
// rl also tests camera movement that includes holding multiple movement keys
// at the same time. The example does not have collision detection so you can
// literally run through the maze.
func rl() {
	rl := &rltag{}
	var err error
	if rl.eng, err = vu.New("Random Levels", 400, 100, 800, 600); err != nil {
		log.Printf("rl: error intitializing engine %s", err)
		return
	}
	rl.run = 5             // move so many cubes worth in one second.
	rl.spin = 270          // spin so many degrees in one second.
	rl.eng.SetDirector(rl) // override user input handling.
	defer rl.eng.Shutdown()
	rl.eng.Action()
}

// Globally unique "tag" that encapsulates example specific data.
type rltag struct {
	eng           vu.Engine
	floors        map[string]*floor // The random grid
	flr           *floor            // The current floor.
	width, height int               // Window size
	holdoff       time.Duration     // Time in milliseconds before swtiching levels.
	levelSwitch   time.Time         // Last time a level switch happened.
	arrow         vu.Part           // Camera/player minimap marker.
	run           float64           // Camera movement speed.
	spin          float64           // Camera spin speed.
}

// Create is the engine intialization callback.
func (rl *rltag) Create(eng vu.Engine) {
	rl.width, rl.height = 800, 600
	rl.floors = make(map[string]*floor)
	rl.setLevel("1")
	rl.holdoff, _ = time.ParseDuration("1000ms")
	rl.levelSwitch = time.Now()

	// set some constant state.
	rl.eng.Enable(vu.BLEND, true)
	rl.eng.Enable(vu.CULL, true)
	rl.eng.Color(0.1, 0.1, 0.1, 1)
	return
}

// Update is the regular engine callback.
func (rl *rltag) Update(input *vu.Input) {
	if input.Resized {
		rl.resize()
	}

	// pre-process user presses.
	// reduce individual move amounts for multiple move requests.
	dt := input.Dt
	moveDelta := dt * 2
	for press, _ := range input.Down {
		switch press {
		case "W", "S", "Q", "E":
			moveDelta *= 0.5
		}
	}

	// process user presses.
	for press, _ := range input.Down {
		switch press {
		case "W":
			rl.flr.plan.MoveView(0, 0, moveDelta*-rl.run)
			rl.arrow.SetLocation(rl.flr.plan.ViewLocation())
		case "S":
			rl.flr.plan.MoveView(0, 0, moveDelta*rl.run)
			rl.arrow.SetLocation(rl.flr.plan.ViewLocation())
		case "Q":
			rl.flr.plan.MoveView(moveDelta*-rl.run, 0, 0)
			rl.arrow.SetLocation(rl.flr.plan.ViewLocation())
		case "E":
			rl.flr.plan.MoveView(moveDelta*rl.run, 0, 0)
			rl.arrow.SetLocation(rl.flr.plan.ViewLocation())
		case "A":
			rl.flr.plan.PanView(vu.YAxis, dt*rl.spin)
			rl.arrow.SetRotation(rl.flr.plan.ViewRotation())
		case "D":
			rl.flr.plan.PanView(vu.YAxis, dt*-rl.spin)
			rl.arrow.SetRotation(rl.flr.plan.ViewRotation())
		case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0":
			if time.Now().After(rl.levelSwitch.Add(rl.holdoff)) {
				rl.setLevel(press)
				rl.levelSwitch = time.Now()
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
		flr.plan.SetPerspective(60, ratio, 0.1, 50)
	}
}

// floor tracks all the information for a given level.
type floor struct {
	layout  grid.Grid // the floor structure.
	arrow   vu.Part   // cam minimap location.
	plan    vu.Scene  // how its drawn.
	mmap    vu.Scene  // how its drawn on the minimap.
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
			"3": grid.New(grid.DENSE_SKIRMISH),
			"4": grid.New(grid.SPARSE_SKIRMISH),
			"5": grid.New(grid.SPARSE_SKIRMISH),
			"6": grid.New(grid.SPARSE_SKIRMISH),
			"7": grid.New(grid.ROOMS_SKIRMISH),
			"8": grid.New(grid.ROOMS_SKIRMISH),
			"9": grid.New(grid.ROOMS_SKIRMISH),
			"0": grid.New(grid.ROOMS_SKIRMISH),
		}
		flr := &floor{}

		// create the scene
		flr.plan = rl.eng.AddScene(vu.VP)
		flr.plan.SetPerspective(60, float64(rl.width)/float64(rl.height), 0.1, 50)
		flr.plan.SetLightLocation(0, 10, 0)
		flr.plan.SetLightColour(0.4, 0.7, 0.9)
		flr.plan.SetVisibleRadius(10)
		flr.plan.SetSorted(true)
		flr.plan.SetVisibleDirection(true)

		// create the overlay
		flr.mmap = rl.eng.AddScene(vu.XZ_XY)
		flr.mmap.Set2D()
		flr.mmap.SetOrthographic(-0.2, 100, -0.2, 75, 0, 10)
		flr.mmap.SetLightLocation(1, 1, 1)
		flr.mmap.SetLightColour(1, 1, 1)
		flr.mapPart = flr.mmap.AddPart()
		flr.mapPart.SetLocation(3, 0, -3)
		flr.plan.SetViewLocation(1, 0, -1)

		// populate the scenes
		lsize := gridSizes[id]
		flr.layout = gridType[id]
		flr.layout.Generate(lsize, lsize)
		width, height := flr.layout.Size()
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				if flr.layout.IsWall(x, y) {
					block := flr.plan.AddPart()
					block.SetFacade("cube", "gouraud").SetMaterial("cube")
					block.SetLocation(float64(x), 0, float64(-y))

					// use flat, gray for the overlay.
					block = flr.mapPart.AddPart()
					block.SetFacade("cube", "flat").SetMaterial("gray")
					block.SetLocation(float64(x), 0, float64(-y))
				}
			}
		}
		flr.arrow = flr.mapPart.AddPart()
		flr.arrow.SetFacade("arrow", "flat").SetMaterial("blue")
		flr.arrow.SetLocation(flr.plan.ViewLocation())
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
