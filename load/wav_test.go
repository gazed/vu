// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package load

import (
	"testing"
)

// Uses vu/eg resource directories.
func TestLoadWave(t *testing.T) {
	snd := &SndData{}
	err := snd.Load("bloop", NewLocator().Dir("WAV", "../eg/audio"))
	if err != nil || int(snd.Attrs.DataSize) != len(snd.Data) {
		t.Errorf("Loading wave failed %s", err)
	}
}
