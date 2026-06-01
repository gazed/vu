// SPDX-FileCopyrightText : © 2022-2024 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

package render

// pass.go

import (
	"github.com/gazed/vu/load"
)

// PassID identifies a render.Pass.
// Lower number render passes are rendered before higher numbers.
type PassID uint8 // upto 256 passes should be sufficient.

const (
	// Main render pass that handles both 3D and 2D models.
	MainRenderPass PassID = iota // currently the only pass.
	// FUTURE: add more render passes as needed.
)

// NewPass initializes a render pass.
// The returned Pass is expected to be reused in render loops.
func NewPass() Pass {
	return Pass{
		Uniforms: map[load.SceneUniform][]byte{},              // map of data
		Lights:   []Light{Light{}, Light{}, Light{}, Light{}}, // max 4 lights
	}
}

// Pass contains a group of Packets for rendering in this render pass.
type Pass struct {

	// Packets are a reusable list of packets, one per model.
	Packets  Packets
	Uniforms map[load.SceneUniform][]byte // Scene uniform data

	// Light position and color information.
	Lights []Light // Lights are reused.
}

// Reset the pass data.
func (rp *Pass) Reset() {
	for i := range load.SceneUniforms {
		d, ok := rp.Uniforms[i]
		if !ok {
			rp.Uniforms[i] = []byte{}
		} else {
			rp.Uniforms[i] = d[:0] // reset keeping memory
		}

		// reset lights.
		for _, l := range rp.Lights {
			l.reset()
		}
	}
	rp.Packets = rp.Packets[:0] // reset packets, keeping allocated memory
}
