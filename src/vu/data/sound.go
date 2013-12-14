// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package data

// Sound is a shared audio resource. Note that the audio data
// has not yet been bound to a sound card.
type Sound struct {
	Name       string // Unique sound name.
	AudioData  []byte // The raw audio data.
	Channels   uint16 // Number of audio channels.
	SampleBits uint16 // 8 bits = 8, 16 bits = 16, etc.
	Frequency  uint32 // 8000, 44100, etc.
	DataSize   uint32 // Size of audio data (total file size minus header size).

	// Buffer is the sound card buffer reference that the sound data is loaded
	// into. Sound buffers can be shared among many sound sources.
	Buffer uint32
}
