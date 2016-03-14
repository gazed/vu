// Copyright © 2015-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"github.com/gazed/vu/audio"
)

// Noise manages sounds associated with a singe Pov. Each noise must be
// loaded with sound data that has been bound to the audio card in order
// for the noise to be played.
type Noise interface {
	Add(sound string) // Loads and adds a sound.
	Play(index int)   // Play. Loaded and bound sounds only.
}

// Noise
// =============================================================================
// noise implements noise.

// noise deals with sounds that are mapped to a location.
type noise struct {
	eng    *engine    // Entity manager.
	eid    uint64     // Entity identifier related to this sound.
	loaded bool       // True if data has been set.
	snds   []*sound   // one or more sounds.
	loads  []*loadReq // Assets waiting to be loaded.
}

// newNoise allocates data structures for a noise.
func newNoise(eng *engine, eid uint64) *noise {
	return &noise{eng: eng, eid: eid}
}

// Add a sound to the noise and mark the noise as needing loading.
func (n *noise) Add(soundName string) {
	n.loaded = false
	n.loads = append(n.loads, &loadReq{data: n, index: len(n.snds), a: newSound(soundName)})
	n.snds = append(n.snds, newSound(soundName))
}

// Play gets the sounds location and generates a play sound request.
// The play request is sent to a goroutine allowing the goroutine to block
// until the machine can service the request.
func (n *noise) Play(index int) {
	if n.loaded && index >= 0 && index < len(n.snds) {
		snd := n.snds[index]
		if p, ok := n.eng.povs[n.eid]; ok {
			x, y, z := p.Location()
			go func(sid uint64, x, y, z float64) {
				n.eng.machine <- &playSound{sid: sid, x: x, y: y, z: z}
			}(snd.sid, x, y, z)
		}
	}
}

// noise
// =============================================================================
// sound

// sound contains sound data.
type sound struct {
	name       string      // Unique name of the sound.
	tag        uint64      // name and type as a number.
	sid        uint64      // Audio card identifier related to sound location.
	did        uint64      // Audio data reference identifier.
	data       *audio.Data // noise data.
	lx, ly, lz float64     // noise location.
}

// newSound allocates space for a texture object.
func newSound(name string) *sound {
	return &sound{name: name, tag: snd + stringHash(name)<<32, data: &audio.Data{}}
}

// label, aid, and bid are used to uniquely identify assets.
func (s *sound) label() string { return s.name }                  // asset name
func (s *sound) aid() uint64   { return s.tag }                   // asset type and name.
func (s *sound) bid() uint64   { return snd + uint64(s.sid)<<32 } // asset type and bind ref.
