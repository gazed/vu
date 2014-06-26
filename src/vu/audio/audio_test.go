// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package audio

import (
	"testing"
	"vu/load"
)

// test that an audio resource can be loaded. Mimics the
// steps taken by the engine.
func TestAudio(t *testing.T) {
	a := &openal{}
	a.Init()
	loader := load.NewLoader()
	loader.SetDir(2, "../eg/audio") // 2 == load.snd
	sound := newSound("bloop")
	if wh, data, err := loader.Wav(sound.Name()); err == nil {
		sound.SetData(wh.Channels, wh.SampleBits, wh.Frequency, wh.DataSize, data)
	}
	if err := sound.Bind(); err != nil || sound.buffer == 0 {
		t.Error("Failed to load audio resource")
	}

	// Don't play noises during normal testing, but if you're interested...
	// Need to "import time" and sleep a bit for the sound to happen.
	// n := NewSoundMaker(sound)
	// n.Play()
	// time.Sleep(500 * time.Millisecond)
	a.Shutdown()
}
