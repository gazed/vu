// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package vu

// PickDir calculates the direction of a pick ray using window coordinates.
// It goes from window coordinates to object coordinates which reverses
// the normal 3D to 2D transform of (From OpenGL FAQ 9.011):
//   1. Object Coordinates are transformed by the ModelView matrix to produce
//      Eye Coordinates.
//   2. Eye Coordinates are transformed by the Projection matrix to produce
//      Clip Coordinates.
//   3. Clip Coordinate X, Y, and Z are divided by Clip Coordinate W to produce
//      Normalized Device Coordinates.
//   4. Normalized Device Coordinates are scaled and translated by the viewport
//      parameters to produce Window Coordinates.
// Create a ray using the returned direction and the cameras current location
// as the ray origin. The ray that can then be used to test for object
// intersections.
//
// Also see
//     http://bookofhook.com/mousepick.pdf
//     http://schabby.de/picking-opengl-ray-tracing/
//     http://antongerdelan.net/opengl4/raycasting.html
//     (opengl FAQ Picking 20.0.010)
//     http://www.opengl.org/archives/resources/faq/technical/selection.htm

import (
	"vu/math/lin"
)

// PickDir calculates the direction of a pick ray using window coordinates.
// It goes from window coordinates to object coordinates which reverses
// the normal 3D to 2D transform.
//
// PickDir does not take scaling into account because it is not doing a full
// calculation of the inverse model matrix.
func PickDir(sx, sy int, fov, w, h, near, far float32, mv *lin.M4) *lin.V4 {
	cx := float32(2*sx)/w - 1
	cy := float32(2*sy)/h - 1
	clip := &lin.V4{cx, cy, -1, 1} // vector direction of ray in clip space

	// use the inverse perspective to go from clip to eye coordinates
	ip := lin.M4PerspectiveI(fov, w/h, near, far)
	eye := clip.MultR(ip)
	eye.Z = -1 // into the screen
	eye.W = 0  // want a vector, not a point

	// Inverse the view transform to go from eye coordinates to object coordinates.
	imv := mv.Clone().IModelView()
	world := eye.MultL(imv)
	return world.Unit()
}
