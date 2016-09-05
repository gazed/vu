// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// Design Notes:
//  o Pov is a location/orientation.
//  o Pov is a transform hierarchy.
//  o Pov filters all user object creation through the engine entity manager.

import (
	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/physics"
	"github.com/gazed/vu/render"
)

// Pov is a Point-of-view which is a combination of position and orientation.
// Pov's are created by the application and may have additional associated
// components like rendered models and physics bodies. The associated
// components use the location and orientation of the Pov.
//
// A Pov can also have a child Pov's whose position and orientation are
// relative to the parent. This hiearchy of parent-child Pov's forms the
// transform scene graph. The hierarchy root is available through Eng.
type Pov struct {
	T *lin.T  // Transform combines location and orientation.
	S *lin.V3 // Per axis scale: >1 to enlarge, fraction<1 to shrink.

	// Cull set to true removes this Pov and its child Pov's
	// from processing. Cull is false by default.
	Cull bool // True to exclude from scene graph processing.

	// Used to register entities with the eng entity manager.
	eng *engine // Entity manager.
	eid uint64  // Unique entity identifier.

	// transform hierarchy. Each pov node can have children which base
	// their position and orientation relative to the parents.
	parent   *Pov   // nil if no parent: nil means root.
	children []*Pov // child transforms.

	// variables for recalculating transforms each update.
	toc float64 // distance to camera.
	mm  *lin.M4 // model transform.
	rot *lin.Q  // scratch rotation/orientation.
}

// newPov allocates and initialzes a point of view transform.
// Called by the engine.
func newPov(eng *engine, eid uint64) *Pov {
	p := &Pov{eng: eng, eid: eid}
	p.T = lin.NewT()
	p.S = &lin.V3{X: 1, Y: 1, Z: 1}

	// allocate scratch variables.
	p.rot = lin.NewQ()
	p.mm = &lin.M4{}
	return p
}

// NewPov creates, attaches, and returns a new child transform Pov.
func (p *Pov) NewPov() *Pov { return p.eng.newPov(p) }

// Dispose deletes and removes either the entire PovNode,
// or one of its attachments, PovModel, PovBody, ...
func (p *Pov) Dispose(kind int) { p.eng.dispose(p, kind) }

// At gets the Pov's location.
// Location is relative to parent. World space if no parent location.
func (p *Pov) At() (x, y, z float64) {
	return p.T.Loc.X, p.T.Loc.Y, p.T.Loc.Z
}

// SetAt sets the Pov at the given location.
// Location is relative to parent. World space if no parent location.
func (p *Pov) SetAt(x, y, z float64) *Pov {
	p.T.Loc.X, p.T.Loc.Y, p.T.Loc.Z = x, y, z
	return p
}

// Move directly affects the location by the given translation amounts
// along the given direction. Physics bodies should use Body.Push which
// affects velocity.
func (p *Pov) Move(x, y, z float64, dir *lin.Q) {
	dx, dy, dz := lin.MultSQ(x, y, z, dir)
	p.T.Loc.X += dx
	p.T.Loc.Y += dy
	p.T.Loc.Z += dz
}

// World returns the world space location.
// This method is valid in the Eng Update callback, not Create.
func (p *Pov) World() (x, y, z float64) {
	// The model matrix, mm, must have been set prior to calling
	v := &lin.V4{X: 0, Y: 0, Z: 0, W: 1}
	v.MultvM(v, p.mm)
	return v.X, v.Y, v.Z
}

// View returns the orientation of the Pov. Orientation combines
// direction and rotation about the direction.
// Orientation is relative to parent. World space if no parent orientation.
func (p *Pov) View() (q *lin.Q) { return p.T.Rot }

// SetView directly sets the Pov's orientation.
// Often used to align this Pov with the orientation of another.
// Orientation is relative to parent. World space if no parent orientation.
func (p *Pov) SetView(q *lin.Q) *Pov {
	r := p.T.Rot
	r.X, r.Y, r.Z, r.W = q.X, q.Y, q.Z, q.W
	return p
}

// Spin rotates x,y,z degrees about the X,Y,Z axis.
// The spins are combined in XYZ order, but generally this
// is used to spin about a single axis at a time.
func (p *Pov) Spin(x, y, z float64) {
	if x != 0 {
		p.rot.SetAa(1, 0, 0, lin.Rad(x))
		p.T.Rot.Mult(p.rot, p.T.Rot)
	}
	if y != 0 {
		p.rot.SetAa(0, 1, 0, lin.Rad(y))
		p.T.Rot.Mult(p.rot, p.T.Rot)
	}
	if z != 0 {
		p.rot.SetAa(0, 0, 1, lin.Rad(z))
		p.T.Rot.Mult(p.rot, p.T.Rot)
	}
}

// Scale retrieves the per-axis scale values at 3 separate XYZ values.
func (p *Pov) Scale() (x, y, z float64) { return p.S.X, p.S.Y, p.S.Z }

// SetScale assigns the XYZ per-axis scale values.
// Scale default is 1, greater than 1 enlarges, a positive fraction shrinks.
func (p *Pov) SetScale(x, y, z float64) *Pov {
	p.S.X, p.S.Y, p.S.Z = x, y, z
	return p
}

// Cam returns nil if there is no camera for this Pov.
func (p *Pov) Cam() *Camera { return p.eng.cam(p) }

// NewCam adds a camera to a Pov. This means that all rendered models
// in this Pov's hierarchy will be viewed with this camera settings.
func (p *Pov) NewCam() *Camera { return p.eng.newCam(p) }

// Body is an optional physics component associated with a Pov.
// Body returns nil if there is no physics body for this Pov.
func (p *Pov) Body() physics.Body { return p.eng.body(p) }

// NewBody creates non-colliding body associated with this Pov.
// Bodies are generally set on top level Pov transforms to ensure
// valid world coordindates.
func (p *Pov) NewBody(b physics.Body) physics.Body { return p.eng.newBody(p, b) }

// SetSolid makes the existing Pov physics Body collide. Nothing happens
// if there is no existing physics Body for this Pov.
func (p *Pov) SetSolid(mass, bounce float64) { p.eng.setSolid(p, mass, bounce) }

// Noise is an optional audio component. Played noises occur at the
// associated Pov's location. Noises that are played will be louder
// as the distance between the played noise and listener decreases.
// Place the single global noise listener at this Pov.
// Returns nil if no sound is associated with this Pov.
func (p *Pov) Noise() Noise { return p.eng.noise(p) }

// NewNoise creates an optional sound associated with this Pov.
func (p *Pov) NewNoise() Noise { return p.eng.newNoise(p) }

// SetListener sets the location of the listener to be the given Pov.
func (p *Pov) SetListener() { p.eng.setListener(p) }

// Light returns nil if no there light for this Pov.
func (p *Pov) Light() *Light { return p.eng.light(p) }

// NewLight creates and associates a light with this Pov.
// Light is optional. It affects lighting calculations for this Pov
// and all child pov's.
func (p *Pov) NewLight() *Light { return p.eng.newLight(p) }

// Layer returns nil if there is no layer for this Pov.
func (p *Pov) Layer() Layer { return p.eng.layer(p) }

// NewLayer creates a rendered texture at this Pov.
func (p *Pov) NewLayer() Layer { return p.eng.newLayer(p, render.ImageBuffer) }

// NewLabel creates a Model that displays a small text phrase.
//   shader      is expected to be a texture mapping shader like "uv".
//   font        is the font mapping file.
//   fontTexture is the bitmapped font texture
func (p *Pov) NewLabel(shader, font, fontTexture string) Labeler {
	return p.NewModel(shader, "fnt:"+font, "tex:"+fontTexture)
}

// NewAnimator creates an animated Model. Animated models have extra
// per-vertex data that repositions the verticies over time to make
// the model look likes its moving.
func (p *Pov) NewAnimator(shader, mod string) Animator {
	m := p.NewModel(shader).(*model)
	m.loadAnim(mod)
	return m
}

// Model returns nil if there is no model for this Pov.
func (p *Pov) Model() Model { return p.eng.model(p) }

// NewModel creates an optional rendered component associated with this Pov.
// Returns nil if a model already exists.
func (p *Pov) NewModel(shader string, attrs ...string) Model {
	return p.eng.newModel(p, shader, attrs...)
}

// remChild is used by a pov removing itself from the hierarchy.
func (p *Pov) remChild(c *Pov) {
	for index, c := range p.children {
		if c.eid == p.eid {
			p.children = append(p.children[:index], p.children[index+1:]...)
			return
		}
	}
}
