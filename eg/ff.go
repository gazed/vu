// Copyright © 2014-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gazed/vu"
	"github.com/gazed/vu/grid"
)

// ff demonstrates flow field path finding by having a bunch of chasers and
// a goal. The goal is randomly reset once all the chasers have reached it.
// Restarting the example will create a different random grid.
// See vu/grid/flow.go for more information on flow fields.
//
// CONTROLS:
//   Sp    : pause while pressed
func ff() {
	ff := &fftag{}
	if err := vu.New(ff, "Flow Field", 400, 100, 750, 750); err != nil {
		log.Printf("ff: error starting engine %s", err)
	}
	defer catchErrors()
}

// Globally unique "tag" that encapsulates example specific data.
type fftag struct {
	top     *vu.Pov    // transform hierarchy root.
	chasers []*chaser  // map chasers.
	goal    *vu.Pov    // chasers goal.
	mmap    *vu.Pov    // allows the main map to be moved around.
	cam     *vu.Camera // how its drawn on the minimap.
	msize   int        // map width and height.
	spots   []int      // unique ids of open spots.
	plan    grid.Grid  // the floor layout.
	flow    grid.Flow  // the flow field.
}

// Create is the engine callback for initial asset creation.
func (ff *fftag) Create(eng vu.Eng, s *vu.State) {
	rand.Seed(time.Now().UTC().UnixNano())

	// create the overlay
	ff.top = eng.Root().NewPov()
	ff.cam = ff.top.NewCam().SetUI()
	ff.mmap = ff.top.NewPov().SetScale(10, 10, 0)
	ff.mmap.SetAt(30, 30, 0)

	// populate the map
	ff.msize = 69
	ff.plan = grid.New(grid.RoomSkirmish)
	ff.plan.Generate(ff.msize, ff.msize)
	width, height := ff.plan.Size()
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if ff.plan.IsOpen(x, y) {
				block := ff.mmap.NewPov()
				block.SetAt(float64(x), float64(y), 0)
				block.NewModel("uv", "msh:icon", "tex:wall")
				ff.spots = append(ff.spots, ff.id(x, y))
			}
		}
	}

	// populate chasers and a goal.
	numChasers := 30
	for cnt := 0; cnt < numChasers; cnt++ {
		ff.chasers = append(ff.chasers, newChaser(ff.mmap))
	}
	ff.goal = ff.mmap.NewPov()
	ff.goal.NewModel("uv", "msh:icon", "tex:goal")
	ff.flow = grid.NewFlow(ff.plan) // flow field for the given plan.
	ff.resetLocations()

	// set non default engine state.
	eng.Set(vu.Color(0.15, 0.15, 0.15, 1))
	ff.resize(s.W, s.H)
}

// Update is the regular engine callback.
func (ff *fftag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	if in.Resized {
		ff.resize(s.W, s.H)
	}

	// pause with space bar.
	if _, ok := in.Down[vu.KSpace]; !ok {
		// move each of the chasers closer to the goal.
		reset := true
		for _, chaser := range ff.chasers {
			if chaser.move(ff.flow) {
				reset = false
			}
		}
		if reset {
			ff.resetLocations()
		}
	}
}

// resetLocations randomly distributes the chasers around the map.
func (ff *fftag) resetLocations() {
	for _, chaser := range ff.chasers {
		spot := ff.spots[rand.Intn(len(ff.spots))] // get open location.
		chaser.gx, chaser.gy = ff.at(spot)         // get map location.
		chaser.nx, chaser.ny = chaser.gx, chaser.gy
		chaser.pov.SetAt(float64(chaser.gx), float64(chaser.gy), 0)
	}
	spot := ff.spots[rand.Intn(len(ff.spots))]
	goalx, goaly := ff.at(spot)
	ff.goal.SetAt(float64(goalx), float64(goaly), 0)

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

// =============================================================================

// chasers move from grid location to grid location until they
// reach the goal.
type chaser struct {
	pov    *vu.Pov // actual location.
	gx, gy int     // old grid location.
	nx, ny int     // next grid location.
	cx, cy int     // optional center to avoid when moving.
}

// chaser moves towards a goal.
func newChaser(parent *vu.Pov) *chaser {
	c := &chaser{}
	c.pov = parent.NewPov()
	c.pov.NewModel("uv", "msh:icon", "tex:token")
	return c
}

// move the chaser a bit closer to its goal.
func (c *chaser) move(flow grid.Flow) (moved bool) {
	sx, sy, _ := c.pov.At() // actual screen location.
	atx := math.Abs(float64(sx-float64(c.nx))) < 0.05
	aty := math.Abs(float64(sy-float64(c.ny))) < 0.05
	if atx && aty { // reached next location.
		c.gx, c.gy = c.nx, c.ny
		nx, ny := flow.Next(c.gx, c.gy)
		if nx == 9 {
			return false // no valid moves for this chaser.
		}
		if nx == 0 && ny == 0 {
			return false // reached the goal
		}
		moved = true
		c.nx, c.ny = c.gx+nx, c.gy+ny
		c.pov.SetAt(float64(c.gx), float64(c.gy), 0)

		// check if the chaser path should go around a corner.
		c.cx, c.cy = 0, 0
		switch {
		case nx == -1 && ny == 1:
			if ax, _ := flow.Next(c.gx-1, c.gy); ax == 9 {
				c.cx, c.cy = c.gx-1, c.gy
			}
			if bx, _ := flow.Next(c.gx, c.gy+1); bx == 9 {
				c.cx, c.cy = c.gx, c.gy+1
			}
		case nx == -1 && ny == -1:
			if ax, _ := flow.Next(c.gx-1, c.gy); ax == 9 {
				c.cx, c.cy = c.gx-1, c.gy
			}
			if bx, _ := flow.Next(c.gx, c.gy-1); bx == 9 {
				c.cx, c.cy = c.gx, c.gy-1
			}
		case nx == 1 && ny == 1:
			if ax, _ := flow.Next(c.gx+1, c.gy); ax == 9 {
				c.cx, c.cy = c.gx+1, c.gy
			}
			if bx, _ := flow.Next(c.gx, c.gy+1); bx == 9 {
				c.cx, c.cy = c.gx, c.gy+1
			}
		case nx == 1 && ny == -1:
			if ax, _ := flow.Next(c.gx+1, c.gy); ax == 9 {
				c.cx, c.cy = c.gx+1, c.gy
			}
			if bx, _ := flow.Next(c.gx, c.gy-1); bx == 9 {
				c.cx, c.cy = c.gx, c.gy-1
			}
		}
	} else {
		moved = true
		speed := 0.1 // move a bit closer to the next spot.
		if c.cx == 0 && c.cy == 0 {
			// move in a straight line.
			if !atx {
				sx += float64(c.nx-c.gx) * speed
			}
			if !aty {
				sy += float64(c.ny-c.gy) * speed
			}
			c.pov.SetAt(sx, sy, 0)
		} else {
			// move in a straight line.
			if !atx {
				sx += float64(c.nx-c.gx) * speed
			}
			if !aty {
				sy += float64(c.ny-c.gy) * speed
			}

			// push the point out so that it moves around the
			// circle radius of center cx, cy.
			dx, dy := sx-float64(c.cx), sy-float64(c.cy)
			vlen := math.Sqrt(dx*dx + dy*dy) // vector length.
			if vlen != 0 {
				dx, dy = dx/vlen, dy/vlen // unit vector.
			}
			radius := 1.0
			sx, sy = float64(c.cx)+dx*radius, float64(c.cy)+dy*radius
			c.pov.SetAt(sx, sy, 0)
		}
	}
	return moved
}
