// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package vu

// scene is a "scene graph" that attempts to organize, optimize, and
// render the model objects (Parts) that are created by the application.

import (
	"sort"
	"vu/math/lin"
	"vu/move"
	"vu/render"
)

// Scene orgnanizes everything that needs to be rendered and is associated with
// a single camera and camera transform. It is also the top level scene node in
// a scene graph. There may be more than one Scene instance to allow for groups
// of objects, like overlays, to use different camera transforms.
type Scene interface {
	AddPart() Part                    // Add a scene node.
	RemPart(p Part)                   // Remove a scene node.
	SetTransform(transform int)       // Set a view transform.
	SetLightLocation(x, y, z float64) // Currently only one light.
	SetLightColour(r, g, b float64)   // Currently only one light.
	Set2D()                           // Overlays are 2D.
	Visible() bool                    // Only visible scenes are rendered.
	SetVisible(visible bool)          // Change whether or not the scene is rendered.

	// Camera.
	ViewLocation() (x, y, z float64)    // Get the view location.
	SetViewLocation(x, y, z float64)    // Set the camera position.
	ViewRotation() (x, y, z, w float64) // Get the view orientation.
	SetViewRotation(x, y, z, w float64) // Set the view orientation.
	MoveView(x, y, z float64)           // Alter camera position based direction.
	PanView(dir int, degrees float64)   // Rotate the camera.
	ViewTilt() (up float64)             // Get the camera tilt angle.
	SetViewTilt(up float64)             // Alter the camera tilt angle.

	// SetSorted orders the parts based on distance so that transparency
	// works. Objects furthest away are drawn first with transparency.
	SetSorted(sorted bool)

	// SetVisibleRadius limits the visible (rendered) objects in the scene to a
	// given radius from the scenes camera. Setting this to 0 or less will turn
	// off the visible radius and everything will be drawn.
	SetVisibleRadius(radius float64) // Restrict rendering to a circle around the camera.

	// SetVisibleDirection is used with SetVisibleRadius such that the visible
	// area is moved radius units in the camera's direction, i.e. the camera is
	// at the edge and facing the visible area.
	SetVisibleDirection(facing bool) // Move the rendering circle in front of the camera.

	// Create the perspective transform for this scene.
	SetPerspective(fov, ratio, near, far float64)                // 3D projection transform.
	SetOrthographic(left, right, bottom, top, near, far float64) // 2D projection transform.
}

// Scene interface
// ===========================================================================
// scene - Scene implementation

// scene implements Scene.
type scene struct {
	parts   []*part       // Each scene node can have 1 or more child nodes.
	P       *lin.M4       // Projection part of MVP matrix.
	L       *render.Light // Scene lighting
	cam     *view         // Camera created during initialization.
	radius  float64       // Set to >0 when objects are to culled by distance.
	facing  bool          // True when objects are to be culled by distance and direction.
	sorted  bool          // True if objects are sorted based on distance.
	is2D    bool          // Whether the scene is 2D or 3D
	visible bool          // Is the scene drawn or not.
	world   move.World    // Injected on creation and passed on to parts.
}

// newScene creates and initializes a scene struct.
func newScene(transform int, world move.World) *scene {
	s := &scene{}
	s.L = &render.Light{}
	s.P = &lin.M4{}
	s.cam = newView()
	s.SetTransform(transform)
	s.parts = []*part{}
	s.visible = true
	s.world = world
	return s
}

// dispose of the scene by disposing all child nodes.
func (s *scene) dispose() {
	for _, part := range s.parts {
		part.Dispose()
	}
}

// Scene interface implementation.
func (s *scene) SetTransform(transform int) { s.cam.vt = getViewTransform(transform) }
func (s *scene) Set2D()                     { s.is2D = true }
func (s *scene) Visible() bool              { return s.visible }
func (s *scene) SetVisible(visible bool)    { s.visible = visible }
func (s *scene) SetLightLocation(x, y, z float64) {
	s.L.X, s.L.Y, s.L.Z = float32(x), float32(y), float32(z)
}
func (s *scene) SetLightColour(r, g, b float64) {
	s.L.Ld.R, s.L.Ld.G, s.L.Ld.B = float32(r), float32(g), float32(b)
}
func (s *scene) Cam() *view                         { return s.cam }
func (s *scene) SetVisibleDirection(on bool)        { s.facing = on }
func (s *scene) SetSorted(sorted bool)              { s.sorted = sorted }
func (s *scene) ViewTilt() (up float64)             { return s.cam.up }
func (s *scene) SetViewTilt(up float64)             { s.cam.up = up }
func (s *scene) ViewRotation() (x, y, z, w float64) { return s.cam.Rotation() }
func (s *scene) SetViewRotation(x, y, z, w float64) { s.cam.SetRotation(x, y, z, w) }
func (s *scene) ViewLocation() (x, y, z float64)    { return s.cam.Location() }
func (s *scene) SetViewLocation(x, y, z float64)    { s.cam.SetLocation(x, y, z) }
func (s *scene) MoveView(x, y, z float64)           { s.cam.Move(x, y, z) }
func (s *scene) PanView(dir int, degrees float64) {
	switch dir {
	case XAxis:
		s.cam.Spin(degrees, 0, 0)
	case YAxis:
		s.cam.Spin(0, degrees, 0)
	case ZAxis:
		s.cam.Spin(0, 0, degrees)
	}
}

// Part interface implementation.
func (s *scene) AddPart() Part {
	np := newPart(s.world)
	s.parts = append(s.parts, np)
	return np
}

// Part interface implementation.
// Find and remove the part (will point to the same record).
// The part is removed from the scene, but not disposed.
func (s *scene) RemPart(child Part) {
	if pt, _ := child.(*part); pt != nil {
		for index, p := range s.parts {
			if p == pt {
				s.parts = append(s.parts[:index], s.parts[index+1:]...)
				return
			}
		}
	}
}

// Scene interface implementation.
func (s *scene) SetVisibleRadius(radius float64) {
	s.radius = 0
	if radius > 0 {
		s.radius = radius
	}
}

// SetPerspective sets the scene to use a 3D perspective
func (s *scene) SetPerspective(fov, ratio, near, far float64) {
	s.P = lin.NewPersp(fov, ratio, near, far)
}

// SetOrthographic sets the scene to use a 2D orthographic perspective.
func (s *scene) SetOrthographic(left, right, bottom, top, near, far float64) {
	s.P = lin.NewOrtho(left, right, bottom, top, near, far)
}

// vt applies the view transform to the scene camera and returns the result.
func (s *scene) vt(vm *lin.M4) *lin.M4 { return s.cam.vt(s.cam, vm) }

// setDistances sets the distance of the part from the camera.  The distance
// is needed for transparency ordering and for culling.
func (s *scene) setDistances(parts Parts) {
	if parts != nil && len(parts) > 0 {
		for _, childPart := range parts {
			childPart.toc = childPart.distanceTo(s.cam.Location())
			s.setDistances(childPart.parts)
		}
	}
}

// cullParts sets the parts distance and marks it culled or not.
func (s *scene) cullParts(parts Parts) {
	if parts != nil && len(parts) > 0 {
		for _, childPart := range parts {
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
			fudgeFactor := float64(0.8) // don't move all the way up.
			lookAt := lin.NewQ().SetAa(1, 0, 0, lin.Rad(s.cam.up))
			lookAt.Mult(s.cam.dir, lookAt)
			cx, cy, cz := lin.MultSQ(0, 0, -s.radius*fudgeFactor, lookAt) // distance
			cx, cy, cz = cx+s.cam.loc.X, cy+s.cam.loc.Y, cz+s.cam.loc.Z   // added to location.

			// cull the part if its to far away.
			p.toc = p.distanceTo(cx, cy, cz)
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
// drawn first.  The parts distance value is expected to have been set.
func (s *scene) sortParts(parts Parts) {
	if parts != nil && len(parts) > 0 {
		for _, childPart := range parts {
			s.sortParts(childPart.parts)
		}
		sort.Sort(parts)
	}
}

// stage lets the scene dictate what is visible and in what order.
// Currently the parts are expected to have been sorted before calling stage.
func (s *scene) stage(vis *[]*render.Vis) {
	parentTransform := lin.NewM4I()
	for _, part := range s.parts {
		if !part.culled {
			part.stage(vis, s, parentTransform)
		}
	}
}
