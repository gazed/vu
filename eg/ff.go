// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/gazed/vu"
	"github.com/gazed/vu/grid"
)

// ff demonstrates flow field path finding by having a bunch of chasers and
// a goal. The goal is randomly reset once all the chasers have reached it.
// Restarting the example will create a different random grid.
func ff() {
	ff := &fftag{}
	var err error
	if ff.eng, err = vu.New("Flow Field", 400, 100, 750, 750); err != nil {
		log.Printf("ff: error intitializing engine %s", err)
		return
	}
	ff.eng.SetDirector(ff) // get user input through Director.Update()
	ff.create()            // create initial assests.
	if err = ff.eng.Verify(); err != nil {
		log.Fatalf("ff: error initializing model :: %s", err)
	}
	defer ff.eng.Shutdown()
	defer catchErrors()
	ff.eng.Action()
}

// Globally unique "tag" that encapsulates example specific data.
type fftag struct {
	eng     vu.Engine
	scene   vu.Scene  // how its drawn on the minimap.
	mmap    vu.Part   // allows the main map to be moved around.
	chasers []vu.Part // map chasers.
	goal    vu.Part   // chasers goal.
	msize   int       // map width and height.
	spots   []int     // unique ids of open spots.
	plan    grid.Grid // the floor layout.
	flow    grid.Flow // the flow field.
}

// create is the startup asset creation.
func (ff *fftag) create() {
	rand.Seed(time.Now().UTC().UnixNano())
	ff.eng.Color(0.15, 0.15, 0.15, 1)

	// create the overlay
	ff.scene = ff.eng.AddScene(vu.VO)
	ff.scene.Set2D()
	ff.mmap = ff.scene.AddPart().SetScale(10, 10, 0)
	ff.mmap.SetLocation(30, 30, 0)

	// populate the map
	ff.msize = 69
	ff.plan = grid.New(grid.ROOMS_SKIRMISH)
	ff.plan.Generate(ff.msize, ff.msize)
	width, height := ff.plan.Size()
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if ff.plan.IsOpen(x, y) {
				block := ff.mmap.AddPart()
				block.SetLocation(float64(x), float64(y), 0)
				block.SetRole("uv").SetMesh("icon").AddTex("wall")
				ff.spots = append(ff.spots, ff.id(x, y))
			}
		}
	}

	// populate chasers and a goal.
	numChasers := 30
	for cnt := 0; cnt < numChasers; cnt++ {
		chaser := ff.mmap.AddPart()
		chaser.SetRole("uv").SetMesh("icon").AddTex("token")
		ff.chasers = append(ff.chasers, chaser)
	}
	ff.goal = ff.mmap.AddPart()
	ff.goal.SetRole("uv").SetMesh("icon").AddTex("goal")

	ff.flow = grid.NewFlow(ff.plan) // flow field for the given plan.
	ff.resetLocations()
	ff.resize() // Size based on current screen size.
	return
}

// Update is the regular engine callback.
func (ff *fftag) Update(in *vu.Input) {
	if in.Resized {
		ff.resize()
	}

	// FUTURE current path follows the edge of the cell.
	// adjust the location to be closer to the center of the grid cell
	// that it is in.

	// move each of the chasers closer to the goal.
	reset := true
	speed := in.Dt * 5
	for _, chaser := range ff.chasers {
		x, y, _ := chaser.Location()
		dx, dy := ff.flow.Next(x, y)
		nextx := x + dx*speed
		nexty := y + dy*speed
		if nx, ny := ff.flow.Next(nextx, nexty); nx == 0 && ny == 0 {
			nextx = x + dx*speed // try only moving in x.
			nexty = y
			if nx, ny := ff.flow.Next(nextx, nexty); nx == 0 && ny == 0 {
				nextx = x
				nexty = y + dy*speed // try only moving in y.
			}
		}
		chaser.SetLocation(nextx, nexty, 0)
		if nextx != x || nexty != y {
			reset = false
		}
	}
	if reset {
		ff.resetLocations()
	}
}

// resize handles user screen/window changes.
func (ff *fftag) resize() {
	x, y, width, height := ff.eng.Size()
	ff.eng.Resize(x, y, width, height)
	ff.scene.Cam().SetOrthographic(0, float64(width), 0, float64(height), 0, 10)
}

// resetLocations randomly distributes the chasers around the map.
func (ff *fftag) resetLocations() {
	for _, chaser := range ff.chasers {
		spot := ff.spots[rand.Intn(len(ff.spots))] // get open location.
		x, y := ff.at(spot)                        // get map location.
		chaser.SetLocation(float64(x), float64(y), 0)
	}
	spot := ff.spots[rand.Intn(len(ff.spots))]
	goalx, goaly := ff.at(spot)
	ff.goal.SetLocation(float64(goalx), float64(goaly), 0)

	// create the flow field with the given goal.
	ff.flow.Create(goalx, goaly)
}

// Turn x,y map indicies to unique identifiers.
func (ff *fftag) id(x, y int) int { return x*ff.msize + y }

// Turn unique identifiers to x,y map indicies.
func (ff *fftag) at(id int) (x, y int) { return id / ff.msize, id % ff.msize }
