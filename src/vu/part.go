// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package vu

import (
	"vu/math/lin"
	"vu/move"
)

// A part is a node in a scene graph. It is responsibile for positioning
// itself and its children in 3D space thus forming a transform hierarchy.
// Note that all child parts will be affected by parent transform changes.
//
// A part can optionally be rendered and/or participate in physics.
//    rendered : SetRole with a shader asset makes the part participate in
//               rendering. Afterwards add the data, textures, and other
//               assets needed by that shader.
//    physics  : SetBody makes the part participate in physics. The Parts
//               location is now controlled by the physics simulation.
type Part interface {
	AddPart() Part  // Create a new subordinate part attached to this one.
	RemPart(p Part) // Remove, dispose, a subordinate part.
	Dispose()       // Remove this part, all sub-parts and any associated info.

	// Moving or Spinning a part that is in physics affects the part's velocity
	// instead of directly updating the location or orientation.
	Location() (x, y, z float64)         // Get, or
	SetLocation(x, y, z float64) Part    // ...Set the current location.
	Move(x, y, z float64)                // Move along indicated direction.
	Rotation() (x, y, z, w float64)      // Get, or
	SetRotation(x, y, z, w float64) Part // ...Set current quaternion rotation.
	Spin(x, y, z float64)                // Rotate degrees about the given axis.
	Scale() (x, y, z float64)            // Get, or
	SetScale(x, y, z float64) Part       // ...Set the current scale.

	// Role is the optional rendered aspect of the Part.
	SetRole(shader string) Role // Creates role if no current role.
	RemRole()                   // Delete the current role.
	Role() Role                 // Return role, or nil if no current role.
	SetCullable(cullable bool)  // False excludes a part from being culled.
	Visible() bool              // Invisible parts are removed from
	SetVisible(visible bool)    // ...rendering without disposing them.

	// Body is the optional physics aspect of the Part.
	SetBody(bod move.Body, m, b float64) // Body of mass m, bounce b.
	RemBody()                            // Delete the physics body.
	Speed() (x, y, z float64)            // Current linear velocity.
	Push(x, y, z float64)                // Change the linear velocity.
	Turn(x, y, z float64)                // Change the angular velocity.
	Stop()                               // Remove all velocity.
	Solid() move.Solid                   // Solid is a lightweight shape
	SetSolid(sol move.Solid)             // ...intended for ray casting.
}

// FUTURE Find a nice, and safe, way to consolidate the location/orientation
//        information. Currently it is needed/duplicated by pov/part, and the
//        physics objects, Body, Solid.

// Part interface
// ===========================================================================
// part - Part implementation

// part implements the Part interface.
type part struct {
	pov                 // Embed the pov location and orientation struct.
	scale    *lin.V3    // Scale, per axis: >1 to enlarge, 0<1 to shrink.
	staged   bool       // True if the the part is visible.
	parts    []*part    // Each part node can be made of one or more parts.
	toc      float64    // Distance to center (to->c) for sorting and culling.
	role     *role      // Render data linking render and asset subsystems.
	body     move.Body  // Motion body used by physics subsystem.
	solid    move.Solid // Solid shape used for ray casting.
	assets   *assets    // Asset manager.
	tracker  feedback   // Feed tracking information out of the hierarchy.
	cullable bool       // Can/can't be culled is under control of the application.
	culled   bool       // Draw or don't under control of engine.
	visible  bool       // Draw or don't under control of application.

	// scratch variables are used each render cycle. They are optimizations that
	// prevent having to create temporary structures each render cycle.
	rotation *lin.Q  // Scratch rotation/orientation.
	mm       *lin.M4 // Scratch model transform.
	pt       *lin.M4 // Scratch parent model transform.
}

// newPart creates and initialzes a part instance.
func newPart(f feedback, a *assets) *part {
	p := &part{}
	p.scale = &lin.V3{1, 1, 1}
	p.tracker = f
	p.loc = &lin.V3{}
	p.dir = &lin.Q{0, 0, 0, 1}
	p.parts = []*part{}
	p.role = nil
	p.assets = a
	p.cullable = true
	p.culled = false
	p.visible = true

	// scratch variables.
	p.rotation = lin.NewQ()
	p.mm = &lin.M4{}
	p.pt = &lin.M4{}
	return p
}

// Part interface implementation.
func (p *part) Dispose() {
	for _, child := range p.parts {
		child.Dispose()
	}
	p.parts = nil
	p.RemRole()
	p.RemBody()
}
func (p *part) SetCullable(cullable bool) { p.cullable = cullable }
func (p *part) Visible() bool             { return p.visible }
func (p *part) SetVisible(visible bool)   { p.visible = visible }

// Part interface implementation.
func (p *part) AddPart() Part {
	np := newPart(p.tracker, p.assets)
	p.parts = append(p.parts, np)
	return np
}

// Part interface implementation.
// Find and remove the part (will point to the same record).
func (p *part) RemPart(child Part) {
	if pt, _ := child.(*part); pt != nil {
		for index, partPtr := range p.parts {
			if partPtr == pt {
				pt.Dispose()
				p.parts = append(p.parts[:index], p.parts[index+1:]...)
				return
			}
		}
	}
}

// SetLocation directly updates the parts location to the given coordinates.
// This is a form of teleportation when the part has an active physics body.
func (p *part) SetLocation(x, y, z float64) Part {
	p.pov.SetLocation(x, y, z)
	if p.body != nil {
		p.body.World().Loc.SetS(x, y, z)
	}
	if p.solid != nil {
		p.solid.World().Loc.SetS(x, y, z)
	}
	return p
}

// SetRotation directly updates the parts rotation to the given direction.
// This is a form of teleportation when the part has an active physics body.
func (p *part) SetRotation(x, y, z, w float64) Part {
	p.pov.SetRotation(x, y, z, w)
	if p.body != nil {
		p.body.World().Rot.SetS(x, y, z, w)
	}
	if p.solid != nil {
		p.solid.World().Rot.SetS(x, y, z, w)
	}
	return p
}

// Spin applies the spin to the parts orientation. It also updates any
// associated Body or Solid
func (p *part) Spin(x, y, z float64) {
	p.pov.Spin(x, y, z)
	if p.body != nil {
		p.body.World().Rot.Set(p.pov.dir)
	}
	if p.solid != nil {
		p.solid.World().Rot.Set(p.pov.dir)
	}
}

func (p *part) Scale() (x, y, z float64) { return p.scale.X, p.scale.Y, p.scale.Z }
func (p *part) SetScale(x, y, z float64) Part {
	p.scale.X, p.scale.Y, p.scale.Z = x, y, z
	if p.role != nil {
		p.role.model.SetScale(x, y, z)
	}
	return p
}

// Move overrides the default movement behaviour so that motion is applied to
// any associated bodies as velocity instead of directly updating the location.
func (p *part) Move(x, y, z float64) {
	if p.body == nil {
		p.pov.Move(x, y, z)
	} else {

		// apply push in the current direction.
		dx, dy, dz := lin.MultSQ(x, y, z, p.dir)
		p.body.Push(dx, dy, dz)
	}
	if p.solid != nil {
		p.solid.World().Loc.Set(p.pov.loc)
	}
}

// Model returns the current rendered model.
func (p *part) Role() Role { return p.role }
func (p *part) SetRole(shader string) Role {
	if p.role == nil {
		p.role = newRole(shader, p.assets)
		p.role.model.SetScale(p.scale.X, p.scale.Y, p.scale.Z)
	}
	return p.role
}
func (p *part) RemRole() {
	if p.role != nil {
		p.role.dispose()
		p.role = nil
	}
}

// Stop removes all linear and angular velocity from a physics body.
// Nothing happens if there is no associated physics body.
func (p *part) Stop() {
	if p.body != nil {
		p.body.Stop()
		p.body.Rest()
	}
}

// Push adds to the bodies current linear velocity.
// Nothing happens if there is no associated physics body.
func (p *part) Push(x, y, z float64) {
	if p.body != nil {
		p.body.Push(x, y, z)
	}
}

// Turn adds to the bodies current linear velocity.
// Nothing happens if there is no associated physics body.
func (p *part) Turn(x, y, z float64) {
	if p.body != nil {
		p.body.Turn(x, y, z)
	}
}

// Speed returns the bodies current linear velocity. Return 0,0,0 if
// there is no associated physics body.
func (p *part) Speed() (x, y, z float64) {
	if p.body != nil {
		return p.body.Speed()
	}
	return 0, 0, 0
}

// Part interface implementation.
func (p *part) SetBody(body move.Body, mass, bounce float64) {
	if p.body != nil {
		p.RemBody()
	}
	p.body = body.SetMaterial(mass, bounce)
	p.body.SetData(p)
	p.body.World().Loc.Set(p.loc)
	p.body.World().Rot.Set(p.dir)
	p.tracker.track(p.body)
}
func (p *part) SetSolid(sol move.Solid) {
	p.solid = sol
	p.solid.World().Loc.Set(p.loc)
	p.solid.World().Rot.Set(p.dir)
}
func (p *part) Solid() move.Solid { return p.solid }

// Part interface implementation.
func (p *part) RemBody() {
	if p.body != nil {
		p.tracker.release(p.body)
		p.body = nil
	}
}

// distanceTo returns the distance squared of the part to the given center.
func (p *part) distanceTo(cenx, ceny, cenz float64) float64 {
	dx := p.loc.X - cenx
	dy := p.loc.Y - ceny
	dz := p.loc.Z - cenz
	return float64(dx*dx + dy*dy + dz*dz)
}

// model transform must be done in scale, rotate, translate order.
func (p *part) modelTransform(m *lin.M4) *lin.M4 {
	p.mm.SetQ(p.rotation.Inv(p.dir))                   // rotation.
	p.mm.ScaleSM(p.Scale())                            // scale is applied first (on left of rotation)
	return p.mm.TranslateMT(p.loc.X, p.loc.Y, p.loc.Z) // translate is applied last (on right of rotation).

}

// stage the part for rendering. Update the part and add it to the list of
// rendered parts.
func (p *part) stage(sc *scene, dt float64) {
	p.pt.Set(lin.M4I)
	p.stagePart(sc, p.pt, dt)
}
func (p *part) stagePart(sc *scene, pt *lin.M4, dt float64) {
	if p.Visible() {
		p.mm = p.modelTransform(p.mm) // updates p.mm (model transform matrix)
		p.mm.Mult(p.mm, pt)           // model transform + parent transform
		if p.role != nil {
			if sc.is2D {
				p.role.Set2D()
			}
			p.role.update(sc.l, p.mm, sc.vm, sc.pm, dt)
			if p.role.effect != nil {
				p.role.effect.Update(p.role.Mesh(), dt)
			}
			p.tracker.stage(p.role.model) // add to the list of rendered parts.
		}

		// render all the parts children
		for _, child := range p.parts {
			if !child.culled {
				p.pt.Set(p.mm) // ensures the original model transform does not change.
				child.stagePart(sc, p.pt, dt)
			}
		}
	}
}

// verify passes the request down to the render.Model.
func (p *part) verify() (err error) {
	if p.role != nil {
		if err = p.role.model.Verify(); err != nil {
			return err
		}
	}
	for _, p := range p.parts {
		if err = p.verify(); err != nil {
			return err
		}
	}
	return
}

// part
// ===========================================================================
// Parts

// Parts is used to sort a slice of parts in order to get transparency working.
// Objects furthest away have to be drawn first.
// This is only public for the sort package and is not for application use.
type Parts []*part

// Sort parts ordered by distance.
func (p Parts) Len() int           { return len(p) }
func (p Parts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Parts) Less(i, j int) bool { return p[i].toc > p[j].toc }
