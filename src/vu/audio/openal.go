// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package audio

import (
	"fmt"
	"log"
	"strings"
	"vu/audio/al"
)

// OSX has OpenAL. To install 64-bit OpenAL on windows use:
//    http://connect.creativelabs.com/openal/Downloads/oalinst.zip

// openal provides sound support for the engine. It exposes the useful parts
// of the underlying OpenAL audio library as well as providing some sound
// utility methods.
type openal struct {
	dev *al.Device  // created on initialization.
	ctx *al.Context // created on initialization.
}

// NewSoundMaker provides an audio generator.
func (a *openal) NewSoundMaker(s Sound) SoundMaker { return newSoundMaker(s.(*sound)) }

// NewSoundListener provides an audio receiver.
func (a *openal) NewSoundListener() SoundListener { return &soundListener{} }

// NewSound provides a place to put sound data.
func (a *openal) NewSound(name string) Sound { return newSound(name) }

// Init runs the one time openal library initialization. It is only expected to
// be called once by the engine on startup.
func (a *openal) Init() (err error) {
	al.Init()
	if err = a.validate(); err != nil {
		return
	}

	// create an openal context for all sounds.
	if a.dev = al.OpenDevice(""); a.dev != nil {
		if a.ctx = al.CreateContext(a.dev, nil); a.ctx != nil {
			al.MakeContextCurrent(a.ctx)
			return // success
		}
	}
	return fmt.Errorf("Could not initialize openal audio")
}

// validate that OpenAL is available. OSX has OpenAL. To install OpenAL on
// windows use:
//    http://connect.creativelabs.com/openal/Downloads/oalinst.zip
func (a *openal) validate() error {
	if report := al.BindingReport(); len(report) > 0 {
		for _, line := range report {
			if strings.Contains(line, "[-]") {
				return fmt.Errorf("OpenAL uninitialized.")
			}
		}
	} else {
		return fmt.Errorf("OpenAL unavailable.")
	}
	return nil
}

// Shutdown closes down the openal library.  This is only expected to be called
// once by the engine when it is shutting down.
func (a *openal) Shutdown() {
	al.MakeContextCurrent(nil)
	if a.ctx != nil {
		al.DestroyContext(a.ctx)
	}
	if a.dev != nil {
		al.CloseDevice(a.dev)
	}
}

// Mute turns the listener gain on/off.
func (a *openal) Mute(mute bool) {
	if mute {
		al.Listenerf(al.GAIN, 0.0)
	} else {
		al.Listenerf(al.GAIN, 1.0)
	}
}

// openal
// ===========================================================================
// soundMaker

// soundMaker enables linking a sound buffer to a sound source.  A noise can have a
// location, orientation, and velocity.  A noise can be played.
type soundMaker struct {
	source uint32
}

// newSoundMaker creates a noise using the given sound.
// It conforms to the the SoundMaker interface.
func newSoundMaker(s *sound) *soundMaker {
	if s == nil {
		return nil
	}
	var source uint32
	al.GenSources(1, &source)
	al.Sourcei(source, al.BUFFER, int32(s.buffer))
	return &soundMaker{source: source}
}

func (sm *soundMaker) SetLocation(x, y, z float64) {
	al.Source3f(sm.source, al.POSITION, float32(x), float32(y), float32(z))
}
func (sm *soundMaker) SetPitch(pitch float64) { al.Sourcef(sm.source, al.PITCH, float32(pitch)) }
func (sm *soundMaker) SetGain(gain float64)   { al.Sourcef(sm.source, al.GAIN, float32(gain)) }
func (sm *soundMaker) Play()                  { al.SourcePlay(sm.source) }

// soundMaker
// ===========================================================================
// soundListener

// soundListener conforms to the the SoundListener interface.
type soundListener struct{}

func (l *soundListener) SetLocation(x, y, z float64) {
	al.Listener3f(al.POSITION, float32(x), float32(y), float32(z))
}
func (l *soundListener) SetVelocity(x, y, z float64) {
	al.Listener3f(al.VELOCITY, float32(x), float32(y), float32(z))
}
func (l *soundListener) SetGain(gain float64) { al.Listenerf(al.GAIN, float32(gain)) }

// soundListener
// ============================================================================
// sound

// sound is a shared audio resource. Note that the audio data
// needs to be bound to a sound card using Bind().
type sound struct {
	name       string // Unique sound name.
	audioData  []byte // The raw audio data.
	channels   uint16 // Number of audio channels.
	sampleBits uint16 // 8 bits = 8, 16 bits = 16, etc.
	frequency  uint32 // 8000, 44100, etc.
	dataSize   uint32 // Size of audio data (total file size minus header size).

	// Buffer is the sound card buffer reference that the sound data is loaded
	// into. Sound buffers can be shared among many sound sources.
	buffer uint32
}

// newSound allocates space for sound data.
func newSound(name string) *sound { return &sound{name: name} }

// Name implements Sound.
func (s *sound) Name() string { return s.name }

// SetData populates sound s with newly loaded data.
func (s *sound) SetData(channels, sampleBits uint16, frequency, dataSize uint32, data []byte) {
	s.channels = channels
	s.sampleBits = sampleBits
	s.frequency = frequency
	s.dataSize = dataSize
	s.audioData = s.audioData[:0] // reset, keeping memory.
	s.audioData = append(s.audioData, data...)
}

// Bind loads the raw audio data into the sound buffer.
func (s *sound) Bind() (err error) {
	if alerr := al.GetError(); alerr != al.NO_ERROR {
		log.Printf("openal.sound.Bind need to find and fix prior error %X", alerr)
	}

	// create the sound buffer and copy the audio data into the buffer
	var format int32
	if format, err = s.format(); err == nil {
		al.GenBuffers(1, &(s.buffer))
		al.BufferData(s.buffer, format, al.Pointer(&(s.audioData[0])), int32(s.dataSize), int32(s.frequency))
		if alerr := al.GetError(); alerr != al.NO_ERROR {
			err = fmt.Errorf("Failed binding sound %s", s.name)
		}
	}
	return
}

// format figures out which of the OpenAL formats to use based on the
// WAVE file information.  A -1 value, and error, is returned if the format
// cannot be determined.
func (s *sound) format() (format int32, err error) {
	format = -1
	if s.channels == 1 && s.sampleBits == 8 {
		format = al.FORMAT_MONO8
	} else if s.channels == 1 && s.sampleBits == 16 {
		format = al.FORMAT_MONO16
	} else if s.channels == 2 && s.sampleBits == 8 {
		format = al.FORMAT_STEREO8
	} else if s.channels == 2 && s.sampleBits == 16 {
		format = al.FORMAT_STEREO16
	}
	if format < 0 {
		err = fmt.Errorf("openal:format cannot recognize audio format")
	}
	return
}
