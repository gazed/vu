// Copyright © 2014-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package grid

import (
	"math/rand"
)

// dungeon is a level comprised of square rooms connected by corridors.
type dungeon struct {
	grid // superclass grid
}

// Generate a dungeon by partioning the given space into randomly sized
// non-overlapping blocks. Reuse the definition of a room from room.go.
func (d *dungeon) Generate(width, depth int) Grid {
	d.create(width, depth, allWalls)
	rooms := d.rooms()
	d.corridors(rooms)

	// connect the rooms with corridors.
	return d
}

// rooms places random non-overlapping square rooms over the given grid.
// The newly created rooms are returned.
func (d *dungeon) rooms() []*room {
	sx, sy := d.Size()
	rooms := []*room{}
	possibleRooms := d.locateRooms(&room{0, 0, sx, sy})
	for _, rm := range possibleRooms {

		// Only use some of the possible spots for rooms.
		if rand.Intn(100) < 75 {
			rooms = append(rooms, rm)

			// randomize the dimensions of larger rooms.
			dx, dy := 1, 1
			if rm.w > 7 && rm.h > 7 {
				dx = rand.Intn(3) + 1
				dy = rand.Intn(3) + 1
			}
			for x := dx; x < rm.w-dx; x++ {
				for y := dy; y < rm.h-dy; y++ {
					d.cells[rm.x+x][rm.y+y].isWall = false
				}
			}
		}
	}
	return rooms
}

// locateRooms randomly and recursively quad-partitions a given room,
// gathering and returning all the generated sub-room dimensions.
func (d *dungeon) locateRooms(rm *room) []*room {
	min, max := 5, 20
	hx, hy := rm.w/2, rm.h/2
	if hx < min || hy < min {
		return []*room{rm} // to small to split.
	}

	// split randomly, or if to large.
	if rm.w > max || rm.h > max || rand.Intn(100) < 50 {
		rooms := []*room{}
		rooms = append(rooms, d.locateRooms(&room{rm.x, rm.y, hx, hy})...)
		rooms = append(rooms, d.locateRooms(&room{rm.x, rm.y + hy, hx, hy})...)
		rooms = append(rooms, d.locateRooms(&room{rm.x + hx, rm.y, hx, hy})...)
		rooms = append(rooms, d.locateRooms(&room{rm.x + hx, rm.y + hy, hx, hy})...)
		return rooms
	}
	return []*room{rm}
}

// corridors links each room to a neighbouring room.
func (d *dungeon) corridors(rooms []*room) {
	for cnt := 0; cnt < len(rooms); cnt++ {
		r0 := rooms[cnt]
		x0, y0 := r0.x+r0.w/2, r0.y+r0.h/2
		if cnt+1 < len(rooms) {
			r1 := rooms[cnt+1]
			x1, y1 := r1.x+r1.w/2, r1.y+r1.h/2
			dx, dy := 1, 1
			if x1-x0 < 0 {
				dx = -1
			}
			if y1-y0 < 0 {
				dy = -1
			}
			newx := x0
			for x := x0; x != x1; x += dx {
				d.cells[x][y0].isWall = false
				newx = x
			}
			for y := y0; y != y1; y += dy {
				d.cells[newx][y].isWall = false
			}
		}
	}
}
