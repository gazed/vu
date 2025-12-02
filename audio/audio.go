// SPDX-FileCopyrightText : © 2014-2022 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

// Package audio plays sounds at 3D locations.
// The expected usage is to initialize the audio system and load sound data.
// Then play sounds that are close enough to the sound listener to be audible.
//
// Package audio is provided as part of the vu (virtual universe) 3D engine.
package audio

// Context is used to initialize and play audio.
// It works with an audioAPI to allow different audio players implementations.
type Context struct {
	player audioAPI // audio device.
}

// New provides the default audio implementation.
func New() *Context { return &Context{player: &openal{}} }

// Init the audio context state. Must be called once on startup.
func (c *Context) Init() error { return c.player.init() }

// Closes and the audio layer, releasing any audio resources.
func (c *Context) Dispose() { c.player.dispose() }

// Volume control: valid values are 0->1.
func (c *Context) SetGain(gain float64) { c.player.setGain(gain) }

// LoadSound copies the sound data to the sound card and returns
// references that can be used to play or dispose the sound.
//
//	sound : updated reference to the bound sound.
//	buff  : updated reference to the sound data buffer.
//	d     : sound data bytes and settings to be bound.
func (c *Context) LoadSound(sound, buff *uint64, d *Data) error {
	return c.player.loadSound(sound, buff, d)
}

// DropSound disposes the audio resources allocated with LoadSound.
// Must be called on a valid audio context, ie: before Dispose()
func (c *Context) DropSound(sound, buff uint64) {
	c.player.dropSound(sound, buff)
}

// PlaceListener control sounds by setting the x,y,z locations
// for the listener. While there is only ever one listener,
// there can be many sounds.
func (c *Context) PlaceListener(x, y, z float64) {
	c.player.placeListener(x, y, z)
}

// Play the given sound.
func (c *Context) PlaySound(sound uint64, x, y, z float64) {
	c.player.playSound(sound, x, y, z)
}

// DisableAudio is used to turn off the audio system when
// there are no supported audio drivers.
func (c *Context) DisableAudio() { c.player = &noAudio{} }

// Loader provides the interface for uploading sound data to the
// audio device.
type Loader interface {
	LoadSound(sound, buff *uint64, d *Data) error
}

// ===========================================================================
// audioAPI must be implemented by the audio device layer.
// Audio interacts with the underlying audio layer which in turn interfaces
// to the sound drivers and hardware. Audio must be initialized once before
// sounds can be bound and played.
type audioAPI interface {
	init() error          // Get the audio layer up and running.
	dispose()             // Closes and cleans up the audio layer.
	setGain(gain float64) // Volume control: valid values are 0->1.

	// LoadSound copies the sound data to the sound card and returns
	// references that can be used to play or dispose the sound.
	//     sound : updated reference to the bound sound.
	//     buff  : updated reference to the sound data buffer.
	//     d     : sound data bytes and settings to be bound.
	loadSound(sound, buff *uint64, d *Data) error
	// DropSound disposes the audio resources allocated with LoadSound.
	// Must be called on a valid audio context, ie: before Dispose()
	dropSound(sound, buff uint64)

	// Control sounds by setting the x,y,z locations for a listener
	// and the played sounds. While there is only ever one listener,
	// there can be many sounds.
	placeListener(x, y, z float64)           // Only ever one listener.
	playSound(sound uint64, x, y, z float64) // Play the bound sound.
}

// ===========================================================================
// noAudio can be used to mock out audio for testing or when audio
// initialization fails.
type noAudio struct{}

func (na *noAudio) init() error                                  { return nil }
func (na *noAudio) dispose()                                     {}
func (na *noAudio) setGain(gain float64)                         {}
func (na *noAudio) loadSound(sound, buff *uint64, d *Data) error { return nil }
func (na *noAudio) dropSound(sound, buff uint64)                 {}
func (na *noAudio) placeListener(x, y, z float64)                {}
func (na *noAudio) playSound(sound uint64, x, y, z float64)      {}

// ===========================================================================

// Data is a shared audio resource that is used to load sound
// data onto a sound card.
type Data struct {
	Name       string // Unique sound name.
	AudioData  []byte // The raw audio data.
	Channels   uint16 // Number of audio channels.
	SampleBits uint16 // 8 bits = 8, 16 bits = 16, etc.
	Frequency  uint32 // 8000, 44100, etc.
	DataSize   uint32 // Audio data size: total file size minus header size.
}

// Set is a convenience method that populates sound data with the
// given information. It attempts to reuse the existing sound data buffer.
func (d *Data) Set(channels, sampleBits uint16, frequency, dataSize uint32, data []byte) {
	d.Channels = channels
	d.SampleBits = sampleBits
	d.Frequency = frequency
	d.DataSize = dataSize
	d.AudioData = d.AudioData[:0] // reset, keeping memory.
	d.AudioData = append(d.AudioData, data...)
}
