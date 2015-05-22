// Copyright © 2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/move"
)

// Points of view, Pov's, are combinations of positions and orientations.
// Created by the application, they may have associated components, like
// rendered models and physics bodies. The associated objects use the
// location and orientation from the Pov. A Pov can also have a child
// Pov whose position and orientation is relative to the parent.
// A hiearchy of parent-child Pov's forms a transform hierarchy.
type Pov interface {
	Location() (x, y, z float64)     // Get, or
	SetLocation(x, y, z float64) Pov // ...Set the current location.
	Rotation() (q *lin.Q)            // Get, or
	SetRotation(q *lin.Q)            // ...Set current quaternion rotation.
	Spin(x, y, z float64)            // Rotate degrees about the given axis.
	Move(x, y, z float64, q *lin.Q)  // Move along indicated direction.

	// Visible affects this transform and any child transforms.
	Visible() bool           // Invisible transforms are removed from
	SetVisible(visible bool) // ...rendering without disposing them.

	// Per axis scale. Normal is 1, greater than 1 to enlarge,
	// positive fraction to shrink.
	Scale() (x, y, z float64)     // Get, or
	SetScale(x, y, z float64) Pov // ...Set the current scale.

	// Create a child POV from this pov.
	NewPov() Pov      // Attaches a new child transform Pov to p.
	Dispose(kind int) // Discard POV, MODEL, BODY, VIEW, or NOISE.

	// Adding a camera to a Pov means that all rendered models in the povs
	// hierarchy will be viewed with the views camera settings.
	View() View    // Nil if no view for this Pov.
	NewView() View // View for the group of this Pov's child models.

	// Model is an optional rendered component associated with a Pov.
	Model() Model                 // Nil if no model.
	NewModel(shader string) Model // Nil if a model already exists.

	// Body is an optional physics component associated with a Pov. Bodies
	// are set at top level Pov transforms to have valid world coordindates.
	Body() move.Body               // Nil if no body.
	NewBody(b move.Body) move.Body // Create non-colliding body.
	SetSolid(mass, bounce float64) // Make existing body collide.

	// Sound is an optional audio component. Played sounds occur at the
	// associated Pov's location. Sounds that are played will be louder
	// as the distance between the played sound and listener decreases.
	Noise() Noise    // Nil if no sound.
	NewNoise() Noise // Create noise. Nil if noise already exists.
	SetListener()    // Place the single global sound listener at this pov.

	// Light is an optional component associated with a Pov.
	Light() Light    // Nil if no light for this Pov.
	NewLight() Light // Create a light at this Pov.
}

// Pov
// =============================================================================
// pov

// pov implements Pov.
type pov struct {
	eng     *engine // Entity manager.
	eid     uint64  // Unique entity identifier.
	at      *lin.T  // point of view: location/orientation.
	scale   *lin.V3 // Per axis scale: >1 to enlarge, 0<1 to shrink.
	visible bool    // True means visible for rendering.

	// Each pov node can have children which base their position and
	// orientation relative to the parents.
	parent   *pov   // nil if no parent.
	children []*pov // child transforms.

	// variables for recalculating transforms each update.
	toc float64 // distance to camera.
	rot *lin.Q  // rotation/orientation.
	mm  *lin.M4 // model transform.
}

// newPov allocates and initialzes a point of view transform.
func newPov(eng *engine, eid uint64) *pov {
	p := &pov{eng: eng, eid: eid, visible: true}
	p.at = lin.NewT()
	p.scale = &lin.V3{X: 1, Y: 1, Z: 1}

	// allocate scratch variables.
	p.rot = lin.NewQ()
	p.mm = &lin.M4{}
	return p
}

// Implement Pov.
func (p *pov) Location() (x, y, z float64) {
	return p.at.Loc.X, p.at.Loc.Y, p.at.Loc.Z
}

// Implement Pov.
func (p *pov) SetLocation(x, y, z float64) Pov {
	p.at.Loc.X, p.at.Loc.Y, p.at.Loc.Z = x, y, z
	return p
}

// Implement Pov.
func (p *pov) Rotation() (q *lin.Q) {
	return p.at.Rot
}

// Implement Pov.
func (p *pov) SetRotation(q *lin.Q) {
	r := p.at.Rot
	r.X, r.Y, r.Z, r.W = q.X, q.Y, q.Z, q.W
}

// Implement Pov.
func (p *pov) Spin(x, y, z float64) {
	if x != 0 {
		p.rot.SetAa(1, 0, 0, lin.Rad(x))
		p.at.Rot.Mult(p.rot, p.at.Rot)
	}
	if y != 0 {
		p.rot.SetAa(0, 1, 0, lin.Rad(y))
		p.at.Rot.Mult(p.rot, p.at.Rot)
	}
	if z != 0 {
		p.rot.SetAa(0, 0, 1, lin.Rad(z))
		p.at.Rot.Mult(p.rot, p.at.Rot)
	}
}

// Move directly affects the location by the given translation amounts
// along the given direction. Physics bodies should use Body.Push which
// affects velocity.
func (p *pov) Move(x, y, z float64, dir *lin.Q) {
	dx, dy, dz := lin.MultSQ(x, y, z, dir)
	p.at.Loc.X += dx
	p.at.Loc.Y += dy
	p.at.Loc.Z += dz
}

// Implement Pov.
func (p *pov) Visible() bool { return p.visible }
func (p *pov) SetVisible(visible bool) {
	p.visible = visible
}

// Implement Pov.
func (p *pov) Scale() (x, y, z float64) { return p.scale.X, p.scale.Y, p.scale.Z }
func (p *pov) SetScale(x, y, z float64) Pov {
	p.scale.X, p.scale.Y, p.scale.Z = x, y, z
	return p
}

// remChild is used by a pov removing itself from the heirarchy.
func (p *pov) remChild(c *pov) {
	for index, c := range p.children {
		if c.eid == p.eid {
			p.children = append(p.children[:index], p.children[index+1:]...)
			return
		}
	}
}

// Implement Pov interface. These convenience methods wrap the entity
// manager methods so that the entity manager doesn't have to be
// referenced anywhere else.
func (p *pov) NewPov() Pov                   { return p.eng.newPov(p) }
func (p *pov) Dispose(kind int)              { p.eng.dispose(p, kind) }
func (p *pov) View() View                    { return p.eng.view(p) }
func (p *pov) NewView() View                 { return p.eng.newView(p) }
func (p *pov) Model() Model                  { return p.eng.model(p) }
func (p *pov) NewModel(shader string) Model  { return p.eng.newModel(p, shader) }
func (p *pov) Light() Light                  { return p.eng.light(p) }
func (p *pov) NewLight() Light               { return p.eng.newLight(p) }
func (p *pov) Body() move.Body               { return p.eng.body(p) }
func (p *pov) NewBody(b move.Body) move.Body { return p.eng.newBody(p, b) }
func (p *pov) SetSolid(mass, bounce float64) { p.eng.setSolid(p, mass, bounce) }
func (p *pov) Noise() Noise                  { return p.eng.noise(p) }
func (p *pov) NewNoise() Noise               { return p.eng.newNoise(p) }
func (p *pov) SetListener()                  { p.eng.setListener(p) }
