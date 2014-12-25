// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"github.com/gazed/vu/audio"
	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/move"
)

// A part is a node in a scene graph. It is responsibile for positioning
// itself and its children in 3D space thus forming a transform hierarchy.
// Note that all child parts will be affected by parent transform changes.
//
// A part can optionally be rendered and/or participate in physics, thus
// providing a link between the two systems.
//    rendered : SetRole with a shader asset makes the part participate in
//               rendering. Afterwards add the data, textures, and other
//               assets needed by that shader.
//    physics  : SetBody makes the part participate in physics. The Parts
//               location is controlled by the physics simulation.
type Part interface {
	AddPart() Part  // Create a new subordinate part attached to this one.
	RemPart(p Part) // Remove, dispose, the given subordinate part.
	Dispose()       // Remove this part and all its subordinate parts.

	// Methods for updating data that affects culling.
	SetCullable(cullable bool) // False excludes a part from being culled.
	Visible() bool             // Invisible parts are removed from
	SetVisible(visible bool)   // ...rendering without disposing them.

	// A Part inherently has a point of view. Note that moving a part with
	// a physics body affects the part's velocity instead of directly
	// updating the location or orientation.
	Location() (x, y, z float64)    // Get, or
	SetLocation(x, y, z float64)    // ...Set the current location.
	Rotation() (x, y, z, w float64) // Get, or
	SetRotation(x, y, z, w float64) // ...Set current quaternion rotation.
	Spin(x, y, z float64)           // Rotate degrees about the given axis.
	Move(x, y, z float64)           // Move along indicated direction.
	Scale() (x, y, z float64)       // Get, or
	SetScale(x, y, z float64) Part  // ...Set the current scale.

	// Role is an optional rendered aspect of the Part. Roles link the
	// part to a shader and shader data.
	Role() Role                 // Return role, or nil if no current role.
	SetRole(shader string) Role // Create role if no current role.
	RemRole()                   // Delete the current role.

	// Body is an optional physics aspect of the Part. Bodies are only to be
	// associated with root level Parts to ensure valid world coordindates.
	Body() move.Body                           // Return nil if no body.
	SetBody(b move.Body, mass, bounce float64) // Body with mass and bounce.
	RemBody()                                  // Delete the physics body.

	// Form is an optional physics shape and location used for raycasting.
	// It does not participate in the physics simulation.
	Form() move.Body     // Return nil if no form.
	SetForm(g move.Body) // Plane or sphere.

	// Sound is an optional audio component. Played sounds occur at the
	// parts current location.
	Sound(name string) audio.SoundMaker    // Return nil if no such sound.
	AddSound(name string) audio.SoundMaker // Associate sound with this part.
}

// Part interface
// =============================================================================
// part - Part implementation

// part implements the Part interface.
type part struct {
	pov         // point of view: location/orientation.
	cull        // Any part can potentially be culled.
	pid  uint32 // Unique part instance identifier.
	sm   *stage // stage manager is injected on creation.

	// Scene graph hierarchy information.
	parts  []*part // Each part node can be made of one or more parts.
	parent uint32  // 0 if no parent.

	// Optional render and transform information.
	role  *role   // Render data linking render and asset subsystems.
	scale *lin.V3 // Scale, per axis: >1 to enlarge, 0<1 to shrink.

	// Optional physics fields are nil if not used.
	body move.Body // Motion body used by physics subsystem.
	form move.Body // Shape and location/orientation used for ray casting.

	// Optional audio field is nil if not used.
	sounds map[string]audio.SoundMaker

	// scratch variables are reused each render cycle.
	rotation *lin.Q  // Scratch rotation/orientation.
	mm       *lin.M4 // Scratch model transform.
	pt       *lin.M4 // Scratch parent model transform.
}

// newPart creates and initialzes a part instance.
func newPart(pid uint32, sm *stage) *part {
	p := &part{}
	p.pid = pid
	p.sm = sm
	p.scale = &lin.V3{1, 1, 1}
	p.pov = newPov()
	p.parts = []*part{}
	p.role = nil
	p.cullable = true
	p.visible = true

	// scratch variables.
	p.rotation = lin.NewQ()
	p.mm = &lin.M4{}
	p.pt = lin.NewM4I()
	return p
}

// Part interface implementation.
func (p *part) Dispose() {
	for _, child := range p.parts {
		child.Dispose()
	}
	p.sm.remPart(p)
	p.parts = nil
	p.RemRole()
	p.RemBody()
}

// Part interface implementation.
func (p *part) AddPart() Part {
	np := p.sm.addPart()
	np.parent = p.pid
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

func (p *part) Location() (x, y, z float64) {
	return p.pov.Loc.X, p.pov.Loc.Y, p.pov.Loc.Z
}
func (p *part) SetLocation(x, y, z float64) {
	p.pov.Loc.X, p.pov.Loc.Y, p.pov.Loc.Z = x, y, z
	if p.sounds != nil {
		for _, sound := range p.sounds {
			sound.SetLocation(x, y, z)
		}
	}
}
func (p *part) Rotation() (x, y, z, w float64) {
	return p.pov.Rot.X, p.pov.Rot.Y, p.pov.Rot.Z, p.pov.Rot.W
}
func (p *part) SetRotation(x, y, z, w float64) {
	p.pov.Rot.X, p.pov.Rot.Y, p.pov.Rot.Z, p.pov.Rot.W = x, y, z, w
}
func (p *part) Spin(x, y, z float64)     { p.pov.Spin(x, y, z) }
func (p *part) Scale() (x, y, z float64) { return p.scale.X, p.scale.Y, p.scale.Z }
func (p *part) SetScale(x, y, z float64) Part {
	p.scale.X, p.scale.Y, p.scale.Z = x, y, z
	if p.role != nil {
		p.role.model.SetScale(x, y, z)
	}
	return p
}

// Move overrides the default movement behaviour so that motion is applied to
// any associated physics bodies as velocity instead of updating the location.
func (p *part) Move(x, y, z float64) {
	if p.body == nil {
		p.pov.Move(x, y, z)
	} else {

		// apply push in the current direction.
		dx, dy, dz := lin.MultSQ(x, y, z, p.Rot)
		p.body.Push(dx, dy, dz)
	}
}

// Model returns the current rendered model.
func (p *part) Role() Role {
	if p.role != nil {
		return p.role
	}
	return nil
}
func (p *part) SetRole(shader string) Role {
	if p.role == nil {
		p.role = newRole(shader, p.sm.assets)
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

// Part interface implementation.
func (p *part) Body() move.Body { return p.body }
func (p *part) SetBody(body move.Body, mass, bounce float64) {
	if p.body != nil {
		p.RemBody()
	}
	p.body = body.SetMaterial(mass, bounce)
	p.body.SetWorld((*lin.T)(&p.pov))
	p.sm.addBody(p.body)
}

// Part interface implementation.
func (p *part) RemBody() {
	if p.body != nil {
		p.sm.remBody(p.body)
		p.body = nil
	}
}

func (p *part) SetForm(form move.Body) {
	p.form = form
	p.form.SetWorld((*lin.T)(&p.pov))
}
func (p *part) Form() move.Body { return p.form }

// AddSound creates a SoundMaker that is linked to this part.
// Nil is returned if there are errors creating the sound.
func (p *part) AddSound(name string) audio.SoundMaker {
	if p.sounds == nil {
		p.sounds = map[string]audio.SoundMaker{}
	}
	s := p.sm.assets.getSound(name)
	sm := p.sm.assets.ac.NewSoundMaker(s)
	p.sounds[name] = sm
	return sm
}
func (p *part) Sound(name string) audio.SoundMaker {
	if p.sounds != nil {
		return p.sounds[name] // returns nil if not found.
	}
	return nil // no sounds have been associated with this part.
}

// model transform must be done in scale, rotate, translate order.
func (p *part) modelTransform(m *lin.M4) *lin.M4 {
	p.mm.SetQ(p.rotation.Inv(p.Rot))                   // rotation.
	p.mm.ScaleSM(p.Scale())                            // scale is applied first (on left of rotation)
	return p.mm.TranslateMT(p.Loc.X, p.Loc.Y, p.Loc.Z) // translate is applied last (on right of rotation).

}

// stage is recursively called down the part hierarchy.
func (p *part) stage(sc *scene, pt *lin.M4, dt float64) {
	p.mm = p.modelTransform(p.mm) // updates p.mm (model transform matrix)
	p.mm.Mult(p.mm, pt)           // model transform + parent transform
	if p.role != nil {
		if sc.is2D {
			p.role.Set2D()
		}
		p.role.update(p.mm, sc.cam.vm, sc.cam.pm, dt)
		p.sm.stage(p.role.model) // add to the list of rendered parts.
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
