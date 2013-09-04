// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package vu

// scene is a "scene graph" that attempts to organize, optimize, and
// render the model objects (Parts) that are created by the application.

import (
	"sort"
	"vu/math/lin"
	"vu/render"
)

// Scene orgnanizes everything that needs to be rendered. It is the top level
// scene node in a scene graph.  There may be more than one Scene instance to
// allow for groups of objects like overlays to use different camera transforms.
type Scene interface {
	AddPart() Part                    // Add a scene node.
	RemPart(p Part)                   // Remove a scene node.
	SetTransform(transform int)       // Any one of a number of view transforms.
	SetLightLocation(x, y, z float32) // Currently only one light.
	SetLightColour(r, g, b float32)   // Currently only one light.
	Set2D()                           // Overlays are 2D.
	Visible() bool                    // Only visible scenes are rendered.
	SetVisible(visible bool)          // Change whether or not the scene is rendered.

	// Expose the camera positioning methods through the scene.
	ViewLocation() (x, y, z float32)         // Get the view location.
	SetViewLocation(x, y, z float32)         // Set the camera position.
	ViewRotation() (x, y, z, w float32)      // Get the view orientation.
	SetViewRotation(x, y, z, w float32)      // Set the view orientation.
	PreviousViewLocation() (x, y, z float32) // Get the previous view location.
	MoveView(x, y, z float32)                // Alter camera position based direction.
	PanView(dir int, degrees float32)        // Rotate the camera.
	ViewTilt() (up float32)                  // Get the camera height.
	SetViewTilt(up float32)                  // Alter the camera height.

	// SetSorted is necessary to order the parts based on distance so that transparency
	// works. The furthest away objects have be drawn first for transparency.
	SetSorted(sorted bool)

	// SetVisibleRadius limits the visible (rendered) objects in the scene to a
	// given amount from the scenes camera.  This is used for minimap type scenes
	// where what needs to be drawn is a given amount around the camera, which is
	// at the center of the visible area.  Setting this to 0 or less will turn off
	// the visible radius and everything will be drawn.
	SetVisibleRadius(radius float32) // Restrict rendering to a circle around the camera.

	// SetVisibleDirection is intended to be used with SetVisibleRadius such that
	// the visible area is moved radius units in the camera direction. This is used
	// for a 3D persons view where the view is in front of them, i.e. the camera
	// is at the edge and facing the visible area.
	SetVisibleDirection(facing bool) // Move the rendering circle in front of the camera.

	SetPerspective(fov, ratio, near, far float32)                // The projection transform.
	SetOrthographic(left, right, bottom, top, near, far float32) // The projection transform.
}

// Scene interface
// ===========================================================================
// scene - Scene implementation

// scene implements Scene.
type scene struct {
	uid     int           // Scene unique id used for comparison and lookup
	eng     *Eng          // Need access to engine subsystems to create new parts.
	parts   []*part       // Each scene node can have 1 or more child nodes.
	P       *lin.M4       // Projection part of MVP matrix.
	L       *render.Light // Scene lighting
	cam     *view         // Camera created during initialization.
	radius  float32       // Set to >0 when objects are to culled by distance.
	facing  bool          // True when objects are to be culled by distance and direction.
	sorted  bool          // True if objects are sorted based on distance.
	is2D    bool          // Whether the scene is 2D or 3D
	visible bool          // Is the scene drawn or not.
}

// newScene creates and initializes a scene struct.
func newScene(eng *Eng, uid, transform int) *scene {
	s := &scene{}
	s.uid = uid
	s.eng = eng
	s.L = &render.Light{}
	s.P = &lin.M4{}
	s.cam = newView()
	s.SetTransform(transform)
	s.parts = []*part{}
	s.visible = true
	return s
}

// dispose of the scene by disposing all child nodes.
func (s *scene) dispose() {
	for _, part := range s.parts {
		part.Dispose()
	}
}

// Scene interface implementation.
func (s *scene) SetTransform(transform int)              { s.cam.vt = getViewTransform(transform) }
func (s *scene) Set2D()                                  { s.is2D = true }
func (s *scene) Visible() bool                           { return s.visible }
func (s *scene) SetVisible(visible bool)                 { s.visible = visible }
func (s *scene) SetLightLocation(x, y, z float32)        { s.L.X, s.L.Y, s.L.Z = x, y, z }
func (s *scene) SetLightColour(r, g, b float32)          { s.L.Ld.R, s.L.Ld.G, s.L.Ld.B = r, g, b }
func (s *scene) Cam() *view                              { return s.cam }
func (s *scene) SetVisibleDirection(on bool)             { s.facing = on }
func (s *scene) SetSorted(sorted bool)                   { s.sorted = sorted }
func (s *scene) ViewTilt() (up float32)                  { return s.cam.up }
func (s *scene) SetViewTilt(up float32)                  { s.cam.up = up }
func (s *scene) ViewRotation() (x, y, z, w float32)      { return s.cam.Rotation() }
func (s *scene) SetViewRotation(x, y, z, w float32)      { s.cam.SetRotation(x, y, z, w) }
func (s *scene) ViewLocation() (x, y, z float32)         { return s.cam.Location() }
func (s *scene) SetViewLocation(x, y, z float32)         { s.cam.SetLocation(x, y, z) }
func (s *scene) PreviousViewLocation() (x, y, z float32) { return s.cam.PreviousLocation() }
func (s *scene) MoveView(x, y, z float32)                { s.cam.Move(x, y, z) }
func (s *scene) PanView(dir int, degrees float32) {
	switch dir {
	case XAxis:
		s.cam.RotateX(degrees)
	case YAxis:
		s.cam.RotateY(degrees)
	case ZAxis:
		s.cam.RotateZ(degrees)
	}
}

// Part interface implementation.
func (s *scene) AddPart() Part {
	np := newPart(s.eng, s.eng.guid(), s.eng.res, s.eng.phy)
	s.parts = append(s.parts, np)
	return np
}

// Part interface implementation.
func (s *scene) RemPart(child Part) {
	if pt, _ := child.(*part); pt != nil {
		for index, p := range s.parts {
			if p.uid == pt.uid {
				s.parts = append(s.parts[:index], s.parts[index+1:]...)
				return
			}
		}
	}
}

// Scene interface implementation.
func (s *scene) SetVisibleRadius(radius float32) {
	s.radius = 0
	if radius > 0 {
		s.radius = radius
	}
}

// SetPerspective sets the scene to use a 3D perspective
func (s *scene) SetPerspective(fov, ratio, near, far float32) {
	s.P = lin.M4Perspective(fov, ratio, near, far)
}

// SetOrthographic sets the scene to use a 2D orthographic perspective.
func (s *scene) SetOrthographic(left, right, bottom, top, near, far float32) {
	s.P = lin.M4Orthographic(left, right, bottom, top, near, far)
}

// vt applies the view transform to the scene camera and returns the result.
func (s *scene) vt() *lin.M4 { return s.cam.vt(s.cam) }

// cullParts sets the parts distance and marks it culled or not.
func (s *scene) cullParts(parts Parts) {
	if parts != nil && len(parts) > 0 {
		for _, childPart := range parts {
			childPart.setDistance(s.cam.Location())
			if childPart.cullable {
				childPart.culled = s.cullPart(childPart)
			}
			s.cullParts(childPart.parts)
		}
	}
}

// cullPart based on distance if the scenes culling variables are set
func (s *scene) cullPart(p *part) bool {
	if s.radius > 0 {

		// base the cull on an area in front of the camera by moving the
		// center up radius units in facing direction.  Don't move it all the way
		// up so that stuff above or below still exists when looking up/down.
		if s.facing {
			toc := p.toc // save the old location.

			// project the camera location along the lookat vector.
			fudgeFactor := float32(0.8) // don't move all the way up.
			updownRot := lin.QAxisAngle(&lin.V3{1, 0, 0}, -s.cam.up)
			lookAt := s.cam.dir.Clone().Mult(updownRot)
			dist := (&lin.V3{0, 0, -s.radius * fudgeFactor}).MultQ(lookAt)
			cen := s.cam.loc.Clone().Add(dist)

			// cull the part if its to far away.
			p.setDistance(cen.X, cen.Y, cen.Z)
			if p.outside(s.radius) {
				p.toc = toc
				return true
			}
			p.toc = toc // restore the old location.
		} else if p.outside(s.radius) {
			return true
		}
	}
	return false
}

// sortParts reorganizes the scene graph based on the distance to the camera.
// This is used for transparent objects so that the further away objects are
// drawn first.
func (s *scene) sortParts(parts Parts) {
	if parts != nil && len(parts) > 0 {
		for _, childPart := range parts {
			s.sortParts(childPart.parts)
		}
		sort.Sort(parts)
	}
}

// render lets the scene dictate what gets sent to the graphics card.
func (s *scene) render(gc render.Renderer) {
	gc.Enable(render.DEPTH, s.is2D)
	parentTransform := lin.M4Identity()
	for _, part := range s.parts {
		if !part.culled {
			part.render(gc, s, parentTransform)
		}
	}
}
