// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package vu

import (
	"log"
	"vu/panel"
)

// overlay links the 2D panel package to the 3D rendering engine.
// Overlay associates a rendering part for each control in the panel.
// All the rendered parts are tracked within a single Scene.
type overlay struct {
	panel  panel.Panel    // The control panel this overlay is based on.
	parts  map[uint]*part // Part map for the panel controls.
	scene  *scene         // Draw the parts in their own scene.
	w, h   int            // Current overlay size.
	staged []*part        // Scratch slice: the visible/rendered parts.
}

// newOverlay expects to be initialized with a control panel.
func newOverlay(scene *scene, p panel.Panel) *overlay {
	o := &overlay{}
	o.scene = scene
	o.panel = p
	o.w, o.h = p.Size()
	o.scene.SetOrthographic(0, float64(o.w), 0, float64(o.h), 0, 10)
	o.scene.vt(o.scene.vm) // calculate view matrix once at start.
	o.scene.Set2D()
	o.parts = map[uint]*part{}
	o.createParts(o.panel)
	for _, child := range p.Dialogs() {
		o.createParts(child)
	}
	o.staged = []*part{}
	return o
}

// dispose ensures that the GPU resources are properly cleaned up.
func (o *overlay) dispose() {
	if o.parts != nil {
		for key, p := range o.parts {
			p.Dispose()
			delete(o.parts, key)
		}
	}
	o.panel = nil
	o.parts = nil
}

// createParts associates a graphic Part with each widget. It is expected
// to be called once on creation of a widgetPainter. Widgets with images are
// created as Facades while widgets without images are flat shaded meshes.
func (o *overlay) createParts(w panel.Widget) {
	p := newPart(o.scene.feed, o.scene.assets)
	switch {
	case w.Img() != "":
		p.SetRole("widget").SetMesh("icon").AddTex(w.Img())
	default:
		p.SetRole("flat").SetMesh("icon")
	}
	o.parts[w.Id()] = p
	if pw, ok := w.(panel.Section); ok {
		for _, child := range pw.Widgets() {
			o.createParts(child)
		}
	}
}

// update is called to refresh the visible panel parts. The top level ortho
// view is constantly adjusted to be identical to the size of the overlay...
// which is expected to match the current screen dimensions in pixels.
// The list of visible parts to be rendered is returned in what is intended
// as a read only list.
func (o *overlay) update() []*part {
	o.staged = o.staged[:0] // reset, keeping underlying memory.
	o.scene.SetVisible(o.panel.Visible())
	if o.scene.Visible() {
		w, h := o.panel.Size()
		if o.w != w || o.h != h {
			o.w, o.h = w, h
			o.scene.SetOrthographic(0, float64(w), 0, float64(h), 0, 10)
		}
		o.updateSection(o.panel, o.panel.Focus())

		// Draw active dialogs over the overlay.
		if dia := o.panel.Dialog(); dia != nil {
			o.updateSection(dia, o.panel.Focus())
		}
	}
	return o.staged
}

// updateSection recursively calls each panel in the hierarchy.
func (o *overlay) updateSection(p panel.Section, focus uint) {
	o.updateWidget(p, focus)

	// Draw any children on top of what's been just been drawn.
	for _, w := range p.Widgets() {
		if pw, ok := w.(panel.Section); ok {
			o.updateSection(pw, focus)
		}
		o.updateWidget(w, focus)
	}
}

// updateWidget updates then render state for the given wiget.
// It auto highlights controls that the mouse is hovering over.
func (o *overlay) updateWidget(w panel.Widget, focus uint) {
	var p *part
	var ok bool
	if p, ok = o.parts[w.Id()]; !ok {
		log.Printf("No part for %d", w.Id())
		return
	}
	p.SetVisible(w.Visible())
	if p.Visible() {
		o.updateImage(p, w, focus)
		cx, cy, hx, hy := w.At()
		depth := -10 + float64(w.Id())*0.1
		p.SetLocation(cx, cy, depth)
		p.SetScale(hx, hy, 0)
		p.stage(o.scene, 0)
		o.staged = append(o.staged, p)
	}
}

// updateImage renders the image for a visible part.
// Only controls can have focus.
func (o *overlay) updateImage(p *part, w panel.Widget, focus uint) {
	p.Role().SetAlpha(0.0)  // hide the widget by default.
	p.Role().SetKd(1, 1, 1) // default white background.
	con, isControl := w.(panel.Control)
	haveFocus := focus == w.Id()
	haveImg := w.Img() != ""
	switch {

	case !haveImg && haveFocus:
		// no image, so just highlight controls that have focus.
		p.Role().SetAlpha(0.25)
		p.Role().SetKd(0.4, 0.4, 0.4)

	case haveImg && isControl && con.ImgOnHover() && !haveFocus:
		// hover images need to be hidden when not in focus

	case haveImg && haveFocus && !con.ImgOnHover():
		p.Role().SetAlpha(1.0)
		p.Role().SetKd(0.9, 0.9, 0.9)

	case haveImg:
		// by default show a section or control with an image.
		p.Role().SetAlpha(1.0)

	default:
		// default is to not show the image
	}

	// check if the image needs updating.
	if haveImg {
		if w.Img() != p.Role().Tex(0) {
			p.Role().UseTex(w.Img(), 0)
		}
	}
}
