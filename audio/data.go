// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package audio

// Data is a shared audio resource that is used to load sound
// data onto a sound card using BindSound().
type Data struct {
	Name       string // Unique sound name.
	AudioData  []byte // The raw audio data.
	Channels   uint16 // Number of audio channels.
	SampleBits uint16 // 8 bits = 8, 16 bits = 16, etc.
	Frequency  uint32 // 8000, 44100, etc.
	DataSize   uint32 // Size of audio data (total file size minus header size).
}

// Set is a convenience method that populates sound data with the
// given information. It attempts to reuse the existing sound buffer.
func (d *Data) Set(channels, sampleBits uint16, frequency, dataSize uint32, data []byte) {
	d.Channels = channels
	d.SampleBits = sampleBits
	d.Frequency = frequency
	d.DataSize = dataSize
	d.AudioData = d.AudioData[:0] // reset, keeping memory.
	d.AudioData = append(d.AudioData, data...)
}
