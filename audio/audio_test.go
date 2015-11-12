// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package audio

import (
	"testing"
	//	"time"

	"github.com/gazed/vu/load"
)

// test that an audio resource can be loaded. Mimics the steps taken
// by the engine. Depends on sound resources from the examples directory.
func TestAudio(t *testing.T) {
	a := audioWrapper()
	a.Init()
	loader := load.NewLoader()
	loader.SetDir(2, "../eg/audio") // 2 == load.snd
	soundData := &Data{}
	if wh, data, err := loader.Wav("bloop"); err == nil {
		soundData.Set(wh.Channels, wh.SampleBits, wh.Frequency, wh.DataSize, data)
	}
	snd, buff := uint64(0), uint64(0)
	if err := a.BindSound(&snd, &buff, soundData); err != nil || buff == 0 || snd == 0 {
		t.Errorf("Failed to load audio resource %d %d : %s", snd, buff, err)
	}

	// Don't play noises during normal testing, but if you're interested...
	// ... then uncomment and "import time" (need to sleep for the sound to happen).
	// 	a.PlaySound(snd, 0, 0, 0)
	// 	time.Sleep(500 * time.Millisecond)
	a.Dispose()
}
