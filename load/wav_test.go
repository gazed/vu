// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package load

import (
	"testing"
)

// Uses vu/eg resource directories.
func TestLoadWave(t *testing.T) {
	load := newLoader().setDir(snd, "../eg/audio")
	if wh, data, err := load.wav("bloop"); err != nil || int(wh.DataSize) != len(data) {
		t.Errorf("Loading wave failed %s", err)
	}
}
