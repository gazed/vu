// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package audio

import (
	"testing"
	"vu/data"
)

// test that an audio resource can be loaded. Mimics the
// steps taken by the engine.
func TestAudio(t *testing.T) {
	a := &openal{}
	a.Init()
	sound := &data.Sound{}
	loader := data.NewLoader()
	loader.SetDir("../eg/audio", sound)
	loader.Load("bloop", &sound)
	err := a.BindSound(sound)
	if err != nil || sound.Buffer == 0 {
		t.Error("Failed to load audio resource")
	}

	// Don't play noises during normal testing, but if you're interested...
	// Need to "import time" and sleep a bit for the sound to happen.
	// 	n := NewSoundMaker(sound)
	// 	n.Play()
	// 	time.Sleep(500 * time.Millisecond)
	a.Shutdown()
}
