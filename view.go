// Copyright © 2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// View dictates how models are rendered. A view is attached to a Pov
// where it renders all models in that Pov's hierarchy. It is possible
// to attach multiple views, even to the same point. All the models
// in each visible view are rendered each update.
type View interface {
	Cam() Camera             // A single camera for a group of Models.
	Visible() bool           // Only visible views are rendered.
	SetVisible(visible bool) // Whether or not the view is rendered.
	SetDepth(enabled bool)   // True for 3D views. 2D views ignore depth.
	SetUI()                  // UI view: 2D, no depth, drawn last.
	SetLast(index int)       // For sequencing UI views. Higher is later.

	// SetCull sets a method that reduces the number of Models rendered
	// each update. It can be engine supplied ie: NewFacingCuller,
	// or application supplied.
	SetCull(c Cull) // Set to nil to turn off culling.
}

// View
// =============================================================================
// view

// view implements View.
type view struct {
	cam     *camera // Camera created during initialization.
	depth   bool    // True for 3D depth processing.
	visible bool    // Is the scene drawn or not.
	cull    Cull    // Set by application.
	overlay int     // Set render bucket with OVERLAY or greater.
}

// newView creates a new structure each time.
func newView() *view {
	v := &view{depth: true}
	v.cam = newCamera()
	v.visible = true
	return v
}

// Implement View interface.
func (v *view) Cam() Camera             { return v.cam }
func (v *view) Visible() bool           { return v.visible }
func (v *view) SetVisible(visible bool) { v.visible = visible }
func (v *view) SetDepth(enabled bool)   { v.depth = enabled }
func (v *view) SetCull(c Cull)          { v.cull = c }
func (v *view) SetLast(index int)       { v.overlay = OVERLAY + index }
func (v *view) SetUI() {
	v.overlay = OVERLAY    // Draw last.
	v.depth = false        // 2D rendering.
	v.cam.SetTransform(VO) // orthographic transform.
}
