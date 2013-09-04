// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package vu

import (
	"vu/math/lin"
	"vu/physics"
	"vu/render"
)

// Part primarly describes where something is placed in a 3D scene. Hence, a part
// is also a node in a scene graph (a part in the scene, yah?). While the main
// resposibility of a part is location and orientation, a part is often associated
// with one or more of the following:
//      facade  : a part can have a visible surface.  The parts location becomes the
//                location of the rendered object. Methods that affect a facade can
//                be directly accessed on a part. SetFacade() and SetBanner() both
//                associate a facade with a part.
//      physics : a part can interface with physics.  The parts location becomes
//                either controlled or influenced by the physics simulation and
//                interactions with other physics bodies.
// Parts are scene graph nodes and can themselves have subordinate parts.
//
// It is possible to have a transform-only part (no facade or physics) in order to
// group other objects in the scene graph.  All subordinate parts will be affected
// by transform changes to parent parts.
type Part interface {
	Id() int // unique id (amongst parts anyways)

	// Dispose, create, and remove parts and sub-parts (bit-parts?).
	Dispose()       // Remove this part, all sub-parts and any associated info.
	AddPart() Part  // Create a new part attached (subordinate) to this one.
	RemPart(p Part) // Remove a subordinate part, without disposing the sub-part.

	// Location, orientation, and scale is the core purpose of a Part.
	Location() (x, y, z float32)    // Get the current location.
	SetLocation(x, y, z float32)    // Set the current location.
	Rotation() (x, y, z, w float32) // Get the current quaternion rotation.
	SetRotation(x, y, z, w float32) // Set the current quaternion rotation.
	Scale() (x, y, z float32)       // Get the size, one value for each axis.
	SetScale(x, y, z float32)       // Set the size, one value for each axis.
	Move(x, y, z float32)           // Move an amount along the current direction.
	RotateX(degrees float32)        // Rotate the part about the X-axis.
	RotateY(degrees float32)        // Rotate the part about the Y-axis.
	RotateZ(degrees float32)        // Rotate the part about the Z-axis.

	// SetFacade associates a facade with a Part.  The facade methods that follow
	// only work when there is a facade associated with the part.
	SetFacade(mesh, shader, material string)     // Create a facade for this part.
	Visible() bool                               // Get the parts visibility.
	SetVisible(visible bool)                     // Set the parts visibility.
	SetCullable(cullable bool)                   // Set the parts cullability.
	Alpha() float32                              // Get the transparency.
	SetAlpha(a float32)                          // Set the transparency.
	SetMaterial(materialName string)             // Set the material.
	SetTexture(textureName string, spin float32) // Set the texture.

	// SetBanner associates a string facade with a Part. The banner methods
	// only work when there is a Banner associated with the part.  Only one of banner
	// or facade may be associated with a part.
	SetBanner(text, shader, glyphs, texture string) // Create a banner for this part.
	UpdateBanner(text string)                       // Change the text for the banner.
	BannerWidth() int                               // Get the banner width in pixels.

	// SetBody associates a Part with physics. The methods grouped with SetBody only
	// work when after SetBody has been called (they are safe to call, but will not
	// do anything).
	SetBody(size, mass float32)         // Create a physics body for this part.
	RemBody()                           // Delete this parts physics body.
	SetShape(shape physics.Shape)       // Only bodies with shape participate in collision.
	LinearMomentum() (x, y, z float32)  // Current linear motion.
	SetLinearMomentum(x, y, z float32)  // Expected to be used to apply initial motion.
	AngularMomentum() (x, y, z float32) // Current angular motion (think of spin).
	SetAngularMomentum(x, y, z float32) // Expected to be used to apply initial motion.
	ResetMomentum()                     // Remove all linear and angular momentum.
	Collide()                           // Run collision on this body.

	// SetResolver allows free bodies to be monitored and have collisions resolved
	// outside of the physics simulation.  Setting a resolver marks the body as
	// "free" meaning that it is not updated by the physics simulation.
	SetResolver(r func(contacts []*physics.Contact))
}

// Part interface
// ===========================================================================
// part - Part implementation

// part implements the Part interface.
type part struct {
	uid      int             // Unique id - only needs to be unique within all parts.
	pov                      // Embed the pov location and orientation struct.
	eng      *Eng            // Injected. Keep a reference for creating new parts.
	res      *roadie         // Injected.
	phy      physics.Physics // Injected.
	staged   bool            // True if the the part is visible.
	scale    *lin.V3         // Scale, per axis: >1 to enlarge, 0<1 to shrink.
	parts    []*part         // Each scene node can have 1 or more scene nodes.
	face     *facade         // Parts can have one optional facade.
	banner   *banner         // Part can have one optional banner instead of a facade.
	body     physics.Body    // Parts can have one motion body.
	visible  bool            // Draw or don't under control of application.
	cullable bool            // Can/can't be culled is under control of the application.
	culled   bool            // Draw or don't under control of engine.
	toc      float64         // Distance to center (to->c) for sorting and culling.
}

// newPart creates and initialzes a part instance.
func newPart(eng *Eng, uid int, res *roadie, phy physics.Physics) *part {
	p := &part{}
	p.uid = uid
	p.eng = eng
	p.res = res
	p.phy = phy
	p.loc = &lin.V3{}
	p.dir = &lin.Q{0, 0, 0, 1}
	p.scale = &lin.V3{1, 1, 1}
	p.parts = []*part{}
	p.visible = true
	p.cullable = true
	p.culled = false
	return p
}

// Part interface implementation.
func (p *part) Dispose() {
	for _, child := range p.parts {
		child.Dispose()
	}
	p.parts = nil
	if p.body != nil {
		p.eng.phy.RemBody(p.body)
	}
}

// Part interface implementation.
func (p *part) AddPart() Part {
	np := newPart(p.eng, p.eng.guid(), p.eng.res, p.eng.phy)
	p.parts = append(p.parts, np)
	copy(p.parts[1:], p.parts[0:])
	p.parts[0] = np
	return np
}

// Part interface implementation.
func (p *part) RemPart(child Part) {
	if pt, _ := child.(*part); pt != nil {
		for index, prt := range p.parts {
			if prt.uid == pt.uid {
				p.parts = append(p.parts[:index], p.parts[index+1:]...)
				return
			}
		}
	}
}

// Part interface implementation.
func (p *part) Id() int                   { return p.uid }
func (p *part) Scale() (x, y, z float32)  { return p.scale.X, p.scale.Y, p.scale.Z }
func (p *part) SetScale(x, y, z float32)  { p.scale.X, p.scale.Y, p.scale.Z = x, y, z }
func (p *part) Visible() bool             { return p.visible }
func (p *part) SetVisible(visible bool)   { p.visible = visible }
func (p *part) SetCullable(cullable bool) { p.cullable = cullable }
func (p *part) SetFacade(mesh, shader, material string) {
	p.face = newFacade(p.uid, p.res, mesh, shader, material)
}

// Part interface implementation.
func (p *part) SetBody(size, mass float32) {
	if p.body != nil {
		p.phy.RemBody(p.body)
	}
	p.body = physics.NewBody(p.uid, size, mass)
	p.body.SetLocation(p.loc)
	p.body.SetRotation(p.dir)
	p.phy.AddBody(p.body)
}

// Part interface implementation.
func (p *part) RemBody() {
	if p.body != nil {
		p.phy.RemBody(p.body)
		p.body = nil
	}
}

// Part interface implementation.
func (p *part) SetShape(shape physics.Shape) {
	if p.body != nil {
		p.body.SetShape(shape)
	}
}

// Part interface implementation.
func (p *part) SetResolver(r func(contacts []*physics.Contact)) {
	if p.body != nil {
		p.body.SetResolver(r)
	}
}

// Part interface implementation.
func (p *part) LinearMomentum() (x, y, z float32) {
	if p.body != nil {
		return p.body.LinearMomentum()
	}
	return 0, 0, 0
}

// Part interface implementation.
func (p *part) SetLinearMomentum(x, y, z float32) {
	if p.body != nil {
		p.body.SetLinearMomentum(x, y, z)
	}
}

// Part interface implementation.
func (p *part) AngularMomentum() (x, y, z float32) {
	if p.body != nil {
		return p.body.AngularMomentum()
	}
	return 0, 0, 0
}

// Part interface implementation.
func (p *part) SetAngularMomentum(x, y, z float32) {
	if p.body != nil {
		p.body.SetAngularMomentum(x, y, z)
	}
}

// Part interface implementation.
func (p *part) ResetMomentum() {
	if p.body != nil {
		p.body.ResetMomentum()
	}
}

func (p *part) Collide() {
	if p.body != nil {
		p.eng.phy.Collide(p.body)
	}
}

// Part interface implementation.
func (p *part) SetMaterial(material string) {
	if p.face != nil {
		p.face.setMaterial(material)
	} else if p.banner != nil {
		p.banner.setMaterial(material)
	}

}

// Part interface implementation.
func (p *part) SetTexture(texture string, rotSpeed float32) {
	if p.face != nil {
		p.face.setTexture(texture, rotSpeed)
	} else if p.banner != nil {
		p.banner.setTexture(texture, rotSpeed)
	}

}

// Part interface implementation.
func (p *part) Alpha() float32 {
	if p.face != nil {
		return p.face.alpha()
	} else if p.banner != nil {
		return p.banner.alpha()
	}

	return 0
}

// Part interface implementation.
func (p *part) SetAlpha(a float32) {
	if p.face != nil {
		p.face.setAlpha(a)
	} else if p.banner != nil {
		p.banner.setAlpha(a)
	}

}

// Part interface implementation.
func (p *part) SetBanner(text, shader, glyphs, texture string) {
	p.banner = newBanner(p.uid, p.res, text, shader, glyphs, texture)
}

// Part interface implementation.
func (p *part) UpdateBanner(text string) {
	if p.banner != nil {
		p.banner.update(text)
	}
}

// Part interface implementation.
func (p *part) BannerWidth() int {
	if p.banner != nil {
		return p.banner.width()
	}
	return 0
}

// model transform must be done in rotation, scale, translate order.
func (p *part) mt() *lin.M4 {
	mt := p.dir.M4()                   // rotation
	mt.ScaleL(p.Scale())               // scale is applied first (on left of rotation)
	return mt.TranslateR(p.Location()) // translate is applied last (on right of rotation).
}

// setDistance calculates and sets the distance of the given part to the
// given center.
func (p *part) setDistance(cenx, ceny, cenz float32) {
	dx := p.loc.X - cenx
	dy := p.loc.Y - ceny
	dz := p.loc.Z - cenz
	p.toc = float64(dx*dx + dy*dy + dz*dz)
}

// outside returns true if the node's distance to center is bigger than the
// given radius.
func (p *part) outside(radius float32) bool { return p.toc > float64(radius*radius) }

// render the part. This takes the parts rendering specific information and passes
// it into the rendering component.
func (p *part) render(gc render.Renderer, sc *scene, parentTransform *lin.M4) {
	if p.visible {
		m := p.mt().Mult(parentTransform) // model transform + parent transform

		// only render nodes with facades or banners.
		// transfer the rendering information in a graphics structure.
		if p.face != nil || p.banner != nil {
			vis := &render.Visible{}
			v := sc.vt()                        // view transform
			vis.Mv = m.Mult(v)                  // model-view transform
			vis.Mvp = vis.Mv.Clone().Mult(sc.P) // model-view-projection transform
			vis.L = sc.L

			// both banners and facades render the same way. Should only be one
			// specified, but prefer a facade over a banner in the case of dev error.
			var face *facade
			if p.face != nil {
				face = p.face
			} else {
				face = &p.banner.facade
			}
			vis.Mesh = face.msh
			vis.Shader = face.shadr
			vis.Mat = face.mat
			vis.Texture = face.tex
			vis.RotSpeed = face.rots
			vis.Scale = p.scale.X

			// Use a large fade default for scenes without radius.
			vis.Fade = 1000
			if sc.radius != 0 {
				vis.Fade = sc.radius
			}
			gc.Render(vis)
		}

		// render all the parts children
		for _, child := range p.parts {
			if !child.culled {
				child.render(gc, sc, m.Clone()) // don't change the parent model transform.
			}
		}
	}
}

// part
// ===========================================================================
// Parts

// Parts is used to sort a slice of parts in order to get transparency working.
// The furthest away things have to be drawn first.
// This is only public for the sort package and is not for use by the application.
type Parts []*part

// Implement the sort interface allowing parts to ordered by the
// value of the t0 field.
func (p Parts) Len() int           { return len(p) }
func (p Parts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Parts) Less(i, j int) bool { return p[i].toc > p[j].toc }
