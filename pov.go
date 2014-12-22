// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// pov combines location and direction (orientation) to give a "point of view".

import (
	"github.com/gazed/vu/math/lin"
)

// pov is a location and an orientation (point of view). It is used for placing
// and rotating objects and cameras in 3D world space.
//    pov.Loc : location/position              - where we are.
//    pov.Rot : rotation/direction/orientation - which way we're facing.
type pov lin.T

func newPov() pov { return pov{&lin.V3{}, &lin.Q{0, 0, 0, 1}} }

// Set (=, copy, clone) assigns all the elements values from transform a to the
// corresponding element values in pov. The updated pov is returned.
func (p *pov) Set(a *lin.T) {
	p.Loc.Set(a.Loc)
	p.Rot.Set(a.Rot)
}

// Move increments the current position with respect to the current
// orientation, i.e. adds the distance travelled in the current direction
// to the current location.
func (p *pov) Move(x, y, z float64) {
	dx, dy, dz := lin.MultSQ(x, y, z, p.Rot)
	p.Loc.X += dx
	p.Loc.Y += dy
	p.Loc.Z += dz
}

// Spin rotates the current direction by the given number degrees around
// each axis.
func (p *pov) Spin(x, y, z float64) {
	if x != 0 {
		rotation := lin.NewQ().SetAa(1, 0, 0, lin.Rad(x))
		p.Rot.Mult(rotation, p.Rot)
	}
	if y != 0 {
		rotation := lin.NewQ().SetAa(0, 1, 0, lin.Rad(y))
		p.Rot.Mult(rotation, p.Rot)
	}
	if z != 0 {
		rotation := lin.NewQ().SetAa(0, 0, 1, lin.Rad(z))
		p.Rot.Mult(rotation, p.Rot)
	}
}
