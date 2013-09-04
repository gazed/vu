// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package audio

import (
	"fmt"
	"log"
	"strings"
	"vu/audio/al"
	"vu/data"
)

// openal provides sound support for the engine. It exposes the useful parts
// of the underlying OpenAL audio library as well as providing some sound
// utility methods.
type openal struct {
	dev *al.Device  // created on initialization.
	ctx *al.Context // created on initialization.
}

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

// BindSound loads the raw audio data into the sound buffer.
func (a *openal) BindSound(s *data.Sound) (err error) {
	if alerr := al.GetError(); alerr != al.NO_ERROR {
		log.Printf("openal:bindSound need to find and fix prior error %X", alerr)
	}

	// create the sound buffer and copy the audio data into the buffer
	var format int32
	if format, err = a.format(s); err == nil {
		al.GenBuffers(1, &(s.Buffer))
		al.BufferData(s.Buffer, format, al.Pointer(&(s.AudioData[0])), int32(s.DataSize), int32(s.Frequency))
		if alerr := al.GetError(); alerr != al.NO_ERROR {
			err = fmt.Errorf("Failed binding sound %s", s.Name)
		}
	}
	return
}

// Mute turns the listener gain on/off.
func (a *openal) Mute(mute bool) {
	if mute {
		al.Listenerf(al.GAIN, 0.0)
	} else {
		al.Listenerf(al.GAIN, 1.0)
	}
}

// format figures out which of the OpenAL formats to use based on the
// WAVE file information.  A -1 value, and error, is returned if the format
// cannot be determined.
func (a *openal) format(s *data.Sound) (format int32, err error) {
	format = -1
	if s.Channels == 1 && s.SampleBits == 8 {
		format = al.FORMAT_MONO8
	} else if s.Channels == 1 && s.SampleBits == 16 {
		format = al.FORMAT_MONO16
	} else if s.Channels == 2 && s.SampleBits == 8 {
		format = al.FORMAT_STEREO8
	} else if s.Channels == 2 && s.SampleBits == 16 {
		format = al.FORMAT_STEREO16
	}
	if format < 0 {
		err = fmt.Errorf("openal:format cannot recognize audio format")
	}
	return
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
func newSoundMaker(s *data.Sound) *soundMaker {
	if s == nil {
		return nil
	}
	var source uint32
	al.GenSources(1, &source)
	al.Sourcei(source, al.BUFFER, int32(s.Buffer))
	return &soundMaker{source: source}
}

func (sm *soundMaker) SetLocation(x, y, z float32) { al.Source3f(sm.source, al.POSITION, x, y, z) }
func (sm *soundMaker) SetPitch(pitch float32)      { al.Sourcef(sm.source, al.PITCH, pitch) }
func (sm *soundMaker) SetGain(gain float32)        { al.Sourcef(sm.source, al.GAIN, gain) }
func (sm *soundMaker) Play()                       { al.SourcePlay(sm.source) }

// soundMaker
// ===========================================================================
// soundListener

// soundListener conforms to the the SoundListener interface.
type soundListener struct{}

func (l *soundListener) SetLocation(x, y, z float32) { al.Listener3f(al.POSITION, x, y, z) }
func (l *soundListener) SetVelocity(x, y, z float32) { al.Listener3f(al.VELOCITY, x, y, z) }
func (l *soundListener) SetGain(gain float32)        { al.Listenerf(al.GAIN, gain) }
