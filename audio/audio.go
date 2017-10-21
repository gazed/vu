// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// Package audio provides access to 3D sound capability.
// The expected usage is to initialize the audio system and load sound data.
// Then play sounds that are close enough to the sound listener to be audible.
//
// Package audio is provided as part of the vu (virtual universe) 3D engine.
package audio

// Audio interacts with the underlying audio layer which in turn interfaces
// to the sound drivers and hardware. Audio must be initialized once before
// sounds can be bound and played.
type Audio interface {
	Init() error          // Get the audio layer up and running.
	Dispose()             // Closes and cleans up the audio layer.
	SetGain(gain float64) // Volume control: valid values are 0->1.

	// BindSound copies the sound data to the sound card and returns
	// references that can be used to dispose of the sound with ReleaseSound.
	//     sound : updated reference to the bound sound.
	//     buff  : updated reference to the sound data buffer.
	//     d     : sound data bytes and settings to be bound.
	BindSound(sound, buff *uint64, d *Data) error
	ReleaseSound(sound uint64)

	// Control sounds by setting the x,y,z locations for a listener
	// and the played sounds. While there is only ever one listener,
	// there can be many sounds.
	PlaceListener(x, y, z float64)           // Only ever one listener.
	PlaySound(sound uint64, x, y, z float64) // Play the bound sound.
}

// Audio
// ===========================================================================
// Provide native implementation.

// New provides a default audio implementation.
func New() Audio { return audioWrapper() }

// ===========================================================================
// Provide mock implementation.

// NoAudio can be used to mock out audio when audio initialization fails.
type NoAudio struct{}

func (na *NoAudio) Init() error                                  { return nil }
func (na *NoAudio) Dispose()                                     {}
func (na *NoAudio) SetGain(gain float64)                         {}
func (na *NoAudio) BindSound(sound, buff *uint64, d *Data) error { return nil }
func (na *NoAudio) ReleaseSound(sound uint64)                    {}
func (na *NoAudio) PlaceListener(x, y, z float64)                {}
func (na *NoAudio) PlaySound(sound uint64, x, y, z float64)      {}
