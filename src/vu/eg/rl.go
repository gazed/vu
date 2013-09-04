// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package main

import (
	"log"
	"time"
	"vu"
	"vu/data"
	"vu/grid"
)

// rl tests the higher level graphics functionality. This includes:
// vu culling/reducing the total objects rendered based on distance:
//		flr.plan.SetVisibleRadius(10)
//		flr.plan.SetVisibleDirection(true)
// The vu/grid package grid generation. Test the following grid types:
//		&grid.Dense{},
//		&grid.Sparse{},
//		&grid.Rooms{},
// The ability to have a 2D overlay scene, in this case a minimap:
//		flr.mmap = rl.eng.AddScene(vu.XZtoXY)
//		flr.mmap.Set2D()
//		flr.mmap.SetOrthographic(-0.2, 100, -0.2, 75, 0, 10)
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
	rl.stagePlay()
	defer rl.eng.Shutdown()
	rl.eng.Action()
}

// Globally unique "tag" for this example.
type rltag struct {
	eng           *vu.Eng
	arrow         vu.Part
	floors        map[string]*floor // the random grid
	flr           *floor            // the current floor.
	width, height int               // window size
	holdoff       time.Duration     // time in milliseconds before a swtiching grids.
	last          time.Time
	run           float32
	spin          float32
}

func (rl *rltag) stagePlay() {
	rl.build()
	rl.holdoff, _ = time.ParseDuration("1000ms")
	rl.last = time.Now()

	// set some constant state.
	rl.eng.Enable(vu.BLEND, true)
	rl.eng.Enable(vu.CULL, true)
	rl.eng.Color(0.1, 0.1, 0.1, 1)
	return
}

type floor struct {
	layout  grid.Grid // the floor structure.
	arrow   vu.Part   // cam minimap location.
	plan    vu.Scene  // how its drawn.
	mmap    vu.Scene  // how its drawn on the minimap.
	mapPart vu.Part   // allows the minimap to be moved around.
}

// Create the maze that can be walked through.
func (rl *rltag) build() {
	rl.width, rl.height = 800, 600
	rl.eng.Load(rl.makeArrowMesh())
	rl.eng.Load(&data.Material{"gray", data.Rgb{0.4, 0.4, 0.4}, data.Rgb{}, data.Rgb{}, 0.5})
	rl.eng.Load(&data.Material{"blue", data.Rgb{0.1, 0.2, 0.8}, data.Rgb{}, data.Rgb{}, 1.0})
	rl.floors = make(map[string]*floor)
	rl.setLevel("1")
}

// Create a custom mesh.
func (rl *rltag) makeArrowMesh() *data.Mesh {
	arrow := &data.Mesh{Name: "arrow"}
	arrow.V = []float32{
		0, 1, -0.75, 1,
		-0.45, 1, 0.45, 1,
		0.45, 1, 0.45, 1,
	}
	arrow.N = []float32{
		0, 0.5, 0.5,
		0, 0.5, 0.5,
		0, 0.5, 0.5,
	}
	arrow.F = []uint16{0, 1, 2}
	rl.eng.BindModel(arrow)
	return arrow
}

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
		flr.plan.SetPerspective(60, float32(rl.width)/float32(rl.height), 0.1, 50)
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
					block.SetFacade("cube", "gouraud", "cube")
					block.SetLocation(float32(x), 0, float32(-y))

					// use flat, gray for the overlay.
					block = flr.mapPart.AddPart()
					block.SetFacade("cube", "flat", "gray")
					block.SetLocation(float32(x), 0, float32(-y))
				}
			}
		}
		flr.arrow = flr.mapPart.AddPart()
		flr.arrow.SetFacade("arrow", "flat", "blue")
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

// Handle engine callbacks.
func (rl *rltag) Focus(focus bool) {}
func (rl *rltag) Resize(x, y, width, height int) {
	rl.eng.ResizeViewport(x, y, width, height)
	rl.width = width
	rl.height = height
	ratio := float32(width) / float32(height)
	for _, flr := range rl.floors {
		flr.plan.SetPerspective(60, ratio, 0.1, 50)
	}
}
func (rl *rltag) Update(pressed []string, gt, dt float32) {

	// pre-process user presses.
	// reduce individual move amounts for multiple move requests.
	moveDelta := rl.eng.Dt * 2
	for _, p := range pressed {
		switch p {
		case "W", "S", "Q", "E":
			moveDelta *= 0.5
		}
	}

	// process user presses.
	for _, p := range pressed {
		switch p {
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
			rl.flr.plan.PanView(vu.YAxis, rl.eng.Dt*rl.spin)
			rl.arrow.SetRotation(rl.flr.plan.ViewRotation())
		case "D":
			rl.flr.plan.PanView(vu.YAxis, rl.eng.Dt*-rl.spin)
			rl.arrow.SetRotation(rl.flr.plan.ViewRotation())
		case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0":
			if time.Now().After(rl.last.Add(rl.holdoff)) {
				rl.setLevel(p)
				rl.last = time.Now()
			}
		}
	}
}
