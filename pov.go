// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// pov.go key application API for creating objects and managing world space.
// DESIGN:
//  o Pov is a location/orientation.
//  o Pov is a transform hierarchy.
//  o Pov passes user object creation requests through the engine entity manager.

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

	// Used to register entities with the eng component manager.
	eng *engine // Overall component manager.
	id  eid     // Unique entity identifier.

	// transform hierarchy. Each pov node can have children which base
	// their position and orientation relative to the parents.
	parent *Pov   // nil if no parent: nil means root.
	kids   []*Pov // child transforms.
	stable bool   // avoid updating non-moving objects.

	// variables for recalculating transforms each update.
	toc float64 // distance to camera.
	mm  *lin.M4 // model matrix world transform.
	rot *lin.Q  // scratch rotation/orientation.
}

// newPov allocates and initialzes a point of view transform.
// Called by the engine.
func newPov(eng *engine, id eid, parent *Pov) *Pov {
	p := &Pov{eng: eng, id: id, parent: parent}
	p.T = lin.NewT()
	p.S = &lin.V3{X: 1, Y: 1, Z: 1}

	// add the pov to its parent.
	if parent != nil {
		parent.kids = append(parent.kids, p)
	}

	// allocate scratch variables.
	p.rot = lin.NewQ()
	p.mm = &lin.M4{}
	return p
}

// NewPov creates, attaches, and returns a new child transform Pov.
func (p *Pov) NewPov() *Pov { return p.eng.newPov(p) }

// Dispose deletes and removes either the entire PovNode,
// or one of its attachments, PovModel, PovBody, ...
func (p *Pov) Dispose(kind int) { p.eng.dispose(p.id, kind) }

// At gets the Pov's location.
// Location is relative to parent. World space if no parent location.
func (p *Pov) At() (x, y, z float64) {
	return p.T.Loc.X, p.T.Loc.Y, p.T.Loc.Z
}

// SetAt sets the Pov at the given location.
// Location is relative to parent. World space if no parent location.
func (p *Pov) SetAt(x, y, z float64) *Pov {
	p.stable = false
	p.T.Loc.X, p.T.Loc.Y, p.T.Loc.Z = x, y, z
	return p
}

// Move directly affects the location by the given translation amounts
// along the given direction. Physics bodies should use Body.Push which
// affects velocity.
func (p *Pov) Move(x, y, z float64, dir *lin.Q) {
	p.stable = false
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
func (p *Pov) View() (q *lin.Q) {
	p.stable = false // referenced sometimes used to change.
	return p.T.Rot
}

// SetView directly sets the Pov's orientation.
// Often used to align this Pov with the orientation of another.
// Orientation is relative to parent. World space if no parent orientation.
func (p *Pov) SetView(q *lin.Q) *Pov {
	p.stable = false
	r := p.T.Rot
	r.X, r.Y, r.Z, r.W = q.X, q.Y, q.Z, q.W
	return p
}

// Spin rotates x,y,z degrees about the X,Y,Z axis.
// The spins are combined in XYZ order, but generally this
// is used to spin about a single axis at a time.
func (p *Pov) Spin(x, y, z float64) {
	p.stable = false
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
// Setting through this method causes world transform update while setting
// the Pov.S directly does not.
func (p *Pov) SetScale(x, y, z float64) *Pov {
	p.stable = false
	p.S.X, p.S.Y, p.S.Z = x, y, z
	return p
}

// Cam returns nil if there is no camera for this Pov.
func (p *Pov) Cam() *Camera { return p.eng.cam(p.id) }

// NewCam adds a camera to a Pov. This means that all rendered models
// in this Pov's hierarchy will be viewed with this camera settings.
func (p *Pov) NewCam() *Camera { return p.eng.newCam(p.id) }

// Body is an optional physics component associated with a Pov.
// Body returns nil if there is no physics body for this Pov.
func (p *Pov) Body() physics.Body { return p.eng.body(p.id) }

// NewBody creates non-colliding body associated with this Pov.
// Bodies are generally set on top level Pov transforms to ensure
// valid world coordindates.
func (p *Pov) NewBody(b physics.Body) physics.Body {
	return p.eng.newBody(p.id, b, p)
}

// SetSolid makes the existing Pov physics Body collide. Nothing happens
// if there is no existing physics Body for this Pov.
func (p *Pov) SetSolid(mass, bounce float64) {
	p.eng.setSolid(p.id, mass, bounce)
}

// AddSound loads audio data associated with this Pov. Played sounds occur
// at the associated Pov's location. Sounds that are played will be louder
// as the distance between the played noise and listener decreases.
// Place the single global noise listener at this Pov.
func (p *Pov) AddSound(name string) { p.eng.addSound(p.id, name) }

// PlaySound plays the sound associated with this Pov.
// The sound index corresponds to the order the sound was added
// to the Pov. The first sound added has index 0.
func (p *Pov) PlaySound(index int) { p.eng.playSound(p.id, index) }

// SetListener sets the location of the listener to be this Pov.
func (p *Pov) SetListener() { p.eng.sounds.setListener(p) }

// Light returns nil if no there light for this Pov.
func (p *Pov) Light() *Light { return p.eng.lights.get(p.id) }

// NewLight creates and associates a light with this Pov.
// Light is optional. It affects lighting calculations for this Pov
// and all child pov's.
func (p *Pov) NewLight() *Light { return p.eng.lights.create(p.id) }

// Layer returns nil if there is no layer for this Pov.
func (p *Pov) Layer() Layer { return p.eng.layer(p.id) }

// NewLayer creates a rendered texture at this Pov.
// The model at this Pov and its children will be rendered to a texture
// for this layer. The default texture size is 1024x1024.
func (p *Pov) NewLayer() Layer { return p.eng.newLayer(p.id, render.ImageBuffer) }

// NewLabel creates a Model that displays a small text phrase.
//   shader      is expected to be a texture mapping shader like "uv".
//   font        is the font mapping file.
//   fontTexture is the bitmapped font texture
func (p *Pov) NewLabel(shader, font, fontTexture string) Labeler {
	return p.NewModel(shader, "fnt:"+font, "tex:"+fontTexture)
}

// Model returns nil if there is no model for this Pov.
func (p *Pov) Model() Model { return p.eng.models.get(p.id) }

// NewModel creates an optional rendered component associated with this Pov.
// Returns nil if a model already exists.
func (p *Pov) NewModel(shader string, attrs ...string) Model {
	return p.eng.models.create(p.id, shader, attrs...)
}

// remChild is used by a pov removing itself from the hierarchy.
func (p *Pov) remChild(c *Pov) {
	for index, c := range p.kids {
		if c.id == p.id {
			p.kids = append(p.kids[:index], p.kids[index+1:]...)
			return
		}
	}
}

// Pov
// =============================================================================
// povs
// FUTURE: break the Pov instance fields into individual arrays
//         as per Data Oriented programming and see if this speeds
//         up transform processing. Benchmark!

// povs manages all the active Pov instances.
// Pov data is kept internally in order to facilitate optimizing
// updateTransform which is called each update.
type povs struct {
	index map[eid]uint32 // Map sparse entity-id to dense slice data.
	data  []*Pov         // Dense array of Pov data...
	eids  []eid          // ...and associated entity identifiers.
}

// newPovs creates a manager for a group of Pov data.
// There is only expected to be once instance created by the engine.
func newPovs() *povs {
	return &povs{data: []*Pov{}, eids: []eid{}, index: map[eid]uint32{}}
}

// create a new Pov. Guarantees that child Pov's appear later in the
// dense data array since children must be created after their parents.
func (ps *povs) create(eng *engine, id eid, parent *Pov) *Pov {
	p := newPov(eng, id, parent)

	// add the pov and update the pov indicies.
	ps.data = append(ps.data, p)
	ps.eids = append(ps.eids, p.id)
	ps.index[p.id] = uint32(len(ps.data)) - 1
	return p
}

// dispose deletes the given Pov and all of its children.
// The children are deleted since they are dependent on the parent
// for their location. This also keeps all children located after
// their parent in the data array.
func (ps *povs) dispose(id eid) {
	delIndex, ok := ps.index[id] // index to item for removal.
	if !ok {
		return // handle deletes on entities that don't exist.
	}
	deletee := ps.data[delIndex] // Pov to be removed.

	// delete the requested item. Preserve order so that parents
	// continue to appear before their children.
	ps.data = append(ps.data[:delIndex], ps.data[delIndex+1:]...)
	ps.eids = append(ps.eids[:delIndex], ps.eids[delIndex+1:]...)
	for cnt, eid := range ps.eids[delIndex:] {
		ps.index[eid] = delIndex + uint32(cnt)
	}

	// Deleting a Pov is the same as deleting an entity.
	// Ensure other components get a chance to update by informing the engine.
	if deletee.parent != nil {
		deletee.parent.remChild(deletee) // remove the one back reference that matters.
	}
	for _, kid := range deletee.kids {
		kid.parent = nil               // avoid unnecessary removing of back references
		deletee.eng.disposePov(kid.id) // may delete other entities.
	}
}

// get the Pov for the given id, returning nil if
// it does not exist.
func (ps *povs) get(id eid) *Pov {
	if index, ok := ps.index[id]; ok {
		return ps.data[index]
	}
	return nil
}

// updateWorldTransforms ensures that world transforms match any changes to location
// and orientation. Child transforms are relative to their parents. Thus parent
// transforms must be, and are, positioned earlier in the slice than their children.
func (ps *povs) updateWorldTransforms() {
	for _, p := range ps.data {
		if p.stable {
			continue // ignore things that haven't moved.
		}
		p.mm.SetQ(p.rot.Inv(p.T.Rot))   // invert model rotation.
		p.mm.ScaleSM(p.S.GetS())        // scale is applied first: left of rotation.
		l := p.T.Loc                    // translate is applied last...
		p.mm.TranslateMT(l.X, l.Y, l.Z) // ...right of rotation.

		// world transform of a child is relative to its parent location.
		if p.parent != nil {
			p.mm.Mult(p.mm, p.parent.mm) // model transform + parent transform
		}

		// mark as processed and ensure all children are updated.
		p.stable = true
		for _, kid := range p.kids {
			kid.stable = false
		}
	}
}

// reset the pov manager dumping old data for garbage collection.
func (ps *povs) reset() {
	ps.data = []*Pov{}
	ps.index = map[eid]uint32{}
}
