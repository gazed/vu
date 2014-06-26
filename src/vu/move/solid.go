// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package move

import (
	"vu/math/lin"
)

// Solid represents an object that can can be tested for ray intersection,
// but which is not intended to be controlled by the physics simulation.
type Solid interface {
	World() *lin.T            // Get/Set world transform which is the...
	SetWorld(world *lin.T)    // ...solids location and direction.
	Data() interface{}        // Get/Set application specific data...
	SetData(data interface{}) // ...often pointer to the scene graph node.
}

// ===========================================================================
// solid implementation.

// solid is the default implementation of the Solid interface.
type solid struct {
	data  interface{} // Unique data set by the calling application.
	shape Shape       // Body shape for collisions.
	world *lin.T      // World transform for the given shape.
	v0    *lin.V3     // Scratch vector.
}

func NewSolid(shape Shape) Solid { return newSolid(shape) }
func newSolid(shape Shape) *solid {
	s := &solid{}
	s.shape = shape
	s.world = lin.NewT().SetI() // world transform
	s.v0 = &lin.V3{}            // scratch vectors
	return s
}

func (s *solid) Data() interface{}        { return s.data }
func (s *solid) SetData(data interface{}) { s.data = data }
func (s *solid) SetWorld(world *lin.T)    { s.world.Set(world) }
func (s *solid) World() *lin.T            { return s.world }
