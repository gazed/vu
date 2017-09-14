// Copyright © 2015-2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// frame.go contains utility methods to help convert visible objects
// into draw calls.

import (
	"math"

	"github.com/gazed/vu/render"
)

// frame links the engine models and cameras to a render frame.
// It generates an ordered list of render system draw data needed for
// one screen. There is always one frame that is being updated while
// the two previous frames are being rendered.
type frame []*render.Draw

// getDraw returns a render.Draw. The frame is grown as needed and draw
// instances are reused if available. Every frame value up to cap(frame)
// is expected to have already been allocated.
func (f frame) getDraw() (frame, **render.Draw) {
	size := len(f)
	switch {
	case size == cap(f):
		f = append(f, render.NewDraw())
	case size < cap(f): // use previously allocated.
		f = f[:size+1]
		if f[size] == nil {
			f[size] = render.NewDraw()
		} else {
			// clean up old data.
			f[size].SetPose(nil)             // Clear animation.
			f[size].SetTex(0, 0, 0, 0, 0, 0) // Clear texture info.
		}
	}
	return f, &f[size]
}

// frame
// =============================================================================
// sorting draw calls.

// setBucket produces a number that is used to order draw calls.
// Higher values are rendered before lower values.
//    xxxxxxxx Pass.... Over.... xxxxxxxA Distance to Camera.................
//    00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000
// Bits marked x are currently unused.
// Pass is for the render target.
//    FBO        for the shadow buffer - any number but 0 for now.
//    FUTURE ... other render targets, g-buffer, etc.
//      0        main back buffer display
// Over drawing groups of objects over other groups of object.
// Generally for handling drawing a UI overlay ontop of a 3D scene.
//    255        Overlay = 0 3D objects with depth - drawn first.
//    254        Overlay = 1 2D UI objects drawn over 3D objects.
//    253        Overlay = 2 ...
// Sorting within the groups of objects.
//      8        Skydome
//      4        Opaque      - currently drawn in order created.
//      1        Transparent - use distance to sort back to front.
func setBucket(pass, overlay uint8) (b uint64) {
	if pass > 0 {
		// Draw shadows and other buffers first.
		// For now just shadow pass. Stick the number in the top 8 bits for
		// all shadow pass objects to be rendered before normal pass display
		// buffer objects.
		b = uint64(pass) << 48
	}

	// Reverse overlay so that an overlay of 0 becomes 255 and is drawn first.
	return b | uint64(math.MaxUint8-overlay)<<40
}

// setDist to camera for sorting transparent objects.
func setDist(bucket uint64, toCam float64) uint64 {
	return bucket | uint64(math.Float32bits(float32(toCam)))
}

// setTransparent marks the object as transparent.
// The distance bits are kept.
func setTransparent(bucket uint64) uint64 {
	return bucket&clearObj | alphaObj // mark transparent.
}

// setSky marks the object as a sky dome.
func setSky(bucket uint64) uint64 {
	return bucket&clearDistance&clearObj | skyObj
}

// setOpaque wipes out the distance so that objects are sorted by
// entity creation order.
// NOTE: this was done to fix bb.go where the billboard and SDF 3D examples
//       would flip distance order with the background based on the distance
//       to center values. Eliminating the distance sort leaves sorting opaque
//       objects under control of the application.
func setOpaque(bucket uint64) (b uint64) {
	return bucket&clearDistance&clearObj | opaqueObj
}

// Usefull bits for setting or clearing the bucket.
const (
	clearDistance = 0xFFFFFFFF00000000
	clearObj      = 0xFFFFFFF0FFFFFFFF
	skyObj        = 0x0000000800000000
	opaqueObj     = 0x0000000400000000
	alphaObj      = 0x0000000100000000
)
