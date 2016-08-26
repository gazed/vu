// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package load

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Wav attempts to load WAV based audio data into SndData.
// The wave PCM soundfile format is from:
//     https://ccrma.stanford.edu/courses/422/projects/WaveFormat
// The Reader r is expected to be opened and closed by the caller.
// A successful import overwrites the data in SndData.
func Wav(r io.Reader, d *SndData) (err error) {
	hdr := &wavHeader{}
	if err = binary.Read(r, binary.LittleEndian, hdr); err != nil {
		return fmt.Errorf("Invalid .wav audio file: %s", err)
	}

	// check that it really is a WAVE file.
	riff, wave := string(hdr.RiffID[:]), string(hdr.WaveID[:])
	if riff != "RIFF" || wave != "WAVE" {
		return fmt.Errorf("Invalid .wav audio file")
	}

	// read the audio data.
	bytesRead := uint32(0)
	data := []byte{}
	inbuff := make([]byte, hdr.DataSize)
	for bytesRead < hdr.DataSize {
		inbytes, readErr := r.Read(inbuff)
		if readErr != nil {
			return fmt.Errorf("Corrupt .wav audio file")
		}
		data = append(data, inbuff...)
		bytesRead += uint32(inbytes)
	}
	if bytesRead != hdr.DataSize {
		return fmt.Errorf("Invalid .wav audio file %d %d", bytesRead, hdr.DataSize)
	}
	attrs := &SndAttributes{Channels: hdr.Channels, Frequency: hdr.Frequency,
		DataSize: hdr.DataSize, SampleBits: hdr.SampleBits}
	d.Attrs, d.Data = attrs, data
	return nil
}

// wavHeader is an internal implementation for loading WAV files.
type wavHeader struct {
	RiffID      [4]byte // "RIFF"
	FileSize    uint32  // Total file size minus 8 bytes.
	WaveID      [4]byte // "WAVE"
	Fmt         [4]byte // "fmt "
	FmtSize     uint32  // Will be 16 for PCM.
	AudioFormat uint16  // Will be 1 for PCM.
	Channels    uint16  // Number of audio channels.
	Frequency   uint32  // 8000, 44100, etc.
	ByteRate    uint32  // SampleRate * NumChannels * BitsPerSample/8.
	BlockAlign  uint16  // NumChannels * BitsPerSample/8.
	SampleBits  uint16  // 8 bits = 8, 16 bits = 16, etc.
	DataID      [4]byte // "data"
	DataSize    uint32  // Size of audio data: total file size minus 44 bytes.
}
