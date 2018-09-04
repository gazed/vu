// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// Package grid is used to generate layout data for random maze or skirmish
// levels. Maze levels have dead ends such that there is only one path to
// get from one spot in the maze to another.  Skirmish levels have
// no dead ends in order to provide plenty of movement options.
//
// Expected usage:
//       maze := grid.New(PRIM_MAZE)  // Create a type of grid.
//       maze.Generate(width, height) // Generate the grid.
//       for x := 0; x < width; x++ {
//          for y := 0; y < height; y++ {
//             if maze.isOpen(x, y) {
//                 // Do something with an open area.
//             } else {
//                 // Do something with a wall.
//             }
//          }
//       }
//
// Package grid is provided as part of the vu (virtual universe) 3D engine.
package grid

import (
	"fmt"
	"math/rand"
	"time"
)

// Grid generates a random floor plan where each grid cell is either a wall
// or a floor. The expected usage is to generate a random level and then to
// use the level by traversing it with the Size and IsOpen methods.
type Grid interface {
	// Size returns the current size of the grid. This will be 0, 0 if
	// Generate has not yet been called.
	Size() (width, depth int) // The current size of the plan.

	// IsOpen returns true if the cell at the given location is traversable.
	// Otherwise the cell is blocked and can be considered a wall.
	IsOpen(x, y int) bool // Return true if the given location is traversable.

	// Seed can be set to generate the same map each time. Leave it unset
	// to get a random map each time.
	Seed(seed int64)

	// Generate creates a grid full of walls and floors based on
	// the given depth and width.
	//
	// The minimum maze dimension is 7x7, and grids must be odd numbered.
	// The given grids width and height are modified, if necessary, to ensure
	// valid values.
	//
	// Generate needs to be called after Seed and before other grid methods.
	// It can be called any time to create a new grid.
	Generate(width, depth int) Grid

	// Band returns the grid depth based on concentric squares. Numbering
	// starts at 0 on the outside and increases towards the center. Using band
	// implies (makes more sense if) the grid width and height are the same.
	Band(x, y int) int
}

// Currently supported grid types that are used as the input to grid.New().
const (
	// PrimMaze is a Randomized Prim's Algorithm (see wikipedia).
	PrimMaze = iota

	// SparseSkirmish creates a skirmish grid by randomly traversing all the
	// grid locations and adding random walls as long as there are more than
	// two ways out of the grid location.
	SparseSkirmish

	// RoomSkirmish is a skirmish grid created by subdividing the area
	// recursively into rooms and chopping random exits in the room walls.
	RoomSkirmish

	// DenseSkirmish is a corridor only skirmish grid. It is a Prim's maze
	// where the dead-ends have been eliminated.  Additionally each side of
	// the grid is guaranteed to have one exit to the level exterior.
	DenseSkirmish

	// Cave produces interconnected non-square areas resembling a large
	// series of caves.
	Cave

	// Dungeon produces interconnected square areas resembling a series
	// of rooms connected by corridors.
	Dungeon
)

// ===========================================================================
// grid implements Grid

// New creates a new grid based on the given gridType.
// Returns nil if the gridType is not recognized.
func New(gridType int) Grid {
	switch gridType {
	case PrimMaze:
		return &primMaze{}
	case SparseSkirmish:
		return &sparse{}
	case RoomSkirmish:
		return &rooms{}
	case DenseSkirmish:
		return &dense{}
	case Cave:
		return &cave{}
	case Dungeon:
		return &dungeon{}
	}
	return nil
}

// The base class for a grid holds an x-by-y group of cells where each
// cell is either a wall or a floor.
type grid struct {
	cells [][]*cell // walls or open areas.
	seed  int64     // used for testing with deterministic grids.
}

// Seed uses the provide seed value to initialize the random source to a
// deterministic state.
func (g *grid) Seed(seed int64) { g.seed = seed }

// Size is the generated width and height of the grid.
func (g *grid) Size() (width, height int) {
	if len(g.cells) > 0 {
		return len(g.cells), len(g.cells[0])
	}
	return 0, 0
}

// IsOpen returns true if the cell at the given location is passable.
// Otherwise the cell is blocked and can be considered a wall.
func (g *grid) IsOpen(x, y int) bool {
	lenx := len(g.cells)
	if x >= 0 && y >= 0 && lenx > 0 && lenx > x && len(g.cells[0]) > y {
		return !g.cells[x][y].isWall
	}
	return false
}

// Band returns the concentrix square number where the outermost square
// is the zeroth band. Un-generated grids and/or invalid input coordinates
// always return 0.
func (g *grid) Band(x, y int) int {
	w, h := g.Size()
	if w > 0 && h > 0 && x >= 0 && x < w && y >= 0 && y < h {
		lowleft := x
		if lowleft > y {
			lowleft = y
		}
		topright := w - x
		if topright > h-y {
			topright = h - y
		}
		if lowleft < topright {
			return lowleft
		}
		return topright
	}
	return 0
}

// cells is the grid represented as a linear slice.
func (g *grid) cellSlice() (cells []*cell) {
	if len(g.cells) > 0 {
		for _, row := range g.cells {
			cells = append(cells, row...)
		}
	}
	return
}

// Used in create to have the default grid made entirely of walls
// or floors. Some algorithms start one way, some the other.
const (
	allFloors = false // Indicates the initial grid is all floors.
	allWalls  = true  // Indicates the initial grid is all walls.
)

// create the space needed by the grid. This is the same for all grid
// implementations. The cellType is expected to be allFloors or allWalls.
func (g *grid) create(width, height int, cellType bool) {
	gridWidth, gridHeight := g.validateSize(width), g.validateSize(height)
	g.cells = make([][]*cell, gridWidth)
	for x := range g.cells {
		g.cells[x] = make([]*cell, gridHeight)
		for y := range g.cells[x] {
			g.cells[x][y] = &cell{x, y, cellType}
		}
	}
	if g.seed == 0 {
		rand.Seed(time.Now().UnixNano())
	} else {
		rand.Seed(g.seed)
	}
}

// validateSize will return valid grid sizes. Given that some of the cells
// are used as walls, and that the grid size must be odd, the minimum grid
// size is 7x7.
func (g *grid) validateSize(size int) (validSize int) {
	validSize = size
	if size < 7 {
		validSize = 7
	} else if size%2 == 0 {
		validSize++
	}
	return
}

// cells connect up and down as north, south, east, west. Note that even though
// the "y-index" can range conceptually from left to right in an array, it is
// used as up and down as per cartesian plane.
func (g *grid) north(u *cell) *cell {
	if u.y+1 < len(g.cells[u.x]) {
		return g.cells[u.x][u.y+1]
	}
	return nil
}
func (g *grid) south(u *cell) *cell {
	if u.y-1 >= 0 {
		return g.cells[u.x][u.y-1]
	}
	return nil
}
func (g *grid) east(u *cell) *cell {
	if u.x+1 < len(g.cells) {
		return g.cells[u.x+1][u.y]
	}
	return nil
}
func (g *grid) west(u *cell) *cell {
	if u.x-1 >= 0 {
		return g.cells[u.x-1][u.y]
	}
	return nil
}

// neighbours returns a list of the floors or walls surrounding the current cell.
// Corners have at most two floors and edge pieces have at most three.
func (g *grid) neighbours(u *cell, isWall bool) []*cell {
	wp := []*cell{}
	neighbours := []*cell{g.north(u), g.south(u), g.west(u), g.east(u)}
	for _, neighbour := range neighbours {
		if neighbour != nil && neighbour.isWall == isWall {
			wp = append(wp, neighbour)
		}
	}
	return wp // walls or floors depending on isWall.
}

// dump prints a grid for debugging purposes.  This expects a fixed width
// font and looks better with some fixed fonts than others.
// The grid is dumped such that the 0,0 is at the bottom left and on
// a console 0,0 will be the first character of the last line dumped.
func (g *grid) dump() {
	width, height := g.Size()
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			if g.IsOpen(x, y) {
				fmt.Print("◽")
			} else {
				fmt.Print("◾")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// grid
// ===========================================================================
// cell

// cell is the building block for a grid.  Each cell knows its position in
// the grid and whether or not its a wall or a floor.
type cell struct {
	x, y   int
	isWall bool
}
