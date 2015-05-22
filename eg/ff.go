// Copyright © 2014-2015 Galvanized Logic Inc.
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
	if err := vu.New(ff, "Flow Field", 400, 100, 750, 750); err != nil {
		log.Printf("ff: error starting engine %s", err)
	}
	defer catchErrors()
}

// Globally unique "tag" that encapsulates example specific data.
type fftag struct {
	top     vu.Pov    // transform hierarchy root.
	chasers []vu.Pov  // map chasers.
	goal    vu.Pov    // chasers goal.
	mmap    vu.Pov    // allows the main map to be moved around.
	cam     vu.Camera // how its drawn on the minimap.
	msize   int       // map width and height.
	spots   []int     // unique ids of open spots.
	plan    grid.Grid // the floor layout.
	flow    grid.Flow // the flow field.
}

// Create is the engine callback for initial asset creation.
func (ff *fftag) Create(eng vu.Eng, s *vu.State) {
	rand.Seed(time.Now().UTC().UnixNano())

	// create the overlay
	ff.top = eng.Root().NewPov()
	view := ff.top.NewView()
	view.SetUI()
	ff.cam = view.Cam()
	ff.mmap = ff.top.NewPov().SetScale(10, 10, 0)
	ff.mmap.SetLocation(30, 30, 0)

	// populate the map
	ff.msize = 69
	ff.plan = grid.New(grid.ROOMS_SKIRMISH)
	ff.plan.Generate(ff.msize, ff.msize)
	width, height := ff.plan.Size()
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if ff.plan.IsOpen(x, y) {
				block := ff.mmap.NewPov()
				block.SetLocation(float64(x), float64(y), 0)
				block.NewModel("uv").LoadMesh("icon").AddTex("wall")
				ff.spots = append(ff.spots, ff.id(x, y))
			}
		}
	}

	// populate chasers and a goal.
	numChasers := 30
	for cnt := 0; cnt < numChasers; cnt++ {
		chaser := ff.mmap.NewPov()
		chaser.NewModel("uv").LoadMesh("icon").AddTex("token")
		ff.chasers = append(ff.chasers, chaser)
	}
	ff.goal = ff.mmap.NewPov()
	ff.goal.NewModel("uv").LoadMesh("icon").AddTex("goal")
	ff.flow = grid.NewFlow(ff.plan) // flow field for the given plan.
	ff.resetLocations()

	// set non default engine state.
	eng.SetColor(0.15, 0.15, 0.15, 1)
	ff.resize(s.W, s.H)
}

// Update is the regular engine callback.
func (ff *fftag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	if in.Resized {
		ff.resize(s.W, s.H)
	}

	// FUTURE: Adjust the chaser locations to be closer to the center.

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

	// create the flow field based on the given goal.
	ff.flow.Create(goalx, goaly)
}

// Turn x,y map indicies to unique identifiers.
func (ff *fftag) id(x, y int) int { return x*ff.msize + y }

// Turn unique identifiers to x,y map indicies.
func (ff *fftag) at(id int) (x, y int) { return id / ff.msize, id % ff.msize }

func (ff *fftag) resize(w, h int) {
	ff.cam.SetOrthographic(0, float64(w), 0, float64(h), 0, 10)
}
