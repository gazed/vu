// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package grid

import (
	"math/rand"
	"time"
)

// rooms is a skirmish grid made up of connected empty spaces.
type rooms struct {
	grid         // superclass grid
	min, max int // room sizes including walls.
}

// room tracks the spaces carved out of the grid. Rooms walls can overlap.
// Note that the code considers 0,0 to be bottom left with width increasing
// x (left to right), and height increasing y (bottom to top).
type room struct {
	x, y int // location of the bottom left corner
	w, h int // size of the room
}

// Specifies if a split happens left/right or top/bottom.
const (
	leftRight = true  // split the room left/right
	topBottom = false // split the room top/bottom.
)

// Generate a skirmish grid that has lots of interconnected rooms.
// Each room has 3 exits.
//
// TODO it is possible to get a bisected maze.  It may be necessary
// to check for a wall across the entire maze and punch a hole in it.
func (rms *rooms) Generate(width, depth int) Grid {
	random := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	rms.create(width, depth, allWalls)
	rms.min = 4 // makes 2x2 rooms (excludes walls), don't want smaller.
	rms.max = 7 // makes 5x5 rooms (excludes walls), don't want larger.
	width, depth = rms.Size()
	initialRoom := &room{0, 0, width, depth}
	dividedRooms := rms.getRooms(random, initialRoom)

	// clear the interior of each room.
	for _, rm := range dividedRooms {
		for x := rm.x + 1; x < rm.x+rm.w-1; x++ {
			for y := rm.y + 1; y < rm.y+rm.h-1; y++ {
				rms.cells[x][y].isWall = allPassages
			}
		}
		rms.ensureExits(random, rm)
	}
	return rms
}

// getRooms recursively subdivides the given rooms into smaller rooms.
// This method randomly chooses to split rooms depending on their size.
// Very large rooms always get split, very small ones don't, and in between
// may or may not be split.
func (rms *rooms) getRooms(random *rand.Rand, rm *room) (dividedRooms []*room) {

	// randomly decide if the room must be split
	maxLR := (rm.w+1)/2 >= rms.min+random.Intn(rms.max/2)
	maxTB := (rm.h+1)/2 >= rms.min+random.Intn(rms.max/2)
	if rm1, rm2 := rms.split(random, rm, maxLR, maxTB); rm1 != nil && rm2 != nil {
		dividedRooms = append(dividedRooms, rms.getRooms(random, rm1)...)
		dividedRooms = append(dividedRooms, rms.getRooms(random, rm2)...)
	} else {
		dividedRooms = append(dividedRooms, rm)
	}
	return
}

// split divides the given room based on the how-to-split booleans.
// The return is either two valid room pointers, or two nil pointers.
func (rms *rooms) split(random *rand.Rand, rm *room, lr, tb bool) (rm1, rm2 *room) {
	if lr || tb {
		splitType := topBottom
		if lr && tb {
			splitType = random.Intn(2) == 0
		} else if lr {
			splitType = leftRight
		}
		splitSpot := rms.splitSpots(rm, splitType)
		splitSpot = random.Intn(splitSpot+1) + rms.min - 1
		r1, r2 := rms.splitRoom(rm, splitSpot, splitType)
		rm1, rm2 = &r1, &r2
	}
	return
}

// splitSpots returns how many choices there are for splitting the indicated
// side in the room.  This is based on the minimum size for the level.
// A return of 0 indicates 1 choice.
func (rms *rooms) splitSpots(rm *room, lr bool) int {
	sideSize := rm.h
	if lr {
		sideSize = rm.w
	}
	return sideSize - rms.min - (rms.min - 1)
}

// splitRoom chops a room into two rooms at the given cut point.
func (rms *rooms) splitRoom(rm *room, cut int, lr bool) (rm1, rm2 room) {
	if lr {
		return room{rm.x, rm.y, cut + 1, rm.h}, room{rm.x + cut, rm.y, rm.w - cut, rm.h}
	}
	return room{rm.x, rm.y, rm.w, cut + 1}, room{rm.x, rm.y + cut, rm.w, rm.h - cut}
}

// ensureExits makes sure there is an outside exit on two sides
// of the room. The corners are left alone.
func (rms *rooms) ensureExits(random *rand.Rand, rm *room) {
	var top, bot, left, right []*cell
	xmax, ymax := rm.w, rm.h
	for x := rm.x + 1; x < rm.x+rm.w-1; x++ {
		u := rms.cells[x][rm.y]
		if u.isWall {
			top = append(top, u)
		}
		u = rms.cells[x][rm.y+rm.h-1]
		if u.isWall {
			bot = append(bot, u)
		}
	}
	for y := rm.y + 1; y < rm.y+rm.h-1; y++ {
		u := rms.cells[rm.x][y]
		if u.isWall {
			left = append(left, u)
		}
		u = rms.cells[rm.x+rm.w-1][y]
		if u.isWall {
			right = append(right, u)
		}
	}

	// randomize which sides get exits.
	if random.Intn(2) == 0 {
		rms.ensureExit(random, top, xmax-2)
		rms.ensureExit(random, left, ymax-2)
	} else {
		rms.ensureExit(random, bot, xmax-2)
		rms.ensureExit(random, right, ymax-2)
	}
}

// ensureExit chops a hole in the given side.  Sometimes the hole chopped
// can be a dead-end.  Chopping an additional hole in the holes neighbouring
// walls guarantees an exit.
func (rms *rooms) ensureExit(random *rand.Rand, side []*cell, max int) {
	if len(side) == max {
		index := random.Intn(len(side))
		u := side[index]
		u.isWall = allPassages
		rms.cells[u.x][u.y].isWall = u.isWall

		// ensure the chop gets into the maze by chopping again if necessary.
		walls := rms.neighbours(u, allWalls)
		if len(walls) == 3 {
			wallIndex := random.Intn(len(walls))
			u := walls[wallIndex]
			u.isWall = allPassages
			rms.cells[u.x][u.y].isWall = u.isWall
		}
	}
}
