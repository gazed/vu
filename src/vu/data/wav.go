// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package data

import (
	"encoding/binary"
	"fmt"
	"io"
)

// waveHeader is used to load a .wav audio file into memory such that it is
// easily usable by the audio library.  The wave PCM soundfile format is:
//    https://ccrma.stanford.edu/courses/422/projects/WaveFormat/
type waveHeader struct {
	RiffId      [4]byte // "RIFF"
	FileSize    uint32  // Total file size minus 8 bytes.
	WaveId      [4]byte // "WAVE"
	Fmt         [4]byte // "fmt "
	FmtSize     uint32  // Will be 16 for PCM.
	AudioFormat uint16  // Will be 1 for PCM.
	Channels    uint16  // Number of audio channels.
	Frequency   uint32  // 8000, 44100, etc.
	ByteRate    uint32  // SampleRate * NumChannels * BitsPerSample/8.
	BlockAlign  uint16  // NumChannels * BitsPerSample/8.
	SampleBits  uint16  // 8 bits = 8, 16 bits = 16, etc.
	DataId      [4]byte // "data"
	DataSize    uint32  // Size of audio data (total file size minus 44 bytes).
}

// wav takes the resource directory and filename of a wav file and attempts to
// load the audio data into a slice of bytes.
func (l loader) wav(sound *Sound, directory, filename string) (err error) {
	var file io.ReadCloser
	if file, err = l.getResource(directory, filename); err == nil {
		defer file.Close()
		var wh *waveHeader
		var data []byte
		if wh, data, err = l.loadWav(file); err == nil {
			sound.Channels = wh.Channels
			sound.SampleBits = wh.SampleBits
			sound.Frequency = wh.Frequency
			sound.DataSize = wh.DataSize
			sound.AudioData = data
		}
	}
	return err
}

// loadWavFile reads a valid wave file into a header and a bunch audio data into bytes.
// Invalid files return a nil header and an empty data slice.
func (l loader) loadWav(file io.ReadCloser) (wh *waveHeader, bytes []byte, err error) {
	wh = &waveHeader{}
	if err = binary.Read(file, binary.LittleEndian, wh); err != nil {
		return nil, []byte{}, fmt.Errorf("Invalid .wav audio file: %s", err)
	}

	// check that it really is a WAVE file.
	riff, wave := string(wh.RiffId[:]), string(wh.WaveId[:])
	if riff != "RIFF" || wave != "WAVE" {
		return nil, []byte{}, fmt.Errorf("Invalid .wav audio file")
	}

	// read the audio data.
	bytesRead := uint32(0)
	data := make([]byte, wh.DataSize)
	for bytesRead < wh.DataSize {
		inbuff := make([]byte, wh.DataSize)
		inbytes, readErr := file.Read(inbuff)
		if readErr != nil {
			return nil, []byte{}, fmt.Errorf("Corrupt .wav audio file")
		}
		for cnt := 0; cnt < inbytes; cnt++ {
			data[bytesRead] = inbuff[cnt]
			bytesRead += 1
		}
	}
	if bytesRead != wh.DataSize {
		return nil, []byte{}, fmt.Errorf("Invalid .wav audio file")
	}
	return wh, data, nil
}
