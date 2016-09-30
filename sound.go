// Copyright © 2015-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// sound.go wraps the audio package and controls all engine sounds.

import (
	"github.com/gazed/vu/audio"
)

// sound is an engine sound asset. Expected to be accessed through
// the sounds component.
type sound struct {
	name       string      // Unique name of the sound.
	tag        aid         // name and type as a number.
	sid        uint64      // Audio card identifier related to sound location.
	did        uint64      // Audio data reference identifier.
	lx, ly, lz float64     // noise location.
	data       *audio.Data // noise data.
}

// newSound allocates space for a texture object.
func newSound(name string) *sound {
	return &sound{name: name, tag: assetID(snd, name), data: &audio.Data{}}
}

// aid is used to uniquely identify assets.
func (s *sound) aid() aid      { return s.tag }  // hashed type and name.
func (s *sound) label() string { return s.name } // asset name

// =============================================================================

// sounds manages audio instances. Each sound must be loaded with sound data
// that has been bound to the audio card in order for the sound to be played.
type sounds struct {
	eng     *engine          // Needed to load and play sounds.
	assets  map[aid]string   // Collect assets for load request.
	loading map[eid][]aid    // New sounds need to be run through loader.
	active  map[eid][]*sound // Sounds that can be played.

	// Sounds are heard by the sound listener at an app set pov.
	soundListener *Pov    // Single listener for all noises. Current location.
	sx, sy, sz    float64 // Last location of the sound listener.
}

// newSounds creates the sound component manager.
// Expected to be called once on startup.
func newSounds(eng *engine) *sounds {
	ss := &sounds{eng: eng}
	ss.assets = map[aid]string{}   // Reused to submit assets for loading.
	ss.loading = map[eid][]aid{}   // Sounds waiting to be loaded.
	ss.active = map[eid][]*sound{} // Sounds that can be played.
	return ss
}

// create a new sound. Allows multiple sounds to be associated with
// an entity.
func (ss *sounds) create(id eid, name string) {
	aid := assetID(snd, name)
	ss.assets[aid] = name                        // all sound assets.
	ss.loading[id] = append(ss.loading[id], aid) // sounds by entity.
}

// dispose all sounds associated with the given entity.
func (ss *sounds) dispose(id eid) {
	delete(ss.active, id)
	delete(ss.loading, id) // Outstanding loads are ignored when they return.
}

// refresh passes new sounds through the loading system.
func (ss *sounds) refresh() {
	if len(ss.assets) > 0 {
		ss.eng.submitLoadReqs(ss.assets)
		ss.assets = map[aid]string{}
	}
}

// finishLoads matches loaded sounds with loading sounds and moves
// loaded sounds to be active.
func (ss *sounds) finishLoads(assets map[aid]asset) {
	for eid, aids := range ss.loading {
		for _, aid := range aids {
			if a, ok := assets[aid]; ok {
				switch s := a.(type) {
				case *sound:
					ss.active[eid] = append(ss.active[eid], s)
				}
			}
		}
		if len(ss.active[eid]) == len(ss.loading[eid]) {
			delete(ss.loading, eid)
		}
	}
}

// play gets the sounds location and generates a play sound request.
// The play request is sent to a goroutine allowing the goroutine to block
// until the machine can service the request.
func (ss *sounds) play(id eid, index int) {
	if snds, ok := ss.active[id]; ok {
		if index < 0 || index >= len(snds) {
			return
		}
		s := snds[index]
		if p := ss.eng.povs.get(id); p != nil {
			x, y, z := p.At()
			go func(sid uint64, x, y, z float64) {
				ss.eng.machine <- &playSound{sid: sid, x: x, y: y, z: z}
			}(s.sid, x, y, z)
		}
	}
}

// repositionSoundListener checks and updates the sound listeners location.
func (ss *sounds) repositionSoundListener() {
	x, y, z := ss.soundListener.At()
	if x != ss.sx || y != ss.sy || z != ss.sz {
		ss.sx, ss.sy, ss.sz = x, y, z
		go func(x, y, z float64) {
			ss.eng.machine <- &placeListener{x: x, y: y, z: z}
		}(x, y, z)
	}
}

// setListener locates the point that can hear sounds.
// There is always only one listener. It is associated with the root pov
// by default. This changes the listener location to the given pov.
func (ss *sounds) setListener(p *Pov) {
	if p != nil {
		ss.soundListener = p
	}
}
