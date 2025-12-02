// SPDX-FileCopyrightText : © 2014-2022 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gazed/vu/internal/audio/al"
	"github.com/gazed/vu/load"
)

// au checks that OpenAL is installed and the bindings are working
// by loading and playing a sound.
// Requires OpenAL32.dll in this directory or a standard system location.
//
// CONTROLS: NA
func au() {
	// map the bindings to the OpenAL dynamic library.
	al.Init()
	if report := al.BindingReport(); len(report) > 0 {
		for _, line := range report {
			if strings.Contains(line, "[ ]") {
				slog.Error("au: OpenAL not available")
				return
			}
		}
	}

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
			aud, err := load.Audio("bloop.wav")
			if err != nil {
				slog.Error("au: error loading audio file bloop.wav", "error", err)
				return
			}

			// copy the audio data into the buffer
			tag := &autag{}
			attrs := aud.Attrs
			format := tag.audioFormat(attrs)
			if format < 0 {
				slog.Error("au: error recognizing audio format")
				return
			}
			al.BufferData(buffer, int32(format), al.Pointer(&(aud.Data[0])),
				int32(attrs.DataSize), int32(attrs.Frequency))

			// attach the source to a buffer.
			al.Sourcei(source, al.BUFFER, int32(buffer))

			// check for any audio library errors that have happened up to this point.
			if openAlErr := al.GetError(); openAlErr != 0 {
				slog.Error("au: OpenAL error", "error", openAlErr)
				return
			}

			// Start playback and give enough time for the playback to finish.
			// OpenAL can throw a SIGABRT if it is shut down while playing.
			al.SourcePlay(source)
			time.Sleep(1000 * time.Millisecond)
			return
		}
		slog.Error("au: error, failed to get a context")
	}
	slog.Error("au: error, failed to get a device")
}

// Globally unique "tag" for this example.
type autag struct{}

// audioFormat figures out which of the OpenAL formats to use based on the
// WAVE file information.
func (a *autag) audioFormat(attrs *load.AudioAttributes) int32 {
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
