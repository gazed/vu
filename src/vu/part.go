// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package vu

import (
	"vu/math/lin"
	"vu/move"
	"vu/render"
)

// A part is a node in a scene graph (a part in the scene, yah?).  Its main
// responsibility is positioning. A part can also be associated with one or more
// of the following:
//      facade  : A part can have a visible surface. The parts location becomes the
//                location of the rendered object. SetFacade() and SetBanner() both
//                associate a facade with a part.
//      physics : A part can interface with physics. The parts location becomes
//                controlled by the physics simulation. SetBody() associates a
//                physics body with a part.
// Parts are scene graph nodes and as such can themselves have subordinate parts.
//
// It is possible to have a transform-only part (no facade or physics) in order to
// group other objects in the scene graph. All subordinate parts will be affected
// by transform changes to parent parts.
type Part interface {

	// Create, remove, and dispose parts and sub-parts (bit-parts?).
	AddPart() Part  // Create a new part attached (subordinate) to this one.
	RemPart(p Part) // Remove a subordinate part, without disposing the sub-part.
	Dispose()       // Remove this part, all sub-parts and any associated info.

	// Location, orientation, and scale.
	Location() (x, y, z float64)    // Get the current location.
	SetLocation(x, y, z float64)    // Set the current location.
	Rotation() (x, y, z, w float64) // Get the current quaternion rotation.
	SetRotation(x, y, z, w float64) // Set the current quaternion rotation.
	Scale() (x, y, z float64)       // Get the size, one value for each axis.
	SetScale(x, y, z float64)       // Set the size, one value for each axis.

	// Move and turn control cameras, players, and AIs. When there is a physics
	// body the velocity is updated. When there is no physics body the position
	// and direction are upated.
	Move(x, y, z float64) // Move an amount along the current direction.
	Spin(x, y, z float64) // Rotate degrees about the given axis.

	// Associates a visible facade with a Part. The facade methods that follow
	// only work when there is a facade associated with the part.
	SetFacade(mesh, shader string) Part   // Create a facade for this part.
	SetMaterial(name string)              // Set the material.
	Visible() bool                        // Get the parts visibility.
	SetVisible(visible bool)              // Set the parts visibility.
	SetCullable(cullable bool)            // Set the parts cullability.
	Alpha() float64                       // Get the transparency.
	SetAlpha(a float64)                   // Set the transparency.
	SetTexture(name string, spin float64) // Set the texture.

	// Associates a text facade with a Part. The banner methods only work
	// when there is a Banner associated with the part.  Only one of banner
	// or facade may be associated with a part.
	SetBanner(text, shader, glyphs, texture string) // Create a banner for this part.
	UpdateBanner(text string)                       // Change the text for the banner.
	BannerWidth() int                               // Get the banner width in pixels.

	// Associates a physics body with a Part. The body methods that follow
	// only work when there is a physics body associated with the part.
	SetBody(body move.Body, mass, bounce float64) // Create a physics body for this part.
	RemBody()                                     // Delete the physics body for this part.
	Speed() (x, y, z float64)                     // Current linear velocity.
	Push(x, y, z float64)                         // Change the linear velocity.
	Turn(x, y, z float64)                         // Change the angular velocity.
	Stop()                                        // Remove all velocity.
}

// Part interface
// ===========================================================================
// part - Part implementation

// part implements the Part interface.
type part struct {
	pov                  // Embed the pov location and orientation struct.
	staged   bool        // True if the the part is visible.
	scale    *lin.V3     // Scale, per axis: >1 to enlarge, 0<1 to shrink.
	parts    []*part     // Each scene node can have 1 or more scene nodes.
	face     *facade     // Parts can have one optional facade.
	cullable bool        // Can/can't be culled is under control of the application.
	culled   bool        // Draw or don't under control of engine.
	visible  bool        // Draw or don't under control of application.
	toc      float64     // Distance to center (to->c) for sorting and culling.
	vis      *render.Vis // Current render data (visible).
	body     move.Body   // Motion body used mostly by physics subsystem.
	world    move.World  // Used to add and removed bodies.

	// scratch variables are used each render cycle. They are optimizations that
	// prevent having to create temporary structures each render cycle.
	model *lin.M4 // Calculates model transform each render cycle.
	pm    *lin.M4 // Scratch parent model transform.
	vm    *lin.M4 // Scratch view model transform.
}

// newPart creates and initialzes a part instance.
func newPart(world move.World) *part {
	p := &part{}
	p.world = world
	p.loc = &lin.V3{}
	p.dir = &lin.Q{0, 0, 0, 1}
	p.scale = &lin.V3{1, 1, 1}
	p.parts = []*part{}
	p.visible = true
	p.cullable = true
	p.culled = false

	// scratch variables.
	p.model = &lin.M4{}
	p.pm = &lin.M4{}
	p.vm = &lin.M4{}

	// allocate the visible structures once and reuse.
	p.vis = render.NewVis()
	return p
}

// Part interface implementation.
func (p *part) Dispose() {
	for _, child := range p.parts {
		child.Dispose()
	}
	p.RemBody()
	p.parts = nil
	p.face = nil
	p.body = nil
	p.scale = nil
}

// Part interface implementation.
func (p *part) AddPart() Part {
	np := newPart(p.world)
	p.parts = append(p.parts, np)
	return np
}

// Part interface implementation.
// Find and remove the part (will point to the same record).
func (p *part) RemPart(child Part) {
	if pt, _ := child.(*part); pt != nil {
		for index, partPtr := range p.parts {
			if partPtr == pt {
				p.RemBody()
				p.parts = append(p.parts[:index], p.parts[index+1:]...)
				return
			}
		}
	}
}

// SetLocation directly updates the parts location to the given coordinates.
// This is a form of teleportation when the part has an active physics body.
func (p *part) SetLocation(x, y, z float64) {
	p.pov.SetLocation(x, y, z)
	if p.body != nil {
		p.body.World().Loc.SetS(x, y, z)
	}
}

// SetRotation directly updates the parts rotation to the given direction.
// This is a form of teleportation when the part has an active physics body.
func (p *part) SetRotation(x, y, z, w float64) {
	p.pov.SetRotation(x, y, z, w)
	if p.body != nil {
		p.body.World().Rot.SetS(x, y, z, w)
	}
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
func (p *part) Scale() (x, y, z float64)  { return p.scale.X, p.scale.Y, p.scale.Z }
func (p *part) SetScale(x, y, z float64)  { p.scale.X, p.scale.Y, p.scale.Z = x, y, z }
func (p *part) Visible() bool             { return p.visible }
func (p *part) SetVisible(visible bool)   { p.visible = visible }
func (p *part) SetCullable(cullable bool) { p.cullable = cullable }
func (p *part) SetFacade(mesh, shader string) Part {
	p.face = newFacade(mesh, shader)
	return p
}

// Part interface implementation.
func (p *part) SetMaterial(material string) {
	if p.face != nil {
		p.face.mat = material
	}
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
	p.world.Add(p.body)
}

// Part interface implementation.
func (p *part) RemBody() {
	if p.body != nil {
		p.world.Rem(p.body)
		p.body = nil
	}
}

// Part interface implementation.
func (p *part) SetTexture(texture string, rotSpeed float64) {
	if p.face != nil {
		p.face.tex, p.face.rots = texture, rotSpeed
	}
}

// Part interface implementation.
func (p *part) Alpha() float64 {
	if p.face != nil {
		return p.face.alpha
	}
	return 0
}

// Part interface implementation.
func (p *part) SetAlpha(a float64) {
	if p.face != nil {
		p.face.alpha = a
	}
}

// Part interface implementation.
func (p *part) SetBanner(text, shader, glyphs, texture string) {
	p.face = newBanner(text, shader, glyphs, texture)
}

// Part interface implementation.
func (p *part) UpdateBanner(text string) {
	if p.face != nil {
		p.face.text = text
	}
}

// Part interface implementation.
func (p *part) BannerWidth() int {
	if p.face != nil {
		return p.vis.GlyphWidth
	}
	return 0
}

// distanceTo returns the distance squared of the part to the given center.
func (p *part) distanceTo(cenx, ceny, cenz float64) float64 {
	dx := p.loc.X - cenx
	dy := p.loc.Y - ceny
	dz := p.loc.Z - cenz
	return float64(dx*dx + dy*dy + dz*dz)
}

// outside returns true if the node's distance to center is bigger than the
// given radius.
func (p *part) outside(radius float64) bool { return p.toc > float64(radius*radius) }

// model transform must be done in rotation, scale, translate order.
func (p *part) mt() *lin.M4 {
	mt := p.model.SetQ(lin.NewQ().Inv(p.dir)) // rotation.
	mt.ScaleSM(p.Scale())                     // scale is applied first (on left of rotation)
	return mt.TranslateMT(p.Location())       // translate is applied last (on right of rotation).
}

// temporary matrix used and reused for model view transforms.
var tm = &lin.M4{}

// stage the part for rendering. This takes the parts rendering specific information
// and copies it into a rendering structure.
func (p *part) stage(visible *[]*render.Vis, sc *scene, parentTransform *lin.M4) {
	if p.visible {
		m := p.mt()
		m.Mult(m, parentTransform) // model transform + parent transform

		// only render nodes with facades or banners.
		// transfer the rendering information in a graphics structure.
		if p.face != nil {
			vis := p.vis
			vis.L = sc.L
			vis.Mv = RenderMatrix(tm.Mult(m, sc.vt(p.vm)), vis.Mv) // generate the model-view transform
			vis.Mvp = RenderMatrix(tm.Mult(tm, sc.P), vis.Mvp)     // generate model-view-projection transform

			// both banners and facades render the same way. Should only be one
			// specified, but prefer a facade over a banner in the case of dev error.
			face := p.face
			vis.Is2D = sc.is2D
			vis.MeshName = face.mesh
			vis.ShaderName = face.shader
			vis.MatName = face.mat
			vis.TexName = face.tex
			vis.RotSpeed = float32(face.rots)
			vis.GlyphName = face.glyphs
			vis.GlyphText = face.text
			vis.Alpha = float32(face.alpha)
			vis.Scale.X, vis.Scale.Y, vis.Scale.Z = float32(p.scale.X), float32(p.scale.Y), float32(p.scale.Z)

			// Use a large fade default for scenes without radius.
			vis.Fade = 1000
			if sc.radius != 0 {
				vis.Fade = float32(sc.radius)
			}
			*visible = append(*visible, vis)
		}

		// render all the parts children
		for _, child := range p.parts {
			if !child.culled {
				p.pm.Set(m) // ensures the original model transform does not change.
				child.stage(visible, sc, p.pm)
			}
		}
	}
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
