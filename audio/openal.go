// Copyright © 2013-2024 Galvanized Logic Inc.

package audio

// openal.go provides the wrapper for the openal bindings.

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/gazed/vu/internal/audio/al"
)

// OpenAL (https://openal.org)
// latest 64-bit version `soft_oal.dll` from https://openal-soft.org/openal-binaries/

// openal provides sound support for the engine. It exposes the useful parts
// of the underlying OpenAL audio library as well as providing some sound
// utility methods.
type openal struct {
	dev al.Device  // created on initialization.
	ctx al.Context // created on initialization.
}

// init runs the one time openal library initialization. It is expected to
// be called once by the engine on startup.
func (a *openal) init() (err error) {
	if err := al.Init(); err != nil {
		return fmt.Errorf("openal init %s", err)
	}
	if err = a.validate(); err != nil {
		return fmt.Errorf("openal validate %s", err)
	}

	// Open the audio device create a context for all sounds.
	if a.dev = al.OpenDevice(""); a.dev == 0 {
		return fmt.Errorf("OpenAL device failed %d", al.GetError())
	}
	if a.ctx = al.CreateContext(a.dev, nil); a.ctx == 0 {
		return fmt.Errorf("OpenAL context failed %d", al.GetError())
	}
	al.MakeContextCurrent(a.ctx)
	return nil // success
}

// validate that OpenAL is available. OSX has OpenAL.
func (a *openal) validate() error {
	if report := al.BindingReport(); len(report) > 0 {
		for _, line := range report {
			if strings.Contains(line, "[ ]") {
				return fmt.Errorf("OpenAL uninitialized")
			}
		}
	} else {
		return fmt.Errorf("OpenAL unavailable")
	}
	return nil
}

// dispose closes down the openal library. This is expected
// to be called once by the engine when it is shutting down.
func (a *openal) dispose() {
	al.MakeContextCurrent(0)
	if a.ctx != 0 {
		al.DestroyContext(a.ctx)
	}
	if a.dev != 0 {
		al.CloseDevice(a.dev)
	}
}

// setGain sets the listener gain to a value between 0 and 1.
// Values outside the 0 to 1 range are ignored.
func (a *openal) setGain(zeroToOne float64) {
	if zeroToOne >= 0 && zeroToOne <= 1 {
		al.Listenerf(al.GAIN, float32(zeroToOne))
	}
}

// loadSound copies sound data to the sound card. If successful then the
// sound reference, snd, and sound data buffer reference, buff are updated
// with valid references.
func (a *openal) loadSound(snd, buff *uint64, d *Data) (err error) {
	if alerr := al.GetError(); alerr != al.NO_ERROR {
		slog.Error("openal.BindSound find and fix prior error", "error", alerr)
	}

	// create the sound buffer and copy the audio data into the buffer
	var buff32, snd32 uint32
	var format int32
	if format, err = a.format(d); err == nil {
		al.GenBuffers(1, &buff32)
		al.BufferData(buff32, format, al.Pointer(&(d.AudioData[0])), int32(d.DataSize), int32(d.Frequency))
		*buff = uint64(buff32)
		if alerr := al.GetError(); alerr != al.NO_ERROR {
			err = fmt.Errorf("Failed binding sound %s", d.Name)
		} else {
			al.GenSources(1, &snd32)
			al.Sourcei(snd32, al.BUFFER, int32(*buff))
			*snd = uint64(snd32)
		}
	}
	return err
}

// Implement audioAPI.
func (a *openal) placeListener(x, y, z float64) {
	al.Listener3f(al.POSITION, float32(x), float32(y), float32(z))
}

// Implement audioAPI.
func (a *openal) playSound(snd uint64, x, y, z float64) {
	al.Source3f(uint32(snd), al.POSITION, float32(x), float32(y), float32(z))
	al.SourcePlay(uint32(snd))
}

// Implement Audio.
func (a *openal) dropSound(snd, buff uint64) {
	snd32 := uint32(snd)
	buff32 := uint32(buff)
	al.DeleteSources(1, &snd32)  // delete source first...
	al.DeleteBuffers(1, &buff32) // ...then delete related buffer.
}

// format figures out which of the OpenAL formats to use based on the
// WAVE file information. A -1 value, and error, is returned if the format
// cannot be determined.
func (a *openal) format(d *Data) (format int32, err error) {
	format = -1
	if d.Channels == 1 && d.SampleBits == 8 {
		format = al.FORMAT_MONO8
	} else if d.Channels == 1 && d.SampleBits == 16 {
		format = al.FORMAT_MONO16
	} else if d.Channels == 2 && d.SampleBits == 8 {
		format = al.FORMAT_STEREO8
	} else if d.Channels == 2 && d.SampleBits == 16 {
		format = al.FORMAT_STEREO16
	}
	if format < 0 {
		err = fmt.Errorf("openal:format cannot recognize audio format")
	}
	return format, err
}
