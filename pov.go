// Copyright © 2013-2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// pov.go controls world transforms and scene graph child parent
//        relationships.

import (
	"log"

	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/render"
)

// AddPart creates a new entity with a point-of-view component (pov).
// A pov adds a location and orientation to an entity. The entity can
// now be positioned and rotated.
//
// The entity is also added to the scene graph so that this entities
// world pov is affected by its parents and will also affect any
// child entities created from this one.
func (e *Ent) AddPart() *Ent {
	eid := e.app.eids.create()
	e.app.povs.create(eid, e.eid)
	return &Ent{app: e.app, eid: eid}
}

// At gets the local space location. This is world space
// if the entity does not have a parent.
//
// Depends on Ent.AddPart. Returns 0,0,0 if there is no part component.
func (e *Ent) At() (x, y, z float64) {
	if p := e.app.povs.get(e.eid); p != nil {
		return p.at()
	}
	log.Printf("At needs AddPart %d", e.eid)
	return 0, 0, 0
}

// SetAt sets the local space location, ie: relative to its parent.
// This is world space if there is no parent location.
//
// Depends on Ent.AddPart.
func (e *Ent) SetAt(x, y, z float64) *Ent {
	if p := e.app.povs.get(e.eid); p != nil {
		p.tn.Loc.X, p.tn.Loc.Y, p.tn.Loc.Z = x, y, z
		return e
	}
	log.Printf("SetAt needs AddPart %d", e.eid)
	return e
}

// World returns the world space coordinates for this entity.
// World space is recalculated after every update.
//
// Depends on Ent.AddPart.
func (e *Ent) World() (wx, wy, wz float64) {
	if p := e.app.povs.get(e.eid); p != nil {
		return p.wx, p.wy, p.wz
	}
	log.Printf("World needs AddPart %d", e.eid)
	return 0, 0, 0
}

// Move directly affects the location by the given translation amounts
// along the given direction. Physics bodies should use Body.Push which
// affects velocity.
//
// Depends on Ent.AddPart.
func (e *Ent) Move(x, y, z float64, dir *lin.Q) {
	if p := e.app.povs.get(e.eid); p != nil {
		dx, dy, dz := lin.MultSQ(x, y, z, dir)
		p.tn.Loc.X += dx
		p.tn.Loc.Y += dy
		p.tn.Loc.Z += dz
		return
	}
	log.Printf("Move missing AddPart %d", e.eid)
}

// View returns the orientation of the pov. Orientation combines
// direction and rotation about the direction.
// Orientation is relative to parent. World space if no parent orientation.
//
// Depends on Ent.AddPart.
func (e *Ent) View() (q *lin.Q) {
	if p := e.app.povs.get(e.eid); p != nil {
		return p.tn.Rot
	}
	log.Printf("View needs AddPart %d", e.eid)
	return lin.NewQ()
}

// SetView directly sets the parts orientation.
// Often used to align this part with the orientation of another.
// Orientation is relative to parent. World space if no parent orientation.
//
// Depends on Ent.AddPart.
func (e *Ent) SetView(q *lin.Q) *Ent {
	if p := e.app.povs.get(e.eid); p != nil {
		r := p.tn.Rot
		r.X, r.Y, r.Z, r.W = q.X, q.Y, q.Z, q.W
		return e
	}
	log.Printf("SetView needs AddPart %d", e.eid)
	return e
}

// Cull sets the culled state.
//
// Depends on Ent.AddPart.
func (e *Ent) Cull(culled bool) {
	if n := e.app.povs.getNode(e.eid); n != nil {
		n.cull = culled
		return
	}
	log.Printf("Cull needs AddPart %d", e.eid)
}

// Culled returns true if entity has been culled from rendering.
//
// Depends on Ent.AddPart. Returns true if there was no part component.
func (e *Ent) Culled() bool {
	if n := e.app.povs.getNode(e.eid); n != nil {
		return n.cull
	}
	log.Printf("Culled needs AddPart %d", e.eid)
	return true
}

// Spin rotates x,y,z degrees about the X,Y,Z axis.
// The spins are combined in XYZ order, but generally this
// is used to spin about a single axis at a time.
//
// Depends on Ent.AddPart.
func (e *Ent) Spin(x, y, z float64) {
	if p := e.app.povs.get(e.eid); p != nil {
		p.spin(e.app.povs.rot, x, y, z)
		return
	}
	log.Printf("Spin needs AddPart %d", e.eid)
}

// SetSpin sets the rotation to 0 before spinning the entity
// like the Spin method.
//
// Depends on Ent.AddPart.
func (e *Ent) SetSpin(x, y, z float64) {
	if p := e.app.povs.get(e.eid); p != nil {
		p.clearSpin()
		p.spin(e.app.povs.rot, x, y, z)
		return
	}
	log.Printf("SetSpin needs AddPart %d", e.eid)
}

// Scale retrieves the local per-axis scale values at 3 separate XYZ values.
// World scale needs to incorporate any parents values.
//
// Depends on Ent.AddPart. Returns 0,0,0 if there is no part component.
func (e *Ent) Scale() (x, y, z float64) {
	if p := e.app.povs.get(e.eid); p != nil {
		return p.scale()
	}
	log.Printf("Scale needs AddPart %d", e.eid)
	return 0, 0, 0
}

// SetScale assigns the XYZ per-axis scale values.
// Scale default is 1, greater than 1 enlarges, a positive fraction shrinks.
//
// Depends on Ent.AddPart.
func (e *Ent) SetScale(x, y, z float64) *Ent {
	if p := e.app.povs.get(e.eid); p != nil {
		p.sn.X, p.sn.Y, p.sn.Z = x, y, z
		return e
	}
	log.Printf("SetScale needs AddPart %d", e.eid)
	return e
}

// pov entity methods
// =============================================================================
// pov data

// pov point-of-view, is a combination of position and orientation.
// A pov is created for each application entity.
//
// A pov's location factors in an update interpolation value to account
// for timing differences between rendering and updating.
//
// FUTURE: Don't use pointers for the transform data so that the pov
//         data is contiguous in memory. Will have to copy the transform
//         data in and out of physics instead of sharing the pointer.
//         In theory contiguous data means fewer cache misses.
//         An initial attempt at this made things slower.
type pov struct {
	eid eid // Unique entity identifier.

	// Local transform is relative to a parent.
	// Effectively a world transform if no parent transform.
	tp, tn     *lin.T  // Transform (prev, now) is location and orientation.
	sp, sn     *lin.V3 // Per axis scale (prev, now): default value 1,1,1.
	wx, wy, wz float64 // World coordinates set each display refresh.
	mm         *lin.M4 // model matrix world transform.
	stable     bool    // avoid updating non-moving objects.
}

// newPov allocates and initialzes a point of view transform.
// Called by the engine.
func newPov(eid eid) *pov {
	p := &pov{eid: eid}
	p.tn = lin.NewT()
	p.tp = lin.NewT()
	p.sn = &lin.V3{X: 1, Y: 1, Z: 1}
	p.sp = &lin.V3{X: 1, Y: 1, Z: 1}
	p.mm = &lin.M4{}
	return p
}

// at gets the pov's local space location.
// Local space since location is relative to parent.
// World space if no parent location.
func (p *pov) at() (x, y, z float64) {
	return p.tn.Loc.X, p.tn.Loc.Y, p.tn.Loc.Z
}

// scale retrieves the local per-axis scale values at 3 separate XYZ values.
// World scale needs to incorporate any parents values.
func (p *pov) scale() (x, y, z float64) { return p.sn.X, p.sn.Y, p.sn.Z }

// clearSpin sets the rotation to default 0.
func (p *pov) clearSpin() { p.tn.Rot.Set(lin.QI) }

// spin using the given x,y,z values and the scratch quaternion.
func (p *pov) spin(rot *lin.Q, x, y, z float64) {
	if x != 0 {
		rot.SetAa(1, 0, 0, lin.Rad(x))
		p.tn.Rot.Mult(rot, p.tn.Rot)
	}
	if y != 0 {
		rot.SetAa(0, 1, 0, lin.Rad(y))
		p.tn.Rot.Mult(rot, p.tn.Rot)
	}
	if z != 0 {
		rot.SetAa(0, 0, 1, lin.Rad(z))
		p.tn.Rot.Mult(rot, p.tn.Rot)
	}
}

// setLocal sets the local space model matrix for the pov.
// Local space ignores the affect of any parent pov.
// The location is an interpolation between the previous values
// and the current value.
func (p *pov) setLocal(lerp float64, rot *lin.Q, v3 *lin.V3) {
	var sx, sy, sz float64
	var lx, ly, lz float64

	// only bother with interpolated values if the two values are
	// relatively close. Otherwise the object has been teleported
	// and should just appear where it is now.
	threshold := 1.0
	if dist := p.tn.Loc.DistSqr(p.tp.Loc); dist > threshold || lerp == 1.0 {
		// get latest value - no sense interpolating.
		sx, sy, sz = p.sn.GetS()     // scale
		lx, ly, lz = p.tn.Loc.GetS() // position
		rot.Set(p.tn.Rot)            // orientation.
	} else {
		// interpolate from the previous values.
		sx, sy, sz = v3.Lerp(p.sp, p.sn, lerp).GetS()         // scale
		lx, ly, lz = v3.Lerp(p.tp.Loc, p.tn.Loc, lerp).GetS() // position
		rot.Nlerp(p.tp.Rot, p.tn.Rot, lerp)                   // orientation.
	}

	// Update the model matrix with
	p.mm.SetQ(rot.Inv(rot))      // invert model rotation.
	p.mm.ScaleSM(sx, sy, sz)     // scale is applied first: left of rotation.
	p.mm.TranslateMT(lx, ly, lz) // ...right of rotation.
}

// draw sets the location render data needed for a single draw call.
// The data is copied into a render.Draw instance. One of the key jobs
// of this method is to put each draw request into a particular
// render bucket so that they are drawn in order once sorted.
func (p *pov) draw(d *render.Draw, mv, mvp *lin.M4, sc *scene, cam *Camera, rt uint32) {
	d.SetMv(mv.Mult(p.mm, cam.vm)) // model-view
	d.SetMvp(mvp.Mult(mv, cam.pm)) // model-view-projection
	d.SetPm(cam.pm)                // projection only.
	d.SetScale(p.scale())
	d.Tag = uint32(p.eid) // for debugging.
	d.Fbo = rt            // 0 for standard display back buffer.
	d.Depth = !sc.is2D()
	d.Bucket = setBucket(uint8(rt), sc.overlay)
}

// pov
// =============================================================================
// povs
// FUTURE: break the pov instance fields into individual arrays
//         as per Data Oriented programming and see if this speeds
//         up transform processing. Benchmark!

// povs is the pov component manager.
type povs struct {
	// Data can change array location without updating eid references.
	// Each of the dense data arrays has indexed data.
	index map[eid]uint32 // Sparse entity-id to ordered slice data.
	povs  []pov          // Dense array of pov data...
	eids  []eid          // ...and associated entity identifiers.
	nodes []node         // Scene graph parent-child data.

	// Scratch for per update tick calculations.
	rot *lin.Q  // scratch rotation/orientation.
	v4  *lin.V4 // scratch vector location.
	v3  *lin.V3 // scratch vector location.
}

// newPovs creates a manager for a group of Pov data.
// There is only expected to be once instance created by the engine.
func newPovs() *povs {
	ps := &povs{}
	ps.povs = []pov{}
	ps.eids = []eid{}
	ps.index = map[eid]uint32{}
	ps.nodes = []node{}

	// allocate scratch variables. These are used each update when
	// updating world positions and rotations.
	ps.rot = lin.NewQ()
	ps.v4 = &lin.V4{}
	ps.v3 = &lin.V3{}
	return ps
}

// create a new pov. Guarantees that child pov's appear later in the
// dense data array since children must be created after their parents.
func (ps *povs) create(eid eid, parent eid) *pov {
	p := newPov(eid)

	// add the pov and update the pov indicies.
	index := len(ps.povs)
	ps.index[p.eid] = uint32(index)
	ps.povs = append(ps.povs, *p)
	ps.eids = append(ps.eids, p.eid)
	ps.nodes = append(ps.nodes, node{}) // node with no parent, no kids.

	// if not root then add the pov to its parent.
	if parent != 0 { // valid entities start at 1.
		(&ps.nodes[index]).parent = parent
		pi := ps.index[parent]
		(&ps.nodes[pi]).kids = append((&ps.nodes[pi]).kids, eid)
	}
	return p
}

// dispose deletes the given pov and all of its children.
// Returns a list of deleted child entities. The returned list does not
// contain eid - the passed in entity id.
func (ps *povs) dispose(id eid, dead []eid) []eid {
	di, ok := ps.index[id] // index to item being deleted.
	delete(ps.index, id)
	if !ok {
		return dead // ignore deletes for entities that do not exist.
	}
	node := ps.nodes[di]

	// delete the requested item. Order is preserved so that
	// parents continue to appear before their children.
	ps.povs = append(ps.povs[:di], ps.povs[di+1:]...)
	ps.eids = append(ps.eids[:di], ps.eids[di+1:]...)
	ps.nodes = append(ps.nodes[:di], ps.nodes[di+1:]...)

	// Fix up map indicies. Remove 1 from each index after the deleted index.
	for _, eid := range ps.eids[di:] {
		ps.index[eid] = ps.index[eid] - 1
	}

	// Remove deleted pov from its parent.
	if pi, ok := ps.index[node.parent]; ok {
		parent := &ps.nodes[pi]
		for cnt, kid := range parent.kids {
			if kid == id {
				parent.kids = append(parent.kids[:cnt], parent.kids[cnt+1:]...)
				break
			}
		}
	}

	// At this point capture orphaned child nodes for deletion.
	for _, kid := range node.kids {
		ki := ps.index[kid]
		(&ps.nodes[ki]).parent = 0 // mark as orphan.
		dead = append(dead, kid)   // child needs to be deleted.
		dead = ps.dispose(kid, dead)
	}
	return dead // entities that may have other components that need deleting.
}

// get the pov for the given id, returning nil if it does not exist.
// Pointer reference only valid for this call.
func (ps *povs) get(id eid) *pov {
	if index, ok := ps.index[id]; ok {
		return &ps.povs[index]
	}
	return nil
}

// getNode returns the scene graph parent child information.
func (ps *povs) getNode(id eid) *node {
	if index, ok := ps.index[id]; ok {
		return &ps.nodes[index]
	}
	return nil
}

// setPrev saves the previous locations and orientations.
// It is called each update. It is needed to interpolate values when
// multiple renders are called between state updates.
//
// FUTURE: immediate update so no updateWorld is necessary.
func (ps *povs) setPrev() {
	for i := 0; i < len(ps.povs); i++ {
		p := &ps.povs[i] // update reference, not copy.
		if !p.tp.Eq(p.tn) || !p.sp.Eq(p.sn) {
			p.stable = false // mark as having moved for updateWorld.
			p.tp.Set(p.tn)
			p.sp.Set(p.sn)
		}
	}
}

// renderAt sets the local world render matrix as an position/rotation
// interpolated between the last 2 updated. Called once per display
// frame update.
func (ps *povs) renderAt(lerp float64) {
	for index := 0; index < len(ps.povs); index++ {
		p := &ps.povs[index] // want to update reference, not copy.
		if p.stable {
			// Anything that hasn't moved already has the correct
			// transform matricies.
			continue
		}
		p.stable = true // mark as processed.

		// ensure all children are updated.
		node := &ps.nodes[index]
		for _, kid := range node.kids {
			if index, ok := ps.index[kid]; ok {
				(&ps.povs[index]).stable = false // update ref, not copy.
			} else {
				log.Printf("WTF Child")
			}
		}

		// update transform and world location.
		p.setLocal(lerp, ps.rot, ps.v3) // update local space transform.
		if node.parent != 0 {
			if pindex, ok := ps.index[node.parent]; ok {
				parent := &ps.povs[pindex] // use ref, not copy.
				// world transform of a child is relative to its parent.
				// Parent's model matrix has already been set because parent
				// pov's appear earlier in ps.data than their children.
				p.mm.Mult(p.mm, parent.mm) // model + parent transform
			} else {
				log.Printf("WTF Parent")
			}
		}

		// save the world location for other calculations.
		v4 := ps.v4
		v4.MultvM(v4.SetS(0, 0, 0, 1), p.mm)

		// FUTURE: Remove side effect with move to immediate world updating.
		p.wx, p.wy, p.wz = v4.X, v4.Y, v4.Z
	}
}

// povs
// =============================================================================
// node - also tracked by the pov component manager.

// node tracks the parent child relationship for an entity.
// It creates a scene graph transform hierarchy. Each node can have
// children which base their position and orientation relative to the parents.
//
// Used as part of the pov component manager to add child parent data
// to data that have position and orientation.
type node struct {
	parent eid   // Parent entity identifier.
	kids   []eid // Child entities.

	// Cull set to true removes this node and its children
	// from scene graph processing. Default false.
	cull bool // True to exclude from scene graph processing.
}
