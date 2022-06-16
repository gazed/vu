// Copyright © 2015-2018 Galvanized Logic Inc.

package vu

// sound.go wraps the audio package and controls all engine sounds.

import (
	"log/slog"
)

// PlaySound plays the given sound at this entities location.
//   - soundID : entity created with Eng.AddSound.
//
// Depends on Entity.AddSound.
func (e *Entity) PlaySound(soundID uint32) {
	if p := e.app.povs.get(e.eid); p != nil {
		e.app.sounds.addNoise(eID(soundID), e.eid)
		return
	}
	slog.Error("PlaySound requires location", "entity", e.eid)
}

// SetListener sets the location of the sound listener to be this entity.
//
// Depends on Entity.AddSound.
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
	all      map[eID]*sound // All sounds assets.
	rebinds  map[eID]*sound // Sounds needing rebind.
	ready    map[eID]*sound // Sounds that can be played.
	noises   map[eID]eID    // Sounds to be played at pov eid.
	listener eID            // Pov listener location.
}

func (ss *sounds) get(eid eID) *sound { return ss.all[eid] }

// newSounds creates the sound component manager.
// Expected to be called once on startup.
func newSounds() *sounds {
	ss := &sounds{}
	ss.all = map[eID]*sound{}     // All sounds.
	ss.ready = map[eID]*sound{}   // Sounds that can be played.
	ss.rebinds = map[eID]*sound{} // Sounds waiting for main thread bind.
	ss.noises = map[eID]eID{}     // Sounds waiting for main thread play.
	return ss
}

// create a new sound. Allows multiple sounds to be associated with
// an entity. Called by the application through the update goroutine.
func (ss *sounds) create(eids *entities, name string) (eid eID) {
	sn := newSound(name)
	for eid, s := range ss.all {
		if sn.tag == s.tag {
			// application error, please fix since eid is now invalid.
			slog.Warn("sound already created", "name", name)
			return eid
		}
	}

	// create a new sound entity.
	eid = eids.create()
	ss.all[eid] = sn // all sound assets.
	return eid
}

// assetLoaded is the asset loader callback.
func (ss *sounds) assetLoaded(eid eID, a asset) {
	switch la := a.(type) {
	case *sound:
		ss.ready[eid] = la   // sound data available.
		ss.rebinds[eid] = la // request rebind before next update.
	default:
		slog.Error("unexepected sound asset", "name", a.label())
	}
}

// rebind is called to bind or rebind sound data.
// This moves the data to the audio device.
func (ss *sounds) rebind(eng *Engine) {
	for eid, s := range ss.rebinds {
		if eng.ac != nil {
			err := eng.ac.LoadSound(&s.sid, &s.did, s.data)
			if err != nil {
				slog.Error("bind sound failed", "name", s.name, "error", err)
				return // dev error - asset should be bindable.
			}
			delete(ss.rebinds, eid)
		}
	}
}

// addNoise saves a sound to be played later on the main thread.
func (ss *sounds) addNoise(soundID, pov eID) {
	ss.noises[soundID] = pov
}

// setListener saves the location of the listener pov.
// so that it can be set later on the main thread.
func (ss *sounds) setListener(pov eID) {
	ss.listener = pov
}

// play the sounds requested during update. Called on the main thread.
func (ss *sounds) play(eng *Engine) {
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

// dispose of sound data associated with the given entity.
func (ss *sounds) dispose(eng *Engine, eid eID) {
	if s := ss.all[eid]; s != nil {
		delete(ss.ready, eid)
		delete(ss.rebinds, eid)
		delete(ss.all, eid)
		delete(ss.noises, eid)

		// delete the sound resources.
		eng.ac.DropSound(s.sid, s.did)
	}
}
