// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

package main

import (
	"log/slog"
	"math"
	"math/rand"
	"time"

	"github.com/gazed/vu"
	"github.com/gazed/vu/grid"
)

// ff demonstrates flow field pathing by having a bunch of chasers and a goal.
// The goal is randomly reset once all the chasers have reached it.
// See vu/grid/flow.go for more information on flow fields.
//
// CONTROLS:
// - Space: pause while pressed
func ff() {
	defer catchErrors()
	ff := &fftag{}
	eng, err := vu.NewEngine(
		vu.Windowed(),
		vu.Title("Flow Field"),
		vu.Size(400, 100, ffScreenSize, ffScreenSize),
		vu.Background(0.15, 0.15, 0.15, 1.0),
	)
	if err != nil {
		slog.Error("ff: engine start", "err", err)
		return
	}

	// Run will call Load once and then call Update each engine tick.
	eng.Run(ff, ff) // does not return while example is running.
}

// Globally unique "tag" that encapsulates example specific data.
type fftag struct {
	ui        *vu.Entity // transform hierarchy root.
	chaseRoot *vu.Entity // parent of chaser instances.
	chasers   []*chaser  // map chasers.
	goal      *vu.Entity // chasers goal.
	level     *vu.Entity // ui level for grid.
	spots     []uint     // unique ids of open spots.
	grid      grid.Grid  // the floor layout.
	flow      *grid.Flow // the flow field.
}

const (
	ffGridSize   = 69   // width and height grid size.
	ffScreenSize = 1200 // screen size in pixels.
)

// Load is the engine callback for initial asset creation.
func (ff *fftag) Load(eng *vu.Engine) error {
	loadTestGrid()

	// import assets from asset files.
	// This creates the assets referenced by the models below.
	// Note that circle and quad meshes engine defaults.
	eng.ImportAssets("color2D.shd", "tint2D.shd") // shaders
	eng.ImportAssets("star.png", "ring.png")      // textures

	// create the 2D overlay
	mazeOffset := 30.0                                   // in pixels.
	mazePixels := float64(ffScreenSize) - mazeOffset*2.0 // screenSize - border
	mazeScale := mazePixels / float64(ffGridSize-1)      // scale by grid size.
	ff.ui = eng.AddScene(vu.Scene2D)
	ff.level = ff.ui.AddPart().SetScale(mazeScale, mazeScale, 0)
	ff.level.SetAt(30, 30, 0)

	// populate the level models
	ff.grid = loadTestGrid()
	size, _ := ff.grid.Size()
	for y := range size {
		for x := range size {
			if ff.grid[y][x] == grid.SOLID {
				block := ff.level.AddModel("shd:color2D", "msh:icon").SetColor(0, 0, 0, 1)
				block.SetAt(float64(x), float64(y), 0)
			} else {
				// open spots used to navigate.
				ff.spots = append(ff.spots, ff.grid.Index(y, x))
			}
		}
	}

	// populate chasers and a goal.
	numChasers := 30
	for range numChasers {
		ff.chasers = append(ff.chasers, newChaser(ff.level))
	}
	ff.goal = ff.level.AddModel("shd:tint2D", "msh:icon", "tex:color:star").SetColor(1, 1, 0, 1)
	ff.resetLocations()
	return nil
}

// FrameUpdate is the pre render frame engine callback
func (ff *fftag) FrameUpdate(eng *vu.Engine, in *vu.Input, delta time.Duration, interpolate float64) {
	// react to one time press events.
	for press := range in.Pressed {
		switch press {
		case vu.KQ, vu.KEsc:
			eng.Shutdown()
			return
		}
	}
}

// FixedUpdate is the pre-physics fixed timestep engine callback.
func (ff *fftag) FixedUpdate(eng *vu.Engine, in *vu.Input, timestep time.Duration) {
	if _, ok := in.Down[vu.KSpace]; ok {
		return // pause with space bar.
	}

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

// resetLocations randomly distributes the chasers and
// the goal around the map.
func (ff *fftag) resetLocations() {
	empty := append([]uint{}, ff.spots...)
	for _, chaser := range ff.chasers {
		index := rand.Intn(len(empty))  //
		spot := empty[index]            // get open location.
		gy, gx := ff.grid.Indices(spot) // get map location.
		chaser.gy, chaser.gx = int(gy), int(gx)
		chaser.nx, chaser.ny = chaser.gx, chaser.gy
		chaser.pov.SetAt(float64(chaser.gx), float64(chaser.gy), 0)
		empty = append(empty[:index], empty[index+1:]...)
	}
	spot := empty[rand.Intn(len(empty))]
	goaly, goalx := ff.grid.Indices(spot)
	ff.goal.SetAt(float64(goalx), float64(goaly), 0)

	// create the flow field based on the given goal.
	ff.flow = grid.NewFlow(ff.grid, goaly, goalx) // grid flow field for the given plan.
	// ff.flow.Dump()
}

// =============================================================================

// chasers move from grid location to grid location until they
// reach the goal.
type chaser struct {
	pov    *vu.Entity // actual location.
	gx, gy int        // old grid location.
	nx, ny int        // next grid location.
	cx, cy int        // optional center to avoid when moving.
}

// chaser moves towards a goal.
func newChaser(parent *vu.Entity) *chaser {
	c := &chaser{}
	c.pov = parent.AddModel("shd:tint2D", "msh:icon", "tex:color:star").SetColor(1, 0, 0, 1)
	return c
}

// move the chaser a bit closer to its goal.
func (c *chaser) move(flow *grid.Flow) (moved bool) {
	sx, sy, _ := c.pov.At() // actual screen location.
	atx := math.Abs(float64(sx-float64(c.nx))) < 0.05
	aty := math.Abs(float64(sy-float64(c.ny))) < 0.05
	if atx && aty { // reached next location.
		c.gx, c.gy = c.nx, c.ny
		nx, ny := flow.Next(c.gy, c.gx)
		if nx == 9 {
			return false // no valid moves for this chaser.
		}
		if nx == 0 && ny == 0 {
			return false // reached the goal
		}

		// chaser starts moving towards next goal.
		c.nx, c.ny = c.gx+nx, c.gy+ny
		c.pov.SetAt(float64(c.gx), float64(c.gy), 0)
	} else {

		// move in a straight line towards next goal.
		speed := 0.05 // move a bit closer to the next spot.
		if !atx {
			sx += float64(c.nx-c.gx) * speed
		}
		if !aty {
			sy += float64(c.ny-c.gy) * speed
		}
		c.pov.SetAt(sx, sy, 0)
	}
	return true // moved
}

// loadTestGrid is a utlity method that generates a grid.Grid from
// the testGrid representation.
func loadTestGrid() grid.Grid {
	solidRune, emptyRune := '◾', '◽'
	g := grid.MakeGrid(ffGridSize, ffGridSize) // testGrid size.
	index := uint(0)                           // sequential grid index.
	for _, c := range testGrid {
		if c == solidRune {
			r, c := g.Indices(index) // grid loction from index.
			g[r][c] = grid.SOLID
			index += 1
		}
		if c == emptyRune {
			index += 1 // grid default is already grid.EMPTY
		}
	}
	return g
}

// size matches ffGridSize == 69
// 69x69 pretty square grid representation from a grid generation algorithm
// in the old vu. Grid generation is currently not part of the core vu engine.
var testGrid = ` 
◾◾◾◾◾◾◽◾◾◽◾◾◾◽◾◾◾◾◾◾◾◾◾◾◾◾◾◾◾◽◾◾◾◾◽◾◾◾◽◽◾◾◾◾◾◾◾◽◾◾◾◾◾◾◽◾◾◾◾◾◾◽◾◾◾◾◾◾◾
◾◽◽◽◾◽◽◽◾◽◽◾◽◽◽◾◽◽◽◽◽◾◽◽◽◽◾◽◽◽◽◽◽◾◽◽◽◽◽◽◽◾◽◽◽◾◽◽◾◽◽◾◽◽◽◽◽◽◽◽◾◽◽◽◾◽◽◽◾
◽◽◽◽◾◽◽◽◽◽◽◾◽◽◽◾◽◽◽◽◽◾◽◽◽◽◾◽◽◽◽◽◽◾◽◽◽◽◾◽◽◽◽◽◽◾◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◾
◾◽◽◾◾◾◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◾◽◽◽◽◾◽◽◽◽◽◽◾◽◽◽◽◾◽◽◽◽◽◽◾◽◽◽◽◽◾◾◾◾◾◾◾◾◾◾◽◽◽◾◽◽◽◾
◾◽◽◽◽◽◽◽◽◽◽◾◽◽◽◾◽◽◽◽◽◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◾◽◽◽◾◽◾◾◽◽◾◽◾◾◽◽◽◽◽◽◽◽◾◽◽◽◾◽◽◽◾
◾◽◽◽◾◽◽◽◾◽◽◽◾◾◽◾◽◽◽◽◽◾◽◽◽◽◾◾◽◽◾◾◾◽◽◽◾◾◾◾◾◽◽◽◽◾◽◽◾◽◽◽◽◽◽◽◾◽◽◽◾◽◽◽◾◽◽◽◾
◾◽◽◽◾◽◽◽◾◽◽◽◽◽◽◾◽◽◽◽◽◾◽◽◽◽◾◽◽◽◽◽◾◽◽◾◽◽◽◽◽◾◽◽◽◾◽◽◾◽◽◾◽◽◽◽◾◽◽◽◾◽◽◽◾◽◽◽◾
◾◽◽◽◾◾◾◾◾◽◽◾◽◽◽◾◽◽◽◽◾◾◽◽◽◽◾◽◽◽◽◽◾◽◽◾◽◽◽◽◽◾◽◽◽◾◾◽◽◽◽◽◾◾◾◽◽◽◾◾◾◾◾◾◽◽◾◾◾
◾◾◾◾◾◽◽◽◾◽◽◾◽◽◽◽◽◽◾◽◽◾◾◾◾◽◽◾◾◾◽◽◾◽◽◾◽◽◽◽◽◾◽◽◽◾◽◽◽◾◾◽◽◽◾◽◽◽◽◽◽◾◽◽◾◽◽◽◾
◾◽◽◽◾◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◽◾◽◽◽◽◽◽◽◾◽◽◾◽◽◽◽◽◽◽◽◾◾◽◾◾◽◽◽◽◽◾◽◽◾◽◽◾◽◽◽◾◽◽◾◽◽◽◾
◾◽◽◽◽◽◽◽◾◽◽◽◾◾◾◾◽◽◽◽◽◽◽◽◽◽◽◽◽◾◽◽◾◽◽◽◽◽◽◽◽◾◽◽◽◾◽◽◾◽◽◾◽◽◾◽◽◾◽◽◽◾◽◽◾◽◽◽◾
◾◽◽◽◾◽◾◾◾◽◽◽◽◽◽◾◽◽◾◽◽◽◾◾◾◾◾◽◽◾◽◾◾◽◾◾◽◽◽◽◽◾◽◽◽◾◾◾◾◽◽◽◽◾◾◽◽◾◽◽◽◾◽◽◾◽◽◽◽
◽◽◽◽◾◽◽◽◽◾◾◾◽◽◽◾◾◾◾◾◾◾◽◽◽◽◾◽◽◾◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◾◽◽◾◽◽◾◽◽◾◽◽◽◽◽◽◾◽◽◾◽◽◽◾
◾◾◽◾◾◽◽◽◽◽◽◾◾◾◽◽◽◽◾◽◽◾◽◽◽◽◾◽◽◾◽◽◽◽◽◽◾◾◾◾◽◾◾◾◽◽◽◽◾◽◽◾◽◽◽◽◽◽◽◽◽◽◽◽◾◽◽◽◾
◾◽◽◽◾◽◽◽◾◽◽◾◽◽◽◽◽◽◽◽◽◾◽◽◽◽◾◽◾◾◽◽◽◽◽◽◽◽◽◽◽◾◽◽◽◾◽◽◾◽◽◾◽◾◽◽◾◾◾◽◾◽◾◽◽◽◽◾◾
◽◽◽◽◾◽◽◾◾◽◽◾◽◽◽◽◽◽◾◽◽◽◽◽◽◽◾◽◽◾◽◽◽◽◽◾◽◽◽◽◽◾◽◽◽◾◽◽◾◽◽◾◽◽◽◽◽◾◽◽◽◽◽◽◽◾◽◽◾
◾◽◽◽◾◽◽◽◾◽◽◾◽◽◽◾◽◽◾◽◾◾◽◾◾◾◾◽◽◾◾◾◾◽◽◾◽◽◽◽◽◾◽◽◽◾◽◽◾◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◽◽◽
◾◽◾◾◾◽◽◽◽◽◽◾◽◽◽◾◽◽◾◽◽◾◽◽◽◽◽◽◽◾◽◽◾◽◽◽◽◽◽◽◽◾◽◽◽◽◽◽◾◽◽◾◽◽◾◽◽◽◽◽◽◽◾◽◽◾◽◽◾
◾◽◽◽◽◾◾◽◽◽◽◾◽◽◽◾◽◽◽◽◽◽◽◽◽◽◾◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◽◾◾◽◽◾◾◽◽◾◽◽◾◽◽◾◽◽◽◽◾◽◽◾◽◽◾
◾◽◽◽◽◽◽◽◽◽◽◾◽◽◽◾◽◽◾◽◾◾◽◾◾◾◾◽◾◽◽◾◾◾◽◽◾◾◽◾◾◾◽◽◽◾◽◽◽◽◽◾◾◾◾◾◽◽◾◾◾◽◾◽◽◾◽◽◾
◾◽◾◾◾◽◽◽◾◽◾◾◽◽◾◾◽◽◽◽◽◽◽◽◽◾◽◽◾◽◽◾◽◽◽◾◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◾◽◽◽◽◽◽◽◾◽◽◾◽◽◾◽◽◾
◾◽◽◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◾◽◽◽◾◽◽◽◽◽◾◾◾◾◾◽◽◾◽◽◾◽◽◽◽◽◽◽◽◽◽◽◽◽◾◽◽◾
◾◽◽◽◽◾◾◽◽◽◽◽◽◽◽◽◽◽◽◽◾◽◽◽◽◾◽◽◾◽◽◾◽◽◽◽◽◽◽◽◽◾◽◽◽◾◾◾◾◾◾◾◾◾◾◽◾◾◾◾◾◾◽◾◾◾◾◽◽
◾◽◽◽◾◽◽◽◾◾◾◾◽◽◽◾◾◾◾◾◾◾◽◾◾◽◾◽◾◽◽◾◽◽◽◾◾◽◽◾◾◽◽◽◽◾◽◽◽◽◾◽◽◽◾◽◽◽◽◽◽◾◽◽◾◽◽◽◽
◾◽◽◽◾◽◽◽◽◽◽◾◾◾◾◽◽◽◽◽◾◽◽◽◽◽◽◽◾◽◽◽◽◽◽◾◽◽◽◽◽◽◾◾◽◾◽◽◽◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◾◽◽◽◾
◾◾◽◾◾◽◽◽◾◽◽◾◽◽◽◽◽◽◽◽◾◽◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◾◽◽◽◾◽◽◽◽◾◽◽◽◾◾◾◾◾◾◾◾◽◾◾◾◾◾◾
◾◽◽◽◾◾◾◾◾◾◽◽◽◽◽◾◽◾◾◾◾◽◾◾◾◾◾◾◾◽◾◾◾◾◾◾◾◾◾◾◾◾◽◽◽◾◽◽◽◽◾◽◽◽◾◽◽◽◽◽◽◽◽◾◽◽◽◽◾
◾◽◽◽◽◽◽◽◽◽◽◾◽◾◾◾◽◽◽◽◾◽◽◽◽◽◽◽◾◽◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◾
◾◽◽◽◽◽◽◽◽◽◽◾◽◽◽◾◽◽◽◽◾◽◽◽◽◽◽◽◾◽◽◽◽◽◽◾◽◽◽◽◽◾◽◽◽◾◽◾◾◾◾◽◾◾◾◽◾◾◾◾◾◽◾◾◾◾◽◾◾
◾◽◽◽◾◽◽◽◽◽◽◾◽◽◽◾◽◽◽◽◾◽◽◽◽◾◽◽◾◽◽◾◽◽◽◽◽◽◽◽◽◾◽◾◾◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◾
◾◾◾◾◾◽◽◽◽◽◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◾◽◽◾◾◽◾◾◽◽◾◾◾◾◾◽◾◽◽◽◽◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◾
◾◽◽◽◾◾◽◾◾◾◾◾◽◽◽◾◽◽◽◽◽◽◽◽◽◾◾◽◽◽◽◽◾◽◽◽◽◽◽◽◽◾◽◽◽◾◾◾◾◾◾◾◾◾◾◾◾◾◾◾◽◽◾◾◾◾◾◽◾
◽◽◽◽◽◽◽◽◽◽◽◾◾◾◽◾◽◽◽◽◾◾◾◾◾◾◽◽◾◽◽◽◽◽◽◽◽◽◽◽◽◾◾◾◾◽◽◽◽◽◾◽◽◽◽◽◽◾◽◽◽◽◽◽◽◾◽◽◾
◾◽◽◽◾◽◽◽◽◽◽◽◽◽◽◾◽◽◽◽◾◽◽◽◽◾◽◽◽◽◽◽◾◽◽◾◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◾◽◽◽◽◽◽◽◾◽◽◾
◾◽◽◽◾◾◾◾◽◽◾◾◽◽◽◾◾◾◽◾◾◽◽◽◽◾◽◽◽◽◽◽◾◽◽◾◽◽◽◽◽◽◽◽◽◾◾◽◽◾◾◽◾◾◾◽◾◾◾◾◽◾◾◽◽◾◽◽◽
◾◾◽◾◽◽◽◽◾◽◽◾◾◽◽◾◽◽◽◽◾◽◽◽◽◾◽◾◾◽◽◽◾◾◾◾◽◽◽◽◽◾◽◾◾◾◽◽◽◽◽◽◽◽◾◽◽◽◾◽◽◽◽◽◽◾◽◽◾
◾◽◽◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◾◽◽◽◽◾◽◽◾◽◽◽◾◽◽◾◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◾◽◽◾◽◽◽◽◽◽◽◾◽◽◾◽◽◾
◾◽◽◽◾◾◽◽◾◾◾◽◽◽◽◽◽◽◽◽◽◽◽◽◽◾◽◽◾◽◽◾◾◽◽◾◾◽◾◾◾◾◽◽◽◽◾◾◾◾◾◾◾◽◽◾◾◽◽◾◾◾◾◾◽◾◾◽◽
◽◽◽◽◾◽◽◽◾◽◽◾◾◽◾◾◾◾◾◾◽◾◽◾◾◾◽◽◽◽◽◽◽◽◽◾◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◽◾◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◽
◽◽◾◾◾◽◽◽◽◽◽◾◽◽◽◽◽◽◾◽◽◽◽◽◽◾◽◽◾◽◽◽◽◽◽◾◽◽◾◽◽◾◽◽◽◾◽◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◾
◾◽◽◽◾◾◾◾◽◽◾◾◽◽◽◽◽◽◾◽◽◽◽◽◽◾◾◾◾◽◽◽◾◾◽◾◽◽◾◽◽◾◾◽◾◾◽◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◽◾◽◽◽◽◾
◾◽◽◽◽◽◽◾◽◽◽◾◽◽◽◾◽◽◾◽◽◽◽◽◽◽◽◽◾◾◾◾◾◽◽◽◽◽◾◽◽◾◽◽◽◾◽◽◽◽◽◽◾◽◽◽◾◽◽◾◾◾◾◾◾◾◽◽◾
◾◽◽◾◾◽◽◾◽◽◽◾◽◽◽◾◽◽◾◽◽◽◾◽◽◾◽◽◾◽◽◽◾◽◽◽◽◽◾◽◾◾◽◽◽◾◽◽◽◽◽◽◾◽◽◽◾◽◽◽◽◽◾◽◽◽◽◽◽
◾◽◽◽◾◽◽◾◽◽◽◽◾◾◾◾◽◾◾◾◾◾◾◾◽◾◽◽◽◽◽◽◾◾◽◾◽◽◾◽◽◾◽◽◽◾◽◽◽◽◽◽◾◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◾
◾◽◽◽◽◽◽◾◾◾◾◽◽◽◽◾◽◽◽◽◽◽◽◽◽◾◽◽◾◾◽◾◾◽◽◽◽◽◾◽◽◾◽◽◽◾◽◽◾◾◾◾◽◽◾◾◾◾◾◾◾◾◾◾◾◾◽◽◾
◾◽◽◾◾◽◽◾◽◽◽◾◽◽◽◾◽◽◾◽◽◽◾◽◽◾◽◾◾◽◽◽◾◽◽◽◽◽◽◽◽◾◽◽◽◽◽◽◾◽◽◾◽◽◽◽◽◽◽◽◽◽◾◽◽◽◽◽◾
◾◽◽◽◽◽◽◽◽◽◽◾◽◽◽◾◽◽◾◽◽◽◾◽◽◾◽◽◾◽◽◽◾◾◾◾◾◾◾◽◾◾◽◽◽◾◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◽◾◽◽◽◽◽◽
◾◽◽◽◽◽◽◾◽◽◽◾◽◽◽◾◽◽◾◽◽◽◾◽◽◾◽◽◾◽◽◽◾◽◽◾◽◽◾◽◽◾◽◽◽◾◾◾◽◾◽◾◽◽◽◽◽◽◾◽◽◽◾◽◽◽◽◽◾
◾◾◾◾◾◽◽◾◾◾◽◾◽◾◾◾◾◾◾◾◾◽◾◽◾◾◽◽◾◽◽◽◽◽◽◾◽◽◾◽◽◽◽◾◾◾◽◽◽◽◽◾◽◽◽◽◽◽◾◽◽◽◾◽◽◽◽◽◾
◾◽◽◽◾◽◽◽◽◽◽◽◽◽◽◾◽◽◽◽◽◽◾◽◽◾◽◽◽◽◽◽◾◽◽◽◽◽◾◽◽◾◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◽◾◽◽◽◾◽◽◽◽◽◾
◾◽◽◽◾◽◽◽◽◽◽◾◽◽◽◾◽◽◽◽◽◽◾◽◽◽◽◽◾◽◽◽◾◽◽◾◽◾◾◽◽◾◽◽◽◽◽◽◽◽◽◾◽◽◽◽◽◽◾◽◽◽◾◽◽◽◽◽◾
◾◽◽◽◾◾◽◾◽◽◾◾◽◽◽◽◽◽◽◽◽◽◾◽◽◽◽◽◾◽◽◽◾◽◽◾◽◽◾◽◽◾◾◾◾◾◾◾◽◾◾◾◽◾◾◾◾◾◾◽◾◾◾◾◾◾◽◾◾
◽◽◽◽◾◽◽◾◽◽◽◽◽◾◾◾◽◽◽◽◽◽◾◽◽◾◽◽◾◽◽◽◾◽◽◾◽◽◽◽◽◽◽◽◽◾◽◽◽◽◾◽◽◾◽◽◽◽◾◽◽◽◾◽◽◾◽◽◾
◾◽◽◽◾◽◽◽◽◽◽◾◽◽◽◾◾◾◾◾◾◽◾◽◾◾◾◽◽◾◾◽◾◽◾◾◾◾◽◾◾◾◽◽◽◾◽◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◾◽◽◾◽◽◾
◾◽◽◽◾◾◽◽◾◾◽◽◽◽◽◽◽◽◾◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◾◾◾◽◽◽◽◽◽◾◽◽◽◽◽◽◽◾◾◽◾◾◽◽◾◽◽◾
◾◽◽◾◾◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◾◽◽◽◾◽◽◽◽◾◽◽◽◽◽◽◾◽◽◾◽◽◽◾◽◽◽◽◾◽◽◾◽◽◽◽◽◽◽◽◾◽◽◾◽◽◾
◾◽◽◽◾◽◽◽◽◽◽◾◽◽◽◾◽◾◽◾◽◾◾◽◽◾◾◾◾◾◾◽◽◾◽◽◾◾◾◾◽◾◽◽◽◾◽◽◽◽◽◽◽◾◽◽◽◽◾◽◽◽◽◽◽◾◽◽◽
◾◽◽◽◾◾◾◾◾◾◾◾◽◽◽◾◽◽◽◽◽◽◽◾◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◾◽◽◽◽◾◾◾◾◽◾◽◾◾◾◽◽◾◾◽◾◾◽◾◾◾◾◾
◾◽◽◽◽◽◽◽◽◽◽◾◾◾◾◾◽◽◾◽◽◽◽◾◽◽◽◽◽◾◽◽◾◽◽◽◽◽◾◽◽◾◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◾◽◽◽◾◽◽◽◽◽◽◾
◾◾◾◾◽◽◽◽◽◽◽◾◽◽◽◾◾◾◾◽◽◽◽◾◽◽◽◽◽◾◽◽◾◽◽◽◽◽◾◽◽◽◽◾◾◾◽◽◽◽◽◽◽◽◽◽◽◾◽◽◽◾◽◽◽◽◽◽◾
◽◽◽◽◾◽◽◽◽◽◽◾◽◽◽◾◽◽◾◽◽◽◽◾◽◽◽◽◽◽◽◽◾◽◽◽◽◽◾◽◽◾◽◽◽◾◽◽◽◽◽◾◽◽◾◽◽◽◽◽◽◾◽◽◽◽◽◽◾
◾◽◽◽◾◾◾◾◾◽◾◾◽◽◽◽◽◽◾◽◽◽◽◾◽◽◽◽◽◽◾◾◾◾◾◾◾◾◾◽◾◾◽◽◽◾◽◽◽◽◽◾◽◽◾◽◽◽◽◽◽◾◽◽◽◽◽◽◾
◾◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◾◽◾◽◾◾◾◾◾◽◽◽◾◾◽◽◽◽◽◾◽◽◽◽◽◾◽◽◽◾◽◾◾◾◾◾◽◽◽◾◾◾◽◽◽◾◽◽◽◽◽◽◽
◾◾◽◾◾◽◽◽◽◽◽◾◽◽◾◽◽◽◾◽◽◾◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◾◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◾◽◽◽◾◽◽◽◽◽◽◾
◾◽◽◽◾◽◽◽◽◽◽◾◽◽◽◾◽◽◾◽◽◽◽◽◽◽◾◽◽◽◽◽◽◽◽◾◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◾◽◽◾◽◽◾◽◽◽◾◽◽◽◽◽◽◾
◾◽◽◽◾◽◽◽◽◽◽◾◽◽◽◽◽◾◾◽◽◽◽◾◾◾◾◾◾◾◾◽◽◾◽◾◾◾◾◾◽◾◽◽◽◽◾◽◾◾◾◾◾◾◾◾◽◽◽◾◾◾◾◾◾◽◽◾◾
◾◽◽◽◽◽◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◾◽◽◾◽◽◽◽◽◽◽◽◾◽◽◽◾◽◽◾◽◽◽◽◾◽◽◽◽◽◽◽◽◾◽◽◽◽◽◾
◽◽◽◽◾◽◽◽◽◽◽◾◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽◾◽◽◽◽◽◽◽◾◽◽◽◽◽◽◽◽◽◽◽◽◽◽◽
◾◾◽◾◾◾◾◾◾◾◽◾◾◾◾◾◾◾◾◾◾◾◾◾◾◾◾◽◾◾◾◾◽◾◾◾◾◾◾◾◾◽◾◾◾◾◾◾◾◽◾◾◾◾◾◾◾◾◾◾◾◾◾◾◾◾◾◾◾
`
