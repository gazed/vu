// Copyright © 2013-2024 Galvanized Logic Inc.

package vu

// pov.go controls world transforms and scene graph child parent
//        relationships.

import (
	"log/slog"
	"time"

	"github.com/gazed/vu/math/lin"
)

// At gets the local space location. This is world space
// if the entity does not have a parent.
//
// Returns 0,0,0 if there is no transform component.
func (e *Entity) At() (x, y, z float64) {
	if p := e.app.povs.get(e.eid); p != nil {
		return p.at()
	}
	slog.Error("At needs transform", "eid", e.eid)
	return 0, 0, 0
}

// SetAt sets the local space location, ie: relative to its parent.
// This is world space if there is no parent location.
//
// Depends on transform.
func (e *Entity) SetAt(x, y, z float64) *Entity {
	if p := e.app.povs.get(e.eid); p != nil {
		if p.tn.Loc.X != x || p.tn.Loc.Y != y || p.tn.Loc.Z != z {
			p.tn.Loc.X, p.tn.Loc.Y, p.tn.Loc.Z = x, y, z
			e.app.povs.updateWorld(p, e.eid)
		}
		return e
	}
	slog.Error("SetAt needs transform", "eid", e.eid)
	return e
}

// World returns the world space coordinates for this entity.
// World space is recalculated immediately on any change.
//
// Depends on transform.
func (e *Entity) World() (wx, wy, wz float64) {
	if p := e.app.povs.get(e.eid); p != nil {
		return p.world()
	}
	slog.Error("World needs transform", "eid", e.eid)
	return 0, 0, 0
}

// WorldRot returns the world rotation for this entity.
// WorldRot space is recalculated immediately on any change
// and returns nil if the entity does not have a part.
//
// Depends on transform.
func (e *Entity) WorldRot() (q *lin.Q) {
	if p := e.app.povs.get(e.eid); p != nil {
		return p.tw.Rot
	}
	slog.Error("WorldRot needs transform", "eid", e.eid)
	return nil
}

// Move directly affects the location by the given translation amounts
// along the given direction. Physics bodies should use Body.Push which
// affects velocity.
//
// Depends on transform.
func (e *Entity) Move(x, y, z float64, dir *lin.Q) {
	if p := e.app.povs.get(e.eid); p != nil {
		dx, dy, dz := lin.MultSQ(x, y, z, dir)
		p.tn.Loc.X += dx
		p.tn.Loc.Y += dy
		p.tn.Loc.Z += dz
		e.app.povs.updateWorld(p, e.eid)
		return
	}
	slog.Error("Move missing transform", "eid", e.eid)
}

// View returns the orientation of the Part. Orientation is a rotation
// about a axis. Orientation is relative to any parent Parts.
// It is world space if there is no parent orientation. Direct updates
// to the rotation matrix must be done with SetView or SetAa.
//
// Depends on transform.
func (e *Entity) View() (q *lin.Q) {
	if p := e.app.povs.get(e.eid); p != nil {
		return p.tn.Rot
	}
	slog.Error("View needs transform", "eid", e.eid)
	return lin.NewQ()
}

// SetView directly sets the parts orientation.
// Often used to align this part with the orientation of another.
// Orientation is relative to parent. World space if no parent orientation.
//
// Depends on transform.
func (e *Entity) SetView(q *lin.Q) *Entity {
	if p := e.app.povs.get(e.eid); p != nil {
		r := p.tn.Rot
		r.X, r.Y, r.Z, r.W = q.X, q.Y, q.Z, q.W
		e.app.povs.updateWorld(p, e.eid)
		return e
	}
	slog.Error("SetView needs transform", "eid", e.eid)
	return e
}

// SetAa sets the orientation using the given axis and angle
// information.
//
// Depends on transform.
func (e *Entity) SetAa(x, y, z, angleInRadians float64) *Entity {
	if p := e.app.povs.get(e.eid); p != nil {
		p.tn.Rot.SetAa(x, y, z, angleInRadians)
		e.app.povs.updateWorld(p, e.eid)
		return e
	}
	slog.Error("SetAa needs transform", "eid", e.eid)
	return e
}

// Cull sets the culled state.
//
// Depends on transform.
func (e *Entity) Cull(culled bool) {
	if n := e.app.povs.getNode(e.eid); n != nil {
		n.cull = culled
		return
	}
	slog.Error("Cull needs transform", "eid", e.eid)
}

// Culled returns true if entity has been culled from rendering.
//
// Depends on transform. Returns true if there was no part component.
func (e *Entity) Culled() bool {
	if n := e.app.povs.getNode(e.eid); n != nil {
		return n.cull
	}
	slog.Error("Culled needs transform", "eid", e.eid)
	return true
}

// Spin rotates x,y,z degrees about the X,Y,Z axis.
// The spins are combined in XYZ order, but generally this
// is used to spin about a single axis at a time.
//
// Depends on transform.
func (e *Entity) Spin(x, y, z float64) *Entity {
	if p := e.app.povs.get(e.eid); p != nil {
		p.spin(e.app.povs.rot, x, y, z)
		e.app.povs.updateWorld(p, e.eid)
		return e
	}
	slog.Error("Spin needs transform", "eid", e.eid)
	return e
}

// SetSpin sets the rotation to 0 before spinning the entity
// like the Spin method.
//
// Depends on transform.
func (e *Entity) SetSpin(x, y, z float64) *Entity {
	if p := e.app.povs.get(e.eid); p != nil {
		p.clearSpin()
		p.spin(e.app.povs.rot, x, y, z)
		e.app.povs.updateWorld(p, e.eid)
		return e
	}
	slog.Error("SetSpin needs transform", "eid", e.eid)
	return e
}

// Scale retrieves the local per-axis scale values at 3 separate XYZ values.
// World scale needs to incorporate any parents values.
//
// Depends on transform. Returns 0,0,0 if there is no part component.
func (e *Entity) Scale() (x, y, z float64) {
	if p := e.app.povs.get(e.eid); p != nil {
		return p.scale()
	}
	slog.Error("Scale needs transform", "eid", e.eid)
	return 0, 0, 0
}

// SetScale assigns the XYZ per-axis scale values.
// Scale default is 1, greater than 1 enlarges, a positive fraction shrinks.
//
// Depends on transform.
func (e *Entity) SetScale(x, y, z float64) *Entity {
	if p := e.app.povs.get(e.eid); p != nil {
		p.sn.X, p.sn.Y, p.sn.Z = x, y, z
		e.app.povs.updateWorld(p, e.eid)
		return e
	}
	slog.Error("SetScale needs transform", "eid", e.eid)
	return e
}

// AddPart creates a new entity with a point-of-view component (pov).
// A pov adds a location and orientation to an entity. The entity can
// now be positioned and rotated.
//
// The entity is also added to the scene graph so that this entities
// world pov is affected by its parents and will also affect any
// child entities created from this one.
func (e *Entity) AddPart() *Entity {
	eid := e.app.eids.create()
	e.app.povs.create(eid, e.eid) // add new entity to parent.
	return &Entity{app: e.app, eid: eid}
}

// =============================================================================
// pov data

// pov point-of-view, is a combination of position and orientation.
// A pov is created for each application entity.
//
// A pov's location factors in an update interpolation value to account
// for timing differences between rendering and updating.
//
// FUTURE: Don't use pointers for the transform data so that the pov
//
//	data is contiguous in memory. Will have to copy the transform
//	data in and out of physics instead of sharing the pointer.
//	In theory contiguous data means fewer cache misses.
//	An initial attempt at this made things slower.
type pov struct {
	eid eID // Unique entity identifier.

	// Local transform is relative to a parent.
	// World transform combine parent transform.
	tp, tn *lin.T  // Local transform (prev, now).
	sp, sn *lin.V3 // Per axis scale (prev, now): default value 1,1,1.
	tw     *lin.T  // World transform. Updated on any change.
	sw     *lin.V3 // World scale. Updated on any change.
	mm, wm *lin.M4 // render model matrix, world matrix.
	stable bool    // avoid updating non-moving objects.
}

// newPov allocates and initialzes a point of view transform.
// Called by the engine.
func newPov(eid eID) *pov {
	p := &pov{eid: eid}
	p.tn = lin.NewT()
	p.tp = lin.NewT()
	p.tw = lin.NewT()
	p.sn = &lin.V3{X: 1, Y: 1, Z: 1}
	p.sp = &lin.V3{X: 1, Y: 1, Z: 1}
	p.sw = &lin.V3{X: 1, Y: 1, Z: 1}
	p.mm = lin.NewM4I()
	p.wm = lin.NewM4I()
	return p
}

// at gets the pov's local space location.
// Local space since location is relative to parent.
// World space if no parent location.
func (p *pov) at() (x, y, z float64) {
	return p.tn.Loc.X, p.tn.Loc.Y, p.tn.Loc.Z
}

// world get the pov's world space.
func (p *pov) world() (x, y, z float64) {
	return p.tw.Loc.X, p.tw.Loc.Y, p.tw.Loc.Z
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

// =============================================================================
// povs

// povs is the pov component manager.
type povs struct {
	// Data can change array location without updating eid references.
	// Each of the dense data arrays has indexed data.
	index map[eID]uint32 // Sparse entity-id to ordered slice data.
	povs  []pov          // Dense array of pov data...
	eids  []eID          // ...and associated entity identifiers.
	nodes []node         // Scene graph parent-child data.

	// Scratch for per update tick calculations.
	rot *lin.Q  // scratch rotation/orientation.
	v4  *lin.V4 // scratch vector location.
	v3  *lin.V3 // scratch vector location.
	m3  *lin.M3 // scratch rotation matrix.
}

// newPovs creates a manager for a group of Pov data.
// There is only expected to be once instance created by the engine.
func newPovs() *povs {
	ps := &povs{}
	ps.povs = []pov{}
	ps.eids = []eID{}
	ps.index = map[eID]uint32{}
	ps.nodes = []node{}

	// allocate scratch variables. These are used each update when
	// updating world positions and rotations.
	ps.rot = lin.NewQ()
	ps.v4 = &lin.V4{}
	ps.v3 = &lin.V3{}
	ps.m3 = &lin.M3{}
	return ps
}

// create a new pov. Guarantees that child pov's appear later in the
// dense data array since children must be created after their parents.
func (ps *povs) create(eid eID, parent eID) *pov {
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
		ps.updateWorld(p, eid) // ensure initial world transforms.
	}
	return p
}

// dispose deletes the given pov and all of its children.
// Returns a list of deleted child entities. The returned list does not
// contain eid - the passed in entity id.
func (ps *povs) dispose(id eID, dead []eID) []eID {
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
func (ps *povs) get(id eID) *pov {
	if index, ok := ps.index[id]; ok {
		return &ps.povs[index]
	}
	return nil
}

// getNode returns the scene graph parent child information.
func (ps *povs) getNode(id eID) *node {
	if index, ok := ps.index[id]; ok {
		return &ps.nodes[index]
	}
	return nil
}

// setPrev saves the previous locations and orientations.
// It is called each update. It is needed to interpolate values when
// multiple renders are called between state updates.
// The eids of the moved povs are returned.
func (ps *povs) setPrev(moved []eID) []eID {
	moved = moved[:0] // reset preserving memory.
	for index := 0; index < len(ps.povs); index++ {
		p := &ps.povs[index] // update reference, not copy.
		if !p.stable {
			moved = append(moved, p.eid)
			p.tp.Set(p.tn)
			p.sp.Set(p.sn)
			p.stable = true
		}
	}
	return moved
}

// setWorldMatrix sets the local world render matrix.
// Called once per render to set the pov.mm model matrix used for rendering.
func (ps *povs) setWorldMatrix(delta time.Duration) {
	for index := 0; index < len(ps.povs); index++ {
		p := &ps.povs[index]

		// FUTURE: interpolate here; or use some other smooth render option.

		// Use the latest transform updated by updateWorld.
		p.mm.Set(p.wm) // copied on first render.
	}
}

// updateWorld sets the world location for the given pov.
// Called immediately on any change to any of the existing transform values.
// Expected to be called for each object update to immediately refresh the
// world transform values.
func (ps *povs) updateWorld(p *pov, eid eID) {
	rot := ps.rot
	if index, ok := ps.index[eid]; ok {
		p.stable = false              // object has changed.
		sx, sy, sz := p.sn.GetS()     // scale
		lx, ly, lz := p.tn.Loc.GetS() // position
		rot.Set(p.tn.Rot)             // orientation.

		// Update the model transform matrix the world space coordinates.
		p.wm.SetQ(rot.Inv(rot))      // invert model rotation.
		p.wm.ScaleSM(sx, sy, sz)     // scale is applied first: left of rotation.
		p.wm.TranslateMT(lx, ly, lz) // translation applied last: right of rotation.

		// Combine with parent transform. The world transform of a child is
		// relative to its parent. Parent's model matrix has already been set
		// because parent pov's appear earlier in ps.data than their children.
		node := &ps.nodes[index]
		if node.parent != 0 {
			if pindex, ok := ps.index[node.parent]; ok {
				parent := &ps.povs[pindex] // use ref, not copy.
				p.wm.Mult(p.wm, parent.wm) // model + parent transform
			} else {
				slog.Error("scene graph missing child", "entity", node.parent) // dev error.
			}
		}

		// Track absolute world transform values.
		// See https://math.stackexchange.com/questions/237369/ and
		// note the limitations when using uneven or negative scales.
		m := p.wm
		p.tw.Loc.SetS(m.Wx, m.Wy, m.Wz) // world space position.
		sx = ps.v3.SetS(m.Xx, m.Xy, p.wm.Xz).Len()
		sy = ps.v3.SetS(m.Yx, m.Yy, p.wm.Yz).Len()
		sz = ps.v3.SetS(m.Zx, m.Zy, p.wm.Zz).Len()
		p.sw.SetS(sx, sy, sz) // world scale
		ps.m3.SetS(
			m.Xx/sx, m.Xy/sx, p.wm.Xz/sx,
			m.Yx/sy, m.Yy/sy, p.wm.Yz/sy,
			m.Zx/sz, m.Zy/sz, p.wm.Zz/sz)
		p.tw.Rot.SetM3(ps.m3)  // world rotation.
		p.tw.Rot.Inv(p.tw.Rot) // Undo model matrix invert.

		// Child nodes must also be updated.
		for _, kid := range node.kids {
			if index, ok := ps.index[kid]; ok {
				ps.updateWorld(&ps.povs[index], kid)
			} else {
				slog.Error("Scene graph missing child.") // dev error.
			}
		}
	}
}

// =============================================================================
// node - also tracked by the pov component manager.

// node tracks the parent child relationship for an entity.
// It creates a scene graph transform hierarchy. Each node can have
// children which base their position and orientation relative to the parents.
//
// Used as part of the pov component manager to add child parent data
// to data that have position and orientation.
type node struct {
	parent eID   // Parent entity identifier.
	kids   []eID // Child entities.

	// Cull set to true removes this node and its children
	// from scene graph processing. Default false.
	cull bool // True to exclude from scene graph processing.
}
