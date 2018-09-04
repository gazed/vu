// Copyright © 2014-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// Package ai provides support for autonomous application unit behaviour.
// This is an experimental package that currently provides A-star and flow
// field path finding algorihms as well as a behaviour tree implementation.
//
// Package ai is provided as part of the vu (virtual universe) 3D engine.
package ai

// More information at:
// https://web.cs.ship.edu/~djmoon/gaming/gaming-notes/ai-movement.pdf
// http://www.raywenderlich.com/24824/introduction-to-ai-programming-for-games
// http://www.hobbygamedev.com/articles/vol8/real-time-videogame-ai/
// http://www.ramalila.net/Adventures/AI/RealTime.html

// Graph is a set of linked points that works with an A* pathing algorithm.
// A graph can be a grids or a set of waypoints or navigation meshes.
type Graph interface {

	// Neighbours returns the locations that can be reached
	// from the current location. The returned data may be
	// overwritten the next time this method is called.
	Neighbours(at Point) []Point

	// Cost for travelling from a to b. Used to value the move from
	// the current location to the next location. Lower cost is better.
	// Impassable areas have very high costs.
	Cost(a, b Point) float64

	// Estimate travel cost from a to b. Used to estimate cost
	// from current location to the goal. Lower cost is better.
	// Estimates that approximate the final path cost are best.
	Estimate(a, b Point) float64
}

// Point is a specific location in a Graph.
// Points with the same ID have the same location.
type Point interface {
	ID() int64 // Unique identifier for a location.
}

// =============================================================================
// Utility classes to attach a priority to a given location.

// priorityPoint wraps a location with a priority value used for sorting.
// May be moved to a interal package.
type priorityPoint struct {
	Priority float64 // Used to sort PriorityPoints.
	Point    Point   // Route location.
}

// PriorityPointHeap is a min-heap of PriorityPoint.
// Based on example from https://golang.org/pkg/container/heap/
type priorityPointHeap []priorityPoint

// Len implements sort.Interface.
func (h priorityPointHeap) Len() int { return len(h) }

// Less implements sort.Interface.
func (h priorityPointHeap) Less(i, j int) bool { return h[i].Priority < h[j].Priority }

// Swap implements sort.Interface.
func (h priorityPointHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

// Push implements heap.Interface.
func (h *priorityPointHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the
	// slice's length, not just its contents.
	*h = append(*h, x.(priorityPoint))
}

// Pop implements heap.Interface.
func (h *priorityPointHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
