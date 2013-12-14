// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

// Package audio provides access to 3D sound capability.
// The package is comprised of three interfaces:
//    • The overall audio system that needs to be initialized on application
//      startup and shutdown on application close. It also needs sound data
//      to be loaded from persistent store and bound to the sound card.
//    • Sound makers that are associated with a sound, location and volume.
//    • Sound listeners that are associated with a location.
//
// The expected usage is to initialize the audio system and load some sounds.
// Then create sound listeners and sound makers. Associate the sound makers
// with sounds. Finally have the sound makers play sounds that are close
// enough to the sound listeners to be audible.
//
// Package audio is provided as part of the vu (virtual universe) 3D engine.
package audio

import (
	"vu/data"
)

// Audio interacts with the underlying audio layer which in turn interfaces to
// the sound drivers and hardware. This must be initialized before SoundMakers
// or SoundListeners can be created and used.
type Audio interface {
	Init() error                         // Get the audio layer up and running.
	Shutdown()                           // Closes and cleans up the audio layer.
	Mute(mute bool)                      // Turns the listener gain on/off.
	BindSound(s *data.Sound) (err error) // Copy sound data to the sound card.
}

// SoundMaker associates a sound with a location and other information.
// A sound maker will only produce an audible sound if there are active sound
// listeners within a reasonable distance.
type SoundMaker interface {
	SetLocation(x, y, z float64) // Where the noise occurs.
	SetPitch(pitch float64)      // Noise pitch.
	SetGain(gain float64)        // Noise volume.
	Play()                       // Make the noise happen now.
}

// SoundListener is the sound receiver.  The listeners location relative to where
// the noise was played determines how much sound comes out of the speakers.
type SoundListener interface {
	SetLocation(x, y, z float64) // Listener location.
	SetVelocity(x, y, z float64) // Listeners movement.
	SetGain(gain float64)        // Crank this up for the hard of hearing.
}

// Audio, SoundMaker, SoundListener interfaces
// ===========================================================================
// Provide default implementations.

// New provides a default audio implementation.
func New() Audio { return &openal{} }

// NewSoundMaker provides an audio generator.
func NewSoundMaker(sound *data.Sound) SoundMaker { return newSoundMaker(sound) }

// NewSoundListener provides an audio receiver.
func NewSoundListener() SoundListener { return &soundListener{} }
