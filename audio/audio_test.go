// SPDX-FileCopyrightText : © 2014-2022 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

package audio

import (
	"testing"
	"time"

	"github.com/gazed/vu/internal/audio/al"
	"github.com/gazed/vu/load"
)

// test that an audio resource can be loaded. Mimics the steps taken
// by the engine. Depends on sound resources from the examples directory.
func TestAudio(t *testing.T) {
	load.SetAssetDir(".wav", "../assets/audio")
	s, err := load.Audio("bloop.wav")
	if err != nil {
		t.Fatalf("audio resource load %s", err)
	}
	a := New()
	err = a.Init()
	if err != nil {
		t.Fatalf("audio init %d : %s", al.GetError(), err)
	}
	at := s.Attrs
	soundData := &Data{}
	soundData.Set(at.Channels, at.SampleBits, at.Frequency, at.DataSize, s.Data)
	snd, buff := uint64(0), uint64(0)
	err = a.LoadSound(&snd, &buff, soundData)
	if err != nil || buff == 0 || snd == 0 {
		t.Fatalf("audio resource bind %d %d : %s", snd, buff, err)
	}
	a.PlaySound(snd, 0, 0, 0)
	time.Sleep(500 * time.Millisecond) // wait for sound to play

	// check that disposing doesn't complain.
	a.DropSound(snd, buff)
	if alerr := al.GetError(); alerr != al.NO_ERROR {
		t.Errorf("DropSound find and fix prior error %d", alerr)
	}
	a.Dispose()
}
