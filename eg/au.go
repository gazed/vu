// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"log"
	"time"

	"github.com/gazed/vu/audio/al"
	"github.com/gazed/vu/load"
)

// au demonstrates basic audio library, vu/audio/al, capabilities.
// It checks that OpenAL is installed and the bindings are working
// by loading and playing a sound.
// See the vu:Pov interface when using sound with the vu engine.
//
// CONTROLS: NA
func au() {
	al.Init() // map the bindings to the OpenAL dynamic library.

	// open the default device.
	if dev := al.OpenDevice(""); dev != 0 {
		defer al.CloseDevice(dev)

		// create a context on the device.
		if ctx := al.CreateContext(dev, nil); ctx != 0 {
			defer al.MakeContextCurrent(0)
			defer al.DestroyContext(ctx)
			al.MakeContextCurrent(ctx)

			// create buffers and sources
			var buffer, source uint32
			al.GenBuffers(1, &buffer)
			al.GenSources(1, &source)

			// read in the audio data.
			sound := &load.SndData{}
			err := sound.Load("bloop", load.NewLocator())
			if err != nil {
				log.Printf("au: error loading audio file %s %s", "bloop", err)
				return
			}

			// copy the audio data into the buffer
			tag := &autag{}
			attrs := sound.Attrs
			format := tag.audioFormat(attrs)
			if format < 0 {
				log.Printf("au: error recognizing audio format")
				return
			}
			al.BufferData(buffer, int32(format), al.Pointer(&(sound.Data[0])),
				int32(attrs.DataSize), int32(attrs.Frequency))

			// attach the source to a buffer.
			al.Sourcei(source, al.BUFFER, int32(buffer))

			// check for any audio library errors that have happened up to this point.
			if openAlErr := al.GetError(); openAlErr != 0 {
				log.Printf("au: OpenAL error %d", openAlErr)
				return
			}

			// Start playback and give enough time for the playback to finish.
			// OpenAL can throw a SIGABRT if it is shut down while playing.
			al.SourcePlay(source)
			time.Sleep(1000 * time.Millisecond)
			return
		}
		log.Printf("au: error, failed to get a context")
	}
	log.Printf("au: error, failed to get a device")
}

// Globally unique "tag" for this example.
type autag struct{}

// audioFormat figures out which of the OpenAL formats to use based on the
// WAVE file information.
func (a *autag) audioFormat(attrs *load.SndAttributes) int32 {
	format := int32(-1)
	if attrs.Channels == 1 && attrs.SampleBits == 8 {
		format = al.FORMAT_MONO8
	} else if attrs.Channels == 1 && attrs.SampleBits == 16 {
		format = al.FORMAT_MONO16
	} else if attrs.Channels == 2 && attrs.SampleBits == 8 {
		format = al.FORMAT_STEREO8
	} else if attrs.Channels == 2 && attrs.SampleBits == 16 {
		format = al.FORMAT_STEREO16
	}
	return format
}
