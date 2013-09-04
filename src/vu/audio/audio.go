// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

// Package audio provides access to 3D sound capability.  It interfaces to the sound
// card through a device layer (currently OpenAL).  The audio package is comprised of
// three main interfaces:
//    1. The overall audio system that needs to be initialized on application
//       startup, and shutdown on application close. It also needs sound data
//       that has been loaded from persistent store and bound to the sound card.
//    2. Sound makers that are associated with a sound, location and volume.
//    3. Sound listeners that are associated with a location.
//
// The expected usage is to initialize the audio system and then load some sounds
// into it. Then setup some sound listeners and some sound makers.  Associate the
// sound makers with the loaded sounds.  Then play sounds that are close enough to
// the SoundListeners to be audible.
//
// OSX has OpenAL. To install 64-bit OpenAL on windows use:
//    http://connect.creativelabs.com/openal/Downloads/oalinst.zip
//
// Package audio is provided as part of the vu (virtual universe) 3D engine.
package audio

import (
	"vu/data"
)

// Audio interfaces to the underling audio layer which in turn interfaces to
// the sound drivers and hardware.  This must be initialized before SoundMakers's
// or SoundListener's can be created and used.
type Audio interface {
	Init() error    // Init gets the audio layer up and running.
	Shutdown()      // Shutdown closes and cleans up the audio layer.
	Mute(mute bool) // Mute turns the listener gain on/off.

	// BindSound copies sound data to the sound card.
	BindSound(s *data.Sound) (err error)
}

// SoundMaker associates a sound with a location and other information.
// A sound maker will only produce an audible sound if there are active sound
// listeners within a reasonable distance.
type SoundMaker interface {
	SetLocation(x, y, z float32) // Where the noise occurs.
	SetPitch(pitch float32)      // The noise pitch.
	SetGain(gain float32)        // The noise volume.
	Play()                       // Make the noise happen now.
}

// SoundListener is the sound receiver.  The listeners location relative to where the
// noise was played determines how much sound comes out of the speakers.
type SoundListener interface {
	SetLocation(x, y, z float32) // The listener location.
	SetVelocity(x, y, z float32) // The listeners movement.
	SetGain(gain float32)        // Crank this up for the hard of hearing.
}

// Audio, SoundMaker, SoundListener interfaces
// ===========================================================================
// Provide default implementations.

// New provides a default Audio implementation.
func New() Audio { return &openal{} }

// NewSoundMaker provides an audio generator.
func NewSoundMaker(sound *data.Sound) SoundMaker { return newSoundMaker(sound) }

// NewSoundListener provides an audio receiver.
func NewSoundListener() SoundListener { return &soundListener{} }
