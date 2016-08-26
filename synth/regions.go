// Copyright © 2014-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package synth

import (
	"math"
	"math/rand"
)

// RegionData is used to associate a data value at each land location.
// For example this can be used to associate a land type or land region
// at each height location.
type RegionData [][]int

// Regions divides a given land. This is attempting to do
// something similar to creating a Voronoi diagram, just using a lot
// less code (also less efficient). The random seed is injected
// so that identical results can be re-created.
func Regions(size, numRegions int, seed int64) RegionData {
	data := make(RegionData, size)
	for x := range data {
		data[x] = make([]int, size)
	}

	// randomly scatter the different seed points for the regions.
	points := map[int]point{}
	rgen := rand.New(rand.NewSource(seed))
	for cnt := 0; cnt < numRegions; cnt++ {
		x, y := rgen.Int31n(int32(size)), rgen.Int31n(int32(size))
		data[x][y] = cnt + 1

		// save the region location.
		points[cnt+1] = point{int(x), int(y)}
	}

	// make unvisited points be part of the closest region type.
	for x := range data {
		for y := range data[x] {
			if data[x][y] == 0 {

				// find the closest region.
				regionType := 0
				shortestDistance := math.MaxInt64
				for rType, xy := range points {
					dsqr := (x-xy.x)*(x-xy.x) + (y-xy.y)*(y-xy.y)
					if dsqr < shortestDistance {
						shortestDistance = dsqr
						regionType = rType
					}
				}
				data[x][y] = regionType
			}
		}
	}
	return data
}

// point is a temporary structure used by Regions
// to save the region seed locations.
type point struct {
	x, y int
}
