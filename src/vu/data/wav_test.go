// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package data

import (
	"testing"
)

// Check that wave audio file can be imported.
func TestLoadWave(t *testing.T) {
	load := &loader{}
	sound := &Sound{Name: "bloop"}
	if err := load.wav(sound, "../eg/audio", "bloop.wav"); err != nil {
		t.Error()
	}
	if int(sound.DataSize) != len(sound.AudioData) {
		t.Error()
	}
}
