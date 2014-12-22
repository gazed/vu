// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// Scene associates a group of parts with a single camera and projection
// transform. For example an overlay often uses an unmoving Orthographic
// camera transform while a main 3D scene uses a moving Perspective camera
// transform. Parts are rendered in the order created unless changed by
// SetLast.
type Scene interface {
	Cam() Camera    // A single camera for this group of Parts.
	AddPart() Part  // Add a root part to the scene.
	RemPart(p Part) // Remove a root part.

	// Visible shows or hides all the scene parts.
	Visible() bool           // Only visible scenes are rendered.
	SetVisible(visible bool) // Whether or not the scene is rendered.
	Set2D()                  // 2D turns off Z-buffer for all parts.

	// SetSorted orders the parts based on distance. Required for transparency.
	SetSorted(sorted bool) // Objects furthest away are drawn first.
	SetLast(p Part)        // Make the part last in the rendered part list.

	// SetCuller is used to reduce the number of Parts rendered.
	// Eg: NewRadiusCuller, NewFacingCuller, or application supplied.
	SetCuller(c Culler) // Set to nil to turn off culling.
}

// Scene interface
// =============================================================================
// scene - Scene implementation

// scene implements Scene. Mostly scene groups parts and acts as a data holder
// for the scene graph operations in stage.
type scene struct {
	sm      *stage  // stage manager is injected on creation.
	parts   []*part // Each scene node can have 1 or more child nodes.
	cam     *camera // Camera created during initialization.
	sorted  bool    // True if objects are to be sorted by distance.
	is2D    bool    // Whether the scene is 2D or 3D
	cull    Culler  // Set by application.
	visible bool    // Is the scene drawn or not.
}

// newScene creates and initializes a scene struct.
func newScene(transform int, sm *stage) *scene {
	s := &scene{}
	s.cam = newCamera()
	s.cam.SetTransform(transform)
	s.parts = []*part{}
	s.visible = true
	s.sm = sm
	return s
}

// dispose of the scene by disposing all child nodes.
func (s *scene) dispose() {
	for _, part := range s.parts {
		part.Dispose()
	}
	s.parts = []*part{} // release references.
}

// Scene interface implementation.
func (s *scene) Set2D()                  { s.is2D = true }
func (s *scene) Visible() bool           { return s.visible }
func (s *scene) SetVisible(visible bool) { s.visible = visible }
func (s *scene) Cam() Camera             { return s.cam }
func (s *scene) SetSorted(sorted bool)   { s.sorted = sorted }
func (s *scene) SetCuller(c Culler)      { s.cull = c }

// Scene interface implementation.
func (s *scene) AddPart() Part {
	np := s.sm.addPart()
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

// Scene interface implementation.
func (s *scene) SetLast(p Part) {
	part := p.(*part)
	for index, pt := range s.parts {
		if pt == part {
			s.parts = append(s.parts[:index], s.parts[index+1:]...)
			s.parts = append(s.parts, part)
		}
	}
}
