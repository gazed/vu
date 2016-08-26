// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package audio

import (
	"testing"
	// "time"

	"github.com/gazed/vu/load"
)

// test that an audio resource can be loaded. Mimics the steps taken
// by the engine. Depends on sound resources from the examples directory.
func TestAudio(t *testing.T) {
	a := audioWrapper()
	a.Init()
	s := &load.SndData{}
	soundData := &Data{}
	if err := s.Load("bloop", load.NewLocator().Dir("WAV", "../eg/audio")); err == nil {
		at := s.Attrs
		soundData.Set(at.Channels, at.SampleBits, at.Frequency, at.DataSize, s.Data)
	}
	snd, buff := uint64(0), uint64(0)
	if err := a.BindSound(&snd, &buff, soundData); err != nil || buff == 0 || snd == 0 {
		t.Errorf("Failed to load audio resource %d %d : %s", snd, buff, err)
	}

	// Don't play noises during normal testing, but if you're interested...
	// ... then uncomment and "import time" (need to sleep for the sound to happen).
	// a.PlaySound(snd, 0, 0, 0)
	// time.Sleep(1000 * time.Millisecond)
	a.Dispose()
}
