// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package vu

import (
	"sort"
	"vu/math/lin"
)

// Scene orgnanizes parts and is associated with a single camera and camera
// transform. There may be more than one Scene instance to allow for groups
// of objects, like overlays, to use different camera transforms.
type Scene interface {
	AddPart() Part                    // Add a scene node.
	RemPart(p Part)                   // Remove a scene node.
	SetTransform(transform int)       // Set a view transform.
	LightLocation() (x, y, z float64) // Currently one light per scene,
	SetLightLocation(x, y, z float64) // ...which has a location,
	SetLightColour(r, g, b float64)   // ...and a colour.
	Set2D()                           // Overlays are 2D.
	Visible() bool                    // Only visible scenes are rendered.
	SetVisible(visible bool)          // Whether or not the scene is rendered.

	// One camera is associated with a scene.
	Location() (x, y, z float64)    // Get, or
	SetLocation(x, y, z float64)    // ...Set the camera location.
	Rotation() (x, y, z, w float64) // Get, or
	SetRotation(x, y, z, w float64) // ...Set the view orientation.
	Tilt() (up float64)             // Get, or
	SetTilt(up float64)             // ...Set the camera tilt angle.
	Move(x, y, z float64)           // Adjust current camera location.
	Spin(dir int, degrees float64)  // Rotate the current camera.

	// SetSorted orders the parts based on distance so that transparency
	// works. Objects furthest away are drawn first with transparency.
	SetSorted(sorted bool)

	// SetCuller is used to reduce the number of Parts rendered.
	SetCuller(c Culler) // Set to nil to turn off culling.

	// Use one of the following to create the scenes projection transforms.
	SetPerspective(fov, ratio, near, far float64)                // 3D.
	SetOrthographic(left, right, bottom, top, near, far float64) // 2D.
}

// FUTURE: implement a dirty flag to only update transforms for
//         parts that have changed position/rotation.
// FUTURE: Move away from combining transforms as matricies? Ie.
//         http://www.euclideanspace.com/maths/geometry/affine/nonMatrix/index.htm

// Scene interface
// =============================================================================
// scene - Scene implementation

// scene implements Scene.
type scene struct {
	parts   []*part  // Each scene node can have 1 or more child nodes.
	l       *light   // Scene lighting
	cam     *view    // Camera created during initialization.
	sorted  bool     // True if objects are sorted based on distance.
	is2D    bool     // Whether the scene is 2D or 3D
	visible bool     // Is the scene drawn or not.
	cull    Culler   // Set by application.
	feed    feedback // Injected on creation and passed on to new parts.
	assets  *assets  // Injected on creation and passed on to new parts.
	vm      *lin.M4  // Scratch: View part of MVP matrix.
	pm      *lin.M4  // Scratch: Projection part of MVP matrix.
}

// newScene creates and initializes a scene struct.
func newScene(transform int, f feedback, a *assets) *scene {
	s := &scene{}
	s.l = &light{}
	s.cam = newView()
	s.SetTransform(transform)
	s.parts = []*part{}
	s.visible = true
	s.feed = f
	s.assets = a
	s.vm = &lin.M4{}
	s.pm = &lin.M4{}
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
func (s *scene) LightLocation() (x, y, z float64) {
	return float64(s.l.x), float64(s.l.y), float64(s.l.z)
}
func (s *scene) SetLightLocation(x, y, z float64) {
	s.l.x, s.l.y, s.l.z = float32(x), float32(y), float32(z)
}
func (s *scene) SetLightColour(r, g, b float64) {
	s.l.ld.R, s.l.ld.G, s.l.ld.B = float32(r), float32(g), float32(b)
}
func (s *scene) Cam() *view                     { return s.cam }
func (s *scene) SetSorted(sorted bool)          { s.sorted = sorted }
func (s *scene) SetCuller(c Culler)             { s.cull = c }
func (s *scene) Tilt() (up float64)             { return s.cam.up }
func (s *scene) SetTilt(up float64)             { s.cam.up = up }
func (s *scene) Rotation() (x, y, z, w float64) { return s.cam.Rotation() }
func (s *scene) SetRotation(x, y, z, w float64) { s.cam.SetRotation(x, y, z, w) }
func (s *scene) Location() (x, y, z float64)    { return s.cam.Location() }
func (s *scene) SetLocation(x, y, z float64)    { s.cam.SetLocation(x, y, z) }
func (s *scene) Move(x, y, z float64)           { s.cam.Move(x, y, z) }
func (s *scene) Spin(dir int, degrees float64) {
	switch dir {
	case XAxis:
		s.cam.Spin(degrees, 0, 0)
	case YAxis:
		s.cam.Spin(0, degrees, 0)
	case ZAxis:
		s.cam.Spin(0, 0, degrees)
	}
}

// Scene interface implementation.
func (s *scene) AddPart() Part {
	np := newPart(s.feed, s.assets)
	s.parts = append(s.parts, np)
	return np
}

// Scene interface implementation.
// Find and remove the part (will point to the same record).
// The part is removed from the scene, but not disposed.
func (s *scene) RemPart(scenePart Part) {
	if pt, _ := scenePart.(*part); pt != nil {
		for index, p := range s.parts {
			if p == pt {
				s.parts = append(s.parts[:index], s.parts[index+1:]...)
				return
			}
		}
	}
}

// SetPerspective sets the scene to use a 3D perspective
func (s *scene) SetPerspective(fov, ratio, near, far float64) {
	s.pm = lin.NewM4().Persp(fov, ratio, near, far)
}

// SetOrthographic sets the scene to use a 2D orthographic perspective.
func (s *scene) SetOrthographic(left, right, bottom, top, near, far float64) {
	s.pm = lin.NewM4().Ortho(left, right, bottom, top, near, far)
}

// vt applies the view transform to the scene camera and returns the result
// in the supplied matrix.
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
				childPart.culled = s.cull.Cull(s, childPart)
			}
			s.cullParts(childPart.parts)
		}
	}
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

// stage lets the scene dictate what is visible and in what order.
// Currently the parts are expected to have been sorted before calling stage.
func (s *scene) stage(dt float64) {
	s.vt(s.vm) // update sc.vm with the latest view matrix (camera transform).
	for _, part := range s.parts {
		if !part.culled {
			part.stage(s, dt)
		}
	}
}

// verify passes the check request down to the parts.
func (s *scene) verify() error {
	for _, p := range s.parts {
		if err := p.verify(); err != nil {
			return err
		}
	}
	return nil
}
