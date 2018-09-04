// Copyright © 2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package ai

import (
	"container/heap" // for priority queue.
)

// Find is an A* path finder. Implementation based on:
//    https://en.wikipedia.org/wiki/A*_search_algorithm
//    https://www.redblobgames.com/pathfinding/a-star/implementation.html
// The input path slice is reset to zero length and filled with the path
// results. An empty path is returned if no route was found.
//
// Computation time increases proportional to the the size and number
// of locations in the graph.
func Find(graph Graph, start, goal Point, path *[]Point) {
	cameFrom := map[int64]Point{start.ID(): start}
	costSoFar := map[int64]float64{start.ID(): 0}
	frontier := &priorityPointHeap{priorityPoint{Point: start, Priority: 0}}
	heap.Init(frontier)
	for frontier.Len() > 0 {
		current := heap.Pop(frontier).(priorityPoint).Point
		if current.ID() == goal.ID() {
			break // success
		}
		for _, next := range graph.Neighbours(current) {
			newCost := costSoFar[current.ID()] + graph.Cost(current, next)
			if csf, ok := costSoFar[next.ID()]; !ok || newCost < csf {
				costSoFar[next.ID()] = newCost
				priority := newCost + graph.Estimate(next, goal)
				heap.Push(frontier, priorityPoint{Point: next, Priority: priority})
				cameFrom[next.ID()] = current
			}
		}
	}

	// unwind cameFrom results to get path from start to goal
	*path = (*path)[:0] // reset to reuse existing memory.
	if last, ok := cameFrom[goal.ID()]; ok {
		*path = append(*path, last, goal)
		for {
			prev := cameFrom[last.ID()]
			*path = append(*path, nil)     // create space at end of slice...
			copy((*path)[1:], (*path)[0:]) // ...move old values up one spot...
			(*path)[0] = prev              // ...and insert new value at beginning.
			if prev.ID() == start.ID() {
				break // reached start.
			}
			last = prev
		}
	}
}
