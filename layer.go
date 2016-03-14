// Copyright © 2015-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"github.com/gazed/vu/math/lin"
)

// FUTURE: Lots more to do for render targets/passes/layers.
//         Need a clean design *and* need to justify the additionaly code
//         complexity against the render output benefits.

// Layer is used to render to a 1024x1024 sized frame buffer based texture.
// A layer represents the output of an extra render pass where objects drawn
// to this off screen texture are used as input for a later render pass.
type Layer interface{}

// Layer
// =============================================================================
// layer implements Layer.

// layer implements Layer.
// Note that the mvp value must be set by the Pov that is taking the picture.
// The values are copied in as the Pov is rendered, and then used later if
// the layer is being used as a shadow map.
type layer struct {
	bid  uint32   // Framebuffer id. Default 0 for default framebuffer.
	db   uint32   // Depth renderbuffer.
	attr int      // What type of layer. Full IMAGE or SHADOW_MAP.
	vp   *lin.M4  // light view-projection layer transform.
	bm   *lin.M4  // bias matrix.
	tex  *texture // place holder for rendered texture. Created on GPU.
}

// newLayer creates the framebuffer needed to render to a texture.
func newLayer(attr int) *layer {
	l := &layer{attr: attr}
	l.vp = &lin.M4{}
	l.bm = &lin.M4{
		Xx: 0.5, Xy: 0.0, Xz: 0.0, Xw: 0.0,
		Yx: 0.0, Yy: 0.5, Yz: 0.0, Yw: 0.0,
		Zx: 0.0, Zy: 0.0, Zz: 0.5, Zw: 0.0,
		Wx: 0.5, Wy: 0.5, Wz: 0.5, Ww: 1.0,
	}
	l.tex = newTexture("rendered")
	l.tex.bound, l.tex.loaded = true, true
	return l
}
