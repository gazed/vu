// SPDX-FileCopyrightText : © 2016-2025 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

package vu

// sound.go wraps the audio package and controls all engine sounds.

import (
	"log/slog"
)

// PlaySound plays the given sound at this entities location.
//   - soundID : entity created with Eng.AddSound.
//
// Depends on Engine.AddSound.
func (e *Entity) PlaySound(eng *Engine, sound *Entity) {
	if p := e.app.povs.get(e.eid); p != nil {
		if s := e.app.sounds.get(sound.eid); s != nil {
			e.app.sounds.play(eng, s, p)
		}
		return
	}
	slog.Error("PlaySound requires location", "entity", e.eid)
}

// SetListener sets the location of the sound listener to be this entity.
//
// Depends on Engine.AddSound.
func (e *Entity) SetListener() {
	if p := e.app.povs.get(e.eid); p != nil {
		e.app.sounds.setListener(e.eid)
		return
	}
	slog.Error("SetListener requires location", "entity", e.eid)
}

// =============================================================================
// sounds: component manager for sound.

// sounds manages audio instances. Each sound must be loaded with sound data
// that has been bound to the audio card in order for the sound to be played.
type sounds struct {
	list     map[eID]*sound // loaded sounds assets.
	listener eID            // Pov listener location.
}

func (ss *sounds) get(eid eID) *sound { return ss.list[eid] }

// newSounds creates the sound component manager.
// Expected to be called once on startup.
func newSounds() *sounds {
	ss := &sounds{}
	ss.list = map[eID]*sound{} // Sounds ready to be played.
	return ss
}

// create a new sound Entity.
func (ss *sounds) create(eids *entities, name string) (eid eID) {
	return eids.create()
}

// assetLoaded associates a loaded sound asset with the given entity
func (ss *sounds) assetLoaded(eid eID, a asset) {
	switch la := a.(type) {
	case *sound:
		ss.list[eid] = la
	}
}

// play the given sound.
func (ss *sounds) play(eng *Engine, sound *sound, pov *pov) {
	if sound != nil && pov != nil {
		x, y, z := pov.at()
		eng.ac.PlaySound(sound.sid, x, y, z)
	}
}

// setListener saves the location of the listener pov.
// so that it can be set later on the main thread.
func (ss *sounds) setListener(pov eID) {
	ss.listener = pov
}

// dispose of sound data associated with the given entity.
func (ss *sounds) dispose(eng *Engine, eid eID) {
	if s := ss.list[eid]; s != nil {
		delete(ss.list, eid)

		// delete the sound resources.
		eng.ac.DropSound(s.sid, s.did)
	}
}
