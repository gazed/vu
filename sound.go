// Copyright © 2015-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// sound.go wraps the audio package and controls all engine sounds.

import (
	"log"

	"github.com/gazed/vu/audio"
)

// PlaySound plays the given sound at this entities location.
//    soundID : entity created with Eng.AddSound.
//
// Depends on Ent.AddPart.
func (e *Ent) PlaySound(soundID uint32) {
	if p := e.app.povs.get(e.eid); p != nil {
		e.app.sounds.addNoise(eid(soundID), e.eid)
		return
	}
	log.Printf("PlaySound needs AddPart %d", e.eid)
}

// SetListener sets the location of the sound listener to be this entity.
//
// Depends on Ent.AddPart.
func (e *Ent) SetListener() {
	if p := e.app.povs.get(e.eid); p != nil {
		e.app.sounds.setListener(e.eid)
		return
	}
	log.Printf("SetListener needs pov %d", e.eid)
}

// pov sound related entity methods
// =============================================================================
// sound

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

// sound
// =============================================================================
// sounds: component manager for sound.

// sounds manages audio instances. Each sound must be loaded with sound data
// that has been bound to the audio card in order for the sound to be played.
type sounds struct {
	all      map[eid]*sound // All sounds assets.
	loading  map[eid]*sound // New sounds need to be run through loader.
	rebinds  map[eid]*sound // Sounds needing rebind.
	ready    map[eid]*sound // Sounds that can be played.
	noises   map[eid]eid    // Sounds to be played at pov eid.
	listener eid            // Pov listener location.
}

// newSounds creates the sound component manager.
// Expected to be called once on startup.
func newSounds() *sounds {
	ss := &sounds{}
	ss.all = map[eid]*sound{}     // All sounds.
	ss.loading = map[eid]*sound{} // Sounds waiting to be loaded.
	ss.ready = map[eid]*sound{}   // Sounds that can be played.
	ss.rebinds = map[eid]*sound{} // Sounds waiting for main thread bind.
	ss.noises = map[eid]eid{}     // Sounds waiting for main thread play.
	return ss
}

// create a new sound. Allows multiple sounds to be associated with
// an entity. Called by the application through the update goroutine.
func (ss *sounds) create(ld *loader, eid eid, name string) {
	sn := newSound(name)
	ss.all[eid] = sn     // all sound assets.
	ss.loading[eid] = sn // sounds need to be run through loader.

	// create a callback closure with the entity id.
	callback := func(a asset) { ss.loaded(eid, a) }
	ld.fetch(newSound(name), callback)
}

// loaded is the asset loader callback. Called on the loader goroutine.
func (ss *sounds) loaded(eid eid, a asset) {
	switch la := a.(type) {
	case *sound:
		if _, ok := ss.loading[eid]; ok {
			delete(ss.loading, eid) // placeholder
			ss.ready[eid] = la      // loaded sound.
			ss.rebinds[eid] = la    // request rebind before next update.
		} else {
			log.Printf("Expected loading sound for: %s", a.label())
		}
	default:
		log.Printf("Unexepected sound asset: %s", a.label())
	}
}

// rebind is called to bind or rebind sound data. This moves the
// data to the audio card. Called on the main thread.
func (ss *sounds) rebind(eng *engine) {
	for eid, s := range ss.rebinds {
		if err := eng.bind(s); err != nil {
			log.Printf("Bind sound %s failed: %s", s.name, err)
			return // dev error - asset should be bindable.
		}
		delete(ss.rebinds, eid)
	}
}

// addNoise saves a sound to be played later on the main thread.
func (ss *sounds) addNoise(soundID, pov eid) {
	ss.noises[soundID] = pov
}

// setListener saves the location of the listener pov.
// so that it can be set later on the main thread.
func (ss *sounds) setListener(pov eid) {
	ss.listener = pov
}

// play the sounds requested during update. Called on the main thread.
func (ss *sounds) play(eng *engine) {
	if ss.listener != 0 {
		// reposition sound listener if necessary.
		if pov := eng.app.povs.get(ss.listener); pov != nil {
			eng.ac.PlaceListener(pov.at())
		}
		ss.listener = 0
	}
	for sid, eid := range ss.noises {
		s := ss.ready[sid]
		pov := eng.app.povs.get(eid)
		if s != nil && pov != nil {
			x, y, z := pov.at()
			eng.ac.PlaySound(s.sid, x, y, z)
		}
		delete(ss.noises, sid)
	}
}

// dispose all sounds associated with the given entity.
func (ss *sounds) dispose(eid eid) {
	delete(ss.loading, eid) // Outstanding loads are ignored when they return.
	delete(ss.ready, eid)
	delete(ss.all, eid)
}
